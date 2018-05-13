// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package collada

import (
	"fmt"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"strings"
)

// NewScene returns a new collada empty scene
func (d *Decoder) NewScene() (core.INode, error) {

	sc := d.dom.Scene
	if sc == nil {
		return nil, fmt.Errorf("No Scene element found")
	}

	ivs := sc.InstanceVisualScene
	if ivs == nil {
		return nil, fmt.Errorf("No InstanceVisualScene element found")
	}

	vs := findVisualScene(&d.dom, ivs.Url)
	if vs == nil {
		return nil, fmt.Errorf("VisualScene id:%s not found", ivs.Url)
	}

	// Creates parent scene
	scene := core.NewNode()
	// Rotate scene if necessary
	if d.dom.Asset.UpAxis == "Z_UP" {
		scene.SetRotationX(-math32.Pi / 2)
	}

	// Creates each node and adds it to the scene
	for _, n := range vs.Node {
		node, err := d.newNode(n)
		if err != nil {
			return nil, err
		}
		scene.Add(node)
	}
	return scene, nil
}

func (d *Decoder) newNode(cnode *Node) (core.INode, error) {

	var node core.INode
	switch nt := cnode.Instance.(type) {
	// Empty Node
	case nil:
		node = core.NewNode()
		// Geometry
	case *InstanceGeometry:

		// Get geometry instance
		geomi, gtype, err := d.GetGeometry(nt.Url)
		if err != nil {
			return nil, err
		}

		switch gtype {
		case gls.TRIANGLES:
			mesh := graphic.NewMesh(geomi, nil)
			geom := geomi.GetGeometry()
			// Associates the material in <bind_material> with the geometry group material
			for _, im := range nt.BindMaterial.TechniqueCommon.InstanceMaterial {
				matid := strings.TrimPrefix(im.Target, "#")
				for i := 0; i < geom.GroupCount(); i++ {
					group := geom.GroupAt(i)
					if group.Matid == matid {
						mat, err := d.GetMaterial(im.Target)
						if err != nil {
							return nil, err
						}
						mesh.AddGroupMaterial(mat, i)
						break
					}
				}
			}
			node = mesh

		case gls.POINTS:
			mat := material.NewPoint(&math32.Color{})
			mat.SetSize(1000)
			node = graphic.NewPoints(geomi, mat)

		case gls.LINES:
			mat := material.NewBasic()
			node = graphic.NewLines(geomi, mat)

		case gls.LINE_STRIP:
			mat := material.NewBasic()
			node = graphic.NewLineStrip(geomi, mat)

		default:
			return nil, fmt.Errorf("primitive not supported")
		}
	default:
		return nil, fmt.Errorf("instance geometry type:%T not supported", nt)
	}

	n := node.GetNode()
	n.SetLoaderID(cnode.Id)

	// Apply transformation elements to the node
	for _, tei := range cnode.TransformationElements {
		switch te := tei.(type) {
		case *Matrix:
			// Get math32.Matrix from the matrix data and
			// transpose to a column matrix
			var m math32.Matrix4
			m.FromArray(te.Data[:], 0)
			m.Transpose()
			// Decompose the transformation matrix
			var position math32.Vector3
			var quaternion math32.Quaternion
			var scale math32.Vector3
			m.Decompose(&position, &quaternion, &scale)
			// Sets the node position, quaternion and scale
			n.SetPositionVec(&position)
			n.SetQuaternionQuat(&quaternion)
			n.SetScaleVec(&scale)
		case *Rotate:
			// Check angle of rotation
			if te.Data[3] == 0 {
				continue
			}
			var q math32.Quaternion
			axis := math32.Vector3{te.Data[0], te.Data[1], te.Data[2]}
			q.SetFromAxisAngle(&axis, math32.DegToRad(te.Data[3]))
			n.SetQuaternionQuat(&q)
		case *Scale:
			n.SetScale(te.Data[0], te.Data[1], te.Data[2])
		case *Translate:
			n.SetPosition(te.Data[0], te.Data[1], te.Data[2])
		default:
			return nil, fmt.Errorf("transformation element not supported")
		}
	}
	// Creates children nodes
	for _, child := range cnode.Node {
		c, err := d.newNode(child)
		if err != nil {
			return nil, err
		}
		n.Add(c)
	}

	return node, nil
}

func findVisualScene(dom *Collada, uri string) *VisualScene {

	id := strings.TrimPrefix(uri, "#")
	for _, vs := range dom.LibraryVisualScenes.VisualScene {
		if vs.Id == id {
			return vs
		}
	}
	return nil
}
