// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package material

import (
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/math32"
)

// Standard material supports the classic lighting model with
// ambient, diffuse, specular and emissive lights.
// The lighting calculation is implemented in the vertex shader.
type Standard struct {
	Material                 // Embedded material
	uni      *gls.Uniform3fv // Uniform array of 3 floats with material properties
}

const (
	vAmbient   = 0              // index for Ambient color in uniform array
	vDiffuse   = 1              // index for Diffuse color in uniform array
	vSpecular  = 2              // index for Specular color in uniform array
	vEmissive  = 3              // index for Emissive color in uniform array
	pShininess = vEmissive * 4  // position for material shininess in uniform array
	pOpacity   = pShininess + 1 // position for material opacity in uniform array
	pSize      = pOpacity + 1   // position for material point size
	pRotationZ = pSize + 1      // position for material point rotation
	uniSize    = 6              // total count of groups 3 floats in uniform
)

// NewStandard creates and returns a pointer to a new standard material
func NewStandard(color *math32.Color) *Standard {

	ms := new(Standard)
	ms.Init("shaderStandard", color)
	return ms
}

// Init initializes the material setting the specified shader and color
// It is used mainly when the material is embedded in another type
func (ms *Standard) Init(shader string, color *math32.Color) {

	ms.Material.Init()
	ms.SetShader(shader)

	// Creates uniforms and set initial values
	ms.uni = gls.NewUniform3fv("Material", uniSize)
	ms.uni.SetColor(vAmbient, color)
	ms.uni.SetColor(vDiffuse, color)
	ms.uni.Set(vSpecular, 0.5, 0.5, 0.5)
	ms.uni.Set(vEmissive, 0, 0, 0)
	ms.uni.SetPos(pShininess, 30.0)
	ms.uni.SetPos(pOpacity, 1.0)
}

// AmbientColor returns the material ambient color reflectivity.
func (ms *Standard) AmbientColor() math32.Color {

	return ms.uni.GetColor(vAmbient)
}

// SetAmbientColor sets the material ambient color reflectivity.
// The default is the same as the diffuse color
func (ms *Standard) SetAmbientColor(color *math32.Color) {

	ms.uni.SetColor(vAmbient, color)
}

// SetColor sets the material diffuse color and also the
// material ambient color reflectivity
func (ms *Standard) SetColor(color *math32.Color) {

	ms.uni.SetColor(vDiffuse, color)
	ms.uni.SetColor(vAmbient, color)
}

// SetEmissiveColor sets the material emissive color
// The default is {0,0,0}
func (ms *Standard) SetEmissiveColor(color *math32.Color) {

	ms.uni.SetColor(vEmissive, color)
}

// EmissiveColor returns the material current emissive color
func (ms *Standard) EmissiveColor() math32.Color {

	return ms.uni.GetColor(vEmissive)
}

// SetSpecularColor sets the material specular color reflectivity.
// The default is {0.5, 0.5, 0.5}
func (ms *Standard) SetSpecularColor(color *math32.Color) {

	ms.uni.SetColor(vSpecular, color)
}

// SetShininess sets the specular highlight factor. Default is 30.
func (ms *Standard) SetShininess(shininess float32) {

	ms.uni.SetPos(pShininess, shininess)
}

// SetOpacity sets the material opacity (alpha). Default is 1.0.
func (ms *Standard) SetOpacity(opacity float32) {

	ms.uni.SetPos(pOpacity, opacity)
}

// RenderSetup is called by the engine before drawing the object
// which uses this material
func (ms *Standard) RenderSetup(gs *gls.GLS) {

	ms.Material.RenderSetup(gs)
	ms.uni.Transfer(gs)
}
