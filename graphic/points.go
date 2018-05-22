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

// Points represents a geometry containing only points
type Points struct {
	Graphic             // Embedded graphic
	uniMVPm gls.Uniform // Model view projection matrix uniform location cache
}

// NewPoints creates and returns a graphic points object with the specified
// geometry and material.
func NewPoints(igeom geometry.IGeometry, imat material.IMaterial) *Points {

	p := new(Points)
	p.Graphic.Init(igeom, gls.POINTS)
	if imat != nil {
		p.AddMaterial(p, imat, 0, 0)
	}
	p.uniMVPm.Init("MVP")
	return p
}

// RenderSetup is called by the engine before rendering this graphic.
func (p *Points) RenderSetup(gs *gls.GLS, rinfo *core.RenderInfo) {

	// Transfer model view projection matrix uniform
	mvpm := p.ModelViewProjectionMatrix()
	location := p.uniMVPm.Location(gs)
	gs.UniformMatrix4fv(location, 1, false, &mvpm[0])
}

// Raycast satisfies the INode interface and checks the intersections
// of this geometry with the specified raycaster.
func (p *Points) Raycast(rc *core.Raycaster, intersects *[]core.Intersect) {

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
		*intersects = append(*intersects, core.Intersect{
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
