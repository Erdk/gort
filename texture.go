package main

import (
	"fmt"
	"image"
	"math"
	"os"

	"path/filepath"

	"image/jpeg"
	"image/png"

	"github.com/Erdk/gort/perlin"
	"github.com/go-gl/mathgl/mgl64"
)

type texture interface {
	value(u, v float64, p mgl64.Vec3) *mgl64.Vec3
}

type constantTexture struct {
	color mgl64.Vec3
}

func (c constantTexture) value(u, v float64, p mgl64.Vec3) *mgl64.Vec3 {
	return &c.color
}

type checkerTexture struct {
	odd, even texture
}

func (c checkerTexture) value(u, v float64, p mgl64.Vec3) *mgl64.Vec3 {
	sines := math.Sin(10.0*p.X()) * math.Sin(10.0*p.Y()) * math.Sin(10.0*p.Z())
	if sines < 0 {
		return c.odd.value(u, v, p)
	}
	return c.even.value(u, v, p)
}

type noiseTexture struct {
	scale float64
}

func (n noiseTexture) value(u, v float64, p mgl64.Vec3) *mgl64.Vec3 {
	a := mgl64.Vec3{1.0, 1.0, 1.0}.Mul(0.5 * (1.0 + math.Sin(n.scale*p.Z()+10.0*perlin.Turbulance(p))))
	return &a
}

type imageTexture struct {
	nx, ny int
	tex    image.Image
}

func (i imageTexture) value(u, v float64, p mgl64.Vec3) *mgl64.Vec3 {
	x := u * float64(i.nx)
	y := (1.0-v)*float64(i.ny) - 0.001
	if x < 0.0 {
		x = 0.0
	}
	if x > float64(i.nx) {
		x = float64(i.nx)
	}
	if y < 0.0 {
		y = 0.0
	}
	if y > float64(i.ny) {
		y = float64(i.ny)
	}
	r, g, b, _ := i.tex.At(int(x), int(y)).RGBA()
	return &mgl64.Vec3{float64(r & 0xFF), float64(g & 0xFF), float64(b & 0xFF)}
}

func getImageTexture(file string) (imageTexture, error) {
	fd, err := os.Open(file)
	if err != nil {
		return imageTexture{}, fmt.Errorf("cannot open texture file: %s", file)
	}

	var iT imageTexture
	switch filepath.Ext(file) {
	case ".png":
		iT.tex, err = png.Decode(fd)
		if err != nil {
			return imageTexture{}, fmt.Errorf("cannot decode texture %s", file)
		}
	case ".jpg", ".jpeg":
		iT.tex, err = jpeg.Decode(fd)
		if err != nil {
			return imageTexture{}, fmt.Errorf("cannot decode texture %s", file)
		}
	default:
		return imageTexture{}, fmt.Errorf("cannot read texture: %s", file)
	}

	bounds := iT.tex.Bounds()
	iT.nx = bounds.Max.X
	iT.ny = bounds.Max.Y

	return iT, nil
}

type isotropicMaterial struct {
	albedo texture
}

func (i isotropicMaterial) scatter(in ray, rec hit) (decision bool, attenuation *mgl64.Vec3, scattered *ray) {
	randVec := randomInUnitSphere()
	scattered = &ray{&rec.p, &randVec, 0.0}
	attenuation = i.albedo.value(rec.u, rec.v, rec.p)

	return true, attenuation, scattered
}

func (i isotropicMaterial) emit(u, v float64, p mgl64.Vec3) *mgl64.Vec3 {
	return &mgl64.Vec3{0.0, 0.0, 0.0}
}
