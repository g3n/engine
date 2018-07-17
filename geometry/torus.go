// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package geometry

import (
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/math32"
	"math"
)

// Torus represents a torus geometry
type Torus struct {
	Geometry                // embedded geometry
	Radius          float64 // Torus radius
	Tube            float64 // Diameter of the torus tube
	RadialSegments  int     // Number of radial segments
	TubularSegments int     // Number of tubular segments
	Arc             float64 // Central angle
}

// NewTorus returns a pointer to a new torus geometry
func NewTorus(radius, tube float64, radialSegments, tubularSegments int, arc float64) *Torus {

	t := new(Torus)
	t.Geometry.Init()

	t.Radius = radius
	t.Tube = tube
	t.RadialSegments = radialSegments
	t.TubularSegments = tubularSegments
	t.Arc = arc

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
			vertex.X = float32((radius + tube*math.Cos(v)) * math.Cos(u))
			vertex.Y = float32((radius + tube*math.Cos(v)) * math.Sin(u))
			vertex.Z = float32(tube * math.Sin(v))
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
