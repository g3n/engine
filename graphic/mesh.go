// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphic

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
)

// Mesh is a Graphic with uniforms for the model, view, projection, and normal matrices.
type Mesh struct {
	Graphic             // Embedded graphic
	uniMm  gls.Uniform  // Model matrix uniform location cache
	uniMVm  gls.Uniform // Model view matrix uniform location cache
	uniMVPm gls.Uniform // Model view projection matrix uniform cache
	uniNm   gls.Uniform // Normal matrix uniform cache
}

// NewMesh creates and returns a pointer to a mesh with the specified geometry and material.
// If the mesh has multi materials, the material specified here must be nil and the
// individual materials must be add using "AddMaterial" or AddGroupMaterial".
func NewMesh(igeom geometry.IGeometry, imat material.IMaterial) *Mesh {

	m := new(Mesh)
	m.Init(igeom, imat)
	return m
}

// Init initializes the Mesh and its uniforms.
func (m *Mesh) Init(igeom geometry.IGeometry, imat material.IMaterial) {

	m.Graphic.Init(igeom, gls.TRIANGLES)

	// Initialize uniforms
	m.uniMm.Init("ModelMatrix")
	m.uniMVm.Init("ModelViewMatrix")
	m.uniMVPm.Init("MVP")
	m.uniNm.Init("NormalMatrix")

	// Adds single material if not nil
	if imat != nil {
		m.AddMaterial(imat, 0, 0)
	}
}

// SetMaterial clears all materials and adds the specified material for all vertices.
func (m *Mesh) SetMaterial(imat material.IMaterial) {

	m.Graphic.ClearMaterials()
	m.Graphic.AddMaterial(m, imat, 0, 0)
}

// AddMaterial adds a material for the specified subset of vertices.
func (m *Mesh) AddMaterial(imat material.IMaterial, start, count int) {

	m.Graphic.AddMaterial(m, imat, start, count)
}

// AddGroupMaterial adds a material for the specified geometry group.
func (m *Mesh) AddGroupMaterial(imat material.IMaterial, gindex int) {

	m.Graphic.AddGroupMaterial(m, imat, gindex)
}

// Clone clones the mesh and satisfies the INode interface.
func (m *Mesh) Clone() core.INode {

	clone := new(Mesh)
	clone.Graphic = *m.Graphic.Clone().(*Graphic)
	clone.SetIGraphic(clone)

	// Initialize uniforms
	clone.uniMm.Init("ModelMatrix")
	clone.uniMVm.Init("ModelViewMatrix")
	clone.uniMVPm.Init("MVP")
	clone.uniNm.Init("NormalMatrix")

	return clone
}

// RenderSetup is called by the engine before drawing the mesh geometry
// It is responsible to updating the current shader uniforms with
// the model matrices.
func (m *Mesh) RenderSetup(gs *gls.GLS, rinfo *core.RenderInfo) {

	// Transfer uniform for model matrix
	mm := m.ModelMatrix()
	location := m.uniMm.Location(gs)
	gs.UniformMatrix4fv(location, 1, false, &mm[0])

	// Transfer uniform for model view matrix
	mvm := m.ModelViewMatrix()
	location = m.uniMVm.Location(gs)
	gs.UniformMatrix4fv(location, 1, false, &mvm[0])

	// Transfer uniform for model view projection matrix
	mvpm := m.ModelViewProjectionMatrix()
	location = m.uniMVPm.Location(gs)
	gs.UniformMatrix4fv(location, 1, false, &mvpm[0])

	// Calculates normal matrix and transfer uniform
	var nm math32.Matrix3
	nm.GetNormalMatrix(mvm)
	location = m.uniNm.Location(gs)
	gs.UniformMatrix3fv(location, 1, false, &nm[0])
}

// Raycast checks intersections between this geometry and the specified raycaster
// and if any found appends it to the specified intersects array.
func (m *Mesh) Raycast(rc *core.Raycaster, intersects *[]core.Intersect) {

	// Transform this mesh geometry bounding sphere from model
	// to world coordinates and checks intersection with raycaster
	geom := m.GetGeometry()
	sphere := geom.BoundingSphere()
	matrixWorld := m.MatrixWorld()
	sphere.ApplyMatrix4(&matrixWorld)
	if !rc.IsIntersectionSphere(&sphere) {
		return
	}

	// Copy ray and transform to model coordinates
	// This ray will will also be used to check intersects with
	// the geometry, as is much less expensive to transform the
	// ray to model coordinates than the geometry to world coordinates.
	var inverseMatrix math32.Matrix4
	inverseMatrix.GetInverse(&matrixWorld)
	var ray math32.Ray
	ray.Copy(&rc.Ray).ApplyMatrix4(&inverseMatrix)
	bbox := geom.BoundingBox()
	if !ray.IsIntersectionBox(&bbox) {
		return
	}

	// Local function to check the intersection of the ray from the raycaster with
	// the specified face defined by three poins.
	checkIntersection := func(mat *material.Material, pA, pB, pC, point *math32.Vector3) *core.Intersect {

		var intersect bool
		switch mat.Side() {
		case material.SideBack:
			intersect = ray.IntersectTriangle(pC, pB, pA, true, point)
		case material.SideFront:
			intersect = ray.IntersectTriangle(pA, pB, pC, true, point)
		case material.SideDouble:
			intersect = ray.IntersectTriangle(pA, pB, pC, false, point)
		}
		if !intersect {
			return nil
		}

		// Transform intersection point from model to world coordinates
		var intersectionPointWorld = *point
		intersectionPointWorld.ApplyMatrix4(&matrixWorld)

		// Calculates the distance from the ray origin to intersection point
		origin := rc.Ray.Origin()
		distance := origin.DistanceTo(&intersectionPointWorld)

		// Checks if distance is between the bounds of the raycaster
		if distance < rc.Near || distance > rc.Far {
			return nil
		}

		return &core.Intersect{
			Distance: distance,
			Point:    intersectionPointWorld,
			Object:   m,
		}
	}

	i := 0
	geom.ReadFaces(func(vA, vB, vC math32.Vector3) bool {
		// Checks intersection of the ray with this face
		mat := m.GetMaterial(i).GetMaterial()
		var point math32.Vector3
		intersect := checkIntersection(mat, &vA, &vB, &vC, &point)
		if intersect != nil {
			intersect.Index = uint32(i)
			*intersects = append(*intersects, *intersect)
		}
		i += 3
		return false
	})
}
