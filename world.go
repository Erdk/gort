package main

import (
	"math/rand"

	"github.com/go-gl/mathgl/mgl64"
)

type world struct {
	Objs []hitable
}

func (w *world) calcHit(r *ray, tMin, tMax float64) (bool, hit) {
	var retRec hit
	hitAnything := false
	closestSoFar := tMax
	for _, v := range w.Objs {
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

func generateWorld(w *world) {
	w.Objs = make([]hitable, 500)
	i := 0

	w.Objs[i] = &sphere{
		Radius:   0.5,
		Center:   mgl64.Vec3{0.0, 0.0, -1.0},
		Material: lambertian{&mgl64.Vec3{0.1, 0.2, 0.5}}}
	i++

	w.Objs[i] = &sphere{
		Radius:   100,
		Center:   mgl64.Vec3{0.0, -100.5, -1.0},
		Material: lambertian{&mgl64.Vec3{0.8, 0.8, 0.0}}}
	i++

	w.Objs[i] = &sphere{
		Radius:   0.5,
		Center:   mgl64.Vec3{1.0, 0.0, -1.0},
		Material: getMetal(mgl64.Vec3{0.8, 0.6, 0.2}, 0.3)}
	i++

	w.Objs[i] = &sphere{
		Radius:   0.5,
		Center:   mgl64.Vec3{-1.0, 0.0, -1.0},
		Material: dielectric{1.5}}
	i++

	w.Objs[i] = &sphere{
		Radius:   -0.45,
		Center:   mgl64.Vec3{-1.0, 0.0, -1.0},
		Material: dielectric{1.5}}
	i++

	w.Objs[i] = &sphere{
		Radius:   1000.0,
		Center:   mgl64.Vec3{0.0, -1000.0, 0.0},
		Material: lambertian{&mgl64.Vec3{0.5, 0.5, 0.5}},
	}
	i++

	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			chooseMat := rand.Float64()
			center := mgl64.Vec3{
				float64(a) + 0.9*rand.Float64(),
				0.2,
				float64(b) + 0.9*rand.Float64()}

			len := center.Sub(mgl64.Vec3{4.0, 0.2, 0.0}).Len()
			if len > 0.9 {
				if chooseMat < 0.8 { // diffuse
					w.Objs[i] = &sphere{
						Radius: 0.2,
						Center: center,
						Material: lambertian{
							&mgl64.Vec3{
								rand.Float64() * rand.Float64(),
								rand.Float64() * rand.Float64(),
								rand.Float64() * rand.Float64(),
							}},
					}
				} else if chooseMat < 0.95 { // metal
					w.Objs[i] = &sphere{
						Radius: 0.2,
						Center: center,
						Material: getMetal(
							mgl64.Vec3{
								0.5 + 0.5*rand.Float64(),
								0.5 + 0.5*rand.Float64(),
								0.5 + 0.5*rand.Float64()},
							0.5*rand.Float64()),
					}
				} else { // glass
					w.Objs[i] = &sphere{
						Radius:   0.2,
						Center:   center,
						Material: dielectric{1.5},
					}
				}

				i++
			}
		}
	}

	w.Objs[i] = &sphere{
		Radius:   1.0,
		Center:   mgl64.Vec3{0.0, 1.0, 0.0},
		Material: dielectric{1.5}}
	i++
	w.Objs[i] = &sphere{
		Radius:   1.0,
		Center:   mgl64.Vec3{-4.0, 1.0, 0.0},
		Material: lambertian{&mgl64.Vec3{0.4, 0.2, 0.1}}}
	i++
	w.Objs[i] = &sphere{
		Radius:   1.0,
		Center:   mgl64.Vec3{4.0, 1.0, 0.0},
		Material: getMetal(mgl64.Vec3{0.7, 0.6, 0.5}, 0.0)}
}
