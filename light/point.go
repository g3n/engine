// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package light

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/math32"
)

type Point struct {
	core.Node                 // Embedded node
	color     math32.Color    // Light color
	intensity float32         // Light intensity
	uni       *gls.Uniform3fv // Uniform with light properties
}

const (
	pColor          = 0                // index of color vector in uniform (0,1,2)
	pPosition       = 1                // index of position vector in uniform (3,4,5)
	pLinearDecay    = 6                // position of scalar linear decay in uniform array
	pQuadraticDecay = pLinearDecay + 1 // position of scalar linear decay in uniform array
	pointUniSize    = 3                // uniform count of 3 float32
)

// NewPoint creates and returns a point light with the specified color and intensity
func NewPoint(color *math32.Color, intensity float32) *Point {

	lp := new(Point)
	lp.Node.Init()
	lp.color = *color
	lp.intensity = intensity

	// Creates uniform and sets initial values
	lp.uni = gls.NewUniform3fv("PointLight", pointUniSize)
	lp.SetColor(color)
	lp.SetIntensity(intensity)
	lp.SetLinearDecay(1.0)
	lp.SetQuadraticDecay(1.0)

	return lp
}

// SetColor sets the color of this light
func (lp *Point) SetColor(color *math32.Color) {

	lp.color = *color
	tmpColor := lp.color
	tmpColor.MultiplyScalar(lp.intensity)
	lp.uni.SetColor(pColor, &tmpColor)
}

// Color returns the current color of this light
func (lp *Point) Color() math32.Color {

	return lp.color
}

// SetIntensity sets the intensity of this  light
func (lp *Point) SetIntensity(intensity float32) {

	lp.intensity = intensity
	tmpColor := lp.color
	tmpColor.MultiplyScalar(lp.intensity)
	lp.uni.SetColor(pColor, &tmpColor)
}

// Intensity returns the current intensity of this light
func (lp *Point) Intensity() float32 {

	return lp.intensity
}

// SetLinearDecay sets the linear decay factor as a function of the distance
func (lp *Point) SetLinearDecay(decay float32) {

	lp.uni.SetPos(pLinearDecay, decay)
}

// LinearDecay returns the current linear decay factor
func (lp *Point) LinearDecay() float32 {

	return lp.uni.GetPos(pLinearDecay)
}

// SetQuadraticDecay sets the quadratic decay factor as a function of the distance
func (lp *Point) SetQuadraticDecay(decay float32) {

	lp.uni.SetPos(pQuadraticDecay, decay)
}

// QuadraticDecay returns the current quadratic decay factor
func (lp *Point) QuadraticDecay() float32 {

	return lp.uni.GetPos(pQuadraticDecay)
}

// RenderSetup is called by the engine before rendering the scene
func (lp *Point) RenderSetup(gs *gls.GLS, rinfo *core.RenderInfo, idx int) {

	// Calculates and updates light position uniform in camera coordinates
	var pos math32.Vector3
	lp.WorldPosition(&pos)
	pos4 := math32.Vector4{pos.X, pos.Y, pos.Z, 1.0}
	pos4.ApplyMatrix4(&rinfo.ViewMatrix)
	lp.uni.SetVector3(pPosition, &math32.Vector3{pos4.X, pos4.Y, pos4.Z})

	// Transfer uniform
	lp.uni.TransferIdx(gs, idx*pointUniSize)
}
