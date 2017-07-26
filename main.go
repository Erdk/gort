package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"sync"

	"time"

	"runtime"

	re "github.com/Erdk/gort/rayengine"
	"github.com/Erdk/gort/util"
	"github.com/pkg/profile"
)

var nx = flag.Int("w", 640, "width of rendered image")
var ny = flag.Int("h", 480, "height of rendered image")
var ns = flag.Int("s", 1000, "samples per pixel")
var nt = flag.Int("t", 1, "number of parallel threads")
var output = flag.String("o", "output", "filename without extension")
var input = flag.String("i", "", "instead of generating world render one from file")
var saveraw = flag.Bool("j", false, "saves generated scene into <output>.json")
var prof = flag.String("prof", "", "generate cpu/mem/block profile, by default none")
var progress = flag.Bool("p", false, "show progress, default false")
var computeUnit = flag.String("cu", "16x16", "unit of computation, format: wxh where w - width of stripe and h is height of stripe, by default '16x16'")

func main() {
	flag.Parse()

	switch *prof {
	case "cpu":
		defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	case "mem":
		defer profile.Start(profile.MemProfile, profile.ProfilePath(".")).Stop()
	case "block":
		defer profile.Start(profile.BlockProfile, profile.ProfilePath(".")).Stop()
	default:
	}

	var progressCounter *util.ProgressCounter
	if *progress {
		progressCounter = util.NewProgressCounter(uint(*nx * *ny))
	}

	// seed random number generator
	rand.Seed(time.Now().UnixNano())

	w := &re.World{}
	if *input != "" {
		w = re.NewWorld(*input, float64(*nx), float64(*ny))
	} else {
		w = re.NewWorld("colorVolWorld", float64(*nx), float64(*ny))
	}

	img := image.NewRGBA(image.Rect(0, 0, int(*nx), int(*ny)))

	xStripe, yStripe, err := util.ParseStripe(*computeUnit)
	if err != nil {
		panic("Wrong stripe format!")
	}
	wQ := util.NewQueue(*nx, *ny, xStripe, yStripe)

	var wg sync.WaitGroup
	if *nt == 0 {
		*nt = runtime.NumCPU()
		runtime.GOMAXPROCS(*nt)
	}
	wg.Add(int(*nt))

	f := func(threadNum int) {
		defer wg.Done()
		currentStripe, continueRun := wQ.GetJob()
		for continueRun {
			randSource := rand.New(rand.NewSource(time.Now().UnixNano()))
			for j := currentStripe.YStart; j < currentStripe.YEnd; j++ {
				for i := currentStripe.XStart; i < currentStripe.XEnd; i++ {
					col := re.ComputeXY(randSource, w, i, j, *nx, *ny, *ns)
					img.Set(i, *ny-j-1, color.RGBA{
						uint8(col[0]),
						uint8(col[1]),
						uint8(col[2]),
						255})
				}
			}

			if *progress {
				progressCounter.IncrementCounter(uint((currentStripe.YEnd - currentStripe.YStart) * (currentStripe.XEnd - currentStripe.XStart)))
			}

			currentStripe, continueRun = wQ.GetJob()
		}
	}

	for i := 0; i < *nt; i++ {
		go f(i)
	}

	fd, _ := os.Create(*output + ".png")
	defer fd.Close()

	wg.Wait()
	png.Encode(fd, img)
	if *progress {
		fmt.Printf("\n")
	}
}
