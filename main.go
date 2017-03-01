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

var nx = flag.Int("w", 320, "width of rendered image")
var ny = flag.Int("h", 240, "height of rendered image")
var ns = flag.Int("s", 200, "samples per pixel")
var nt = flag.Int("t", 1, "number of parallel threads")
var output = flag.String("o", "output", "filename without extension")
var input = flag.String("i", "", "instead of generating world render one from file")
var saveraw = flag.Bool("j", false, "saves generated scene into <output>.json")
var prof = flag.String("prof", "", "generate cpu/mem/block profile, by default none")
var progress = flag.Bool("p", false, "show progress, default false")

type progressCounter struct {
	counter, max, lastPrinted int
	mtx                       *sync.Mutex
}

var progCounter *progressCounter

func (p *progressCounter) incrementCounter() {
	p.mtx.Lock()
	p.counter++
	newPrinted := int(float64(p.counter) / float64(p.max) * 100)
	if newPrinted > p.lastPrinted {
		p.lastPrinted = newPrinted
		fmt.Printf("\r%d%%", p.lastPrinted)
	}
	p.mtx.Unlock()
}

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
		progCounter = &progressCounter{}
		progCounter.counter = 0
		progCounter.max = *nx * *ny
		progCounter.lastPrinted = 0
		progCounter.mtx = &sync.Mutex{}
	}

	// seed random number generator
	rand.Seed(time.Now().UnixNano())

	// lookfrom := mgl64.Vec3{26.0, 2.0, 4.0}
	// lookat := mgl64.Vec3{0.0, 4.0, 0.0}
	// distToFocus := 10.0
	// aperture := 0.0
	// vp := newVP(lookfrom, lookat, mgl64.Vec3{0.0, 1.0, 0.0}, 20.0, float64(*nx)/float64(*ny), aperture, distToFocus, 0.0, 1.0)

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
		cornellBox(w)

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

	img := image.NewRGBA(image.Rect(0, 0, *nx, *ny))

	var wg sync.WaitGroup
	if *nt == 0 {
		*nt = runtime.NumCPU()
	}
	wg.Add(*nt)

	f := func(threadNum, x1, x2 int) {
		defer wg.Done()
		for j := 0; j < *ny; j++ {
			for i := x1; i < x2; i++ {
				col := computeXY(w, vp, i, j)
				// if col.X() < 0.0 || col.X() > 255.0 || col.Y() < 0.0 || col.Y() > 255.0 || col.Z() < 0.0 || col.Z() > 255.0 {
				// 	fmt.Printf("WRONG COLOUR! %f %f %f", col.X(), col.Y(), col.Z())
				// }
				img.Set(i, *ny-j, color.RGBA{
					uint8(col.X()),
					uint8(col.Y()),
					uint8(col.Z()),
					255})
				if *progress {
					progCounter.incrementCounter()
				}
			}
		}
	}

	for i := 0; i < *nt; i++ {
		go f(i, *nx / *nt * i, *nx / *nt * (i+1))
	}

	fd, _ := os.Create(*output + ".png")
	defer fd.Close()

	wg.Wait()
	png.Encode(fd, img)
	if *progress {
		fmt.Printf("\n")
	}
}
