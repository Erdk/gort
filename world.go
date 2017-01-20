package main

import (
	"math/rand"

	"github.com/go-gl/mathgl/mgl64"
)

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

func generateWorld(w *world) {
	w.objs = make([]hitable, 500)
	i := 0

	w.objs[i] = &sphere{
		r: 0.5,
		c: mgl64.Vec3{0.0, 0.0, -1.0},
		m: lambertian{&mgl64.Vec3{0.1, 0.2, 0.5}}}
	i++

	w.objs[i] = &sphere{
		r: 100,
		c: mgl64.Vec3{0.0, -100.5, -1.0},
		m: lambertian{&mgl64.Vec3{0.8, 0.8, 0.0}}}
	i++

	w.objs[i] = &sphere{
		r: 0.5,
		c: mgl64.Vec3{1.0, 0.0, -1.0},
		m: getMetal(mgl64.Vec3{0.8, 0.6, 0.2}, 0.3)}
	i++

	w.objs[i] = &sphere{
		r: 0.5,
		c: mgl64.Vec3{-1.0, 0.0, -1.0},
		m: dielectric{1.5}}
	i++

	w.objs[i] = &sphere{
		r: -0.45,
		c: mgl64.Vec3{-1.0, 0.0, -1.0},
		m: dielectric{1.5}}
	i++

	w.objs[i] = &sphere{
		r: 1000.0,
		c: mgl64.Vec3{0.0, -1000.0, 0.0},
		m: lambertian{&mgl64.Vec3{0.5, 0.5, 0.5}},
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
					w.objs[i] = &sphere{
						r: 0.2,
						c: center,
						m: lambertian{
							&mgl64.Vec3{
								rand.Float64() * rand.Float64(),
								rand.Float64() * rand.Float64(),
								rand.Float64() * rand.Float64(),
							}},
					}
				} else if chooseMat < 0.95 { // metal
					w.objs[i] = &sphere{
						r: 0.2,
						c: center,
						m: getMetal(
							mgl64.Vec3{
								0.5 + 0.5*rand.Float64(),
								0.5 + 0.5*rand.Float64(),
								0.5 + 0.5*rand.Float64()},
							0.5*rand.Float64()),
					}
				} else { // glass
					w.objs[i] = &sphere{
						r: 0.2,
						c: center,
						m: dielectric{1.5},
					}
				}

				i++
			}
		}
	}

	w.objs[i] = &sphere{
		r: 1.0,
		c: mgl64.Vec3{0.0, 1.0, 0.0},
		m: dielectric{1.5}}
	i++
	w.objs[i] = &sphere{
		r: 1.0,
		c: mgl64.Vec3{-4.0, 1.0, 0.0},
		m: lambertian{&mgl64.Vec3{0.4, 0.2, 0.1}}}
	i++
	w.objs[i] = &sphere{
		r: 1.0,
		c: mgl64.Vec3{4.0, 1.0, 0.0},
		m: getMetal(mgl64.Vec3{0.7, 0.6, 0.5}, 0.0)}
}
