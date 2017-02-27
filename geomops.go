package main

import (
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

type translate struct {
	Objs   hitable
	Offset mgl64.Vec3
}

func (t translate) calcHit(r *ray, min, max float64) (bool, hit) {
	movedOrigin := r.origin.Sub(t.Offset)
	movedRay := ray{&movedOrigin, r.direction, r.time}

	if decision, hit := t.Objs.calcHit(&movedRay, min, max); decision {
		hit.p = hit.p.Add(t.Offset)
		return true, hit
	}

	return false, hit{}
}

func (t translate) boundingBox(t0, t1 float64) (bool, aabb) {
	if decision, box := t.Objs.boundingBox(t0, t1); decision {
		return true, aabb{box.min.Add(t.Offset), box.max.Add(t.Offset)}
	}

	return false, aabb{}
}

type rotateY struct {
	Obj                hitable
	SinTheta, CosTheta float64
	HasBox             bool
	Box                aabb
}

func NewRotateY(obj hitable, angle float64) rotateY {
	ry := rotateY{}
	ry.Obj = obj
	radians := math.Pi / 180.0 * angle
	ry.SinTheta = math.Sin(radians)
	ry.CosTheta = math.Cos(radians)
	ry.HasBox, ry.Box = obj.boundingBox(0.0, 1.0)
	min := mgl64.Vec3{math.MaxFloat64, math.MaxFloat64, math.MaxFloat64}
	max := mgl64.Vec3{-math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64}

	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			for k := 0; k < 2; k++ {
				x := float64(i)*ry.Box.max.X() + (1.0 - float64(i)*ry.Box.min.X())
				y := float64(j)*ry.Box.max.Y() + (1.0 - float64(j)*ry.Box.min.Y())
				z := float64(k)*ry.Box.max.Z() + (1.0 - float64(k)*ry.Box.min.Z())
				newx := ry.CosTheta*x + ry.SinTheta*z
				newz := -ry.SinTheta*x + ry.CosTheta*z
				tester := mgl64.Vec3{newx, y, newz}
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

	ry.Box = aabb{min, max}

	return ry
}

func (ry rotateY) calcHit(r *ray, min, max float64) (bool, hit) {
	origin := *r.origin
	direction := *r.direction
	origin[0] = ry.CosTheta*r.origin[0] - ry.SinTheta*r.origin[2]
	origin[2] = ry.SinTheta*r.origin[0] + ry.CosTheta*r.origin[2]
	direction[0] = ry.CosTheta*r.direction[0] - ry.SinTheta*r.direction[2]
	direction[2] = ry.SinTheta*r.direction[0] + ry.CosTheta*r.direction[2]
	rotatedRay := ray{&origin, &direction, r.time}
	if decision, rec := ry.Obj.calcHit(&rotatedRay, min, max); decision {
		p := rec.p
		n := rec.n
		p[0] = ry.CosTheta*rec.p[0] + ry.SinTheta*rec.p[2]
		p[2] = -ry.SinTheta*rec.p[0] + ry.CosTheta*rec.p[2]
		n[0] = ry.CosTheta*rec.n[0] + ry.SinTheta*rec.n[2]
		n[2] = -ry.SinTheta*rec.n[0] + ry.CosTheta*rec.n[2]
		rec.p = p
		rec.n = n
		return true, rec
	}

	return false, hit{}
}

func (ry rotateY) boundingBox(t0, t1 float64) (bool, aabb) {
	return ry.HasBox, ry.Box
}
