// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/window"
)

/*********************************************

 Window panel
 +-----------------------------------------+
 | Title panel                             |
 +-----------------------------------------+
 |  Content panel                          |
 |  +-----------------------------------+  |
 |  |                                   |  |
 |  |                                   |  |
 |  |                                   |  |
 |  |                                   |  |
 |  |                                   |  |
 |  |                                   |  |
 |  +-----------------------------------+  |
 |                                         |
 +-----------------------------------------+

*********************************************/

// Window represents a window GUI element
type Window struct {
	Panel      // Embedded Panel
	styles     *WindowStyles
	title      *WindowTitle // internal optional title panel
	client     Panel        // internal client panel
	resizable  ResizeBorders
	overBorder string
	drag       bool
	mouseX     float32
	mouseY     float32
}

// WindowStyle contains the styling of a Window
type WindowStyle struct {
	Border           RectBounds
	Paddings         RectBounds
	BorderColor      math32.Color4
	TitleBorders     RectBounds
	TitleBorderColor math32.Color4
	TitleBgColor     math32.Color4
	TitleFgColor     math32.Color4
}

// WindowStyles contains a WindowStyle for each valid GUI state
type WindowStyles struct {
	Normal   WindowStyle
	Over     WindowStyle
	Focus    WindowStyle
	Disabled WindowStyle
}

// ResizeBorders specifies which window borders can be resized
type ResizeBorders int

// Resizing can be allowed or disallowed on each window edge
const (
	ResizeTop = ResizeBorders(1 << (iota + 1))
	ResizeRight
	ResizeBottom
	ResizeLeft
	ResizeAll = ResizeTop | ResizeRight | ResizeBottom | ResizeLeft
)

// NewWindow creates and returns a pointer to a new window with the
// specified dimensions
func NewWindow(width, height float32) *Window {

	w := new(Window)
	w.styles = &StyleDefault().Window

	w.Panel.Initialize(width, height)
	w.Panel.Subscribe(OnMouseDown, w.onMouse)
	w.Panel.Subscribe(OnMouseUp, w.onMouse)
	w.Panel.Subscribe(OnCursor, w.onCursor)
	w.Panel.Subscribe(OnCursorEnter, w.onCursor)
	w.Panel.Subscribe(OnCursorLeave, w.onCursor)
	w.Panel.Subscribe(OnResize, func(evname string, ev interface{}) { w.recalc() })

	w.client.Initialize(0, 0)
	w.Panel.Add(&w.client)

	w.recalc()
	w.update()
	return w
}

// SetResizable set the borders which are resizable
func (w *Window) SetResizable(res ResizeBorders) {

	w.resizable = res
}

// SetTitle sets the title of this window
func (w *Window) SetTitle(text string) {

	if w.title == nil {
		w.title = newWindowTitle(w, text)
		w.Panel.Add(w.title)
	} else {
		w.title.label.SetText(text)
	}
	w.update()
	w.recalc()
}

// Add adds a child panel to the client area of this window
func (w *Window) Add(ichild IPanel) *Window {

	w.client.Add(ichild)
	return w
}

// SetLayout set the layout of this window content area
func (w *Window) SetLayout(layout ILayout) {

	w.client.SetLayout(layout)
}

// onMouse process subscribed mouse events over the window
func (w *Window) onMouse(evname string, ev interface{}) {

	mev := ev.(*window.MouseEvent)
	switch evname {
	case OnMouseDown:
		par := w.Parent().(IPanel).GetPanel()
		par.SetTopChild(w)
		if w.overBorder != "" {
			w.drag = true
			w.mouseX = mev.Xpos
			w.mouseY = mev.Ypos
			w.root.SetMouseFocus(w)
		}
	case OnMouseUp:
		w.drag = false
		w.root.SetCursorNormal()
		w.root.SetMouseFocus(nil)
	default:
		return
	}
	w.root.StopPropagation(StopAll)
}

// onCursor process subscribed cursor events over the window
func (w *Window) onCursor(evname string, ev interface{}) {

	if evname == OnCursor {
		cev := ev.(*window.CursorEvent)
		if !w.drag {
			cx := cev.Xpos - w.pospix.X
			cy := cev.Ypos - w.pospix.Y
			if cy <= w.borderSizes.Top {
				if w.resizable&ResizeTop != 0 {
					w.overBorder = "top"
					w.root.SetCursorVResize()
				}
			} else if cy >= w.height-w.borderSizes.Bottom {
				if w.resizable&ResizeBottom != 0 {
					w.overBorder = "bottom"
					w.root.SetCursorVResize()
				}
			} else if cx <= w.borderSizes.Left {
				if w.resizable&ResizeLeft != 0 {
					w.overBorder = "left"
					w.root.SetCursorHResize()
				}
			} else if cx >= w.width-w.borderSizes.Right {
				if w.resizable&ResizeRight != 0 {
					w.overBorder = "right"
					w.root.SetCursorHResize()
				}
			} else {
				if w.overBorder != "" {
					w.root.SetCursorNormal()
					w.overBorder = ""
				}
			}
		} else {
			switch w.overBorder {
			case "top":
				delta := cev.Ypos - w.mouseY
				w.mouseY = cev.Ypos
				newHeight := w.Height() - delta
				if newHeight < w.MinHeight() {
					return
				}
				w.SetPositionY(w.Position().Y + delta)
				w.SetHeight(newHeight)
			case "right":
				delta := cev.Xpos - w.mouseX
				w.mouseX = cev.Xpos
				newWidth := w.Width() + delta
				w.SetWidth(newWidth)
			case "bottom":
				delta := cev.Ypos - w.mouseY
				w.mouseY = cev.Ypos
				newHeight := w.Height() + delta
				w.SetHeight(newHeight)
			case "left":
				delta := cev.Xpos - w.mouseX
				w.mouseX = cev.Xpos
				newWidth := w.Width() - delta
				if newWidth < w.MinWidth() {
					return
				}
				w.SetPositionX(w.Position().X + delta)
				w.SetWidth(newWidth)
			}
		}
	} else if evname == OnCursorLeave {
		if !w.drag {
			w.root.SetCursorNormal()
		}
	}
	w.root.StopPropagation(StopAll)
}

// update updates the button visual state
func (w *Window) update() {

	if !w.Enabled() {
		w.applyStyle(&w.styles.Disabled)
		return
	}
	w.applyStyle(&w.styles.Normal)
}

func (w *Window) applyStyle(s *WindowStyle) {

	w.SetBordersColor4(&s.BorderColor)
	w.SetBordersFrom(&s.Border)
	w.SetPaddingsFrom(&s.Paddings)
	if w.title != nil {
		w.title.applyStyle(s)
	}
}

// recalc recalculates the sizes and positions of the internal panels
// from the outside to the inside.
func (w *Window) recalc() {

	// Window title
	height := w.content.Height
	width := w.content.Width
	cx := float32(0)
	cy := float32(0)
	if w.title != nil {
		w.title.SetWidth(w.content.Width)
		w.title.recalc()
		height -= w.title.height
		cy = w.title.height
	}

	// Content area
	w.client.SetPosition(cx, cy)
	w.client.SetSize(width, height)
}

// WindowTitle represents the title bar of a Window
type WindowTitle struct {
	Panel   // Embedded panel
	win     *Window
	label   Label
	pressed bool
	drag    bool
	mouseX  float32
	mouseY  float32
}

// newWindowTitle creates and returns a pointer to a window title panel
func newWindowTitle(win *Window, text string) *WindowTitle {

	wt := new(WindowTitle)
	wt.win = win

	wt.Panel.Initialize(0, 0)
	wt.label.initialize(text, StyleDefault().Font)
	wt.Panel.Add(&wt.label)

	wt.Subscribe(OnMouseDown, wt.onMouse)
	wt.Subscribe(OnMouseUp, wt.onMouse)
	wt.Subscribe(OnCursor, wt.onCursor)
	wt.Subscribe(OnCursorEnter, wt.onCursor)
	wt.Subscribe(OnCursorLeave, wt.onCursor)

	wt.recalc()
	return wt
}

// onMouse process subscribed mouse button events over the window title
func (wt *WindowTitle) onMouse(evname string, ev interface{}) {

	mev := ev.(*window.MouseEvent)
	switch evname {
	case OnMouseDown:
		wt.pressed = true
		wt.mouseX = mev.Xpos
		wt.mouseY = mev.Ypos
		wt.win.root.SetMouseFocus(wt)
	case OnMouseUp:
		wt.pressed = false
		wt.win.root.SetMouseFocus(nil)
	default:
		return
	}
	wt.win.root.StopPropagation(Stop3D)
}

// onCursor process subscribed cursor events over the window title
func (wt *WindowTitle) onCursor(evname string, ev interface{}) {

	if evname == OnCursorEnter {
		wt.win.root.SetCursorHand()
	} else if evname == OnCursorLeave {
		wt.win.root.SetCursorNormal()
	} else if evname == OnCursor {
		if !wt.pressed {
			wt.win.root.StopPropagation(Stop3D)
			return
		}
		cev := ev.(*window.CursorEvent)
		dy := wt.mouseY - cev.Ypos
		dx := wt.mouseX - cev.Xpos
		wt.mouseX = cev.Xpos
		wt.mouseY = cev.Ypos
		posX := wt.win.Position().X - dx
		posY := wt.win.Position().Y - dy
		wt.win.SetPosition(posX, posY)
	}
	wt.win.root.StopPropagation(Stop3D)
}

// applyStyles sets the specified window title style
func (wt *WindowTitle) applyStyle(s *WindowStyle) {

	wt.SetBordersFrom(&s.TitleBorders)
	wt.SetBordersColor4(&s.TitleBorderColor)
	wt.SetColor4(&s.TitleBgColor)
	wt.label.SetColor4(&s.TitleFgColor)
}

// recalc recalculates the height and position of the label in the title bar.
func (wt *WindowTitle) recalc() {

	xpos := (wt.width - wt.label.width) / 2
	wt.label.SetPositionX(xpos)
	wt.SetContentHeight(wt.label.Height())
}
