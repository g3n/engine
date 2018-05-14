// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/g3n/engine/window"
)

// Consolidate window events plus GUI events
const (
	OnClick       = "gui.OnClick"       // Widget clicked by mouse left button or key
	OnCursor      = window.OnCursor     // cursor (mouse) position events
	OnCursorEnter = "gui.OnCursorEnter" // cursor enters the panel area
	OnCursorLeave = "gui.OnCursorLeave" // cursor leaves the panel area
	OnMouseDown   = window.OnMouseDown  // any mouse button is pressed
	OnMouseUp     = window.OnMouseUp    // any mouse button is released
	OnMouseOut    = "gui.OnMouseOut"    // mouse button pressed outside of the panel
	OnKeyDown     = window.OnKeyDown    // key is pressed
	OnKeyUp       = window.OnKeyUp      // key is released
	OnKeyRepeat   = window.OnKeyRepeat  // key is pressed again by automatic repeat
	OnChar        = window.OnChar       // key is pressed and has unicode
	OnResize      = "gui.OnResize"      // panel size changed (no parameters)
	OnEnable      = "gui.OnEnable"      // panel enabled state changed (no parameters)
	OnChange      = "gui.OnChange"      // onChange is emitted by List, DropDownList, CheckBox and Edit
	OnScroll      = window.OnScroll     // scroll event
	OnChild       = "gui.OnChild"       // child added to or removed from panel
	OnRadioGroup  = "gui.OnRadioGroup"  // radio button from a group changed state
	OnRightClick  = "gui.OnRightClick"  // Widget clicked by mouse right button
)
