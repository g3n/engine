// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package gls

// #include <stdlib.h>
// #include "glcorearb.h"
// #include "glapi.h"
import "C"

import (
	"fmt"
	"math"
	"reflect"
	"unsafe"
)

// GLS encapsulates the state of an OpenGL context and contains
// methods to call OpenGL functions.
type GLS struct {
	stats               Stats             // statistics
	prog                *Program          // current active shader program
	programs            map[*Program]bool // shader programs cache
	checkErrors         bool              // check openGL API errors flag
	activeTexture       uint32            // cached last set active texture unit
	viewportX           int32             // cached last set viewport x
	viewportY           int32             // cached last set viewport y
	viewportWidth       int32             // cached last set viewport width
	viewportHeight      int32             // cached last set viewport height
	lineWidth           float32           // cached last set line width
	sideView            int               // cached last set triangle side view mode
	frontFace           uint32            // cached last set glFrontFace value
	depthFunc           uint32            // cached last set depth function
	depthMask           int               // cached last set depth mask
	capabilities        map[int]int       // cached capabilities (Enable/Disable)
	blendEquation       uint32            // cached last set blend equation value
	blendSrc            uint32            // cached last set blend src value
	blendDst            uint32            // cached last set blend equation destination value
	blendEquationRGB    uint32            // cached last set blend equation rgb value
	blendEquationAlpha  uint32            // cached last set blend equation alpha value
	blendSrcRGB         uint32            // cached last set blend src rgb
	blendSrcAlpha       uint32            // cached last set blend src alpha value
	blendDstRGB         uint32            // cached last set blend destination rgb value
	blendDstAlpha       uint32            // cached last set blend destination alpha value
	polygonModeFace     uint32            // cached last set polygon mode face
	polygonModeMode     uint32            // cached last set polygon mode mode
	polygonOffsetFactor float32           // cached last set polygon offset factor
	polygonOffsetUnits  float32           // cached last set polygon offset units
	gobuf               []byte            // conversion buffer with GO memory
	cbuf                []byte            // conversion buffer with C memory
}

// Stats contains counters of OpenGL resources being used as well
// the cummulative numbers of some OpenGL calls for performance evaluation.
type Stats struct {
	Shaders    int    // Current number of shader programs
	Vaos       int    // Number of Vertex Array Objects
	Buffers    int    // Number of Buffer Objects
	Textures   int    // Number of Textures
	Caphits    uint64 // Cummulative number of hits for Enable/Disable
	UnilocHits uint64 // Cummulative number of uniform location cache hits
	UnilocMiss uint64 // Cummulative number of uniform location cache misses
	Unisets    uint64 // Cummulative number of uniform sets
	Drawcalls  uint64 // Cummulative number of draw calls
}

// Polygon side view.
const (
	FrontSide = iota + 1
	BackSide
	DoubleSide
)

const (
	capUndef    = 0
	capDisabled = 1
	capEnabled  = 2
	uintUndef   = math.MaxUint32
	intFalse    = 0
	intTrue     = 1
)

// New creates and returns a new instance of an GLS object
// which encapsulates the state of an OpenGL context
// This should be called only after an active OpenGL context
// was established, such as by creating a new window.
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

	// Preallocates conversion buffers
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

// ChecksErrors returns if error checking is enabled or not.
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

	C.glClearColor(0, 0, 0, 1)
	C.glClearDepth(1)
	C.glClearStencil(0)
	gs.Enable(DEPTH_TEST)
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

func (gs *GLS) ActiveTexture(texture uint32) {

	if gs.activeTexture == texture {
		return
	}
	C.glActiveTexture(C.GLenum(texture))
	gs.activeTexture = texture
}

func (gs *GLS) AttachShader(program, shader uint32) {

	C.glAttachShader(C.GLuint(program), C.GLuint(shader))
}

func (gs *GLS) BindBuffer(target int, vbo uint32) {

	C.glBindBuffer(C.GLenum(target), C.GLuint(vbo))
}

func (gs *GLS) BindTexture(target int, tex uint32) {

	C.glBindTexture(C.GLenum(target), C.GLuint(tex))
}

func (gs *GLS) BindVertexArray(vao uint32) {

	C.glBindVertexArray(C.GLuint(vao))
}

func (gs *GLS) BlendEquation(mode uint32) {

	if gs.blendEquation == mode {
		return
	}
	C.glBlendEquation(C.GLenum(mode))
	gs.blendEquation = mode
}

func (gs *GLS) BlendEquationSeparate(modeRGB uint32, modeAlpha uint32) {

	if gs.blendEquationRGB == modeRGB && gs.blendEquationAlpha == modeAlpha {
		return
	}
	C.glBlendEquationSeparate(C.GLenum(modeRGB), C.GLenum(modeAlpha))
	gs.blendEquationRGB = modeRGB
	gs.blendEquationAlpha = modeAlpha
}

func (gs *GLS) BlendFunc(sfactor, dfactor uint32) {

	if gs.blendSrc == sfactor && gs.blendDst == dfactor {
		return
	}
	C.glBlendFunc(C.GLenum(sfactor), C.GLenum(dfactor))
	gs.blendSrc = sfactor
	gs.blendDst = dfactor
}

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

func (gs *GLS) BufferData(target uint32, size int, data interface{}, usage uint32) {

	C.glBufferData(C.GLenum(target), C.GLsizeiptr(size), ptr(data), C.GLenum(usage))
}

func (gs *GLS) ClearColor(r, g, b, a float32) {

	C.glClearColor(C.GLfloat(r), C.GLfloat(g), C.GLfloat(b), C.GLfloat(a))
}

func (gs *GLS) Clear(mask uint) {

	C.glClear(C.GLbitfield(mask))
}

func (gs *GLS) CompileShader(shader uint32) {

	C.glCompileShader(C.GLuint(shader))
}

func (gs *GLS) CreateProgram() uint32 {

	p := C.glCreateProgram()
	return uint32(p)
}

func (gs *GLS) CreateShader(stype uint32) uint32 {

	h := C.glCreateShader(C.GLenum(stype))
	return uint32(h)
}

func (gs *GLS) DeleteBuffers(bufs ...uint32) {

	C.glDeleteBuffers(C.GLsizei(len(bufs)), (*C.GLuint)(&bufs[0]))
	gs.stats.Buffers -= len(bufs)
}

func (gs *GLS) DeleteShader(shader uint32) {

	C.glDeleteShader(C.GLuint(shader))
}

func (gs *GLS) DeleteProgram(program uint32) {

	C.glDeleteProgram(C.GLuint(program))
}

func (gs *GLS) DeleteTextures(tex ...uint32) {

	C.glDeleteTextures(C.GLsizei(len(tex)), (*C.GLuint)(&tex[0]))
	gs.stats.Textures -= len(tex)
}

func (gs *GLS) DeleteVertexArrays(vaos ...uint32) {

	C.glDeleteVertexArrays(C.GLsizei(len(vaos)), (*C.GLuint)(&vaos[0]))
	gs.stats.Vaos -= len(vaos)
}

func (gs *GLS) DepthFunc(mode uint32) {

	if gs.depthFunc == mode {
		return
	}
	C.glDepthFunc(C.GLenum(mode))
	gs.depthFunc = mode
}

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

func (gs *GLS) DrawArrays(mode uint32, first int32, count int32) {

	C.glDrawArrays(C.GLenum(mode), C.GLint(first), C.GLsizei(count))
	gs.stats.Drawcalls++
}

func (gs *GLS) DrawBuffer(mode uint32) {

	C.glDrawBuffer(C.GLenum(mode))
}

func (gs *GLS) DrawElements(mode uint32, count int32, itype uint32, start uint32) {

	C.glDrawElements(C.GLenum(mode), C.GLsizei(count), C.GLenum(itype), unsafe.Pointer(uintptr(start)))
	gs.stats.Drawcalls++
}

func (gs *GLS) Enable(cap int) {

	if gs.capabilities[cap] == capEnabled {
		gs.stats.Caphits++
		return
	}
	C.glEnable(C.GLenum(cap))
	gs.capabilities[cap] = capEnabled
}

func (gs *GLS) EnableVertexAttribArray(index uint32) {

	C.glEnableVertexAttribArray(C.GLuint(index))
}

func (gs *GLS) Disable(cap int) {

	if gs.capabilities[cap] == capDisabled {
		gs.stats.Caphits++
		return
	}
	C.glDisable(C.GLenum(cap))
	gs.capabilities[cap] = capDisabled
}

func (gs *GLS) CullFace(mode uint32) {

	C.glCullFace(C.GLenum(mode))
}

func (gs *GLS) FrontFace(mode uint32) {

	if gs.frontFace == mode {
		return
	}
	C.glFrontFace(C.GLenum(mode))
	gs.frontFace = mode
}

func (gs *GLS) GenBuffer() uint32 {

	var buf uint32
	C.glGenBuffers(1, (*C.GLuint)(&buf))
	gs.stats.Buffers++
	return buf
}

func (gs *GLS) GenerateMipmap(target uint32) {

	C.glGenerateMipmap(C.GLenum(target))
}

func (gs *GLS) GenTexture() uint32 {

	var tex uint32
	C.glGenTextures(1, (*C.GLuint)(&tex))
	gs.stats.Textures++
	return tex
}

func (gs *GLS) GenVertexArray() uint32 {

	var vao uint32
	C.glGenVertexArrays(1, (*C.GLuint)(&vao))
	gs.stats.Vaos++
	return vao
}

func (gs *GLS) GetAttribLocation(program uint32, name string) int32 {

	loc := C.glGetAttribLocation(C.GLuint(program), gs.gobufStr(name))
	return int32(loc)
}

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

func (gs *GLS) GetString(name uint32) string {

	cs := C.glGetString(C.GLenum(name))
	return C.GoString((*C.char)(unsafe.Pointer(cs)))
}

// GetUniformLocation returns the location of a uniform variable for the specified program.
func (gs *GLS) GetUniformLocation(program uint32, name string) int32 {

	loc := C.glGetUniformLocation(C.GLuint(program), gs.gobufStr(name))
	return int32(loc)
}

func (gs *GLS) GetViewport() (x, y, width, height int32) {

	return gs.viewportX, gs.viewportY, gs.viewportWidth, gs.viewportHeight
}

func (gs *GLS) LineWidth(width float32) {

	if gs.lineWidth == width {
		return
	}
	C.glLineWidth(C.GLfloat(width))
	gs.lineWidth = width
}

func (gs *GLS) LinkProgram(program uint32) {

	C.glLinkProgram(C.GLuint(program))
}

func (gs *GLS) GetShaderiv(shader, pname uint32, params *int32) {

	C.glGetShaderiv(C.GLuint(shader), C.GLenum(pname), (*C.GLint)(params))
}

func (gs *GLS) Scissor(x, y int32, width, height uint32) {

	C.glScissor(C.GLint(x), C.GLint(y), C.GLsizei(width), C.GLsizei(height))
}

func (gs *GLS) ShaderSource(shader uint32, src string) {

	csource := gs.cbufStr(src)
	C.glShaderSource(C.GLuint(shader), 1, (**C.GLchar)(unsafe.Pointer(&csource)), nil)
}

func (gs *GLS) TexImage2D(target uint32, level int32, iformat int32, width int32, height int32, border int32, format uint32, itype uint32, data interface{}) {

	C.glTexImage2D(C.GLenum(target),
		C.GLint(level),
		C.GLint(iformat),
		C.GLsizei(width),
		C.GLsizei(height),
		C.GLint(border),
		C.GLenum(format),
		C.GLenum(itype),
		ptr(data))
}

func (gs *GLS) TexParameteri(target uint32, pname uint32, param int32) {

	C.glTexParameteri(C.GLenum(target), C.GLenum(pname), C.GLint(param))
}

func (gs *GLS) PolygonMode(face, mode uint32) {

	if gs.polygonModeFace == face && gs.polygonModeMode == mode {
		return
	}
	C.glPolygonMode(C.GLenum(face), C.GLenum(mode))
	gs.polygonModeFace = face
	gs.polygonModeMode = mode
}

func (gs *GLS) PolygonOffset(factor float32, units float32) {

	if gs.polygonOffsetFactor == factor && gs.polygonOffsetUnits == units {
		return
	}
	C.glPolygonOffset(C.GLfloat(factor), C.GLfloat(units))
	gs.polygonOffsetFactor = factor
	gs.polygonOffsetUnits = units
}

func (gs *GLS) Uniform1i(location int32, v0 int32) {

	C.glUniform1i(C.GLint(location), C.GLint(v0))
	gs.stats.Unisets++
}

func (gs *GLS) Uniform1f(location int32, v0 float32) {

	C.glUniform1f(C.GLint(location), C.GLfloat(v0))
	gs.stats.Unisets++
}

func (gs *GLS) Uniform2f(location int32, v0, v1 float32) {

	C.glUniform2f(C.GLint(location), C.GLfloat(v0), C.GLfloat(v1))
	gs.stats.Unisets++
}

func (gs *GLS) Uniform3f(location int32, v0, v1, v2 float32) {

	C.glUniform3f(C.GLint(location), C.GLfloat(v0), C.GLfloat(v1), C.GLfloat(v2))
	gs.stats.Unisets++
}

func (gs *GLS) Uniform4f(location int32, v0, v1, v2, v3 float32) {

	C.glUniform4f(C.GLint(location), C.GLfloat(v0), C.GLfloat(v1), C.GLfloat(v2), C.GLfloat(v3))
	gs.stats.Unisets++
}

func (gs *GLS) UniformMatrix3fv(location int32, count int32, transpose bool, pm *float32) {

	C.glUniformMatrix3fv(C.GLint(location), C.GLsizei(count), bool2c(transpose), (*C.GLfloat)(pm))
	gs.stats.Unisets++
}

func (gs *GLS) UniformMatrix4fv(location int32, count int32, transpose bool, pm *float32) {

	C.glUniformMatrix4fv(C.GLint(location), C.GLsizei(count), bool2c(transpose), (*C.GLfloat)(pm))
	gs.stats.Unisets++
}

func (gs *GLS) Uniform1fv(location int32, count int32, v []float32) {

	C.glUniform1fv(C.GLint(location), C.GLsizei(count), (*C.GLfloat)(&v[0]))
	gs.stats.Unisets++
}

func (gs *GLS) Uniform2fv(location int32, count int32, v *float32) {

	C.glUniform2fv(C.GLint(location), C.GLsizei(count), (*C.GLfloat)(v))
	gs.stats.Unisets++
}

func (gs *GLS) Uniform2fvUP(location int32, count int32, v unsafe.Pointer) {

	C.glUniform2fv(C.GLint(location), C.GLsizei(count), (*C.GLfloat)(v))
	gs.stats.Unisets++
}

func (gs *GLS) Uniform3fv(location int32, count int32, v *float32) {

	C.glUniform3fv(C.GLint(location), C.GLsizei(count), (*C.GLfloat)(v))
	gs.stats.Unisets++
}

func (gs *GLS) Uniform3fvUP(location int32, count int32, v unsafe.Pointer) {

	C.glUniform3fv(C.GLint(location), C.GLsizei(count), (*C.GLfloat)(v))
	gs.stats.Unisets++
}

func (gs *GLS) Uniform4fv(location int32, count int32, v []float32) {

	C.glUniform4fv(C.GLint(location), C.GLsizei(count), (*C.GLfloat)(&v[0]))
	gs.stats.Unisets++
}

func (gs *GLS) Uniform4fvUP(location int32, count int32, v unsafe.Pointer) {

	C.glUniform4fv(C.GLint(location), C.GLsizei(count), (*C.GLfloat)(v))
	gs.stats.Unisets++
}

func (gs *GLS) VertexAttribPointer(index uint32, size int32, xtype uint32, normalized bool, stride int32, offset uint32) {

	C.glVertexAttribPointer(C.GLuint(index), C.GLint(size), C.GLenum(xtype), bool2c(normalized), C.GLsizei(stride), unsafe.Pointer(uintptr(offset)))
}

func (gs *GLS) Viewport(x, y, width, height int32) {

	C.glViewport(C.GLint(x), C.GLint(y), C.GLsizei(width), C.GLsizei(height))
	gs.viewportX = x
	gs.viewportY = y
	gs.viewportWidth = width
	gs.viewportHeight = height
}

// Use set this program as the current program.
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
