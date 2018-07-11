// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gltf

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"unsafe"

	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/texture"
)

// ParseJSON parses the glTF data from the specified JSON file
// and returns a pointer to the parsed structure.
func ParseJSON(filename string) (*GLTF, error) {

	// Open file
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	// Extract path from file
	path := filepath.Dir(filename)
	defer f.Close()
	return ParseJSONReader(f, path)
}

// ParseJSONReader parses the glTF JSON data from the specified reader
// and returns a pointer to the parsed structure
func ParseJSONReader(r io.Reader, path string) (*GLTF, error) {

	g := new(GLTF)
	g.path = path

	dec := json.NewDecoder(r)
	err := dec.Decode(g)
	if err != nil {
		return nil, err
	}
	return g, nil
}

// ParseBin parses the glTF data from the specified binary file
// and returns a pointer to the parsed structure.
func ParseBin(filename string) (*GLTF, error) {

	// Open file
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	// Extract path from file
	path := filepath.Dir(filename)
	defer f.Close()
	return ParseBinReader(f, path)
}

// ParseBinReader parses the glTF data from the specified binary reader
// and returns a pointer to the parsed structure
func ParseBinReader(r io.Reader, path string) (*GLTF, error) {

	// Reads header
	var header GLBHeader
	err := binary.Read(r, binary.LittleEndian, &header)
	if err != nil {
		return nil, err
	}

	// Checks magic and version
	if header.Magic != GLBMagic {
		return nil, fmt.Errorf("Invalid GLB Magic field")
	}
	if header.Version < 2 {
		return nil, fmt.Errorf("GLB version:%v not supported", header.Version)
	}

	// Read first chunk (JSON)
	buf, err := readChunk(r, GLBJson)
	if err != nil {
		return nil, err
	}

	// Parse JSON into gltf object
	bb := bytes.NewBuffer(buf)
	gltf, err := ParseJSONReader(bb, path)
	if err != nil {
		return nil, err
	}

	// Check for and read second chunk (binary, optional)
	data, err := readChunk(r, GLBBin)
	if err != nil {
		return nil, err
	}

	gltf.data = data

	return gltf, nil
}

// readChunk reads a GLB chunk with the specified type and returns the data in a byte array.
func readChunk(r io.Reader, chunkType uint32) ([]byte, error) {

	// Read chunk header
	var chunk GLBChunk
	err := binary.Read(r, binary.LittleEndian, &chunk)
	if err != nil {
		if err == io.EOF {
			return nil, nil
		}
		return nil, err
	}

	// Check chunk type
	if chunk.Type != chunkType {
		return nil, fmt.Errorf("expected GLB chunk type [%v] but obtained chunk type [%v]", chunkType, chunk.Type)
	}

	// Read chunk data
	data := make([]byte, chunk.Length)
	err = binary.Read(r, binary.LittleEndian, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// NewScene creates a parent Node which contains all nodes contained by
// the specified scene index from the GLTF Scenes array
func (g *GLTF) NewScene(si int) (core.INode, error) {

	// Check if provided scene index is valid
	if si < 0 || si >= len(g.Scenes) {
		return nil, fmt.Errorf("invalid scene index")
	}
	s := g.Scenes[si]

	scene := core.NewNode()
	scene.SetName(s.Name)
	for i := 0; i < len(s.Nodes); i++ {
		child, err := g.NewNode(i)
		if err != nil {
			return nil, err
		}
		scene.Add(child)
	}
	return scene, nil
}

// NewNode creates and returns a new Node described by the specified index
// in the decoded GLTF Nodes array.
func (g *GLTF) NewNode(i int) (core.INode, error) {

	var in core.INode
	var err error
	node := g.Nodes[i]

	// Check if the node is a Mesh (triangles, lines, etc...)
	if node.Mesh != nil {
		in, err = g.loadMesh(*node.Mesh)
		if err != nil {
			return nil, err
		}
		// Check if the node is Camera
	} else if node.Camera != nil {
		in, err = g.loadCamera(*node.Camera)
		if err != nil {
			return nil, err
		}
		// Other cases, return empty node
	} else {
		in = core.NewNode()
	}

	// Get *core.Node from core.INode
	n := in.GetNode()
	n.SetName(node.Name)

	// If defined, set node local transformation matrix
	if node.Matrix != nil {
		n.SetMatrix((*math32.Matrix4)(node.Matrix))
		// Otherwise, check rotation, scale and translation fields
	} else {
		// Rotation quaternion
		if node.Rotation != nil {
			log.Error("Rotation:%v", node.Translation)
			n.SetQuaternion(node.Rotation[0], node.Rotation[1], node.Rotation[2], node.Rotation[3])
		}
		// Scale
		if node.Scale != nil {
			log.Error("Scale:%v", node.Translation)
			n.SetScale(node.Scale[0], node.Scale[1], node.Scale[2])
		}
		// Translation
		if node.Translation != nil {
			log.Error("Translation:%v", node.Translation)
			n.SetPosition(node.Translation[0], node.Translation[1], node.Translation[2])
		}
	}

	// Recursively load node children  and add them to the parent
	for _, ci := range node.Children {
		child, err := g.NewNode(ci)
		if err != nil {
			return nil, err
		}
		n.Add(child)
	}

	return in, nil
}

// loadCamera creates and returns a Camera Node
// from the specified GLTF.Cameras index
func (g *GLTF) loadCamera(ci int) (core.INode, error) {

	camDesc := g.Cameras[ci]
	if camDesc.Type == "perspective" {
		desc := camDesc.Perspective
		fov := 360 * (desc.Yfov) / 2 * math32.Pi
		aspect := float32(2) // TODO how to get the current aspect ratio of the viewport from here ?
		if desc.AspectRatio != nil {
			aspect = *desc.AspectRatio
		}
		far := float32(2E6)
		if desc.Zfar != nil {
			far = *desc.Zfar
		}
		cam := camera.NewPerspective(fov, aspect, desc.Znear, far)
		return cam, nil
	}

	if camDesc.Type == "orthographic" {
		desc := camDesc.Orthographic
		cam := camera.NewOrthographic(desc.Xmag/-2, desc.Xmag/2, desc.Ymag/2, desc.Ymag/-2, desc.Znear, desc.Zfar)
		return cam, nil

	}

	return nil, fmt.Errorf("unsupported camera type: %s", camDesc.Type)
}

// loadMesh creates and returns a Graphic Node (graphic.Mesh, graphic.Lines, graphic.Points, etc)
// from the specified GLTF Mesh index
func (g *GLTF) loadMesh(mi int) (core.INode, error) {

	var err error
	m := g.Meshes[mi]

	// Create container node
	meshNode := core.NewNode()

	for i := 0; i < len(m.Primitives); i++ {

		// Get primitive information
		p := m.Primitives[i]

		// Indexed Geometry
		indices := math32.NewArrayU32(0, 0)
		if p.Indices != nil {
			pidx, err := g.loadIndices(*p.Indices)
			if err != nil {
				return nil, err
			}
			indices = append(indices, pidx...)
		} else {
			// Non-indexed primitive
			// indices array stay empty
		}

		// Load primitive material
		var grMat material.IMaterial
		if p.Material != nil {
			grMat, err = g.loadMaterial(*p.Material)
			if err != nil {
				return nil, err
			}
		} else {
			grMat = g.newDefaultMaterial()
		}

		// Create geometry
		geom := geometry.NewGeometry()

		// Load primitive attributes
		for name, aci := range p.Attributes {
			// TODO
			//interleaved := g.isInterleaved(aci)
			//if interleaved {
			//	buf, err := g.loadBufferView(*g.Accessors[aci].BufferView)
			//	if err != nil {
			//		return nil, err
			//	}
			//}
			if name == "POSITION" {
				ppos, err := g.loadVec3(aci)
				if err != nil {
					return nil, err
				}
				vbo := gls.NewVBO()
				vbo.AddAttrib("VertexPosition", 3)
				vbo.SetBuffer(ppos)
				geom.AddVBO(vbo)
				continue
			}
			if name == "NORMAL" {
				pnorms, err := g.loadVec3(aci)
				if err != nil {
					return nil, err
				}
				vbo := gls.NewVBO()
				vbo.AddAttrib("VertexNormal", 3)
				vbo.SetBuffer(pnorms)
				geom.AddVBO(vbo)
				continue
			}
			if name == "TEXCOORD_0" {
				puvs, err := g.loadVec2(aci)
				if err != nil {
					return nil, err
				}
				vbo := gls.NewVBO()
				vbo.AddAttrib("VertexTexcoord", 2)
				vbo.SetBuffer(puvs)
				geom.AddVBO(vbo)
				continue
			}
		}

		// Creates Geometry and add attribute VBO
		if len(indices) > 0 {
			geom.SetIndices(indices)
		}

		//log.Error("positions:%v", positions)
		//log.Error("indices..:%v", indices)
		//log.Error("normals..:%v", normals)
		//log.Error("uvs0.....:%v", uvs0)
		//log.Error("VBUF size in number of floats:%v", len(vbuf))

		// Default mode is 4 (TRIANGLES)
		mode := TRIANGLES
		if p.Mode != nil {
			mode = *p.Mode
		}

		// Create Mesh
		if mode == TRIANGLES {
			primitiveMesh := graphic.NewMesh(geom, nil)
			primitiveMesh.AddMaterial(grMat, 0, 0)
			meshNode.Add(primitiveMesh)
		}
		// Create Lines
		if mode == LINES {
			primitiveMesh := graphic.NewLines(geom, grMat)
			meshNode.Add(primitiveMesh)
		}
		// Create LineStrip
		if mode == LINE_STRIP {
			primitiveMesh := graphic.NewLineStrip(geom, grMat)
			meshNode.Add(primitiveMesh)
		}
		// Create Points
		if mode == POINTS {
			primitiveMesh := graphic.NewPoints(geom, grMat)
			meshNode.Add(primitiveMesh)
		}
		// TODO error when primitive is unsupported
		// fmt.Errorf("Unsupported primitive:%v", mode)
	}

	return meshNode, nil
}

func (g *GLTF) newDefaultMaterial() material.IMaterial {

	return material.NewStandard(&math32.Color{0.5, 0.5, 0.5})
}

// loadMaterials loads the material specified by the material index
func (g *GLTF) loadMaterial(mi int) (material.IMaterial, error) {

	mat := g.Materials[mi]
	// Checks for material extensions
	if mat.Extensions != nil {
		for ext, v := range mat.Extensions {
			if ext == "KHR_materials_common" {
				return g.loadMaterialCommon(v)
			} else {
				return nil, fmt.Errorf("Unsupported extension:%s", ext)
			}
		}
		return nil, fmt.Errorf("Empty material extensions")
		// Material should be PBR
	} else {
		return g.loadMaterialPBR(&mat)
	}
}

// loadTextureInfo loads the texture specified by the TextureInfo pointer
func (g *GLTF) loadTextureInfo(ti *TextureInfo) (*texture.Texture2D, error) {

	return g.loadTexture(ti.Index)
}

// loadTexture loads the texture specified by its index
func (g *GLTF) loadTexture(texi int) (*texture.Texture2D, error) {

	// loads texture image
	texDesc := g.Textures[texi]
	img, err := g.loadImage(texDesc.Source)
	if err != nil {
		return nil, err
	}
	tex := texture.NewTexture2DFromRGBA(img)

	// Get sampler and apply texture parameters
	samp := g.Samplers[texDesc.Sampler]

	// Magnification filter
	magFilter := gls.NEAREST
	if samp.MagFilter != nil {
		magFilter = *samp.MagFilter
	}
	tex.SetMagFilter(uint32(magFilter))

	// Minification filter
	minFilter := gls.NEAREST
	if samp.MinFilter != nil {
		minFilter = *samp.MinFilter
	}
	tex.SetMinFilter(uint32(minFilter))

	// S coordinate wrapping mode
	wrapS := gls.REPEAT
	if samp.WrapS != nil {
		wrapS = *samp.WrapS
	}
	tex.SetWrapS(uint32(wrapS))

	// T coordinate wrapping mode
	wrapT := gls.REPEAT
	if samp.WrapT != nil {
		wrapT = *samp.WrapT
	}
	tex.SetWrapT(uint32(wrapT))

	return tex, nil
}

// loadImage loads the image specified by the index of GLTF.Images
// Image can be loaded from binary chunk file or data URI or external file.
func (g *GLTF) loadImage(ii int) (*image.RGBA, error) {

	log.Error("loadImage:%v", ii)
	imgDesc := g.Images[ii]
	var data []byte
	var err error
	// If Uri is empty, load image from GLB binary chunk
	if imgDesc.Uri == "" {
		bvi := imgDesc.BufferView
		if bvi == nil {
			return nil, fmt.Errorf("Image has empty URI and no BufferView")
		}
		bv := g.BufferViews[*bvi]
		offset := 0
		if bv.ByteOffset != nil {
			offset = *bv.ByteOffset
		}
		data = g.data[offset : offset+bv.ByteLength]
		// Checks if image URI is data URL
	} else if isDataURL(imgDesc.Uri) {
		data, err = loadDataURL(imgDesc.Uri)
		if err != nil {
			return nil, err
		}
		// Load image data from file
	} else {
		fpath := filepath.Join(g.path, imgDesc.Uri)
		f, err := os.Open(fpath)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		data, err = ioutil.ReadAll(f)
		if err != nil {
			return nil, err
		}
	}

	// Decodes image data
	bb := bytes.NewBuffer(data)
	img, _, err := image.Decode(bb)
	if err != nil {
		return nil, err
	}

	// Converts image to RGBA format
	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return nil, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	return rgba, nil
}

// loadVec3 load array of float32 values from the specified accessor index.
// The acesssor must have type of VEC3 and component type of FLOAT
func (g *GLTF) loadVec3(ai int) (math32.ArrayF32, error) {

	// Get Accessor for the specified index
	ac := g.Accessors[ai]
	if ac.BufferView == nil {
		return nil, fmt.Errorf("Accessor.BufferView == nil NOT SUPPORTED")
	}

	// Checks acessor ComponentType
	if ac.ComponentType != FLOAT {
		return nil, fmt.Errorf("Accessor.ComponentType != FLOAT NOT SUPPORTED")
	}

	// Checks acessor Type
	if ac.Type != VEC3 {
		return nil, fmt.Errorf("Accessor.ComponentType != VEC3 NOT SUPPORTED")
	}

	// Loads data from associated BufferView
	data, err := g.loadBufferView(*ac.BufferView)
	if err != nil {
		return nil, err
	}

	// Accessor offset into BufferView
	offset := 0
	if ac.ByteOffset != nil {
		offset = *ac.ByteOffset
	}
	data = data[offset:]

	arr := (*[1 << 30]float32)(unsafe.Pointer(&data[0]))[:ac.Count*3]
	return math32.ArrayF32(arr), nil
}

// loadVec2 load array of Vector2 from the specified accessor index
func (g *GLTF) loadVec2(ai int) (math32.ArrayF32, error) {

	// Get Accessor for the specified index
	ac := g.Accessors[ai]
	if ac.BufferView == nil {
		return nil, fmt.Errorf("Accessor.BufferView == nil NOT SUPPORTED")
	}

	// Checks acessor ComponentType
	if ac.ComponentType != FLOAT {
		return nil, fmt.Errorf("Accessor.ComponentType != FLOAT NOT SUPPORTED")
	}

	// Checks acessor Type
	if ac.Type != VEC2 {
		return nil, fmt.Errorf("Accessor.ComponentType != VEC2 NOT SUPPORTED")
	}

	// Loads data from associated BufferView
	data, err := g.loadBufferView(*ac.BufferView)
	if err != nil {
		return nil, err
	}

	// Accessor offset into BufferView
	offset := 0
	if ac.ByteOffset != nil {
		offset = *ac.ByteOffset
	}
	data = data[offset:]

	arr := (*[1 << 30]float32)(unsafe.Pointer(&data[0]))[:ac.Count*2]
	return math32.ArrayF32(arr), nil
}

// loadIndices load the indices array specified by the Accessor index.
func (g *GLTF) loadIndices(ai int) (math32.ArrayU32, error) {

	// Get Accessor for the specified index
	ac := g.Accessors[ai]
	if ac.BufferView == nil {
		return nil, fmt.Errorf("Accessor.BufferView == nil NOT SUPPORTED YET")
	}

	// Loads indices data from associated BufferView
	data, err := g.loadBufferView(*ac.BufferView)
	if err != nil {
		return nil, err
	}

	// Accessor offset into BufferView
	offset := 0
	if ac.ByteOffset != nil {
		offset = *ac.ByteOffset
	}
	data = data[offset:]

	// If index component is UNSIGNED_INT nothing to do
	if ac.ComponentType == UNSIGNED_INT {
		arr := (*[1 << 30]uint32)(unsafe.Pointer(&data[0]))[:ac.Count]
		return math32.ArrayU32(arr), nil
	}

	// Converts UNSIGNED_SHORT indices to UNSIGNED_INT
	if ac.ComponentType == UNSIGNED_SHORT {
		indices := math32.NewArrayU32(ac.Count, ac.Count)
		for i := 0; i < ac.Count; i++ {
			indices[i] = uint32(data[i*2]) + uint32(data[i*2+1])*256
		}
		return indices, nil
	}

	// Converts UNSIGNED_BYTE indices to UNSIGNED_INT
	if ac.ComponentType == UNSIGNED_BYTE {
		indices := math32.NewArrayU32(ac.Count, ac.Count)
		for i := 0; i < ac.Count; i++ {
			indices[i] = uint32(data[i])
		}
		return indices, nil
	}
	return nil, fmt.Errorf("Unsupported Accessor ComponentType:%v", ac.ComponentType)
}

// isInterleaves checks if the BufferView used by the specified Accessor index is
// interleaved or not
func (g *GLTF) isInterleaved(aci int) bool {

	// Get the Accessor's BufferView
	accessor := g.Accessors[aci]
	if accessor.BufferView == nil {
		return false
	}
	bv := g.BufferViews[*accessor.BufferView]

	// Calculates the size in bytes of a complete attribute
	itemSize := TypeSizes[accessor.Type]
	itemBytes := int(gls.FloatSize) * itemSize

	// If the BufferView stride is equal to the item size, the buffer is not interleaved
	if bv.ByteStride == nil {
		return false
	}
	if *bv.ByteStride == itemBytes {
		return false
	}
	return true
}

// loadBufferView loads and returns a byte slice with data from the specified
// BufferView index
func (g *GLTF) loadBufferView(bvi int) ([]byte, error) {

	bv := g.BufferViews[bvi]
	buf, err := g.loadBuffer(bv.Buffer)
	if err != nil {
		return nil, err
	}

	offset := 0
	if bv.ByteOffset != nil {
		offset = *bv.ByteOffset
	}
	return buf[offset : offset+bv.ByteLength], nil
}

// loadBuffer loads and returns the data from the specified GLTF Buffer index
func (g *GLTF) loadBuffer(bi int) ([]byte, error) {

	buf := &g.Buffers[bi]
	// If Buffer URI uses the chunk data field
	if buf.Uri == "" {
		return g.data, nil
	}

	// If buffer already loaded:
	log.Error("loadBuffer cache:%v", len(buf.data))
	if len(buf.data) > 0 {
		return buf.data, nil
	}

	// Checks if buffer URI is a data URI
	var data []byte
	var err error
	if isDataURL(buf.Uri) {
		data, err = loadDataURL(buf.Uri)
		if err != nil {
			return nil, err
		}
		// Loads external buffer file
	} else {
		log.Error("loadBuffer: loading file")
		// Try to load buffer from file
		fpath := filepath.Join(g.path, buf.Uri)
		f, err := os.Open(fpath)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		data, err = ioutil.ReadAll(f)
		if err != nil {
			return nil, err
		}
	}
	// Checks data length
	if len(data) != buf.ByteLength {
		return nil, fmt.Errorf("Buffer:%d read data length:%d expected:%d", bi, len(data), buf.ByteLength)
	}
	// Cache buffer data
	buf.data = data
	log.Error("cache data:%v", len(buf.data))
	return data, nil
}

// dataURL describes a decoded data url string
type dataURL struct {
	MediaType string
	Encoding  string
	Data      string
}

const (
	dataURLprefix = "data:"
	mimeBIN       = "application/octet-stream"
	mimePNG       = "image/png"
	mimeJPEG      = "image/jpeg"
)

var validMediaTypes = []string{mimeBIN, mimePNG, mimeJPEG}

// isDataURL checks if the specified string has the prefix of data URL
func isDataURL(url string) bool {

	if strings.HasPrefix(url, dataURLprefix) {
		return true
	}
	return false
}

// loadDataURL decodes the specified data URI string (base64)
func loadDataURL(url string) ([]byte, error) {

	var du dataURL
	err := parseDataURL(url, &du)
	if err != nil {
		return nil, err
	}

	// Checks for valid media type
	found := false
	for i := 0; i < len(validMediaTypes); i++ {
		if validMediaTypes[i] == du.MediaType {
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("Data URI media type:%s not supported", du.MediaType)
	}

	// Checks encoding
	if du.Encoding != "base64" {
		return nil, fmt.Errorf("Data URI encoding:%s not supported", du.Encoding)
	}

	// Decodes data from BASE64
	data, err := base64.StdEncoding.DecodeString(du.Data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// parseDataURL tries to parse the specified string as a data URL with the format:
// data:[<mediatype>][;base64],<data>
// and if successfull returns true and updates the specified pointer with the parsed fields.
func parseDataURL(url string, du *dataURL) error {

	// Checks prefix
	if !isDataURL(url) {
		return fmt.Errorf("Specified string is not a data URL")
	}

	// Separates header from data
	body := url[len(dataURLprefix):]
	parts := strings.Split(body, ",")
	if len(parts) != 2 {
		return fmt.Errorf("Data URI contains more than one ','")
	}
	du.Data = parts[1]

	// Separates media type from optional encoding
	res := strings.Split(parts[0], ";")
	du.MediaType = res[0]
	if len(res) < 2 {
		return nil
	}
	if len(res) >= 2 {
		du.Encoding = res[1]
	}
	return nil
}
