// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package material contains several types of materials which
// can be used to set the appearance of graphic object
package material

import (
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/texture"
)

// Material visible side(s)
type Side int

const (
	SideFront  Side = 0
	SideBack   Side = 1
	SideDouble Side = 2
)

// Blending
type Blending int

const (
	BlendingNone        Blending = 0
	BlendingNormal      Blending = 1
	BlendingAdditive    Blending = 2
	BlendingSubtractive Blending = 3
	BlendingMultiply    Blending = 4
	BlendingCustom      Blending = 5
)

// Use lights flags
type UseLights int

const (
	UseLightNone        UseLights = 0x00
	UseLightAmbient     UseLights = 0x01
	UseLightDirectional UseLights = 0x02
	UseLightPoint       UseLights = 0x04
	UseLightSpot        UseLights = 0x08
	UseLightAll         UseLights = 0xFF
)

// Interface for all materials
type IMaterial interface {
	GetMaterial() *Material
	RenderSetup(gs *gls.GLS)
	Dispose()
}

//
// Base Material
//
type Material struct {
	refcount         int                  // Current number of references
	shader           string               // Shader name
	shaderUnique     bool                 // shader has only one instance (does not depend on lights or textures)
	uselights        UseLights            // consider lights for shader selection
	sidevis          Side                 // sides visible
	wireframe        bool                 // show as wirefrme
	depthMask        bool                 // Enable writing into the depth buffer
	depthTest        bool                 // Enable depth buffer test
	depthFunc        uint32               // Actvie depth test function
	blending         Blending             // blending mode
	blendRGB         uint32               // separate blend equation for RGB
	blendAlpha       uint32               // separate blend equation for Alpha
	blendSrcRGB      uint32               // separate blend func source RGB
	blendDstRGB      uint32               // separate blend func dest RGB
	blendSrcAlpha    uint32               // separate blend func source Alpha
	blendDstAlpha    uint32               // separate blend func dest Alpha
	lineWidth        float32              // line width for lines and mesh wireframe
	polyOffsetFactor float32              // polygon offset factor
	polyOffsetUnits  float32              // polygon offset units
	textures         []*texture.Texture2D // List of textures
}

// NewMaterial returns a pointer to a new material
func NewMaterial() *Material {

	mat := new(Material)
	return mat.Init()
}

func (mat *Material) Init() *Material {

	mat.refcount = 1
	mat.uselights = UseLightAll
	mat.sidevis = SideFront
	mat.wireframe = false
	mat.depthMask = true
	mat.depthFunc = gls.LEQUAL
	mat.depthTest = true
	mat.blending = BlendingNormal
	mat.lineWidth = 1.0
	mat.polyOffsetFactor = 0
	mat.polyOffsetUnits = 0
	mat.textures = make([]*texture.Texture2D, 0)

	return mat
}

// GetMaterial satisfies the IMaterial interface
func (mat *Material) GetMaterial() *Material {

	return mat
}

// Incref increments the reference count for this material
// and returns a pointer to the material.
// It should be used when this material is shared by another
// Graphic object.
func (mat *Material) Incref() *Material {

	mat.refcount++
	return mat
}

// Dispose decrements this material reference count and
// if necessary releases OpenGL resources, C memory
// and textures associated with this material.
func (mat *Material) Dispose() {

	if mat.refcount > 1 {
		mat.refcount--
		return
	}
	for i := 0; i < len(mat.textures); i++ {
		mat.textures[i].Dispose()
	}
	mat.Init()
}

// SetShader sets the name of the shader program for this material
func (mat *Material) SetShader(sname string) {

	mat.shader = sname
}

// Shader returns the current name of the shader program for this material
func (mat *Material) Shader() string {

	return mat.shader
}

// SetShaderUnique sets indication that this material shader is unique and
// does not depend on the number of lights in the scene and/or the
// number of textures in the material.
func (mat *Material) SetShaderUnique(unique bool) {

	mat.shaderUnique = unique
}

// ShaderUnique returns this material shader is unique.
func (mat *Material) ShaderUnique() bool {

	return mat.shaderUnique
}

// SetUseLights sets the material use lights bit mask specifying which
// light types will be used when rendering the material
// By default the material will use all lights
func (mat *Material) SetUseLights(lights UseLights) {

	mat.uselights = lights
}

// UseLights returns the current use lights bitmask
func (mat *Material) UseLights() UseLights {

	return mat.uselights
}

// Sets the visible side(s) (SideFront | SideBack | SideDouble)
func (mat *Material) SetSide(side Side) {

	mat.sidevis = side
}

// Side returns the current side visibility for this material
func (mat *Material) Side() Side {

	return mat.sidevis
}

func (mat *Material) SetWireframe(state bool) {

	mat.wireframe = state
}

func (mat *Material) SetDepthMask(state bool) {

	mat.depthMask = state
}

func (mat *Material) SetDepthTest(state bool) {

	mat.depthTest = state
}

func (mat *Material) SetBlending(blending Blending) {

	mat.blending = blending
}

func (mat *Material) SetLineWidth(width float32) {

	mat.lineWidth = width
}

func (mat *Material) SetPolygonOffset(factor, units float32) {

	mat.polyOffsetFactor = factor
	mat.polyOffsetUnits = units
}

func (mat *Material) RenderSetup(gs *gls.GLS) {

	// Sets triangle side view mode
	switch mat.sidevis {
	case SideFront:
		gs.Enable(gls.CULL_FACE)
		gs.FrontFace(gls.CCW)
	case SideBack:
		gs.Enable(gls.CULL_FACE)
		gs.FrontFace(gls.CW)
	case SideDouble:
		gs.Disable(gls.CULL_FACE)
		gs.FrontFace(gls.CCW)
	}

	if mat.depthTest {
		gs.Enable(gls.DEPTH_TEST)
	} else {
		gs.Disable(gls.DEPTH_TEST)
	}
	gs.DepthMask(mat.depthMask)
	gs.DepthFunc(mat.depthFunc)

	if mat.wireframe {
		gs.PolygonMode(gls.FRONT_AND_BACK, gls.LINE)
	} else {
		gs.PolygonMode(gls.FRONT_AND_BACK, gls.FILL)
	}

	// Set polygon offset if requested
	gs.PolygonOffset(mat.polyOffsetFactor, mat.polyOffsetUnits)

	// Sets line width
	gs.LineWidth(mat.lineWidth)

	// Sets blending
	switch mat.blending {
	case BlendingNone:
		gs.Disable(gls.BLEND)
	case BlendingNormal:
		gs.Enable(gls.BLEND)
		gs.BlendEquationSeparate(gls.FUNC_ADD, gls.FUNC_ADD)
		gs.BlendFunc(gls.SRC_ALPHA, gls.ONE_MINUS_SRC_ALPHA)
	case BlendingAdditive:
		gs.Enable(gls.BLEND)
		gs.BlendEquation(gls.FUNC_ADD)
		gs.BlendFunc(gls.SRC_ALPHA, gls.ONE)
	case BlendingSubtractive:
		gs.Enable(gls.BLEND)
		gs.BlendEquation(gls.FUNC_ADD)
		gs.BlendFunc(gls.ZERO, gls.ONE_MINUS_SRC_COLOR)
		break
	case BlendingMultiply:
		gs.Enable(gls.BLEND)
		gs.BlendEquation(gls.FUNC_ADD)
		gs.BlendFunc(gls.ZERO, gls.SRC_COLOR)
		break
	case BlendingCustom:
		gs.BlendEquationSeparate(mat.blendRGB, mat.blendAlpha)
		gs.BlendFuncSeparate(mat.blendSrcRGB, mat.blendDstRGB, mat.blendSrcAlpha, mat.blendDstAlpha)
		break
	default:
		panic("Invalid blending")
	}

	// Render textures
	for idx, tex := range mat.textures {
		tex.RenderSetup(gs, idx)
	}
}

// AddTexture adds the specified Texture2d to the material
func (mat *Material) AddTexture(tex *texture.Texture2D) {

	mat.textures = append(mat.textures, tex)
}

// RemoveTexture removes the specified Texture2d from the material
func (mat *Material) RemoveTexture(tex *texture.Texture2D) {

	for pos, curr := range mat.textures {
		if curr == tex {
			copy(mat.textures[pos:], mat.textures[pos+1:])
			mat.textures[len(mat.textures)-1] = nil
			mat.textures = mat.textures[:len(mat.textures)-1]
			break
		}
	}

}

// HasTexture checks if the material contains the specified texture
func (mat *Material) HasTexture(tex *texture.Texture2D) bool {

	for _, curr := range mat.textures {
		if curr == tex {
			return true
		}
	}
	return false
}

// TextureCount returns the current number of textures
func (mat *Material) TextureCount() int {

	return len(mat.textures)
}
