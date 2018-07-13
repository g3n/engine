// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/window"
)

// Scroller is the GUI element that allows scrolling of a target IPanel.
// A scroller can have up to two scrollbars, one vertical and one horizontal.
// The vertical scrollbar, if any, can be located either on the left or on the right.
// The horizontal scrollbar, if any, can be located either on the top or on the bottom.
// The interlocking of the scrollbars (which happens when both scrollbars are visible) can be configured.
// Whether each scrollbar overlaps the content can also be configured (useful for transparent UIs).
type Scroller struct {
	Panel                        // Embedded panel
	mode          ScrollMode     // ScrollMode specifies which scroll directions are allowed
	target        IPanel         // The IPanel that will be scrolled through
	hscroll       *ScrollBar     // Horizontal scrollbar (may be nil)
	vscroll       *ScrollBar     // Vertical scrollbar (may be nil)
	style         *ScrollerStyle // The current style
	corner        *Panel         // The optional corner panel (can be visible when scrollMode==Both, interlocking==None, corner=true)
	cursorOver    bool           // Cursor is over the scroller
	modKeyPressed bool           // Modifier key is pressed
}

// ScrollMode specifies which scroll directions are allowed.
type ScrollMode int

// The various scroll modes.
const (
	ScrollNone       = ScrollMode(0x00)                              // No scrolling allowed
	ScrollVertical   = ScrollMode(0x01)                              // Vertical scrolling allowed
	ScrollHorizontal = ScrollMode(0x02)                              // Horizontal scrolling allowed
	ScrollBoth       = ScrollMode(ScrollVertical | ScrollHorizontal) // Both vertical and horizontal scrolling allowed
)

// ScrollbarInterlocking specifies what happens where the vertical and horizontal scrollbars meet.
type ScrollbarInterlocking int

// The three scrollbar interlocking types.
const (
	ScrollbarInterlockingNone       = ScrollbarInterlocking(iota) // No scrollbar interlocking
	ScrollbarInterlockingVertical                                 // Vertical scrollbar takes precedence
	ScrollbarInterlockingHorizontal                               // Horizontal scrollbar takes precedence
)

// ScrollbarPosition specifies where the scrollbar is located.
// For the vertical scrollbar it specifies whether it's added to the left or to the right.
// For the horizontal scrollbar it specifies whether it's added to the top or to the bottom.
type ScrollbarPosition int

// The four possible scrollbar positions.
const (
	ScrollbarLeft   = ScrollbarPosition(iota) // Scrollbar is positioned on the left of the scroller
	ScrollbarRight                            // Scrollbar is positioned on the right of the scroller
	ScrollbarTop                              // Scrollbar is positioned on the top of the scroller
	ScrollbarBottom                           // Scrollbar is positioned on the bottom of the scroller
)

// ScrollerStyle contains the styling of a Scroller
type ScrollerStyle struct {
	PanelStyle                                   // Embedded PanelStyle
	VerticalScrollbar     ScrollerScrollbarStyle // The style of the vertical scrollbar
	HorizontalScrollbar   ScrollerScrollbarStyle // The style of the horizontal scrollbar
	CornerPanel           PanelStyle             // The style of the corner panel
	ScrollbarInterlocking ScrollbarInterlocking  // Specifies what happens where the vertical and horizontal scrollbars meet
	CornerCovered         bool                   // True indicates that the corner panel should be visible when appropriate
}

// ScrollerScrollbarStyle is the set of style options for a scrollbar that is part of a scroller.
type ScrollerScrollbarStyle struct {
	ScrollBarStyle                   // Embedded ScrollBarStyle (TODO, should be ScrollBarStyle*S*, implement style logic)
	Position       ScrollbarPosition // Specifies the positioning of the scrollbar
	Broadness      float32           // Broadness of the scrollbar
	OverlapContent bool              // Specifies whether the scrollbar is shown above the content area
	AutoSizeButton bool              // Specifies whether the scrollbar button size is adjusted based on content/view proportion
}

// TODO these configuration variables could be made part of a global engine configuration object in the future
// They should not be added to style since they are not style changes and not to the struct since they are global

// ScrollModifierKey is the Key that changes the scrolling direction from vertical to horizontal
const ScrollModifierKey = window.KeyLeftShift

// NewScroller creates and returns a pointer to a new Scroller with the specified
// target IPanel and ScrollMode.
func NewScroller(width, height float32, mode ScrollMode, target IPanel) *Scroller {

	s := new(Scroller)
	s.initialize(width, height, mode, target)
	return s
}

// initialize initializes this scroller and can be called by other types which embed a scroller
func (s *Scroller) initialize(width, height float32, mode ScrollMode, target IPanel) {

	s.Panel.Initialize(width, height)
	s.style = &StyleDefault().Scroller
	s.target = target
	s.Panel.Add(s.target)
	s.mode = mode

	s.Subscribe(OnCursorEnter, s.onCursor)
	s.Subscribe(OnCursorLeave, s.onCursor)
	s.Subscribe(OnScroll, s.onScroll)
	s.Subscribe(OnKeyDown, s.onKey)
	s.Subscribe(OnKeyUp, s.onKey)
	s.Subscribe(OnResize, s.onResize)

	s.Update()
}

// SetScrollMode sets the scroll mode
func (s *Scroller) SetScrollMode(mode ScrollMode) {

	s.mode = mode
	s.Update()
}

// ScrollMode returns the current scroll mode
func (s *Scroller) ScrollMode() ScrollMode {

	return s.mode
}

// SetScrollbarInterlocking sets the scrollbar interlocking mode
func (s *Scroller) SetScrollbarInterlocking(interlocking ScrollbarInterlocking) {

	s.style.ScrollbarInterlocking = interlocking
	s.Update()
}

// ScrollbarInterlocking returns the current scrollbar interlocking mode
func (s *Scroller) ScrollbarInterlocking() ScrollbarInterlocking {

	return s.style.ScrollbarInterlocking
}

// SetCornerCovered specifies whether the corner covering panel is shown when appropriate
func (s *Scroller) SetCornerCovered(state bool) {

	s.style.CornerCovered = state
	s.Update()
}

// CornerCovered returns whether the corner covering panel is being shown when appropriate
func (s *Scroller) CornerCovered() bool {

	return s.style.CornerCovered
}

// SetVerticalScrollbarPosition sets the position of the vertical scrollbar (i.e. left or right)
func (s *Scroller) SetVerticalScrollbarPosition(pos ScrollbarPosition) {

	s.style.VerticalScrollbar.Position = pos
	s.recalc()
}

// VerticalScrollbarPosition returns the current position of the vertical scrollbar (i.e. left or right)
func (s *Scroller) VerticalScrollbarPosition() ScrollbarPosition {

	return s.style.VerticalScrollbar.Position
}

// SetHorizontalScrollbarPosition sets the position of the horizontal scrollbar (i.e. top or bottom)
func (s *Scroller) SetHorizontalScrollbarPosition(pos ScrollbarPosition) {

	s.style.HorizontalScrollbar.Position = pos
	s.recalc()
}

// HorizontalScrollbarPosition returns the current position of the horizontal scrollbar (i.e. top or bottom)
func (s *Scroller) HorizontalScrollbarPosition() ScrollbarPosition {

	return s.style.HorizontalScrollbar.Position
}

// SetVerticalScrollbarOverlapping specifies whether the vertical scrollbar overlaps the content area
func (s *Scroller) SetVerticalScrollbarOverlapping(state bool) {

	s.style.VerticalScrollbar.OverlapContent = state
	s.Update()
}

// VerticalScrollbarOverlapping returns whether the vertical scrollbar overlaps the content area
func (s *Scroller) VerticalScrollbarOverlapping() bool {

	return s.style.VerticalScrollbar.OverlapContent
}

// SetHorizontalScrollbarOverlapping specifies whether the horizontal scrollbar overlaps the content area
func (s *Scroller) SetHorizontalScrollbarOverlapping(state bool) {

	s.style.HorizontalScrollbar.OverlapContent = state
	s.Update()
}

// HorizontalScrollbarOverlapping returns whether the horizontal scrollbar overlaps the content area
func (s *Scroller) HorizontalScrollbarOverlapping() bool {

	return s.style.HorizontalScrollbar.OverlapContent
}

// SetVerticalScrollbarAutoSizeButton specifies whether the vertical scrollbar button is sized automatically
func (s *Scroller) SetVerticalScrollbarAutoSizeButton(state bool) {

	s.style.VerticalScrollbar.AutoSizeButton = state
	if s.vscroll != nil {
		if state == false {
			s.vscroll.SetButtonSize(s.style.VerticalScrollbar.ScrollBarStyle.ButtonLength)
		}
		s.recalc()
	}
}

// VerticalScrollbarAutoSizeButton returns whether the vertical scrollbar button is sized automatically
func (s *Scroller) VerticalScrollbarAutoSizeButton() bool {

	return s.style.VerticalScrollbar.AutoSizeButton
}

// SetHorizontalScrollbarAutoSizeButton specifies whether the horizontal scrollbar button is sized automatically
func (s *Scroller) SetHorizontalScrollbarAutoSizeButton(state bool) {

	s.style.HorizontalScrollbar.AutoSizeButton = state
	if s.hscroll != nil {
		if state == false {
			s.hscroll.SetButtonSize(s.style.HorizontalScrollbar.ScrollBarStyle.ButtonLength)
		}
		s.recalc()
	}
}

// HorizontalScrollbarAutoSizeButton returns whether the horizontal scrollbar button is sized automatically
func (s *Scroller) HorizontalScrollbarAutoSizeButton() bool {

	return s.style.HorizontalScrollbar.AutoSizeButton
}

// SetVerticalScrollbarBroadness sets the broadness of the vertical scrollbar
func (s *Scroller) SetVerticalScrollbarBroadness(broadness float32) {

	s.style.VerticalScrollbar.Broadness = broadness
	if s.vscroll != nil {
		s.vscroll.SetWidth(broadness)
		s.Update()
	}
}

// VerticalScrollbarBroadness returns the broadness of the vertical scrollbar
func (s *Scroller) VerticalScrollbarBroadness() float32 {

	return s.style.VerticalScrollbar.Broadness
}

// SetHorizontalScrollbarBroadness sets the broadness of the horizontal scrollbar
func (s *Scroller) SetHorizontalScrollbarBroadness(broadness float32) {

	s.style.HorizontalScrollbar.Broadness = broadness
	if s.hscroll != nil {
		s.hscroll.SetHeight(broadness)
		s.Update()
	}
}

// HorizontalScrollbarBroadness returns the broadness of the horizontal scrollbar
func (s *Scroller) HorizontalScrollbarBroadness() float32 {

	return s.style.HorizontalScrollbar.Broadness
}

// ScrollTo scrolls the target panel such that the specified target point is centered on the scroller's view area
func (s *Scroller) ScrollTo(x, y float32) {
	// TODO
}

// onCursor receives subscribed cursor events over the panel
func (s *Scroller) onCursor(evname string, ev interface{}) {

	switch evname {
	case OnCursorEnter:
		s.root.SetScrollFocus(s)
		s.root.SetKeyFocus(s)
		s.cursorOver = true
	case OnCursorLeave:
		s.root.SetScrollFocus(nil)
		s.root.SetKeyFocus(nil)
		s.cursorOver = false
	}
	s.root.StopPropagation(Stop3D)
}

// onScroll receives mouse scroll events when this scroller has the scroll focus (set by OnMouseEnter)
func (s *Scroller) onScroll(evname string, ev interface{}) {

	sev := ev.(*window.ScrollEvent)

	vScrollVisible := (s.vscroll != nil) && s.vscroll.Visible()
	hScrollVisible := (s.hscroll != nil) && s.hscroll.Visible()

	mult := float32(1) / float32(10)
	offsetX := sev.Xoffset * mult
	offsetY := sev.Yoffset * mult

	// If modifier key is pressed (left shift by default) - then scroll in the horizontal direction
	if s.modKeyPressed {
		if math32.Abs(offsetY) > math32.Abs(offsetX) {
			offsetX = offsetY
		}
		offsetY = 0
	}

	log.Error("X: %v, Y: %v", offsetX, offsetY)

	if vScrollVisible {
		if hScrollVisible {
			// Both scrollbars are present - scroll both
			s.vscroll.SetValue(float32(s.vscroll.Value()) - offsetY)
			s.hscroll.SetValue(float32(s.hscroll.Value()) - offsetX)
		} else {
			// Only vertical scrollbar present - scroll it
			s.vscroll.SetValue(float32(s.vscroll.Value()) - offsetY)
		}
	} else if hScrollVisible {
		// Only horizontal scrollbar present - scroll it
		s.hscroll.SetValue(float32(s.hscroll.Value()) - offsetX)
	}

	s.recalc()
	s.root.StopPropagation(Stop3D)
}

// onKey receives key events
func (s *Scroller) onKey(evname string, ev interface{}) {

	key := ev.(*window.KeyEvent)
	log.Error("Key %v", key)
	if key.Keycode == ScrollModifierKey {
		if evname == OnKeyDown {
			s.modKeyPressed = true
			log.Error("true")
		} else if evname == OnKeyUp {
			log.Error("false")
			s.modKeyPressed = false
		}
	}
	s.root.StopPropagation(Stop3D)
}

// onResize receives resize events
func (s *Scroller) onResize(evname string, ev interface{}) {

	s.Update()
}

// setVerticalScrollbarVisible sets the vertical scrollbar visible, creating and initializing it if it's the first time
func (s *Scroller) setVerticalScrollbarVisible() {
	if s.vscroll == nil {
		s.vscroll = NewVScrollBar(s.style.VerticalScrollbar.Broadness, 0)
		s.vscroll.applyStyle(&s.style.VerticalScrollbar.ScrollBarStyle)
		s.vscroll.Subscribe(OnChange, s.onScrollBarEvent)
		s.Add(s.vscroll)
	}
	s.vscroll.SetVisible(true)
}

// setVerticalScrollbarVisible sets the horizontal scrollbar visible, creating and initializing it if it's the first time
func (s *Scroller) setHorizontalScrollbarVisible() {
	if s.hscroll == nil {
		s.hscroll = NewHScrollBar(0, s.style.HorizontalScrollbar.Broadness)
		s.hscroll.applyStyle(&s.style.HorizontalScrollbar.ScrollBarStyle)
		s.hscroll.Subscribe(OnChange, s.onScrollBarEvent)
		s.Add(s.hscroll)
	}
	s.hscroll.SetVisible(true)
}

// updateScrollbarsVisibility updates the visibility of the scrollbars and corner panel, creating them if necessary.
// This method should be called when either the target panel changes size or when either the scroll mode or
// style of the Scroller changes.
func (s *Scroller) updateScrollbarsVisibility() {

	// Obtain the size of the target panel
	targetWidth := s.target.TotalWidth()
	targetHeight := s.target.TotalHeight()

	// If vertical scrolling is enabled and the vertical scrollbar should be visible
	if (s.mode&ScrollVertical > 0) && (targetHeight > s.content.Height) {
		s.setVerticalScrollbarVisible()
	} else if s.vscroll != nil {
		s.vscroll.SetVisible(false)
		s.vscroll.SetValue(0)
	}

	// If horizontal scrolling is enabled and the horizontal scrollbar should be visible
	if (s.mode&ScrollHorizontal > 0) && (targetWidth > s.content.Width) {
		s.setHorizontalScrollbarVisible()
	} else if s.hscroll != nil {
		s.hscroll.SetVisible(false)
		s.hscroll.SetValue(0)
	}

	// If both scrollbars can be visible we need to check whether we should show the corner panel and also whether
	// any scrollbar's presence caused the other to be required. The latter is a literal and figurative edge case
	// that happens when the target panel is larger than the content in one dimension but smaller than the content
	// in the other, and in the dimension that it is smaller than the content, the difference is less than the width
	// of the scrollbar. In that case we need to show both scrollbars to allow viewing of the complete target panel.
	if s.mode == ScrollBoth {

		vScrollVisible := (s.vscroll != nil) && s.vscroll.Visible()
		hScrollVisible := (s.hscroll != nil) && s.hscroll.Visible()

		// Check if adding any of the scrollbars ended up covering an edge of the target. If that's the case,
		// then show the other scrollbar as well (if the covering scrollbar's style is set to non-overlapping).

		// If the vertical scrollbar is visible and covering the target (and its style is not set to overlap)
		if vScrollVisible && (targetWidth > (s.content.Width - s.vscroll.width)) && !s.style.VerticalScrollbar.OverlapContent {
			s.setHorizontalScrollbarVisible() // Show the other scrollbar too
		}
		// If the horizontal scrollbar is visible and covering the target (and its style is not set to overlap)
		if hScrollVisible && (targetHeight > (s.content.Height - s.hscroll.height)) && !s.style.HorizontalScrollbar.OverlapContent {
			s.setVerticalScrollbarVisible() // Show the other scrollbar too
		}

		// Update visibility variables since they may have changed
		vScrollVisible = (s.vscroll != nil) && s.vscroll.Visible()
		hScrollVisible = (s.hscroll != nil) && s.hscroll.Visible()

		// If both vertical and horizontal scrolling is enabled, and the style specifies no interlocking
		// and a corner panel, and both scrollbars are visible - then the corner panel should be visible
		if (s.style.ScrollbarInterlocking == ScrollbarInterlockingNone) && s.style.CornerCovered && vScrollVisible && hScrollVisible {
			if s.corner == nil {
				s.corner = NewPanel(s.vscroll.width, s.hscroll.height)
				s.corner.ApplyStyle(&s.style.CornerPanel)
				s.Add(s.corner)
			}
			s.corner.SetVisible(true)
		} else if s.corner != nil {
			s.corner.SetVisible(false)
		}
	}
}

// onScrollEvent is called when the scrollbar value changes
func (s *Scroller) onScrollBarEvent(evname string, ev interface{}) {

	s.recalc()
}

// recalc recalculates the positions and sizes of the scrollbars and corner panel,
// updates the size of the scrollbar buttons, and repositions the target panel
func (s *Scroller) recalc() {

	// The multipliers of the scrollbars' [0,1] values.
	// After applied, they will give the correct target panel position.
	// They can be thought of as the range of motion of the target panel in each axis
	multHeight := s.target.TotalHeight() - s.content.Height
	multWidth := s.target.TotalWidth() - s.content.Width

	var targetX, targetY float32
	var offsetX, offsetY float32

	vScrollVisible := (s.mode&ScrollVertical > 0) && (s.vscroll != nil) && s.vscroll.Visible()
	hScrollVisible := (s.mode&ScrollHorizontal > 0) && (s.hscroll != nil) && s.hscroll.Visible()

	// If the vertical scrollbar is visible
	if vScrollVisible {
		s.recalcV() // Recalculate scrollbar size/position (checks for the other scrollbar's presence)
		targetY = -float32(s.vscroll.Value())
		// If we don't want it to overlap the content area
		if s.style.VerticalScrollbar.OverlapContent == false {
			// Increase the target's range of X motion by the width of the vertical scrollbar
			multWidth += s.vscroll.width
			// If the vertical scrollbar is on the left we also want to add an offset to the target panel
			if s.style.VerticalScrollbar.Position == ScrollbarLeft {
				offsetX += s.vscroll.width
			}
		}
	}

	// If the horizontal scrollbar is visible
	if hScrollVisible {
		s.recalcH() // Recalculate scrollbar size/position (checks for the other scrollbar's presence)
		targetX = -float32(s.hscroll.Value())
		// If we don't want it to overlap the content area
		if s.style.HorizontalScrollbar.OverlapContent == false {
			// Increase the target's range of Y motion by the height of the horizontal scrollbar
			multHeight += s.hscroll.height
			// If the horizontal scrollbar is on the top we also want to add an offset to the target panel
			if s.style.HorizontalScrollbar.Position == ScrollbarTop {
				offsetY += s.hscroll.height
			}
		}
	}

	// Reposition the target panel
	s.target.SetPosition(targetX*multWidth+offsetX, targetY*multHeight+offsetY)

	// If the corner panel should be visible, update its position and size
	if (s.mode == ScrollBoth) && (s.style.ScrollbarInterlocking == ScrollbarInterlockingNone) &&
		(s.style.CornerCovered == true) && vScrollVisible && hScrollVisible {
		s.corner.SetPosition(s.vscroll.Position().X, s.hscroll.Position().Y)
		s.corner.SetSize(s.vscroll.width, s.hscroll.height)
	}

}

// recalcV recalculates the size and position of the vertical scrollbar
func (s *Scroller) recalcV() {

	// Position the vertical scrollbar horizontally according to the style
	var vscrollPosX float32 // = 0 (ScrollbarLeft)
	if s.style.VerticalScrollbar.Position == ScrollbarRight {
		vscrollPosX = s.ContentWidth() - s.vscroll.width
	}

	// Start with the default Y position and height of the vertical scrollbar
	var vscrollPosY float32
	vscrollHeight := s.ContentHeight()
	viewHeight := s.ContentHeight()

	// If the horizontal scrollbar is present - reduce the viewHeight ...
	if (s.hscroll != nil) && s.hscroll.Visible() {
		if s.style.HorizontalScrollbar.OverlapContent == false {
			viewHeight -= s.hscroll.height
		}
		// If the interlocking style doesn't give precedence to the vertical scrollbar - reduce the scrollbar height ...
		if s.style.ScrollbarInterlocking != ScrollbarInterlockingVertical {
			vscrollHeight -= s.hscroll.height
			// If the horizontal scrollbar is on top - offset the vertical scrollbar vertically
			if s.style.HorizontalScrollbar.Position == ScrollbarTop {
				vscrollPosY = s.hscroll.height
			}
		}
	}

	// Adjust the scrollbar button size to the correct proportion proportion according to the style
	if s.style.VerticalScrollbar.AutoSizeButton {
		s.vscroll.SetButtonSize(vscrollHeight * viewHeight / s.target.TotalHeight())
	}

	// Update the position and height of the vertical scrollbar
	s.vscroll.SetPosition(vscrollPosX, vscrollPosY)
	s.vscroll.SetHeight(vscrollHeight)
}

// recalcH recalculates the size and position of the horizontal scrollbar
func (s *Scroller) recalcH() {

	// Position the horizontal scrollbar vertically according to the style
	var hscrollPosY float32 // = 0 (ScrollbarTop)
	if s.style.HorizontalScrollbar.Position == ScrollbarBottom {
		hscrollPosY = s.ContentHeight() - s.hscroll.height
	}

	// Start with default X position and width of the horizontal scrollbar
	var hscrollPosX float32
	hscrollWidth := s.ContentWidth()
	viewWidth := s.ContentWidth()

	// If the vertical scrollbar is present - reduce the viewWidth ...
	if (s.vscroll != nil) && s.vscroll.Visible() {
		if s.style.VerticalScrollbar.OverlapContent == false {
			viewWidth -= s.vscroll.width
		}
		// If the interlocking style doesn't give precedence to the horizontal scrollbar - reduce the scrollbar width ...
		if s.style.ScrollbarInterlocking != ScrollbarInterlockingHorizontal {
			hscrollWidth -= s.vscroll.width
			// If the vertical scrollbar is on the left - offset the horizontal scrollbar horizontally
			if s.style.VerticalScrollbar.Position == ScrollbarLeft {
				hscrollPosX = s.vscroll.width
			}
		}
	}

	// Adjust the scrollbar button size to the correct proportion proportion according to the style
	if s.style.HorizontalScrollbar.AutoSizeButton {
		s.hscroll.SetButtonSize(hscrollWidth * viewWidth / s.target.TotalWidth())
	}

	// Update the position and width of the horizontal scrollbar
	s.hscroll.SetPosition(hscrollPosX, hscrollPosY)
	s.hscroll.SetWidth(hscrollWidth)
}

// Update updates the visibility of the scrollbars, corner panel, and then recalculates
func (s *Scroller) Update() {

	s.updateScrollbarsVisibility()
	s.recalc()
}

// TODO - if the style is changed this needs to be called to update the scrollbars and corner panel
func (s *Scroller) applyStyle(ss *ScrollerStyle) {

	s.style = ss

	s.vscroll.applyStyle(&s.style.VerticalScrollbar.ScrollBarStyle)
	s.hscroll.applyStyle(&s.style.HorizontalScrollbar.ScrollBarStyle)
	s.corner.ApplyStyle(&s.style.CornerPanel)

	s.Update()
}
