// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/g3n/engine/window"
)

/***************************************

 Button Panel
 +-------------------------------+
 |  Image/Icon      Label        |
 |  +----------+   +----------+  |
 |  |          |   |          |  |
 |  |          |   |          |  |
 |  +----------+   +----------+  |
 +-------------------------------+

****************************************/

// Button represents a button GUI element
type Button struct {
	*Panel                  // Embedded Panel
	Label     *Label        // Label panel
	image     *Image        // pointer to button image (may be nil)
	icon      *Label        // pointer to button icon (may be nil
	styles    *ButtonStyles // pointer to current button styles
	mouseOver bool          // true if mouse is over button
	pressed   bool          // true if button is pressed
}

// ButtonStyle contains the styling of a Button
type ButtonStyle BasicStyle

// ButtonStyles contains one ButtonStyle for each possible button state
type ButtonStyles struct {
	Normal   ButtonStyle
	Over     ButtonStyle
	Focus    ButtonStyle
	Pressed  ButtonStyle
	Disabled ButtonStyle
}

// NewButton creates and returns a pointer to a new button widget
// with the specified text for the button label.
func NewButton(text string) *Button {

	b := new(Button)
	b.styles = &StyleDefault().Button

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

// SetIcon sets the button icon from the default Icon font.
// If there is currently a selected image, it is removed
func (b *Button) SetIcon(icode string) {

	ico := NewIcon(icode)
	if b.image != nil {
		b.Panel.Remove(b.image)
		b.image = nil
	}
	if b.icon != nil {
		b.Panel.Remove(b.icon)
	}
	b.icon = ico
	b.icon.SetFontSize(b.Label.FontSize() * 1.4)
	b.Panel.Add(b.icon)

	b.recalc()
	b.update()
}

// SetImage sets the button left image from the specified filename
// If there is currently a selected icon, it is removed
func (b *Button) SetImage(imgfile string) error {

	img, err := NewImage(imgfile)
	if err != nil {
		return err
	}
	if b.image != nil {
		b.Panel.Remove(b.image)
	}
	b.image = img
	b.Panel.Add(b.image)
	b.recalc()
	return nil
}

// SetStyles set the button styles overriding the default style
func (b *Button) SetStyles(bs *ButtonStyles) {

	b.styles = bs
	b.update()
}

// onCursor process subscribed cursor events
func (b *Button) onCursor(evname string, ev interface{}) {

	switch evname {
	case OnCursorEnter:
		b.mouseOver = true
		b.update()
	case OnCursorLeave:
		b.pressed = false
		b.mouseOver = false
		b.update()
	}
	b.root.StopPropagation(StopAll)
}

// onMouseEvent process subscribed mouse events
func (b *Button) onMouse(evname string, ev interface{}) {

	switch evname {
	case OnMouseDown:
		b.root.SetKeyFocus(b)
		b.pressed = true
		b.update()
		b.Dispatch(OnClick, nil)
	case OnMouseUp:
		b.pressed = false
		b.update()
	default:
		return
	}
	b.root.StopPropagation(StopAll)
}

// onKey processes subscribed key events
func (b *Button) onKey(evname string, ev interface{}) {

	kev := ev.(*window.KeyEvent)
	if evname == OnKeyDown && kev.Keycode == window.KeyEnter {
		b.pressed = true
		b.update()
		b.Dispatch(OnClick, nil)
		b.root.StopPropagation(Stop3D)
		return
	}
	if evname == OnKeyUp && kev.Keycode == window.KeyEnter {
		b.pressed = false
		b.update()
		b.root.StopPropagation(Stop3D)
		return
	}
	return
}

// update updates the button visual state
func (b *Button) update() {

	if !b.Enabled() {
		b.applyStyle(&b.styles.Disabled)
		return
	}
	if b.pressed {
		b.applyStyle(&b.styles.Pressed)
		return
	}
	if b.mouseOver {
		b.applyStyle(&b.styles.Over)
		return
	}
	b.applyStyle(&b.styles.Normal)
}

// applyStyle applies the specified button style
func (b *Button) applyStyle(bs *ButtonStyle) {

	b.Panel.ApplyStyle(&bs.PanelStyle)
	if b.icon != nil {
		b.icon.SetColor4(&bs.FgColor)
	}
	b.Label.SetColor4(&bs.FgColor)
}

// recalc recalculates all dimensions and position from inside out
func (b *Button) recalc() {

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
