// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package geometry

import (
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/math32"
)

// Box represents the geometry of a rectangular cuboid.
// See https://en.wikipedia.org/wiki/Cuboid#Rectangular_cuboid for more details.
// A Box geometry is defined by its width, height, and length and also by the number
// of segments in each dimension.
type Box struct {
	Geometry
	Width          float32
	Height         float32
	Length         float32
	WidthSegments  int // > 0
	HeightSegments int // > 0
	LengthSegments int // > 0
}

// NewCube creates a new cube geometry of the specified size.
func NewCube(size float32) *Box {
	return NewSegmentedBox(size, size, size, 1, 1, 1)
}

// NewSegmentedCube creates a cube geometry of the specified size and number of segments.
func NewSegmentedCube(size float32, segments int) *Box {
	return NewSegmentedBox(size, size, size, segments, segments, segments)
}

// NewBox creates a box geometry of the specified width, height, and length.
func NewBox(width, height, length float32) *Box {
	return NewSegmentedBox(width, height, length, 1, 1, 1)
}

// NewSegmentedBox creates a box geometry of the specified size and with the specified number
// of segments in each dimension. This is the Box constructor with most tunable parameters.
func NewSegmentedBox(width, height, length float32, widthSegments, heightSegments, lengthSegments int) *Box {

	box := new(Box)
	box.Geometry.Init()

	// Validate arguments
	if widthSegments <= 0 || heightSegments <= 0 || lengthSegments <= 0 {
		panic("Invalid argument(s). All segment quantities should be greater than zero.")
	}

	box.Width = width
	box.Height = height
	box.Length = length
	box.WidthSegments = widthSegments
	box.HeightSegments = heightSegments
	box.LengthSegments = lengthSegments

	// Create buffers
	positions := math32.NewArrayF32(0, 16)
	normals := math32.NewArrayF32(0, 16)
	uvs := math32.NewArrayF32(0, 16)
	indices := math32.NewArrayU32(0, 16)

	// Internal function to build each of the six box planes
	buildPlane := func(u, v string, udir, vdir int, width, height, length float32, materialIndex uint) {

		offset := positions.Len() / 3
		gridX := box.WidthSegments
		gridY := box.HeightSegments
		var w string

		if (u == "x" && v == "y") || (u == "y" && v == "x") {
			w = "z"
		} else if (u == "x" && v == "z") || (u == "z" && v == "x") {
			w = "y"
			gridY = box.LengthSegments
		} else if (u == "z" && v == "y") || (u == "y" && v == "z") {
			w = "x"
			gridX = box.LengthSegments
		}

		var normal math32.Vector3
		if length > 0 {
			normal.SetByName(w, 1)
		} else {
			normal.SetByName(w, -1)
		}

		wHalf := width / 2
		hHalf := height / 2
		gridX1 := gridX + 1
		gridY1 := gridY + 1
		segmentWidth := width / float32(gridX)
		segmentHeight := height / float32(gridY)

		// Generate the plane vertices, normals, and uv coordinates
		for iy := 0; iy < gridY1; iy++ {
			for ix := 0; ix < gridX1; ix++ {
				var vector math32.Vector3
				vector.SetByName(u, (float32(ix)*segmentWidth-wHalf)*float32(udir))
				vector.SetByName(v, (float32(iy)*segmentHeight-hHalf)*float32(vdir))
				vector.SetByName(w, length)
				positions.AppendVector3(&vector)
				normals.AppendVector3(&normal)
				uvs.Append(float32(ix)/float32(gridX), float32(1)-(float32(iy)/float32(gridY)))
			}
		}

		// Generate the indices for the vertices, normals and uv coordinates
		gstart := indices.Size()
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
		box.AddGroup(gstart, gcount, int(materialIndex))
	}

	wHalf := box.Width / 2
	hHalf := box.Height / 2
	lHalf := box.Length / 2

	buildPlane("z", "y", -1, -1, box.Length, box.Height, wHalf, 0) // px
	buildPlane("z", "y", 1, -1, box.Length, box.Height, -wHalf, 1) // nx
	buildPlane("x", "z", 1, 1, box.Width, box.Length, hHalf, 2)    // py
	buildPlane("x", "z", 1, -1, box.Width, box.Length, -hHalf, 3)  // ny
	buildPlane("x", "y", 1, -1, box.Width, box.Height, lHalf, 4)   // pz
	buildPlane("x", "y", -1, -1, box.Width, box.Height, -lHalf, 5) // nz

	box.SetIndices(indices)
	box.AddVBO(gls.NewVBO(positions).AddAttrib(gls.VertexPosition))
	box.AddVBO(gls.NewVBO(normals).AddAttrib(gls.VertexNormal))
	box.AddVBO(gls.NewVBO(uvs).AddAttrib(gls.VertexTexcoord))

	// Update bounding box
	box.boundingBox.Min = math32.Vector3{-wHalf, -hHalf, -lHalf}
	box.boundingBox.Max = math32.Vector3{wHalf, hHalf, lHalf}
	box.boundingBoxValid = true

	// Update bounding sphere
	box.boundingSphere.Radius = math32.Sqrt(math32.Pow(width/2,2) + math32.Pow(height/2,2) + math32.Pow(length/2,2))
	box.boundingSphereValid = true

	// Update area
	box.area = 2*width + 2*height + 2*length
	box.areaValid = true

	// Update volume
	box.volume = width * height * length
	box.volumeValid = true

	return box
}
