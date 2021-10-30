// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package gltf
package gltf

import (
	"image"

	"github.com/g3n/engine/animation"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
)

// glTF Extensions.
const (
	KhrDracoMeshCompression           = "KHR_draco_mesh_compression"
	KhrMaterialsUnlit                 = "KHR_materials_unlit"
	KhrMaterialsCommon                = "KHR_materials_common" // TODO this is officially part of glTF 1.0 (remove?)
	KhrMaterialsPbrSpecularGlossiness = "KHR_materials_pbrSpecularGlossiness"
)

// GLTF is the root object for a glTF asset.
type GLTF struct {
	ExtensionsUsed     []string               // Names of glTF extensions used somewhere in this asset. Not required.
	ExtensionsRequired []string               // Names of glTF extensions required to properly load this asset. Not required.
	Accessors          []Accessor             // An array of accessors. Not required.
	Animations         []Animation            // An array of keyframe animations. Not required.
	Asset              Asset                  // Metadata about the glTF asset. Required.
	Buffers            []Buffer               // An array of buffers. Not required.
	BufferViews        []BufferView           // An array of bufferViews. Not required.
	Cameras            []Camera               // An array of cameras. Not required.
	Images             []Image                // An array of images. Not required.
	Materials          []Material             // An array of materials. Not required.
	Meshes             []Mesh                 // An array of meshes. Not required.
	Nodes              []Node                 // An array of nodes. Not required.
	Samplers           []Sampler              // An array of samplers. Not required.
	Scene              *int                   // The index of the default scene. Not required.
	Scenes             []Scene                // An array of scenes. Not required.
	Skins              []Skin                 // An array of skins. Not required.
	Textures           []Texture              // An array of textures. Not required.
	Extensions         map[string]interface{} // Dictionary object with extension-specific objects. Not required.
	Extras             interface{}            // Application-specific data. Not required.

	path string // File path for resources.
	data []byte // Binary file Chunk 1 data.
}

// Accessor is a typed view into a BufferView.
type Accessor struct {
	BufferView    *int                   // The index of the buffer view. Not required.
	ByteOffset    *int                   // The offset relative to the start of the BufferView in bytes. Not required. Default is 0.
	ComponentType int                    // The data type of components in the attribute. Required.
	Normalized    bool                   // Specifies whether integer data values should be normalized. Not required. Default is false.
	Count         int                    // The number of attributes referenced by this accessor. Required.
	Type          string                 // Specifies if the attribute is a scalar, vector or matrix. Required.
	Max           []float32              // Maximum value of each component in this attribute. Not required.
	Min           []float32              // Minimum value of each component in this attribute. Not required.
	Sparse        *Sparse                // Sparse storage attribute that deviates from their initialization value. Not required.
	Name          string                 // The user-defined name of this object. Not required.
	Extensions    map[string]interface{} // Dictionary object with extension specific objects. Not required.
	Extras        interface{}            // Application-specific data. Not required.

	cache math32.ArrayF32 // TODO implement caching
}

// Animation is a keyframe animation.
type Animation struct {
	Channels   []Channel              // An array of channels, each of which targets an animation's sampler at a node's property. Different channels of the same animation can't have equal targets. Required.
	Samplers   []AnimationSampler     // An array of samplers that combines input and output accessors with an interpolation algorithm to define a keyframe graph (but not its target). Required.
	Name       string                 // The user-defined name of this object. Not required.
	Extensions map[string]interface{} // Dictionary object with extension specific objects. Not required.
	Extras     interface{}            // Application-specific data. Not required.

	cache *animation.Animation // Cached Animation. // TODO
}

// AnimationSample combines input and output accessors with an interpolation algorithm to define a keyframe graph (but not its target).
type AnimationSampler struct {
	Input         int                    // The index of an accessor containing keyframe input values, e.g., time. Required.
	Interpolation string                 // Interpolation algorithm. Not required. Default is "LINEAR".
	Output        int                    // The index of an accessor, containing keyframe output values. Required.
	Extensions    map[string]interface{} // Dictionary object with extension-specific objects. Not required.
	Extras        interface{}            // Application-specific data. Not required.
}

// Asset contains metadata about the glTF asset.
type Asset struct {
	Copyright  string                 // A copyright message suitable for display to credit the content creator. Not required.
	Generator  string                 // Tool that generated this glTF model. Useful for debugging. Not required.
	Version    string                 // The glTF version that this asset targets. Required.
	MinVersion string                 // The minimum glTF version that this asset targets. Not required.
	Extensions map[string]interface{} // Dictionary object with extension-specific objects. Not required.
	Extras     interface{}            // Application-specific data. Not required.
}

// Buffer points to binary geometry, animation, or skins.
type Buffer struct {
	Uri        string                 // The URI of the buffer. Not required.
	ByteLength int                    // The length of the buffer in bytes. Required.
	Name       string                 // The user-defined name of this object. Not required.
	Extensions map[string]interface{} // Dictionary object with extension-specific objects. Not required.
	Extras     interface{}            // Application-specific data. Not required.

	cache []byte // Cached buffer data.
}

// BufferView is a view into a buffer generally representing a subset of the buffer.
type BufferView struct {
	Buffer     int                    // The index of the buffer. Required.
	ByteOffset *int                   // The offset into the buffer, in bytes. Not required. Default is 0.
	ByteLength int                    // The length of the buffer view, in bytes. Required.
	ByteStride *int                   // The stride, in bytes. Not required.
	Target     *int                   // The target that the GPU buffer should be bound to. Not required.
	Name       string                 // The user-defined name of this object. Not required.
	Extensions map[string]interface{} // Dictionary object with extension specific objects. Not required.
	Extras     interface{}            // Application-specific data. Not required.

	cache []byte // Cached buffer view data.
}

// Camera is a camera's projection.
// A node can reference a camera to apply a transform to place the camera in the scene.
type Camera struct {
	Orthographic *Orthographic          // An orthographic camera containing properties to create an orthographic projection matrix. Not required.
	Perspective  *Perspective           // A perspective camera containing properties to create a perspective projection matrix. Not required.
	Type         string                 // Specifies if the camera uses a perspective or orthographic projection. Required.
	Name         string                 // The user-defined name of this object. Not required.
	Extensions   map[string]interface{} // Dictionary object with extension-specific objects. Not required.
	Extras       interface{}            // Application-specific data. Not required.

	cache camera.ICamera // Cached ICamera. // TODO
}

// Channel targets an animation's sampler at a node's property.
type Channel struct {
	Sampler    int                    // The index of a sampler in this animation used to compute the value for the target. Required.
	Target     Target                 // The index of the node and TRS property to target. Required.
	Extensions map[string]interface{} // Dictionary object with extension-specific objects. Not required.
	Extras     interface{}            // Application-specific data. Not required.
}

// Image data used to create a texture.
// Image can be referenced by URI or bufferView index. mimeType is required in the latter case.
type Image struct {
	Uri        string                 // The URI of the image. Not required.
	MimeType   string                 // The image's MIME type. Not required.
	BufferView *int                   // The index of the bufferView that contains the image. Use this instead of the image's uri property. Not required.
	Name       string                 // The user-defined name of this object. Not required.
	Extensions map[string]interface{} // Dictionary object with extension-specific objects. Not required.
	Extras     interface{}            // Application-specific data. Not required.

	cache *image.RGBA // Cached image.
}

// Indices of those attributes that deviate from their initialization value.
type Indices struct {
	BufferView    int                    // The index of the bufferView with sparse indices. Referenced bufferView can't have ARRAY_BUFFER or ELEMENT_ARRAY_BUFFER target. Required.
	ByteOffset    int                    // The offset relative to the start of the bufferView in bytes. Must be aligned. Not required. Default is 0.
	ComponentType int                    // The indices data type. Required.
	Extensions    map[string]interface{} // Dictionary object with extension-specific objects. Not required.
	Extras        interface{}            // Application-specific data. Not required.
}

// Material describes the material appearance of a primitive.
type Material struct {
	Name                 string                 // The user-defined name of this object. Not required.
	PbrMetallicRoughness *PbrMetallicRoughness  // A set of parameter values that are used to define the metallic-roughness material model from Physically-Based Rendering (PBR) methodology. When not specified, all the default values of pbrMetallicRoughness apply. Not required.
	NormalTexture        *NormalTextureInfo     // The normal map texture. Not required.
	OcclusionTexture     *OcclusionTextureInfo  // The occlusion map texture. Not required.
	EmissiveTexture      *TextureInfo           // The emissive map texture. Not required.
	EmissiveFactor       *[3]float32            // The emissive color of the material. Not required. Default is [0,0,0]
	AlphaMode            string                 // The alpha rendering mode of the material. Not required. Default is OPAQUE.
	AlphaCutoff          float32                // The alpha cutoff value of the material. Not required. Default is 0.5.
	DoubleSided          bool                   // Specifies whether the material is double sided. Not required. Default is false.
	Extensions           map[string]interface{} // Dictionary object with extension-specific objects. Not required.
	Extras               interface{}            // Application-specific data. Not required.

	cache material.IMaterial // Cached IMaterial.
}

// Mesh is a set of primitives to be rendered.
// A node can contain one mesh. A node's transform places the mesh in the scene.
type Mesh struct {
	Primitives []Primitive            // An array of primitives, each defining geometry to be rendered with a material. Required.
	Weights    []float32              // Array of weights to be applied to the Morph Targets. Not required.
	Name       string                 // The user-defined name of this object. Not required.
	Extensions map[string]interface{} // Dictionary object with extension-specific objects. Not required.
	Extras     interface{}            // Application-specific data. Not required.

	cache core.INode // Cached INode. We don't cache an IGraphic here because a glTFL mesh can contain multiple primitive IGraphics.
}

// Node is a node in the node hierarchy.
// When the node contains skin, all mesh.primitives must contain JOINTS_0 and WEIGHTS_0 attributes.
// A node can have either a matrix or any combination of translation/rotation/scale (TRS) properties.
// TRS properties are converted to matrices and postmultiplied in the T * R * S order to compose the transformation matrix; first the scale is applied to the vertices, then the rotation, and then the translation.
// If none are provided, the transform is the identity.
// When a node is targeted for animation (referenced by an animation.channel.target), only TRS properties may be present; matrix will not be present.
type Node struct {
	Camera      *int                   // Index of the camera referenced by this node. Not required.
	Children    []int                  // The indices of this node's children. Not required.
	Skin        *int                   // The index of the skin referenced by this node. Not required.
	Matrix      *[16]float32           // Floating point 4x4 transformation matrix in column-major order. Not required. Default is the identity matrix.
	Mesh        *int                   // The index of the mesh in this node. Not required.
	Rotation    *[4]float32            // The node's unit quaternion rotation in the order (x, y, z, w), where w is the scalar. Not required. Default is [0,0,0,1].
	Scale       *[3]float32            // The node's non-uniform scale, given as the scaling factors along the x, y, and z axes. Not required. Default is [1,1,1].
	Translation *[3]float32            // The node's translation along the x, y, and z axes. Not required. Default is [0,0,0].
	Weights     []float32              // The weights of the instantiated Morph Target. Number of elements must match number of Morph Targets of used mesh. Not required.
	Name        string                 // The user-defined name of this object. Not required.
	Extensions  map[string]interface{} // Dictionary object with extension-specific objects. Not required.
	Extras      interface{}            // Application-specific data. Not required.

	cache core.INode // Cached INode.
}

// TODO Why not combine NormalTextureInfo and OcclusionTextureInfo ? Or simply add Scale to TextureInfo and use only TextureInfo?
// Propose this change to the official specification!

// NormalTextureInfo is a reference to a texture.
type NormalTextureInfo struct {
	Index      int                    // The index of the texture. Required.
	TexCoord   int                    // The set index of texture's TEXCOORD attribute used for texture coordinate mapping. Not required. Default is 0.
	Scale      float32                // The scalar multiplier applied to each normal vector of the normal texture. Not required. Default is 1.
	Extensions map[string]interface{} // Dictionary object with extension-specific objects. Not required.
	Extras     interface{}            // Application-specific data. Not required.
}

// OcclusionTextureInfo is a reference to a texture.
type OcclusionTextureInfo struct {
	Index      int                    // The index of the texture. Required.
	TexCoord   int                    // The set index of texture's TEXCOORD attribute used for texture coordinate mapping. Not required. Default is 0.
	Strength   float32                // The scalar multiplier controlling the amount of occlusion applied. Not required. Default is 1.
	Extensions map[string]interface{} // Dictionary object with extension-specific objects. Not required.
	Extras     interface{}            // Application-specific data. Not required.
}

// Orthographic is an orthographic camera containing properties to create an orthographic projection matrix.
type Orthographic struct {
	Xmag       float32                // The floating-point horizontal magnification of the view. Required.
	Ymag       float32                // The floating-point vertical magnification of the view. Required.
	Zfar       float32                // The floating-point distance to the far clipping plane. Zfar must be greater than Znear. Required.
	Znear      float32                // The floating-point distance to the near clipping plane. Required.
	Extensions map[string]interface{} // Dictionary object with extension-specific objects. Not required.
	Extras     interface{}            // Application-specific data. Not required.
}

// PbrMetallicRoughness is a set of parameter values that are used to define the metallic-roughness material model from Physically-Based Rendering (PBR) methodology.
type PbrMetallicRoughness struct {
	BaseColorFactor          *[4]float32            // The material's base color factor. Not required. Default is [1,1,1,1]
	BaseColorTexture         *TextureInfo           // The base color texture. Not required.
	MetallicFactor           *float32               // The metalness of the material. Not required. Default is 1.
	RoughnessFactor          *float32               // The roughness of the material. Not required. Default is 1.
	MetallicRoughnessTexture *TextureInfo           // The metallic-roughness texture. Not required.
	Extensions               map[string]interface{} // Dictionary object with extension-specific objects. Not required.
	Extras                   interface{}            // Application-specific data. Not required.
}

// Perspective is a perspective camera containing properties to create a perspective projection matrix.
type Perspective struct {
	AspectRatio *float32               // The floating-point aspect ratio of the field of view. Not required.
	Yfov        float32                // The floating-point vertical field of view in radians. Required.
	Zfar        *float32               // The floating-point distance to the far clipping plane. Not required.
	Znear       float32                // The floating-point distance to the near clipping plane. Required.
	Extensions  map[string]interface{} // Dictionary object with extension-specific objects. Not required.
	Extras      interface{}            // Application-specific data. Not required.
}

// Primitive represents geometry to be rendered with the given material.
type Primitive struct {
	Attributes map[string]int         // A dictionary object, where each key corresponds to mesh attribute semantic and each value is the index of the accessor containing attribute's data. Required.
	Indices    *int                   // The index of the accessor that contains the indices. Not required.
	Material   *int                   // The index of the material to apply to this primitive when rendering. Not required.
	Mode       *int                   // The type of primitives to render. Not required. Default is 4 (TRIANGLES).
	Targets    []map[string]int       // An array of Morph Targets. Each Morph Target is a dictionary mapping attributes (only POSITION, NORMAL, and TANGENT supported) to their deviations in the Morph Target.
	Extensions map[string]interface{} // Dictionary object with extension-specific objects. Not required.
	Extras     interface{}            // Application-specific data. Not required.
}

// Sampler represents a texture sampler with properties for filtering and wrapping modes.
type Sampler struct {
	MagFilter  *int                   // Magnification filter. Not required.
	MinFilter  *int                   // Minification filter. Not required.
	WrapS      *int                   // s coordinate wrapping mode. Not required. Default is 10497 (REPEAT).
	WrapT      *int                   // t coordinate wrapping mode. Not required. Default is 10497 (REPEAT).
	Name       string                 // The user-defined name of this object. Not required.
	Extensions map[string]interface{} // Dictionary object with extension-specific objects. Not required.
	Extras     interface{}            // Application-specific data. Not required.
}

// Scene contains root nodes.
type Scene struct {
	Nodes      []int                  // The indices of the root nodes. Not required.
	Name       string                 // The user-defined name of this object. Not required.
	Extensions map[string]interface{} // Dictionary object with extension-specific objects. Not required. Not required.
	Extras     interface{}            // Application-specific data. Not required. Not required.
}

// Joints and matrices defining a skin.
type Skin struct {
	InverseBindMatrices int                    // The index of the accessor containing the floating-point 4x4 inverse-bind matrices. The default is that each matrix is a 4x4 identity matrix, which implies that inverse-bind matrices were pre-applied. Not required.
	Skeleton            *int                   // The index of the node used as a skeleton root. When undefined, joints transforms resolve to scene root. Not required.
	Joints              []int                  // Indices of skeleton nodes, used as joints in this skin. Required.
	Name                string                 // The user-define named of this object. Not required.
	Extensions          map[string]interface{} // Dictionary object with extension-specific objects. Not required.
	Extras              interface{}            // Application-specific data. Not required.

	cache *graphic.Skeleton // Cached skin.
}

// Sparse storage of attributes that deviate from their initialization value.
type Sparse struct {
	Count      int                    // Number of entries stored in the sparse array. Required.
	Indices    []int                  // Index array of size count that points to those accessor attributes that deviate from their initialization value. Indices must strictly increase. Required.
	Values     []int                  // Array of size count times number of components, storing the displaced accessor attributes pointed by indices. Substituted values must have the same componentType and number of components as the base accessor. Required.
	Extensions map[string]interface{} // Dictionary object with extension-specific objects. Not required.
	Extras     interface{}            // Application-specific data. Not required.
}

// Target represents the index of the node and TRS property than an animation channel targets.
type Target struct {
	Node int    // The index of the node to target. Not required.
	Path string // The name of the node's TRS property to modify, or the "weights" of the Morph Targets it instantiates. Required.
	// For the "translation" property, the values that are provided by the sampler are the translation along the x, y, and z axes.
	// For the "rotation" property, the values are a quaternion in the order (x, y, z, w), where w is the scalar.
	// For the "scale" property, the values are the scaling factors along the x, y, and z axes.
	Extensions map[string]interface{} // Dictionary object with extension-specific objects. Not required.
	Extras     interface{}            // Application-specific data. Not required.
}

// Texture represents a texture and its sampler.
type Texture struct {
	Sampler    *int                   // The index of the sampler used by this texture. When undefined, a sampler with REPEAT wrapping and AUTO filtering should be used. Not required.
	Source     int                    // The index of the image used by this texture. Not required.
	Name       string                 // The user-defined name of this object. Not required.
	Extensions map[string]interface{} // Dictionary object with extension-specific objects. Not required. Not required.
	Extras     interface{}            // Application-specific data. Not required. Not required.
}

// TextureInfo is a reference to a texture.
type TextureInfo struct {
	Index      int                    // The index of the texture. Required.
	TexCoord   int                    // The set index of texture's TEXCOORD attribute used for texture coordinate mapping. Not required. Default is 0.
	Extensions map[string]interface{} // Dictionary object with extension-specific objects. Not required.
	Extras     interface{}            // Application-specific data. Not required.
}

// Values is an array of size accessor.sparse.count times number of components storing the displaced accessor attributes pointed by accessor.sparse.indices.
type Values struct {
	BufferView int                    // The index of the bufferView with sparse values. Referenced bufferView can't have ARRAY_BUFFER or ELEMENT_ARRAY_BUFFER target. Required.
	ByteOffset int                    // The offset relative to the start of the bufferView in bytes. Must be aligned. Not required. Default is 0.
	Extensions map[string]interface{} // Dictionary object with extension-specific objects. Not required.
	Extras     interface{}            // Application-specific data. Not required.
}

// Primitive types.
const (
	POINTS         = 0
	LINES          = 1
	LINE_LOOP      = 2
	LINE_STRIP     = 3
	TRIANGLES      = 4
	TRIANGLE_STRIP = 5
	TRIANGLE_FAN   = 6
)

// OpenGL array types.
const (
	ARRAY_BUFFER         = 34962 // For vertex attributes
	ELEMENT_ARRAY_BUFFER = 34963 // For indices
)

// Texture filtering modes.
const (
	NEAREST                = 9728
	LINEAR                 = 9729
	NEAREST_MIPMAP_NEAREST = 9984
	LINEAR_MIPMAP_NEAREST  = 9985
	NEAREST_MIPMAP_LINEAR  = 9986
	LINEAR_MIPMAP_LINEAR   = 9987
)

// Texture sampling modes.
const (
	CLAMP_TO_EDGE   = 33071
	MIRRORED_REPEAT = 33648
	REPEAT          = 10497
)

// Possible componentType values.
const (
	BYTE           = 5120
	UNSIGNED_BYTE  = 5121
	SHORT          = 5122
	UNSIGNED_SHORT = 5123
	UNSIGNED_INT   = 5125
	FLOAT          = 5126
)

// Attribute element types.
const (
	SCALAR = "SCALAR"
	VEC2   = "VEC2"
	VEC3   = "VEC3"
	VEC4   = "VEC4"
	MAT2   = "MAT2"
	MAT3   = "MAT3"
	MAT4   = "MAT4"
)

// TypeSizes maps an attribute element type to the number of components it contains.
var TypeSizes = map[string]int{
	SCALAR: 1,
	VEC2:   2,
	VEC3:   3,
	VEC4:   4,
	MAT2:   4,
	MAT3:   9,
	MAT4:   16,
}

// AttributeName maps the glTF attribute name to the internal g3n attribute type.
var AttributeName = map[string]gls.AttribType{
	"POSITION":   gls.VertexPosition,
	"NORMAL":     gls.VertexNormal,
	"TANGENT":    gls.VertexTangent,
	"TEXCOORD_0": gls.VertexTexcoord,
	"TEXCOORD_1": gls.VertexTexcoord2,
	"COLOR_0":    gls.VertexColor,
	"JOINTS_0":   gls.SkinIndex,
	"WEIGHTS_0":  gls.SkinWeight,
}

type GLB struct {
	Header GLBHeader
	JSON   GLBChunk
	Data   GLBChunk
}

type GLBHeader struct {
	Magic   uint32
	Version uint32
	Length  uint32 // Not used directly
}

type GLBChunk struct {
	Length uint32
	Type   uint32
}

const (
	GLBMagic = 0x46546C67
	GLBJson  = 0x4E4F534A
	GLBBin   = 0x004E4942
)
