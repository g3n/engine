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

// Sprite is a potentially animated image positioned in space that always faces the camera.
type Sprite struct {
	Graphic             // Embedded graphic
	uniMVPM gls.Uniform // Model view projection matrix uniform location cache
}

// NewSprite creates and returns a pointer to a sprite with the specified dimensions and material
func NewSprite(width, height float32, imat material.IMaterial) *Sprite {

	s := new(Sprite)

	// Creates geometry
	geom := geometry.NewGeometry()
	w := width / 2
	h := height / 2

	// Builds array with vertex positions and texture coordinates
	positions := math32.NewArrayF32(0, 12)
	positions.Append(
		-w, -h, 0, 0, 0,
		w, -h, 0, 1, 0,
		w, h, 0, 1, 1,
		-w, h, 0, 0, 1,
	)
	// Builds array of indices
	indices := math32.NewArrayU32(0, 6)
	indices.Append(0, 1, 2, 0, 2, 3)

	// Set geometry buffers
	geom.SetIndices(indices)
	geom.AddVBO(
		gls.NewVBO(positions).
			AddAttrib(gls.VertexPosition).
			AddAttrib(gls.VertexTexcoord),
	)

	s.Graphic.Init(s, geom, gls.TRIANGLES)
	s.AddMaterial(s, imat, 0, 0)

	s.uniMVPM.Init("MVP")
	return s
}

// RenderSetup sets up the rendering of the sprite.
func (s *Sprite) RenderSetup(gs *gls.GLS, rinfo *core.RenderInfo) {

	// Calculates model view matrix
	mw := s.MatrixWorld()
	var mvm math32.Matrix4
	mvm.MultiplyMatrices(&rinfo.ViewMatrix, &mw)

	// Decomposes model view matrix
	var position math32.Vector3
	var quaternion math32.Quaternion
	var scale math32.Vector3
	mvm.Decompose(&position, &quaternion, &scale)

	// Removes any rotation in X and Y axes and compose new model view matrix
	rotation := s.Rotation()
	actualScale := s.Scale()
	if actualScale.X >= 0 {
		rotation.Y = 0
	} else {
		rotation.Y = math32.Pi
	}
	if actualScale.Y >= 0 {
		rotation.X = 0
	} else {
		rotation.X = math32.Pi
	}
	quaternion.SetFromEuler(&rotation)
	var mvmNew math32.Matrix4
	mvmNew.Compose(&position, &quaternion, &scale)

	// Calculates final MVP and updates uniform
	var mvpm math32.Matrix4
	mvpm.MultiplyMatrices(&rinfo.ProjMatrix, &mvmNew)
	location := s.uniMVPM.Location(gs)
	gs.UniformMatrix4fv(location, 1, false, &mvpm[0])
}
