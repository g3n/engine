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
// drawn using the current default font.
func NewLabel(msg string) *Label {

	l := new(Label)
	l.initialize(msg, StyleDefault.Font)
	return l
}

// NewIconLabel creates and returns a label panel using the specified text
// drawn using the default icon font.
func NewIconLabel(msg string) *Label {

	l := new(Label)
	l.initialize(msg, StyleDefault.FontIcon)
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
	l.fgColor = math32.Black4
	l.SetText(msg)
}

// SetText draws the label text using the current font
func (l *Label) SetText(msg string) {

	// Do not allow empty labels
	str := msg
	if len(msg) == 0 {
		str = " "
	}

	// Set font properties
	l.font.SetSize(l.fontSize)
	l.font.SetDPI(l.fontDPI)
	l.font.SetLineSpacing(l.lineSpacing)
	l.font.SetBgColor4(&l.bgColor)
	l.font.SetFgColor4(&l.fgColor)

	// Measure text
	width, height := l.font.MeasureText(str)
	// Create image canvas with the exact size of the texture
	// and draw the text.
	canvas := text.NewCanvas(width, height, &l.bgColor)
	canvas.DrawText(0, 0, str, l.font)

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
	l.currentText = str
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
