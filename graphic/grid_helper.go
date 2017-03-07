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

type GridHelper struct {
	Lines
}

// NewGridHelper creates and returns a pointer to a new grid help object
// with the specified size and step
func NewGridHelper(size, step float32, color *math32.Color) *GridHelper {

	grid := new(GridHelper)

	half_size := size / 2
	positions := math32.NewArrayF32(0, 0)
	for i := -half_size; i <= half_size; i += step {
		positions.Append(
			-half_size, 0, i, color.R, color.G, color.B,
			half_size, 0, i, color.R, color.G, color.B,
			i, 0, -half_size, color.R, color.G, color.B,
			i, 0, half_size, color.R, color.G, color.B,
		)
	}

	// Creates geometry
	geom := geometry.NewGeometry()
	geom.AddVBO(
		gls.NewVBO().
			AddAttrib("VertexPosition", 3).
			AddAttrib("VertexColor", 3).
			SetBuffer(positions),
	)

	// Creates material
	mat := material.NewBasic()
	mat.SetLineWidth(1.0)

	// Initialize lines with the specified geometry and material
	grid.Lines.Init(geom, mat)
	return grid
}
