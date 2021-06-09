// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package geometry

import (
	"math"

	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/math32"
)

// NewDisk creates a disk (filled circle) geometry with the specified
// radius and number of radial segments/triangles (minimum 3).
func NewDisk(radius float64, segments int) *Geometry {
	return NewDiskSector(radius, segments, 0, 2*math.Pi)
}

// NewDiskSector creates a disk (filled circle) or disk sector geometry with the specified radius,
// number of radial segments/triangles (minimum 3), sector start angle in radians, and sector size angle in radians.
// The center of the disk is at the origin, and theta runs counter-clockwise on the XY plane, starting at (x,y,z)=(1,0,0).
func NewDiskSector(radius float64, segments int, thetaStart, thetaLength float64) *Geometry {

	d := NewGeometry()

	// Validate arguments
	if segments < 3 {
		panic("Invalid argument: segments. The number of segments needs to be greater or equal to 3.")
	}

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

	d.SetIndices(indices)
	d.AddVBO(gls.NewVBO(positions).AddAttrib(gls.VertexPosition))
	d.AddVBO(gls.NewVBO(normals).AddAttrib(gls.VertexNormal))
	d.AddVBO(gls.NewVBO(uvs).AddAttrib(gls.VertexTexcoord))

	// Update volume
	d.volume = 0
	d.volumeValid = true

	return d
}
