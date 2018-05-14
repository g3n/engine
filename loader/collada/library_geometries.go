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
// LibraryGeometries
//
type LibraryGeometries struct {
	Asset    *Asset
	Geometry []*Geometry
}

// Dump prints out information about the LibraryGeometries
func (lg *LibraryGeometries) Dump(out io.Writer, indent int) {

	if lg == nil {
		return
	}
	fmt.Fprintf(out, "%sLibraryGeometries:\n", sIndent(indent))
	ind := indent + step
	if lg.Asset != nil {
		lg.Asset.Dump(out, ind)
	}
	for _, g := range lg.Geometry {
		g.Dump(out, ind)
	}
}

//
// Geometry
//
type Geometry struct {
	Id               string      // Geometry id (optional)
	Name             string      // Geometry name (optional)
	GeometricElement interface{} // Geometry type object (Mesh|others)
}

// Dump prints out information about the Geometry
func (g *Geometry) Dump(out io.Writer, indent int) {

	fmt.Fprintf(out, "%sGeometry id:%s name:%s\n", sIndent(indent), g.Id, g.Name)
	ind := indent + step
	switch gt := g.GeometricElement.(type) {
	case *Mesh:
		gt.Dump(out, ind)
		break
	}

}

//
//  Mesh
//
type Mesh struct {
	Source            []*Source     // One or more sources Sources
	Vertices          Vertices      // Vertices positions
	PrimitiveElements []interface{} // Geometry primitives (polylist|others)
}

// Dump prints out information about the Mesh
func (m *Mesh) Dump(out io.Writer, indent int) {

	fmt.Fprintf(out, "%sMesh:\n", sIndent(indent))
	ind := indent + step
	for _, s := range m.Source {
		s.Dump(out, ind)
	}
	m.Vertices.Dump(out, ind)
	for _, pe := range m.PrimitiveElements {
		switch pt := pe.(type) {
		case *Lines:
			pt.Dump(out, ind)
		case *Polylist:
			pt.Dump(out, ind)
		}
	}
}

//
// Vertices
//
type Vertices struct {
	Id    string
	Name  string
	Input []Input
}

// Dump prints out information about the Vertices
func (v *Vertices) Dump(out io.Writer, indent int) {

	fmt.Fprintf(out, "%sVertices id:%s name:%s\n", sIndent(indent), v.Id, v.Name)
	for _, inp := range v.Input {
		inp.Dump(out, indent+step)
	}
}

//
// Input
//
type Input struct {
	Semantic string
	Source   string // source URL
}

// Dump prints out information about the Input
func (i *Input) Dump(out io.Writer, indent int) {

	fmt.Fprintf(out, "%sInput semantic:%s source:%s\n", sIndent(indent), i.Semantic, i.Source)
}

//
// Polylist
//
type Polylist struct {
	Name     string
	Count    int
	Material string
	Input    []InputShared
	Vcount   []int
	P        []int
}

// Dump prints out information about the Polylist
func (pl *Polylist) Dump(out io.Writer, indent int) {

	fmt.Fprintf(out, "%sPolylist name:%s count:%d material:%s\n", sIndent(indent), pl.Name, pl.Count, pl.Material)
	ind := indent + step
	for _, is := range pl.Input {
		is.Dump(out, ind)
	}
	fmt.Fprintf(out, "%sVcount(%d):%v\n", sIndent(ind), len(pl.Vcount), intsToString(pl.Vcount, 20))
	fmt.Fprintf(out, "%sP(%d):%v\n", sIndent(ind), len(pl.P), intsToString(pl.P, 20))
}

//
// InputShared
//
type InputShared struct {
	Offset   int
	Semantic string
	Source   string // source URL
	Set      int
}

//
// Triangles
//
type Triangles struct {
	Name     string
	Count    int
	Material string
	Input    []InputShared
	P        []int
}

//
// Lines
//
type Lines struct {
	Name     string
	Count    int
	Material string
	Input    []InputShared
	P        []int
}

// Dump prints out information about the Lines
func (ln *Lines) Dump(out io.Writer, indent int) {

	fmt.Fprintf(out, "%sLines name:%s count:%d material:%s\n", sIndent(indent), ln.Name, ln.Count, ln.Material)
	ind := indent + step
	for _, is := range ln.Input {
		is.Dump(out, ind)
	}
	fmt.Fprintf(out, "%sP(%d):%v\n", sIndent(ind), len(ln.P), intsToString(ln.P, 20))
}

//
// LineStrips
//
type LineStrips struct {
	Name     string
	Count    int
	Material string
	Input    []InputShared
	P        []int
}

//
// Trifans
//
type Trifans struct {
	Name     string
	Count    int
	Material string
	Input    []InputShared
	P        []int
}

//
// Tristrips
//
type Tristrips struct {
	Name     string
	Count    int
	Material string
	Input    []InputShared
	P        []int
}

// Dump prints out information about the Tristrips
func (is *InputShared) Dump(out io.Writer, indent int) {

	fmt.Fprintf(out, "%sInputShared offset:%d semantic:%s source:%s set:%d\n",
		sIndent(indent), is.Offset, is.Semantic, is.Source, is.Set)
}

// Decodes "library_geometry" children
func (d *Decoder) decLibraryGeometries(start xml.StartElement, dom *Collada) error {

	lg := new(LibraryGeometries)
	dom.LibraryGeometries = lg

	for {
		// Get next child START
		child, _, err := d.decNextChild(start)
		if err != nil || child.Name.Local == "" {
			return err
		}
		// Decode optional asset
		if child.Name.Local == "asset" {
			lg.Asset = new(Asset)
			err = d.decAsset(child, lg.Asset)
			if err != nil {
				return err
			}
			continue
		}
		//log.Debug("decLibraryGeometries: %s", child.Name.Local)
		// Decode geometry
		if child.Name.Local == "geometry" {
			err = d.decGeometry(child, lg)
			if err != nil {
				return err
			}
			continue
		}
	}
}

// decGeometry receives the start element of a geometry and
// decodes all its children and appends the decoded geometry
// to the specified slice.
func (d *Decoder) decGeometry(start xml.StartElement, lg *LibraryGeometries) error {

	// Get geometry id and name attributes
	geom := &Geometry{}
	geom.Id = findAttrib(start, "id").Value
	geom.Name = findAttrib(start, "name").Value
	lg.Geometry = append(lg.Geometry, geom)

	// Decodes geometry children
	for {
		// Get next child
		child, _, err := d.decNextChild(start)
		if err != nil || child.Name.Local == "" {
			return err
		}
		// Decode mesh
		if child.Name.Local == "mesh" {
			err = d.decMesh(child, geom)
			if err != nil {
				return err
			}
			continue
		}
	}
}

// decMesh decodes the mesh from the specified geometry
func (d *Decoder) decMesh(start xml.StartElement, geom *Geometry) error {

	// Associates this mesh to the parent geometry
	mesh := &Mesh{}
	geom.GeometricElement = mesh

	// Decodes mesh children
	for {
		// Get next child
		child, _, err := d.decNextChild(start)
		if err != nil || child.Name.Local == "" {
			return err
		}
		//log.Debug("decMesh(%s): %s", start.Name.Local, child.Name.Local)
		// Decodes source
		if child.Name.Local == "source" {
			source, err := d.decSource(child)
			if err != nil {
				return err
			}
			mesh.Source = append(mesh.Source, source)
			continue
		}
		// Decodes vertices
		if child.Name.Local == "vertices" {
			err = d.decVertices(child, mesh)
			if err != nil {
				return err
			}
			continue
		}
		// Decodes lines
		if child.Name.Local == "lines" {
			err = d.decLines(child, mesh)
			if err != nil {
				return err
			}
			continue
		}
		// Decodes polylist
		if child.Name.Local == "polylist" {
			err = d.decPolylist(child, mesh)
			if err != nil {
				return err
			}
			continue
		}
	}
}

func (d *Decoder) decVertices(start xml.StartElement, mesh *Mesh) error {

	mesh.Vertices.Id = findAttrib(start, "id").Value
	mesh.Vertices.Name = findAttrib(start, "name").Value

	for {
		// Get next child
		child, _, err := d.decNextChild(start)
		if err != nil || child.Name.Local == "" {
			return err
		}
		// input
		if child.Name.Local == "input" {
			inp, err := d.decInput(child)
			if err != nil {
				return err
			}
			mesh.Vertices.Input = append(mesh.Vertices.Input, inp)
		}
	}
}

func (d *Decoder) decInput(start xml.StartElement) (Input, error) {

	var inp Input
	inp.Semantic = findAttrib(start, "semantic").Value
	inp.Source = findAttrib(start, "source").Value
	return inp, nil
}

func (d *Decoder) decLines(start xml.StartElement, mesh *Mesh) error {

	ln := &Lines{}
	ln.Name = findAttrib(start, "name").Value
	ln.Count, _ = strconv.Atoi(findAttrib(start, "count").Value)
	ln.Material = findAttrib(start, "material").Value
	mesh.PrimitiveElements = append(mesh.PrimitiveElements, ln)

	for {
		// Get next child
		child, data, err := d.decNextChild(start)
		if err != nil || child.Name.Local == "" {
			return err
		}
		// Decode input shared
		if child.Name.Local == "input" {
			inp, err := d.decInputShared(child)
			if err != nil {
				return err
			}
			ln.Input = append(ln.Input, inp)
			continue
		}
		// Decode p (primitive)
		if child.Name.Local == "p" {
			p, err := d.decPrimitive(child, data)
			if err != nil {
				return err
			}
			ln.P = p
		}
	}
}

func (d *Decoder) decPolylist(start xml.StartElement, mesh *Mesh) error {

	pl := &Polylist{}
	pl.Name = findAttrib(start, "name").Value
	pl.Count, _ = strconv.Atoi(findAttrib(start, "count").Value)
	pl.Material = findAttrib(start, "material").Value
	mesh.PrimitiveElements = append(mesh.PrimitiveElements, pl)

	for {
		// Get next child
		child, data, err := d.decNextChild(start)
		if err != nil || child.Name.Local == "" {
			return err
		}
		// Decode input shared
		if child.Name.Local == "input" {
			inp, err := d.decInputShared(child)
			if err != nil {
				return err
			}
			pl.Input = append(pl.Input, inp)
			continue
		}
		// Decode vcount
		if child.Name.Local == "vcount" {
			vc, err := d.decVcount(child, data, pl.Count)
			if err != nil {
				return err
			}
			pl.Vcount = vc
			continue
		}
		// Decode p (primitive)
		if child.Name.Local == "p" {
			p, err := d.decPrimitive(child, data)
			if err != nil {
				return err
			}
			pl.P = p
		}
	}
}

func (d *Decoder) decInputShared(start xml.StartElement) (InputShared, error) {

	var inp InputShared
	inp.Offset, _ = strconv.Atoi(findAttrib(start, "offset").Value)
	inp.Semantic = findAttrib(start, "semantic").Value
	inp.Source = findAttrib(start, "source").Value
	inp.Set, _ = strconv.Atoi(findAttrib(start, "set").Value)
	return inp, nil
}

func (d *Decoder) decVcount(start xml.StartElement, data []byte, size int) ([]int, error) {

	vcount := make([]int, size)
	var br bytesReader
	br.Init(data)
	idx := 0
	for {
		tok := br.TokenNext()
		if tok == nil {
			break
		}
		v, err := strconv.Atoi(string(tok))
		if err != nil {
			return nil, err
		}
		vcount[idx] = v
		idx++
	}
	return vcount, nil
}

func (d *Decoder) decPrimitive(start xml.StartElement, data []byte) ([]int, error) {

	p := make([]int, 0)
	var br bytesReader
	br.Init(data)
	idx := 0
	for {
		tok := br.TokenNext()
		if tok == nil {
			break
		}
		v, err := strconv.Atoi(string(tok))
		if err != nil {
			return nil, err
		}
		p = append(p, v)
		idx++
	}
	return p, nil
}
