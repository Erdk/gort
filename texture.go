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
	minx, miny int
	maxx, maxy int
	tex        image.Image
}

func (i imageTexture) value(u, v float64, p mgl64.Vec3) *mgl64.Vec3 {
	x := u*float64(i.maxx-i.minx) + float64(i.minx)
	y := (1.0-v)*float64(i.maxy-i.miny) + float64(i.minx) - 0.001
	if x < float64(i.minx) {
		x = float64(i.minx)
	}
	if x > float64(i.maxx-1.0) {
		x = float64(i.maxx - 1.0)
	}
	if y < float64(i.miny) {
		y = float64(i.miny)
	}
	if y > float64(i.maxy-1.0) {
		y = float64(i.maxy - 1.0)
	}
	r, g, b, _ := i.tex.At(int(x), int(y)).RGBA()
	return &mgl64.Vec3{float64(r) / float64(65536), float64(g) / float64(65536), float64(b) / float64(65536)}
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
	iT.minx = bounds.Min.X
	iT.maxx = bounds.Max.X
	iT.miny = bounds.Min.Y
	iT.maxy = bounds.Max.Y

	return iT, nil
}
