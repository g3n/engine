// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/g3n/engine/gui/assets/icon"
	"github.com/g3n/engine/window"
)

// DropDown represents a dropdown GUI element.
type DropDown struct {
	Panel                        // Embedded panel
	icon         *Label          // internal label with icon
	list         *List           // internal list
	styles       *DropDownStyles // pointer to dropdown styles
	litem        *ImageLabel     // Item shown in drop box (copy of selected)
	selItem      *ImageLabel     // selected item from list
	overDropdown bool
	overList     bool
	focus        bool
	clickOut     bool
}

// DropDownStyle contains the styling of a DropDown.
type DropDownStyle BasicStyle

// DropDownStyles contains a DropDownStyle for each valid GUI state.
type DropDownStyles struct {
	Normal   DropDownStyle
	Over     DropDownStyle
	Focus    DropDownStyle
	Disabled DropDownStyle
}

// NewDropDown creates and returns a pointer to a new drop down widget with the specified width.
func NewDropDown(width float32, item *ImageLabel) *DropDown {

	dd := new(DropDown)
	dd.styles = &StyleDefault().DropDown
	dd.litem = item

	dd.Panel.Initialize(width, 0)
	dd.Panel.Subscribe(OnKeyDown, dd.onKeyEvent)
	dd.Panel.Subscribe(OnMouseDown, dd.onMouse)
	dd.Panel.Subscribe(OnCursorEnter, dd.onCursor)
	dd.Panel.Subscribe(OnCursorLeave, dd.onCursor)
	dd.Panel.Subscribe(OnResize, func(name string, ev interface{}) { dd.recalc() })

	// ListItem
	dd.Panel.Add(dd.litem)

	// Create icon
	dd.icon = NewIcon(" ")
	dd.icon.SetFontSize(StyleDefault().Label.PointSize * 1.3)
	dd.icon.SetText(string(icon.ArrowDropDown))
	dd.Panel.Add(dd.icon)

	/// Create list
	dd.list = NewVList(0, 0)
	dd.list.bounded = false
	dd.list.dropdown = true
	dd.list.SetVisible(false)

	dd.list.Subscribe(OnMouseDown, dd.onListMouse)
	dd.list.Subscribe(OnMouseOut, dd.onListMouse)
	dd.list.Subscribe(OnChange, dd.onListChangeEvent)
	dd.list.Subscribe(OnCursor, func(evname string, ev interface{}) { dd.root.StopPropagation(StopAll) })
	dd.Panel.Add(dd.list)

	dd.update()
	// This will trigger recalc()
	dd.Panel.SetContentHeight(item.Height())
	return dd
}

// Add adds a list item at the end of the list
func (dd *DropDown) Add(item *ImageLabel) {

	dd.list.Add(item)
}

// InsertAt inserts a list item at the specified position
// Returs true if the item was successfully inserted
func (dd *DropDown) InsertAt(pos int, item *ImageLabel) {

	dd.list.InsertAt(pos, item)
}

// RemoveAt removes the list item from the specified position
// Returs true if the item was successfully removed
func (dd *DropDown) RemoveAt(pos int) {

	dd.list.RemoveAt(pos)
}

// ItemAt returns the list item at the specified position
func (dd *DropDown) ItemAt(pos int) *ImageLabel {

	return dd.list.ItemAt(pos).(*ImageLabel)
}

// Len returns the number of items in the dropdown's list.
func (dd *DropDown) Len() int {

	return dd.list.Len()
}

// Selected returns the currently selected item or nil if no item was selected
func (dd *DropDown) Selected() *ImageLabel {

	return dd.selItem
}

// SelectedPos returns the currently selected position or -1 if no item was selected
func (dd *DropDown) SelectedPos() int {
	return dd.list.selected()
}

// SetSelected sets the selected item
func (dd *DropDown) SetSelected(item *ImageLabel) {
	dd.list.SetSelected(dd.selItem, false)
	dd.list.SetSelected(item, true)
	dd.copySelected()
	dd.update()
}

// SelectPos selects the item at the specified position
func (dd *DropDown) SelectPos(pos int) {
    dd.list.SetSelected(dd.selItem, false)
	dd.list.SelectPos(pos, true)
    dd.Dispatch(OnChange, nil)
}

// SetStyles sets the drop down styles overriding the default style
func (dd *DropDown) SetStyles(dds *DropDownStyles) {

	dd.styles = dds
	dd.update()
}

// onKeyEvent is called when key event is received when this dropdown has the key focus.
func (dd *DropDown) onKeyEvent(evname string, ev interface{}) {

	kev := ev.(*window.KeyEvent)
	switch kev.Keycode {
	case window.KeyF1:
		if dd.list.Visible() {
			dd.list.SetVisible(false)
		}
	default:
		return
	}
}

// onMouse receives subscribed mouse events over the dropdown
func (dd *DropDown) onMouse(evname string, ev interface{}) {

	if evname == OnMouseDown {
		// If clickOut list already closed
		if dd.clickOut {
			dd.clickOut = false
			return
		}
		dd.list.SetVisible(true)
		dd.root.SetKeyFocus(dd.list)
		return
	}
}

// onCursor receives subscribed cursor events over the dropdown
func (dd *DropDown) onCursor(evname string, ev interface{}) {

	if evname == OnCursorEnter {
		dd.overDropdown = true
	}
	if evname == OnCursorLeave {
		dd.overDropdown = false
	}
	dd.update()
	dd.root.StopPropagation(StopAll)
}

// onListMouseEvent receives mouse events over the list
func (dd *DropDown) onListMouse(evname string, ev interface{}) {

	mev := ev.(*window.MouseEvent)
	// List was clicked
	if evname == OnMouseDown {
		// If click occurred inside the list scrollbar ignore it
		if dd.list.vscroll != nil {
			if dd.list.vscroll.InsideBorders(mev.Xpos, mev.Ypos) {
				return
			}
		}
		// Otherwise, closes the list
		dd.list.SetVisible(false)
		//dd.copySelected()
		dd.overList = false
		dd.update()
		return
	}
	// Hide list when clicked out
	if evname == OnMouseOut {
		if dd.list.Visible() {
			dd.list.SetVisible(false)
		}
		// If list clickout occurred inside the dropdown, set 'clickOut' to
		// indicate that the list was already closed
		if dd.Panel.InsideBorders(mev.Xpos, mev.Ypos) {
			dd.clickOut = true
		}
	}
}

// onListCursor receives subscribed events over the list
//func (dd *DropDown) onListCursor(evname string, ev interface{}) {
//
//	if evname == OnCursorEnter {
//		dd.overList = true
//		dd.update()
//		return
//	}
//	if evname == OnCursorLeave {
//		dd.overList = false
//		dd.update()
//		return
//	}
//	dd.root.StopPropagation(StopAll)
//}

// copySelected copy to the dropdown panel the selected item
// from the list.
func (dd *DropDown) copySelected() {
    selected := dd.list.Selected()
    if len(selected) > 0 {
        dd.selItem = selected[0].(*ImageLabel)
        dd.litem.CopyFields(dd.selItem)
        dd.litem.SetWidth(dd.selItem.Width())
        dd.recalc()
        dd.Dispatch(OnChange, nil)
    } else {
        return
    }
}

// onListChangeEvent is called when an item in the list is selected
func (dd *DropDown) onListChangeEvent(evname string, ev interface{}) {

	dd.copySelected()
}

// recalc recalculates the dimensions and positions of the dropdown
// panel, children and list
func (dd *DropDown) recalc() {

	// Dropdown icon position
	posx := dd.Panel.ContentWidth() - dd.icon.Width()
	dd.icon.SetPosition(posx, 0)

	// List item position and width
	ipan := dd.litem.GetPanel()
	ipan.SetPosition(0, 0)
	height := ipan.Height()

	// List position
	dd.list.SetWidth(dd.Panel.Width())
	dd.list.SetHeight(6*height + 1)
	dd.list.SetPositionX(0)
	dd.list.SetPositionY(dd.Panel.Height())
}

// update updates the visual state
func (dd *DropDown) update() {

	if dd.overDropdown || dd.overList {
		dd.applyStyle(&dd.styles.Over)
		dd.list.ApplyStyle(StyleOver)
		return
	}
	if dd.focus {
		dd.applyStyle(&dd.styles.Focus)
		dd.list.ApplyStyle(StyleFocus)
		return
	}
	dd.applyStyle(&dd.styles.Normal)
	dd.list.ApplyStyle(StyleNormal)
}

// applyStyle applies the specified style
func (dd *DropDown) applyStyle(s *DropDownStyle) {

	dd.Panel.ApplyStyle(&s.PanelStyle)
}
