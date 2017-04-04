// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

type MenuBar struct {
}

type Menu struct {
	Panel             // embedded panel
	items []*MenuItem // menu items
}

type MenuItem struct {
	Panel             // embedded panel
	label   *Label    // optional internal label (nil for separators)
	image   *Image    // optional left internal image
	licon   *Label    // optional left internal icon label
	ricon   *Label    // optional right internal icon label for submenu
	icode   int       // icon code (if icon is set)
	subm    *MenuItem // optional pointer to sub menu
	shorcut int32     // shortcut code
	enabled bool      // enabled state
}

// NewMenu creates and returns a pointer to a new empty menu
func NewMenu() *Menu {

	m := new(Menu)
	m.items = make([]*MenuItem, 0)
	return m
}

// AddItem creates and adds a new menu item to this menu and returns the pointer
// to the created item.
func (m *Menu) AddItem(text string) *MenuItem {

	mi := new(MenuItem)
	mi.label = NewLabel(text)
	mi.Panel.Add(mi.label)

	m.items = append(m.items, mi)
	return mi
}

// AddSeparator creates and adds a new separator to the menu
func (m *Menu) AddSeparator() *MenuItem {

	mi := new(MenuItem)
	return mi
}

// RemoveItem removes the specified menu item from this menu
func (m *Menu) RemoveItem(mi *MenuItem) {

}

// recalc recalculates the positions of this menu internal items
func (m *Menu) recalc() {

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

// recalc recalculates the positions of this menu item internal panels
func (mi *MenuItem) recalc() {
}
