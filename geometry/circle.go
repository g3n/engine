// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package geometry

import (
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/math32"
	"math"
)

type Circle struct {
	Geometry
	Radius      float64
	Segments    int
	ThetaStart  float64
	ThetaLength float64
}

// NewCircle creates and returns a pointer to a new Circle geometry object.
// The geometry is defined by its radius, the number of segments (triangles), minimum = 3,
// the start angle in radians for the first segment (thetaStart) and
// the central angle in radians (thetaLength) of the circular sector.
func NewCircle(radius float64, segments int, thetaStart, thetaLength float64) *Circle {

	circ := new(Circle)
	circ.Geometry.Init()

	circ.Radius = radius
	circ.Segments = segments
	circ.ThetaStart = thetaStart
	circ.ThetaLength = thetaLength

	if segments < 3 {
		segments = 3
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

	// Append circle center uv coord
	centerUV := math32.NewVector2(0.5, 0.5)
	uvs.AppendVector2(centerUV)

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
	circ.AddVBO(gls.NewVBO().AddAttrib("VertexPosition", 3).SetBuffer(positions))
	circ.AddVBO(gls.NewVBO().AddAttrib("VertexNormal", 3).SetBuffer(normals))
	circ.AddVBO(gls.NewVBO().AddAttrib("VertexTexcoord", 2).SetBuffer(uvs))

	//circ.BoundingSphere = math32.NewSphere(math32.NewVector3(0,0,0), float32(radius))

	return circ
}
