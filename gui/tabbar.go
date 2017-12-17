// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import "github.com/g3n/engine/math32"

// TabBar is a panel which can contain other panels contained in Tabs.
// Only one panel is visible at one time.
// To show another panel the corresponding Tab must be selected.
type TabBar struct {
	Panel                    // Embedded panel
	styles     *TabBarStyles // Pointer to current styles
	tabs       []*Tab        // Array of tabs
	selected   int           // Index of the selected tab
	cursorOver bool          // Cursor over flag
}

// Tab describes an individual tab from the TabBar
type Tab struct {
	header  Panel  // Tab header
	label   *Label // Tab optional label
	icon    *Label // Tab optional icon
	img     *Image // Tab optional image
	content Panel  // User content panel
}

// TabBarStyle describes the style
type TabBarStyle struct {
	Border      BorderSizes
	Paddings    BorderSizes
	BorderColor math32.Color4
	BgColor     math32.Color4
	FgColor     math32.Color
}

type TabBarStyles struct {
	Normal   TabBarStyle
	Over     TabBarStyle
	Focus    TabBarStyle
	Disabled TabBarStyle
}

// NewTabBar creates and returns a pointer to a new TabBar widget
// with the specified width and height
func NewTabBar(width, height float32) *TabBar {

	// Creates new TabBar
	tb := new(TabBar)
	tb.Initialize(width, height)
	tb.styles = &StyleDefault().TabBar
	tb.tabs = make([]*Tab, 0)
	tb.selected = -1

	// Subscribe to panel events
	tb.Subscribe(OnCursorEnter, tb.onCursor)
	tb.Subscribe(OnCursorLeave, tb.onCursor)
	tb.Subscribe(OnEnable, func(name string, ev interface{}) { tb.update() })
	tb.Subscribe(OnResize, func(name string, ev interface{}) { tb.recalc() })

	tb.recalc()
	tb.update()
	return tb
}

// AddTab creates and adds a new Tab panel with the specified header text
// at end of this TabBar list of tabs.
// Returns the pointer to thew new Tab.
func (tb *TabBar) AddTab(text string) *Tab {

	return tb.InsertTab(text, len(tb.tabs))
}

// InsertTab creates and inserts a new Tab panel with the specified header text
// at the specified position in the TabBar from left to right.
// Returns the pointer to the new Tab or nil if the position is invalid.
func (tb *TabBar) InsertTab(text string, pos int) *Tab {

	// Checks position to insert into
	if pos < 0 || pos > len(tb.tabs) {
		return nil
	}

	// Creates and initializes the new Tab
	tab := new(Tab)
	if text != "" {
		tab.label = NewLabel(text)
	}
	tab.content.Initialize(0, 0)

	// Inserts created Tab at the specified position
	tb.tabs = append(tb.tabs, nil)
	copy(tb.tabs[pos+1:], tb.tabs[pos:])
	tb.tabs[pos] = tab

	tb.update()
	tb.recalc()
	return tab
}

// RemoveTab removes the tab at the specified position in the TabBar.
// Returns the pointer of the removed tab or nil if the position is invalid.
func (tb *TabBar) RemoveTab(pos int) *Tab {

	if pos < 0 || pos >= len(tb.tabs) {
		return nil
	}
	tab := tb.tabs[pos]
	// Remove tab from array
	copy(tb.tabs[pos:], tb.tabs[pos+1:])
	tb.tabs[len(tb.tabs)-1] = nil
	tb.tabs = tb.tabs[:len(tb.tabs)-1]
	// Checks if tab was selected
	if tb.selected == pos {

	}
	return tab
}

// TabCount returns the current number of tabs
func (tb *TabBar) TabCount() int {

	return len(tb.tabs)
}

// TabAt returns the pointer of the Tab object at the specified index.
// Return nil if the index is invalid
func (tb *TabBar) TabAt(idx int) *Tab {

	if idx < 0 || idx >= len(tb.tabs) {
		return nil
	}
	return tb.tabs[idx]
}

// SetSelected sets the selected tab of the TabBar to the tab with the specified position.
// Returns the pointer of the selected tab or nil if the position is invalid.
func (tb *TabBar) SetSelected(pos int) *Tab {

	if pos < 0 || pos >= len(tb.tabs) {
		return nil
	}

	tb.selected = pos
	return tb.tabs[pos]
}

// onCursor process subscribed cursor events
func (tb *TabBar) onCursor(evname string, ev interface{}) {

	switch evname {
	case OnCursorEnter:
		tb.cursorOver = true
		tb.update()
	case OnCursorLeave:
		tb.cursorOver = false
		tb.update()
	default:
		return
	}
	tb.root.StopPropagation(StopAll)
}

// applyStyle applies the specified TabBar style
func (tb *TabBar) applyStyle(s *TabBarStyle) {

	tb.SetBordersFrom(&s.Border)
	tb.SetBordersColor4(&s.BorderColor)
	tb.SetPaddingsFrom(&s.Paddings)
	tb.SetColor4(&s.BgColor)
}

// recalc recalculates and updates the positions of all tabs
func (tb *TabBar) recalc() {

}

// update...
func (tb *TabBar) update() {

	if !tb.Enabled() {
		tb.applyStyle(&tb.styles.Disabled)
		return
	}
	if tb.cursorOver {
		tb.applyStyle(&tb.styles.Over)
		return
	}
	tb.applyStyle(&tb.styles.Normal)
}

//
// Tab methods
//

// SetText sets the text of the tab header
func (tab *Tab) SetText(text string) *Tab {

	return tab
}

// SetIcon sets the icon of the tab header
func (tab *Tab) SetIcon(icon string) *Tab {

	return tab
}

// Panel returns a pointer to the specified tab content panel
func (tab *Tab) Panel() *Panel {

	return &tab.content
}

// recalc recalculates the positions of the Tab header internal panels
func (tab *Tab) recalc() {

	width := tab.header.Width()

}
