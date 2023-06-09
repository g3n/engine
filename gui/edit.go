// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"strings"
	"time"

	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/text"
	"github.com/g3n/engine/window"
)

// Edit represents a text edit box GUI element
type Edit struct {
	Label              // Embedded label
	MaxLength   int    // Maximum number of characters
	width       int    // edit width in pixels
	placeHolder string // place holder string
	text        string // current edit text
	col         int    // current column
	selStart    int    // start column of selection. always < selEnd. if selStart == selEnd then nothing is selected.
	selEnd      int    // end column of selection. always > selStart. if selStart == selEnd then nothing is selected.
	focus       bool   // key focus flag
	cursorOver  bool
	mouseDrag   bool // true when the mouse is moved while left mouse button is down. Used for selecting text via mouse
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
	ed.selStart = 0
	ed.selEnd = 0
	ed.focus = false

	ed.Label.initialize("", StyleDefault().Font)
	ed.Label.Subscribe(OnKeyDown, ed.onKey)
	ed.Label.Subscribe(OnKeyRepeat, ed.onKey)
	ed.Label.Subscribe(OnChar, ed.onChar)
	ed.Label.Subscribe(OnMouseDown, ed.onMouseDown)
	ed.Label.Subscribe(OnMouseUp, ed.onMouseUp)
	ed.Label.Subscribe(OnCursorEnter, ed.onCursor)
	ed.Label.Subscribe(OnCursorLeave, ed.onCursor)
	ed.Label.Subscribe(OnCursor, ed.onCursor)
	ed.Label.Subscribe(OnEnable, func(evname string, ev interface{}) { ed.update() })
	ed.Subscribe(OnFocusLost, ed.OnFocusLost)

	ed.update()
	return ed
}

// SetText sets this edit text
func (ed *Edit) SetText(newText string) *Edit {

	// Remove new lines from text
	ed.text = strings.Replace(newText, "\n", "", -1)
	ed.col = text.StrCount(ed.text)
	ed.selStart = ed.col
	ed.selEnd = ed.col
	ed.update()
	return ed
}

// Text returns the current edited text
func (ed *Edit) Text() string {

	return ed.text
}

// SelectedText returns the currently selected text
// or empty string when nothing is selected
func (ed *Edit) SelectedText() string {

	if ed.selStart == ed.selEnd {
		return ""
	}

	s := ""
	charNum := 0
	for _, currentRune := range ed.text {
		if charNum >= ed.selEnd {
			break
		}
		if charNum >= ed.selStart {
			s += string(currentRune)
		}
		charNum++
	}
	return s
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
func (ed *Edit) OnFocusLost(evname string, ev interface{}) {

	ed.focus = false
	ed.update()
	Manager().ClearTimeout(ed.blinkID)
}

// CursorPos sets the position of the cursor at the
// specified  column if possible
func (ed *Edit) CursorPos(col int) {

	if col <= text.StrCount(ed.text) {
		ed.col = col
		ed.selStart = col
		ed.selEnd = col
		ed.redraw(ed.focus)
	}
}

// SetSelection selects the text between start and end
func (ed *Edit) SetSelection(start, end int) {

	// make sure end is bigger than start
	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if end > text.StrCount(ed.text) {
		end = text.StrCount(ed.text)
	}

	ed.selStart = start
	ed.selEnd = end
	ed.col = end
	ed.redraw(ed.focus)
}

// CursorLeft moves the edit cursor one character left if possible
// If text is selected the cursor is moved to the beginning of the selection instead
// and the selection is removed
func (ed *Edit) CursorLeft() {

	if ed.selStart == ed.selEnd {
		// no selection
		// move cursor to the left if possible
		if ed.col > 0 {
			ed.col--
			ed.selStart = ed.col
			ed.selEnd = ed.col
			ed.redraw(ed.focus)
		}
	} else {
		// reset selection and move cursor to start of selection
		ed.col = ed.selStart
		ed.selStart = ed.col
		ed.selEnd = ed.col
		ed.redraw(ed.focus)
	}
}

// CursorRight moves the edit cursor one character right if possible
// If text is selected the cursor is moved to the end of the selection instead
// and the selection is removed
func (ed *Edit) CursorRight() {

	if ed.selStart == ed.selEnd {
		// no selection
		// move cursor to the right if possible
		if ed.col < text.StrCount(ed.text) {
			ed.col++
			ed.selStart = ed.col
			ed.selEnd = ed.col
			ed.redraw(ed.focus)
		}
	} else {
		// reset selection and move cursor to end of selection
		ed.col = ed.selEnd
		ed.selStart = ed.col
		ed.selEnd = ed.col
		ed.redraw(ed.focus)
	}
}

// SelectLeft expands/shrinks the selection to the left if possible
func (ed *Edit) SelectLeft() {

	if ed.col > 0 {
		if ed.col == ed.selStart {
			// cursor is at the start of selection
			// expand selection to the left
			ed.col--
			ed.selStart = ed.col
			ed.redraw(ed.focus)
		} else {
			// cursor is at the end of selection:
			// remove selection from the end
			ed.col--
			ed.selEnd = ed.col
			ed.redraw(ed.focus)
		}
	}
}

// SelectRight expands/shrinks the selection to the right if possible
func (ed *Edit) SelectRight() {

	if ed.col < text.StrCount(ed.text) {
		if ed.col == ed.selEnd {
			// cursor is at the end of selection:
			// expand selection to the right if possible
			ed.col++
			ed.selEnd = ed.col
			ed.redraw(ed.focus)
		} else {
			// cursor is at the start of selection:
			// remove selection from the start
			ed.col++
			ed.selStart = ed.col
			ed.redraw(ed.focus)
		}
	}
}

// SelectHome expands the selection to the left to the beginning of the text
func (ed *Edit) SelectHome() {

	if ed.selStart < ed.col {
		ed.selEnd = ed.selStart
	}
	ed.col = 0
	ed.selStart = 0
	ed.redraw(ed.focus)
}

// SelectEnd expands the selection to the right to the end of the text
func (ed *Edit) SelectEnd() {

	if ed.selEnd > ed.col {
		ed.selStart = ed.selEnd
	}
	ed.col = text.StrCount(ed.text)
	ed.selEnd = ed.col
	ed.redraw(ed.focus)
}

// SelectAll selects all text
func (ed *Edit) SelectAll() {

	ed.selStart = 0
	ed.selEnd = text.StrCount(ed.text)
	ed.col = ed.selEnd
	ed.redraw(ed.focus)
}

// CursorBack either deletes the character at left of the cursor if possible
// Or if text is selected the selected text is removed all at once
func (ed *Edit) CursorBack() {

	if ed.selStart == ed.selEnd {
		if ed.col > 0 {
			ed.col--
			ed.selStart = ed.col
			ed.selEnd = ed.col
			ed.text = text.StrRemove(ed.text, ed.col)
			ed.redraw(ed.focus)
			ed.Dispatch(OnChange, nil)
		}
	} else {
		ed.DeleteSelection()
	}
}

// CursorHome moves the edit cursor to the beginning of the text
func (ed *Edit) CursorHome() {

	ed.col = 0
	ed.selStart = ed.col
	ed.selEnd = ed.col
	ed.redraw(ed.focus)
}

// CursorEnd moves the edit cursor to the end of the text
func (ed *Edit) CursorEnd() {

	ed.col = text.StrCount(ed.text)
	ed.selStart = ed.col
	ed.selEnd = ed.col
	ed.redraw(ed.focus)
}

// CursorDelete either deletes the character at the right of the cursor if possible
// Or if text is selected the selected text is removed all at once
func (ed *Edit) CursorDelete() {

	if ed.selStart == ed.selEnd {
		if ed.col < text.StrCount(ed.text) {
			ed.text = text.StrRemove(ed.text, ed.col)
			ed.redraw(ed.focus)
			ed.Dispatch(OnChange, nil)
		}
	} else {
		ed.DeleteSelection()
	}
}

// DeleteSelection deletes the selected characters. Does nothing if nothing is selected.
func (ed *Edit) DeleteSelection() {

	if ed.selStart == ed.selEnd {
		return
	}

	changed := false
	ed.col = ed.selStart
	for ed.selEnd > ed.selStart {
		if ed.col < text.StrCount(ed.text) {
			changed = true
			ed.text = text.StrRemove(ed.text, ed.col)
			ed.selEnd--
		}
	}
	if changed {
		ed.Dispatch(OnChange, nil)
		ed.redraw(ed.focus)
	}
}

// CursorInput inserts the specified string at the current cursor position
// If text is selected the selected text gets overwritten
func (ed *Edit) CursorInput(s string) {

	if ed.selStart != ed.selEnd {
		ed.DeleteSelection()
	}
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
	ed.selStart = ed.col
	ed.selEnd = ed.col

	ed.Dispatch(OnChange, nil)
	ed.redraw(ed.focus)
}

// redraw redraws the text showing the caret if specified
// the selection caret is always shown (when text is selected)
func (ed *Edit) redraw(caret bool) {

	line := 0
	ed.Label.setTextCaret(ed.text, editMarginX, int(ed.Width()), caret, line, ed.col, ed.selStart, ed.selEnd)
}

// onKey receives subscribed key events
func (ed *Edit) onKey(evname string, ev interface{}) {

	kev := ev.(*window.KeyEvent)
	if kev.Mods != window.ModShift && kev.Mods != window.ModControl {
		switch kev.Key {
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
	} else if kev.Mods == window.ModShift {
		switch kev.Key {
		case window.KeyLeft:
			ed.SelectLeft()
		case window.KeyRight:
			ed.SelectRight()
		case window.KeyHome:
			ed.SelectHome()
		case window.KeyEnd:
			ed.SelectEnd()
		case window.KeyBackspace:
			ed.CursorBack()
		case window.KeyDelete:
			ed.SelectAll()
			ed.DeleteSelection()
		default:
			return
		}
	} else if kev.Mods == window.ModControl {
		switch kev.Key {
		case window.KeyA:
			ed.SelectAll()
		}
	}
}

// onChar receives subscribed char events
func (ed *Edit) onChar(evname string, ev interface{}) {

	cev := ev.(*window.CharEvent)
	ed.CursorInput(string(cev.Char))
}

// onMouseDown receives subscribed mouse down events
func (ed *Edit) onMouseDown(evname string, ev interface{}) {

	e := ev.(*window.MouseEvent)
	if e.Button != window.MouseButtonLeft {
		return
	}

	// set caret to clicked position
	ed.handleMouse(e.Xpos, false)

	ed.mouseDrag = true

	// Set key focus to this panel
	// Set the focus AFTER the mouse selection is handled
	// Otherwise the OnFocus event would fire before the cursor is set.
	// That way the OnFocus handler could NOT influence the selection
	// Because it would be overridden/cleared directly afterwards.
	Manager().SetKeyFocus(ed)
}

// handleMouse is setting the caret when the mouse is clicked
// or setting the text selection when the mouse is dragged
func (ed *Edit) handleMouse(mouseX float32, dragged bool) {

	// Find clicked column
	var nchars int
	for nchars = 1; nchars <= text.StrCount(ed.text); nchars++ {
		width, _ := ed.Label.font.MeasureText(text.StrPrefix(ed.text, nchars))
		posx := mouseX - ed.pospix.X
		if posx < editMarginX+float32(width) {
			break
		}
	}
	if !ed.focus {
		ed.focus = true
		ed.blinkID = Manager().SetInterval(750*time.Millisecond, nil, ed.blink)
	}
	if !dragged {
		ed.CursorPos(nchars - 1)
	} else {
		newPos := nchars - 1
		if newPos > ed.col {
			distance := newPos - ed.col
			for i := 0; i < distance; i++ {
				ed.SelectRight()
			}
		} else if newPos < ed.col {
			distance := ed.col - newPos
			for i := 0; i < distance; i++ {
				ed.SelectLeft()
			}
		}
	}
}

// onMouseEvent receives subscribed mouse up events
func (ed *Edit) onMouseUp(evname string, ev interface{}) {

	ed.mouseDrag = false
}

// onCursor receives subscribed cursor events
func (ed *Edit) onCursor(evname string, ev interface{}) {

	if evname == OnCursorEnter {
		window.Get().SetCursor(window.IBeamCursor)
		ed.cursorOver = true
		ed.update()
		return
	}
	if evname == OnCursorLeave {
		window.Get().SetCursor(window.ArrowCursor)
		ed.cursorOver = false
		ed.mouseDrag = false
		ed.update()
		return
	}
	if ed.mouseDrag {
		e := ev.(*window.CursorEvent)
		// select text based on mouse position
		ed.handleMouse(e.Xpos, true)
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
		ed.Label.setTextCaret(ed.placeHolder, editMarginX, int(ed.Width()), false, -1, ed.col, ed.selStart, ed.selEnd)
	} else {
		ed.Label.SetColor4(&s.FgColor)
		ed.redraw(ed.focus)
	}
}
