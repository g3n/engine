// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math32

import ()

type Color4 struct {
	R float32
	G float32
	B float32
	A float32
}

var Black4 = Color4{0, 0, 0, 1}
var White4 = Color4{1, 1, 1, 1}
var Red4 = Color4{1, 0, 0, 1}
var Green4 = Color4{0, 1, 0, 1}
var Blue4 = Color4{0, 0, 1, 1}
var Gray4 = Color4{0.5, 0.5, 0.5, 1}

func NewColor4(r, g, b, a float32) *Color4 {

	return &Color4{R: r, G: g, B: b, A: a}
}

// Set sets this color individual R,G,B,A components
func (c *Color4) Set(r, g, b, a float32) *Color4 {

	c.R = r
	c.G = g
	c.B = b
	c.A = b
	return c
}

// SetHex sets the color RGB components from the
// specified integer interpreted as a color hex number
// Alpha component is not modified
func (c *Color4) SetHex(value uint) *Color4 {

	c.R = float32((value >> 16 & 255)) / 255
	c.G = float32((value >> 8 & 255)) / 255
	c.B = float32((value & 255)) / 255
	return c
}

// SetName sets the color RGB components from the
// specified HTML color name
// Alpha component is not modified
func (c *Color4) SetName(name string) *Color4 {

	return c.SetHex(colorKeywords[name])
}

func (c *Color4) MultiplyScalar(v float32) *Color4 {

	c.R *= v
	c.G *= v
	c.B *= v
	return c
}

// FromColor sets this Color4 fields from Color and an alpha
func (c *Color4) FromColor(other *Color, alpha float32) {

	c.R = other.R
	c.G = other.G
	c.B = other.B
	c.A = alpha
}

// ToColor returns a Color with this Color4 RGB components
func (c *Color4) ToColor() Color {

	return Color{c.R, c.G, c.B}
}
