package main

import (
	"math"
	"math/rand"

	"github.com/go-gl/mathgl/mgl64"
)

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

func computeXY(w *world, vp *viewport, nx, ny, ns, x, y int) mgl64.Vec3 {
	col := mgl64.Vec3{0.0, 0.0, 0.0}
	for s := 0; s < ns; s++ {
		u := (float64(x) + rand.Float64()) / float64(nx)
		v := (float64(y) + rand.Float64()) / float64(ny)
		r := vp.getRay(u, v)
		col = col.Add(retColor(&r, w, 0))
	}

	col = col.Mul(1.0 / float64(ns))
	col = mgl64.Vec3{
		math.Sqrt(col.X()) * 255.99,
		math.Sqrt(col.Y()) * 255.99,
		math.Sqrt(col.Z()) * 255.99,
	}

	return col
}
