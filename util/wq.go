package util

import "sync"

type Stripe struct {
	XStart, XEnd, YStart, YEnd int
}

func ParseStripe(stripeString string) (int, int, error) {
	return 16, 16, nil
}

type queue struct {
	q   []Stripe
	mtx *sync.Mutex
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func NewQueue(xMax, yMax, xStripe, yStripe int) *queue {
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

	q.q = make([]Stripe, numXStripes*numYStripes)
	for i := 0; i < numXStripes; i++ {
		for j := 0; j < numYStripes; j++ {
			q.q[i*numYStripes+j] = Stripe{
				XStart: i * xStripe,
				XEnd:   min((i+1)*xStripe, xMax),
				YStart: j * yStripe,
				YEnd:   min((j+1)*yStripe, yMax),
			}
		}
	}

	return q
}

func (q *queue) GetJob() (Stripe, bool) {
	if len(q.q) == 0 {
		return Stripe{}, false
	}

	q.mtx.Lock()
	retStripe := q.q[0]
	q.q = q.q[1:len(q.q)]
	q.mtx.Unlock()

	return retStripe, true
}
