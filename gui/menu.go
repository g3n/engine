// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/g3n/engine/gui/assets"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/window"
)

type MenuBar struct {
}

type Menu struct {
	Panel              // embedded panel
	styles *MenuStyles // pointer to current styles
	items  []*MenuItem // menu items
	active bool        // menu active state
	mitem  *MenuItem   // parent menu item for sub menus
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
	menu        *Menu              // pointer to container menu
	licon       *Label             // optional left icon label
	label       *Label             // optional text label (nil for separators)
	shortcut    *Label             // optional shorcut text label
	ricon       *Label             // optional right internal icon label for submenu
	icode       int                // icon code (if icon is set)
	subm        *Menu              // optional pointer to sub menu
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

// NewMenu creates and returns a pointer to a new empty menu
func NewMenu() *Menu {

	m := new(Menu)
	m.Panel.Initialize(0, 0)
	m.styles = &StyleDefault.Menu
	m.items = make([]*MenuItem, 0)
	m.Panel.Subscribe(OnCursorEnter, m.onCursor)
	m.Panel.Subscribe(OnCursorLeave, m.onCursor)
	m.Panel.Subscribe(OnKeyDown, m.onKey)
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
	mi.subm = subm
	mi.subm.SetVisible(false)
	mi.subm.SetBounded(false)
	mi.subm.mitem = mi
	mi.menu = m
	mi.ricon = NewIconLabel(string(assets.ChevronRight))
	mi.Panel.Add(mi.ricon)
	mi.Panel.Add(mi.subm)
	mi.update()
	m.recalc()
	return nil
}

// RemoveItem removes the specified menu item from this menu
func (m *Menu) RemoveItem(mi *MenuItem) {

}

// onCursor process subscribed cursor events
func (m *Menu) onCursor(evname string, ev interface{}) {

	log.Error("evname:%s / %v", evname, ev)
	if evname == OnCursorEnter {
		m.root.SetKeyFocus(m)
		m.active = true
	} else if evname == OnCursorLeave {
		m.root.SetKeyFocus(nil)
		m.active = false
		if m.mitem != nil && !m.mitem.selected {
			m.SetVisible(false)
		}
	}
	m.root.StopPropagation(StopAll)
}

// onKey process subscribed key events
func (m *Menu) onKey(evname string, ev interface{}) {

	prevsel := m.selectedItem()
	var sel int
	kev := ev.(*window.KeyEvent)
	switch kev.Keycode {
	case window.KeyDown:
		sel = m.nextItem(prevsel)
		m.setSelected(sel)
	case window.KeyUp:
		sel = m.prevItem(prevsel)
		m.setSelected(sel)
	case window.KeyEnter:
	default:
		return
	}
}

// setSelected sets the menu item at the specified index as selected
// and all others as not selected.
func (m *Menu) setSelected(idx int) {

	for i := 0; i < len(m.items); i++ {
		mi := m.items[i]
		if i == idx {
			mi.selected = true
		} else {
			mi.selected = false
		}
		mi.update()
	}
}

// selectedItem returns the index of the current selected menu item
// Returns -1 if no item selected
func (m *Menu) selectedItem() int {

	for i := 0; i < len(m.items); i++ {
		mi := m.items[i]
		if mi.selected {
			return i
		}
	}
	return -1
}

// nextItem returns the index of the next enabled option from the
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

// prevItem returns the index of previous enabled menu item from
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
func (mi *MenuItem) SetImage(img *Image) *MenuItem {

	return mi
}

// SetText sets the text of this menu item
func (mi *MenuItem) SetText(text string) *MenuItem {

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

func (mi *MenuItem) onCursor(evname string, ev interface{}) {

	switch evname {
	case OnCursorEnter:
		mi.selected = true
		mi.update()
	case OnCursorLeave:
		mi.selected = false
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
	if mi.disabled {
		mi.applyStyle(&mi.styles.Disabled)
		return
	}
	if mi.selected {
		mi.applyStyle(&mi.styles.Over)
		if mi.subm != nil {
			mi.menu.SetTopChild(mi)
			mi.subm.SetVisible(true)
			mi.subm.SetPosition(mi.Width()-4, 0)
		}
		return
	}
	if mi.subm != nil && !mi.subm.active {
		mi.subm.SetVisible(false)
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
