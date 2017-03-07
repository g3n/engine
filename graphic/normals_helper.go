// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphic

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
)

type NormalsHelper struct {
	Lines
	size   float32
	target *core.Node
	tgeom  *geometry.Geometry
}

// NewNormalsHelper creates, initializes and returns a pointer to Normals helper object.
// This helper shows the normal vectors of the specified object.
func NewNormalsHelper(ig IGraphic, size float32, color *math32.Color, lineWidth float32) *NormalsHelper {

	// Creates new Normals helper
	nh := new(NormalsHelper)
	nh.size = size

	// Saves the object to show the normals
	nh.target = ig.GetNode()

	// Get the geometry of the target object
	nh.tgeom = ig.GetGeometry()

	// Get the number of target vertex positions
	vertices := nh.tgeom.VBO("VertexPosition")
	n := vertices.Buffer().Size() * 2

	// Creates this helper geometry
	geom := geometry.NewGeometry()
	positions := math32.NewArrayF32(n, n)
	geom.AddVBO(gls.NewVBO().AddAttrib("VertexPosition", 3).SetBuffer(positions))

	// Creates this helper material
	mat := material.NewStandard(color)
	mat.SetLineWidth(lineWidth)

	// Initialize graphic
	nh.Lines.Init(geom, mat)

	nh.Update()
	return nh
}

// Update should be called in the render loop to update the normals from the
// target object
func (nh *NormalsHelper) Update() {

	var v1 math32.Vector3
	var v2 math32.Vector3
	var normalMatrix math32.Matrix3

	// Updates the target object matrix and get its normal matrix
	matrixWorld := nh.target.MatrixWorld()
	normalMatrix.GetNormalMatrix(&matrixWorld)

	// Get the target positions and normals buffers
	tposvbo := nh.tgeom.VBO("VertexPosition")
	tpositions := tposvbo.Buffer()
	tnormvbo := nh.tgeom.VBO("VertexNormal")
	tnormals := tnormvbo.Buffer()

	// Get this object positions buffer
	geom := nh.GetGeometry()
	posvbo := geom.VBO("VertexPosition")
	positions := posvbo.Buffer()

	// For each target object vertex position:
	for pos := 0; pos < tpositions.Size(); pos += 3 {
		// Get the target vertex position and apply the current world matrix transform
		// to get the base for this normal line segment.
		tpositions.GetVector3(pos, &v1)
		v1.ApplyMatrix4(&matrixWorld)

		// Calculates the end position of the normal line segment
		tnormals.GetVector3(pos, &v2)
		v2.ApplyMatrix3(&normalMatrix).Normalize().MultiplyScalar(nh.size).Add(&v1)

		// Sets the line segment representing the normal of the current target position
		// at this helper VBO
		positions.SetVector3(2*pos, &v1)
		positions.SetVector3(2*pos+3, &v2)
	}
	posvbo.Update()
}
