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

type Mesh struct {
	Graphic             // Embedded graphic
	uniMVM  gls.Uniform // Model view matrix uniform location cache
	uniMVPM gls.Uniform // Model view projection matrix uniform cache
	uniNM   gls.Uniform // Normal matrix uniform cache
}

// NewMesh creates and returns a pointer to a mesh with the specified geometry and material
// If the mesh has multi materials, the material specified here must be nil and the
// individual materials must be add using "AddMateria" or AddGroupMaterial"
func NewMesh(igeom geometry.IGeometry, imat material.IMaterial) *Mesh {

	m := new(Mesh)
	m.Init(igeom, imat)
	return m
}

func (m *Mesh) Init(igeom geometry.IGeometry, imat material.IMaterial) {

	m.Graphic.Init(igeom, gls.TRIANGLES)

	// Initialize uniforms
	m.uniMVM.Init("ModelViewMatrix")
	m.uniMVPM.Init("MVP")
	m.uniNM.Init("NormalMatrix")

	// Adds single material if not nil
	if imat != nil {
		m.AddMaterial(imat, 0, 0)
	}
}

func (m *Mesh) AddMaterial(imat material.IMaterial, start, count int) {

	m.Graphic.AddMaterial(m, imat, start, count)
}

// Add group material
func (m *Mesh) AddGroupMaterial(imat material.IMaterial, gindex int) {

	m.Graphic.AddGroupMaterial(m, imat, gindex)
}

// RenderSetup is called by the engine before drawing the mesh geometry
// It is responsible to updating the current shader uniforms with
// the model matrices.
func (m *Mesh) RenderSetup(gs *gls.GLS, rinfo *core.RenderInfo) {

	// Calculates model view matrix and transfer uniform
	mw := m.MatrixWorld()
	var mvm math32.Matrix4
	mvm.MultiplyMatrices(&rinfo.ViewMatrix, &mw)
	location := m.uniMVM.Location(gs)
	gs.UniformMatrix4fv(location, 1, false, &mvm[0])

	// Calculates model view projection matrix and updates uniform
	var mvpm math32.Matrix4
	mvpm.MultiplyMatrices(&rinfo.ProjMatrix, &mvm)
	location = m.uniMVPM.Location(gs)
	gs.UniformMatrix4fv(location, 1, false, &mvpm[0])

	// Calculates normal matrix and updates uniform
	var nm math32.Matrix3
	nm.GetNormalMatrix(&mvm)
	location = m.uniNM.Location(gs)
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
	inverseMatrix.GetInverse(&matrixWorld, true)
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

	// Get buffer with position vertices
	vboPos := geom.VBO("VertexPosition")
	if vboPos == nil {
		panic("mesh.Raycast(): VertexPosition VBO not found")
	}
	positions := vboPos.Buffer()
	indices := geom.Indices()

	var vA math32.Vector3
	var vB math32.Vector3
	var vC math32.Vector3

	// Geometry has indexed vertices
	if indices.Size() > 0 {
		for i := 0; i < indices.Size(); i += 3 {
			// Get face indices
			a := indices[i]
			b := indices[i+1]
			c := indices[i+2]
			// Get face position vectors
			positions.GetVector3(int(3*a), &vA)
			positions.GetVector3(int(3*b), &vB)
			positions.GetVector3(int(3*c), &vC)
			// Checks intersection of the ray with this face
			mat := m.GetMaterial(i).GetMaterial()
			var point math32.Vector3
			intersect := checkIntersection(mat, &vA, &vB, &vC, &point)
			if intersect != nil {
				intersect.Index = uint32(i)
				*intersects = append(*intersects, *intersect)
			}
		}
		// Geometry has NO indexed vertices
	} else {
		for i := 0; i < positions.Size(); i += 9 {
			// Get face indices
			a := i / 3
			b := a + 1
			c := a + 2
			// Set face position vectors
			positions.GetVector3(int(3*a), &vA)
			positions.GetVector3(int(3*b), &vB)
			positions.GetVector3(int(3*c), &vC)
			// Checks intersection of the ray with this face
			mat := m.GetMaterial(i).GetMaterial()
			var point math32.Vector3
			intersect := checkIntersection(mat, &vA, &vB, &vC, &point)
			if intersect != nil {
				intersect.Index = uint32(a)
				*intersects = append(*intersects, *intersect)
			}
		}
	}
}
