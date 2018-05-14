// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package collada

import (
	"fmt"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/texture"
	"path/filepath"
	"strings"
)

// GetMaterial returns a pointer to an instance of the material
// with the specified id in the Collada document and an error.
// If no previous instance of the material was found it is created.
func (d *Decoder) GetMaterial(id string) (material.IMaterial, error) {

	// If material already created, returns it
	mat := d.materials[id]
	if mat != nil {
		return mat, nil
	}

	// Creates material and saves it associated with its id
	mat, err := d.NewMaterial(id)
	if err != nil {
		return nil, err
	}
	d.materials[id] = mat
	return mat, nil
}

// NewMaterial creates and returns a pointer to a new material
// from the specified material id/url in the dom
func (d *Decoder) NewMaterial(id string) (material.IMaterial, error) {

	id = strings.TrimPrefix(id, "#")
	// Looks for material with specified id
	mat := findMaterial(&d.dom, id)
	if mat == nil {
		return nil, fmt.Errorf("Material id:%s not found", id)
	}

	// Looks for associated effect
	effect := findEffect(&d.dom, mat.InstanceEffect.Url)
	if effect == nil {
		return nil, fmt.Errorf("Effect id:%s not found", mat.InstanceEffect.Url)
	}

	// Looks for ProfileCOMMON
	pc := findProfileCOMMON(effect)
	if pc == nil {
		return nil, fmt.Errorf("ProfileCOMMON not found")
	}

	switch se := pc.Technique.ShaderElement.(type) {
	case *Blinn:
		return d.newBlinnMaterial(se)
	case *Constant:
		return d.newConstantMaterial(se)
	case *Lambert:
		return d.newLambertMaterial(se)
	case *Phong:
		return d.newPhongMaterial(se)
	default:
		return nil, fmt.Errorf("Invalid shader element")
	}
}

// GetTexture2D returns a pointer to an instance of the Texture2D
// with the specified id in the Collada document and an error.
// If no previous instance of the texture was found it is created.
func (d *Decoder) GetTexture2D(id string) (*texture.Texture2D, error) {

	// If texture already created, returns it
	tex := d.tex2D[id]
	if tex != nil {
		return tex, nil
	}

	// Creates texture and saves it associated with its id
	tex, err := d.NewTexture2D(id)
	if err != nil {
		return nil, err
	}
	d.tex2D[id] = tex
	return tex, nil
}

// NewTexture2D creates and returns a pointer to a new Texture2D
// from the specified sampler2D id/url in the dom
func (d *Decoder) NewTexture2D(id string) (*texture.Texture2D, error) {

	// Find newparam in all effects profiles with the specified id
	np := findNewparam(&d.dom, id)
	if np == nil {
		return nil, fmt.Errorf("Texture id:%s not found", id)
	}

	// Checks if parameter is a Sampler2D
	sampler2D, ok := np.ParameterType.(*Sampler2D)
	if !ok {
		return nil, fmt.Errorf("Texture id:%s is not a sampler2D", id)
	}

	// Get the parameter for the Sampler2D source
	np = findNewparam(&d.dom, sampler2D.Source)
	if np == nil {
		return nil, fmt.Errorf("Sampler2D source:%s not found", id)
	}

	// Checks if parameter is a surface
	surface, ok := np.ParameterType.(*Surface)
	if !ok {
		return nil, fmt.Errorf("Sampler2D source:%s is not a Surface", id)
	}

	// Checks if surface Init is InitFrom
	initFrom, ok := surface.Init.(InitFrom)
	if !ok {
		return nil, fmt.Errorf("Surface:%s init is not InitFrom", sampler2D.Source)
	}

	// Find image
	img := findImage(&d.dom, initFrom.Uri)
	if img == nil {
		return nil, fmt.Errorf("Image:%s not found", initFrom.Uri)
	}

	// Get image init from
	imgInitFrom, ok := img.ImageSource.(InitFrom)
	if !ok {
		return nil, fmt.Errorf("Image:%s source is not InitFrom", initFrom.Uri)
	}

	// Builds image file path and try to create texture
	filepath := filepath.Join(d.dirImages, filepath.Base(imgInitFrom.Uri))
	tex, err := texture.NewTexture2DFromImage(filepath)
	if err != nil {
		return nil, err
	}
	return tex, nil
}

func (d *Decoder) newBlinnMaterial(se *Blinn) (material.IMaterial, error) {

	return nil, fmt.Errorf("Not implemented")
}

func (d *Decoder) newConstantMaterial(se *Constant) (material.IMaterial, error) {

	return nil, fmt.Errorf("Not implemented")
}

func (d *Decoder) newLambertMaterial(se *Lambert) (material.IMaterial, error) {

	return nil, fmt.Errorf("Not implemented")
}

func (d *Decoder) newPhongMaterial(se *Phong) (material.IMaterial, error) {

	// Creates material with default color
	m := material.NewPhong(&math32.Color{0.5, 0.5, 0.5})

	// If "diffuse" is Color set its value in the material
	_, ok := se.Diffuse.(*Color)
	if ok {
		color := getColor(se.Diffuse)
		m.SetColor(&color)
	} else {
		// Diffuse must be a Texture
		tex, ok := se.Diffuse.(*Texture)
		if !ok {
			return nil, fmt.Errorf("diffuse is not Color nor Texture")
		}
		// Get texture 2D
		tex2D, err := d.GetTexture2D(tex.Texture)
		if err != nil {
			return nil, err
		}
		// Add texture to this material
		m.AddTexture(tex2D)
	}

	emission := getColor(se.Emission)
	m.SetEmissiveColor(&emission)

	//ambient := getColor(se.Ambient)
	//m.SetAmbientColor(&ambient)

	specular := getColor(se.Specular)
	m.SetSpecularColor(&specular)

	shininess := getFloatOrParam(se.Shininess)
	m.SetShininess(shininess)

	//m.SetOpacity(opacity float32) {
	//m.SetWireframe(true)
	m.SetSide(material.SideDouble)

	return m, nil
}

func getColor(ci interface{}) math32.Color {

	switch c := ci.(type) {
	case *Color:
		return math32.Color{c.Data[0], c.Data[1], c.Data[2]}
	}
	return math32.Color{}
}

func getColor4(ci interface{}) math32.Color4 {

	switch c := ci.(type) {
	case Color:
		return math32.Color4{c.Data[0], c.Data[1], c.Data[2], c.Data[3]}
	}
	return math32.Color4{}
}

func getFloatOrParam(vi interface{}) float32 {

	switch v := vi.(type) {
	case *Float:
		return v.Data
	}
	return 0
}

func findMaterial(dom *Collada, id string) *Material {

	for _, m := range dom.LibraryMaterials.Material {
		if m.Id == id {
			return m
		}
	}
	return nil
}

func findEffect(dom *Collada, uri string) *Effect {

	id := strings.TrimPrefix(uri, "#")
	for _, effect := range dom.LibraryEffects.Effect {
		if effect.Id == id {
			return effect
		}
	}
	return nil
}

func findProfileCOMMON(ef *Effect) *ProfileCOMMON {

	for _, pi := range ef.Profile {
		pc, ok := pi.(*ProfileCOMMON)
		if ok {
			return pc
		}
	}
	return nil
}

func findNewparam(dom *Collada, uri string) *Newparam {

	id := strings.TrimPrefix(uri, "#")
	for _, effect := range dom.LibraryEffects.Effect {
		for _, prof := range effect.Profile {
			pc, ok := prof.(*ProfileCOMMON)
			if !ok {
				continue
			}
			for _, np := range pc.Newparam {
				if np.Sid == id {
					return np
				}
			}
		}
	}
	return nil
}

func findImage(dom *Collada, uri string) *Image {

	id := strings.TrimPrefix(uri, "#")
	for _, img := range dom.LibraryImages.Image {
		if img.Id == id {
			return img
		}
	}
	return nil
}
