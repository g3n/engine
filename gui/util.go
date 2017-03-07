// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

type BorderSizes struct {
	Top    float32
	Right  float32
	Bottom float32
	Left   float32
}

func (bs *BorderSizes) Set(top, right, bottom, left float32) {

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

type Rect struct {
	X      float32
	Y      float32
	Width  float32
	Height float32
}

func (r *Rect) Contains(x, y float32) bool {

	if x < r.X || x > r.X+r.Width {
		return false
	}
	if y < r.Y || y > r.Y+r.Height {
		return false
	}
	return true
}
