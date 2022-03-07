// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// TODO: UVs, Caps

package geometry

import (
	"math"

	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/math32"
)

func CalculateNormals(indices math32.ArrayU32, positions, normals math32.ArrayF32) math32.ArrayF32 {
	var x1x2, y1y2, z1z2, x3x2, y3y2, z3z2, x, y, z, l float32
	var x1, y1, z1, x2, y2, z2, x3, y3, z3 int // position indexes

	for i := 0; i < len(indices)/3; i++ {
		x1 = int(indices[i*3] * uint32(3))
		y1 = x1 + 1
		z1 = x1 + 2
		x2 = int(indices[i*3+1] * uint32(3))
		y2 = x2 + 1
		z2 = x2 + 2
		x3 = int(indices[i*3+2] * uint32(3))
		y3 = x3 + 1
		z3 = x3 + 2

		x1x2 = positions[x1] - positions[x2]
		y1y2 = positions[y1] - positions[y2]
		z1z2 = positions[z1] - positions[z2]
		x3x2 = positions[x3] - positions[x2]
		y3y2 = positions[y3] - positions[y2]
		z3z2 = positions[z3] - positions[z2]

		x = y1y2*z3z2 - z1z2*y3y2
		y = z1z2*x3x2 - x1x2*z3z2
		z = x1x2*y3y2 - y1y2*x3x2

		l = float32(math.Sqrt(float64(x)*float64(x) + float64(y)*float64(y) + float64(z)*float64(z)))
		if l == 0 {
			l = 1.0
		}

		normals[x1] += x / l
		normals[y1] += y / l
		normals[z1] += z / l
		normals[x2] += x / l
		normals[y2] += y / l
		normals[z2] += z / l
		normals[x3] += x / l
		normals[y3] += y / l
		normals[z3] += z / l
	}
	for i := 0; i < len(normals)/3; i++ {
		x = normals[i*3]
		y = normals[i*3+1]
		z = normals[i*3+2]
		l = float32(math.Sqrt(float64(x)*float64(x) + float64(y)*float64(y) + float64(z)*float64(z)))
		if l == 0 {
			l = 1.0
		}
		normals[i*3] = x / l
		normals[i*3+1] = y / l
		normals[i*3+2] = z / l
	}
	return normals
}

func NewRibbon(paths [][]math32.Vector3, close bool) *Geometry {
	/*
		if len(paths) < 3 {
			close = false
		}
	*/
	c := NewGeometry()

	var ls, is []int // path lengths, path indexes
	positions := math32.NewArrayF32(0, 0)
	indices := math32.NewArrayU32(0, 0)

	i := 0
	for p := 0; p < len(paths); p++ {
		path := paths[p]
		l := len(path)
		ls = append(ls, l)
		is = append(is, i)
		for j := 0; j < l; j++ {
			positions.AppendVector3(&path[j])
		}
		i += l
	}

	l1 := ls[0] - 1 // path1 length
	l2 := ls[1] - 1 // path2 length
	min := l2
	if l1 < l2 {
		min = l1
	}
	p := 0
	i = 0
	for i <= min && p < len(ls)-1 {
		t := is[p+1] - is[p]

		indices.Append(uint32(i), uint32(i+t), uint32(i+1))
		indices.Append(uint32(i+t+1), uint32(i+1), uint32(i+t))
		i++
		if i == min {
			if close {
				indices.Append(uint32(i), uint32(i+t), uint32(is[p]))
				indices.Append(uint32(is[p]+t), uint32(is[p]), uint32(i+t))
			}
			p++
			if p == len(ls)-1 {
				break
			}
			l1 = ls[p] - 1
			l2 = ls[p+1] - 1
			i = is[p]
			if l1 < l2 {
				min = l1 + i
			} else {
				min = l2 + i
			}
		}
	}

	normals := math32.NewArrayF32(positions.Size(), positions.Size())
	normals = CalculateNormals(indices, positions, normals)

	c.SetIndices(indices)
	c.AddVBO(gls.NewVBO(positions).AddAttrib(gls.VertexPosition))
	c.AddVBO(gls.NewVBO(normals).AddAttrib(gls.VertexNormal))

	return c
}

func NewTube(path []math32.Vector3, radius float32, radialSegments int, close bool) *Geometry {
	l := len(path)

	var tangents, normals, binormals []math32.Vector3
	tangents = make([]math32.Vector3, l)
	normals = make([]math32.Vector3, l)
	binormals = make([]math32.Vector3, l)

	tangents[0] = *path[1].Clone().Sub(&path[0])
	tangents[0].Normalize()
	tangents[l-1] = *path[l-1].Clone().Sub(&path[l-2])
	tangents[l-1].Normalize()

	var tmpVertex *math32.Vector3
	if tangents[0].X != 1 {
		tmpVertex = math32.NewVector3(1, 0, 0)
	} else if tangents[0].Y != 1 {
		tmpVertex = math32.NewVector3(0, 1, 0)
	} else if tangents[0].Z != 1 {
		tmpVertex = math32.NewVector3(0, 0, 1)
	}

	normals[0] = *tangents[0].Clone().Cross(tmpVertex)
	normals[0].Normalize()
	binormals[0] = *tangents[0].Clone().Cross(&normals[0])
	binormals[0].Normalize()

	for i := 1; i < l; i++ {
		prev := *path[i].Clone().Sub(&path[i-1])
		if i < l-1 {
			cur := *path[i+1].Clone().Sub(&path[i])
			tangents[i] = *prev.Clone().Add(&cur)
			tangents[i].Normalize()

		}
		normals[i] = *binormals[i-1].Clone().Cross(&tangents[i])
		normals[i].Normalize()
		binormals[i] = *tangents[i].Clone().Cross(&normals[i])
		binormals[i].Normalize()
	}

	pi2 := math.Pi * 2
	step := pi2 / float64(radialSegments)

	var radialPaths [][]math32.Vector3
	for i := 0; i < l; i++ {
		var radialPath []math32.Vector3
		var ang float32
		for ang = 0.0; ang < float32(pi2); ang += float32(step) {
			matrix := math32.NewMatrix4()
			matrix.MakeRotationAxis(&tangents[i], ang)

			x := normals[i].X
			y := normals[i].Y
			z := normals[i].Z
			rw := 1 / (x*matrix[3] + y*matrix[7] + z*matrix[11] + matrix[15])
			newX := (x*matrix[0] + y*matrix[4] + z*matrix[8] + matrix[12]) * rw
			newY := (x*matrix[1] + y*matrix[5] + z*matrix[9] + matrix[13]) * rw
			newZ := (x*matrix[2] + y*matrix[6] + z*matrix[10] + matrix[14]) * rw

			rotated := math32.NewVector3(newX, newY, newZ).MultiplyScalar(radius).Add(&path[i])
			radialPath = append(radialPath, *rotated)
		}
		radialPaths = append(radialPaths, radialPath)
	}

	return NewRibbon(radialPaths, close)
}
