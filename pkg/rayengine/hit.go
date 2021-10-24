// gort renderer
// Copyright (C) 2017 Erdk <mr.erdk@gmail.com>
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
// Copyright © 2017 Erdk <mr.erdk@gmail.com>

package rayengine

import (
	"math/rand"
)

type hit struct {
	t         float64
	u, v      float64
	p, normal *Vec
	m         material
}

type hitable interface {
	// r: casted ray
	// min: begin of time slice
	// max: end of time slice
	calcHit(randSource *rand.Rand, r *ray, min, max float64) (bool, hit)
	boundingBox(t0, t1 float64) (bool, *aabb)
}

type hitlist []hitable

func (h hitlist) calcHit(randSource *rand.Rand, r *ray, min, max float64) (bool, hit) {
	var pointOfHit hit
	hitAnything := false
	closestSoFar := max
	for _, v := range h {
		if v == nil {
			continue
		}
		if hitted, rec := v.calcHit(randSource, r, min, closestSoFar); hitted {
			if !InCloseRange(r.origin[0], rec.p[0]) && !InCloseRange(r.origin[1], rec.p[1]) && !InCloseRange(r.origin[2], rec.p[2]) {
				hitAnything = true
				closestSoFar = rec.t
				pointOfHit = rec
			}
		}
	}

	return hitAnything, pointOfHit
}

func (h hitlist) boundingBox(t0, t1 float64) (bool, *aabb) {
	if len(h) < 1 {
		return false, nil
	}

	firstTrue, tempBox := h[0].boundingBox(t0, t1)
	if !firstTrue {
		return false, tempBox
	}
	box := tempBox
	for i := 1; i < len(h); i++ {
		if h[i] == nil {
			break
		}
		ok, tempBox := h[i].boundingBox(t0, t1)
		if ok {
			box = surroundingBox(box, tempBox)
		} else {
			return false, box
		}
	}

	return true, box
}
