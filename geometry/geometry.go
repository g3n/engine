// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package geometry

import (
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/math32"
)

// Interface for all geometries
type IGeometry interface {
	GetGeometry() *Geometry
	RenderSetup(gs *gls.GLS)
	Dispose()
}

type Geometry struct {
	refcount            int             // Current number of references
	vbos                []*gls.VBO      // Array of VBOs
	groups              []Group         // Array geometry groups
	indices             math32.ArrayU32 // Buffer with indices
	gs                  *gls.GLS        // Pointer to gl context. Valid after first render setup
	handleVAO           uint32          // Handle to OpenGL VAO
	handleIndices       uint32          // Handle to OpenGL buffer for indices
	updateIndices       bool            // Flag to indicate that indices must be transferred
	boundingBox         math32.Box3     // Last calculated bounding box
	boundingBoxValid    bool            // Indicates if last calculated bounding box is valid
	boundingSphere      math32.Sphere   // Last calculated bounding sphere
	boundingSphereValid bool            // Indicates if last calculated bounding sphere is valid
}

// Geometry group object
type Group struct {
	Start    int    // Index of first element of the group
	Count    int    // Number of elements in the group
	Matindex int    // Material index for this group
	Matid    string // Material id used when loading external models
}

func NewGeometry() *Geometry {

	g := new(Geometry)
	g.Init()
	return g
}

// Init initializes the geometry
func (g *Geometry) Init() {

	g.refcount = 1
	g.vbos = make([]*gls.VBO, 0)
	g.groups = make([]Group, 0)
	g.gs = nil
	g.handleVAO = 0
	g.handleIndices = 0
	g.updateIndices = true
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
// if necessary releases OpenGL resources, C memory
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

func (g *Geometry) GetGeometry() *Geometry {

	return g
}

// AddGroup adds a geometry group (for multimaterial)
func (g *Geometry) AddGroup(start, count, matIndex int) *Group {

	g.groups = append(g.groups, Group{start, count, matIndex, ""})
	return &g.groups[len(g.groups)-1]
}

// AddGroupList adds the specified list of groups to this geometry
func (g *Geometry) AddGroupList(groups []Group) {

	for _, group := range groups {
		g.groups = append(g.groups, group)
	}
}

// GroupCount returns the number of geometry groups (for multimaterial)
func (g *Geometry) GroupCount() int {

	return len(g.groups)
}

// GroupAt returns pointer to geometry group at the specified index
func (g *Geometry) GroupAt(idx int) *Group {

	return &g.groups[idx]
}

// SetIndices sets the indices array for this geometry
func (g *Geometry) SetIndices(indices math32.ArrayU32) {

	g.indices = indices
	g.boundingBoxValid = false
	g.boundingSphereValid = false
}

// Indices returns this geometry indices array
func (g *Geometry) Indices() math32.ArrayU32 {

	return g.indices
}

// AddVBO adds a Vertex Buffer Object for this geometry
func (g *Geometry) AddVBO(vbo *gls.VBO) {

	g.vbos = append(g.vbos, vbo)
}

// VBO returns a pointer to this geometry VBO for the specified attribute.
// Returns nil if the VBO is not found.
func (g *Geometry) VBO(attrib string) *gls.VBO {

	for _, vbo := range g.vbos {
		if vbo.Attrib(attrib) != nil {
			return vbo
		}
	}
	return nil
}

// Returns the number of items in the first VBO
// (The number of items should be same for all VBOs)
// An item is a complete group of attributes in the VBO buffer
func (g *Geometry) Items() int {

	if len(g.vbos) == 0 {
		return 0
	}
	vbo := g.vbos[0]
	if vbo.AttribCount() == 0 {
		return 0
	}
	return vbo.Buffer().Bytes() / vbo.Stride()
}

// BoundingBox computes the bounding box of the geometry if necessary
// and returns is value
func (g *Geometry) BoundingBox() math32.Box3 {

	// If valid, returns its value
	if g.boundingBoxValid {
		return g.boundingBox
	}

	// Get buffer with position vertices
	vbPos := g.VBO("VertexPosition")
	if vbPos == nil {
		return g.boundingBox
	}
	positions := vbPos.Buffer()

	// Calculates bounding box
	var vertex math32.Vector3
	g.boundingBox.Min.Set(0, 0, 0)
	g.boundingBox.Max.Set(0, 0, 0)
	for i := 0; i < positions.Size(); i += 3 {
		positions.GetVector3(i, &vertex)
		g.boundingBox.ExpandByPoint(&vertex)
	}
	g.boundingBoxValid = true
	return g.boundingBox
}

// BoundingSphere computes the bounding sphere of this geometry
// if necessary and returns its value.
func (g *Geometry) BoundingSphere() math32.Sphere {

	// if valid, returns its value
	if g.boundingSphereValid {
		return g.boundingSphere
	}

	// Get buffer with position vertices
	vbPos := g.VBO("VertexPosition")
	if vbPos == nil {
		return g.boundingSphere
	}
	positions := vbPos.Buffer()

	// Get/calculates the bounding box
	box := g.BoundingBox()

	// Sets the center of the bounding sphere the same as the center of the bounding box.
	box.Center(&g.boundingSphere.Center)
	center := g.boundingSphere.Center

	// Find the radius of the bounding sphere
	maxRadiusSq := float32(0.0)
	for i := 0; i < positions.Size(); i += 3 {
		var vertex math32.Vector3
		positions.GetVector3(i, &vertex)
		maxRadiusSq = math32.Max(maxRadiusSq, center.DistanceToSquared(&vertex))
	}
	radius := math32.Sqrt(maxRadiusSq)
	if math32.IsNaN(radius) {
		panic("geometry.BoundingSphere: computed radius is NaN")
	}
	g.boundingSphere.Radius = float32(radius)
	g.boundingSphereValid = true
	return g.boundingSphere
}

// ApplyMatrix multiplies each of the geometry position vertices
// by the specified matrix and apply the correspondent normal
// transform matrix to the geometry normal vectors.
// The geometry's bounding box and sphere are recomputed if needed.
func (g *Geometry) ApplyMatrix(m *math32.Matrix4) {

	// Get positions buffer
	vboPos := g.VBO("VertexPosition")
	if vboPos == nil {
		return
	}
	positions := vboPos.Buffer()
	// Apply matrix to all position vertices
	for i := 0; i < positions.Size(); i += 3 {
		var vertex math32.Vector3
		positions.GetVector3(i, &vertex)
		vertex.ApplyMatrix4(m)
		positions.SetVector3(i, &vertex)
	}
	vboPos.Update()

	// Get normals buffer
	vboNormals := g.VBO("VertexNormal")
	if vboNormals == nil {
		return
	}
	normals := vboNormals.Buffer()
	// Apply normal matrix to all normal vectors
	var normalMatrix math32.Matrix3
	normalMatrix.GetNormalMatrix(m)
	for i := 0; i < normals.Size(); i += 3 {
		var vertex math32.Vector3
		normals.GetVector3(i, &vertex)
		vertex.ApplyMatrix3(&normalMatrix).Normalize()
		normals.SetVector3(i, &vertex)
	}
	vboNormals.Update()
}

// RenderSetup is called by the renderer before drawing the geometry
func (g *Geometry) RenderSetup(gs *gls.GLS) {

	// First time initialization
	if g.gs == nil {
		// Generates VAO and binds it
		g.handleVAO = gs.GenVertexArray()
		gs.BindVertexArray(g.handleVAO)
		// Generates VBO for indices
		g.handleIndices = gs.GenBuffer()
		// Saves pointer to gl indicating initialization was done.
		g.gs = gs
	}

	// Update VBOs
	gs.BindVertexArray(g.handleVAO)
	for _, vbo := range g.vbos {
		vbo.Transfer(gs)
	}

	// Updates Indices buffer if necessary
	if g.indices.Size() > 0 && g.updateIndices {
		gs.BindBuffer(gls.ELEMENT_ARRAY_BUFFER, g.handleIndices)
		gs.BufferData(gls.ELEMENT_ARRAY_BUFFER, g.indices.Bytes(), g.indices, gls.STATIC_DRAW)
		g.updateIndices = false
	}
}
