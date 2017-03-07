// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package geometry

import (
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/math32"
)

type Plane struct {
	Geometry
	Width          float32
	Height         float32
	WidthSegments  int
	HeightSegments int
}

// NewPlane creates and returns a pointer to a Plane Geometry.
// The plane is defined by its width, height and the number of width and height segments.
// The minimum number of segments for the width and/or the height is 1.
// The plane is generated centered in the XY plane with Z=0.
func NewPlane(width, height float32, widthSegments, heightSegments int) *Plane {

	plane := new(Plane)
	plane.Geometry.Init()

	plane.Width = width
	plane.Height = height
	plane.WidthSegments = widthSegments
	plane.HeightSegments = heightSegments

	width_half := width / 2
	height_half := height / 2
	gridX := widthSegments
	gridY := heightSegments
	gridX1 := gridX + 1
	gridY1 := gridY + 1
	segment_width := width / float32(gridX)
	segment_height := height / float32(gridY)

	// Create buffers
	positions := math32.NewArrayF32(0, 16)
	normals := math32.NewArrayF32(0, 16)
	uvs := math32.NewArrayF32(0, 16)
	indices := math32.NewArrayU32(0, 16)

	// Generate plane vertices, vertices normals and vertices texture mappings.
	for iy := 0; iy < gridY1; iy++ {
		y := float32(iy)*segment_height - height_half
		for ix := 0; ix < gridX1; ix++ {
			x := float32(ix)*segment_width - width_half
			positions.Append(float32(x), float32(-y), 0)
			normals.Append(0, 0, 1)
			uvs.Append(float32(float64(ix)/float64(gridX)), float32(float64(1)-(float64(iy)/float64(gridY))))
		}
	}

	// Generate plane vertices indices for the faces
	for iy := 0; iy < gridY; iy++ {
		for ix := 0; ix < gridX; ix++ {
			a := ix + gridX1*iy
			b := ix + gridX1*(iy+1)
			c := (ix + 1) + gridX1*(iy+1)
			d := (ix + 1) + gridX1*iy
			indices.Append(uint32(a), uint32(b), uint32(d))
			indices.Append(uint32(b), uint32(c), uint32(d))
		}
	}

	plane.SetIndices(indices)
	plane.AddVBO(gls.NewVBO().AddAttrib("VertexPosition", 3).SetBuffer(positions))
	plane.AddVBO(gls.NewVBO().AddAttrib("VertexNormal", 3).SetBuffer(normals))
	plane.AddVBO(gls.NewVBO().AddAttrib("VertexTexcoord", 2).SetBuffer(uvs))

	return plane
}
