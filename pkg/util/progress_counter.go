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

import (
	"fmt"
	"sync"
)

// ProgressCounter holds mutex and state of current computed pixels
type ProgressCounter struct {
	counter, max, lastPrinted uint
	mtx                       *sync.Mutex
}

// NewProgressCounter #constructor for ProgressCounter
func NewProgressCounter(pixelNum uint) *ProgressCounter {
	pC := &ProgressCounter{}
	pC.counter = 0
	pC.max = pixelNum
	pC.lastPrinted = 0
	pC.mtx = &sync.Mutex{}

	return pC
}

// IncrementCounter increments current counter by 1 and prints result percent done to the stdut if it's greater that previously printed value
func (p *ProgressCounter) IncrementCounter(count uint) {
	p.mtx.Lock()
	p.counter += count
	newPrinted := uint(float64(p.counter) / float64(p.max) * 100)
	if newPrinted > p.lastPrinted {
		p.lastPrinted = newPrinted
		fmt.Printf("\r%d%%", p.lastPrinted)
	}
	p.mtx.Unlock()
}
