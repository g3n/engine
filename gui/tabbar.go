// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

// TabBar is a panel which can contain other panels organized
// as horizontal tabs.
type TabBar struct {
	Panel           // Embedded panel
	tabs     []*Tab // Array of tabs
	selected int    // Index of the selected tab
}

// Tab describes and individual tab from the TabBar
type Tab struct {
	header  Panel  // Tab header
	label   *Label // Tab optional label
	icon    *Label // Tab optional icon
	img     *Image // Tab optional image
	content Panel  // User content panel
}

// NewTabBar creates and returns a pointer to a new TabBar widget
// with the specified width and height
func NewTabBar(width, height float32) *TabBar {

	return nil
}

// AddTab creates and adds a new tab with the specified text
// at end of this TabBar list of tabs.
func (tb *TabBar) AddTab(text string) *Tab {

	return tb.InsertTab(text, len(tb.tabs))
}

// InsertTab creates and inserts a new tab at the specified position
// from left to right.
// Return nil if the position is invalid
func (tb *TabBar) InsertTab(text string, pos int) *Tab {

	if pos < 0 || pos > len(tb.tabs) {
		return nil
	}
	tab := new(Tab)

	tb.tabs = append(tb.tabs, nil)
	copy(tb.tabs[pos+1:], tb.tabs[pos:])
	tb.tabs[pos] = tab
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

// recalc recalculates and updates the positions of all tabs
func (tb *TabBar) recalc() {

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
