// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

// RectBounds specifies the size of the boundaries of a rectangle.
// It can represent the thickness of the borders, the margins, or the padding of a rectangle.
type RectBounds struct {
	Top    float32
	Right  float32
	Bottom float32
	Left   float32
}

// Set sets the values of the RectBounds.
func (bs *RectBounds) Set(top, right, bottom, left float32) {

	if top >= 0 {
		bs.Top = top
	}
	if right >= 0 {
		bs.Right = right
	}
	if bottom >= 0 {
		bs.Bottom = bottom
	}
	if left >= 0 {
		bs.Left = left
	}
}

// Rect represents a rectangle.
type Rect struct {
	X      float32
	Y      float32
	Width  float32
	Height float32
}

// Contains determines whether a 2D point is inside the Rect.
func (r *Rect) Contains(x, y float32) bool {

	if x < r.X || x > r.X+r.Width {
		return false
	}
	if y < r.Y || y > r.Y+r.Height {
		return false
	}
	return true
}
