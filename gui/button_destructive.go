// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

//import (
//	"github.com/g3n/engine/window"
//)

/***************************************

 ButtonDestructive Panel
 +-------------------------------+
 |  Image/Icon      Label        |
 |  +----------+   +----------+  |
 |  |          |   |          |  |
 |  |          |   |          |  |
 |  +----------+   +----------+  |
 +-------------------------------+

****************************************/

// ButtonDestructive represents a destructive button GUI element
type ButtonDestructive struct {
	Button                           // Embedded Button
	styles  *ButtonDestructiveStyles // pointer to current destructive button styles
	pressed int                      // the number of times this button has been pressed
}

// ButtonDestructiveStyle contains the styling of a ButtonDestructive
type ButtonDestructiveStyle ButtonStyle

// ButtonDestructiveStyles contains one ButtonDestructiveStyle for each possible button state
type ButtonDestructiveStyles struct {
	Normal       ButtonDestructiveStyle
	Over         ButtonDestructiveStyle
	PressedOnce  ButtonDestructiveStyle
	PressedTwice ButtonDestructiveStyle
	Disabled     ButtonDestructiveStyle
}

// NewButtonDestructive creates and returns a pointer to a new button widget
// with the specified text for the button label.
func NewButtonDestructive(text string) *ButtonDestructive {

	b := new(ButtonDestructive)
	b.Button.styles = &StyleDefault().Button // unnecessary, I think
	b.styles = &StyleDefault().ButtonDestructive

	// Initializes the button panel
	b.Panel = NewPanel(0, 0)

	// Subscribe to panel events
	b.Subscribe(OnKeyDown, b.onKey)
	b.Subscribe(OnKeyUp, b.onKey)
	b.Subscribe(OnMouseUp, b.onMouse)
	b.Subscribe(OnMouseDown, b.onMouse)
	b.Subscribe(OnCursor, b.onCursor)
	b.Subscribe(OnCursorEnter, b.onCursor)
	b.Subscribe(OnCursorLeave, b.onCursor)
	b.Subscribe(OnEnable, func(name string, ev interface{}) { b.update() })
	b.Subscribe(OnResize, func(name string, ev interface{}) { b.recalc() })

	// Creates label
	b.Label = NewLabel(text)
	b.Label.Subscribe(OnResize, func(name string, ev interface{}) { b.recalc() })
	b.Panel.Add(b.Label)

	b.recalc() // recalc first then update!
	b.update()
	return b
}

// SetStyles set the destructive button styles overriding the default style
func (b *ButtonDestructive) SetStyles(bs *ButtonDestructiveStyles) {

	b.styles = bs
	b.update()
}

// onCursor process subscribed cursor events
func (b *ButtonDestructive) onCursor(evname string, ev interface{}) {

	switch evname {
	case OnCursorEnter:
		b.mouseOver = true
		b.update()
	case OnCursorLeave:
		b.pressed = 0
		b.mouseOver = false
		b.update()
	}
	b.root.StopPropagation(StopAll)
}

// onMouseEvent process subscribed mouse events
func (b *ButtonDestructive) onMouse(evname string, ev interface{}) {

	switch evname {
	case OnMouseDown:
		b.root.SetKeyFocus(b)
		b.pressed++
		b.update()
		b.Dispatch(OnClick, nil)
	case OnMouseUp:
		b.update()
	default:
		return
	}
	b.root.StopPropagation(StopAll)
}

// onKey processes subscribed key events
func (b *ButtonDestructive) onKey(evname string, ev interface{}) {

	kev := ev.(*window.KeyEvent)
	if evname == OnKeyDown && kev.Keycode == window.KeyEnter {
		b.pressed++
		b.update()
		b.Dispatch(OnClick, nil)
		b.root.StopPropagation(Stop3D)
		return
	}
	if evname == OnKeyUp && kev.Keycode == window.KeyEnter {
		b.update()
		b.root.StopPropagation(Stop3D)
		return
	}
	return
}

// update updates the button visual state
func (b *ButtonDestructive) update() {

	if !b.Enabled() {
		b.applyStyle(&b.styles.Disabled)
		return
	}
	switch {
	case b.pressed >= 2:
		b.applyStyle(&b.styles.PressedTwice)
		return
	case b.pressed == 1:
		b.applyStyle(&b.styles.PressedOnce)
		return
	}
	if b.mouseOver {
		b.applyStyle(&b.styles.Over)
		return
	}
	b.applyStyle(&b.styles.Normal)
}

// applyStyle applies the specified button style
func (b *ButtonDestructive) applyStyle(bs *ButtonDestructiveStyle) {

	b.Panel.ApplyStyle(&bs.PanelStyle)
	if b.icon != nil {
		b.icon.SetColor4(&bs.FgColor)
	}
	b.Label.SetColor4(&bs.FgColor)
}

// recalc recalculates all dimensions and position from inside out
func (b *ButtonDestructive) recalc() {

	// Current width and height of button content area
	width := b.Panel.ContentWidth()
	height := b.Panel.ContentHeight()

	// Image or icon width
	imgWidth := float32(0)
	spacing := float32(4)
	if b.image != nil {
		imgWidth = b.image.Width()
	} else if b.icon != nil {
		imgWidth = b.icon.Width()
	}
	if imgWidth == 0 {
		spacing = 0
	}

	// If the label is empty and an icon of image was defined ignore the label widthh
	// to centralize the icon/image in the button
	labelWidth := spacing + b.Label.Width()
	if b.Label.Text() == "" && imgWidth > 0 {
		labelWidth = 0
	}

	// Sets new content width and height if necessary
	minWidth := imgWidth + labelWidth
	minHeight := b.Label.Height()
	resize := false
	if width < minWidth {
		width = minWidth
		resize = true
	}
	if height < minHeight {
		height = minHeight
		resize = true
	}
	if resize {
		b.SetContentSize(width, height)
	}

	// Centralize horizontally
	px := (width - minWidth) / 2

	// Set label position
	ly := (height - b.Label.Height()) / 2
	b.Label.SetPosition(px+imgWidth+spacing, ly)

	// Image/icon position
	if b.image != nil {
		iy := (height - b.image.height) / 2
		b.image.SetPosition(px, iy)
	} else if b.icon != nil {
		b.icon.SetPosition(px, ly)
	}
}
