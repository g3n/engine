// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package light

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/math32"
)

// Directional represents a directional, positionless light
type Directional struct {
	core.Node              // Embedded node
	color     math32.Color // Light color
	intensity float32      // Light intensity
	uni       gls.Uniform  // Uniform location cache
	udata     struct {     // Combined uniform data in 2 vec3:
		color    math32.Color   // Light color
		position math32.Vector3 // Light position
	}
}

// NewDirectional creates and returns a pointer of a new directional light
// the specified color and intensity.
func NewDirectional(color *math32.Color, intensity float32) *Directional {

	ld := new(Directional)
	ld.Node.Init(ld)

	ld.color = *color
	ld.intensity = intensity
	ld.uni.Init("DirLight")
	ld.SetColor(color)
	return ld
}

// SetColor sets the color of this light
func (ld *Directional) SetColor(color *math32.Color) {

	ld.color = *color
	ld.udata.color = ld.color
	ld.udata.color.MultiplyScalar(ld.intensity)
}

// Color returns the current color of this light
func (ld *Directional) Color() math32.Color {

	return ld.color
}

// SetIntensity sets the intensity of this light
func (ld *Directional) SetIntensity(intensity float32) {

	ld.intensity = intensity
	ld.udata.color = ld.color
	ld.udata.color.MultiplyScalar(ld.intensity)
}

// Intensity returns the current intensity of this light
func (ld *Directional) Intensity() float32 {

	return ld.intensity
}

// RenderSetup is called by the engine before rendering the scene
func (ld *Directional) RenderSetup(gs *gls.GLS, rinfo *core.RenderInfo, idx int) {

	// Calculates light position in camera coordinates and updates uniform
	var pos math32.Vector3
	ld.WorldPosition(&pos)
	pos4 := math32.Vector4{pos.X, pos.Y, pos.Z, 0.0}
	pos4.ApplyMatrix4(&rinfo.ViewMatrix)
	ld.udata.position.X = pos4.X
	ld.udata.position.Y = pos4.Y
	ld.udata.position.Z = pos4.Z

	// Transfer uniform data
	const vec3count = 2
	location := ld.uni.LocationIdx(gs, vec3count*int32(idx))
	gs.Uniform3fv(location, vec3count, &ld.udata.color.R)
}
