package main

import (
	"math"
	"math/rand"

	"github.com/go-gl/mathgl/mgl64"
)

type sphere struct {
	Center   mgl64.Vec3
	Radius   float64
	Material material
}

func (s *sphere) calcHit(randSource *rand.Rand, r *ray, tMin, tMax float64) (bool, hit) {
	oc := r.origin.Sub(s.Center)
	a := r.direction.Dot(*r.direction)
	b := oc.Dot(*r.direction)
	c := oc.Dot(oc) - s.Radius*s.Radius

	discriminant := b*b - a*c
	if discriminant > 0 {
		var rec hit
		bbac := math.Sqrt(b*b - a*c)
		temp := (-b - bbac) / a
		if temp < tMax && temp > tMin {
			rec.t = temp
			rec.p = r.pointAtParam(rec.t)
			rec.n = rec.p.Sub(s.Center)
			rec.n = rec.n.Mul(1.0 / s.Radius)
			rec.m = s.Material
			rec.u, rec.v = getSphereUV((rec.p.Sub(s.Center)).Mul(1.0 / s.Radius))
			//fmt.Printf("sphere u: %v v: %v\n", rec.u, rec.v)
			return true, rec
		}

		temp = (-b + bbac) / a
		if temp < tMax && temp > tMin {
			rec.t = temp
			rec.p = r.pointAtParam(rec.t)
			rec.n = rec.p.Sub(s.Center)
			rec.n = rec.n.Mul(1.0 / s.Radius)
			rec.m = s.Material
			rec.u, rec.v = getSphereUV((rec.p.Sub(s.Center)).Mul(1.0 / s.Radius))
			//fmt.Printf("sphere u: %v v: %v\n", rec.u, rec.v)
			return true, rec
		}
	}

	return false, hit{}
}

func (s *sphere) boundingBox(t0, ti float64) (bool, aabb) {
	return true, aabb{s.Center.Sub(mgl64.Vec3{s.Radius, s.Radius, s.Radius}), s.Center.Add(mgl64.Vec3{s.Radius, s.Radius, s.Radius})}
}

func getSphereUV(p mgl64.Vec3) (u, v float64) {
	phi := math.Atan2(p.Z(), p.X())
	theta := math.Asin(p.Y())
	u = 1.0 - (phi+math.Pi)/(2.0*math.Pi)
	v = (theta + math.Pi/2.0) / math.Pi
	return
}

type movingSphere struct {
	Center0, Center1 mgl64.Vec3
	Time0, Time1     float64
	Radius           float64
	Material         material
}

func (m *movingSphere) calcHit(randSource *rand.Rand, r *ray, tMin, tMax float64) (bool, hit) {
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

func (m *movingSphere) boundingBox(t0, ti float64) (bool, aabb) {
	aabb0 := aabb{m.Center0.Sub(mgl64.Vec3{m.Radius, m.Radius, m.Radius}), m.Center0.Add(mgl64.Vec3{m.Radius, m.Radius, m.Radius})}
	aabb1 := aabb{m.Center1.Sub(mgl64.Vec3{m.Radius, m.Radius, m.Radius}), m.Center1.Add(mgl64.Vec3{m.Radius, m.Radius, m.Radius})}

	return true, surroundingBox(aabb0, aabb1)
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

func surroundingBox(box0, box1 aabb) aabb {
	small := mgl64.Vec3{
		math.Min(box0.min.X(), box1.min.X()),
		math.Min(box0.min.Y(), box1.min.Y()),
		math.Min(box0.min.Z(), box1.min.Z()),
	}
	large := mgl64.Vec3{
		math.Max(box0.max.X(), box1.max.X()),
		math.Max(box0.max.Y(), box1.max.Y()),
		math.Max(box0.max.Z(), box1.max.Z()),
	}
	return aabb{small, large}
}
