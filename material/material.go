// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package material contains virtual materials which
// specify the appearance of graphic objects.
package material

import (
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/texture"
	"github.com/g3n/engine/util/logger"
)

// Package logger
var log = logger.New("MATERIAL", logger.Default)

// Side represents the material's visible side(s).
type Side int

// The face side(s) to be rendered. The non-rendered side will be culled to improve performance.
const (
	SideFront = Side(iota)
	SideBack
	SideDouble
)

// Blending specifies the blending mode.
type Blending int

// The various blending types.
const (
	BlendNone = Blending(iota)
	BlendNormal
	BlendAdditive
	BlendSubtractive
	BlendMultiply
	BlendCustom
)

// UseLights is a bitmask that specifies which types of lights affect the material.
type UseLights int

// The possible UseLights values.
const (
	UseLightNone        UseLights = 0x00
	UseLightAmbient     UseLights = 0x01
	UseLightDirectional UseLights = 0x02
	UseLightPoint       UseLights = 0x04
	UseLightSpot        UseLights = 0x08
	UseLightAll         UseLights = 0xFF
)

// IMaterial is the interface for all materials.
type IMaterial interface {
	GetMaterial() *Material
	RenderSetup(gs *gls.GLS)
	Dispose()
}

// Material is the base material.
type Material struct {
	refcount int // Current number of references

	// Shader specification
	shader        string            // Shader name
	shaderUnique  bool              // shader has only one instance (does not depend on lights or textures)
	ShaderDefines gls.ShaderDefines // shader defines

	side        Side                 // Face side(s) visibility
	blending    Blending             // Blending mode
	useLights   UseLights            // Which light types to consider
	transparent bool                 // Whether at all transparent
	wireframe   bool                 // Whether to render only the wireframe
	lineWidth   float32              // Line width for lines and wireframe
	Textures    []*texture.Texture2D // List of Textures

	polyOffsetFactor float32 // polygon offset factor
	polyOffsetUnits  float32 // polygon offset units

	depthMask bool   // Enable writing into the depth buffer
	depthTest bool   // Enable depth buffer test
	depthFunc uint32 // Active depth test function

	// TODO stencil properties

	// Equations used for custom blending (when blending=BlendCustom) // TODO implement methods
	blendRGB      uint32 // separate blending equation for RGB
	blendAlpha    uint32 // separate blending equation for Alpha
	blendSrcRGB   uint32 // separate blending func source RGB
	blendDstRGB   uint32 // separate blending func dest RGB
	blendSrcAlpha uint32 // separate blending func source Alpha
	blendDstAlpha uint32 // separate blending func dest Alpha
}

// NewMaterial creates and returns a pointer to a new Material.
func NewMaterial() *Material {

	mat := new(Material)
	return mat.Init()
}

// Init initializes the material.
func (mat *Material) Init() *Material {

	mat.refcount = 1
	mat.useLights = UseLightAll
	mat.side = SideFront
	mat.transparent = false
	mat.wireframe = false
	mat.depthMask = true
	mat.depthFunc = gls.LEQUAL
	mat.depthTest = true
	mat.blending = BlendNormal
	mat.lineWidth = 1.0
	mat.polyOffsetFactor = 0
	mat.polyOffsetUnits = 0
	mat.Textures = make([]*texture.Texture2D, 0)

	// Setup shader defines and add default values
	mat.ShaderDefines = *gls.NewShaderDefines()

	return mat
}

// GetMaterial satisfies the IMaterial interface.
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
// and Textures associated with this material.
func (mat *Material) Dispose() {

	// Only dispose if last
	if mat.refcount > 1 {
		mat.refcount--
		return
	}
	// Delete Textures
	for i := 0; i < len(mat.Textures); i++ {
		mat.Textures[i].Dispose()
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
// number of Textures in the material.
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

	mat.useLights = lights
}

// UseLights returns the current use lights bitmask
func (mat *Material) UseLights() UseLights {

	return mat.useLights
}

// SetSide sets the visible side(s) (SideFront | SideBack | SideDouble)
func (mat *Material) SetSide(side Side) {

	mat.side = side
}

// Side returns the current side visibility for this material
func (mat *Material) Side() Side {

	return mat.side
}

// SetTransparent sets whether this material is transparent.
func (mat *Material) SetTransparent(state bool) {

	mat.transparent = state
}

// Transparent returns whether this material has been set as transparent.
func (mat *Material) Transparent() bool {

	return mat.transparent
}

// SetWireframe sets whether only the wireframe is rendered.
func (mat *Material) SetWireframe(state bool) {

	mat.wireframe = state
}

// Wireframe returns whether only the wireframe is rendered.
func (mat *Material) Wireframe() bool {

	return mat.wireframe
}

func (mat *Material) SetDepthMask(state bool) {

	mat.depthMask = state
}

func (mat *Material) SetDepthTest(state bool) {

	mat.depthTest = state
}

func (mat *Material) SetDepthFunc(state uint32) {

	mat.depthFunc = state
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

// RenderSetup is called by the renderer before drawing objects with this material.
func (mat *Material) RenderSetup(gs *gls.GLS) {

	// Sets triangle side view mode
	switch mat.side {
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
	case BlendNone:
		gs.Disable(gls.BLEND)
	case BlendNormal:
		gs.Enable(gls.BLEND)
		gs.BlendEquation(gls.FUNC_ADD)
		gs.BlendFunc(gls.SRC_ALPHA, gls.ONE_MINUS_SRC_ALPHA)
	case BlendAdditive:
		gs.Enable(gls.BLEND)
		gs.BlendEquation(gls.FUNC_ADD)
		gs.BlendFunc(gls.SRC_ALPHA, gls.ONE)
	case BlendSubtractive:
		gs.Enable(gls.BLEND)
		gs.BlendEquation(gls.FUNC_ADD)
		gs.BlendFunc(gls.ZERO, gls.ONE_MINUS_SRC_COLOR)
		break
	case BlendMultiply:
		gs.Enable(gls.BLEND)
		gs.BlendEquation(gls.FUNC_ADD)
		gs.BlendFunc(gls.ZERO, gls.SRC_COLOR)
		break
	case BlendCustom:
		gs.BlendEquationSeparate(mat.blendRGB, mat.blendAlpha)
		gs.BlendFuncSeparate(mat.blendSrcRGB, mat.blendDstRGB, mat.blendSrcAlpha, mat.blendDstAlpha)
		break
	default:
		panic("Invalid blending")
	}

	// Render Textures
	// Keep track of counts of unique sampler names to correctly index sampler arrays
	samplerCounts := make(map[string]int)
	for slotIdx, tex := range mat.Textures {
		samplerName, _ := tex.GetUniformNames()
		uniIdx, _ := samplerCounts[samplerName]
		tex.RenderSetup(gs, slotIdx, uniIdx)
		samplerCounts[samplerName] = uniIdx + 1
	}
}

// AddTexture adds the specified Texture2d to the material
func (mat *Material) AddTexture(tex *texture.Texture2D) {

	mat.Textures = append(mat.Textures, tex)
}

// RemoveTexture removes the specified Texture2d from the material
func (mat *Material) RemoveTexture(tex *texture.Texture2D) {

	for pos, curr := range mat.Textures {
		if curr == tex {
			copy(mat.Textures[pos:], mat.Textures[pos+1:])
			mat.Textures[len(mat.Textures)-1] = nil
			mat.Textures = mat.Textures[:len(mat.Textures)-1]
			break
		}
	}
}

// HasTexture checks if the material contains the specified texture
func (mat *Material) HasTexture(tex *texture.Texture2D) bool {

	for _, curr := range mat.Textures {
		if curr == tex {
			return true
		}
	}
	return false
}

// TextureCount returns the current number of Textures
func (mat *Material) TextureCount() int {

	return len(mat.Textures)
}
