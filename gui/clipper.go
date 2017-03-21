// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
)

// Clipper is a 2D graphic which optionally clips its children inside its boundary
type Clipper struct {
	graphic.Graphic                     // Embedded graphic
	root            *Root               // pointer to root container
	width           float32             // external width in pixels
	height          float32             // external height in pixels
	mat             *material.Material  // panel material
	modelMatrixUni  gls.UniformMatrix4f // pointer to model matrix uniform
	pospix          math32.Vector3      // absolute position in pixels
	xmin            float32             // minimum absolute x this panel can use
	xmax            float32             // maximum absolute x this panel can use
	ymin            float32             // minimum absolute y this panel can use
	ymax            float32             // maximum absolute y this panel can use
	bounded         bool                // panel is bounded by its parent
	enabled         bool                // enable event processing
	cursorEnter     bool                // mouse enter dispatched
	layout          ILayout             // current layout for children
	layoutParams    interface{}         // current layout parameters used by container panel
}

// NewClipper creates and returns a pointer to a new clipper with the
// specified dimensions in pixels
func NewClipper(width, height float32, geom *geometry.Geometry, mat *material.Material, mode uint32) *Clipper {

	c := new(Clipper)
	c.Initialize(width, height, geom, mat, mode)
	return c
}

// Initialize initializes this panel with a different geometry, material and OpenGL primitive
func (c *Clipper) Initialize(width, height float32, geom *geometry.Geometry, mat *material.Material, mode uint32) {

	c.width = width
	c.height = height

	// Initialize graphic
	c.Graphic.Init(geom, mode)
	c.AddMaterial(c, mat, 0, 0)

	// Creates and adds uniform
	c.modelMatrixUni.Init("ModelMatrix")

	// Set defaults
	c.bounded = true
	c.enabled = true
	//c.resize(width, height)
}

// RenderSetup is called by the Engine before drawing the object
func (c *Clipper) RenderSetup(gl *gls.GLS, rinfo *core.RenderInfo) {

	// Sets model matrix
	var mm math32.Matrix4
	c.SetModelMatrix(gl, &mm)
	c.modelMatrixUni.SetMatrix4(&mm)

	// Transfer model matrix uniform
	c.modelMatrixUni.Transfer(gl)
}

// SetModelMatrix calculates and sets the specified matrix with the model matrix for this panel
func (c *Clipper) SetModelMatrix(gl *gls.GLS, mm *math32.Matrix4) {

	// Get the current viewport width and height
	_, _, width, height := gl.GetViewport()
	fwidth := float32(width)
	fheight := float32(height)

	// Scale the quad for the viewport so it has fixed dimensions in pixels.
	fw := float32(c.width) / fwidth
	fh := float32(c.height) / fheight
	var scale math32.Vector3
	scale.Set(2*fw, 2*fh, 1)

	// Convert absolute position in pixel coordinates from the top/left to
	// standard OpenGL clip coordinates of the quad center
	var posclip math32.Vector3
	posclip.X = (c.pospix.X - fwidth/2) / (fwidth / 2)
	posclip.Y = -(c.pospix.Y - fheight/2) / (fheight / 2)
	posclip.Z = c.Position().Z

	// Calculates the model matrix
	var quat math32.Quaternion
	quat.SetIdentity()
	mm.Compose(&posclip, &quat, &scale)
}
