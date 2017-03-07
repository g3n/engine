// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package geometry

import (
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/math32"
)

type Box struct {
	Geometry
	Width          float64
	Height         float64
	Depth          float64
	WidthSegments  int
	HeightSegments int
	DepthSegments  int
}

// NewBox creates and returns a pointer to a new Box geometry object.
// The geometry is defined by its width, height, depth and the number of
// segments of each dimension (minimum = 1).
func NewBox(width, height, depth float64, widthSegments, heightSegments, depthSegments int) *Box {

	box := new(Box)
	box.Geometry.Init()

	box.Width = width
	box.Height = height
	box.Depth = depth
	box.WidthSegments = widthSegments
	box.HeightSegments = heightSegments
	box.DepthSegments = depthSegments

	// Create buffers
	positions := math32.NewArrayF32(0, 16)
	normals := math32.NewArrayF32(0, 16)
	uvs := math32.NewArrayF32(0, 16)
	indices := math32.NewArrayU32(0, 16)

	width_half := width / 2
	height_half := height / 2
	depth_half := depth / 2

	// Internal function to build each box plane
	buildPlane := func(u, v string, udir, vdir int, width, height, depth float64, materialIndex uint) {

		gridX := widthSegments
		gridY := heightSegments
		width_half := width / 2
		height_half := height / 2
		offset := positions.Len() / 3
		var w string

		if (u == "x" && v == "y") || (u == "y" && v == "x") {
			w = "z"
		} else if (u == "x" && v == "z") || (u == "z" && v == "x") {
			w = "y"
			gridY = depthSegments
		} else if (u == "z" && v == "y") || (u == "y" && v == "z") {
			w = "x"
			gridX = depthSegments
		}

		gridX1 := gridX + 1
		gridY1 := gridY + 1
		segment_width := width / float64(gridX)
		segment_height := height / float64(gridY)
		var normal math32.Vector3
		if depth > 0 {
			normal.SetByName(w, 1)
		} else {
			normal.SetByName(w, -1)
		}

		// Generates the plane vertices, normals and uv coordinates.
		for iy := 0; iy < gridY1; iy++ {
			for ix := 0; ix < gridX1; ix++ {
				var vector math32.Vector3
				vector.SetByName(u, float32((float64(ix)*segment_width-width_half)*float64(udir)))
				vector.SetByName(v, float32((float64(iy)*segment_height-height_half)*float64(vdir)))
				vector.SetByName(w, float32(depth))
				positions.AppendVector3(&vector)
				normals.AppendVector3(&normal)
				uvs.Append(float32(float64(ix)/float64(gridX)), float32(float64(1)-(float64(iy)/float64(gridY))))
			}
		}

		gstart := indices.Size()
		matIndex := materialIndex
		// Generates the indices for the vertices, normals and uvs
		for iy := 0; iy < gridY; iy++ {
			for ix := 0; ix < gridX; ix++ {
				a := ix + gridX1*iy
				b := ix + gridX1*(iy+1)
				c := (ix + 1) + gridX1*(iy+1)
				d := (ix + 1) + gridX1*iy
				indices.Append(uint32(a+offset), uint32(b+offset), uint32(d+offset), uint32(b+offset), uint32(c+offset), uint32(d+offset))
			}
		}
		gcount := indices.Size() - gstart
		box.AddGroup(gstart, gcount, int(matIndex))
	}

	buildPlane("z", "y", -1, -1, depth, height, width_half, 0)  // px
	buildPlane("z", "y", 1, -1, depth, height, -width_half, 1)  // nx
	buildPlane("x", "z", 1, 1, width, depth, height_half, 2)    // py
	buildPlane("x", "z", 1, -1, width, depth, -height_half, 3)  // ny
	buildPlane("x", "y", 1, -1, width, height, depth_half, 4)   // pz
	buildPlane("x", "y", -1, -1, width, height, -depth_half, 5) // nz

	box.SetIndices(indices)
	box.AddVBO(gls.NewVBO().AddAttrib("VertexPosition", 3).SetBuffer(positions))
	box.AddVBO(gls.NewVBO().AddAttrib("VertexNormal", 3).SetBuffer(normals))
	box.AddVBO(gls.NewVBO().AddAttrib("VertexTexcoord", 2).SetBuffer(uvs))

	return box
}
