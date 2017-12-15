// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/text"
	"github.com/g3n/engine/texture"
)

// Label is a panel which contains a texture for rendering text
// The content size of the label panel is the exact size of texture
type Label struct {
	Panel       // Embedded panel
	fontSize    float64
	fontDPI     float64
	lineSpacing float64
	bgColor     math32.Color4
	fgColor     math32.Color4
	font        *text.Font
	tex         *texture.Texture2D // Pointer to texture with drawed text
	currentText string
}

// NewLabel creates and returns a label panel with the specified text
// drawn using the current default text font.
// If icon is true the text is drawn using the default icon font
func NewLabel(msg string, icon ...bool) *Label {

	l := new(Label)
	if len(icon) > 0 && icon[0] {
		l.initialize(msg, StyleDefault().FontIcon)
	} else {
		l.initialize(msg, StyleDefault().Font)
	}
	return l
}

// initialize initializes this label and is normally used by other
// gui types which contains a label.
func (l *Label) initialize(msg string, font *text.Font) {

	l.font = font
	l.Panel.Initialize(0, 0)
	l.fontSize = 14
	l.fontDPI = 72
	l.lineSpacing = 1.0
	l.bgColor = math32.Color4{0, 0, 0, 0}
	l.fgColor = math32.Color4{0, 0, 0, 1}
	l.SetText(msg)
}

// SetText draws the label text using the current font
func (l *Label) SetText(msg string) {

	// Need at least a character to get dimensions
	l.currentText = msg
	if msg == "" {
		msg = " "
	}

	// Set font properties
	l.font.SetSize(l.fontSize)
	l.font.SetDPI(l.fontDPI)
	l.font.SetLineSpacing(l.lineSpacing)
	l.font.SetBgColor4(&l.bgColor)
	l.font.SetFgColor4(&l.fgColor)

	// Measure text
	width, height := l.font.MeasureText(msg)
	// Create image canvas with the exact size of the texture
	// and draw the text.
	canvas := text.NewCanvas(width, height, &l.bgColor)
	canvas.DrawText(0, 0, msg, l.font)

	// Creates texture if if doesnt exist.
	if l.tex == nil {
		l.tex = texture.NewTexture2DFromRGBA(canvas.RGBA)
		l.tex.SetMagFilter(gls.NEAREST)
		l.tex.SetMinFilter(gls.NEAREST)
		l.Panel.Material().AddTexture(l.tex)
		// Otherwise update texture with new image
	} else {
		l.tex.SetFromRGBA(canvas.RGBA)
	}

	// Updates label panel dimensions
	l.Panel.SetContentSize(float32(width), float32(height))
}

// Text returns the current label text
func (l *Label) Text() string {

	return l.currentText
}

// SetColor sets the color of the label text
// The color alpha is set to 1.0
func (l *Label) SetColor(color *math32.Color) *Label {

	l.fgColor.FromColor(color, 1.0)
	l.SetText(l.currentText)
	return l
}

// SetColor4 sets the color4 of the label text
func (l *Label) SetColor4(color4 *math32.Color4) *Label {

	l.fgColor = *color4
	l.SetText(l.currentText)
	return l
}

// Color returns the current color of the label text
func (l *Label) Color() math32.Color4 {

	return l.fgColor
}

// SetBgColor sets the color of the label background
// The color alpha is set to 1.0
func (l *Label) SetBgColor(color *math32.Color) *Label {

	l.bgColor.FromColor(color, 1.0)
	l.Panel.SetColor4(&l.bgColor)
	l.SetText(l.currentText)
	return l
}

// SetBgColor4 sets the color4 of the label background
func (l *Label) SetBgColor4(color *math32.Color4) *Label {

	l.bgColor = *color
	l.Panel.SetColor4(&l.bgColor)
	l.SetText(l.currentText)
	return l
}

// BgColor returns the current color the label background
func (l *Label) BgColor() math32.Color4 {

	return l.bgColor
}

// SetFont sets this label text or icon font
func (l *Label) SetFont(f *text.Font) {

	l.font = f
	l.SetText(l.currentText)
}

// SetFontSize sets label font size
func (l *Label) SetFontSize(size float64) *Label {

	l.fontSize = size
	l.SetText(l.currentText)
	return l
}

// FontSize returns the current label font size
func (l *Label) FontSize() float64 {

	return l.fontSize
}

// SetFontDPI sets the font dots per inch
func (l *Label) SetFontDPI(dpi float64) *Label {

	l.fontDPI = dpi
	l.SetText(l.currentText)
	return l
}

// SetLineSpacing sets the spacing between lines.
// The default value is 1.0
func (l *Label) SetLineSpacing(spacing float64) *Label {

	l.lineSpacing = spacing
	l.SetText(l.currentText)
	return l
}

// setTextCaret sets the label text and draws a caret at the
// specified line and column.
// It is normally used by the Edit widget.
func (l *Label) setTextCaret(msg string, mx, width, line, col int) {

	// Set font properties
	l.font.SetSize(l.fontSize)
	l.font.SetDPI(l.fontDPI)
	l.font.SetLineSpacing(l.lineSpacing)
	l.font.SetBgColor4(&l.bgColor)
	l.font.SetFgColor4(&l.fgColor)

	// Create canvas and draw text
	_, height := l.font.MeasureText(msg)
	canvas := text.NewCanvas(width, height, &l.bgColor)
	canvas.DrawTextCaret(mx, 0, msg, l.font, line, col)

	// Creates texture if if doesnt exist.
	if l.tex == nil {
		l.tex = texture.NewTexture2DFromRGBA(canvas.RGBA)
		l.Panel.Material().AddTexture(l.tex)
		// Otherwise update texture with new image
	} else {
		l.tex.SetFromRGBA(canvas.RGBA)
	}
	// Set texture filtering parameters for text
	l.tex.SetMagFilter(gls.NEAREST)
	l.tex.SetMinFilter(gls.NEAREST)

	// Updates label panel dimensions
	l.Panel.SetContentSize(float32(width), float32(height))
	l.currentText = msg
}
