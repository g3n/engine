// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package text

import (
	"bufio"
	"fmt"
	"github.com/g3n/engine/math32"
	"image"
	"image/png"
	"os"
	"unicode/utf8"
)

// CharInfo contains the information to locate a character in an Atlas
type CharInfo struct {
	X      int // Position X in pixels in the sheet image from left to right
	Y      int // Position Y in pixels in the sheet image from top to bottom
	Width  int // Char width in pixels
	Height int // Char heigh in pixels
	// Normalized position of char in the image
	OffsetX float32
	OffsetY float32
	RepeatX float32
	RepeatY float32
}

// Atlas represents an image containing characters and the information about their location in the image
type Atlas struct {
	Chars   []CharInfo
	Image   *image.RGBA
	Height  int // Recommended vertical space between two lines of text
	Ascent  int // Distance from the top of a line to its base line
	Descent int // Distance from the bottom of a line to its baseline
}

// NewAtlas returns a pointer to a new Atlas object
func NewAtlas(font *Font, first, last rune) *Atlas {

	a := new(Atlas)
	a.Chars = make([]CharInfo, last+1)

	// Get font metrics
	metrics := font.Metrics()
	a.Height = int(metrics.Height >> 6)
	a.Ascent = int(metrics.Ascent >> 6)
	a.Descent = int(metrics.Descent >> 6)

	const cols = 16
	col := 0
	encoded := make([]byte, 4)
	line := []byte{}
	lines := ""
	maxWidth := 0
	lastX := 0
	lastY := a.Descent
	nlines := 0
	for code := first; code <= last; code++ {
		// Encodes rune into UTF8 and appends to current line
		count := utf8.EncodeRune(encoded, code)
		line = append(line, encoded[:count]...)

		// Measure current line
		width, _ := font.MeasureText(string(line))

		// Sets current code fields
		cinfo := &a.Chars[code]
		cinfo.X = lastX
		cinfo.Y = lastY
		cinfo.Width = width - lastX - 1
		cinfo.Height += a.Height
		lastX = width
		fmt.Printf("%c: cinfo:%+v\n", code, cinfo)

		// Checks end of the current line
		col++
		if col >= cols || code == last {
			nlines++
			lines += string(line) + "\n"
			line = []byte{}
			// Checks max width
			if width > maxWidth {
				maxWidth = width
			}
			if code == last {
				break
			}
			col = 0
			lastX = 0
			lastY += a.Height
		}
	}
	height := (nlines * a.Height) + a.Descent

	// Draw atlas image
	canvas := NewCanvas(maxWidth, height, &math32.Color4{1, 1, 1, 1})
	canvas.DrawText(0, 0, lines, font)
	a.Image = canvas.RGBA

	// Calculate normalized char positions in the image
	fWidth := float32(maxWidth)
	fHeight := float32(height)
	for i := 0; i < len(a.Chars); i++ {
		char := &a.Chars[i]
		char.OffsetX = float32(char.X) / fWidth
		char.OffsetY = float32(char.Y) / fHeight
		char.RepeatX = float32(char.Width) / fWidth
		char.RepeatY = float32(char.Height) / fHeight
	}

	a.SavePNG("atlas.png")
	return a
}

// SavePNG saves the current atlas image as a PNG image file
func (a *Atlas) SavePNG(filename string) error {

	// Save that RGBA image to disk.
	outFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outFile.Close()

	b := bufio.NewWriter(outFile)
	err = png.Encode(b, a.Image)
	if err != nil {
		return err
	}
	err = b.Flush()
	if err != nil {
		return err
	}
	return nil
}
