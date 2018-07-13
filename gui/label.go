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

// Label is a panel which contains a texture with text.
// The content size of the label panel is the exact size of the texture.
type Label struct {
	Panel                    // Embedded Panel
	font  *text.Font         // TrueType font face
	tex   *texture.Texture2D // Texture with text
	style *LabelStyle        // The style of the panel and font attributes
	text  string             // Text being displayed
}

// LabelStyle contains all the styling attributes of a Label.
// It's essentially a BasicStyle combined with FontAttributes.
type LabelStyle struct {
	PanelStyle
	text.FontAttributes
	FgColor math32.Color4
}

// NewLabel creates and returns a label panel with
// the specified text drawn using the default text font.
func NewLabel(text string) *Label {
	return NewLabelWithFont(text, StyleDefault().Font)
}

// NewIcon creates and returns a label panel with
// the specified text drawn using the default icon font.
func NewIcon(icon string) *Label {
	return NewLabelWithFont(icon, StyleDefault().FontIcon)
}

// NewLabelWithFont creates and returns a label panel with
// the specified text drawn using the specified font.
func NewLabelWithFont(msg string, font *text.Font) *Label {

	l := new(Label)
	l.initialize(msg, font)
	return l
}

// initialize initializes this label and is normally used by other
// components which contain a label.
func (l *Label) initialize(msg string, font *text.Font) {

	l.font = font
	l.Panel.Initialize(0, 0)

	// TODO: Remove this hack in an elegant way e.g. set the label style depending of if it's an icon or text label and have two defaults (one for icon labels one for text tabels)
	if font != StyleDefault().FontIcon {
		l.Panel.SetPaddings(2, 0, 2, 0)
	}

	// Copy the style based on the default Label style
	styleCopy := StyleDefault().Label
	l.style = &styleCopy

	l.SetText(msg)
}

// SetText sets and draws the label text using the font.
func (l *Label) SetText(text string) {

	// Need at least a character to get dimensions
	l.text = text
	if text == "" {
		text = " "
	}

	// Set font properties
	l.font.SetAttributes(&l.style.FontAttributes)
	l.font.SetColor(&l.style.FgColor)

	// Create an image with the text
	textImage := l.font.DrawText(text)

	// Create texture if it doesn't exist yet
	if l.tex == nil {
		l.tex = texture.NewTexture2DFromRGBA(textImage)
		l.tex.SetMagFilter(gls.NEAREST)
		l.tex.SetMinFilter(gls.NEAREST)
		l.Panel.Material().AddTexture(l.tex)
		// Otherwise update texture with new image
	} else {
		l.tex.SetFromRGBA(textImage)
	}

	// Update label panel dimensions
	l.Panel.SetContentSize(float32(textImage.Rect.Dx()), float32(textImage.Rect.Dy()))
}

// Text returns the label text.
func (l *Label) Text() string {

	return l.text
}

// SetColor sets the text color.
// Alpha is set to 1 (opaque).
func (l *Label) SetColor(color *math32.Color) *Label {

	l.style.FgColor.FromColor(color, 1.0)
	l.SetText(l.text)
	return l
}

// SetColor4 sets the text color.
func (l *Label) SetColor4(color4 *math32.Color4) *Label {

	l.style.FgColor = *color4
	l.SetText(l.text)
	return l
}

// Color returns the text color.
func (l *Label) Color() math32.Color4 {

	return l.style.FgColor
}

// SetBgColor sets the background color.
// The color alpha is set to 1.0
func (l *Label) SetBgColor(color *math32.Color) *Label {

	l.style.BgColor.FromColor(color, 1.0)
	l.Panel.SetColor4(&l.style.BgColor)
	l.SetText(l.text)
	return l
}

// SetBgColor4 sets the background color.
func (l *Label) SetBgColor4(color *math32.Color4) *Label {

	l.style.BgColor = *color
	l.Panel.SetColor4(&l.style.BgColor)
	l.SetText(l.text)
	return l
}

// BgColor returns returns the background color.
func (l *Label) BgColor() math32.Color4 {

	return l.style.BgColor
}

// SetFont sets the font.
func (l *Label) SetFont(f *text.Font) {

	l.font = f
	l.SetText(l.text)
}

// Font returns the font.
func (l *Label) Font() *text.Font {

	return l.font
}

// SetFontSize sets the point size of the font.
func (l *Label) SetFontSize(size float64) *Label {

	l.style.PointSize = size
	l.SetText(l.text)
	return l
}

// FontSize returns the point size of the font.
func (l *Label) FontSize() float64 {

	return l.style.PointSize
}

// SetFontDPI sets the resolution of the font in dots per inch (DPI).
func (l *Label) SetFontDPI(dpi float64) *Label {

	l.style.DPI = dpi
	l.SetText(l.text)
	return l
}

// FontDPI returns the resolution of the font in dots per inch (DPI).
func (l *Label) FontDPI() float64 {

	return l.style.DPI
}

// SetLineSpacing sets the spacing between lines.
func (l *Label) SetLineSpacing(spacing float64) *Label {

	l.style.LineSpacing = spacing
	l.SetText(l.text)
	return l
}

// LineSpacing returns the spacing between lines.
func (l *Label) LineSpacing() float64 {

	return l.style.LineSpacing
}

// setTextCaret sets the label text and draws a caret at the
// specified line and column.
// It is normally used by the Edit widget.
func (l *Label) setTextCaret(msg string, mx, width, line, col int) {

	// Set font properties
	l.font.SetAttributes(&l.style.FontAttributes)
	l.font.SetColor(&l.style.FgColor)

	// Create canvas and draw text
	_, height := l.font.MeasureText(msg)
	canvas := text.NewCanvas(width, height, &l.style.BgColor)
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
	l.text = msg
}
