// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/g3n/engine/gui/assets"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/text"
)

func init() {

	setupDefaultStyle()
}

// Pointer to default style
var StyleDefault *Style

// All styles
type Style struct {
	Font          *text.Font
	FontIcon      *text.Font
	Button        ButtonStyles
	CheckRadio    CheckRadioStyles
	Edit          EditStyles
	ScrollBar     ScrollBarStyle
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
}

const (
	defaultFont     = "fonts/FreeSans.ttf"
	defaultFontBold = "fonts/FreeSansBold.ttf"
	defaultFontIcon = "fonts/MaterialIcons-Regular.ttf"
)

const (
	OverStyle = iota + 1
	FocusStyle
	DisabledStyle
	NormalStyle
	DefaultStyle
)

// setupDefaultStyle initializes the default Gui global styles
func setupDefaultStyle() {

	StyleDefault = &Style{}

	// Creates Default Font
	fontData := assets.MustAsset(defaultFont)
	font, err := text.NewFontFromData(fontData)
	if err != nil {
		panic(err)
	}
	font.SetLineSpacing(1.0)
	font.SetSize(14)
	font.SetDPI(72)
	font.SetFgColor4(&math32.Color4{0, 0, 0, 1})
	font.SetBgColor4(&math32.Color4{1, 1, 1, 0})
	StyleDefault.Font = font

	// Creates Icon Font
	fontIconData := assets.MustAsset(defaultFontIcon)
	fontIcon, err := text.NewFontFromData(fontIconData)
	if err != nil {
		panic(err)
	}
	fontIcon.SetLineSpacing(1.0)
	fontIcon.SetSize(14)
	fontIcon.SetDPI(72)
	fontIcon.SetFgColor4(&math32.Color4{0, 0, 0, 1})
	fontIcon.SetBgColor4(&math32.Color4{1, 1, 1, 1})
	StyleDefault.FontIcon = fontIcon

	borderSizes := BorderSizes{1, 1, 1, 1}
	borderColor := math32.Color4{0, 0, 0, 1}
	borderColorDis := math32.Color4{0.4, 0.4, 0.4, 1}

	bgColor := math32.Color{0.85, 0.85, 0.85}
	bgColor4 := math32.Color4{0, 0, 0, 0}
	bgColorOver := math32.Color{0.9, 0.9, 0.9}
	bgColor4Over := math32.Color4{1, 1, 1, 0.5}
	bgColor4Sel := math32.Color4{0.6, 0.6, 0.6, 1}

	fgColor := math32.Color{0, 0, 0}
	fgColorSel := math32.Color{0, 0, 0}
	fgColorDis := math32.Color{0.4, 0.4, 0.4}

	// Button styles
	StyleDefault.Button = ButtonStyles{
		Normal: ButtonStyle{
			Border:      borderSizes,
			Paddings:    BorderSizes{2, 4, 2, 4},
			BorderColor: borderColor,
			BgColor:     bgColor,
			FgColor:     fgColor,
		},
		Over: ButtonStyle{
			Border:      borderSizes,
			Paddings:    BorderSizes{2, 4, 2, 4},
			BorderColor: borderColor,
			BgColor:     bgColorOver,
			FgColor:     fgColor,
		},
		Focus: ButtonStyle{
			Border:      borderSizes,
			Paddings:    BorderSizes{2, 4, 2, 4},
			BorderColor: borderColor,
			BgColor:     bgColorOver,
			FgColor:     fgColor,
		},
		Pressed: ButtonStyle{
			Border:      BorderSizes{2, 2, 2, 2},
			Paddings:    BorderSizes{2, 4, 2, 4},
			BorderColor: borderColor,
			BgColor:     bgColorOver,
			FgColor:     fgColor,
		},
		Disabled: ButtonStyle{
			Border:      borderSizes,
			Paddings:    BorderSizes{2, 4, 2, 4},
			BorderColor: borderColorDis,
			BgColor:     bgColor,
			FgColor:     fgColorDis,
		},
	}

	// CheckRadio styles
	StyleDefault.CheckRadio = CheckRadioStyles{
		Normal: CheckRadioStyle{
			Border:      BorderSizes{0, 0, 0, 0},
			Paddings:    BorderSizes{0, 0, 0, 0},
			BorderColor: borderColor,
			BgColor:     bgColor4,
			FgColor:     fgColor,
		},
		Over: CheckRadioStyle{
			Border:      BorderSizes{0, 0, 0, 0},
			Paddings:    BorderSizes{0, 0, 0, 0},
			BorderColor: borderColor,
			BgColor:     bgColor4Over,
			FgColor:     fgColor,
		},
		Focus: CheckRadioStyle{
			Border:      BorderSizes{0, 0, 0, 0},
			Paddings:    BorderSizes{0, 0, 0, 0},
			BorderColor: borderColor,
			BgColor:     bgColor4Over,
			FgColor:     fgColor,
		},
		Disabled: CheckRadioStyle{
			Border:      BorderSizes{0, 0, 0, 0},
			Paddings:    BorderSizes{0, 0, 0, 0},
			BorderColor: borderColor,
			BgColor:     bgColor4,
			FgColor:     fgColorDis,
		},
	}

	// Edit styles
	StyleDefault.Edit = EditStyles{
		Normal: EditStyle{
			Border:      BorderSizes{1, 1, 1, 1},
			Paddings:    BorderSizes{0, 0, 0, 0},
			BorderColor: borderColor,
			BgColor:     bgColor,
			BgAlpha:     1.0,
			FgColor:     fgColor,
			HolderColor: math32.Color{0.4, 0.4, 0.4},
		},
		Over: EditStyle{
			Border:      BorderSizes{1, 1, 1, 1},
			Paddings:    BorderSizes{0, 0, 0, 0},
			BorderColor: borderColor,
			BgColor:     bgColorOver,
			BgAlpha:     1.0,
			FgColor:     fgColor,
			HolderColor: math32.Color{0.4, 0.4, 0.4},
		},
		Focus: EditStyle{
			Border:      BorderSizes{1, 1, 1, 1},
			Paddings:    BorderSizes{0, 0, 0, 0},
			BorderColor: borderColor,
			BgColor:     bgColorOver,
			BgAlpha:     1.0,
			FgColor:     fgColor,
			HolderColor: math32.Color{0.4, 0.4, 0.4},
		},
		Disabled: EditStyle{
			Border:      BorderSizes{1, 1, 1, 1},
			Paddings:    BorderSizes{0, 0, 0, 0},
			BorderColor: borderColor,
			BgColor:     bgColor,
			BgAlpha:     1.0,
			FgColor:     fgColorDis,
			HolderColor: math32.Color{0.4, 0.4, 0.4},
		},
	}

	// ScrollBar style
	StyleDefault.ScrollBar = ScrollBarStyle{
		Paddings:     BorderSizes{1, 1, 1, 1},
		Borders:      BorderSizes{1, 1, 1, 1},
		BordersColor: borderColor,
		Color:        math32.Color{0.8, 0.8, 0.8},
		Button: ScrollBarButtonStyle{
			Borders:      BorderSizes{1, 1, 1, 1},
			BordersColor: borderColor,
			Color:        math32.Color{0.5, 0.5, 0.5},
			Size:         30,
		},
	}

	// Slider styles
	StyleDefault.Slider = SliderStyles{
		Normal: SliderStyle{
			Border:      borderSizes,
			BorderColor: borderColor,
			Paddings:    BorderSizes{0, 0, 0, 0},
			BgColor:     math32.Color4{0.8, 0.8, 0.8, 1},
			FgColor:     math32.Color4{0, 0.8, 0, 1},
		},
		Over: SliderStyle{
			Border:      borderSizes,
			BorderColor: borderColor,
			Paddings:    BorderSizes{0, 0, 0, 0},
			BgColor:     math32.Color4{1, 1, 1, 1},
			FgColor:     math32.Color4{0, 1, 0, 1},
		},
		Focus: SliderStyle{
			Border:      borderSizes,
			BorderColor: borderColor,
			Paddings:    BorderSizes{0, 0, 0, 0},
			BgColor:     math32.Color4{1, 1, 1, 1},
			FgColor:     math32.Color4{0, 1, 0, 1},
		},
		Disabled: SliderStyle{
			Border:      borderSizes,
			BorderColor: borderColor,
			Paddings:    BorderSizes{0, 0, 0, 0},
			BgColor:     math32.Color4{0.8, 0.8, 0.8, 1},
			FgColor:     math32.Color4{0, 0.8, 0, 1},
		},
	}

	// Splitter styles
	StyleDefault.Splitter = SplitterStyles{
		Normal: SplitterStyle{
			SpacerBorderColor: borderColor,
			SpacerColor:       bgColor,
			SpacerSize:        6,
		},
		Over: SplitterStyle{
			SpacerBorderColor: borderColor,
			SpacerColor:       bgColorOver,
			SpacerSize:        6,
		},
		Drag: SplitterStyle{
			SpacerBorderColor: borderColor,
			SpacerColor:       bgColorOver,
			SpacerSize:        6,
		},
	}

	StyleDefault.Window = WindowStyles{
		Normal: WindowStyle{
			Border:           BorderSizes{4, 4, 4, 4},
			Paddings:         BorderSizes{0, 0, 0, 0},
			BorderColor:      math32.Color4{0.2, 0.2, 0.2, 1},
			TitleBorders:     BorderSizes{0, 0, 1, 0},
			TitleBorderColor: math32.Color4{0, 0, 0, 1},
			TitleBgColor:     math32.Color4{0, 1, 0, 1},
			TitleFgColor:     math32.Color4{0, 0, 0, 1},
		},
		Over: WindowStyle{
			Border:           BorderSizes{4, 4, 4, 4},
			Paddings:         BorderSizes{0, 0, 0, 0},
			BorderColor:      math32.Color4{0.2, 0.2, 0.2, 1},
			TitleBorders:     BorderSizes{0, 0, 1, 0},
			TitleBorderColor: math32.Color4{0, 0, 0, 1},
			TitleBgColor:     math32.Color4{0, 1, 0, 1},
			TitleFgColor:     math32.Color4{0, 0, 0, 1},
		},
		Focus: WindowStyle{
			Border:           BorderSizes{4, 4, 4, 4},
			Paddings:         BorderSizes{0, 0, 0, 0},
			BorderColor:      math32.Color4{0.2, 0.2, 0.2, 1},
			TitleBorders:     BorderSizes{0, 0, 1, 0},
			TitleBorderColor: math32.Color4{0, 0, 0, 1},
			TitleBgColor:     math32.Color4{0, 1, 0, 1},
			TitleFgColor:     math32.Color4{0, 0, 0, 1},
		},
		Disabled: WindowStyle{
			Border:           BorderSizes{4, 4, 4, 4},
			Paddings:         BorderSizes{0, 0, 0, 0},
			BorderColor:      math32.Color4{0.2, 0.2, 0.2, 1},
			TitleBorders:     BorderSizes{0, 0, 1, 0},
			TitleBorderColor: math32.Color4{0, 0, 0, 1},
			TitleBgColor:     math32.Color4{0, 1, 0, 1},
			TitleFgColor:     math32.Color4{0, 0, 0, 1},
		},
	}

	// Scroller styles
	StyleDefault.Scroller = ScrollerStyles{
		Normal: ScrollerStyle{
			Border:      BorderSizes{1, 1, 1, 1},
			Paddings:    BorderSizes{0, 0, 0, 0},
			BorderColor: borderColor,
			BgColor:     bgColor,
			FgColor:     fgColor,
		},
		Over: ScrollerStyle{
			Border:      BorderSizes{1, 1, 1, 1},
			Paddings:    BorderSizes{0, 0, 0, 0},
			BorderColor: borderColor,
			BgColor:     bgColorOver,
			FgColor:     fgColor,
		},
		Focus: ScrollerStyle{
			Border:      BorderSizes{1, 1, 1, 1},
			Paddings:    BorderSizes{0, 0, 0, 0},
			BorderColor: borderColor,
			BgColor:     bgColorOver,
			FgColor:     fgColor,
		},
		Disabled: ScrollerStyle{
			Border:      BorderSizes{1, 1, 1, 1},
			Paddings:    BorderSizes{0, 0, 0, 0},
			BorderColor: borderColor,
			BgColor:     bgColor,
			FgColor:     fgColor,
		},
	}

	// List styles
	StyleDefault.List = ListStyles{
		Scroller: &ScrollerStyles{
			Normal: ScrollerStyle{
				Border:      BorderSizes{1, 1, 1, 1},
				Paddings:    BorderSizes{0, 0, 0, 0},
				BorderColor: borderColor,
				BgColor:     bgColor,
				FgColor:     fgColor,
			},
			Over: ScrollerStyle{
				Border:      BorderSizes{1, 1, 1, 1},
				Paddings:    BorderSizes{0, 0, 0, 0},
				BorderColor: borderColor,
				BgColor:     bgColorOver,
				FgColor:     fgColor,
			},
			Focus: ScrollerStyle{
				Border:      BorderSizes{1, 1, 1, 1},
				Paddings:    BorderSizes{0, 0, 0, 0},
				BorderColor: borderColor,
				BgColor:     bgColorOver,
				FgColor:     fgColor,
			},
			Disabled: ScrollerStyle{
				Border:      BorderSizes{1, 1, 1, 1},
				Paddings:    BorderSizes{0, 0, 0, 0},
				BorderColor: borderColor,
				BgColor:     bgColor,
				FgColor:     fgColor,
			},
		},
		Item: &ListItemStyles{
			Normal: ListItemStyle{
				Border:      BorderSizes{1, 0, 1, 0},
				Paddings:    BorderSizes{0, 0, 0, 2},
				BorderColor: math32.Color4{0, 0, 0, 0},
				BgColor:     bgColor4,
				FgColor:     fgColor,
			},
			Selected: ListItemStyle{
				Border:      BorderSizes{1, 0, 1, 0},
				Paddings:    BorderSizes{0, 0, 0, 2},
				BorderColor: math32.Color4{0, 0, 0, 0},
				BgColor:     bgColor4Sel,
				FgColor:     fgColorSel,
			},
			Highlighted: ListItemStyle{
				Border:      BorderSizes{1, 0, 1, 0},
				Paddings:    BorderSizes{0, 0, 0, 2},
				BorderColor: math32.Color4{0, 0, 0, 1},
				BgColor:     bgColor4Over,
				FgColor:     fgColor,
			},
			SelHigh: ListItemStyle{
				Border:      BorderSizes{1, 0, 1, 0},
				Paddings:    BorderSizes{0, 0, 0, 2},
				BorderColor: math32.Color4{0, 0, 0, 1},
				BgColor:     bgColor4Sel,
				FgColor:     fgColorSel,
			},
		},
	}

	StyleDefault.DropDown = DropDownStyles{
		Normal: &DropDownStyle{
			Border:      BorderSizes{1, 1, 1, 1},
			Paddings:    BorderSizes{0, 0, 0, 2},
			BorderColor: borderColor,
			BgColor:     bgColor,
			FgColor:     fgColor,
		},
		Over: &DropDownStyle{
			Border:      BorderSizes{1, 1, 1, 1},
			Paddings:    BorderSizes{0, 0, 0, 2},
			BorderColor: borderColor,
			BgColor:     bgColorOver,
			FgColor:     fgColor,
		},
		Focus: &DropDownStyle{
			Border:      BorderSizes{1, 1, 1, 1},
			Paddings:    BorderSizes{0, 0, 0, 2},
			BorderColor: borderColor,
			BgColor:     bgColorOver,
			FgColor:     fgColor,
		},
		Disabled: &DropDownStyle{
			Border:      BorderSizes{1, 1, 1, 1},
			Paddings:    BorderSizes{0, 0, 0, 2},
			BorderColor: borderColor,
			BgColor:     bgColor,
			FgColor:     fgColor,
		},
	}

	StyleDefault.Folder = FolderStyles{
		Normal: &FolderStyle{
			Margins:     BorderSizes{0, 0, 0, 0},
			Border:      BorderSizes{1, 1, 1, 1},
			Paddings:    BorderSizes{2, 0, 2, 2},
			BorderColor: borderColor,
			BgColor:     bgColor,
			FgColor:     fgColor,
			Icons:       [2]int{assets.ExpandMore, assets.ExpandLess},
		},
		Over: &FolderStyle{
			Margins:     BorderSizes{0, 0, 0, 0},
			Border:      BorderSizes{1, 1, 1, 1},
			Paddings:    BorderSizes{2, 0, 2, 2},
			BorderColor: borderColor,
			BgColor:     bgColorOver,
			FgColor:     fgColor,
			Icons:       [2]int{assets.ExpandMore, assets.ExpandLess},
		},
		Focus: &FolderStyle{
			Margins:     BorderSizes{0, 0, 0, 0},
			Border:      BorderSizes{1, 1, 1, 1},
			Paddings:    BorderSizes{2, 2, 2, 2},
			BorderColor: borderColor,
			BgColor:     bgColorOver,
			FgColor:     fgColor,
			Icons:       [2]int{assets.ExpandMore, assets.ExpandLess},
		},
		Disabled: &FolderStyle{
			Margins:     BorderSizes{0, 0, 0, 0},
			Border:      BorderSizes{1, 1, 1, 1},
			Paddings:    BorderSizes{2, 2, 2, 2},
			BorderColor: borderColor,
			BgColor:     bgColorOver,
			FgColor:     fgColor,
			Icons:       [2]int{assets.ExpandMore, assets.ExpandLess},
		},
	}

	StyleDefault.Tree = TreeStyles{
		List: &ListStyles{
			Scroller: &ScrollerStyles{
				Normal: ScrollerStyle{
					Border:      BorderSizes{1, 1, 1, 1},
					Paddings:    BorderSizes{0, 0, 0, 0},
					BorderColor: borderColor,
					BgColor:     bgColor,
					FgColor:     fgColor,
				},
				Over: ScrollerStyle{
					Border:      BorderSizes{1, 1, 1, 1},
					Paddings:    BorderSizes{0, 0, 0, 0},
					BorderColor: borderColor,
					BgColor:     bgColorOver,
					FgColor:     fgColor,
				},
				Focus: ScrollerStyle{
					Border:      BorderSizes{1, 1, 1, 1},
					Paddings:    BorderSizes{0, 0, 0, 0},
					BorderColor: borderColor,
					BgColor:     bgColorOver,
					FgColor:     fgColor,
				},
				Disabled: ScrollerStyle{
					Border:      BorderSizes{1, 1, 1, 1},
					Paddings:    BorderSizes{0, 0, 0, 0},
					BorderColor: borderColor,
					BgColor:     bgColor,
					FgColor:     fgColor,
				},
			},
			Item: &ListItemStyles{
				Normal: ListItemStyle{
					Border:      BorderSizes{1, 0, 1, 0},
					Paddings:    BorderSizes{0, 0, 0, 2},
					BorderColor: math32.Color4{0, 0, 0, 0},
					BgColor:     bgColor4,
					FgColor:     fgColor,
				},
				Selected: ListItemStyle{
					Border:      BorderSizes{1, 0, 1, 0},
					Paddings:    BorderSizes{0, 0, 0, 2},
					BorderColor: math32.Color4{0, 0, 0, 0},
					BgColor:     bgColor4Sel,
					FgColor:     fgColorSel,
				},
				Highlighted: ListItemStyle{
					Border:      BorderSizes{1, 0, 1, 0},
					Paddings:    BorderSizes{0, 0, 0, 2},
					BorderColor: math32.Color4{0, 0, 0, 1},
					BgColor:     bgColor4Over,
					FgColor:     fgColor,
				},
				SelHigh: ListItemStyle{
					Border:      BorderSizes{1, 0, 1, 0},
					Paddings:    BorderSizes{0, 0, 0, 2},
					BorderColor: math32.Color4{0, 0, 0, 1},
					BgColor:     bgColor4Sel,
					FgColor:     fgColorSel,
				},
			},
		},
		Node: &TreeNodeStyles{
			Normal: TreeNodeStyle{
				Border:      BorderSizes{0, 0, 0, 0},
				Paddings:    BorderSizes{0, 0, 0, 0},
				BorderColor: borderColor,
				BgColor:     bgColor,
				FgColor:     fgColor,
				Icons:       [2]int{assets.ExpandMore, assets.ExpandLess},
			},
		},
		Padlevel: 16.0,
	}

	StyleDefault.ControlFolder = ControlFolderStyles{
		Folder: &FolderStyles{
			Normal: &FolderStyle{
				Margins:     BorderSizes{0, 0, 0, 0},
				Border:      BorderSizes{1, 1, 1, 1},
				Paddings:    BorderSizes{2, 0, 2, 2},
				BorderColor: math32.Color4{0, 0, 0, 0},
				BgColor:     math32.Color{0, 0.5, 1},
				FgColor:     fgColor,
				Icons:       [2]int{assets.ExpandMore, assets.ExpandLess},
			},
			Over: &FolderStyle{
				Margins:     BorderSizes{0, 0, 0, 0},
				Border:      BorderSizes{1, 1, 1, 1},
				Paddings:    BorderSizes{2, 0, 2, 2},
				BorderColor: math32.Color4{0, 0, 0, 0},
				BgColor:     math32.Color{0, 0.5, 1},
				FgColor:     fgColor,
				Icons:       [2]int{assets.ExpandMore, assets.ExpandLess},
			},
			Focus: &FolderStyle{
				Margins:     BorderSizes{0, 0, 0, 0},
				Border:      BorderSizes{1, 1, 1, 1},
				Paddings:    BorderSizes{2, 2, 2, 2},
				BorderColor: math32.Color4{0, 0, 0, 0},
				BgColor:     math32.Color{0, 0.5, 1},
				FgColor:     fgColor,
				Icons:       [2]int{assets.ExpandMore, assets.ExpandLess},
			},
			Disabled: &FolderStyle{
				Margins:     BorderSizes{0, 0, 0, 0},
				Border:      BorderSizes{1, 1, 1, 1},
				Paddings:    BorderSizes{2, 2, 2, 2},
				BorderColor: math32.Color4{0, 0, 0, 0},
				BgColor:     math32.Color{0, 0.5, 1},
				FgColor:     fgColor,
				Icons:       [2]int{assets.ExpandMore, assets.ExpandLess},
			},
		},
		Tree: &TreeStyles{
			List: &ListStyles{
				Scroller: &ScrollerStyles{
					Normal: ScrollerStyle{
						Border:      BorderSizes{1, 1, 1, 1},
						Paddings:    BorderSizes{0, 2, 0, 0},
						BorderColor: borderColor,
						BgColor:     bgColor,
						FgColor:     fgColor,
					},
					Over: ScrollerStyle{
						Border:      BorderSizes{1, 1, 1, 1},
						Paddings:    BorderSizes{0, 2, 0, 0},
						BorderColor: borderColor,
						BgColor:     bgColorOver,
						FgColor:     fgColor,
					},
					Focus: ScrollerStyle{
						Border:      BorderSizes{1, 1, 1, 1},
						Paddings:    BorderSizes{0, 2, 0, 0},
						BorderColor: borderColor,
						BgColor:     bgColorOver,
						FgColor:     fgColor,
					},
					Disabled: ScrollerStyle{
						Border:      BorderSizes{1, 1, 1, 1},
						Paddings:    BorderSizes{0, 2, 0, 0},
						BorderColor: borderColor,
						BgColor:     bgColor,
						FgColor:     fgColor,
					},
				},
				Item: &ListItemStyles{
					Normal: ListItemStyle{
						Border:      BorderSizes{1, 0, 1, 0},
						Paddings:    BorderSizes{0, 0, 0, 2},
						BorderColor: math32.Color4{0, 0, 0, 0},
						BgColor:     bgColor4,
						FgColor:     fgColor,
					},
					Selected: ListItemStyle{
						Border:      BorderSizes{1, 0, 1, 0},
						Paddings:    BorderSizes{0, 0, 0, 2},
						BorderColor: math32.Color4{0, 0, 0, 0},
						BgColor:     bgColor4,
						FgColor:     fgColor,
					},
					Highlighted: ListItemStyle{
						Border:      BorderSizes{1, 0, 1, 0},
						Paddings:    BorderSizes{0, 0, 0, 2},
						BorderColor: math32.Color4{0, 0, 0, 1},
						BgColor:     bgColor4Over,
						FgColor:     fgColor,
					},
					SelHigh: ListItemStyle{
						Border:      BorderSizes{1, 0, 1, 0},
						Paddings:    BorderSizes{0, 0, 0, 2},
						BorderColor: math32.Color4{0, 0, 0, 1},
						BgColor:     bgColor4Sel,
						FgColor:     fgColorSel,
					},
				},
			},
			Node: &TreeNodeStyles{
				Normal: TreeNodeStyle{
					Border:      BorderSizes{0, 0, 0, 0},
					Paddings:    BorderSizes{0, 0, 0, 0},
					BorderColor: borderColor,
					BgColor:     bgColor,
					FgColor:     fgColor,
					Icons:       [2]int{assets.ExpandMore, assets.ExpandLess},
				},
			},
			Padlevel: 2.0,
		},
	}

	// Menu styles
	StyleDefault.Menu = MenuStyles{
		Body: &MenuBodyStyles{
			Normal: MenuBodyStyle{
				Border:      BorderSizes{1, 1, 1, 1},
				Paddings:    BorderSizes{2, 2, 2, 2},
				BorderColor: borderColor,
				BgColor:     bgColor,
				FgColor:     fgColor,
			},
			Over: MenuBodyStyle{
				Border:      BorderSizes{1, 1, 1, 1},
				Paddings:    BorderSizes{2, 2, 2, 2},
				BorderColor: borderColor,
				BgColor:     bgColorOver,
				FgColor:     fgColor,
			},
			Focus: MenuBodyStyle{
				Border:      BorderSizes{1, 1, 1, 1},
				Paddings:    BorderSizes{2, 2, 2, 2},
				BorderColor: borderColor,
				BgColor:     bgColorOver,
				FgColor:     fgColor,
			},
			Disabled: MenuBodyStyle{
				Border:      BorderSizes{1, 1, 1, 1},
				Paddings:    BorderSizes{2, 2, 2, 2},
				BorderColor: borderColor,
				BgColor:     bgColor,
				FgColor:     fgColor,
			},
		},
		Item: &MenuItemStyles{
			Normal: MenuItemStyle{
				Border:           BorderSizes{0, 0, 0, 0},
				Paddings:         BorderSizes{2, 4, 2, 2},
				BorderColor:      borderColor,
				BgColor:          bgColor,
				FgColor:          fgColor,
				IconPaddings:     BorderSizes{0, 6, 0, 4},
				ShortcutPaddings: BorderSizes{0, 0, 0, 10},
				RiconPaddings:    BorderSizes{2, 0, 0, 4},
			},
			Over: MenuItemStyle{
				Border:           BorderSizes{0, 0, 0, 0},
				Paddings:         BorderSizes{2, 4, 2, 2},
				BorderColor:      borderColor,
				BgColor:          math32.Color{0.6, 0.6, 0.6},
				FgColor:          fgColor,
				IconPaddings:     BorderSizes{0, 6, 0, 4},
				ShortcutPaddings: BorderSizes{0, 0, 0, 10},
				RiconPaddings:    BorderSizes{2, 0, 0, 4},
			},
			Disabled: MenuItemStyle{
				Border:           BorderSizes{0, 0, 0, 0},
				Paddings:         BorderSizes{2, 4, 2, 2},
				BorderColor:      borderColor,
				BgColor:          bgColor,
				FgColor:          fgColorDis,
				IconPaddings:     BorderSizes{0, 6, 0, 4},
				ShortcutPaddings: BorderSizes{0, 0, 0, 10},
				RiconPaddings:    BorderSizes{2, 0, 0, 4},
			},
			Separator: MenuItemStyle{
				Border:      BorderSizes{2, 2, 2, 2},
				Paddings:    BorderSizes{0, 0, 0, 0},
				BorderColor: math32.Color4{0, 0, 0, 0},
				BgColor:     math32.Color{0.6, 0.6, 0.6},
				FgColor:     fgColor,
			},
		},
	}

	// Table styles
	StyleDefault.Table = TableStyles{
		Header: &TableHeaderStyle{
			Border:      BorderSizes{0, 1, 1, 0},
			Paddings:    BorderSizes{2, 2, 2, 2},
			BorderColor: borderColor,
			BgColor:     math32.Color{0.7, 0.7, 0.7},
			FgColor:     fgColor,
		},
		Row: &TableRowStyles{
			Normal: TableRowStyle{
				Border:      BorderSizes{0, 1, 1, 0},
				Paddings:    BorderSizes{2, 2, 2, 2},
				BorderColor: math32.Color4{0.6, 0.6, 0.6, 1},
				BgColor:     bgColor,
				FgColor:     fgColor,
			},
			Selected: TableRowStyle{
				Border:      BorderSizes{0, 1, 1, 0},
				Paddings:    BorderSizes{2, 2, 2, 2},
				BorderColor: math32.Color4{0.6, 0.6, 0.6, 1},
				BgColor:     math32.Color{0.7, 0.7, 0.7},
				FgColor:     fgColor,
			},
		},
		Status: &TableStatusStyle{
			Border:      BorderSizes{1, 0, 0, 0},
			Paddings:    BorderSizes{2, 2, 2, 2},
			BorderColor: borderColor,
			BgColor:     math32.Color{0.9, 0.9, 0.9},
			FgColor:     fgColor,
		},
	}
}
