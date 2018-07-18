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

// AxisHelper is the visual representation of the three axes
type AxisHelper struct {
	Lines
}

// NewAxisHelper returns a pointer to a new AxisHelper object
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
	geom.AddVBO(gls.NewVBO(positions).AddAttrib(gls.VertexPosition))
	geom.AddVBO(gls.NewVBO(colors).AddAttrib(gls.VertexColor))

	// Creates line material
	mat := material.NewBasic()

	// Initialize lines with the specified geometry and material
	axis.Lines.Init(geom, mat)
	return axis
}
