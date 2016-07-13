package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
	"sync"

	"github.com/go-gl/mathgl/mgl64"
)

type ray struct {
	origin, direction *mgl64.Vec3
}

func (r *ray) point_at_param(t float64) mgl64.Vec3 {
	return r.origin.Add(r.direction.Mul(t))
}

type hit struct {
	t    float64
	p, n mgl64.Vec3
}

type hitable interface {
	calc_hit(r *ray, t_min, t_max float64) (bool, hit)
}

type sphere struct {
	c mgl64.Vec3
	r float64
}

func (s *sphere) calc_hit(r *ray, t_min, t_max float64) (bool, hit) {
	oc := r.origin.Sub(s.c)
	a := r.direction.Dot(*r.direction)
	b := oc.Dot(*r.direction)
	c := oc.Dot(oc) - s.r*s.r

	discriminant := b*b - a*c
	if discriminant > 0 {
		temp := (-b - math.Sqrt(b*b-a*c)) / a
		if temp < t_max && temp > t_min {
			var rec hit
			rec.t = temp
			rec.p = r.point_at_param(rec.t)
			rec.n = rec.p.Sub(s.c)
			rec.n = rec.n.Mul(1.0 / s.r)
			return true, rec
		}

		temp = (-b + math.Sqrt(b*b-a*c)) / a
		if temp < t_max && temp > t_min {
			var rec hit
			rec.t = temp
			rec.p = r.point_at_param(rec.t)
			rec.n = rec.p.Sub(s.c)
			rec.n = rec.n.Mul(1.0 / s.r)
			return true, rec
		}
	}

	return false, hit{}
}

type world struct {
	objs []hitable
}

func (w *world) calc_hit(r *ray, t_min, t_max float64) (bool, hit) {
	var ret_rec hit
	hit_anything := false
	closest_so_far := t_max
	for _, v := range w.objs {
		if h, rec := v.calc_hit(r, t_min, closest_so_far); h {
			hit_anything = true
			closest_so_far = rec.t
			ret_rec = rec
		}
	}

	if hit_anything {
		return true, ret_rec
	} else {
		return false, hit{}
	}
}

func random_in_unit_sphere() mgl64.Vec3 {
	p := mgl64.Vec3{2.0 * rand.Float64(), 2.0 * rand.Float64(), 2.0 * rand.Float64()}
	p = p.Sub(mgl64.Vec3{1.0, 1.0, 1.0})
	for p.Len()*p.Len() >= 1.0 {
		p = mgl64.Vec3{2.0 * rand.Float64(), 2.0 * rand.Float64(), 2.0 * rand.Float64()}
		p = p.Sub(mgl64.Vec3{1.0, 1.0, 1.0})
	}

	return p
}

func ret_color(r *ray, w *world) mgl64.Vec3 {
	if h, rec := w.calc_hit(r, 0.001, math.MaxFloat64); h {
		target := rec.p.Add(rec.n.Add(random_in_unit_sphere()))
		tmp := target.Sub(rec.p)
		ret := ret_color(&ray{&rec.p, &tmp}, w)
		//return mgl64.Vec3{0.5 * (rec.n.X() + 1.0), 0.5 * (rec.n.Y() + 1.0), 0.5 * (rec.n.Z() + 1)}
		return ret.Mul(0.5)
	}

	uv := r.direction.Normalize()
	t := 0.5 * (uv.Y() + 1.0)
	ret := mgl64.Vec3{1.0 - t, 1.0 - t, 1.0 - t}
	tmp := mgl64.Vec3{0.5 * t, 0.7 * t, 1.0 * t}
	return ret.Add(tmp)
}

type viewport struct {
	lower_left_corner, horizontal, vertical, origin mgl64.Vec3
}

func NewVP() *viewport {
	var vp viewport
	vp.lower_left_corner = mgl64.Vec3{-2.0, -1.0, -1.0}
	vp.horizontal = mgl64.Vec3{4.0, 0.0, 0.0}
	vp.vertical = mgl64.Vec3{0.0, 2.0, 0.0}
	vp.origin = mgl64.Vec3{0.0, 0.0, 0.0}

	return &vp
}

func (vp *viewport) get_ray(u, v float64) ray {
	tmp := vp.lower_left_corner.Add(vp.horizontal.Mul(u))
	tmp = tmp.Add(vp.vertical.Mul(v)) //+ v * vertical - origin
	tmp = tmp.Add(vp.origin.Mul(-1))
	return ray{origin: &vp.origin, direction: &tmp}
}

const THREADS = 4

func main() {

	if len(os.Args) != 2 {
		fmt.Printf("Usage: %v <filename>\n", os.Args[0])
		os.Exit(1)
	}

	nx := 1280
	ny := 640
	ns := 50

	vp := NewVP()

	w := &world{}
	w.objs = append(w.objs, &sphere{r: 0.5, c: mgl64.Vec3{0.0, 0.0, -1.0}}, &sphere{r: 100, c: mgl64.Vec3{0.0, -100.5, -1.0}})

	img := image.NewRGBA(image.Rect(0, 0, 1280, 640))

	var wg sync.WaitGroup
	wg.Add(THREADS)

	f := func(x1, x2 int) {
		defer wg.Done()
		for j := 0; j < ny; j++ {
			for i := x1; i < x2; i++ {
				col := mgl64.Vec3{0.0, 0.0, 0.0}
				for s := 0; s < ns; s++ {
					u := (float64(i) + rand.Float64()) / float64(nx)
					v := (float64(j) + rand.Float64()) / float64(ny)
					r := vp.get_ray(u, v)
					col = col.Add(ret_color(&r, w))
				}

				col = col.Mul(1.0 / float64(ns))
				col = mgl64.Vec3{math.Sqrt(col.X()) * 255.99, math.Sqrt(col.Y()) * 255.99, math.Sqrt(col.Z()) * 255.99}
				img.Set(i, ny-j, color.RGBA{uint8(col.X()), uint8(col.Y()), uint8(col.Z()), 255})
			}
		}
	}

	for i := 0; i < THREADS; i++ {
		go f(nx/THREADS*i, nx/THREADS*(i+1))
	}

	fd, _ := os.Create(os.Args[1] + ".png")
	defer fd.Close()

	wg.Wait()
	png.Encode(fd, img)
}
