package main

import (
	"math"
	"math/rand"

	"github.com/go-gl/mathgl/mgl64"
)

type material interface {
	scatter(in ray, rec hit) (decision bool, attenuation *mgl64.Vec3, scattered *ray)
}

type lambertian struct {
	Albedo *mgl64.Vec3
}

func getLambertian(v mgl64.Vec3) lambertian {
	return lambertian{Albedo: &v}
}

func (l lambertian) scatter(in ray, rec hit) (decision bool, attenuation *mgl64.Vec3, scattered *ray) {
	//vec3 target = rec.p + rec.normal + random_in_unit_sphere()
	//scattered =ray(rec.p, target - rec.p)
	target := rec.p.Add(rec.n.Add(randomInUnitSphere()))
	tmp := target.Sub(rec.p)
	scattered = &ray{&rec.p, &tmp, in.time}
	attenuation = l.Albedo
	decision = true

	return
}

func reflectvec(v, n mgl64.Vec3) mgl64.Vec3 {
	tmp := n.Mul(2.0 * v.Dot(n))
	return v.Sub(tmp)
}

type metal struct {
	Albedo *mgl64.Vec3
	Fuzz   float64
}

func getMetal(v mgl64.Vec3, f float64) metal {
	if f >= 1.0 {
		f = 1.0
	}

	return metal{Albedo: &v, Fuzz: f}
}

func (m metal) scatter(in ray, rec hit) (decision bool, attenuation *mgl64.Vec3, scattered *ray) {
	reflected := reflectvec(in.direction.Normalize(), rec.n)
	tmp := randomInUnitSphere()
	tmp = tmp.Mul(m.Fuzz)
	reflected = reflected.Add(tmp)
	scattered = &ray{&rec.p, &reflected, in.time}
	attenuation = m.Albedo
	decision = scattered.direction.Dot(rec.n) > 0.0
	return
}

func refract(v, n mgl64.Vec3, niOverNt float64) (bool, mgl64.Vec3) {
	uv := v.Normalize()
	dt := uv.Dot(n)
	discriminant := 1.0 - niOverNt*niOverNt*(1-dt*dt)
	if discriminant > 0.0 {
		// ni_over_nt * (uv - n * dt) - n * sqrt(discriminant)
		uv = uv.Sub(n.Mul(dt)).Mul(niOverNt)
		uv = uv.Sub(n.Mul(math.Sqrt(discriminant)))
		return true, uv
	}

	return false, mgl64.Vec3{}
}

func schlick(cosine, refIdx float64) float64 {
	r0 := (1.0 - refIdx) / (1 + refIdx)
	r0 = r0 * r0
	return r0 + (1.0-r0)*math.Pow((1.0-cosine), 5.0)
}

type dielectric struct {
	RefIdx float64
}

func (d dielectric) scatter(in ray, rec hit) (decision bool, attenuation *mgl64.Vec3, scattered *ray) {
	reflected := reflectvec(*in.direction, rec.n)
	attenuation = &mgl64.Vec3{1.0, 1.0, 1.0}
	var niOverNt float64
	var cosine float64
	var reflectProbe float64
	var outwardNormal mgl64.Vec3
	if in.direction.Dot(rec.n) > 0.0 {
		outwardNormal = rec.n.Mul(-1.0)
		niOverNt = d.RefIdx
		cosine = d.RefIdx * in.direction.Dot(rec.n) / in.direction.Len()
	} else {
		outwardNormal = rec.n
		niOverNt = 1.0 / d.RefIdx
		cosine = -1.0 * in.direction.Dot(rec.n) * in.direction.Len()
	}

	ifRefract, refracted := refract(*in.direction, outwardNormal, niOverNt)
	if ifRefract {
		reflectProbe = schlick(cosine, d.RefIdx)
	} else {
		reflectProbe = 1.0
	}

	if rand.Float64() < reflectProbe {
		scattered = &ray{&rec.p, &reflected, in.time}
	} else {
		scattered = &ray{&rec.p, &refracted, in.time}
	}

	decision = true
	return
}
