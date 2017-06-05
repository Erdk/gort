package rayengine

import (
	"math"
	"math/rand"

	. "github.com/Erdk/gort/rayengine/types"
)

const maxdepth = 5

func randomInUnitSphere(randSource *rand.Rand) *Vec {
	p := &Vec{2.0*randSource.Float64() - 1.0, 2.0*randSource.Float64() - 1.0, 2.0*randSource.Float64() - 1.0}

	for p.Len()*p.Len() >= 1.0 {
		p[0] = 2.0*randSource.Float64() - 1.0
		p[1] = 2.0*randSource.Float64() - 1.0
		p[2] = 2.0*randSource.Float64() - 1.0
	}

	return p
}

func retColor(randSource *rand.Rand, r *ray, w *World, depth int) (float64, float64, float64) {
	if h, rec := w.calcHit(randSource, r, 0.000001, math.MaxFloat64); h {
		emitR, emitG, emitB := rec.m.emit(rec.u, rec.v, rec.p)
		if decision, attenuationR, attenuationG, attenuationB, scattered := rec.m.scatter(randSource, r, rec); decision && depth < maxdepth {
			tmpR, tmpG, tmpB := retColor(randSource, scattered, w, depth+1)
			return emitR + attenuationR*tmpR, emitG + attenuationG*tmpG, emitB + attenuationB*tmpB
		}

		return emitR, emitG, emitB
	}
	return 0.0, 0.0, 0.0
}

func ComputeXY(randSource *rand.Rand, w *World, cam *camera, x, y, nx, ny, ns int) *Vec {
	col := &Vec{0.0, 0.0, 0.0}
	for s := 0; s < ns; s++ {
		u := (float64(x) + randSource.Float64()) / float64(nx)
		v := (float64(y) + randSource.Float64()) / float64(ny)
		r := cam.getRay(randSource, u, v)
		rR, rG, rB := retColor(randSource, &r, w, 0)
		col[0] += rR
		col[1] += rG
		col[2] += rB
	}

	col = col.MulSI(1.0 / float64(ns))
	col[0] = math.Sqrt(col[0]) * 255.99
	col[1] = math.Sqrt(col[1]) * 255.99
	col[2] = math.Sqrt(col[2]) * 255.99

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
