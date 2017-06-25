// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gls

import (
	"github.com/g3n/engine/math32"
	"unsafe"
)

// VBO abstracts an OpenGL Vertex Buffer Object
type VBO struct {
	gs      *GLS            // reference to GLS state
	handle  uint32          // OpenGL handle for this VBO
	usage   uint32          // Expected usage patter of the buffer
	update  bool            // Update flag
	buffer  math32.ArrayF32 // Data buffer
	attribs []VBOattrib     // List of attributes
}

// VBOattrib describes one attribute of an OpenGL Vertex Buffer Object
type VBOattrib struct {
	Name     string // Name of of the attribute
	ItemSize int32  // Number of elements for each item
}

// NewVBO creates and returns a pointer to a new OpenGL Vertex Buffer Object
func NewVBO() *VBO {

	vbo := new(VBO)
	vbo.init()
	return vbo
}

// init initializes this VBO
func (vbo *VBO) init() {

	vbo.gs = nil
	vbo.handle = 0
	vbo.usage = STATIC_DRAW
	vbo.update = true
	vbo.attribs = make([]VBOattrib, 0)
}

// AddAttrib adds a new attribute to this VBO
func (vbo *VBO) AddAttrib(name string, itemSize int32) *VBO {

	vbo.attribs = append(vbo.attribs, VBOattrib{
		Name:     name,
		ItemSize: itemSize,
	})
	return vbo
}

// Attrib finds and returns pointer the attribute with the specified name
// or nil if not found
func (vbo *VBO) Attrib(name string) *VBOattrib {

	for _, attr := range vbo.attribs {
		if attr.Name == name {
			return &attr
		}
	}
	return nil
}

// AttribAt returns pointer to the VBO attribute at the specified index
func (vbo *VBO) AttribAt(idx int) *VBOattrib {

	return &vbo.attribs[idx]
}

// AttribCount returns the current number of attributes for this VBO
func (vbo *VBO) AttribCount() int {

	return len(vbo.attribs)
}

// Dispose disposes of this VBO OpenGL resources
// As currently the VBO is used only by the Geometry object
// it is not referenced counted.
func (vbo *VBO) Dispose() {

	if vbo.gs != nil {
		vbo.gs.DeleteBuffers(vbo.handle)
	}
	vbo.gs = nil
}

// Sets the VBO buffer
func (vbo *VBO) SetBuffer(buffer math32.ArrayF32) *VBO {

	vbo.buffer = buffer
	vbo.update = true
	return vbo
}

// Sets the expected usage pattern of the buffer.
// The default value is GL_STATIC_DRAW.
func (vbo *VBO) SetUsage(usage uint32) {

	vbo.usage = usage
}

// Buffer returns pointer to the VBO buffer
func (vbo *VBO) Buffer() *math32.ArrayF32 {

	return &vbo.buffer
}

// Updates sets the update flag to force the VBO update
func (vbo *VBO) Update() {

	vbo.update = true
}

// Stride returns the stride of this VBO which is the number of bytes
// used by one group attributes side by side in the buffer
// Ex: x y z r g b x y z r g b ...x y z r g b
// The stride will be: sizeof(float) * 6 = 24
func (vbo *VBO) Stride() int {

	stride := 0
	elsize := int(unsafe.Sizeof(float32(0)))
	for _, attrib := range vbo.attribs {
		stride += elsize * int(attrib.ItemSize)
	}
	return stride
}

// Transfer is called internally and transfer the data in the VBO buffer to OpenGL if necessary
func (vbo *VBO) Transfer(gs *GLS) {

	// If the VBO buffer is empty, ignore
	if vbo.buffer.Bytes() == 0 {
		return
	}

	// First time initialization
	if vbo.gs == nil {
		vbo.handle = gs.GenBuffer()
		gs.BindBuffer(ARRAY_BUFFER, vbo.handle)
		// Calculates stride
		stride := vbo.Stride()
		// For each attribute
		var items uint32 = 0
		var offset uint32 = 0
		elsize := int32(unsafe.Sizeof(float32(0)))
		for _, attrib := range vbo.attribs {
			// Get attribute location in the current program
			loc := gs.prog.GetAttribLocation(attrib.Name)
			if loc < 0 {
				continue
			}
			// Enables attribute and sets its stride and offset in the buffer
			gs.EnableVertexAttribArray(uint32(loc))
			gs.VertexAttribPointer(uint32(loc), attrib.ItemSize, FLOAT, false, int32(stride), offset)
			items += uint32(attrib.ItemSize)
			offset = uint32(elsize) * items
		}
		vbo.gs = gs // this indicates that the vbo was initialized
	}
	if !vbo.update {
		return
	}
	// Transfer the VBO data to OpenGL
	gs.BindBuffer(ARRAY_BUFFER, vbo.handle)
	gs.BufferData(ARRAY_BUFFER, vbo.buffer.Bytes(), &vbo.buffer[0], vbo.usage)
	vbo.update = false
}
