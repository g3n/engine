// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math32

import (
	"strings"
)

// Color4 describes an RGBA color
type Color4 struct {
	R float32
	G float32
	B float32
	A float32
}

// NewColor4 creates and returns a pointer to a new Color4
// with the specified standard web color name (case insensitive)
// and an optional alpha channel value.
// Returns nil if the specified color name not found
func NewColor4(name string, alpha ...float32) *Color4 {

	c, ok := mapColorNames[strings.ToLower(name)]
	if !ok {
		return nil
	}
	a := float32(1)
	if len(alpha) > 0 {
		a = alpha[0]
	}
	return &Color4{c.R, c.G, c.B, a}
}

// Color4Name returns a Color4 with the specified standard web color name
// and an optional alpha channel value.
func Color4Name(name string, alpha ...float32) Color4 {

	c := mapColorNames[strings.ToLower(name)]
	a := float32(1)
	if len(alpha) > 0 {
		a = alpha[0]
	}
	return Color4{c.R, c.G, c.B, a}
}

// Set sets this color individual R,G,B,A components
// Returns pointer to this updated color
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
// Returns pointer to this updated color
func (c *Color4) SetHex(value uint) *Color4 {

	c.R = float32((value >> 16 & 255)) / 255
	c.G = float32((value >> 8 & 255)) / 255
	c.B = float32((value & 255)) / 255
	return c
}

// SetName sets the color RGB components from the
// specified standard web color name
// Returns pointer to this updated color
func (c *Color4) SetName(name string) *Color4 {

	*c = Color4Name(name, 1)
	return c
}

// Add adds to each RGBA component of this color the correspondent component of other color
// Returns pointer to this updated color
func (c *Color4) Add(other *Color4) *Color4 {

	c.R += other.R
	c.G += other.G
	c.B += other.B
	c.A += other.A
	return c
}

// MultiplyScalar multiplies each RGBA component of this color by the specified scalar.
// Returns pointer to this updated color
func (c *Color4) MultiplyScalar(v float32) *Color4 {

	c.R *= v
	c.G *= v
	c.B *= v
	c.A *= v
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
