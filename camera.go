package main

import (
	"math"
	"math/rand"

	"github.com/go-gl/mathgl/mgl64"
)

type viewport struct {
	lowerLeftCorner mgl64.Vec3
	horizontal      mgl64.Vec3
	vertical        mgl64.Vec3
	origin          mgl64.Vec3
	u, v, w         mgl64.Vec3
	lensRadius      float64
}

func randomInUnitDisk() mgl64.Vec3 {
	p := mgl64.Vec3{2.0*rand.Float64() - 1.0, 2.0*rand.Float64() - 1.0, 0.0}
	for p.Dot(p) >= 1.0 {
		p = mgl64.Vec3{2.0*rand.Float64() - 1.0, 2.0*rand.Float64() - 1.0, 0.0}
	}

	return p
}

func newVP(lookfrom, lookat, vup mgl64.Vec3, vfov, aspect, aperture, focusDist float64) *viewport {
	var vp viewport
	vp.lensRadius = aperture / 2.0
	theta := vfov * math.Pi / 180.0
	halfHeight := math.Tan(theta / 2.0)
	halfWidth := aspect * halfHeight

	vp.origin = lookfrom
	vp.w = lookfrom.Add(lookat.Mul(-1.0))
	vp.w = vp.w.Normalize()
	vp.u = vup.Cross(vp.w)
	vp.u = vp.u.Normalize()
	vp.v = vp.w.Cross(vp.u)
	//vp.lowerLeftCorner = mgl64.Vec3{-half_width, -half_height, -1.0}
	vp.lowerLeftCorner = vp.origin.Add(vp.u.Mul(-halfWidth * focusDist))
	vp.lowerLeftCorner = vp.lowerLeftCorner.Add(vp.v.Mul(-halfHeight * focusDist))
	vp.lowerLeftCorner = vp.lowerLeftCorner.Add(vp.w.Mul(-focusDist))
	vp.horizontal = vp.u.Mul(2.0 * halfWidth * focusDist)
	vp.vertical = vp.v.Mul(2.0 * halfHeight * focusDist)

	return &vp
}

func (vp *viewport) getRay(s, t float64) ray {
	rd := randomInUnitDisk()
	rd = rd.Mul(vp.lensRadius)
	offset := vp.u.Mul(rd.X())
	offset = offset.Add(vp.v.Mul(rd.Y()))

	rayOri := vp.origin.Add(offset)

	rayDir := vp.lowerLeftCorner.Add(vp.horizontal.Mul(s))
	rayDir = rayDir.Add(vp.vertical.Mul(t)) //+ v * vertical - origin
	rayDir = rayDir.Add(vp.origin.Mul(-1.0))
	rayDir = rayDir.Add(offset.Mul(-1.0))
	return ray{origin: &rayOri, direction: &rayDir}
}
