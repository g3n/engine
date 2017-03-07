// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package material

import (
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/math32"
)

type Point struct {
	Material                // Embedded base material
	emissive  gls.Uniform3f // point emissive uniform
	size      gls.Uniform1f // point size uniform
	opacity   gls.Uniform1f // point opacity uniform
	rotationZ gls.Uniform1f // point z rotation
}

// NewPoint creates and returns a pointer to a new point material
func NewPoint(color *math32.Color) *Point {

	pm := new(Point)
	pm.Material.Init()
	pm.SetShader("shaderPoint")

	// Creates color uniform
	pm.emissive.Init("MatEmissiveColor")
	pm.emissive.SetColor(color)

	// Creates point size uniform
	pm.size.Init("PointSize")
	pm.size.Set(1.0)

	// Creates point opacity uniform
	pm.opacity.Init("MatOpacity")
	pm.opacity.Set(1.0)

	// Creates point rotation Z uniform
	pm.rotationZ.Init("RotationZ")
	pm.rotationZ.Set(0)

	return pm
}

// SetEmissiveColor sets the material emissive color
// The default is {0,0,0}
func (pm *Point) SetEmissiveColor(color *math32.Color) {

	pm.emissive.SetColor(color)
}

// EmissiveColor returns the material current emissive color
func (pm *Point) EmissiveColor() math32.Color {

	return pm.emissive.GetColor()
}

func (pm *Point) SetSize(size float32) {

	pm.size.Set(size)
}

func (pm *Point) SetOpacity(opacity float32) {

	pm.opacity.Set(opacity)
}

func (pm *Point) SetRotationZ(rot float32) {

	pm.rotationZ.Set(rot)
}

func (pm *Point) RenderSetup(gs *gls.GLS) {

	pm.Material.RenderSetup(gs)

	pm.emissive.Transfer(gs)
	pm.size.Transfer(gs)
	pm.opacity.Transfer(gs)
	pm.rotationZ.Transfer(gs)
}
