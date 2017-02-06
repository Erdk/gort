package main

import "github.com/go-gl/mathgl/mgl64"

type ray struct {
	origin, direction *mgl64.Vec3
	time              float64
}

func (r *ray) pointAtParam(t float64) mgl64.Vec3 {
	return r.origin.Add(r.direction.Mul(t))
}
