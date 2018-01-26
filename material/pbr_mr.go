// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package material

import (
	"unsafe"

	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/texture"
)

// PbrMr is a physically based rendered material which uses the metallic-roughness model.
type PbrMr struct {
	Material                                // Embedded material
	baseColorTex         *texture.Texture2D // Optional base color texture
	metallicRoughnessTex *texture.Texture2D // Optional metallic-roughness
	normalTex            *texture.Texture2D // Optional normal texture
	occlusionTex         *texture.Texture2D // Optional occlusion texture
	emissiveTex          *texture.Texture2D // Optional emissive texture
	uni                  gls.Uniform        // Uniform location cache
	udata                struct {           // Combined uniform data
		baseColorFactor math32.Color4
		emissiveFactor  math32.Color4
		metallicFactor  float32
		roughnessFactor float32
	}
}

// Number of glsl shader vec4 elements used by uniform data
const pbrMrVec4Count = 3

// NewPbrMr creates and returns a pointer to a new PbrMr material.
func NewPbrMr() *PbrMr {

	m := new(PbrMr)
	m.Material.Init()
	m.SetShader("pbr_mr")

	// Creates uniform and set defaulf values
	m.uni.Init("Material")
	m.udata.baseColorFactor = math32.Color4{1, 1, 1, 1}
	m.udata.emissiveFactor = math32.Color4{0, 0, 0, 0}
	m.udata.metallicFactor = 1.0
	m.udata.roughnessFactor = 1.0
	return m
}

// SetBaseColorFactor sets this material base color.
// Its default value is {1,1,1,1}.
// Returns pointer to this updated material.
func (m *PbrMr) SetBaseColorFactor(c *math32.Color4) *PbrMr {

	m.udata.baseColorFactor = *c
	return m
}

// SetBaseColorTexture sets this material optional texture base color.
// Returns pointer to this updated material.
func (m *PbrMr) SetBaseColorTexture(tex *texture.Texture2D) *PbrMr {

	m.baseColorTex = tex
	if m.baseColorTex != nil {
		m.baseColorTex.SetUniformNames("uBaseColorSampler", "uBaseColorTexParams")
		m.SetShaderDefine("HAS_BASECOLORMAP", "")
	} else {
		m.UnsetShaderDefine("HAS_BASECOLORMAP")
	}
	return m
}

// SetEmissiveFactor sets the emissive color of the material.
// Its default is {0, 0, 0}.
// Returns pointer to this updated material.
func (m *PbrMr) SetEmissiveFactor(c *math32.Color) *PbrMr {

	m.udata.emissiveFactor.R = c.R
	m.udata.emissiveFactor.G = c.G
	m.udata.emissiveFactor.B = c.B
	return m
}

// SetMetallicFactor sets this material metallic factor.
// Its default value is 1.0
// Returns pointer to this updated material.
func (m *PbrMr) SetMetallicFactor(v float32) *PbrMr {

	m.udata.metallicFactor = v
	return m
}

// SetMetallicRoughnessTexture sets this material optional metallic-roughness texture.
// Returns pointer to this updated material.
func (m *PbrMr) SetMetallicRoughnessTexture(tex *texture.Texture2D) *PbrMr {

	m.metallicRoughnessTex = tex
	if m.metallicRoughnessTex != nil {
		m.metallicRoughnessTex.SetUniformNames("uMetallicRoughnessSampler", "uMetallicRoughnessTexParams")
		m.SetShaderDefine("HAS_METALROUGHNESSMAP", "")
	} else {
		m.UnsetShaderDefine("HAS_METALROUGHNESSMAP")
	}
	return m
}

// SetNormalTexture sets this material optional normal texture.
// Returns pointer to this updated material.
func (m *PbrMr) SetNormalTexture(tex *texture.Texture2D) *PbrMr {

	m.normalTex = tex
	if m.normalTex != nil {
		m.normalTex.SetUniformNames("uNormalSampler", "uNormalSamplerTexParams")
		m.SetShaderDefine("HAS_NORMALMAP", "")
	} else {
		m.UnsetShaderDefine("HAS_NORMALMAP")
	}
	return m
}

// SetOcclusionTexture sets this material optional occlusion texture.
// Returns pointer to this updated material.
func (m *PbrMr) SetOcclusionTexture(tex *texture.Texture2D) *PbrMr {

	m.occlusionTex = tex
	if m.occlusionTex != nil {
		m.occlusionTex.SetUniformNames("uOcclusionSampler", "uOcclusionTexParams")
		m.SetShaderDefine("HAS_OCCLUSIONMAP", "")
	} else {
		m.UnsetShaderDefine("HAS_OCCLUSIONMAP")
	}
	return m
}

// SetEmissiveTexture sets this material optional emissive texture.
// Returns pointer to this updated material.
func (m *PbrMr) SetEmissiveTexture(tex *texture.Texture2D) *PbrMr {

	m.emissiveTex = tex
	if m.emissiveTex != nil {
		m.emissiveTex.SetUniformNames("uEmissiveSampler", "uEmissiveSamplerTexParams")
		m.SetShaderDefine("HAS_EMISSIVEMAP", "")
	} else {
		m.UnsetShaderDefine("HAS_EMISSIVEMAP")
	}
	return m
}

// SetRoughnessFactor sets this material roughness factor.
// Its default value is 1.0
// Returns pointer to this updated material.
func (m *PbrMr) SetRoughnessFactor(v float32) *PbrMr {

	m.udata.roughnessFactor = v
	return m
}

// RenderSetup transfer this material uniforms and textures to the shader
func (m *PbrMr) RenderSetup(gl *gls.GLS) {

	m.Material.RenderSetup(gl)
	location := m.uni.Location(gl)
	gl.Uniform4fvUP(location, pbrMrVec4Count, unsafe.Pointer(&m.udata))

	// Transfer optional textures
	if m.baseColorTex != nil {
		m.baseColorTex.RenderSetup(gl, 0)
	}
	if m.metallicRoughnessTex != nil {
		m.metallicRoughnessTex.RenderSetup(gl, 0)
	}
	if m.normalTex != nil {
		m.normalTex.RenderSetup(gl, 0)
	}
	if m.occlusionTex != nil {
		m.occlusionTex.RenderSetup(gl, 0)
	}
	if m.emissiveTex != nil {
		m.emissiveTex.RenderSetup(gl, 0)
	}
}
