// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package collada

import (
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
)

//
// LibraryEffects
//
type LibraryEffects struct {
	Id     string
	Name   string
	Asset  *Asset
	Effect []*Effect
}

// Dump prints out information about the LibraryEffects
func (le *LibraryEffects) Dump(out io.Writer, indent int) {

	if le == nil {
		return
	}
	fmt.Fprintf(out, "%sLibraryEffects id:%s name:%s\n", sIndent(indent), le.Id, le.Name)
	for _, ef := range le.Effect {
		ef.Dump(out, indent+step)
	}
}

//
// Effect
//
type Effect struct {
	Id      string
	Name    string
	Asset   *Asset
	Profile []interface{}
}

// Dump prints out information about the Effect
func (ef *Effect) Dump(out io.Writer, indent int) {

	fmt.Printf("%sEffect id:%s name:%s\n", sIndent(indent), ef.Id, ef.Name)
	ind := indent + step
	for _, p := range ef.Profile {
		switch pt := p.(type) {
		case *ProfileCOMMON:
			pt.Dump(out, ind)
			break
		}
	}
}

//
// ProfileCOMMON
//
type ProfileCOMMON struct {
	Id        string
	Asset     *Asset
	Newparam  []*Newparam
	Technique struct {
		Id            string
		Sid           string
		Asset         *Asset
		ShaderElement interface{} // Blinn|Constant|Lambert|Phong
	}
}

// Dump prints out information about the ProfileCOMMON
func (pc *ProfileCOMMON) Dump(out io.Writer, indent int) {

	fmt.Printf("%sProfileCOMMON id:%s\n", sIndent(indent), pc.Id)
	ind := indent + step

	for _, np := range pc.Newparam {
		np.Dump(out, ind)
	}

	fmt.Printf("%sTechnique id:%s sid:%s\n", sIndent(ind), pc.Technique.Id, pc.Technique.Sid)
	ind += step
	switch sh := pc.Technique.ShaderElement.(type) {
	case *Phong:
		sh.Dump(out, ind)
		break
	}
}

//
// Newparam
//
type Newparam struct {
	Sid           string
	Semantic      string
	ParameterType interface{}
}

// Dump prints out information about the Newparam
func (np *Newparam) Dump(out io.Writer, indent int) {

	fmt.Printf("%sNewparam sid:%s\n", sIndent(indent), np.Sid)
	ind := indent + step
	switch pt := np.ParameterType.(type) {
	case *Surface:
		pt.Dump(out, ind)
	case *Sampler2D:
		pt.Dump(out, ind)
	}
}

//
// Surface
//
type Surface struct {
	Type string
	Init interface{}
}

// Dump prints out information about the Surface
func (sf *Surface) Dump(out io.Writer, indent int) {

	fmt.Printf("%sSurface type:%s\n", sIndent(indent), sf.Type)
	ind := indent + step
	switch it := sf.Init.(type) {
	case InitFrom:
		it.Dump(out, ind)
	}
}

//
// Sampler2D
//
type Sampler2D struct {
	Source string
}

// Dump prints out information about the Sampler2D
func (sp *Sampler2D) Dump(out io.Writer, indent int) {

	fmt.Printf("%sSampler2D\n", sIndent(indent))
	ind := indent + step
	fmt.Printf("%sSource:%s\n", sIndent(ind), sp.Source)
}

//
// Blinn
//
type Blinn struct {
	Emission          interface{}
	Ambient           interface{}
	Diffuse           interface{}
	Specular          interface{}
	Shininess         interface{}
	Reflective        interface{}
	Reflectivity      interface{}
	Transparent       interface{}
	Transparency      interface{}
	IndexOfRefraction interface{}
}

// Dump prints out information about the Blinn
func (bl *Blinn) Dump(out io.Writer, indent int) {

	fmt.Fprintf(out, "%sBlinn\n", sIndent(indent))
	ind := indent + step
	DumpColorOrTexture("Emssion", bl.Emission, out, ind)
	DumpColorOrTexture("Ambient", bl.Ambient, out, ind)
	DumpColorOrTexture("Diffuse", bl.Diffuse, out, ind)
	DumpColorOrTexture("Specular", bl.Specular, out, ind)
	DumpFloatOrParam("Shininess", bl.Shininess, out, ind)
	DumpColorOrTexture("Reflective", bl.Reflective, out, ind)
	DumpFloatOrParam("Reflectivity", bl.Reflectivity, out, ind)
	DumpColorOrTexture("Transparent", bl.Transparent, out, ind)
	DumpFloatOrParam("Transparency", bl.Transparency, out, ind)
	DumpFloatOrParam("IndexOfRefraction", bl.IndexOfRefraction, out, ind)
}

//
// Constant
//
type Constant struct {
	Emission          interface{}
	Reflective        interface{}
	Reflectivity      interface{}
	Transparent       interface{}
	Transparency      interface{}
	IndexOfRefraction interface{}
}

//
// Lambert
//
type Lambert struct {
	Emission          interface{}
	Ambient           interface{}
	Diffuse           interface{}
	Reflective        interface{}
	Reflectivity      interface{}
	Transparent       interface{}
	Transparency      interface{}
	IndexOfRefraction interface{}
}

//
// Phong
//
type Phong struct {
	Emission          interface{}
	Ambient           interface{}
	Diffuse           interface{}
	Specular          interface{}
	Shininess         interface{}
	Reflective        interface{}
	Reflectivity      interface{}
	Transparent       interface{}
	Transparency      interface{}
	IndexOfRefraction interface{}
}

// Dump prints out information about the Phong
func (ph *Phong) Dump(out io.Writer, indent int) {

	fmt.Fprintf(out, "%sPhong\n", sIndent(indent))
	ind := indent + step
	DumpColorOrTexture("Emission", ph.Emission, out, ind)
	DumpColorOrTexture("Ambient", ph.Ambient, out, ind)
	DumpColorOrTexture("Diffuse", ph.Diffuse, out, ind)
	DumpColorOrTexture("Specular", ph.Specular, out, ind)
	DumpFloatOrParam("Shininess", ph.Shininess, out, ind)
	DumpColorOrTexture("Reflective", ph.Reflective, out, ind)
	DumpFloatOrParam("Reflectivity", ph.Reflectivity, out, ind)
	DumpColorOrTexture("Transparent", ph.Transparent, out, ind)
	DumpFloatOrParam("Transparency", ph.Transparency, out, ind)
	DumpFloatOrParam("IndexOfRefraction", ph.IndexOfRefraction, out, ind)
}

// DumpColorOrTexture prints out information about the Color or Texture
func DumpColorOrTexture(name string, v interface{}, out io.Writer, indent int) {

	if v == nil {
		return
	}
	fmt.Fprintf(out, "%s%s\n", sIndent(indent), name)
	ind := indent + step
	switch vt := v.(type) {
	case *Color:
		vt.Dump(out, ind)
	case *Texture:
		vt.Dump(out, ind)
	}
}

// DumpFloatOrParam prints out information about the Float or Param
func DumpFloatOrParam(name string, v interface{}, out io.Writer, indent int) {

	if v == nil {
		return
	}
	fmt.Fprintf(out, "%s%s\n", sIndent(indent), name)
	ind := indent + step
	switch vt := v.(type) {
	case *Float:
		vt.Dump(out, ind)
		break
	}
}

//
// Color
//
type Color struct {
	Sid  string
	Data [4]float32
}

// Dump prints out information about the Color
func (c *Color) Dump(out io.Writer, indent int) {

	fmt.Fprintf(out, "%sColor sid:%s data:%v\n", sIndent(indent), c.Sid, c.Data)
}

//
// Float
//
type Float struct {
	Sid  string
	Data float32
}

// Dump prints out information about the Float
func (f *Float) Dump(out io.Writer, indent int) {

	fmt.Fprintf(out, "%sFloat sid:%s data:%v\n", sIndent(indent), f.Sid, f.Data)
}

//
// Texture
//
type Texture struct {
	Texture  string
	Texcoord string
}

// Dump prints out information about the Texture
func (t *Texture) Dump(out io.Writer, indent int) {

	fmt.Fprintf(out, "%sTexture texture:%s texcoord:%v\n", sIndent(indent), t.Texture, t.Texcoord)
}

func (d *Decoder) decLibraryEffects(start xml.StartElement, dom *Collada) error {

	le := new(LibraryEffects)
	dom.LibraryEffects = le
	le.Id = findAttrib(start, "id").Value
	le.Name = findAttrib(start, "name").Value

	for {
		// Get next child element
		child, _, err := d.decNextChild(start)
		if err != nil || child.Name.Local == "" {
			return err
		}
		// Decodes <effect>
		if child.Name.Local == "effect" {
			err := d.decEffect(child, le)
			if err != nil {
				return err
			}
			continue
		}
	}
}

func (d *Decoder) decEffect(start xml.StartElement, le *LibraryEffects) error {

	e := new(Effect)
	e.Id = findAttrib(start, "id").Value
	e.Name = findAttrib(start, "name").Value
	le.Effect = append(le.Effect, e)

	for {
		// Get next child element
		child, _, err := d.decNextChild(start)
		if err != nil || child.Name.Local == "" {
			return err
		}
		if child.Name.Local == "profile_COMMON" {
			err := d.decEffectProfileCommon(child, e)
			if err != nil {
				return err
			}
			continue
		}
	}
}

func (d *Decoder) decEffectProfileCommon(start xml.StartElement, e *Effect) error {

	pc := new(ProfileCOMMON)
	pc.Id = findAttrib(start, "id").Value
	e.Profile = append(e.Profile, pc)

	for {
		// Get next child element
		child, _, err := d.decNextChild(start)
		if err != nil || child.Name.Local == "" {
			return err
		}
		if child.Name.Local == "newparam" {
			err := d.decProfileCommonNewparam(child, pc)
			if err != nil {
				return err
			}
			continue
		}
		if child.Name.Local == "technique" {
			err := d.decProfileCommonTechnique(child, pc)
			if err != nil {
				return err
			}
			continue
		}
	}
}

func (d *Decoder) decProfileCommonNewparam(start xml.StartElement, pc *ProfileCOMMON) error {

	np := new(Newparam)
	np.Sid = findAttrib(start, "sid").Value
	pc.Newparam = append(pc.Newparam, np)

	for {
		child, _, err := d.decNextChild(start)
		if err != nil || child.Name.Local == "" {
			return err
		}
		if child.Name.Local == "surface" {
			err := d.decSurface(child, np)
			if err != nil {
				return err
			}
			continue
		}
		if child.Name.Local == "sampler2D" {
			err := d.decSampler2D(child, np)
			if err != nil {
				return err
			}
			continue
		}
	}
}

func (d *Decoder) decSurface(start xml.StartElement, np *Newparam) error {

	sf := new(Surface)
	sf.Type = findAttrib(start, "type").Value
	np.ParameterType = sf

	for {
		child, data, err := d.decNextChild(start)
		if err != nil || child.Name.Local == "" {
			return err
		}
		if child.Name.Local == "init_from" {
			sf.Init = InitFrom{string(data)}
			continue
		}
	}
}

func (d *Decoder) decSampler2D(start xml.StartElement, np *Newparam) error {

	sp := new(Sampler2D)
	np.ParameterType = sp

	for {
		child, data, err := d.decNextChild(start)
		if err != nil || child.Name.Local == "" {
			return err
		}
		if child.Name.Local == "source" {
			sp.Source = string(data)
			continue
		}
	}
}

func (d *Decoder) decProfileCommonTechnique(start xml.StartElement, pc *ProfileCOMMON) error {

	pc.Technique.Id = findAttrib(start, "id").Value
	pc.Technique.Sid = findAttrib(start, "sid").Value

	for {
		child, _, err := d.decNextChild(start)
		if err != nil || child.Name.Local == "" {
			return err
		}
		if child.Name.Local == "blinn" {
			err := d.decBlinn(child, pc)
			if err != nil {
				return err
			}
			continue
		}
		if child.Name.Local == "constant" {
			log.Warn("CONSTANT not implemented yet")
			continue
		}
		if child.Name.Local == "lambert" {
			log.Warn("LAMBERT not implemented yet")
			continue
		}
		if child.Name.Local == "phong" {
			err := d.decPhong(child, pc)
			if err != nil {
				return err
			}
			continue
		}
	}
}

func (d *Decoder) decBlinn(start xml.StartElement, pc *ProfileCOMMON) error {

	bl := new(Blinn)
	pc.Technique.ShaderElement = bl

	for {
		child, _, err := d.decNextChild(start)
		if err != nil || child.Name.Local == "" {
			return err
		}
		if child.Name.Local == "emission" {
			err := d.decColorOrTexture(child, &bl.Emission)
			if err != nil {
				return err
			}
			continue
		}
		if child.Name.Local == "ambient" {
			err := d.decColorOrTexture(child, &bl.Ambient)
			if err != nil {
				return err
			}
			continue
		}
		if child.Name.Local == "diffuse" {
			err := d.decColorOrTexture(child, &bl.Diffuse)
			if err != nil {
				return err
			}
			continue
		}
		if child.Name.Local == "specular" {
			err := d.decColorOrTexture(child, &bl.Specular)
			if err != nil {
				return err
			}
			continue
		}
		if child.Name.Local == "shininess" {
			err := d.decFloatOrParam(child, &bl.Shininess)
			if err != nil {
				return err
			}
			continue
		}
		if child.Name.Local == "reflective" {
			err := d.decFloatOrParam(child, &bl.Reflective)
			if err != nil {
				return err
			}
			continue
		}
		if child.Name.Local == "reflectivity" {
			err := d.decFloatOrParam(child, &bl.Reflectivity)
			if err != nil {
				return err
			}
			continue
		}
		if child.Name.Local == "transparent" {
			// not supported
			continue
		}
		if child.Name.Local == "transparency" {
			err := d.decFloatOrParam(child, &bl.Transparency)
			if err != nil {
				return err
			}
			continue
		}
		if child.Name.Local == "index_of_refraction" {
			err := d.decFloatOrParam(child, &bl.IndexOfRefraction)
			if err != nil {
				return err
			}
			continue
		}
	}
}

func (d *Decoder) decPhong(start xml.StartElement, pc *ProfileCOMMON) error {

	ph := new(Phong)
	pc.Technique.ShaderElement = ph

	for {
		child, _, err := d.decNextChild(start)
		if err != nil || child.Name.Local == "" {
			return err
		}
		if child.Name.Local == "emission" {
			err := d.decColorOrTexture(child, &ph.Emission)
			if err != nil {
				return err
			}
			continue
		}
		if child.Name.Local == "ambient" {
			err := d.decColorOrTexture(child, &ph.Ambient)
			if err != nil {
				return err
			}
			continue
		}
		if child.Name.Local == "diffuse" {
			err := d.decColorOrTexture(child, &ph.Diffuse)
			if err != nil {
				return err
			}
			continue
		}
		if child.Name.Local == "specular" {
			err := d.decColorOrTexture(child, &ph.Specular)
			if err != nil {
				return err
			}
			continue
		}
		if child.Name.Local == "shininess" {
			err := d.decFloatOrParam(child, &ph.Shininess)
			if err != nil {
				return err
			}
			continue
		}
		if child.Name.Local == "reflective" {
			err := d.decFloatOrParam(child, &ph.Reflective)
			if err != nil {
				return err
			}
			continue
		}
		if child.Name.Local == "reflectivity" {
			err := d.decFloatOrParam(child, &ph.Reflectivity)
			if err != nil {
				return err
			}
			continue
		}
		if child.Name.Local == "transparent" {
			// not supported
			continue
		}
		if child.Name.Local == "transparency" {
			err := d.decFloatOrParam(child, &ph.Transparency)
			if err != nil {
				return err
			}
			continue
		}
		if child.Name.Local == "index_of_refraction" {
			err := d.decFloatOrParam(child, &ph.IndexOfRefraction)
			if err != nil {
				return err
			}
			continue
		}
	}
}

func (d *Decoder) decColorOrTexture(start xml.StartElement, dest *interface{}) error {

	child, cdata, err := d.decNextChild(start)
	if err != nil || child.Name.Local == "" {
		return err
	}
	if child.Name.Local == "color" {
		c := &Color{}
		c.Sid = findAttrib(child, "sid").Value
		*dest = c
		var br bytesReader
		br.Init(cdata)
		idx := 0
		for {
			tok := br.TokenNext()
			if tok == nil || len(tok) == 0 {
				break
			}
			v, err := strconv.ParseFloat(string(tok), 32)
			if err != nil {
				return err
			}
			c.Data[idx] = float32(v)
			idx++
		}
		return nil
	}
	if child.Name.Local == "texture" {
		t := &Texture{}
		t.Texture = findAttrib(child, "texture").Value
		t.Texcoord = findAttrib(child, "texcoord").Value
		*dest = t
		return nil
	}
	if child.Name.Local == "param" {
		return fmt.Errorf("not supported")
	}
	return nil
}

func (d *Decoder) decFloatOrParam(start xml.StartElement, dest *interface{}) error {

	child, cdata, err := d.decNextChild(start)
	if err != nil || child.Name.Local == "" {
		return err
	}
	if child.Name.Local == "float" {
		f := &Float{}
		f.Sid = findAttrib(child, "sid").Value
		*dest = f
		v, err := strconv.ParseFloat(string(cdata), 32)
		if err != nil {
			return err
		}
		f.Data = float32(v)
		return nil
	}

	return nil
}
