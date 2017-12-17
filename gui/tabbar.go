// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import "github.com/g3n/engine/math32"

// TabBar is a panel which can contain other panels arranged in horizontal Tabs.
// Only one panel is visible at a time.
// To show another panel the corresponding Tab must be selected.
type TabBar struct {
	Panel                    // Embedded panel
	styles     *TabBarStyles // Pointer to current styles
	tabs       []*Tab        // Array of tabs
	selected   int           // Index of the selected tab
	cursorOver bool          // Cursor over flag
}

// TabBarStyle describes the style of the TabBar
type TabBarStyle struct {
	Border      BorderSizes   // Border sizes
	Paddings    BorderSizes   // Padding sizes
	BorderColor math32.Color4 // Border color
	BgColor     math32.Color4 // Background color
}

// TabBarStyles describes all the TabBarStyles
type TabBarStyles struct {
	Normal   TabBarStyle // Style for normal exhibition
	Over     TabBarStyle // Style when cursor is over the TabBar
	Focus    TabBarStyle // Style when the TabBar has key focus
	Disabled TabBarStyle // Style when the TabBar is disabled
	Tab      TabStyles   // Style for Tabs
}

// TabStyle describes the style of the individual Tabs
type TabStyle struct {
	Border      BorderSizes
	Paddings    BorderSizes
	BorderColor math32.Color4
	BgColor     math32.Color4
	FgColor     math32.Color
}

// TabStyles describes all Tab styles
type TabStyles struct {
	MinWidth float32  // Minimum Tab header width
	Normal   TabStyle // Style for normal exhibition
	Over     TabStyle // Style when cursor is over the Tab
	Focus    TabStyle // Style when the Tab has key focus
	Disabled TabStyle // Style when the Tab is disabled
	Selected TabStyle // Style when the Tab is selected
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

	// Inserts created Tab at the specified position
	tab := newTab(text, &tb.styles.Tab)
	tb.tabs = append(tb.tabs, nil)
	copy(tb.tabs[pos+1:], tb.tabs[pos:])
	tb.tabs[pos] = tab
	tb.Add(&tab.header)
	tb.Add(&tab.content)

	tb.update()
	tb.recalc()
	return tab
}

// RemoveTab removes the tab at the specified position in the TabBar.
// Returns the pointer of the removed tab or nil if the position is invalid.
func (tb *TabBar) RemoveTab(pos int) *Tab {

	// Check position to remove from
	if pos < 0 || pos >= len(tb.tabs) {
		return nil
	}

	// Remove tab from array
	tab := tb.tabs[pos]
	copy(tb.tabs[pos:], tb.tabs[pos+1:])
	tb.tabs[len(tb.tabs)-1] = nil
	tb.tabs = tb.tabs[:len(tb.tabs)-1]

	// Checks if removed tab was selected
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

// Selected returns the position of the selected Tab.
// Returns value < 0 if there is no selected Tab.
func (tb *TabBar) Selected() int {

	return tb.selected
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

	maxWidth := tb.ContentWidth() / float32(len(tb.tabs))
	headerx := float32(0)
	for i := 0; i < len(tb.tabs); i++ {
		tab := tb.tabs[i]
		tab.recalc(maxWidth)
		tab.header.SetPosition(headerx, 0)
		// Sets size and position of the Tab content panel
		contentx := float32(0)
		contenty := tab.header.Height()
		tab.content.SetWidth(tb.ContentWidth())
		tab.content.SetHeight(tb.ContentHeight() - tab.header.Height())
		tab.content.SetPosition(contentx, contenty)
		headerx += tab.header.Width()
	}
}

// update updates the TabBar visual state
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
// Tab
//

// Tab describes an individual tab of the TabBar
type Tab struct {
	styles     *TabStyles // Pointer to Tab current styles
	header     Panel      // Tab header
	label      *Label     // Tab label
	icon       *Label     // Tab optional icon
	img        *Image     // Tab optional image
	content    Panel      // User content panel
	cursorOver bool
}

// newTab creates and returns a pointer to a new Tab
func newTab(text string, styles *TabStyles) *Tab {

	tab := new(Tab)
	tab.styles = styles
	tab.header.Initialize(0, 0)
	tab.label = NewLabel(text)
	tab.header.Add(tab.label)
	tab.content.Initialize(0, 0)

	// Subscribe to header events
	tab.header.Subscribe(OnCursorEnter, tab.onCursor)
	tab.header.Subscribe(OnCursorLeave, tab.onCursor)

	tab.update()
	return tab
}

// onCursor process subscribed cursor events
func (tab *Tab) onCursor(evname string, ev interface{}) {

	switch evname {
	case OnCursorEnter:
		tab.cursorOver = true
		tab.update()
	case OnCursorLeave:
		tab.cursorOver = false
		tab.update()
	default:
		return
	}
	tab.header.root.StopPropagation(StopAll)
}

// SetText sets the text of the tab header
func (tab *Tab) SetText(text string) *Tab {

	return tab
}

// SetIcon sets the icon of the tab header
func (tab *Tab) SetIcon(icon string) *Tab {

	return tab
}

// Content returns a pointer to the specified Tab content panel
func (tab *Tab) Content() *Panel {

	return &tab.content
}

// applyStyle applies the specified Tab style to the Tab header
func (tab *Tab) applyStyle(s *TabStyle) {

	tab.header.SetBordersFrom(&s.Border)
	tab.header.SetBordersColor4(&s.BorderColor)
	tab.header.SetPaddingsFrom(&s.Paddings)
	tab.header.SetColor4(&s.BgColor)
}

func (tab *Tab) update() {

	if !tab.header.Enabled() {
		tab.applyStyle(&tab.styles.Disabled)
		return
	}
	if tab.cursorOver {
		tab.applyStyle(&tab.styles.Over)
		return
	}
	tab.applyStyle(&tab.styles.Normal)
}

// recalc recalculates the size of the Tab header and the size
// and positions of the Taheader internal panels
func (tab *Tab) recalc(maxWidth float32) {

	height := tab.label.Height()
	//iconWidth := float32(0)
	//if tab.icon != nil {
	//	tab.icon.SetPosition(0, 0)
	//	iconWidth = tab.icon.Width()
	//} else if tab.img != nil {
	//	tab.img.SetPosition(0, 0)
	//	iconWidth = tab.img.Width()
	//}

	width := tab.label.Width()
	tab.header.SetContentSize(width, height)

}
