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
	"math"
	"strings"
)

type Font struct {
	ttf         *truetype.Font
	face        font.Face
	fgColor     math32.Color4
	bgColor     math32.Color4
	fontSize    float64
	fontDPI     float64
	lineSpacing float64
	fg          *image.Uniform
	bg          *image.Uniform
	hinting     font.Hinting
	changed     bool
}

const (
	HintingNone     = font.HintingNone
	HintingVertical = font.HintingVertical
	HintingFull     = font.HintingFull
)

// NewFont creates and returns a new font object using the specified
// truetype font file
func NewFont(fontfile string) (*Font, error) {

	// Reads font bytes
	fontBytes, err := ioutil.ReadFile(fontfile)
	if err != nil {
		return nil, err
	}
	return NewFontFromData(fontBytes)
}

// NewFontFromData creates and returns a new font object from the
// specified data
func NewFontFromData(fontData []byte) (*Font, error) {

	// Parses the font data
	ff, err := truetype.Parse(fontData)
	if err != nil {
		return nil, err
	}

	f := new(Font)
	f.ttf = ff
	f.fontSize = 12
	f.fontDPI = 72
	f.lineSpacing = 1.0
	f.hinting = font.HintingNone
	f.SetFgColor4(&math32.Color4{0, 0, 0, 1})
	f.SetBgColor4(&math32.Color4{1, 1, 1, 0})
	f.changed = false

	// Creates font face
	f.face = truetype.NewFace(f.ttf, &truetype.Options{
		Size:    f.fontSize,
		DPI:     f.fontDPI,
		Hinting: f.hinting,
	})
	return f, nil
}

// SetDPI sets the current DPI for the font
func (f *Font) SetDPI(dpi float64) {

	if dpi == f.fontDPI {
		return
	}
	f.fontDPI = dpi
	f.changed = true
}

// SetSize sets the size of the font
func (f *Font) SetSize(size float64) {

	if size == f.fontSize {
		return
	}
	f.fontSize = size
	f.changed = true
}

// SetHinting sets the hinting type
func (f *Font) SetHinting(hinting font.Hinting) {

	if hinting == f.hinting {
		return
	}
	f.hinting = hinting
	f.changed = true
}

// Size returns the current font size
func (f *Font) Size() float64 {

	return f.fontSize
}

// DPI returns the current font DPI
func (f *Font) DPI() float64 {

	return f.fontDPI
}

// SetFgColor sets the current foreground color of the font
// The alpha value is set to 1 (opaque)
func (f *Font) SetFgColor(color *math32.Color) {

	f.fgColor.R = color.R
	f.fgColor.G = color.G
	f.fgColor.B = color.B
	f.fgColor.A = 1.0
	f.fg = image.NewUniform(Color4NRGBA(&f.fgColor))
}

// SetFgColor4 sets the current foreground color of the font
func (f *Font) SetFgColor4(color *math32.Color4) {

	f.fgColor = *color
	f.fg = image.NewUniform(Color4NRGBA(color))
}

// FgColor returns the current foreground color
func (f *Font) FgColor4() math32.Color4 {

	return f.fgColor
}

// SetBgColor sets the current foreground color of the font
// The alpha value is set to 1 (opaque)
func (f *Font) SetBgColor(color *math32.Color) {

	f.bgColor.R = color.R
	f.bgColor.G = color.G
	f.bgColor.B = color.B
	f.bgColor.A = 1.0
	f.bg = image.NewUniform(Color4NRGBA(&f.fgColor))
}

// SetBgColor sets the current background color of the font
func (f *Font) SetBgColor4(color *math32.Color4) {

	f.bgColor = *color
	f.bg = image.NewUniform(Color4NRGBA(color))
}

// BgColor returns the current background color
func (f *Font) BgColor4() math32.Color4 {

	return f.bgColor
}

// SetLineSpacing sets the spacing between lines
func (f *Font) SetLineSpacing(spacing float64) {

	f.lineSpacing = spacing
}

// MeasureText returns the maximum width and height in pixels
// necessary for an image to contain the specified text.
// The supplied text string can contain line break escape sequences (\n).
func (f *Font) MeasureText(text string) (int, int) {

	// Creates drawer
	f.updateFace()
	d := &font.Drawer{Dst: nil, Src: f.fg, Face: f.face}

	// Draw text
	width := 0
	py := int(math.Ceil(f.fontSize * f.fontDPI / 72))
	dy := int(math.Ceil(f.fontSize * f.lineSpacing * f.fontDPI / 72))
	lines := strings.Split(text, "\n")
	for _, s := range lines {
		d.Dot = fixed.P(0, py)
		lfixed := d.MeasureString(s)
		lw := int(lfixed >> 6)
		if lw > width {
			width = lw
		}
		py += dy
	}
	height := py - dy/2
	return width, height
}

func (f *Font) Metrics() font.Metrics {

	f.updateFace()
	return f.face.Metrics()
}

// Canvas is an image to draw text
type Canvas struct {
	RGBA    *image.RGBA
	bgColor *image.Uniform
}

// Update font face if font parameters changed
func (f *Font) updateFace() {

	if f.changed {
		f.face = truetype.NewFace(f.ttf, &truetype.Options{
			Size:    f.fontSize,
			DPI:     f.fontDPI,
			Hinting: f.hinting,
		})
		f.changed = false
	}
}

// NewCanvas creates and returns a pointer to a new canvas with the
// specified width and height in pixels and background color
func NewCanvas(width, height int, bgColor *math32.Color4) *Canvas {

	c := new(Canvas)
	c.RGBA = image.NewRGBA(image.Rect(0, 0, width, height))

	// Creates the image.Uniform for the background color
	c.bgColor = image.NewUniform(Color4NRGBA(bgColor))

	// Draw image
	draw.Draw(c.RGBA, c.RGBA.Bounds(), c.bgColor, image.ZP, draw.Src)
	return c
}

// DrawText draws text at the specified position (in pixels)
// of this canvas, using the specified font.
// The supplied text string can contain line break escape sequences (\n).
func (c Canvas) DrawText(x, y int, text string, f *Font) {

	// Creates drawer
	f.updateFace()
	d := &font.Drawer{Dst: c.RGBA, Src: f.fg, Face: f.face}

	// Draw text
	py := y + int(math.Ceil(f.fontSize*f.fontDPI/72))
	dy := int(math.Ceil(f.fontSize * f.lineSpacing * f.fontDPI / 72))
	lines := strings.Split(text, "\n")
	for _, s := range lines {
		d.Dot = fixed.P(x, py)
		d.DrawString(s)
		py += dy
	}
}

// DrawTextCaret draws text at the specified position (in pixels)
// of this canvas, using the specified font, and also a caret at
// the specified line and column.
// The supplied text string can contain line break escape sequences (\n).
func (c Canvas) DrawTextCaret(x, y int, text string, f *Font, line, col int) error {

	// Creates drawer
	f.updateFace()
	d := &font.Drawer{Dst: c.RGBA, Src: f.fg, Face: f.face}

	py := y + int(math.Ceil(f.fontSize*f.fontDPI/72))
	dy := int(math.Ceil(f.fontSize * f.lineSpacing * f.fontDPI / 72))
	lines := strings.Split(text, "\n")
	for l, s := range lines {
		d.Dot = fixed.P(x, py)
		d.DrawString(s)
		// Checks for caret position
		if l == line && col <= StrCount(s) {
			width, _ := f.MeasureText(StrPrefix(s, col))
			// Draw caret vertical line
			caretH := int(f.fontSize) + 2
			//caretY := int(pt.Y>>6) - int(f.fontSize) + 2
			caretY := int(d.Dot.Y>>6) - int(f.fontSize) + 2
			color := Color4NRGBA(&math32.Color4{0, 0, 0, 1})
			for j := caretY; j < caretY+caretH; j++ {
				c.RGBA.Set(x+width, j, color)
			}
		}
		py += dy
	}
	//	pt := freetype.Pt(font.marginX+x, font.marginY+y+int(font.ctx.PointToFixed(font.fontSize)>>6))
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
	//			caretH := int(font.fontSize) + 2
	//			caretY := int(pt.Y>>6) - int(font.fontSize) + 2
	//			color := Color4NRGBA(&math32.Color4{0, 0, 0, 1})
	//			for j := caretY; j < caretY+caretH; j++ {
	//				c.RGBA.Set(x+width, j, color)
	//			}
	//		}
	//		// Increment y coordinate
	//		pt.Y += font.ctx.PointToFixed(font.fontSize * font.lineSpacing)
	//	}
	return nil
}

// Color4NRGBA converts a math32.Color4 to Go's image/color.NRGBA
// NON pre-multiplied alpha.
func Color4NRGBA(c *math32.Color4) color.NRGBA {

	red := uint8(c.R * 0xFF)
	green := uint8(c.G * 0xFF)
	blue := uint8(c.B * 0xFF)
	al := uint8(c.A * 0xFF)
	return color.NRGBA{red, green, blue, al}
}
