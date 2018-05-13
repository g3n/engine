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
// LibraryVisualScenes
//
type LibraryVisualScenes struct {
	Asset       *Asset
	VisualScene []*VisualScene
}

// Dump prints out information about the LibraryVisualScenes
func (lv *LibraryVisualScenes) Dump(out io.Writer, indent int) {

	if lv == nil {
		return
	}
	fmt.Fprintf(out, "%sLibraryVisualScenes:\n", sIndent(indent))
	ind := indent + step
	if lv.Asset != nil {
		lv.Asset.Dump(out, ind)
	}
	for _, vs := range lv.VisualScene {
		vs.Dump(out, ind)
	}
}

//
// VisualScene contains all the nodes of a visual scene
//
type VisualScene struct {
	Id   string
	Name string
	Node []*Node // Array of nodes
}

// Dump prints out information about the VisualScene
func (vs *VisualScene) Dump(out io.Writer, indent int) {

	fmt.Fprintf(out, "%sVisualScene id:%s name:%s\n", sIndent(indent), vs.Id, vs.Name)
	for _, n := range vs.Node {
		n.Dump(out, indent+step)
	}
}

//
// Node is embedded in each node instance
//
type Node struct {
	Id                     string
	Name                   string
	Sid                    string
	Type                   string
	Layer                  []string
	TransformationElements []interface{} // Node instance type (may be nil)
	Instance               interface{}
	Node                   []*Node // Array of children nodes
}

// Dump prints out information about the Node
func (n *Node) Dump(out io.Writer, indent int) {

	fmt.Fprintf(out, "%sNode id:%s name:%s sid:%s type:%s layer:%v\n",
		sIndent(indent), n.Id, n.Name, n.Sid, n.Type, n.Layer)
	// Dump transformation elements
	for _, te := range n.TransformationElements {
		switch tt := te.(type) {
		case *Matrix:
			tt.Dump(out, indent+step)
		case *Rotate:
			tt.Dump(out, indent+step)
		case *Scale:
			tt.Dump(out, indent+step)
		case *Translate:
			tt.Dump(out, indent+step)
		}
	}
	// Dump instance type
	switch it := n.Instance.(type) {
	case *InstanceGeometry:
		it.Dump(out, indent+step)
	}
	// Dump node children
	for _, n := range n.Node {
		n.Dump(out, indent+step)
	}
}

//
// Matrix
//
type Matrix struct {
	Sid  string
	Data [16]float32
}

// Dump prints out information about the Matrix
func (m *Matrix) Dump(out io.Writer, indent int) {

	fmt.Fprintf(out, "%sMatrix sid:%s data:%v\n", sIndent(indent), m.Sid, m.Data)
}

//
// Rotate
//
type Rotate struct {
	Sid  string
	Data [4]float32
}

// Dump prints out information about the Rotate
func (r *Rotate) Dump(out io.Writer, indent int) {

	fmt.Fprintf(out, "%sRotate sid:%s data:%v\n", sIndent(indent), r.Sid, r.Data)
}

//
// Translate
//
type Translate struct {
	Sid  string
	Data [3]float32
}

// Dump prints out information about the Translate
func (t *Translate) Dump(out io.Writer, indent int) {

	fmt.Fprintf(out, "%sTranslate sid:%s data:%v\n", sIndent(indent), t.Sid, t.Data)
}

//
// Scale
//
type Scale struct {
	Sid  string
	Data [3]float32
}

// Dump prints out information about the Scale
func (s *Scale) Dump(out io.Writer, indent int) {

	fmt.Fprintf(out, "%sScale sid:%s data:%v\n", sIndent(indent), s.Sid, s.Data)
}

//
// InstanceGeometry
//
type InstanceGeometry struct {
	Url          string // Geometry URL (required) references the ID of a Geometry
	Name         string // name of this element (optional)
	BindMaterial *BindMaterial
}

// Dump prints out information about the InstanceGeometry
func (ig *InstanceGeometry) Dump(out io.Writer, indent int) {

	fmt.Fprintf(out, "%sInstanceGeometry url:%s name:%s\n", sIndent(indent), ig.Url, ig.Name)
	if ig.BindMaterial != nil {
		ig.BindMaterial.Dump(out, indent+step)
	}
}

//
// BindMaterial
//
type BindMaterial struct {
	Params          []Param
	TechniqueCommon struct {
		InstanceMaterial []*InstanceMaterial
	}
}

// Dump prints out information about the BindMaterial
func (bm *BindMaterial) Dump(out io.Writer, indent int) {

	fmt.Fprintf(out, "%sBindMaterial\n", sIndent(indent))
	ind := indent + step
	fmt.Fprintf(out, "%sTechniqueCommon\n", sIndent(ind))
	ind += step
	for _, im := range bm.TechniqueCommon.InstanceMaterial {
		im.Dump(out, ind)
	}
}

//
// InstanceMaterial
//
type InstanceMaterial struct {
	Sid             string
	Name            string
	Target          string
	Symbol          string
	Bind            []Bind
	BindVertexInput []BindVertexInput
}

// Dump prints out information about the InstanceMaterial
func (im *InstanceMaterial) Dump(out io.Writer, indent int) {

	fmt.Fprintf(out, "%sInstanceMaterial sid:%s name:%s target:%s symbol:%s\n",
		sIndent(indent), im.Sid, im.Name, im.Target, im.Symbol)
	ind := indent + step
	for _, bvi := range im.BindVertexInput {
		bvi.Dump(out, ind)
	}
}

//
// Bind
//
type Bind struct {
	Semantic string
	Target   string
}

//
// BindVertexInput
//
type BindVertexInput struct {
	Semantic      string
	InputSemantic string
	InputSet      uint
}

// Dump prints out information about the BindVertexInput
func (bvi *BindVertexInput) Dump(out io.Writer, indent int) {

	fmt.Fprintf(out, "%sBindVertexInput semantic:%s InputSemantic:%s InputSet:%d\n",
		sIndent(indent), bvi.Semantic, bvi.InputSemantic, bvi.InputSet)
}

func (d *Decoder) decLibraryVisualScenes(start xml.StartElement, dom *Collada) error {

	lv := new(LibraryVisualScenes)
	dom.LibraryVisualScenes = lv
	for {
		// Get next child element
		child, _, err := d.decNextChild(start)
		if err != nil || child.Name.Local == "" {
			return err
		}
		// Decodes VisualScene
		if child.Name.Local == "visual_scene" {
			err = d.decVisualScene(child, lv)
			if err != nil {
				return err
			}
			continue
		}
	}
}

func (d *Decoder) decVisualScene(vsStart xml.StartElement, lv *LibraryVisualScenes) error {

	// Get attributes and appends new visual scene
	vs := &VisualScene{}
	vs.Id = findAttrib(vsStart, "id").Value
	vs.Name = findAttrib(vsStart, "name").Value
	vs.Node = make([]*Node, 0)
	lv.VisualScene = append(lv.VisualScene, vs)

	// Decodes visual scene children
	for {
		// Get next child element
		child, _, err := d.decNextChild(vsStart)
		if err != nil || child.Name.Local == "" {
			return err
		}
		// Decodes Node
		if child.Name.Local == "node" {
			err = d.decNode(child, &vs.Node)
			if err != nil {
				return err
			}
		}
	}
}

func (d *Decoder) decNode(nodeStart xml.StartElement, parent *[]*Node) error {

	// Get node attributes and appends the new node to its parent
	n := &Node{}
	n.Id = findAttrib(nodeStart, "id").Value
	n.Name = findAttrib(nodeStart, "name").Value
	n.Sid = findAttrib(nodeStart, "name").Value
	n.Type = findAttrib(nodeStart, "type").Value
	n.Node = make([]*Node, 0)
	*parent = append(*parent, n)

	// Decodes node children
	for {
		// Get next child element
		child, data, err := d.decNextChild(nodeStart)
		if err != nil || child.Name.Local == "" {
			return err
		}
		if child.Name.Local == "matrix" {
			err = d.decMatrix(data, n)
			if err != nil {
				return err
			}
			continue
		}
		if child.Name.Local == "rotate" {
			err = d.decRotate(data, n)
			if err != nil {
				return err
			}
			continue
		}
		if child.Name.Local == "scale" {
			err = d.decScale(data, n)
			if err != nil {
				return err
			}
			continue
		}
		if child.Name.Local == "translate" {
			err = d.decTranslate(data, n)
			if err != nil {
				return err
			}
			continue
		}
		if child.Name.Local == "instance_geometry" {
			err = d.decInstanceGeometry(child, n)
			if err != nil {
				return err
			}
			continue
		}
		// Decodes child node recursively
		if child.Name.Local == "node" {
			err = d.decNode(child, &n.Node)
			if err != nil {
				return err
			}
			continue
		}
	}
}

func (d *Decoder) decMatrix(cdata []byte, n *Node) error {

	mat := new(Matrix)
	n.TransformationElements = append(n.TransformationElements, mat)

	err := decFloat32Sequence(cdata, mat.Data[0:16])
	if err != nil {
		return err
	}
	return nil
}

func (d *Decoder) decRotate(cdata []byte, n *Node) error {

	rot := new(Rotate)
	n.TransformationElements = append(n.TransformationElements, rot)

	err := decFloat32Sequence(cdata, rot.Data[0:4])
	if err != nil {
		return err
	}
	return nil
}

func (d *Decoder) decTranslate(cdata []byte, n *Node) error {

	tr := new(Translate)
	n.TransformationElements = append(n.TransformationElements, tr)

	err := decFloat32Sequence(cdata, tr.Data[0:3])
	if err != nil {
		return err
	}
	return nil
}

func (d *Decoder) decScale(cdata []byte, n *Node) error {

	s := new(Scale)
	n.TransformationElements = append(n.TransformationElements, s)

	err := decFloat32Sequence(cdata, s.Data[0:3])
	if err != nil {
		return err
	}
	return nil
}

func (d *Decoder) decInstanceGeometry(start xml.StartElement, n *Node) error {

	// Creates new InstanceGeometry,sets its attributes and associates with node
	ig := new(InstanceGeometry)
	ig.Url = findAttrib(start, "url").Value
	ig.Name = findAttrib(start, "name").Value
	n.Instance = ig

	// Decodes instance geometry children
	for {
		// Get next child element
		child, _, err := d.decNextChild(start)
		if err != nil || child.Name.Local == "" {
			return err
		}
		// Decodes bind_material
		if child.Name.Local == "bind_material" {
			err := d.decBindMaterial(child, &ig.BindMaterial)
			if err != nil {
				return err
			}
			continue
		}
	}
}

func (d *Decoder) decBindMaterial(start xml.StartElement, dest **BindMaterial) error {

	*dest = new(BindMaterial)
	for {
		// Get next child element
		child, _, err := d.decNextChild(start)
		if err != nil || child.Name.Local == "" {
			return err
		}
		if child.Name.Local == "technique_common" {
			err := d.decBindMaterialTechniqueCommon(child, *dest)
			if err != nil {
				return err
			}
			continue
		}
	}
}

func (d *Decoder) decBindMaterialTechniqueCommon(start xml.StartElement, bm *BindMaterial) error {

	for {
		// Get next child element
		child, _, err := d.decNextChild(start)
		if err != nil || child.Name.Local == "" {
			return err
		}
		if child.Name.Local == "instance_material" {
			err := d.decInstanceMaterial(child, bm)
			if err != nil {
				return err
			}
			continue
		}
	}
}

func (d *Decoder) decInstanceMaterial(start xml.StartElement, bm *BindMaterial) error {

	im := new(InstanceMaterial)
	im.Sid = findAttrib(start, "sid").Value
	im.Name = findAttrib(start, "name").Value
	im.Target = findAttrib(start, "target").Value
	im.Symbol = findAttrib(start, "symbol").Value
	bm.TechniqueCommon.InstanceMaterial = append(bm.TechniqueCommon.InstanceMaterial, im)

	for {
		child, _, err := d.decNextChild(start)
		if err != nil || child.Name.Local == "" {
			return err
		}
		if child.Name.Local == "bind_vertex_input" {
			var bvi BindVertexInput
			bvi.Semantic = findAttrib(child, "semantic").Value
			bvi.InputSemantic = findAttrib(child, "input_semantic").Value
			v, err := strconv.Atoi(findAttrib(child, "input_set").Value)
			if err != nil {
				return err
			}
			bvi.InputSet = uint(v)
			im.BindVertexInput = append(im.BindVertexInput, bvi)
			continue
		}
	}
}
