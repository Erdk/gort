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

// Vec is a vector with implemented main math ops, both mutable and unmutable
type Vec [3]float64

// Copy returns copy of vector
func (v *Vec) Copy() *Vec {
	nv := *v
	return &nv
}

// AddSM adds scalar 'a' to the member of vector in place and returns pointer to the result
func (v *Vec) AddSM(a float64) *Vec {
	v[0] += a
	v[1] += a
	v[2] += a
	return v
}

// AddSI copies vector, adds scalar 'a' to the each member of the new vector and returns pointer to it
func (v *Vec) AddSI(a float64) *Vec {
	nv := *v
	nv[0] += a
	nv[1] += a
	nv[2] += a
	return &nv
}

// AddVM adds vector 'v2' to vector in place and returns pointer to the result
func (v *Vec) AddVM(v2 *Vec) *Vec {
	v[0] += v2[0]
	v[1] += v2[1]
	v[2] += v2[2]
	return v
}

// AddVI copies vector, adds vector 'v2' to the new vector and returns pointer to the result
func (v *Vec) AddVI(v2 *Vec) *Vec {
	nv := *v
	nv[0] += v2[0]
	nv[1] += v2[1]
	nv[2] += v2[2]
	return &nv
}

// SubSM subtracts scalar 'a' from each member of the vector in place and returns pointer to it
func (v *Vec) SubSM(a float64) *Vec {
	v[0] -= a
	v[1] -= a
	v[2] -= a
	return v
}

// SubSI copies vector, subtracts scalar 'a' from the each member and returns pointer to the result
func (v *Vec) SubSI(a float64) *Vec {
	nv := *v
	nv[0] -= a
	nv[1] -= a
	nv[2] -= a
	return &nv
}

// SubVM subtracts vector 'v2' from vector in place and returns pointer to the result
func (v *Vec) SubVM(v2 *Vec) *Vec {
	v[0] -= v2[0]
	v[1] -= v2[1]
	v[2] -= v2[2]
	return v
}

// SubVI copies vector, subtracts vector 'v2' from copy and returns pointer to the result
func (v *Vec) SubVI(v2 *Vec) *Vec {
	nv := *v
	nv[0] -= v2[0]
	nv[1] -= v2[1]
	nv[2] -= v2[2]
	return &nv
}

// MulSM multiplies each member by scalar 'a' in place and returns pointer to the result
func (v *Vec) MulSM(a float64) *Vec {
	v[0] *= a
	v[1] *= a
	v[2] *= a
	return v
}

// MulSI copies vector, multiply each member by scalar 'a' and returns pointer to the result
func (v *Vec) MulSI(a float64) *Vec {
	nv := *v
	nv[0] *= a
	nv[1] *= a
	nv[2] *= a
	return &nv
}

// NegM negates vector in place and returns pointer to it
func (v *Vec) NegM() *Vec {
	v[0] = -v[0]
	v[1] = -v[1]
	v[2] = -v[2]
	return v
}

// NegI copies vector, negates each member and returns pointer to the result
func (v *Vec) NegI() *Vec {
	nv := *v
	nv[0] = -nv[0]
	nv[1] = -nv[1]
	nv[2] = -nv[2]
	return &nv
}

// DivSM divides vector by scalar 'a' in place and returns pointer to the result
func (v *Vec) DivSM(a float64) *Vec {
	v[0] /= a
	v[1] /= a
	v[2] /= a
	return v
}

// Len returns length of the vector
func (v *Vec) Len() float64 {
	return math.Sqrt(v[0]*v[0] + v[1]*v[1] + v[2]*v[2])
}

// LenSQ return squared length of the vector
func (v *Vec) LenSQ() float64 {
	return v[0]*v[0] + v[1]*v[1] + v[2]*v[2]
}

// Normalize normalizes vector in place and returns pointer to the result
func (v *Vec) Normalize() *Vec {
	len := v.Len()
	v[0] /= len
	v[1] /= len
	v[2] /= len
	return v
}

// NormalizeI copies vector, normalizes the copy and returns pointer to it
func (v *Vec) NormalizeI() *Vec {
	len := v.Len()
	nv := &Vec{}
	nv[0] = v[0] / len
	nv[1] = v[1] / len
	nv[2] = v[2] / len
	return nv
}

// Dot returns dot product of vector and vector 'v2'
func (v *Vec) Dot(v2 *Vec) float64 {
	return v[0]*v2[0] + v[1]*v2[1] + v[2]*v2[2]
}

// CrossI returns pointer to the  new vector with a cross product of vector and vector 'v2'
func (v *Vec) CrossI(v2 *Vec) *Vec {
	nv := &Vec{}
	nv[0] = v[1]*v2[2] - v[2]*v2[1]
	nv[1] = v[2]*v2[0] - v[0]*v2[2]
	nv[2] = v[0]*v2[1] - v[1]*v2[0]
	return nv
}
