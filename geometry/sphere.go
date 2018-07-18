// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package geometry

import (
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/math32"
	"math"
)

// Sphere represents a sphere geometry
type Sphere struct {
	Geometry
	Radius         float64
	WidthSegments  int
	HeightSegments int
	PhiStart       float64
	PhiLength      float64
	ThetaStart     float64
	ThetaLength    float64
}

// NewSphere returns a pointer to a new Sphere geometry object
func NewSphere(radius float64, widthSegments, heightSegments int, phiStart, phiLength, thetaStart, thetaLength float64) *Sphere {

	s := new(Sphere)
	s.Geometry.Init()

	s.Radius = radius
	s.WidthSegments = widthSegments
	s.HeightSegments = heightSegments
	s.PhiStart = phiStart
	s.PhiLength = phiLength
	s.ThetaStart = thetaStart

	thetaEnd := thetaStart + thetaLength
	vertexCount := (widthSegments + 1) * (heightSegments + 1)

	// Create buffers
	positions := math32.NewArrayF32(vertexCount*3, vertexCount*3)
	normals := math32.NewArrayF32(vertexCount*3, vertexCount*3)
	uvs := math32.NewArrayF32(vertexCount*2, vertexCount*2)
	indices := math32.NewArrayU32(0, vertexCount)

	index := 0
	vertices := make([][]uint32, 0)
	var normal math32.Vector3

	for y := 0; y <= heightSegments; y++ {
		verticesRow := make([]uint32, 0)
		v := float64(y) / float64(heightSegments)
		for x := 0; x <= widthSegments; x++ {
			u := float64(x) / float64(widthSegments)
			px := -radius * math.Cos(phiStart+u*phiLength) * math.Sin(thetaStart+v*thetaLength)
			py := radius * math.Cos(thetaStart+v*thetaLength)
			pz := radius * math.Sin(phiStart+u*phiLength) * math.Sin(thetaStart+v*thetaLength)
			normal.Set(float32(px), float32(py), float32(pz)).Normalize()

			positions.Set(index*3, float32(px), float32(py), float32(pz))
			normals.SetVector3(index*3, &normal)
			uvs.Set(index*2, float32(u), float32(v))
			verticesRow = append(verticesRow, uint32(index))
			index++
		}
		vertices = append(vertices, verticesRow)
	}

	for y := 0; y < heightSegments; y++ {
		for x := 0; x < widthSegments; x++ {
			v1 := vertices[y][x+1]
			v2 := vertices[y][x]
			v3 := vertices[y+1][x]
			v4 := vertices[y+1][x+1]
			if y != 0 || thetaStart > 0 {
				indices.Append(v1, v2, v4)
			}
			if y != heightSegments-1 || thetaEnd < math.Pi {
				indices.Append(v2, v3, v4)
			}
		}
	}

	s.SetIndices(indices)
	s.AddVBO(gls.NewVBO(positions).AddAttrib(gls.VertexPosition))
	s.AddVBO(gls.NewVBO(normals).AddAttrib(gls.VertexNormal))
	s.AddVBO(gls.NewVBO(uvs).AddAttrib(gls.VertexTexcoord))

	r := float32(radius)

	// Update bounding sphere
	s.boundingSphere.Radius = 3
	s.boundingSphereValid = true

	// Update bounding box
	s.boundingBox = math32.Box3{math32.Vector3{-r, -r, -r}, math32.Vector3{r, r, r}}
	s.boundingBoxValid = true

	return s
}
