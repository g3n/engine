// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/text"
)

// Style contains the styles for all GUI elements
type Style struct {
	Color         ColorStyle
	Font          *text.Font
	FontIcon      *text.Font
	Label         LabelStyle
	Button        ButtonStyles
	CheckRadio    CheckRadioStyles
	Edit          EditStyles
	ScrollBar     ScrollBarStyles
	Slider        SliderStyles
	Splitter      SplitterStyles
	Window        WindowStyles
	ItemScroller  ItemScrollerStyles
	Scroller      ScrollerStyle
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

// ColorStyle defines the main colors used.
type ColorStyle struct {
	BgDark    math32.Color4
	BgMed     math32.Color4
	BgNormal  math32.Color4
	BgOver    math32.Color4
	Highlight math32.Color4
	Select    math32.Color4
	Text      math32.Color4
	TextDis   math32.Color4
}

// States that a GUI element can be in
const (
	StyleOver = iota + 1
	StyleFocus
	StyleDisabled
	StyleNormal
	StyleDef
)
