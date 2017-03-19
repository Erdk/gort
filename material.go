package main

import (
	"math"
	"math/rand"

	. "github.com/Erdk/gort/types"
)

type material interface {
	scatter(randSource *rand.Rand, in *ray, rec hit) (decision bool, attenuation *Vec, scattered *ray)
	emit(u, v float64, p *Vec) *Vec
}

type lambertian struct {
	Albedo texture
}

func newLambertianRGB(r, g, b float64) material {
	return &lambertian{constantTexture{&Vec{r, g, b}}}
}

func (l *lambertian) scatter(randSource *rand.Rand, in *ray, rec hit) (decision bool, attenuation *Vec, scattered *ray) {
	in.origin = rec.p
	in.direction = rec.n.AddVM(randomInUnitSphere(randSource))
	scattered = in
	attenuation = l.Albedo.value(rec.u, rec.v, rec.p)
	decision = true
	return
}

func (l *lambertian) emit(u, v float64, p *Vec) *Vec {
	return &Vec{}
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

func (m *metal) scatter(randSource *rand.Rand, in *ray, rec hit) (decision bool, attenuation *Vec, scattered *ray) {
	reflected := reflectvec(in.direction.Normalize(), rec.n).AddVM(randomInUnitSphere(randSource).MulSM(m.Fuzz))
	in.origin = rec.p
	in.direction = reflected
	scattered = in
	attenuation = m.Albedo
	decision = scattered.direction.Dot(rec.n) > 0.0
	return
}

func (m *metal) emit(u, v float64, p *Vec) *Vec {
	return &Vec{}
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

func (d *dielectric) scatter(randSource *rand.Rand, in *ray, rec hit) (decision bool, attenuation *Vec, scattered *ray) {
	reflected := reflectvec(in.direction, rec.n)
	attenuation = d.Attenuation
	var niOverNt float64
	var cosine float64
	var reflectProbe float64
	var outwardNormal *Vec
	if in.direction.Dot(rec.n) > 0.0 {
		outwardNormal = rec.n.MulSI(-1.0)
		niOverNt = d.RefIdx
		cosine = d.RefIdx * in.direction.Dot(rec.n) / in.direction.Len()
	} else {
		outwardNormal = rec.n
		niOverNt = 1.0 / d.RefIdx
		cosine = -1.0 * in.direction.Dot(rec.n) * in.direction.Len()
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

func (d *dielectric) emit(u, v float64, p *Vec) *Vec {
	return &Vec{0.0, 0.0, 0.0}
}

type diffuseLight struct {
	emitTexture texture
}

func newDiffuseLightRGB(r, g, b float64) material {
	return &diffuseLight{constantTexture{&Vec{r, g, b}}}
}

func (d *diffuseLight) scatter(randSource *rand.Rand, in *ray, rec hit) (decision bool, attenuation *Vec, scattered *ray) {
	return false, nil, nil
}

func (d *diffuseLight) emit(u, v float64, p *Vec) *Vec {
	return d.emitTexture.value(u, v, p)
}

type isotropicMaterial struct {
	Albedo texture
}

func newIsotropicMaterialRGB(r, g, b float64) material {
	return &isotropicMaterial{constantTexture{&Vec{r, g, b}}}
}

func (i isotropicMaterial) scatter(randSource *rand.Rand, in *ray, rec hit) (decision bool, attenuation *Vec, scattered *ray) {
	randVec := randomInUnitSphere(randSource)
	scattered = &ray{rec.p, randVec, 0.0}
	attenuation = i.Albedo.value(rec.u, rec.v, rec.p)

	return true, attenuation, scattered
}

func (i isotropicMaterial) emit(u, v float64, p *Vec) *Vec {
	return &Vec{0.0, 0.0, 0.0}
}
