// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/g3n/engine/texture"
	"image"
)

type Image struct {
	Panel                    // Embedded panel
	tex   *texture.Texture2D // pointer to image texture
}

// NewImage creates and returns an image panel with the image
// from the specified image used as a texture.
// Initially the size of the panel content area is the exact size of the image.
func NewImage(imgfile string) (image *Image, err error) {

	tex, err := texture.NewTexture2DFromImage(imgfile)
	if err != nil {
		return nil, err
	}
	return NewImageFromTex(tex), nil
}

// NewImageFromRGBA creates and returns an image panel from the
// specified image
func NewImageFromRGBA(rgba *image.RGBA) *Image {

	tex := texture.NewTexture2DFromRGBA(rgba)
	return NewImageFromTex(tex)
}

// NewImageFromTex creates and returns an image panel from the specified texture2D
func NewImageFromTex(tex *texture.Texture2D) *Image {

	i := new(Image)
	i.Panel.Initialize(0, 0)
	i.tex = tex
	i.Panel.SetContentSize(float32(i.tex.Width()), float32(i.tex.Height()))
	i.Material().AddTexture(i.tex)
	return i
}

//func (i *Image) Clone() *Image {
//
//	return NewImageFromTex(i.tex.Clone())
//}
