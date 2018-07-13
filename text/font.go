// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package text

import (
	"github.com/g3n/engine/math32"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"strings"
)

// Font represents a TrueType font face.
// Attributes must be set prior to drawing.
type Font struct {
	ttf     *truetype.Font // The TrueType font
	face    font.Face      // The font face
	attrib  FontAttributes // Internal attribute cache
	fg      *image.Uniform // Text color cache
	bg      *image.Uniform // Background color cache
	changed bool           // Whether attributes have changed and the font face needs to be recreated
}

// FontAttributes contains tunable attributes of a font.
type FontAttributes struct {
	PointSize   float64      // Point size of the font
	DPI         float64      // Resolution of the font in dots per inch
	LineSpacing float64      // Spacing between lines (in terms of font height)
	Hinting     font.Hinting // Font hinting
}

// Font Hinting types.
const (
	HintingNone     = font.HintingNone
	HintingVertical = font.HintingVertical
	HintingFull     = font.HintingFull
)

// NewFont creates and returns a new font object using the specified TrueType font file.
func NewFont(ttfFile string) (*Font, error) {

	// Reads font bytes
	fontBytes, err := ioutil.ReadFile(ttfFile)
	if err != nil {
		return nil, err
	}
	return NewFontFromData(fontBytes)
}

// NewFontFromData creates and returns a new font object from the specified TTF data.
func NewFontFromData(fontData []byte) (*Font, error) {

	// Parses the font data
	ttf, err := truetype.Parse(fontData)
	if err != nil {
		return nil, err
	}

	f := new(Font)
	f.ttf = ttf

	// Initialize with default values
	f.attrib = FontAttributes{}
	f.attrib.PointSize = 12
	f.attrib.DPI = 72
	f.attrib.LineSpacing = 1.0
	f.attrib.Hinting = font.HintingNone
	f.SetColor(&math32.Color4{0, 0, 0, 1})

	// Create font face
	f.face = truetype.NewFace(f.ttf, &truetype.Options{
		Size:    f.attrib.PointSize,
		DPI:     f.attrib.DPI,
		Hinting: f.attrib.Hinting,
	})

	return f, nil
}

// SetPointSize sets the point size of the font.
func (f *Font) SetPointSize(size float64) {

	if size == f.attrib.PointSize {
		return
	}
	f.attrib.PointSize = size
	f.changed = true
}

// SetDPI sets the resolution of the font in dots per inches (DPI).
func (f *Font) SetDPI(dpi float64) {

	if dpi == f.attrib.DPI {
		return
	}
	f.attrib.DPI = dpi
	f.changed = true
}

// SetLineSpacing sets the amount of spacing between lines (in terms of font height).
func (f *Font) SetLineSpacing(spacing float64) {

	if spacing == f.attrib.LineSpacing {
		return
	}
	f.attrib.LineSpacing = spacing
	f.changed = true
}

// SetHinting sets the hinting type.
func (f *Font) SetHinting(hinting font.Hinting) {

	if hinting == f.attrib.Hinting {
		return
	}
	f.attrib.Hinting = hinting
	f.changed = true
}

// SetFgColor sets the text color.
func (f *Font) SetFgColor(color *math32.Color4) {

	f.fg = image.NewUniform(Color4RGBA(color))
}

// SetBgColor sets the background color.
func (f *Font) SetBgColor(color *math32.Color4) {

	f.bg = image.NewUniform(Color4RGBA(color))
}

// SetColor sets the text color to the specified value and makes the background color transparent.
// Note that for perfect transparency in the anti-aliased region it's important that the RGB components
// of the text and background colors match. This method handles that for the user.
func (f *Font) SetColor(fg *math32.Color4) {

	f.fg = image.NewUniform(Color4RGBA(fg))
	f.bg = image.NewUniform(Color4RGBA(&math32.Color4{fg.R, fg.G, fg.B, 0}))
}

// SetAttributes sets the font attributes.
func (f *Font) SetAttributes(fa *FontAttributes) {

	f.SetPointSize(fa.PointSize)
	f.SetDPI(fa.DPI)
	f.SetLineSpacing(fa.LineSpacing)
	f.SetHinting(fa.Hinting)
}

// updateFace updates the font face if parameters have changed.
func (f *Font) updateFace() {

	if f.changed {
		f.face = truetype.NewFace(f.ttf, &truetype.Options{
			Size:    f.attrib.PointSize,
			DPI:     f.attrib.DPI,
			Hinting: f.attrib.Hinting,
		})
		f.changed = false
	}
}

// MeasureText returns the minimum width and height in pixels necessary for an image to contain
// the specified text. The supplied text string can contain line break escape sequences (\n).
func (f *Font) MeasureText(text string) (int, int) {

	// Create font drawer
	f.updateFace()
	d := &font.Drawer{Dst: nil, Src: f.fg, Face: f.face}

	// Draw text
	var width, height int
	metrics := f.face.Metrics()
	lineHeight := (metrics.Ascent + metrics.Descent).Ceil()
	lineGap := int((f.attrib.LineSpacing - float64(1)) * float64(lineHeight))

	lines := strings.Split(text, "\n")
	for i, s := range lines {
		d.Dot = fixed.P(0, height)
		lineWidth := d.MeasureString(s).Ceil()
		if lineWidth > width {
			width = lineWidth
		}
		height += lineHeight
		if i > 1 {
			height += lineGap
		}
	}
	return width, height
}

// Metrics returns the font metrics.
func (f *Font) Metrics() font.Metrics {

	f.updateFace()
	return f.face.Metrics()
}

// DrawText draws the specified text on a new, tightly fitting image, and returns a pointer to the image.
func (f *Font) DrawText(text string) *image.RGBA {

	width, height := f.MeasureText(text)
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), f.bg, image.ZP, draw.Src)
	f.DrawTextOnImage(text, 0, 0, img)

	return img
}

// DrawTextOnImage draws the specified text on the specified image at the specified coordinates.
func (f *Font) DrawTextOnImage(text string, x, y int, dst *image.RGBA) {

	f.updateFace()
	d := &font.Drawer{Dst: dst, Src: f.fg, Face: f.face}

	// Draw text
	metrics := f.face.Metrics()
	py := y + metrics.Ascent.Round()
	lineHeight := (metrics.Ascent + metrics.Descent).Ceil()
	lineGap := int((f.attrib.LineSpacing - float64(1)) * float64(lineHeight))
	lines := strings.Split(text, "\n")
	for i, s := range lines {
		d.Dot = fixed.P(x, py)
		d.DrawString(s)
		py += lineHeight
		if i > 1 {
			py += lineGap
		}
	}
}

// Canvas is an image to draw on.
type Canvas struct {
	RGBA    *image.RGBA
	bgColor *image.Uniform
}

// NewCanvas creates and returns a pointer to a new canvas with the
// specified width and height in pixels and background color
func NewCanvas(width, height int, bgColor *math32.Color4) *Canvas {

	c := new(Canvas)
	c.RGBA = image.NewRGBA(image.Rect(0, 0, width, height))

	// Creates the image.Uniform for the background color
	c.bgColor = image.NewUniform(Color4RGBA(bgColor))

	// Draw image
	draw.Draw(c.RGBA, c.RGBA.Bounds(), c.bgColor, image.ZP, draw.Src)
	return c
}

// DrawText draws text at the specified position (in pixels)
// of this canvas, using the specified font.
// The supplied text string can contain line break escape sequences (\n).
func (c Canvas) DrawText(x, y int, text string, f *Font) {

	f.DrawTextOnImage(text, x, y, c.RGBA)
}

// DrawTextCaret draws text at the specified position (in pixels)
// of this canvas, using the specified font, and also a caret at
// the specified line and column.
// The supplied text string can contain line break escape sequences (\n).
// TODO Implement caret as a gui.Panel in gui.Edit
func (c Canvas) DrawTextCaret(x, y int, text string, f *Font, line, col int) error {

	// Creates drawer
	f.updateFace()
	d := &font.Drawer{Dst: c.RGBA, Src: f.fg, Face: f.face}

	// Draw text
	metrics := f.face.Metrics()
	py := y + metrics.Ascent.Round()
	lineHeight := (metrics.Ascent + metrics.Descent).Ceil()
	lineGap := int((f.attrib.LineSpacing - float64(1)) * float64(lineHeight))
	lines := strings.Split(text, "\n")
	for l, s := range lines {
		d.Dot = fixed.P(x, py)
		d.DrawString(s)
		// Checks for caret position
		if l == line && col <= StrCount(s) {
			width, _ := f.MeasureText(StrPrefix(s, col))
			// Draw caret vertical line
			caretH := int(f.attrib.PointSize) + 2
			caretY := int(d.Dot.Y>>6) - int(f.attrib.PointSize) + 2
			color := Color4RGBA(&math32.Color4{0, 0, 0, 1}) // Hardcoded to black
			for j := caretY; j < caretY+caretH; j++ {
				c.RGBA.Set(x+width, j, color)
			}
		}
		py += lineHeight
		if l > 1 {
			py += lineGap
		}
	}

	// TODO remove ?
	//	pt := freetype.Pt(font.marginX+x, font.marginY+y+int(font.ctx.PointToFixed(font.attrib.PointSize)>>6))
	//	for l, s := range lines {
	//		// Draw string
	//		_, err := font.ctx.DrawString(s, pt)
	//		if err != nil {
	//			return err
	//		}
	//		// Checks for caret position
	//		if l == line && col <= StrCount(s) {
	//			width, _, err := font.MeasureText(StrPrefix(s, col))
	//			if err != nil {
	//				return err
	//			}
	//			// Draw caret vertical line
	//			caretH := int(font.PointSize) + 2
	//			caretY := int(pt.Y>>6) - int(font.PointSize) + 2
	//			color := Color4RGBA(&math32.Color4{0, 0, 0, 1})
	//			for j := caretY; j < caretY+caretH; j++ {
	//				c.RGBA.Set(x+width, j, color)
	//			}
	//		}
	//		// Increment y coordinate
	//		pt.Y += font.ctx.PointToFixed(font.PointSize * font.LineSpacing)
	//	}
	return nil
}

// Color4RGBA converts a math32.Color4 to Go's color.RGBA.
func Color4RGBA(c *math32.Color4) color.RGBA {

	red := uint8(c.R * 0xFF)
	green := uint8(c.G * 0xFF)
	blue := uint8(c.B * 0xFF)
	alpha := uint8(c.A * 0xFF)
	return color.RGBA{red, green, blue, alpha}
}
