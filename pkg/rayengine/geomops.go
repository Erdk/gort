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

type translate struct {
	Objs   hitable
	Offset *Vec
}

func (t *translate) calcHit(randSource *rand.Rand, r *ray, min, max float64) (bool, hit) {
	movedOrigin := r.origin.SubVI(t.Offset)
	movedRay := ray{movedOrigin, r.direction, r.time}

	if decision, hit := t.Objs.calcHit(randSource, &movedRay, min, max); decision {
		hit.p.AddVM(t.Offset)
		return true, hit
	}

	return false, hit{}
}

func (t *translate) boundingBox(t0, t1 float64) (bool, *aabb) {
	if decision, box := t.Objs.boundingBox(t0, t1); decision {
		return true, &aabb{box.min.AddVI(t.Offset), box.max.AddVI(t.Offset)}
	}

	return false, nil
}

// RotateY holds reference to original object along precomputed parameters for computing rotation
type RotateY struct {
	Obj                hitable
	SinTheta, CosTheta float64
	HasBox             bool
	Box                *aabb
}

// NewRotateY returns object containing original object rotated by 'angle' over y axis at center of object's bounding box
func NewRotateY(obj hitable, angle float64) *RotateY {
	ry := RotateY{}
	ry.Obj = obj
	radians := math.Pi / 180.0 * angle
	ry.SinTheta = math.Sin(radians)
	ry.CosTheta = math.Cos(radians)
	ry.HasBox, ry.Box = obj.boundingBox(0.0, 1.0)
	min := &Vec{math.MaxFloat64, math.MaxFloat64, math.MaxFloat64}
	max := &Vec{-math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64}

	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			for k := 0; k < 2; k++ {
				x := float64(i)*ry.Box.max[0] + (1.0 - float64(i)*ry.Box.min[0])
				y := float64(j)*ry.Box.max[1] + (1.0 - float64(j)*ry.Box.min[1])
				z := float64(k)*ry.Box.max[2] + (1.0 - float64(k)*ry.Box.min[2])
				newx := ry.CosTheta*x + ry.SinTheta*z
				newz := -ry.SinTheta*x + ry.CosTheta*z
				tester := Vec{newx, y, newz}
				for c := 0; c < 3; c++ {
					if tester[c] > max[c] {
						max[c] = tester[c]
					}
					if tester[c] < min[c] {
						min[c] = tester[c]
					}
				}
			}
		}
	}

	ry.Box = &aabb{min, max}

	return &ry
}

func (ry *RotateY) calcHit(randSource *rand.Rand, r *ray, min, max float64) (bool, hit) {
	origin := *r.origin
	direction := *r.direction
	origin[0] = ry.CosTheta*r.origin[0] - ry.SinTheta*r.origin[2]
	origin[2] = ry.SinTheta*r.origin[0] + ry.CosTheta*r.origin[2]
	direction[0] = ry.CosTheta*r.direction[0] - ry.SinTheta*r.direction[2]
	direction[2] = ry.SinTheta*r.direction[0] + ry.CosTheta*r.direction[2]
	rotatedRay := ray{&origin, &direction, r.time}
	if decision, rec := ry.Obj.calcHit(randSource, &rotatedRay, min, max); decision {
		p := rec.p
		normal := rec.normal
		p[0] = ry.CosTheta*rec.p[0] + ry.SinTheta*rec.p[2]
		p[2] = -ry.SinTheta*rec.p[0] + ry.CosTheta*rec.p[2]
		normal[0] = ry.CosTheta*rec.normal[0] + ry.SinTheta*rec.normal[2]
		normal[2] = -ry.SinTheta*rec.normal[0] + ry.CosTheta*rec.normal[2]
		rec.p = p
		rec.normal = normal
		return true, rec
	}

	return false, hit{}
}

func (ry *RotateY) boundingBox(t0, t1 float64) (bool, *aabb) {
	return ry.HasBox, ry.Box
}
