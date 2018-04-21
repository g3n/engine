// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/text"
	"github.com/g3n/engine/window"
	"strings"
	"time"
)

// Edit represents a text edit box GUI element
type Edit struct {
	Label              // Embedded label
	MaxLength   int    // Maximum number of characters
	width       int    // edit width in pixels
	placeHolder string // place holder string
	text        string // current edit text
	col         int    // current column
	focus       bool   // key focus flag
	cursorOver  bool
	blinkID     int
	caretOn     bool
	styles      *EditStyles
}

// EditStyle contains the styling of an Edit
type EditStyle struct {
	Border      RectBounds
	Paddings    RectBounds
	BorderColor math32.Color4
	BgColor     math32.Color4
	BgAlpha     float32
	FgColor     math32.Color4
	HolderColor math32.Color4
}

// EditStyles contains an EditStyle for each valid GUI state
type EditStyles struct {
	Normal   EditStyle
	Over     EditStyle
	Focus    EditStyle
	Disabled EditStyle
}

const (
	editMarginX = 4
	blinkTime   = 1000
)

// NewEdit creates and returns a pointer to a new edit widget
func NewEdit(width int, placeHolder string) *Edit {

	ed := new(Edit)
	ed.width = width
	ed.placeHolder = placeHolder

	ed.styles = &StyleDefault().Edit
	ed.text = ""
	ed.MaxLength = 80
	ed.col = 0
	ed.focus = false

	ed.Label.initialize("", StyleDefault().Font)
	ed.Label.Subscribe(OnKeyDown, ed.onKey)
	ed.Label.Subscribe(OnKeyRepeat, ed.onKey)
	ed.Label.Subscribe(OnChar, ed.onChar)
	ed.Label.Subscribe(OnMouseDown, ed.onMouse)
	ed.Label.Subscribe(OnCursorEnter, ed.onCursor)
	ed.Label.Subscribe(OnCursorLeave, ed.onCursor)
	ed.Label.Subscribe(OnEnable, func(evname string, ev interface{}) { ed.update() })

	ed.update()
	return ed
}

// SetText sets this edit text
func (ed *Edit) SetText(text string) *Edit {

	// Remove new lines from text
	ed.text = strings.Replace(text, "\n", "", -1)
	ed.update()
	return ed
}

// Text returns the current edited text
func (ed *Edit) Text() string {

	return ed.text
}

// SetFontSize sets label font size (overrides Label.SetFontSize)
func (ed *Edit) SetFontSize(size float64) *Edit {

	ed.Label.SetFontSize(size)
	ed.redraw(ed.focus)
	return ed
}

// SetStyles set the button styles overriding the default style
func (ed *Edit) SetStyles(es *EditStyles) {

	ed.styles = es
	ed.update()
}

// LostKeyFocus satisfies the IPanel interface and is called by gui root
// container when the panel loses the key focus
func (ed *Edit) LostKeyFocus() {

	ed.focus = false
	ed.update()
	ed.root.ClearTimeout(ed.blinkID)
}

// CursorPos sets the position of the cursor at the
// specified  column if possible
func (ed *Edit) CursorPos(col int) {

	if col <= text.StrCount(ed.text) {
		ed.col = col
		ed.redraw(ed.focus)
	}
}

// CursorLeft moves the edit cursor one character left if possible
func (ed *Edit) CursorLeft() {

	if ed.col > 0 {
		ed.col--
		ed.redraw(ed.focus)
	}
}

// CursorRight moves the edit cursor one character right if possible
func (ed *Edit) CursorRight() {

	if ed.col < text.StrCount(ed.text) {
		ed.col++
		ed.redraw(ed.focus)
	}
}

// CursorBack deletes the character at left of the cursor if possible
func (ed *Edit) CursorBack() {

	if ed.col > 0 {
		ed.col--
		ed.text = text.StrRemove(ed.text, ed.col)
		ed.redraw(ed.focus)
		ed.Dispatch(OnChange, nil)
	}
}

// CursorHome moves the edit cursor to the beginning of the text
func (ed *Edit) CursorHome() {

	ed.col = 0
	ed.redraw(ed.focus)
}

// CursorEnd moves the edit cursor to the end of the text
func (ed *Edit) CursorEnd() {

	ed.col = text.StrCount(ed.text)
	ed.redraw(ed.focus)
}

// CursorDelete deletes the character at the right of the cursor if possible
func (ed *Edit) CursorDelete() {

	if ed.col < text.StrCount(ed.text) {
		ed.text = text.StrRemove(ed.text, ed.col)
		ed.redraw(ed.focus)
		ed.Dispatch(OnChange, nil)
	}
}

// CursorInput inserts the specified string at the current cursor position
func (ed *Edit) CursorInput(s string) {

	if text.StrCount(ed.text) >= ed.MaxLength {
		return
	}

	// Set new text with included input
	var newText string
	if ed.col < text.StrCount(ed.text) {
		newText = text.StrInsert(ed.text, s, ed.col)
	} else {
		newText = ed.text + s
	}

	// Checks if new text exceeds edit width
	width, _ := ed.Label.font.MeasureText(newText)
	if float32(width)+editMarginX+float32(1) >= ed.Label.ContentWidth() {
		return
	}

	ed.text = newText
	ed.col++

	ed.Dispatch(OnChange, nil)
	ed.redraw(ed.focus)
}

// redraw redraws the text showing the caret if specified
func (ed *Edit) redraw(caret bool) {

	line := 0
	if !caret {
		line = -1
	}
	ed.Label.setTextCaret(ed.text, editMarginX, ed.width, line, ed.col)
}

// onKey receives subscribed key events
func (ed *Edit) onKey(evname string, ev interface{}) {

	kev := ev.(*window.KeyEvent)
	switch kev.Keycode {
	case window.KeyLeft:
		ed.CursorLeft()
	case window.KeyRight:
		ed.CursorRight()
	case window.KeyHome:
		ed.CursorHome()
	case window.KeyEnd:
		ed.CursorEnd()
	case window.KeyBackspace:
		ed.CursorBack()
	case window.KeyDelete:
		ed.CursorDelete()
	default:
		return
	}
	ed.root.StopPropagation(Stop3D)
}

// onChar receives subscribed char events
func (ed *Edit) onChar(evname string, ev interface{}) {

	cev := ev.(*window.CharEvent)
	ed.CursorInput(string(cev.Char))
}

// onMouseEvent receives subscribed mouse down events
func (ed *Edit) onMouse(evname string, ev interface{}) {

	e := ev.(*window.MouseEvent)
	if e.Button != window.MouseButtonLeft {
		return
	}

	// Set key focus to this panel
	ed.root.SetKeyFocus(ed)

	// Find clicked column
	var nchars int
	for nchars = 1; nchars <= text.StrCount(ed.text); nchars++ {
		width, _ := ed.Label.font.MeasureText(text.StrPrefix(ed.text, nchars))
		posx := e.Xpos - ed.pospix.X
		if posx < editMarginX+float32(width) {
			break
		}
	}
	if !ed.focus {
		ed.focus = true
		ed.blinkID = ed.root.SetInterval(750*time.Millisecond, nil, ed.blink)
	}
	ed.CursorPos(nchars - 1)
	ed.root.StopPropagation(Stop3D)
}

// onCursor receives subscribed cursor events
func (ed *Edit) onCursor(evname string, ev interface{}) {

	if evname == OnCursorEnter {
		ed.root.SetCursorText()
		ed.cursorOver = true
		ed.update()
		ed.root.StopPropagation(Stop3D)
		return
	}
	if evname == OnCursorLeave {
		ed.root.SetCursorNormal()
		ed.cursorOver = false
		ed.update()
		ed.root.StopPropagation(Stop3D)
		return
	}
}

// blink blinks the caret
func (ed *Edit) blink(arg interface{}) {

	if !ed.focus {
		return
	}
	if !ed.caretOn {
		ed.caretOn = true
	} else {
		ed.caretOn = false
	}
	ed.redraw(ed.caretOn)
}

// update updates the visual state
func (ed *Edit) update() {

	if !ed.Enabled() {
		ed.applyStyle(&ed.styles.Disabled)
		return
	}
	if ed.cursorOver {
		ed.applyStyle(&ed.styles.Over)
		return
	}
	if ed.focus {
		ed.applyStyle(&ed.styles.Focus)
		return
	}
	ed.applyStyle(&ed.styles.Normal)
}

// applyStyle applies the specified style
func (ed *Edit) applyStyle(s *EditStyle) {

	ed.SetBordersFrom(&s.Border)
	ed.SetBordersColor4(&s.BorderColor)
	ed.SetPaddingsFrom(&s.Paddings)
	ed.Label.SetColor4(&s.FgColor)
	ed.Label.SetBgColor4(&s.BgColor)
	//ed.Label.SetBgAlpha(s.BgAlpha)

	if !ed.focus && len(ed.text) == 0 && len(ed.placeHolder) > 0 {
		ed.Label.SetColor4(&s.HolderColor)
		ed.Label.setTextCaret(ed.placeHolder, editMarginX, ed.width, -1, ed.col)
	} else {
		ed.Label.SetColor4(&s.FgColor)
		ed.redraw(ed.focus)
	}
}
