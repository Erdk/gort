package perlin

import (
	"math/rand"
	"time"

	"math"

	. "github.com/Erdk/gort/rayengine/types"
)

var randSource *rand.Rand

func Noise(p *Vec) float64 {
	u := p[0] - math.Floor(p[0])
	v := p[1] - math.Floor(p[1])
	w := p[2] - math.Floor(p[2])
	i := int(math.Floor(p[0]))
	j := int(math.Floor(p[1]))
	k := int(math.Floor(p[2]))
	var c [2][2][2]*Vec
	for di := 0; di < 2; di++ {
		for dj := 0; dj < 2; dj++ {
			for dk := 0; dk < 2; dk++ {
				c[di][dj][dk] = ranVec[permX[(i+di)&255]^permY[(j+dj)&255]^permZ[(k+dk)&255]]
			}
		}
	}

	return perlinInterpolation(c, u, v, w)
}

func Turbulance(p *Vec, depth *int) float64 {
	turbDepth := 7
	if depth != nil {
		turbDepth = *depth
	}

	accum := 0.0
	weight := 1.0
	tempP := p.Copy()
	for i := 0; i < turbDepth; i++ {
		accum += weight * Noise(tempP)
		weight *= 0.5
		tempP.MulSM(2.0)
	}

	return math.Abs(accum)
}

func perlinInterpolation(c [2][2][2]*Vec, u, v, w float64) float64 {
	accum := 0.0
	uu := u * u * (3.0 - 2.0*u)
	vv := v * v * (3.0 - 2.0*v)
	ww := w * w * (3.0 - 2.0*w)
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			for k := 0; k < 2; k++ {
				weightV := &Vec{u - float64(i), v - float64(j), w - float64(k)}
				accum +=
					(float64(i)*uu + (1.0-float64(i))*(1.0-uu)) *
						(float64(j)*vv + (1.0-float64(j))*(1.0-vv)) *
						(float64(k)*ww + (1.0-float64(k))*(1.0-ww)) * c[i][j][k].Dot(weightV)
			}
		}
	}

	return accum
}

var ranVec [256]*Vec
var permX [256]int
var permY [256]int
var permZ [256]int

func init() {
	randSource = rand.New(rand.NewSource(time.Now().UnixNano()))
	// Perlin generate
	for i := range ranVec {
		ranVec[i] = (&Vec{-1.0 + 2.0*randSource.Float64(), -1.0 + 2.0*randSource.Float64(), -1.0 + 2.0*randSource.Float64()}).Normalize()
	}

	// init
	for i := range permX {
		permX[i] = i
		permY[i] = i
		permZ[i] = i
	}

	// shuffle perm X
	for i := len(permX) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		permX[i], permX[j] = permX[j], permX[i]
	}
	// shuffle perm Y
	for i := len(permY) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		permY[i], permY[j] = permY[j], permY[i]
	}
	// shuffle perm Z
	for i := len(permZ) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		permZ[i], permZ[j] = permZ[j], permZ[i]
	}
}
