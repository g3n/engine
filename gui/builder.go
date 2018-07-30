// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"fmt"
	"io/ioutil"
	"os"
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
	am       map[string]interface{}     // parsed attribute map
	builders map[string]BuilderFunc     // map of builder functions by type
	attribs  map[string]AttribCheckFunc // map of attribute name with check functions
	layouts  map[string]IBuilderLayout  // map of layout type to layout builder
	imgpath  string                     // base path for image panels files
}

// IBuilderLayout is the interface for all layout builders
type IBuilderLayout interface {
	BuildLayout(b *Builder, am map[string]interface{}) (ILayout, error)
	BuildParams(b *Builder, am map[string]interface{}) (interface{}, error)
}

// BuilderFunc is type for functions which build a gui object from an attribute map
type BuilderFunc func(*Builder, map[string]interface{}) (IPanel, error)

// AttribCheckFunc is the type for all attribute check functions
type AttribCheckFunc func(b *Builder, am map[string]interface{}, fname string) error

// IgnoreSuffix specifies the suffix of ignored keys
const IgnoreSuffix = "_"

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
	TypeChart       = "chart"
	TypeTable       = "table"
	TypeTabBar      = "tabbar"
	TypeHBoxLayout  = "hbox"
	TypeVBoxLayout  = "vbox"
	TypeGridLayout  = "grid"
	TypeDockLayout  = "dock"
)

// Common attribute names
const (
	AttribAlignv         = "alignv"        // Align
	AttribAlignh         = "alignh"        // Align
	AttribAspectHeight   = "aspectheight"  // float32
	AttribAspectWidth    = "aspectwidth"   // float32
	AttribBgColor        = "bgcolor"       // Color4
	AttribBorders        = "borders"       // RectBounds
	AttribBorderColor    = "bordercolor"   // Color4
	AttribChecked        = "checked"       // bool
	AttribColor          = "color"         // Color4
	AttribCols           = "cols"          // int GridLayout
	AttribColSpan        = "colspan"       // int GridLayout
	AttribColumns        = "columns"       // []map[string]interface{} Table
	AttribContent        = "content"       // map[string]interface{} Table
	AttribCountStepx     = "countstepx"    // float32
	AttribEdge           = "edge"          // int
	AttribEnabled        = "enabled"       // bool
	AttribExpand         = "expand"        // float32
	AttribExpandh        = "expandh"       // bool
	AttribExpandv        = "expandv"       // bool
	AttribFirstx         = "firstx"        // float32
	AttribFontColor      = "fontcolor"     // Color4
	AttribFontDPI        = "fontdpi"       // float32
	AttribFontSize       = "fontsize"      // float32
	AttribFormat         = "format"        // string
	AttribGroup          = "group"         // string
	AttribHeader         = "header"        // string
	AttribHeight         = "height"        // float32
	AttribHidden         = "hidden"        // bool Table
	AttribId             = "id"            // string
	AttribIcon           = "icon"          // string
	AttribImageFile      = "imagefile"     // string
	AttribImageLabel     = "imagelabel"    // []map[string]interface{}
	AttribItems          = "items"         // []map[string]interface{}
	AttribLayout         = "layout"        // map[string]interface{}
	AttribLayoutParams   = "layoutparams"  // map[string]interface{}
	AttribLineSpacing    = "linespacing"   // float32
	AttribLines          = "lines"         // int
	AttribMargin         = "margin"        // float32
	AttribMargins        = "margins"       // RectBounds
	AttribMinwidth       = "minwidth"      // float32 Table
	AttribAutoHeight     = "autoheight"    // bool
	AttribAutoWidth      = "autowidth"     // bool
	AttribName           = "name"          // string
	AttribPaddings       = "paddings"      // RectBounds
	AttribPanel0         = "panel0"        // map[string]interface{}
	AttribPanel1         = "panel1"        // map[string]interface{}
	AttribParentInternal = "parent_"       // string (internal attribute)
	AttribPinned         = "pinned"        // bool
	AttribPlaceHolder    = "placeholder"   // string
	AttribPosition       = "position"      // []float32
	AttribRangeAuto      = "rangeauto"     // bool
	AttribRangeMin       = "rangemin"      // float32
	AttribRangeMax       = "rangemax"      // float32
	AttribRender         = "render"        // bool
	AttribResizeBorders  = "resizeborders" // Resizable
	AttribResize         = "resize"        // bool Table
	AttribScaleFactor    = "scalefactor"   // float32
	AttribScalex         = "scalex"        // map[string]interface{}
	AttribScaley         = "scaley"        // map[string]interface{}
	AttribShortcut       = "shortcut"      // []int
	AttribShowHeader     = "showheader"    // bool
	AttribSortType       = "sorttype"      // TableSortType Table
	AttribSpacing        = "spacing"       // float32
	AttribSplit          = "split"         // float32
	AttribStepx          = "stepx"         // float32
	AttribText           = "text"          // string
	AttribTitle          = "title"         // string
	AttribType           = "type"          // string
	AttribUserData       = "userdata"      // interface{}
	AttribWidth          = "width"         // float32
	AttribValue          = "value"         // float32
	AttribVisible        = "visible"       // bool
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

// maps table sort type to value
var mapTableSortType = map[string]TableSortType{
	"none":   TableSortNone,
	"string": TableSortString,
	"number": TableSortNumber,
}

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
		TypeWindow:      buildWindow,
		TypeChart:       buildChart,
		TypeTable:       buildTable,
		TypeTabBar:      buildTabBar,
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
		AttribAlignv:        AttribCheckAlign,
		AttribAlignh:        AttribCheckAlign,
		AttribAspectWidth:   AttribCheckFloat,
		AttribAspectHeight:  AttribCheckFloat,
		AttribHeight:        AttribCheckFloat,
		AttribBgColor:       AttribCheckColor,
		AttribBorders:       AttribCheckBorderSizes,
		AttribBorderColor:   AttribCheckColor,
		AttribChecked:       AttribCheckBool,
		AttribColor:         AttribCheckColor,
		AttribCols:          AttribCheckInt,
		AttribColSpan:       AttribCheckInt,
		AttribColumns:       AttribCheckListMap,
		AttribContent:       AttribCheckMap,
		AttribCountStepx:    AttribCheckFloat,
		AttribEdge:          AttribCheckEdge,
		AttribEnabled:       AttribCheckBool,
		AttribExpand:        AttribCheckFloat,
		AttribExpandh:       AttribCheckBool,
		AttribExpandv:       AttribCheckBool,
		AttribFirstx:        AttribCheckFloat,
		AttribFontColor:     AttribCheckColor,
		AttribFontDPI:       AttribCheckFloat,
		AttribFontSize:      AttribCheckFloat,
		AttribFormat:        AttribCheckString,
		AttribGroup:         AttribCheckString,
		AttribHeader:        AttribCheckString,
		AttribHidden:        AttribCheckBool,
		AttribIcon:          AttribCheckIcons,
		AttribId:            AttribCheckString,
		AttribImageFile:     AttribCheckString,
		AttribImageLabel:    AttribCheckMap,
		AttribItems:         AttribCheckListMap,
		AttribLayout:        AttribCheckLayout,
		AttribLayoutParams:  AttribCheckMap,
		AttribLineSpacing:   AttribCheckFloat,
		AttribLines:         AttribCheckInt,
		AttribMargin:        AttribCheckFloat,
		AttribMargins:       AttribCheckBorderSizes,
		AttribMinwidth:      AttribCheckFloat,
		AttribAutoHeight:    AttribCheckBool,
		AttribAutoWidth:     AttribCheckBool,
		AttribName:          AttribCheckString,
		AttribPaddings:      AttribCheckBorderSizes,
		AttribPanel0:        AttribCheckMap,
		AttribPanel1:        AttribCheckMap,
		AttribPinned:        AttribCheckBool,
		AttribPlaceHolder:   AttribCheckString,
		AttribPosition:      AttribCheckPosition,
		AttribRangeAuto:     AttribCheckBool,
		AttribRangeMin:      AttribCheckFloat,
		AttribRangeMax:      AttribCheckFloat,
		AttribRender:        AttribCheckBool,
		AttribResizeBorders: AttribCheckResizeBorders,
		AttribResize:        AttribCheckBool,
		AttribScaleFactor:   AttribCheckFloat,
		AttribScalex:        AttribCheckMap,
		AttribScaley:        AttribCheckMap,
		AttribShortcut:      AttribCheckMenuShortcut,
		AttribShowHeader:    AttribCheckBool,
		AttribSortType:      AttribCheckTableSortType,
		AttribSpacing:       AttribCheckFloat,
		AttribSplit:         AttribCheckFloat,
		AttribStepx:         AttribCheckFloat,
		AttribText:          AttribCheckString,
		AttribTitle:         AttribCheckString,
		AttribType:          AttribCheckStringLower,
		AttribUserData:      AttribCheckInterface,
		AttribValue:         AttribCheckFloat,
		AttribVisible:       AttribCheckBool,
		AttribWidth:         AttribCheckFloat,
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
			if par != nil {
				ms[AttribParentInternal] = par
			}
			for k, v := range vt {
				// Checks key
				ks, ok := k.(string)
				if !ok {
					return nil, fmt.Errorf("Keys must be strings")
				}
				ks = strings.ToLower(ks)
				// Ignores keys suffixed by IgnoreSuffix
				if strings.HasSuffix(ks, IgnoreSuffix) {
					continue
				}
				// Checks value
				vi, err := visitor(v, ms)
				if err != nil {
					return nil, err
				}
				ms[ks] = vi
				// If has parent or is a single top level panel, checks attributes
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
			return ms, nil
		}
		return v, nil
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
	// Single object
	if b.am[AttribType] != nil {
		objs = append(objs, "")
		return objs
	}
	// Multiple objects
	for name := range b.am {
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
		return b.build(b.am, nil)
	}
	// Map of gui objects
	am, ok := b.am[name]
	if !ok {
		return nil, fmt.Errorf("Object name:%s not found", name)
	}
	return b.build(am.(map[string]interface{}), nil)
}

// SetImagepath Sets the path for image panels relative image files
func (b *Builder) SetImagepath(path string) {

	b.imgpath = path
}

// AddBuilderPanel adds a panel builder function for the specified type name.
// If the type name already exists it is replaced.
func (b *Builder) AddBuilderPanel(typename string, bf BuilderFunc) {

	b.builders[typename] = bf
}

// AddBuilderLayout adds a layout builder object for the specified type name.
// If the type name already exists it is replaced.
func (b *Builder) AddBuilderLayout(typename string, bl IBuilderLayout) {

	b.layouts[typename] = bl
}

// AddAttrib adds an attribute type and its checker/converte
// If the attribute type name already exists it is replaced.
func (b *Builder) AddAttrib(typename string, acf AttribCheckFunc) {

	b.attribs[typename] = acf
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

// SetAttribs sets common attributes from the description to the specified panel
func (b *Builder) SetAttribs(am map[string]interface{}, ipan IPanel) error {

	panel := ipan.GetPanel()
	// Set optional position
	if am[AttribPosition] != nil {
		va := am[AttribPosition].([]float32)
		panel.SetPosition(va[0], va[1])
	}

	// Set optional panel width
	if am[AttribWidth] != nil {
		panel.SetWidth(am[AttribWidth].(float32))
	}

	// Sets optional panel height
	if am[AttribHeight] != nil {
		panel.SetHeight(am[AttribHeight].(float32))
	}

	// Set optional margin sizes
	if am[AttribMargins] != nil {
		panel.SetMarginsFrom(am[AttribMargins].(*RectBounds))
	}

	// Set optional border sizes
	if am[AttribBorders] != nil {
		panel.SetBordersFrom(am[AttribBorders].(*RectBounds))
	}

	// Set optional border color
	if am[AttribBorderColor] != nil {
		panel.SetBordersColor4(am[AttribBorderColor].(*math32.Color4))
	}

	// Set optional paddings sizes
	if am[AttribPaddings] != nil {
		panel.SetPaddingsFrom(am[AttribPaddings].(*RectBounds))
	}

	// Set optional panel color
	if am[AttribColor] != nil {
		panel.SetColor4(am[AttribColor].(*math32.Color4))
	}

	if am[AttribName] != nil {
		panel.SetName(am[AttribName].(string))
	}

	if am[AttribVisible] != nil {
		panel.SetVisible(am[AttribVisible].(bool))
	}

	if am[AttribEnabled] != nil {
		panel.SetEnabled(am[AttribEnabled].(bool))
	}

	if am[AttribRender] != nil {
		panel.SetRenderable(am[AttribRender].(bool))
	}

	if am[AttribUserData] != nil {
		panel.SetUserData(am[AttribUserData])
	}

	// Sets optional layout (must pass IPanel not *Panel)
	err := b.setLayout(am, ipan)
	if err != nil {
		return nil
	}

	// Sets optional layout params
	err = b.setLayoutParams(am, panel)
	return err
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

// setLayoutParams sets the optional layout params of the specified panel and its attributes
func (b *Builder) setLayoutParams(am map[string]interface{}, ipan IPanel) error {

	// Get layout params attributes
	lpi := am[AttribLayoutParams]
	if lpi == nil {
		return nil
	}
	lp := lpi.(map[string]interface{})

	// Checks if layout param specifies the layout type
	// This is useful when the panel has no parent yet
	var ltype string
	if v := lp[AttribType]; v != nil {
		ltype = v.(string)
	} else {
		// Get layout type from parent
		pi := am[AttribParentInternal]
		if pi == nil {
			return b.err(am, AttribType, "Panel has no parent")
		}
		par := pi.(map[string]interface{})
		v := par[AttribLayout]
		if v == nil {
			return b.err(am, AttribType, "Parent has no layout")
		}
		playout := v.(map[string]interface{})
		ltype = playout[AttribType].(string)
	}

	// Get layout builder and builds layout params
	lbuilder := b.layouts[ltype]
	params, err := lbuilder.BuildParams(b, lp)
	if err != nil {
		return err
	}
	ipan.GetPanel().SetLayoutParams(params)
	return nil
}

// AttribCheckTableSortType checks and converts attribute table column sort type
func AttribCheckTableSortType(b *Builder, am map[string]interface{}, fname string) error {

	// If attribute not found, ignore
	v := am[fname]
	if v == nil {
		return nil
	}
	vs, ok := v.(string)
	if !ok {
		return b.err(am, fname, "Invalid attribute")
	}
	tstype, ok := mapTableSortType[vs]
	if !ok {
		return b.err(am, fname, "Invalid attribute")
	}
	am[fname] = tstype
	return nil
}

// AttribCheckResizeBorders checks and converts attribute with list of window resizable borders
func AttribCheckResizeBorders(b *Builder, am map[string]interface{}, fname string) error {

	// If attribute not found, ignore
	v := am[fname]
	if v == nil {
		return nil
	}

	// Attribute must be string
	vs, ok := v.(bool)
	if !ok {
		return b.err(am, fname, "Invalid resizable attribute")
	}

	am[fname] = vs
	return nil
}

// AttribCheckEdge checks and converts attribute with name of layout edge
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

// AttribCheckLayout checks and converts layout attribute
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

// AttribCheckAlign checks and converts layout align* attribute
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

// AttribCheckMenuShortcut checks and converts attribute describing menu shortcut key
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

// AttribCheckListMap checks and converts attribute to []map[string]interface{}
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

// AttribCheckMap checks and converts attribute to map[string]interface{}
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

// AttribCheckIcons checks and converts attribute with a list of icon names or codepoints
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

// AttribCheckColor checks and converts attribute with color name or color component values
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

// AttribCheckBorderSizes checks and convert attribute with border sizes
func AttribCheckBorderSizes(b *Builder, am map[string]interface{}, fname string) error {

	va, err := b.parseFloats(am, fname, 1, 4)
	if err != nil {
		return err
	}
	if va == nil {
		return nil
	}
	if len(va) == 1 {
		am[fname] = &RectBounds{va[0], va[0], va[0], va[0]}
		return nil
	}
	am[fname] = &RectBounds{va[0], va[1], va[2], va[3]}
	return nil
}

// AttribCheckPosition checks and convert attribute with x and y position
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

// AttribCheckStringLower checks and convert string attribute to lower case
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

// AttribCheckFloat checks and convert attribute to float32
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
	}
	return b.err(am, fname, fmt.Sprintf("Not a number:%T", v))
}

// AttribCheckInt checks and convert attribute to int
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

// AttribCheckString checks and convert attribute to string
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

// AttribCheckBool checks and convert attribute to bool
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

// AttribCheckInterface accepts any attribute value
func AttribCheckInterface(b *Builder, am map[string]interface{}, fname string) error {

	return nil
}

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

	// Get path of objects till the error
	names := []string{}
	var name string
	for {
		if v := am[AttribName]; v != nil {
			name = v.(string)
		} else {
			name = "?"
		}
		names = append(names, name)
		var par interface{}
		if par = am[AttribParentInternal]; par == nil {
			break
		}
		am = par.(map[string]interface{})
	}
	path := []string{}
	for i := len(names) - 1; i >= 0; i-- {
		path = append(path, names[i])
	}

	return fmt.Errorf("Error in object:%s field:%s -> %s", strings.Join(path, "/"), fname, msg)
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
			if mk == AttribParentInternal {
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
