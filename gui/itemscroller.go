// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/g3n/engine/window"
	"math"
)

// ItemScroller is the GUI element that allows scrolling of IPanels
type ItemScroller struct {
	Panel                              // Embedded panel
	vert           bool                // vertical/horizontal scroller flag
	styles         *ItemScrollerStyles // pointer to current styles
	items          []IPanel            // list of panels in the scroller
	hscroll        *ScrollBar          // horizontal scroll bar
	vscroll        *ScrollBar          // vertical scroll bar
	maxAutoWidth   float32             // maximum auto width (if 0, auto width disabled)
	maxAutoHeight  float32             // maximum auto height (if 0, auto width disabled)
	first          int                 // first visible item position
	adjustItem     bool                // adjust item to width or height
	focus          bool                // has keyboard focus
	cursorOver     bool                // mouse is over the list
	autoButtonSize bool                // scroll button size is adjusted relative to content/view
	scrollBarEvent bool
}

// ItemScrollerStyle contains the styling of a ItemScroller
type ItemScrollerStyle BasicStyle

// ItemScrollerStyles contains a ItemScrollerStyle for each valid GUI state
type ItemScrollerStyles struct {
	Normal   ItemScrollerStyle
	Over     ItemScrollerStyle
	Focus    ItemScrollerStyle
	Disabled ItemScrollerStyle
}

// NewVScroller creates and returns a pointer to a new vertical scroller panel
// with the specified dimensions.
func NewVScroller(width, height float32) *ItemScroller {

	return newScroller(true, width, height)
}

// NewHScroller creates and returns a pointer to a new horizontal scroller panel
// with the specified dimensions.
func NewHScroller(width, height float32) *ItemScroller {

	return newScroller(false, width, height)
}

// newScroller creates and returns a pointer to a new ItemScroller panel
// with the specified layout orientation and initial dimensions
func newScroller(vert bool, width, height float32) *ItemScroller {

	s := new(ItemScroller)
	s.initialize(vert, width, height)
	return s
}

// Clear removes and disposes of all the scroller children
func (s *ItemScroller) Clear() {

	s.Panel.DisposeChildren(true)
	s.first = 0
	s.hscroll = nil
	s.vscroll = nil
	s.items = s.items[0:0]
	s.update()
	s.recalc()
}

// Len return the number of items in the scroller
func (s *ItemScroller) Len() int {

	return len(s.items)
}

// Add appends the specified item to the end of the scroller
func (s *ItemScroller) Add(item IPanel) {

	s.InsertAt(len(s.items), item)
}

// InsertAt inserts an item at the specified position
func (s *ItemScroller) InsertAt(pos int, item IPanel) {

	// Validates position
	if pos < 0 || pos > len(s.items) {
		panic("ItemScroller.InsertAt(): Invalid position")
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
func (s *ItemScroller) RemoveAt(pos int) IPanel {

	// Validates position
	if pos < 0 || pos >= len(s.items) {
		panic("ItemScroller.RemoveAt(): Invalid position")
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

// Remove removes the specified item from the ItemScroller
func (s *ItemScroller) Remove(item IPanel) {

	for p, curr := range s.items {
		if curr == item {
			s.RemoveAt(p)
			return
		}
	}
}

// ItemAt returns the item at the specified position.
// Returns nil if the position is invalid.
func (s *ItemScroller) ItemAt(pos int) IPanel {

	if pos < 0 || pos >= len(s.items) {
		return nil
	}
	return s.items[pos]
}

// ItemPosition returns the position of the specified item in
// the scroller of -1 if not found
func (s *ItemScroller) ItemPosition(item IPanel) int {

	for pos := 0; pos < len(s.items); pos++ {
		if s.items[pos] == item {
			return pos
		}
	}
	return -1
}

// First returns the position of the first visible item
func (s *ItemScroller) First() int {

	return s.first
}

// SetFirst set the position of first visible if possible
func (s *ItemScroller) SetFirst(pos int) {

	if pos >= 0 && pos <= s.maxFirst() {
		s.first = pos
		s.recalc()
	}
}

// ScrollDown scrolls the list down one item if possible
func (s *ItemScroller) ScrollDown() {

	max := s.maxFirst()
	if s.first >= max {
		return
	}
	s.first++
	s.recalc()
}

// ScrollUp scrolls the list up one item if possible
func (s *ItemScroller) ScrollUp() {

	if s.first == 0 {
		return
	}
	s.first--
	s.recalc()
}

// ItemVisible returns indication if the item at the specified
// position is completely visible or not
func (s *ItemScroller) ItemVisible(pos int) bool {

	if pos < s.first {
		return false
	}

	// Vertical scroller
	if s.vert {
		var height float32
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
	}

	// Horizontal scroller
	var width float32
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

// SetStyles set the scroller styles overriding the default style
func (s *ItemScroller) SetStyles(ss *ItemScrollerStyles) {

	s.styles = ss
	s.update()
}

// ApplyStyle applies the specified style to the ItemScroller
func (s *ItemScroller) ApplyStyle(style int) {

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

// SetAutoWidth sets the maximum automatic width
func (s *ItemScroller) SetAutoWidth(maxWidth float32) {

	s.maxAutoWidth = maxWidth
}

// SetAutoHeight sets the maximum automatic height
func (s *ItemScroller) SetAutoHeight(maxHeight float32) {

	s.maxAutoHeight = maxHeight
}

// SetAutoButtonSize specified whether the scrollbutton size should be adjusted relative to the size of the content/view
func (s *ItemScroller) SetAutoButtonSize(autoButtonSize bool) {

	s.autoButtonSize = autoButtonSize
}

// initialize initializes this scroller and is normally used by other types which contains a scroller
func (s *ItemScroller) initialize(vert bool, width, height float32) {

	s.vert = vert
	s.autoButtonSize = true
	s.Panel.Initialize(width, height)
	s.styles = &StyleDefault().ItemScroller

	s.Panel.Subscribe(OnCursorEnter, s.onCursor)
	s.Panel.Subscribe(OnCursorLeave, s.onCursor)
	s.Panel.Subscribe(OnScroll, s.onScroll)
	s.Panel.Subscribe(OnResize, s.onResize)

	s.update()
	s.recalc()
}

// onCursor receives subscribed cursor events over the panel
func (s *ItemScroller) onCursor(evname string, ev interface{}) {

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
func (s *ItemScroller) onScroll(evname string, ev interface{}) {

	sev := ev.(*window.ScrollEvent)
	if sev.Yoffset > 0 {
		s.ScrollUp()
	} else if sev.Yoffset < 0 {
		s.ScrollDown()
	}
	s.root.StopPropagation(Stop3D)
}

// onResize receives resize events
func (s *ItemScroller) onResize(evname string, ev interface{}) {

	s.recalc()
}

// autoSize resizes the scroller if necessary
func (s *ItemScroller) autoSize() {

	if s.maxAutoWidth == 0 && s.maxAutoHeight == 0 {
		return
	}

	var width float32
	var height float32
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
func (s *ItemScroller) recalc() {

	if s.vert {
		s.vRecalc()
	} else {
		s.hRecalc()
	}
}

// vRecalc recalculates for the vertical scroller
func (s *ItemScroller) vRecalc() {

	// Checks if scroll bar should be visible or not
	scroll := false
	if s.first > 0 {
		scroll = true
	} else {
		var posY float32
		for _, item := range s.items[s.first:] {
			posY += item.TotalHeight()
			if posY > s.height {
				scroll = true
				break
			}
		}
	}
	s.setVScrollBar(scroll)

	// Compute size of scroll button
	if scroll && s.autoButtonSize {
		var totalHeight float32
		for _, item := range s.items {
			// TODO OPTIMIZATION
			// Break when the view/content proportion becomes smaller than the minimum button size
			totalHeight += item.TotalHeight()
		}
		s.vscroll.SetButtonSize(s.height * s.height / totalHeight)
	}

	// Items width
	width := s.ContentWidth()
	if scroll {
		width -= s.vscroll.Width()
	}

	var posY float32
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
func (s *ItemScroller) hRecalc() {

	// Checks if scroll bar should be visible or not
	scroll := false
	if s.first > 0 {
		scroll = true
	} else {
		var posX float32
		for _, item := range s.items[s.first:] {
			posX += item.GetPanel().Width()
			if posX > s.width {
				scroll = true
				break
			}
		}
	}
	s.setHScrollBar(scroll)

	// Compute size of scroll button
	if scroll && s.autoButtonSize {
		var totalWidth float32
		for _, item := range s.items {
			// TODO OPTIMIZATION
			// Break when the view/content proportion becomes smaller than the minimum button size
			totalWidth += item.GetPanel().Width()
		}
		s.hscroll.SetButtonSize(s.width * s.width / totalWidth)
	}

	// Items height
	height := s.ContentHeight()
	if scroll {
		height -= s.hscroll.Height()
	}

	var posX float32
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
func (s *ItemScroller) maxFirst() int {

	// Vertical scroller
	if s.vert {
		var height float32
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
	}

	// Horizontal scroller
	var width float32
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

// setVScrollBar sets the visibility state of the vertical scrollbar
func (s *ItemScroller) setVScrollBar(state bool) {

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
func (s *ItemScroller) setHScrollBar(state bool) {

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
func (s *ItemScroller) onScrollBarEvent(evname string, ev interface{}) {

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
func (s *ItemScroller) update() {

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
func (s *ItemScroller) applyStyle(st *ItemScrollerStyle) {

	s.Panel.ApplyStyle(&st.PanelStyle)
}
