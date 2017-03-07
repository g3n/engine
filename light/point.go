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
	core.Node                     // Embedded node
	color           math32.Color  // Light color
	intensity       float32       // Light intensity
	uColor          gls.Uniform3f // PointLightColor uniform
	uPosition       gls.Uniform3f // PointLightPosition uniform
	uLinearDecay    gls.Uniform1f // PointLightLinearDecay uniform
	uQuadraticDecay gls.Uniform1f // PointLightQuadraticDecay uniform
}

// NewPoint creates and returns a point light with the specified color and intensity
func NewPoint(color *math32.Color, intensity float32) *Point {

	lp := new(Point)
	lp.Node.Init()
	lp.color = *color
	lp.intensity = intensity

	// Creates uniforms
	lp.uColor.Init("PointLightColor")
	lp.uPosition.Init("PointLightPosition")
	lp.uLinearDecay.Init("PointLightLinearDecay")
	lp.uQuadraticDecay.Init("PointLightQuadraticDecay")

	// Set initial values
	lp.SetColor(color)
	lp.uPosition.Set(0, 0, 0)
	lp.uLinearDecay.Set(1.0)
	lp.uQuadraticDecay.Set(1.0)
	return lp
}

// SetColor sets the color of this light
func (lp *Point) SetColor(color *math32.Color) {

	lp.color = *color
	tmpColor := lp.color
	tmpColor.MultiplyScalar(lp.intensity)
	lp.uColor.SetColor(&tmpColor)
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
	lp.uColor.SetColor(&tmpColor)
}

// Intensity returns the current intensity of this light
func (lp *Point) Intensity() float32 {

	return lp.intensity
}

// SetLinearDecay sets the linear decay factor as a function of the distance
func (lp *Point) SetLinearDecay(decay float32) {

	lp.uLinearDecay.Set(decay)
}

// LinearDecay returns the current linear decay factor
func (lp *Point) LinearDecay() float32 {

	return lp.uLinearDecay.Get()
}

// SetQuadraticDecay sets the quadratic decay factor as a function of the distance
func (lp *Point) SetQuadraticDecay(decay float32) {

	lp.uQuadraticDecay.Set(decay)
}

// QuadraticDecay returns the current quadratic decay factor
func (lp *Point) QuadraticDecay() float32 {

	return lp.uQuadraticDecay.Get()
}

// RenderSetup is called by the engine before rendering the scene
func (lp *Point) RenderSetup(gs *gls.GLS, rinfo *core.RenderInfo, idx int) {

	// Transfer uniforms
	lp.uColor.TransferIdx(gs, idx)
	lp.uLinearDecay.TransferIdx(gs, idx)
	lp.uQuadraticDecay.TransferIdx(gs, idx)

	// Calculates and updates light position uniform in camera coordinates
	var pos math32.Vector3
	lp.WorldPosition(&pos)
	pos4 := math32.Vector4{pos.X, pos.Y, pos.Z, 1.0}
	pos4.ApplyMatrix4(&rinfo.ViewMatrix)
	lp.uPosition.SetVector3(&math32.Vector3{pos4.X, pos4.Y, pos4.Z})
	lp.uPosition.TransferIdx(gs, idx)
}
