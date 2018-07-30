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
	s.Font = font

	// Creates icon font
	fontIconData := assets.MustAsset(iconName)
	fontIcon, err := text.NewFontFromData(fontIconData)
	if err != nil {
		panic(err)
	}
	s.FontIcon = fontIcon

	zeroBounds := RectBounds{0, 0, 0, 0}
	oneBounds := RectBounds{1, 1, 1, 1}
	twoBounds := RectBounds{2, 2, 2, 2}

	borderColor := math32.Color4Name("DimGray")
	borderColorDis := math32.Color4Name("LightGray")

	bgColor := math32.Color4{0.85, 0.85, 0.85, 1}
	bgColor4 := math32.Color4{0, 0, 0, 0}
	bgColorOver := math32.Color4{0.9, 0.9, 0.9, 1}
	bgColor4Over := math32.Color4{1, 1, 1, 0.5}
	bgColor4Sel := math32.Color4{0.6, 0.6, 0.6, 1}

	fgColor := math32.Color4{0, 0, 0, 1}
	fgColorSel := math32.Color4{0, 0, 0, 1}
	fgColorDis := math32.Color4{0.4, 0.4, 0.4, 1}

	// Label style
	s.Label = LabelStyle{}
	s.Label.FontAttributes = text.FontAttributes{}
	s.Label.FontAttributes.PointSize = 14
	s.Label.FontAttributes.DPI = 72
	s.Label.FontAttributes.Hinting = text.HintingNone
	s.Label.FontAttributes.LineSpacing = 1.0
	s.Label.BgColor = math32.Color4{0, 0, 0, 0}
	s.Label.FgColor = math32.Color4{0, 0, 0, 1}

	// Button styles
	s.Button = ButtonStyles{}
	s.Button.Normal = ButtonStyle{}
	s.Button.Normal.Border = oneBounds
	s.Button.Normal.Padding = RectBounds{2, 4, 2, 4}
	s.Button.Normal.BorderColor = borderColor
	s.Button.Normal.BgColor = bgColor
	s.Button.Normal.FgColor = fgColor
	s.Button.Over = s.Button.Normal
	s.Button.Over.BgColor = bgColorOver
	s.Button.Focus = s.Button.Over
	s.Button.Pressed = s.Button.Over
	s.Button.Pressed.Border = RectBounds{2, 2, 2, 2}
	s.Button.Pressed.Padding = RectBounds{2, 2, 0, 4}
	s.Button.Disabled = s.Button.Normal
	s.Button.Disabled.BorderColor = borderColorDis
	s.Button.Disabled.FgColor = fgColorDis

	// CheckRadio styles
	s.CheckRadio = CheckRadioStyles{}
	s.CheckRadio.Normal = CheckRadioStyle{}
	s.CheckRadio.Normal.BorderColor = borderColor
	s.CheckRadio.Normal.BgColor = bgColor4
	s.CheckRadio.Normal.FgColor = fgColor
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
		HolderColor: math32.Color4{0.4, 0.4, 0.4, 1},
	}
	s.Edit.Over = s.Edit.Normal
	s.Edit.Over.BgColor = bgColorOver
	s.Edit.Focus = s.Edit.Over
	s.Edit.Disabled = s.Edit.Normal
	s.Edit.Disabled.FgColor = fgColorDis

	// ScrollBar styles
	s.ScrollBar = ScrollBarStyles{}
	s.ScrollBar.Normal = ScrollBarStyle{}
	s.ScrollBar.Normal.Padding = oneBounds
	s.ScrollBar.Normal.Border = oneBounds
	s.ScrollBar.Normal.BorderColor = borderColor
	s.ScrollBar.Normal.BgColor = math32.Color4{0.8, 0.8, 0.8, 1}
	s.ScrollBar.Normal.ButtonLength = 32
	s.ScrollBar.Normal.Button = PanelStyle{
		Border:      oneBounds,
		BorderColor: borderColor,
		BgColor:     math32.Color4{0.5, 0.5, 0.5, 1},
	}
	s.ScrollBar.Over = s.ScrollBar.Normal
	s.ScrollBar.Disabled = s.ScrollBar.Normal

	// Slider styles
	s.Slider = SliderStyles{}
	s.Slider.Normal = SliderStyle{}
	s.Slider.Normal.Border = oneBounds
	s.Slider.Normal.BorderColor = borderColor
	s.Slider.Normal.BgColor = math32.Color4{0.8, 0.8, 0.8, 1}
	s.Slider.Normal.FgColor = math32.Color4{0, 0.8, 0, 1}
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
	s.Window.Normal = WindowStyle{}
	s.Window.Normal.Border = RectBounds{4, 4, 4, 4}
	s.Window.Normal.Padding = zeroBounds
	s.Window.Normal.BorderColor = math32.Color4{0.2, 0.2, 0.2, 1}
	s.Window.Normal.TitleStyle = WindowTitleStyle{}
	s.Window.Normal.TitleStyle.Border = RectBounds{0, 0, 1, 0}
	s.Window.Normal.TitleStyle.BorderColor = math32.Color4{0, 0, 0, 1}
	s.Window.Normal.TitleStyle.BgColor = math32.Color4{0, 1, 0, 1} // s.Color.Select
	s.Window.Normal.TitleStyle.FgColor = math32.Color4{0, 0, 0, 1} // s.Color.Text
	s.Window.Over = s.Window.Normal
	s.Window.Focus = s.Window.Normal
	s.Window.Disabled = s.Window.Normal

	// ItemScroller styles
	s.Scroller = ScrollerStyle{}
	s.Scroller.VerticalScrollbar = ScrollerScrollbarStyle{}
	s.Scroller.VerticalScrollbar.ScrollBarStyle = s.ScrollBar.Normal
	s.Scroller.VerticalScrollbar.Broadness = 16
	s.Scroller.VerticalScrollbar.Position = ScrollbarRight
	s.Scroller.VerticalScrollbar.OverlapContent = false
	s.Scroller.VerticalScrollbar.AutoSizeButton = true
	s.Scroller.HorizontalScrollbar = s.Scroller.VerticalScrollbar
	s.Scroller.HorizontalScrollbar.Position = ScrollbarBottom
	s.Scroller.ScrollbarInterlocking = ScrollbarInterlockingNone
	s.Scroller.CornerCovered = true
	s.Scroller.CornerPanel = PanelStyle{}
	s.Scroller.CornerPanel.BgColor = math32.Color4Name("silver")
	s.Scroller.Border = oneBounds
	s.Scroller.BorderColor = borderColor
	s.Scroller.BgColor = bgColor

	// ItemScroller styles
	s.ItemScroller = ItemScrollerStyles{}
	s.ItemScroller.Normal = ItemScrollerStyle{}
	s.ItemScroller.Normal.Border = oneBounds
	s.ItemScroller.Normal.BorderColor = borderColor
	s.ItemScroller.Normal.BgColor = bgColor
	s.ItemScroller.Normal.FgColor = fgColor
	s.ItemScroller.Over = s.ItemScroller.Normal
	s.ItemScroller.Over.BgColor = bgColorOver
	s.ItemScroller.Focus = s.ItemScroller.Over
	s.ItemScroller.Disabled = s.ItemScroller.Normal

	// List styles
	s.List = ListStyles{}
	s.List.Scroller = &s.ItemScroller
	s.List.Item = &ListItemStyles{}
	s.List.Item.Normal = ListItemStyle{}
	s.List.Item.Normal.Border = RectBounds{0, 0, 1, 0}
	s.List.Item.Normal.Padding = RectBounds{0, 0, 0, 2}
	s.List.Item.Normal.BorderColor = math32.Color4{0, 0, 0, 0}
	s.List.Item.Normal.BgColor = bgColor4
	s.List.Item.Normal.FgColor = fgColor
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
	s.DropDown.Normal = DropDownStyle{}
	s.DropDown.Normal.Border = oneBounds
	s.DropDown.Normal.Padding = RectBounds{0, 0, 0, 2}
	s.DropDown.Normal.BorderColor = borderColor
	s.DropDown.Normal.BgColor = bgColor
	s.DropDown.Normal.FgColor = fgColor
	s.DropDown.Over = s.DropDown.Normal
	s.DropDown.Over.BgColor = bgColorOver
	s.DropDown.Focus = s.DropDown.Over
	s.DropDown.Disabled = s.DropDown.Normal

	// Folder styles
	s.Folder = FolderStyles{}
	s.Folder.Normal = FolderStyle{}
	s.Folder.Normal.Border = oneBounds
	s.Folder.Normal.Padding = RectBounds{2, 0, 2, 2}
	s.Folder.Normal.BorderColor = borderColor
	s.Folder.Normal.BgColor = bgColor
	s.Folder.Normal.FgColor = fgColor
	s.Folder.Normal.Icons = [2]string{icon.ExpandMore, icon.ExpandLess}
	s.Folder.Over = s.Folder.Normal
	s.Folder.Over.BgColor = bgColorOver
	s.Folder.Focus = s.Folder.Over
	s.Folder.Focus.Padding = twoBounds
	s.Folder.Disabled = s.Folder.Focus

	// Tree styles
	s.Tree = TreeStyles{}
	s.Tree.Padlevel = 16.0
	s.Tree.List = &s.List
	s.Tree.Node = &TreeNodeStyles{}
	s.Tree.Node.Normal = TreeNodeStyle{}
	s.Tree.Node.Normal.BorderColor = borderColor
	s.Tree.Node.Normal.BgColor = bgColor4
	s.Tree.Node.Normal.FgColor = fgColor
	s.Tree.Node.Normal.Icons = [2]string{icon.ExpandMore, icon.ExpandLess}

	// ControlFolder styles
	s.ControlFolder = ControlFolderStyles{}
	s.ControlFolder.Folder = &FolderStyles{}
	s.ControlFolder.Folder.Normal = s.Folder.Normal
	s.ControlFolder.Folder.Normal.BorderColor = math32.Color4{0, 0, 0, 0}
	s.ControlFolder.Folder.Normal.BgColor = math32.Color4{0, 0.5, 1, 1}
	s.ControlFolder.Folder.Over = s.ControlFolder.Folder.Normal
	s.ControlFolder.Folder.Focus = s.ControlFolder.Folder.Normal
	s.ControlFolder.Folder.Focus.Padding = twoBounds
	s.ControlFolder.Folder.Disabled = s.ControlFolder.Folder.Focus
	s.ControlFolder.Tree = &TreeStyles{}
	s.ControlFolder.Tree.Padlevel = 2.0
	s.ControlFolder.Tree.List = &ListStyles{}
	scrollerStylesCopy := *s.List.Scroller
	s.ControlFolder.Tree.List.Scroller = &scrollerStylesCopy
	s.ControlFolder.Tree.List.Scroller.Normal.Padding = RectBounds{0, 2, 0, 0}
	s.ControlFolder.Tree.List.Scroller.Over.Padding = RectBounds{0, 2, 0, 0}
	s.ControlFolder.Tree.List.Scroller.Focus.Padding = RectBounds{0, 2, 0, 0}
	s.ControlFolder.Tree.List.Scroller.Disabled.Padding = RectBounds{0, 2, 0, 0}
	s.ControlFolder.Tree.List.Item = s.List.Item
	s.ControlFolder.Tree.Node = &TreeNodeStyles{}
	s.ControlFolder.Tree.Node.Normal = s.Tree.Node.Normal

	// Menu styles
	s.Menu = MenuStyles{}
	s.Menu.Body = &MenuBodyStyles{}
	s.Menu.Body.Normal = MenuBodyStyle{}
	s.Menu.Body.Normal.Border = oneBounds
	s.Menu.Body.Normal.Padding = twoBounds
	s.Menu.Body.Normal.BorderColor = borderColor
	s.Menu.Body.Normal.BgColor = bgColor
	s.Menu.Body.Normal.FgColor = fgColor
	s.Menu.Body.Over = s.Menu.Body.Normal
	s.Menu.Body.Over.BgColor = bgColorOver
	s.Menu.Body.Focus = s.Menu.Body.Over
	s.Menu.Body.Disabled = s.Menu.Body.Normal
	s.Menu.Item = &MenuItemStyles{}
	s.Menu.Item.Normal = MenuItemStyle{}
	s.Menu.Item.Normal.Padding = RectBounds{2, 4, 2, 2}
	s.Menu.Item.Normal.BorderColor = borderColor
	s.Menu.Item.Normal.BgColor = bgColor
	s.Menu.Item.Normal.FgColor = fgColor
	s.Menu.Item.Normal.IconPaddings = RectBounds{0, 6, 0, 4}
	s.Menu.Item.Normal.ShortcutPaddings = RectBounds{0, 0, 0, 10}
	s.Menu.Item.Normal.RiconPaddings = RectBounds{2, 0, 0, 4}
	s.Menu.Item.Over = s.Menu.Item.Normal
	s.Menu.Item.Over.BgColor = math32.Color4{0.6, 0.6, 0.6, 1}
	s.Menu.Item.Disabled = s.Menu.Item.Normal
	s.Menu.Item.Disabled.FgColor = fgColorDis
	s.Menu.Item.Separator = MenuItemStyle{}
	s.Menu.Item.Separator.Border = twoBounds
	s.Menu.Item.Separator.Padding = zeroBounds
	s.Menu.Item.Separator.BorderColor = math32.Color4{0, 0, 0, 0}
	s.Menu.Item.Separator.BgColor = math32.Color4{0.6, 0.6, 0.6, 1}
	s.Menu.Item.Separator.FgColor = fgColor

	// Table styles
	s.Table = TableStyles{}
	s.Table.Header = TableHeaderStyle{}
	s.Table.Header.Border = RectBounds{0, 1, 1, 0}
	s.Table.Header.Padding = twoBounds
	s.Table.Header.BorderColor = borderColor
	s.Table.Header.BgColor = math32.Color4{0.7, 0.7, 0.7, 1}
	s.Table.Header.FgColor = fgColor
	s.Table.RowEven = TableRowStyle{}
	s.Table.RowEven.Border = RectBounds{0, 1, 1, 0}
	s.Table.RowEven.Padding = twoBounds
	s.Table.RowEven.BorderColor = math32.Color4{0.6, 0.6, 0.6, 1}
	s.Table.RowEven.BgColor = math32.Color4{0.90, 0.90, 0.90, 1}
	s.Table.RowEven.FgColor = fgColor
	s.Table.RowOdd = s.Table.RowEven
	s.Table.RowOdd.BgColor = math32.Color4{0.88, 0.88, 0.88, 1}
	s.Table.RowCursor = s.Table.RowEven
	s.Table.RowCursor.BgColor = math32.Color4{0.75, 0.75, 0.75, 1}
	s.Table.RowSel = s.Table.RowEven
	s.Table.RowSel.BgColor = math32.Color4{0.70, 0.70, 0.70, 1}
	s.Table.Status = TableStatusStyle{}
	s.Table.Status.Border = RectBounds{1, 0, 0, 0}
	s.Table.Status.Padding = twoBounds
	s.Table.Status.BorderColor = borderColor
	s.Table.Status.BgColor = math32.Color4{0.9, 0.9, 0.9, 1}
	s.Table.Status.FgColor = fgColor
	s.Table.Resizer = TableResizerStyle{
		Width:       4,
		Border:      zeroBounds,
		BorderColor: borderColor,
		BgColor:     math32.Color4{0.4, 0.4, 0.4, 0.6},
	}

	// ImageButton styles
	s.ImageButton = ImageButtonStyles{}
	s.ImageButton.Normal = ImageButtonStyle{}
	s.ImageButton.Normal.Border = oneBounds
	s.ImageButton.Normal.BorderColor = borderColor
	s.ImageButton.Normal.BgColor = bgColor4
	s.ImageButton.Normal.FgColor = fgColor
	s.ImageButton.Over = s.ImageButton.Normal
	s.ImageButton.Over.BgColor = bgColor4Over
	s.ImageButton.Focus = s.ImageButton.Over
	s.ImageButton.Pressed = s.ImageButton.Over
	s.ImageButton.Pressed.Border = twoBounds
	s.ImageButton.Disabled = s.ImageButton.Normal
	s.ImageButton.Disabled.FgColor = fgColorDis

	// TabBar styles
	s.TabBar = TabBarStyles{
		SepHeight:          1,
		ListButtonIcon:     icon.MoreVert,
		ListButtonPaddings: RectBounds{2, 4, 0, 0},
	}
	s.TabBar.Normal = TabBarStyle{}
	s.TabBar.Normal.Border = oneBounds
	s.TabBar.Normal.Padding = RectBounds{2, 0, 0, 0}
	s.TabBar.Normal.BorderColor = borderColor
	s.TabBar.Normal.BgColor = math32.Color4{0.7, 0.7, 0.7, 1}
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
	s.TabBar.Tab.Normal = TabStyle{}
	s.TabBar.Tab.Normal.Margin = RectBounds{0, 2, 0, 2}
	s.TabBar.Tab.Normal.Border = RectBounds{1, 1, 0, 1}
	s.TabBar.Tab.Normal.Padding = twoBounds
	s.TabBar.Tab.Normal.BorderColor = borderColor
	s.TabBar.Tab.Normal.BgColor = bgColor4
	s.TabBar.Tab.Normal.FgColor = fgColor
	s.TabBar.Tab.Over = s.TabBar.Tab.Normal
	s.TabBar.Tab.Over.BgColor = bgColor4Over
	s.TabBar.Tab.Focus = s.TabBar.Tab.Normal
	s.TabBar.Tab.Focus.BgColor = bgColor4
	s.TabBar.Tab.Disabled = s.TabBar.Tab.Focus
	s.TabBar.Tab.Selected = s.TabBar.Tab.Normal
	s.TabBar.Tab.Selected.BgColor = math32.Color4{0.85, 0.85, 0.85, 1}

	return s
}
