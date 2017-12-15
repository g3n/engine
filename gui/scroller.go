// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/window"
	"math"
)

type Scroller struct {
	Panel                          // Embedded panel
	vert           bool            // vertical/horizontal scroller flag
	styles         *ScrollerStyles // pointer to current styles
	items          []IPanel        // list of panels in the scroller
	hscroll        *ScrollBar      // horizontal scroll bar
	vscroll        *ScrollBar      // vertical scroll bar
	maxAutoWidth   float32         // maximum auto width (if 0, auto width disabled)
	maxAutoHeight  float32         // maximum auto height (if 0, auto width disabled)
	first          int             // first visible item position
	adjustItem     bool            // adjust item to width or height
	focus          bool            // has keyboard focus
	cursorOver     bool            // mouse is over the list
	scrollBarEvent bool
}

type ScrollerStyle struct {
	Border      BorderSizes
	Paddings    BorderSizes
	BorderColor math32.Color4
	BgColor     math32.Color
	FgColor     math32.Color
}

type ScrollerStyles struct {
	Normal   ScrollerStyle
	Over     ScrollerStyle
	Focus    ScrollerStyle
	Disabled ScrollerStyle
}

// NewVScroller creates and returns a pointer to a new vertical scroller panel
// with the specified dimensions.
func NewVScroller(width, height float32) *Scroller {

	return newScroller(true, width, height)
}

// NewHScroller creates and returns a pointer to a new horizontal scroller panel
// with the specified dimensions.
func NewHScroller(width, height float32) *Scroller {

	return newScroller(false, width, height)
}

// newScroller creates and returns a pointer to a new Scroller panel
// with the specified layout orientation and initial dimensions
func newScroller(vert bool, width, height float32) *Scroller {

	s := new(Scroller)
	s.initialize(vert, width, height)
	return s
}

// Clear removes and disposes of all the scroller children
func (s *Scroller) Clear() {

	s.Panel.DisposeChildren(true)
	s.first = 0
	s.hscroll = nil
	s.vscroll = nil
	s.items = s.items[0:0]
	s.update()
	s.recalc()
}

// Len return the number of items in the scroller
func (s *Scroller) Len() int {

	return len(s.items)
}

// Add appends the specified item to the end of the scroller
func (s *Scroller) Add(item IPanel) {

	s.InsertAt(len(s.items), item)
}

// InsertAt inserts an item at the specified position
func (s *Scroller) InsertAt(pos int, item IPanel) {

	// Validates position
	if pos < 0 || pos > len(s.items) {
		panic("Scroller.InsertAt(): Invalid position")
	}
	item.GetPanel().SetVisible(false)

	// Insert item in the items array
	s.items = append(s.items, nil)
	copy(s.items[pos+1:], s.items[pos:])
	s.items[pos] = item

	// Insert item in the scroller
	s.Panel.Add(item)
	s.autoSize()
	s.recalc()

	// Scroll bar should be on the foreground,
	// in relation of all the other child panels.
	if s.vscroll != nil {
		s.Panel.SetTopChild(s.vscroll)
	}
	if s.hscroll != nil {
		s.Panel.SetTopChild(s.hscroll)
	}
}

// RemoveAt removes item from the specified position
func (s *Scroller) RemoveAt(pos int) IPanel {

	// Validates position
	if pos < 0 || pos >= len(s.items) {
		panic("Scroller.RemoveAt(): Invalid position")
	}

	// Remove event listener
	item := s.items[pos]

	// Remove item from the items array
	copy(s.items[pos:], s.items[pos+1:])
	s.items[len(s.items)-1] = nil
	s.items = s.items[:len(s.items)-1]

	// Remove item from the scroller children
	s.Panel.Remove(item)
	s.autoSize()
	s.recalc()
	return item
}

// Remove removes the specified item from the Scroller
func (s *Scroller) Remove(item IPanel) {

	for p, curr := range s.items {
		if curr == item {
			s.RemoveAt(p)
			return
		}
	}
}

// GetItem returns the item at the specified position.
// Returns nil if the position is invalid.
func (s *Scroller) ItemAt(pos int) IPanel {

	if pos < 0 || pos >= len(s.items) {
		return nil
	}
	return s.items[pos]
}

// ItemPosition returns the position of the specified item in
// the scroller of -1 if not found
func (s *Scroller) ItemPosition(item IPanel) int {

	for pos := 0; pos < len(s.items); pos++ {
		if s.items[pos] == item {
			return pos
		}
	}
	return -1
}

// First returns the position of the first visible item
func (s *Scroller) First() int {

	return s.first
}

// SetFirst set the position of first visible if possible
func (s *Scroller) SetFirst(pos int) {

	if pos >= 0 && pos <= s.maxFirst() {
		s.first = pos
		s.recalc()
	}
}

// ScrollDown scrolls the list down one item if possible
func (s *Scroller) ScrollDown() {

	max := s.maxFirst()
	if s.first >= max {
		return
	}
	s.first++
	s.recalc()
}

// ScrollUp scrolls the list up one item if possible
func (s *Scroller) ScrollUp() {

	if s.first == 0 {
		return
	}
	s.first--
	s.recalc()
}

// ItemVisible returns indication if the item at the specified
// position is completely visible or not
func (s *Scroller) ItemVisible(pos int) bool {

	if pos < s.first {
		return false
	}

	// Vertical scroller
	if s.vert {
		var height float32 = 0
		for i := s.first; i < len(s.items); i++ {
			item := s.items[pos]
			height += item.GetPanel().height
			if height > s.height {
				return false
			}
			if pos == i {
				return true
			}
		}
		return false
		// Horizontal scroller
	} else {
		var width float32 = 0
		for i := s.first; i < len(s.items); i++ {
			item := s.items[pos]
			width += item.GetPanel().width
			if width > s.width {
				return false
			}
			if pos == i {
				return true
			}
		}
		return false
	}
}

// SetStyles set the scroller styles overriding the default style
func (s *Scroller) SetStyles(ss *ScrollerStyles) {

	s.styles = ss
	s.update()
}

func (s *Scroller) ApplyStyle(style int) {

	switch style {
	case StyleOver:
		s.applyStyle(&s.styles.Over)
	case StyleFocus:
		s.applyStyle(&s.styles.Focus)
	case StyleNormal:
		s.applyStyle(&s.styles.Normal)
	case StyleDef:
		s.update()
	}
}

func (s *Scroller) SetAutoWidth(maxWidth float32) {

	s.maxAutoWidth = maxWidth
}

func (s *Scroller) SetAutoHeight(maxHeight float32) {

	s.maxAutoHeight = maxHeight
}

// initialize initializes this scroller and is normally used by other types which contains a scroller
func (s *Scroller) initialize(vert bool, width, height float32) {

	s.vert = vert
	s.Panel.Initialize(width, height)
	s.styles = &StyleDefault().Scroller

	s.Panel.Subscribe(OnCursorEnter, s.onCursor)
	s.Panel.Subscribe(OnCursorLeave, s.onCursor)
	s.Panel.Subscribe(OnScroll, s.onScroll)
	s.Panel.Subscribe(OnResize, s.onResize)

	s.update()
	s.recalc()
}

// onCursor receives subscribed cursor events over the panel
func (s *Scroller) onCursor(evname string, ev interface{}) {

	switch evname {
	case OnCursorEnter:
		s.root.SetScrollFocus(s)
		s.cursorOver = true
		s.update()
	case OnCursorLeave:
		s.root.SetScrollFocus(nil)
		s.cursorOver = false
		s.update()
	}
	s.root.StopPropagation(Stop3D)
}

// onScroll receives subscriber mouse scroll events when this scroller has
// the scroll focus (set by OnMouseEnter)
func (s *Scroller) onScroll(evname string, ev interface{}) {

	sev := ev.(*window.ScrollEvent)
	if sev.Yoffset > 0 {
		s.ScrollUp()
	} else if sev.Yoffset < 0 {
		s.ScrollDown()
	}
	s.root.StopPropagation(Stop3D)
}

// onScroll receives resize events
func (s *Scroller) onResize(evname string, ev interface{}) {

	s.recalc()
}

// autoSize resizes the scroller if necessary
func (s *Scroller) autoSize() {

	if s.maxAutoWidth == 0 && s.maxAutoHeight == 0 {
		return
	}

	var width float32 = 0
	var height float32 = 0
	for _, item := range s.items {
		panel := item.GetPanel()
		if panel.Width() > width {
			width = panel.Width()
		}
		height += panel.TotalHeight()
	}

	// If auto maximum width enabled
	if s.maxAutoWidth > 0 {
		if width <= s.maxAutoWidth {
			s.SetContentWidth(width)
		}
	}
	// If auto maximum height enabled
	if s.maxAutoHeight > 0 {
		if height <= s.maxAutoHeight {
			s.SetContentHeight(height)
		}
	}
}

// recalc recalculates the positions and visibilities of all the items
func (s *Scroller) recalc() {

	if s.vert {
		s.vRecalc()
	} else {
		s.hRecalc()
	}
}

// vRecalc recalculates for the vertical scroller
func (s *Scroller) vRecalc() {

	// Checks if scroll bar should be visible or not
	scroll := false
	if s.first > 0 {
		scroll = true
	} else {
		var posY float32 = 0
		for _, item := range s.items[s.first:] {
			posY += item.TotalHeight()
			if posY >= s.height {
				scroll = true
				break
			}
		}
	}
	s.setVScrollBar(scroll)
	// Items width
	width := s.ContentWidth()
	if scroll {
		width -= s.vscroll.Width()
	}

	var posY float32 = 0
	// Sets positions of all items
	for pos, ipan := range s.items {
		item := ipan.GetPanel()
		if pos < s.first {
			item.SetVisible(false)
			continue
		}
		// If item is after last visible, sets not visible
		if posY > s.height {
			item.SetVisible(false)
			continue
		}
		// Sets item position
		item.SetVisible(true)
		item.SetPosition(0, posY)
		if s.adjustItem {
			item.SetWidth(width)
		}
		posY += ipan.TotalHeight()
	}

	// Set scroll bar value if recalc was not due by scroll event
	if scroll && !s.scrollBarEvent {
		s.vscroll.SetValue(float32(s.first) / float32(s.maxFirst()))
	}
	s.scrollBarEvent = false
}

// hRecalc recalculates for the horizontal scroller
func (s *Scroller) hRecalc() {

	// Checks if scroll bar should be visible or not
	scroll := false
	if s.first > 0 {
		scroll = true
	} else {
		var posX float32 = 0
		for _, item := range s.items[s.first:] {
			posX += item.GetPanel().Width()
			if posX >= s.width {
				scroll = true
				break
			}
		}
	}
	s.setHScrollBar(scroll)
	// Items height
	height := s.ContentHeight()
	if scroll {
		height -= s.hscroll.Height()
	}

	var posX float32 = 0
	// Sets positions of all items
	for pos, ipan := range s.items {
		item := ipan.GetPanel()
		// If item is before first visible, sets not visible
		if pos < s.first {
			item.SetVisible(false)
			continue
		}
		// If item is after last visible, sets not visible
		if posX > s.width {
			item.SetVisible(false)
			continue
		}
		// Sets item position
		item.SetVisible(true)
		item.SetPosition(posX, 0)
		if s.adjustItem {
			item.SetHeight(height)
		}
		posX += item.Width()
	}

	// Set scroll bar value if recalc was not due by scroll event
	if scroll && !s.scrollBarEvent {
		s.hscroll.SetValue(float32(s.first) / float32(s.maxFirst()))
	}
	s.scrollBarEvent = false
}

// maxFirst returns the maximum position of the first visible item
func (s *Scroller) maxFirst() int {

	// Vertical scroller
	if s.vert {
		var height float32 = 0
		pos := len(s.items) - 1
		if pos < 0 {
			return 0
		}
		for {
			item := s.items[pos]
			height += item.GetPanel().Height()
			if height > s.Height() {
				break
			}
			pos--
			if pos < 0 {
				break
			}
		}
		return pos + 1
		// Horizontal scroller
	} else {
		var width float32 = 0
		pos := len(s.items) - 1
		if pos < 0 {
			return 0
		}
		for {
			item := s.items[pos]
			width += item.GetPanel().Width()
			if width > s.Width() {
				break
			}
			pos--
			if pos < 0 {
				break
			}
		}
		return pos + 1
	}
}

// setVScrollBar sets the visibility state of the vertical scrollbar
func (s *Scroller) setVScrollBar(state bool) {

	// Visible
	if state {
		var scrollWidth float32 = 20
		if s.vscroll == nil {
			s.vscroll = NewVScrollBar(0, 0)
			s.vscroll.SetBorders(0, 0, 0, 1)
			s.vscroll.Subscribe(OnChange, s.onScrollBarEvent)
			s.Panel.Add(s.vscroll)
		}
		s.vscroll.SetSize(scrollWidth, s.ContentHeight())
		s.vscroll.SetPositionX(s.ContentWidth() - scrollWidth)
		s.vscroll.SetPositionY(0)
		s.vscroll.recalc()
		s.vscroll.SetVisible(true)
		// Not visible
	} else {
		if s.vscroll != nil {
			s.vscroll.SetVisible(false)
		}
	}
}

// setHScrollBar sets the visibility state of the horizontal scrollbar
func (s *Scroller) setHScrollBar(state bool) {

	// Visible
	if state {
		var scrollHeight float32 = 20
		if s.hscroll == nil {
			s.hscroll = NewHScrollBar(0, 0)
			s.hscroll.SetBorders(1, 0, 0, 0)
			s.hscroll.Subscribe(OnChange, s.onScrollBarEvent)
			s.Panel.Add(s.hscroll)
		}
		s.hscroll.SetSize(s.ContentWidth(), scrollHeight)
		s.hscroll.SetPositionX(0)
		s.hscroll.SetPositionY(s.ContentHeight() - scrollHeight)
		s.hscroll.recalc()
		s.hscroll.SetVisible(true)
		// Not visible
	} else {
		if s.hscroll != nil {
			s.hscroll.SetVisible(false)
		}
	}
}

// onScrollEvent is called when the list scrollbar value changes
func (s *Scroller) onScrollBarEvent(evname string, ev interface{}) {

	var pos float64
	if s.vert {
		pos = s.vscroll.Value()
	} else {
		pos = s.hscroll.Value()
	}

	first := int(math.Floor((float64(s.maxFirst()) * pos) + 0.5))
	if first == s.first {
		return
	}
	s.scrollBarEvent = true
	s.first = first
	s.recalc()
}

// update updates the visual state the list and its items
func (s *Scroller) update() {

	if s.cursorOver {
		s.applyStyle(&s.styles.Over)
		return
	}
	if s.focus {
		s.applyStyle(&s.styles.Focus)
		return
	}
	s.applyStyle(&s.styles.Normal)
}

// applyStyle sets the specified style
func (s *Scroller) applyStyle(st *ScrollerStyle) {

	s.SetBordersFrom(&st.Border)
	s.SetBordersColor4(&st.BorderColor)
	s.SetPaddingsFrom(&st.Paddings)
	s.SetColor(&st.BgColor)
}
