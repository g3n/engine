// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package light

import (
	"unsafe"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/math32"
)

// Spot represents a spotlight
type Spot struct {
	core.Node              // Embedded node
	color     math32.Color // Light color
	intensity float32      // Light intensity
	uni       gls.Uniform  // Uniform location cache
	udata     struct {     // Combined uniform data in 5 vec3:
		color          math32.Color   // Light color
		position       math32.Vector3 // Light position
		direction      math32.Vector3 // Light direction
		angularDecay   float32        // Angular decay factor
		cutoffAngle    float32        // Cut off angle
		linearDecay    float32        // Distance linear decay
		quadraticDecay float32        // Distance quadratic decay
		dummy1         float32        // Completes 5*vec3
		dummy2         float32        // Completes 5*vec3
	}
}

// NewSpot creates and returns a spot light with the specified color and intensity
func NewSpot(color *math32.Color, intensity float32) *Spot {

	l := new(Spot)
	l.Node.Init()
	l.color = *color
	l.intensity = intensity
	l.uni.Init("SpotLight")
	l.SetColor(color)
	l.SetAngularDecay(15.0)
	l.SetCutoffAngle(45.0)
	l.SetLinearDecay(1.0)
	l.SetQuadraticDecay(1.0)
	return l
}

// SetColor sets the color of this light
func (l *Spot) SetColor(color *math32.Color) {

	l.color = *color
	l.udata.color = l.color
	l.udata.color.MultiplyScalar(l.intensity)
}

// Color returns the current color of this light
func (l *Spot) Color() math32.Color {

	return l.color
}

// SetIntensity sets the intensity of this light
func (l *Spot) SetIntensity(intensity float32) {

	l.intensity = intensity
	l.udata.color = l.color
	l.udata.color.MultiplyScalar(l.intensity)
}

// Intensity returns the current intensity of this light
func (l *Spot) Intensity() float32 {

	return l.intensity
}

// SetCutoffAngle sets the cutoff angle in degrees from 0 to 90
func (l *Spot) SetCutoffAngle(angle float32) {

	l.udata.cutoffAngle = angle
}

// CutoffAngle returns the current cutoff angle in degrees from 0 to 90
func (l *Spot) CutoffAngle() float32 {

	return l.udata.cutoffAngle
}

// SetAngularDecay sets the angular decay exponent
func (l *Spot) SetAngularDecay(decay float32) {

	l.udata.angularDecay = decay
}

// AngularDecay returns the current angular decay exponent
func (l *Spot) AngularDecay() float32 {

	return l.udata.angularDecay
}

// SetLinearDecay sets the linear decay factor as a function of the distance
func (l *Spot) SetLinearDecay(decay float32) {

	l.udata.linearDecay = decay
}

// LinearDecay returns the current linear decay factor
func (l *Spot) LinearDecay() float32 {

	return l.udata.linearDecay
}

// SetQuadraticDecay sets the quadratic decay factor as a function of the distance
func (l *Spot) SetQuadraticDecay(decay float32) {

	l.udata.quadraticDecay = decay
}

// QuadraticDecay returns the current quadratic decay factor
func (l *Spot) QuadraticDecay() float32 {

	return l.udata.quadraticDecay
}

// RenderSetup is called by the engine before rendering the scene
func (l *Spot) RenderSetup(gs *gls.GLS, rinfo *core.RenderInfo, idx int) {

	// Calculates and updates light position uniform in camera coordinates
	var pos math32.Vector3
	l.WorldPosition(&pos)
	var pos4 math32.Vector4
	pos4.SetVector3(&pos, 1.0)
	pos4.ApplyMatrix4(&rinfo.ViewMatrix)
	l.udata.position.X = pos4.X
	l.udata.position.Y = pos4.Y
	l.udata.position.Z = pos4.Z

	// Calculates and updates light direction uniform in camera coordinates
	var dir math32.Vector3
	l.WorldDirection(&dir)
	pos4.SetVector3(&dir, 0.0)
	pos4.ApplyMatrix4(&rinfo.ViewMatrix)
	l.udata.direction.X = pos4.X
	l.udata.direction.Y = pos4.Y
	l.udata.direction.Z = pos4.Z

	// Transfer uniform data
	const vec3count = 5
	location := l.uni.LocationIdx(gs, vec3count*int32(idx))
	gs.Uniform3fvUP(location, vec3count, unsafe.Pointer(&l.udata))
}
