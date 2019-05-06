package rayengine

import (
	"math/rand"
)

var triangleDebugCount int = 0

type triangle struct {
	X, Y, Z  *Vec
	Material material
	FlipN    bool
}

// https://www.khronos.org/opengl/wiki/Calculating_a_Surface_Normal
func (tri *triangle) calculateNormal() *Vec {
	u := tri.Y.SubVI(tri.X)
	v := tri.Z.SubVI(tri.X)

	return u.CrossI(v)
}

func (tri *triangle) calcHit(randSource *rand.Rand, r *ray, min, max float64) (bool, hit) {
	triangleDebugCount++
	var rec hit

	e1 := tri.Y.SubVI(tri.X)
	e2 := tri.Z.SubVI(tri.X)

	pvec := r.direction.CrossI(e2)
	det := e1.Dot(pvec)

	if InCloseRange(det, 0) {
		return false, rec
	}

	invDet := 1.0 / det
	tvec := r.origin.SubVI(tri.X)
	u := tvec.Dot(pvec) * invDet
	if u < 0.0 || u > 1.0 {
		return false, rec
	}

	qvec := tvec.CrossI(e1)
	v := r.direction.Dot(qvec) * invDet
	if v < 0.0 || u+v > 1.0 {
		return false, rec
	}

	rec.t = e2.Dot(qvec) * invDet
	rec.normal = e1.CrossI(e2).Normalize().NegM()
	rec.u = 0
	rec.v = 0
	rec.m = tri.Material
	rec.p = r.pointAtParam(rec.t)

	return true, rec
}

func (tri *triangle) boundingBox(t0, t1 float64) (bool, *aabb) {
	var min, max Vec

	for i := 0; i < 3; i++ {
		if tri.X[i] < tri.Y[i] {
			if tri.X[i] < tri.Z[i] {
				// X < Y, Z?
				min[i] = tri.X[i]
				if tri.Z[i] < tri.Y[i] {
					max[i] = tri.Y[i]
				} else {
					max[i] = tri.Z[i]
				}
			} else {
				// Z < X < Y
				min[i] = tri.Z[i]
				max[i] = tri.Y[i]
			}
		} else {
			// Y < X , Z?
			if tri.X[i] < tri.Z[i] {
				min[i] = tri.Y[i]
				max[i] = tri.Z[i]
				// Y < X < Z
			} else {
				// Y ? Z < X
				max[i] = tri.X[i]
				if tri.Y[i] < tri.Z[i] {
					min[i] = tri.Y[i]
				} else {
					min[i] = tri.Z[i]
				}
			}
		}
	}

	for i := 0; i < 3; i++ {
		min[i] -= 2 * EPSILON
		max[i] += 2 * EPSILON
	}

	return true, &aabb{min: &min, max: &max}
}
