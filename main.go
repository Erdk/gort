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

	"github.com/go-gl/mathgl/mgl64"
	"github.com/pkg/profile"
)

var nx = flag.Int("w", 640, "width of rendered image")
var ny = flag.Int("h", 480, "height of rendered image")
var ns = flag.Int("s", 400, "samples per pixel")
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

	if *progress {
		progCounter = newProgressCounter(uint(*nx * *ny))
	}

	// seed random number generator
	rand.Seed(time.Now().UnixNano())

	lookfrom := mgl64.Vec3{278.0, 278.0, -800}
	lookat := mgl64.Vec3{278.0, 278.0, 0.0}
	distToFocus := 10.0
	aperture := 0.0
	vfov := 40.0
	vp := newVP(lookfrom, lookat, mgl64.Vec3{0.0, 1.0, 0.0}, vfov, float64(*nx)/float64(*ny), aperture, distToFocus, 0.0, 1.0)

	w := &world{}

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
		testTexture(w)

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

	xStripe, yStripe, err := parseStripe(*computeUnit)
	if err != nil {
		panic("Wrong stripe format!")
	}
	wQ := newQueue(*nx, *ny, xStripe, yStripe)

	var wg sync.WaitGroup
	if *nt == 0 {
		*nt = runtime.NumCPU()
	}
	wg.Add(int(*nt))

	f := func(threadNum int) {
		defer wg.Done()
		currentStripe, continueRun := wQ.getJob()
		for continueRun {
			randSource := rand.New(rand.NewSource(time.Now().UnixNano()))
			for j := currentStripe.yStart; j < currentStripe.yEnd; j++ {
				for i := currentStripe.xStart; i < currentStripe.xEnd; i++ {
					col := computeXY(randSource, w, vp, i, j)
					img.Set(i, *ny-j-1, color.RGBA{
						uint8(col.X()),
						uint8(col.Y()),
						uint8(col.Z()),
						255})
				}
			}

			if *progress {
				progCounter.incrementCounter(uint((currentStripe.yEnd - currentStripe.yStart) * (currentStripe.xEnd - currentStripe.xStart)))
			}

			currentStripe, continueRun = wQ.getJob()
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
