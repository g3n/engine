// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/window"
)

type Splitter struct {
	Panel                     // Embedded panel
	P0        Panel           // Left/Top panel
	P1        Panel           // Right/Bottom panel
	styles    *SplitterStyles // pointer to current styles
	spacer    Panel           // spacer panel
	horiz     bool            // horizontal or vertical splitter
	pos       float32         // central position of the spacer panel bar in pixels
	posLast   float32         // last position of the mouse cursor when dragging
	pressed   bool            // mouse button is pressed and dragging
	mouseOver bool            // mouse is over the spacer panel
}

type SplitterStyle struct {
	SpacerBorderColor math32.Color4
	SpacerColor       math32.Color
	SpacerSize        float32
}

type SplitterStyles struct {
	Normal SplitterStyle
	Over   SplitterStyle
	Drag   SplitterStyle
}

// NewHSplitter creates and returns a pointer to a new horizontal splitter
// widget with the specified initial dimensions
func NewHSplitter(width, height float32) *Splitter {

	return newSplitter(true, width, height)
}

// NewVSplitter creates and returns a pointer to a new vertical splitter
// widget with the specified initial dimensions
func NewVSplitter(width, height float32) *Splitter {

	return newSplitter(false, width, height)
}

// newSpliter creates and returns a pointer of a new splitter with
// the specified orientation and initial dimensions.
func newSplitter(horiz bool, width, height float32) *Splitter {

	s := new(Splitter)
	s.horiz = horiz
	s.styles = &StyleDefault.Splitter
	s.Panel.Initialize(width, height)

	// Initialize left/top panel
	s.P0.Initialize(0, 0)
	s.Panel.Add(&s.P0)

	// Initialize right/bottom panel
	s.P1.Initialize(0, 0)
	s.Panel.Add(&s.P1)

	// Initialize spacer panel
	s.spacer.Initialize(0, 0)
	s.Panel.Add(&s.spacer)

	if horiz {
		s.spacer.SetBorders(0, 1, 0, 1)
		s.spacer.SetWidth(6)
		s.pos = width / 2
	} else {
		s.spacer.SetBorders(1, 0, 1, 0)
		s.spacer.SetHeight(6)
		s.pos = height / 2
	}

	s.spacer.Subscribe(OnMouseDown, s.onMouse)
	s.spacer.Subscribe(OnMouseUp, s.onMouse)
	s.spacer.Subscribe(OnCursor, s.onCursor)
	s.spacer.Subscribe(OnCursorEnter, s.onCursor)
	s.spacer.Subscribe(OnCursorLeave, s.onCursor)
	s.update()
	s.recalc()
	return s
}

// SetSplit sets the position of the splitter bar.
// It accepts a value from 0.0 to 1.0
func (s *Splitter) SetSplit(pos float32) {

	if s.horiz {
		s.setSplit(pos * s.Width())
	} else {
		s.setSplit(pos * s.Height())
	}
	s.recalc()
}

// Split returns the current position of the splitter bar.
// It returns a value from 0.0 to 1.0
func (s *Splitter) Split() float32 {

	var pos float32
	if s.horiz {
		pos = s.pos / s.Width()
	} else {
		pos = s.pos / s.Height()
	}
	return pos
}

// onMouse receives subscribed mouse events over the spacer panel
func (s *Splitter) onMouse(evname string, ev interface{}) {

	mev := ev.(*window.MouseEvent)
	switch evname {
	case OnMouseDown:
		s.pressed = true
		if s.horiz {
			s.posLast = mev.Xpos
		} else {
			s.posLast = mev.Ypos
		}
		s.root.SetMouseFocus(&s.spacer)
	case OnMouseUp:
		s.pressed = false
		s.root.SetCursorNormal()
		s.root.SetMouseFocus(nil)
	default:
	}
	s.root.StopPropagation(Stop3D)
}

// onCursor receives subscribed cursor events over the spacer panel
func (s *Splitter) onCursor(evname string, ev interface{}) {

	if evname == OnCursorEnter {
		if s.horiz {
			s.root.SetCursorHResize()
		} else {
			s.root.SetCursorVResize()
		}
		s.mouseOver = true
		s.update()
	} else if evname == OnCursorLeave {
		s.root.SetCursorNormal()
		s.mouseOver = false
		s.update()
	} else if evname == OnCursor {
		if !s.pressed {
			return
		}
		cev := ev.(*window.CursorEvent)
		var delta float32
		if s.horiz {
			delta = cev.Xpos - s.posLast
			s.posLast = cev.Xpos
		} else {
			delta = cev.Ypos - s.posLast
			s.posLast = cev.Ypos
		}
		s.setSplit(s.pos + delta)
		s.recalc()
	}
	s.root.StopPropagation(Stop3D)
}

// setSplit sets the validated and clamped split position from the received value.
func (s *Splitter) setSplit(pos float32) {

	var max float32
	var halfspace float32
	if s.horiz {
		halfspace = s.spacer.Width() / 2
		max = s.Width()
	} else {
		halfspace = s.spacer.Height() / 2
		max = s.Height()
	}

	if pos > max-halfspace {
		s.pos = max - halfspace
	} else if pos <= halfspace {
		s.pos = halfspace
	} else {
		s.pos = pos
	}
}

// update updates the splitter visual state
func (s *Splitter) update() {

	if s.pressed {
		s.applyStyle(&s.styles.Drag)
		return
	}
	if s.mouseOver {
		s.applyStyle(&s.styles.Over)
		return
	}
	s.applyStyle(&s.styles.Normal)
}

// applyStyle applies the specified splitter style
func (s *Splitter) applyStyle(ss *SplitterStyle) {

	s.spacer.SetBordersColor4(&ss.SpacerBorderColor)
	s.spacer.SetColor(&ss.SpacerColor)
	if s.horiz {
		s.spacer.SetWidth(ss.SpacerSize)
	} else {
		s.spacer.SetHeight(ss.SpacerSize)
	}
}

// recalc relcalculates the position and sizes of the internal panels
func (s *Splitter) recalc() {

	width := s.Width()
	height := s.Height()
	if s.horiz {
		halfspace := s.spacer.Width() / 2
		// First panel
		s.P0.SetPosition(0, 0)
		s.P0.SetSize(s.pos-halfspace, height)
		// Spacer panel
		s.spacer.SetPosition(s.pos-halfspace, 0)
		s.spacer.SetHeight(height)
		// Second panel
		s.P1.SetPosition(s.pos+halfspace, 0)
		s.P1.SetSize(width-s.pos-halfspace, height)
	} else {
		halfspace := s.spacer.Height() / 2
		// First panel
		s.P0.SetPosition(0, 0)
		s.P0.SetSize(width, s.pos-halfspace)
		// Spacer panel
		s.spacer.SetPosition(0, s.pos-halfspace)
		s.spacer.SetWidth(width)
		// Second panel
		s.P1.SetPosition(0, s.pos+halfspace)
		s.P1.SetSize(width, height-s.pos-halfspace)
	}
}
