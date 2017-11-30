// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/g3n/engine/math32"
	"gopkg.in/yaml.v2"
)

// Builder builds GUI objects from a declarative description in YAML format
type Builder struct {
	desc   map[string]*panelDesc
	panels []IPanel // first level panels
}

type panelDesc struct {
	Type        string
	Posx        *float32
	Posy        *float32
	Width       float32
	Height      float32
	Margins     string
	Borders     string
	BorderColor string
	Paddings    string
	Color       string
	Enabled     bool
	Visible     bool
	Renderable  bool
	Children    []panelDesc
	Layout      layoutAttr
	Text        string
	FontSize    float32
	FontDPI     float32
	PlaceHolder string
	MaxLength   uint
}

type layoutAttr struct {
	Type string
}

const (
	descPanel        = "Panel"
	descLabel        = "Label"
	descEdit         = "Edit"
	fieldMargins     = "margins"
	fieldBorders     = "borders"
	fieldBorderColor = "bordercolor"
	fieldPaddings    = "paddings"
	fieldColor       = "color"
	objSingle        = "--single--"
)

//
// NewBuilder creates and returns a pointer to a new gui Builder object
//
func NewBuilder() *Builder {

	b := new(Builder)

	return b
}

//
// BuildFromData builds gui objects from the specified description
// in YAML format
//
func (b *Builder) BuildFromString(desc string) error {

	// Try assuming the description contains a single root panel
	var pd panelDesc
	err := yaml.Unmarshal([]byte(desc), &pd)
	if err != nil {
		return err
	}
	if pd.Type != "" {
		pan, err := b.build(&pd, "", nil)
		if err != nil {
			return err
		}
		b.panels = append(b.panels, pan)
		return nil
	}

	// Try assuming the description is a map of panels
	var pdm map[string]panelDesc
	err = yaml.Unmarshal([]byte(desc), &pdm)
	if err != nil {
		return err
	}
	// Builds each panel in the map
	for name, pd := range pdm {
		pan, err := b.build(&pd, name, nil)
		if err != nil {
			return err
		}
		b.panels = append(b.panels, pan)
	}
	return nil
}

//
// BuildFromFile builds gui objects from the specified file which
// must contain objects descriptions in YAML format
//
func (b *Builder) BuildFromFile(filepath string) error {

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
	return b.BuildFromString(string(data))
}

func (b *Builder) Panels() []IPanel {

	return b.panels
}

func (b *Builder) Label(path string) (*Label, error) {

	return nil, nil
}

//
// build builds gui objects from the specified description and its children recursively
//
func (b *Builder) build(pd *panelDesc, pname string, parent *Panel) (IPanel, error) {

	fmt.Printf("\n%+v\n\n", pd)
	var err error
	var pan IPanel
	switch pd.Type {
	case descPanel:
		pan, err = b.buildPanel(pd, pname)
	case descLabel:
		pan, err = b.buildLabel(pd, pname)
	case descEdit:
		pan, err = b.buildEdit(pd, pname)
	default:
		err = fmt.Errorf("Invalid panel type:%s", pd.Type)
	}
	if err != nil {
		return nil, err
	}
	if parent != nil {
		parent.Add(pan)
	}
	return pan, nil
}

func (b *Builder) buildPanel(pd *panelDesc, pname string) (IPanel, error) {

	log.Error("buildPanel:[%s]", pd.Borders)
	pan := NewPanel(pd.Width, pd.Height)

	// Set margin sizes
	bs, err := b.parseBorderSizes(pname, fieldMargins, pd.Margins)
	if err != nil {
		return nil, err
	}
	if bs != nil {
		pan.SetMarginsFrom(bs)
	}

	// Set border sizes
	bs, err = b.parseBorderSizes(pname, fieldBorders, pd.Borders)
	if err != nil {
		return nil, err
	}
	if bs != nil {
		pan.SetBordersFrom(bs)
	}

	// Set border color
	c, err := b.parseColor(pname, fieldBorderColor, pd.BorderColor)
	if err != nil {
		return nil, err
	}
	if c != nil {
		pan.SetBordersColor4(c)
	}

	// Set paddings sizes
	bs, err = b.parseBorderSizes(pname, fieldPaddings, pd.Paddings)
	if err != nil {
		return nil, err
	}
	if bs != nil {
		pan.SetPaddingsFrom(bs)
	}

	// Set color
	c, err = b.parseColor(pname, fieldColor, pd.Color)
	if err != nil {
		return nil, err
	}
	if c != nil {
		pan.SetColor4(c)
	}

	return pan, nil
}

func (b *Builder) buildLabel(pa *panelDesc, name string) (IPanel, error) {

	label := NewLabel("df")

	return label, nil
}

func (b *Builder) buildEdit(pa *panelDesc, name string) (IPanel, error) {

	return nil, nil
}

//
// parseBorderSizes parses a string field which can contain one float value or
// float values. In the first case all borders has the same width
//
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
		parts = strings.Split(field, " ")
	} else {
		parts = strings.Split(field, ",")
	}
	if len(parts) < min || len(parts) > max {
		return nil, b.err(pname, fname, "Invalid number of float32 values")
	}

	// Parse each field value and appends to slice
	var values []float32
	for i := 0; i < len(parts); i++ {
		val, err := strconv.ParseFloat(parts[i], 32)
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
