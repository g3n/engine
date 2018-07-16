// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gls

import (
	"github.com/g3n/engine/math32"
	"unsafe"
)

// VBO abstracts an OpenGL Vertex Buffer Object.
type VBO struct {
	gs      *GLS            // Reference to OpenGL state
	handle  uint32          // OpenGL handle for this VBO
	usage   uint32          // Expected usage pattern of the buffer
	update  bool            // Update flag
	buffer  math32.ArrayF32 // Data buffer
	attribs []VBOattrib     // List of attributes
}

// VBOattrib describes one attribute of an OpenGL Vertex Buffer Object.
type VBOattrib struct {
	Type     AttribType // Type of the attribute
	Name     string     // Name of the attribute
	ItemSize int32      // Number of elements
}

// AttribType is the functional type of a vbo attribute.
type AttribType int

const (
	Undefined = AttribType(iota)
	VertexPosition
	VertexNormal
	VertexColor
	VertexTexcoord
)

// Map from attribute type to default attribute name.
var attribTypeMap = map[AttribType]string{
	VertexPosition: "VertexPosition",
	VertexNormal:   "VertexNormal",
	VertexColor:    "VertexColor",
	VertexTexcoord: "VertexTexcoord",
}

// NewVBO creates and returns a pointer to a new OpenGL Vertex Buffer Object.
func NewVBO(buffer math32.ArrayF32) *VBO {

	vbo := new(VBO)
	vbo.init()
	vbo.SetBuffer(buffer)
	return vbo
}

// init initializes the VBO.
func (vbo *VBO) init() {

	vbo.gs = nil
	vbo.handle = 0
	vbo.usage = STATIC_DRAW
	vbo.update = true
	vbo.attribs = make([]VBOattrib, 0)
}

// AddAttrib adds a new attribute to the VBO by attribute type.
func (vbo *VBO) AddAttrib(atype AttribType, itemSize int32) *VBO {

	vbo.attribs = append(vbo.attribs, VBOattrib{
		Type:     atype,
		Name:     attribTypeMap[atype],
		ItemSize: itemSize,
	})
	return vbo
}

// AddAttribName adds a new attribute to the VBO by name.
func (vbo *VBO) AddAttribName(name string, itemSize int32) *VBO {

	vbo.attribs = append(vbo.attribs, VBOattrib{
		Type:     Undefined,
		Name:     name,
		ItemSize: itemSize,
	})
	return vbo
}

// Attrib finds and returns a pointer to the VBO attribute with the specified type.
// Returns nil if not found.
func (vbo *VBO) Attrib(atype AttribType) *VBOattrib {

	for i := range vbo.attribs {
		if vbo.attribs[i].Type == atype {
			return &vbo.attribs[i]
		}
	}
	return nil
}

// AttribName finds and returns a pointer to the VBO attribute with the specified name.
// Returns nil if not found.
func (vbo *VBO) AttribName(name string) *VBOattrib {

	for i := range vbo.attribs {
		if vbo.attribs[i].Name == name {
			return &vbo.attribs[i]
		}
	}
	return nil
}

// AttribAt returns a pointer to the VBO attribute at the specified index.
func (vbo *VBO) AttribAt(idx int) *VBOattrib {

	return &vbo.attribs[idx]
}

// AttribCount returns the current number of attributes for this VBO.
func (vbo *VBO) AttribCount() int {

	return len(vbo.attribs)
}

// Attributes returns the attributes for this VBO.
func (vbo *VBO) Attributes() []VBOattrib {

	return vbo.attribs
}

// Dispose disposes of the OpenGL resources used by the VBO.
// As currently the VBO is used only by Geometry objects
// it is not referenced counted.
func (vbo *VBO) Dispose() {

	if vbo.gs != nil {
		vbo.gs.DeleteBuffers(vbo.handle)
	}
	vbo.gs = nil
}

// SetBuffer sets the VBO buffer.
func (vbo *VBO) SetBuffer(buffer math32.ArrayF32) *VBO {

	vbo.buffer = buffer
	vbo.update = true
	return vbo
}

// SetUsage sets the expected usage pattern of the buffer.
// The default value is GL_STATIC_DRAW.
func (vbo *VBO) SetUsage(usage uint32) {

	vbo.usage = usage
}

// Buffer returns a pointer to the VBO buffer.
func (vbo *VBO) Buffer() *math32.ArrayF32 {

	return &vbo.buffer
}

// Update sets the update flag to force the VBO update.
func (vbo *VBO) Update() {

	vbo.update = true
}

// AttribOffset returns the total number of elements from
// all attributes preceding the attribute specified by type.
func (vbo *VBO) AttribOffset(attribType AttribType) int {

	elementCount := 0
	for _, attr := range vbo.attribs {
		if attr.Type == attribType {
			return elementCount
		}
		elementCount += int(attr.ItemSize)
	}
	return elementCount
}

// AttribOffsetName returns the total number of elements from
// all attributes preceding the attribute specified by name.
func (vbo *VBO) AttribOffsetName(name string) int {

	elementCount := 0
	for _, attr := range vbo.attribs {
		if attr.Name == name {
			return elementCount
		}
		elementCount += int(attr.ItemSize)
	}
	return elementCount
}

// Stride returns the stride of the VBO, which is the number of elements in
// one complete set of group attributes. E.g. for an interleaved VBO with two attributes:
// "VertexPosition" (3 elements) and "VertexTexcoord" (2 elements), the stride would be 5:
// [X, Y, Z, U, V], X, Y, Z, U, V, X, Y, Z, U, V... X, Y, Z, U, V.
func (vbo *VBO) Stride() int {

	stride := 0
	for _, attrib := range vbo.attribs {
		stride += int(attrib.ItemSize)
	}
	return stride
}

// StrideSize returns the number of bytes used by one complete set of group attributes.
// E.g. for an interleaved VBO with two attributes: "VertexPosition" (3 elements)
// and "VertexTexcoord" (2 elements), the stride would be 5:
// [X, Y, Z, U, V], X, Y, Z, U, V, X, Y, Z, U, V... X, Y, Z, U, V
// and the stride size would be: sizeof(float)*stride = 4*5 = 20
func (vbo *VBO) StrideSize() int {

	stride := vbo.Stride()
	elsize := int(unsafe.Sizeof(float32(0)))
	return stride * elsize
}

// Transfer (called internally) transfers the data from the VBO buffer to OpenGL if necessary.
func (vbo *VBO) Transfer(gs *GLS) {

	// If the VBO buffer is empty, ignore
	if vbo.buffer.Bytes() == 0 {
		return
	}

	// First time initialization
	if vbo.gs == nil {
		vbo.handle = gs.GenBuffer()
		gs.BindBuffer(ARRAY_BUFFER, vbo.handle)
		// Calculates stride size
		strideSize := vbo.StrideSize()
		// For each attribute
		var items uint32
		var offset uint32
		elsize := int32(unsafe.Sizeof(float32(0)))
		for _, attrib := range vbo.attribs {
			// Get attribute location in the current program
			loc := gs.prog.GetAttribLocation(attrib.Name)
			if loc < 0 {
				log.Warn("Attribute not found: %v", attrib.Name)
				continue
			}
			// Enables attribute and sets its stride and offset in the buffer
			gs.EnableVertexAttribArray(uint32(loc))
			gs.VertexAttribPointer(uint32(loc), attrib.ItemSize, FLOAT, false, int32(strideSize), offset)
			items += uint32(attrib.ItemSize)
			offset = uint32(elsize) * items
		}
		vbo.gs = gs // this indicates that the vbo was initialized
	}

	// If nothing has changed, no need to transfer data to OpenGL
	if !vbo.update {
		return
	}

	// Transfer the VBO data to OpenGL
	gs.BindBuffer(ARRAY_BUFFER, vbo.handle)
	gs.BufferData(ARRAY_BUFFER, vbo.buffer.Bytes(), &vbo.buffer[0], vbo.usage)
	vbo.update = false
}

// OperateOnVectors3 iterates over all 3-float32 items for the specified attribute
// and calls the specified callback function with a pointer to each item as a Vector3.
// The vector pointers can be modified inside the callback and the modifications will be applied to the buffer at each iteration.
// The callback function returns false to continue or true to break.
func (vbo *VBO) OperateOnVectors3(attribType AttribType, cb func(vec *math32.Vector3) bool) {

	stride := vbo.Stride()
	offset := vbo.AttribOffset(attribType)
	buffer := vbo.Buffer()

	// Call callback for each vector3, updating the buffer afterward
	var vec math32.Vector3
	for i := offset; i < vbo.buffer.Size(); i += stride {
		buffer.GetVector3(i, &vec)
		brk := cb(&vec)
		buffer.SetVector3(i, &vec)
		if brk {
			break
		}
	}
	vbo.Update()
}

// ReadVectors3 iterates over all 3-float32 items for the specified attribute
// and calls the specified callback function with the value of each item as a Vector3.
// The callback function returns false to continue or true to break.
func (vbo *VBO) ReadVectors3(attribType AttribType, cb func(vec math32.Vector3) bool) {

	stride := vbo.Stride()
	offset := vbo.AttribOffset(attribType)
	positions := vbo.Buffer()

	// Call callback for each vector3
	var vec math32.Vector3
	for i := offset; i < positions.Size(); i += stride {
		positions.GetVector3(i, &vec)
		brk := cb(vec)
		if brk {
			break
		}
	}
}


// Read3Vectors3 iterates over all 3-float32 items (3 items at a time) for the specified attribute
// and calls the specified callback function with the value of each of the 3 items as Vector3.
// The callback function returns false to continue or true to break.
func (vbo *VBO) ReadTripleVectors3(attribType AttribType, cb func(vec1, vec2, vec3 math32.Vector3) bool) {

	stride := vbo.Stride()
	offset := vbo.AttribOffset(attribType)
	positions := vbo.Buffer()

	doubleStride := 2*stride
	loopStride := 3*stride

	// Call callback for each vector3 triple
	var vec1, vec2, vec3 math32.Vector3
	for i := offset; i < positions.Size(); i += loopStride {
		positions.GetVector3(i, &vec1)
		positions.GetVector3(i + stride, &vec2)
		positions.GetVector3(i + doubleStride, &vec3)
		brk := cb(vec1, vec2, vec3)
		if brk {
			break
		}
	}
}