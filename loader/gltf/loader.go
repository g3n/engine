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
	"github.com/g3n/engine/animation"
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

	// Read header
	var header GLBHeader
	err := binary.Read(r, binary.LittleEndian, &header)
	if err != nil {
		return nil, err
	}

	// Check magic and version
	if header.Magic != GLBMagic {
		return nil, fmt.Errorf("invalid GLB Magic field")
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
		return nil, fmt.Errorf("expected GLB chunk type [%v] but found [%v]", chunkType, chunk.Type)
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
// the specified scene index from the GLTF Scenes array.
func (g *GLTF) NewScene(si int) (core.INode, error) {

	log.Debug("Loading Scene %d", si)

	// Check if provided scene index is valid
	if si < 0 || si >= len(g.Scenes) {
		return nil, fmt.Errorf("invalid scene index")
	}
	s := g.Scenes[si]

	scene := core.NewNode()
	scene.SetName(s.Name)
	for _, ni := range s.Nodes {
		child, err := g.NewNode(ni)
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

	log.Debug("Loading Node %d", i)

	var in core.INode
	var err error
	nodeData := g.Nodes[i]

	// Check if the node is a Mesh (triangles, lines, etc...)
	if nodeData.Mesh != nil {
		in, err = g.NewMesh(*nodeData.Mesh)
		if err != nil {
			return nil, err
		}
		// Check if the node is Camera
	} else if nodeData.Camera != nil {
		in, err = g.NewCamera(*nodeData.Camera)
		if err != nil {
			return nil, err
		}
		// Other cases, return empty node
	} else {
		log.Debug("Empty Node")
		in = core.NewNode()
	}

	// Cache inode in nodeData
	g.Nodes[i].node = in

	// Get *core.Node from core.INode
	node := in.GetNode()
	node.SetName(nodeData.Name)

	// If defined, set node local transformation matrix
	if nodeData.Matrix != nil {
		node.SetMatrix((*math32.Matrix4)(nodeData.Matrix))
		// Otherwise, check rotation, scale and translation fields
	} else {
		// Rotation quaternion
		if nodeData.Rotation != nil {
			node.SetQuaternion(nodeData.Rotation[0], nodeData.Rotation[1], nodeData.Rotation[2], nodeData.Rotation[3])
		}
		// Scale
		if nodeData.Scale != nil {
			node.SetScale(nodeData.Scale[0], nodeData.Scale[1], nodeData.Scale[2])
		}
		// Translation
		if nodeData.Translation != nil {
			node.SetPosition(nodeData.Translation[0], nodeData.Translation[1], nodeData.Translation[2])
		}
	}

	// Recursively load node children and add them to the parent
	for _, ci := range nodeData.Children {
		child, err := g.NewNode(ci)
		if err != nil {
			return nil, err
		}
		node.Add(child)
	}

	return in, nil
}

// NewAnimation creates an Animation for the specified
// the animation index from the GLTF Animations array.
func (g *GLTF) NewAnimation(i int) (*animation.Animation, error) {

	log.Debug("Loading Animation %d", i)

	// Check if provided scene index is valid
	if i < 0 || i >= len(g.Animations) {
		return nil, fmt.Errorf("invalid animation index")
	}
	a := g.Animations[i]

	anim := animation.NewAnimation()
	anim.SetName(a.Name)
	for i := 0; i < len(a.Channels); i++ {

		chData := a.Channels[i]
		target := chData.Target
		sampler := a.Samplers[chData.Sampler]
		node := g.Nodes[target.Node].node
		// TODO Instantiate node if not exists ?

		var validTypes []string
		var validComponentTypes []int

		var ch animation.IChannel
		if target.Path == "translation" {
			ch = animation.NewPositionChannel(node)
			validTypes = []string{VEC3}
			validComponentTypes = []int{FLOAT}
		} else if target.Path == "rotation" {
			ch = animation.NewRotationChannel(node)
			validTypes = []string{VEC4}
			validComponentTypes = []int{FLOAT, BYTE, UNSIGNED_BYTE, SHORT, UNSIGNED_SHORT}
		} else if target.Path == "scale" {
			ch = animation.NewScaleChannel(node)
			validTypes = []string{VEC3}
			validComponentTypes = []int{FLOAT}
		} else if target.Path == "weights" {
			validTypes = []string{SCALAR}
			validComponentTypes = []int{FLOAT, BYTE, UNSIGNED_BYTE, SHORT, UNSIGNED_SHORT}
			return nil, fmt.Errorf("morph animation (with 'weights' path) not supported yet")
			// TODO
			//	for _, child := range node.GetNode().Children() {
			//		gr, ok := child.(graphic.Graphic)
			//		if ok {
			//			gr.geom
			//		}
			//	}
			//	ch = animation.NewMorphChannel(TODO)
		}

		// TODO what if Input and Output accessors are interleaved? probably de-interleave in these 2 cases

		keyframes, err := g.loadAccessorF32(sampler.Input, "Input", []string{SCALAR}, []int{FLOAT})
		if err != nil {
			return nil, err
		}
		values, err := g.loadAccessorF32(sampler.Output, "Output", validTypes, validComponentTypes)
		if err != nil {
			return nil, err
		}
		ch.SetBuffers(keyframes, values)
		ch.SetInterpolationType(animation.InterpolationType(sampler.Interpolation))
		anim.AddChannel(ch)
	}
	return anim, nil
}

// NewCamera creates and returns a Camera Node
// from the specified GLTF.Cameras index.
func (g *GLTF) NewCamera(ci int) (core.INode, error) {

	log.Debug("Loading Camera %d", ci)

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

// NewMesh creates and returns a Graphic Node (graphic.Mesh, graphic.Lines, graphic.Points, etc)
// from the specified GLTF.Meshes index.
func (g *GLTF) NewMesh(mi int) (core.INode, error) {

	log.Debug("Loading Mesh %d", mi)

	var err error
	meshData := g.Meshes[mi]

	// Create container node
	meshNode := core.NewNode()

	for i := 0; i < len(meshData.Primitives); i++ {

		// Get primitive information
		p := meshData.Primitives[i]

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
			grMat, err = g.NewMaterial(*p.Material)
			if err != nil {
				return nil, err
			}
		} else {
			grMat = g.newDefaultMaterial()
		}

		// Create geometry
		geom := geometry.NewGeometry()

		// Indices of buffer views
		interleavedVBOs := make(map[int]*gls.VBO, 0)

		// Load primitive attributes
		for name, aci := range p.Attributes {
			accessor := g.Accessors[aci]

			// Validate that accessor is compatible with attribute
			err = g.validateAccessorAttribute(accessor, name)
			if err != nil {
				return nil, err
			}

			// Load data and add it to geometry's VBO
			if g.isInterleaved(accessor) {
				bvIdx := *accessor.BufferView
				// Check if we already loaded this buffer view
				vbo, ok := interleavedVBOs[bvIdx]
				if ok {
					// Already created VBO for this buffer view
					// Add attribute with correct byteOffset
					g.addAttributeToVBO(vbo, name, uint32(*accessor.ByteOffset))
				} else {
					// Load data and create vbo
					buf, err := g.loadBufferView(g.BufferViews[bvIdx])
					if err != nil {
						return nil, err
					}
					data := g.bytesToArrayF32(buf, g.BufferViews[bvIdx].ByteLength)
					vbo := gls.NewVBO(data)
					g.addAttributeToVBO(vbo, name, 0)
					// Save reference to VBO keyed by index of the buffer view
					interleavedVBOs[bvIdx] = vbo
					// Add VBO to geometry
					geom.AddVBO(vbo)
				}
			} else {
				buf, err := g.loadAccessorBytes(accessor)
				if err != nil {
					return nil, err
				}
				data := g.bytesToArrayF32(buf, accessor.Count*TypeSizes[accessor.Type])
				vbo := gls.NewVBO(data)
				g.addAttributeToVBO(vbo, name, 0)
				// Add VBO to geometry
				geom.AddVBO(vbo)
			}
		}

		// Creates Geometry and add attribute VBO
		if len(indices) > 0 {
			geom.SetIndices(indices)
		}

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
		} else if mode == LINES {
			meshNode.Add(graphic.NewLines(geom, grMat))
		} else if mode == LINE_STRIP {
			meshNode.Add(graphic.NewLineStrip(geom, grMat))
		} else if mode == POINTS {
			meshNode.Add(graphic.NewPoints(geom, grMat))
		} else {
			return nil, fmt.Errorf("Unsupported primitive:%v", mode)
		}
	}

	return meshNode, nil
}

// loadIndices loads the indices stored in the specified accessor.
func (g *GLTF) loadIndices(ai int) (math32.ArrayU32, error) {

	return g.loadAccessorU32(ai, "indices", []string{SCALAR}, []int{UNSIGNED_BYTE, UNSIGNED_SHORT, UNSIGNED_INT}) // TODO check that it's ELEMENT_ARRAY_BUFFER
}

// addAttributeToVBO adds the appropriate attribute to the provided vbo based on the glTF attribute name.
func (g *GLTF) addAttributeToVBO(vbo *gls.VBO, attribName string, byteOffset uint32) {

	if attribName == "POSITION" {
		vbo.AddAttribOffset(gls.VertexPosition, byteOffset)
	} else if attribName == "NORMAL" {
		vbo.AddAttribOffset(gls.VertexNormal, byteOffset)
	} else if attribName == "TANGENT" {
		vbo.AddAttribOffset(gls.VertexTangent, byteOffset)
	} else if attribName == "TEXCOORD_0" {
		vbo.AddAttribOffset(gls.VertexTexcoord, byteOffset)
	} else if attribName == "COLOR_0" {	// TODO glTF spec says COLOR can be VEC3 or VEC4
		vbo.AddAttribOffset(gls.VertexColor, byteOffset)
	} else if attribName == "JOINTS_0" {
		// TODO
	} else if attribName == "WEIGHTS_0" {
		// TODO
	} else {
		panic(fmt.Sprintf("Attribute %v is not supported!", attribName))
	}
}

// validateAccessorAttribute validates the specified accessor for the given attribute name.
func (g *GLTF) validateAccessorAttribute(ac Accessor, attribName string) error {

	parts := strings.Split(attribName, "_")
	semantic := parts[0]

	usage := "attribute " + attribName

	if attribName == "POSITION" {
		return g.validateAccessor(ac, usage, []string{VEC3}, []int{FLOAT})
	} else if attribName == "NORMAL" {
		return g.validateAccessor(ac, usage, []string{VEC3}, []int{FLOAT})
	} else if attribName == "TANGENT" {
		return g.validateAccessor(ac, usage, []string{VEC4}, []int{FLOAT})
	} else if semantic == "TEXCOORD" {
		return g.validateAccessor(ac, usage, []string{VEC2}, []int{FLOAT, UNSIGNED_BYTE, UNSIGNED_SHORT})
	} else if semantic == "COLOR" {
		return g.validateAccessor(ac, usage, []string{VEC3, VEC4}, []int{FLOAT, UNSIGNED_BYTE, UNSIGNED_SHORT})
	} else if semantic == "JOINTS" {
		return g.validateAccessor(ac, usage, []string{VEC4}, []int{UNSIGNED_BYTE, UNSIGNED_SHORT})
	} else if semantic == "WEIGHTS" {
		return g.validateAccessor(ac, usage, []string{VEC4}, []int{FLOAT, UNSIGNED_BYTE, UNSIGNED_SHORT})
	} else {
		return fmt.Errorf("attribute %v is not supported", attribName)
	}
}

// validateAccessor validates the specified attribute accessor with the specified allowed types and component types.
func (g *GLTF) validateAccessor(ac Accessor, usage string, validTypes []string, validComponentTypes []int) error {

	// Validate accessor type
	validType := false
	for _, vType := range validTypes {
		if ac.Type == vType {
			validType = true
			break
		}
	}
	if !validType {
		return fmt.Errorf("invalid Accessor.Type %v for %s", ac.Type, usage)
	}

	// Validate accessor component type
	validComponentType := false
	for _, vComponentType := range validComponentTypes {
		if ac.ComponentType == vComponentType {
			validComponentType = true
			break
		}
	}
	if !validComponentType {
		return fmt.Errorf("invalid Accessor.ComponentType %v for %s", ac.ComponentType, usage)
	}

	return nil
}

// newDefaultMaterial creates and returns the default material.
func (g *GLTF) newDefaultMaterial() material.IMaterial {

	return material.NewStandard(&math32.Color{0.5, 0.5, 0.5})
}

// NewMaterial creates and returns a new material based on the material data with the specified index.
func (g *GLTF) NewMaterial(mi int) (material.IMaterial, error) {

	log.Debug("Loading Material %d", mi)

	matData := g.Materials[mi]
	// Check for material extensions
	if matData.Extensions != nil {
		for ext, v := range matData.Extensions {
			if ext == "KHR_materials_common" {
				return g.loadMaterialCommon(v)
			} else {
				return nil, fmt.Errorf("unsupported extension:%s", ext)
			}
		}
		return nil, fmt.Errorf("empty material extensions")
	} else {
		// Material is normally PBR
		return g.loadMaterialPBR(&matData)
	}
}

// NewTexture loads the texture specified by its index.
func (g *GLTF) NewTexture(texi int) (*texture.Texture2D, error) {

	log.Debug("Loading Texture %d", texi)

	// loads texture image
	texDesc := g.Textures[texi]
	img, err := g.NewImage(texDesc.Source)
	if err != nil {
		return nil, err
	}
	tex := texture.NewTexture2DFromRGBA(img)

	// Get sampler and apply texture parameters
	if texDesc.Sampler != nil {
		sampler := g.Samplers[*texDesc.Sampler]
		g.applySampler(sampler, tex)
	}

	return tex, nil
}

// applySamplers applies the specified Sampler to the provided texture.
func (g *GLTF) applySampler(sampler Sampler, tex *texture.Texture2D) {

	// Magnification filter
	magFilter := gls.LINEAR
	if sampler.MagFilter != nil {
		magFilter = *sampler.MagFilter
	}
	tex.SetMagFilter(uint32(magFilter))

	// Minification filter
	minFilter := gls.LINEAR_MIPMAP_LINEAR
	if sampler.MinFilter != nil {
		minFilter = *sampler.MinFilter
	}
	tex.SetMinFilter(uint32(minFilter))

	// S coordinate wrapping mode
	wrapS := gls.REPEAT
	if sampler.WrapS != nil {
		wrapS = *sampler.WrapS
	}
	tex.SetWrapS(uint32(wrapS))

	// T coordinate wrapping mode
	wrapT := gls.REPEAT
	if sampler.WrapT != nil {
		wrapT = *sampler.WrapT
	}
	tex.SetWrapT(uint32(wrapT))
}

// NewImage loads the image specified by the index of GLTF.Images.
// Image can be loaded from binary chunk file or data URI or external file..
func (g *GLTF) NewImage(ii int) (*image.RGBA, error) {

	log.Debug("Loading Image %d", ii)

	imgDesc := g.Images[ii]
	var data []byte
	var err error
	// If Uri is empty, load image from GLB binary chunk
	if imgDesc.Uri == "" {
		bvi := imgDesc.BufferView
		if bvi == nil {
			return nil, fmt.Errorf("image has empty URI and no BufferView")
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

// bytesToArrayU32 converts a byte array to ArrayU32.
func (g *GLTF) bytesToArrayU32(data []byte, componentType, count int) (math32.ArrayU32, error) {

	// If component is UNSIGNED_INT nothing to do
	if componentType == UNSIGNED_INT {
		arr := (*[1 << 30]uint32)(unsafe.Pointer(&data[0]))[:count]
		return math32.ArrayU32(arr), nil
	}

	// Converts UNSIGNED_SHORT to UNSIGNED_INT
	if componentType == UNSIGNED_SHORT {
		out := math32.NewArrayU32(count, count)
		for i := 0; i < count; i++ {
			out[i] = uint32(data[i*2]) + uint32(data[i*2+1])*256
		}
		return out, nil
	}

	// Converts UNSIGNED_BYTE indices to UNSIGNED_INT
	if componentType == UNSIGNED_BYTE {
		out := math32.NewArrayU32(count, count)
		for i := 0; i < count; i++ {
			out[i] = uint32(data[i])
		}
		return out, nil
	}

	return nil, fmt.Errorf("Unsupported Accessor ComponentType:%v", componentType)
}

// bytesToArrayF32 converts a byte array to ArrayF32.
func (g *GLTF) bytesToArrayF32(data []byte, size int) math32.ArrayF32 {

	return (*[1 << 30]float32)(unsafe.Pointer(&data[0]))[:size]
}

// loadAccessorU32 loads data from the specified accessor and performs validation of the Type and ComponentType.
func (g *GLTF) loadAccessorU32(ai int, usage string, validTypes []string, validComponentTypes []int) (math32.ArrayU32, error) {

	// Get Accessor for the specified index
	ac := g.Accessors[ai]
	if ac.BufferView == nil {
		return nil, fmt.Errorf("Accessor.BufferView == nil NOT SUPPORTED YET") // TODO
	}

	// Validate type and component type
	err := g.validateAccessor(ac, usage, validTypes, validComponentTypes)
	if err != nil {
		return nil, err
	}

	// Load bytes
	data, err := g.loadAccessorBytes(ac)
	if err != nil {
		return nil, err
	}

	return g.bytesToArrayU32(data, ac.ComponentType, ac.Count)
}

// loadAccessorF32 loads data from the specified accessor and performs validation of the Type and ComponentType.
func (g *GLTF) loadAccessorF32(ai int, usage string, validTypes []string, validComponentTypes []int) (math32.ArrayF32, error) {

	// Get Accessor for the specified index
	ac := g.Accessors[ai]
	if ac.BufferView == nil {
		return nil, fmt.Errorf("Accessor.BufferView == nil NOT SUPPORTED")
	}

	// Validate type and component type
	err := g.validateAccessor(ac, usage, validTypes, validComponentTypes)
	if err != nil {
		return nil, err
	}

	// Load bytes
	data, err := g.loadAccessorBytes(ac)
	if err != nil {
		return nil, err
	}

	return g.bytesToArrayF32(data, ac.Count*TypeSizes[ac.Type]), nil
}

// loadAccessorBytes returns the base byte array used by an accessor.
func (g *GLTF) loadAccessorBytes(ac Accessor) ([]byte, error) {

	// Get the Accessor's BufferView
	if ac.BufferView == nil {
		return nil, fmt.Errorf("Accessor has nil BufferView") // TODO
	}
	bv := g.BufferViews[*ac.BufferView]

	// Loads data from associated BufferView
	data, err := g.loadBufferView(bv)
	if err != nil {
		return nil, err
	}

	// Accessor offset into BufferView
	offset := 0
	if ac.ByteOffset != nil {
		offset = *ac.ByteOffset
	}
	data = data[offset:]

	// Check if interleaved and de-interleave if necessary ? TODO

	// Calculate the size in bytes of a complete attribute
	itemSize := TypeSizes[ac.Type]
	itemBytes := int(gls.FloatSize) * itemSize

	// If the BufferView stride is equal to the item size, the buffer is not interleaved
	if (bv.ByteStride != nil) && (*bv.ByteStride != itemBytes) {
		// BufferView data is interleaved, de-interleave
		// TODO
		return nil, fmt.Errorf("data is interleaved - not supported for animation yet")
	}

	return data, nil
}

// isInterleaves returns whether the BufferView used by the provided accessor is interleaved.
func (g *GLTF) isInterleaved(accessor Accessor) bool {

	// Get the Accessor's BufferView
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

// loadBufferView loads and returns a byte slice with data from the specified BufferView.
func (g *GLTF) loadBufferView(bv BufferView) ([]byte, error) {

	// Load buffer view buffer
	buf, err := g.loadBuffer(bv.Buffer)
	if err != nil {
		return nil, err
	}

	// Establish offset
	offset := 0
	if bv.ByteOffset != nil {
		offset = *bv.ByteOffset
	}

	// Compute and return offset slice
	return buf[offset : offset+bv.ByteLength], nil
}

// loadBuffer loads and returns the data from the specified GLTF Buffer index
func (g *GLTF) loadBuffer(bi int) ([]byte, error) {

	buf := &g.Buffers[bi]
	// If Buffer URI use the chunk data field
	if buf.Uri == "" {
		return g.data, nil
	}

	// If buffer already loaded:
	log.Debug("loadBuffer cache:%v", len(buf.data))
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
		log.Debug("loadBuffer: loading file")
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
	log.Debug("cache data:%v", len(buf.data))
	return data, nil
}

// dataURL describes a decoded data url string.
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

// isDataURL checks if the specified string has the prefix of data URL.
func isDataURL(url string) bool {

	if strings.HasPrefix(url, dataURLprefix) {
		return true
	}
	return false
}

// loadDataURL decodes the specified data URI string (base64).
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
		return nil, fmt.Errorf("data URI media type:%s not supported", du.MediaType)
	}

	// Checks encoding
	if du.Encoding != "base64" {
		return nil, fmt.Errorf("data URI encoding:%s not supported", du.Encoding)
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

	// Check prefix
	if !isDataURL(url) {
		return fmt.Errorf("specified string is not a data URL")
	}

	// Separate header from data
	body := url[len(dataURLprefix):]
	parts := strings.Split(body, ",")
	if len(parts) != 2 {
		return fmt.Errorf("data URI contains more than one ','")
	}
	du.Data = parts[1]

	// Separate media type from optional encoding
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
