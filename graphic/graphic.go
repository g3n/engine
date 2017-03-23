// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphic

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/material"
)

// Graphic is a Node which has a visible representation in the scene.
// It has an associated geometry and one or more materials.
// It is the base type used by other graphics such as lines, line_strip,
// points and meshes.
type Graphic struct {
	core.Node                     // Embedded Node
	igeom      geometry.IGeometry // Associated IGeometry
	materials  []GraphicMaterial  // Materials
	mode       uint32             // OpenGL primitive
	renderable bool               // Renderable flag
}

// GraphicMaterial specifies the material to be used for
// a subset of vertices from the Graphic geometry
// A Graphic object has at least one GraphicMaterial
type GraphicMaterial struct {
	imat     material.IMaterial // Associated material
	start    int                // Index of first element in the geometry
	count    int                // Number of elements
	igraphic IGraphic           // Graphic which contains this GraphicMaterial
}

// Interface for all Graphics
type IGraphic interface {
	core.INode
	GetGraphic() *Graphic
	GetGeometry() *geometry.Geometry
	Renderable() bool
	SetRenderable(bool)
	RenderSetup(gs *gls.GLS, rinfo *core.RenderInfo)
}

// NewGraphic creates and returns a pointer to a new graphic object with
// the specified geometry and OpenGL primitive.
// The created graphic object, though, has not materials.
func NewGraphic(igeom geometry.IGeometry, mode uint32) *Graphic {

	gr := new(Graphic)
	return gr.Init(igeom, mode)
}

// Init initializes a Graphic type embedded in another type
// with the specified geometry and OpenGL mode.
func (gr *Graphic) Init(igeom geometry.IGeometry, mode uint32) *Graphic {

	gr.Node.Init()
	gr.igeom = igeom
	gr.mode = mode
	gr.materials = make([]GraphicMaterial, 0)
	gr.renderable = true
	return gr
}

// GetGraphic satisfies the IGraphic interface and
// returns pointer to the base Graphic
func (gr *Graphic) GetGraphic() *Graphic {

	return gr
}

// GetGeometry satisfies the IGraphic interface and returns
// a pointer to the geometry associated with this graphic
func (gr *Graphic) GetGeometry() *geometry.Geometry {

	return gr.igeom.GetGeometry()
}

// Dispose overrides the embedded Node Dispose method
func (gr *Graphic) Dispose() {

	gr.igeom.Dispose()
	for i := 0; i < len(gr.materials); i++ {
		gr.materials[i].imat.Dispose()
	}
}

// SetRenderable satisfies the IGraphic interface and
// sets the renderable state of this Graphic (default = true)
func (gr *Graphic) SetRenderable(state bool) {

	gr.renderable = state
}

// Renderable satisfies the IGraphic interface and
// returns the renderable state of this graphic
func (gr *Graphic) Renderable() bool {

	return gr.renderable
}

// Add material for the specified subset of vertices.
// If the material applies to all vertices, start and count must be 0.
func (gr *Graphic) AddMaterial(igr IGraphic, imat material.IMaterial, start, count int) {

	gmat := GraphicMaterial{
		imat:     imat,
		start:    start,
		count:    count,
		igraphic: igr,
	}
	gr.materials = append(gr.materials, gmat)
}

// Add group material
func (gr *Graphic) AddGroupMaterial(igr IGraphic, imat material.IMaterial, gindex int) {

	geom := gr.igeom.GetGeometry()
	if gindex < 0 || gindex >= geom.GroupCount() {
		panic("Invalid group index")
	}
	group := geom.GroupAt(gindex)
	gr.AddMaterial(igr, imat, group.Start, group.Count)
}

// Materials returns slice with this graphic materials
func (gr *Graphic) Materials() []GraphicMaterial {

	return gr.materials
}

// GetMaterial returns the  material associated with the specified vertex position
func (gr *Graphic) GetMaterial(vpos int) material.IMaterial {

	for _, gmat := range gr.materials {
		// Case for unimaterial
		if gmat.count == 0 {
			return gmat.imat
		}
		if gmat.start >= vpos && gmat.start+gmat.count <= vpos {
			return gmat.imat
		}
	}
	return nil
}

func (grmat *GraphicMaterial) GetMaterial() material.IMaterial {

	return grmat.imat
}

// Render is called by the renderer to render this graphic material
func (grmat *GraphicMaterial) Render(gs *gls.GLS, rinfo *core.RenderInfo) {

	// Setup the associated material (set states and transfer material uniforms and textures)
	grmat.imat.RenderSetup(gs)

	// Setup the associated geometry (set VAO and transfer VBOS)
	gr := grmat.igraphic.GetGraphic()
	gr.igeom.RenderSetup(gs)

	// Setup current graphic (transfer matrices)
	grmat.igraphic.RenderSetup(gs, rinfo)

	// Get the number of vertices for the current material
	count := grmat.count

	geom := gr.igeom.GetGeometry()
	indices := geom.Indices()
	// Indexed geometry
	if indices.Size() > 0 {
		if count == 0 {
			count = indices.Size()
		}
		gs.DrawElements(gr.mode, int32(count), gls.UNSIGNED_INT, 4*uint32(grmat.start))
		// Non indexed geometry
	} else {
		if count == 0 {
			count = geom.Items()
		}
		gs.DrawArrays(gr.mode, int32(grmat.start), int32(count))
	}
}
