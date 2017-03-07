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

	// Creates geometry with three orthogonal lines
	// starting at the origin
	geom := geometry.NewGeometry()
	positions := math32.NewArrayF32(0, 18)
	positions.Append(
		0, 0, 0, size, 0, 0,
		0, 0, 0, 0, size, 0,
		0, 0, 0, 0, 0, size,
	)
	colors := math32.NewArrayF32(0, 18)
	colors.Append(
		1, 0, 0, 1, 0.6, 0,
		0, 1, 0, 0.6, 1, 0,
		0, 0, 1, 0, 0.6, 1,
	)
	geom.AddVBO(gls.NewVBO().AddAttrib("VertexPosition", 3).SetBuffer(positions))
	geom.AddVBO(gls.NewVBO().AddAttrib("VertexColor", 3).SetBuffer(colors))

	// Creates line material
	mat := material.NewBasic()
	mat.SetLineWidth(2.0)

	// Initialize lines with the specified geometry and material
	axis.Lines.Init(geom, mat)
	return axis
}
