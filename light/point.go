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

// Point is an omnidirectional light source
type Point struct {
	core.Node              // Embedded node
	color     math32.Color // Light color
	intensity float32      // Light intensity
	uni       gls.Uniform  // Uniform location cache
	udata     struct {     // Combined uniform data in 3 vec3:
		color          math32.Color   // Light color
		position       math32.Vector3 // Light position
		linearDecay    float32        // Distance linear decay factor
		quadraticDecay float32        // Distance quadratic decay factor
		dummy          float32        // Completes 3*vec3
	}
}

// NewPoint creates and returns a point light with the specified color and intensity
func NewPoint(color *math32.Color, intensity float32) *Point {

	lp := new(Point)
	lp.Node.Init()
	lp.color = *color
	lp.intensity = intensity

	// Creates uniform and sets initial values
	lp.uni.Init("PointLight")
	lp.SetColor(color)
	lp.SetIntensity(intensity)
	lp.SetLinearDecay(1.0)
	lp.SetQuadraticDecay(1.0)
	return lp
}

// SetColor sets the color of this light
func (lp *Point) SetColor(color *math32.Color) {

	lp.color = *color
	lp.udata.color = lp.color
	lp.udata.color.MultiplyScalar(lp.intensity)
}

// Color returns the current color of this light
func (lp *Point) Color() math32.Color {

	return lp.color
}

// SetIntensity sets the intensity of this  light
func (lp *Point) SetIntensity(intensity float32) {

	lp.intensity = intensity
	lp.udata.color = lp.color
	lp.udata.color.MultiplyScalar(lp.intensity)
}

// Intensity returns the current intensity of this light
func (lp *Point) Intensity() float32 {

	return lp.intensity
}

// SetLinearDecay sets the linear decay factor as a function of the distance
func (lp *Point) SetLinearDecay(decay float32) {

	lp.udata.linearDecay = decay
}

// LinearDecay returns the current linear decay factor
func (lp *Point) LinearDecay() float32 {

	return lp.udata.linearDecay
}

// SetQuadraticDecay sets the quadratic decay factor as a function of the distance
func (lp *Point) SetQuadraticDecay(decay float32) {

	lp.udata.quadraticDecay = decay
}

// QuadraticDecay returns the current quadratic decay factor
func (lp *Point) QuadraticDecay() float32 {

	return lp.udata.quadraticDecay
}

// RenderSetup is called by the engine before rendering the scene
func (lp *Point) RenderSetup(gs *gls.GLS, rinfo *core.RenderInfo, idx int) {

	// Calculates light position in camera coordinates and updates uniform
	var pos math32.Vector3
	lp.WorldPosition(&pos)
	pos4 := math32.Vector4{pos.X, pos.Y, pos.Z, 1.0}
	pos4.ApplyMatrix4(&rinfo.ViewMatrix)
	lp.udata.position.X = pos4.X
	lp.udata.position.Y = pos4.Y
	lp.udata.position.Z = pos4.Z

	// Transfer uniform data
	const vec3count = 3
	location := lp.uni.LocationIdx(gs, vec3count*int32(idx))
	gs.Uniform3fvUP(location, vec3count, unsafe.Pointer(&lp.udata))
}
