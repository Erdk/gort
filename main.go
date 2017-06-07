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

	"encoding/json"

	"runtime"

	re "github.com/Erdk/gort/rayengine"
	. "github.com/Erdk/gort/rayengine/types"
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

	var progCounter *util.ProgressCounter
	if *progress {
		progCounter = util.NewProgressCounter(uint(*nx * *ny))
	}

	// seed random number generator
	rand.Seed(time.Now().UnixNano())

	w := &re.World{}

	lookfrom := &Vec{278.0, 278.0, -700}
	lookat := &Vec{278.0, 278.0, 0.0}
	distToFocus := 10.0
	aperture := 0.0
	vfov := 40.0
	w.Cam = re.NewCamera(lookfrom, lookat, &Vec{0.0, 1.0, 0.0}, vfov,
		float64(*nx)/float64(*ny), aperture, distToFocus, 0.0, 1.0)

	if *input != "" {
		fd, _ := os.Open(*input)
		defer fd.Close()
		dec := json.NewDecoder(fd)
		if err := dec.Decode(w); err != nil {
			fmt.Printf("Cannot read scene: %s", err.Error())
			os.Exit(1)
		}
	} else {
		//generateWorld(w)
		//perlinTest(w)
		//lightAndRectTest(w)
		//cornellBox(w)
		//generateWorld2(w)
		//testTexture(w)
		re.ColorVolWorld(w)

		if *saveraw {
			raw, err := json.Marshal(*w)
			if err != nil {
				fmt.Printf("Cannot marshall scene: %s", err.Error())
				os.Exit(1)
			}

			fd, _ := os.Create(*output + ".json")
			defer fd.Close()
			if _, err := fd.Write(raw); err != nil {
				fmt.Printf("Cannot write raw json with scene: %s\n", err.Error())
			}
		}
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
				progCounter.IncrementCounter(uint((currentStripe.YEnd - currentStripe.YStart) * (currentStripe.XEnd - currentStripe.XStart)))
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
