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
	"gopkg.in/yaml.v2"
)

// Builder builds GUI objects from a declarative description in YAML format
type Builder struct {
	desc    map[string]*panelDesc // parsed descriptions
	imgpath string                // base path for image panels files
	objpath strStack              // stack of object names being built
}

type panelStyle struct {
	Borders     string
	Paddings    string
	BorderColor string
	BgColor     string
	FgColor     string
}

type panelStyles struct {
	Normal   *panelStyle
	Over     *panelStyle
	Focus    *panelStyle
	Pressed  *panelStyle
	Disabled *panelStyle
}

type panelDesc struct {
	Type         string   // Gui object type: Panel, Label, Edit, etc ...
	Name         string   // Optional name for identification
	Position     string   // Optional position as: x y | x,y
	Width        float32  // Optional width (default = 0)
	Height       float32  // Optional height (default = 0)
	AspectWidth  *float32 // Optional aspectwidth (default = nil)
	AspectHeight *float32 // Optional aspectwidth (default = nil)
	Margins      string   // Optional margins as 1 or 4 float values
	Borders      string   // Optional borders as 1 or 4 float values
	BorderColor  string   // Optional border color as name or 3 or 4 float values
	Paddings     string   // Optional paddings as 1 or 4 float values
	Color        string   // Optional color as 1 or 4 float values
	Enabled      bool
	Visible      bool
	Renderable   bool
	Imagefile    string // For Panel, Button
	Children     []*panelDesc
	Layout       layoutAttr
	Styles       *panelStyles
	Text         string   // Label, Button
	Icons        string   // Label
	BgColor      string   // Label
	FontColor    string   // Label
	FontSize     *float32 // Label
	FontDPI      *float32 // Label
	LineSpacing  *float32 // Label
	PlaceHolder  string   // Edit
	MaxLength    *uint    // Edit
	Icon         string   // Button
	Group        string   // RadioButton
}

type layoutAttr struct {
	Type string
}

const (
	descTypePanel       = "Panel"
	descTypeImagePanel  = "ImagePanel"
	descTypeLabel       = "Label"
	descTypeImageLabel  = "ImageLabel"
	descTypeButton      = "Button"
	descTypeCheckBox    = "CheckBox"
	descTypeRadioButton = "RadioButton"
	descTypeEdit        = "Edit"
	descTypeVList       = "VList"
	descTypeHList       = "HList"
	fieldMargins        = "margins"
	fieldBorders        = "borders"
	fieldBorderColor    = "bordercolor"
	fieldPaddings       = "paddings"
	fieldColor          = "color"
	fieldBgColor        = "bgcolor"
)

// NewBuilder creates and returns a pointer to a new gui Builder object
func NewBuilder() *Builder {

	return new(Builder)
}

// ParseString parses a string with gui objects descriptions in YAML format
// It there was a previously parsed description, it is cleared.
func (b *Builder) ParseString(desc string) error {

	// Try assuming the description contains a single root panel
	var pd panelDesc
	err := yaml.Unmarshal([]byte(desc), &pd)
	if err != nil {
		return err
	}
	if pd.Type != "" {
		b.desc = make(map[string]*panelDesc)
		b.desc[""] = &pd
		return nil
	}

	// Try assuming the description is a map of panels
	var pdm map[string]*panelDesc
	err = yaml.Unmarshal([]byte(desc), &pdm)
	if err != nil {
		return err
	}
	b.desc = pdm
	return nil
}

// ParseFile builds gui objects from the specified file which
// must contain objects descriptions in YAML format
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
	for name, pd := range b.desc {
		if pd.Type != "" {
			objs = append(objs, name)
		}
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

	pd, ok := b.desc[name]
	if !ok {
		return nil, fmt.Errorf("Object name:%s not found", name)
	}
	b.objpath.clear()
	b.objpath.push(pd.Name)
	return b.build(pd, nil)
}

// Sets the path for image panels relative image files
func (b *Builder) SetImagepath(path string) {

	b.imgpath = path
}

// build builds the gui object from the specified description.
// All its children are also built recursively
// Returns the built object or an error
func (b *Builder) build(pd *panelDesc, iparent IPanel) (IPanel, error) {

	var err error
	var pan IPanel
	switch pd.Type {
	case descTypePanel:
		pan, err = b.buildPanel(pd)
	case descTypeImagePanel:
		pan, err = b.buildImagePanel(pd)
	case descTypeLabel:
		pan, err = b.buildLabel(pd)
	case descTypeImageLabel:
		pan, err = b.buildImageLabel(pd)
	case descTypeButton:
		pan, err = b.buildButton(pd)
	case descTypeCheckBox:
		pan, err = b.buildCheckBox(pd)
	case descTypeRadioButton:
		pan, err = b.buildRadioButton(pd)
	case descTypeEdit:
		pan, err = b.buildEdit(pd)
	case descTypeVList:
		pan, err = b.buildVList(pd)
	case descTypeHList:
		pan, err = b.buildHList(pd)
	default:
		err = fmt.Errorf("Invalid panel type:%s", pd.Type)
	}
	if err != nil {
		return nil, err
	}
	if iparent != nil {
		iparent.GetPanel().Add(pan)
	}
	return pan, nil
}

// buildPanel builds a gui object of type: "Panel"
func (b *Builder) buildPanel(pd *panelDesc) (IPanel, error) {

	// Builds panel and set common attributes
	pan := NewPanel(pd.Width, pd.Height)
	err := b.setCommon(pd, pan)
	if err != nil {
		return nil, err
	}

	// Builds panel children recursively
	for i := 0; i < len(pd.Children); i++ {
		b.objpath.push(pd.Children[i].Name)
		child, err := b.build(pd.Children[i], pan)
		b.objpath.pop()
		if err != nil {
			return nil, err
		}
		pan.Add(child)
	}
	return pan, nil
}

// buildImagePanel builds a gui object of type: "ImagePanel"
func (b *Builder) buildImagePanel(pd *panelDesc) (IPanel, error) {

	// Imagefile must be supplied
	if pd.Imagefile == "" {
		return nil, b.err("Imagefile", "Imagefile must be supplied")
	}

	// If path is not absolute join with user supplied image base path
	path := pd.Imagefile
	if !filepath.IsAbs(path) {
		path = filepath.Join(b.imgpath, path)
	}

	// Builds panel and set common attributes
	panel, err := NewImage(path)
	if err != nil {
		return nil, err
	}
	err = b.setCommon(pd, panel)
	if err != nil {
		return nil, err
	}

	// AspectWidth and AspectHeight attributes
	if pd.AspectWidth != nil {
		panel.SetContentAspectWidth(*pd.AspectWidth)
	}
	if pd.AspectHeight != nil {
		panel.SetContentAspectHeight(*pd.AspectHeight)
	}

	// Builds panel children recursively
	for i := 0; i < len(pd.Children); i++ {
		b.objpath.push(pd.Children[i].Name)
		child, err := b.build(pd.Children[i], panel)
		b.objpath.pop()
		if err != nil {
			return nil, err
		}
		panel.Add(child)
	}
	return panel, nil
}

// buildLabel builds a gui object of type: "Label"
func (b *Builder) buildLabel(pd *panelDesc) (IPanel, error) {

	// Builds label with icon or text font
	var label *Label
	icons, err := b.parseIconNames("icons", pd.Icons)
	if err != nil {
		return nil, err
	}
	if icons != "" {
		label = NewLabel(icons, true)
	} else {
		label = NewLabel(pd.Text)
	}
	// Sets common attributes
	err = b.setCommon(pd, label)
	if err != nil {
		return nil, err
	}

	// Set optional background color
	c, err := b.parseColor(fieldBgColor, pd.BgColor)
	if err != nil {
		return nil, err
	}
	if c != nil {
		label.SetBgColor4(c)
	}

	// Set optional font color
	c, err = b.parseColor("fontcolor", pd.FontColor)
	if err != nil {
		return nil, err
	}
	if c != nil {
		label.SetColor4(c)
	}

	// Sets optional font size
	if pd.FontSize != nil {
		label.SetFontSize(float64(*pd.FontSize))
	}

	// Sets optional font dpi
	if pd.FontDPI != nil {
		label.SetFontDPI(float64(*pd.FontDPI))
	}

	// Sets optional line spacing
	if pd.LineSpacing != nil {
		label.SetLineSpacing(float64(*pd.LineSpacing))
	}

	return label, nil
}

// buildImageLabel builds a gui object of type: ImageLabel
func (b *Builder) buildImageLabel(pd *panelDesc) (IPanel, error) {

	// Builds image label and set common attributes
	imglabel := NewImageLabel(pd.Text)
	err := b.setCommon(pd, imglabel)
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

	return imglabel, nil
}

// buildButton builds a gui object of type: Button
func (b *Builder) buildButton(pd *panelDesc) (IPanel, error) {

	// Builds button and set commont attributes
	button := NewButton(pd.Text)
	err := b.setCommon(pd, button)
	if err != nil {
		return nil, err
	}

	// Sets optional icon
	if pd.Icon != "" {
		cp, err := b.parseIconName("icon", pd.Icon)
		if err != nil {
			return nil, err
		}
		button.SetIcon(int(cp))
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
func (b *Builder) buildCheckBox(pd *panelDesc) (IPanel, error) {

	// Builds check box and set commont attributes
	cb := NewCheckBox(pd.Text)
	err := b.setCommon(pd, cb)
	if err != nil {
		return nil, err
	}

	return cb, nil
}

// buildRadioButton builds a gui object of type: RadioButton
func (b *Builder) buildRadioButton(pd *panelDesc) (IPanel, error) {

	// Builds check box and set commont attributes
	rb := NewRadioButton(pd.Text)
	err := b.setCommon(pd, rb)
	if err != nil {
		return nil, err
	}

	// Sets optional radio button group
	if pd.Group != "" {
		rb.SetGroup(pd.Group)
	}
	return rb, nil
}

// buildEdit builds a gui object of type: "Edit"
func (b *Builder) buildEdit(pd *panelDesc) (IPanel, error) {

	// Builds button and set commont attributes
	edit := NewEdit(int(pd.Width), pd.PlaceHolder)
	err := b.setCommon(pd, edit)
	if err != nil {
		return nil, err
	}
	edit.SetText(pd.Text)
	return edit, nil
}

// buildVList builds a gui object of type: VList
func (b *Builder) buildVList(pd *panelDesc) (IPanel, error) {

	// Builds list and set commont attributes
	list := NewVList(pd.Width, pd.Height)
	err := b.setCommon(pd, list)
	if err != nil {
		return nil, err
	}

	// Builds list children
	for i := 0; i < len(pd.Children); i++ {
		b.objpath.push(pd.Children[i].Name)
		child, err := b.build(pd.Children[i], list)
		b.objpath.pop()
		if err != nil {
			return nil, err
		}
		list.Add(child)
	}
	return list, nil
}

// buildHList builds a gui object of type: VList
func (b *Builder) buildHList(pd *panelDesc) (IPanel, error) {

	// Builds list and set commont attributes
	list := NewHList(pd.Width, pd.Height)
	err := b.setCommon(pd, list)
	if err != nil {
		return nil, err
	}

	// Builds list children
	for i := 0; i < len(pd.Children); i++ {
		b.objpath.push(pd.Children[i].Name)
		child, err := b.build(pd.Children[i], list)
		b.objpath.pop()
		if err != nil {
			return nil, err
		}
		list.Add(child)
	}
	return list, nil
}

// setCommon sets the common attributes in the description to the specified panel
func (b *Builder) setCommon(pd *panelDesc, ipan IPanel) error {

	// Set optional position
	panel := ipan.GetPanel()
	if pd.Position != "" {
		va, err := b.parseFloats("position", pd.Position, 2, 2)
		if va == nil || err != nil {
			return err
		}
		panel.SetPosition(va[0], va[1])
	}

	// Set optional margin sizes
	bs, err := b.parseBorderSizes(fieldMargins, pd.Margins)
	if err != nil {
		return err
	}
	if bs != nil {
		panel.SetMarginsFrom(bs)
	}

	// Set optional border sizes
	bs, err = b.parseBorderSizes(fieldBorders, pd.Borders)
	if err != nil {
		return err
	}
	if bs != nil {
		panel.SetBordersFrom(bs)
	}

	// Set optional border color
	c, err := b.parseColor(fieldBorderColor, pd.BorderColor)
	if err != nil {
		return err
	}
	if c != nil {
		panel.SetBordersColor4(c)
	}

	// Set optional paddings sizes
	bs, err = b.parseBorderSizes(fieldPaddings, pd.Paddings)
	if err != nil {
		return err
	}
	if bs != nil {
		panel.SetPaddingsFrom(bs)
	}

	// Set optional color
	c, err = b.parseColor(fieldColor, pd.Color)
	if err != nil {
		return err
	}
	if c != nil {
		panel.SetColor4(c)
	}
	return nil
}

// parseBorderSizes parses a string field which can contain one float value or
// float values. In the first case all borders has the same width
func (b *Builder) parseBorderSizes(fname, field string) (*BorderSizes, error) {

	va, err := b.parseFloats(fname, field, 1, 4)
	if va == nil || err != nil {
		return nil, err
	}
	if len(va) == 1 {
		return &BorderSizes{va[0], va[0], va[0], va[0]}, nil
	}
	return &BorderSizes{va[0], va[1], va[2], va[3]}, nil
}

// parseColor parses a string field which can contain a color name or
// a list of 3 or 4 float values for the color components
func (b *Builder) parseColor(fname, field string) (*math32.Color4, error) {

	// Checks if field is empty
	field = strings.Trim(field, " ")
	if field == "" {
		return nil, nil
	}

	// If string has 1 or 2 fields it must be a color name and optional alpha
	parts := strings.Fields(field)
	if len(parts) == 1 || len(parts) == 2 {
		// First part must be a color name
		if !math32.IsColor(parts[0]) {
			return nil, b.err(fname, fmt.Sprintf("Invalid color name:%s", parts[0]))
		}
		c := math32.ColorName(parts[0])
		c4 := math32.Color4{c.R, c.G, c.B, 1}
		if len(parts) == 2 {
			val, err := strconv.ParseFloat(parts[1], 32)
			if err != nil {
				return nil, b.err(fname, fmt.Sprintf("Invalid float32 value:%s", parts[1]))
			}
			c4.A = float32(val)
		}
		return &c4, nil
	}

	// Accept 3 or 4 floats values
	va, err := b.parseFloats(fname, field, 3, 4)
	if err != nil {
		return nil, err
	}
	if len(va) == 3 {
		return &math32.Color4{va[0], va[1], va[2], 1}, nil
	}
	return &math32.Color4{va[0], va[1], va[2], va[3]}, nil
}

// parseIconNames parses a string with a list of icon names or codepoints and
// returns a string with the icons codepoints encoded in UTF8
func (b *Builder) parseIconNames(fname, field string) (string, error) {

	text := ""
	parts := strings.Fields(field)
	for i := 0; i < len(parts); i++ {
		cp, err := b.parseIconName(fname, parts[i])
		if err != nil {
			return "", err
		}
		text = text + string(cp)
	}
	return text, nil
}

// parseIconName parses a string with an icon name or codepoint in hex
// and returns the icon codepoints value and an error
func (b *Builder) parseIconName(fname, field string) (uint, error) {

	// Try name first
	cp := icon.Codepoint(field)
	if cp != 0 {
		return cp, nil
	}

	// Try to parse as hex value
	cp2, err := strconv.ParseUint(field, 16, 32)
	if err != nil {
		return 0, b.err(fname, fmt.Sprintf("Invalid icon codepoint value/name:%v", field))
	}
	return uint(cp2), nil
}

// parseFloats parses a string with a list of floats with the specified size
// and returns a slice. The specified size is 0 any number of floats is allowed.
// The individual values can be separated by spaces or commas
func (b *Builder) parseFloats(fname, field string, min, max int) ([]float32, error) {

	// Checks if field is empty
	field = strings.Trim(field, " ")
	if field == "" {
		return nil, nil
	}

	// Separate individual fields
	var parts []string
	if strings.Index(field, ",") < 0 {
		parts = strings.Fields(field)
	} else {
		parts = strings.Split(field, ",")
	}
	if len(parts) < min || len(parts) > max {
		return nil, b.err(fname, "Invalid number of float32 values")
	}

	// Parse each field value and appends to slice
	var values []float32
	for i := 0; i < len(parts); i++ {
		val, err := strconv.ParseFloat(strings.Trim(parts[i], " "), 32)
		if err != nil {
			return nil, b.err(fname, err.Error())
		}
		values = append(values, float32(val))
	}
	return values, nil
}

// err creates and returns an error for the current object, field name and with the specified message
func (b *Builder) err(fname, msg string) error {

	return fmt.Errorf("Error in object:%s field:%s -> %s", b.objpath.path(), fname, msg)
}

// strStack is a stack of strings
type strStack struct {
	stack []string
}

// clear removes all elements from the stack
func (ss *strStack) clear() {

	ss.stack = []string{}
}

// push pushes a string to the top of the stack
func (ss *strStack) push(v string) {

	ss.stack = append(ss.stack, v)
}

// pop removes and returns the string at the top of the stack.
// Returns an empty string if the stack is empty
func (ss *strStack) pop() string {

	if len(ss.stack) == 0 {
		return ""
	}
	length := len(ss.stack)
	v := ss.stack[length-1]
	ss.stack = ss.stack[:length-1]
	return v
}

// path returns a string composed of all the strings in the
// stack separated by a forward slash.
func (ss *strStack) path() string {

	return strings.Join(ss.stack, "/")
}
