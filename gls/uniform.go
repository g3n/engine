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

	loc := gs.prog.GetUniformLocation(uni.name)
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
	loc := gs.prog.GetUniformLocation(uni.nameidx)
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

// NewUniformMatrix3 creates and returns a pointer to a new UniformMatrix3f
// with the specified name
func NewUniformMatrix3f(name string) *UniformMatrix3f {

	uni := new(UniformMatrix3f)
	uni.Init(name)
	return uni
}

// Init initializes this uniform the specified name
// It is normally used when the uniform is embedded in another object.
func (uni *UniformMatrix3f) Init(name string) {

	uni.name = name
}

// SetMatrix3 sets the matrix stored by the uniform
func (uni *UniformMatrix3f) SetMatrix3(m *math32.Matrix3) {

	uni.v = *m
}

// GetMatrix3 gets the matrix stored by the uniform
func (uni *UniformMatrix3f) GetMatrix3() math32.Matrix3 {

	return uni.v
}

// SetElement sets the value of the matrix element at the specified column and row
func (uni *UniformMatrix3f) SetElement(col, row int, v float32) {

	uni.v[col*3+row] = v
}

// GetElement gets the value of the matrix element at the specified column and row
func (uni *UniformMatrix3f) GetElement(col, row int, v float32) float32 {

	return uni.v[col*3+row]
}

// Set sets the value of the matrix element by its position starting
// from 0 for col0, row0 to 8 from col2, row2.
// This way the matrix can be considered as a vector of 9 elements
func (uni *UniformMatrix3f) Set(pos int, v float32) {

	uni.v[pos] = v
}

// Get gets the value of the matrix element by its position starting
// from 0 for col0, row0 to 8 from col2, row2.
// This way the matrix can be considered as a vector of 9 elements
func (uni *UniformMatrix3f) Get(pos int) float32 {

	return uni.v[pos]
}

// Transfer transfer the uniform matrix data to the graphics library
func (uni *UniformMatrix3f) Transfer(gl *GLS) {

	gl.UniformMatrix3fv(uni.Location(gl), 1, false, &uni.v[0])
}

// TransferIdx transfer the uniform matrix data to a specified destination index
// of an uniform array to the graphics library
func (uni *UniformMatrix3f) TransferIdx(gl *GLS, idx int) {

	gl.UniformMatrix3fv(uni.LocationIdx(gl, idx), 1, false, &uni.v[0])
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

	gl.UniformMatrix4fv(uni.Location(gl), 1, false, &uni.v[0])
}

func (uni *UniformMatrix4f) TransferIdx(gl *GLS, idx int) {

	gl.UniformMatrix4fv(uni.LocationIdx(gl, idx), 1, false, &uni.v[0])
}

//
// Type Uniform1fv is a uniform containing an array of float32 values
//
type Uniform1fv struct {
	Uniform           // embedded uniform
	v       []float32 // array of values
}

// NewUniform1fv creates and returns an uniform with array of float32 values
// with the specified size
func NewUniform1fv(name string, count int) *Uniform1fv {

	uni := new(Uniform1fv)
	uni.Init(name, count)
	return uni
}

// Init initializes an Uniform1fv object with the specified name and count of float32 values.
// It is normally used when the uniform is embedded in another object.
func (uni *Uniform1fv) Init(name string, count int) {

	uni.name = name
	uni.v = make([]float32, count)
}

// SetVector3 sets the value of the elements of uniform starting at pos
// from the specified Vector3 object.
func (uni *Uniform1fv) SetVector3(pos int, v *math32.Vector3) {

	uni.v[pos] = v.X
	uni.v[pos+1] = v.Y
	uni.v[pos+2] = v.Z
}

// GetVector3 gets the value of the elements of the uniform starting at
// pos as a Vector3 object.
func (uni *Uniform1fv) GetVector3(pos int) math32.Vector3 {

	return math32.Vector3{uni.v[pos], uni.v[pos+1], uni.v[pos+2]}
}

// SetColor sets the value of the elements of the uniform starting at pos
// from the specified Color object.
func (uni *Uniform1fv) SetColor(pos int, color *math32.Color) {

	uni.v[pos] = color.R
	uni.v[pos+1] = color.G
	uni.v[pos+2] = color.B
}

// GetColor gets the value of the elements of the uniform starting at pos
// as a Color object.
func (uni *Uniform1fv) GetColor(pos int) math32.Color {

	return math32.Color{uni.v[pos], uni.v[pos+1], uni.v[pos+2]}
}

// Set sets the value of the element at the specified position
func (uni *Uniform1fv) Set(pos int, v float32) {

	uni.v[pos] = v
}

// Get gets the value of the element at the specified position
func (uni *Uniform1fv) Get(pos int, v float32) float32 {

	return uni.v[pos]
}

// TransferIdx transfer a block of 'count' values of this uniform starting at 'pos'.
func (uni *Uniform1fv) TransferIdx(gl *GLS, pos, count int) {

	gl.Uniform1fv(uni.LocationIdx(gl, pos), 1, uni.v[pos:])
}

//
// Type Uniform3fv is a uniform containing an array of three float32 values
//
type Uniform3fv struct {
	Uniform           // embedded uniform
	count   int       // number of groups of 3 float32 values
	v       []float32 // array of values
}

// NewUniform3fv creates and returns an uniform array with the specified size
// of 3 float values
func NewUniform3fv(name string, count int) *Uniform3fv {

	uni := new(Uniform3fv)
	uni.Init(name, count)
	return uni
}

// Init initializes an Uniform3fv object with the specified name and count of 3 float32 groups.
// It is normally used when the uniform is embedded in another object.
func (uni *Uniform3fv) Init(name string, count int) {

	uni.name = name
	uni.count = count
	uni.v = make([]float32, count*3)
}

// Set sets the value of all elements of the specified group of 3 floats for this uniform array
func (uni *Uniform3fv) Set(idx int, v0, v1, v2 float32) {

	pos := idx * 3
	uni.v[pos] = v0
	uni.v[pos+1] = v1
	uni.v[pos+2] = v2
}

// Get gets the value of all elements of the specified group of 3 floats for this uniform array
func (uni *Uniform3fv) Get(idx int) (v0, v1, v2 float32) {

	pos := idx * 3
	return uni.v[pos], uni.v[pos+1], uni.v[pos+2]
}

// SetVector3 sets the value of all elements for the specified group of 3 float for this uniform array
// from the specified Vector3 object.
func (uni *Uniform3fv) SetVector3(idx int, v *math32.Vector3) {

	pos := idx * 3
	uni.v[pos] = v.X
	uni.v[pos+1] = v.Y
	uni.v[pos+2] = v.Z
}

// GetVector3 gets the value of all elements of the specified group of 3 float for this uniform array
// as a Vector3 object.
func (uni *Uniform3fv) GetVector3(idx int) math32.Vector3 {

	pos := idx * 3
	return math32.Vector3{uni.v[pos], uni.v[pos+1], uni.v[pos+2]}
}

// SetColor sets the value of all elements of the specified group of 3 floats for this uniform array
// form the specified Color object.
func (uni *Uniform3fv) SetColor(idx int, color *math32.Color) {

	pos := idx * 3
	uni.v[pos] = color.R
	uni.v[pos+1] = color.G
	uni.v[pos+2] = color.B
}

// GetColor gets the value of all elements of the specified group of 3 float for this uniform array
// as a Color object.
func (uni *Uniform3fv) GetColor(idx int) math32.Color {

	pos := idx * 3
	return math32.Color{uni.v[pos], uni.v[pos+1], uni.v[pos+2]}
}

// SetPos sets the value at the specified position in the uniform array.
func (uni *Uniform3fv) SetPos(pos int, v float32) {

	uni.v[pos] = v
}

// GetPos gets the value at the specified position in the uniform array.
func (uni *Uniform3fv) GetPos(pos int) float32 {

	return uni.v[pos]
}

// Transfer transfers the current values of this uniform to the current shader program
func (uni *Uniform3fv) Transfer(gl *GLS) {

	gl.Uniform3fv(uni.Location(gl), int32(uni.count), uni.v)
}

// Transfer transfers the current values of this uniform to a specified
// start position in the target uniform array.
func (uni *Uniform3fv) TransferIdx(gl *GLS, idx int) {

	gl.Uniform3fv(uni.LocationIdx(gl, idx), int32(uni.count), uni.v)
}

//
// Type Uniform4fv is a Uniform containing an array of four float32 values
//
type Uniform4fv struct {
	Uniform           // embedded uniform
	count   int       // number of group of 4 float32 values
	v       []float32 // array of values
}

// NewUniform4fv creates and returns an uniform array with the specified size
// of 4 float values
func NewUniform4fv(name string, count int) *Uniform4fv {

	uni := new(Uniform4fv)
	uni.Init(name, count)
	return uni
}

// Init initializes an Uniform4fv object with the specified name and count of 4 float32 groups.
// It is normally used when the uniform is embedded in another object.
func (uni *Uniform4fv) Init(name string, count int) {

	uni.name = name
	uni.count = count
	uni.v = make([]float32, count*4)
}

// Set sets the value of all elements of the specified group of 4 floats for this uniform array
func (uni *Uniform4fv) Set(idx int, v0, v1, v2, v3 float32) {

	pos := idx * 4
	uni.v[pos] = v0
	uni.v[pos+1] = v1
	uni.v[pos+2] = v2
	uni.v[pos+3] = v3
}

// Get gets the value of all elements of the specified group of 4 floats for this uniform array
func (uni *Uniform4fv) Get(idx int) (v0, v1, v2, v3 float32) {

	pos := idx * 4
	return uni.v[pos], uni.v[pos+1], uni.v[pos+2], uni.v[pos+3]
}

// SetVector4 sets the value of all elements for the specified group of 4 float for this uniform array
// from the specified Vector4 object.
func (uni *Uniform4fv) SetVector4(idx int, v *math32.Vector4) {

	pos := idx * 4
	uni.v[pos] = v.X
	uni.v[pos+1] = v.Y
	uni.v[pos+2] = v.Z
	uni.v[pos+3] = v.W
}

// GetVector4 gets the value of all elements of the specified group of 4 float for this uniform array
// as a Vector4 object.
func (uni *Uniform4fv) GetVector4(idx int) math32.Vector4 {

	pos := idx * 4
	return math32.Vector4{uni.v[pos], uni.v[pos+1], uni.v[pos+2], uni.v[pos+3]}
}

// SetColor4 sets the value of all elements of the specified group of 4 floats for this uniform array
// form the specified Color4 object.
func (uni *Uniform4fv) SetColor4(idx int, color *math32.Color4) {

	pos := idx * 4
	uni.v[pos] = color.R
	uni.v[pos+1] = color.G
	uni.v[pos+2] = color.B
	uni.v[pos+3] = color.A
}

// GetColor4 gets the value of all elements of the specified group of 4 float for this uniform array
// as a Color4 object.
func (uni *Uniform4fv) GetColor4(idx int) math32.Color4 {

	pos := idx * 4
	return math32.Color4{uni.v[pos], uni.v[pos+1], uni.v[pos+2], uni.v[pos+3]}
}

// Transfer transfers the current values of this uniform to the current shader program
func (uni *Uniform4fv) Transfer(gl *GLS) {

	gl.Uniform4fv(uni.Location(gl), int32(uni.count), uni.v)
}
