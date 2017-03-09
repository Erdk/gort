package main

import (
	"math/rand"

	"github.com/go-gl/mathgl/mgl64"
)

type box struct {
	Min, Max mgl64.Vec3
	Faces    hitlist
}

// NewBox returns box bounded by two points, p0 and p1
func NewBox(p0, p1 mgl64.Vec3, m material) *box {
	var b box
	b.Min = p0
	b.Max = p1
	b.Faces = make([]hitable, 6)

	b.Faces[0] = xyrect{p0.X(), p1.X(), p0.Y(), p1.Y(), p1.Z(), m}
	b.Faces[1] = flipNormals{xyrect{p0.X(), p1.X(), p0.Y(), p1.Y(), p0.Z(), m}}

	b.Faces[2] = xzrect{p0.X(), p1.X(), p0.Z(), p1.Z(), p1.Y(), m}
	b.Faces[3] = flipNormals{xzrect{p0.X(), p1.X(), p0.Z(), p1.Z(), p0.Y(), m}}

	b.Faces[4] = yzrect{p0.Y(), p1.Y(), p0.Z(), p1.Z(), p1.X(), m}
	b.Faces[5] = flipNormals{yzrect{p0.Y(), p1.Y(), p0.Z(), p1.Z(), p0.X(), m}}

	return &b
}

func (b box) calcHit(randSource *rand.Rand, r *ray, min, max float64) (bool, hit) {
	return b.Faces.calcHit(randSource, r, min, max)
}

func (b box) boundingBox(t0, t1 float64) (bool, aabb) {
	return b.Faces.boundingBox(t0, t1)
}
