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
	"math/rand"
)

type box struct {
	Min, Max *Vec
	Faces    hitlist
}

// NewBox returns box bounded by two points, p0 and p1
func NewBox(p0, p1 *Vec, m material) *box {
	var b box
	b.Min = p0
	b.Max = p1
	b.Faces = make([]hitable, 6)

	b.Faces[0] = &xyrect{p0[0], p1[0], p0[1], p1[1], p1[2], m}
	b.Faces[1] = &flipNormals{&xyrect{p0[0], p1[0], p0[1], p1[1], p0[2], m}}

	b.Faces[2] = &xzrect{p0[0], p1[0], p0[2], p1[2], p1[1], m}
	b.Faces[3] = &flipNormals{&xzrect{p0[0], p1[0], p0[2], p1[2], p0[1], m}}

	b.Faces[4] = &yzrect{p0[1], p1[1], p0[2], p1[2], p1[0], m}
	b.Faces[5] = &flipNormals{&yzrect{p0[1], p1[1], p0[2], p1[2], p0[0], m}}

	return &b
}

func (b *box) calcHit(randSource *rand.Rand, r *ray, min, max float64) (bool, hit) {
	return b.Faces.calcHit(randSource, r, min, max)
}

func (b *box) boundingBox(t0, t1 float64) (bool, *aabb) {
	return b.Faces.boundingBox(t0, t1)
}