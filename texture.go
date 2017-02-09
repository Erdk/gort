package main

import (
	"math"

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
