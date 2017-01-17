package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
	"sync"

	"time"

	"github.com/go-gl/mathgl/mgl64"
)

type ray struct {
	origin, direction *mgl64.Vec3
}

func (r *ray) pointAtParam(t float64) mgl64.Vec3 {
	return r.origin.Add(r.direction.Mul(t))
}

type hit struct {
	t    float64
	p, n mgl64.Vec3
	m    material
}

type hitable interface {
	calcHit(r *ray, tMin, tMax float64) (bool, hit)
}

type sphere struct {
	c mgl64.Vec3
	r float64
	m material
}

func (s *sphere) calcHit(r *ray, tMin, tMax float64) (bool, hit) {
	oc := r.origin.Sub(s.c)
	a := r.direction.Dot(*r.direction)
	b := oc.Dot(*r.direction)
	c := oc.Dot(oc) - s.r*s.r

	discriminant := b*b - a*c
	if discriminant > 0 {
		bbac := math.Sqrt(b*b - a*c)
		temp := (-b - bbac) / a
		if temp < tMax && temp > tMin {
			var rec hit
			rec.t = temp
			rec.p = r.pointAtParam(rec.t)
			rec.n = rec.p.Sub(s.c)
			rec.n = rec.n.Mul(1.0 / s.r)
			rec.m = s.m
			return true, rec
		}

		temp = (-b + bbac) / a
		if temp < tMax && temp > tMin {
			var rec hit
			rec.t = temp
			rec.p = r.pointAtParam(rec.t)
			rec.n = rec.p.Sub(s.c)
			rec.n = rec.n.Mul(1.0 / s.r)
			rec.m = s.m
			return true, rec
		}
	}

	return false, hit{}
}

type world struct {
	objs []hitable
}

func (w *world) calcHit(r *ray, tMin, tMax float64) (bool, hit) {
	var retRec hit
	hitAnything := false
	closestSoFar := tMax
	for _, v := range w.objs {
		if v == nil {
			continue
		}
		if h, rec := v.calcHit(r, tMin, closestSoFar); h {
			hitAnything = true
			closestSoFar = rec.t
			retRec = rec
		}
	}

	if hitAnything {
		return true, retRec
	}

	return false, hit{}
}

func randomInUnitSphere() mgl64.Vec3 {
	p := mgl64.Vec3{2.0*rand.Float64() - 1.0, 2.0*rand.Float64() - 1.0, 2.0*rand.Float64() - 1.0}
	for p.Len()*p.Len() >= 1.0 {
		p = mgl64.Vec3{2.0*rand.Float64() - 1.0, 2.0*rand.Float64() - 1.0, 2.0*rand.Float64() - 1.0}
	}

	return p
}

func retColor(r *ray, w *world, depth int) mgl64.Vec3 {
	if h, rec := w.calcHit(r, 0.001, math.MaxFloat64); h {
		if decision, attenuation, scattered := rec.m.scatter(*r, rec); decision && depth < 50 {
			tmp := retColor(scattered, w, depth+1)
			return mgl64.Vec3{
				attenuation.X() * tmp.X(),
				attenuation.Y() * tmp.Y(),
				attenuation.Z() * tmp.Z(),
			}
		}

		return mgl64.Vec3{0.0, 0.0, 0.0}
	}

	uv := r.direction.Normalize()
	t := 0.5 * (uv.Y() + 1.0)
	ret := mgl64.Vec3{1.0 - t, 1.0 - t, 1.0 - t}
	tmp := mgl64.Vec3{0.5 * t, 0.7 * t, 1.0 * t}
	return ret.Add(tmp)
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
				col := mgl64.Vec3{0.0, 0.0, 0.0}
				for s := 0; s < ns; s++ {
					u := (float64(i) + rand.Float64()) / float64(nx)
					v := (float64(j) + rand.Float64()) / float64(ny)
					r := vp.getRay(u, v)
					col = col.Add(retColor(&r, w, 0))
				}

				col = col.Mul(1.0 / float64(ns))
				col = mgl64.Vec3{
					math.Sqrt(col.X()) * 255.99,
					math.Sqrt(col.Y()) * 255.99,
					math.Sqrt(col.Z()) * 255.99,
				}
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
