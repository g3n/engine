// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package gls allows access to the OpenGL functions.
package gls

import (
	"github.com/g3n/engine/util/logger"
	"github.com/go-gl/gl/v3.3-core/gl"
	"math"
)

// GLS allows access to the OpenGL functions and keeps state to
// minimize functions calling.
// It also keeps some statistics of some OpenGL objects currently allocated
type GLS struct {
	// Statistics
	Stats struct {
		Vaos     int // Number of Vertex Array Objects
		Vbos     int // Number of Vertex Buffer Objects
		Textures int // Number of Textures
	}
	Prog               *Program          // Current active program
	programs           map[*Program]bool // Programs cache
	checkErrors        bool              // Check openGL API errors flag
	viewportX          int32
	viewportY          int32
	viewportWidth      int32
	viewportHeight     int32
	lineWidth          float32
	sideView           int
	depthFunc          uint32
	depthMask          int
	capabilities       map[int]int
	blendEquation      uint32
	blendSrc           uint32
	blendDst           uint32
	blendEquationRGB   uint32
	blendEquationAlpha uint32
	blendSrcRGB        uint32
	blendSrcAlpha      uint32
	blendDstRGB        uint32
	blendDstAlpha      uint32
}

const (
	capUndef    = 0
	capDisabled = 1
	capEnabled  = 2
)
const (
	uintUndef = math.MaxUint32
	intFalse  = 0
	intTrue   = 1
)

// Polygon side view.
const (
	FrontSide = iota + 1
	BackSide
	DoubleSide
)

// Package logger
var log = logger.New("GLS", logger.Default)

// New creates and returns a new instance of an GLS object
// which encapsulates the state of an OpenGL context
// This should be called only after an active OpenGL context
// was established, such as by creating a new window.
func New() (*GLS, error) {

	gs := new(GLS)
	gs.Reset()

	// Initialize GL
	err := gl.Init()
	if err != nil {
		return nil, err
	}
	gs.SetDefaultState()
	gs.checkErrors = true
	return gs, nil
}

// SetCheckErrors enables/disables checking for errors after the
// call of any OpenGL function. It is enabled by default but
// could be disabled after an application is stable to improve the performance.
func (gs *GLS) SetCheckErrors(enable bool) {

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
}

func (gs *GLS) SetDefaultState() {

	gl.ClearColor(0, 0, 0, 1)
	gl.ClearDepth(1)
	gl.ClearStencil(0)
	gs.Enable(gl.DEPTH_TEST)
	gs.DepthFunc(gl.LEQUAL)
	gl.FrontFace(gl.CCW)
	gl.CullFace(gl.BACK)
	gs.Enable(gl.CULL_FACE)
	gs.Enable(gl.BLEND)
	gs.BlendEquation(gl.FUNC_ADD)
	gs.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gs.Enable(gl.VERTEX_PROGRAM_POINT_SIZE)
	gs.Enable(gl.PROGRAM_POINT_SIZE)
	gs.Enable(gl.MULTISAMPLE)
}

func (gs *GLS) ActiveTexture(texture uint32) {

	gl.ActiveTexture(texture)
	gs.checkError("ActiveTexture")
}

func (gs *GLS) BindBuffer(target int, vbo uint32) {

	gl.BindBuffer(uint32(target), vbo)
	gs.checkError("BindBuffer")
}

func (gs *GLS) BindTexture(target int, tex uint32) {

	gl.BindTexture(uint32(target), tex)
	gs.checkError("BindTexture")
}

func (gs *GLS) BindVertexArray(vao uint32) {

	gl.BindVertexArray(vao)
	gs.checkError("BindVertexArray")
}

func (gs *GLS) BlendEquation(mode uint32) {

	if gs.blendEquation == mode {
		return
	}
	gl.BlendEquation(mode)
	gs.checkError("BlendEquation")
	gs.blendEquation = mode
}

func (gs *GLS) BlendEquationSeparate(modeRGB uint32, modeAlpha uint32) {

	if gs.blendEquationRGB == modeRGB && gs.blendEquationAlpha == modeAlpha {
		return
	}
	gl.BlendEquationSeparate(uint32(modeRGB), uint32(modeAlpha))
	gs.checkError("BlendEquationSeparate")
	gs.blendEquationRGB = modeRGB
	gs.blendEquationAlpha = modeAlpha
}

func (gs *GLS) BlendFunc(sfactor, dfactor uint32) {

	if gs.blendSrc == sfactor && gs.blendDst == dfactor {
		return
	}
	gl.BlendFunc(sfactor, dfactor)
	gs.checkError("BlendFunc")
	gs.blendSrc = sfactor
	gs.blendDst = dfactor
}

func (gs *GLS) BlendFuncSeparate(srcRGB uint32, dstRGB uint32, srcAlpha uint32, dstAlpha uint32) {

	if gs.blendSrcRGB == srcRGB && gs.blendDstRGB == dstRGB &&
		gs.blendSrcAlpha == srcAlpha && gs.blendDstAlpha == dstAlpha {
		return
	}
	gl.BlendFuncSeparate(srcRGB, dstRGB, srcAlpha, dstAlpha)
	gs.checkError("BlendFuncSeparate")
	gs.blendSrcRGB = srcRGB
	gs.blendDstRGB = dstRGB
	gs.blendSrcAlpha = srcAlpha
	gs.blendDstAlpha = dstAlpha
}

func (gs *GLS) BufferData(target uint32, size int, data interface{}, usage uint32) {

	gl.BufferData(target, size, gl.Ptr(data), usage)
	gs.checkError("BufferData")
}

func (gs *GLS) ClearColor(r, g, b, a float32) {

	gl.ClearColor(r, g, b, a)
}

func (gs *GLS) Clear(mask int) {

	gl.Clear(uint32(mask))
}

func (gs *GLS) DeleteBuffers(vbos ...uint32) {

	gl.DeleteBuffers(int32(len(vbos)), &vbos[0])
	gs.checkError("DeleteBuffers")
}

func (gs *GLS) DeleteTextures(tex ...uint32) {

	gl.DeleteTextures(int32(len(tex)), &tex[0])
	gs.checkError("DeleteTextures")
	gs.Stats.Textures -= len(tex)
}

func (gs *GLS) DeleteVertexArrays(vaos ...uint32) {

	gl.DeleteVertexArrays(int32(len(vaos)), &vaos[0])
	gs.checkError("DeleteVertexArrays")
}

func (gs *GLS) DepthFunc(mode uint32) {

	if gs.depthFunc == mode {
		return
	}
	gl.DepthFunc(mode)
	gs.checkError("DepthFunc")
	gs.depthFunc = mode
}

func (gs *GLS) DepthMask(flag bool) {

	if gs.depthMask == intTrue && flag {
		return
	}
	if gs.depthMask == intFalse && !flag {
		return
	}
	gl.DepthMask(flag)
	gs.checkError("DepthMask")
	if flag {
		gs.depthMask = intTrue
	} else {
		gs.depthMask = intFalse
	}
}

func (gs *GLS) DrawArrays(mode uint32, first int32, count int32) {

	gl.DrawArrays(mode, first, count)
	gs.checkError("DrawArrays")
}

func (gs *GLS) DrawElements(mode uint32, count int32, itype uint32, start uint32) {

	gl.DrawElements(mode, int32(count), itype, gl.PtrOffset(int(start)))
	gs.checkError("DrawElements")
}

func (gs *GLS) Enable(cap int) {

	if gs.capabilities[cap] == capEnabled {
		return
	}
	gl.Enable(uint32(cap))
	gs.checkError("Enable")
	gs.capabilities[cap] = capEnabled
}

func (gs *GLS) EnableVertexAttribArray(index uint32) {

	gl.EnableVertexAttribArray(index)
	gs.checkError("EnableVertexAttribArray")
}

func (gs *GLS) Disable(cap int) {

	if gs.capabilities[cap] == capDisabled {
		return
	}
	gl.Disable(uint32(cap))
	gs.checkError("Disable")
	gs.capabilities[cap] = capDisabled
}

func (gs *GLS) FrontFace(mode uint32) {

	gl.FrontFace(mode)
	gs.checkError("FrontFace")
}

func (gs *GLS) GenBuffer() uint32 {

	var buf uint32
	gl.GenBuffers(1, &buf)
	gs.checkError("GenBuffers")
	gs.Stats.Vbos++
	return buf
}

func (gs *GLS) GenerateMipmap(target uint32) {

	gl.GenerateMipmap(target)
	gs.checkError("GenerateMipmap")
}

func (gs *GLS) GenTexture() uint32 {

	var tex uint32
	gl.GenTextures(1, &tex)
	gs.checkError("GenTextures")
	gs.Stats.Textures++
	return tex
}

func (gs *GLS) GenVertexArray() uint32 {

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gs.checkError("GenVertexArrays")
	gs.Stats.Vaos++
	return vao
}

func (gs *GLS) GetString(name uint32) string {

	cstr := gl.GetString(name)
	return gl.GoStr(cstr)
}

func (gs *GLS) GetViewport() (x, y, width, height int32) {

	return gs.viewportX, gs.viewportY, gs.viewportWidth, gs.viewportHeight
}

func (gs *GLS) LineWidth(width float32) {

	if gs.lineWidth == width {
		return
	}
	gl.LineWidth(width)
	gs.checkError("LineWidth")
	gs.lineWidth = width
}

func (gs *GLS) SetDepthTest(mode bool) {

	if mode {
		gs.Enable(gl.DEPTH_TEST)
	} else {
		gs.Disable(gl.DEPTH_TEST)
	}
}

func (gs *GLS) SetSideView(mode int) {

	if gs.sideView == mode {
		return
	}
	switch mode {
	// Default: show only the front size
	case FrontSide:
		gs.Enable(gl.CULL_FACE)
		gl.FrontFace(gl.CCW)
	// Show only the back side
	case BackSide:
		gs.Enable(gl.CULL_FACE)
		gl.FrontFace(gl.CW)
	// Show both sides
	case DoubleSide:
		gs.Disable(gl.CULL_FACE)
	default:
		panic("SetSideView() invalid mode")
	}
	gs.sideView = mode
}

func (gs *GLS) TexImage2D(target uint32, level int32, iformat int32, width int32, height int32, border int32, format uint32, itype uint32, data interface{}) {

	gl.TexImage2D(uint32(target), int32(level), int32(iformat), int32(width), int32(height), int32(border), uint32(format), uint32(itype), gl.Ptr(data))
	gs.checkError("TexImage2D")
}

func (gs *GLS) TexStorage2D(target int, levels int, iformat int, width, height int) {

	gl.TexStorage2D(uint32(target), int32(levels), uint32(iformat), int32(width), int32(height))
	gs.checkError("TexStorage2D")
}

func (gs *GLS) TexParameteri(target uint32, pname uint32, param int32) {

	gl.TexParameteri(target, pname, param)
	gs.checkError("TexParameteri")
}

func (gs *GLS) PolygonMode(face, mode int) {

	gl.PolygonMode(uint32(face), uint32(mode))
	gs.checkError("PolygonMode")
}

func (gs *GLS) PolygonOffset(factor float32, units float32) {

	gl.PolygonOffset(factor, units)
	gs.checkError("PolygonOffset")
}

func (gs *GLS) Uniform1i(location int32, v0 int32) {

	gl.Uniform1i(location, v0)
	gs.checkError("Uniform1i")
}

func (gs *GLS) Uniform1f(location int32, v0 float32) {

	gl.Uniform1f(location, v0)
	gs.checkError("Uniform1f")
}

func (gs *GLS) Uniform2f(location int32, v0, v1 float32) {

	gl.Uniform2f(location, v0, v1)
	gs.checkError("Uniform2f")
}

func (gs *GLS) Uniform3f(location int32, v0, v1, v2 float32) {

	gl.Uniform3f(location, v0, v1, v2)
	gs.checkError("Uniform3f")
}

func (gs *GLS) Uniform4f(location int32, v0, v1, v2, v3 float32) {

	gl.Uniform4f(location, v0, v1, v2, v3)
	gs.checkError("Uniform4f")
}

func (gs *GLS) UniformMatrix3fv(location int32, count int32, transpose bool, v []float32) {

	gl.UniformMatrix3fv(location, count, transpose, &v[0])
	gs.checkError("UniformMatrix3fv")
}

func (gs *GLS) UniformMatrix4fv(location int32, count int32, transpose bool, v []float32) {

	gl.UniformMatrix4fv(location, count, transpose, &v[0])
	gs.checkError("UniformMatrix4fv")
}

// Use set this program as the current program.
func (gs *GLS) UseProgram(prog *Program) {

	if prog.handle == 0 {
		panic("Invalid program")
	}
	gl.UseProgram(prog.handle)
	gs.checkError("UseProgram")
	gs.Prog = prog

	// Inserts program in cache if not already there.
	if !gs.programs[prog] {
		gs.programs[prog] = true
		log.Debug("New Program activated. Total: %d", len(gs.programs))
	}
}

func (gs *GLS) VertexAttribPointer(index uint32, size int32, xtype uint32, normalized bool, stride int32, offset uint32) {

	gl.VertexAttribPointer(index, size, xtype, normalized, stride, gl.PtrOffset(int(offset)))
	gs.checkError("VertexAttribPointer")
}

func (gs *GLS) Viewport(x, y, width, height int32) {

	gl.Viewport(x, y, width, height)
	gs.checkError("Viewport")
	gs.viewportX = x
	gs.viewportY = y
	gs.viewportWidth = width
	gs.viewportHeight = height
}

// checkError checks the error code of the previously called OpenGL function
func (gls *GLS) checkError(fname string) {

	if gls.checkErrors {
		ecode := gl.GetError()
		if ecode != 0 {
			log.Fatal("Error:%d calling:%s()", ecode, fname)
		}
	}
}
