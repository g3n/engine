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

// Sprite is a potentially animated image positioned in space that always faces the camera.
type Sprite struct {
	Graphic             // Embedded graphic
	uniMVPM gls.Uniform // Model view projection matrix uniform location cache
}

// NewSprite creates and returns a pointer to a sprite with the specified dimensions and material
func NewSprite(width, height float32, imat material.IMaterial) *Sprite {

	s := new(Sprite)

	// Creates geometry
	geom := geometry.NewGeometry()
	w := width / 2
	h := height / 2

	// Builds array with vertex positions and texture coordinates
	positions := math32.NewArrayF32(0, 12)
	positions.Append(
		-w, -h, 0, 0, 0,
		w, -h, 0, 1, 0,
		w, h, 0, 1, 1,
		-w, h, 0, 0, 1,
	)
	// Builds array of indices
	indices := math32.NewArrayU32(0, 6)
	indices.Append(0, 1, 2, 0, 2, 3)

	// Set geometry buffers
	geom.SetIndices(indices)
	geom.AddVBO(
		gls.NewVBO(positions).
			AddAttrib(gls.VertexPosition).
			AddAttrib(gls.VertexTexcoord),
	)

	s.Graphic.Init(geom, gls.TRIANGLES)
	s.AddMaterial(s, imat, 0, 0)

	s.uniMVPM.Init("MVP")
	return s
}

// RenderSetup sets up the rendering of the sprite.
func (s *Sprite) RenderSetup(gs *gls.GLS, rinfo *core.RenderInfo) {

	// Calculates model view matrix
	mw := s.MatrixWorld()
	var mvm math32.Matrix4
	mvm.MultiplyMatrices(&rinfo.ViewMatrix, &mw)

	// Decomposes model view matrix
	var position math32.Vector3
	var quaternion math32.Quaternion
	var scale math32.Vector3
	mvm.Decompose(&position, &quaternion, &scale)

	// Removes any rotation in X and Y axes and compose new model view matrix
	rotation := s.Rotation()
	rotation.X = 0
	rotation.Y = 0
	quaternion.SetFromEuler(&rotation)
	var mvmNew math32.Matrix4
	mvmNew.Compose(&position, &quaternion, &scale)

	// Calculates final MVP and updates uniform
	var mvpm math32.Matrix4
	mvpm.MultiplyMatrices(&rinfo.ProjMatrix, &mvmNew)
	location := s.uniMVPM.Location(gs)
	gs.UniformMatrix4fv(location, 1, false, &mvpm[0])
}

// Raycast checks intersections between this geometry and the specified raycaster
// and if any found appends it to the specified intersects array.
func (s *Sprite) Raycast(rc *core.Raycaster, intersects *[]core.Intersect) {

	// Copy and convert ray to camera coordinates
	var ray math32.Ray
	ray.Copy(&rc.Ray).ApplyMatrix4(&rc.ViewMatrix)

	// Calculates ViewMatrix * MatrixWorld
	var mv math32.Matrix4
	matrixWorld := s.MatrixWorld()
	mv.MultiplyMatrices(&rc.ViewMatrix, &matrixWorld)

	// Decompose transformation matrix in its components
	var position math32.Vector3
	var quaternion math32.Quaternion
	var scale math32.Vector3
	mv.Decompose(&position, &quaternion, &scale)

	// Remove any rotation in X and Y axis and
	// compose new transformation matrix
	rotation := s.Rotation()
	rotation.X = 0
	rotation.Y = 0
	quaternion.SetFromEuler(&rotation)
	mv.Compose(&position, &quaternion, &scale)

	// Get buffer with vertices and uvs
	geom := s.GetGeometry()
	vboPos := geom.VBO(gls.VertexPosition)
	if vboPos == nil {
		panic("sprite.Raycast(): VertexPosition VBO not found")
	}
	// Get vertex positions, transform to camera coordinates and
	// checks intersection with ray
	buffer := vboPos.Buffer()
	indices := geom.Indices()
	var v1 math32.Vector3
	var v2 math32.Vector3
	var v3 math32.Vector3
	var point math32.Vector3
	intersect := false
	for i := 0; i < indices.Size(); i += 3 {
		pos := indices[i]
		buffer.GetVector3(int(pos*5), &v1)
		v1.ApplyMatrix4(&mv)
		pos = indices[i+1]
		buffer.GetVector3(int(pos*5), &v2)
		v2.ApplyMatrix4(&mv)
		pos = indices[i+2]
		buffer.GetVector3(int(pos*5), &v3)
		v3.ApplyMatrix4(&mv)
		if ray.IntersectTriangle(&v1, &v2, &v3, false, &point) {
			intersect = true
			break
		}
	}
	if !intersect {
		return
	}
	// Get distance from intersection point
	origin := ray.Origin()
	distance := origin.DistanceTo(&point)

	// Checks if distance is between the bounds of the raycaster
	if distance < rc.Near || distance > rc.Far {
		return
	}

	// Appends intersection to received parameter.
	*intersects = append(*intersects, core.Intersect{
		Distance: distance,
		Point:    point,
		Object:   s,
	})
}
