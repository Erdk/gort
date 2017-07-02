package rayengine

import (
	"math/rand"

	"github.com/Erdk/gort/rayengine/types"
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

func PerlinTest(w *World) {
	w.Objs = make([]hitable, 6)

	// materials
	white := newLambertianRGB(0.73, 0.73, 0.73)
	//light := newDiffuseLightRGB(0.0, 0.0, 0.0)
	perlin := &lambertian{&noiseTexture{0.05}, &noiseTexture{0.05}}

	// objects

	// room
	w.Objs[0] = &flipNormals{&yzrect{0.0, 555.0, 0.0, 555.0, 555.0, white}}
	w.Objs[1] = &yzrect{0.0, 555.0, 0.0, 555.0, 0.0, white}
	w.Objs[2] = &flipNormals{&xzrect{0.0, 555.0, 0.0, 555.0, 555.0, white}}
	w.Objs[3] = &xzrect{0.0, 555.0, 0.0, 555.0, 0.0, white}
	w.Objs[4] = &flipNormals{&xyrect{0.0, 555.0, 0.0, 555.0, 555.0, white}}

	// centered sphere
	w.Objs[5] = &sphere{&types.Vec{278.0, 278.0, 278.0}, 130, perlin}
}

func LightAndRectTest(w *World) {
	w.Objs = make([]hitable, 4)
	perlinTex := &noiseTexture{4.0}
	w.Objs[0] = &sphere{
		Radius:   1000,
		Center:   &types.Vec{0.0, -1000.0, 0.0},
		Material: &lambertian{perlinTex, &constantTexture{&types.Vec{0.0, 0.0, 0.0}}},
	}
	w.Objs[1] = &sphere{
		Radius:   2,
		Center:   &types.Vec{0.0, 2.0, 0.0},
		Material: &lambertian{perlinTex, &constantTexture{&types.Vec{0.0, 0.0, 0.0}}},
	}
	w.Objs[2] = &sphere{
		Radius:   2,
		Center:   &types.Vec{0.0, 7.0, 0.0},
		Material: newDiffuseLightRGB(4.0, 4.0, 4.0),
	}
	w.Objs[3] = &xyrect{3.0, 5.0, 1.0, 3.0, -2.0, newDiffuseLightRGB(4.0, 4.0, 4.0)}
}

func CornellBox(w *World) {
	w.Objs = make([]hitable, 8)
	red := newLambertianRGB(0.65, 0.05, 0.05)
	white := newLambertianRGB(0.73, 0.73, 0.73)
	green := newLambertianRGB(0.12, 0.45, 0.15)
	light := newDiffuseLightRGB(7.0, 7.0, 7.0)
	w.Objs[0] = &flipNormals{&yzrect{0.0, 555.0, 0.0, 555.0, 555.0, green}}
	w.Objs[1] = &yzrect{0.0, 555.0, 0.0, 555.0, 0.0, red}
	w.Objs[2] = &xzrect{213.0, 343.0, 227.0, 332.0, 554.0, light}
	w.Objs[3] = &flipNormals{&xzrect{0.0, 555.0, 0.0, 555.0, 555.0, white}}
	w.Objs[4] = &xzrect{0.0, 555.0, 0.0, 555.0, 0.0, white}
	w.Objs[5] = &flipNormals{&xyrect{0.0, 555.0, 0.0, 555.0, 555.0, white}}

	b1 := &translate{NewRotateY(
		NewBox(&types.Vec{0.0, 0.0, 0.0}, &types.Vec{165.0, 165.0, 165.0}, white), -18.0),
		&types.Vec{130.0, 0.0, 65.0}}
	b2 := &translate{NewRotateY(
		NewBox(&types.Vec{0.0, 0.0, 0.0}, &types.Vec{165.0, 330.0, 165.0}, white), 15.0),
		&types.Vec{265.0, 0.0, 295.0}}
	w.Objs[6] = b1 //&constantMedium{b1, 0.01, newIsotropicMaterialRGB(1.0, 1.0, 1.0)}
	w.Objs[7] = b2 //&constantMedium{b2, 0.01, newIsotropicMaterialRGB(0.0, 0.0, 0.0)}
}

func TestTexture(w *World) {
	w.Objs = make([]hitable, 7)

	// materials
	red := newLambertianRGB(0.65, 0.05, 0.05)
	white := newLambertianRGB(0.73, 0.73, 0.73)
	green := newLambertianRGB(0.12, 0.45, 0.15)
	light := newDiffuseLightRGB(7.0, 7.0, 7.0)
	texture, err := getImageTexture("static/earthmap.jpg")
	if err != nil {
		panic("CANNOT LOAD TEXTURE!")
	}
	textureMaterial := lambertian{texture, &constantTexture{&types.Vec{0.0, 0.0, 0.0}}}

	// objects

	// room
	w.Objs[0] = &flipNormals{&yzrect{0.0, 555.0, 0.0, 555.0, 555.0, green}}
	w.Objs[1] = &yzrect{0.0, 555.0, 0.0, 555.0, 0.0, red}
	w.Objs[2] = &xzrect{113.0, 443.0, 127.0, 432.0, 554.0, light}
	w.Objs[3] = &flipNormals{&xzrect{0.0, 555.0, 0.0, 555.0, 555.0, white}}
	w.Objs[4] = &xzrect{0.0, 555.0, 0.0, 555.0, 0.0, white}
	w.Objs[5] = &flipNormals{&xyrect{0.0, 555.0, 0.0, 555.0, 555.0, white}}

	// centered sphere
	w.Objs[6] = &sphere{&types.Vec{278.0, 278.0, 278.0}, 130, &textureMaterial}
}

// colorVolWorld: generates scene with room and 3 dielectric spheres, middle one contains volume object
func ColorVolWorld(w *World) {
	w.Objs = make([]hitable, 10)

	// materials

	white := newLambertianRGB(0.73, 0.73, 0.73)
	//red := newLambertianRGB(0.65, 0.05, 0.05)
	//blue := newLambertianRGB(0.05, 0.05, 0.65)
	//green := newLambertianRGB(0.12, 0.45, 0.15)

	light := newDiffuseLightRGB(1.0, 1.0, 1.0)

	//earthTexture, err := getImageTexture("static/earthmap.jpg")
	//if err != nil {
	//	panic("CANNOT LOAD TEXTURE!")
	//}
	//earthMat := &lambertian{earthTexture, earthTexture}

	//moonTexture, err := getImageTexture("static/moonmap.jpg")
	//if err != nil {
	//	panic("CANNOT LOAD TEXTURE!")
	//}
	//moonMat := &lambertian{moonTexture, moonTexture}

	// objects

	// room
	w.Objs[0] = &flipNormals{&yzrect{0.0, 555.0, 0.0, 555.0, 555.0, white}}
	w.Objs[1] = &yzrect{0.0, 555.0, 0.0, 555.0, 0.0, white}
	//	w.Objs[2] = &xzrect{113.0, 443.0, 127.0, 432.0, 554.0, light}
	w.Objs[2] = &flipNormals{&xzrect{0.0, 555.0, 0.0, 555.0, 555.0, light}}
	w.Objs[3] = &xzrect{0.0, 555.0, 0.0, 555.0, 0.0, white}
	w.Objs[4] = &flipNormals{&xyrect{0.0, 555.0, 0.0, 555.0, 555.0, white}}
	// Earth
	//w.Objs[5] = &sphere{&Vec{208.0, 208.0, 208.0}, 140, earthMat}
	//	w.Objs[5] = &sphere{&Vec{278.0, 208.0, 278.0}, 140, light}
	// Moom
	//w.Objs[6] = &sphere{&Vec{417.0, 417.0, 417.0}, 40, moonMat}
	//w.Objs[6] = &sphere{&Vec{417.0, 317.0, 217.0}, 40, red}
	//w.Objs[7] = &sphere{&Vec{317.0, 417.0, 217.0}, 40, green}
	//w.Objs[8] = &sphere{&Vec{217.0, 417.0, 317.0}, 40, blue}

	//w.Objs[10] = &triangle{
	//	&Vec{0.0, 0.0, 554.0},
	//	&Vec{278.0, 555.0, 554.0},
	//	&Vec{555.0, 0.0, 554.0},
	//	red,
	//	true,
	//}

	// "mist"
	//boxBoundary := NewBox(&Vec{0.0, 0.0, 0.0}, &Vec{555.0, 555.0, 555.0}, newDielectric(1.5))
	//w.Objs[7] = &constantMedium{boxBoundary, 0.0005, newIsotropicMaterialRGB(0.3, 0.3, 0.3)}
}

func GenerateWorld(w *World) {
	w.Objs = make([]hitable, 500)
	i := 0

	w.Objs[i] = &sphere{
		Radius:   0.5,
		Center:   &types.Vec{0.0, 0.0, -1.0},
		Material: newLambertianRGB(0.1, 0.2, 0.5)}
	i++

	w.Objs[i] = &sphere{
		Radius:   100,
		Center:   &types.Vec{0.0, -100.5, -1.0},
		Material: newLambertianRGB(0.8, 0.8, 0.0)}
	i++

	w.Objs[i] = &sphere{
		Radius:   0.5,
		Center:   &types.Vec{1.0, 0.0, -1.0},
		Material: newMetalRGB(0.3, 0.8, 0.6, 0.2)}
	i++

	w.Objs[i] = &sphere{
		Radius:   0.5,
		Center:   &types.Vec{-1.0, 0.0, -1.0},
		Material: newDielectric(1.5)}
	i++

	w.Objs[i] = &sphere{
		Radius:   -0.45,
		Center:   &types.Vec{-1.0, 0.0, -1.0},
		Material: newDielectric(1.5)}
	i++

	cTexture := &checkerTexture{&constantTexture{&types.Vec{0.2, 0.3, 0.1}}, &constantTexture{&types.Vec{0.9, 0.9, 0.9}}}
	w.Objs[i] = &sphere{
		Radius:   1000.0,
		Center:   &types.Vec{0.0, -1000.0, 0.0},
		Material: &lambertian{cTexture, &constantTexture{&types.Vec{0.0, 0.0, 0.0}}},
	}
	i++

	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			chooseMat := rand.Float64()
			center := &types.Vec{
				float64(a) + 0.9*rand.Float64(),
				0.2,
				float64(b) + 0.9*rand.Float64()}

			len := center.SubVI(&types.Vec{4.0, 0.2, 0.0}).Len()
			if len > 0.9 {
				if chooseMat < 0.8 { // diffuse
					w.Objs[i] = &movingSphere{
						Radius:  0.2,
						Center0: center,
						Center1: center.AddVI(&types.Vec{0.0, 0.5 * rand.Float64(), 0.0}),
						Time0:   0.0,
						Time1:   1.0,
						Material: newLambertianRGB(
							rand.Float64()*rand.Float64(),
							rand.Float64()*rand.Float64(),
							rand.Float64()*rand.Float64(),
						),
					}
				} else if chooseMat < 0.95 { // metal
					w.Objs[i] = &sphere{
						Radius: 0.2,
						Center: center,
						Material: newMetalRGB(
							0.5*rand.Float64(),
							0.5+0.5*rand.Float64(),
							0.5+0.5*rand.Float64(),
							0.5+0.5*rand.Float64(),
						),
					}
				} else { // glass
					w.Objs[i] = &sphere{
						Radius:   0.2,
						Center:   center,
						Material: newDielectric(1.5),
					}
				}

				i++
			}
		}
	}

	w.Objs[i] = &sphere{
		Radius:   1.0,
		Center:   &types.Vec{0.0, 1.0, 0.0},
		Material: newDielectric(1.5)}
	i++
	w.Objs[i] = &sphere{
		Radius:   1.0,
		Center:   &types.Vec{-4.0, 1.0, 0.0},
		Material: newLambertianRGB(0.4, 0.2, 0.1)}
	i++
	w.Objs[i] = &sphere{
		Radius:   1.0,
		Center:   &types.Vec{4.0, 1.0, 0.0},
		Material: newMetalRGB(0.0, 0.7, 0.6, 0.5)}
}

func GenerateWorld2(w *World) {
	const nb = 20
	l := 0
	w.Objs = make([]hitable, 30)

	white := newLambertianRGB(0.73, 0.73, 0.73)
	ground := newLambertianRGB(0.48, 0.83, 0.53)
	light := newDiffuseLightRGB(7.0, 7.0, 7.0)

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
			boxlist[b] = NewBox(&types.Vec{x0, y0, z0}, &types.Vec{x1, y1, z1}, ground)
			b = b + 1
		}
	}

	// floor
	w.Objs[l] = bvhNodeInit(boxlist, b, 0.0, 1.0)
	l++

	// light
	w.Objs[l] = &xzrect{123.0, 423.0, 147.0, 412.0, 554.0, light}
	l++

	center := &types.Vec{400.0, 400.0, 400.0}

	w.Objs[l] = &movingSphere{center, center.AddVI(&types.Vec{30.0, 0.0, 0.0}), 0.0, 1.0, 50.0, newLambertianRGB(0.7, 0.3, 0.1)}
	l++

	w.Objs[l] = &sphere{&types.Vec{260.0, 150.0, 45.0}, 50.0, newDielectric(1.5)}
	l++

	w.Objs[l] = &sphere{&types.Vec{0.0, 150.0, 145.0}, 50.0, newMetalRGB(10, 0.8, 0.8, 0.9)}
	l++

	boundary := &sphere{&types.Vec{360.0, 150.0, 145.0}, 70.0, newDielectric(1.5)}
	w.Objs[l] = boundary
	l++

	w.Objs[l] = &constantMedium{boundary, 0.2, newIsotropicMaterialRGB(0.2, 0.4, 0.9)}
	l++

	boundary2 := &sphere{&types.Vec{0.0, 0.0, 0.0}, 5000.0, newDielectric(1.5)}
	w.Objs[l] = &constantMedium{boundary2, 0.0001, newIsotropicMaterialRGB(1.0, 1.0, 1.0)}
	l++

	texture, err := getImageTexture("static/earthmap.jpg")
	if err != nil {
		panic("CANNOT LOAD TEXTURE!")
	}
	textureMaterial := lambertian{texture, &constantTexture{&types.Vec{0.0, 0.0, 0.0}}}
	w.Objs[l] = &sphere{&types.Vec{400.0, 200.0, 400.0}, 100, &textureMaterial}
	l++

	pertex := lambertian{&noiseTexture{0.1}, &constantTexture{&types.Vec{0.0, 0.0, 0.0}}}
	w.Objs[l] = &sphere{&types.Vec{220.0, 280.0, 300.0}, 80.0, &pertex}
	l++

	boxlist2 := make([]hitable, 1000)
	for j := 0; j < 1000; j++ {
		boxlist2[j] = &sphere{&types.Vec{160.0 * rand.Float64(), 160.0 * rand.Float64(), 160.0 * rand.Float64()}, 10.0, white}
	}

	w.Objs[l] = &translate{NewRotateY(bvhNodeInit(boxlist2, 1000, 0.0, 1.0), 15.0), &types.Vec{-100.0, 270.0, 395.0}}
}
