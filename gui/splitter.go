// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/window"
)

// Splitter is a GUI element that splits two panels and can be adjusted
type Splitter struct {
	Panel                     // Embedded panel
	P0        Panel           // Left/Top panel
	P1        Panel           // Right/Bottom panel
	styles    *SplitterStyles // pointer to current styles
	spacer    Panel           // spacer panel
	horiz     bool            // horizontal or vertical splitter
	pos       float32         // relative position of the center of the spacer panel (0 to 1)
	posLast   float32         // last position in pixels of the mouse cursor when dragging
	pressed   bool            // mouse button is pressed and dragging
	mouseOver bool            // mouse is over the spacer panel
}

// SplitterStyle contains the styling of a Splitter
type SplitterStyle struct {
	SpacerBorderColor math32.Color4
	SpacerColor       math32.Color4
	SpacerSize        float32
}

// SplitterStyles contains a SplitterStyle for each valid GUI state
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
	s.styles = &StyleDefault().Splitter
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
		s.pos = 0.5
	} else {
		s.spacer.SetBorders(1, 0, 1, 0)
		s.pos = 0.5
	}

	s.Subscribe(OnResize, s.onResize)
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

	s.setSplit(pos)
	s.recalc()
}

// Split returns the current position of the splitter bar.
// It returns a value from 0.0 to 1.0
func (s *Splitter) Split() float32 {

	return s.pos
}

// onResize receives subscribed resize events for the whole splitter panel
func (s *Splitter) onResize(evname string, ev interface{}) {

	s.recalc()
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
		pos := s.pos
		if s.horiz {
			delta = cev.Xpos - s.posLast
			s.posLast = cev.Xpos
			pos += delta / s.ContentWidth()
		} else {
			delta = cev.Ypos - s.posLast
			s.posLast = cev.Ypos
			pos += delta / s.ContentHeight()
		}
		s.setSplit(pos)
		s.recalc()
	}
	s.root.StopPropagation(Stop3D)
}

// setSplit sets the validated and clamped split position from the received value.
func (s *Splitter) setSplit(pos float32) {

	if pos < 0 {
		s.pos = 0
	} else if pos > 1 {
		s.pos = 1
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
	s.spacer.SetColor4(&ss.SpacerColor)
	if s.horiz {
		s.spacer.SetWidth(ss.SpacerSize)
	} else {
		s.spacer.SetHeight(ss.SpacerSize)
	}
}

// recalc relcalculates the position and sizes of the internal panels
func (s *Splitter) recalc() {

	width := s.ContentWidth()
	height := s.ContentHeight()
	if s.horiz {
		// Calculate x position for spacer panel
		spx := width*s.pos - s.spacer.Width()/2
		if spx < 0 {
			spx = 0
		} else if spx > width-s.spacer.Width() {
			spx = width - s.spacer.Width()
		}
		// Left panel
		s.P0.SetPosition(0, 0)
		s.P0.SetSize(spx, height)
		// Spacer panel
		s.spacer.SetPosition(spx, 0)
		s.spacer.SetHeight(height)
		// Right panel
		s.P1.SetPosition(spx+s.spacer.Width(), 0)
		s.P1.SetSize(width-spx-s.spacer.Width(), height)
	} else {
		// Calculate y position for spacer panel
		spy := height*s.pos - s.spacer.Height()/2
		if spy < 0 {
			spy = 0
		} else if spy > height-s.spacer.Height() {
			spy = height - s.spacer.Height()
		}
		// Top panel
		s.P0.SetPosition(0, 0)
		s.P0.SetSize(width, spy)
		// Spacer panel
		s.spacer.SetPosition(0, spy)
		s.spacer.SetWidth(width)
		// Bottom panel
		s.P1.SetPosition(0, spy+s.spacer.Height())
		s.P1.SetSize(width, height-spy-s.spacer.Height())
	}
}
