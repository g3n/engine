package gls

// // Platform build flags
// #cgo freebsd CFLAGS: -DGL_GLEXT_PROTOTYPES
// #cgo freebsd LDFLAGS: -ldl -lGL
//
// #cgo linux CFLAGS: -DGL_GLEXT_PROTOTYPES
// #cgo linux LDFLAGS: -ldl -lGL
//
// #cgo windows CFLAGS: -DGL_GEXT_PROTOTYPES
// #cgo windows LDFLAGS: -lopengl32
//
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

type GLS struct {
	stats               Stats             // statistics
	Prog                *Program          // Current active program
	programs            map[*Program]bool // Programs cache
	checkErrors         bool              // check openGL API errors flag
	viewportX           int32             // cached last set viewport x
	viewportY           int32             // cached last set viewport y
	viewportWidth       int32             // cached last set viewport width
	viewportHeight      int32             // cached last set viewport height
	lineWidth           float32           // cached last set line width
	sideView            int               // cached last set triangle side view mode
	depthFunc           uint32            // cached last set depth function
	depthMask           int               // cached last set depth mask
	capabilities        map[int]int       // cached capabilities (Enable/Disable)
	blendEquation       uint32
	blendSrc            uint32
	blendDst            uint32
	blendEquationRGB    uint32
	blendEquationAlpha  uint32
	blendSrcRGB         uint32
	blendSrcAlpha       uint32
	blendDstRGB         uint32
	blendDstAlpha       uint32
	polygonOffsetFactor float32
	polygonOffsetUnits  float32
	logBuf              []byte // pre allocated buffer for program/shader logs
}

type Stats struct {
	Vaos     int // Number of Vertex Array Objects
	Vbos     int // Number of Vertex Buffer Objects
	Textures int // Number of Textures
	// Cummulative fields
	Caphits   int // Number of hits for Enable/Disable
	Unisets   int // Number of uniform sets
	Drawcalls int // Number of draw calls
}

const (
	capUndef    = 0
	capDisabled = 1
	capEnabled  = 2
	uintUndef   = math.MaxUint32
	intFalse    = 0
	intTrue     = 1
	maxLogBuf   = 32 * 1024
)

// Polygon side view.
const (
	FrontSide = iota + 1
	BackSide
	DoubleSide
)

// New creates and returns a new instance of an GLS object
// which encapsulates the state of an OpenGL context
// This should be called only after an active OpenGL context
// was established, such as by creating a new window.
func New() (*GLS, error) {

	gs := new(GLS)
	gs.Reset()

	// Initialize GL
	err := C.glapiLoad()
	if err != 0 {
		return nil, fmt.Errorf("Error loading OpenGL")
	}
	gs.SetDefaultState()
	gs.checkErrors = true
	gs.logBuf = make([]byte, maxLogBuf)
	return gs, nil
}

// SetCheckErrors enables/disables checking for errors after the
// call of any OpenGL function. It is enabled by default but
// could be disabled after an application is stable to improve the performance.
func (gs *GLS) SetCheckErrors(enable bool) {

	if enable {
		C.glapiCheckError(1)
	} else {
		C.glapiCheckError(1)
	}
	gs.checkErrors = enable
}

// ChecksErrors returns if error checking is enabled or not.
func (gs *GLS) CheckErrors() bool {

	return gs.checkErrors
}

// Reset resets the internal state kept of the OpenGL
func (gs *GLS) Reset() {

	gs.lineWidth = 0.0
	gs.sideView = uintUndef
	gs.depthFunc = 0
	gs.depthMask = uintUndef
	gs.capabilities = make(map[int]int)
	gs.programs = make(map[*Program]bool)
	gs.Prog = nil

	gs.blendEquation = uintUndef
	gs.blendSrc = uintUndef
	gs.blendDst = uintUndef
	gs.blendEquationRGB = 0
	gs.blendEquationAlpha = 0
	gs.blendSrcRGB = uintUndef
	gs.blendSrcAlpha = uintUndef
	gs.blendDstRGB = uintUndef
	gs.blendDstAlpha = uintUndef
	gs.polygonOffsetFactor = -1
	gs.polygonOffsetUnits = -1
}

func (gs *GLS) SetDefaultState() {

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

func (gs *GLS) ActiveTexture(texture uint32) {

	C.glActiveTexture(C.GLenum(texture))
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

func (gs *GLS) DeleteBuffers(vbos ...uint32) {

	C.glDeleteBuffers(C.GLsizei(len(vbos)), (*C.GLuint)(&vbos[0]))
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

func (gs *GLS) DrawElements(mode uint32, count int32, itype uint32, start uint32) {

	C.glDrawElements(C.GLenum(mode), C.GLsizei(count), C.GLenum(itype), ptrOffset(int(start)))
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

	C.glFrontFace(C.GLenum(mode))
}

func (gs *GLS) GenBuffer() uint32 {

	var buf uint32
	C.glGenBuffers(1, (*C.GLuint)(&buf))
	gs.stats.Vbos++
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

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	res := C.glGetAttribLocation(C.GLuint(program), (*C.GLchar)(cname))
	return int32(res)
}

func (gs *GLS) GetProgramiv(program, pname uint32, params *int32) {

	C.glGetProgramiv(C.GLuint(program), C.GLenum(pname), (*C.GLint)(params))
}

func (gs *GLS) GetProgramInfoLog(program uint32) string {

	// Get length of program info log buffer
	var logLength int32
	gs.GetProgramiv(program, INFO_LOG_LENGTH, &logLength)
	if logLength == 0 {
		return ""
	}
	C.glGetProgramInfoLog(C.GLuint(program), maxLogBuf, nil, (*C.GLchar)(unsafe.Pointer(&gs.logBuf[0])))
	return string(gs.logBuf)
}

func (gs *GLS) GetShaderInfoLog(shader uint32) string {

	// Get length of shaderinfo log buffer
	var logLength int32
	gs.GetShaderiv(shader, INFO_LOG_LENGTH, &logLength)
	if logLength == 0 {
		return ""
	}
	buf := make([]byte, logLength)
	C.glGetShaderInfoLog(C.GLuint(shader), C.GLsizei(logLength), nil, (*C.GLchar)(unsafe.Pointer(&gs.logBuf[0])))
	return string(buf)
}

func (gs *GLS) GetString(name uint32) string {

	cstr := C.glGetString(C.GLenum(name))
	return goStr((*uint8)(cstr))
}

func (gs *GLS) GetUniformLocation(program uint32, name string) int32 {

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	loc := C.glGetUniformLocation(C.GLuint(program), (*C.GLchar)(cname))
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

func (gs *GLS) SetDepthTest(mode bool) {

	if mode {
		gs.Enable(DEPTH_TEST)
	} else {
		gs.Disable(DEPTH_TEST)
	}
}

func (gs *GLS) SetSideView(mode int) {

	if gs.sideView == mode {
		return
	}
	switch mode {
	// Default: show only the front size
	case FrontSide:
		gs.Enable(CULL_FACE)
		C.glFrontFace(CCW)
	// Show only the back side
	case BackSide:
		gs.Enable(CULL_FACE)
		C.glFrontFace(CW)
	// Show both sides
	case DoubleSide:
		gs.Disable(CULL_FACE)
	default:
		panic("SetSideView() invalid mode")
	}
	gs.sideView = mode
}

func (gs *GLS) GetShaderiv(shader, pname uint32, params *int32) {

	C.glGetShaderiv(C.GLuint(shader), C.GLenum(pname), (*C.GLint)(params))
}

func (gs *GLS) ShaderSource(shader uint32, src string) {

	csource := C.CString(src)
	defer C.free(unsafe.Pointer(csource))
	C.glShaderSource(C.GLuint(shader), 1, (**C.GLchar)(unsafe.Pointer(csource)), nil)
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

func (gs *GLS) TexStorage2D(target int, levels int, iformat int, width, height int) {

	C.glTexStorage2D(C.GLenum(target), C.GLsizei(levels), C.GLenum(iformat), C.GLsizei(width), C.GLsizei(height))
}

func (gs *GLS) TexParameteri(target uint32, pname uint32, param int32) {

	C.glTexParameteri(C.GLenum(target), C.GLenum(pname), C.GLint(param))
}

func (gs *GLS) PolygonMode(face, mode int) {

	C.glPolygonMode(C.GLenum(face), C.GLenum(mode))
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

func (gs *GLS) Uniform2fv(location int32, count int32, v []float32) {

	C.glUniform2fv(C.GLint(location), C.GLsizei(count), (*C.GLfloat)(&v[0]))
	gs.stats.Unisets++
}

func (gs *GLS) Uniform3fv(location int32, count int32, v []float32) {

	C.glUniform3fv(C.GLint(location), C.GLsizei(count), (*C.GLfloat)(&v[0]))
	gs.stats.Unisets++
}

func (gs *GLS) Uniform4fv(location int32, count int32, v []float32) {

	C.glUniform4fv(C.GLint(location), C.GLsizei(count), (*C.GLfloat)(&v[0]))
	gs.stats.Unisets++
}

func (gs *GLS) VertexAttribPointer(index uint32, size int32, xtype uint32, normalized bool, stride int32, offset uint32) {

	C.glVertexAttribPointer(C.GLuint(index), C.GLint(size), C.GLenum(xtype), bool2c(normalized), C.GLsizei(stride), ptrOffset(int(offset)))
}

func (gs *GLS) Viewport(x, y, width, height int32) {

	C.glViewport(C.GLint(x), C.GLint(y), C.GLsizei(width), C.GLsizei(height))
	gs.viewportX = x
	gs.viewportY = y
	gs.viewportWidth = width
	gs.viewportHeight = height
}

// Use set this program as the current program.
//func (gs *GLS) UseProgram(prog *Program) {
//
//	if prog.handle == 0 {
//		panic("Invalid program")
//	}
//	C.glUseProgram(prog.handle)
//	gs.Prog = prog
//
//	// Inserts program in cache if not already there.
//	if !gs.programs[prog] {
//		gs.programs[prog] = true
//		log.Debug("New Program activated. Total: %d", len(gs.programs))
//	}
//}

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

// ptrOffset takes a pointer offset and returns a GL-compatible pointer.
// Useful for functions such as glVertexAttribPointer that take pointer
// parameters indicating an offset rather than an absolute memory address.
func ptrOffset(offset int) unsafe.Pointer {

	return unsafe.Pointer(uintptr(offset))
}

//// Str takes a null-terminated Go string and returns its GL-compatible address.
//// This function reaches into Go string storage in an unsafe way so the caller
//// must ensure the string is not garbage collected.
//func Str(str string) *uint8 {
//	if !strings.HasSuffix(str, "\x00") {
//		panic("str argument missing null terminator: " + str)
//	}
//	header := (*reflect.StringHeader)(unsafe.Pointer(&str))
//	return (*uint8)(unsafe.Pointer(header.Data))
//}

// goStr takes a null-terminated string returned by OpenGL and constructs a
// corresponding Go string.
func goStr(cstr *uint8) string {

	return C.GoString((*C.char)(unsafe.Pointer(cstr)))
}
