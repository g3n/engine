// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build wasm

package gls

import (
	"fmt"
	"syscall/js"
	"unsafe"
)

// GLS encapsulates the state of a WebGL context and contains
// methods to call WebGL functions.
type GLS struct {
	stats       Stats             // statistics
	prog        *Program          // current active shader program
	programs    map[*Program]bool // shader programs cache
	checkErrors bool              // check openGL API errors flag

	// Cache WebGL state to avoid making unnecessary API calls
	activeTexture       uint32      // cached last set active texture unit
	viewportX           int32       // cached last set viewport x
	viewportY           int32       // cached last set viewport y
	viewportWidth       int32       // cached last set viewport width
	viewportHeight      int32       // cached last set viewport height
	lineWidth           float32     // cached last set line width
	sideView            int         // cached last set triangle side view mode
	frontFace           uint32      // cached last set glFrontFace value
	depthFunc           uint32      // cached last set depth function
	depthMask           int         // cached last set depth mask
	capabilities        map[int]int // cached capabilities (Enable/Disable)
	blendEquation       uint32      // cached last set blend equation value
	blendSrc            uint32      // cached last set blend src value
	blendDst            uint32      // cached last set blend equation destination value
	blendEquationRGB    uint32      // cached last set blend equation rgb value
	blendEquationAlpha  uint32      // cached last set blend equation alpha value
	blendSrcRGB         uint32      // cached last set blend src rgb
	blendSrcAlpha       uint32      // cached last set blend src alpha value
	blendDstRGB         uint32      // cached last set blend destination rgb value
	blendDstAlpha       uint32      // cached last set blend destination alpha value
	polygonModeFace     uint32      // cached last set polygon mode face
	polygonModeMode     uint32      // cached last set polygon mode mode
	polygonOffsetFactor float32     // cached last set polygon offset factor
	polygonOffsetUnits  float32     // cached last set polygon offset units

	// js.Value storage maps
	programMap      map[uint32]js.Value
	shaderMap       map[uint32]js.Value
	bufferMap       map[uint32]js.Value
	framebufferMap  map[uint32]js.Value
	renderbufferMap map[uint32]js.Value
	textureMap      map[uint32]js.Value
	uniformMap      map[uint32]js.Value
	vertexArrayMap  map[uint32]js.Value

	// Next free index to be used for each map
	programMapIndex      uint32
	shaderMapIndex       uint32
	bufferMapIndex       uint32
	framebufferMapIndex  uint32
	renderbufferMapIndex uint32
	textureMapIndex      uint32
	uniformMapIndex      uint32
	vertexArrayMapIndex  uint32

	// Canvas and WebGL Context
	canvas js.Value
	gl     js.Value
}

// New creates and returns a new instance of a GLS object,
// which encapsulates the state of an WebGL context.
// This should be called only after an active WebGL context
// is established, such as by creating a new window.
func New(webglCtx js.Value) (*GLS, error) {

	gs := new(GLS)
	gs.reset()
	gs.checkErrors = false
	gs.gl = webglCtx

	// Create js.Value storage maps
	gs.programMap = make(map[uint32]js.Value)
	gs.shaderMap = make(map[uint32]js.Value)
	gs.bufferMap = make(map[uint32]js.Value)
	gs.framebufferMap = make(map[uint32]js.Value)
	gs.renderbufferMap = make(map[uint32]js.Value)
	gs.textureMap = make(map[uint32]js.Value)
	gs.uniformMap = make(map[uint32]js.Value)
	gs.vertexArrayMap = make(map[uint32]js.Value)

	// Initialize indexes to be used with the maps above
	gs.programMapIndex = 1
	gs.shaderMapIndex = 1
	gs.bufferMapIndex = 1
	gs.framebufferMapIndex = 1
	gs.renderbufferMapIndex = 1
	gs.textureMapIndex = 1
	gs.uniformMapIndex = 1
	gs.vertexArrayMapIndex = 1

	gs.setDefaultState()
	return gs, nil
}

// SetCheckErrors enables/disables checking for errors after the
// call of any WebGL function. It is enabled by default but
// could be disabled after an application is stable to improve the performance.
func (gs *GLS) SetCheckErrors(enable bool) {

	gs.checkErrors = enable
}

// CheckErrors returns if error checking is enabled or not.
func (gs *GLS) CheckErrors() bool {

	return gs.checkErrors
}

// reset resets the internal state kept of the WebGL
func (gs *GLS) reset() {

	gs.lineWidth = 0.0
	gs.sideView = uintUndef
	gs.frontFace = 0
	gs.depthFunc = 0
	gs.depthMask = uintUndef
	gs.capabilities = make(map[int]int)
	gs.programs = make(map[*Program]bool)
	gs.prog = nil

	gs.activeTexture = uintUndef
	gs.blendEquation = uintUndef
	gs.blendSrc = uintUndef
	gs.blendDst = uintUndef
	gs.blendEquationRGB = 0
	gs.blendEquationAlpha = 0
	gs.blendSrcRGB = uintUndef
	gs.blendSrcAlpha = uintUndef
	gs.blendDstRGB = uintUndef
	gs.blendDstAlpha = uintUndef
	gs.polygonModeFace = 0
	gs.polygonModeMode = 0
	gs.polygonOffsetFactor = -1
	gs.polygonOffsetUnits = -1
}

// setDefaultState is used internally to set the initial state of WebGL
// for this context.
func (gs *GLS) setDefaultState() {

	gs.ClearColor(0, 0, 0, 1)
	gs.ClearDepth(1)
	gs.ClearStencil(0)
	gs.Enable(DEPTH_TEST)
	gs.DepthFunc(LEQUAL)
	gs.FrontFace(CCW)
	gs.CullFace(BACK)
	gs.Enable(CULL_FACE)
	gs.Enable(BLEND)
	gs.BlendEquation(FUNC_ADD)
	gs.BlendFunc(SRC_ALPHA, ONE_MINUS_SRC_ALPHA)

	// TODO commented constants not available in WebGL
	//gs.Enable(VERTEX_PROGRAM_POINT_SIZE)
	//gs.Enable(PROGRAM_POINT_SIZE)
	//gs.Enable(MULTISAMPLE)
	gs.Enable(POLYGON_OFFSET_FILL)
	//gs.Enable(POLYGON_OFFSET_LINE)
	//gs.Enable(POLYGON_OFFSET_POINT)
}

// Stats copy the current values of the internal statistics structure
// to the specified pointer.
func (gs *GLS) Stats(s *Stats) {

	*s = gs.stats
	s.Shaders = len(gs.programs)
}

// ActiveTexture selects which texture unit subsequent texture state calls
// will affect. The number of texture units an implementation supports is
// implementation dependent, but must be at least 48 in GL 3.3.
func (gs *GLS) ActiveTexture(texture uint32) {

	if gs.activeTexture == texture {
		return
	}
	gs.gl.Call("activeTexture", int(texture))
	gs.checkError("ActiveTexture")
	gs.activeTexture = texture
}

// AttachShader attaches the specified shader object to the specified program object.
func (gs *GLS) AttachShader(program, shader uint32) {

	gs.gl.Call("attachShader", gs.programMap[program], gs.shaderMap[shader])
	gs.checkError("AttachShader")
}

// BindBuffer binds a buffer object to the specified buffer binding point.
func (gs *GLS) BindBuffer(target int, vbo uint32) {

	gs.gl.Call("bindBuffer", target, gs.bufferMap[vbo])
	gs.checkError("BindBuffer")
}

// BindTexture lets you create or use a named texture.
func (gs *GLS) BindTexture(target int, tex uint32) {

	gs.gl.Call("bindTexture", target, gs.textureMap[tex])
	gs.checkError("BindTexture")
}

// BindVertexArray binds the vertex array object.
func (gs *GLS) BindVertexArray(vao uint32) {

	gs.gl.Call("bindVertexArray", gs.vertexArrayMap[vao])
	gs.checkError("BindVertexArray")
}

// BlendEquation sets the blend equations for all draw buffers.
func (gs *GLS) BlendEquation(mode uint32) {

	if gs.blendEquation == mode {
		return
	}
	gs.gl.Call("blendEquation", int(mode))
	gs.checkError("BlendEquation")
	gs.blendEquation = mode
}

// BlendEquationSeparate sets the blend equations for all draw buffers
// allowing different equations for the RGB and alpha components.
func (gs *GLS) BlendEquationSeparate(modeRGB uint32, modeAlpha uint32) {

	if gs.blendEquationRGB == modeRGB && gs.blendEquationAlpha == modeAlpha {
		return
	}
	gs.gl.Call("blendEquationSeparate", int(modeRGB), int(modeAlpha))
	gs.checkError("BlendEquationSeparate")
	gs.blendEquationRGB = modeRGB
	gs.blendEquationAlpha = modeAlpha
}

// BlendFunc defines the operation of blending for
// all draw buffers when blending is enabled.
func (gs *GLS) BlendFunc(sfactor, dfactor uint32) {

	if gs.blendSrc == sfactor && gs.blendDst == dfactor {
		return
	}
	gs.gl.Call("blendFunc", int(sfactor), int(dfactor))
	gs.checkError("BlendFunc")
	gs.blendSrc = sfactor
	gs.blendDst = dfactor
}

// BlendFuncSeparate defines the operation of blending for all draw buffers when blending
// is enabled, allowing different operations for the RGB and alpha components.
func (gs *GLS) BlendFuncSeparate(srcRGB uint32, dstRGB uint32, srcAlpha uint32, dstAlpha uint32) {

	if gs.blendSrcRGB == srcRGB && gs.blendDstRGB == dstRGB &&
		gs.blendSrcAlpha == srcAlpha && gs.blendDstAlpha == dstAlpha {
		return
	}
	gs.gl.Call("blendFuncSeparate", int(srcRGB), int(dstRGB), int(srcAlpha), int(dstAlpha))
	gs.checkError("BlendFuncSeparate")
	gs.blendSrcRGB = srcRGB
	gs.blendDstRGB = dstRGB
	gs.blendSrcAlpha = srcAlpha
	gs.blendDstAlpha = dstAlpha
}

// BufferData creates a new data store for the buffer object currently
// bound to target, deleting any pre-existing data store.
func (gs *GLS) BufferData(target uint32, size int, data interface{}, usage uint32) {

	dataTA := js.TypedArrayOf(data)
	gs.gl.Call("bufferData", int(target), dataTA, int(usage))
	gs.checkError("BufferData")
	dataTA.Release()
}

// ClearColor specifies the red, green, blue, and alpha values
// used by glClear to clear the color buffers.
func (gs *GLS) ClearColor(r, g, b, a float32) {

	gs.gl.Call("clearColor", r, g, b, a)
	gs.checkError("ClearColor")
}

// ClearDepth specifies the depth value used by Clear to clear the depth buffer.
func (gs *GLS) ClearDepth(v float32) {

	gs.gl.Call("clearDepth", v)
	gs.checkError("ClearDepth")
}

// ClearStencil specifies the index used by Clear to clear the stencil buffer.
func (gs *GLS) ClearStencil(v int32) {

	gs.gl.Call("clearStencil", int(v))
	gs.checkError("ClearStencil")
}

// Clear sets the bitplane area of the window to values previously
// selected by ClearColor, ClearDepth, and ClearStencil.
func (gs *GLS) Clear(mask uint) {

	gs.gl.Call("clear", int(mask))
	gs.checkError("Clear")
}

// CompileShader compiles the source code strings that
// have been stored in the specified shader object.
func (gs *GLS) CompileShader(shader uint32) {

	gs.gl.Call("compileShader", gs.shaderMap[shader])
	gs.checkError("CompileShader")
}

// CreateProgram creates an empty program object and returns
// a non-zero value by which it can be referenced.
func (gs *GLS) CreateProgram() uint32 {

	gs.programMap[gs.programMapIndex] = gs.gl.Call("createProgram")
	gs.checkError("CreateProgram")
	idx := gs.programMapIndex
	gs.programMapIndex++
	return idx
}

// CreateShader creates an empty shader object and returns
// a non-zero value by which it can be referenced.
func (gs *GLS) CreateShader(stype uint32) uint32 {

	gs.shaderMap[gs.shaderMapIndex] = gs.gl.Call("createShader", int(stype))
	gs.checkError("CreateShader")
	idx := gs.shaderMapIndex
	gs.shaderMapIndex++
	return idx
}

// DeleteBuffers deletes n​buffer objects named
// by the elements of the provided array.
func (gs *GLS) DeleteBuffers(bufs ...uint32) {

	for _, buf := range bufs {
		gs.gl.Call("deleteBuffer", gs.bufferMap[buf])
		gs.checkError("DeleteBuffers")
		gs.stats.Buffers--
		delete(gs.bufferMap, buf)
	}
}

// DeleteShader frees the memory and invalidates the name
// associated with the specified shader object.
func (gs *GLS) DeleteShader(shader uint32) {

	gs.gl.Call("deleteShader", gs.shaderMap[shader])
	gs.checkError("DeleteShader")
	delete(gs.shaderMap, shader)
}

// DeleteProgram frees the memory and invalidates the name
// associated with the specified program object.
func (gs *GLS) DeleteProgram(program uint32) {

	gs.gl.Call("deleteProgram", gs.programMap[program])
	gs.checkError("DeleteProgram")
	delete(gs.programMap, program)
}

// DeleteTextures deletes n​textures named
// by the elements of the provided array.
func (gs *GLS) DeleteTextures(tex ...uint32) {

	for _, t := range tex {
		gs.gl.Call("deleteTexture", gs.textureMap[t])
		gs.checkError("DeleteTextures")
		delete(gs.textureMap, t)
		gs.stats.Textures--
	}
}

// DeleteVertexArrays deletes n​vertex array objects named
// by the elements of the provided array.
func (gs *GLS) DeleteVertexArrays(vaos ...uint32) {

	for _, v := range vaos {
		gs.gl.Call("deleteVertexArray", gs.vertexArrayMap[v])
		gs.checkError("DeleteVertexArrays")
		delete(gs.vertexArrayMap, v)
		gs.stats.Vaos--
	}
}

// TODO ReadPixels

// DepthFunc specifies the function used to compare each incoming pixel
// depth value with the depth value present in the depth buffer.
func (gs *GLS) DepthFunc(mode uint32) {

	if gs.depthFunc == mode {
		return
	}
	gs.gl.Call("depthFunc", int(mode))
	gs.checkError("DepthFunc")
	gs.depthFunc = mode
}

// DepthMask enables or disables writing into the depth buffer.
func (gs *GLS) DepthMask(flag bool) {

	if gs.depthMask == intTrue && flag {
		return
	}
	if gs.depthMask == intFalse && !flag {
		return
	}
	gs.gl.Call("depthMask", flag)
	gs.checkError("DepthMask")
	if flag {
		gs.depthMask = intTrue
	} else {
		gs.depthMask = intFalse
	}
}

// DrawArrays renders primitives from array data.
func (gs *GLS) DrawArrays(mode uint32, first int32, count int32) {

	gs.gl.Call("drawArrays", int(mode), first, count)
	gs.checkError("DrawArrays")
	gs.stats.Drawcalls++
}

// DrawElements renders primitives from array data.
func (gs *GLS) DrawElements(mode uint32, count int32, itype uint32, start uint32) {

	gs.gl.Call("drawElements", int(mode), count, int(itype), start)
	gs.checkError("DrawElements")
	gs.stats.Drawcalls++
}

// Enable enables the specified capability.
func (gs *GLS) Enable(cap int) {

	if gs.capabilities[cap] == capEnabled {
		gs.stats.Caphits++
		return
	}
	gs.gl.Call("enable", int32(cap))
	gs.checkError("Enable")
	gs.capabilities[cap] = capEnabled
}

// Disable disables the specified capability.
func (gs *GLS) Disable(cap int) {

	if gs.capabilities[cap] == capDisabled {
		gs.stats.Caphits++
		return
	}
	gs.gl.Call("disable", cap)
	gs.checkError("Disable")
	gs.capabilities[cap] = capDisabled
}

// EnableVertexAttribArray enables a generic vertex attribute array.
func (gs *GLS) EnableVertexAttribArray(index uint32) {

	gs.gl.Call("enableVertexAttribArray", index)
	gs.checkError("EnableVertexAttribArray")
}

// CullFace specifies whether front- or back-facing facets can be culled.
func (gs *GLS) CullFace(mode uint32) {

	gs.gl.Call("cullFace", int(mode))
	gs.checkError("CullFace")
}

// FrontFace defines front- and back-facing polygons.
func (gs *GLS) FrontFace(mode uint32) {

	if gs.frontFace == mode {
		return
	}
	gs.gl.Call("frontFace", int(mode))
	gs.checkError("FrontFace")
	gs.frontFace = mode
}

// GenBuffer generates a ​buffer object name.
func (gs *GLS) GenBuffer() uint32 {

	gs.bufferMap[gs.bufferMapIndex] = gs.gl.Call("createBuffer")
	gs.checkError("CreateBuffer")
	idx := gs.bufferMapIndex
	gs.bufferMapIndex++
	gs.stats.Buffers++
	return idx
}

// GenerateMipmap generates mipmaps for the specified texture target.
func (gs *GLS) GenerateMipmap(target uint32) {

	gs.gl.Call("generateMipmap", int(target))
	gs.checkError("GenerateMipmap")
}

// GenTexture generates a texture object name.
func (gs *GLS) GenTexture() uint32 {

	gs.textureMap[gs.textureMapIndex] = gs.gl.Call("createTexture")
	gs.checkError("GenTexture")
	idx := gs.textureMapIndex
	gs.textureMapIndex++
	gs.stats.Textures++
	return idx
}

// GenVertexArray generates a vertex array object name.
func (gs *GLS) GenVertexArray() uint32 {

	gs.vertexArrayMap[gs.vertexArrayMapIndex] = gs.gl.Call("createVertexArray")
	gs.checkError("GenVertexArray")
	idx := gs.vertexArrayMapIndex
	gs.vertexArrayMapIndex++
	gs.stats.Vaos++
	return idx
}

// GetAttribLocation returns the location of the specified attribute variable.
func (gs *GLS) GetAttribLocation(program uint32, name string) int32 {

	loc := gs.gl.Call("getAttribLocation", gs.programMap[program], name).Int()
	gs.checkError("GetAttribLocation")
	return int32(loc)
}

// GetProgramiv returns the specified parameter from the specified program object.
func (gs *GLS) GetProgramiv(program, pname uint32, params *int32) {

	sparam := gs.gl.Call("getProgramParameter", gs.programMap[program], int(pname))
	gs.checkError("GetProgramiv")
	switch pname {
	case DELETE_STATUS, LINK_STATUS, VALIDATE_STATUS:
		if sparam.Bool() {
			*params = TRUE
		} else {
			*params = FALSE
		}
	default:
		*params = int32(sparam.Int())
	}
}

// GetProgramInfoLog returns the information log for the specified program object.
func (gs *GLS) GetProgramInfoLog(program uint32) string {

	res := gs.gl.Call("getProgramInfoLog", gs.programMap[program]).String()
	gs.checkError("GetProgramInfoLog")
	return res
}

// GetShaderInfoLog returns the information log for the specified shader object.
func (gs *GLS) GetShaderInfoLog(shader uint32) string {

	res := gs.gl.Call("getShaderInfoLog", gs.shaderMap[shader]).String()
	gs.checkError("GetShaderInfoLog")
	return res
}

// GetString returns a string describing the specified aspect of the current GL connection.
func (gs *GLS) GetString(name uint32) string {

	res := gs.gl.Call("getParameter", int(name)).String()
	gs.checkError("GetString")
	return res
}

// GetUniformLocation returns the location of a uniform variable for the specified program.
func (gs *GLS) GetUniformLocation(program uint32, name string) int32 {

	loc := gs.gl.Call("getUniformLocation", gs.programMap[program], name)
	if loc == js.Null() {
		return -1
	}
	gs.uniformMap[gs.uniformMapIndex] = loc
	gs.checkError("GetUniformLocation")
	idx := gs.uniformMapIndex
	gs.uniformMapIndex++
	return int32(idx)
}

// GetViewport returns the current viewport information.
func (gs *GLS) GetViewport() (x, y, width, height int32) {

	return gs.viewportX, gs.viewportY, gs.viewportWidth, gs.viewportHeight
}

// LineWidth specifies the rasterized width of both aliased and antialiased lines.
func (gs *GLS) LineWidth(width float32) {

	if gs.lineWidth == width {
		return
	}
	gs.gl.Call("lineWidth", width)
	gs.checkError("LineWidth")
	gs.lineWidth = width
}

// LinkProgram links the specified program object.
func (gs *GLS) LinkProgram(program uint32) {

	gs.gl.Call("linkProgram", gs.programMap[program])
	gs.checkError("LinkProgram")
}

// GetShaderiv returns the specified parameter from the specified shader object.
func (gs *GLS) GetShaderiv(shader, pname uint32, params *int32) {

	sparam := gs.gl.Call("getShaderParameter", gs.shaderMap[shader], int(pname))
	gs.checkError("GetShaderiv")
	switch pname {
	case DELETE_STATUS, COMPILE_STATUS:
		if sparam.Bool() {
			*params = TRUE
		} else {
			*params = FALSE
		}
	default:
		*params = int32(sparam.Int())
	}
}

// Scissor defines the scissor box rectangle in window coordinates.
func (gs *GLS) Scissor(x, y int32, width, height uint32) {

	gs.gl.Call("scissor", x, y, int(width), int(height))
	gs.checkError("Scissor")
}

// ShaderSource sets the source code for the specified shader object.
func (gs *GLS) ShaderSource(shader uint32, src string) {

	gs.gl.Call("shaderSource", gs.shaderMap[shader], src)
	gs.checkError("ShaderSource")
}

// TexImage2D specifies a two-dimensional texture image.
func (gs *GLS) TexImage2D(target uint32, level int32, iformat int32, width int32, height int32, format uint32, itype uint32, data interface{}) {

	dataTA := js.TypedArrayOf(data)
	gs.gl.Call("texImage2D", int(target), level, iformat, width, height, 0, int(format), int(itype), dataTA)
	gs.checkError("TexImage2D")
	dataTA.Release()
}

// TexParameteri sets the specified texture parameter on the specified texture.
func (gs *GLS) TexParameteri(target uint32, pname uint32, param int32) {

	gs.gl.Call("texParameteri", int(target), int(pname), param)
	gs.checkError("TexParameteri")
}

// PolygonMode controls the interpretation of polygons for rasterization.
func (gs *GLS) PolygonMode(face, mode uint32) {

	log.Warn("PolygonMode not available in WebGL")
}

// PolygonOffset sets the scale and units used to calculate depth values.
func (gs *GLS) PolygonOffset(factor float32, units float32) {

	if gs.polygonOffsetFactor == factor && gs.polygonOffsetUnits == units {
		return
	}
	gs.gl.Call("polygonOffset", factor, units)
	gs.checkError("PolygonOffset")
	gs.polygonOffsetFactor = factor
	gs.polygonOffsetUnits = units
}

// Uniform1i sets the value of an int uniform variable for the current program object.
func (gs *GLS) Uniform1i(location int32, v0 int32) {

	gs.gl.Call("uniform1i", gs.uniformMap[uint32(location)], v0)
	gs.checkError("Uniform1i")
	gs.stats.Unisets++
}

// Uniform1f sets the value of a float uniform variable for the current program object.
func (gs *GLS) Uniform1f(location int32, v0 float32) {

	gs.gl.Call("uniform1f", gs.uniformMap[uint32(location)], v0)
	gs.checkError("Uniform1f")
	gs.stats.Unisets++
}

// Uniform2f sets the value of a vec2 uniform variable for the current program object.
func (gs *GLS) Uniform2f(location int32, v0, v1 float32) {

	gs.gl.Call("uniform2f", gs.uniformMap[uint32(location)], v0, v1)
	gs.checkError("Uniform2f")
	gs.stats.Unisets++
}

// Uniform3f sets the value of a vec3 uniform variable for the current program object.
func (gs *GLS) Uniform3f(location int32, v0, v1, v2 float32) {

	gs.gl.Call("uniform3f", gs.uniformMap[uint32(location)], v0, v1, v2)
	gs.checkError("Uniform3f")
	gs.stats.Unisets++
}

// Uniform4f sets the value of a vec4 uniform variable for the current program object.
func (gs *GLS) Uniform4f(location int32, v0, v1, v2, v3 float32) {

	gs.gl.Call("uniform4f", gs.uniformMap[uint32(location)], v0, v1, v2, v3)
	gs.checkError("Uniform4f")
	gs.stats.Unisets++
}

//// UniformMatrix3fv sets the value of one or many 3x3 float matrices for the current program object.
func (gs *GLS) UniformMatrix3fv(location int32, count int32, transpose bool, pm *float32) {

	data := (*[1 << 30]float32)(unsafe.Pointer(pm))[:9*count]
	dataTA := js.TypedArrayOf(data)
	gs.gl.Call("uniformMatrix3fv", gs.uniformMap[uint32(location)], transpose, dataTA)
	dataTA.Release()
	gs.checkError("UniformMatrix3fv")
	gs.stats.Unisets++
}

// UniformMatrix4fv sets the value of one or many 4x4 float matrices for the current program object.
func (gs *GLS) UniformMatrix4fv(location int32, count int32, transpose bool, pm *float32) {

	data := (*[1 << 30]float32)(unsafe.Pointer(pm))[:16*count]
	dataTA := js.TypedArrayOf(data)
	gs.gl.Call("uniformMatrix4fv", gs.uniformMap[uint32(location)], transpose, dataTA)
	dataTA.Release()
	gs.checkError("UniformMatrix4fv")
	gs.stats.Unisets++
}

// Uniform1fv sets the value of one or many float uniform variables for the current program object.
func (gs *GLS) Uniform1fv(location int32, count int32, v *float32) {

	data := (*[1 << 30]float32)(unsafe.Pointer(v))[:count]
	dataTA := js.TypedArrayOf(data)
	gs.gl.Call("uniform1fv", gs.uniformMap[uint32(location)], dataTA)
	dataTA.Release()
	gs.checkError("Uniform1fv")
	gs.stats.Unisets++
}

// Uniform2fv sets the value of one or many vec2 uniform variables for the current program object.
func (gs *GLS) Uniform2fv(location int32, count int32, v *float32) {

	data := (*[1 << 30]float32)(unsafe.Pointer(v))[:2*count]
	dataTA := js.TypedArrayOf(data)
	gs.gl.Call("uniform2fv", gs.uniformMap[uint32(location)], dataTA)
	dataTA.Release()
	gs.checkError("Uniform2fv")
	gs.stats.Unisets++
}

// Uniform3fv sets the value of one or many vec3 uniform variables for the current program object.
func (gs *GLS) Uniform3fv(location int32, count int32, v *float32) {

	data := (*[1 << 30]float32)(unsafe.Pointer(v))[:3*count]
	dataTA := js.TypedArrayOf(data)
	gs.gl.Call("uniform3fv", gs.uniformMap[uint32(location)], dataTA)
	dataTA.Release()
	gs.checkError("Uniform3fv")
	gs.stats.Unisets++
}

// Uniform4fv sets the value of one or many vec4 uniform variables for the current program object.
func (gs *GLS) Uniform4fv(location int32, count int32, v *float32) {

	data := (*[1 << 30]float32)(unsafe.Pointer(v))[:4*count]
	dataTA := js.TypedArrayOf(data)
	gs.gl.Call("uniform4fv", gs.uniformMap[uint32(location)], dataTA)
	dataTA.Release()
	gs.checkError("Uniform4fv")
	gs.stats.Unisets++
}

// VertexAttribPointer defines an array of generic vertex attribute data.
func (gs *GLS) VertexAttribPointer(index uint32, size int32, xtype uint32, normalized bool, stride int32, offset uint32) {

	gs.gl.Call("vertexAttribPointer", index, size, int(xtype), normalized, stride, offset)
	gs.checkError("VertexAttribPointer")
}

// Viewport sets the viewport.
func (gs *GLS) Viewport(x, y, width, height int32) {

	gs.gl.Call("viewport", x, y, width, height)
	gs.checkError("Viewport")
	gs.viewportX = x
	gs.viewportY = y
	gs.viewportWidth = width
	gs.viewportHeight = height
}

// UseProgram sets the specified program as the current program.
func (gs *GLS) UseProgram(prog *Program) {

	if prog.handle == 0 {
		panic("Invalid program")
	}

	gs.gl.Call("useProgram", gs.programMap[prog.handle])
	gs.checkError("UseProgram")
	gs.prog = prog

	// Inserts program in cache if not already there.
	if !gs.programs[prog] {
		gs.programs[prog] = true
		log.Debug("New Program activated. Total: %d", len(gs.programs))
	}
}

// checkError checks if there are any WebGL errors and panics if so.
func (gs *GLS) checkError(name string) {

	if !gs.checkErrors {
		return
	}
	err := gs.gl.Call("getError")
	if err.Int() != NO_ERROR {
		panic(fmt.Sprintf("%s error: %v", name, err))
	}
}
