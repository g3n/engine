// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package obj
package obj

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/texture"
)

// Decoder contains all decoded data from the obj and mtl files
type Decoder struct {
	Objects       []Object             // decoded objects
	Matlib        string               // name of the material lib
	Materials     map[string]*Material // maps material name to object
	Vertices      math32.ArrayF32      // vertices positions array
	Normals       math32.ArrayF32      // vertices normals
	Uvs           math32.ArrayF32      // vertices texture coordinates
	Warnings      []string             // warning messages
	line          uint                 // current line number
	objCurrent    *Object              // current object
	matCurrent    *Material            // current material
	smoothCurrent bool                 // current smooth state
	mtlDir        string               // Directory of material file
}

// Object contains all information about one decoded object
type Object struct {
	Name      string   // Object name
	Faces     []Face   // Faces
	materials []string // Materials used in this object
}

// Face contains all information about an object face
type Face struct {
	Vertices []int  // Indices to the face vertices
	Uvs      []int  // Indices to the face UV coordinates
	Normals  []int  // Indices to the face normals
	Material string // Material name
	Smooth   bool   // Smooth face
}

// Material contains all information about an object material
type Material struct {
	Name       string       // Material name
	Illum      int          // Illumination model
	Opacity    float32      // Opacity factor
	Refraction float32      // Refraction factor
	Shininess  float32      // Shininess (specular exponent)
	Ambient    math32.Color // Ambient color reflectivity
	Diffuse    math32.Color // Diffuse color reflectivity
	Specular   math32.Color // Specular color reflectivity
	Emissive   math32.Color // Emissive color
	MapKd      string       // Texture file linked to diffuse color
}

// Local constants
const (
	blanks   = "\r\n\t "
	invINDEX = math.MaxUint32
	objType  = "obj"
	mtlType  = "mtl"
)

// Decode decodes the specified obj and mtl files returning a decoder
// object and an error.
func Decode(objpath string, mtlpath string) (*Decoder, error) {

	// Opens obj file
	fobj, err := os.Open(objpath)
	if err != nil {
		return nil, err
	}
	defer fobj.Close()

	// If path of material file not supplied,
	// try to use the base name of the obj file
	if len(mtlpath) == 0 {
		dir, objfile := filepath.Split(objpath)
		ext := filepath.Ext(objfile)
		mtlpath = dir + objfile[:len(objfile)-len(ext)] + ".mtl"
	}

	// Opens mtl file
	fmtl, err := os.Open(mtlpath)
	if err != nil {
		return nil, err
	}
	defer fmtl.Close()

	dec, err := DecodeReader(fobj, fmtl)
	if err != nil {
		return nil, err
	}
	dec.mtlDir = filepath.Dir(objpath)
	return dec, nil
}

// DecodeReader decodes the specified obj and mtl readers returning a decoder
// object and an error.
func DecodeReader(objreader, mtlreader io.Reader) (*Decoder, error) {

	dec := new(Decoder)
	dec.Objects = make([]Object, 0)
	dec.Warnings = make([]string, 0)
	dec.Materials = make(map[string]*Material)
	dec.Vertices = math32.NewArrayF32(0, 0)
	dec.Normals = math32.NewArrayF32(0, 0)
	dec.Uvs = math32.NewArrayF32(0, 0)
	dec.line = 1

	// Parses obj lines
	err := dec.parse(objreader, dec.parseObjLine)
	if err != nil {
		return nil, err
	}

	// Parses mtl lines
	dec.matCurrent = nil
	dec.line = 1
	err = dec.parse(mtlreader, dec.parseMtlLine)
	if err != nil {
		return nil, err
	}
	return dec, nil
}

// NewGroup creates and returns a group containing as children meshes
// with all the decoded objects.
// A group is returned even if there is only one object decoded.
func (dec *Decoder) NewGroup() (*core.Node, error) {

	group := core.NewNode()
	for i := 0; i < len(dec.Objects); i++ {
		mesh, err := dec.NewMesh(&dec.Objects[i])
		if err != nil {
			return nil, err
		}
		group.Add(mesh)
	}
	return group, nil
}

// NewMesh creates and returns a mesh from an specified decoded object.
func (dec *Decoder) NewMesh(obj *Object) (*graphic.Mesh, error) {

	// Creates object geometry
	geom, err := dec.NewGeometry(obj)
	if err != nil {
		return nil, err
	}

	// Single material
	if geom.GroupCount() == 1 {
		matName := obj.materials[0]
		matDesc := dec.Materials[matName]
		// Creates material
		mat := material.NewPhong(&matDesc.Diffuse)
		ambientColor := mat.AmbientColor()
		mat.SetAmbientColor(ambientColor.Multiply(&matDesc.Ambient))
		mat.SetSpecularColor(&matDesc.Specular)
		mat.SetShininess(matDesc.Shininess)
		// Loads material textures if specified
		err = dec.loadTex(&mat.Material, matDesc)
		if err != nil {
			return nil, err
		}
		return graphic.NewMesh(geom, mat), nil
	}

	// Multi material
	mesh := graphic.NewMesh(geom, nil)
	for idx := 0; idx < geom.GroupCount(); idx++ {
		group := geom.GroupAt(idx)
		matName := obj.materials[group.Matindex]
		matDesc := dec.Materials[matName]
		// Creates material
		matGroup := material.NewPhong(&matDesc.Diffuse)
		ambientColor := matGroup.AmbientColor()
		matGroup.SetAmbientColor(ambientColor.Multiply(&matDesc.Ambient))
		matGroup.SetSpecularColor(&matDesc.Specular)
		matGroup.SetShininess(matDesc.Shininess)
		// Loads material textures if specified
		err = dec.loadTex(&matGroup.Material, matDesc)
		if err != nil {
			return nil, err
		}
		mesh.AddGroupMaterial(matGroup, idx)
	}
	return mesh, nil
}

// NewGeometry generates and returns a geometry from the specified object
func (dec *Decoder) NewGeometry(obj *Object) (*geometry.Geometry, error) {

	geom := geometry.NewGeometry()

	// Create buffers
	positions := math32.NewArrayF32(0, 0)
	normals := math32.NewArrayF32(0, 0)
	uvs := math32.NewArrayF32(0, 0)
	indices := math32.NewArrayU32(0, 0)

	// copy all vertex info from the decoded Object, face and index to the geometry
	copyVertex := func(face *Face, idx int) {
		var vec3 math32.Vector3
		var vec2 math32.Vector2

		pos := positions.Size() / 3
		// Copy vertex position and append to geometry
		dec.Vertices.GetVector3(3*face.Vertices[idx], &vec3)
		positions.AppendVector3(&vec3)
		// Copy vertex normal and append to geometry
		if face.Normals[idx] != invINDEX {
			dec.Normals.GetVector3(3*face.Normals[idx], &vec3)
			normals.AppendVector3(&vec3)
		}
		// Copy vertex uv and append to geometry
		if face.Uvs[idx] != invINDEX {
			dec.Uvs.GetVector2(2*face.Uvs[idx], &vec2)
			uvs.AppendVector2(&vec2)
		}
		indices.Append(uint32(pos))
	}

	var group *geometry.Group
	matName := ""
	matIndex := 0
	for _, face := range obj.Faces {
		// If face material changed, starts a new group
		if face.Material != matName {
			group = geom.AddGroup(indices.Size(), 0, matIndex)
			matName = face.Material
			matIndex++
		}
		// Copy face vertices to geometry
		for idx := 1; idx < len(face.Vertices)-1; idx++ {
			copyVertex(&face, 0)
			copyVertex(&face, idx)
			copyVertex(&face, idx+1)
			group.Count += 3
		}
	}

	geom.SetIndices(indices)
	geom.AddVBO(gls.NewVBO(positions).AddAttrib(gls.VertexPosition))
	geom.AddVBO(gls.NewVBO(normals).AddAttrib(gls.VertexNormal))
	geom.AddVBO(gls.NewVBO(uvs).AddAttrib(gls.VertexTexcoord))

	return geom, nil
}

// loadTex loads textures described in the material descriptor into the
// specified material
func (dec *Decoder) loadTex(mat *material.Material, desc *Material) error {

	// Checks if material descriptor specified texture
	if desc.MapKd == "" {
		return nil
	}

	// Get texture file path
	// If texture file path is not absolute assumes it is relative
	// to the directory of the material file
	var texPath string
	if filepath.IsAbs(desc.MapKd) {
		texPath = desc.MapKd
	} else {
		texPath = filepath.Join(dec.mtlDir, desc.MapKd)
	}

	// Try to load texture from image file
	tex, err := texture.NewTexture2DFromImage(texPath)
	if err != nil {
		return err
	}
	mat.AddTexture(tex)
	return nil
}

// parse reads the lines from the specified reader and dispatch them
// to the specified line parser.
func (dec *Decoder) parse(reader io.Reader, parseLine func(string) error) error {

	bufin := bufio.NewReader(reader)
	dec.line = 1
	for {
		// Reads next line and abort on errors (not EOF)
		line, err := bufin.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}
		// Parses the line
		line = strings.Trim(line, blanks)
		perr := parseLine(line)
		if perr != nil {
			return perr
		}
		// If EOF ends of parsing.
		if err == io.EOF {
			break
		}
		dec.line++
	}
	return nil
}

// Parses obj file line, dispatching to specific parsers
func (dec *Decoder) parseObjLine(line string) error {

	// Ignore empty lines
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return nil
	}
	// Ignore comment lines
	ltype := fields[0]
	if strings.HasPrefix(ltype, "#") {
		return nil
	}
	switch ltype {
	// Material library
	case "mtllib":
		return dec.parseMatlib(fields[1:])
	// Object name
	case "o":
		return dec.parseObject(fields[1:])
	// Group names. We are considering "group" the same as "object"
	// This may not be right
	case "g":
		return dec.parseObject(fields[1:])
	// Vertex coordinate
	case "v":
		return dec.parseVertex(fields[1:])
	// Vertex normal coordinate
	case "vn":
		return dec.parseNormal(fields[1:])
	// Vertex texture coordinate
	case "vt":
		return dec.parseTex(fields[1:])
	// Face vertex
	case "f":
		return dec.parseFace(fields[1:])
	// Use material
	case "usemtl":
		return dec.parseUsemtl(fields[1:])
	// Smooth
	case "s":
		return dec.parseSmooth(fields[1:])
	default:
		dec.appendWarn(objType, "field not supported: "+ltype)
	}
	return nil
}

// Parses a mtllib line:
// mtllib <name>
func (dec *Decoder) parseMatlib(fields []string) error {

	if len(fields) < 1 {
		return errors.New("Object line (o) with less than 2 fields")
	}
	dec.Matlib = fields[0]
	return nil
}

// Parses an object line:
// o <name>
func (dec *Decoder) parseObject(fields []string) error {

	if len(fields) < 1 {
		return errors.New("Object line (o) with less than 2 fields")
	}
	var ob Object
	ob.Name = fields[0]
	ob.Faces = make([]Face, 0)
	ob.materials = make([]string, 0)
	dec.Objects = append(dec.Objects, ob)
	dec.objCurrent = &dec.Objects[len(dec.Objects)-1]
	return nil
}

// Parses a vertex position line
// v <x> <y> <z> [w]
func (dec *Decoder) parseVertex(fields []string) error {

	if len(fields) < 3 {
		return errors.New("Less than 3 vertices in 'v' line")
	}
	for _, f := range fields[:3] {
		val, err := strconv.ParseFloat(f, 32)
		if err != nil {
			return err
		}
		dec.Vertices.Append(float32(val))
	}
	return nil
}

// Parses a vertex normal line
// vn <x> <y> <z>
func (dec *Decoder) parseNormal(fields []string) error {

	if len(fields) < 3 {
		return errors.New("Less than 3 normals in 'vn' line")
	}
	for _, f := range fields[:3] {
		val, err := strconv.ParseFloat(f, 32)
		if err != nil {
			return err
		}
		dec.Normals.Append(float32(val))
	}
	return nil
}

// Parses a vertex texture coordinate line:
// vt <u> <v> <w>
func (dec *Decoder) parseTex(fields []string) error {

	if len(fields) < 2 {
		return errors.New("Less than 2 texture coords. in 'vt' line")
	}
	for _, f := range fields[:2] {
		val, err := strconv.ParseFloat(f, 32)
		if err != nil {
			return err
		}
		dec.Uvs.Append(float32(val))
	}
	return nil
}

// parseFace parses a face decription line:
// f v1[/vt1][/vn1] v2[/vt2][/vn2] v3[/vt3][/vn3] ...
func (dec *Decoder) parseFace(fields []string) error {

	// If current object has no material, appends last material if defined
	if len(dec.objCurrent.materials) == 0 && dec.matCurrent != nil {
		dec.objCurrent.materials = append(dec.objCurrent.materials, dec.matCurrent.Name)
	}

	if len(fields) < 3 {
		return dec.formatError("Face line with less 3 fields")
	}
	if dec.matCurrent == nil {
		return dec.formatError("No material defined")
	}
	var face Face
	face.Vertices = make([]int, len(fields))
	face.Uvs = make([]int, len(fields))
	face.Normals = make([]int, len(fields))
	face.Material = dec.matCurrent.Name
	face.Smooth = dec.smoothCurrent

	for pos, f := range fields {

		// Separate the current field in its components: v vt vn
		vfields := strings.Split(f, "/")
		if len(vfields) < 1 {
			return dec.formatError("Face field with no parts")
		}

		// Get the index of this vertex position (must always exist)
		val, err := strconv.ParseInt(vfields[0], 10, 32)
		if err != nil {
			return err
		}

		// Positive index is an absolute vertex index
		if val > 0 {
			face.Vertices[pos] = int(val - 1)
			// Negative vertex index is relative to the last parsed vertex
		} else if val < 0 {
			current := (len(dec.Vertices) / 3) - 1
			face.Vertices[pos] = current + int(val) + 1
			// Vertex index could never be 0
		} else {
			return dec.formatError("Face vertex index value equal to 0")
		}

		// Get the index of this vertex UV coordinate (optional)
		if len(vfields) > 1 && len(vfields[1]) > 0 {
			val, err := strconv.ParseInt(vfields[1], 10, 32)
			if err != nil {
				return err
			}

			// Positive index is an absolute UV index
			if val > 0 {
				face.Uvs[pos] = int(val - 1)
				// Negative vertex index is relative to the last parsed uv
			} else if val < 0 {
				current := (len(dec.Uvs) / 2) - 1
				face.Uvs[pos] = current + int(val) + 1
				// UV index could never be 0
			} else {
				return dec.formatError("Face uv index value equal to 0")
			}
		} else {
			face.Uvs[pos] = invINDEX
		}

		// Get the index of this vertex normal (optional)
		if len(vfields) >= 3 {
			val, err = strconv.ParseInt(vfields[2], 10, 32)
			if err != nil {
				return err
			}

			// Positive index is an absolute normal index
			if val > 0 {
				face.Normals[pos] = int(val - 1)
				// Negative vertex index is relative to the last parsed normal
			} else if val < 0 {
				current := (len(dec.Normals) / 3) - 1
				face.Normals[pos] = current + int(val) + 1
				// Normal index could never be 0
			} else {
				return dec.formatError("Face normal index value equal to 0")
			}
		} else {
			face.Normals[pos] = invINDEX
		}
	}
	// Appends this face to the current object
	dec.objCurrent.Faces = append(dec.objCurrent.Faces, face)
	return nil
}

// parseUsemtl parses a "usemtl" decription line:
// usemtl <name>
func (dec *Decoder) parseUsemtl(fields []string) error {

	if len(fields) < 1 {
		return dec.formatError("Usemtl with no fields")
	}

	// Checks if this material has already been parsed
	name := fields[0]
	mat := dec.Materials[name]
	// Creates material descriptor
	if mat == nil {
		mat = new(Material)
		mat.Name = name
		dec.Materials[name] = mat
	}
	dec.objCurrent.materials = append(dec.objCurrent.materials, name)
	// Set this as the current material
	dec.matCurrent = mat
	return nil
}

// parseSmooth parses a "s" decription line:
// s <0|1>
func (dec *Decoder) parseSmooth(fields []string) error {

	if len(fields) < 1 {
		return dec.formatError("'s' with no fields")
	}

	if fields[0] == "0" || fields[0] == "off" {
		dec.smoothCurrent = false
		return nil
	}
	if fields[0] == "1" || fields[0] == "on" {
		dec.smoothCurrent = true
		return nil
	}
	return dec.formatError("'s' with invalid value")
}

/******************************************************************************
mtl parse functions
*/

// Parses material file line, dispatching to specific parsers
func (dec *Decoder) parseMtlLine(line string) error {

	// Ignore empty lines
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return nil
	}
	// Ignore comment lines
	ltype := fields[0]
	if strings.HasPrefix(ltype, "#") {
		return nil
	}
	switch ltype {
	case "newmtl":
		return dec.parseNewmtl(fields[1:])
	case "d":
		return dec.parseDissolve(fields[1:])
	case "Ka":
		return dec.parseKa(fields[1:])
	case "Kd":
		return dec.parseKd(fields[1:])
	case "Ke":
		return dec.parseKe(fields[1:])
	case "Ks":
		return dec.parseKs(fields[1:])
	case "Ni":
		return dec.parseNi(fields[1:])
	case "Ns":
		return dec.parseNs(fields[1:])
	case "illum":
		return dec.parseIllum(fields[1:])
	case "map_Kd":
		return dec.parseMapKd(fields[1:])
	default:
		dec.appendWarn(mtlType, "field not supported: "+ltype)
	}
	return nil
}

// Parses new material definition
// newmtl <mat_name>
func (dec *Decoder) parseNewmtl(fields []string) error {

	if len(fields) < 1 {
		return dec.formatError("newmtl with no fields")
	}
	// Checks if material has already been seen
	name := fields[0]
	mat := dec.Materials[name]
	// Creates material descriptor
	if mat == nil {
		mat = new(Material)
		mat.Name = name
		dec.Materials[name] = mat
	}
	dec.matCurrent = mat
	return nil
}

// Parses the dissolve factor (opacity)
// d <factor>
func (dec *Decoder) parseDissolve(fields []string) error {

	if len(fields) < 1 {
		return dec.formatError("'d' with no fields")
	}
	val, err := strconv.ParseFloat(fields[0], 32)
	if err != nil {
		return dec.formatError("'d' parse float error")
	}
	dec.matCurrent.Opacity = float32(val)
	return nil
}

// Parses ambient reflectivity:
// Ka r g b
func (dec *Decoder) parseKa(fields []string) error {

	if len(fields) < 3 {
		return dec.formatError("'Ka' with less than 3 fields")
	}
	var colors [3]float32
	for pos, f := range fields[:3] {
		val, err := strconv.ParseFloat(f, 32)
		if err != nil {
			return err
		}
		colors[pos] = float32(val)
	}
	dec.matCurrent.Ambient.Set(colors[0], colors[1], colors[2])
	return nil
}

// Parses diffuse reflectivity:
// Kd r g b
func (dec *Decoder) parseKd(fields []string) error {

	if len(fields) < 3 {
		return dec.formatError("'Kd' with less than 3 fields")
	}
	var colors [3]float32
	for pos, f := range fields[:3] {
		val, err := strconv.ParseFloat(f, 32)
		if err != nil {
			return err
		}
		colors[pos] = float32(val)
	}
	dec.matCurrent.Diffuse.Set(colors[0], colors[1], colors[2])
	return nil
}

// Parses emissive color:
// Ke r g b
func (dec *Decoder) parseKe(fields []string) error {

	if len(fields) < 3 {
		return dec.formatError("'Ke' with less than 3 fields")
	}
	var colors [3]float32
	for pos, f := range fields[:3] {
		val, err := strconv.ParseFloat(f, 32)
		if err != nil {
			return err
		}
		colors[pos] = float32(val)
	}
	dec.matCurrent.Emissive.Set(colors[0], colors[1], colors[2])
	return nil
}

// Parses specular reflectivity:
// Ks r g b
func (dec *Decoder) parseKs(fields []string) error {

	if len(fields) < 3 {
		return dec.formatError("'Ks' with less than 3 fields")
	}
	var colors [3]float32
	for pos, f := range fields[:3] {
		val, err := strconv.ParseFloat(f, 32)
		if err != nil {
			return err
		}
		colors[pos] = float32(val)
	}
	dec.matCurrent.Specular.Set(colors[0], colors[1], colors[2])
	return nil
}

// Parses optical density, also known as index of refraction
// Ni <optical_density>
func (dec *Decoder) parseNi(fields []string) error {

	if len(fields) < 1 {
		return dec.formatError("'Ni' with no fields")
	}
	val, err := strconv.ParseFloat(fields[0], 32)
	if err != nil {
		return dec.formatError("'d' parse float error")
	}
	dec.matCurrent.Refraction = float32(val)
	return nil
}

// Parses specular exponent
// Ns <specular_exponent>
func (dec *Decoder) parseNs(fields []string) error {

	if len(fields) < 1 {
		return dec.formatError("'Ns' with no fields")
	}
	val, err := strconv.ParseFloat(fields[0], 32)
	if err != nil {
		return dec.formatError("'d' parse float error")
	}
	dec.matCurrent.Shininess = float32(val)
	return nil
}

// Parses illumination model (0 to 10)
// illum <ilum_#>
func (dec *Decoder) parseIllum(fields []string) error {

	if len(fields) < 1 {
		return dec.formatError("'illum' with no fields")
	}
	val, err := strconv.ParseUint(fields[0], 10, 32)
	if err != nil {
		return dec.formatError("'d' parse int error")
	}
	dec.matCurrent.Illum = int(val)
	return nil
}

// Parses color texture linked to the diffuse reflectivity of the material
// map_Kd [-options] <filename>
func (dec *Decoder) parseMapKd(fields []string) error {

	if len(fields) < 1 {
		return dec.formatError("No fields")
	}
	dec.matCurrent.MapKd = fields[0]
	return nil
}

func (dec *Decoder) formatError(msg string) error {

	return fmt.Errorf("%s in line:%d", msg, dec.line)
}

func (dec *Decoder) appendWarn(ftype string, msg string) {

	wline := fmt.Sprintf("%s(%d): %s", ftype, dec.line, msg)
	dec.Warnings = append(dec.Warnings, wline)
}
