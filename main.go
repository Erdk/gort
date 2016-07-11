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

type hit_rec struct {
	t    float64
	p, n mgl64.Vec3
}

type hitable interface {
	hit(r *ray, t_min, t_max float64) (bool, hit_rec)
}

type sphere struct {
	c mgl64.Vec3
	r float64
}

func (s *sphere) hit(r *ray, t_min, t_max float64) (bool, hit_rec) {
	oc := r.origin.Sub(s.c)
	a := r.direction.Dot(*r.direction)
	b := oc.Dot(*r.direction)
	c := oc.Dot(oc) - s.r*s.r

	discriminant := b*b - a*c
	if discriminant > 0 {
		temp := (-b - math.Sqrt(b*b-a*c)) / a
		if temp < t_max && temp > t_min {
			var rec hit_rec
			rec.t = temp
			rec.p = r.point_at_param(rec.t)
			rec.n = rec.p.Sub(s.c)
			rec.n = rec.n.Mul(1.0 / s.r)
			return true, rec
		}

		temp = (-b + math.Sqrt(b*b-a*c)) / a
		if temp < t_max && temp > t_min {
			var rec hit_rec
			rec.t = temp
			rec.p = r.point_at_param(rec.t)
			rec.n = rec.p.Sub(s.c)
			rec.n = rec.n.Mul(1.0 / s.r)
			return true, rec
		}
	}

	return false, hit_rec{}
}

/*
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
*/
type world struct {
	objs []hitable
}

func (w *world) hit(r *ray, t_min, t_max float64) (bool, hit_rec) {
	var ret_rec hit_rec
	hit_anything := false
	closest_so_far := t_max
	for _, v := range w.objs {
		//		fmt.Printf("world obj: %v\n", v)
		if h, rec := v.hit(r, t_min, closest_so_far); h {
			hit_anything = true
			closest_so_far = rec.t
			ret_rec = rec
		}
	}

	if hit_anything {
		return true, ret_rec
	} else {
		return false, hit_rec{}
	}
}

func color(r *ray, w *world) mgl64.Vec3 {
	if h, rec := w.hit(r, 0.0, math.MaxFloat64); h {
		return mgl64.Vec3{0.5 * (rec.n.X() + 1.0), 0.5 * (rec.n.Y() + 1.0), 0.5 * (rec.n.Z() + 1)}
	}

	uv := r.direction.Normalize()
	t := 0.5 * (uv.Y() + 1.0)
	ret := mgl64.Vec3{1.0 - t, 1.0 - t, 1.0 - t}
	tmp := mgl64.Vec3{0.5 * t, 0.7 * t, 1.0 * t}
	return ret.Add(tmp)
}

func main() {
	nx := 1280
	ny := 640

	fmt.Printf("P3\n%d %d\n255\n", nx, ny)

	lower_left_corner := mgl64.Vec3{-2.0, -1.0, -1.0}
	horizontal := mgl64.Vec3{4.0, 0.0, 0.0}
	vertical := mgl64.Vec3{0.0, 2.0, 0.0}
	origin := &mgl64.Vec3{0.0, 0.0, 0.0}

	w := &world{}
	w.objs = append(w.objs, &sphere{r: 0.5, c: mgl64.Vec3{0.0, 0.0, -1.0}}, &sphere{r: 100, c: mgl64.Vec3{0.0, -100.5, -1.0}})

	for j := ny - 1; j >= 0; j-- {
		for i := 0; i < nx; i++ {
			u := float64(i) / float64(nx)
			v := float64(j) / float64(ny)

			temp := lower_left_corner.Add(horizontal.Mul(u))
			temp = temp.Add(vertical.Mul(v))
			r := &ray{origin: origin, direction: &temp}

			col := color(r, w)
			col = col.Mul(255.99)
			fmt.Printf("%v %v %v\n", int(col.X()), int(col.Y()), int(col.Z()))
		}
	}
}
