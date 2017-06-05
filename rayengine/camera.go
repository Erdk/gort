package rayengine

import (
	"math"
	"math/rand"

	. "github.com/Erdk/gort/rayengine/types"
)

type camera struct {
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

func NewCamera(lookfrom, lookat, vup *Vec, vfov, aspect, aperture, focusDist, t0, t1 float64) *camera {
	var cam camera
	cam.lensRadius = aperture / 2.0

	theta := vfov * math.Pi / 180.0
	halfHeight := math.Tan(theta / 2.0)
	halfWidth := aspect * halfHeight

	cam.time0 = t0
	cam.time1 = t1
	cam.origin = lookfrom
	cam.w = lookfrom.SubVI(lookat).Normalize()
	cam.u = vup.CrossI(cam.w).Normalize()
	cam.v = cam.w.CrossI(cam.u)
	cam.lowerLeftCorner = cam.origin.AddVI(cam.u.MulSI(-halfWidth * focusDist))
	cam.lowerLeftCorner.
		AddVM(cam.v.MulSI(-halfHeight * focusDist)).
		AddVM(cam.w.MulSI(-focusDist))
	cam.horizontal = cam.u.MulSI(2.0 * halfWidth * focusDist)
	cam.vertical = cam.v.MulSI(2.0 * halfHeight * focusDist)

	return &cam
}

func (cam *camera) getRay(randSource *rand.Rand, s, t float64) ray {
	rd := randomInUnitDisk(randSource)
	rd = rd.MulSI(cam.lensRadius)
	offset := cam.u.MulSI(rd[0]).AddVM(cam.v.MulSI(rd[1]))

	rayOri := cam.origin.AddVI(offset)

	rayDir := cam.lowerLeftCorner.
		AddVI(cam.horizontal.MulSI(s)).
		AddVM(cam.vertical.MulSI(t)).
		AddVM(cam.origin.MulSI(-1.0)).
		AddVM(offset.MulSM(-1.0))
	rayTime := cam.time0 + randSource.Float64()*(cam.time1-cam.time0)
	return ray{rayOri, rayDir, rayTime}
}
