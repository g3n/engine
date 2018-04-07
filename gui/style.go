// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/g3n/engine/text"
)

// Style contains the styles for all GUI elements
type Style struct {
	Font          *text.Font
	FontIcon      *text.Font
	Button        ButtonStyles
	CheckRadio    CheckRadioStyles
	Edit          EditStyles
	ScrollBar     ScrollBarStyles
	Slider        SliderStyles
	Splitter      SplitterStyles
	Window        WindowStyles
	Scroller      ScrollerStyles
	List          ListStyles
	DropDown      DropDownStyles
	Folder        FolderStyles
	Tree          TreeStyles
	ControlFolder ControlFolderStyles
	Menu          MenuStyles
	Table         TableStyles
	ImageButton   ImageButtonStyles
	TabBar        TabBarStyles
}

// States that a GUI element can be in
const (
	StyleOver = iota + 1
	StyleFocus
	StyleDisabled
	StyleNormal
	StyleDef
)
