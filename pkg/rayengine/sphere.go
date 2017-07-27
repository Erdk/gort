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

type sphere struct {
	Center   *Vec
	Radius   float64
	Material material
}

func (s *sphere) calcHit(randSource *rand.Rand, r *ray, tMin, tMax float64) (bool, hit) {
	oc := r.origin.SubVI(s.Center)
	a := r.direction.Dot(r.direction)
	b := oc.Dot(r.direction)
	c := oc.Dot(oc) - s.Radius*s.Radius

	discriminant := b*b - a*c
	if discriminant > 0 {
		var rec hit
		bbac := math.Sqrt(b*b - a*c)
		temp := (-b - bbac) / a
		if temp < tMax && temp > tMin {
			rec.t = temp
			rec.p = r.pointAtParam(rec.t)
			rec.normal = rec.p.SubVI(s.Center).DivSM(s.Radius)
			rec.m = s.Material
			rec.u, rec.v = getSphereUV(rec.normal)
			return true, rec
		}

		temp = (-b + bbac) / a
		if temp < tMax && temp > tMin {
			rec.t = temp
			rec.p = r.pointAtParam(rec.t)
			rec.normal = rec.p.SubVI(s.Center).DivSM(s.Radius)
			rec.m = s.Material
			rec.u, rec.v = getSphereUV(rec.normal)
			return true, rec
		}
	}

	return false, hit{}
}

func (s *sphere) boundingBox(t0, ti float64) (bool, *aabb) {
	return true, &aabb{s.Center.SubSI(s.Radius), s.Center.AddSI(s.Radius)}
}

func getSphereUV(p *Vec) (u, v float64) {
	phi := math.Atan2(p[2], p[0])
	theta := math.Asin(p[1])
	u = 1.0 - (phi+math.Pi)/(2.0*math.Pi)
	v = (theta + math.Pi/2.0) / math.Pi
	return
}

type movingSphere struct {
	Center0, Center1 *Vec
	Time0, Time1     float64
	Radius           float64
	Material         material
}

func (m *movingSphere) calcHit(randSource *rand.Rand, r *ray, tMin, tMax float64) (bool, hit) {
	oc := r.origin.SubVI(m.center(r.time))
	a := r.direction.Dot(r.direction)
	b := oc.Dot(r.direction)
	c := oc.Dot(oc) - m.Radius*m.Radius

	discriminant := b*b - a*c
	if discriminant > 0 {
		bbac := math.Sqrt(b*b - a*c)
		temp := (-b - bbac) / a
		if temp < tMax && temp > tMin {
			var rec hit
			rec.t = temp
			rec.p = r.pointAtParam(rec.t)
			rec.normal = rec.p.SubVI(m.center(r.time)).DivSM(m.Radius)
			rec.m = m.Material
			return true, rec
		}

		temp = (-b + bbac) / a
		if temp < tMax && temp > tMin {
			var rec hit
			rec.t = temp
			rec.p = r.pointAtParam(rec.t)
			rec.normal = rec.p.SubVI(m.center(r.time)).DivSM(m.Radius)
			rec.m = m.Material
			return true, rec
		}
	}

	return false, hit{}
}

func (m *movingSphere) boundingBox(t0, ti float64) (bool, *aabb) {
	aabb0 := &aabb{m.Center0.SubSI(m.Radius), m.Center0.AddSI(m.Radius)}
	aabb1 := &aabb{m.Center1.SubSI(m.Radius), m.Center1.AddSI(m.Radius)}

	return true, surroundingBox(aabb0, aabb1)
}

// EPSILON constant for checking if two float numbers are the same
const EPSILON = 0.0000001

func (m *movingSphere) center(time float64) *Vec {
	if m.Time1-m.Time0 > EPSILON {
		tmp := (time - m.Time0) / (m.Time1 - m.Time0)
		tmpVec := m.Center1.SubVI(m.Center0).MulSM(tmp)
		return m.Center0.AddVI(tmpVec)
	}

	return m.Center0
}

func surroundingBox(box0, box1 *aabb) *aabb {
	small := &Vec{
		math.Min(box0.min[0], box1.min[0]),
		math.Min(box0.min[1], box1.min[1]),
		math.Min(box0.min[2], box1.min[2]),
	}
	large := &Vec{
		math.Max(box0.max[0], box1.max[0]),
		math.Max(box0.max[1], box1.max[1]),
		math.Max(box0.max[2], box1.max[2]),
	}
	return &aabb{small, large}
}
