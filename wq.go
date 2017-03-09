package main

import "sync"

type stripe struct {
	xStart, xEnd, yStart, yEnd int
}

func parseStripe(stripeString string) (int, int, error) {
	return 16, 16, nil
}

type queue struct {
	q   []stripe
	mtx *sync.Mutex
}

func newQueue(xMax, yMax, xStripe, yStripe int) *queue {
	q := &queue{}
	q.mtx = &sync.Mutex{}

	numXStripes := xMax / xStripe
	if xMax%xStripe != 0 {
		numXStripes++
	}

	numYStripes := yMax / yStripe
	if yMax%yStripe != 0 {
		numYStripes++
	}

	q.q = make([]stripe, numXStripes*numYStripes)
	for i := 0; i < numXStripes; i++ {
		for j := 0; j < numYStripes; j++ {
			xBound := (i + 1) * xStripe
			if xBound > xMax {
				xBound = xMax
			}
			yBound := (j + 1) * yStripe
			if yBound > yMax {
				yBound = yMax
			}
			q.q[i*numYStripes+j] = stripe{
				xStart: i * xStripe,
				xEnd:   xBound,
				yStart: j * yStripe,
				yEnd:   yBound,
			}
		}
	}

	return q
}

func (q *queue) getJob() (stripe, bool) {
	if len(q.q) == 0 {
		return stripe{}, false
	}

	q.mtx.Lock()
	retStripe := q.q[0]
	q.q = q.q[1:len(q.q)]
	q.mtx.Unlock()

	return retStripe, true
}
