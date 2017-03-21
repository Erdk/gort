package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"testing"
	"time"
)

func BenchmarkGort(*testing.B) {
	flag.Parse()
	*nx = 400
	*ny = 300
	*ns = 300
	*nt = 0
	*output = "bench_output"
	*progress = true

	rand.Seed(time.Now().UnixNano())

	if *progress {
		progCounter = newProgressCounter(uint(*nx * *ny))
	}

	lookfrom := mgl64.Vec3{278.0, 278.0, -800}
	lookat := mgl64.Vec3{278.0, 278.0, 0.0}
	distToFocus := 10.0
	aperture := 0.0
	vfov := 40.0
	vp := newVP(lookfrom, lookat, mgl64.Vec3{0.0, 1.0, 0.0}, vfov, float64(*nx)/float64(*ny), aperture, distToFocus, 0.0, 1.0)

	w := &world{}
	cornellBox(w)

	img := image.NewRGBA(image.Rect(0, 0, *nx, *ny))

	xStripe := 32
	yStripe := 32
	wQ := newQueue(*nx, *ny, xStripe, yStripe)

	var wg sync.WaitGroup
	if *nt == 0 {
		*nt = runtime.NumCPU()
	}
	wg.Add(*nt)

	f := func(threadNum int) {
		defer wg.Done()
		randSource := rand.New(rand.NewSource(time.Now().UnixNano()))
		currentStripe, continueRun := wQ.getJob()
		for continueRun {
			for j := currentStripe.YStart; j < currentStripe.YEnd; j++ {
				for i := currentStripe.XStart; i < currentStripe.XEnd; i++ {
					col := computeXY(randSource, w, vp, i, j)
					img.Set(i, *ny-j, color.RGBA{
						uint8(col.X()),
						uint8(col.Y()),
						uint8(col.Z()),
						255})
				}
			}
			if *progress {
				progCounter.incrementCounter(uint((currentStripe.YEnd - currentStripe.YStart) * (currentStripe.XEnd - currentStripe.XStart)))
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
