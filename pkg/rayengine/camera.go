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
	"math/rand"
)

// Camera describes properties of virtual camera which be used as parameters for rendering process
type Camera struct {
	LowerLeftCorner *Vec
	Horizontal      *Vec
	Vertical        *Vec
	Origin          *Vec
	U, V, W         *Vec
	LensRadius      float64
	Time0, Time1    float64
}

func randomInUnitDisk(randSource *rand.Rand) *Vec {
	p := &Vec{2.0*randSource.Float64() - 1.0, 2.0*randSource.Float64() - 1.0, 0.0}
	for p.Dot(p) >= 1.0 {
		p[0] = 2.0*randSource.Float64() - 1.0
		p[1] = 2.0*randSource.Float64() - 1.0
	}

	return p
}

// NewCamera prepares instance of object describing viewport and rendering properties
func NewCamera(lookfrom, lookat, vup *Vec, vfov, aspect, aperture, focusDist, t0, t1 float64) Camera {
	var cam Camera
	cam.LensRadius = aperture / 2.0

	theta := vfov * math.Pi / 180.0
	halfHeight := math.Tan(theta / 2.0)
	halfWidth := aspect * halfHeight

	cam.Time0 = t0
	cam.Time1 = t1
	cam.Origin = lookfrom
	cam.W = lookfrom.SubVI(lookat).Normalize()
	cam.U = vup.CrossI(cam.W).Normalize()
	cam.V = cam.W.CrossI(cam.U)
	cam.LowerLeftCorner = cam.Origin.AddVI(cam.U.MulSI(-halfWidth * focusDist))
	cam.LowerLeftCorner.
		AddVM(cam.V.MulSI(-halfHeight * focusDist)).
		AddVM(cam.W.MulSI(-focusDist))
	cam.Horizontal = cam.U.MulSI(2.0 * halfWidth * focusDist)
	cam.Vertical = cam.V.MulSI(2.0 * halfHeight * focusDist)

	return cam
}

func (cam *Camera) getRay(randSource *rand.Rand, s, t float64) ray {
	rd := randomInUnitDisk(randSource)
	rd = rd.MulSI(cam.LensRadius)
	offset := cam.U.MulSI(rd[0]).AddVM(cam.V.MulSI(rd[1]))

	rayOri := cam.Origin.AddVI(offset)

	rayDir := cam.LowerLeftCorner.
		AddVI(cam.Horizontal.MulSI(s)).
		AddVM(cam.Vertical.MulSI(t)).
		AddVM(cam.Origin.MulSI(-1.0)).
		AddVM(offset.MulSM(-1.0))
	rayTime := cam.Time0 + randSource.Float64()*(cam.Time1-cam.Time0)
	return ray{rayOri, rayDir, rayTime}
}
