// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !wasm

package gls

// #include <stdlib.h>
// #include "glcorearb.h"
// #include "glapi.h"
import "C"

import (
	"fmt"
	"reflect"
	"unsafe"
)

// GLS encapsulates the state of an OpenGL context and contains
// methods to call OpenGL functions.
type GLS struct {
	stats       Stats             // statistics
	prog        *Program          // current active shader program
	programs    map[*Program]bool // shader programs cache
	checkErrors bool              // check openGL API errors flag

	// Cache OpenGL state to avoid making unnecessary API calls
	activeTexture  uint32  // cached last set active texture unit
	viewportX      int32   // cached last set viewport x
	viewportY      int32   // cached last set viewport y
	viewportWidth  int32   // cached last set viewport width
	viewportHeight int32   // cached last set viewport height
	lineWidth      float32 // cached last set line width
	sideView       int     // cached last set triangle side view mode
	frontFace      uint32  // cached last set glFrontFace value
	depthFunc      uint32  // cached last set depth function
	depthMask      int     // cached last set depth mask
	//stencilFunc
	stencilMask         uint32      // cached last set stencil mask
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
	gobuf               []byte      // conversion buffer with GO memory
	cbuf                []byte      // conversion buffer with C memory
}

// New creates and returns a new instance of a GLS object,
// which encapsulates the state of an OpenGL context.
// This should be called only after an active OpenGL context
// is established, such as by creating a new window.
func New() (*GLS, error) {

	gs := new(GLS)
	gs.reset()

	// Load OpenGL functions
	err := C.glapiLoad()
	if err != 0 {
		return nil, fmt.Errorf("Error loading OpenGL")
	}
	gs.setDefaultState()
	gs.checkErrors = true

	// Preallocate conversion buffers
	size := 1 * 1024
	gs.gobuf = make([]byte, size)
	p := C.malloc(C.size_t(size))
	gs.cbuf = (*[1 << 30]byte)(unsafe.Pointer(p))[:size:size]

	return gs, nil
}

// SetCheckErrors enables/disables checking for errors after the
// call of any OpenGL function. It is enabled by default but
// could be disabled after an application is stable to improve the performance.
func (gs *GLS) SetCheckErrors(enable bool) {

	if enable {
		C.glapiCheckError(1)
	} else {
		C.glapiCheckError(0)
	}
	gs.checkErrors = enable
}

// CheckErrors returns if error checking is enabled or not.
func (gs *GLS) CheckErrors() bool {

	return gs.checkErrors
}

// reset resets the internal state kept of the OpenGL
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

// setDefaultState is used internally to set the initial state of OpenGL
// for this context.
func (gs *GLS) setDefaultState() {

	gs.ClearColor(0, 0, 0, 1)
	gs.ClearDepth(1)
	gs.ClearStencil(0)
	gs.Enable(DEPTH_TEST)
	//gs.DepthMask(true)
	gs.DepthFunc(LEQUAL)
	gs.FrontFace(CCW)
	gs.CullFace(BACK)
	gs.Enable(CULL_FACE)
	gs.Enable(BLEND)
	gs.BlendEquation(FUNC_ADD)
	gs.BlendFunc(SRC_ALPHA, ONE_MINUS_SRC_ALPHA)

	gs.Enable(VERTEX_PROGRAM_POINT_SIZE)
	gs.Enable(PROGRAM_POINT_SIZE)
	gs.Enable(MULTISAMPLE)
	gs.Enable(POLYGON_OFFSET_FILL)
	gs.Enable(POLYGON_OFFSET_LINE)
	gs.Enable(POLYGON_OFFSET_POINT)
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
	C.glActiveTexture(C.GLenum(texture))
	gs.activeTexture = texture
}

// AttachShader attaches the specified shader object to the specified program object.
func (gs *GLS) AttachShader(program, shader uint32) {

	C.glAttachShader(C.GLuint(program), C.GLuint(shader))
}

// BindBuffer binds a buffer object to the specified buffer binding point.
func (gs *GLS) BindBuffer(target int, vbo uint32) {

	C.glBindBuffer(C.GLenum(target), C.GLuint(vbo))
}

// BindTexture lets you create or use a named texture.
func (gs *GLS) BindTexture(target int, tex uint32) {

	C.glBindTexture(C.GLenum(target), C.GLuint(tex))
}

// BindVertexArray binds the vertex array object.
func (gs *GLS) BindVertexArray(vao uint32) {

	C.glBindVertexArray(C.GLuint(vao))
}

// BlendEquation sets the blend equations for all draw buffers.
func (gs *GLS) BlendEquation(mode uint32) {

	if gs.blendEquation == mode {
		return
	}
	C.glBlendEquation(C.GLenum(mode))
	gs.blendEquation = mode
}

// BlendEquationSeparate sets the blend equations for all draw buffers
// allowing different equations for the RGB and alpha components.
func (gs *GLS) BlendEquationSeparate(modeRGB uint32, modeAlpha uint32) {

	if gs.blendEquationRGB == modeRGB && gs.blendEquationAlpha == modeAlpha {
		return
	}
	C.glBlendEquationSeparate(C.GLenum(modeRGB), C.GLenum(modeAlpha))
	gs.blendEquationRGB = modeRGB
	gs.blendEquationAlpha = modeAlpha
}

// BlendFunc defines the operation of blending for
// all draw buffers when blending is enabled.
func (gs *GLS) BlendFunc(sfactor, dfactor uint32) {

	if gs.blendSrc == sfactor && gs.blendDst == dfactor {
		return
	}
	C.glBlendFunc(C.GLenum(sfactor), C.GLenum(dfactor))
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
	C.glBlendFuncSeparate(C.GLenum(srcRGB), C.GLenum(dstRGB), C.GLenum(srcAlpha), C.GLenum(dstAlpha))
	gs.blendSrcRGB = srcRGB
	gs.blendDstRGB = dstRGB
	gs.blendSrcAlpha = srcAlpha
	gs.blendDstAlpha = dstAlpha
}

// BufferData creates a new data store for the buffer object currently
// bound to target, deleting any pre-existing data store.
func (gs *GLS) BufferData(target uint32, size int, data interface{}, usage uint32) {

	C.glBufferData(C.GLenum(target), C.GLsizeiptr(size), ptr(data), C.GLenum(usage))
}

// ClearColor specifies the red, green, blue, and alpha values
// used by glClear to clear the color buffers.
func (gs *GLS) ClearColor(r, g, b, a float32) {

	C.glClearColor(C.GLfloat(r), C.GLfloat(g), C.GLfloat(b), C.GLfloat(a))
}

// ClearDepth specifies the depth value used by Clear to clear the depth buffer.
func (gs *GLS) ClearDepth(v float32) {

	C.glClearDepth(C.GLclampd(v))
}

// ClearStencil specifies the index used by Clear to clear the stencil buffer.
func (gs *GLS) ClearStencil(v int32) {

	C.glClearStencil(C.GLint(v))
}

// Clear sets the bitplane area of the window to values previously
// selected by ClearColor, ClearDepth, and ClearStencil.
func (gs *GLS) Clear(mask uint) {

	C.glClear(C.GLbitfield(mask))
}

// CompileShader compiles the source code strings that
// have been stored in the specified shader object.
func (gs *GLS) CompileShader(shader uint32) {

	C.glCompileShader(C.GLuint(shader))
}

// CreateProgram creates an empty program object and returns
// a non-zero value by which it can be referenced.
func (gs *GLS) CreateProgram() uint32 {

	p := C.glCreateProgram()
	return uint32(p)
}

// CreateShader creates an empty shader object and returns
// a non-zero value by which it can be referenced.
func (gs *GLS) CreateShader(stype uint32) uint32 {

	h := C.glCreateShader(C.GLenum(stype))
	return uint32(h)
}

// DeleteBuffers deletes n​buffer objects named
// by the elements of the provided array.
func (gs *GLS) DeleteBuffers(bufs ...uint32) {

	C.glDeleteBuffers(C.GLsizei(len(bufs)), (*C.GLuint)(&bufs[0]))
	gs.stats.Buffers -= len(bufs)
}

// DeleteShader frees the memory and invalidates the name
// associated with the specified shader object.
func (gs *GLS) DeleteShader(shader uint32) {

	C.glDeleteShader(C.GLuint(shader))
}

// DeleteProgram frees the memory and invalidates the name
// associated with the specified program object.
func (gs *GLS) DeleteProgram(program uint32) {

	C.glDeleteProgram(C.GLuint(program))
}

// DeleteTextures deletes n​textures named
// by the elements of the provided array.
func (gs *GLS) DeleteTextures(tex ...uint32) {

	C.glDeleteTextures(C.GLsizei(len(tex)), (*C.GLuint)(&tex[0]))
	gs.stats.Textures -= len(tex)
}

// DeleteVertexArrays deletes n​vertex array objects named
// by the elements of the provided array.
func (gs *GLS) DeleteVertexArrays(vaos ...uint32) {

	C.glDeleteVertexArrays(C.GLsizei(len(vaos)), (*C.GLuint)(&vaos[0]))
	gs.stats.Vaos -= len(vaos)
}

// ReadPixels returns the current rendered image.
// x, y: specifies the window coordinates of the first pixel that is read from the frame buffer.
// width, height: specifies the dimensions of the pixel rectangle.
// format: specifies the format of the pixel data.
// format_type: specifies the data type of the pixel data.
// more information: http://docs.gl/gl3/glReadPixels
func (gs *GLS) ReadPixels(x, y, width, height, format, formatType int) []byte {
	size := uint32((width - x) * (height - y) * 4)
	C.glReadPixels(C.GLint(x), C.GLint(y), C.GLsizei(width), C.GLsizei(height), C.GLenum(format), C.GLenum(formatType), unsafe.Pointer(gs.gobufSize(size)))
	return gs.gobuf[:size]
}

// DepthFunc specifies the function used to compare each incoming pixel
// depth value with the depth value present in the depth buffer.
func (gs *GLS) DepthFunc(mode uint32) {

	if gs.depthFunc == mode {
		return
	}
	C.glDepthFunc(C.GLenum(mode))
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
	C.glDepthMask(bool2c(flag))
	if flag {
		gs.depthMask = intTrue
	} else {
		gs.depthMask = intFalse
	}
}

func (gs *GLS) StencilOp(fail, zfail, zpass uint32) {

	// TODO save state
	C.glStencilOp(C.GLenum(fail), C.GLenum(zfail), C.GLenum(zpass))
}

func (gs *GLS) StencilFunc(mode uint32, ref int32, mask uint32) {

	// TODO save state
	C.glStencilFunc(C.GLenum(mode), C.GLint(ref), C.GLuint(mask))
}

// TODO doc
// StencilMask enables or disables writing into the stencil buffer.
func (gs *GLS) StencilMask(mask uint32) {

	if gs.stencilMask == mask {
		return
	}
	C.glStencilMask(C.GLuint(mask))
	gs.stencilMask = mask
}

// DrawArrays renders primitives from array data.
func (gs *GLS) DrawArrays(mode uint32, first int32, count int32) {

	C.glDrawArrays(C.GLenum(mode), C.GLint(first), C.GLsizei(count))
	gs.stats.Drawcalls++
}

// DrawElements renders primitives from array data.
func (gs *GLS) DrawElements(mode uint32, count int32, itype uint32, start uint32) {

	C.glDrawElements(C.GLenum(mode), C.GLsizei(count), C.GLenum(itype), unsafe.Pointer(uintptr(start)))
	gs.stats.Drawcalls++
}

// Enable enables the specified capability.
func (gs *GLS) Enable(cap int) {

	if gs.capabilities[cap] == capEnabled {
		gs.stats.Caphits++
		return
	}
	C.glEnable(C.GLenum(cap))
	gs.capabilities[cap] = capEnabled
}

// Disable disables the specified capability.
func (gs *GLS) Disable(cap int) {

	if gs.capabilities[cap] == capDisabled {
		gs.stats.Caphits++
		return
	}
	C.glDisable(C.GLenum(cap))
	gs.capabilities[cap] = capDisabled
}

// EnableVertexAttribArray enables a generic vertex attribute array.
func (gs *GLS) EnableVertexAttribArray(index uint32) {

	C.glEnableVertexAttribArray(C.GLuint(index))
}

// CullFace specifies whether front- or back-facing facets can be culled.
func (gs *GLS) CullFace(mode uint32) {

	C.glCullFace(C.GLenum(mode))
}

// FrontFace defines front- and back-facing polygons.
func (gs *GLS) FrontFace(mode uint32) {

	if gs.frontFace == mode {
		return
	}
	C.glFrontFace(C.GLenum(mode))
	gs.frontFace = mode
}

// GenBuffer generates a ​buffer object name.
func (gs *GLS) GenBuffer() uint32 {

	var buf uint32
	C.glGenBuffers(1, (*C.GLuint)(&buf))
	gs.stats.Buffers++
	return buf
}

// GenerateMipmap generates mipmaps for the specified texture target.
func (gs *GLS) GenerateMipmap(target uint32) {

	C.glGenerateMipmap(C.GLenum(target))
}

// GenTexture generates a texture object name.
func (gs *GLS) GenTexture() uint32 {

	var tex uint32
	C.glGenTextures(1, (*C.GLuint)(&tex))
	gs.stats.Textures++
	return tex
}

// GenVertexArray generates a vertex array object name.
func (gs *GLS) GenVertexArray() uint32 {

	var vao uint32
	C.glGenVertexArrays(1, (*C.GLuint)(&vao))
	gs.stats.Vaos++
	return vao
}

// GetAttribLocation returns the location of the specified attribute variable.
func (gs *GLS) GetAttribLocation(program uint32, name string) int32 {

	loc := C.glGetAttribLocation(C.GLuint(program), gs.gobufStr(name))
	return int32(loc)
}

// GetProgramiv returns the specified parameter from the specified program object.
func (gs *GLS) GetProgramiv(program, pname uint32, params *int32) {

	C.glGetProgramiv(C.GLuint(program), C.GLenum(pname), (*C.GLint)(params))
}

// GetProgramInfoLog returns the information log for the specified program object.
func (gs *GLS) GetProgramInfoLog(program uint32) string {

	var length int32
	gs.GetProgramiv(program, INFO_LOG_LENGTH, &length)
	if length == 0 {
		return ""
	}
	C.glGetProgramInfoLog(C.GLuint(program), C.GLsizei(length), nil, gs.gobufSize(uint32(length)))
	return string(gs.gobuf[:length])
}

// GetShaderInfoLog returns the information log for the specified shader object.
func (gs *GLS) GetShaderInfoLog(shader uint32) string {

	var length int32
	gs.GetShaderiv(shader, INFO_LOG_LENGTH, &length)
	if length == 0 {
		return ""
	}
	C.glGetShaderInfoLog(C.GLuint(shader), C.GLsizei(length), nil, gs.gobufSize(uint32(length)))
	return string(gs.gobuf[:length])
}

// GetString returns a string describing the specified aspect of the current GL connection.
func (gs *GLS) GetString(name uint32) string {

	cs := C.glGetString(C.GLenum(name))
	return C.GoString((*C.char)(unsafe.Pointer(cs)))
}

// GetUniformLocation returns the location of a uniform variable for the specified program.
func (gs *GLS) GetUniformLocation(program uint32, name string) int32 {

	loc := C.glGetUniformLocation(C.GLuint(program), gs.gobufStr(name))
	return int32(loc)
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
	C.glLineWidth(C.GLfloat(width))
	gs.lineWidth = width
}

// LinkProgram links the specified program object.
func (gs *GLS) LinkProgram(program uint32) {

	C.glLinkProgram(C.GLuint(program))
}

// GetShaderiv returns the specified parameter from the specified shader object.
func (gs *GLS) GetShaderiv(shader, pname uint32, params *int32) {

	C.glGetShaderiv(C.GLuint(shader), C.GLenum(pname), (*C.GLint)(params))
}

// Scissor defines the scissor box rectangle in window coordinates.
func (gs *GLS) Scissor(x, y int32, width, height uint32) {

	C.glScissor(C.GLint(x), C.GLint(y), C.GLsizei(width), C.GLsizei(height))
}

// ShaderSource sets the source code for the specified shader object.
func (gs *GLS) ShaderSource(shader uint32, src string) {

	csource := gs.cbufStr(src)
	C.glShaderSource(C.GLuint(shader), 1, (**C.GLchar)(unsafe.Pointer(&csource)), nil)
}

// TexImage2D specifies a two-dimensional texture image.
func (gs *GLS) TexImage2D(target uint32, level int32, iformat int32, width int32, height int32, format uint32, itype uint32, data interface{}) {

	C.glTexImage2D(C.GLenum(target),
		C.GLint(level),
		C.GLint(iformat),
		C.GLsizei(width),
		C.GLsizei(height),
		C.GLint(0),
		C.GLenum(format),
		C.GLenum(itype),
		ptr(data))
}

// TexParameteri sets the specified texture parameter on the specified texture.
func (gs *GLS) TexParameteri(target uint32, pname uint32, param int32) {

	C.glTexParameteri(C.GLenum(target), C.GLenum(pname), C.GLint(param))
}

// PolygonMode controls the interpretation of polygons for rasterization.
func (gs *GLS) PolygonMode(face, mode uint32) {

	if gs.polygonModeFace == face && gs.polygonModeMode == mode {
		return
	}
	C.glPolygonMode(C.GLenum(face), C.GLenum(mode))
	gs.polygonModeFace = face
	gs.polygonModeMode = mode
}

// PolygonOffset sets the scale and units used to calculate depth values.
func (gs *GLS) PolygonOffset(factor float32, units float32) {

	if gs.polygonOffsetFactor == factor && gs.polygonOffsetUnits == units {
		return
	}
	C.glPolygonOffset(C.GLfloat(factor), C.GLfloat(units))
	gs.polygonOffsetFactor = factor
	gs.polygonOffsetUnits = units
}

// Uniform1i sets the value of an int uniform variable for the current program object.
func (gs *GLS) Uniform1i(location int32, v0 int32) {

	C.glUniform1i(C.GLint(location), C.GLint(v0))
	gs.stats.Unisets++
}

// Uniform1f sets the value of a float uniform variable for the current program object.
func (gs *GLS) Uniform1f(location int32, v0 float32) {

	C.glUniform1f(C.GLint(location), C.GLfloat(v0))
	gs.stats.Unisets++
}

// Uniform2f sets the value of a vec2 uniform variable for the current program object.
func (gs *GLS) Uniform2f(location int32, v0, v1 float32) {

	C.glUniform2f(C.GLint(location), C.GLfloat(v0), C.GLfloat(v1))
	gs.stats.Unisets++
}

// Uniform3f sets the value of a vec3 uniform variable for the current program object.
func (gs *GLS) Uniform3f(location int32, v0, v1, v2 float32) {

	C.glUniform3f(C.GLint(location), C.GLfloat(v0), C.GLfloat(v1), C.GLfloat(v2))
	gs.stats.Unisets++
}

// Uniform4f sets the value of a vec4 uniform variable for the current program object.
func (gs *GLS) Uniform4f(location int32, v0, v1, v2, v3 float32) {

	C.glUniform4f(C.GLint(location), C.GLfloat(v0), C.GLfloat(v1), C.GLfloat(v2), C.GLfloat(v3))
	gs.stats.Unisets++
}

// UniformMatrix3fv sets the value of one or many 3x3 float matrices for the current program object.
func (gs *GLS) UniformMatrix3fv(location int32, count int32, transpose bool, pm *float32) {

	C.glUniformMatrix3fv(C.GLint(location), C.GLsizei(count), bool2c(transpose), (*C.GLfloat)(pm))
	gs.stats.Unisets++
}

// UniformMatrix4fv sets the value of one or many 4x4 float matrices for the current program object.
func (gs *GLS) UniformMatrix4fv(location int32, count int32, transpose bool, pm *float32) {

	C.glUniformMatrix4fv(C.GLint(location), C.GLsizei(count), bool2c(transpose), (*C.GLfloat)(pm))
	gs.stats.Unisets++
}

// Uniform1fv sets the value of one or many float uniform variables for the current program object.
func (gs *GLS) Uniform1fv(location int32, count int32, v *float32) {

	C.glUniform1fv(C.GLint(location), C.GLsizei(count), (*C.GLfloat)(v))
	gs.stats.Unisets++
}

// Uniform2fv sets the value of one or many vec2 uniform variables for the current program object.
func (gs *GLS) Uniform2fv(location int32, count int32, v *float32) {

	C.glUniform2fv(C.GLint(location), C.GLsizei(count), (*C.GLfloat)(v))
	gs.stats.Unisets++
}

// Uniform3fv sets the value of one or many vec3 uniform variables for the current program object.
func (gs *GLS) Uniform3fv(location int32, count int32, v *float32) {

	C.glUniform3fv(C.GLint(location), C.GLsizei(count), (*C.GLfloat)(v))
	gs.stats.Unisets++
}

// Uniform4fv sets the value of one or many vec4 uniform variables for the current program object.
func (gs *GLS) Uniform4fv(location int32, count int32, v *float32) {

	C.glUniform4fv(C.GLint(location), C.GLsizei(count), (*C.GLfloat)(v))
	gs.stats.Unisets++
}

// VertexAttribPointer defines an array of generic vertex attribute data.
func (gs *GLS) VertexAttribPointer(index uint32, size int32, xtype uint32, normalized bool, stride int32, offset uint32) {

	C.glVertexAttribPointer(C.GLuint(index), C.GLint(size), C.GLenum(xtype), bool2c(normalized), C.GLsizei(stride), unsafe.Pointer(uintptr(offset)))
}

// Viewport sets the viewport.
func (gs *GLS) Viewport(x, y, width, height int32) {

	C.glViewport(C.GLint(x), C.GLint(y), C.GLsizei(width), C.GLsizei(height))
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
	C.glUseProgram(C.GLuint(prog.handle))
	gs.prog = prog

	// Inserts program in cache if not already there.
	if !gs.programs[prog] {
		gs.programs[prog] = true
		log.Debug("New Program activated. Total: %d", len(gs.programs))
	}
}

// Ptr takes a slice or pointer (to a singular scalar value or the first
// element of an array or slice) and returns its GL-compatible address.
//
// For example:
//
// 	var data []uint8
// 	...
// 	gl.TexImage2D(gl.TEXTURE_2D, ..., gl.UNSIGNED_BYTE, gl.Ptr(&data[0]))
func ptr(data interface{}) unsafe.Pointer {
	if data == nil {
		return unsafe.Pointer(nil)
	}
	var addr unsafe.Pointer
	v := reflect.ValueOf(data)
	switch v.Type().Kind() {
	case reflect.Ptr:
		e := v.Elem()
		switch e.Kind() {
		case
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			addr = unsafe.Pointer(e.UnsafeAddr())
		default:
			panic(fmt.Errorf("unsupported pointer to type %s; must be a slice or pointer to a singular scalar value or the first element of an array or slice", e.Kind()))
		}
	case reflect.Uintptr:
		addr = unsafe.Pointer(v.Pointer())
	case reflect.Slice:
		addr = unsafe.Pointer(v.Index(0).UnsafeAddr())
	default:
		panic(fmt.Errorf("unsupported type %s; must be a slice or pointer to a singular scalar value or the first element of an array or slice", v.Type()))
	}
	return addr
}

// bool2c convert a Go bool to C.GLboolean
func bool2c(b bool) C.GLboolean {

	if b {
		return C.GLboolean(1)
	}
	return C.GLboolean(0)
}

// gobufSize returns a pointer to static buffer with the specified size not including the terminator.
// If there is available space, there is no memory allocation.
func (gs *GLS) gobufSize(size uint32) *C.GLchar {

	if size+1 > uint32(len(gs.gobuf)) {
		gs.gobuf = make([]byte, size+1)
	}
	return (*C.GLchar)(unsafe.Pointer(&gs.gobuf[0]))
}

// gobufStr converts a Go String to a C string by copying it to a static buffer
// and returning a pointer to the start of the buffer.
// If there is available space, there is no memory allocation.
func (gs *GLS) gobufStr(s string) *C.GLchar {

	p := gs.gobufSize(uint32(len(s) + 1))
	copy(gs.gobuf, s)
	gs.gobuf[len(s)] = 0
	return p
}

// cbufSize returns a pointer to static buffer with C memory
// If there is available space, there is no memory allocation.
func (gs *GLS) cbufSize(size uint32) *C.GLchar {

	if size > uint32(len(gs.cbuf)) {
		if len(gs.cbuf) > 0 {
			C.free(unsafe.Pointer(&gs.cbuf[0]))
		}
		p := C.malloc(C.size_t(size))
		gs.cbuf = (*[1 << 30]byte)(unsafe.Pointer(p))[:size:size]
	}
	return (*C.GLchar)(unsafe.Pointer(&gs.cbuf[0]))
}

// cbufStr converts a Go String to a C string by copying it to a single pre-allocated buffer
// using C memory and returning a pointer to the start of the buffer.
// If there is available space, there is no memory allocation.
func (gs *GLS) cbufStr(s string) *C.GLchar {

	p := gs.cbufSize(uint32(len(s) + 1))
	copy(gs.cbuf, s)
	gs.cbuf[len(s)] = 0
	return p
}
