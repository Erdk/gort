package main

import . "github.com/Erdk/gort/types"

type ray struct {
	origin, direction *Vec
	time              float64
}

func (r *ray) pointAtParam(t float64) *Vec {
	v := &Vec{}
	v[0] = r.origin[0] + r.direction[0]*t
	v[1] = r.origin[1] + r.direction[1]*t
	v[2] = r.origin[2] + r.direction[2]*t
	return v
}
