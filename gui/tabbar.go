// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"fmt"

	"github.com/g3n/engine/math32"
)

// TabBar is a panel which can contain other panels arranged in horizontal Tabs.
// Only one panel is visible at a time.
// To show another panel the corresponding Tab must be selected.
type TabBar struct {
	Panel                    // Embedded panel
	styles     *TabBarStyles // Pointer to current styles
	tabs       []*Tab        // Array of tabs
	separator  Panel         // Separator Panel
	iconList   *Label        // Icon for tab list button
	list       *List         // List for not visible tabs
	selected   int           // Index of the selected tab
	cursorOver bool          // Cursor over TabBar panel flag
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
	SepHeight    float32 // Separator width
	IconList     string  // Icon for showing tab list
	IconPaddings BorderSizes
	Normal       TabBarStyle // Style for normal exhibition
	Over         TabBarStyle // Style when cursor is over the TabBar
	Focus        TabBarStyle // Style when the TabBar has key focus
	Disabled     TabBarStyle // Style when the TabBar is disabled
	Tab          TabStyles   // Style for Tabs
}

// TabStyle describes the style of the individual Tabs
type TabStyle struct {
	Margins     BorderSizes
	Border      BorderSizes
	Paddings    BorderSizes
	BorderColor math32.Color4
	BgColor     math32.Color4
	FgColor     math32.Color
}

// TabStyles describes all Tab styles
type TabStyles struct {
	IconClose string   // Codepoint for close icon in Tab header
	Normal    TabStyle // Style for normal exhibition
	Over      TabStyle // Style when cursor is over the Tab
	Focus     TabStyle // Style when the Tab has key focus
	Disabled  TabStyle // Style when the Tab is disabled
	Selected  TabStyle // Style when the Tab is selected
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

	// Creates separator panel
	tb.separator.Initialize(0, 0)
	tb.Add(&tb.separator)

	// Create list
	tb.list = NewVList(0, 0)
	tb.list.Subscribe(OnMouseOut, func(evname string, ev interface{}) {
		tb.list.SetVisible(false)
	})
	tb.list.Subscribe(OnChange, tb.onListChange)
	tb.Add(tb.list)

	// Creates list icon button
	tb.iconList = NewLabel(tb.styles.IconList, true)
	tb.iconList.SetPaddingsFrom(&tb.styles.IconPaddings)
	tb.iconList.Subscribe(OnMouseDown, tb.onListButton)
	tb.Add(tb.iconList)

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
// at the end of this TabBar list of tabs.
// Returns the pointer to thew new Tab.
func (tb *TabBar) AddTab(text string) *Tab {

	tab := tb.InsertTab(text, len(tb.tabs))
	tb.SetSelected(len(tb.tabs) - 1)
	return tab
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
	tab := newTab(text, tb, &tb.styles.Tab)
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
// Returns an error if the position is invalid.
func (tb *TabBar) RemoveTab(pos int) error {

	// Check position to remove from
	if pos < 0 || pos >= len(tb.tabs) {
		return fmt.Errorf("Invalid tab position:%d", pos)
	}

	// Remove tab from TabBar panel
	tab := tb.tabs[pos]
	tb.Remove(&tab.header)
	tb.Remove(&tab.content)

	// Remove tab from tabbar array
	copy(tb.tabs[pos:], tb.tabs[pos+1:])
	tb.tabs[len(tb.tabs)-1] = nil
	tb.tabs = tb.tabs[:len(tb.tabs)-1]

	// Checks if removed tab was selected
	if tb.selected == pos {
		// TODO
	}

	tb.update()
	tb.recalc()
	return nil
}

// TabCount returns the current number of Tabs in the TabBar
func (tb *TabBar) TabCount() int {

	return len(tb.tabs)
}

// TabAt returns the pointer of the Tab object at the specified position.
// Return nil if the position is invalid
func (tb *TabBar) TabAt(pos int) *Tab {

	if pos < 0 || pos >= len(tb.tabs) {
		return nil
	}
	return tb.tabs[pos]
}

// TabPosition returns the position of the Tab specified by its pointer
func (tb *TabBar) TabPosition(tab *Tab) int {

	for i := 0; i < len(tb.tabs); i++ {
		if tb.tabs[i] == tab {
			return i
		}
	}
	return -1
}

// SetSelected sets the selected tab of the TabBar to the tab with the specified position.
// Returns the pointer of the selected tab or nil if the position is invalid.
func (tb *TabBar) SetSelected(pos int) *Tab {

	if pos < 0 || pos >= len(tb.tabs) {
		return nil
	}
	for i := 0; i < len(tb.tabs); i++ {
		if i == pos {
			tb.tabs[i].setSelected(true)
		} else {
			tb.tabs[i].setSelected(false)
		}
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

// onListButtonMouse process subscribed MouseButton events over the list button
func (tb *TabBar) onListButton(evname string, ev interface{}) {

	switch evname {
	case OnMouseDown:
		if !tb.list.Visible() {
			tb.list.SetVisible(true)
		}
	default:
		return
	}
	tb.root.StopPropagation(StopAll)
}

// onListChange process OnChange event from the tab list
func (tb *TabBar) onListChange(evname string, ev interface{}) {

	selected := tb.list.Selected()
	pos := selected[0].GetPanel().UserData().(int)
	tb.SetSelected(pos)
	tb.list.SetVisible(false)
}

// applyStyle applies the specified TabBar style
func (tb *TabBar) applyStyle(s *TabBarStyle) {

	tb.SetBordersFrom(&s.Border)
	tb.SetBordersColor4(&s.BorderColor)
	tb.SetPaddingsFrom(&s.Paddings)
	tb.SetColor4(&s.BgColor)
	tb.separator.SetColor4(&s.BorderColor)
}

// recalc recalculates and updates the positions of all tabs
func (tb *TabBar) recalc() {

	// Determines how many tabs could be fully shown
	iconWidth := tb.iconList.Width()
	availWidth := tb.ContentWidth() - iconWidth
	var tabWidth float32
	var totalWidth float32
	var count int
	for i := 0; i < len(tb.tabs); i++ {
		tab := tb.tabs[i]
		minw := tab.minWidth()
		if minw > tabWidth {
			tabWidth = minw
		}
		totalWidth = float32(count+1) * tabWidth
		if totalWidth > availWidth {
			break
		}
		count++
	}

	tb.list.Clear()
	if count < len(tb.tabs) {
		// Sets the list button visible andposition
		tb.iconList.SetVisible(true)
		height := tb.tabs[0].header.Height()
		iy := (height - tb.iconList.Height()) / 2
		tb.iconList.SetPosition(availWidth, iy)
		// Sets the tab list position and size
		listWidth := float32(200)
		lx := tb.ContentWidth() - listWidth
		ly := height + 1
		tb.list.SetPosition(lx, ly)
		tb.list.SetSize(listWidth, 200)
		tb.SetTopChild(tb.list)

	} else {
		tb.iconList.SetVisible(false)
		tb.list.SetVisible(false)
	}

	var headerx float32
	twidth := availWidth / float32(count)
	for i := 0; i < len(tb.tabs); i++ {
		tab := tb.tabs[i]
		// If Tab can be shown
		if i < count {
			tab.recalc(twidth)
			tab.header.SetPosition(headerx, 0)
			// Sets size and position of the Tab content panel
			contentx := float32(0)
			contenty := tab.header.Height() + tb.styles.SepHeight
			tab.content.SetWidth(tb.ContentWidth())
			tab.content.SetHeight(tb.ContentHeight() - contenty)
			tab.content.SetPosition(contentx, contenty)
			headerx += tab.header.Width()
			tab.header.SetVisible(true)
			continue
			// Tab cannot be shown, insert into vertical list
		} else {
			tab.header.SetVisible(false)
			item := NewImageLabel(tab.label.Text())
			item.SetUserData(i)
			tb.list.Add(item)
		}
	}

	// Sets the separator size, position and visibility
	if len(tb.tabs) > 0 {
		tb.separator.SetSize(tb.ContentWidth(), tb.styles.SepHeight)
		tb.separator.SetPositionY(tb.tabs[0].header.Height())
		tb.separator.SetVisible(true)
	} else {
		tb.separator.SetVisible(false)
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
	tb         *TabBar    // Pointer to parent *TabBar
	styles     *TabStyles // Pointer to Tab current styles
	header     Panel      // Tab header
	label      *Label     // Tab user label
	iconClose  *Label     // Tab close icon
	icon       *Label     // Tab optional user icon
	img        *Image     // Tab optional user image
	bottom     Panel
	content    Panel // User content panel
	cursorOver bool
	selected   bool
}

// newTab creates and returns a pointer to a new Tab
func newTab(text string, tb *TabBar, styles *TabStyles) *Tab {

	tab := new(Tab)
	tab.tb = tb
	tab.styles = styles
	// Setup the header panel
	tab.header.Initialize(0, 0)
	tab.label = NewLabel(text)
	tab.iconClose = NewLabel(styles.IconClose, true)
	tab.header.Add(tab.label)
	tab.header.Add(tab.iconClose)
	// Creates the bottom panel
	tab.bottom.Initialize(0, 0)
	tab.bottom.SetBounded(false)
	tab.bottom.SetColor4(&tab.styles.Selected.BgColor)
	tab.header.Add(&tab.bottom)
	tab.content.Initialize(0, 0)

	// Subscribe to header panel events
	tab.header.Subscribe(OnCursorEnter, tab.onCursor)
	tab.header.Subscribe(OnCursorLeave, tab.onCursor)
	tab.header.Subscribe(OnMouseDown, tab.onMouseHeader)
	tab.iconClose.Subscribe(OnMouseDown, tab.onMouseIcon)

	tab.update()
	return tab
}

// onCursor process subscribed cursor events over the tab header
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

// onMouse process subscribed mouse events over the tab header
func (tab *Tab) onMouseHeader(evname string, ev interface{}) {

	switch evname {
	case OnMouseDown:
		tab.tb.SetSelected(tab.tb.TabPosition(tab))
	default:
		return
	}
	tab.header.root.StopPropagation(StopAll)
}

// onMouseIcon process subscribed mouse events over the tab close icon
func (tab *Tab) onMouseIcon(evname string, ev interface{}) {

	switch evname {
	case OnMouseDown:
		tab.tb.RemoveTab(tab.tb.TabPosition(tab))
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

func (tab *Tab) setSelected(selected bool) {

	tab.selected = selected
	tab.content.SetVisible(selected)
	tab.bottom.SetVisible(selected)
	tab.update()
	tab.setBottomPanel()
}

// minWidth returns the minimum width of this Tab header to allow
// all of its elements to be shown in full.
func (tab *Tab) minWidth() float32 {

	var minWidth float32
	if tab.icon != nil {
		minWidth = tab.icon.Width()
	} else if tab.img != nil {
		minWidth = tab.img.Width()
	}
	minWidth += tab.label.Width()
	minWidth += tab.iconClose.Width()
	return minWidth + tab.header.MinWidth()
}

// applyStyle applies the specified Tab style to the Tab header
func (tab *Tab) applyStyle(s *TabStyle) {

	tab.header.SetMarginsFrom(&s.Margins)
	tab.header.SetBordersFrom(&s.Border)
	tab.header.SetBordersColor4(&s.BorderColor)
	tab.header.SetPaddingsFrom(&s.Paddings)
	tab.header.SetColor4(&s.BgColor)
}

// update updates the Tab header visual style
func (tab *Tab) update() {

	if !tab.header.Enabled() {
		tab.applyStyle(&tab.styles.Disabled)
		return
	}
	if tab.selected {
		tab.applyStyle(&tab.styles.Selected)
		return
	}
	if tab.cursorOver {
		tab.applyStyle(&tab.styles.Over)
		return
	}
	tab.applyStyle(&tab.styles.Normal)
}

// setBottomPanel sets the position and size of the Tab bottom panel
// to cover the Tabs separator
func (tab *Tab) setBottomPanel() {

	if tab.selected {
		bwidth := tab.header.ContentWidth() + tab.header.Paddings().Left + tab.header.Paddings().Right
		bx := tab.styles.Selected.Margins.Left + tab.styles.Selected.Border.Left
		tab.bottom.SetSize(bwidth, tab.tb.styles.SepHeight)
		tab.bottom.SetPosition(bx, tab.header.Height())
	}
}

// recalc recalculates the size of the Tab header and the size
// and positions of the Tab header internal panels
func (tab *Tab) recalc(width float32) {

	height := tab.label.Height()
	tab.header.SetContentHeight(height)
	tab.header.SetWidth(width)

	labx := float32(0)
	if tab.icon != nil {
		tab.icon.SetPosition(0, 0)
		labx = tab.icon.Width()
	} else if tab.img != nil {
		tab.img.SetPosition(0, 0)
		labx = tab.img.Width()
	}
	tab.label.SetPosition(labx, 0)

	// Sets the close icon position
	icx := tab.header.ContentWidth() - tab.iconClose.Width()
	icy := (tab.header.ContentHeight() - tab.iconClose.Height()) / 2
	tab.iconClose.SetPosition(icx, icy)

	// Sets the position of the bottom panel to cover separator
	tab.setBottomPanel()
}
