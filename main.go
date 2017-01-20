package main

import (
	"flag"
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

var nx = flag.Int("w", 640, "width of rendered image")
var ny = flag.Int("h", 480, "height of rendered image")
var ns = flag.Int("s", 200, "samples per pixel")
var nt = flag.Int("t", 2, "number of parallel threads")
var output = flag.String("s", "output", "filename without extension")

func main() {
	flag.Parse()

	// seed random number generator
	rand.Seed(time.Now().UnixNano())

	lookfrom := mgl64.Vec3{13.0, 2.0, 3.0}
	lookat := mgl64.Vec3{0.0, 0.0, 0.0}
	distToFocus := 10.0
	aperture := 0.1
	vp := newVP(lookfrom, lookat, mgl64.Vec3{0.0, 1.0, 0.0}, 20.0, float64(*nx)/float64(*ny), aperture, distToFocus)

	w := &world{}
	generateWorld(w)

	img := image.NewRGBA(image.Rect(0, 0, *nx, *ny))

	var wg sync.WaitGroup
	wg.Add(*nt)

	f := func(threadNum, x1, x2 int) {
		defer wg.Done()
		for j := 0; j < *ny; j++ {
			for i := x1; i < x2; i++ {
				col := computeXY(w, vp, *nx, *ny, *ns, i, j)
				img.Set(i, *ny-j, color.RGBA{
					uint8(col.X()),
					uint8(col.Y()),
					uint8(col.Z()),
					255})
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
}
