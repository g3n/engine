// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/g3n/engine/window"
)

// List represents a list GUI element
type List struct {
	ItemScroller             // Embedded scroller
	styles       *ListStyles // Pointer to styles
	single       bool        // Single selection flag (default is true)
	focus        bool        // has keyboard focus
	dropdown     bool        // this is used as dropdown
	keyNext      window.Key  // Code of key to select next item
	keyPrev      window.Key  // Code of key to select previous item
}

// ListItem encapsulates each item inserted into the list
type ListItem struct {
	Panel               // Container panel
	item        IPanel  // Original item
	selected    bool    // Item selected flag
	highlighted bool    // Item highlighted flag
	padLeft     float32 // Additional left padding
	list        *List   // Pointer to list
}

// ListStyles encapsulates a set of styles for the list and item.
type ListStyles struct {
	Scroller *ItemScrollerStyles
	Item     *ListItemStyles
}

// ListItemStyles contains one ListItemStyle for each possible item state.
type ListItemStyles struct {
	Normal      ListItemStyle
	Over        ListItemStyle
	Selected    ListItemStyle
	Highlighted ListItemStyle
	SelHigh     ListItemStyle
}

// ListItemStyle contains the styling of a list item.
type ListItemStyle BasicStyle

// OnListItemResize is the identifier of the event dispatched when a ListItem's child panel is resized
const OnListItemResize = "gui.OnListItemResize"

// NewVList creates and returns a pointer to a new vertical list panel
// with the specified dimensions
func NewVList(width, height float32) *List {

	return newList(true, width, height)
}

// NewHList creates and returns a pointer to a new horizontal list panel
// with the specified dimensions
func NewHList(width, height float32) *List {

	return newList(false, width, height)
}

// newList creates and returns a pointer to a new list panel
// with the specified orientation and dimensions
func newList(vert bool, width, height float32) *List {

	li := new(List)
	li.initialize(vert, width, height)
	return li
}

func (li *List) initialize(vert bool, width, height float32) {

	li.styles = &StyleDefault().List
	li.single = true

	li.ItemScroller.initialize(vert, width, height)
	li.ItemScroller.SetStyles(li.styles.Scroller)
	li.ItemScroller.adjustItem = true
	li.ItemScroller.Subscribe(OnMouseDown, li.onMouseEvent)
	li.ItemScroller.Subscribe(OnKeyDown, li.onKeyEvent)
	li.ItemScroller.Subscribe(OnKeyRepeat, li.onKeyEvent)

	if vert {
		li.keyNext = window.KeyDown
		li.keyPrev = window.KeyUp
	} else {
		li.keyNext = window.KeyRight
		li.keyPrev = window.KeyLeft
	}

	li.update()
}

// SetSingle sets the single/multiple selection flag of the list
func (li *List) SetSingle(state bool) {

	li.single = state
}

// Single returns the current state of the single/multiple selection flag
func (li *List) Single() bool {

	return li.single
}

// SetStyles set the listr styles overriding the default style
func (li *List) SetStyles(s *ListStyles) {

	li.styles = s
	li.ItemScroller.SetStyles(li.styles.Scroller)
	li.update()
}

// Add add a list item at the end of the list
func (li *List) Add(item IPanel) {

	li.InsertAt(len(li.items), item)
}

// InsertAt inserts a list item at the specified position
// Returs true if the item was successfully inserted
func (li *List) InsertAt(pos int, item IPanel) {

	litem := newListItem(li, item)
	li.ItemScroller.InsertAt(pos, litem)
	litem.Panel.Subscribe(OnMouseDown, litem.onMouse)
	litem.Panel.Subscribe(OnCursorEnter, litem.onCursor)
}

// RemoveAt removes the list item from the specified position
func (li *List) RemoveAt(pos int) IPanel {

	// Remove the list item from the internal scroller
	pan := li.ItemScroller.RemoveAt(pos)
	litem := pan.(*ListItem)

	// Remove item from the list item children and disposes of the list item panel
	item := litem.item
	litem.Remove(item)
	litem.Dispose()
	return item
}

// Remove removes the specified item from the list
func (li *List) Remove(item IPanel) {

	for p, curr := range li.items {
		if curr.(*ListItem).item == item {
			li.RemoveAt(p)
			return
		}
	}
}

// ItemAt returns the list item at the specified position
func (li *List) ItemAt(pos int) IPanel {

	item := li.ItemScroller.ItemAt(pos)
	if item == nil {
		return nil
	}
	litem := item.(*ListItem)
	return litem.item
}

// ItemPosition returns the position of the specified item in
// the list or -1 if not found
func (li *List) ItemPosition(item IPanel) int {

	for pos := 0; pos < len(li.items); pos++ {
		if li.items[pos].(*ListItem).item == item {
			return pos
		}
	}
	return -1
}

// Selected returns list with the currently selected items
func (li *List) Selected() []IPanel {

	sel := []IPanel{}
	for _, item := range li.items {
		litem := item.(*ListItem)
		if litem.selected {
			sel = append(sel, litem.item)
		}
	}
	return sel
}

// SetSelected selects or unselects the specified item
func (li *List) SetSelected(item IPanel, state bool) {

	for _, curr := range li.items {
		litem := curr.(*ListItem)
		if litem.item == item {
			litem.SetSelected(state)
			li.update()
			li.Dispatch(OnChange, nil)
			return
		}
	}
}

// SelectPos selects or unselects the item at the specified position
func (li *List) SelectPos(pos int, state bool) {

	if pos < 0 || pos >= len(li.items) {
		return
	}
	litem := li.items[pos].(*ListItem)
	if litem.selected == state {
		return
	}
	litem.SetSelected(state)
	li.update()
	li.Dispatch(OnChange, nil)
}

// SetItemPadLeftAt sets the additional left padding for this item
// It is used mainly by the tree control
func (li *List) SetItemPadLeftAt(pos int, pad float32) {

	if pos < 0 || pos >= len(li.items) {
		return
	}
	litem := li.items[pos].(*ListItem)
	litem.padLeft = pad
	litem.update()
}

// selNext selects or highlights the next item, if possible
func (li *List) selNext(sel bool, update bool) *ListItem {

	// Checks for empty list
	if len(li.items) == 0 {
		return nil
	}
	// Find currently selected item
	var pos int
	if sel {
		pos = li.selected()
	} else {
		pos = li.highlighted()
	}

	var newItem *ListItem
	newSel := true
	// If no item found, returns first.
	if pos < 0 {
		newItem = li.items[0].(*ListItem)
		if sel {
			newItem.SetSelected(true)
		} else {
			newItem.SetHighlighted(true)
		}
	} else {
		item := li.items[pos].(*ListItem)
		// Item is not the last, get next
		if pos < len(li.items)-1 {
			newItem = li.items[pos+1].(*ListItem)
			if sel {
				item.SetSelected(false)
				newItem.SetSelected(true)
			} else {
				item.SetHighlighted(false)
				newItem.SetHighlighted(true)
			}
			if !li.ItemVisible(pos + 1) {
				li.ScrollDown()
			}
			// Item is the last, don't change
		} else {
			newItem = item
			newSel = false
		}
	}

	if update {
		li.update()
	}
	if sel && newSel {
		li.Dispatch(OnChange, nil)
	}
	return newItem
}

// selPrev selects or highlights the next item, if possible
func (li *List) selPrev(sel bool, update bool) *ListItem {

	// Check for empty list
	if len(li.items) == 0 {
		return nil
	}

	// Find first selected item
	var pos int
	if sel {
		pos = li.selected()
	} else {
		pos = li.highlighted()
	}

	var newItem *ListItem
	newSel := true
	// If no item found, returns first.
	if pos < 0 {
		newItem = li.items[0].(*ListItem)
		if sel {
			newItem.SetSelected(true)
		} else {
			newItem.SetHighlighted(true)
		}
	} else {
		item := li.items[pos].(*ListItem)
		if pos == 0 {
			newItem = item
			newSel = false
		} else {
			newItem = li.items[pos-1].(*ListItem)
			if sel {
				item.SetSelected(false)
				newItem.SetSelected(true)
			} else {
				item.SetHighlighted(false)
				newItem.SetHighlighted(true)
			}
			if (pos - 1) < li.first {
				li.ScrollUp()
			}
		}
	}
	if update {
		li.update()
	}
	if sel && newSel {
		li.Dispatch(OnChange, nil)
	}
	return newItem
}

// selected returns the position of first selected item
func (li *List) selected() (pos int) {

	for pos, item := range li.items {
		if item.(*ListItem).selected {
			return pos
		}
	}
	return -1
}

// highlighted returns the position of first highlighted item
func (li *List) highlighted() (pos int) {

	for pos, item := range li.items {
		if item.(*ListItem).highlighted {
			return pos
		}
	}
	return -1
}

// onMouseEvent receives subscribed mouse events for the list
func (li *List) onMouseEvent(evname string, ev interface{}) {

	li.root.SetKeyFocus(li)
	li.root.StopPropagation(StopAll)
}

// onKeyEvent receives subscribed key events for the list
func (li *List) onKeyEvent(evname string, ev interface{}) {

	kev := ev.(*window.KeyEvent)
	// Dropdown mode
	if li.dropdown {
		switch kev.Keycode {
		case li.keyNext:
			li.selNext(true, true)
		case li.keyPrev:
			li.selPrev(true, true)
		case window.KeyEnter:
			li.SetVisible(false)
		default:
			return
		}
		li.root.StopPropagation(Stop3D)
		return
	}

	// Listbox mode single selection
	if li.single {
		switch kev.Keycode {
		case li.keyNext:
			li.selNext(true, true)
		case li.keyPrev:
			li.selPrev(true, true)
		default:
			return
		}
		li.root.StopPropagation(Stop3D)
		return
	}

	// Listbox mode multiple selection
	switch kev.Keycode {
	case li.keyNext:
		li.selNext(false, true)
	case li.keyPrev:
		li.selPrev(false, true)
	case window.KeySpace:
		pos := li.highlighted()
		if pos >= 0 {
			litem := li.items[pos].(*ListItem)
			li.setSelection(litem, !litem.selected, true, true)
		}
	default:
		return
	}
	li.root.StopPropagation(Stop3D)
}

// setSelection sets the selected state of the specified item
// updating the visual appearance of the list if necessary
func (li *List) setSelection(litem *ListItem, state bool, force bool, dispatch bool) {

	// If already at this state, nothing to do
	if litem.selected == state && !force {
		return
	}
	litem.SetSelected(state)

	// If single selection, deselects all other items
	if li.single {
		for _, curr := range li.items {
			if curr.(*ListItem) != litem {
				curr.(*ListItem).SetSelected(false)
				curr.(*ListItem).SetHighlighted(false)
			}
		}
	}
	li.update()
	if dispatch {
		li.Dispatch(OnChange, nil)
	}
}

// update updates the visual state the list and its items
func (li *List) update() {

	// Update the list items styles
	for _, item := range li.items {
		item.(*ListItem).update()
	}
}

//
// ListItem methods
//

func newListItem(list *List, item IPanel) *ListItem {

	litem := new(ListItem)
	litem.Panel.Initialize(0, 0)
	litem.item = item
	litem.list = list
	litem.Panel.Add(item)
	litem.SetContentWidth(item.GetPanel().Width())
	litem.SetContentHeight(item.GetPanel().Height())
	// If this list item is resized, sends event to its child panel
	litem.Subscribe(OnResize, func(evname string, ev interface{}) {
		item.GetPanel().Dispatch(OnListItemResize, nil)
	})
	litem.update()
	return litem
}

// onMouse receives mouse button events over the list item
func (litem *ListItem) onMouse(evname string, ev interface{}) {

	if litem.list.single {
		litem.list.setSelection(litem, true, true, true)
	} else {
		litem.list.setSelection(litem, !litem.selected, true, true)
	}
}

// onCursor receives subscribed cursor events over the list item
func (litem *ListItem) onCursor(evname string, ev interface{}) {

	if litem.list.dropdown {
		litem.list.setSelection(litem, true, true, false)
		return
	}
}

// SetSelected sets this item selected state
func (litem *ListItem) SetSelected(state bool) {

	litem.selected = state
	//litem.item.SetSelected2(state)
}

// SetHighlighted sets this item selected state
func (litem *ListItem) SetHighlighted(state bool) {

	litem.highlighted = state
	//litem.item.SetHighlighted2(state)
}

// updates the list item visual style accordingly to its current state
func (litem *ListItem) update() {

	list := litem.list
	if litem.selected && !litem.highlighted {
		litem.applyStyle(&list.styles.Item.Selected)
		return
	}
	if !litem.selected && litem.highlighted {
		litem.applyStyle(&list.styles.Item.Highlighted)
		return
	}
	if litem.selected && litem.highlighted {
		litem.applyStyle(&list.styles.Item.SelHigh)
		return
	}
	litem.applyStyle(&list.styles.Item.Normal)
}

// applyStyle applies the specified style to this ListItem
func (litem *ListItem) applyStyle(s *ListItemStyle) {

	styleCopy := s.PanelStyle
	styleCopy.Padding.Left += litem.padLeft
	litem.Panel.ApplyStyle(&styleCopy)
}
