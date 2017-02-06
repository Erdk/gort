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

type movingSphere struct {
	Center0, Center1 mgl64.Vec3
	Time0, Time1     float64
	Radius           float64
	Material         material
}

func (m *movingSphere) calcHit(r *ray, tMin, tMax float64) (bool, hit) {
	oc := r.origin.Sub(m.center(r.time))
	a := r.direction.Dot(*r.direction)
	b := oc.Dot(*r.direction)
	c := oc.Dot(oc) - m.Radius*m.Radius

	discriminant := b*b - a*c
	if discriminant > 0 {
		bbac := math.Sqrt(b*b - a*c)
		temp := (-b - bbac) / a
		if temp < tMax && temp > tMin {
			var rec hit
			rec.t = temp
			rec.p = r.pointAtParam(rec.t)
			rec.n = rec.p.Sub(m.center(r.time))
			rec.n = rec.n.Mul(1.0 / m.Radius)
			rec.m = m.Material
			return true, rec
		}

		temp = (-b + bbac) / a
		if temp < tMax && temp > tMin {
			var rec hit
			rec.t = temp
			rec.p = r.pointAtParam(rec.t)
			rec.n = rec.p.Sub(m.center(r.time))
			rec.n = rec.n.Mul(1.0 / m.Radius)
			rec.m = m.Material
			return true, rec
		}
	}

	return false, hit{}
}

// EPSILON constant for checking if two float numbers are the same
const EPSILON = 0.0000001

func (m *movingSphere) center(time float64) mgl64.Vec3 {
	if m.Time1-m.Time0 > EPSILON {
		tmp := (time - m.Time0) / (m.Time1 - m.Time0)
		tmpVec := m.Center1.Sub(m.Center0)
		tmpVec = tmpVec.Mul(tmp)
		return m.Center0.Add(tmpVec)
	}

	return m.Center0
}
