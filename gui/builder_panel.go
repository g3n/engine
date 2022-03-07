// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"path/filepath"

	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/window"
)

// buildPanel builds an object of type Panel
func buildPanel(b *Builder, am map[string]interface{}) (IPanel, error) {

	pan := NewPanel(0, 0)
	err := b.SetAttribs(am, pan)
	if err != nil {
		return nil, err
	}

	// Builds children recursively
	if am[AttribItems] != nil {
		items := am[AttribItems].([]map[string]interface{})
		for i := 0; i < len(items); i++ {
			item := items[i]
			child, err := b.build(item, pan)
			if err != nil {
				return nil, err
			}
			pan.Add(child)
		}
	}
	return pan, nil
}

// buildImagePanel builds a gui object of type ImagePanel
func buildImagePanel(b *Builder, am map[string]interface{}) (IPanel, error) {

	// Checks imagefile attribute
	if am[AttribImageFile] == nil {
		return nil, b.err(am, AttribImageFile, "Must be supplied")
	}

	// If path is not absolute join with user supplied image base path
	imagefile := am[AttribImageFile].(string)
	if !filepath.IsAbs(imagefile) {
		imagefile = filepath.Join(b.imgpath, imagefile)
	}

	// Builds panel and set common attributes
	panel, err := NewImage(imagefile)
	if err != nil {
		return nil, err
	}
	err = b.SetAttribs(am, panel)
	if err != nil {
		return nil, err
	}

	// Sets optional AspectWidth attribute
	if aw := am[AttribAspectWidth]; aw != nil {
		panel.SetContentAspectWidth(aw.(float32))
	}

	// Sets optional AspectHeight attribute
	if ah := am[AttribAspectHeight]; ah != nil {
		panel.SetContentAspectHeight(ah.(float32))
	}

	// Builds children recursively
	if am[AttribItems] != nil {
		items := am[AttribItems].([]map[string]interface{})
		for i := 0; i < len(items); i++ {
			item := items[i]
			child, err := b.build(item, panel)
			if err != nil {
				return nil, err
			}
			panel.Add(child)
		}
	}
	return panel, nil
}

// buildLabel builds a gui object of type Label
func buildLabel(b *Builder, am map[string]interface{}) (IPanel, error) {

	var label *Label
	if am[AttribIcon] != nil {
		label = NewIcon(am[AttribIcon].(string))
	} else if am[AttribText] != nil {
		label = NewLabel(am[AttribText].(string))
	} else {
		label = NewLabel("")
	}

	// Sets common attributes
	err := b.SetAttribs(am, label)
	if err != nil {
		return nil, err
	}

	// Set optional background color
	if bgc := am[AttribBgColor]; bgc != nil {
		label.SetBgColor4(bgc.(*math32.Color4))
	}

	// Set optional font color
	if fc := am[AttribFontColor]; fc != nil {
		label.SetColor4(fc.(*math32.Color4))
	}

	// Sets optional font size
	if fs := am[AttribFontSize]; fs != nil {
		label.SetFontSize(float64(fs.(float32)))
	}

	// Sets optional font dpi
	if fdpi := am[AttribFontDPI]; fdpi != nil {
		label.SetFontDPI(float64(fdpi.(float32)))
	}

	// Sets optional line spacing
	if ls := am[AttribLineSpacing]; ls != nil {
		label.SetLineSpacing(float64(ls.(float32)))
	}

	return label, nil
}

// buildImageLabel builds a gui object of type: ImageLabel
func buildImageLabel(b *Builder, am map[string]interface{}) (IPanel, error) {

	// Builds image label and set common attributes
	var text string
	if am[AttribText] != nil {
		text = am[AttribText].(string)
	}
	imglabel := NewImageLabel(text)
	err := b.SetAttribs(am, imglabel)
	if err != nil {
		return nil, err
	}

	// Sets optional icon(s)
	if icon := am[AttribIcon]; icon != nil {
		imglabel.SetIcon(icon.(string))
	}

	// Sets optional image from file
	// If path is not absolute join with user supplied image base path
	if imgf := am[AttribImageFile]; imgf != nil {
		path := imgf.(string)
		if !filepath.IsAbs(path) {
			path = filepath.Join(b.imgpath, path)
		}
		err := imglabel.SetImageFromFile(path)
		if err != nil {
			return nil, err
		}
	}

	return imglabel, nil
}

// buildButton builds a gui object of type: Button
func buildButton(b *Builder, am map[string]interface{}) (IPanel, error) {

	// Builds button and set commont attributes
	var text string
	if am[AttribText] != nil {
		text = am[AttribText].(string)
	}
	button := NewButton(text)
	err := b.SetAttribs(am, button)
	if err != nil {
		return nil, err
	}

	// Sets optional icon(s)
	if icon := am[AttribIcon]; icon != nil {
		button.SetIcon(icon.(string))
	}

	// Sets optional image from file
	// If path is not absolute join with user supplied image base path
	if imgf := am[AttribImageFile]; imgf != nil {
		path := imgf.(string)
		if !filepath.IsAbs(path) {
			path = filepath.Join(b.imgpath, path)
		}
		err := button.SetImage(path)
		if err != nil {
			return nil, err
		}
	}

	return button, nil
}

// buildEdit builds a gui object of type: "Edit"
func buildEdit(b *Builder, am map[string]interface{}) (IPanel, error) {

	// Builds button and set attributes
	var width float32
	var placeholder string
	if aw := am[AttribWidth]; aw != nil {
		width = aw.(float32)
	}
	if ph := am[AttribPlaceHolder]; ph != nil {
		placeholder = ph.(string)
	}
	edit := NewEdit(int(width), placeholder)
	err := b.SetAttribs(am, edit)
	if err != nil {
		return nil, err
	}
	return edit, nil
}

// buildCheckBox builds a gui object of type: CheckBox
func buildCheckBox(b *Builder, am map[string]interface{}) (IPanel, error) {

	// Builds check box and set commont attributes
	var text string
	if am[AttribText] != nil {
		text = am[AttribText].(string)
	}
	cb := NewCheckBox(text)
	err := b.SetAttribs(am, cb)
	if err != nil {
		return nil, err
	}

	// Sets optional checked value
	if checked := am[AttribChecked]; checked != nil {
		cb.SetValue(checked.(bool))
	}
	return cb, nil
}

// buildRadioButton builds a gui object of type: RadioButton
func buildRadioButton(b *Builder, am map[string]interface{}) (IPanel, error) {

	// Builds check box and set commont attributes
	var text string
	if am[AttribText] != nil {
		text = am[AttribText].(string)
	}
	rb := NewRadioButton(text)
	err := b.SetAttribs(am, rb)
	if err != nil {
		return nil, err
	}

	// Sets optional radio button group
	if gr := am[AttribGroup]; gr != nil {
		rb.SetGroup(gr.(string))
	}

	// Sets optional checked value
	if checked := am[AttribChecked]; checked != nil {
		rb.SetValue(checked.(bool))
	}
	return rb, nil
}

// buildVList builds a gui object of type: VList
func buildVList(b *Builder, am map[string]interface{}) (IPanel, error) {

	// Builds list and set commont attributes
	list := NewVList(0, 0)
	err := b.SetAttribs(am, list)
	if err != nil {
		return nil, err
	}

	// Builds children
	if am[AttribItems] != nil {
		items := am[AttribItems].([]map[string]interface{})
		for i := 0; i < len(items); i++ {
			item := items[i]
			child, err := b.build(item, list)
			if err != nil {
				return nil, err
			}
			list.Add(child)
		}
	}
	return list, nil
}

// buildHList builds a gui object of type: VList
func buildHList(b *Builder, am map[string]interface{}) (IPanel, error) {

	// Builds list and set commont attributes
	list := NewHList(0, 0)
	err := b.SetAttribs(am, list)
	if err != nil {
		return nil, err
	}

	// Builds children
	if am[AttribItems] != nil {
		items := am[AttribItems].([]map[string]interface{})
		for i := 0; i < len(items); i++ {
			item := items[i]
			child, err := b.build(item, list)
			if err != nil {
				return nil, err
			}
			list.Add(child)
		}
	}
	return list, nil
}

// buildDropDown builds a gui object of type: DropDown
func buildDropDown(b *Builder, am map[string]interface{}) (IPanel, error) {

	// If image label attribute defined use it, otherwise
	// uses default value.
	var imglabel *ImageLabel
	if iv := am[AttribImageLabel]; iv != nil {
		imgl := iv.(map[string]interface{})
		imgl[AttribType] = TypeImageLabel
		ipan, err := b.build(imgl, nil)
		if err != nil {
			return nil, err
		}
		imglabel = ipan.(*ImageLabel)
	} else {
		imglabel = NewImageLabel("")
	}

	// Builds drop down and set common attributes
	dd := NewDropDown(0, imglabel)
	err := b.SetAttribs(am, dd)
	if err != nil {
		return nil, err
	}

	// Builds children
	if am[AttribItems] != nil {
		items := am[AttribItems].([]map[string]interface{})
		for i := 0; i < len(items); i++ {
			item := items[i]
			child, err := b.build(item, dd)
			if err != nil {
				return nil, err
			}
			dd.Add(child.(*ImageLabel))
		}
	}
	return dd, nil
}

// buildMenu builds a gui object of type: Menu or MenuBar from the
// specified panel descriptor.
func buildMenu(b *Builder, am map[string]interface{}) (IPanel, error) {

	// Builds menu bar or menu
	var menu *Menu
	if am[AttribType].(string) == TypeMenuBar {
		menu = NewMenuBar()
	} else {
		menu = NewMenu()
	}

	// Only sets attribs for top level menus
	if pi := am[AttribParentInternal]; pi != nil {
		par := pi.(map[string]interface{})
		ptype := ""
		if ti := par[AttribType]; ti != nil {
			ptype = ti.(string)
		}
		if ptype != TypeMenu && ptype != TypeMenuBar {
			err := b.SetAttribs(am, menu)
			if err != nil {
				return nil, err
			}
		}
	}

	// Builds and adds menu items
	if am[AttribItems] != nil {
		items := am[AttribItems].([]map[string]interface{})
		for i := 0; i < len(items); i++ {
			// Get the item optional type and text
			item := items[i]
			itype := ""
			itext := ""
			if iv := item[AttribType]; iv != nil {
				itype = iv.(string)
			}
			if iv := item[AttribText]; iv != nil {
				itext = iv.(string)
			}
			// Item is another menu
			if itype == TypeMenu {
				subm, err := buildMenu(b, item)
				if err != nil {
					return nil, err
				}
				menu.AddMenu(itext, subm.(*Menu))
				continue
			}
			// Item is a separator
			if itext == TypeSeparator {
				menu.AddSeparator()
				continue
			}
			// Item must be a menu option
			mi := menu.AddOption(itext)
			// Set item optional icon(s)
			if icon := item[AttribIcon]; icon != nil {
				mi.SetIcon(icon.(string))
			}
			// Sets optional menu item shortcut
			if sci := item[AttribShortcut]; sci != nil {
				sc := sci.([]int)
				mi.SetShortcut(window.ModifierKey(sc[0]), window.Key(sc[1]))
			}
		}
	}
	return menu, nil
}

// buildSlider builds a gui object of type: HSlider or VSlider
func buildSlider(b *Builder, am map[string]interface{}) (IPanel, error) {

	// Builds horizontal or vertical slider
	var slider *Slider
	if am[AttribType].(string) == TypeHSlider {
		slider = NewHSlider(0, 0)
	} else {
		slider = NewVSlider(0, 0)
	}

	// Sets common attributes
	err := b.SetAttribs(am, slider)
	if err != nil {
		return nil, err
	}

	// Sets optional text
	if itext := am[AttribText]; itext != nil {
		slider.SetText(itext.(string))
	}
	// Sets optional scale factor
	if isf := am[AttribScaleFactor]; isf != nil {
		slider.SetScaleFactor(isf.(float32))
	}
	// Sets optional value
	if iv := am[AttribValue]; iv != nil {
		slider.SetValue(iv.(float32))
	}
	return slider, nil
}

// buildSplitter builds a gui object of type: HSplitterr or VSplitter
func buildSplitter(b *Builder, am map[string]interface{}) (IPanel, error) {

	// Builds horizontal or vertical splitter
	var splitter *Splitter
	if am[AttribType].(string) == TypeHSplitter {
		splitter = NewHSplitter(0, 0)
	} else {
		splitter = NewVSplitter(0, 0)
	}

	// Sets common attributes
	err := b.SetAttribs(am, splitter)
	if err != nil {
		return nil, err
	}

	// Sets optional split value
	if v := am[AttribSplit]; v != nil {
		splitter.SetSplit(v.(float32))
	}

	// Internal function to set each of the splitter's panel attributes and items
	setpan := func(attrib string, pan *Panel) error {

		// Get internal panel attributes
		ipattribs := am[attrib]
		if ipattribs == nil {
			return nil
		}
		pattr := ipattribs.(map[string]interface{})
		// Set panel attributes
		err := b.SetAttribs(pattr, pan)
		if err != nil {
			return nil
		}
		// Builds panel children
		if pattr[AttribItems] != nil {
			items := pattr[AttribItems].([]map[string]interface{})
			for i := 0; i < len(items); i++ {
				item := items[i]
				child, err := b.build(item, pan)
				if err != nil {
					return err
				}
				pan.Add(child)
			}
		}
		return nil
	}

	// Set optional splitter panel's attributes
	err = setpan(AttribPanel0, &splitter.P0)
	if err != nil {
		return nil, err
	}
	err = setpan(AttribPanel1, &splitter.P1)
	if err != nil {
		return nil, err
	}

	return splitter, nil
}

// buildTree builds a gui object of type: Tree
func buildTree(b *Builder, am map[string]interface{}) (IPanel, error) {

	// Builds tree and sets its common attributes
	tree := NewTree(0, 0)
	err := b.SetAttribs(am, tree)
	if err != nil {
		return nil, err
	}

	// Internal function to build tree nodes recursively
	var buildItems func(am map[string]interface{}, pnode *TreeNode) error
	buildItems = func(am map[string]interface{}, pnode *TreeNode) error {

		v := am[AttribItems]
		if v == nil {
			return nil
		}
		items := v.([]map[string]interface{})

		for i := 0; i < len(items); i++ {
			// Get the item type
			item := items[i]
			itype := ""
			if v := item[AttribType]; v != nil {
				itype = v.(string)
			}
			itext := ""
			if v := item[AttribText]; v != nil {
				itext = v.(string)
			}

			// Item is a tree node
			if itype == "" || itype == TypeTreeNode {
				var node *TreeNode
				if pnode == nil {
					node = tree.AddNode(itext)
				} else {
					node = pnode.AddNode(itext)
				}
				err := buildItems(item, node)
				if err != nil {
					return err
				}
				continue
			}
			// Other controls
			ipan, err := b.build(item, nil)
			if err != nil {
				return err
			}
			if pnode == nil {
				tree.Add(ipan)
			} else {
				pnode.Add(ipan)
			}
		}
		return nil
	}

	// Build nodes
	err = buildItems(am, nil)
	if err != nil {
		return nil, err
	}
	return tree, nil
}

// buildWindow builds a gui object of type: Window
func buildWindow(b *Builder, am map[string]interface{}) (IPanel, error) {

	// Builds window and sets its common attributes
	win := NewWindow(0, 0)
	err := b.SetAttribs(am, win)
	if err != nil {
		return nil, err
	}

	// Sets optional title
	if title := am[AttribTitle]; title != nil {
		win.SetTitle(title.(string))
	}

	// Set optional resizable borders
	if resiz := am[AttribResizeBorders]; resiz != nil {
		win.SetResizable(resiz.(bool))
	}

	// Builds window children
	if v := am[AttribItems]; v != nil {
		items := v.([]map[string]interface{})
		for i := 0; i < len(items); i++ {
			item := items[i]
			child, err := b.build(item, win)
			if err != nil {
				return nil, err
			}
			win.Add(child)
		}
	}
	return win, nil
}

// buildChart builds a gui object of type: Chart
func buildChart(b *Builder, am map[string]interface{}) (IPanel, error) {

	// Builds window and sets its common attributes
	chart := NewChart(0, 0)
	err := b.SetAttribs(am, chart)
	if err != nil {
		return nil, err
	}

	// Sets optional title
	if title := am[AttribTitle]; title != nil {
		chart.SetTitle(title.(string), 14)
	}

	// Sets x scale attibutes
	if v := am[AttribScalex]; v != nil {
		sx := v.(map[string]interface{})
		// Sets optional x scale margin
		if mx := sx[AttribMargin]; mx != nil {
			chart.SetMarginX(mx.(float32))
		}
		// Sets optional x scale format
		if fx := sx[AttribFormat]; fx != nil {
			chart.SetFormatX(fx.(string))
		}
		// Sets optional x scale font size
		if fsize := sx[AttribFontSize]; fsize != nil {
			chart.SetFontSizeX(float64(fsize.(float32)))
		}
		// Number of lines
		lines := 4
		if v := sx[AttribLines]; v != nil {
			lines = v.(int)
		}
		// Lines color
		color4 := math32.NewColor4("Black")
		if v := sx[AttribColor]; v != nil {
			color4 = v.(*math32.Color4)
		}
		color := color4.ToColor()
		chart.SetScaleX(lines, &color)
		// Range first
		firstX := float32(0)
		if v := sx[AttribFirstx]; v != nil {
			firstX = v.(float32)
		}
		// Range step
		stepX := float32(1)
		if v := sx[AttribStepx]; v != nil {
			stepX = v.(float32)
		}
		// Range count step
		countStepX := float32(1)
		if v := sx[AttribStepx]; v != nil {
			countStepX = v.(float32)
		}
		chart.SetRangeX(firstX, stepX, countStepX)
	}

	// Sets y scale attibutes
	if v := am[AttribScaley]; v != nil {
		sy := v.(map[string]interface{})
		// Sets optional y scale margin
		if my := sy[AttribMargin]; my != nil {
			chart.SetMarginY(my.(float32))
		}
		// Sets optional y scale format
		if fy := sy[AttribFormat]; fy != nil {
			chart.SetFormatY(fy.(string))
		}
		// Sets optional y scale font size
		if fsize := sy[AttribFontSize]; fsize != nil {
			chart.SetFontSizeY(float64(fsize.(float32)))
		}
		// Number of lines
		lines := 4
		if v := sy[AttribLines]; v != nil {
			lines = v.(int)
		}
		// Lines color
		color4 := math32.NewColor4("Black")
		if v := sy[AttribColor]; v != nil {
			color4 = v.(*math32.Color4)
		}
		color := color4.ToColor()
		chart.SetScaleY(lines, &color)
		// Range min
		rmin := float32(-10)
		if v := sy[AttribRangeMin]; v != nil {
			rmin = v.(float32)
		}
		// Range max
		rmax := float32(10)
		if v := sy[AttribRangeMax]; v != nil {
			rmax = v.(float32)
		}
		chart.SetRangeY(rmin, rmax)
		// Range auto
		if rauto := sy[AttribRangeAuto]; rauto != nil {
			chart.SetRangeYauto(v.(bool))
		}
	}

	return chart, nil
}

// buildTable builds a gui object of type: Table
func buildTable(b *Builder, am map[string]interface{}) (IPanel, error) {

	// Internal function to build a TableColumn from its attribute map
	buildTableCol := func(b *Builder, am map[string]interface{}) (*TableColumn, error) {
		tc := &TableColumn{}
		if iv := am[AttribId]; iv != nil {
			tc.Id = iv.(string)
		}
		if iv := am[AttribHeader]; iv != nil {
			tc.Header = iv.(string)
		}
		if iv := am[AttribWidth]; iv != nil {
			tc.Width = iv.(float32)
		}
		if iv := am[AttribMinwidth]; iv != nil {
			tc.Minwidth = iv.(float32)
		}
		if iv := am[AttribHidden]; iv != nil {
			tc.Hidden = iv.(bool)
		}
		if iv := am[AttribFormat]; iv != nil {
			tc.Format = iv.(string)
		}
		if iv := am[AttribExpand]; iv != nil {
			tc.Expand = iv.(float32)
		}
		if iv := am[AttribResize]; iv != nil {
			tc.Resize = iv.(bool)
		}
		if iv := am[AttribSortType]; iv != nil {
			tc.Sort = iv.(TableSortType)
		}
		return tc, nil
	}

	// Builds table columns array
	tableCols := []TableColumn{}
	if iv := am[AttribColumns]; iv != nil {
		cols := iv.([]map[string]interface{})
		for _, c := range cols {
			tc, err := buildTableCol(b, c)
			if err != nil {
				return nil, err
			}
			tableCols = append(tableCols, *tc)
		}
	}

	// Creates table and set common attributes
	table, err := NewTable(0, 0, tableCols)
	if err != nil {
		return nil, err
	}
	err = b.SetAttribs(am, table)
	if err != nil {
		return nil, err
	}

	// Sets optional show header attribute
	if show := am[AttribShowHeader]; show != nil {
		table.ShowHeader(show.(bool))
	}

	return table, nil
}

// buildTabBar builds a gui object of type: TabBare
func buildTabBar(b *Builder, am map[string]interface{}) (IPanel, error) {

	// Creates TabBar and set common attributes
	tabbar := NewTabBar(0, 0)
	err := b.SetAttribs(am, tabbar)
	if err != nil {
		return nil, err
	}
	v := am[AttribItems]

	// For each tab
	if v != nil {
		items := v.([]map[string]interface{})
		for _, item := range items {
			// Creates Tab
			text := ""
			if v := item[AttribText]; v != nil {
				text = v.(string)
			}
			tab := tabbar.AddTab(text)
			// Sets optional icon
			if v := item[AttribIcon]; v != nil {
				tab.SetIcon(v.(string))
			}
			// Sets optional image
			if v := item[AttribImageFile]; v != nil {
				// If path is not absolute join with user supplied image base path
				imagefile := v.(string)
				if !filepath.IsAbs(imagefile) {
					imagefile = filepath.Join(b.imgpath, imagefile)
				}
				err := tab.SetImage(imagefile)
				if err != nil {
					return nil, err
				}
			}
			// Sets content panel
			if v := item[AttribContent]; v != nil {
				am := v.(map[string]interface{})
				content, err := b.build(am, nil)
				if err != nil {
					return nil, err
				}
				tab.SetContent(content)
			}
			// Sets pinned state
			if v := item[AttribPinned]; v != nil {
				tab.SetPinned(v.(bool))
			}
		}
	}

	return tabbar, nil
}
