// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package gui

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/g3n/engine/gui/assets/icon"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/window"
	"gopkg.in/yaml.v2"
)

// Builder builds GUI objects from a declarative description in YAML format
type Builder struct {
	am       map[string]interface{}     // parsed map with gui object atttributes
	imgpath  string                     // base path for image panels files
	builders map[string]BuilderFunc     // map of builder functions by type
	attribs  map[string]AttribCheckFunc // map of attribute name with check functions
	layouts  map[string]IBuilderLayout  // map of layout type to layout builder
}

// IBuilderLayout is the interface for all layout builders
type IBuilderLayout interface {
	BuildLayout(b *Builder, am map[string]interface{}) (ILayout, error)
	BuildParams(b *Builder, am map[string]interface{}) (interface{}, error)
}

// BuilderFunc is type for functions which build a gui object from an attribute map
type BuilderFunc func(*Builder, map[string]interface{}) (IPanel, error)

// BuilderFunc is type for functions which builds a layout object from an attribute map
type LayoutFunc func(*Builder, map[string]interface{}) (ILayout, error)

//// descLayout contains all layout attributes
//type descLayout struct {
//	Type      string  // Type of the layout: HBox, VBox, Grid, Dock, others...
//	Cols      int     // Number of columns for Grid layout
//	Spacing   float32 // Spacing in pixels for HBox and VBox
//	AlignH    string  // HBox group alignment type
//	AlignV    string  // VBox group alignment type
//	MinHeight bool    // HBox, VBox minimum height flag
//	MinWidth  bool    // HBox, VBox minimum width flag
//	ExpandH   bool    // Grid
//	ExpandV   bool    // Grid
//}
//
//// descLayoutParam describes all layout parameters types
//type descLayoutParams struct {
//	Expand  *float32 // HBox, VBox expand factor
//	ColSpan int      // Grid layout colspan
//	AlignH  string   // horizontal alignment
//	AlignV  string   // vertical alignment
//	Edge    string   // Dock layout edge: top,right,bottom,left,center
//}

//// descPanel describes all panel attributes
//type descPanel struct {
//	Type         string            // Gui object type: Panel, Label, Edit, etc ...
//	Name         string            // Optional name for identification
//	Position     string            // Optional position as: x y | x,y
//	Width        *float32          // Optional width (default = 0)
//	Height       *float32          // Optional height (default = 0)
//	AspectWidth  *float32          // Optional aspectwidth (default = nil)
//	AspectHeight *float32          // Optional aspectwidth (default = nil)
//	Margins      string            // Optional margins as 1 or 4 float values
//	Borders      string            // Optional borders as 1 or 4 float values
//	BorderColor  string            // Optional border color as name or 3 or 4 float values
//	Paddings     string            // Optional paddings as 1 or 4 float values
//	Color        string            // Optional color as 1 or 4 float values
//	Enabled      *bool             // All:
//	Visible      *bool             // All:
//	Renderable   *bool             // All:
//	Imagefile    string            // For Panel, Button
//	Layout       *descLayout       // Optional pointer to layout
//	LayoutParams *descLayoutParams // Optional layout parameters
//	Text         string            // Label, Button
//	Icons        string            // Label
//	BgColor      string            // Label
//	FontColor    string            // Label
//	FontSize     *float32          // Label
//	FontDPI      *float32          // Label
//	LineSpacing  *float32          // Label
//	PlaceHolder  string            // Edit
//	MaxLength    *uint             // Edit
//	Icon         string            // Button
//	Group        string            // RadioButton
//	Checked      bool              // CheckBox, RadioButton
//	ImageLabel   *descPanel        // DropDown
//	Items        []*descPanel      // Menu, MenuBar
//	Shortcut     string            // Menu
//	Value        *float32          // Slider
//	ScaleFactor  *float32          // Slider
//	Title        string            // Window
//	Resizable    string            // Window resizable borders
//	P0           *descPanel        // Splitter panel 0
//	P1           *descPanel        // Splitter panel 1
//	Split        *float32          // Splitter split value
//	parent       *descPanel        // used internally
//}

// Panel and layout types
const (
	TypePanel       = "panel"
	TypeImagePanel  = "imagepanel"
	TypeLabel       = "label"
	TypeImageLabel  = "imagelabel"
	TypeButton      = "button"
	TypeCheckBox    = "checkbox"
	TypeRadioButton = "radiobutton"
	TypeEdit        = "edit"
	TypeVList       = "vlist"
	TypeHList       = "hlist"
	TypeDropDown    = "dropdown"
	TypeHSlider     = "hslider"
	TypeVSlider     = "vslider"
	TypeHSplitter   = "hsplitter"
	TypeVSplitter   = "vsplitter"
	TypeSeparator   = "separator"
	TypeTree        = "tree"
	TypeTreeNode    = "node"
	TypeMenuBar     = "menubar"
	TypeMenu        = "menu"
	TypeWindow      = "window"
	TypeHBoxLayout  = "hbox"
	TypeVBoxLayout  = "vbox"
	TypeGridLayout  = "grid"
	TypeDockLayout  = "dock"
)

// Common attribute names
const (
	AttribAlignv       = "alignv"       // Align
	AttribAlignh       = "alignh"       // Align
	AttribAspectHeight = "aspectheight" // float32
	AttribAspectWidth  = "aspectwidth"  // float32
	AttribBgColor      = "bgcolor"      // Color4
	AttribBorders      = "borders"      // BorderSizes
	AttribBorderColor  = "bordercolor"  // Color4
	AttribChecked      = "checked"      // bool
	AttribColor        = "color"        // Color4
	AttribCols         = "cols"         // Int
	AttribColSpan      = "colspan"      // Int
	AttribEdge         = "edge"         // int
	AttribEnabled      = "enabled"      // bool
	AttribExpand       = "expand"       // float32
	AttribExpandh      = "expandh"      // bool
	AttribExpandv      = "expandv"      // bool
	AttribFontColor    = "fontcolor"    // Color4
	AttribFontDPI      = "fontdpi"      // float32
	AttribFontSize     = "fontsize"     // float32
	AttribGroup        = "group"        // string
	AttribHeight       = "height"       // float32
	AttribIcon         = "icon"         // string
	AttribImageFile    = "imagefile"    // string
	AttribImageLabel   = "imagelabel"   // []map[string]interface{}
	AttribItems        = "items"        // []map[string]interface{}
	AttribLayout       = "layout"       // map[string]interface{}
	AttribLayoutParams = "layoutparams" // map[string]interface{}
	AttribLineSpacing  = "linespacing"  // float32
	AttribMinHeight    = "minheight"    // bool
	AttribMinWidth     = "minwidth"     // bool
	AttribMargins      = "margins"      // BorderSizes
	AttribName         = "name"         // string
	AttribPaddings     = "paddings"     // BorderSizes
	AttribPanel0       = "panel0"       // map[string]interface{}
	AttribPanel1       = "panel1"       // map[string]interface{}
	AttribParent_      = "parent_"      // string (internal attribute)
	AttribPlaceHolder  = "placeholder"  // string
	AttribPosition     = "position"     // []float32
	AttribRender       = "render"       // bool
	AttribScaleFactor  = "scalefactor"  // float32
	AttribShortcut     = "shortcut"     // []int
	AttribSpacing      = "spacing"      // float32
	AttribText         = "text"         // string
	AttribType         = "type"         // string
	AttribWidth        = "width"        // float32
	AttribValue        = "value"        // float32
	AttribVisible      = "visible"      // bool
)

const (
	aPOS         = 1 << iota                                  // attribute position
	aSIZE        = 1 << iota                                  // attribute size
	aNAME        = 1 << iota                                  // attribute name
	aMARGINS     = 1 << iota                                  // attribute margins widths
	aBORDERS     = 1 << iota                                  // attribute borders widths
	aBORDERCOLOR = 1 << iota                                  // attribute border color
	aPADDINGS    = 1 << iota                                  // attribute paddings widths
	aCOLOR       = 1 << iota                                  // attribute panel bgcolor
	aENABLED     = 1 << iota                                  // attribute enabled for events
	aRENDER      = 1 << iota                                  // attribute renderable
	aVISIBLE     = 1 << iota                                  // attribute visible
	asPANEL      = 0xFF                                       // attribute set for panels
	asWIDGET     = aPOS | aNAME | aSIZE | aENABLED | aVISIBLE // attribute set for widgets
)

// maps align name with align parameter
var mapAlignh = map[string]Align{
	"none":   AlignNone,
	"left":   AlignLeft,
	"right":  AlignRight,
	"width":  AlignWidth,
	"center": AlignCenter,
}

// maps align name with align parameter
var mapAlignv = map[string]Align{
	"none":   AlignNone,
	"top":    AlignTop,
	"bottom": AlignBottom,
	"height": AlignHeight,
	"center": AlignCenter,
}

// maps edge name (dock layout) with edge parameter
var mapEdgeName = map[string]int{
	"top":    DockTop,
	"right":  DockRight,
	"bottom": DockBottom,
	"left":   DockLeft,
	"center": DockCenter,
}

// maps resize border name (window) with parameter value
var mapResizable = map[string]Resizable{
	"top":    ResizeTop,
	"right":  ResizeRight,
	"bottom": ResizeBottom,
	"left":   ResizeLeft,
	"all":    ResizeAll,
}

type AttribCheckFunc func(b *Builder, am map[string]interface{}, fname string) error

// NewBuilder creates and returns a pointer to a new gui Builder object
func NewBuilder() *Builder {

	b := new(Builder)
	// Sets map of object type to builder function
	b.builders = map[string]BuilderFunc{
		TypePanel:       buildPanel,
		TypeImagePanel:  buildImagePanel,
		TypeLabel:       buildLabel,
		TypeImageLabel:  buildImageLabel,
		TypeButton:      buildButton,
		TypeEdit:        buildEdit,
		TypeCheckBox:    buildCheckBox,
		TypeRadioButton: buildRadioButton,
		TypeVList:       buildVList,
		TypeHList:       buildHList,
		TypeDropDown:    buildDropDown,
		TypeMenu:        buildMenu,
		TypeMenuBar:     buildMenu,
		TypeHSlider:     buildSlider,
		TypeVSlider:     buildSlider,
		TypeHSplitter:   buildSplitter,
		TypeVSplitter:   buildSplitter,
		TypeTree:        buildTree,
	}
	// Sets map of layout type name to layout function
	b.layouts = map[string]IBuilderLayout{
		TypeHBoxLayout: &BuilderLayoutHBox{},
		TypeVBoxLayout: &BuilderLayoutVBox{},
		TypeGridLayout: &BuilderLayoutGrid{},
		TypeDockLayout: &BuilderLayoutDock{},
	}
	// Sets map of attribute name to check function
	b.attribs = map[string]AttribCheckFunc{
		AttribAlignv:       AttribCheckAlign,
		AttribAlignh:       AttribCheckAlign,
		AttribAspectWidth:  AttribCheckFloat,
		AttribAspectHeight: AttribCheckFloat,
		AttribHeight:       AttribCheckFloat,
		AttribMargins:      AttribCheckBorderSizes,
		AttribBgColor:      AttribCheckColor,
		AttribBorders:      AttribCheckBorderSizes,
		AttribBorderColor:  AttribCheckColor,
		AttribChecked:      AttribCheckBool,
		AttribColor:        AttribCheckColor,
		AttribCols:         AttribCheckInt,
		AttribColSpan:      AttribCheckInt,
		AttribEdge:         AttribCheckEdge,
		AttribEnabled:      AttribCheckBool,
		AttribExpand:       AttribCheckFloat,
		AttribExpandh:      AttribCheckBool,
		AttribExpandv:      AttribCheckBool,
		AttribFontColor:    AttribCheckColor,
		AttribFontDPI:      AttribCheckFloat,
		AttribFontSize:     AttribCheckFloat,
		AttribGroup:        AttribCheckString,
		AttribIcon:         AttribCheckIcons,
		AttribImageFile:    AttribCheckString,
		AttribImageLabel:   AttribCheckMap,
		AttribItems:        AttribCheckListMap,
		AttribLayout:       AttribCheckLayout,
		AttribLayoutParams: AttribCheckMap,
		AttribLineSpacing:  AttribCheckFloat,
		AttribMinHeight:    AttribCheckBool,
		AttribMinWidth:     AttribCheckBool,
		AttribName:         AttribCheckString,
		AttribPaddings:     AttribCheckBorderSizes,
		AttribPanel0:       AttribCheckMap,
		AttribPanel1:       AttribCheckMap,
		AttribPlaceHolder:  AttribCheckString,
		AttribPosition:     AttribCheckPosition,
		AttribRender:       AttribCheckBool,
		AttribScaleFactor:  AttribCheckFloat,
		AttribShortcut:     AttribCheckMenuShortcut,
		AttribSpacing:      AttribCheckFloat,
		AttribText:         AttribCheckString,
		AttribType:         AttribCheckStringLower,
		AttribValue:        AttribCheckFloat,
		AttribVisible:      AttribCheckBool,
		AttribWidth:        AttribCheckFloat,
	}
	return b
}

// ParseString parses a string with gui objects descriptions in YAML format
// It there was a previously parsed description, it is cleared.
func (b *Builder) ParseString(desc string) error {

	// Parses descriptor string in YAML format saving result in
	// a map of interface{} to interface{} as YAML allows numeric keys.
	var mii map[interface{}]interface{}
	err := yaml.Unmarshal([]byte(desc), &mii)
	if err != nil {
		return err
	}

	// If all the values of the top level map keys are other maps,
	// then it is a description of several objects, otherwise it is
	// a description of a single object.
	single := false
	for _, v := range mii {
		_, ok := v.(map[interface{}]interface{})
		if !ok {
			single = true
			break
		}
	}
	log.Error("single:%v", single)

	// Internal function which converts map[interface{}]interface{} to
	// map[string]interface{} recursively and lower case of all map keys.
	// It also sets a field named "parent_", which pointer to the parent map
	// This field causes a circular reference in the result map which prevents
	// the use of Go's Printf to print the result map.
	var visitor func(v, par interface{}) (interface{}, error)
	visitor = func(v, par interface{}) (interface{}, error) {

		switch vt := v.(type) {
		case []interface{}:
			ls := []interface{}{}
			for _, item := range vt {
				ci, err := visitor(item, par)
				if err != nil {
					return nil, err
				}
				ls = append(ls, ci)
			}
			return ls, nil

		case map[interface{}]interface{}:
			ms := make(map[string]interface{})
			for k, v := range vt {
				// Checks key
				ks, ok := k.(string)
				if !ok {
					return nil, fmt.Errorf("Keys must be strings")
				}
				ks = strings.ToLower(ks)
				// Checks value
				vi, err := visitor(v, ms)
				if err != nil {
					return nil, err
				}
				ms[ks] = vi
				// If has panel has parent or is a single top level panel, checks attributes
				if par != nil || single {
					// Get attribute check function
					acf, ok := b.attribs[ks]
					if !ok {
						return nil, fmt.Errorf("Invalid attribute:%s", ks)
					}
					// Checks attribute
					err = acf(b, ms, ks)
					if err != nil {
						return nil, err
					}
				}
			}
			if par != nil {
				ms[AttribParent_] = par
			}
			return ms, nil

		default:
			return v, nil
		}
		return nil, nil
	}

	// Get map[string]interface{} with lower case keys from parsed descritor
	res, err := visitor(mii, nil)
	if err != nil {
		return err
	}
	msi, ok := res.(map[string]interface{})
	if !ok {
		return fmt.Errorf("Parsed result is not a map")
	}
	b.am = msi
	//b.debugPrint(b.am, 1)
	return nil
}

// ParseFile parses a file with gui objects descriptions in YAML format
// It there was a previously parsed description, it is cleared.
func (b *Builder) ParseFile(filepath string) error {

	// Reads all file data
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}

	// Parses file data
	return b.ParseString(string(data))
}

// Names returns a sorted list of names of top level previously parsed objects.
// Only objects with defined types are returned.
// If there is only a single object with no name, its name is returned
// as an empty string
func (b *Builder) Names() []string {

	var objs []string
	if b.am[AttribType] != nil {
		objs = append(objs, "")
		return objs
	}
	for name, _ := range b.am {
		objs = append(objs, name)
	}
	sort.Strings(objs)
	return objs
}

// Build builds a gui object and all its children recursively.
// The specified name should be a top level name from a
// from a previously parsed description
// If the descriptions contains a single object with no name,
// It should be specified the empty string to build this object.
func (b *Builder) Build(name string) (IPanel, error) {

	// Only one object
	if name == "" {
		log.Error("TYPE:---------->%T", b.am)
		return b.build(b.am, nil)
	}
	// Map of gui objects
	am, ok := b.am[name]
	if !ok {
		return nil, fmt.Errorf("Object name:%s not found", name)
	}
	return b.build(am.(map[string]interface{}), nil)
}

// Sets the path for image panels relative image files
func (b *Builder) SetImagepath(path string) {

	b.imgpath = path
}

func (b *Builder) AddBuilder(typename string, bf BuilderFunc) {

	b.builders[typename] = bf
}

// build builds the gui object from the specified description.
// All its children are also built recursively
// Returns the built object or an error
func (b *Builder) build(am map[string]interface{}, iparent IPanel) (IPanel, error) {

	// Get panel type
	itype := am[AttribType]
	if itype == nil {
		return nil, fmt.Errorf("Type not specified")
	}
	typename := itype.(string)

	// Get builder function for this type name
	builder := b.builders[typename]
	if builder == nil {
		return nil, fmt.Errorf("Invalid type:%v", typename)
	}

	// Builds panel
	pan, err := builder(b, am)
	if err != nil {
		return nil, err
	}
	// Adds built panel to parent
	if iparent != nil {
		iparent.GetPanel().Add(pan)
	}
	return pan, nil
}

// buildPanel builds an object of type Panel
func buildPanel(b *Builder, am map[string]interface{}) (IPanel, error) {

	pan := NewPanel(0, 0)
	err := b.setAttribs(am, pan, asPANEL)
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
	err = b.setAttribs(am, panel, asPANEL)
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
		label = NewLabel(am[AttribIcon].(string), true)
	} else if am[AttribText] != nil {
		label = NewLabel(am[AttribText].(string))
	} else {
		label = NewLabel("")
	}

	// Sets common attributes
	err := b.setAttribs(am, label, asPANEL)
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
	err := b.setAttribs(am, imglabel, asPANEL)
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
	err := b.setAttribs(am, button, asWIDGET)
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
	err := b.setAttribs(am, edit, asWIDGET)
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
	err := b.setAttribs(am, cb, asWIDGET)
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
	err := b.setAttribs(am, rb, asWIDGET)
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
	err := b.setAttribs(am, list, asWIDGET)
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
	err := b.setAttribs(am, list, asWIDGET)
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
	err := b.setAttribs(am, dd, asWIDGET)
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
	if pi := am[AttribParent_]; pi != nil {
		par := pi.(map[string]interface{})
		ptype := ""
		if ti := par[AttribType]; ti != nil {
			ptype = ti.(string)
		}
		if ptype != TypeMenu && ptype != TypeMenuBar {
			err := b.setAttribs(am, menu, asWIDGET)
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
	err := b.setAttribs(am, slider, asWIDGET)
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
	err := b.setAttribs(am, splitter, asWIDGET)
	if err != nil {
		return nil, err
	}

	// Sets optional split value
	if iv := am[AttribValue]; iv != nil {
		splitter.SetSplit(iv.(float32))
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
		err := b.setAttribs(pattr, pan, asPANEL)
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
	err := b.setAttribs(am, tree, asWIDGET)
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

// setLayout sets the optional layout of the specified panel
func (b *Builder) setLayout(am map[string]interface{}, ipan IPanel) error {

	// Get layout type
	lai := am[AttribLayout]
	if lai == nil {
		return nil
	}
	lam := lai.(map[string]interface{})
	ltype := lam[AttribType]
	if ltype == nil {
		return b.err(am, AttribType, "Layout must have a type")
	}

	// Get layout builder
	lbuilder := b.layouts[ltype.(string)]
	if lbuilder == nil {
		return b.err(am, AttribType, "Invalid layout type")
	}

	// Builds layout builder and set to panel
	layout, err := lbuilder.BuildLayout(b, lam)
	if err != nil {
		return err
	}
	ipan.SetLayout(layout)
	return nil
}

func (b *Builder) setLayoutParams(am map[string]interface{}, ipan IPanel) error {

	// Get layout params attributes
	lpi := am[AttribLayoutParams]
	if lpi == nil {
		return nil
	}
	lp := lpi.(map[string]interface{})

	// Get layout type from parent
	pi := am[AttribParent_]
	if pi == nil {
		return b.err(am, AttribType, "Panel has no parent")
	}
	par := pi.(map[string]interface{})
	v := par[AttribLayout]
	if v == nil {
		return nil
	}
	playout := v.(map[string]interface{})
	pltype := playout[AttribType].(string)

	// Get layout builder and builds layout params
	lbuilder := b.layouts[pltype]
	params, err := lbuilder.BuildParams(b, lp)
	if err != nil {
		return err
	}
	ipan.GetPanel().SetLayoutParams(params)
	return nil
}

func AttribCheckEdge(b *Builder, am map[string]interface{}, fname string) error {

	v := am[fname]
	if v == nil {
		return nil
	}
	vs, ok := v.(string)
	if !ok {
		return b.err(am, fname, "Invalid edge name")
	}
	edge, ok := mapEdgeName[vs]
	if !ok {
		return b.err(am, fname, "Invalid edge name")
	}
	am[fname] = edge
	return nil
}

func AttribCheckLayout(b *Builder, am map[string]interface{}, fname string) error {

	v := am[fname]
	if v == nil {
		return nil
	}
	msi, ok := v.(map[string]interface{})
	if !ok {
		return b.err(am, fname, "Not a map")
	}
	lti := msi[AttribType]
	if lti == nil {
		return b.err(am, fname, "Layout must have a type")
	}
	lfunc := b.layouts[lti.(string)]
	if lfunc == nil {
		return b.err(am, fname, "Invalid layout type")
	}
	return nil
}

func AttribCheckAlign(b *Builder, am map[string]interface{}, fname string) error {

	v := am[fname]
	if v == nil {
		return nil
	}
	vs, ok := v.(string)
	if !ok {
		return b.err(am, fname, "Invalid alignment")
	}
	var align Align
	if fname == AttribAlignh {
		align, ok = mapAlignh[vs]
	} else {
		align, ok = mapAlignv[vs]
	}
	if !ok {
		return b.err(am, fname, "Invalid alignment")
	}
	am[fname] = align
	return nil
}

func AttribCheckMenuShortcut(b *Builder, am map[string]interface{}, fname string) error {

	v := am[fname]
	if v == nil {
		return nil
	}
	vs, ok := v.(string)
	if !ok {
		return b.err(am, fname, "Not a string")
	}
	sc := strings.Trim(vs, " ")
	if sc == "" {
		return nil
	}
	parts := strings.Split(sc, "+")
	var mods window.ModifierKey
	for i := 0; i < len(parts)-1; i++ {
		switch parts[i] {
		case "Shift":
			mods |= window.ModShift
		case "Ctrl":
			mods |= window.ModControl
		case "Alt":
			mods |= window.ModAlt
		default:
			return b.err(am, fname, "Invalid shortcut:"+sc)
		}
	}
	// The last part must be a key
	keyname := parts[len(parts)-1]
	var keycode int
	found := false
	for kcode, kname := range mapKeyText {
		if kname == keyname {
			keycode = int(kcode)
			found = true
			break
		}
	}
	if !found {
		return b.err(am, fname, "Invalid shortcut:"+sc)
	}
	am[fname] = []int{int(mods), keycode}
	return nil
}

func AttribCheckListMap(b *Builder, am map[string]interface{}, fname string) error {

	v := am[fname]
	if v == nil {
		return nil
	}
	li, ok := v.([]interface{})
	if !ok {
		return b.err(am, fname, "Not a list")
	}
	lmsi := make([]map[string]interface{}, 0)
	for i := 0; i < len(li); i++ {
		item := li[i]
		msi, ok := item.(map[string]interface{})
		if !ok {
			return b.err(am, fname, "Item is not a map")
		}
		lmsi = append(lmsi, msi)
	}
	am[fname] = lmsi
	return nil
}

func AttribCheckMap(b *Builder, am map[string]interface{}, fname string) error {

	v := am[fname]
	if v == nil {
		return nil
	}
	msi, ok := v.(map[string]interface{})
	if !ok {
		return b.err(am, fname, "Not a map")
	}
	am[fname] = msi
	return nil
}

func AttribCheckIcons(b *Builder, am map[string]interface{}, fname string) error {

	v := am[fname]
	if v == nil {
		return nil
	}
	fs, ok := v.(string)
	if !ok {
		return b.err(am, fname, "Not a string")
	}
	text := ""
	parts := strings.Fields(fs)
	for i := 0; i < len(parts); i++ {
		// Try name first
		cp := icon.Codepoint(parts[i])
		if cp != "" {
			text += string(cp)
			continue
		}

		// Try to parse as hex value
		val, err := strconv.ParseUint(parts[i], 16, 32)
		if err != nil {
			return b.err(am, fname, fmt.Sprintf("Invalid icon codepoint value/name:%v", parts[i]))
		}
		text += string(val)
	}
	am[fname] = text
	return nil
}

func AttribCheckColor(b *Builder, am map[string]interface{}, fname string) error {

	// Checks if field is nil
	v := am[fname]
	if v == nil {
		return nil
	}

	// Converts to string
	fs, ok := v.(string)
	if !ok {
		return b.err(am, fname, "Not a string")
	}

	// Checks if string field is empty
	fs = strings.Trim(fs, " ")
	if fs == "" {
		return nil
	}

	// If string has 1 or 2 fields it must be a color name and optional alpha
	parts := strings.Fields(fs)
	if len(parts) == 1 || len(parts) == 2 {
		// First part must be a color name
		c, ok := math32.IsColorName(parts[0])
		if !ok {
			return b.err(am, fname, fmt.Sprintf("Invalid color name:%s", parts[0]))
		}
		c4 := math32.Color4{c.R, c.G, c.B, 1}
		if len(parts) == 2 {
			val, err := strconv.ParseFloat(parts[1], 32)
			if err != nil {
				return b.err(am, fname, fmt.Sprintf("Invalid float32 value:%s", parts[1]))
			}
			c4.A = float32(val)
		}
		am[fname] = &c4
		return nil
	}

	// Accept 3 or 4 floats values
	va, err := b.parseFloats(am, fname, 3, 4)
	if err != nil {
		return err
	}
	if len(va) == 3 {
		am[fname] = &math32.Color4{va[0], va[1], va[2], 1}
		return nil
	}
	am[fname] = &math32.Color4{va[0], va[1], va[2], va[3]}
	return nil
}

func AttribCheckBorderSizes(b *Builder, am map[string]interface{}, fname string) error {

	va, err := b.parseFloats(am, fname, 1, 4)
	if err != nil {
		return err
	}
	if va == nil {
		return nil
	}
	if len(va) == 1 {
		am[fname] = &BorderSizes{va[0], va[0], va[0], va[0]}
		return nil
	}
	am[fname] = &BorderSizes{va[0], va[1], va[2], va[3]}
	return nil
}

func AttribCheckPosition(b *Builder, am map[string]interface{}, fname string) error {

	v := am[fname]
	if v == nil {
		return nil
	}
	af, err := b.parseFloats(am, fname, 2, 2)
	if err != nil {
		return err
	}
	am[fname] = af
	return nil
}

func AttribCheckStringLower(b *Builder, am map[string]interface{}, fname string) error {

	err := AttribCheckString(b, am, fname)
	if err != nil {
		return err
	}
	if v := am[fname]; v != nil {
		am[fname] = strings.ToLower(v.(string))
	}
	return nil
}

func AttribCheckFloat(b *Builder, am map[string]interface{}, fname string) error {

	v := am[fname]
	if v == nil {
		return nil
	}
	switch n := v.(type) {
	case int:
		am[fname] = float32(n)
		return nil
	case float64:
		am[fname] = float32(n)
		return nil
	default:
		return b.err(am, fname, fmt.Sprintf("Not a number:%T", v))
	}
	return nil
}

func AttribCheckInt(b *Builder, am map[string]interface{}, fname string) error {

	v := am[fname]
	if v == nil {
		return nil
	}
	vint, ok := v.(int)
	if !ok {
		return b.err(am, fname, "Not an integer")
	}
	am[fname] = vint
	return nil
}

func AttribCheckString(b *Builder, am map[string]interface{}, fname string) error {

	v := am[fname]
	if v == nil {
		return nil
	}
	s, ok := v.(string)
	if !ok {
		return b.err(am, fname, "Not a string")
	}
	am[fname] = s
	return nil
}

func AttribCheckBool(b *Builder, am map[string]interface{}, fname string) error {

	v := am[fname]
	if v == nil {
		return nil
	}
	bv, ok := v.(bool)
	if !ok {
		return b.err(am, fname, "Not a bool")
	}
	am[fname] = bv
	return nil
}

/***


// buildImageLabel builds a gui object of type: ImageLabel
func (b *Builder) buildImageLabel(pd *descPanel) (IPanel, error) {

	// Builds image label and set common attributes
	imglabel := NewImageLabel(pd.Text)
	err := b.setAttribs(pd, imglabel, asPANEL)
	if err != nil {
		return nil, err
	}

	// Sets optional icon(s)
	icons, err := b.parseIconNames("icons", pd.Icons)
	if err != nil {
		return nil, err
	}
	if icons != "" {
		imglabel.SetIcon(icons)
	}

	// Sets optional image from file
	// If path is not absolute join with user supplied image base path
	if pd.Imagefile != "" {
		path := pd.Imagefile
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
func (b *Builder) buildButton(pd *descPanel) (IPanel, error) {

	// Builds button and set commont attributes
	button := NewButton(pd.Text)
	err := b.setAttribs(pd, button, asBUTTON)
	if err != nil {
		return nil, err
	}

	// Sets optional icon
	if pd.Icon != "" {
		cp, err := b.parseIconName("icon", pd.Icon)
		if err != nil {
			return nil, err
		}
		button.SetIcon(cp)
	}

	// Sets optional image from file
	// If path is not absolute join with user supplied image base path
	if pd.Imagefile != "" {
		path := pd.Imagefile
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

// buildCheckBox builds a gui object of type: CheckBox
func (b *Builder) buildCheckBox(pd *descPanel) (IPanel, error) {

	// Builds check box and set commont attributes
	cb := NewCheckBox(pd.Text)
	err := b.setAttribs(pd, cb, asWIDGET)
	if err != nil {
		return nil, err
	}
	cb.SetValue(pd.Checked)
	return cb, nil
}

// buildRadioButton builds a gui object of type: RadioButton
func (b *Builder) buildRadioButton(pd *descPanel) (IPanel, error) {

	// Builds check box and set commont attributes
	rb := NewRadioButton(pd.Text)
	err := b.setAttribs(pd, rb, asWIDGET)
	if err != nil {
		return nil, err
	}

	// Sets optional radio button group
	if pd.Group != "" {
		rb.SetGroup(pd.Group)
	}
	rb.SetValue(pd.Checked)
	return rb, nil
}

// buildEdit builds a gui object of type: "Edit"
func (b *Builder) buildEdit(dp *descPanel) (IPanel, error) {

	// Builds button and set attributes
	width, _ := b.size(dp)
	edit := NewEdit(int(width), dp.PlaceHolder)
	err := b.setAttribs(dp, edit, asWIDGET)
	if err != nil {
		return nil, err
	}
	edit.SetText(dp.Text)
	return edit, nil
}

// buildVList builds a gui object of type: VList
func (b *Builder) buildVList(dp *descPanel) (IPanel, error) {

	// Builds list and set commont attributes
	width, height := b.size(dp)
	list := NewVList(width, height)
	err := b.setAttribs(dp, list, asWIDGET)
	if err != nil {
		return nil, err
	}

	// Builds list children
	for i := 0; i < len(dp.Items); i++ {
		item := dp.Items[i]
		b.objpath.push(item.Name)
		child, err := b.build(item, list)
		b.objpath.pop()
		if err != nil {
			return nil, err
		}
		list.Add(child)
	}
	return list, nil
}

// buildHList builds a gui object of type: VList
func (b *Builder) buildHList(dp *descPanel) (IPanel, error) {

	// Builds list and set commont attributes
	width, height := b.size(dp)
	list := NewHList(width, height)
	err := b.setAttribs(dp, list, asWIDGET)
	if err != nil {
		return nil, err
	}

	// Builds list children
	for i := 0; i < len(dp.Items); i++ {
		item := dp.Items[i]
		b.objpath.push(item.Name)
		child, err := b.build(item, list)
		b.objpath.pop()
		if err != nil {
			return nil, err
		}
		list.Add(child)
	}
	return list, nil
}

// buildDropDown builds a gui object of type: DropDown
func (b *Builder) buildDropDown(pd *descPanel) (IPanel, error) {

	// If image label attribute defined use it, otherwise
	// uses default value.
	var imglabel *ImageLabel
	if pd.ImageLabel != nil {
		pd.ImageLabel.Type = descTypeImageLabel
		ipan, err := b.build(pd.ImageLabel, nil)
		if err != nil {
			return nil, err
		}
		imglabel = ipan.(*ImageLabel)
	} else {
		imglabel = NewImageLabel("")
	}

	// Builds drop down and set common attributes
	width, _ := b.size(pd)
	dd := NewDropDown(width, imglabel)
	err := b.setAttribs(pd, dd, asWIDGET)
	if err != nil {
		return nil, err
	}

	// Builds drop down children
	for i := 0; i < len(pd.Items); i++ {
		item := pd.Items[i]
		item.Type = descTypeImageLabel
		b.objpath.push(item.Name)
		child, err := b.build(item, dd)
		b.objpath.pop()
		if err != nil {
			return nil, err
		}
		dd.Add(child.(*ImageLabel))
	}
	return dd, nil
}

// buildSlider builds a gui object of type: HSlider or VSlider
func (b *Builder) buildSlider(pd *descPanel, horiz bool) (IPanel, error) {

	// Builds slider and sets its position
	width, height := b.size(pd)
	var slider *Slider
	if horiz {
		slider = NewHSlider(width, height)
	} else {
		slider = NewVSlider(width, height)
	}
	err := b.setAttribs(pd, slider, asWIDGET)
	if err != nil {
		return nil, err
	}

	// Sets optional text
	if pd.Text != "" {
		slider.SetText(pd.Text)
	}
	// Sets optional scale factor
	if pd.ScaleFactor != nil {
		slider.SetScaleFactor(*pd.ScaleFactor)
	}
	// Sets optional value
	if pd.Value != nil {
		slider.SetValue(*pd.Value)
	}
	return slider, nil
}

// buildSplitter builds a gui object of type: HSplitterr or VSplitter
func (b *Builder) buildSplitter(pd *descPanel, horiz bool) (IPanel, error) {

	// Builds splitter and sets its common attributes
	width, height := b.size(pd)
	var splitter *Splitter
	if horiz {
		splitter = NewHSplitter(width, height)
	} else {
		splitter = NewVSplitter(width, height)
	}
	err := b.setAttribs(pd, splitter, asWIDGET)
	if err != nil {
		return nil, err
	}

	// Optional split value
	if pd.Split != nil {
		splitter.SetSplit(*pd.Split)
	}

	// Splitter panel 0 attributes and items
	if pd.P0 != nil {
		err := b.setAttribs(pd.P0, &splitter.P0, asPANEL)
		if err != nil {
			return nil, err
		}
		err = b.addPanelItems(pd.P0, &splitter.P0)
		if err != nil {
			return nil, err
		}
	}

	// Splitter panel 1 attributes and items
	if pd.P1 != nil {
		err := b.setAttribs(pd.P1, &splitter.P1, asPANEL)
		if err != nil {
			return nil, err
		}
		err = b.addPanelItems(pd.P1, &splitter.P1)
		if err != nil {
			return nil, err
		}
	}

	return splitter, nil
}

// buildTree builds a gui object of type: Tree
func (b *Builder) buildTree(dp *descPanel) (IPanel, error) {

	// Builds tree and sets its common attributes
	width, height := b.size(dp)
	tree := NewTree(width, height)
	err := b.setAttribs(dp, tree, asWIDGET)
	if err != nil {
		return nil, err
	}

	// Internal function to build tree nodes recursively
	var buildItems func(dp *descPanel, pnode *TreeNode) error
	buildItems = func(dp *descPanel, pnode *TreeNode) error {
		for i := 0; i < len(dp.Items); i++ {
			item := dp.Items[i]
			// Item is a tree node
			if item.Type == "" || item.Type == descTypeTreeNode {
				var node *TreeNode
				if pnode == nil {
					node = tree.AddNode(item.Text)
				} else {
					node = pnode.AddNode(item.Text)
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
	err = buildItems(dp, nil)
	if err != nil {
		return nil, err
	}
	return tree, nil
}

// buildMenu builds a gui object of type: Menu or MenuBar from the
// specified panel descriptor.
func (b *Builder) buildMenu(pd *descPanel, child, bar bool) (IPanel, error) {

	// Builds menu bar or menu
	var menu *Menu
	if bar {
		menu = NewMenuBar()
	} else {
		menu = NewMenu()
	}
	// Only sets attribs for top level menus
	if !child {
		err := b.setAttribs(pd, menu, asWIDGET)
		if err != nil {
			return nil, err
		}
	}

	// Builds and adds menu items
	for i := 0; i < len(pd.Items); i++ {
		item := pd.Items[i]
		// Item is another menu
		if item.Type == descTypeMenu {
			subm, err := b.buildMenu(item, true, false)
			if err != nil {
				return nil, err
			}
			menu.AddMenu(item.Text, subm.(*Menu))
			continue
		}
		// Item is a separator
		if item.Type == "Separator" {
			menu.AddSeparator()
			continue
		}
		// Item must be a menu option
		mi := menu.AddOption(item.Text)
		// Set item optional icon(s)
		icons, err := b.parseIconNames("icon", item.Icon)
		if err != nil {
			return nil, err
		}
		if icons != "" {
			mi.SetIcon(string(icons))
		}
		// Sets optional menu item shortcut
		err = b.setMenuShortcut(mi, "shortcut", item.Shortcut)
		if err != nil {
			return nil, err
		}
	}
	return menu, nil
}

// buildWindow builds a gui object of type: Window from the
// specified panel descriptor.
func (b *Builder) buildWindow(dp *descPanel) (IPanel, error) {

	// Builds window and sets its common attributes
	width, height := b.size(dp)
	win := NewWindow(width, height)
	err := b.setAttribs(dp, win, asWIDGET)
	if err != nil {
		return nil, err
	}

	// Title attribute
	win.SetTitle(dp.Title)

	// Parse resizable borders
	if dp.Resizable != "" {
		parts := strings.Fields(dp.Resizable)
		var res Resizable
		for _, name := range parts {
			v, ok := mapResizable[name]
			if !ok {
				return nil, b.err("resizable", "Invalid resizable name:"+name)
			}
			res |= v
		}
		win.SetResizable(res)
	}

	// Builds window client panel children recursively
	for i := 0; i < len(dp.Items); i++ {
		item := dp.Items[i]
		b.objpath.push(item.Name)
		child, err := b.build(item, win)
		b.objpath.pop()
		if err != nil {
			return nil, err
		}
		win.Add(child)
	}
	return win, nil
}

// addPanelItems adds the items in the panel descriptor to the specified panel
func (b *Builder) addPanelItems(dp *descPanel, ipan IPanel) error {

	pan := ipan.GetPanel()
	for i := 0; i < len(dp.Items); i++ {
		item := dp.Items[i]
		b.objpath.push(item.Name)
		child, err := b.build(item, pan)
		b.objpath.pop()
		if err != nil {
			return err
		}
		pan.Add(child)
	}
	return nil
}

// setAttribs sets common attributes from the description to the specified panel
// The attributes which are set can be specified by the specified bitmask.
func (b *Builder) setAttribs(pd *descPanel, ipan IPanel, attr uint) error {

	panel := ipan.GetPanel()
	// Set optional position
	if attr&aPOS != 0 && pd.Position != "" {
		va, err := b.parseFloats("position", pd.Position, 2, 2)
		if va == nil || err != nil {
			return err
		}
		panel.SetPosition(va[0], va[1])
	}

	// Set optional size
	if attr&aSIZE != 0 {
		if pd.Width != nil {
			panel.SetWidth(*pd.Width)
		}
		if pd.Height != nil {
			panel.SetHeight(*pd.Height)
		}
	}

	// Set optional margin sizes
	if attr&aMARGINS != 0 {
		bs, err := b.parseBorderSizes(fieldMargins, pd.Margins)
		if err != nil {
			return err
		}
		if bs != nil {
			panel.SetMarginsFrom(bs)
		}
	}

	// Set optional border sizes
	if attr&aBORDERS != 0 {
		bs, err := b.parseBorderSizes(fieldBorders, pd.Borders)
		if err != nil {
			return err
		}
		if bs != nil {
			panel.SetBordersFrom(bs)
		}
	}

	// Set optional border color
	if attr&aBORDERCOLOR != 0 {
		c, err := b.parseColor(fieldBorderColor, pd.BorderColor)
		if err != nil {
			return err
		}
		if c != nil {
			panel.SetBordersColor4(c)
		}
	}

	// Set optional paddings sizes
	if attr&aPADDINGS != 0 {
		bs, err := b.parseBorderSizes(fieldPaddings, pd.Paddings)
		if err != nil {
			return err
		}
		if bs != nil {
			panel.SetPaddingsFrom(bs)
		}
	}

	// Set optional color
	if attr&aCOLOR != 0 {
		c, err := b.parseColor(fieldColor, pd.Color)
		if err != nil {
			return err
		}
		if c != nil {
			panel.SetColor4(c)
		}
	}

	if attr&aNAME != 0 && pd.Name != "" {
		panel.SetName(pd.Name)
	}
	if attr&aVISIBLE != 0 && pd.Visible != nil {
		panel.SetVisible(*pd.Visible)
	}
	if attr&aENABLED != 0 && pd.Enabled != nil {
		panel.SetEnabled(*pd.Enabled)
	}
	if attr&aRENDER != 0 && pd.Renderable != nil {
		panel.SetRenderable(*pd.Renderable)
	}

	err := b.setLayoutParams(pd, ipan)
	if err != nil {
		return err
	}
	return b.setLayout(pd, ipan)
}

// setLayoutParams sets the optional layout params attribute for specified the panel
func (b *Builder) setLayoutParams(dp *descPanel, ipan IPanel) error {

	// If layout params not declared, nothing to do
	if dp.LayoutParams == nil {
		return nil
	}

	// Get the parent layout
	if dp.parent == nil {
		return b.err("layoutparams", "No parent defined")
	}
	playout := dp.parent.Layout
	if playout == nil {
		return b.err("layoutparams", "Parent does not have layout")
	}
	panel := ipan.GetPanel()
	dlp := dp.LayoutParams

	// HBoxLayout parameters
	if playout.Type == descTypeHBoxLayout {
		// Creates layout parameter
		params := HBoxLayoutParams{Expand: 0, AlignV: AlignTop}
		// Sets optional expand parameter
		if dlp.Expand != nil {
			params.Expand = *dlp.Expand
		}
		// Sets optional align parameter
		if dlp.AlignV != "" {
			align, ok := mapAlignName[dlp.AlignV]
			if !ok {
				return b.err("align", "Invalid align name:"+dlp.AlignV)
			}
			params.AlignV = align
		}
		panel.SetLayoutParams(&params)
		return nil
	}

	// VBoxLayout parameters
	if playout.Type == descTypeVBoxLayout {
		// Creates layout parameter
		params := VBoxLayoutParams{Expand: 0, AlignH: AlignLeft}
		// Sets optional expand parameter
		if dlp.Expand != nil {
			params.Expand = *dlp.Expand
		}
		// Sets optional align parameter
		if dlp.AlignH != "" {
			align, ok := mapAlignName[dlp.AlignH]
			if !ok {
				return b.err("align", "Invalid align name:"+dlp.AlignH)
			}
			params.AlignH = align
		}
		panel.SetLayoutParams(&params)
		return nil
	}

	// GridLayout parameters
	if playout.Type == descTypeGridLayout {
		// Creates layout parameter
		params := GridLayoutParams{
			ColSpan: 0,
			AlignH:  AlignNone,
			AlignV:  AlignNone,
		}
		params.ColSpan = dlp.ColSpan
		// Sets optional alignh parameter
		if dlp.AlignH != "" {
			align, ok := mapAlignName[dlp.AlignH]
			if !ok {
				return b.err("alignh", "Invalid align name:"+dlp.AlignH)
			}
			params.AlignH = align
		}
		// Sets optional alignv parameter
		if dlp.AlignV != "" {
			align, ok := mapAlignName[dlp.AlignV]
			if !ok {
				return b.err("alignv", "Invalid align name:"+dlp.AlignV)
			}
			params.AlignV = align
		}
		panel.SetLayoutParams(&params)
		return nil
	}

	// DockLayout parameters
	if playout.Type == descTypeDockLayout {
		if dlp.Edge != "" {
			edge, ok := mapEdgeName[dlp.Edge]
			if !ok {
				return b.err("edge", "Invalid edge name:"+dlp.Edge)
			}
			params := DockLayoutParams{Edge: edge}
			panel.SetLayoutParams(&params)
			return nil
		}
	}

	return b.err("layoutparams", "Invalid parent layout:"+playout.Type)
}


// setLayout sets the optional panel layout and layout parameters
func (b *Builder) setLayout(dp *descPanel, ipan IPanel) error {

	// If layout types not declared, nothing to do
	if dp.Layout == nil {
		return nil
	}
	dl := dp.Layout

	// HBox layout
	if dl.Type == descTypeHBoxLayout {
		hbl := NewHBoxLayout()
		hbl.SetSpacing(dl.Spacing)
		if dl.AlignH != "" {
			align, ok := mapAlignName[dl.AlignH]
			if !ok {
				return b.err("align", "Invalid align name:"+dl.AlignV)
			}
			hbl.SetAlignH(align)
		}
		hbl.SetMinHeight(dl.MinHeight)
		hbl.SetMinWidth(dl.MinWidth)
		ipan.SetLayout(hbl)
		return nil
	}

	// VBox layout
	if dl.Type == descTypeVBoxLayout {
		vbl := NewVBoxLayout()
		vbl.SetSpacing(dl.Spacing)
		if dl.AlignV != "" {
			align, ok := mapAlignName[dl.AlignV]
			if !ok {
				return b.err("align", "Invalid align name:"+dl.AlignV)
			}
			vbl.SetAlignV(align)
		}
		vbl.SetMinHeight(dl.MinHeight)
		vbl.SetMinWidth(dl.MinWidth)
		ipan.SetLayout(vbl)
		return nil
	}

	// Grid layout
	if dl.Type == descTypeGridLayout {
		// Number of columns
		if dl.Cols == 0 {
			return b.err("cols", "Invalid number of columns:"+dl.AlignH)
		}
		grl := NewGridLayout(dl.Cols)
		// Global horizontal alignment
		if dl.AlignH != "" {
			alignh, ok := mapAlignName[dl.AlignH]
			if !ok {
				return b.err("alignh", "Invalid horizontal align:"+dl.AlignH)
			}
			grl.SetAlignH(alignh)
		}
		// Global vertical alignment
		if dl.AlignV != "" {
			alignv, ok := mapAlignName[dl.AlignV]
			if !ok {
				return b.err("alignv", "Invalid vertical align:"+dl.AlignH)
			}
			grl.SetAlignV(alignv)
		}
		// Expansion flags
		grl.SetExpandH(dl.ExpandH)
		grl.SetExpandV(dl.ExpandV)
		ipan.SetLayout(grl)
		return nil
	}

	// Dock layout
	if dl.Type == descTypeDockLayout {
		dockl := NewDockLayout()
		ipan.SetLayout(dockl)
		return nil
	}

	return b.err("layout", "Invalid layout type:"+dl.Type)
}
****/

// setAttribs sets common attributes from the description to the specified panel
// The attributes which are set can be specified by the specified bitmask.
func (b *Builder) setAttribs(am map[string]interface{}, ipan IPanel, attr uint) error {

	panel := ipan.GetPanel()
	// Set optional position
	if attr&aPOS != 0 && am[AttribPosition] != nil {
		va := am[AttribPosition].([]float32)
		panel.SetPosition(va[0], va[1])
	}

	// Set optional panel width
	if attr&aSIZE != 0 && am[AttribWidth] != nil {
		panel.SetWidth(am[AttribWidth].(float32))
		log.Error("set width:%v", am[AttribWidth])
	}

	// Sets optional panel height
	if attr&aSIZE != 0 && am[AttribHeight] != nil {
		panel.SetHeight(am[AttribHeight].(float32))
	}

	// Set optional margin sizes
	if attr&aMARGINS != 0 && am[AttribMargins] != nil {
		panel.SetMarginsFrom(am[AttribMargins].(*BorderSizes))
	}

	// Set optional border sizes
	if attr&aBORDERS != 0 && am[AttribBorders] != nil {
		panel.SetBordersFrom(am[AttribBorders].(*BorderSizes))
	}

	// Set optional border color
	if attr&aBORDERCOLOR != 0 && am[AttribBorderColor] != nil {
		panel.SetBordersColor4(am[AttribBorderColor].(*math32.Color4))
	}

	// Set optional paddings sizes
	if attr&aPADDINGS != 0 && am[AttribPaddings] != nil {
		panel.SetPaddingsFrom(am[AttribPaddings].(*BorderSizes))
	}

	// Set optional panel color
	if attr&aCOLOR != 0 && am[AttribColor] != nil {
		panel.SetColor4(am[AttribColor].(*math32.Color4))
	}

	if attr&aNAME != 0 && am[AttribName] != nil {
		panel.SetName(am[AttribName].(string))
	}

	if attr&aVISIBLE != 0 && am[AttribVisible] != nil {
		panel.SetVisible(am[AttribVisible].(bool))
	}

	if attr&aENABLED != 0 && am[AttribEnabled] != nil {
		panel.SetEnabled(am[AttribEnabled].(bool))
	}
	if attr&aRENDER != 0 && am[AttribRender] != nil {
		panel.SetRenderable(am[AttribRender].(bool))
	}

	// Sets optional layout
	err := b.setLayout(am, panel)
	if err != nil {
		return nil
	}

	// Sets optional layout params
	err = b.setLayoutParams(am, panel)
	return err
}

//func (b *Builder) setMenuShortcut(mi *MenuItem, fname, field string) error {
//
//	field = strings.Trim(field, " ")
//	if field == "" {
//		return nil
//	}
//	parts := strings.Split(field, "+")
//	var mods window.ModifierKey
//	for i := 0; i < len(parts)-1; i++ {
//		switch parts[i] {
//		case "Shift":
//			mods |= window.ModShift
//		case "Ctrl":
//			mods |= window.ModControl
//		case "Alt":
//			mods |= window.ModAlt
//		default:
//			return b.err(am, fname, "Invalid shortcut:"+field)
//		}
//	}
//	// The last part must be a key
//	key := parts[len(parts)-1]
//	for kcode, kname := range mapKeyText {
//		if kname == key {
//			mi.SetShortcut(mods, kcode)
//			return nil
//		}
//	}
//	return b.err(fname, "Invalid shortcut:"+field)
//}

// parseFloats parses a string with a list of floats with the specified size
// and returns a slice. The specified size is 0 any number of floats is allowed.
// The individual values can be separated by spaces or commas
func (b *Builder) parseFloats(am map[string]interface{}, fname string, min, max int) ([]float32, error) {

	// Checks if field is empty
	v := am[fname]
	if v == nil {
		return nil, nil
	}

	// If field has only one value, it is an int or a float64
	switch ft := v.(type) {
	case int:
		return []float32{float32(ft)}, nil
	case float64:
		return []float32{float32(ft)}, nil
	}

	// Converts to string
	fs, ok := v.(string)
	if !ok {
		return nil, b.err(am, fname, "Not a string")
	}

	// Checks if string field is empty
	fs = strings.Trim(fs, " ")
	if fs == "" {
		return nil, nil
	}

	// Separate individual fields
	var parts []string
	if strings.Index(fs, ",") < 0 {
		parts = strings.Fields(fs)
	} else {
		parts = strings.Split(fs, ",")
	}
	if len(parts) < min || len(parts) > max {
		return nil, b.err(am, fname, "Invalid number of float32 values")
	}

	// Parse each field value and appends to slice
	var values []float32
	for i := 0; i < len(parts); i++ {
		val, err := strconv.ParseFloat(strings.Trim(parts[i], " "), 32)
		if err != nil {
			return nil, b.err(am, fname, err.Error())
		}
		values = append(values, float32(val))
	}
	return values, nil
}

// err creates and returns an error for the current object, field name and with the specified message
func (b *Builder) err(am map[string]interface{}, fname, msg string) error {

	return fmt.Errorf("Error in object:%s field:%s -> %s", am[AttribName], fname, msg)
}

// debugPrint prints the internal attribute map of the builder for debugging.
// This map cannot be printed by fmt.Printf() because it has cycles.
// A map contains a key: _parent, which pointer to is parent map, if any.
func (b *Builder) debugPrint(v interface{}, level int) {

	switch vt := v.(type) {
	case map[string]interface{}:
		level += 3
		fmt.Printf("\n")
		for mk, mv := range vt {
			if mk == AttribParent_ {
				continue
			}
			fmt.Printf("%s%s:", strings.Repeat(" ", level), mk)
			b.debugPrint(mv, level)
		}
	case []map[string]interface{}:
		for _, v := range vt {
			b.debugPrint(v, level)
		}
	default:
		fmt.Printf(" %v (%T)\n", vt, vt)
	}
}
