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

type LineStrip struct {
	Graphic             // Embedded graphic object
	uniMVP  gls.Uniform // Model view projection matrix uniform location cache
}

// NewLineStrip creates and returns a pointer to a new LineStrip graphic
// with the specified geometry and material
func NewLineStrip(igeom geometry.IGeometry, imat material.IMaterial) *LineStrip {

	l := new(LineStrip)
	l.Graphic.Init(igeom, gls.LINE_STRIP)
	l.AddMaterial(l, imat, 0, 0)
	l.uniMVP.Init("MVP")
	return l
}

// RenderSetup is called by the engine before drawing this geometry
func (l *LineStrip) RenderSetup(gs *gls.GLS, rinfo *core.RenderInfo) {

	// Calculates model view projection matrix and updates uniform
	mw := l.MatrixWorld()
	var mvpm math32.Matrix4
	mvpm.MultiplyMatrices(&rinfo.ViewMatrix, &mw)
	mvpm.MultiplyMatrices(&rinfo.ProjMatrix, &mvpm)

	// Transfer mvpm uniform
	location := l.uniMVP.Location(gs)
	gs.UniformMatrix4fv(location, 1, false, &mvpm[0])
}

// Raycast satisfies the INode interface and checks the intersections
// of this geometry with the specified raycaster
func (l *LineStrip) Raycast(rc *core.Raycaster, intersects *[]core.Intersect) {

	lineRaycast(l, rc, intersects, 1)
}
