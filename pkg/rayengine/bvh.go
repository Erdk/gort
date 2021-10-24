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
	"math/rand"
	"sort"
)

type bvhNode struct {
	left, right hitable
	box         *aabb
}

type by func(o1, o2 hitable) bool

func (sortDirection by) Sort(objs []hitable) {
	bs := &bvhSorter{objs, sortDirection}
	sort.Sort(bs)
}

func byX(o1, o2 hitable) bool {
	o1Bound, o1Box := o1.boundingBox(0.0, 0.0)
	o2Bound, o2Box := o2.boundingBox(0.0, 0.0)

	if !o1Bound || !o2Bound {
		fmt.Printf("byX: No bounding box for o1: %v and o2: %v\n", o1, o2)
	}

	return o1Box.min[0] < o2Box.min[0]
}
func byY(o1, o2 hitable) bool {
	o1Bound, o1Box := o1.boundingBox(0.0, 0.0)
	o2Bound, o2Box := o2.boundingBox(0.0, 0.0)

	if !o1Bound || !o2Bound {
		fmt.Printf("byY: No bounding box for o1: %v and o2: %v\n", o1, o2)
	}

	return o1Box.min[1] < o2Box.min[1]
}
func byZ(o1, o2 hitable) bool {
	o1Bound, o1Box := o1.boundingBox(0.0, 0.0)
	o2Bound, o2Box := o2.boundingBox(0.0, 0.0)

	if !o1Bound || !o2Bound {
		fmt.Printf("byZ: No bounding box for o1: %v and o2: %v\n", o1, o2)
	}

	return o1Box.min[2] < o2Box.min[2]
}

type bvhSorter struct {
	objs []hitable
	by   func(o1, o2 hitable) bool
}

func (b *bvhSorter) Len() int {
	return len(b.objs)
}

func (b *bvhSorter) Swap(i, j int) {
	b.objs[i], b.objs[j] = b.objs[j], b.objs[i]
}

func (b *bvhSorter) Less(i, j int) bool {
	return b.by(b.objs[i], b.objs[j])
}

func bvhNodeInit(objs []hitable, n int, time0, time1 float64) *bvhNode {
	b := &bvhNode{}
	axis := int(3 * rand.Float64())

	// sort by chosen axis
	switch axis {
	case 0:
		by(byX).Sort(objs)
	case 1:
		by(byY).Sort(objs)
	default:
		by(byZ).Sort(objs)
	}

	// build tree, "real"" obejcts are on leafs, internal nodes are representing bunding boxes
	switch n {
	case 1:
		b.left = objs[0]
		b.right = objs[0]
	case 2:
		b.left = objs[0]
		b.right = objs[1]
	default:
		b.left = bvhNodeInit(objs[:n/2], n/2, time0, time1)
		b.right = bvhNodeInit(objs[n/2:], n-n/2, time0, time1)
	}

	leftBound, leftBox := b.left.boundingBox(time0, time1)
	rightBound, rightBox := b.right.boundingBox(time0, time1)

	if !leftBound || !rightBound {
		fmt.Printf("No bounding box in bvhNodeInit")
	}

	b.box = surroundingBox(leftBox, rightBox)
	return b
}

func (b *bvhNode) calcHit(randSource *rand.Rand, r *ray, min, max float64) (bool, hit) {
	if b.box.hit(r, min, max) {
		hitLeft, recLeft := b.left.calcHit(randSource, r, min, max)
		hitRight, recRight := b.right.calcHit(randSource, r, min, max)

		if hitLeft && hitRight {
			if recLeft.t < recRight.t {
				return true, recLeft
			}

			return true, recRight
		}

		if hitLeft {
			return true, recLeft
		}

		if hitRight {
			return true, recRight
		}
	}

	return false, hit{}
}

func (b *bvhNode) boundingBox(t0, t1 float64) (bool, *aabb) {
	return true, b.box
}
