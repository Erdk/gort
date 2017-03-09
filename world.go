package main

import (
	"math/rand"

	"github.com/go-gl/mathgl/mgl64"
)

type world struct {
	Objs hitlist
}

func (w *world) calcHit(randSource *rand.Rand, r *ray, tMin, tMax float64) (bool, hit) {
	return w.Objs.calcHit(randSource, r, tMin, tMax)
}

func (w *world) boundingBox(t0, t1 float64) (bool, aabb) {
	return w.Objs.boundingBox(t0, t1)
}

func perlinTest(w *world) {
	w.Objs = make([]hitable, 2)
	perlinTex := noiseTexture{1.0}
	w.Objs[0] = &sphere{
		Radius:   1000,
		Center:   mgl64.Vec3{0.0, -1000.0, 0.0},
		Material: &lambertian{perlinTex},
	}
	w.Objs[1] = &sphere{
		Radius:   2.0,
		Center:   mgl64.Vec3{0.0, 2.0, 0.0},
		Material: &lambertian{perlinTex},
	}
}

func lightAndRectTest(w *world) {
	w.Objs = make([]hitable, 4)
	perlinTex := noiseTexture{4.0}
	w.Objs[0] = &sphere{
		Radius:   1000,
		Center:   mgl64.Vec3{0.0, -1000.0, 0.0},
		Material: &lambertian{perlinTex},
	}
	w.Objs[1] = &sphere{
		Radius:   2,
		Center:   mgl64.Vec3{0.0, 2.0, 0.0},
		Material: &lambertian{perlinTex},
	}
	w.Objs[2] = &sphere{
		Radius:   2,
		Center:   mgl64.Vec3{0.0, 7.0, 0.0},
		Material: &diffuseLight{constantTexture{mgl64.Vec3{4.0, 4.0, 4.0}}},
	}
	w.Objs[3] = &xyrect{3.0, 5.0, 1.0, 3.0, -2.0, &diffuseLight{constantTexture{mgl64.Vec3{4.0, 4.0, 4.0}}}}
}

func cornellBox(w *world) {
	w.Objs = make([]hitable, 8)
	red := &lambertian{constantTexture{mgl64.Vec3{0.65, 0.05, 0.05}}}
	white := &lambertian{constantTexture{mgl64.Vec3{0.73, 0.73, 0.73}}}
	green := &lambertian{constantTexture{mgl64.Vec3{0.12, 0.45, 0.15}}}
	light := &diffuseLight{constantTexture{mgl64.Vec3{7.0, 7.0, 7.0}}}
	w.Objs[0] = &flipNormals{yzrect{0.0, 555.0, 0.0, 555.0, 555.0, green}}
	w.Objs[1] = &yzrect{0.0, 555.0, 0.0, 555.0, 0.0, red}
	w.Objs[2] = &xzrect{113.0, 443.0, 127.0, 432.0, 554.0, light}
	w.Objs[3] = &flipNormals{xzrect{0.0, 555.0, 0.0, 555.0, 555.0, white}}
	w.Objs[4] = &xzrect{0.0, 555.0, 0.0, 555.0, 0.0, white}
	w.Objs[5] = &flipNormals{xyrect{0.0, 555.0, 0.0, 555.0, 555.0, white}}

	b1 := &translate{NewRotateY(
		NewBox(mgl64.Vec3{0.0, 0.0, 0.0}, mgl64.Vec3{165.0, 165.0, 165.0}, white), -18.0),
		mgl64.Vec3{130.0, 0.0, 65.0}}
	b2 := &translate{NewRotateY(
		NewBox(mgl64.Vec3{0.0, 0.0, 0.0}, mgl64.Vec3{165.0, 330.0, 165.0}, white), 15.0),
		mgl64.Vec3{265.0, 0.0, 295.0}}
	w.Objs[6] = &constantMedium{b1, 0.01, isotropicMaterial{constantTexture{mgl64.Vec3{1.0, 1.0, 1.0}}}}
	w.Objs[7] = &constantMedium{b2, 0.01, isotropicMaterial{constantTexture{mgl64.Vec3{0.0, 0.0, 0.0}}}}
}

func testTexture(w *world) {
	w.Objs = make([]hitable, 7)

	red := &lambertian{constantTexture{mgl64.Vec3{0.65, 0.05, 0.05}}}
	white := &lambertian{constantTexture{mgl64.Vec3{0.73, 0.73, 0.73}}}
	green := &lambertian{constantTexture{mgl64.Vec3{0.12, 0.45, 0.15}}}
	light := &diffuseLight{constantTexture{mgl64.Vec3{4.0, 4.0, 4.0}}}
	w.Objs[0] = &flipNormals{yzrect{0.0, 555.0, 0.0, 555.0, 555.0, green}}
	w.Objs[1] = &yzrect{0.0, 555.0, 0.0, 555.0, 0.0, red}
	w.Objs[2] = &xzrect{113.0, 443.0, 127.0, 432.0, 554.0, light}
	w.Objs[3] = &flipNormals{xzrect{0.0, 555.0, 0.0, 555.0, 555.0, white}}
	w.Objs[4] = &xzrect{0.0, 555.0, 0.0, 555.0, 0.0, white}

	texture, err := getImageTexture("static/earthmap.jpg")
	if err != nil {
		panic("CANNOT LOAD TEXTURE!")
	}
	textureMaterial := lambertian{texture}
	w.Objs[5] = &flipNormals{xyrect{0.0, 555.0, 0.0, 555.0, 555.0, white}}
	w.Objs[6] = &sphere{mgl64.Vec3{275.0, 275.0, 275.0}, 100, &textureMaterial}
}

func generateWorld(w *world) {
	w.Objs = make([]hitable, 500)
	i := 0

	cTexture := checkerTexture{constantTexture{mgl64.Vec3{0.2, 0.3, 0.1}}, constantTexture{mgl64.Vec3{0.9, 0.9, 0.9}}}

	w.Objs[i] = &sphere{
		Radius:   0.5,
		Center:   mgl64.Vec3{0.0, 0.0, -1.0},
		Material: &lambertian{constantTexture{mgl64.Vec3{0.1, 0.2, 0.5}}}}
	i++

	w.Objs[i] = &sphere{
		Radius:   100,
		Center:   mgl64.Vec3{0.0, -100.5, -1.0},
		Material: getLambertian(mgl64.Vec3{0.8, 0.8, 0.0})}
	i++

	w.Objs[i] = &sphere{
		Radius:   0.5,
		Center:   mgl64.Vec3{1.0, 0.0, -1.0},
		Material: getMetal(mgl64.Vec3{0.8, 0.6, 0.2}, 0.3)}
	i++

	w.Objs[i] = &sphere{
		Radius:   0.5,
		Center:   mgl64.Vec3{-1.0, 0.0, -1.0},
		Material: &dielectric{1.5}}
	i++

	w.Objs[i] = &sphere{
		Radius:   -0.45,
		Center:   mgl64.Vec3{-1.0, 0.0, -1.0},
		Material: &dielectric{1.5}}
	i++

	w.Objs[i] = &sphere{
		Radius:   1000.0,
		Center:   mgl64.Vec3{0.0, -1000.0, 0.0},
		Material: &lambertian{cTexture},
	}
	i++

	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			chooseMat := rand.Float64()
			center := mgl64.Vec3{
				float64(a) + 0.9*rand.Float64(),
				0.2,
				float64(b) + 0.9*rand.Float64()}

			len := center.Sub(mgl64.Vec3{4.0, 0.2, 0.0}).Len()
			if len > 0.9 {
				if chooseMat < 0.8 { // diffuse
					w.Objs[i] = &movingSphere{
						Radius:  0.2,
						Center0: center,
						Center1: center.Add(mgl64.Vec3{0.0, 0.5 * rand.Float64(), 0.0}),
						Time0:   0.0,
						Time1:   1.0,
						Material: getLambertian(
							mgl64.Vec3{
								rand.Float64() * rand.Float64(),
								rand.Float64() * rand.Float64(),
								rand.Float64() * rand.Float64(),
							}),
					}
				} else if chooseMat < 0.95 { // metal
					w.Objs[i] = &sphere{
						Radius: 0.2,
						Center: center,
						Material: getMetal(
							mgl64.Vec3{
								0.5 + 0.5*rand.Float64(),
								0.5 + 0.5*rand.Float64(),
								0.5 + 0.5*rand.Float64()},
							0.5*rand.Float64()),
					}
				} else { // glass
					w.Objs[i] = &sphere{
						Radius:   0.2,
						Center:   center,
						Material: &dielectric{1.5},
					}
				}

				i++
			}
		}
	}

	w.Objs[i] = &sphere{
		Radius:   1.0,
		Center:   mgl64.Vec3{0.0, 1.0, 0.0},
		Material: &dielectric{1.5}}
	i++
	w.Objs[i] = &sphere{
		Radius:   1.0,
		Center:   mgl64.Vec3{-4.0, 1.0, 0.0},
		Material: getLambertian(mgl64.Vec3{0.4, 0.2, 0.1})}
	i++
	w.Objs[i] = &sphere{
		Radius:   1.0,
		Center:   mgl64.Vec3{4.0, 1.0, 0.0},
		Material: getMetal(mgl64.Vec3{0.7, 0.6, 0.5}, 0.0)}
}

func generateWorld2(w *world) {
	const nb = 20
	l := 0
	w.Objs = make([]hitable, 30)
	white := getLambertian(mgl64.Vec3{0.73, 0.73, 0.73})
	ground := getLambertian(mgl64.Vec3{0.48, 0.83, 0.53})
	light := diffuseLight{constantTexture{mgl64.Vec3{7.0, 7.0, 7.0}}}

	b := 0
	boxlist := make([]hitable, nb*nb)
	for i := 0; i < nb; i++ {
		for j := 0; j < nb; j++ {
			w := 100.0
			x0 := -1000.0 + float64(i)*w
			z0 := -1000.0 + float64(j)*w
			y0 := 0.0
			x1 := x0 + w
			y1 := 100.0 * (rand.Float64() + 0.01)
			z1 := z0 + w
			boxlist[b] = NewBox(mgl64.Vec3{x0, y0, z0}, mgl64.Vec3{x1, y1, z1}, ground)
			b = b + 1
		}
	}

	w.Objs[l] = bvhNodeInit(boxlist, b, 0.0, 1.0)
	l++

	w.Objs[l] = &xzrect{123.0, 423.0, 147.0, 412.0, 554.0, &light}
	l++

	center := mgl64.Vec3{400.0, 400.0, 400.0}

	w.Objs[l] = &movingSphere{center, center.Add(mgl64.Vec3{30.0, 0.0, 0.0}), 0.0, 1.0, 50.0, getLambertian(mgl64.Vec3{0.7, 0.3, 0.1})}
	l++

	w.Objs[l] = &sphere{mgl64.Vec3{260.0, 150.0, 45.0}, 50.0, &dielectric{1.5}}
	l++

	w.Objs[l] = &sphere{mgl64.Vec3{0.0, 150.0, 145.0}, 50.0, getMetal(mgl64.Vec3{0.8, 0.8, 0.9}, 10.0)}
	l++

	boundary := &sphere{mgl64.Vec3{360.0, 150.0, 145.0}, 70.0, &dielectric{1.5}}
	w.Objs[l] = boundary
	l++

	w.Objs[l] = &constantMedium{boundary, 0.2, isotropicMaterial{constantTexture{mgl64.Vec3{0.2, 0.4, 0.9}}}}
	l++

	boundary2 := &sphere{mgl64.Vec3{0.0, 0.0, 0.0}, 5000.0, &dielectric{1.5}}
	w.Objs[l] = &constantMedium{boundary2, 0.0001, isotropicMaterial{constantTexture{mgl64.Vec3{1.0, 1.0, 1.0}}}}
	l++

	texture, err := getImageTexture("static/earthmap.jpg")
	if err != nil {
		panic("CANNOT LOAD TEXTURE!")
	}
	textureMaterial := lambertian{texture}
	w.Objs[l] = &sphere{mgl64.Vec3{400.0, 200.0, 400.0}, 100, &textureMaterial}
	l++

	pertex := lambertian{noiseTexture{0.1}}
	w.Objs[l] = &sphere{mgl64.Vec3{220.0, 280.0, 300.0}, 80.0, &pertex}
	l++

	boxlist2 := make([]hitable, 1000)
	for j := 0; j < 1000; j++ {
		boxlist2[j] = &sphere{mgl64.Vec3{160.0 * rand.Float64(), 160.0 * rand.Float64(), 160.0 * rand.Float64()}, 10.0, white}
	}

	w.Objs[l] = translate{NewRotateY(bvhNodeInit(boxlist2, 1000, 0.0, 1.0), 15.0), mgl64.Vec3{-100.0, 270.0, 395.0}}
}
