// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package collision

import (
	"sort"

	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
)

// Raycaster represents an empty object that can cast rays and check for ray intersections.
type Raycaster struct {
	// The distance from the ray origin to the intersected points
	// must be greater than the value of this field to be considered.
	// The defaul value is 0.0
	Near float32
	// The distance from the ray origin to the intersected points
	// must be less than the value of this field to be considered.
	// The defaul value is +Infinity.
	Far float32
	// Minimum distance in world coordinates between the ray and
	// a line segment when checking intersects with lines.
	// The default value is 0.1
	LinePrecision float32
	// Minimum distance in world coordinates between the ray and
	// a point when checking intersects with points.
	// The default value is 0.1
	PointPrecision float32
	// This field must be set with the camera view matrix used
	// when checking for sprite intersections.
	// It is set automatically when using camera.SetRaycaster
	ViewMatrix math32.Matrix4
	// Embedded ray
	math32.Ray
}

// Intersect describes the intersection between a ray and an object
type Intersect struct {
	// Distance between the origin of the ray and the intersect
	Distance float32
	// Point of intersection in world coordinates
	Point math32.Vector3
	// Intersected node
	Object core.INode
	// If the geometry has indices, this field is the
	// index in the Indices buffer of the vertex intersected
	// or the first vertex of the intersected face.
	// If the geometry doesn't have indices, this field is the
	// index in the positions buffer of the vertex intersected
	// or the first vertex of the insersected face.
	Index uint32
}

// NewRaycaster creates and returns a pointer to a new raycaster object
// with the specified origin and direction.
func NewRaycaster(origin, direction *math32.Vector3) *Raycaster {

	rc := new(Raycaster)
	rc.Ray.Set(origin, direction)
	rc.Near = 0
	rc.Far = math32.Inf(1)
	rc.LinePrecision = 0.1
	rc.PointPrecision = 0.1
	return rc
}

// IntersectObject checks intersections between this raycaster and
// and the specified node. If recursive is true, it also checks
// the intersection with the node's children.
// Intersections are returned sorted by distance, closest first.
func (rc *Raycaster) IntersectObject(inode core.INode, recursive bool) []Intersect {

	intersects := []Intersect{}
	rc.intersectObject(inode, &intersects, recursive)
	sort.Slice(intersects, func(i, j int) bool {
		return intersects[i].Distance < intersects[j].Distance
	})
	return intersects
}

// IntersectObjects checks intersections between this raycaster and
// the specified array of scene nodes. If recursive is true, it also checks
// the intersection with each nodes' children.
// Intersections are returned sorted by distance, closest first.
func (rc *Raycaster) IntersectObjects(inodes []core.INode, recursive bool) []Intersect {

	intersects := []Intersect{}
	for _, inode := range inodes {
		rc.intersectObject(inode, &intersects, recursive)
	}
	sort.Slice(intersects, func(i, j int) bool {
		return intersects[i].Distance < intersects[j].Distance
	})
	return intersects
}

func (rc *Raycaster) intersectObject(inode core.INode, intersects *[]Intersect, recursive bool) {

	node := inode.GetNode()
	if !node.Visible() {
		return
	}

	switch in := inode.(type) {
	case *graphic.Sprite:
		rc.RaycastSprite(in, intersects)
	case *graphic.Points:
		rc.RaycastPoints(in, intersects)
	case *graphic.Mesh:
		rc.RaycastMesh(in, intersects)
	case *graphic.Lines:
		rc.RaycastLines(in, intersects)
	case *graphic.LineStrip:
		rc.RaycastLineStrip(in, intersects)
	}

	if recursive {
		for _, child := range node.Children() {
			rc.intersectObject(child, intersects, true)
		}
	}
	return
}

// SetRaycaster sets the specified raycaster with this camera position in world coordinates
// pointing to the direction defined by the specified coordinates unprojected using this camera.
func (rc *Raycaster) SetFromCamera(cam *camera.Camera, sx, sy float32) error { // TODO maybe use ICamera

	var origin, direction math32.Vector3
	matrixWorld := cam.MatrixWorld()
	origin.SetFromMatrixPosition(&matrixWorld)
	direction.Set(sx, sy, 0.5)

	unproj := cam.Unproject(&direction) // unproj = direction after this point TODO improve clarity
	unproj.Sub(&origin).Normalize()
	rc.Set(&origin, &direction)
	cam.ViewMatrix(&rc.ViewMatrix) // Update the view matrix of the raycaster
	return nil
}

// RaycastSprite checks intersections between the raycaster and the specified sprite
// and if any found appends it to the specified intersects array.
func (rc *Raycaster) RaycastSprite(s *graphic.Sprite, intersects *[]Intersect) {

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
	*intersects = append(*intersects, Intersect{
		Distance: distance,
		Point:    point,
		Object:   s,
	})
}

// RaycastPoints
func (rc *Raycaster) RaycastPoints(p *graphic.Points, intersects *[]Intersect) {

	// Checks intersection with the bounding sphere transformed to world coordinates
	geom := p.GetGeometry()
	sphere := geom.BoundingSphere()
	matrixWorld := p.MatrixWorld()
	sphere.ApplyMatrix4(&matrixWorld)
	if !rc.IsIntersectionSphere(&sphere) {
		return
	}

	// Copy ray and transforms to model coordinates
	var inverseMatrix math32.Matrix4
	var ray math32.Ray
	inverseMatrix.GetInverse(&matrixWorld)
	ray.Copy(&rc.Ray).ApplyMatrix4(&inverseMatrix)

	// Checks intersection with all points
	scale := p.Scale()
	localThreshold := rc.PointPrecision / ((scale.X + scale.Y + scale.Z) / 3)
	localThresholdSq := localThreshold * localThreshold

	// internal function to check intersection with a point
	testPoint := func(point *math32.Vector3, index int) {

		// Get distance from ray to point and if greater than threshold,
		// nothing to do.
		rayPointDistanceSq := ray.DistanceSqToPoint(point)
		if rayPointDistanceSq >= localThresholdSq {
			return
		}
		var intersectPoint math32.Vector3
		ray.ClosestPointToPoint(point, &intersectPoint)
		intersectPoint.ApplyMatrix4(&matrixWorld)
		origin := rc.Ray.Origin()
		distance := origin.DistanceTo(&intersectPoint)
		if distance < rc.Near || distance > rc.Far {
			return
		}
		// Appends intersection of raycaster with this point
		*intersects = append(*intersects, Intersect{
			Distance: distance,
			Point:    intersectPoint,
			Index:    uint32(index),
			Object:   p,
		})
	}

	i := 0
	geom.ReadVertices(func(vertex math32.Vector3) bool {
		testPoint(&vertex, i)
		i++
		return false
	})
}

// RaycastMesh
func (rc *Raycaster) RaycastMesh(m *graphic.Mesh, intersects *[]Intersect) {

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
	checkIntersection := func(mat *material.Material, pA, pB, pC, point *math32.Vector3) *Intersect {

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

		return &Intersect{
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

// RaycastLines
func (rc *Raycaster) RaycastLines(l *graphic.Lines, intersects *[]Intersect) {

	lineRaycast(l, rc, intersects, 2)
}

// RaycastLineStrip
func (rc *Raycaster) RaycastLineStrip(l *graphic.LineStrip, intersects *[]Intersect) {

	lineRaycast(l, rc, intersects, 1)
}

// Internal function used by raycasting for Lines and LineStrip.
func lineRaycast(igr graphic.IGraphic, rc *Raycaster, intersects *[]Intersect, step int) {

	// Get the bounding sphere
	gr := igr.GetGraphic()
	geom := igr.GetGeometry()
	sphere := geom.BoundingSphere()

	// Transform bounding sphere from model to world coordinates and
	// checks intersection with raycaster
	matrixWorld := gr.MatrixWorld()
	sphere.ApplyMatrix4(&matrixWorld)
	if !rc.IsIntersectionSphere(&sphere) {
		return
	}

	// Copy ray and transform to model coordinates
	// This ray will will also be used to check intersects with
	// the geometry, as is much less expensive to transform the
	// ray to model coordinates than the geometry to world coordinates.
	var inverseMatrix math32.Matrix4
	var ray math32.Ray
	inverseMatrix.GetInverse(&matrixWorld)
	ray.Copy(&rc.Ray).ApplyMatrix4(&inverseMatrix)

	var vstart math32.Vector3
	var vend math32.Vector3
	var interSegment math32.Vector3
	var interRay math32.Vector3

	// Get geometry positions and indices buffers
	vboPos := geom.VBO(gls.VertexPosition)
	if vboPos == nil {
		return
	}
	positions := vboPos.Buffer()
	indices := geom.Indices()
	precisionSq := rc.LinePrecision * rc.LinePrecision

	// Checks intersection with individual lines for indexed geometry
	if indices.Size() > 0 {
		for i := 0; i < indices.Size()-1; i += step {
			// Calculates distance from ray to this line segment
			a := indices[i]
			b := indices[i+1]
			positions.GetVector3(int(3*a), &vstart)
			positions.GetVector3(int(3*b), &vend)
			distSq := ray.DistanceSqToSegment(&vstart, &vend, &interRay, &interSegment)
			if distSq > precisionSq {
				continue
			}
			// Move back to world coordinates for distance calculation
			interRay.ApplyMatrix4(&matrixWorld)
			origin := rc.Ray.Origin()
			distance := origin.DistanceTo(&interRay)
			if distance < rc.Near || distance > rc.Far {
				continue
			}

			interSegment.ApplyMatrix4(&matrixWorld)
			*intersects = append(*intersects, Intersect{
				Distance: distance,
				Point:    interSegment,
				Index:    uint32(i),
				Object:   igr,
			})
		}
		// Checks intersection with individual lines for NON indexed geometry
	} else {
		for i := 0; i < positions.Size()/3-1; i += step {
			positions.GetVector3(int(3*i), &vstart)
			positions.GetVector3(int(3*i+3), &vend)
			distSq := ray.DistanceSqToSegment(&vstart, &vend, &interRay, &interSegment)
			if distSq > precisionSq {
				continue
			}

			// Move back to world coordinates for distance calculation
			interRay.ApplyMatrix4(&matrixWorld)
			origin := rc.Ray.Origin()
			distance := origin.DistanceTo(&interRay)
			if distance < rc.Near || distance > rc.Far {
				continue
			}

			interSegment.ApplyMatrix4(&matrixWorld)
			*intersects = append(*intersects, Intersect{
				Distance: distance,
				Point:    interSegment,
				Index:    uint32(i),
				Object:   igr,
			})
		}
	}
}
