// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package material

import (
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/math32"
)

type Standard struct {
	Material                 // Embedded material
	uni      *gls.Uniform3fv // Uniform array of 3 floats with material properties
	//ambient   *gls.Uniform3f // Ambient color uniform
	//diffuse   *gls.Uniform3f // Diffuse color uniform
	//specular  *gls.Uniform3f // Specular color uniform
	//emissive  *gls.Uniform3f // Emissive color uniform
	//shininess *gls.Uniform1f // Shininess exponent uniform
	//opacity   *gls.Uniform1f // Opacity (alpha)uniform
}

const (
	vAmbient   = 0              // index for Ambient color in uniform array
	vDiffuse   = 1              // index for Diffuse color in uniform array
	vSpecular  = 2              // index for Specular color in uniform array
	vEmissive  = 3              // index for Emissive color in uniform array
	pShininess = vEmissive * 4  // position for material shininess in uniform array
	pOpacity   = pShininess + 1 // position for material opacitiy in uniform array
	uniSize    = 5              // total count of groups 3 floats in uniform
)

// NewStandard creates and returns a pointer to a new standard material
func NewStandard(color *math32.Color) *Standard {

	ms := new(Standard)
	ms.Init("shaderStandard", color)
	return ms
}

func (ms *Standard) Init(shader string, color *math32.Color) {

	ms.Material.Init()
	ms.SetShader(shader)

	// Creates uniforms and adds to material
	ms.uni = gls.NewUniform3fv("Material", uniSize)
	//ms.emissive = gls.NewUniform3f("MatEmissiveColor")
	//ms.ambient = gls.NewUniform3f("MatAmbientColor")
	//ms.diffuse = gls.NewUniform3f("MatDiffuseColor")
	//ms.specular = gls.NewUniform3f("MatSpecularColor")
	//ms.shininess = gls.NewUniform1f("MatShininess")
	//ms.opacity = gls.NewUniform1f("MatOpacity")

	// Set initial values
	ms.uni.SetColor(vAmbient, color)
	ms.uni.SetColor(vDiffuse, color)
	ms.uni.Set(vSpecular, 0.5, 0.5, 0.5)
	ms.uni.Set(vEmissive, 0, 0, 0)
	ms.uni.SetPos(pShininess, 30.0)
	ms.uni.SetPos(pOpacity, 1.0)

	//ms.ambient.SetColor(color)
	//ms.diffuse.SetColor(color)
	//ms.emissive.Set(0, 0, 0)
	//ms.specular.Set(0.5, 0.5, 0.5)
	//ms.shininess.Set(30.0)
	//ms.opacity.Set(1.0)
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
	//	ms.diffuse.SetColor(color)
	//	ms.ambient.SetColor(color)
}

// SetEmissiveColor sets the material emissive color
// The default is {0,0,0}
func (ms *Standard) SetEmissiveColor(color *math32.Color) {

	ms.uni.SetColor(vEmissive, color)
	//	ms.emissive.SetColor(color)
}

// EmissiveColor returns the material current emissive color
func (ms *Standard) EmissiveColor() math32.Color {

	return ms.uni.GetColor(vEmissive)
	//return ms.emissive.GetColor()
}

// SetSpecularColor sets the material specular color reflectivity.
// The default is {0.5, 0.5, 0.5}
func (ms *Standard) SetSpecularColor(color *math32.Color) {

	ms.uni.SetColor(vSpecular, color)
	//ms.specular.SetColor(color)
}

// SetShininess sets the specular highlight factor. Default is 30.
func (ms *Standard) SetShininess(shininess float32) {

	//ms.shininess.Set(shininess)
	ms.uni.SetPos(pShininess, shininess)
}

// SetOpacity sets the material opacity (alpha). Default is 1.0.
func (ms *Standard) SetOpacity(opacity float32) {

	ms.uni.SetPos(pOpacity, opacity)
	//ms.opacity.Set(opacity)
}

func (ms *Standard) RenderSetup(gs *gls.GLS) {

	ms.Material.RenderSetup(gs)
	ms.uni.Transfer(gs)
	//ms.emissive.Transfer(gs)
	//ms.ambient.Transfer(gs)
	//ms.diffuse.Transfer(gs)
	//ms.specular.Transfer(gs)
	//ms.shininess.Transfer(gs)
	//ms.opacity.Transfer(gs)
}
