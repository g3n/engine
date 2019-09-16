// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphic

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/material"
)

// Points represents a geometry containing only points
type Points struct {
	Graphic             // Embedded graphic
	uniMVPm gls.Uniform // Model view projection matrix uniform location cache
}

// NewPoints creates and returns a graphic points object with the specified
// geometry and material.
func NewPoints(igeom geometry.IGeometry, imat material.IMaterial) *Points {

	p := new(Points)
	p.Graphic.Init(p, igeom, gls.POINTS)
	if imat != nil {
		p.AddMaterial(p, imat, 0, 0)
	}
	p.uniMVPm.Init("MVP")
	return p
}

// RenderSetup is called by the engine before rendering this graphic.
func (p *Points) RenderSetup(gs *gls.GLS, rinfo *core.RenderInfo) {

	// Transfer model view projection matrix uniform
	mvpm := p.ModelViewProjectionMatrix()
	location := p.uniMVPm.Location(gs)
	gs.UniformMatrix4fv(location, 1, false, &mvpm[0])
}
