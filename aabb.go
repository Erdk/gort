package main

import (
	"math"

	. "github.com/Erdk/gort/types"
)

type aabb struct {
	min, max *Vec
}

func (a *aabb) hit(r *ray, tmin, tmax float64) bool {
	for i := 0; i < 3; i++ {
		t0 := math.Min((a.min[i]-r.origin[i])/r.direction[i],
			(a.max[i]-r.origin[i])/r.direction[i])
		t1 := math.Max((a.min[i]-r.origin[i])/r.direction[i],
			(a.max[i]-r.origin[i])/r.direction[i])
		tmin = math.Max(t0, tmin)
		tmax = math.Min(t1, tmax)
		if tmax <= tmin {
			return false
		}

	}

	return true
}
