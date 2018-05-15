// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package texture

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"unsafe"

	"github.com/g3n/engine/gls"
)

// Texture2D represents a texture
type Texture2D struct {
	gs           *gls.GLS    // Pointer to OpenGL state
	refcount     int         // Current number of references
	texname      uint32      // Texture handle
	magFilter    uint32      // magnification filter
	minFilter    uint32      // minification filter
	wrapS        uint32      // wrap mode for s coordinate
	wrapT        uint32      // wrap mode for t coordinate
	iformat      int32       // internal format
	width        int32       // texture width in pixels
	height       int32       // texture height in pixels
	format       uint32      // format of the pixel data
	formatType   uint32      // type of the pixel data
	updateData   bool        // texture data needs to be sent
	updateParams bool        // texture parameters needs to be sent
	genMipmap    bool        // generate mipmaps flag
	data         interface{} // array with texture data
	uniUnit      gls.Uniform // Texture unit uniform location cache
	uniInfo      gls.Uniform // Texture info uniform location cache
	udata        struct {    // Combined uniform data in 3 vec2:
		offsetX float32
		offsetY float32
		repeatX float32
		repeatY float32
		flipY   float32
		visible float32
	}
}

func newTexture2D() *Texture2D {

	t := new(Texture2D)
	t.gs = nil
	t.refcount = 1
	t.texname = 0
	t.magFilter = gls.LINEAR
	t.minFilter = gls.LINEAR_MIPMAP_LINEAR
	t.wrapS = gls.CLAMP_TO_EDGE
	t.wrapT = gls.CLAMP_TO_EDGE
	t.updateData = false
	t.updateParams = true
	t.genMipmap = true

	// Initialize Uniform elements
	t.uniUnit.Init("MatTexture")
	t.uniInfo.Init("MatTexinfo")
	t.SetOffset(0, 0)
	t.SetRepeat(1, 1)
	t.SetFlipY(true)
	t.SetVisible(true)
	return t
}

// NewTexture2DFromImage creates and returns a pointer to a new Texture2D
// using the specified image file as data.
// Supported image formats are: PNG, JPEG and GIF.
func NewTexture2DFromImage(imgfile string) (*Texture2D, error) {

	// Decodes image file into RGBA8
	rgba, err := DecodeImage(imgfile)
	if err != nil {
		return nil, err
	}

	t := newTexture2D()
	t.SetFromRGBA(rgba)
	return t, nil
}

// NewTexture2DFromRGBA creates a new texture from a pointer to an RGBA image object.
func NewTexture2DFromRGBA(rgba *image.RGBA) *Texture2D {

	t := newTexture2D()
	t.SetFromRGBA(rgba)
	return t
}

// NewTexture2DFromData creates a new texture from data
func NewTexture2DFromData(width, height int, format int, formatType, iformat int, data interface{}) *Texture2D {

	t := newTexture2D()
	t.SetData(width, height, format, formatType, iformat, data)
	return t
}

// Incref increments the reference count for this texture
// and returns a pointer to the geometry.
// It should be used when this texture is shared by another
// material.
func (t *Texture2D) Incref() *Texture2D {

	t.refcount++
	return t
}

// Dispose decrements this texture reference count and
// if necessary releases OpenGL resources and C memory
// associated with this texture.
func (t *Texture2D) Dispose() {

	if t.refcount > 1 {
		t.refcount--
		return
	}
	if t.gs != nil {
		t.gs.DeleteTextures(t.texname)
		t.gs = nil
	}
}

// SetUniformNames sets the names of the uniforms in the shader for sampler and texture info.
func (t *Texture2D) SetUniformNames(sampler, info string) {

	t.uniUnit.Init(sampler)
	t.uniInfo.Init(info)
}

// GetUniformNames returns the names of the uniforms in the shader for sampler and texture info.
func (t *Texture2D) GetUniformNames() (sampler, info string) {

	return t.uniUnit.Name(), t.uniInfo.Name()
}

// SetImage sets a new image for this texture
func (t *Texture2D) SetImage(imgfile string) error {

	// Decodes image file into RGBA8
	rgba, err := DecodeImage(imgfile)
	if err != nil {
		return err
	}
	t.SetFromRGBA(rgba)
	return nil
}

// SetFromRGBA sets the texture data from the specified image.RGBA object
func (t *Texture2D) SetFromRGBA(rgba *image.RGBA) {

	t.SetData(
		rgba.Rect.Size().X,
		rgba.Rect.Size().Y,
		gls.RGBA,
		gls.UNSIGNED_BYTE,
		gls.RGBA8,
		rgba.Pix,
	)
}

// SetData sets the texture data
func (t *Texture2D) SetData(width, height int, format int, formatType, iformat int, data interface{}) {

	t.width = int32(width)
	t.height = int32(height)
	t.format = uint32(format)
	t.formatType = uint32(formatType)
	t.iformat = int32(iformat)
	t.data = data
	t.updateData = true
}

// SetVisible sets the visibility state of the texture
func (t *Texture2D) SetVisible(state bool) {

	if state {
		t.udata.visible = 1
	} else {
		t.udata.visible = 0
	}
}

// Visible returns the current visibility state of the texture
func (t *Texture2D) Visible() bool {

	if t.udata.visible == 0 {
		return false
	}

	return true
}

// SetMagFilter sets the filter to be applied when the texture element
// covers more than on pixel. The default value is gls.Linear.
func (t *Texture2D) SetMagFilter(magFilter uint32) {

	t.magFilter = magFilter
	t.updateParams = true
}

// SetMinFilter sets the filter to be applied when the texture element
// covers less than on pixel. The default value is gls.Linear.
func (t *Texture2D) SetMinFilter(minFilter uint32) {

	t.minFilter = minFilter
	t.updateParams = true
}

// SetWrapS set the wrapping mode for texture S coordinate
// The default value is GL_CLAMP_TO_EDGE;
func (t *Texture2D) SetWrapS(wrapS uint32) {

	t.wrapS = wrapS
	t.updateParams = true
}

// SetWrapT set the wrapping mode for texture T coordinate
// The default value is GL_CLAMP_TO_EDGE;
func (t *Texture2D) SetWrapT(wrapT uint32) {

	t.wrapT = wrapT
	t.updateParams = true
}

// SetRepeat set the repeat factor
func (t *Texture2D) SetRepeat(x, y float32) {

	t.udata.repeatX = x
	t.udata.repeatY = y
}

// Repeat returns the current X and Y repeat factors
func (t *Texture2D) Repeat() (float32, float32) {

	return t.udata.repeatX, t.udata.repeatY
}

// SetOffset sets the offset factor
func (t *Texture2D) SetOffset(x, y float32) {

	t.udata.offsetX = x
	t.udata.offsetY = y
}

// Offset returns the current X and Y offset factors
func (t *Texture2D) Offset() (float32, float32) {

	return t.udata.offsetX, t.udata.offsetY
}

// SetFlipY set the state for flipping the Y coordinate
func (t *Texture2D) SetFlipY(state bool) {

	if state {
		t.udata.flipY = 1
	} else {
		t.udata.flipY = 0
	}
}

// Width returns the texture width in pixels
func (t *Texture2D) Width() int {

	return int(t.width)
}

// Height returns the texture height in pixels
func (t *Texture2D) Height() int {

	return int(t.height)
}

// DecodeImage reads and decodes the specified image file into RGBA8.
// The supported image files are PNG, JPEG and GIF.
func DecodeImage(imgfile string) (*image.RGBA, error) {

	// Open image file
	file, err := os.Open(imgfile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decodes image
	img, _, err := image.Decode(file)
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

// RenderSetup is called by the material render setup
func (t *Texture2D) RenderSetup(gs *gls.GLS, slotIdx, uniIdx int) { // Could have as input - TEXTURE0 (slot) and uni location

	// One time initialization
	if t.gs == nil {
		t.texname = gs.GenTexture()
		t.gs = gs
	}

	// Sets the texture unit for this texture
	gs.ActiveTexture(uint32(gls.TEXTURE0 + slotIdx))
	gs.BindTexture(gls.TEXTURE_2D, t.texname)

	// Transfer texture data to OpenGL if necessary
	if t.updateData {
		gs.TexImage2D(
			gls.TEXTURE_2D, // texture type
			0,              // level of detail
			t.iformat,      // internal format
			t.width,        // width in texels
			t.height,       // height in texels
			0,              // border must be 0
			t.format,       // format of supplied texture data
			t.formatType,   // type of external format color component
			t.data,         // image data
		)
		// Generates mipmaps if requested
		if t.genMipmap {
			gs.GenerateMipmap(gls.TEXTURE_2D)
		}
		// No data to send
		t.updateData = false
	}

	// Sets texture parameters if needed
	if t.updateParams {
		gs.TexParameteri(gls.TEXTURE_2D, gls.TEXTURE_MAG_FILTER, int32(t.magFilter))
		gs.TexParameteri(gls.TEXTURE_2D, gls.TEXTURE_MIN_FILTER, int32(t.minFilter))
		gs.TexParameteri(gls.TEXTURE_2D, gls.TEXTURE_WRAP_S, int32(t.wrapS))
		gs.TexParameteri(gls.TEXTURE_2D, gls.TEXTURE_WRAP_T, int32(t.wrapT))
		t.updateParams = false
	}

	// Transfer texture unit uniform
	var location int32
	if uniIdx == 0 {
		location = t.uniUnit.Location(gs)
	} else {
		location = t.uniUnit.LocationIdx(gs, int32(uniIdx))
	}
	gs.Uniform1i(location, int32(slotIdx))

	// Transfer texture info combined uniform
	const vec2count = 3
	location = t.uniInfo.LocationIdx(gs, vec2count*int32(uniIdx))
	gs.Uniform2fvUP(location, vec2count, unsafe.Pointer(&t.udata))
}
