package main

import (
	"math/rand"

	. "github.com/Erdk/gort/types"
)

type hit struct {
	t    float64
	u, v float64
	p, n *Vec
	m    material
}

type hitable interface {
	// r: casted ray
	// min: begin of time slice
	// max: end of time slice
	calcHit(randSource *rand.Rand, r *ray, min, max float64) (bool, hit)
	boundingBox(t0, t1 float64) (bool, aabb)
}

type hitlist []hitable

func (h hitlist) calcHit(randSource *rand.Rand, r *ray, min, max float64) (bool, hit) {
	var retRec hit
	hitAnything := false
	closestSoFar := max
	for _, v := range h {
		if v == nil {
			continue
		}
		if h, rec := v.calcHit(randSource, r, min, closestSoFar); h {
			hitAnything = true
			closestSoFar = rec.t
			retRec = rec
		}
	}

	if hitAnything {
		return true, retRec
	}

	return false, hit{}
}

func (h hitlist) boundingBox(t0, t1 float64) (bool, aabb) {
	if len(h) < 1 {
		return false, aabb{}
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
