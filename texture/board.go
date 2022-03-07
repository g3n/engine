// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package texture

import (
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/math32"
)

// NewBoard creates and returns a pointer to a new checker board 2D texture.
// A checker board texture contains up to 4 different colors arranged in
// the following order:
//  +------+------+
//  |      |      |
//  |  c3  |  c4  |
//  |      |      |
//  +------+------+
//  |      |      |
//  |  c1  |  c2  | height (pixels)
//  |      |      |
//  +------+------+
//    width
//  (pixels)
//
func NewBoard(width, height int, c1, c2, c3, c4 *math32.Color, alpha float32) *Texture2D {

	// Generates texture data
	data := make([]float32, width*height*4*4)
	colorData := func(sx, sy int, c *math32.Color) {
		for y := sy; y < sy+height; y++ {
			for x := sx; x < sx+width; x++ {
				pos := (x + y*2*width) * 4
				data[pos] = c.R
				data[pos+1] = c.G
				data[pos+2] = c.B
				data[pos+3] = alpha
			}
		}
	}
	colorData(0, 0, c1)
	colorData(width, 0, c2)
	colorData(0, height, c3)
	colorData(width, height, c4)

	// Creates, initializes and returns board texture object
	return NewTexture2DFromData(width*2, height*2, gls.RGBA, gls.FLOAT, gls.RGBA8, data)
}
