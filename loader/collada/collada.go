// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package collada
package collada

import (
	"encoding/xml"
	"fmt"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/texture"
	"io"
	"os"
)

// Decoder contains all decoded data from collada file
type Decoder struct {
	xmldec     *xml.Decoder                  // xml decoder used internally
	lastToken  interface{}                   // last token read
	dom        Collada                       // Collada dom
	dirImages  string                        // Base directory for images
	geometries map[string]geomInstance       // Instanced geometries by id
	materials  map[string]material.IMaterial // Instanced materials by id
	tex2D      map[string]*texture.Texture2D // Instanced textures 2D by id
}

type geomInstance struct {
	geom  geometry.IGeometry
	ptype uint32
}

// Decode decodes the specified collada file returning a decoder object and an error.
func Decode(filepath string) (*Decoder, error) {

	// Opens file
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return DecodeReader(f)
}

// DecodeReader decodes the specified collada reader returning a decoder object and an error.
func DecodeReader(f io.Reader) (*Decoder, error) {

	d := new(Decoder)
	d.xmldec = xml.NewDecoder(f)
	d.geometries = make(map[string]geomInstance)
	d.materials = make(map[string]material.IMaterial)
	d.tex2D = make(map[string]*texture.Texture2D)

	err := d.decCollada(&d.dom)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (d *Decoder) SetDirImages(path string) {

	d.dirImages = path
}

//
// Collada DOM root
//
type Collada struct {
	Version             string
	Asset               Asset
	LibraryAnimations   *LibraryAnimations
	LibraryImages       *LibraryImages
	LibraryLights       *LibraryLights
	LibraryEffects      *LibraryEffects
	LibraryMaterials    *LibraryMaterials
	LibraryGeometries   *LibraryGeometries
	LibraryVisualScenes *LibraryVisualScenes
	Scene               *Scene
}

//
// Dump writes to the specified writer a text dump of the decoded Collada DOM
// to aid debugging.
//
func (d *Decoder) Dump(out io.Writer, indent int) {

	fmt.Fprintf(out, "%sCollada version:%s\n", sIndent(indent), d.dom.Version)
	d.dom.Asset.Dump(out, indent+step)
	d.dom.LibraryAnimations.Dump(out, indent+step)
	d.dom.LibraryImages.Dump(out, indent+step)
	d.dom.LibraryLights.Dump(out, indent+step)
	d.dom.LibraryEffects.Dump(out, indent+step)
	d.dom.LibraryMaterials.Dump(out, indent+step)
	d.dom.LibraryGeometries.Dump(out, indent+step)
	d.dom.LibraryVisualScenes.Dump(out, indent+step)
	d.dom.Scene.Dump(out, indent+step)
}

//
// Contributor
//
type Contributor struct {
	Author        string
	AuthorEmail   string
	AuthorWebsite string
	AuthoringTool string
	Comments      string
	Copyright     string
	SourceData    string
}

// Dump prints out information about the Contributor
func (c *Contributor) Dump(out io.Writer, indent int) {

	fmt.Fprintf(out, "%sContributor:\n", sIndent(indent))
	ind := indent + step
	if len(c.Author) > 0 {
		fmt.Fprintf(out, "%sAuthor:%s\n", sIndent(ind), c.Author)
	}
	if len(c.AuthorEmail) > 0 {
		fmt.Fprintf(out, "%sAuthorEmail:%s\n", sIndent(ind), c.AuthorEmail)
	}
	if len(c.AuthorWebsite) > 0 {
		fmt.Fprintf(out, "%sAuthorWebsite:%s\n", sIndent(ind), c.AuthorWebsite)
	}
	if len(c.AuthoringTool) > 0 {
		fmt.Fprintf(out, "%sAuthoringTool:%s\n", sIndent(ind), c.AuthoringTool)
	}
	if len(c.Comments) > 0 {
		fmt.Fprintf(out, "%sComments:%s\n", sIndent(ind), c.Comments)
	}
	if len(c.Copyright) > 0 {
		fmt.Fprintf(out, "%sCopyright:%s\n", sIndent(ind), c.Copyright)
	}
	if len(c.SourceData) > 0 {
		fmt.Fprintf(out, "%sSourceData:%s\n", sIndent(ind), c.SourceData)
	}
}

//
// Asset
//
type Asset struct {
	Contributor Contributor
	Created     string
	Modified    string
	UpAxis      string
}

// Dump prints out information about the Asset
func (a *Asset) Dump(out io.Writer, indent int) {

	fmt.Fprintf(out, "%sAsset:\n", sIndent(indent))
	ind := indent + step
	a.Contributor.Dump(out, ind)
	fmt.Fprintf(out, "%sCreated:%s\n", sIndent(ind), a.Created)
	fmt.Fprintf(out, "%sModified:%s\n", sIndent(ind), a.Modified)
	fmt.Fprintf(out, "%sUpAxis:%s\n", sIndent(ind), a.UpAxis)
}

//
// Scene
//
type Scene struct {
	InstanceVisualScene *InstanceVisualScene
}

// Dump prints out information about the Scene
func (s *Scene) Dump(out io.Writer, indent int) {

	fmt.Fprintf(out, "%sScene:\n", sIndent(indent))
	ind := indent + step
	s.InstanceVisualScene.Dump(out, ind)
}

//
// InstanceVisualScene
//
type InstanceVisualScene struct {
	Sid  string
	Name string
	Url  string
}

// Dump prints out information about the InstanceVisualScene
func (ivs *InstanceVisualScene) Dump(out io.Writer, indent int) {

	if ivs == nil {
		return
	}
	fmt.Fprintf(out, "%sInstanceVisualScene sid:%s name:%s url:%s\n",
		sIndent(indent), ivs.Sid, ivs.Name, ivs.Url)
}

func (d *Decoder) decCollada(dom *Collada) error {

	// Loop to read all first level elements
	var tok interface{}
	var err error
	first := true
	for {
		// Reads next token
		tok, err = d.xmldec.Token()
		if err == io.EOF {
			return nil
		}
		// If not a start element ignore and continue
		start, ok := tok.(xml.StartElement)
		if !ok {
			continue
		}
		// First element must be "COLLADA"
		if first {
			if start.Name.Local != "COLLADA" {
				return fmt.Errorf("Not a COLLADA file")
			}
			first = false
			dom.Version = findAttrib(start, "version").Value
			continue
		}

		// Decode specified start elements
		if start.Name.Local == "asset" {
			err = d.decAsset(start, &dom.Asset)
			if err != nil {
				break
			}
			continue
		}
		if start.Name.Local == "library_animations" {
			err = d.decLibraryAnimations(start, dom)
			if err != nil {
				break
			}
			continue
		}
		if start.Name.Local == "library_images" {
			err = d.decLibraryImages(start, dom)
			if err != nil {
				break
			}
			continue
		}
		if start.Name.Local == "library_effects" {
			err = d.decLibraryEffects(start, dom)
			if err != nil {
				break
			}
			continue
		}
		if start.Name.Local == "library_lights" {
			err = d.decLibraryLights(start, dom)
			if err != nil {
				break
			}
			continue
		}
		if start.Name.Local == "library_materials" {
			err = d.decLibraryMaterials(start, dom)
			if err != nil {
				break
			}
			continue
		}
		if start.Name.Local == "library_geometries" {
			err = d.decLibraryGeometries(start, dom)
			if err != nil {
				break
			}
			continue
		}
		if start.Name.Local == "library_visual_scenes" {
			err = d.decLibraryVisualScenes(start, dom)
			if err != nil {
				break
			}
			continue
		}
		if start.Name.Local == "scene" {
			err = d.decScene(start, dom)
			if err != nil {
				break
			}
			continue
		}
	}
	return err
}

func (d *Decoder) decContributor(start xml.StartElement, c *Contributor) error {

	for {
		child, data, err := d.decNextChild(start)
		if err != nil || child.Name.Local == "" {
			return err
		}
		if child.Name.Local == "author" {
			c.Author = string(data)
			continue
		}
		if child.Name.Local == "author_email" {
			c.AuthorEmail = string(data)
			continue
		}
		if child.Name.Local == "author_website" {
			c.AuthorWebsite = string(data)
			continue
		}
		if child.Name.Local == "authoring_tool" {
			c.AuthoringTool = string(data)
			continue
		}
		if child.Name.Local == "comments" {
			c.Comments = string(data)
			continue
		}
		if child.Name.Local == "copyright" {
			c.Comments = string(data)
			continue
		}
		if child.Name.Local == "source_data" {
			c.SourceData = string(data)
			continue
		}
	}
}

func (d *Decoder) decAsset(assetStart xml.StartElement, a *Asset) error {

	for {
		child, data, err := d.decNextChild(assetStart)
		if err != nil || child.Name.Local == "" {
			return err
		}
		if child.Name.Local == "contributor" {
			err := d.decContributor(child, &a.Contributor)
			if err != nil {
				return err
			}
			continue
		}
		if child.Name.Local == "created" {
			a.Created = string(data)
			continue
		}
		if child.Name.Local == "modified" {
			a.Modified = string(data)
			continue
		}
		if child.Name.Local == "up_axis" {
			a.UpAxis = string(data)
			continue
		}
	}
}

func (d *Decoder) decScene(start xml.StartElement, dom *Collada) error {

	dom.Scene = new(Scene)
	for {
		child, _, err := d.decNextChild(start)
		if err != nil || child.Name.Local == "" {
			return err
		}
		if child.Name.Local == "instance_visual_scene" {
			err := d.decInstanceVisualScene(child, dom.Scene)
			if err != nil {
				return err
			}
			continue
		}
	}
}

func (d *Decoder) decInstanceVisualScene(start xml.StartElement, s *Scene) error {

	vs := new(InstanceVisualScene)
	s.InstanceVisualScene = vs
	vs.Sid = findAttrib(start, "sid").Value
	vs.Name = findAttrib(start, "name").Value
	vs.Url = findAttrib(start, "url").Value
	return nil
}
