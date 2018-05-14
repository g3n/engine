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
// LibraryLights
//
type LibraryLights struct {
	Id    string
	Name  string
	Asset *Asset
	Light []*Light
}

// Dump prints out information about the LibraryLights
func (ll *LibraryLights) Dump(out io.Writer, indent int) {

	if ll == nil {
		return
	}
	fmt.Fprintf(out, "%sLibraryLights id:%s name:%s\n", sIndent(indent), ll.Id, ll.Name)
	for _, light := range ll.Light {
		light.Dump(out, indent+step)
	}
}

//
// Light
//
type Light struct {
	Id              string
	Name            string
	TechniqueCommon struct {
		Type interface{}
	}
}

// Dump prints out information about the Light
func (li *Light) Dump(out io.Writer, indent int) {

	fmt.Fprintf(out, "%sLights id:%s name:%s\n", sIndent(indent), li.Id, li.Name)
	ind := indent + step
	fmt.Fprintf(out, "%sTechniqueCommon\n", sIndent(ind))
	ind += step
	switch lt := li.TechniqueCommon.Type.(type) {
	case *Ambient:
		lt.Dump(out, ind)
	case *Directional:
		lt.Dump(out, ind)
	case *Point:
		lt.Dump(out, ind)
	case *Spot:
		lt.Dump(out, ind)
	}
}

//
// Ambient
//
type Ambient struct {
	Color LightColor
}

// Dump prints out information about the Ambient
func (amb *Ambient) Dump(out io.Writer, indent int) {

	fmt.Fprintf(out, "%sAmbient\n", sIndent(indent))
	ind := indent + step
	amb.Color.Dump(out, ind)
}

//
// Directional
//
type Directional struct {
	Color LightColor
}

// Dump prints out information about the Directional
func (dir *Directional) Dump(out io.Writer, indent int) {

	fmt.Fprintf(out, "%sDirectional\n", sIndent(indent))
	ind := indent + step
	dir.Color.Dump(out, ind)
}

//
// Point
//
type Point struct {
	Color                LightColor
	ConstantAttenuation  *FloatValue
	LinearAttenuation    *FloatValue
	QuadraticAttenuation *FloatValue
}

// Dump prints out information about the Point
func (pl *Point) Dump(out io.Writer, indent int) {

	fmt.Fprintf(out, "%sPoint\n", sIndent(indent))
	indent += step
	pl.Color.Dump(out, indent)
	pl.ConstantAttenuation.Dump("ConstantAttenuation", out, indent)
	pl.LinearAttenuation.Dump("LinearAttenuation", out, indent)
	pl.QuadraticAttenuation.Dump("QuadraticAttenuation", out, indent)
}

//
// Spot
//
type Spot struct {
	Color                LightColor
	ConstantAttenuation  *FloatValue
	LinearAttenuation    *FloatValue
	QuadraticAttenuation *FloatValue
	FalloffAngle         *FloatValue
	FalloffExponent      *FloatValue
}

// Dump prints out information about the Spot
func (sl *Spot) Dump(out io.Writer, indent int) {

	fmt.Fprintf(out, "%sSpot\n", sIndent(indent))
	indent += step
	sl.Color.Dump(out, indent)
	sl.ConstantAttenuation.Dump("ConstantAttenuation", out, indent)
	sl.LinearAttenuation.Dump("LinearAttenuation", out, indent)
	sl.QuadraticAttenuation.Dump("QuadraticAttenuation", out, indent)
	sl.FalloffAngle.Dump("FalloffAngle", out, indent)
	sl.FalloffExponent.Dump("FalloffExponent", out, indent)
}

//
// FloatValue
//
type FloatValue struct {
	Sid   string
	Value float32
}

// Dump prints out information about the FloatValue
func (fv *FloatValue) Dump(name string, out io.Writer, indent int) {

	if fv == nil {
		return
	}
	fmt.Fprintf(out, "%s%s sid:%s value:%v\n", sIndent(indent), name, fv.Sid, fv.Value)
}

//
// LightColor
//
type LightColor struct {
	Sid  string
	Data [3]float32
}

// Dump prints out information about the LightColor
func (lc *LightColor) Dump(out io.Writer, indent int) {

	fmt.Fprintf(out, "%sColor sid:%s data:%v\n", sIndent(indent), lc.Sid, lc.Data)
}

func (d *Decoder) decLibraryLights(start xml.StartElement, dom *Collada) error {

	ll := new(LibraryLights)
	dom.LibraryLights = ll
	ll.Id = findAttrib(start, "id").Value
	ll.Name = findAttrib(start, "name").Value

	for {
		child, _, err := d.decNextChild(start)
		if err != nil || child.Name.Local == "" {
			return err
		}
		if child.Name.Local == "light" {
			err := d.decLight(child, ll)
			if err != nil {
				return err
			}
			continue
		}
	}
}

func (d *Decoder) decLight(start xml.StartElement, ll *LibraryLights) error {

	light := new(Light)
	ll.Light = append(ll.Light, light)
	light.Id = findAttrib(start, "id").Value
	light.Name = findAttrib(start, "name").Value

	for {
		child, _, err := d.decNextChild(start)
		if err != nil || child.Name.Local == "" {
			return err
		}
		if child.Name.Local == "technique_common" {
			err := d.decLightTechniqueCommon(child, light)
			if err != nil {
				return err
			}
			continue
		}
	}
}

func (d *Decoder) decLightTechniqueCommon(start xml.StartElement, li *Light) error {

	child, _, err := d.decNextChild(start)
	if err != nil || child.Name.Local == "" {
		return err
	}
	if child.Name.Local == "ambient" {
		err := d.decAmbient(child, li)
		if err != nil {
			return err
		}
	}
	if child.Name.Local == "directional" {
		err := d.decDirectional(child, li)
		if err != nil {
			return err
		}
	}
	if child.Name.Local == "point" {
		err := d.decPoint(child, li)
		if err != nil {
			return err
		}
	}
	if child.Name.Local == "spot" {
		err := d.decSpot(child, li)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Decoder) decAmbient(start xml.StartElement, li *Light) error {

	amb := new(Ambient)
	li.TechniqueCommon.Type = amb

	child, cdata, err := d.decNextChild(start)
	if err != nil || child.Name.Local == "" {
		return err
	}
	if child.Name.Local == "color" {
		err := d.decLightColor(child, cdata, &amb.Color)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Decoder) decDirectional(start xml.StartElement, li *Light) error {

	dir := new(Directional)
	li.TechniqueCommon.Type = dir

	child, cdata, err := d.decNextChild(start)
	if err != nil || child.Name.Local == "" {
		return err
	}
	if child.Name.Local == "color" {
		err := d.decLightColor(child, cdata, &dir.Color)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Decoder) decPoint(start xml.StartElement, li *Light) error {

	pl := new(Point)
	li.TechniqueCommon.Type = pl

	for {
		child, cdata, err := d.decNextChild(start)
		if err != nil || child.Name.Local == "" {
			return err
		}
		if child.Name.Local == "color" {
			err := d.decLightColor(child, cdata, &pl.Color)
			if err != nil {
				return err
			}
			continue
		}
		if child.Name.Local == "constant_attenuation" {
			fv, err := d.decFloatValue(child, cdata)
			if err != nil {
				return err
			}
			pl.ConstantAttenuation = fv
			continue
		}
		if child.Name.Local == "linear_attenuation" {
			fv, err := d.decFloatValue(child, cdata)
			if err != nil {
				return err
			}
			pl.LinearAttenuation = fv
			continue
		}
		if child.Name.Local == "quadratic_attenuation" {
			fv, err := d.decFloatValue(child, cdata)
			if err != nil {
				return err
			}
			pl.QuadraticAttenuation = fv
			continue
		}
	}
}

func (d *Decoder) decSpot(start xml.StartElement, li *Light) error {

	sl := new(Spot)
	li.TechniqueCommon.Type = sl

	for {
		child, cdata, err := d.decNextChild(start)
		if err != nil || child.Name.Local == "" {
			return err
		}
		if child.Name.Local == "color" {
			err := d.decLightColor(child, cdata, &sl.Color)
			if err != nil {
				return err
			}
			continue
		}
		if child.Name.Local == "constant_attenuation" {
			fv, err := d.decFloatValue(child, cdata)
			if err != nil {
				return err
			}
			sl.ConstantAttenuation = fv
			continue
		}
		if child.Name.Local == "linear_attenuation" {
			fv, err := d.decFloatValue(child, cdata)
			if err != nil {
				return err
			}
			sl.LinearAttenuation = fv
			continue
		}
		if child.Name.Local == "quadratic_attenuation" {
			fv, err := d.decFloatValue(child, cdata)
			if err != nil {
				return err
			}
			sl.QuadraticAttenuation = fv
			continue
		}
		if child.Name.Local == "falloff_angle" {
			fv, err := d.decFloatValue(child, cdata)
			if err != nil {
				return err
			}
			sl.FalloffAngle = fv
			continue
		}
		if child.Name.Local == "falloff_exponent" {
			fv, err := d.decFloatValue(child, cdata)
			if err != nil {
				return err
			}
			sl.FalloffExponent = fv
			continue
		}
	}
}

func (d *Decoder) decFloatValue(start xml.StartElement, cdata []byte) (*FloatValue, error) {

	fv := new(FloatValue)
	fv.Sid = findAttrib(start, "sid").Value
	v, err := strconv.ParseFloat(string(cdata), 32)
	if err != nil {
		return nil, err
	}
	fv.Value = float32(v)
	return fv, nil
}

func (d *Decoder) decLightColor(start xml.StartElement, cdata []byte, lc *LightColor) error {

	lc.Sid = findAttrib(start, "sid").Value
	return decFloat32Sequence(cdata, lc.Data[:])
}
