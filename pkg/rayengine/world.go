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

// World holds camera parameters and list of objects
type World struct {
	Cam  Camera
	Objs hitlist
}

func (w *World) calcHit(randSource *rand.Rand, r *ray, tMin, tMax float64) (bool, hit) {
	return w.Objs.calcHit(randSource, r, tMin, tMax)
}

func (w *World) boundingBox(t0, t1 float64) (bool, *aabb) {
	return w.Objs.boundingBox(t0, t1)
}
