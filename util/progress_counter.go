package util

import (
	"fmt"
	"sync"
)

type ProgressCounter struct {
	counter, max, lastPrinted uint
	mtx                       *sync.Mutex
}

func NewProgressCounter(pixelNum uint) *ProgressCounter {
	pC := &ProgressCounter{}
	pC.counter = 0
	pC.max = pixelNum
	pC.lastPrinted = 0
	pC.mtx = &sync.Mutex{}

	return pC
}

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
