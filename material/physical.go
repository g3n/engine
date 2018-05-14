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
	m.udata.emissiveFactor = math32.Color4{0, 0, 0, 0}
	m.udata.metallicFactor = 0.5
	m.udata.roughnessFactor = 0.5
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
// Its default value is 1.0
// Returns pointer to this updated material.
func (m *Physical) SetMetallicFactor(v float32) *Physical {

	m.udata.metallicFactor = v
	return m
}

// SetRoughnessFactor sets this material roughness factor.
// Its default value is 1.0
// Returns pointer to this updated material.
func (m *Physical) SetRoughnessFactor(v float32) *Physical {

	m.udata.roughnessFactor = v
	return m
}

// SetEmissiveFactor sets the emissive color of the material.
// Its default is {0, 0, 0}.
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
		m.SetShaderDefine("HAS_BASECOLORMAP", "")
	} else {
		m.UnsetShaderDefine("HAS_BASECOLORMAP")
	}
	return m
}

// SetMetallicRoughnessMap sets this material optional metallic-roughness texture.
// Returns pointer to this updated material.
func (m *Physical) SetMetallicRoughnessMap(tex *texture.Texture2D) *Physical {

	m.metallicRoughnessTex = tex
	if m.metallicRoughnessTex != nil {
		m.metallicRoughnessTex.SetUniformNames("uMetallicRoughnessSampler", "uMetallicRoughnessTexParams")
		m.SetShaderDefine("HAS_METALROUGHNESSMAP", "")
	} else {
		m.UnsetShaderDefine("HAS_METALROUGHNESSMAP")
	}
	return m
}

// TODO add SetNormalMap (and SetSpecularMap) to StandardMaterial.
// SetNormalMap sets this material optional normal texture.
// Returns pointer to this updated material.
func (m *Physical) SetNormalMap(tex *texture.Texture2D) *Physical {

	m.normalTex = tex
	if m.normalTex != nil {
		m.normalTex.SetUniformNames("uNormalSampler", "uNormalSamplerTexParams")
		m.SetShaderDefine("HAS_NORMALMAP", "")
	} else {
		m.UnsetShaderDefine("HAS_NORMALMAP")
	}
	return m
}

// SetOcclusionMap sets this material optional occlusion texture.
// Returns pointer to this updated material.
func (m *Physical) SetOcclusionMap(tex *texture.Texture2D) *Physical {

	m.occlusionTex = tex
	if m.occlusionTex != nil {
		m.occlusionTex.SetUniformNames("uOcclusionSampler", "uOcclusionTexParams")
		m.SetShaderDefine("HAS_OCCLUSIONMAP", "")
	} else {
		m.UnsetShaderDefine("HAS_OCCLUSIONMAP")
	}
	return m
}

// SetEmissiveMap sets this material optional emissive texture.
// Returns pointer to this updated material.
func (m *Physical) SetEmissiveMap(tex *texture.Texture2D) *Physical {

	m.emissiveTex = tex
	if m.emissiveTex != nil {
		m.emissiveTex.SetUniformNames("uEmissiveSampler", "uEmissiveSamplerTexParams")
		m.SetShaderDefine("HAS_EMISSIVEMAP", "")
	} else {
		m.UnsetShaderDefine("HAS_EMISSIVEMAP")
	}
	return m
}

// RenderSetup transfer this material uniforms and textures to the shader
func (m *Physical) RenderSetup(gl *gls.GLS) {

	m.Material.RenderSetup(gl)
	location := m.uni.Location(gl)
	log.Error("Physical RenderSetup location:%v udata:%+v", location, m.udata)
	gl.Uniform4fvUP(location, physicalVec4Count, unsafe.Pointer(&m.udata))

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
