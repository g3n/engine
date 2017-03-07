// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package light

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/math32"
)

type Spot struct {
	core.Node                      // Embedded node
	color           math32.Color   // Light color
	intensity       float32        // Light intensity
	direction       math32.Vector3 // Direction in world coordinates
	uColor          gls.Uniform3f  // Uniform for light color
	uPosition       gls.Uniform3f  // Uniform for position in camera coordinates
	uDirection      gls.Uniform3f  // Uniform for direction in camera coordinates
	uAngularDecay   gls.Uniform1f  // Uniform for angular attenuation exponent
	uCutoffAngle    gls.Uniform1f  // Uniform for cutoff angle from 0 to 90 degrees
	uLinearDecay    gls.Uniform1f  // Uniform for linear distance decay
	uQuadraticDecay gls.Uniform1f  // Uniform for quadratic distance decay
}

// NewSpot creates and returns a spot light with the specified color and intensity
func NewSpot(color *math32.Color, intensity float32) *Spot {

	sp := new(Spot)
	sp.Node.Init()

	sp.color = *color
	sp.intensity = intensity

	// Creates uniforms
	sp.uColor.Init("SpotLightColor")
	sp.uPosition.Init("SpotLightPosition")
	sp.uDirection.Init("SpotLightDirection")
	sp.uAngularDecay.Init("SpotLightAngularDecay")
	sp.uCutoffAngle.Init("SpotLightCutoffAngle")
	sp.uLinearDecay.Init("SpotLightLinearDecay")
	sp.uQuadraticDecay.Init("SpotQuadraticDecay")

	// Set initial values
	sp.intensity = intensity
	sp.SetColor(color)
	sp.uAngularDecay.Set(15.0)
	sp.uCutoffAngle.Set(45.0)
	sp.uLinearDecay.Set(1.0)
	sp.uQuadraticDecay.Set(1.0)
	return sp
}

// SetColor sets the color of this light
func (sl *Spot) SetColor(color *math32.Color) {

	sl.color = *color
	tmpColor := sl.color
	tmpColor.MultiplyScalar(sl.intensity)
	sl.uColor.SetColor(&tmpColor)
}

// Color returns the current color of this light
func (sl *Spot) Color() math32.Color {

	return sl.color
}

// SetIntensity sets the intensity of this light
func (sl *Spot) SetIntensity(intensity float32) {

	sl.intensity = intensity
	tmpColor := sl.color
	tmpColor.MultiplyScalar(sl.intensity)
	sl.uColor.SetColor(&tmpColor)
}

// Intensity returns the current intensity of this light
func (sl *Spot) Intensity() float32 {

	return sl.intensity
}

// SetDirection sets the direction of the spot light in world coordinates
func (sp *Spot) SetDirection(direction *math32.Vector3) {

	sp.direction = *direction
}

// Direction returns the current direction of this spot light in world coordinates
func (sp *Spot) Direction(direction *math32.Vector3) math32.Vector3 {

	return sp.direction
}

// SetCutoffAngle sets the cutoff angle in degrees from 0 to 90
func (sl *Spot) SetCutoffAngle(angle float32) {

	sl.uCutoffAngle.Set(angle)
}

// CutoffAngle returns the current cutoff angle in degrees from 0 to 90
func (sl *Spot) CutoffAngle() float32 {

	return sl.uCutoffAngle.Get()
}

// SetAngularDecay sets the angular decay exponent
func (sl *Spot) SetAngularDecay(decay float32) {

	sl.uAngularDecay.Set(decay)
}

// AngularDecay returns the current angular decay exponent
func (sl *Spot) AngularDecay() float32 {

	return sl.uAngularDecay.Get()
}

// SetLinearDecay sets the linear decay factor as a function of the distance
func (sl *Spot) SetLinearDecay(decay float32) {

	sl.uLinearDecay.Set(decay)
}

// LinearDecay returns the current linear decay factor
func (sl *Spot) LinearDecay() float32 {

	return sl.uLinearDecay.Get()
}

// SetQuadraticDecay sets the quadratic decay factor as a function of the distance
func (sl *Spot) SetQuadraticDecay(decay float32) {

	sl.uQuadraticDecay.Set(decay)
}

// QuadraticDecay returns the current quadratic decay factor
func (sl *Spot) QuadraticDecay() float32 {

	return sl.uQuadraticDecay.Get()
}

// RenderSetup is called by the engine before rendering the scene
func (sl *Spot) RenderSetup(gs *gls.GLS, rinfo *core.RenderInfo, idx int) {

	sl.uColor.TransferIdx(gs, idx)
	sl.uAngularDecay.TransferIdx(gs, idx)
	sl.uCutoffAngle.TransferIdx(gs, idx)
	sl.uLinearDecay.TransferIdx(gs, idx)
	sl.uQuadraticDecay.TransferIdx(gs, idx)

	// Calculates and updates light position uniform in camera coordinates
	var pos math32.Vector3
	sl.WorldPosition(&pos)
	var pos4 math32.Vector4
	pos4.SetVector3(&pos, 1.0)
	pos4.ApplyMatrix4(&rinfo.ViewMatrix)
	sl.uPosition.SetVector3(&math32.Vector3{pos4.X, pos4.Y, pos4.Z})
	sl.uPosition.TransferIdx(gs, idx)

	// Calculates and updates light direction uniform in camera coordinates
	pos4.SetVector3(&sl.direction, 0.0)
	pos4.ApplyMatrix4(&rinfo.ViewMatrix)
	// Normalize here ??
	sl.uDirection.SetVector3(&math32.Vector3{pos4.X, pos4.Y, pos4.Z})
	sl.uDirection.TransferIdx(gs, idx)
}
