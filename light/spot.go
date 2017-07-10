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
	core.Node                 // Embedded node
	color     math32.Color    // Light color
	intensity float32         // Light intensity
	direction math32.Vector3  // Direction in world coordinates
	uni       *gls.Uniform3fv // Uniform with spot light properties
}

const (
	sColor          = 0  // index of color vector in uniform (0,1,2)
	sPosition       = 1  // index of position vector in uniform (3,4,5)
	sDirection      = 2  // index of position vector in uniform (6,7,8)
	sAngularDecay   = 9  // position of scalar angular decay in uniform array
	sCutoffAngle    = 10 // position of cutoff angle in uniform array
	sLinearDecay    = 11 // position of scalar linear decay in uniform array
	sQuadraticDecay = 12 // position of scalar quadratic decay in uniform array
	spotUniSize     = 5  // uniform count of 5 float32
)

// NewSpot creates and returns a spot light with the specified color and intensity
func NewSpot(color *math32.Color, intensity float32) *Spot {

	sl := new(Spot)
	sl.Node.Init()

	sl.color = *color
	sl.intensity = intensity

	// Creates uniforms and set initial values
	sl.uni = gls.NewUniform3fv("SpotLight", spotUniSize)
	sl.SetColor(color)
	sl.SetAngularDecay(15.0)
	sl.SetCutoffAngle(45.0)
	sl.SetLinearDecay(1.0)
	sl.SetQuadraticDecay(1.0)

	return sl
}

// SetColor sets the color of this light
func (sl *Spot) SetColor(color *math32.Color) {

	sl.color = *color
	tmpColor := sl.color
	tmpColor.MultiplyScalar(sl.intensity)
	sl.uni.SetColor(sColor, &tmpColor)
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
	sl.uni.SetColor(sColor, &tmpColor)
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

	sl.uni.SetPos(sCutoffAngle, angle)
}

// CutoffAngle returns the current cutoff angle in degrees from 0 to 90
func (sl *Spot) CutoffAngle() float32 {

	return sl.uni.GetPos(sCutoffAngle)
}

// SetAngularDecay sets the angular decay exponent
func (sl *Spot) SetAngularDecay(decay float32) {

	sl.uni.SetPos(sAngularDecay, decay)
}

// AngularDecay returns the current angular decay exponent
func (sl *Spot) AngularDecay() float32 {

	return sl.uni.GetPos(sAngularDecay)
}

// SetLinearDecay sets the linear decay factor as a function of the distance
func (sl *Spot) SetLinearDecay(decay float32) {

	sl.uni.SetPos(sLinearDecay, decay)
}

// LinearDecay returns the current linear decay factor
func (sl *Spot) LinearDecay() float32 {

	return sl.uni.GetPos(sLinearDecay)
}

// SetQuadraticDecay sets the quadratic decay factor as a function of the distance
func (sl *Spot) SetQuadraticDecay(decay float32) {

	sl.uni.SetPos(sQuadraticDecay, decay)
}

// QuadraticDecay returns the current quadratic decay factor
func (sl *Spot) QuadraticDecay() float32 {

	return sl.uni.GetPos(sQuadraticDecay)
}

// RenderSetup is called by the engine before rendering the scene
func (sl *Spot) RenderSetup(gs *gls.GLS, rinfo *core.RenderInfo, idx int) {

	// Calculates and updates light position uniform in camera coordinates
	var pos math32.Vector3
	sl.WorldPosition(&pos)
	var pos4 math32.Vector4
	pos4.SetVector3(&pos, 1.0)
	pos4.ApplyMatrix4(&rinfo.ViewMatrix)
	sl.uni.SetVector3(sPosition, &math32.Vector3{pos4.X, pos4.Y, pos4.Z})

	// Calculates and updates light direction uniform in camera coordinates
	pos4.SetVector3(&sl.direction, 0.0)
	pos4.ApplyMatrix4(&rinfo.ViewMatrix)
	// Normalize here ??
	sl.uni.SetVector3(sDirection, &math32.Vector3{pos4.X, pos4.Y, pos4.Z})

	// Transfer uniform
	sl.uni.TransferIdx(gs, idx*spotUniSize)
}
