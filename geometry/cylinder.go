// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package geometry

import (
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/math32"
	"math"
)

type Cylinder struct {
	Geometry
	RadiusTop      float64
	RadiusBottom   float64
	Height         float64
	RadialSegments int
	HeightSegments int
	ThetaStart     float64
	ThetaLength    float64
	Top            bool
	Bottom         bool
}

// NewCylinder creates and returns a pointer to a new Cylinder geometry object.
func NewCylinder(radiusTop, radiusBottom, height float64,
	radialSegments, heightSegments int,
	thetaStart, thetaLength float64, top, bottom bool) *Cylinder {

	c := new(Cylinder)
	c.Geometry.Init()

	c.RadiusTop = radiusTop
	c.RadiusBottom = radiusBottom
	c.Height = height
	c.RadialSegments = radialSegments
	c.HeightSegments = heightSegments
	c.ThetaStart = thetaStart
	c.ThetaLength = thetaLength
	c.Top = top
	c.Bottom = bottom

	heightHalf := height / 2
	vertices := [][]int{}
	uvs := [][]math32.Vector2{}

	// Create buffer for vertex positions
	Positions := math32.NewArrayF32(0, 0)

	for y := 0; y <= heightSegments; y++ {
		var verticesRow = []int{}
		var uvsRow = []math32.Vector2{}
		v := float64(y) / float64(heightSegments)
		radius := v*(radiusBottom-radiusTop) + radiusTop
		for x := 0; x <= radialSegments; x++ {
			u := float64(x) / float64(radialSegments)
			var vertex math32.Vector3
			vertex.X = float32(radius * math.Sin(u*thetaLength+thetaStart))
			vertex.Y = float32(-v*height + heightHalf)
			vertex.Z = float32(radius * math.Cos(u*thetaLength+thetaStart))
			Positions.AppendVector3(&vertex)
			verticesRow = append(verticesRow, Positions.Size()/3-1)
			uvsRow = append(uvsRow, math32.Vector2{float32(u), 1.0 - float32(v)})
		}
		vertices = append(vertices, verticesRow)
		uvs = append(uvs, uvsRow)
	}

	tanTheta := (radiusBottom - radiusTop) / height
	var na, nb math32.Vector3

	// Create preallocated buffers for normals and uvs and buffer for indices
	npos := Positions.Size()
	Normals := math32.NewArrayF32(npos, npos)
	Uvs := math32.NewArrayF32(2*npos/3, 2*npos/3)
	Indices := math32.NewArrayU32(0, 0)

	for x := 0; x < radialSegments; x++ {
		if radiusTop != 0 {
			Positions.GetVector3(3*vertices[0][x], &na)
			Positions.GetVector3(3*vertices[0][x+1], &nb)
		} else {
			Positions.GetVector3(3*vertices[1][x], &na)
			Positions.GetVector3(3*vertices[1][x+1], &nb)
		}

		na.SetY(float32(math.Sqrt(float64(na.X*na.X+na.Z*na.Z)) * tanTheta)).Normalize()
		nb.SetY(float32(math.Sqrt(float64(nb.X*nb.X+nb.Z*nb.Z)) * tanTheta)).Normalize()

		for y := 0; y < heightSegments; y++ {
			v1 := vertices[y][x]
			v2 := vertices[y+1][x]
			v3 := vertices[y+1][x+1]
			v4 := vertices[y][x+1]

			n1 := na
			n2 := na
			n3 := nb
			n4 := nb

			uv1 := uvs[y][x]
			uv2 := uvs[y+1][x]
			uv3 := uvs[y+1][x+1]
			uv4 := uvs[y][x+1]

			Indices.Append(uint32(v1), uint32(v2), uint32(v4))
			Normals.SetVector3(3*v1, &n1)
			Normals.SetVector3(3*v2, &n2)
			Normals.SetVector3(3*v4, &n4)

			Indices.Append(uint32(v2), uint32(v3), uint32(v4))
			Normals.SetVector3(3*v2, &n2)
			Normals.SetVector3(3*v3, &n3)
			Normals.SetVector3(3*v4, &n4)

			Uvs.SetVector2(2*v1, &uv1)
			Uvs.SetVector2(2*v2, &uv2)
			Uvs.SetVector2(2*v3, &uv3)
			Uvs.SetVector2(2*v4, &uv4)
		}
	}
	// First group is the body of the cylinder
	// without the caps
	c.AddGroup(0, Indices.Size(), 0)
	nextGroup := Indices.Size()

	// Top cap
	if top && radiusTop > 0 {

		// Array of vertex indices to build used to build the faces.
		indices := []uint32{}
		nextidx := Positions.Size() / 3

		// Appends top segments vertices and builds array of its indices
		var uv1, uv2, uv3 math32.Vector2
		for x := 0; x < radialSegments; x++ {
			uv1 = uvs[0][x]
			uv2 = uvs[0][x+1]
			uv3 = math32.Vector2{uv2.X, 0}
			// Appends CENTER with its own UV.
			Positions.Append(0, float32(heightHalf), 0)
			Normals.Append(0, 1, 0)
			Uvs.AppendVector2(&uv3)
			indices = append(indices, uint32(nextidx))
			nextidx++
			// Appends vertex
			v := math32.Vector3{}
			vi := vertices[0][x]
			Positions.GetVector3(3*vi, &v)
			Positions.AppendVector3(&v)
			Normals.Append(0, 1, 0)
			Uvs.AppendVector2(&uv1)
			indices = append(indices, uint32(nextidx))
			nextidx++
		}
		// Appends copy of first vertex (center)
		var vertex, normal math32.Vector3
		var uv math32.Vector2
		Positions.GetVector3(3*int(indices[0]), &vertex)
		Normals.GetVector3(3*int(indices[0]), &normal)
		Uvs.GetVector2(2*int(indices[0]), &uv)
		Positions.AppendVector3(&vertex)
		Normals.AppendVector3(&normal)
		Uvs.AppendVector2(&uv)
		indices = append(indices, uint32(nextidx))
		nextidx++

		// Appends copy of second vertex (v1) USING LAST UV2
		Positions.GetVector3(3*int(indices[1]), &vertex)
		Normals.GetVector3(3*int(indices[1]), &normal)
		Positions.AppendVector3(&vertex)
		Normals.AppendVector3(&normal)
		Uvs.AppendVector2(&uv2)
		indices = append(indices, uint32(nextidx))
		nextidx++

		// Append faces indices
		for x := 0; x < radialSegments; x++ {
			pos := 2 * x
			i1 := indices[pos]
			i2 := indices[pos+1]
			i3 := indices[pos+3]
			Indices.Append(uint32(i1), uint32(i2), uint32(i3))
		}
		// Second group is optional top cap of the cylinder
		c.AddGroup(nextGroup, Indices.Size()-nextGroup, 1)
		nextGroup = Indices.Size()
	}

	// Bottom cap
	if bottom && radiusBottom > 0 {

		// Array of vertex indices to build used to build the faces.
		indices := []uint32{}
		nextidx := Positions.Size() / 3

		// Appends top segments vertices and builds array of its indices
		var uv1, uv2, uv3 math32.Vector2
		for x := 0; x < radialSegments; x++ {
			uv1 = uvs[heightSegments][x]
			uv2 = uvs[heightSegments][x+1]
			uv3 = math32.Vector2{uv2.X, 1}
			// Appends CENTER with its own UV.
			Positions.Append(0, float32(-heightHalf), 0)
			Normals.Append(0, -1, 0)
			Uvs.AppendVector2(&uv3)
			indices = append(indices, uint32(nextidx))
			nextidx++
			// Appends vertex
			v := math32.Vector3{}
			vi := vertices[heightSegments][x]
			Positions.GetVector3(3*vi, &v)
			Positions.AppendVector3(&v)
			Normals.Append(0, -1, 0)
			Uvs.AppendVector2(&uv1)
			indices = append(indices, uint32(nextidx))
			nextidx++
		}

		// Appends copy of first vertex (center)
		var vertex, normal math32.Vector3
		var uv math32.Vector2
		Positions.GetVector3(3*int(indices[0]), &vertex)
		Normals.GetVector3(3*int(indices[0]), &normal)
		Uvs.GetVector2(2*int(indices[0]), &uv)
		Positions.AppendVector3(&vertex)
		Normals.AppendVector3(&normal)
		Uvs.AppendVector2(&uv)
		indices = append(indices, uint32(nextidx))
		nextidx++

		// Appends copy of second vertex (v1) USING LAST UV2
		Positions.GetVector3(3*int(indices[1]), &vertex)
		Normals.GetVector3(3*int(indices[1]), &normal)
		Positions.AppendVector3(&vertex)
		Normals.AppendVector3(&normal)
		Uvs.AppendVector2(&uv2)
		indices = append(indices, uint32(nextidx))
		nextidx++

		// Appends faces indices
		for x := 0; x < radialSegments; x++ {
			pos := 2 * x
			i1 := indices[pos]
			i2 := indices[pos+3]
			i3 := indices[pos+1]
			Indices.Append(uint32(i1), uint32(i2), uint32(i3))
		}
		// Third group is optional bottom cap of the cylinder
		c.AddGroup(nextGroup, Indices.Size()-nextGroup, 2)
	}

	c.SetIndices(Indices)
	c.AddVBO(gls.NewVBO().AddAttrib("VertexPosition", 3).SetBuffer(Positions))
	c.AddVBO(gls.NewVBO().AddAttrib("VertexNormal", 3).SetBuffer(Normals))
	c.AddVBO(gls.NewVBO().AddAttrib("VertexTexcoord", 2).SetBuffer(Uvs))

	return c
}
