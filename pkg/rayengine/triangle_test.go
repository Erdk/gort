package rayengine

import (
	"math/rand"
	"testing"
	"time"
)

func TestTriangleCalcHit(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	tri := &triangle{
		&Vec{0.0, 0.0, 0.0},
		&Vec{1.0, 0.0, 0.0},
		&Vec{0.5, 1.0, 0.0},
		newLambertianRGB(1.0, 1.0, 1.0),
		false,
	}

	// perpendicular ray
	r := &ray{
		&Vec{0.5, 0.5, -1.0}, // origin
		&Vec{0, 0, 1},        // direction
		0,                    // time
	}
	ishit, result := tri.calcHit(rand.New(rand.NewSource(time.Now().UnixNano())), r, 0, 10)

	if !ishit {
		// there should be hit
		t.Errorf("no hit for: %v %v", r.origin, r.direction)
	} else {
		// hit should be at point (0.5, 0.5, 0)
		if result.p[0] != 0.5 || result.p[1] != 0.5 || result.p[2] != 0.0 {
			t.Error("hit point don't match: ", result.p)
		}

		// normal should be at (0, 0, 1)
		if result.normal[0] != 0.0 || result.normal[1] != 0.0 || result.normal[2] != -1.0 {
			t.Error("normal don't match: ", result.normal)
		}

		if result.t != 1.0 {
			t.Error("incorrect t: ", result.t)
		}
	}

	// ray from up
	r = &ray{
		&Vec{0.5, 1.0, -1.0}, // origin
		&Vec{0, -0.5, 1},     // direction
		0,                    // time
	}
	ishit, result = tri.calcHit(rand.New(rand.NewSource(time.Now().UnixNano())), r, 0, 10)

	if !ishit {
		// there should be hit
		t.Errorf("no hit for: %v %v ", r.origin, r.direction)
	} else {

		// hit should be at point (0.5, 0.5, 0)
		if result.p[0] != 0.5 || result.p[1] != 0.5 || result.p[2] != 0.0 {
			t.Error("hit point don't match: ", result.p)
		}

		// normal should be at (0, 0, 1)
		if result.normal[0] != 0.0 || result.normal[1] != 0.0 || result.normal[2] != -1.0 {
			t.Error("normal don't match: ", result.normal)
		}

		if result.t != 1.0 {
			t.Error("incorrect t: ", result.t)
		}
	}

	// ray from down
	r = &ray{
		&Vec{0.5, 0.0, -1.0}, // origin
		&Vec{0.0, 0.5, 1},    // direction
		0,                    // time
	}
	ishit, result = tri.calcHit(rand.New(rand.NewSource(time.Now().UnixNano())), r, 0, 10)

	if !ishit {
		// there should be hit
		t.Errorf("no hit for: %v %v", r.origin, r.direction)
	} else {
		// hit should be at point (0.5, 0.5, 0)
		if result.p[0] != 0.5 || result.p[1] != 0.5 || result.p[2] != 0.0 {
			t.Error("hit point don't match: ", result.p)
		}

		// normal should be at (0, 0, 1)
		if result.normal[0] != 0.0 || result.normal[1] != 0.0 || result.normal[2] != -1.0 {
			t.Error("normal don't match: ", result.normal)
		}

		if result.t != 1.0 {
			t.Error("incorrect t: ", result.t)
		}
	}
	// ray from left
	// ray from right
}

func TestTriangleBoundingBox(t *testing.T) {
	tri := &triangle{
		&Vec{0.0, 0.0, 554.0},
		&Vec{555.0, 0, 554.0},
		&Vec{278.5, 555, 554.0},
		newLambertianRGB(1.0, 1.0, 1.0),
		false,
	}

	_, bb := tri.boundingBox(0, 1)

	// X bounds
	if InCloseRange(bb.min[0], 0.0) || InCloseRange(bb.max[0], 555.0) {
		t.Errorf("X bounds not right: %v %v", bb.min[0], bb.max[0])
	}

	// Y bounds
	if InCloseRange(bb.min[1], 0.0) || InCloseRange(bb.max[1], 555.0) {
		t.Errorf("Y bounds not right: %v %v", bb.min[1], bb.max[1])
	}

	// Z bounds
	if InCloseRange(bb.min[2], 554.0) || InCloseRange(bb.max[2], 554.0) {
		t.Errorf("Z bounds not right: %v %v", bb.min[2], bb.max[2])
	}
}
