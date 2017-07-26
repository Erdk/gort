package rayengine

import (
	"math/rand"
)

type World struct {
	Cam  camera
	Objs hitlist
}

func (w *World) calcHit(randSource *rand.Rand, r *ray, tMin, tMax float64) (bool, hit) {
	return w.Objs.calcHit(randSource, r, tMin, tMax)
}

func (w *World) boundingBox(t0, t1 float64) (bool, *aabb) {
	return w.Objs.boundingBox(t0, t1)
}
