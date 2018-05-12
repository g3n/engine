// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"github.com/g3n/engine/math32"
	"sort"
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
	Object INode
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
func (rc *Raycaster) IntersectObject(inode INode, recursive bool) []Intersect {

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
func (rc *Raycaster) IntersectObjects(inodes []INode, recursive bool) []Intersect {

	intersects := []Intersect{}
	for _, inode := range inodes {
		rc.intersectObject(inode, &intersects, recursive)
	}
	sort.Slice(intersects, func(i, j int) bool {
		return intersects[i].Distance < intersects[j].Distance
	})
	return intersects
}

func (rc *Raycaster) intersectObject(inode INode, intersects *[]Intersect, recursive bool) {

	node := inode.GetNode()
	if !node.Visible() {
		return
	}
	inode.Raycast(rc, intersects)
	if recursive {
		for _, child := range node.Children() {
			rc.intersectObject(child, intersects, true)
		}
	}
	return
}
