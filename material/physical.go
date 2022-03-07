// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package material

import (
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/texture"
)

// Physical is a physically based rendered material which uses the metallic-roughness model.
type Physical struct {
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
const physicalVec4Count = 3

// NewPhysical creates and returns a pointer to a new Physical material.
func NewPhysical() *Physical {

	m := new(Physical)
	m.Material.Init()
	m.SetShader("physical")

	// Creates uniform and set default values
	m.uni.Init("Material")
	m.udata.baseColorFactor = math32.Color4{1, 1, 1, 1}
	m.udata.emissiveFactor = math32.Color4{0, 0, 0, 1}
	m.udata.metallicFactor = 1
	m.udata.roughnessFactor = 1
	return m
}

// SetBaseColorFactor sets this material base color.
// Its default value is {1,1,1,1}.
// Returns pointer to this updated material.
func (m *Physical) SetBaseColorFactor(c *math32.Color4) *Physical {

	m.udata.baseColorFactor = *c
	return m
}

// SetMetallicFactor sets this material metallic factor.
// Its default value is 1.
// Returns pointer to this updated material.
func (m *Physical) SetMetallicFactor(v float32) *Physical {

	m.udata.metallicFactor = v
	return m
}

// SetRoughnessFactor sets this material roughness factor.
// Its default value is 1.
// Returns pointer to this updated material.
func (m *Physical) SetRoughnessFactor(v float32) *Physical {

	m.udata.roughnessFactor = v
	return m
}

// SetEmissiveFactor sets the emissive color of the material.
// Its default is {1, 1, 1}.
// Returns pointer to this updated material.
func (m *Physical) SetEmissiveFactor(c *math32.Color) *Physical {

	m.udata.emissiveFactor.R = c.R
	m.udata.emissiveFactor.G = c.G
	m.udata.emissiveFactor.B = c.B
	return m
}

// SetBaseColorMap sets this material optional texture base color.
// Returns pointer to this updated material.
func (m *Physical) SetBaseColorMap(tex *texture.Texture2D) *Physical {

	m.baseColorTex = tex
	if m.baseColorTex != nil {
		m.baseColorTex.SetUniformNames("uBaseColorSampler", "uBaseColorTexParams")
		m.ShaderDefines.Set("HAS_BASECOLORMAP", "")
		m.AddTexture(m.baseColorTex)
	} else {
		m.ShaderDefines.Unset("HAS_BASECOLORMAP")
		m.RemoveTexture(m.baseColorTex)
	}
	return m
}

// SetMetallicRoughnessMap sets this material optional metallic-roughness texture.
// Returns pointer to this updated material.
func (m *Physical) SetMetallicRoughnessMap(tex *texture.Texture2D) *Physical {

	m.metallicRoughnessTex = tex
	if m.metallicRoughnessTex != nil {
		m.metallicRoughnessTex.SetUniformNames("uMetallicRoughnessSampler", "uMetallicRoughnessTexParams")
		m.ShaderDefines.Set("HAS_METALROUGHNESSMAP", "")
		m.AddTexture(m.metallicRoughnessTex)
	} else {
		m.ShaderDefines.Unset("HAS_METALROUGHNESSMAP")
		m.RemoveTexture(m.metallicRoughnessTex)
	}
	return m
}

// SetNormalMap sets this material optional normal texture.
// Returns pointer to this updated material.
// TODO add SetNormalMap (and SetSpecularMap) to StandardMaterial.
func (m *Physical) SetNormalMap(tex *texture.Texture2D) *Physical {

	m.normalTex = tex
	if m.normalTex != nil {
		m.normalTex.SetUniformNames("uNormalSampler", "uNormalTexParams")
		m.ShaderDefines.Set("HAS_NORMALMAP", "")
		m.AddTexture(m.normalTex)
	} else {
		m.ShaderDefines.Unset("HAS_NORMALMAP")
		m.RemoveTexture(m.normalTex)
	}
	return m
}

// SetOcclusionMap sets this material optional occlusion texture.
// Returns pointer to this updated material.
func (m *Physical) SetOcclusionMap(tex *texture.Texture2D) *Physical {

	m.occlusionTex = tex
	if m.occlusionTex != nil {
		m.occlusionTex.SetUniformNames("uOcclusionSampler", "uOcclusionTexParams")
		m.ShaderDefines.Set("HAS_OCCLUSIONMAP", "")
		m.AddTexture(m.occlusionTex)
	} else {
		m.ShaderDefines.Unset("HAS_OCCLUSIONMAP")
		m.RemoveTexture(m.occlusionTex)
	}
	return m
}

// SetEmissiveMap sets this material optional emissive texture.
// Returns pointer to this updated material.
func (m *Physical) SetEmissiveMap(tex *texture.Texture2D) *Physical {

	m.emissiveTex = tex
	if m.emissiveTex != nil {
		m.emissiveTex.SetUniformNames("uEmissiveSampler", "uEmissiveTexParams")
		m.ShaderDefines.Set("HAS_EMISSIVEMAP", "")
		m.AddTexture(m.emissiveTex)
	} else {
		m.ShaderDefines.Unset("HAS_EMISSIVEMAP")
		m.RemoveTexture(m.emissiveTex)
	}
	return m
}

// RenderSetup transfer this material uniforms and textures to the shader
func (m *Physical) RenderSetup(gl *gls.GLS) {

	m.Material.RenderSetup(gl)
	location := m.uni.Location(gl)
	gl.Uniform4fv(location, physicalVec4Count, &m.udata.baseColorFactor.R)
}
