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

	"github.com/g3n/engine/math32"
	"gopkg.in/yaml.v2"
)

// Builder builds GUI objects from a declarative description in YAML format
type Builder struct {
	desc    map[string]*panelDesc // parsed descriptions
	imgpath string                // base path for image panels files
}

type panelStyle struct {
	Borders     string
	Paddings    string
	BorderColor string
	BgColor     string
	FgColor     string
}

type panelStyles struct {
	Normal   panelStyle
	Over     panelStyle
	Focus    panelStyle
	Pressed  panelStyle
	Disabled panelStyle
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
	Imagefile    string // Optional image filepath for ImagePanel
	Children     []*panelDesc
	Layout       layoutAttr
	Styles       *panelStyles
	Text         string
	BgColor      string
	FontColor    string // Optional
	FontSize     *float32
	FontDPI      *float32
	LineSpacing  *float32
	PlaceHolder  string
	MaxLength    *uint
}

type layoutAttr struct {
	Type string
}

const (
	descTypePanel      = "Panel"
	descTypeImagePanel = "ImagePanel"
	descTypeLabel      = "Label"
	descTypeEdit       = "Edit"
	fieldMargins       = "margins"
	fieldBorders       = "borders"
	fieldBorderColor   = "bordercolor"
	fieldPaddings      = "paddings"
	fieldColor         = "color"
)

//
// NewBuilder creates and returns a pointer to a new gui Builder object
//
func NewBuilder() *Builder {

	return new(Builder)
}

//
// ParseString parses a string with gui objects descriptions in YAML format
// It there was a previously parsed description, it is cleared.
//
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
		fmt.Printf("\n%+v\n", b.desc)
		return nil
	}

	// Try assuming the description is a map of panels
	var pdm map[string]*panelDesc
	err = yaml.Unmarshal([]byte(desc), &pdm)
	if err != nil {
		return err
	}
	b.desc = pdm
	fmt.Printf("\n%+v\n", b.desc)
	return nil
}

//
// ParseFile builds gui objects from the specified file which
// must contain objects descriptions in YAML format
//
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

//
// Names returns a sorted list of names of top level previously parsed objects.
// If there is only a single object with no name, its name is returned
// as an empty string
//
func (b *Builder) Names() []string {

	var objs []string
	for name, _ := range b.desc {
		objs = append(objs, name)
	}
	sort.Strings(objs)
	return objs
}

//
// Build builds a gui object and all its children recursively.
// The specified name should be a top level name from a
// from a previously parsed description
// If the descriptions contains a single object with no name,
// It should be specified the empty string to build this object.
//
func (b *Builder) Build(name string) (IPanel, error) {

	pd, ok := b.desc[name]
	if !ok {
		return nil, fmt.Errorf("Object name:%s not found", name)
	}
	return b.build(pd, nil)
}

// Sets the path for image panels relative image files
func (b *Builder) SetImagepath(path string) {

	b.imgpath = path
}

//
// build builds the gui object from the specified description.
// All its children are also built recursively
// Returns the built object or an error
//
func (b *Builder) build(pd *panelDesc, iparent IPanel) (IPanel, error) {

	fmt.Printf("\n%+v\n\n", pd)
	var err error
	var pan IPanel
	switch pd.Type {
	case descTypePanel:
		pan, err = b.buildPanel(pd)
	case descTypeImagePanel:
		pan, err = b.buildImagePanel(pd)
	case descTypeLabel:
		pan, err = b.buildLabel(pd)
	case descTypeEdit:
		pan, err = b.buildEdit(pd)
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
		child, err := b.build(pd.Children[i], pan)
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
		return nil, b.err(pd.Name, "Imagefile", "Imagefile must be supplied")
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
		child, err := b.build(pd.Children[i], panel)
		if err != nil {
			return nil, err
		}
		panel.Add(child)
	}
	return panel, nil
}

func (b *Builder) buildLabel(pd *panelDesc) (IPanel, error) {

	// Builds panel and set common attributes
	label := NewLabel(pd.Text)
	err := b.setCommon(pd, label)
	if err != nil {
		return nil, err
	}

	// Set optional background color
	c, err := b.parseColor(pd.Name, "bgcolor", pd.BgColor)
	if err != nil {
		return nil, err
	}
	if c != nil {
		label.SetBgColor4(c)
	}

	// Set optional font color
	c, err = b.parseColor(pd.Name, "fontcolor", pd.FontColor)
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

func (b *Builder) buildEdit(pa *panelDesc) (IPanel, error) {

	return nil, nil
}

// setCommon sets the common attributes in the description to the specified panel
func (b *Builder) setCommon(pd *panelDesc, ipan IPanel) error {

	// Set optional position
	panel := ipan.GetPanel()
	if pd.Position != "" {
		va, err := b.parseFloats(pd.Name, "position", pd.Position, 2, 2)
		if va == nil || err != nil {
			return err
		}
		panel.SetPosition(va[0], va[1])
	}

	// Set optional margin sizes
	bs, err := b.parseBorderSizes(pd.Name, fieldMargins, pd.Margins)
	if err != nil {
		return err
	}
	if bs != nil {
		panel.SetMarginsFrom(bs)
	}

	// Set optional border sizes
	bs, err = b.parseBorderSizes(pd.Name, fieldBorders, pd.Borders)
	if err != nil {
		return err
	}
	if bs != nil {
		panel.SetBordersFrom(bs)
	}

	// Set optional border color
	c, err := b.parseColor(pd.Name, fieldBorderColor, pd.BorderColor)
	if err != nil {
		return err
	}
	if c != nil {
		panel.SetBordersColor4(c)
	}

	// Set optional paddings sizes
	bs, err = b.parseBorderSizes(pd.Name, fieldPaddings, pd.Paddings)
	if err != nil {
		return err
	}
	if bs != nil {
		panel.SetPaddingsFrom(bs)
	}

	// Set optional color
	c, err = b.parseColor(pd.Name, fieldColor, pd.Color)
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
func (b *Builder) parseBorderSizes(pname, fname, field string) (*BorderSizes, error) {

	va, err := b.parseFloats(pname, fname, field, 1, 4)
	if va == nil || err != nil {
		return nil, err
	}
	if len(va) == 1 {
		return &BorderSizes{va[0], va[0], va[0], va[0]}, nil
	}
	return &BorderSizes{va[0], va[1], va[2], va[3]}, nil
}

//
// parseColor parses a string field which can contain a color name or
// a list of 3 or 4 float values for the color components
//
func (b *Builder) parseColor(pname, fname, field string) (*math32.Color4, error) {

	// Checks if field is empty
	field = strings.Trim(field, " ")
	if field == "" {
		return nil, nil
	}

	// Checks if field is a color name
	value := math32.ColorUint(field)
	if value != 0 {
		var c math32.Color
		c.SetName(field)
		return &math32.Color4{c.R, c.G, c.B, 1}, nil
	}

	// Accept 3 or 4 floats values
	va, err := b.parseFloats(pname, fname, field, 3, 4)
	if err != nil {
		return nil, err
	}
	if len(va) == 3 {
		return &math32.Color4{va[0], va[1], va[2], 1}, nil
	}
	return &math32.Color4{va[0], va[1], va[2], va[3]}, nil
}

//
// parseFloats parses a string with a list of floats with the specified size
// and returns a slice. The specified size is 0 any number of floats is allowed.
// The individual values can be separated by spaces or commas
//
func (b *Builder) parseFloats(pname, fname, field string, min, max int) ([]float32, error) {

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
		return nil, b.err(pname, fname, "Invalid number of float32 values")
	}

	// Parse each field value and appends to slice
	var values []float32
	for i := 0; i < len(parts); i++ {
		val, err := strconv.ParseFloat(strings.Trim(parts[i], " "), 32)
		if err != nil {
			return nil, fmt.Errorf("Error parsing float32 field:[%s]: %s", field, err)
		}
		values = append(values, float32(val))
	}
	return values, nil
}

func (b *Builder) err(pname, fname, msg string) error {

	return fmt.Errorf("Error in object:%s field:%s -> %s", pname, fname, msg)
}
