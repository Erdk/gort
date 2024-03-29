// gort renderer
// Copyright (C) 2017 Erdk <mr.erdk@gmail.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
// Copyright © 2017 Erdk <mr.erdk@gmail.com>

package rayengine

import (
	"math"
	"math/rand"
)

const maxdepth = 7

func randomInUnitSphere(randSource *rand.Rand) *Vec {
	p := &Vec{
		2.0*randSource.Float64() - 1.0,
		2.0*randSource.Float64() - 1.0,
		2.0*randSource.Float64() - 1.0,
	}

	for p.Len()*p.Len() <= EPSILON {
		p[0] = 2.0*randSource.Float64() - 1.0
		p[1] = 2.0*randSource.Float64() - 1.0
		p[2] = 2.0*randSource.Float64() - 1.0
	}

	return p.Normalize()
}

func retColor(randSource *rand.Rand, r *ray, w *World, depth int) (float64, float64, float64) {
	if h, pointOfHit := w.calcHit(randSource, r, EPSILON, math.MaxFloat64); h {
		emitR, emitG, emitB := pointOfHit.m.emit(pointOfHit.u, pointOfHit.v, pointOfHit.p)
		if decision, attenuationR, attenuationG, attenuationB, scattered := pointOfHit.m.scatter(randSource, r, pointOfHit); decision && depth < maxdepth {
			tmpR, tmpG, tmpB := retColor(randSource, scattered, w, depth+1)
			return emitR + attenuationR*tmpR, emitG + attenuationG*tmpG, emitB + attenuationB*tmpB
		}

		return emitR, emitG, emitB
	}
	return 0.0, 0.0, 0.0
}

// ComputeXY returns color at (X, Y) position of camera plane, computed ns times and averaged
func ComputeXY(randSource *rand.Rand, w *World, x, y, nx, ny, ns uint) *Vec {
	col := &Vec{0.0, 0.0, 0.0}
	for s := uint(0); s < ns; s++ {
		u := (float64(x) + randSource.Float64()) / float64(nx)
		v := (float64(y) + randSource.Float64()) / float64(ny)
		r := w.Cam.getRay(randSource, u, v)
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
