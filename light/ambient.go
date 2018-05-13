// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package light

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/math32"
)

// Ambient represents an ambient light
type Ambient struct {
	core.Node              // Embedded node
	color     math32.Color // Light color
	intensity float32      // Light intensity
	uni       gls.Uniform  // Uniform location cache
}

// NewAmbient returns a pointer to a new ambient color with the specified
// color and intensity
func NewAmbient(color *math32.Color, intensity float32) *Ambient {

	la := new(Ambient)
	la.Node.Init()
	la.color = *color
	la.intensity = intensity
	la.uni.Init("AmbientLightColor")
	return la
}

// SetColor sets the color of this light
func (la *Ambient) SetColor(color *math32.Color) {

	la.color = *color
}

// Color returns the current color of this light
func (la *Ambient) Color() math32.Color {

	return la.color
}

// SetIntensity sets the intensity of this light
func (la *Ambient) SetIntensity(intensity float32) {

	la.intensity = intensity
}

// Intensity returns the current intensity of this light
func (la *Ambient) Intensity() float32 {

	return la.intensity
}

// RenderSetup is called by the engine before rendering the scene
func (la *Ambient) RenderSetup(gs *gls.GLS, rinfo *core.RenderInfo, idx int) {

	color := la.color
	color.MultiplyScalar(la.intensity)
	location := la.uni.LocationIdx(gs, int32(idx))
	gs.Uniform3f(location, color.R, color.G, color.B)
}
