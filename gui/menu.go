// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/g3n/engine/math32"
)

type MenuBar struct {
}

type Menu struct {
	Panel              // embedded panel
	styles *MenuStyles // pointer to current styles
}

// MenuBodyStyle describes the style of the menu body
type MenuBodyStyle struct {
	Border      BorderSizes
	Paddings    BorderSizes
	BorderColor math32.Color4
	BgColor     math32.Color
	FgColor     math32.Color
}

// MenuBodyStyles describes all styles of the menu body
type MenuBodyStyles struct {
	Normal   MenuBodyStyle
	Over     MenuBodyStyle
	Focus    MenuBodyStyle
	Disabled MenuBodyStyle
}

// MenuStyles describes all styles of the menu (body and item)
type MenuStyles struct {
	Body *MenuBodyStyles // Menu body styles
	Item *MenuItemStyles // Menu item styles
}

// MenuItem is an option of a Menu
type MenuItem struct {
	Panel                     // embedded panel
	styles    *MenuItemStyles // pointer to current styles
	label     *Label          // optional internal label (nil for separators)
	image     *Image          // optional left internal image
	licon     *Label          // optional left internal icon label
	ricon     *Label          // optional right internal icon label for submenu
	icode     int             // icon code (if icon is set)
	subm      *MenuItem       // optional pointer to sub menu
	shorcut   int32           // shortcut code
	enabled   bool            // enabled state
	mouseOver bool
}

// MenuItemStyle describes the style of a menu item
type MenuItemStyle struct {
	Border      BorderSizes
	Paddings    BorderSizes
	BorderColor math32.Color4
	BgColor     math32.Color
	FgColor     math32.Color
}

// MenuItemStyles describes all the menu item styles
type MenuItemStyles struct {
	Normal    MenuItemStyle
	Over      MenuItemStyle
	Disabled  MenuItemStyle
	Separator MenuItemStyle
}

// NewMenu creates and returns a pointer to a new empty menu
func NewMenu() *Menu {

	m := new(Menu)
	m.Panel.Initialize(0, 0)
	m.styles = &StyleDefault.Menu
	m.update()
	return m
}

// AddItem creates and adds a new menu item to this menu and returns the pointer
// to the created item.
func (m *Menu) AddItem(text string) *MenuItem {

	mi := newMenuItem(text, m.styles.Item)
	m.Panel.Add(mi)
	m.recalc()
	return mi
}

// AddSeparator creates and adds a new separator to the menu
func (m *Menu) AddSeparator() *MenuItem {

	mi := newMenuItem("", m.styles.Item)
	m.Panel.Add(mi)
	m.recalc()
	return mi
}

// RemoveItem removes the specified menu item from this menu
func (m *Menu) RemoveItem(mi *MenuItem) {

}

// update updates the menu visual state
func (m *Menu) update() {

	//if s.cursorOver {
	//	s.applyStyle(&s.styles.Over)
	//	return
	//}
	//if s.focus {
	//	s.applyStyle(&s.styles.Focus)
	//	return
	//}
	m.applyStyle(&m.styles.Body.Normal)
}

// applyStyle applies the specified menu body style
func (m *Menu) applyStyle(mbs *MenuBodyStyle) {

	m.SetBordersFrom(&mbs.Border)
	m.SetBordersColor4(&mbs.BorderColor)
	m.SetPaddingsFrom(&mbs.Paddings)
	m.SetColor(&mbs.BgColor)
}

// recalc recalculates the positions of this menu internal items
// and the content width and height of the menu
func (m *Menu) recalc() {

	// Find maximum item width
	maxw := float32(0)
	for i := 0; i < len(m.Children()); i++ {
		minw := m.Children()[i].(*MenuItem).minWidth()
		if minw > maxw {
			maxw = minw
		}
	}

	// Sets the position and width of the menu items
	// The height is defined by the menu item itself
	px := float32(0)
	py := float32(0)
	for i := 0; i < len(m.Children()); i++ {
		mi := m.Children()[i].(*MenuItem)
		mi.SetPosition(px, py)
		py += mi.height
		mi.SetWidth(maxw)
		log.Error("width:%v", maxw)
		mi.recalc()
	}
	m.SetContentSize(maxw, py)
}

// newMenuItem creates and returns a pointer to a new menu item
// with the specified text.
func newMenuItem(text string, styles *MenuItemStyles) *MenuItem {

	mi := new(MenuItem)
	mi.Panel.Initialize(0, 0)
	mi.styles = styles
	if text != "" {
		mi.label = NewLabel(text)
		mi.Panel.Add(mi.label)
		mi.Panel.Subscribe(OnCursorEnter, mi.onCursor)
		mi.Panel.Subscribe(OnCursorLeave, mi.onCursor)
	}
	mi.recalc()
	mi.update()
	return mi
}

// SetIcon sets the left icon of this menu item
// If an image was previously set it is replaced by this icon
func (mi *MenuItem) SetIcon(icode int) *MenuItem {

	return mi
}

// SetImage sets the left image of this menu item
// If an icon was previously set it is replaced by this image
func (mi *MenuItem) SetImage(img *Image) *MenuItem {

	return mi
}

// SetText sets the text of this menu item
func (mi *MenuItem) SetText(text string) *MenuItem {

	return mi
}

// SetShortcut sets the keyboard shortcut of this menu item
func (mi *MenuItem) SetShortcut(text string) *MenuItem {

	return mi
}

// SetSubmenu sets an associated sub menu item for this menu item
func (mi *MenuItem) SetSubmenu(smi *MenuItem) *MenuItem {

	return mi
}

// SetEnabled sets the enabled state of this menu item
func (mi *MenuItem) SetEnabled(enabled bool) *MenuItem {

	return mi
}

func (mi *MenuItem) onCursor(evname string, ev interface{}) {

	switch evname {
	case OnCursorEnter:
		mi.mouseOver = true
		mi.update()
	case OnCursorLeave:
		mi.mouseOver = false
		mi.update()
	}
}

// update updates the menu item visual state
func (mi *MenuItem) update() {

	// Separator
	if mi.label == nil {
		mi.applyStyle(&mi.styles.Separator)
		return
	}
	if mi.mouseOver {
		mi.applyStyle(&mi.styles.Over)
		return
	}
	mi.applyStyle(&mi.styles.Normal)
}

// applyStyle applies the specified menu item style
func (mi *MenuItem) applyStyle(mis *MenuItemStyle) {

	mi.SetBordersFrom(&mis.Border)
	mi.SetBordersColor4(&mis.BorderColor)
	mi.SetPaddingsFrom(&mis.Paddings)
	mi.SetColor(&mis.BgColor)
}

// recalc recalculates the positions of this menu item internal panels
// and the total height of the menu item.
func (mi *MenuItem) recalc() {

	// Separator
	if mi.label == nil {
		mi.SetHeight(4)
		return
	}
	h := mi.label.height
	w := mi.label.width
	mi.SetContentSize(w, h)
	log.Error("menuitem size: %v/%v", w, h)
}

// minWidth returns the minimum width of this menu item
func (mi *MenuItem) minWidth() float32 {

	mw := mi.MinWidth()
	if mi.licon != nil {
		mw += mi.licon.width
	}
	if mi.image != nil {
		mw += mi.image.width
	}
	if mi.label != nil {
		mw += mi.label.width
	}
	if mi.ricon != nil {
		mw += mi.ricon.width
	}
	return mw
}

// minHeight returns the minimum height of this menu item
func (mi *MenuItem) minHeight() float32 {

	mh := mi.MinHeight()
	if mi.label == nil {
		return mh + 4
	}
	mh += mi.label.height
	return mh
}
