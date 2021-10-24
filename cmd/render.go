// gort renderer
// Copyright (C) 2017 Erdk <mr.erdk@gmail.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
// Copyright © 2017 Erdk <mr.erdk@gmail.com>

package cmd

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"time"

	re "github.com/Erdk/gort/pkg/rayengine"
	"github.com/Erdk/gort/pkg/util"
	"github.com/pkg/profile"
	"github.com/spf13/cobra"
)

var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "Render scene",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		render()
	},
}

var nx uint
var ny uint
var ns uint
var nt uint
var output string
var scene string
var prof string
var progress bool
var computeUnit string
var printJSON bool
var loadJSON string

func init() {
	rootCmd.AddCommand(renderCmd)

	// Flags
	renderCmd.Flags().UintVarP(&nx, "width", "w", 640, "width of rendered image")
	renderCmd.Flags().UintVarP(&ny, "height", "e", 480, "width of rendered image")
	renderCmd.Flags().UintVarP(&ns, "samples", "s", 500, "width of rendered image")
	renderCmd.Flags().UintVarP(&nt, "threads", "t", 1, "width of rendered image, 0 to launch one thread per CPU")
	renderCmd.Flags().StringVarP(&output, "output", "o", "", "filename without extension, default: output_<timestamp>")
	renderCmd.Flags().StringVar(&scene, "scene", "defRoomOneTriangle", "chose scene to render")
	renderCmd.Flags().StringVarP(&prof, "profile", "r", "", "generate cpu/mem/block profile")
	renderCmd.Flags().BoolVarP(&progress, "progress", "p", false, "show progress")
	renderCmd.Flags().StringVar(&computeUnit, "computeunit", "16x16", "unit of computation, format: wxh where w - width of stripe and h is height of stripe")
	renderCmd.Flags().BoolVarP(&printJSON, "print", "", false, "(DEBUG) print JSON profile of renderer to stdout")
	renderCmd.Flags().StringVarP(&loadJSON, "load", "l", "", "renders scene from JSON file")
}

func render() {
	switch prof {
	case "cpu":
		defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	case "mem":
		defer profile.Start(profile.MemProfile, profile.ProfilePath(".")).Stop()
	case "block":
		defer profile.Start(profile.BlockProfile, profile.ProfilePath(".")).Stop()
	default:
	}

	var progressCounter *util.ProgressCounter
	if progress {
		progressCounter = util.NewProgressCounter(nx * ny)
	}

	// seed random number generator
	rand.Seed(time.Now().UnixNano())

	var w *re.World
	if loadJSON != "" {
		// placeholder for proper scene loader
		w = re.NewWorld(scene, float64(nx), float64(ny))
	} else {
		w = re.NewWorld(scene, float64(nx), float64(ny))
	}

	// prints back scene to stdout, for debug purposes only
	if printJSON {
		b, _ := json.Marshal(w)
		fmt.Printf("scene:\n %v", string(b))
		return
	}

	img := image.NewRGBA(image.Rect(0, 0, int(nx), int(ny)))

	xPatch, yPatch, err := util.StringToPatch(computeUnit)
	if err != nil {
		panic("Wrong stripe format!")
	}
	wQ := util.NewQueue(nx, ny, xPatch, yPatch)

	var wg sync.WaitGroup
	if nt == 0 {
		nt = uint(runtime.NumCPU())
		runtime.GOMAXPROCS(int(nt))
	}
	wg.Add(int(nt))

	f := func(threadNum int) {
		defer wg.Done()
		currentPatch, continueRun := wQ.GetJob()
		for continueRun {
			randSource := rand.New(rand.NewSource(time.Now().UnixNano()))
			for j := currentPatch.YStart; j < currentPatch.YEnd; j++ {
				for i := currentPatch.XStart; i < currentPatch.XEnd; i++ {
					col := re.ComputeXY(randSource, w, i, j, nx, ny, ns)
					img.Set(int(i), int(ny-j-1), color.RGBA{
						uint8(col[0]),
						uint8(col[1]),
						uint8(col[2]),
						255})
				}
			}

			if progress {
				progressCounter.IncrementCounter(uint((currentPatch.YEnd - currentPatch.YStart) * (currentPatch.XEnd - currentPatch.XStart)))
			}

			currentPatch, continueRun = wQ.GetJob()
		}
	}

	for i := 0; i < int(nt); i++ {
		go f(i)
	}

	wg.Wait()

	if output == "" {
		t := time.Now()
		output = "output_" + t.Format("20060102150405")
	}
	fd, _ := os.Create(output + ".png")
	defer fd.Close()

	png.Encode(fd, img)
	if progress {
		fmt.Printf("\n")
	}
}
