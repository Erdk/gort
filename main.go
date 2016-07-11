package main

import (
	"fmt"
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

type ray struct {
	origin, direction *mgl64.Vec3
}

func (r *ray) point_at_param(t float64) mgl64.Vec3 {
	return r.origin.Add(r.direction.Mul(t))
}

func hit_sphere(center *mgl64.Vec3, rad float64, r *ray) float64 {
	oc := r.origin.Sub(*center)
	a := r.direction.Dot(*r.direction)
	b := 2.0 * oc.Dot(*r.direction)
	c := oc.Dot(oc) - rad*rad
	discriminant := b*b - 4*a*c
	if discriminant < 0 {
		return -1.0
	} else {
		return (-b - math.Sqrt(discriminant)) / (2.0 * a)
	}
}

func color(r *ray) mgl64.Vec3 {
	t := hit_sphere(&mgl64.Vec3{0.0, 0.0, -1.0}, 0.5, r)

	if t > 0.0 {
		pap := r.point_at_param(t)
		pap = pap.Sub(mgl64.Vec3{0.0, 0.0, -1.0})
		N := pap.Normalize()

		return mgl64.Vec3{0.5 * (N.X() + 1.0), 0.5 * (N.Y() + 1.0), 0.5 * (N.Z() + 1)}
	}

	uv := r.direction.Normalize()
	t = 0.5 * (uv.Y() + 1.0)
	ret := mgl64.Vec3{1.0 * (1.0 - t), 1.0 * (1.0 - t), 1.0 * (1.0 - t)}
	tmp := mgl64.Vec3{0.5 * t, 0.7 * t, 1.0 * t}
	return ret.Add(tmp)
}

func main() {
	nx := 1280
	ny := 640

	fmt.Printf("P3\n%d %d\n255\n", nx, ny)

	lower_left_corner = mgl64.Vec3{-2.0, -1.0, -1.0}
	horizontal        = mgl64.Vec3{4.0, 0.0, 0.0}
	vertical          = mgl64.Vec3{0.0, 2.0, 0.0}
	origin := &mgl64.Vec3{}

	for j := ny - 1; j >= 0; j-- {
		for i := 0; i < nx; i++ {
			u := float64(i) / float64(nx)
			v := float64(j) / float64(ny)

			temp := lower_left_corner.Add(horizontal.Mul(u))
			temp = temp.Add(vertical.Mul(v))
			r := &ray{origin: origin, direction: &temp}
			col := color(r)
			col = col.Mul(255.99)
			fmt.Printf("%v %v %v\n", int(col.X()), int(col.Y()), int(col.Z()))
		}
	}
}
