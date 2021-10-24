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

type ray struct {
	origin, direction *Vec
	time              float64
}

func (r *ray) pointAtParam(t float64) *Vec {
	v := &Vec{}
	v[0] = r.origin[0] + r.direction[0]*t
	v[1] = r.origin[1] + r.direction[1]*t
	v[2] = r.origin[2] + r.direction[2]*t
	return v
}
