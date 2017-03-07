// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package light

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/math32"
)

type Directional struct {
	core.Node                // Embedded node
	color      math32.Color  // Light color
	intensity  float32       // Light intensity
	uColor     gls.Uniform3f // Light color uniform (color * intensity)
	uDirection gls.Uniform3f // Light direction uniform
}

func NewDirectional(color *math32.Color, intensity float32) *Directional {

	ld := new(Directional)
	ld.Node.Init()

	ld.color = *color
	ld.intensity = intensity
	ld.uColor.Init("DirLightColor")
	ld.uDirection.Init("DirLightPosition")
	ld.SetColor(color)
	return ld
}

// SetColor sets the color of this light
func (ld *Directional) SetColor(color *math32.Color) {

	ld.color = *color
	tmpColor := ld.color
	tmpColor.MultiplyScalar(ld.intensity)
	ld.uColor.SetColor(&tmpColor)
}

// Color returns the current color of this light
func (ld *Directional) Color() math32.Color {

	return ld.color
}

// SetIntensity sets the intensity of this light
func (ld *Directional) SetIntensity(intensity float32) {

	ld.intensity = intensity
	tmpColor := ld.color
	tmpColor.MultiplyScalar(ld.intensity)
	ld.uColor.SetColor(&tmpColor)
}

// Intensity returns the current intensity of this light
func (ld *Directional) Intensity() float32 {

	return ld.intensity
}

// RenderSetup is called by the engine before rendering the scene
func (ld *Directional) RenderSetup(gs *gls.GLS, rinfo *core.RenderInfo, idx int) {

	// Sets color
	ld.uColor.TransferIdx(gs, idx)

	// Calculates and updates light direction uniform in camera coordinates
	var pos math32.Vector3
	ld.WorldPosition(&pos)
	pos4 := math32.Vector4{pos.X, pos.Y, pos.Z, 0.0}
	pos4.ApplyMatrix4(&rinfo.ViewMatrix)
	ld.uDirection.SetVector3(&math32.Vector3{pos4.X, pos4.Y, pos4.Z})
	ld.uDirection.TransferIdx(gs, idx)
}
