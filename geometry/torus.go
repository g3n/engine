// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package geometry

import (
	"math"

	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/math32"
)

// NewTorus creates a torus geometry with the specified revolution radius, tube radius,
// number of radial segments, number of tubular segments, and arc length angle in radians.
// TODO instead of 'arc' have thetaStart and thetaLength for consistency with other generators
// TODO then rename this to NewTorusSector and add a NewTorus constructor
func NewTorus(radius, tubeRadius float64, radialSegments, tubularSegments int, arc float64) *Geometry {

	t := NewGeometry()

	// Create buffers
	positions := math32.NewArrayF32(0, 0)
	normals := math32.NewArrayF32(0, 0)
	uvs := math32.NewArrayF32(0, 0)
	indices := math32.NewArrayU32(0, 0)

	var center math32.Vector3
	for j := 0; j <= radialSegments; j++ {
		for i := 0; i <= tubularSegments; i++ {
			u := float64(i) / float64(tubularSegments) * arc
			v := float64(j) / float64(radialSegments) * math.Pi * 2

			center.X = float32(radius * math.Cos(u))
			center.Y = float32(radius * math.Sin(u))

			var vertex math32.Vector3
			vertex.X = float32((radius + tubeRadius*math.Cos(v)) * math.Cos(u))
			vertex.Y = float32((radius + tubeRadius*math.Cos(v)) * math.Sin(u))
			vertex.Z = float32(tubeRadius * math.Sin(v))
			positions.AppendVector3(&vertex)

			uvs.Append(float32(float64(i)/float64(tubularSegments)), float32(float64(j)/float64(radialSegments)))
			normals.AppendVector3(vertex.Sub(&center).Normalize())
		}
	}

	for j := 1; j <= radialSegments; j++ {
		for i := 1; i <= tubularSegments; i++ {
			a := (tubularSegments+1)*j + i - 1
			b := (tubularSegments+1)*(j-1) + i - 1
			c := (tubularSegments+1)*(j-1) + i
			d := (tubularSegments+1)*j + i
			indices.Append(uint32(a), uint32(b), uint32(d), uint32(b), uint32(c), uint32(d))
		}
	}

	t.SetIndices(indices)
	t.AddVBO(gls.NewVBO(positions).AddAttrib(gls.VertexPosition))
	t.AddVBO(gls.NewVBO(normals).AddAttrib(gls.VertexNormal))
	t.AddVBO(gls.NewVBO(uvs).AddAttrib(gls.VertexTexcoord))

	return t
}
