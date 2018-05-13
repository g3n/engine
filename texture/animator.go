// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package texture

import (
	"github.com/g3n/engine/gls"
	"time"
)

// Animator can generate a texture animation based on a texture sheet
type Animator struct {
	tex       *Texture2D    // pointer to texture being displayed
	dispTime  time.Duration // disply duration of each tile (default = 1.0/30.0)
	maxCycles int           // maximum number of cycles (default = 0 - continuous)
	cycles    int           // current number of complete cycles
	columns   int           // number of columns
	rows      int           // number of rows
	row       int           // current row
	col       int           // current column
	tileTime  time.Time     // time when tile started to be displayed
}

// NewAnimator creates and returns a texture sheet animator for the specified texture.
func NewAnimator(tex *Texture2D, htiles, vtiles int) *Animator {

	a := new(Animator)
	a.tex = tex
	a.columns = htiles
	a.rows = vtiles
	a.dispTime = time.Millisecond * 16
	a.maxCycles = 0

	// Sets texture properties
	tex.SetWrapS(gls.REPEAT)
	tex.SetWrapT(gls.REPEAT)
	tex.SetRepeat(1/float32(a.columns), 1/float32(a.rows))

	// Initial state
	a.Restart()
	return a
}

// SetDispTime sets the display time of each tile in milliseconds.
// The default value is: 1.0/30.0 = 16.6.ms
func (a *Animator) SetDispTime(dtime time.Duration) {

	a.dispTime = dtime
}

// SetMaxCycles sets the number of complete cycles to display.
// The default value is: 0 (display continuously)
func (a *Animator) SetMaxCycles(maxCycles int) {

	a.maxCycles = maxCycles
}

// Cycles returns the number of complete cycles displayed
func (a *Animator) Cycles() int {

	return a.cycles
}

// Restart restart the animator
func (a *Animator) Restart() {

	// Time of the currently displayed image
	a.tileTime = time.Now()
	// Position of current tile to display
	a.row = 0
	a.col = 0
	// Number of cycles displayed
	a.cycles = 0
}

// Update prepares the next tile to be rendered.
// Must be called with the current time
func (a *Animator) Update(now time.Time) {

	// Checks maximum number of cycles
	if a.maxCycles > 0 && a.cycles >= a.maxCycles {
		return
	}

	// If current tile time not reached, do nothing
	if now.Sub(a.tileTime) < a.dispTime {
		return
	}
	a.tileTime = now

	// Sets the position of the next tile to show
	a.tex.SetOffset(float32(a.col)/float32(a.columns), float32(a.row)/float32(a.rows))
	a.col += 1
	if a.col >= a.columns {
		a.col = 0
		a.row += 1
		if a.row >= a.rows {
			a.row = 0
			a.cycles += 1
		}
	}
}
