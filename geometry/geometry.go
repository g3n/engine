// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package geometry

import (
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/math32"
	"strconv"
)

// IGeometry is the interface for all geometries.
type IGeometry interface {
	GetGeometry() *Geometry
	RenderSetup(gs *gls.GLS)
	Dispose()
}

// Geometry encapsulates a three-dimensional vertex-based geometry.
type Geometry struct {
	gs            *gls.GLS          // Pointer to gl context. Valid after first render setup
	groups        []Group           // Array geometry groups
	refcount      int               // Current number of references
	vbos          []*gls.VBO        // Array of VBOs
	handleVAO     uint32            // Handle to OpenGL VAO
	indices       math32.ArrayU32   // Buffer with indices
	handleIndices uint32            // Handle to OpenGL buffer for indices
	updateIndices bool              // Flag to indicate that indices must be transferred
	ShaderDefines gls.ShaderDefines // Geometry-specific shader defines

	// Geometric properties
	boundingBox    math32.Box3    // Last calculated bounding box
	boundingSphere math32.Sphere  // Last calculated bounding sphere
	area           float32        // Last calculated area
	volume         float32        // Last calculated volume
	rotInertia     math32.Matrix3 // Last calculated rotational inertia matrix

	// Flags indicating whether geometric properties are valid
	boundingBoxValid    bool // Indicates if last calculated bounding box is valid
	boundingSphereValid bool // Indicates if last calculated bounding sphere is valid
	areaValid           bool // Indicates if last calculated area is valid
	volumeValid         bool // Indicates if last calculated volume is valid
	rotInertiaValid     bool // Indicates if last calculated rotational inertia matrix is valid
}

// Group is a geometry group object.
type Group struct {
	Start    int    // Index of first element of the group
	Count    int    // Number of elements in the group
	Matindex int    // Material index for this group
	Matid    string // Material id used when loading external models
}

// NewGeometry creates and returns a pointer to a new Geometry.
func NewGeometry() *Geometry {

	g := new(Geometry)
	g.Init()
	return g
}

// Init initializes the geometry.
func (g *Geometry) Init() {

	g.refcount = 1
	g.vbos = make([]*gls.VBO, 0)
	g.groups = make([]Group, 0)
	g.gs = nil
	g.handleVAO = 0
	g.handleIndices = 0
	g.updateIndices = true
	g.ShaderDefines = *gls.NewShaderDefines()
}

// Incref increments the reference count for this geometry
// and returns a pointer to the geometry.
// It should be used when this geometry is shared by another
// Graphic object.
func (g *Geometry) Incref() *Geometry {

	g.refcount++
	return g
}

// Dispose decrements this geometry reference count and
// if possible releases OpenGL resources, C memory
// and VBOs associated with this geometry.
func (g *Geometry) Dispose() {

	if g.refcount > 1 {
		g.refcount--
		return
	}

	// Delete VAO and indices buffer
	if g.gs != nil {
		g.gs.DeleteVertexArrays(g.handleVAO)
		g.gs.DeleteBuffers(g.handleIndices)
	}
	// Delete this geometry VBO buffers
	for i := 0; i < len(g.vbos); i++ {
		g.vbos[i].Dispose()
	}
	g.Init()
}

// GetGeometry satisfies the IGeometry interface.
func (g *Geometry) GetGeometry() *Geometry {

	return g
}

// AddGroup adds a geometry group (for multimaterial).
func (g *Geometry) AddGroup(start, count, matIndex int) *Group {

	g.groups = append(g.groups, Group{start, count, matIndex, ""})
	return &g.groups[len(g.groups)-1]
}

// AddGroupList adds the specified list of groups to this geometry.
func (g *Geometry) AddGroupList(groups []Group) {

	for _, group := range groups {
		g.groups = append(g.groups, group)
	}
}

// GroupCount returns the number of geometry groups (for multimaterial).
func (g *Geometry) GroupCount() int {

	return len(g.groups)
}

// GroupAt returns pointer to geometry group at the specified index.
func (g *Geometry) GroupAt(idx int) *Group {

	return &g.groups[idx]
}

// SetIndices sets the indices array for this geometry.
func (g *Geometry) SetIndices(indices math32.ArrayU32) {

	g.indices = indices
	g.updateIndices = true
	g.boundingBoxValid = false
	g.boundingSphereValid = false
}

// Indices returns the indices array for this geometry.
func (g *Geometry) Indices() math32.ArrayU32 {

	return g.indices
}

// SetVAO sets the Vertex Array Object handle associated with this geometry.
func (g *Geometry) SetVAO(handle uint32) {

	g.handleVAO = handle
}

// VAO returns the Vertex Array Object handle associated with this geometry.
func (g *Geometry) VAO() uint32 {

	return g.handleVAO
}

// AddVBO adds a Vertex Buffer Object for this geometry.
func (g *Geometry) AddVBO(vbo *gls.VBO) {

	// Check that the provided VBO doesn't have conflicting attributes with existing VBOs
	for _, existingVbo := range g.vbos {
		for _, attrib := range vbo.Attributes() {
			if existingVbo.AttribName(attrib.Name) != nil {
				panic("Geometry.AddVBO: geometry already has a VBO with attribute name:" + attrib.Name)
			}
			if attrib.Type != gls.Undefined && existingVbo.Attrib(attrib.Type) != nil {
				panic("Geometry.AddVBO: geometry already has a VBO with attribute type:" + strconv.Itoa(int(attrib.Type)))
			}
		}
	}

	g.vbos = append(g.vbos, vbo)
}

// VBO returns a pointer to this geometry's VBO which contain the specified attribute.
// Returns nil if the VBO is not found.
func (g *Geometry) VBO(atype gls.AttribType) *gls.VBO {

	for _, vbo := range g.vbos {
		if vbo.Attrib(atype) != nil {
			return vbo
		}
	}
	return nil
}

// VBOName returns a pointer to this geometry's VBO which contain the specified attribute.
// Returns nil if the VBO is not found.
func (g *Geometry) VBOName(name string) *gls.VBO {

	for _, vbo := range g.vbos {
		if vbo.AttribName(name) != nil {
			return vbo
		}
	}
	return nil
}

// VBOs returns all of this geometry's VBOs.
func (g *Geometry) VBOs() []*gls.VBO {

	return g.vbos
}

// Items returns the number of items in the first VBO.
// (The number of items should be same for all VBOs)
// An item is a complete group of attributes in the VBO buffer.
func (g *Geometry) Items() int {

	if len(g.vbos) == 0 {
		return 0
	}
	vbo := g.vbos[0]
	if vbo.AttribCount() == 0 {
		return 0
	}
	return vbo.Buffer().Bytes() / vbo.StrideSize()
}

// SetAttributeName sets the name of the VBO attribute associated with the provided attribute type.
func (g *Geometry) SetAttributeName(atype gls.AttribType, attribName string) {

	vbo := g.VBO(atype)
	if vbo != nil {
		vbo.Attrib(atype).Name = attribName
	}
}

// AttributeName returns the name of the VBO attribute associated with the provided attribute type.
func (g *Geometry) AttributeName(atype gls.AttribType) string {

	return g.VBO(atype).Attrib(atype).Name
}

// OperateOnVertices iterates over all the vertices and calls
// the specified callback function with a pointer to each vertex.
// The vertex pointers can be modified inside the callback and
// the modifications will be applied to the buffer at each iteration.
// The callback function returns false to continue or true to break.
func (g *Geometry) OperateOnVertices(cb func(vertex *math32.Vector3) bool) {

	// Get buffer with position vertices
	vbo := g.VBO(gls.VertexPosition)
	if vbo == nil {
		return
	}
	vbo.OperateOnVectors3(gls.VertexPosition, cb)

	// Geometric properties may have changed
	g.boundingBoxValid = false
	g.boundingSphereValid = false
	g.areaValid = false
	g.volumeValid = false
	g.rotInertiaValid = false
}

// ReadVertices iterates over all the vertices and calls
// the specified callback function with the value of each vertex.
// The callback function returns false to continue or true to break.
func (g *Geometry) ReadVertices(cb func(vertex math32.Vector3) bool) {

	// Get buffer with position vertices
	vbo := g.VBO(gls.VertexPosition)
	if vbo == nil {
		return
	}
	vbo.ReadVectors3(gls.VertexPosition, cb)
}

// OperateOnVertexNormals iterates over all the vertex normals
// and calls the specified callback function with a pointer to each normal.
// The vertex pointers can be modified inside the callback and
// the modifications will be applied to the buffer at each iteration.
// The callback function returns false to continue or true to break.
func (g *Geometry) OperateOnVertexNormals(cb func(normal *math32.Vector3) bool) {

	// Get buffer with position vertices
	vbo := g.VBO(gls.VertexNormal)
	if vbo == nil {
		return
	}
	vbo.OperateOnVectors3(gls.VertexNormal, cb)
}

// ReadVertexNormals iterates over all the vertex normals and calls
// the specified callback function with the value of each normal.
// The callback function returns false to continue or true to break.
func (g *Geometry) ReadVertexNormals(cb func(vertex math32.Vector3) bool) {

	// Get buffer with position vertices
	vbo := g.VBO(gls.VertexNormal)
	if vbo == nil {
		return
	}
	vbo.ReadVectors3(gls.VertexNormal, cb)
}

// ReadFaces iterates over all the vertices and calls
// the specified callback function with face-forming vertex triples.
// The callback function returns false to continue or true to break.
func (g *Geometry) ReadFaces(cb func(vA, vB, vC math32.Vector3) bool) {

	// Get buffer with position vertices
	vbo := g.VBO(gls.VertexPosition)
	if vbo == nil {
		return
	}

	// If geometry has indexed vertices need to loop over indexes
	if g.Indexed() {
		var vA, vB, vC math32.Vector3
		positions := vbo.Buffer()
		for i := 0; i < g.indices.Size(); i += 3 {
			// Get face vertices
			positions.GetVector3(int(3*g.indices[i]), &vA)
			positions.GetVector3(int(3*g.indices[i+1]), &vB)
			positions.GetVector3(int(3*g.indices[i+2]), &vC)
			// Call callback with face vertices
			brk := cb(vA, vB, vC)
			if brk {
				break
			}
		}
	} else {
		// Geometry does NOT have indexed vertices - can read vertices in sequence
		vbo.ReadTripleVectors3(gls.VertexPosition, cb)
	}
}

// TODO Read and Operate on Texcoords, Faces, Edges, FaceNormals, etc...

// Indexed returns whether the geometry is indexed or not.
func (g *Geometry) Indexed() bool {

	return g.indices.Size() > 0
}

// BoundingBox computes the bounding box of the geometry if necessary
// and returns is value.
func (g *Geometry) BoundingBox() math32.Box3 {

	// If valid, return its value
	if g.boundingBoxValid {
		return g.boundingBox
	}

	// Reset bounding box
	g.boundingBox.Min.Set(0, 0, 0)
	g.boundingBox.Max.Set(0, 0, 0)

	// Expand bounding box by each vertex
	g.ReadVertices(func(vertex math32.Vector3) bool {
		g.boundingBox.ExpandByPoint(&vertex)
		return false
	})
	g.boundingBoxValid = true
	return g.boundingBox
}

// BoundingSphere computes the bounding sphere of this geometry
// if necessary and returns its value.
func (g *Geometry) BoundingSphere() math32.Sphere {

	// If valid, return its value
	if g.boundingSphereValid {
		return g.boundingSphere
	}

	// Reset radius, calculate bounding box and copy center
	g.boundingSphere.Radius = float32(0)
	box := g.BoundingBox()
	box.Center(&g.boundingSphere.Center)

	// Find the radius of the bounding sphere
	maxRadiusSq := float32(0)
	g.ReadVertices(func(vertex math32.Vector3) bool {
		maxRadiusSq = math32.Max(maxRadiusSq, g.boundingSphere.Center.DistanceToSquared(&vertex))
		return false
	})
	g.boundingSphere.Radius = float32(math32.Sqrt(maxRadiusSq))
	g.boundingSphereValid = true
	return g.boundingSphere
}

// Area returns the surface area.
// NOTE: This only works for triangle-based meshes.
func (g *Geometry) Area() float32 {

	// If valid, return its value
	if g.areaValid {
		return g.area
	}

	// Reset area
	g.area = 0

	// Sum area of all triangles
	g.ReadFaces(func(vA, vB, vC math32.Vector3) bool {
		vA.Sub(&vC)
		vB.Sub(&vC)
		vC.CrossVectors(&vA, &vB)
		g.area += vC.Length() / 2.0
		return false
	})
	g.areaValid = true
	return g.area
}

// Volume returns the volume.
// NOTE: This only works for closed triangle-based meshes.
func (g *Geometry) Volume() float32 {

	// If valid, return its value
	if g.volumeValid {
		return g.volume
	}

	// Reset volume
	g.volume = 0

	// Calculate volume of all tetrahedrons
	g.ReadFaces(func(vA, vB, vC math32.Vector3) bool {
		vA.Sub(&vC)
		vB.Sub(&vC)
		g.volume += vC.Dot(vA.Cross(&vB)) / 6.0
		return false
	})
	g.volumeValid = true
	return g.volume
}

// RotationalInertia returns the rotational inertia tensor, also known as the moment of inertia.
// This assumes constant density of 1 (kg/m^2).
// To adjust for a different constant density simply scale the returning matrix by the density.
func (g *Geometry) RotationalInertia(mass float32) math32.Matrix3 {

	// If valid, return its value
	if g.rotInertiaValid {
		return g.rotInertia
	}

	// Reset rotational inertia
	g.rotInertia.Zero()

	// For now approximate result based on bounding box
	b := math32.NewVec3()
	box := g.BoundingBox()
	box.Size(b)
	multiplier := mass / 12.0

	x := (b.Y*b.Y + b.Z*b.Z) * multiplier
	y := (b.X*b.X + b.Z*b.Z) * multiplier
	z := (b.Y*b.Y + b.X*b.X) * multiplier

	g.rotInertia.Set(
		x, 0, 0,
		0, y, 0,
		0, 0, z,
	)
	return g.rotInertia
}

// ProjectOntoAxis projects the geometry onto the specified axis,
// effectively squashing it into a line passing through the local origin.
// Returns the maximum and the minimum values on that line (i.e. signed distances from the local origin).
func (g *Geometry) ProjectOntoAxis(localAxis *math32.Vector3) (float32, float32) {

	var max, min float32
	g.ReadVertices(func(vertex math32.Vector3) bool {
		val := vertex.Dot(localAxis)
		if val > max {
			max = val
		}
		if val < min {
			min = val
		}
		return false
	})
	return max, min
}

// TODO:
// https://stackoverflow.com/questions/21640545/how-to-check-for-convexity-of-a-3d-mesh
// func (g *Geometry) IsConvex() bool {
//
// {

// ApplyMatrix multiplies each of the geometry position vertices
// by the specified matrix and apply the correspondent normal
// transform matrix to the geometry normal vectors.
// The geometry's bounding box and sphere are recomputed if needed.
func (g *Geometry) ApplyMatrix(m *math32.Matrix4) {

	// Apply matrix to all vertices
	g.OperateOnVertices(func(vertex *math32.Vector3) bool {
		vertex.ApplyMatrix4(m)
		return false
	})

	// Apply normal matrix to all normal vectors
	var normalMatrix math32.Matrix3
	normalMatrix.GetNormalMatrix(m)
	g.OperateOnVertexNormals(func(normal *math32.Vector3) bool {
		normal.ApplyMatrix3(&normalMatrix).Normalize()
		return false
	})
}

// RenderSetup is called by the renderer before drawing the geometry.
func (g *Geometry) RenderSetup(gs *gls.GLS) {

	// First time initialization
	if g.gs == nil {
		if g.handleVAO == 0 {
			// Generate VAO and bind it
			g.handleVAO = gs.GenVertexArray()
		}
		// Generate VBO for indices
		g.handleIndices = gs.GenBuffer()
		// Save pointer to gs indicating initialization was done
		g.gs = gs
	}

	// Update VBOs
	gs.BindVertexArray(g.handleVAO)
	for _, vbo := range g.vbos {
		vbo.Transfer(gs)
	}

	// Update Indices buffer if necessary
	if g.indices.Size() > 0 && g.updateIndices {
		gs.BindBuffer(gls.ELEMENT_ARRAY_BUFFER, g.handleIndices)
		gs.BufferData(gls.ELEMENT_ARRAY_BUFFER, g.indices.Bytes(), g.indices, gls.STATIC_DRAW)
		g.updateIndices = false
	}
}
