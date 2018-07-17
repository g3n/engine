// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package collada

import (
	"fmt"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/math32"
	"reflect"
	"strings"
)

// GetGeometry returns a pointer to an instance of the geometry
// with the specified id in the Collada document, its primitive type
// and an error. If no previous instance of the geometry was found
// the geometry is created
func (d *Decoder) GetGeometry(id string) (geometry.IGeometry, uint32, error) {

	// If geometry already created, returns it
	ginst, ok := d.geometries[id]
	if ok {
		return ginst.geom, ginst.ptype, nil
	}

	// Creates geometry and saves it associated with its id
	geom, ptype, err := d.NewGeometry(id)
	if err != nil {
		return nil, 0, err
	}
	d.geometries[id] = geomInstance{geom, ptype}

	return geom, ptype, nil
}

// NewGeometry creates and returns a pointer to a new instance of the geometry
// with the specified id in the Collada document, its primitive type and and error.
func (d *Decoder) NewGeometry(id string) (geometry.IGeometry, uint32, error) {

	id = strings.TrimPrefix(id, "#")
	// Look for geometry with specified id in the dom
	var geo *Geometry
	for _, g := range d.dom.LibraryGeometries.Geometry {
		if g.Id == id {
			geo = g
			break
		}
	}
	if geo == nil {
		return nil, 0, fmt.Errorf("Geometry:%s not found", id)
	}

	// Geometry type
	switch gt := geo.GeometricElement.(type) {
	// Collada mesh category includes points, lines, linestrips, triangles,
	// triangle fans, triangle strips and polygons.
	case *Mesh:
		return newMesh(gt)
		// B-Spline
		// Bezier
		// NURBS
		// Patch
	default:
		return nil, 0, fmt.Errorf("GeometryElement:%T not supported", gt)
	}
}

func newMesh(m *Mesh) (*geometry.Geometry, uint32, error) {

	// If no primitive elements present, it is a mesh of points
	if len(m.PrimitiveElements) == 0 {
		return newMeshPoints(m)
	}

	// All the primitive elements must be of the same type
	var etype reflect.Type
	for i := 0; i < len(m.PrimitiveElements); i++ {
		el := m.PrimitiveElements[i]
		if i == 0 {
			etype = reflect.TypeOf(el)
		} else {
			if reflect.TypeOf(el) != etype {
				return nil, 0, fmt.Errorf("primitive elements of different types")
			}
		}
	}

	pei := m.PrimitiveElements[0]
	switch pet := pei.(type) {
	case *Polylist:
		return newMeshPolylist(m, m.PrimitiveElements)
	case *Triangles:
		return newMeshTriangles(m, pet)
	case *Lines:
		return newMeshLines(m, pet)
	case *LineStrips:
		return newMeshLineStrips(m, pet)
	case *Trifans:
		return newMeshTrifans(m, pet)
	case *Tristrips:
		return newMeshTristrips(m, pet)
	default:
		return nil, 0, fmt.Errorf("PrimitiveElement:%T not supported", pet)
	}
}

// Creates a geometry from a polylist
// Only triangles are supported
func newMeshPolylist(m *Mesh, pels []interface{}) (*geometry.Geometry, uint32, error) {

	// Get vertices positions
	if len(m.Vertices.Input) != 1 {
		return nil, 0, fmt.Errorf("Mesh.Vertices.Input length not supported")
	}
	vinp := m.Vertices.Input[0]
	if vinp.Semantic != "POSITION" {
		return nil, 0, fmt.Errorf("Mesh.Vertices.Input.Semantic:%s not supported", vinp.Semantic)
	}

	// Get vertices input source
	inps := getMeshSource(m, vinp.Source)
	if inps == nil {
		return nil, 0, fmt.Errorf("Source:%s not found", vinp.Source)
	}

	// Get vertices input float array
	// Ignore Accessor (??)
	posArray, ok := inps.ArrayElement.(*FloatArray)
	if !ok {
		return nil, 0, fmt.Errorf("Mesh.Vertices.Input.Source not FloatArray")
	}

	// Creates buffers
	positions := math32.NewArrayF32(0, 0)
	normals := math32.NewArrayF32(0, 0)
	uvs := math32.NewArrayF32(0, 0)
	indices := math32.NewArrayU32(0, 0)

	// Creates vertices attributes map for reusing indices
	mVindex := make(map[[8]float32]uint32)
	var index uint32
	geomGroups := make([]geometry.Group, 0)
	groupMatindex := 0
	// For each Polylist
	for _, pel := range pels {
		// Checks if element is Polylist
		pl, ok := pel.(*Polylist)
		if !ok {
			return nil, 0, fmt.Errorf("Element is not a Polylist")
		}
		// If Polylist has not inputs, ignore
		if pl.Input == nil || len(pl.Input) == 0 {
			continue
		}
		// Checks if all Vcount elements are triangles
		for _, v := range pl.Vcount {
			if v != 3 {
				return nil, 0, fmt.Errorf("Only triangles are supported in Polylist")
			}
		}
		// Get VERTEX input
		inpVertex := getInputSemantic(pl.Input, "VERTEX")
		if inpVertex == nil {
			return nil, 0, fmt.Errorf("VERTEX input not found")
		}

		// Get optional NORMAL input
		inpNormal := getInputSemantic(pl.Input, "NORMAL")
		var normArray *FloatArray
		if inpNormal != nil {
			// Get normals source
			source := getMeshSource(m, inpNormal.Source)
			if source == nil {
				return nil, 0, fmt.Errorf("NORMAL source:%s not found", inpNormal.Source)
			}
			// Get normals source float array
			normArray, ok = source.ArrayElement.(*FloatArray)
			if !ok {
				return nil, 0, fmt.Errorf("NORMAL source:%s not float array", inpNormal.Source)
			}
		}

		// Get optional TEXCOORD input
		inpTexcoord := getInputSemantic(pl.Input, "TEXCOORD")
		var texArray *FloatArray
		if inpTexcoord != nil {
			// Get texture coordinates source
			source := getMeshSource(m, inpTexcoord.Source)
			if source == nil {
				return nil, 0, fmt.Errorf("TEXCOORD source:%s not found", inpTexcoord.Source)
			}
			// Get texture coordinates source float array
			texArray, ok = source.ArrayElement.(*FloatArray)
			if !ok {
				return nil, 0, fmt.Errorf("TEXCOORD source:%s not float array", inpTexcoord.Source)
			}
		}

		// Initialize geometry group
		groupStart := indices.Size()
		// For each primitive index
		inputCount := len(pl.Input)
		for i := 0; i < len(pl.P); i += inputCount {
			// Vertex attributes: position(3) + normal(3) + uv(2)
			var vx [8]float32

			// Vertex position
			posIndex := pl.P[i+inpVertex.Offset] * 3
			// Get position vector and appends to its buffer
			vx[0] = posArray.Data[posIndex]
			vx[1] = posArray.Data[posIndex+1]
			vx[2] = posArray.Data[posIndex+2]

			// Optional vertex normal
			if inpNormal != nil {
				// Get normal index from P
				normIndex := pl.P[i+inpNormal.Offset] * 3
				// Get normal vector and appends to its buffer
				vx[3] = normArray.Data[normIndex]
				vx[4] = normArray.Data[normIndex+1]
				vx[5] = normArray.Data[normIndex+2]
			}

			// Optional vertex texture coordinate
			if inpTexcoord != nil {
				// Get normal index from P
				texIndex := pl.P[i+inpTexcoord.Offset] * 2
				// Get normal vector and appends to its buffer
				vx[6] = texArray.Data[texIndex]
				vx[7] = texArray.Data[texIndex+1]
			}

			// If this vertex and its attributes has already been appended,
			// reuse it, adding its index to the index buffer
			// to reuse its index
			idx, ok := mVindex[vx]
			if ok {
				indices.Append(idx)
				continue
			}
			// Appends new vertex position and attributes to its buffers
			positions.Append(vx[0], vx[1], vx[2])
			if inpNormal != nil {
				normals.Append(vx[3], vx[4], vx[5])
			}
			if inpTexcoord != nil {
				uvs.Append(vx[6], vx[7])
			}
			indices.Append(index)
			// Save the index to this vertex position and attributes for
			// future reuse
			mVindex[vx] = index
			index++
		}
		// Adds this geometry group to the list
		geomGroups = append(geomGroups, geometry.Group{
			Start:    groupStart,
			Count:    indices.Size() - groupStart,
			Matindex: groupMatindex,
			Matid:    pl.Material,
		})
		groupMatindex++
	}

	// Debug dump
	//for i := 0; i < positions.Size()/3; i++ {
	//    vidx := i*3
	//    msg := fmt.Sprintf("i:%2d position:%v %v %v",
	//        i, positions.Get(vidx), positions.Get(vidx+1), positions.Get(vidx+2))
	//    if normals.Size() > 0 {
	//        msg += fmt.Sprintf("\tnormal:%v %v %v",
	//            normals.Get(vidx), normals.Get(vidx+1), normals.Get(vidx+2))
	//    }
	//    if uvs.Size() > 0 {
	//	    msg += fmt.Sprintf("\tuv:%v %v", uvs.Get(i*2), uvs.Get(i*2+1))
	//    }
	//    log.Debug("%s", msg)
	//}
	//log.Debug("indices(%d):%v", indices.Size(), indices)
	//log.Debug("groups:%v", geomGroups)

	// Creates geometry
	geom := geometry.NewGeometry()

	// Creates VBO with vertex positions
	geom.AddVBO(gls.NewVBO(positions).AddAttrib(gls.VertexPosition))

	// Creates VBO with vertex normals
	if normals.Size() > 0 {
		geom.AddVBO(gls.NewVBO(normals).AddAttrib(gls.VertexNormal))
	}

	// Creates VBO with uv coordinates
	if uvs.Size() > 0 {
		geom.AddVBO(gls.NewVBO(uvs).AddAttrib(gls.VertexTexcoord))
	}

	// Sets the geometry indices buffer
	geom.SetIndices(indices)

	// Add material groups to the geometry
	geom.AddGroupList(geomGroups)

	return geom, gls.TRIANGLES, nil
}

func newMeshTriangles(m *Mesh, tr *Triangles) (*geometry.Geometry, uint32, error) {

	return nil, 0, fmt.Errorf("not implemented yet")
}

func newMeshLines(m *Mesh, ln *Lines) (*geometry.Geometry, uint32, error) {

	if ln.Input == nil || len(ln.Input) == 0 {
		return nil, 0, fmt.Errorf("No inputs in lines")
	}

	// Get vertices positions
	if len(m.Vertices.Input) != 1 {
		return nil, 0, fmt.Errorf("Mesh.Vertices.Input length not supported")
	}
	vinp := m.Vertices.Input[0]
	if vinp.Semantic != "POSITION" {
		return nil, 0, fmt.Errorf("Mesh.Vertices.Input.Semantic:%s not supported", vinp.Semantic)
	}

	// Get vertices input source
	inps := getMeshSource(m, vinp.Source)
	if inps == nil {
		return nil, 0, fmt.Errorf("Source:%s not found", vinp.Source)
	}

	// Get vertices input float array
	// Ignore Accessor (??)
	posArray, ok := inps.ArrayElement.(*FloatArray)
	if !ok {
		return nil, 0, fmt.Errorf("Mesh.Vertices.Input.Source not FloatArray")
	}

	// Get VERTEX input
	inpVertex := getInputSemantic(ln.Input, "VERTEX")
	if inpVertex == nil {
		return nil, 0, fmt.Errorf("VERTEX input not found")
	}

	// Creates buffers
	positions := math32.NewArrayF32(0, 0)
	indices := math32.NewArrayU32(0, 0)

	mVindex := make(map[[3]float32]uint32)
	inputCount := len(ln.Input)
	var index uint32
	for i := 0; i < len(ln.P); i += inputCount {
		// Vertex position
		var vx [3]float32
		// Get position index from P
		posIndex := ln.P[i+inpVertex.Offset] * 3
		// Get position vector and appends to its buffer
		vx[0] = posArray.Data[posIndex]
		vx[1] = posArray.Data[posIndex+1]
		vx[2] = posArray.Data[posIndex+2]

		// If this vertex and its attributes has already been appended,
		// reuse it, adding its index to the index buffer
		// to reuse its index
		idx, ok := mVindex[vx]
		if ok {
			indices.Append(idx)
			continue
		}
		// Appends new vertex position and attributes to its buffers
		positions.Append(vx[0], vx[1], vx[2])
		indices.Append(index)
		// Save the index to this vertex position and attributes for
		// future reuse
		//mVindex[vx] = index
		index++
	}

	// Creates geometry
	geom := geometry.NewGeometry()

	// Creates VBO with vertex positions
	geom.AddVBO(gls.NewVBO(positions).AddAttrib(gls.VertexPosition))

	// Sets the geometry indices buffer
	geom.SetIndices(indices)
	return geom, gls.LINES, nil
}

func newMeshLineStrips(m *Mesh, ls *LineStrips) (*geometry.Geometry, uint32, error) {

	return nil, 0, fmt.Errorf("not implemented yet")
}

func newMeshTrifans(m *Mesh, ls *Trifans) (*geometry.Geometry, uint32, error) {

	return nil, 0, fmt.Errorf("not implemented yet")
}

func newMeshTristrips(m *Mesh, ls *Tristrips) (*geometry.Geometry, uint32, error) {

	return nil, 0, fmt.Errorf("not implemented yet")
}

// Creates and returns pointer to a new geometry for POINTS
func newMeshPoints(m *Mesh) (*geometry.Geometry, uint32, error) {

	// Get vertices positions
	if len(m.Vertices.Input) != 1 {
		return nil, 0, fmt.Errorf("Mesh.Vertices.Input length not supported")
	}
	vinp := m.Vertices.Input[0]
	if vinp.Semantic != "POSITION" {
		return nil, 0, fmt.Errorf("Mesh.Vertices.Input.Semantic:%s not supported", vinp.Semantic)
	}

	// Get vertices input source
	inps := getMeshSource(m, vinp.Source)
	if inps == nil {
		return nil, 0, fmt.Errorf("Source:%s not found", vinp.Source)
	}

	// Get vertices input float array
	// Ignore Accessor (??)
	posArray, ok := inps.ArrayElement.(*FloatArray)
	if !ok {
		return nil, 0, fmt.Errorf("Mesh.Vertices.Input.Source not FloatArray")
	}
	// Creates buffer and copy data
	positions := math32.NewArrayF32(posArray.Count, posArray.Count)
	//positions.CopyFrom(posArray.Data)
	copy(positions, posArray.Data)

	// Creates geometry and add VBO with vertex positions
	geom := geometry.NewGeometry()
	geom.AddVBO(gls.NewVBO(positions).AddAttrib(gls.VertexPosition))
	return geom, gls.POINTS, nil
}

func getMeshSource(m *Mesh, uri string) *Source {

	id := strings.TrimPrefix(uri, "#")
	for _, s := range m.Source {
		if s.Id == id {
			return s
		}
	}
	return nil
}

func getInputSemantic(inps []InputShared, semantic string) *InputShared {

	for i := 0; i < len(inps); i++ {
		if inps[i].Semantic == semantic {
			return &inps[i]
		}
	}
	return nil
}
