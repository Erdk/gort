// gort renderer
// Copyright (C) 2017 Erdk <mr.erdk@gmail.com>
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
// Copyright © 2017 Erdk <mr.erdk@gmail.com>

package rayengine

import (
	"fmt"
	"image"
	"math"
	"os"

	"path/filepath"

	"image/jpeg"
	"image/png"
)

type texture interface {
	value(u, v float64, p *Vec) (float64, float64, float64)
}

type constantTexture struct {
	color *Vec
}

func (c *constantTexture) value(u, v float64, p *Vec) (float64, float64, float64) {
	return c.color[0], c.color[1], c.color[2]
}

type checkerTexture struct {
	odd, even texture
}

func (c *checkerTexture) value(u, v float64, p *Vec) (float64, float64, float64) {
	sines := math.Sin(10.0*p[0]) * math.Sin(10.0*p[1]) * math.Sin(10.0*p[2])
	if sines < 0 {
		return c.odd.value(u, v, p)
	}
	return c.even.value(u, v, p)
}

type noiseTexture struct {
	scale float64
}

func (n *noiseTexture) value(u, v float64, p *Vec) (float64, float64, float64) {
	pT := perlinTurbulance(p.MulSI(n.scale), nil)
	// if n.scale == 0.01 {
	// 	fmt.Printf("u: %v v: %v p: %v perlinT: %v\n", u, v, p, pT)
	// }
	ret := &Vec{0.5, 0.5, 0.5}
	ret.MulSM(1.0 + math.Sin(n.scale*p[1]+5.0*pT))
	return ret[0], ret[1], ret[2]
}

type imageTexture struct {
	minx, miny int
	maxx, maxy int
	tex        image.Image
}

func (i *imageTexture) value(u, v float64, p *Vec) (float64, float64, float64) {
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
	return float64(r) / float64(65536), float64(g) / float64(65536), float64(b) / float64(65536)
}

func getImageTexture(file string) (*imageTexture, error) {
	fd, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("cannot open texture file: %s", file)
	}

	var iT imageTexture
	switch filepath.Ext(file) {
	case ".png":
		iT.tex, err = png.Decode(fd)
		if err != nil {
			return nil, fmt.Errorf("cannot decode texture %s", file)
		}
	case ".jpg", ".jpeg":
		iT.tex, err = jpeg.Decode(fd)
		if err != nil {
			return nil, fmt.Errorf("cannot decode texture %s", file)
		}
	default:
		return nil, fmt.Errorf("cannot read texture: %s", file)
	}

	bounds := iT.tex.Bounds()
	iT.minx = bounds.Min.X
	iT.maxx = bounds.Max.X
	iT.miny = bounds.Min.Y
	iT.maxy = bounds.Max.Y

	return &iT, nil
}
