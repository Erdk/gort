package rayengine

import (
	"fmt"
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

/*
#include <glm/glm.hpp>
using namespace::glm;

// orig and dir defines the ray. v0, v1, v2 defines the triangle.
// returns the distance from the ray origin to the intersection or 0.
float
triangle_intersection(const vec3& orig,
                      const vec3& dir,
                      const vec3& v0,
                      const vec3& v1,
                      const vec3& v2) {
    vec3 e1 = v1 - v0;
    vec3 e2 = v2 - v0;
    // Calculate planes normal vector
    vec3 pvec = cross(dir, e2);
    float det = dot(e1, pvec);

    // Ray is parallel to plane
    if (det < 1e-8 && det > -1e-8) {
        return 0;
    }

    float inv_det = 1 / det;
    vec3 tvec = orig - v0;
    float u = dot(tvec, pvec) * inv_det;
    if (u < 0 || u > 1) {
        return 0;
    }

    vec3 qvec = cross(tvec, e1);
    float v = dot(dir, qvec) * inv_det;
    if (v < 0 || u + v > 1) {
        return 0;
    }
    return dot(e2, qvec) * inv_det;
}

*/

func (tri *triangle) calcHit(randSource *rand.Rand, r *ray, min, max float64) (bool, hit) {
	triangleDebugCount++
	var rec hit

	e1 := tri.Y.SubVI(tri.X)
	e2 := tri.Z.SubVI(tri.X)

	pvec := r.direction.CrossI(e2)
	det := e1.Dot(pvec)

	if det < 0.00000001 && det > -0.00000001 {
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
	rec.normal = e1.CrossI(e2).Normalize()
	rec.u = 0
	rec.v = 0
	rec.m = tri.Material
	rec.p = r.pointAtParam(rec.t)

	diff := rec.p[2] - tri.X[2]
	if diff > 0.00000001 || diff < -0.00000001 {
		fmt.Printf("ERROR ERROR rec.p.Z: %v", rec.p)
	}

	return true, rec
}

/*	triangleDebugCount++
	var rec hit
	normal := tri.calculateNormal()
	if tri.FlipN {
		normal.NegM()
	}
	//area2 := planeNormal.Len()

	// check if plane with triangle is parallel to the ray
	nDotRayDirection := normal.Dot(r.direction)
	if math.Abs(nDotRayDirection) < 0.001 {
		if tri.Debug {
			fmt.Println("ray and triangle are parallel, normal: ", normal,
				" direction: ", r.direction, " dot ray direction: ", nDotRayDirection)
		}
		return false, rec // ray and triangle are parallel
	}

	// compute d parameter
	d := normal.Dot(tri.X)

	// compute t
	t := (normal.Dot(r.origin) + d) / nDotRayDirection
	if t < min || t > max {
		//triangle is "behind" the ray
		return false, rec
	}

	var C *Vec // vector perpendicular to the triangle's plane

	// intersection point
	p := r.origin.AddVI(r.direction.MulSI(t))

	edge := tri.Y.SubVI(tri.X)
	vp := p.SubVI(tri.X)
	C = edge.CrossI(vp)
	if int(normal.Dot(C)) < 0 {
		return false, rec
	}

	edge = tri.Z.SubVI(tri.Y)
	vp = p.SubVI(tri.Y)
	C = edge.CrossI(vp)
	if int(normal.Dot(C)) < 0 {
		return false, rec
	}

	edge = tri.X.SubVI(tri.Z)
	vp = p.SubVI(tri.Z)
	C = edge.CrossI(vp)
	if int(normal.Dot(C)) < 0 {
		return false, rec
	}

	if t > 1.0 {
		fmt.Println("n: ", normal, " origin: ", r.origin, " d: ", d, " nDtoRayDir: ", nDotRayDirection)
		fmt.Printf("origin: [%0.2f %0.2f %0.2f] dir [%0.2f %0.2f %0.2f] n: [%0.2f %0.2f %0.2f] t: %.2f p: [%0.2f %0.2f %0.2f]\n", r.origin[0], r.origin[1], r.origin[2], r.direction[0], r.direction[1], r.direction[2], normal[0], normal[1], normal[2], t, p[0], p[1], p[2])
	}
	rec.m = tri.Material
	rec.normal = normal.Normalize()
	rec.t = t
	rec.p = r.pointAtParam(rec.t)
	rec.u = 0
	rec.v = 0

	return true, rec
}
*/

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

	fmt.Println("triangle aabb.min: ", min, " aabb.max: ", max)

	return true, &aabb{min: &min, max: &max}
}
