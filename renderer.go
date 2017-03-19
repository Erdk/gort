package main

import (
	"math"
	"math/rand"

	. "github.com/Erdk/gort/types"
)

func randomInUnitSphere(randSource *rand.Rand) *Vec {
	p := &Vec{2.0*randSource.Float64() - 1.0, 2.0*randSource.Float64() - 1.0, 2.0*randSource.Float64() - 1.0}

	for p.Len()*p.Len() >= 1.0 {
		p[0] = 2.0*randSource.Float64() - 1.0
		p[1] = 2.0*randSource.Float64() - 1.0
		p[2] = 2.0*randSource.Float64() - 1.0
	}

	return p
}

func retColor(randSource *rand.Rand, r *ray, w *world, depth int) *Vec {
	if h, rec := w.calcHit(randSource, r, 0.001, math.MaxFloat64); h {
		emit := rec.m.emit(rec.u, rec.v, rec.p)
		if decision, attenuation, scattered := rec.m.scatter(randSource, r, rec); decision && depth < 50 {
			tmp := retColor(randSource, scattered, w, depth+1)
			return &Vec{
				emit[0] + attenuation[0]*tmp[0],
				emit[1] + attenuation[1]*tmp[1],
				emit[2] + attenuation[2]*tmp[2],
			}
		}

		return emit
	}
	return &Vec{0.0, 0.0, 0.0}
}

func computeXY(randSource *rand.Rand, w *world, vp *viewport, x, y int) *Vec {
	col := &Vec{0.0, 0.0, 0.0}
	for s := 0; s < *ns; s++ {
		u := (float64(x) + randSource.Float64()) / float64(*nx)
		v := (float64(y) + randSource.Float64()) / float64(*ny)
		r := vp.getRay(randSource, u, v)
		col.AddVM(retColor(randSource, &r, w, 0))
	}

	col = col.MulSI(1.0 / float64(*ns))
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
