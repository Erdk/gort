package main

import (
	"math"
	"math/rand"

	. "github.com/Erdk/gort/types"
)

type material interface {
	scatter(randSource *rand.Rand, in *ray, rec hit) (decision bool, attenuationR, attenuationG, attenuationB float64, scattered *ray)
	emit(u, v float64, p *Vec) (float64, float64, float64)
}

type lambertian struct {
	Albedo texture
	Emit   texture
}

func newLambertianRGB(r, g, b float64) material {
	return &lambertian{&constantTexture{&Vec{r, g, b}}, &constantTexture{&Vec{0.0, 0.0, 0.0}}}
}

func (l *lambertian) scatter(randSource *rand.Rand, in *ray, rec hit) (decision bool, attenuationR, attenuationG, attenuationB float64, scattered *ray) {
	in.origin = rec.p
	in.direction = rec.normal.AddVM(randomInUnitSphere(randSource))
	scattered = in
	attenuationR, attenuationG, attenuationB = l.Albedo.value(rec.u, rec.v, rec.p)
	decision = true
	return
}

func (l *lambertian) emit(u, v float64, p *Vec) (float64, float64, float64) {
	return l.Emit.value(u, v, p)
}

func reflectvec(v, n *Vec) *Vec {
	return v.Copy().SubVM(n.MulSI(2.0 * v.Dot(n)))
}

type metal struct {
	Albedo *Vec
	Fuzz   float64
}

func newMetalRGB(fuzz, r, g, b float64) material {
	if fuzz >= 1.0 {
		fuzz = 1.0
	}

	return &metal{Albedo: &Vec{r, g, b}, Fuzz: fuzz}
}

func (m *metal) scatter(randSource *rand.Rand, in *ray, rec hit) (decision bool, attenuationR, attenuationG, attenuationB float64, scattered *ray) {
	reflected := reflectvec(in.direction.Normalize(), rec.normal).AddVM(randomInUnitSphere(randSource).MulSM(m.Fuzz))
	in.origin = rec.p
	in.direction = reflected
	scattered = in
	attenuationR = m.Albedo[0]
	attenuationG = m.Albedo[1]
	attenuationB = m.Albedo[2]
	decision = scattered.direction.Dot(rec.normal) > 0.0
	return
}

func (m *metal) emit(u, v float64, p *Vec) (float64, float64, float64) {
	return 0.0, 0.0, 0.0
}

func refract(v, n *Vec, niOverNt float64) (bool, *Vec) {
	uv := v.NormalizeI()
	dt := uv.Dot(n)
	discriminant := 1.0 - niOverNt*niOverNt*(1-dt*dt)
	if discriminant > 0.0 {
		// ni_over_nt * (uv - n * dt) - n * sqrt(discriminant)
		uv = uv.SubVM(n.MulSI(dt)).MulSM(niOverNt)
		uv = uv.SubVM(n.MulSI(math.Sqrt(discriminant)))
		return true, uv
	}

	return false, &Vec{}
}

func schlick(cosine, refIdx float64) float64 {
	r0 := (1.0 - refIdx) / (1 + refIdx)
	r0 = r0 * r0
	return r0 + (1.0-r0)*math.Pow((1.0-cosine), 5.0)
}

type dielectric struct {
	RefIdx      float64
	Attenuation *Vec
}

func newDielectric(refIdx float64) material {
	return &dielectric{refIdx, &Vec{1.0, 1.0, 1.0}}
}

func newDielectricRGB(refIdx, r, g, b float64) material {
	return &dielectric{refIdx, &Vec{r, g, b}}
}

func (d *dielectric) scatter(randSource *rand.Rand, in *ray, rec hit) (decision bool, attenuationR, attenuationG, attenuationB float64, scattered *ray) {
	reflected := reflectvec(in.direction, rec.normal)
	attenuationR = d.Attenuation[0]
	attenuationG = d.Attenuation[1]
	attenuationB = d.Attenuation[2]
	var niOverNt float64
	var cosine float64
	var reflectProbe float64
	var outwardNormal *Vec
	if in.direction.Dot(rec.normal) > 0.0 {
		outwardNormal = rec.normal.NegI()
		niOverNt = d.RefIdx
		cosine = d.RefIdx * in.direction.Dot(rec.normal) / in.direction.Len()
	} else {
		outwardNormal = rec.normal
		niOverNt = 1.0 / d.RefIdx
		cosine = -1.0 * in.direction.Dot(rec.normal) / in.direction.Len()
	}

	ifRefract, refracted := refract(in.direction, outwardNormal, niOverNt)
	if ifRefract {
		reflectProbe = schlick(cosine, d.RefIdx)
	} else {
		reflectProbe = 1.0
	}

	if randSource.Float64() < reflectProbe {
		scattered = &ray{rec.p, reflected, in.time}
	} else {
		scattered = &ray{rec.p, refracted, in.time}
	}

	decision = true
	return
}

func (d *dielectric) emit(u, v float64, p *Vec) (float64, float64, float64) {
	return 0.0, 0.0, 0.0
}

type diffuseLight struct {
	emitTexture texture
}

func newDiffuseLightRGB(r, g, b float64) material {
	return &diffuseLight{&constantTexture{&Vec{r, g, b}}}
}

func (d *diffuseLight) scatter(randSource *rand.Rand, in *ray, rec hit) (decision bool, attenuationR, attenuationG, attenuationB float64, scattered *ray) {
	return false, 0.0, 0.0, 0.0, nil
}

func (d *diffuseLight) emit(u, v float64, p *Vec) (float64, float64, float64) {
	return d.emitTexture.value(u, v, p)
}

type isotropicMaterial struct {
	Albedo texture
}

func newIsotropicMaterialRGB(r, g, b float64) material {
	return &isotropicMaterial{&constantTexture{&Vec{r, g, b}}}
}

func (i *isotropicMaterial) scatter(randSource *rand.Rand, in *ray, rec hit) (decision bool, attenuationR, attenuationG, attenuationB float64, scattered *ray) {
	randVec := randomInUnitSphere(randSource)
	scattered = &ray{rec.p, randVec, 0.0}
	attenuationR, attenuationG, attenuationB = i.Albedo.value(rec.u, rec.v, rec.p)

	return
}

func (i *isotropicMaterial) emit(u, v float64, p *Vec) (float64, float64, float64) {
	return 0.0, 0.0, 0.0
}
