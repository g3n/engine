// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gltf

// GLTF is the root object for a glTF asset
type GLTF struct {
	ExtensionsUsed     []string
	ExtensionsRequired []string
	Accessors          []Accessor
	Animations         []Animation
	Asset              Asset
	Buffers            []Buffer
	BufferViews        []BufferView
	Cameras            []Camera
	Images             []Image
	Materials          []Material
	Meshes             []Mesh
	Nodes              []Node
	Samplers           []Sampler
	Scene              *int
	Scenes             []Scene
	Skins              []Skin
	Textures           []Texture
	Extensions         map[string]interface{}
	Extras             interface{}
	path               string // file path for resources
	data               []byte // binary file Chunk 1 data
}

// Accessor describes a view into a BufferView
type Accessor struct {
	BufferView    *int                   // The index of the buffer view
	ByteOffset    *int                   // The offset relative to the start of the BufferView in bytes
	ComponentType int                    // The datatype of components in the attribute
	Normalized    bool                   // Specifies whether integer data values should be normalized
	Count         int                    // The number of attributes referenced by this accessor
	Type          string                 // Specifies if the attribute is a scalar, vector or matrix
	Max           []float32              // Maximum value of each component in this attribute
	Min           []float32              // Minimum value of each component in this attribute
	Sparse        *Sparse                // Sparse storage attribute that deviates from their initialization value
	Name          string                 // The user-defined name of this object
	Extensions    map[string]interface{} // Dictionary object with extension specific objects
	Extras        interface{}            // Application-specific data
}

// A Keyframe animation
type Animation struct {
	Channels   []Channel              // An array of Channels
	Samplers   []Sampler              // An array of samplers that combines input and output accessors with an interpolation algorithm to define a keyframe graph
	Name       string                 // The user-defined name of this object
	Extensions map[string]interface{} // Dictionary object with extension specific objects
	Extras     interface{}            // Application-specific data
}

// Combines input and output accessors with an interpolation algorithm to define a keyframe graph
type AnimationSampler struct {
	Input         int                    // The index of the accessor containing keyframe input values
	Interpolation string                 // Interpolation algorithm
	Output        int                    // The index of an accessor containing keyframe output values
	Extensions    map[string]interface{} // Dictionary object with extension specific objects
	Extras        interface{}            // Application-specific data
}

// Metadata about the glTF asset.
type Asset struct {
	Copyright  string                 // A copyright message suitable for display to credit the content creator
	Generator  string                 // Tool that generated this glTF model. Useful for debugging
	Version    string                 // The glTF version that this asset targets
	MinVersion string                 // The minimum glTF version that this asset targets
	Extensions map[string]interface{} // Dictionary object with extension specific objects
	Extras     interface{}            // Application-specific data
}

// A Buffer points to binary geometry, animation or skins
type Buffer struct {
	Uri        string                 // The URI of the buffer
	ByteLength int                    // The length of the buffer in bytes
	Name       string                 // The user-defined name of this object
	Extensions map[string]interface{} // Dictionary object with extension specific objects
	Extras     interface{}            // Application-specific data
	data       []byte                 // cached buffer data
}

// A view into a buffer generally representing a subset of the buffer.
type BufferView struct {
	Buffer     int                    // The index of the buffer
	ByteOffset *int                   // The offset into the buffer in bytes
	ByteLength int                    // The length of the buffer view in bytes
	ByteStride *int                   // The stride in bytes
	Target     *int                   // The target that the GPU buffer should be bound to
	Name       string                 // The user-defined name of this object
	Extensions map[string]interface{} // Dictionary object with extension specific objects
	Extras     interface{}            // Application-specific data
}

// A camera's projection.
type Camera struct {
	Orthographic *Orthographic          // An orthographic camera containing properties to create an orthographic projection matrix
	Perspective  *Perspective           // A perspective camera containing properties to create a perspective projection matrix
	Type         string                 // Specifies if the camera uses a perspective or orthographic projection
	Name         string                 // The user-defined name of this object
	Extensions   map[string]interface{} // Dictionary object with extension specific objects
	Extras       interface{}            // Application-specific data
}

// Targets an animation's sampler at a node's property
type Channel struct {
	Sampler    int                    // The index of a sampler in this animation used to compute the value of the target
	Target     Target                 // The index of the node and TRS property to target
	Extensions map[string]interface{} // Dictionary object with extension specific objects
	Extras     interface{}            // Application-specific data
}

// Image data used to create a texture
type Image struct {
	Uri        string                 // The URI of the image
	MimeType   string                 // The image's MIME type
	BufferView *int                   // The index of the BufferView the contains the image
	Name       string                 // The user-defined name of this object
	Extensions map[string]interface{} // Dictionary object with extension specific objects
	Extras     interface{}            // Application-specific data
}

// Indices of those attributes that deviate from their initialization value.
type Indices struct {
	BufferView    int                    // The index of the BufferView with sparse indices
	ByteOffset    int                    // The offset relative to the start of the BufferView in bytes
	ComponentType int                    // The indices data type
	Extensions    map[string]interface{} // Dictionary object with extension specific objects
	Extras        interface{}            // Application-specific data
}

// Material describes the material appearance of a primitive
type Material struct {
	Name                 string                 // The user-defined name of this object
	Extensions           map[string]interface{} // Dictionary object with extension specific objects
	Extras               interface{}            // Application-specific data
	PbrMetallicRoughness *PbrMetallicRoughness
	NormalTexture        *NormalTextureInfo    // The normal map texture
	OcclusionTexture     *OcclusionTextureInfo // The occlusion map texture
	EmissiveTexture      *TextureInfo          // The emissive map texture
	EmissiveFactor       [3]float32            // The emissive color of the material
	AlphaMode            string                // The alpha rendering mode of the material
	AlphaCutoff          float32               // The alpha cutoff value of the material
	DoubleSided          bool                  // Specifies whether the material is double sided
}

// Mesh is a set of primitives to be rendered.
type Mesh struct {
	Primitives []Primitive            // Array of primitives
	Weights    []float32              // Array of weights to be applied to the Morph Targets
	Name       string                 // The user-define name of this object
	Extensions map[string]interface{} // Dictionary object with extension specific objects
	Extras     interface{}            // Application-specific data
}

// A Node in the hierarchy
type Node struct {
	Camera      *int                   // Index of the camera referenced by this node
	Children    []int                  // The indices of this node's children
	Skin        *int                   // The index of the skin referenced by this node
	Matrix      *[16]float32           // Floating point 4x4 transformation matrix in column-major order
	Mesh        *int                   // The index of the mesh in this node
	Rotation    *[4]float32            // The node's unit quaternion rotation in the order x,y,z,w
	Scale       *[3]float32            // The node's non-uniform scale
	Translation *[3]float32            // The node's translation
	Weights     []float32              // The weight's of the instantiated Morph Target
	Name        string                 // User-defined name of this object
	Extensions  map[string]interface{} // Dictionary object with extension specific objects
	Extras      interface{}            // Application-specific data
}

// Reference to a texture
type NormalTextureInfo struct {
	Index      int                    // The index of the texture
	TexCoord   int                    // The set index of texture's TEXCOORD attribute used for texture coordinate mapping
	Scale      float32                // The scalar multiplier applied to each normal vector of the normal texture
	Extensions map[string]interface{} // Dictionary object with extension specific objects
	Extras     interface{}            // Application-specific data
}

// Reference to a texture
type OcclusionTextureInfo struct {
	Index      int                    // The index of the texture
	TexCoord   int                    // The set index of texture's TEXCOORD attribute used for texture coordinate mapping
	Strength   float32                // A scalar multiplier controlling the amount of occlusion applied
	Extensions map[string]interface{} // Dictionary object with extension specific objects
	Extras     interface{}            // Application-specific data
}

// An orthographic camera containing properties to create an orthographic projection matrix
type Orthographic struct {
	Xmag       float32                // The floating-point horizontal magnification of the view
	Ymag       float32                // The floating-point vertical magnification of the view
	Zfar       float32                // The floating-point distance to the far clipping plane. zfar must be greater than znear
	Znear      float32                // The floating-point distance to the near clipping plane
	Extensions map[string]interface{} // Dictionary object with extension specific objects
	Extras     interface{}            // Application-specific data
}

// A set of parameter values that are used to define the metallic-roughness material model
// from Physically-Based Rendering (PBR) methodology.
type PbrMetallicRoughness struct {
	BaseColorFactor          [4]float32             // The material's base color factor
	BaseColorTexture         *TextureInfo           // The base color texture
	MetallicFactor           float32                // The metalness of the material
	RoughnessFactor          float32                // The roughness of the material
	MetallicRoughnessTexture *TextureInfo           // The metallic-roughness texture
	Extensions               map[string]interface{} // Dictionary object with extension specific objects
	Extras                   interface{}            // Application-specific data
}

// A perspective camera containing properties to create a perspective projection matrix
type Perspective struct {
	AspectRatio *float32               // The floating-point aspect ratio of the field of view
	Yfov        float32                // The floating-point vertical field of view in radians
	Zfar        *float32               // The floating-point distance to the far clipping plane
	Znear       float32                // The floating-point distance to the near clipping plane.
	Extensions  map[string]interface{} // Dictionary object with extension specific objects
	Extras      interface{}            // Application-specific data
}

// Geometry to be rendered with the given material
type Primitive struct {
	Attributes map[string]int         // A dictionary object, where each key corresponds to mesh attribute semantic and each value is the index of the accessor containing attribute's data
	Indices    *int                   // The index of the accessor that contains the indices
	Material   *int                   // The index of the material to apply to this primitive when rendering
	Mode       *int                   // The type of primitive to render
	Targets    []map[string]int       // An array of Morph Targets
	Extensions map[string]interface{} // Dictionary object with extension specific objects
	Extras     interface{}            // Application-specific data
}

// Texture sampler properties for filtering and wrapping modes
type Sampler struct {
	MagFilter  *int                   // Magnification filter
	MinFilter  *int                   // Minification filter
	WrapS      *int                   // s coordinate wrapping mode
	WrapT      *int                   // t coordinate wrapping mode
	Name       string                 // The user-define name for this object
	Extensions map[string]interface{} // Dictionary object with extension specific objects
	Extras     interface{}            // Application-specific data
}

// The root nodes of a scene
type Scene struct {
	Nodes      []int                  // The indices of each root node
	Name       string                 // The user-define name for this object
	Extensions map[string]interface{} // Dictionary object with extension specific objects
	Extras     interface{}            // Application-specific data
}

// Joints and matrices defining a skin.
type Skin struct {
	InverseBindMatrices int                    // The index of the accessor containing the 4x4 inverse-bind matrices
	Skeleton            int                    // The index of the node used as a skeleton root
	Joints              []int                  // Indices of skeleton nodes, used as joints in this skin
	Name                string                 // The user-define name for this object
	Extensions          map[string]interface{} // Dictionary object with extension specific objects
	Extras              interface{}            // Application-specific data
}

// Sparse storage of attributes that deviate from their initialization value.
type Sparse struct {
	Count int // Number of entries stored in the sparse array
	//Indices
	//Values
	Extensions map[string]interface{} // Dictionary object with extension specific objects
	Extras     interface{}            // Application-specific data
}

// The index of the node and TRS property than an animation channel targets
type Target struct {
	Node       int                    // The index of the node to target
	Path       string                 // The name of the node's TRS property to modify
	Extensions map[string]interface{} // Dictionary object with extension specific objects
	Extras     interface{}            // Application-specific data
}

// A texture and its sampler.
type Texture struct {
	Sampler    int                    // The index of the sampler used by this texture
	Source     int                    // The index of the image used by this texture
	Name       string                 // The user-define name for this object
	Extensions map[string]interface{} // Dictionary object with extension specific objects
	Extras     interface{}            // Application-specific data
}

// Reference to a texture.
type TextureInfo struct {
	Index      int                    // The index of the texture
	TexCoord   int                    // The set index of texture's TEXCOORD attribute used for texture coordinate mapping
	Extensions map[string]interface{} // Dictionary object with extension specific objects
	Extras     interface{}            // Application-specific data
}

// Array of size accessor.sparse.count times number of components storing
// the displaced accessor attributes pointed by accessor.sparse.indices.
type Values struct {
	BufferView int                    // The index of the bufferView with sparse values
	ByteOffset int                    // he offset relative to the start of the bufferView in bytes
	Extensions map[string]interface{} // Dictionary object with extension specific objects
	Extras     interface{}            // Application-specific data
}

const (
	POINTS                 = 0
	LINES                  = 1
	LINE_LOOP              = 2
	LINE_STRIP             = 3
	TRIANGLES              = 4
	TRIANGLE_STRIP         = 5
	TRIANGLE_FAN           = 6
	ARRAY_BUFFER           = 34962
	ELEMENT_ARRAY_BUFFER   = 34963
	NEAREST                = 9728
	LINEAR                 = 9729
	NEAREST_MIPMAP_NEAREST = 9984
	LINEAR_MIPMAP_NEAREST  = 9985
	NEAREST_MIPMAP_LINEAR  = 9986
	LINEAR_MIPMAP_LINEAR   = 9987
	CLAMP_TO_EDGE          = 33071
	MIRRORED_REPEAT        = 33648
	REPEAT                 = 10497
	UNSIGNED_BYTE          = 5121
	UNSIGNED_SHORT         = 5123
	UNSIGNED_INT           = 5125
	FLOAT                  = 5126
)

const (
	POSITION   = "POSITION"
	NORMAL     = "NORMAL"
	TANGENT    = "TANGENT"
	TEXCOORD_0 = "TEXCOORD_0"
	TEXCOORD_1 = "TEXCOORD_1"
	COLOR_0    = "COLOR_0"
	JOINTS_0   = "JOINTS_0"
	WEIGHTS_0  = "WEIGHTS_0"
	SCALAR     = "SCALAR"
	VEC2       = "VEC2"
	VEC3       = "VEC3"
	VEC4       = "VEC4"
	MAT2       = "MAT2"
	MAT3       = "MAT3"
	MAT4       = "MAT4"
)

// TypeSizes maps the attribute type to the number of its elements
var TypeSizes = map[string]int{
	SCALAR: 1,
	VEC2:   2,
	VEC3:   3,
	VEC4:   4,
	MAT2:   4,
	MAT3:   9,
	MAT4:   16,
}

type GLB struct {
	Header GLBHeader
	JSON   GLBChunk
	Data   GLBChunk
}

type GLBHeader struct {
	Magic   uint32
	Version uint32
	Length  uint32
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
