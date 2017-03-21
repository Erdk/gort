package main

import (
	"math"
	"math/rand"

	. "github.com/Erdk/gort/types"
)

type viewport struct {
	lowerLeftCorner *Vec
	horizontal      *Vec
	vertical        *Vec
	origin          *Vec
	u, v, w         *Vec
	lensRadius      float64
	time0, time1    float64
}

func randomInUnitDisk(randSource *rand.Rand) *Vec {
	p := &Vec{2.0*randSource.Float64() - 1.0, 2.0*randSource.Float64() - 1.0, 0.0}
	for p.Dot(p) >= 1.0 {
		p[0] = 2.0*randSource.Float64() - 1.0
		p[1] = 2.0*randSource.Float64() - 1.0
	}

	return p
}

func newVP(lookfrom, lookat, vup *Vec, vfov, aspect, aperture, focusDist, t0, t1 float64) *viewport {
	var vp viewport
	vp.lensRadius = aperture / 2.0
	theta := vfov * math.Pi / 180.0
	halfHeight := math.Tan(theta / 2.0)
	halfWidth := aspect * halfHeight

	vp.time0 = t0
	vp.time1 = t1
	vp.origin = lookfrom
	vp.w = lookfrom.SubVI(lookat).Normalize()
	vp.u = vup.CrossI(vp.w).Normalize()
	vp.v = vp.w.CrossI(vp.u)
	vp.lowerLeftCorner = vp.origin.AddVI(vp.u.MulSI(-halfWidth * focusDist))
	vp.lowerLeftCorner.AddVM(vp.v.MulSI(-halfHeight * focusDist)).AddVM(vp.w.MulSI(-focusDist))
	vp.horizontal = vp.u.MulSI(2.0 * halfWidth * focusDist)
	vp.vertical = vp.v.MulSI(2.0 * halfHeight * focusDist)

	return &vp
}

func (vp *viewport) getRay(randSource *rand.Rand, s, t float64) ray {
	rd := randomInUnitDisk(randSource)
	rd = rd.MulSI(vp.lensRadius)
	offset := vp.u.MulSI(rd[0]).AddVM(vp.v.MulSI(rd[1]))

	rayOri := vp.origin.AddVI(offset)

	rayDir := vp.lowerLeftCorner.AddVI(vp.horizontal.MulSI(s)).AddVM(vp.vertical.MulSI(t)).AddVM(vp.origin.MulSI(-1.0)).AddVM(offset.MulSM(-1.0))
	rayTime := vp.time0 + randSource.Float64()*(vp.time1-vp.time0)
	return ray{rayOri, rayDir, rayTime}
}
