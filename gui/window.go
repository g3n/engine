// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/window"
	"github.com/g3n/engine/gui/assets/icon"
)

/*********************************************

 Window panel
 +-------------------------------------+---+
 |              Title panel            | X |
 +-------------------------------------+---+
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
	Panel       // Embedded Panel
	styles      *WindowStyles
	title       *WindowTitle // internal optional title panel
	client      Panel        // internal client panel
	resizable   bool         // Specifies whether the window is resizable
	drag        bool    // Whether the mouse buttons is pressed (i.e. when dragging)
	dragPadding float32 // Extra width used to resize (in addition to border sizes)

	// To keep track of which window borders the cursor is over
	overTop    bool
	overRight  bool
	overBottom bool
	overLeft   bool

	// Minimum and maximum sizes
	minSize math32.Vector2
	maxSize math32.Vector2
}

// WindowStyle contains the styling of a Window
type WindowStyle struct {
	PanelStyle
	TitleStyle WindowTitleStyle
}

// WindowStyles contains a WindowStyle for each valid GUI state
type WindowStyles struct {
	Normal   WindowStyle
	Over     WindowStyle
	Focus    WindowStyle
	Disabled WindowStyle
}

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

	w.dragPadding = 5

	w.recalc()
	w.update()
	return w
}

// SetResizable sets whether the window is resizable.
func (w *Window) SetResizable(state bool) {

	w.resizable = state
}

// SetCloseButton sets whether the window has a close button on the top right.
func (w *Window) SetCloseButton(state bool) {

	w.title.setCloseButton(state)
}

// SetTitle sets the title of the window.
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

// SetLayout sets the layout of the client panel.
func (w *Window) SetLayout(layout ILayout) {

	w.client.SetLayout(layout)
}

// onMouse process subscribed mouse events over the window
func (w *Window) onMouse(evname string, ev interface{}) {

	switch evname {
	case OnMouseDown:
		// Move the window above everything contained in its parent
		parent := w.Parent().(IPanel).GetPanel()
		parent.SetTopChild(w)
		// If the click happened inside the draggable area, then set drag to true
		if w.overTop || w.overRight || w.overBottom || w.overLeft {
			w.drag = true
			w.root.SetMouseFocus(w)
		}
	case OnMouseUp:
		w.drag = false
		w.root.SetMouseFocus(nil)
	default:
		return
	}
	w.root.StopPropagation(StopAll)
}

// onCursor process subscribed cursor events over the window
func (w *Window) onCursor(evname string, ev interface{}) {

	// If the window is not resizable we are not interested in cursor movements
	if !w.resizable {
		return
	}
	if evname == OnCursor {
		cev := ev.(*window.CursorEvent)
		// If already dragging - update window size and position depending
		// on the cursor position and the borders being dragged
		if w.drag {
			if w.overTop {
				delta := cev.Ypos - w.pospix.Y
				newHeight := w.Height() - delta
				minHeight := w.title.height
				if newHeight >= minHeight {
					w.SetPositionY(w.Position().Y + delta)
					w.SetHeight(math32.Max(newHeight, minHeight))
				} else {
					w.SetPositionY(w.Position().Y + w.Height() - minHeight)
					w.SetHeight(w.title.height)
				}
			}
			if w.overRight {
				delta := cev.Xpos - (w.pospix.X + w.width)
				newWidth := w.Width() + delta
				w.SetWidth(math32.Max(newWidth, w.title.label.Width() + w.title.closeButton.Width()))
			}
			if w.overBottom {
				delta := cev.Ypos - (w.pospix.Y + w.height)
				newHeight := w.Height() + delta
				w.SetHeight(math32.Max(newHeight, w.title.height))
			}
			if w.overLeft {
				delta := cev.Xpos - w.pospix.X
				newWidth := w.Width() - delta
				minWidth := w.title.label.Width() + w.title.closeButton.Width()
				if newWidth >= minWidth {
					w.SetPositionX(w.Position().X + delta)
					w.SetWidth(math32.Max(newWidth, minWidth))
				} else {
					w.SetPositionX(w.Position().X + w.Width() - minWidth)
					w.SetWidth(minWidth)
				}
			}
		} else {
			// Obtain cursor position relative to window
			cx := cev.Xpos - w.pospix.X
			cy := cev.Ypos - w.pospix.Y
			// Check if cursor is on the top of the window (border + drag margin)
			if cy <= w.borderSizes.Top {
				w.overTop = true
				w.root.SetCursorVResize()
			} else {
				w.overTop = false
			}
			// Check if cursor is on the bottom of the window (border + drag margin)
			if cy >= w.height-w.borderSizes.Bottom - w.dragPadding {
				w.overBottom = true
			} else {
				w.overBottom = false
			}
			// Check if cursor is on the left of the window (border + drag margin)
			if cx <= w.borderSizes.Left + w.dragPadding {
				w.overLeft = true
				w.root.SetCursorHResize()
			} else {
				w.overLeft = false
			}
			// Check if cursor is on the right of the window (border + drag margin)
			if cx >= w.width-w.borderSizes.Right - w.dragPadding {
				w.overRight = true
				w.root.SetCursorHResize()
			} else {
				w.overRight = false
			}
			// Update cursor image based on cursor position
			if (w.overTop || w.overBottom) && !w.overRight && !w.overLeft {
				w.root.SetCursorVResize()
			} else if (w.overRight || w.overLeft) && !w.overTop && !w.overBottom {
				w.root.SetCursorHResize()
			} else if (w.overRight && w.overTop) || (w.overBottom && w.overLeft) {
				w.root.SetCursorDiagResize1()
			} else if (w.overRight && w.overBottom) || (w.overTop && w.overLeft) {
				w.root.SetCursorDiagResize2()
			}
			// If cursor is not near the border of the window then reset the cursor
			if !w.overTop && !w.overRight && !w.overBottom && !w.overLeft {
				w.root.SetCursorNormal()
			}
		}
	} else if evname == OnCursorLeave {
		w.root.SetCursorNormal()
		w.drag = false
	}
	w.root.StopPropagation(StopAll)
}

// update updates the window's visual state.
func (w *Window) update() {

	if !w.Enabled() {
		w.applyStyle(&w.styles.Disabled)
		return
	}
	w.applyStyle(&w.styles.Normal)
}

// applyStyle applies a window style to the window.
func (w *Window) applyStyle(s *WindowStyle) {

	w.SetBordersColor4(&s.BorderColor)
	w.SetBordersFrom(&s.Border)
	w.SetPaddingsFrom(&s.Padding)
	w.client.SetMarginsFrom(&s.Margin)
	w.client.SetColor4(&s.BgColor)
	if w.title != nil {
		w.title.applyStyle(&s.TitleStyle)
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
	Panel              // Embedded panel
	win                *Window // Window to which this title belongs
	label              Label   // Label for the title
	pressed            bool    // Whether the left mouse button is pressed
	closeButton        *Button // The close button on the top right corner
	closeButtonVisible bool    // Whether the close button is present

	// Last mouse coordinates
	mouseX             float32
	mouseY             float32
}

// WindowTitleStyle contains the styling for a window title.
type WindowTitleStyle struct {
	PanelStyle
	FgColor     math32.Color4
}

// newWindowTitle creates and returns a pointer to a window title panel.
func newWindowTitle(win *Window, text string) *WindowTitle {

	wt := new(WindowTitle)
	wt.win = win

	wt.Panel.Initialize(0, 0)
	wt.label.initialize(text, StyleDefault().Font)
	wt.Panel.Add(&wt.label)

	wt.closeButton = NewButton("")
	wt.closeButton.SetIcon(icon.Close)
	wt.closeButton.Subscribe(OnCursorEnter, func(s string, i interface{}) {
		wt.win.root.SetCursorNormal()
	})
	wt.closeButton.Subscribe(OnClick, func(s string, i interface{}) {
		wt.win.Parent().GetNode().Remove(wt.win)
		wt.win.Dispose()
		wt.win.Dispatch("gui.OnWindowClose", nil)
	})
	wt.Panel.Add(wt.closeButton)
	wt.closeButtonVisible = true

	wt.Subscribe(OnMouseDown, wt.onMouse)
	wt.Subscribe(OnMouseUp, wt.onMouse)
	wt.Subscribe(OnCursor, wt.onCursor)
	wt.Subscribe(OnCursorEnter, wt.onCursor)
	wt.Subscribe(OnCursorLeave, wt.onCursor)

	wt.recalc()
	return wt
}

// setCloseButton sets whether the close button is present on the top right corner.
func (wt *WindowTitle) setCloseButton(state bool) {

	if state {
		wt.closeButtonVisible = true
		wt.Panel.Add(wt.closeButton)
	} else {
		wt.closeButtonVisible = false
		wt.Panel.Remove(wt.closeButton)
	}
}

// onMouse process subscribed mouse button events over the window title.
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

// onCursor process subscribed cursor events over the window title.
func (wt *WindowTitle) onCursor(evname string, ev interface{}) {

	if evname == OnCursorLeave {
		wt.win.root.SetCursorNormal()
		wt.pressed = false
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

// applyStyle applies the specified WindowTitleStyle.
func (wt *WindowTitle) applyStyle(s *WindowTitleStyle) {

	wt.Panel.ApplyStyle(&s.PanelStyle)
	wt.label.SetColor4(&s.FgColor)
}

// recalc recalculates the height and position of the label in the title bar.
func (wt *WindowTitle) recalc() {

	xpos := (wt.width - wt.label.width) / 2
	wt.label.SetPositionX(xpos)
	wt.SetContentHeight(wt.closeButton.Height())

	if wt.closeButtonVisible {
		wt.closeButton.SetPositionX(wt.width - wt.closeButton.width)
	}
}
