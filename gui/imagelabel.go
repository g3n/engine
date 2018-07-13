// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/g3n/engine/math32"
)

/***************************************

 ImageLabel
 +--------------------------------+
 |  Image or Icon   Label         |
 |  +-----------+   +----------+  |
 |  |           |   |          |  |
 |  |           |   |          |  |
 |  +-----------+   +----------+  |
 +--------------------------------+

****************************************/

// ImageLabel is a panel which can contain an Image or Icon plus a Label side by side.
type ImageLabel struct {
	Panel        // Embedded panel
	label Label  // internal label
	image *Image // optional internal image
	icon  *Label // optional internal icon label
}

// ImageLabelStyle contains the styling of an ImageLabel.
type ImageLabelStyle BasicStyle

// NewImageLabel creates and returns a pointer to a new image label widget
// with the specified text for the label and no image/icon
func NewImageLabel(text string) *ImageLabel {

	il := new(ImageLabel)

	// Initializes the panel
	il.Panel.Initialize(0, 0)
	il.Panel.Subscribe(OnResize, func(evname string, ev interface{}) { il.recalc() })

	// Initializes the label
	il.label.initialize(text, StyleDefault().Font)
	il.label.Subscribe(OnResize, func(evname string, ev interface{}) { il.recalc() })
	il.Panel.Add(&il.label)

	il.recalc()
	return il
}

// SetText sets the text of the image label
func (il *ImageLabel) SetText(text string) {

	il.label.SetText(text)
}

// Text returns the current label text
func (il *ImageLabel) Text() string {

	return il.label.Text()
}

// SetIcon sets the image label icon from the default Icon font.
// If there is currently a selected image, it is removed
func (il *ImageLabel) SetIcon(icon string) {

	if il.image != nil {
		il.Panel.Remove(il.image)
		il.image = nil
	}
	if il.icon == nil {
		il.icon = NewIcon(icon)
		il.icon.SetFontSize(StyleDefault().Label.PointSize * 1.4)
		il.Panel.Add(il.icon)
	}
	il.icon.SetText(icon)
	il.recalc()
}

// SetImageVisible sets the image visibility
func (il *ImageLabel) SetImageVisible(vis bool) {

	if il.image == nil {
		return
	}
	il.image.SetVisible(vis)
}

// ImageVisible returns the image visibility
func (il *ImageLabel) ImageVisible() bool {

	if il.image == nil {
		return false
	}
	return il.image.Visible()
}

// SetImage sets the image label image
func (il *ImageLabel) SetImage(img *Image) {

	if il.icon != nil {
		il.Panel.Remove(il.icon)
		il.icon = nil
	}
	if il.image != nil {
		il.Panel.Remove(il.image)
	}
	il.image = img
	if img != nil {
		il.Panel.Add(il.image)
	}
	il.recalc()
}

// SetImageFromFile sets the image label image from the specified filename
// If there is currently a selected icon, it is removed
func (il *ImageLabel) SetImageFromFile(imgfile string) error {

	img, err := NewImage(imgfile)
	if err != nil {
		return err
	}
	il.SetImage(img)
	return nil
}

// SetColor sets the color of the label and icon text
func (il *ImageLabel) SetColor(color *math32.Color) {

	il.label.SetColor(color)
	if il.icon != nil {
		il.icon.SetColor(color)
	}
}

// SetColor4 sets the color4 of the label and icon
func (il *ImageLabel) SetColor4(color *math32.Color4) {

	il.label.SetColor4(color)
	if il.icon != nil {
		il.icon.SetColor4(color)
	}
}

// SetBgColor sets the color of the image label background
// The color alpha is set to 1.0
func (il *ImageLabel) SetBgColor(color *math32.Color) {

	il.Panel.SetColor(color)
	if il.icon != nil {
		il.icon.SetColor(color)
	}
	il.label.SetBgColor(color)
}

// SetBgColor4 sets the color4 of the image label background
func (il *ImageLabel) SetBgColor4(color *math32.Color4) {

	il.Panel.SetColor4(color)
	if il.icon != nil {
		il.icon.SetColor4(color)
	}
	il.label.SetBgColor4(color)
}

// SetFontSize sets the size of the image label font size
func (il *ImageLabel) SetFontSize(size float64) {

	il.label.SetFontSize(size)
}

// CopyFields copies another image label icon/image and text to this one
func (il *ImageLabel) CopyFields(other *ImageLabel) {

	il.label.SetText(other.label.Text())
	if other.icon != nil {
		il.SetIcon(other.icon.Text())
	}
	if other.image != nil {
		// TODO li.SetImage(other.image.Clone())
	}
	il.recalc()
}

// applyStyle applies the specified image label style
func (il *ImageLabel) applyStyle(s *ImageLabelStyle) {

	il.Panel.ApplyStyle(&s.PanelStyle)
	if il.icon != nil {
		il.icon.SetColor4(&s.FgColor)
	}
	il.label.SetColor4(&s.FgColor)
}

// recalc recalculates dimensions and positions from inside out
func (il *ImageLabel) recalc() {

	// Current width and height the content area
	width := il.Panel.ContentWidth()
	height := il.Panel.ContentHeight()

	// Image or icon width
	var imgWidth float32
	var spacing float32
	if il.image != nil {
		imgWidth = il.image.Width()
		spacing = 4
	} else if il.icon != nil {
		imgWidth = il.icon.Width()
		spacing = 4
	}

	// Sets new content width and height if necessary
	minWidth := imgWidth + spacing + il.label.Width()
	minHeight := il.label.Height()
	resize := false
	if width < minWidth {
		width = minWidth
		resize = true
	}
	if height < minHeight {
		height = minHeight
		resize = true
	}
	if resize {
		il.Panel.SetContentSize(width, height)
	}

	// Centralize horizontally
	px := (width - minWidth) / 2

	// Set label position
	ly := (height - il.label.Height()) / 2
	il.label.SetPosition(px+imgWidth+spacing, ly)

	// Image/icon position
	if il.image != nil {
		iy := (height - il.image.height) / 2
		il.image.SetPosition(px, iy)
	} else if il.icon != nil {
		il.icon.SetPosition(px, ly)
	}
}
