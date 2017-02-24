package main

import "github.com/go-gl/mathgl/mgl64"

type hit struct {
	t    float64
	u, v float64
	p, n mgl64.Vec3
	m    material
}

type hitable interface {
	calcHit(r *ray, min, max float64) (bool, hit)
	boundingBox(t0, t1 float64) (bool, aabb)
}
