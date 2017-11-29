// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package material

import (
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/math32"
	"unsafe"
)

// Standard material supports the classic lighting model with
// ambient, diffuse, specular and emissive lights.
// The lighting calculation is implemented in the vertex shader.
type Standard struct {
	Material             // Embedded material
	uni      gls.Uniform // Uniform location cache
	udata    struct {    // Combined uniform data in 6 vec3:
		ambient    math32.Color // Ambient color reflectivity
		diffuse    math32.Color // Diffuse color reflectivity
		specular   math32.Color // Specular color reflectivity
		emissive   math32.Color // Emissive color
		shininess  float32      // Specular shininess factor
		opacity    float32      // Opacity
		psize      float32      // Point size
		protationZ float32      // Point rotation around Z axis
	}
}

// Number of glsl shader vec3 elements used by uniform data
const standardVec3Count = 6

// NewStandard creates and returns a pointer to a new standard material
func NewStandard(color *math32.Color) *Standard {

	ms := new(Standard)
	ms.Init("standard", color)
	return ms
}

// Init initializes the material setting the specified shader and color
// It is used mainly when the material is embedded in another type
func (ms *Standard) Init(shader string, color *math32.Color) {

	ms.Material.Init()
	ms.SetShader(shader)

	// Creates uniforms and set initial values
	ms.uni.Init("Material")
	ms.SetColor(color)
	ms.SetSpecularColor(&math32.Color{0.5, 0.5, 0.5})
	ms.SetEmissiveColor(&math32.Color{0, 0, 0})
	ms.SetShininess(30.0)
	ms.SetOpacity(1.0)
}

// AmbientColor returns the material ambient color reflectivity.
func (ms *Standard) AmbientColor() math32.Color {

	return ms.udata.ambient
}

// SetAmbientColor sets the material ambient color reflectivity.
// The default is the same as the diffuse color
func (ms *Standard) SetAmbientColor(color *math32.Color) {

	ms.udata.ambient = *color
}

// SetColor sets the material diffuse color and also the
// material ambient color reflectivity
func (ms *Standard) SetColor(color *math32.Color) {

	ms.udata.diffuse = *color
	ms.udata.ambient = *color
}

// SetEmissiveColor sets the material emissive color
// The default is {0,0,0}
func (ms *Standard) SetEmissiveColor(color *math32.Color) {

	ms.udata.emissive = *color
}

// EmissiveColor returns the material current emissive color
func (ms *Standard) EmissiveColor() math32.Color {

	return ms.udata.emissive
}

// SetSpecularColor sets the material specular color reflectivity.
// The default is {0.5, 0.5, 0.5}
func (ms *Standard) SetSpecularColor(color *math32.Color) {

	ms.udata.specular = *color
}

// SetShininess sets the specular highlight factor. Default is 30.
func (ms *Standard) SetShininess(shininess float32) {

	ms.udata.shininess = shininess
}

// SetOpacity sets the material opacity (alpha). Default is 1.0.
func (ms *Standard) SetOpacity(opacity float32) {

	ms.udata.opacity = opacity
}

// RenderSetup is called by the engine before drawing the object
// which uses this material
func (ms *Standard) RenderSetup(gs *gls.GLS) {

	ms.Material.RenderSetup(gs)
	location := ms.uni.Location(gs)
	gs.Uniform3fvUP(location, standardVec3Count, unsafe.Pointer(&ms.udata))
}
