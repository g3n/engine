// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/g3n/engine/gui/assets"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/window"
)

type Menu struct {
	Panel                // embedded panel
	styles   *MenuStyles // pointer to current styles
	bar      bool        // true for menu bar
	items    []*MenuItem // menu items
	active   bool        // menu active state
	autoOpen bool        // open sub menus when mouse over if true
	mitem    *MenuItem   // parent menu item for sub menu
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

// MenuStyles describes all styles of the menu body and menu item
type MenuStyles struct {
	Body *MenuBodyStyles // Menu body styles
	Item *MenuItemStyles // Menu item styles
}

// MenuItem is an option of a Menu
type MenuItem struct {
	Panel                          // embedded panel
	styles      *MenuItemStyles    // pointer to current styles
	menu        *Menu              // pointer to parent menu
	licon       *Label             // optional left icon label
	label       *Label             // optional text label (nil for separators)
	shortcut    *Label             // optional shorcut text label
	ricon       *Label             // optional right internal icon label for submenu
	icode       int                // icon code (if icon is set)
	submenu     *Menu              // pointer to optional associated sub menu
	keyModifier window.ModifierKey // shortcut key modifier
	keyCode     window.Key         // shortcut key code
	disabled    bool               // item disabled state
	selected    bool               // selection state
}

// MenuItemStyle describes the style of a menu item
type MenuItemStyle struct {
	Border           BorderSizes
	Paddings         BorderSizes
	BorderColor      math32.Color4
	BgColor          math32.Color
	FgColor          math32.Color
	IconPaddings     BorderSizes
	ShortcutPaddings BorderSizes
	RiconPaddings    BorderSizes
}

// MenuItemStyles describes all the menu item styles
type MenuItemStyles struct {
	Normal    MenuItemStyle
	Over      MenuItemStyle
	Disabled  MenuItemStyle
	Separator MenuItemStyle
}

var mapKeyModifier = map[window.ModifierKey]string{
	window.ModShift:   "Shift",
	window.ModControl: "Ctrl",
	window.ModAlt:     "Alt",
}
var mapKeyText = map[window.Key]string{
	window.KeyApostrophe: "'",
	window.KeyComma:      ",",
	window.KeyMinus:      "-",
	window.KeyPeriod:     ".",
	window.KeySlash:      "/",
	window.Key0:          "0",
	window.Key1:          "1",
	window.Key2:          "2",
	window.Key3:          "3",
	window.Key4:          "4",
	window.Key5:          "5",
	window.Key6:          "6",
	window.Key7:          "7",
	window.Key8:          "8",
	window.Key9:          "9",
	window.KeySemicolon:  ";",
	window.KeyEqual:      "=",
	window.KeyA:          "A",
	window.KeyB:          "B",
	window.KeyC:          "C",
	window.KeyD:          "D",
	window.KeyE:          "E",
	window.KeyF:          "F",
	window.KeyG:          "G",
	window.KeyH:          "H",
	window.KeyI:          "I",
	window.KeyJ:          "J",
	window.KeyK:          "K",
	window.KeyL:          "L",
	window.KeyM:          "M",
	window.KeyN:          "N",
	window.KeyO:          "O",
	window.KeyP:          "P",
	window.KeyQ:          "Q",
	window.KeyR:          "R",
	window.KeyS:          "S",
	window.KeyT:          "T",
	window.KeyU:          "U",
	window.KeyV:          "V",
	window.KeyW:          "W",
	window.KeyX:          "X",
	window.KeyY:          "Y",
	window.KeyZ:          "Z",
	window.KeyF1:         "F1",
	window.KeyF2:         "F2",
	window.KeyF3:         "F3",
	window.KeyF4:         "F4",
	window.KeyF5:         "F5",
	window.KeyF6:         "F6",
	window.KeyF7:         "F7",
	window.KeyF8:         "F8",
	window.KeyF9:         "F9",
	window.KeyF10:        "F10",
	window.KeyF11:        "F11",
	window.KeyF12:        "F12",
}

func NewMenuBar() *Menu {

	m := NewMenu()
	m.bar = true
	return m
}

// NewMenu creates and returns a pointer to a new empty menu
func NewMenu() *Menu {

	m := new(Menu)
	m.Panel.Initialize(0, 0)
	m.styles = &StyleDefault.Menu
	m.items = make([]*MenuItem, 0)
	m.Panel.Subscribe(OnCursorEnter, m.onCursor)
	m.Panel.Subscribe(OnCursorLeave, m.onCursor)
	m.Panel.Subscribe(OnKeyDown, m.onKey)
	m.Panel.Subscribe(OnMouseOut, m.onMouse)
	m.update()
	return m
}

// AddOption creates and adds a new menu item to this menu with the
// specified text and returns the pointer to the created menu item.
func (m *Menu) AddOption(text string) *MenuItem {

	mi := newMenuItem(text, m.styles.Item)
	m.Panel.Add(mi)
	m.items = append(m.items, mi)
	mi.menu = m
	m.recalc()
	return mi
}

// AddSeparator creates and adds a new separator to the menu
func (m *Menu) AddSeparator() *MenuItem {

	mi := newMenuItem("", m.styles.Item)
	m.Panel.Add(mi)
	m.items = append(m.items, mi)
	mi.menu = m
	m.recalc()
	return mi
}

// AddMenu creates and adds a new menu item to this menu with the
// specified text and sub menu.
// Returns the pointer to the created menu item.
func (m *Menu) AddMenu(text string, subm *Menu) *MenuItem {

	mi := newMenuItem(text, m.styles.Item)
	m.Panel.Add(mi)
	m.items = append(m.items, mi)
	mi.submenu = subm
	mi.submenu.SetVisible(false)
	mi.submenu.SetBounded(false)
	mi.submenu.mitem = mi
	mi.submenu.autoOpen = true
	mi.menu = m
	mi.ricon = NewIconLabel(string(assets.ChevronRight))
	mi.Panel.Add(mi.ricon)
	mi.Panel.Add(mi.submenu)
	mi.update()
	m.recalc()
	return nil
}

// RemoveItem removes the specified menu item from this menu
func (m *Menu) RemoveItem(mi *MenuItem) {

}

// onCursor process subscribed cursor events
func (m *Menu) onCursor(evname string, ev interface{}) {

	if evname == OnCursorEnter {
		m.root.SetKeyFocus(m)
		m.active = true
	} else if evname == OnCursorLeave {
		m.active = false
		// If this is a sub menu and the parent menu item is not selected
		// hides this sub menu
		//if m.mitem != nil && !m.mitem.selected {
		//	m.SetVisible(false)
		//}
	}
	m.root.StopPropagation(StopAll)
}

// onKey process subscribed key events
func (m *Menu) onKey(evname string, ev interface{}) {

	sel := m.selectedPos()
	kev := ev.(*window.KeyEvent)
	switch kev.Keycode {
	// Select next enabled menu item
	case window.KeyDown:
		next := m.nextItem(sel)
		m.setSelectedPos(next)
	// Select previous enabled menu item
	case window.KeyUp:
		prev := m.prevItem(sel)
		m.setSelectedPos(prev)
	// Return to previous menu
	case window.KeyLeft:
		if m.mitem != nil {
			m.active = false
			m.mitem.menu.setSelectedItem(m.mitem)
			m.root.SetKeyFocus(m.mitem.menu)
		}
	// Enter into sub menu
	case window.KeyRight:
		if sel < 0 {
			return
		}
		mi := m.items[sel]
		if mi.submenu != nil {
			m.root.SetKeyFocus(mi.submenu)
			mi.submenu.setSelectedPos(0)
		}
	case window.KeyEnter:
	default:
		return
	}
}

// onMouse process subscribed mouse events for the menu
func (m *Menu) onMouse(evname string, ev interface{}) {

	if evname == OnMouseOut {
		if m.bar {
			m.autoOpen = false
			m.setSelectedPos(-1)
		}
	}
}

// setSelectedPos sets the menu item at the specified position as selected
// and all others as not selected.
func (m *Menu) setSelectedPos(pos int) {

	for i := 0; i < len(m.items); i++ {
		mi := m.items[i]
		if i == pos {
			mi.selected = true
		} else {
			mi.selected = false
		}
		// If menu item has a sub menu, unselects the sub menu options recursively
		if mi.submenu != nil {
			mi.submenu.setSelectedPos(-1)
		}
		mi.update()
	}
}

// setSelectedItem sets the specified menu item as selected
// and all others as not selected
func (m *Menu) setSelectedItem(mitem *MenuItem) {

	for i := 0; i < len(m.items); i++ {
		mi := m.items[i]
		if mi == mitem {
			mi.selected = true
		} else {
			mi.selected = false
		}
		// If menu item has a sub menu, unselects the sub menu options recursively
		if mi.submenu != nil {
			mi.submenu.setSelectedItem(nil)
		}
		mi.update()
	}
}

// selectedPos returns the position of the current selected menu item
// Returns -1 if no item selected
func (m *Menu) selectedPos() int {

	for i := 0; i < len(m.items); i++ {
		mi := m.items[i]
		if mi.selected {
			return i
		}
	}
	return -1
}

// nextItem returns the position of the next enabled option from the
// specified position
func (m *Menu) nextItem(pos int) int {

	res := 0
	for i := pos + 1; i < len(m.items); i++ {
		mi := m.items[i]
		if mi.disabled || mi.label == nil {
			continue
		}
		res = i
		break
	}
	return res
}

// prevItem returns the position of previous enabled menu item from
// the specified position
func (m *Menu) prevItem(pos int) int {

	res := len(m.items) - 1
	for i := pos - 1; i >= 0 && i < len(m.items); i-- {
		mi := m.items[i]
		if mi.disabled || mi.label == nil {
			continue
		}
		res = i
		break
	}
	return res
}

// update updates the menu visual state
func (m *Menu) update() {

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

	if m.bar {
		m.recalcBar()
		return
	}

	// Find the maximum icon and label widths
	minWidth := float32(0)
	iconWidth := float32(0)
	labelWidth := float32(0)
	shortcutWidth := float32(0)
	riconWidth := float32(0)
	for i := 0; i < len(m.items); i++ {
		mi := m.items[i]
		minWidth = mi.MinWidth()
		// Separator
		if mi.label == nil {
			continue
		}
		// Left icon width
		if mi.licon != nil && mi.licon.width > iconWidth {
			iconWidth = mi.licon.width
		}
		// Option label width
		if mi.label.width > labelWidth {
			labelWidth = mi.label.width
		}
		// Shortcut label width
		if mi.shortcut != nil && mi.shortcut.width > shortcutWidth {
			shortcutWidth = mi.shortcut.width
		}
		// Right icon (submenu indicator) width
		if mi.ricon != nil && mi.ricon.width > riconWidth {
			riconWidth = mi.ricon.width
		}
	}
	width := minWidth + iconWidth + labelWidth + shortcutWidth + riconWidth

	// Sets the position and width of the menu items
	// The height is defined by the menu item itself
	px := float32(0)
	py := float32(0)
	for i := 0; i < len(m.items); i++ {
		mi := m.items[i]
		mi.SetPosition(px, py)
		mh := mi.minHeight()
		py += mh
		mi.SetSize(width, mh)
		mi.recalc(iconWidth, labelWidth, shortcutWidth)
	}
	m.SetContentSize(width, py)
}

// recalcBar recalculates the positions of this MenuBar internal items
// and the content width and height of the menu
func (m *Menu) recalcBar() {

	height := float32(0)
	for i := 0; i < len(m.items); i++ {
		mi := m.items[i]
		if mi.minHeight() > height {
			height = mi.minHeight()
		}
	}

	px := float32(0)
	for i := 0; i < len(m.items); i++ {
		mi := m.items[i]
		mi.SetPosition(px, 0)
		width := float32(0)
		width = mi.minWidth()
		mi.SetSize(width, height)
		px += mi.Width()
	}
	m.SetContentSize(px, height)
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
		mi.Panel.Subscribe(OnMouseUp, mi.onMouse)
		mi.Panel.Subscribe(OnMouseDown, mi.onMouse)
	}
	mi.update()
	return mi
}

// SetIcon sets the left icon of this menu item
// If an image was previously set it is replaced by this icon
func (mi *MenuItem) SetIcon(icode int) *MenuItem {

	mi.licon = NewIconLabel(string(icode))
	mi.Panel.Add(mi.licon)
	mi.update()
	return mi
}

// SetImage sets the left image of this menu item
// If an icon was previously set it is replaced by this image
func (mi *MenuItem) SetImage(img *Image) {

}

// SetText sets the text of this menu item
func (mi *MenuItem) SetText(text string) *MenuItem {

	if mi.label == nil {
		return mi
	}
	mi.label.SetText(text)
	mi.update()
	mi.menu.recalc()
	return mi
}

// SetShortcut sets the keyboard shortcut of this menu item
func (mi *MenuItem) SetShortcut(mod window.ModifierKey, key window.Key) *MenuItem {

	if mapKeyText[key] == "" {
		panic("Invalid menu shortcut key")
	}
	mi.keyModifier = mod
	mi.keyCode = key
	text := ""
	if mi.keyModifier&window.ModShift != 0 {
		text = mapKeyModifier[window.ModShift]
	}
	if mi.keyModifier&window.ModControl != 0 {
		if text != "" {
			text += "+"
		}
		text += mapKeyModifier[window.ModControl]
	}
	if mi.keyModifier&window.ModAlt != 0 {
		if text != "" {
			text += "+"
		}
		text += mapKeyModifier[window.ModAlt]
	}
	if text != "" {
		text += "+"
	}
	text += mapKeyText[key]

	mi.shortcut = NewLabel(text)
	mi.Panel.Add(mi.shortcut)
	mi.update()
	mi.menu.recalc()
	return mi
}

// SetSubmenu sets an associated sub menu item for this menu item
func (mi *MenuItem) SetSubmenu(smi *MenuItem) *MenuItem {

	return mi
}

// SetEnabled sets the enabled state of this menu item
func (mi *MenuItem) SetEnabled(enabled bool) *MenuItem {

	mi.disabled = !enabled
	mi.update()
	return mi
}

// onCursor processes subscribed cursor events over the menu item
func (mi *MenuItem) onCursor(evname string, ev interface{}) {

	switch evname {
	case OnCursorEnter:
		mi.menu.setSelectedItem(mi)
	case OnCursorLeave:
		//if mi.submenu != nil && mi.submenu.active {
		//	return
		//}
		//mi.menu.setSelectedItem(nil)
	}
}

// onMouse processes subscribed mouse events over the menu item
func (mi *MenuItem) onMouse(evname string, ev interface{}) {

	switch evname {
	case OnMouseDown:
		// MenuBar option
		if mi.menu.bar {
			mi.menu.autoOpen = !mi.menu.autoOpen
			if mi.submenu != nil && mi.submenu.Visible() {
				mi.submenu.SetVisible(false)
				return
			}
			mi.update()
			//if mi.submenu != nil {
			//	if !mi.submenu.Visible() {
			//		mi.submenu.SetVisible(true)
			//		mi.submenu.SetPosition(0, mi.Height()-2)
			//	} else {
			//		mi.submenu.SetVisible(false)
			//	}
			//} else {
			//	// Dispatch on click
			//}
		} else {

		}
	case OnMouseUp:
	}
}

// update updates the menu item visual state
func (mi *MenuItem) update() {

	// Separator
	if mi.label == nil {
		mi.applyStyle(&mi.styles.Separator)
		return
	}
	// Disabled item
	if mi.disabled {
		mi.applyStyle(&mi.styles.Disabled)
		return
	}
	// Selected item
	if mi.selected {
		mi.applyStyle(&mi.styles.Over)
		if mi.submenu != nil && mi.menu.autoOpen {
			mi.menu.SetTopChild(mi)
			mi.submenu.SetVisible(true)
			if mi.menu != nil && mi.menu.bar {
				mi.submenu.SetPosition(0, mi.Height()-2)
			} else {
				mi.submenu.SetPosition(mi.Width()-2, 0)
			}
		}
		return
	}
	// If this menu item has a sub menu and the sub menu is not active,
	// hides the sub menu
	if mi.submenu != nil {
		mi.submenu.SetVisible(false)
	}
	mi.applyStyle(&mi.styles.Normal)
}

// applyStyle applies the specified menu item style
func (mi *MenuItem) applyStyle(mis *MenuItemStyle) {

	mi.SetBordersFrom(&mis.Border)
	mi.SetBordersColor4(&mis.BorderColor)
	mi.SetPaddingsFrom(&mis.Paddings)
	mi.SetColor(&mis.BgColor)
	if mi.licon != nil {
		mi.licon.SetPaddingsFrom(&mis.IconPaddings)
	}
	if mi.label != nil {
		mi.label.SetColor(&mis.FgColor)
	}
	if mi.shortcut != nil {
		mi.shortcut.SetPaddingsFrom(&mis.ShortcutPaddings)
	}
	if mi.ricon != nil {
		mi.ricon.SetPaddingsFrom(&mis.RiconPaddings)
	}
}

// recalc recalculates the positions of this menu item internal panels
func (mi *MenuItem) recalc(iconWidth, labelWidth, shortcutWidth float32) {

	// Separator
	if mi.label == nil {
		return
	}
	if mi.licon != nil {
		py := (mi.label.height - mi.licon.height) / 2
		mi.licon.SetPosition(0, py)
	}
	mi.label.SetPosition(iconWidth, 0)
	if mi.shortcut != nil {
		mi.shortcut.SetPosition(iconWidth+labelWidth, 0)
	}
	if mi.ricon != nil {
		mi.ricon.SetPosition(iconWidth+labelWidth+shortcutWidth, 0)
	}
}

// minHeight returns the minimum height of this menu item
func (mi *MenuItem) minHeight() float32 {

	mh := mi.MinHeight()
	if mi.label == nil {
		return mh + 1
	}
	mh += mi.label.height
	return mh
}

// minWidth returns the minimum width of this menu item
func (mi *MenuItem) minWidth() float32 {

	mw := mi.MinWidth()
	if mi.label == nil {
		return mw + 1
	}
	mw += mi.label.width
	return mw
}
