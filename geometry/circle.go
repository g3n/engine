// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package geometry

import (
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/math32"
	"math"
)

// Circle represents the geometry of a filled circle (i.e. a disk)
// The center of the circle is at the origin, and theta runs counter-clockwise
// on the XY plane, starting at (x,y,z)=(1,0,0).
type Circle struct {
	Geometry
	Radius      float64
	Segments    int // >= 3
	ThetaStart  float64
	ThetaLength float64
}

// NewCircle creates a new circle geometry with the specified radius
// and number of radial segments/triangles (minimum 3).
func NewCircle(radius float64, segments int) *Circle {
	return NewCircleSector(radius, segments, 0, 2*math.Pi)
}

// NewCircleSector creates a new circle or circular sector geometry with the specified radius,
// number of radial segments/triangles (minimum 3), sector start angle in radians (thetaStart),
// and sector size angle in radians (thetaLength). This is the Circle constructor with most tunable parameters.
func NewCircleSector(radius float64, segments int, thetaStart, thetaLength float64) *Circle {

	circ := new(Circle)
	circ.Geometry.Init()

	// Validate arguments
	if segments < 3 {
		panic("Invalid argument: segments. The number of segments needs to be greater or equal to 3.")
	}

	circ.Radius = radius
	circ.Segments = segments
	circ.ThetaStart = thetaStart
	circ.ThetaLength = thetaLength

	// Create buffers
	positions := math32.NewArrayF32(0, 16)
	normals := math32.NewArrayF32(0, 16)
	uvs := math32.NewArrayF32(0, 16)
	indices := math32.NewArrayU32(0, 16)

	// Append circle center position
	center := math32.NewVector3(0, 0, 0)
	positions.AppendVector3(center)

	// Append circle center normal
	var normal math32.Vector3
	normal.Z = 1
	normals.AppendVector3(&normal)

	// Append circle center uv coordinate
	centerUV := math32.NewVector2(0.5, 0.5)
	uvs.AppendVector2(centerUV)

	// Generate the segments
	for i := 0; i <= segments; i++ {
		segment := thetaStart + float64(i)/float64(segments)*thetaLength

		vx := float32(radius * math.Cos(segment))
		vy := float32(radius * math.Sin(segment))

		// Appends vertex position, normal and uv coordinates
		positions.Append(vx, vy, 0)
		normals.AppendVector3(&normal)
		uvs.Append((vx/float32(radius)+1)/2, (vy/float32(radius)+1)/2)
	}

	for i := 1; i <= segments; i++ {
		indices.Append(uint32(i), uint32(i)+1, 0)
	}

	circ.SetIndices(indices)
	circ.AddVBO(gls.NewVBO(positions).AddAttrib(gls.VertexPosition))
	circ.AddVBO(gls.NewVBO(normals).AddAttrib(gls.VertexNormal))
	circ.AddVBO(gls.NewVBO(uvs).AddAttrib(gls.VertexTexcoord))

	// Update volume
	circ.volume = 0
	circ.volumeValid = true

	return circ
}
