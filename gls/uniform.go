// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gls

import (
	"fmt"
	"github.com/g3n/engine/math32"
)

//
// Type Uniform is the type for all uniforms
//
type Uniform struct {
	name    string // original name
	nameidx string // cached indexed name
	idx     int    // index value of indexed name
}

// Location returns the current location of the uniform
// for the current active program
func (uni *Uniform) Location(gs *GLS) int32 {

	loc := gs.Prog.GetUniformLocation(uni.name)
	return loc
}

// Location returns the current location of the uniform
// for the current active program and index
func (uni *Uniform) LocationIdx(gs *GLS, idx int) int32 {

	// Rebuilds uniform indexed name if necessary
	if uni.nameidx == "" || uni.idx != idx {
		uni.nameidx = fmt.Sprintf("%s[%d]", uni.name, idx)
		uni.idx = idx
	}
	//log.Debug("Location(%s, %d)", uni.name, idx)
	loc := gs.Prog.GetUniformLocation(uni.nameidx)
	return loc
}

//
// Type Uniform1i is a Uniform containing one int value
//
type Uniform1i struct {
	Uniform
	v0 int32
}

func NewUniform1i(name string) *Uniform1i {

	uni := new(Uniform1i)
	uni.name = name
	return uni
}

func (uni *Uniform1i) Init(name string) {

	uni.name = name
}

func (uni *Uniform1i) Set(v int32) {

	uni.v0 = v
}

func (uni *Uniform1i) Get() int32 {

	return uni.v0
}

func (uni *Uniform1i) Transfer(gs *GLS) {

	gs.Uniform1i(uni.Location(gs), uni.v0)
}

func (uni *Uniform1i) TransferIdx(gs *GLS, idx int) {

	gs.Uniform1i(uni.LocationIdx(gs, idx), uni.v0)
}

//
// Type Uniform1f is a Uniform containing one float32 value
//
type Uniform1f struct {
	Uniform
	v0 float32
}

func NewUniform1f(name string) *Uniform1f {

	uni := new(Uniform1f)
	uni.Init(name)
	return uni
}

func (uni *Uniform1f) Init(name string) {

	uni.name = name
}

func (uni *Uniform1f) Set(v float32) {

	uni.v0 = v
}

func (uni *Uniform1f) Get() float32 {

	return uni.v0
}

func (uni *Uniform1f) Transfer(gs *GLS) {

	gs.Uniform1f(uni.Location(gs), uni.v0)
}

func (uni *Uniform1f) TransferIdx(gs *GLS, idx int) {

	gs.Uniform1f(uni.LocationIdx(gs, idx), uni.v0)
}

//
// Type Uniform2f is a Uniform containing two float32 values
//
type Uniform2f struct {
	Uniform
	v0 float32
	v1 float32
}

func NewUniform2f(name string) *Uniform2f {

	uni := new(Uniform2f)
	uni.Init(name)
	return uni
}

func (uni *Uniform2f) Init(name string) {

	uni.name = name
}

func (uni *Uniform2f) Set(v0, v1 float32) {

	uni.v0 = v0
	uni.v1 = v1
}

func (uni *Uniform2f) Get() (float32, float32) {

	return uni.v0, uni.v1
}

func (uni *Uniform2f) SetVector2(v *math32.Vector2) {

	uni.v0 = v.X
	uni.v1 = v.Y
}

func (uni *Uniform2f) GetVector2() math32.Vector2 {

	return math32.Vector2{uni.v0, uni.v1}
}

func (uni *Uniform2f) Transfer(gs *GLS) {

	gs.Uniform2f(uni.Location(gs), uni.v0, uni.v1)
}

func (uni *Uniform2f) TransferIdx(gs *GLS, idx int) {

	gs.Uniform2f(uni.LocationIdx(gs, idx), uni.v0, uni.v1)
}

//
// Type Uniform3f is a Uniform containing three float32 values
//
type Uniform3f struct {
	Uniform
	v0 float32
	v1 float32
	v2 float32
}

func NewUniform3f(name string) *Uniform3f {

	uni := new(Uniform3f)
	uni.Init(name)
	return uni
}

func (uni *Uniform3f) Init(name string) {

	uni.name = name
}

func (uni *Uniform3f) Set(v0, v1, v2 float32) {

	uni.v0 = v0
	uni.v1 = v1
	uni.v2 = v2
}

func (uni *Uniform3f) Get() (float32, float32, float32) {

	return uni.v0, uni.v1, uni.v2
}

func (uni *Uniform3f) SetVector3(v *math32.Vector3) {

	uni.v0 = v.X
	uni.v1 = v.Y
	uni.v2 = v.Z
}

func (uni *Uniform3f) GetVector3() math32.Vector3 {

	return math32.Vector3{uni.v0, uni.v1, uni.v2}
}

func (uni *Uniform3f) SetColor(color *math32.Color) {

	uni.v0 = color.R
	uni.v1 = color.G
	uni.v2 = color.B
}

func (uni *Uniform3f) GetColor() math32.Color {

	return math32.Color{uni.v0, uni.v1, uni.v2}
}

func (uni *Uniform3f) Transfer(gl *GLS) {

	loc := uni.Location(gl)
	gl.Uniform3f(loc, uni.v0, uni.v1, uni.v2)
	//log.Debug("Uniform3f: %s (%v) -> %v,%v,%v", uni.name, loc, uni.v0, uni.v1, uni.v2)
}

func (uni *Uniform3f) TransferIdx(gl *GLS, idx int) {

	loc := uni.LocationIdx(gl, idx)
	gl.Uniform3f(loc, uni.v0, uni.v1, uni.v2)
	//log.Debug("Uniform3f: %s -> %v,%v,%v", uni.nameidx, uni.v0, uni.v1, uni.v2)
}

//
// Type Uniform4f is a Uniform containing four float32 values
//
type Uniform4f struct {
	Uniform
	v0 float32
	v1 float32
	v2 float32
	v3 float32
}

func NewUniform4f(name string) *Uniform4f {

	uni := new(Uniform4f)
	uni.Init(name)
	return uni
}

func (uni *Uniform4f) Init(name string) {

	uni.name = name
}

func (uni *Uniform4f) Set(v0, v1, v2, v3 float32) {

	uni.v0 = v0
	uni.v1 = v1
	uni.v2 = v2
	uni.v3 = v3
}

func (uni *Uniform4f) Get() (float32, float32, float32, float32) {

	return uni.v0, uni.v1, uni.v2, uni.v3
}

func (uni *Uniform4f) SetVector4(v *math32.Vector4) {

	uni.v0 = v.X
	uni.v1 = v.Y
	uni.v2 = v.Z
	uni.v3 = v.W
}

func (uni *Uniform4f) GetVector4() math32.Vector4 {

	return math32.Vector4{uni.v0, uni.v1, uni.v2, uni.v3}
}

func (uni *Uniform4f) SetColor4(c *math32.Color4) {

	uni.v0 = c.R
	uni.v1 = c.G
	uni.v2 = c.B
	uni.v3 = c.A
}

func (uni *Uniform4f) GetColor4() math32.Color4 {

	return math32.Color4{uni.v0, uni.v1, uni.v2, uni.v3}
}

func (uni *Uniform4f) Transfer(gl *GLS) {

	//log.Debug("Uniform4f.Transfer: %s %d", uni.name, uni.Location(gl))
	gl.Uniform4f(uni.Location(gl), uni.v0, uni.v1, uni.v2, uni.v3)
}

func (uni *Uniform4f) TransferIdx(gl *GLS, idx int) {

	gl.Uniform4f(uni.LocationIdx(gl, idx), uni.v0, uni.v1, uni.v2, uni.v3)
}

//
// Type UniformMatrix3f is a Uniform containing nine float32 values
// organized as 3x3 matrix
//
type UniformMatrix3f struct {
	Uniform
	v [9]float32
}

func NewUniformMatrix3f(name string) *UniformMatrix3f {

	uni := new(UniformMatrix3f)
	uni.Init(name)
	return uni
}

func (uni *UniformMatrix3f) Init(name string) {

	uni.name = name
}

func (uni *UniformMatrix3f) SetMatrix3(m *math32.Matrix3) {

	uni.v = *m
}

func (uni *UniformMatrix3f) GetMatrix3() math32.Matrix3 {

	return uni.v
}

func (uni *UniformMatrix3f) Transfer(gl *GLS) {

	gl.UniformMatrix3fv(uni.Location(gl), 1, false, uni.v[0:9])
}

func (uni *UniformMatrix3f) TransferIdx(gl *GLS, idx int) {

	gl.UniformMatrix3fv(uni.LocationIdx(gl, idx), 1, false, uni.v[0:9])
}

//
// Type UniformMatrix4f is a Uniform containing sixteen float32 values
// organized as 4x4 matrix
//
type UniformMatrix4f struct {
	Uniform
	v [16]float32
}

func NewUniformMatrix4f(name string) *UniformMatrix4f {

	uni := new(UniformMatrix4f)
	uni.Init(name)
	return uni
}

func (uni *UniformMatrix4f) Init(name string) {

	uni.name = name
}

func (uni *UniformMatrix4f) SetMatrix4(m *math32.Matrix4) {

	uni.v = *m
}

func (uni *UniformMatrix4f) GetMatrix4() math32.Matrix4 {

	return uni.v
}

func (uni *UniformMatrix4f) Transfer(gl *GLS) {

	gl.UniformMatrix4fv(uni.Location(gl), 1, false, uni.v[0:16])
}

func (uni *UniformMatrix4f) TransferIdx(gl *GLS, idx int) {

	gl.UniformMatrix4fv(uni.LocationIdx(gl, idx), 1, false, uni.v[0:16])
}
