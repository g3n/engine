// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package gui

import (
	"github.com/g3n/engine/gui/assets"
	"github.com/g3n/engine/gui/assets/icon"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/text"
)

// NewLightStyle creates and returns a pointer to the a new "light" style
func NewLightStyle() *Style {

	// Fonts to use
	const fontName = "fonts/FreeSans.ttf"
	const iconName = "fonts/MaterialIcons-Regular.ttf"
	s := new(Style)

	// Creates text font
	fontData := assets.MustAsset(fontName)
	font, err := text.NewFontFromData(fontData)
	if err != nil {
		panic(err)
	}
	font.SetLineSpacing(1.0)
	font.SetSize(14)
	font.SetDPI(72)
	font.SetFgColor4(math32.NewColor4("black"))
	font.SetBgColor4(math32.NewColor4("black", 0))
	s.Font = font

	// Creates icon font
	fontIconData := assets.MustAsset(iconName)
	fontIcon, err := text.NewFontFromData(fontIconData)
	if err != nil {
		panic(err)
	}
	fontIcon.SetLineSpacing(1.0)
	fontIcon.SetSize(14)
	fontIcon.SetDPI(72)
	fontIcon.SetFgColor4(math32.NewColor4("black"))
	fontIcon.SetBgColor4(math32.NewColor4("white", 0))
	s.FontIcon = fontIcon

	zeroBounds := RectBounds{0, 0, 0, 0}
	oneBounds := RectBounds{1, 1, 1, 1}
	twoBounds := RectBounds{2, 2, 2, 2}

	borderColor := math32.Color4Name("DimGray")
	borderColorDis := math32.Color4Name("LightGray")

	bgColor := math32.Color{0.85, 0.85, 0.85}
	bgColor4 := math32.Color4{0, 0, 0, 0}
	bgColorOver := math32.Color{0.9, 0.9, 0.9}
	bgColor4Over := math32.Color4{1, 1, 1, 0.5}
	bgColor4Sel := math32.Color4{0.6, 0.6, 0.6, 1}

	fgColor := math32.Color{0, 0, 0}
	fgColorSel := math32.Color{0, 0, 0}
	fgColorDis := math32.Color{0.4, 0.4, 0.4}

	// Button styles
	s.Button = ButtonStyles{}
	s.Button.Normal = ButtonStyle{
		Border:      oneBounds,
		Paddings:    RectBounds{2, 4, 2, 4},
		BorderColor: borderColor,
		BgColor:     bgColor,
		FgColor:     fgColor,
	}
	s.Button.Over = s.Button.Normal
	s.Button.Over.BgColor = bgColorOver
	s.Button.Focus = s.Button.Over
	s.Button.Pressed = s.Button.Over
	s.Button.Pressed.Border = twoBounds
	s.Button.Disabled = s.Button.Normal
	s.Button.Disabled.BorderColor = borderColorDis
	s.Button.Disabled.FgColor = fgColorDis

	// CheckRadio styles
	s.CheckRadio = CheckRadioStyles{}
	s.CheckRadio.Normal = CheckRadioStyle{
		Border:      zeroBounds,
		Paddings:    zeroBounds,
		BorderColor: borderColor,
		BgColor:     bgColor4,
		FgColor:     fgColor,
	}
	s.CheckRadio.Over = s.CheckRadio.Normal
	s.CheckRadio.Over.BgColor = bgColor4Over
	s.CheckRadio.Focus = s.CheckRadio.Over
	s.CheckRadio.Disabled = s.CheckRadio.Normal
	s.CheckRadio.Disabled.FgColor = fgColorDis

	// Edit styles
	s.Edit = EditStyles{}
	s.Edit.Normal = EditStyle{
		Border:      oneBounds,
		Paddings:    zeroBounds,
		BorderColor: borderColor,
		BgColor:     bgColor,
		BgAlpha:     1.0,
		FgColor:     fgColor,
		HolderColor: math32.Color{0.4, 0.4, 0.4},
	}
	s.Edit.Over = s.Edit.Normal
	s.Edit.Over.BgColor = bgColorOver
	s.Edit.Focus = s.Edit.Over
	s.Edit.Disabled = s.Edit.Normal
	s.Edit.Disabled.FgColor = fgColorDis

	// ScrollBar styles
	s.ScrollBar = ScrollBarStyles{}
	s.ScrollBar.Normal = ScrollBarStyle{
		Paddings:     oneBounds,
		Borders:      oneBounds,
		BordersColor: borderColor,
		Color:        math32.Color{0.8, 0.8, 0.8},
		Button: ScrollBarButtonStyle{
			Borders:      oneBounds,
			BordersColor: borderColor,
			Color:        math32.Color{0.5, 0.5, 0.5},
			Size:         30,
		},
	}
	s.ScrollBar.Over = s.ScrollBar.Normal
	s.ScrollBar.Disabled = s.ScrollBar.Normal

	// Slider styles
	s.Slider = SliderStyles{}
	s.Slider.Normal = SliderStyle{
		Border:      oneBounds,
		BorderColor: borderColor,
		Paddings:    zeroBounds,
		BgColor:     math32.Color4{0.8, 0.8, 0.8, 1},
		FgColor:     math32.Color4{0, 0.8, 0, 1},
	}
	s.Slider.Over = s.Slider.Normal
	s.Slider.Over.BgColor = math32.Color4{1, 1, 1, 1}
	s.Slider.Over.FgColor = math32.Color4{0, 1, 0, 1}
	s.Slider.Focus = s.Slider.Over
	s.Slider.Disabled = s.Slider.Normal

	// Splitter styles
	s.Splitter = SplitterStyles{}
	s.Splitter.Normal = SplitterStyle{
		SpacerBorderColor: borderColor,
		SpacerColor:       bgColor,
		SpacerSize:        6,
	}
	s.Splitter.Over = s.Splitter.Normal
	s.Splitter.Over.SpacerColor = bgColorOver
	s.Splitter.Drag = s.Splitter.Over

	// Window styles
	s.Window = WindowStyles{}
	s.Window.Normal = WindowStyle{
		Border:           RectBounds{4, 4, 4, 4},
		Paddings:         zeroBounds,
		BorderColor:      math32.Color4{0.2, 0.2, 0.2, 1},
		TitleBorders:     RectBounds{0, 0, 1, 0},
		TitleBorderColor: math32.Color4{0, 0, 0, 1},
		TitleBgColor:     math32.Color4{0, 1, 0, 1},
		TitleFgColor:     math32.Color4{0, 0, 0, 1},
	}
	s.Window.Over = s.Window.Normal
	s.Window.Focus = s.Window.Normal
	s.Window.Disabled = s.Window.Normal

	// Scroller styles
	s.Scroller = ScrollerStyles{}
	s.Scroller.Normal = ScrollerStyle{
		Border:      oneBounds,
		Paddings:    zeroBounds,
		BorderColor: borderColor,
		BgColor:     bgColor,
		FgColor:     fgColor,
	}
	s.Scroller.Over = s.Scroller.Normal
	s.Scroller.Over.BgColor = bgColorOver
	s.Scroller.Focus = s.Scroller.Over
	s.Scroller.Disabled = s.Scroller.Normal

	// List styles
	s.List = ListStyles{}
	s.List.Scroller = &s.Scroller
	s.List.Item = &ListItemStyles{}
	s.List.Item.Normal = ListItemStyle{
		Border:      RectBounds{1, 0, 1, 0},
		Paddings:    RectBounds{0, 0, 0, 2},
		BorderColor: math32.Color4{0, 0, 0, 0},
		BgColor:     bgColor4,
		FgColor:     fgColor,
	}
	s.List.Item.Selected = s.List.Item.Normal
	s.List.Item.Selected.BgColor = bgColor4Sel
	s.List.Item.Selected.FgColor = fgColorSel
	s.List.Item.Highlighted = s.List.Item.Normal
	s.List.Item.Highlighted.BorderColor = math32.Color4{0, 0, 0, 1}
	s.List.Item.Highlighted.BgColor = bgColor4Over
	s.List.Item.Highlighted.FgColor = fgColor
	s.List.Item.SelHigh = s.List.Item.Highlighted
	s.List.Item.SelHigh.BgColor = bgColor4Sel
	s.List.Item.SelHigh.FgColor = fgColorSel

	// DropDown styles
	s.DropDown = DropDownStyles{}
	s.DropDown.Normal = DropDownStyle{
		Border:      oneBounds,
		Paddings:    RectBounds{0, 0, 0, 2},
		BorderColor: borderColor,
		BgColor:     bgColor,
		FgColor:     fgColor,
	}
	s.DropDown.Over = s.DropDown.Normal
	s.DropDown.Over.BgColor = bgColorOver
	s.DropDown.Focus = s.DropDown.Over
	s.DropDown.Disabled = s.DropDown.Normal

	// Folder styles
	s.Folder = FolderStyles{}
	s.Folder.Normal = FolderStyle{
		Margins:     zeroBounds,
		Border:      oneBounds,
		Paddings:    RectBounds{2, 0, 2, 2},
		BorderColor: borderColor,
		BgColor:     bgColor,
		FgColor:     fgColor,
		Icons:       [2]string{icon.ExpandMore, icon.ExpandLess},
	}
	s.Folder.Over = s.Folder.Normal
	s.Folder.Over.BgColor = bgColorOver
	s.Folder.Focus = s.Folder.Over
	s.Folder.Focus.Paddings = twoBounds
	s.Folder.Disabled = s.Folder.Focus

	// Tree styles
	s.Tree = TreeStyles{}
	s.Tree.Padlevel = 16.0
	s.Tree.List = &s.List
	s.Tree.Node = &TreeNodeStyles{}
	s.Tree.Node.Normal = TreeNodeStyle{
		Border:      zeroBounds,
		Paddings:    zeroBounds,
		BorderColor: borderColor,
		BgColor:     bgColor4,
		FgColor:     fgColor,
		Icons:       [2]string{icon.ExpandMore, icon.ExpandLess},
	}

	// ControlFolder styles
	s.ControlFolder = ControlFolderStyles{}
	s.ControlFolder.Folder = &FolderStyles{}
	s.ControlFolder.Folder.Normal = s.Folder.Normal
	s.ControlFolder.Folder.Normal.BorderColor = math32.Color4{0, 0, 0, 0}
	s.ControlFolder.Folder.Normal.BgColor = math32.Color{0, 0.5, 1}
	s.ControlFolder.Folder.Over = s.ControlFolder.Folder.Normal
	s.ControlFolder.Folder.Focus = s.ControlFolder.Folder.Normal
	s.ControlFolder.Folder.Focus.Paddings = twoBounds
	s.ControlFolder.Folder.Disabled = s.ControlFolder.Folder.Focus
	s.ControlFolder.Tree = &TreeStyles{}
	s.ControlFolder.Tree.Padlevel = 2.0
	s.ControlFolder.Tree.List = &ListStyles{}
	scrollerStylesCopy := *s.List.Scroller
	s.ControlFolder.Tree.List.Scroller = &scrollerStylesCopy
	s.ControlFolder.Tree.List.Scroller.Normal.Paddings = RectBounds{0, 2, 0, 0}
	s.ControlFolder.Tree.List.Scroller.Over.Paddings = RectBounds{0, 2, 0, 0}
	s.ControlFolder.Tree.List.Scroller.Focus.Paddings = RectBounds{0, 2, 0, 0}
	s.ControlFolder.Tree.List.Scroller.Disabled.Paddings = RectBounds{0, 2, 0, 0}
	s.ControlFolder.Tree.List.Item = s.List.Item
	s.ControlFolder.Tree.Node = &TreeNodeStyles{}
	s.ControlFolder.Tree.Node.Normal = s.Tree.Node.Normal

	// Menu styles
	s.Menu = MenuStyles{}
	s.Menu.Body = &MenuBodyStyles{}
	s.Menu.Body.Normal = MenuBodyStyle{
		Border:      oneBounds,
		Paddings:    twoBounds,
		BorderColor: borderColor,
		BgColor:     bgColor,
		FgColor:     fgColor,
	}
	s.Menu.Body.Over = s.Menu.Body.Normal
	s.Menu.Body.Over.BgColor = bgColorOver
	s.Menu.Body.Focus = s.Menu.Body.Over
	s.Menu.Body.Disabled = s.Menu.Body.Normal
	s.Menu.Item = &MenuItemStyles{}
	s.Menu.Item.Normal = MenuItemStyle{
		Border:           zeroBounds,
		Paddings:         RectBounds{2, 4, 2, 2},
		BorderColor:      borderColor,
		BgColor:          bgColor,
		FgColor:          fgColor,
		IconPaddings:     RectBounds{0, 6, 0, 4},
		ShortcutPaddings: RectBounds{0, 0, 0, 10},
		RiconPaddings:    RectBounds{2, 0, 0, 4},
	}
	s.Menu.Item.Over = s.Menu.Item.Normal
	s.Menu.Item.Over.BgColor = math32.Color{0.6, 0.6, 0.6}
	s.Menu.Item.Disabled = s.Menu.Item.Normal
	s.Menu.Item.Disabled.FgColor = fgColorDis
	s.Menu.Item.Separator = MenuItemStyle{
		Border:      twoBounds,
		Paddings:    zeroBounds,
		BorderColor: math32.Color4{0, 0, 0, 0},
		BgColor:     math32.Color{0.6, 0.6, 0.6},
		FgColor:     fgColor,
	}

	// Table styles
	s.Table = TableStyles{}
	s.Table.Header = TableHeaderStyle{
		Border:      RectBounds{0, 1, 1, 0},
		Paddings:    twoBounds,
		BorderColor: borderColor,
		BgColor:     math32.Color{0.7, 0.7, 0.7},
		FgColor:     fgColor,
	}
	s.Table.RowEven = TableRowStyle{
		Border:      RectBounds{0, 1, 1, 0},
		Paddings:    twoBounds,
		BorderColor: math32.Color4{0.6, 0.6, 0.6, 1},
		BgColor:     math32.Color{0.90, 0.90, 0.90},
		FgColor:     fgColor,
	}
	s.Table.RowOdd = s.Table.RowEven
	s.Table.RowOdd.BgColor = math32.Color{0.88, 0.88, 0.88}
	s.Table.RowCursor = s.Table.RowEven
	s.Table.RowCursor.BgColor = math32.Color{0.75, 0.75, 0.75}
	s.Table.RowSel = s.Table.RowEven
	s.Table.RowSel.BgColor = math32.Color{0.70, 0.70, 0.70}
	s.Table.Status = TableStatusStyle{
		Border:      RectBounds{1, 0, 0, 0},
		Paddings:    twoBounds,
		BorderColor: borderColor,
		BgColor:     math32.Color{0.9, 0.9, 0.9},
		FgColor:     fgColor,
	}
	s.Table.Resizer = TableResizerStyle{
		Width:       4,
		Border:      zeroBounds,
		BorderColor: borderColor,
		BgColor:     math32.Color4{0.4, 0.4, 0.4, 0.6},
	}

	// ImageButton styles
	s.ImageButton = ImageButtonStyles{}
	s.ImageButton.Normal = ImageButtonStyle{
		Border:      oneBounds,
		Paddings:    zeroBounds,
		BorderColor: borderColor,
		BgColor:     bgColor4,
		FgColor:     fgColor,
	}
	s.ImageButton.Over = s.ImageButton.Normal
	s.ImageButton.Over.BgColor = bgColor4Over
	s.ImageButton.Focus = s.ImageButton.Over
	s.ImageButton.Pressed = s.ImageButton.Over
	s.ImageButton.Disabled = s.ImageButton.Normal
	s.ImageButton.Disabled.FgColor = fgColorDis

	// TabBar styles
	s.TabBar = TabBarStyles{
		SepHeight:          1,
		ListButtonIcon:     icon.MoreVert,
		ListButtonPaddings: RectBounds{2, 4, 0, 0},
	}
	s.TabBar.Normal = TabBarStyle{
		Border:      oneBounds,
		Paddings:    RectBounds{2, 0, 0, 0},
		BorderColor: borderColor,
		BgColor:     math32.Color4{0.7, 0.7, 0.7, 1},
	}
	s.TabBar.Over = s.TabBar.Normal
	s.TabBar.Over.BgColor = bgColor4Over
	s.TabBar.Focus = s.TabBar.Normal
	s.TabBar.Focus.BgColor = bgColor4
	s.TabBar.Disabled = s.TabBar.Focus
	s.TabBar.Tab = TabStyles{
		IconPaddings:  RectBounds{2, 2, 0, 0},
		ImagePaddings: RectBounds{0, 2, 0, 0},
		IconClose:     icon.Clear,
	}
	s.TabBar.Tab.Normal = TabStyle{
		Margins:     RectBounds{0, 2, 0, 2},
		Border:      RectBounds{1, 1, 0, 1},
		Paddings:    twoBounds,
		BorderColor: borderColor,
		BgColor:     bgColor4,
		FgColor:     fgColor,
	}
	s.TabBar.Tab.Over = s.TabBar.Tab.Normal
	s.TabBar.Tab.Over.BgColor = bgColor4Over
	s.TabBar.Tab.Focus = s.TabBar.Tab.Normal
	s.TabBar.Tab.Focus.BgColor = bgColor4
	s.TabBar.Tab.Disabled = s.TabBar.Tab.Focus
	s.TabBar.Tab.Selected = s.TabBar.Tab.Normal
	s.TabBar.Tab.Selected.BgColor = math32.Color4{0.85, 0.85, 0.85, 1}

	return s
}
