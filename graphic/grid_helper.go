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

// GridHelper is the visual representation of a grid
type GridHelper struct {
	Lines
}

// NewGridHelper creates and returns a pointer to a new grid help object
// with the specified size and step
func NewGridHelper(size, step float32, color *math32.Color) *GridHelper {

	grid := new(GridHelper)

	half := size / 2
	positions := math32.NewArrayF32(0, 0)
	for i := -half; i <= half; i += step {
		positions.Append(
			-half, 0, i, color.R, color.G, color.B,
			half, 0, i, color.R, color.G, color.B,
			i, 0, -half, color.R, color.G, color.B,
			i, 0, half, color.R, color.G, color.B,
		)
	}

	// Creates geometry
	geom := geometry.NewGeometry()
	geom.AddVBO(
		gls.NewVBO(positions).
			AddAttrib(gls.VertexPosition).
			AddAttrib(gls.VertexColor),
	)

	// Creates material
	mat := material.NewBasic()

	// Initialize lines with the specified geometry and material
	grid.Lines.Init(geom, mat)
	return grid
}
