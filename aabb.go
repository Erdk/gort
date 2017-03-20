package main

import "github.com/go-gl/mathgl/mgl64"
import "math"

type aabb struct {
	min, max mgl64.Vec3
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
