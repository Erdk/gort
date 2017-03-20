package main

import (
	"math/rand"

	. "github.com/Erdk/gort/types"
)

type xyrect struct {
	X0, X1, Y0, Y1, K float64
	Material          material
}

func (xy *xyrect) calcHit(randSource *rand.Rand, r *ray, min, max float64) (bool, hit) {
	var rec hit
	t := (xy.K - r.origin[2]) / r.direction[2]
	if t < min || t > max {
		return false, rec
	}

	x := r.origin[0] + t*r.direction[0]
	y := r.origin[1] + t*r.direction[1]
	if x < xy.X0 || x > xy.X1 || y < xy.Y0 || y > xy.Y1 {
		return false, rec
	}

	rec.u = (x - xy.X0) / (xy.X1 - xy.X0)
	rec.v = (y - xy.Y0) / (xy.Y1 - xy.Y0)
	rec.t = t
	rec.m = xy.Material
	rec.p = r.pointAtParam(t)
	rec.n = &Vec{0.0, 0.0, 1.0}

	return true, rec
}

func (xy *xyrect) boundingBox(t0, t1 float64) (bool, aabb) {
	return true, aabb{min: &Vec{xy.X0, xy.Y0, xy.K - 0.0001},
		max: &Vec{xy.X1, xy.Y1, xy.K + 0.0001}}
}

type xzrect struct {
	X0, X1, Z0, Z1, K float64
	Material          material
}

func (xz *xzrect) calcHit(randSource *rand.Rand, r *ray, min, max float64) (bool, hit) {
	var rec hit
	t := (xz.K - r.origin[1]) / r.direction[1]
	if t < min || t > max {
		return false, rec
	}

	x := r.origin[0] + t*r.direction[0]
	z := r.origin[2] + t*r.direction[2]
	if x < xz.X0 || x > xz.X1 || z < xz.Z0 || z > xz.Z1 {
		return false, rec
	}

	rec.u = (x - xz.X0) / (xz.X1 - xz.X0)
	rec.v = (z - xz.Z0) / (xz.Z1 - xz.Z0)
	rec.t = t
	rec.m = xz.Material
	rec.p = r.pointAtParam(t)
	rec.n = &Vec{0.0, 1.0, 0.0}

	return true, rec
}

func (xz *xzrect) boundingBox(t0, t1 float64) (bool, aabb) {
	return true, aabb{min: &Vec{xz.X0, xz.Z0, xz.K - 0.0001},
		max: &Vec{xz.X1, xz.Z1, xz.K + 0.0001}}
}

type yzrect struct {
	Y0, Y1, Z0, Z1, K float64
	Material          material
}

func (yz *yzrect) calcHit(randSource *rand.Rand, r *ray, min, max float64) (bool, hit) {
	var rec hit
	t := (yz.K - r.origin[0]) / r.direction[0]
	if t < min || t > max {
		return false, rec
	}

	y := r.origin[1] + t*r.direction[1]
	z := r.origin[2] + t*r.direction[2]
	if y < yz.Y0 || y > yz.Y1 || z < yz.Z0 || z > yz.Z1 {
		return false, rec
	}

	rec.u = (y - yz.Y0) / (yz.Y1 - yz.Y0)
	rec.v = (z - yz.Z0) / (yz.Z1 - yz.Z0)
	rec.t = t
	rec.m = yz.Material
	rec.p = r.pointAtParam(t)
	rec.n = &Vec{1.0, 0.0, 0.0}

	return true, rec
}

func (yz *yzrect) boundingBox(t0, t1 float64) (bool, aabb) {
	return true, aabb{min: &Vec{yz.Y0, yz.Z0, yz.K - 0.0001},
		max: &Vec{yz.Y1, yz.Z1, yz.K + 0.0001}}
}

type flipNormals struct {
	h hitable
}

func (f *flipNormals) calcHit(randSource *rand.Rand, r *ray, min, max float64) (bool, hit) {
	dec, rec := f.h.calcHit(randSource, r, min, max)
	if dec {
		rec.n.MulSM(-1.0)
		return dec, rec
	}

	return false, hit{}
}

func (f *flipNormals) boundingBox(t0, t1 float64) (bool, aabb) {
	return f.h.boundingBox(t0, t1)
}
