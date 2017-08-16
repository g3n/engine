// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphic

import (
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
)

type AxisHelper struct {
	Lines
}

func NewAxisHelper(size float32) *AxisHelper {

	axis := new(AxisHelper)

	vertexCount := 2

	// Creates geometry with three orthogonal lines
	// starting at the origin
	geom := geometry.NewGeometry()

	zero := math32.Vector3Zero
	up := math32.Vector3Up.MultiplyScalar(size)
	right := math32.Vector3Right.MultiplyScalar(size)
	back := math32.Vector3Back.MultiplyScalar(size)

	positions := math32.NewArrayF32(vertexCount*3, vertexCount*3)
	positions.AppendVector3(
		&zero, up,
		&zero, right,
		&zero, back,
	)

	colors := math32.NewArrayF32(vertexCount*3, vertexCount*3)
	colors.AppendColor(
		&math32.Red, &math32.Red,
		&math32.Green, &math32.Green,
		&math32.Blue, &math32.Blue,
	)
	geom.AddVBO(gls.NewVBO().AddAttrib("VertexPosition", 3).SetBuffer(positions))
	geom.AddVBO(gls.NewVBO().AddAttrib("VertexColor", 3).SetBuffer(colors))

	// Creates line material
	mat := material.NewScreenSpaceLine()

	// Initialize lines with the specified geometry and material
	axis.Lines.Init(geom, mat)
	return axis
}
