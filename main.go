package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"sync"

	"time"

	"github.com/go-gl/mathgl/mgl64"
)

type hit struct {
	t    float64
	p, n mgl64.Vec3
	m    material
}

type hitable interface {
	calcHit(r *ray, tMin, tMax float64) (bool, hit)
}

const cTHREADS = 6

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %v <filename>\n", os.Args[0])
		os.Exit(1)
	}

	// seed random number generator
	rand.Seed(time.Now().UnixNano())

	// hardcoded image dimenssions
	nx := 1920
	ny := 1080
	ns := 200

	lookfrom := mgl64.Vec3{13.0, 2.0, 3.0}
	lookat := mgl64.Vec3{0.0, 0.0, 0.0}
	distToFocus := 10.0
	aperture := 0.1
	vp := newVP(lookfrom, lookat, mgl64.Vec3{0.0, 1.0, 0.0}, 20.0, float64(nx)/float64(ny), aperture, distToFocus)

	w := &world{}
	generateWorld(w)

	img := image.NewRGBA(image.Rect(0, 0, nx, ny))

	var wg sync.WaitGroup
	wg.Add(cTHREADS)

	f := func(threadNum, x1, x2 int) {
		defer wg.Done()
		for j := 0; j < ny; j++ {
			for i := x1; i < x2; i++ {
				//fmt.Printf("Thread %v: x: %v y: %v\n", threadNum, i, j)
				col := computeXY(w, vp, nx, ny, ns, i, j)
				img.Set(i, ny-j, color.RGBA{
					uint8(col.X()),
					uint8(col.Y()),
					uint8(col.Z()),
					255})
			}
		}
	}

	for i := 0; i < cTHREADS; i++ {
		go f(i, nx/cTHREADS*i, nx/cTHREADS*(i+1))
	}

	fd, _ := os.Create(os.Args[1] + ".png")
	defer fd.Close()

	wg.Wait()
	png.Encode(fd, img)
}
