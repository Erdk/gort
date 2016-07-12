package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
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

func ret_color(r *ray, w *world) mgl64.Vec3 {
	if h, rec := w.calc_hit(r, 0.0, math.MaxFloat64); h {
		return mgl64.Vec3{0.5 * (rec.n.X() + 1.0), 0.5 * (rec.n.Y() + 1.0), 0.5 * (rec.n.Z() + 1)}
	}

	uv := r.direction.Normalize()
	t := 0.5 * (uv.Y() + 1.0)
	ret := mgl64.Vec3{1.0 - t, 1.0 - t, 1.0 - t}
	tmp := mgl64.Vec3{0.5 * t, 0.7 * t, 1.0 * t}
	return ret.Add(tmp)
}

func main() {

	if len(os.Args) != 2 {
		fmt.Printf("Usage: %v <filename>\n", os.Args[0])
		os.Exit(1)
	}

	nx := 1280
	ny := 640

	lower_left_corner := mgl64.Vec3{-2.0, -1.0, -1.0}
	horizontal := mgl64.Vec3{4.0, 0.0, 0.0}
	vertical := mgl64.Vec3{0.0, 2.0, 0.0}
	origin := &mgl64.Vec3{0.0, 0.0, 0.0}

	w := &world{}
	w.objs = append(w.objs, &sphere{r: 0.5, c: mgl64.Vec3{0.0, 0.0, -1.0}}, &sphere{r: 100, c: mgl64.Vec3{0.0, -100.5, -1.0}})

	img := image.NewRGBA(image.Rect(0, 0, 1280, 640))

	var wg sync.WaitGroup
	wg.Add(4)

	f := func(x1, x2 int) {
		defer wg.Done()
		for j := 0; j < ny; j++ {
			for i := x1; i < x2; i++ {
				u := float64(i) / float64(nx)
				v := float64(j) / float64(ny)

				temp := lower_left_corner.Add(horizontal.Mul(u))
				temp = temp.Add(vertical.Mul(v))
				r := &ray{origin: origin, direction: &temp}

				col := ret_color(r, w)
				col = col.Mul(255.99)
				img.Set(i, ny-j, color.RGBA{uint8(col.X()), uint8(col.Y()), uint8(col.Z()), 255})
			}
		}
	}

	for i := 0; i < 4; i++ {
		go f(nx/4*i, nx/4*(i+1))
	}

	fd, _ := os.Create(os.Args[1] + ".png")
	defer fd.Close()

	wg.Wait()
	png.Encode(fd, img)
}
