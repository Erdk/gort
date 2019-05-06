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

// EPSILON constant for checking if two float numbers are the same
const EPSILON = 0.000000001

// InCloseRange checks if two numbers are identical, for floating-point computations
func InCloseRange(a, b float64) bool {
	if a-b < EPSILON && a-b > -EPSILON {
		return true
	}
	return false
}

// Vec is a vector with implemented main math ops, both mutable and unmutable
type Vec [3]float64

func (v *Vec) Copy() *Vec {
	nv := *v
	return &nv
}

// AddSM adds scalar to vector, result will be changed base vector
func (v *Vec) AddSM(a float64) *Vec {
	v[0] += a
	v[1] += a
	v[2] += a
	return v
}

func (v *Vec) AddSI(a float64) *Vec {
	nv := *v
	nv[0] += a
	nv[1] += a
	nv[2] += a
	return &nv
}

// AddVM adds vector to vector, result will be changed base vector
func (v *Vec) AddVM(v2 *Vec) *Vec {
	v[0] += v2[0]
	v[1] += v2[1]
	v[2] += v2[2]
	return v
}

func (v *Vec) AddVI(v2 *Vec) *Vec {
	nv := *v
	nv[0] += v2[0]
	nv[1] += v2[1]
	nv[2] += v2[2]
	return &nv
}

func (v *Vec) SubSM(a float64) *Vec {
	v[0] -= a
	v[1] -= a
	v[2] -= a
	return v
}

func (v *Vec) SubSI(a float64) *Vec {
	nv := *v
	nv[0] -= a
	nv[1] -= a
	nv[2] -= a
	return &nv
}

func (v *Vec) SubVM(v2 *Vec) *Vec {
	v[0] -= v2[0]
	v[1] -= v2[1]
	v[2] -= v2[2]
	return v
}

func (v *Vec) SubVI(v2 *Vec) *Vec {
	nv := *v
	nv[0] -= v2[0]
	nv[1] -= v2[1]
	nv[2] -= v2[2]
	return &nv
}

func (v *Vec) MulSM(a float64) *Vec {
	v[0] *= a
	v[1] *= a
	v[2] *= a
	return v
}

func (v *Vec) MulSI(a float64) *Vec {
	nv := *v
	nv[0] *= a
	nv[1] *= a
	nv[2] *= a
	return &nv
}

func (v *Vec) NegM() *Vec {
	v[0] = -v[0]
	v[1] = -v[1]
	v[2] = -v[2]
	return v
}

func (v *Vec) NegI() *Vec {
	nv := *v
	nv[0] = -nv[0]
	nv[1] = -nv[1]
	nv[2] = -nv[2]
	return &nv
}

func (v *Vec) DivSM(a float64) *Vec {
	v[0] /= a
	v[1] /= a
	v[2] /= a
	return v
}

func (v *Vec) Len() float64 {
	return math.Sqrt(v[0]*v[0] + v[1]*v[1] + v[2]*v[2])
}

func (v *Vec) LenSQ() float64 {
	return v[0]*v[0] + v[1]*v[1] + v[2]*v[2]
}

func (v *Vec) Normalize() *Vec {
	len := v.Len()
	v[0] /= len
	v[1] /= len
	v[2] /= len
	return v
}

func (v *Vec) NormalizeI() *Vec {
	len := v.Len()
	nv := &Vec{}
	nv[0] = v[0] / len
	nv[1] = v[1] / len
	nv[2] = v[2] / len
	return nv
}

func (v *Vec) Dot(v2 *Vec) float64 {
	return v[0]*v2[0] + v[1]*v2[1] + v[2]*v2[2]
}

func (v *Vec) CrossI(v2 *Vec) *Vec {
	nv := &Vec{}
	nv[0] = v[1]*v2[2] - v[2]*v2[1]
	nv[1] = v[2]*v2[0] - v[0]*v2[2]
	nv[2] = v[0]*v2[1] - v[1]*v2[0]

	return nv
}
