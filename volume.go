package main

import (
	"math"
	"math/rand"

	. "github.com/Erdk/gort/types"
)

type constantMedium struct {
	Boundary hitable
	Density  float64
	Material material
}

func (c *constantMedium) calcHit(randSource *rand.Rand, r *ray, min, max float64) (bool, hit) {
	var rec hit
	if decision, rec1 := c.Boundary.calcHit(randSource, r, -math.MaxFloat64, math.MaxFloat64); decision {
		if decision, rec2 := c.Boundary.calcHit(randSource, r, rec1.t+0.0001, math.MaxFloat64); decision {
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
