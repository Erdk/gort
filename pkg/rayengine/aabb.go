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
)

type aabb struct {
	min, max *Vec
}

func (a *aabb) hit(r *ray, tmin, tmax float64) bool {
	for i := 0; i < 3; i++ {
		t0 := math.Min((a.min[i]-r.origin[i])/r.direction[i],
			(a.max[i]-r.origin[i])/r.direction[i])
		t1 := math.Max((a.min[i]-r.origin[i])/r.direction[i],
			(a.max[i]-r.origin[i])/r.direction[i])
		tmin = math.Max(t0, tmin)
		tmax = math.Min(t1, tmax)
		if tmax <= tmin {
			return false
		}

	}

	return true
}
