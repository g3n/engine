// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package material

import (
	"github.com/g3n/engine/gls"
)

type ScreenSpaceLine struct {
	Material                    // Embedded material
	thickness    *gls.Uniform1f // thickness properties
	viewPortSize *gls.Uniform2f // view port size properties
}

func NewScreenSpaceLine() *ScreenSpaceLine {

	m := new(ScreenSpaceLine)
	m.Material.Init()
	m.thickness = gls.NewUniform1f("thickness")
	m.viewPortSize = gls.NewUniform2f("viewportSize")

	m.SetShader("shaderScreenSpaceLine")
	m.SetDepthTest(false)
	m.SetSide(SideDouble)
	m.SetThickness(4.0)

	return m
}

// SetThickness sets the line thickness
// The default is 4.0
func (m *ScreenSpaceLine) SetThickness(thickness float32) {

	m.thickness.Set(thickness)
}

// RenderSetup is called by the engine before drawing the object
// which uses this material
func (m *ScreenSpaceLine) RenderSetup(gs *gls.GLS) {

	m.Material.RenderSetup(gs)

	_, _, width, height := gs.GetViewport()
	m.viewPortSize.Set(float32(width), float32(height))

	m.thickness.Transfer(gs)
	m.viewPortSize.Transfer(gs)
}
