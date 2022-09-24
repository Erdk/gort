package util

import (
	"testing"
)

// Dummy test to test pipeline
func TestPreogressCounter(t *testing.T) {
	var pixel_num uint = 100000
	pc := NewProgressCounter(pixel_num)

	if pc.max != pixel_num {
		t.Fatalf(`pc.max expected %v, got %v`, pixel_num, pc.max)
	}
}
