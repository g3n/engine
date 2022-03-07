// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/g3n/engine/window"
)

// Core events sent by the GUI manager.
// The target panel is the panel immediately under the mouse cursor.
const (
	// Events sent to target panel's lowest subscribed ancestor
	OnMouseDown = window.OnMouseDown // Any mouse button is pressed
	OnMouseUp   = window.OnMouseUp   // Any mouse button is released
	OnScroll    = window.OnScroll    // Scrolling mouse wheel

	// Events sent to all panels except the ancestors of the target panel
	OnMouseDownOut = "gui.OnMouseDownOut" // Any mouse button is pressed
	OnMouseUpOut   = "gui.OnMouseUpOut"   // Any mouse button is released

	// Event sent to new target panel and all of its ancestors up to (not including) the common ancestor of the new and old targets
	OnCursorEnter = "gui.OnCursorEnter" // Cursor entered the panel or a descendant
	// Event sent to old target panel and all of its ancestors up to (not including) the common ancestor of the new and old targets
	OnCursorLeave = "gui.OnCursorLeave" // Cursor left the panel or a descendant
	// Event sent to the cursor-focused IDispatcher if any, else sent to target panel's lowest subscribed ancestor
	OnCursor = window.OnCursor // Cursor is over the panel

	// Event sent to the new key-focused IDispatcher, specified on a call to gui.Manager().SetKeyFocus(core.IDispatcher)
	OnFocus = "gui.OnFocus" // All keyboard events will be exclusively sent to the receiving IDispatcher
	// Event sent to the previous key-focused IDispatcher when another panel is key-focused
	OnFocusLost = "gui.OnFocusLost" // Keyboard events will stop being sent to the receiving IDispatcher

	// Events sent to the key-focused IDispatcher
	OnKeyDown   = window.OnKeyDown   // A key is pressed
	OnKeyUp     = window.OnKeyUp     // A key is released
	OnKeyRepeat = window.OnKeyRepeat // A key was pressed and is now automatically repeating
	OnChar      = window.OnChar      // A unicode key is pressed
)

const (
	OnResize     = "gui.OnResize"     // Panel size changed (no parameters)
	OnEnable     = "gui.OnEnable"     // Panel enabled/disabled (no parameters)
	OnClick      = "gui.OnClick"      // Widget clicked by mouse left button or via key press
	OnChange     = "gui.OnChange"     // Value was changed. Emitted by List, DropDownList, CheckBox and Edit
	OnRadioGroup = "gui.OnRadioGroup" // Radio button within a group changed state
)
