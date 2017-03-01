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
		emit := rec.m.emit(rec.u, rec.v, rec.p)
		if decision, attenuation, scattered := rec.m.scatter(*r, rec); decision && depth < 50 {
			tmp := retColor(scattered, w, depth+1)
			return mgl64.Vec3{
				emit.X() + attenuation.X()*tmp.X(),
				emit.Y() + attenuation.Y()*tmp.Y(),
				emit.Z() + attenuation.Z()*tmp.Z(),
			}
		}

		return *emit
	}
	return mgl64.Vec3{0.0, 0.0, 0.0}
}

func computeXY(w *world, vp *viewport, x, y int) mgl64.Vec3 {
	col := mgl64.Vec3{0.0, 0.0, 0.0}
	for s := 0; s < *ns; s++ {
		u := (float64(x) + rand.Float64()) / float64(*nx)
		v := (float64(y) + rand.Float64()) / float64(*ny)
		r := vp.getRay(u, v)
		col = col.Add(retColor(&r, w, 0))
	}

	col = col.Mul(1.0 / float64(*ns))
	col = mgl64.Vec3{
		math.Sqrt(col.X()) * 255.99,
		math.Sqrt(col.Y()) * 255.99,
		math.Sqrt(col.Z()) * 255.99,
	}

	// "normalize" colours
	for i := range col {
		if col[i] > 255.0 {
			col[i] = 255.0
		}

		if col[i] < 0.0 {
			col[i] = 0.0
		}
	}
	return col
}
