// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

//import (
//	"github.com/g3n/engine/window"
//)

/***************************************

 ButtonDestructive Panel
 +-------------------------------+
 |  Image/Icon      Label        |
 |  +----------+   +----------+  |
 |  |          |   |          |  |
 |  |          |   |          |  |
 |  +----------+   +----------+  |
 +-------------------------------+

****************************************/

// ButtonDestructive represents a destructive button GUI element
type ButtonDestructive struct {
	Button                          // Embedded Button
	styles *ButtonDestructiveStyles // pointer to current destructive button styles
}

// ButtonDestructiveStyle contains the styling of a ButtonDestructive
type ButtonDestructiveStyle ButtonStyle

// ButtonDestructiveStyles contains one ButtonDestructiveStyle for each possible button state
type ButtonDestructiveStyles struct {
	Normal       ButtonDestructiveStyle
	Over         ButtonDestructiveStyle
	PressedOnce  ButtonDestructiveStyle
	PressedTwice ButtonDestructiveStyle
	Disabled     ButtonDestructiveStyle
}
