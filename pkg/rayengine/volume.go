// gort renderer
// Copyright (C) 2017 Łukasz 'Erdk' Redynk <mr.erdk@gmail.com>
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
// Copyright © 2017 Łukasz 'Erdk' Redynk <mr.erdk@gmail.com>

package rayengine

import (
	"math"
	"math/rand"
)

type constantMedium struct {
	Boundary hitable
	Density  float64
	Material material
}

func (c *constantMedium) calcHit(randSource *rand.Rand, r *ray, min, max float64) (bool, hit) {
	var rec hit
	if decision, rec1 := c.Boundary.calcHit(randSource, r, -math.MaxFloat64, math.MaxFloat64); decision {
		if decision, rec2 := c.Boundary.calcHit(randSource, r, rec1.t+EPSILON*10, math.MaxFloat64); decision {
			if rec1.t < min {
				rec1.t = min
			}

			if rec2.t > max {
				rec2.t = max
			}

			if rec1.t >= rec2.t {
				return false, hit{}
			}

			if rec1.t < 0 {
				rec1.t = 0
			}

			distanceInsideBoundary := (rec2.t - rec1.t) * r.direction.Len()
			hitDistance := -(1.0 / c.Density) * math.Log(randSource.Float64())

			if hitDistance < distanceInsideBoundary {
				rec.t = rec1.t + hitDistance/r.direction.Len()
				rec.p = r.pointAtParam(rec.t)
				rec.normal = &Vec{0.0, 0.0, 0.0}
				rec.m = c.Material

				return true, rec
			}
		}
	}

	return false, hit{}
}

func (c *constantMedium) boundingBox(t0, t1 float64) (bool, *aabb) {
	return c.Boundary.boundingBox(t0, t1)
}
