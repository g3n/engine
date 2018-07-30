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

// NewDarkStyle creates and returns a pointer to the a new "dark" style
func NewDarkStyle() *Style {

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

	s.Color.BgDark = math32.Color4{43.0 / 256.0, 43.0 / 256.0, 43.0 / 256.0, 1}
	s.Color.BgMed = math32.Color4{49.0 / 256.0, 51.0 / 256.0, 53.0 / 256.0, 1}
	s.Color.BgNormal = math32.Color4{60.0 / 256.0, 63.0 / 256.0, 65.0 / 256.0, 1}
	s.Color.BgOver = math32.Color4{70.0 / 256.0, 74.0 / 256.0, 77.0 / 256.0, 1}
	s.Color.Highlight = math32.Color4{75.0 / 256.0, 110.0 / 256.0, 175.0 / 256.0, 1}
	s.Color.Select = math32.Color4{13.0 / 256.0, 41.0 / 256.0, 62.0 / 256.0, 1}
	s.Color.Text = math32.Color4{1, 1, 1, 1}
	s.Color.TextDis = math32.Color4{0.4, 0.4, 0.4, 1}

	borderColor := s.Color.BgDark
	transparent := math32.Color4{0, 0, 0, 0}

	//bgColorBlue := math32.Color4{59.0/256.0, 71.0/256.0, 84.0/256.0, 1}
	//bgColorBlue2 := math32.Color4{59.0/256.0, 71.0/256.0, 120.0/256.0, 1}
	//bgColorBlueDark := math32.Color4{49.0/256.0, 59.0/256.0, 69.0/256.0, 1}
	//bgColorGrey := math32.Color4{85.0/256.0, 85.0/256.0, 85.0/256.0, 1}

	//bgColorOld := math32.Color4{0.85, 0.85, 0.85, 1}

	// Label style
	s.Label = LabelStyle{}
	s.Label.FontAttributes = text.FontAttributes{}
	s.Label.FontAttributes.PointSize = 14
	s.Label.FontAttributes.DPI = 72
	s.Label.FontAttributes.Hinting = text.HintingNone
	s.Label.FontAttributes.LineSpacing = 1.0
	s.Label.BgColor = math32.Color4{1, 1, 1, 0}
	s.Label.FgColor = math32.Color4{1, 1, 1, 1}

	// Button styles
	s.Button = ButtonStyles{}
	s.Button.Normal = ButtonStyle{}
	s.Button.Normal.Border = oneBounds
	s.Button.Normal.Padding = RectBounds{2, 4, 2, 4}
	s.Button.Normal.BorderColor = s.Color.BgDark
	s.Button.Normal.BgColor = s.Color.BgMed
	s.Button.Normal.FgColor = s.Color.Text
	s.Button.Over = s.Button.Normal
	s.Button.Over.BgColor = s.Color.BgOver
	s.Button.Focus = s.Button.Over
	s.Button.Pressed = s.Button.Over
	s.Button.Pressed.Border = RectBounds{2, 2, 2, 2}
	s.Button.Pressed.Padding = RectBounds{2, 2, 0, 4}
	s.Button.Disabled = s.Button.Normal
	s.Button.Disabled.BorderColor = s.Color.TextDis
	s.Button.Disabled.FgColor = s.Color.TextDis

	// CheckRadio styles
	s.CheckRadio = CheckRadioStyles{}
	s.CheckRadio.Normal = CheckRadioStyle{}
	s.CheckRadio.Normal.BorderColor = borderColor
	s.CheckRadio.Normal.BgColor = transparent
	s.CheckRadio.Normal.FgColor = s.Color.Text
	s.CheckRadio.Over = s.CheckRadio.Normal
	s.CheckRadio.Over.BgColor = s.Color.BgOver
	s.CheckRadio.Focus = s.CheckRadio.Over
	s.CheckRadio.Disabled = s.CheckRadio.Normal
	s.CheckRadio.Disabled.FgColor = s.Color.TextDis

	// Edit styles
	s.Edit = EditStyles{}
	s.Edit.Normal = EditStyle{
		Border:      oneBounds,
		Paddings:    zeroBounds,
		BorderColor: borderColor,
		BgColor:     s.Color.BgMed,
		BgAlpha:     1.0,
		FgColor:     s.Color.Text,
		HolderColor: math32.Color4{0.4, 0.4, 0.4, 1},
	}
	s.Edit.Over = s.Edit.Normal
	s.Edit.Over.BgColor = s.Color.BgNormal
	s.Edit.Focus = s.Edit.Normal
	s.Edit.Disabled = s.Edit.Normal
	s.Edit.Disabled.FgColor = s.Color.TextDis

	// ScrollBar styles
	s.ScrollBar = ScrollBarStyles{}
	s.ScrollBar.Normal = ScrollBarStyle{}
	s.ScrollBar.Normal.Padding = oneBounds
	s.ScrollBar.Normal.BgColor = math32.Color4{0, 0, 0, 0.2}
	s.ScrollBar.Normal.ButtonLength = 32
	s.ScrollBar.Normal.Button = PanelStyle{
		BgColor: math32.Color4{0.8, 0.8, 0.8, 0.5},
	}
	s.ScrollBar.Over = s.ScrollBar.Normal
	s.ScrollBar.Disabled = s.ScrollBar.Normal

	// Slider styles
	s.Slider = SliderStyles{}
	s.Slider.Normal = SliderStyle{}
	s.Slider.Normal.Border = oneBounds
	s.Slider.Normal.BorderColor = borderColor
	s.Slider.Normal.BgColor = s.Color.BgDark
	s.Slider.Normal.FgColor = s.Color.Highlight //bgColorBlue2 //math32.Color4{0, 0.4, 0, 1}
	s.Slider.Normal.FgColor.A = 0.5
	s.Slider.Over = s.Slider.Normal
	s.Slider.Over.BgColor = s.Color.BgNormal
	s.Slider.Over.FgColor = s.Color.Highlight //math32.Color4{0, 0.5, 0, 1}
	s.Slider.Focus = s.Slider.Over
	s.Slider.Disabled = s.Slider.Normal

	// Splitter styles
	s.Splitter = SplitterStyles{}
	s.Splitter.Normal = SplitterStyle{
		SpacerBorderColor: borderColor,
		SpacerColor:       s.Color.BgNormal,
		SpacerSize:        6,
	}
	s.Splitter.Over = s.Splitter.Normal
	s.Splitter.Over.SpacerColor = s.Color.BgOver
	s.Splitter.Drag = s.Splitter.Over

	// Window styles
	s.Window = WindowStyles{}
	s.Window.Normal = WindowStyle{}
	s.Window.Normal.Border = RectBounds{2, 2, 2, 2}
	s.Window.Normal.Padding = zeroBounds
	s.Window.Normal.BorderColor = s.Color.BgDark
	s.Window.Normal.TitleStyle = WindowTitleStyle{}
	s.Window.Normal.TitleStyle.Border = RectBounds{0, 0, 1, 0}
	s.Window.Normal.TitleStyle.BorderColor = math32.Color4{0, 0, 0, 1}
	s.Window.Normal.TitleStyle.BgColor = s.Color.Select
	s.Window.Normal.TitleStyle.FgColor = s.Color.Text
	s.Window.Over = s.Window.Normal
	s.Window.Focus = s.Window.Normal
	s.Window.Disabled = s.Window.Normal

	// ItemScroller styles
	s.Scroller = ScrollerStyle{}
	s.Scroller.VerticalScrollbar = ScrollerScrollbarStyle{}
	s.Scroller.VerticalScrollbar.ScrollBarStyle = s.ScrollBar.Normal
	s.Scroller.VerticalScrollbar.Broadness = 12
	s.Scroller.VerticalScrollbar.Position = ScrollbarRight
	s.Scroller.VerticalScrollbar.OverlapContent = true
	s.Scroller.VerticalScrollbar.AutoSizeButton = true
	s.Scroller.HorizontalScrollbar = s.Scroller.VerticalScrollbar
	s.Scroller.HorizontalScrollbar.Position = ScrollbarBottom
	s.Scroller.ScrollbarInterlocking = ScrollbarInterlockingNone
	s.Scroller.CornerCovered = true
	s.Scroller.CornerPanel = PanelStyle{}
	s.Scroller.CornerPanel.BgColor = math32.Color4{0, 0, 0, 0.2}
	s.Scroller.Border = oneBounds
	s.Scroller.BorderColor = borderColor
	s.Scroller.BgColor = s.Color.BgNormal

	// ItemScroller styles
	s.ItemScroller = ItemScrollerStyles{}
	s.ItemScroller.Normal = ItemScrollerStyle{}
	s.ItemScroller.Normal.Border = oneBounds
	s.ItemScroller.Normal.BorderColor = borderColor
	s.ItemScroller.Normal.BgColor = s.Color.BgNormal
	s.ItemScroller.Normal.FgColor = s.Color.Text
	s.ItemScroller.Over = s.ItemScroller.Normal
	//s.ItemScroller.Over.BgColor = bgColorOver
	s.ItemScroller.Focus = s.ItemScroller.Over
	s.ItemScroller.Disabled = s.ItemScroller.Normal

	// ItemList styles
	s.List = ListStyles{}
	s.List.Scroller = &s.ItemScroller
	s.List.Item = &ListItemStyles{}
	s.List.Item.Normal = ListItemStyle{}
	s.List.Item.Normal.Border = RectBounds{0, 0, 1, 0}
	s.List.Item.Normal.Padding = RectBounds{0, 0, 0, 2}
	s.List.Item.Normal.BorderColor = math32.Color4{0, 0, 0, 0}
	s.List.Item.Normal.BgColor = transparent
	s.List.Item.Normal.FgColor = s.Color.Text
	s.List.Item.Over = s.List.Item.Normal
	s.List.Item.Over.BgColor = s.Color.BgOver
	s.List.Item.Over.FgColor = s.Color.Select
	s.List.Item.Selected = s.List.Item.Normal
	s.List.Item.Selected.BgColor = s.Color.Highlight
	s.List.Item.Selected.FgColor = s.Color.Select
	s.List.Item.Highlighted = s.List.Item.Normal
	s.List.Item.Highlighted.BorderColor = math32.Color4{0, 0, 0, 1}
	s.List.Item.Highlighted.BgColor = s.Color.BgOver
	s.List.Item.Highlighted.FgColor = s.Color.Text
	s.List.Item.SelHigh = s.List.Item.Highlighted
	s.List.Item.SelHigh.BgColor = s.Color.BgNormal
	s.List.Item.SelHigh.FgColor = s.Color.Select

	// DropDown styles
	s.DropDown = DropDownStyles{}
	s.DropDown.Normal = DropDownStyle{}
	s.DropDown.Normal.Border = oneBounds
	s.DropDown.Normal.Padding = RectBounds{0, 0, 0, 2}
	s.DropDown.Normal.BorderColor = borderColor
	s.DropDown.Normal.BgColor = s.Color.BgNormal
	s.DropDown.Normal.FgColor = s.Color.Text
	s.DropDown.Over = s.DropDown.Normal
	s.DropDown.Over.BgColor = s.Color.BgOver
	s.DropDown.Focus = s.DropDown.Over
	s.DropDown.Disabled = s.DropDown.Normal

	// Folder styles
	s.Folder = FolderStyles{}
	s.Folder.Normal = FolderStyle{}
	s.Folder.Normal.Border = oneBounds
	s.Folder.Normal.Padding = RectBounds{2, 0, 2, 2}
	s.Folder.Normal.BorderColor = borderColor
	s.Folder.Normal.BgColor = s.Color.BgNormal
	s.Folder.Normal.FgColor = s.Color.Text
	s.Folder.Normal.Icons = [2]string{icon.ExpandMore, icon.ExpandLess}
	s.Folder.Over = s.Folder.Normal
	s.Folder.Over.BgColor = s.Color.BgOver
	s.Folder.Focus = s.Folder.Over
	s.Folder.Focus.Padding = twoBounds
	s.Folder.Disabled = s.Folder.Focus

	// Tree styles
	s.Tree = TreeStyles{}
	s.Tree.Padlevel = 28.0
	s.Tree.List = &s.List
	s.Tree.Node = &TreeNodeStyles{}
	s.Tree.Node.Normal = TreeNodeStyle{}
	s.Tree.Node.Normal.BorderColor = borderColor
	s.Tree.Node.Normal.BgColor = transparent
	s.Tree.Node.Normal.FgColor = s.Color.Text
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
	s.Menu.Body.Normal.BgColor = s.Color.BgNormal
	s.Menu.Body.Normal.FgColor = s.Color.Text
	s.Menu.Body.Over = s.Menu.Body.Normal
	s.Menu.Body.Over.BgColor = s.Color.BgOver
	s.Menu.Body.Focus = s.Menu.Body.Over
	s.Menu.Body.Disabled = s.Menu.Body.Normal
	s.Menu.Item = &MenuItemStyles{}
	s.Menu.Item.Normal = MenuItemStyle{}
	s.Menu.Item.Normal.Padding = RectBounds{2, 4, 2, 2}
	s.Menu.Item.Normal.BorderColor = borderColor
	s.Menu.Item.Normal.BgColor = s.Color.BgNormal
	s.Menu.Item.Normal.FgColor = s.Color.Text
	s.Menu.Item.Normal.IconPaddings = RectBounds{0, 6, 0, 4}
	s.Menu.Item.Normal.ShortcutPaddings = RectBounds{0, 0, 0, 10}
	s.Menu.Item.Normal.RiconPaddings = RectBounds{2, 0, 0, 4}
	s.Menu.Item.Over = s.Menu.Item.Normal
	s.Menu.Item.Over.BgColor = s.Color.Highlight
	s.Menu.Item.Disabled = s.Menu.Item.Normal
	s.Menu.Item.Disabled.FgColor = s.Color.TextDis
	s.Menu.Item.Separator = MenuItemStyle{}
	s.Menu.Item.Separator.Border = twoBounds
	s.Menu.Item.Separator.Padding = zeroBounds
	s.Menu.Item.Separator.BorderColor = math32.Color4{0, 0, 0, 0}
	s.Menu.Item.Separator.BgColor = math32.Color4{0.6, 0.6, 0.6, 1}
	s.Menu.Item.Separator.FgColor = s.Color.Text

	// Table styles
	s.Table = TableStyles{}
	s.Table.Header = TableHeaderStyle{}
	s.Table.Header.Border = RectBounds{0, 1, 1, 0}
	s.Table.Header.Padding = twoBounds
	s.Table.Header.BorderColor = s.Color.BgNormal
	s.Table.Header.BgColor = s.Color.BgDark
	s.Table.Header.FgColor = s.Color.Text
	s.Table.RowEven = TableRowStyle{}
	s.Table.RowEven.Border = RectBounds{0, 1, 1, 0}
	s.Table.RowEven.Padding = twoBounds
	s.Table.RowEven.BorderColor = s.Color.BgDark
	s.Table.RowEven.BgColor = s.Color.BgNormal
	s.Table.RowEven.FgColor = s.Color.Text
	s.Table.RowOdd = s.Table.RowEven
	s.Table.RowOdd.BgColor = s.Color.BgMed
	s.Table.RowCursor = s.Table.RowEven
	s.Table.RowCursor.BgColor = s.Color.Highlight
	s.Table.RowSel = s.Table.RowEven
	s.Table.RowSel.BgColor = s.Color.Select
	s.Table.Status = TableStatusStyle{}
	s.Table.Status.Border = RectBounds{1, 0, 0, 0}
	s.Table.Status.Padding = twoBounds
	s.Table.Status.BorderColor = borderColor
	s.Table.Status.BgColor = s.Color.BgDark
	s.Table.Status.FgColor = s.Color.Text
	s.Table.Resizer = TableResizerStyle{
		Width:       4,
		Border:      zeroBounds,
		BorderColor: borderColor,
		BgColor:     math32.Color4{0.4, 0.4, 0.4, 0.6},
	}

	// ImageButton styles
	s.ImageButton = ImageButtonStyles{}
	s.ImageButton.Normal = ImageButtonStyle{}
	s.ImageButton.Normal.BgColor = transparent
	s.ImageButton.Normal.FgColor = s.Color.Text
	s.ImageButton.Over = s.ImageButton.Normal
	s.ImageButton.Over.BgColor = s.Color.BgOver
	s.ImageButton.Focus = s.ImageButton.Over
	s.ImageButton.Pressed = s.ImageButton.Over
	s.ImageButton.Pressed.Border = oneBounds
	s.ImageButton.Disabled = s.ImageButton.Normal
	s.ImageButton.Disabled.FgColor = s.Color.TextDis

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
	s.TabBar.Normal.BgColor = s.Color.BgMed
	s.TabBar.Over = s.TabBar.Normal
	//s.TabBar.Over.BgColor = s.Color.BgOver
	s.TabBar.Focus = s.TabBar.Normal
	s.TabBar.Focus.BgColor = transparent
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
	s.TabBar.Tab.Normal.BgColor = s.Color.BgNormal
	s.TabBar.Tab.Normal.FgColor = s.Color.Text
	s.TabBar.Tab.Over = s.TabBar.Tab.Normal
	s.TabBar.Tab.Over.BgColor = s.Color.BgOver
	s.TabBar.Tab.Focus = s.TabBar.Tab.Normal
	s.TabBar.Tab.Focus.BgColor = transparent
	s.TabBar.Tab.Disabled = s.TabBar.Tab.Focus
	s.TabBar.Tab.Selected = s.TabBar.Tab.Normal
	s.TabBar.Tab.Selected.BgColor = s.Color.BgOver

	return s
}
