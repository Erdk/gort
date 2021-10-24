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

package util

import "sync"

// Patch represents current computing patch of rendered image
type Patch struct {
	XStart, XEnd, YStart, YEnd uint
}

// StringToPatch #wip for now returns 16x16 cmputing dimmensions, in future will parse and validate user provided patch size
func StringToPatch(stripeString string) (uint, uint, error) {
	return 16, 16, nil
}

// Queue holds array patches to compute, single patch -> render coords to compute by single thread
type Queue struct {
	q   []Patch
	mtx *sync.Mutex
}

func min(a, b uint) uint {
	if a < b {
		return a
	}

	return b
}

// NewQueue #constructor for queue, considering compute patch size and whole render image size returns queue with work items with particular coordinates to compute for single thread
func NewQueue(xMax, yMax, xStripe, yStripe uint) *Queue {
	q := &Queue{}
	q.mtx = &sync.Mutex{}

	numXStripes := xMax / xStripe
	if xMax%xStripe != 0 {
		numXStripes++
	}

	numYStripes := yMax / yStripe
	if yMax%yStripe != 0 {
		numYStripes++
	}

	q.q = make([]Patch, numXStripes*numYStripes)
	for i := uint(0); i < numXStripes; i++ {
		for j := uint(0); j < numYStripes; j++ {
			q.q[i*numYStripes+j] = Patch{
				XStart: i * xStripe,
				XEnd:   min((i+1)*xStripe, xMax),
				YStart: j * yStripe,
				YEnd:   min((j+1)*yStripe, yMax),
			}
		}
	}

	return q
}

// GetJob returns next patch to process by thread
func (q *Queue) GetJob() (Patch, bool) {
	q.mtx.Lock()
	defer q.mtx.Unlock()

	if len(q.q) == 0 {
		return Patch{}, false
	}

	retStripe := q.q[0]
	q.q = q.q[1:len(q.q)]
	return retStripe, true
}
