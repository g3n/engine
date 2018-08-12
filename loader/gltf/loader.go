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

	// TODO Check for extensions used and extensions required

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

// LoadScene creates a parent Node which contains all nodes contained by
// the specified scene index from the GLTF Scenes array.
func (g *GLTF) LoadScene(sceneIdx int) (core.INode, error) {

	// Check if provided scene index is valid
	if sceneIdx < 0 || sceneIdx >= len(g.Scenes) {
		return nil, fmt.Errorf("invalid scene index")
	}
	log.Debug("Loading Scene %d", sceneIdx)
	sceneData := g.Scenes[sceneIdx]

	scene := core.NewNode()
	scene.SetName(sceneData.Name)

	// Load all nodes
	for _, ni := range sceneData.Nodes {
		child, err := g.LoadNode(ni)
		if err != nil {
			return nil, err
		}
		scene.Add(child)
	}
	return scene, nil
}

// LoadNode creates and returns a new Node described by the specified index
// in the decoded GLTF Nodes array.
func (g *GLTF) LoadNode(nodeIdx int) (core.INode, error) {

	// Check if provided node index is valid
	if nodeIdx < 0 || nodeIdx >= len(g.Nodes) {
		return nil, fmt.Errorf("invalid node index")
	}
	nodeData := g.Nodes[nodeIdx]
	// Return cached if available
	if nodeData.cache != nil {
		log.Debug("Fetching Node %d (cached)", nodeIdx)
		return nodeData.cache, nil
	}
	log.Debug("Loading Node %d", nodeIdx)

	var in core.INode
	var err error
	// Check if the node is a Mesh (triangles, lines, etc...)
	if nodeData.Mesh != nil {
		in, err = g.LoadMesh(*nodeData.Mesh)
		if err != nil {
			return nil, err
		}

		if nodeData.Skin != nil {
			children := in.GetNode().Children()
			if len(children) > 1 {
				//log.Error("skinning/rigging meshes with more than a single primitive is not supported")
				return nil, fmt.Errorf("skinning/rigging meshes with more than a single primitive is not supported")
			}
			mesh := children[0].(*graphic.Mesh)
			// Create RiggedMesh
			rm := graphic.NewRiggedMesh(mesh)
			skeleton, err := g.LoadSkin(*nodeData.Skin)
			if err != nil {
				return nil, err
			}
			rm.SetSkeleton(skeleton)
			in = rm
		}

		// Check if the node is Camera
	} else if nodeData.Camera != nil {
		in, err = g.LoadCamera(*nodeData.Camera)
		if err != nil {
			return nil, err
		}
		// Other cases, return empty node
	} else {
		log.Debug("Empty Node")
		in = core.NewNode()
	}

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

	// Cache node
	g.Nodes[nodeIdx].cache = in

	// Recursively load node children and add them to the parent
	for _, ci := range nodeData.Children {
		child, err := g.LoadNode(ci)
		if err != nil {
			return nil, err
		}
		node.Add(child)
	}

	return in, nil
}

// LoadSkin loads the skin with specified index.
func (g *GLTF) LoadSkin(skinIdx int) (*graphic.Skeleton, error) {

	// Check if provided skin index is valid
	if skinIdx < 0 || skinIdx >= len(g.Skins) {
		return nil, fmt.Errorf("invalid skin index")
	}
	skinData := g.Skins[skinIdx]
	// Return cached if available
	if skinData.cache != nil {
		log.Debug("Fetching Skin %d (cached)", skinIdx)
		return skinData.cache, nil
	}
	log.Debug("Loading Skin %d", skinIdx)

	// Create Skeleton and set it on Rigged mesh
	skeleton := graphic.NewSkeleton()

	// Load inverseBindMatrices
	ibmData, err := g.loadAccessorF32(skinData.InverseBindMatrices, "ibm", []string{MAT4}, []int{FLOAT})
	if err != nil {
		return nil, err
	}

	// Add bones
	for i := range skinData.Joints {
		jointNode, err := g.LoadNode(skinData.Joints[i])
		if err != nil {
			return nil, err
		}
		var ibm math32.Matrix4
		ibmData.GetMatrix4(16 * i, &ibm)
		skeleton.AddBone(jointNode.GetNode(), &ibm)
	}

	// Cache skin
	g.Skins[skinIdx].cache = skeleton

	return skeleton, nil
}

// LoadAnimationByName loads the animations with specified name.
// If there are multiple animations with the same name it loads the first occurrence.
func (g *GLTF) LoadAnimationByName(animName string) (*animation.Animation, error) {

	for i := range g.Animations {
		if g.Animations[i].Name == animName {
			return g.LoadAnimation(i)
		}
	}
	return nil, fmt.Errorf("could not find animation named %v", animName)
}

// LoadAnimation creates an Animation for the specified
// animation index from the GLTF Animations array.
func (g *GLTF) LoadAnimation(animIdx int) (*animation.Animation, error) {

	// Check if provided animation index is valid
	if animIdx < 0 || animIdx >= len(g.Animations) {
		return nil, fmt.Errorf("invalid animation index")
	}
	log.Debug("Loading Animation %d", animIdx)
	animData := g.Animations[animIdx]

	anim := animation.NewAnimation()
	anim.SetName(animData.Name)
	for i := 0; i < len(animData.Channels); i++ {

		chData := animData.Channels[i]
		target := chData.Target
		sampler := animData.Samplers[chData.Sampler]
		node, err := g.LoadNode(target.Node)
		if err != nil {
			return nil, err
		}

		var validTypes []string
		var validComponentTypes []int

		var ch animation.IChannel
		if target.Path == "translation" {
			validTypes = []string{VEC3}
			validComponentTypes = []int{FLOAT}
			ch = animation.NewPositionChannel(node)
		} else if target.Path == "rotation" {
			validTypes = []string{VEC4}
			validComponentTypes = []int{FLOAT, BYTE, UNSIGNED_BYTE, SHORT, UNSIGNED_SHORT}
			ch = animation.NewRotationChannel(node)
		} else if target.Path == "scale" {
			validTypes = []string{VEC3}
			validComponentTypes = []int{FLOAT}
			ch = animation.NewScaleChannel(node)
		} else if target.Path == "weights" {
			validTypes = []string{SCALAR}
			validComponentTypes = []int{FLOAT, BYTE, UNSIGNED_BYTE, SHORT, UNSIGNED_SHORT}
			children := node.GetNode().Children()
			if len(children) > 1 {
				return nil, fmt.Errorf("animating meshes with more than a single primitive is not supported")
			}
			morphGeom := children[0].(graphic.IGraphic).IGeometry().(*geometry.MorphGeometry)
			ch = animation.NewMorphChannel(morphGeom)
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

// LoadCamera creates and returns a Camera Node
// from the specified GLTF.Cameras index.
func (g *GLTF) LoadCamera(camIdx int) (core.INode, error) {

	// Check if provided camera index is valid
	if camIdx < 0 || camIdx >= len(g.Cameras) {
		return nil, fmt.Errorf("invalid camera index")
	}
	log.Debug("Loading Camera %d", camIdx)
	camData := g.Cameras[camIdx]

	if camData.Type == "perspective" {
		desc := camData.Perspective
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

	if camData.Type == "orthographic" {
		desc := camData.Orthographic
		cam := camera.NewOrthographic(desc.Xmag/-2, desc.Xmag/2, desc.Ymag/2, desc.Ymag/-2, desc.Znear, desc.Zfar)
		return cam, nil

	}

	return nil, fmt.Errorf("unsupported camera type: %s", camData.Type)
}

// LoadMesh creates and returns a Graphic Node (graphic.Mesh, graphic.Lines, graphic.Points, etc)
// from the specified GLTF.Meshes index.
func (g *GLTF) LoadMesh(meshIdx int) (core.INode, error) {

	// Check if provided mesh index is valid
	if meshIdx < 0 || meshIdx >= len(g.Meshes) {
		return nil, fmt.Errorf("invalid mesh index")
	}
	meshData := g.Meshes[meshIdx]
	// Return cached if available
	if meshData.cache != nil {
		// TODO CLONE/REINSTANCE INSTEAD
		//log.Debug("Instancing Mesh %d (from cached)", meshIdx)
		//return meshData.cache, nil
	}
	log.Debug("Loading Mesh %d", meshIdx)

	var err error

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
			grMat, err = g.LoadMaterial(*p.Material)
			if err != nil {
				return nil, err
			}
		} else {
			grMat = g.newDefaultMaterial()
		}

		// Create geometry
		var igeom geometry.IGeometry
		igeom = geometry.NewGeometry()
		geom := igeom.GetGeometry()

		err = g.loadAttributes(geom, p.Attributes, indices)
		if err != nil {
			return nil, err
		}

		// If primitive has targets then the geometry should be a morph geometry
		if len(p.Targets) > 0 {
			morphGeom := geometry.NewMorphGeometry(geom)

			// TODO Load morph target names if present in extras under "targetNames"
			// TODO Update morph target weights if present in Mesh.Weights

			// Load targets
			for i := range p.Targets {
				tGeom := geometry.NewGeometry()
				attributes := p.Targets[i]
				err = g.loadAttributes(tGeom, attributes, indices)
				if err != nil {
					return nil, err
				}
				morphGeom.AddMorphTargetDeltas(tGeom)
			}

			igeom = morphGeom
		}

		// Default mode is 4 (TRIANGLES)
		mode := TRIANGLES
		if p.Mode != nil {
			mode = *p.Mode
		}

		// Create Mesh
		// TODO materials for LINES, etc need to be different...
		if mode == TRIANGLES {
			meshNode.Add(graphic.NewMesh(igeom, grMat))
		} else if mode == LINES {
			meshNode.Add(graphic.NewLines(igeom, grMat))
		} else if mode == LINE_STRIP {
			meshNode.Add(graphic.NewLineStrip(igeom, grMat))
		} else if mode == POINTS {
			meshNode.Add(graphic.NewPoints(igeom, grMat))
		} else {
			return nil, fmt.Errorf("unsupported primitive:%v", mode)
		}
	}

	// Cache mesh
	g.Meshes[meshIdx].cache = meshNode

	return meshNode, nil
}

// loadAttributes loads the provided list of vertex attributes as VBO(s) into the specified geometry.
func (g *GLTF) loadAttributes(geom *geometry.Geometry, attributes map[string]int, indices math32.ArrayU32) error {

	// Indices of buffer views
	interleavedVBOs := make(map[int]*gls.VBO, 0)

	// Load primitive attributes
	for name, aci := range attributes {
		accessor := g.Accessors[aci]

		// Validate that accessor is compatible with attribute
		err := g.validateAccessorAttribute(accessor, name)
		if err != nil {
			return err
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
				buf, err := g.loadBufferView(bvIdx)
				if err != nil {
					return err
				}
				data, err := g.bytesToArrayF32(buf, accessor.ComponentType, accessor.Count*TypeSizes[accessor.Type])
				if err != nil {
					return err
				}
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
				return err
			}
			data, err := g.bytesToArrayF32(buf, accessor.ComponentType, accessor.Count*TypeSizes[accessor.Type])
			if err != nil {
				return err
			}
			vbo := gls.NewVBO(data)
			g.addAttributeToVBO(vbo, name, 0)
			// Add VBO to geometry
			geom.AddVBO(vbo)
		}
	}

	// Set indices
	if len(indices) > 0 {
		geom.SetIndices(indices)
	}

	return nil
}

// loadIndices loads the indices stored in the specified accessor.
func (g *GLTF) loadIndices(ai int) (math32.ArrayU32, error) {

	return g.loadAccessorU32(ai, "indices", []string{SCALAR}, []int{UNSIGNED_BYTE, UNSIGNED_SHORT, UNSIGNED_INT}) // TODO verify that it's ELEMENT_ARRAY_BUFFER
}

// addAttributeToVBO adds the appropriate attribute to the provided vbo based on the glTF attribute name.
func (g *GLTF) addAttributeToVBO(vbo *gls.VBO, attribName string, byteOffset uint32) {

	aType, ok := AttributeName[attribName]
	if !ok {
		log.Warn(fmt.Sprintf("Attribute %v is not supported!", attribName))
		return
	}
	vbo.AddAttribOffset(aType, byteOffset)
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
		// Note that morph targets only support VEC3 whereas normal attributes only support VEC4.
		return g.validateAccessor(ac, usage, []string{VEC3, VEC4}, []int{FLOAT})
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

// LoadMaterial creates and returns a new material based on the material data with the specified index.
func (g *GLTF) LoadMaterial(matIdx int) (material.IMaterial, error) {

	// Check if provided material index is valid
	if matIdx < 0 || matIdx >= len(g.Materials) {
		return nil, fmt.Errorf("invalid material index")
	}
	matData := g.Materials[matIdx]
	// Return cached if available
	if matData.cache != nil {
		log.Debug("Fetching Material %d (cached)", matIdx)
		return matData.cache, nil
	}
	log.Debug("Loading Material %d", matIdx)

	var err error
	var imat material.IMaterial

	// Check for material extensions
	if matData.Extensions != nil {
		for ext, extData := range matData.Extensions {
			if ext == KhrMaterialsCommon {
				imat, err = g.loadMaterialCommon(extData)
			} else if ext == KhrMaterialsUnlit {
				//imat, err = g.loadMaterialUnlit(matData, extData)
			//} else if ext == KhrMaterialsPbrSpecularGlossiness {
			} else {
				return nil, fmt.Errorf("unsupported extension:%s", ext)
			}
		}
	} else {
		// Material is normally PBR
		imat, err = g.loadMaterialPBR(&matData)
	}

	// Cache material
	g.Materials[matIdx].cache = imat

	return imat, err
}

// LoadTexture loads the texture specified by its index.
func (g *GLTF) LoadTexture(texIdx int) (*texture.Texture2D, error) {

	// Check if provided texture index is valid
	if texIdx < 0 || texIdx >= len(g.Textures) {
		return nil, fmt.Errorf("invalid texture index")
	}
	texData := g.Textures[texIdx]
	// NOTE: Textures can't be cached because they have their own uniforms
	log.Debug("Loading Texture %d", texIdx)

	// Load texture image
	img, err := g.LoadImage(texData.Source)
	if err != nil {
		return nil, err
	}
	tex := texture.NewTexture2DFromRGBA(img)

	// Get sampler and apply texture parameters
	if texData.Sampler != nil {
		err = g.applySampler(*texData.Sampler, tex)
		if err != nil {
			return nil, err
		}
	}

	return tex, nil
}

// applySamplers applies the specified Sampler to the provided texture.
func (g *GLTF) applySampler(samplerIdx int, tex *texture.Texture2D) error {

	log.Debug("Applying Sampler %d", samplerIdx)
	// Check if provided sampler index is valid
	if samplerIdx < 0 || samplerIdx >= len(g.Samplers) {
		return fmt.Errorf("invalid sampler index")
	}
	sampler := g.Samplers[samplerIdx]

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

	return nil
}

// LoadImage loads the image specified by the index of GLTF.Images.
// Image can be loaded from binary chunk file or data URI or external file..
func (g *GLTF) LoadImage(imgIdx int) (*image.RGBA, error) {

	// Check if provided image index is valid
	if imgIdx < 0 || imgIdx >= len(g.Images) {
		return nil, fmt.Errorf("invalid image index")
	}
	imgData := g.Images[imgIdx]
	// Return cached if available
	if imgData.cache != nil {
		log.Debug("Fetching Image %d (cached)", imgIdx)
		return imgData.cache, nil
	}
	log.Debug("Loading Image %d", imgIdx)

	var data []byte
	var err error
	// If Uri is empty, load image from GLB binary chunk
	if imgData.Uri == "" {
		if imgData.BufferView == nil {
			return nil, fmt.Errorf("image has empty URI and no BufferView")
		}
		data, err = g.loadBufferView(*imgData.BufferView)
	} else if isDataURL(imgData.Uri) {
		// Checks if image URI is data URL
		data, err = loadDataURL(imgData.Uri)
	} else {
		// Load image data from file
		data, err = g.loadFileBytes(imgData.Uri)
	}

	if err != nil {
		return nil, err
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

	// Cache image
	g.Images[imgIdx].cache = rgba

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

	return nil, fmt.Errorf("unsupported Accessor ComponentType:%v", componentType)
}

// bytesToArrayF32 converts a byte array to ArrayF32.
func (g *GLTF) bytesToArrayF32(data []byte, componentType, count int) (math32.ArrayF32, error) {

	// If component is UNSIGNED_INT nothing to do
	if componentType == UNSIGNED_INT {
		arr := (*[1 << 30]float32)(unsafe.Pointer(&data[0]))[:count]
		return math32.ArrayF32(arr), nil
	}

	// Converts UNSIGNED_SHORT to UNSIGNED_INT
	if componentType == UNSIGNED_SHORT {
		out := math32.NewArrayF32(count, count)
		for i := 0; i < count; i++ {
			out[i] = float32(data[i*2]) + float32(data[i*2+1])*256
		}
		return out, nil
	}

	// Converts UNSIGNED_BYTE indices to UNSIGNED_INT
	if componentType == UNSIGNED_BYTE {
		out := math32.NewArrayF32(count, count)
		for i := 0; i < count; i++ {
			out[i] = float32(data[i])
		}
		return out, nil
	}

	return (*[1 << 30]float32)(unsafe.Pointer(&data[0]))[:count], nil
}

// loadAccessorU32 loads data from the specified accessor and performs validation of the Type and ComponentType.
func (g *GLTF) loadAccessorU32(ai int, usage string, validTypes []string, validComponentTypes []int) (math32.ArrayU32, error) {

	// Get Accessor for the specified index
	ac := g.Accessors[ai]
	if ac.BufferView == nil {
		return nil, fmt.Errorf("accessor.BufferView == nil NOT SUPPORTED YET") // TODO
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

	return g.bytesToArrayU32(data, ac.ComponentType, ac.Count*TypeSizes[ac.Type])
}

// loadAccessorF32 loads data from the specified accessor and performs validation of the Type and ComponentType.
func (g *GLTF) loadAccessorF32(ai int, usage string, validTypes []string, validComponentTypes []int) (math32.ArrayF32, error) {

	// Get Accessor for the specified index
	ac := g.Accessors[ai]
	if ac.BufferView == nil {
		return nil, fmt.Errorf("accessor.BufferView == nil NOT SUPPORTED YET") // TODO
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

	return g.bytesToArrayF32(data, ac.ComponentType, ac.Count*TypeSizes[ac.Type])
}

// loadAccessorBytes returns the base byte array used by an accessor.
func (g *GLTF) loadAccessorBytes(ac Accessor) ([]byte, error) {

	// Get the Accessor's BufferView
	if ac.BufferView == nil {
		return nil, fmt.Errorf("accessor.BufferView == nil NOT SUPPORTED YET") // TODO
	}
	bv := g.BufferViews[*ac.BufferView]

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

	// TODO check if interleaved and de-interleave if necessary?

	// Calculate the size in bytes of a complete attribute
	itemSize := TypeSizes[ac.Type]
	itemBytes := int(gls.FloatSize) * itemSize

	// If the BufferView stride is equal to the item size, the buffer is not interleaved
	if (bv.ByteStride != nil) && (*bv.ByteStride != itemBytes) {
		// BufferView data is interleaved, de-interleave
		// TODO
		return nil, fmt.Errorf("data is interleaved - not supported for animation yet")
	}

	// TODO Sparse accessor

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
func (g *GLTF) loadBufferView(bvIdx int) ([]byte, error) {

	// Check if provided buffer view index is valid
	if bvIdx < 0 || bvIdx >= len(g.BufferViews) {
		return nil, fmt.Errorf("invalid buffer view index")
	}
	bvData := g.BufferViews[bvIdx]
	// Return cached if available
	if bvData.cache != nil {
		log.Debug("Fetching BufferView %d (cached)", bvIdx)
		return bvData.cache, nil
	}
	log.Debug("Loading BufferView %d", bvIdx)

	// Load buffer view buffer
	buf, err := g.loadBuffer(bvData.Buffer)
	if err != nil {
		return nil, err
	}

	// Establish offset
	offset := 0
	if bvData.ByteOffset != nil {
		offset = *bvData.ByteOffset
	}

	// Compute and return offset slice
	bvBytes := buf[offset : offset+bvData.ByteLength]

	// Cache buffer view
	g.BufferViews[bvIdx].cache = bvBytes

	return bvBytes, nil
}

// loadBuffer loads and returns the data from the specified GLTF Buffer index
func (g *GLTF) loadBuffer(bufIdx int) ([]byte, error) {

	// Check if provided buffer index is valid
	if bufIdx < 0 || bufIdx >= len(g.Buffers) {
		return nil, fmt.Errorf("invalid buffer index")
	}
	bufData := &g.Buffers[bufIdx]
	// Return cached if available
	if bufData.cache != nil {
		log.Debug("Fetching Buffer %d (cached)", bufIdx)
		return bufData.cache, nil
	}
	log.Debug("Loading Buffer %d", bufIdx)

	// If buffer URI use the chunk data field
	if bufData.Uri == "" {
		return g.data, nil
	}

	// Checks if buffer URI is a data URI
	var data []byte
	var err error
	if isDataURL(bufData.Uri) {
		data, err = loadDataURL(bufData.Uri)
	} else {
		// Try to load buffer from file
		data, err = g.loadFileBytes(bufData.Uri)
	}
	if err != nil {
		return nil, err
	}

	// Checks data length
	if len(data) != bufData.ByteLength {
		return nil, fmt.Errorf("buffer:%d read data length:%d expected:%d", bufIdx, len(data), bufData.ByteLength)
	}
	// Cache buffer data
	g.Buffers[bufIdx].cache = data
	log.Debug("cache data:%v", len(bufData.cache))
	return data, nil
}

// loadFileBytes loads the file with specified path as a byte array.
func (g *GLTF) loadFileBytes(uri string) ([]byte, error) {

	log.Debug("Loading File: %v", uri)

	fpath := filepath.Join(g.path, uri)
	f, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
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
