// gort renderer
// Copyright (C) 2017 Łukasz 'Erdk' Redynk <mr.erdk@gmail.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
// Copyright © 2017 Łukasz 'Erdk' Redynk <mr.erdk@gmail.com>

package rayengine

import (
	"math/rand"
)

// AvailableWorlds returns presets with scenes
func AvailableWorlds() []string {
	return []string{
		"perlin",
		"lightAndRectTest",
		"cornellBox",
		"testTexture",
		"defRoomOneTriangle",
		"generateWorld",
		"generateWorld2"}
}

// NewWorld creates world with camera and objects from preset
func NewWorld(preset string, nx, ny float64) *World {
	w := &World{}
	switch preset {
	case "perlin":
		perlinTest(w, nx, ny)
		return w
	case "lightAndRectTest":
		lightAndRectTest(w)
		return w
	case "cornellBox":
		cornellBox(w, nx, ny)
		return w
	case "testTexture":
		testTexture(w, nx, ny)
		return w
	case "defRoomOneTriangle":
		defRoomOneTriangle(w, nx, ny)
		return w
	case "genWorld":
		generateWorld(w)
		return w
	case "genWorld2":
		generateWorld2(w)
		return w
	}
	return nil
}

/*
 * createDefaultRoom creates simple room with dimmensions
 * 555x555x555, camera set and the middle of XY scene, looking from -700z at 0z
 */
func createDefaultRoom(w *World, nx, ny float64) {
	// setup default camera
	lookFrom := &Vec{278.0, 278.0, -900}
	lookAt := &Vec{278.0, 278.0, 0.0}
	distToFocus := 10.0
	aperture := 0.0
	vFov := 40.0
	w.Cam = NewCamera(lookFrom, lookAt, &Vec{0.0, 1.0, 0.0}, vFov,
		float64(nx)/float64(ny), aperture, distToFocus, 0.0, 1.0)

	// setup room with top light, white walls, floor and ceiling
	w.Objs = make([]hitable, 5)

	// materials
	white := newLambertianRGB(0.73, 0.73, 0.73)
	light := newDiffuseLightRGB(1.0, 1.0, 1.0)

	// room
	w.Objs[0] = &flipNormals{&yzrect{0.0, 555.0, 0.0, 555.0, 555.0, white}}
	w.Objs[1] = &yzrect{0.0, 555.0, 0.0, 555.0, 0.0, white}
	w.Objs[2] = &flipNormals{&xzrect{0.0, 555.0, 0.0, 555.0, 555.0, light}}
	w.Objs[3] = &xzrect{0.0, 555.0, 0.0, 555.0, 0.0, white}
	w.Objs[4] = &flipNormals{&xyrect{0.0, 555.0, 0.0, 555.0, 555.0, white}}
}

func perlinTest(w *World, nx, ny float64) {
	createDefaultRoom(w, nx, ny)

	perlin := &lambertian{&noiseTexture{0.05}, &noiseTexture{0.05}}
	// centered sphere
	w.Objs = append(w.Objs, &sphere{&Vec{278.0, 278.0, 278.0}, 130, perlin})
}

func lightAndRectTest(w *World) {
	w.Objs = make([]hitable, 4)
	perlinTex := &noiseTexture{4.0}
	w.Objs[0] = &sphere{
		Radius:   1000,
		Center:   &Vec{0.0, -1000.0, 0.0},
		Material: &lambertian{perlinTex, &constantTexture{&Vec{0.0, 0.0, 0.0}}},
	}
	w.Objs[1] = &sphere{
		Radius:   2,
		Center:   &Vec{0.0, 2.0, 0.0},
		Material: &lambertian{perlinTex, &constantTexture{&Vec{0.0, 0.0, 0.0}}},
	}
	w.Objs[2] = &sphere{
		Radius:   2,
		Center:   &Vec{0.0, 7.0, 0.0},
		Material: newDiffuseLightRGB(4.0, 4.0, 4.0),
	}
	w.Objs[3] = &xyrect{3.0, 5.0, 1.0, 3.0, -2.0, newDiffuseLightRGB(4.0, 4.0, 4.0)}
}

func cornellBox(w *World, nx, ny float64) {
	lookFrom := &Vec{278.0, 278.0, -700}
	lookAt := &Vec{278.0, 278.0, 0.0}
	distToFocus := 10.0
	aperture := 0.0
	vFov := 40.0
	w.Cam = NewCamera(lookFrom, lookAt, &Vec{0.0, 1.0, 0.0}, vFov, float64(nx)/float64(ny), aperture, distToFocus, 0.0, 1.0)

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
		NewBox(&Vec{0.0, 0.0, 0.0}, &Vec{165.0, 165.0, 165.0}, white), -18.0),
		&Vec{130.0, 0.0, 65.0}}
	b2 := &translate{NewRotateY(
		NewBox(&Vec{0.0, 0.0, 0.0}, &Vec{165.0, 330.0, 165.0}, white), 15.0),
		&Vec{265.0, 0.0, 295.0}}
	w.Objs[6] = b1
	w.Objs[7] = b2
}

func testTexture(w *World, nx, ny float64) {
	createDefaultRoom(w, nx, ny)

	// materials
	texture, err := getImageTexture("static/earthmap.jpg")
	if err != nil {
		panic("CANNOT LOAD TEXTURE!")
	}
	textureMaterial := lambertian{texture, &constantTexture{&Vec{0.0, 0.0, 0.0}}}

	// centered sphere
	w.Objs = append(w.Objs, &sphere{&Vec{278.0, 278.0, 278.0}, 130, &textureMaterial})
}

// Creates test scene with a triangle
func defRoomOneTriangle(w *World, nx, ny float64) {
	createDefaultRoom(w, nx, ny)

	// materials
	red := newLambertianRGB(0.85, 0.05, 0.05)

	// triangle
	// bottom
	v11 := &Vec{60.0, 0.0, 278.5}
	v12 := &Vec{205.0, 0.0, 278.5}
	v13 := &Vec{350.0, 0.0, 278.5}
	v14 := &Vec{495.0, 0.0, 278.5}
	// second tier
	v21 := &Vec{132.5, 125.57, 278.5}
	v22 := &Vec{277.5, 125.57, 278.5}
	v23 := &Vec{422.5, 125.57, 278.5}
	// third tier
	v31 := &Vec{205.0, 251.14, 278.5}
	v32 := &Vec{350.0, 251.14, 278.5}
	// forth tier
	v41 := &Vec{277.5, 376.72, 278.5}

	triangleList := make([]hitable, 6)
	// bottom 3 triangles
	triangleList[0] = &triangle{v11, v12, v21, red, true}
	triangleList[1] = &triangle{v12, v13, v22, red, true}
	triangleList[2] = &triangle{v13, v14, v23, red, true}
	// middle 2 triangles
	triangleList[3] = &triangle{v21, v22, v31, red, true}
	triangleList[4] = &triangle{v22, v23, v32, red, true}
	// top triangle
	triangleList[5] = &triangle{v31, v32, v41, red, true}

	w.Objs = append(w.Objs, bvhNodeInit(triangleList, 6, 0.0, 1.0))
}

func generateWorld(w *World) {
	w.Objs = make([]hitable, 500)
	i := 0

	w.Objs[i] = &sphere{
		Radius:   0.5,
		Center:   &Vec{0.0, 0.0, -1.0},
		Material: newLambertianRGB(0.1, 0.2, 0.5)}
	i++

	w.Objs[i] = &sphere{
		Radius:   100,
		Center:   &Vec{0.0, -100.5, -1.0},
		Material: newLambertianRGB(0.8, 0.8, 0.0)}
	i++

	w.Objs[i] = &sphere{
		Radius:   0.5,
		Center:   &Vec{1.0, 0.0, -1.0},
		Material: newMetalRGB(0.3, 0.8, 0.6, 0.2)}
	i++

	w.Objs[i] = &sphere{
		Radius:   0.5,
		Center:   &Vec{-1.0, 0.0, -1.0},
		Material: newDielectric(1.5)}
	i++

	w.Objs[i] = &sphere{
		Radius:   -0.45,
		Center:   &Vec{-1.0, 0.0, -1.0},
		Material: newDielectric(1.5)}
	i++

	cTexture := &checkerTexture{&constantTexture{&Vec{0.2, 0.3, 0.1}}, &constantTexture{&Vec{0.9, 0.9, 0.9}}}
	w.Objs[i] = &sphere{
		Radius:   1000.0,
		Center:   &Vec{0.0, -1000.0, 0.0},
		Material: &lambertian{cTexture, &constantTexture{&Vec{0.0, 0.0, 0.0}}},
	}
	i++

	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			chooseMat := rand.Float64()
			center := &Vec{
				float64(a) + 0.9*rand.Float64(),
				0.2,
				float64(b) + 0.9*rand.Float64()}

			len := center.SubVI(&Vec{4.0, 0.2, 0.0}).Len()
			if len > 0.9 {
				if chooseMat < 0.8 { // diffuse
					w.Objs[i] = &movingSphere{
						Radius:  0.2,
						Center0: center,
						Center1: center.AddVI(&Vec{0.0, 0.5 * rand.Float64(), 0.0}),
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
		Center:   &Vec{0.0, 1.0, 0.0},
		Material: newDielectric(1.5)}
	i++
	w.Objs[i] = &sphere{
		Radius:   1.0,
		Center:   &Vec{-4.0, 1.0, 0.0},
		Material: newLambertianRGB(0.4, 0.2, 0.1)}
	i++
	w.Objs[i] = &sphere{
		Radius:   1.0,
		Center:   &Vec{4.0, 1.0, 0.0},
		Material: newMetalRGB(0.0, 0.7, 0.6, 0.5)}
}

func generateWorld2(w *World) {
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
			boxlist[b] = NewBox(&Vec{x0, y0, z0}, &Vec{x1, y1, z1}, ground)
			b = b + 1
		}
	}

	// floor
	w.Objs[l] = bvhNodeInit(boxlist, b, 0.0, 1.0)
	l++

	// light
	w.Objs[l] = &xzrect{123.0, 423.0, 147.0, 412.0, 554.0, light}
	l++

	center := &Vec{400.0, 400.0, 400.0}

	w.Objs[l] = &movingSphere{center, center.AddVI(&Vec{30.0, 0.0, 0.0}), 0.0, 1.0, 50.0, newLambertianRGB(0.7, 0.3, 0.1)}
	l++

	w.Objs[l] = &sphere{&Vec{260.0, 150.0, 45.0}, 50.0, newDielectric(1.5)}
	l++

	w.Objs[l] = &sphere{&Vec{0.0, 150.0, 145.0}, 50.0, newMetalRGB(10, 0.8, 0.8, 0.9)}
	l++

	boundary := &sphere{&Vec{360.0, 150.0, 145.0}, 70.0, newDielectric(1.5)}
	w.Objs[l] = boundary
	l++

	w.Objs[l] = &constantMedium{boundary, 0.2, newIsotropicMaterialRGB(0.2, 0.4, 0.9)}
	l++

	boundary2 := &sphere{&Vec{0.0, 0.0, 0.0}, 5000.0, newDielectric(1.5)}
	w.Objs[l] = &constantMedium{boundary2, 0.0001, newIsotropicMaterialRGB(1.0, 1.0, 1.0)}
	l++

	texture, err := getImageTexture("static/earthmap.jpg")
	if err != nil {
		panic("CANNOT LOAD TEXTURE!")
	}
	textureMaterial := lambertian{texture, &constantTexture{&Vec{0.0, 0.0, 0.0}}}
	w.Objs[l] = &sphere{&Vec{400.0, 200.0, 400.0}, 100, &textureMaterial}
	l++

	pertex := lambertian{&noiseTexture{0.1}, &constantTexture{&Vec{0.0, 0.0, 0.0}}}
	w.Objs[l] = &sphere{&Vec{220.0, 280.0, 300.0}, 80.0, &pertex}
	l++

	boxlist2 := make([]hitable, 1000)
	for j := 0; j < 1000; j++ {
		boxlist2[j] = &sphere{&Vec{160.0 * rand.Float64(), 160.0 * rand.Float64(), 160.0 * rand.Float64()}, 10.0, white}
	}

	w.Objs[l] = &translate{NewRotateY(bvhNodeInit(boxlist2, 1000, 0.0, 1.0), 15.0), &Vec{-100.0, 270.0, 395.0}}
}
