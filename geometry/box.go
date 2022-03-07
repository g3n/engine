// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package geometry

import (
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/math32"
)

// NewCube creates a cube geometry with the specified size.
func NewCube(size float32) *Geometry {
	return NewSegmentedCube(size, 1)
}

// NewSegmentedCube creates a segmented cube geometry with the specified size and number of segments.
func NewSegmentedCube(size float32, segments int) *Geometry {
	return NewSegmentedBox(size, size, size, segments, segments, segments)
}

// NewBox creates a box geometry with the specified width, height, and length.
func NewBox(width, height, length float32) *Geometry {
	return NewSegmentedBox(width, height, length, 1, 1, 1)
}

// NewSegmentedBox creates a segmented box geometry with the specified width, height, length, and number of segments in each dimension.
func NewSegmentedBox(width, height, length float32, widthSegments, heightSegments, lengthSegments int) *Geometry {

	box := NewGeometry()

	// Validate arguments
	if widthSegments <= 0 || heightSegments <= 0 || lengthSegments <= 0 {
		panic("Invalid argument(s). All segment quantities should be greater than zero.")
	}

	// Create buffers
	positions := math32.NewArrayF32(0, 16)
	normals := math32.NewArrayF32(0, 16)
	uvs := math32.NewArrayF32(0, 16)
	indices := math32.NewArrayU32(0, 16)

	// Internal function to build each of the six box planes
	buildPlane := func(u, v string, udir, vdir int, width, height, length float32, materialIndex uint) {

		offset := positions.Len() / 3
		gridX := widthSegments
		gridY := heightSegments
		var w string

		if (u == "x" && v == "y") || (u == "y" && v == "x") {
			w = "z"
		} else if (u == "x" && v == "z") || (u == "z" && v == "x") {
			w = "y"
			gridY = lengthSegments
		} else if (u == "z" && v == "y") || (u == "y" && v == "z") {
			w = "x"
			gridX = lengthSegments
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

	wHalf := width / 2
	hHalf := height / 2
	lHalf := length / 2

	buildPlane("z", "y", -1, -1, length, height, wHalf, 0) // px
	buildPlane("z", "y", 1, -1, length, height, -wHalf, 1) // nx
	buildPlane("x", "z", 1, 1, width, length, hHalf, 2)    // py
	buildPlane("x", "z", 1, -1, width, length, -hHalf, 3)  // ny
	buildPlane("x", "y", 1, -1, width, height, lHalf, 4)   // pz
	buildPlane("x", "y", -1, -1, width, height, -lHalf, 5) // nz

	box.SetIndices(indices)
	box.AddVBO(gls.NewVBO(positions).AddAttrib(gls.VertexPosition))
	box.AddVBO(gls.NewVBO(normals).AddAttrib(gls.VertexNormal))
	box.AddVBO(gls.NewVBO(uvs).AddAttrib(gls.VertexTexcoord))

	// Update bounding box
	box.boundingBox.Min = math32.Vector3{-wHalf, -hHalf, -lHalf}
	box.boundingBox.Max = math32.Vector3{wHalf, hHalf, lHalf}
	box.boundingBoxValid = true

	// Update bounding sphere
	box.boundingSphere.Radius = math32.Sqrt(math32.Pow(width/2, 2) + math32.Pow(height/2, 2) + math32.Pow(length/2, 2))
	box.boundingSphereValid = true

	// Update area
	box.area = 2*width + 2*height + 2*length
	box.areaValid = true

	// Update volume
	box.volume = width * height * length
	box.volumeValid = true

	return box
}
