package main

import (
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

type sphere struct {
	Center   mgl64.Vec3
	Radius   float64
	Material material
}

func (s *sphere) calcHit(r *ray, tMin, tMax float64) (bool, hit) {
	oc := r.origin.Sub(s.Center)
	a := r.direction.Dot(*r.direction)
	b := oc.Dot(*r.direction)
	c := oc.Dot(oc) - s.Radius*s.Radius

	discriminant := b*b - a*c
	if discriminant > 0 {
		bbac := math.Sqrt(b*b - a*c)
		temp := (-b - bbac) / a
		if temp < tMax && temp > tMin {
			var rec hit
			rec.t = temp
			rec.p = r.pointAtParam(rec.t)
			rec.n = rec.p.Sub(s.Center)
			rec.n = rec.n.Mul(1.0 / s.Radius)
			rec.m = s.Material
			return true, rec
		}

		temp = (-b + bbac) / a
		if temp < tMax && temp > tMin {
			var rec hit
			rec.t = temp
			rec.p = r.pointAtParam(rec.t)
			rec.n = rec.p.Sub(s.Center)
			rec.n = rec.n.Mul(1.0 / s.Radius)
			rec.m = s.Material
			return true, rec
		}
	}

	return false, hit{}
}
