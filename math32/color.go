// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package math32

import (
	"strings"
)

// Type color describes an RGB color
type Color struct {
	R float32
	G float32
	B float32
}

// New creates and returns a pointer to a new Color
// with the specified web standard color name (case insensitive).
// Returns nil if the color name not found
func NewColor(name string) *Color {

	c, ok := mapColorNames[strings.ToLower(name)]
	if !ok {
		return nil
	}
	return &c
}

// Name returns a Color with the specified standard web color
// name (case insensitive).
// Returns black color if the specified color name not found
func ColorName(name string) Color {

	return mapColorNames[strings.ToLower(name)]
}

// NewColorHex creates and returns a pointer to a new color
// with its RGB components from the specified hex value
func NewColorHex(color uint) *Color {

	return (&Color{}).SetHex(color)
}

// Set sets this color individual R,G,B components
func (c *Color) Set(r, g, b float32) *Color {

	c.R = r
	c.G = g
	c.B = b
	return c
}

// SetHex sets the color RGB components from the
// specified integer interpreted as a color hex number
func (c *Color) SetHex(value uint) *Color {

	c.R = float32((value >> 16 & 255)) / 255
	c.G = float32((value >> 8 & 255)) / 255
	c.B = float32((value & 255)) / 255
	return c
}

// SetName sets the color RGB components from the
// specified standard web color name
func (c *Color) SetName(name string) *Color {

	color, ok := mapColorNames[strings.ToLower(name)]
	if ok {
		*c = color
	}
	return c
}

func (c *Color) Add(other *Color) *Color {

	c.R += other.R
	c.G += other.G
	c.B += other.B
	return c
}

func (c *Color) AddColors(color1, color2 *Color) *Color {

	c.R = color1.R + color2.R
	c.G = color1.G + color2.G
	c.B = color1.B + color2.B
	return c
}

func (c *Color) AddScalar(s float32) *Color {

	c.R += s
	c.G += s
	c.B += s
	return c
}

func (c *Color) Multiply(other *Color) *Color {

	c.R *= other.R
	c.G *= other.G
	c.B *= other.B
	return c
}

func (c *Color) MultiplyScalar(v float32) *Color {

	c.R *= v
	c.G *= v
	c.B *= v
	return c
}

func (c *Color) Lerp(color *Color, alpha float32) *Color {

	c.R += (color.R - c.R) * alpha
	c.G += (color.G - c.G) * alpha
	c.B += (color.B - c.B) * alpha
	return c
}

func (c *Color) Equals(other *Color) bool {

	return (c.R == other.R) && (c.G == other.G) && (c.B == other.B)
}

func IsColorName(name string) (Color, bool) {

	c, ok := mapColorNames[strings.ToLower(name)]
	return c, ok
}

// mapColorNames maps standard web color names to a Color with
// the standard web color's RGB component values
var mapColorNames = map[string]Color{
	"aliceblue":            Color{0.941, 0.973, 1.000},
	"antiquewhite":         Color{0.980, 0.922, 0.843},
	"aqua":                 Color{0.000, 1.000, 1.000},
	"aquamarine":           Color{0.498, 1.000, 0.831},
	"azure":                Color{0.941, 1.000, 1.000},
	"beige":                Color{0.961, 0.961, 0.863},
	"bisque":               Color{1.000, 0.894, 0.769},
	"black":                Color{0.000, 0.000, 0.000},
	"blanchedalmond":       Color{1.000, 0.922, 0.804},
	"blue":                 Color{0.000, 0.000, 1.000},
	"blueviolet":           Color{0.541, 0.169, 0.886},
	"brown":                Color{0.647, 0.165, 0.165},
	"burlywood":            Color{0.871, 0.722, 0.529},
	"cadetblue":            Color{0.373, 0.620, 0.627},
	"chartreuse":           Color{0.498, 1.000, 0.000},
	"chocolate":            Color{0.824, 0.412, 0.118},
	"coral":                Color{1.000, 0.498, 0.314},
	"cornflowerblue":       Color{0.392, 0.584, 0.929},
	"cornsilk":             Color{1.000, 0.973, 0.863},
	"crimson":              Color{0.863, 0.078, 0.235},
	"cyan":                 Color{0.000, 1.000, 1.000},
	"darkblue":             Color{0.000, 0.000, 0.545},
	"darkcyan":             Color{0.000, 0.545, 0.545},
	"darkgoldenrod":        Color{0.722, 0.525, 0.043},
	"darkgray":             Color{0.663, 0.663, 0.663},
	"darkgreen":            Color{0.000, 0.392, 0.000},
	"darkgrey":             Color{0.663, 0.663, 0.663},
	"darkkhaki":            Color{0.741, 0.718, 0.420},
	"darkmagenta":          Color{0.545, 0.000, 0.545},
	"darkolivegreen":       Color{0.333, 0.420, 0.184},
	"darkorange":           Color{1.000, 0.549, 0.000},
	"darkorchid":           Color{0.600, 0.196, 0.800},
	"darkred":              Color{0.545, 0.000, 0.000},
	"darksalmon":           Color{0.914, 0.588, 0.478},
	"darkseagreen":         Color{0.561, 0.737, 0.561},
	"darkslateblue":        Color{0.282, 0.239, 0.545},
	"darkslategray":        Color{0.184, 0.310, 0.310},
	"darkslategrey":        Color{0.184, 0.310, 0.310},
	"darkturquoise":        Color{0.000, 0.808, 0.820},
	"darkviolet":           Color{0.580, 0.000, 0.827},
	"deeppink":             Color{1.000, 0.078, 0.576},
	"deepskyblue":          Color{0.000, 0.749, 1.000},
	"dimgray":              Color{0.412, 0.412, 0.412},
	"dimgrey":              Color{0.412, 0.412, 0.412},
	"dodgerblue":           Color{0.118, 0.565, 1.000},
	"firebrick":            Color{0.698, 0.133, 0.133},
	"floralwhite":          Color{1.000, 0.980, 0.941},
	"forestgreen":          Color{0.133, 0.545, 0.133},
	"fuchsia":              Color{1.000, 0.000, 1.000},
	"gainsboro":            Color{0.863, 0.863, 0.863},
	"ghostwhite":           Color{0.973, 0.973, 1.000},
	"gold":                 Color{1.000, 0.843, 0.000},
	"goldenrod":            Color{0.855, 0.647, 0.125},
	"gray":                 Color{0.502, 0.502, 0.502},
	"green":                Color{0.000, 0.502, 0.000},
	"greenyellow":          Color{0.678, 1.000, 0.184},
	"grey":                 Color{0.502, 0.502, 0.502},
	"honeydew":             Color{0.941, 1.000, 0.941},
	"hotpink":              Color{1.000, 0.412, 0.706},
	"indianred":            Color{0.804, 0.361, 0.361},
	"indigo":               Color{0.294, 0.000, 0.510},
	"ivory":                Color{1.000, 1.000, 0.941},
	"khaki":                Color{0.941, 0.902, 0.549},
	"lavender":             Color{0.902, 0.902, 0.980},
	"lavenderblush":        Color{1.000, 0.941, 0.961},
	"lawngreen":            Color{0.486, 0.988, 0.000},
	"lemonchiffon":         Color{1.000, 0.980, 0.804},
	"lightblue":            Color{0.678, 0.847, 0.902},
	"lightcoral":           Color{0.941, 0.502, 0.502},
	"lightcyan":            Color{0.878, 1.000, 1.000},
	"lightgoldenrodyellow": Color{0.980, 0.980, 0.824},
	"lightgray":            Color{0.827, 0.827, 0.827},
	"lightgreen":           Color{0.565, 0.933, 0.565},
	"lightgrey":            Color{0.827, 0.827, 0.827},
	"lightpink":            Color{1.000, 0.714, 0.757},
	"lightsalmon":          Color{1.000, 0.627, 0.478},
	"lightseagreen":        Color{0.125, 0.698, 0.667},
	"lightskyblue":         Color{0.529, 0.808, 0.980},
	"lightslategray":       Color{0.467, 0.533, 0.600},
	"lightslategrey":       Color{0.467, 0.533, 0.600},
	"lightsteelblue":       Color{0.690, 0.769, 0.871},
	"lightyellow":          Color{1.000, 1.000, 0.878},
	"lime":                 Color{0.000, 1.000, 0.000},
	"limegreen":            Color{0.196, 0.804, 0.196},
	"linen":                Color{0.980, 0.941, 0.902},
	"magenta":              Color{1.000, 0.000, 1.000},
	"maroon":               Color{0.502, 0.000, 0.000},
	"mediumaquamarine":     Color{0.400, 0.804, 0.667},
	"mediumblue":           Color{0.000, 0.000, 0.804},
	"mediumorchid":         Color{0.729, 0.333, 0.827},
	"mediumpurple":         Color{0.576, 0.439, 0.859},
	"mediumseagreen":       Color{0.235, 0.702, 0.443},
	"mediumslateblue":      Color{0.482, 0.408, 0.933},
	"mediumspringgreen":    Color{0.000, 0.980, 0.604},
	"mediumturquoise":      Color{0.282, 0.820, 0.800},
	"mediumvioletred":      Color{0.780, 0.082, 0.522},
	"midnightblue":         Color{0.098, 0.098, 0.439},
	"mintcream":            Color{0.961, 1.000, 0.980},
	"mistyrose":            Color{1.000, 0.894, 0.882},
	"moccasin":             Color{1.000, 0.894, 0.710},
	"navajowhite":          Color{1.000, 0.871, 0.678},
	"navy":                 Color{0.000, 0.000, 0.502},
	"oldlace":              Color{0.992, 0.961, 0.902},
	"olive":                Color{0.502, 0.502, 0.000},
	"olivedrab":            Color{0.420, 0.557, 0.137},
	"orange":               Color{1.000, 0.647, 0.000},
	"orangered":            Color{1.000, 0.271, 0.000},
	"orchid":               Color{0.855, 0.439, 0.839},
	"palegoldenrod":        Color{0.933, 0.910, 0.667},
	"palegreen":            Color{0.596, 0.984, 0.596},
	"paleturquoise":        Color{0.686, 0.933, 0.933},
	"palevioletred":        Color{0.859, 0.439, 0.576},
	"papayawhip":           Color{1.000, 0.937, 0.835},
	"peachpuff":            Color{1.000, 0.855, 0.725},
	"peru":                 Color{0.804, 0.522, 0.247},
	"pink":                 Color{1.000, 0.753, 0.796},
	"plum":                 Color{0.867, 0.627, 0.867},
	"powderblue":           Color{0.690, 0.878, 0.902},
	"purple":               Color{0.502, 0.000, 0.502},
	"red":                  Color{1.000, 0.000, 0.000},
	"rosybrown":            Color{0.737, 0.561, 0.561},
	"royalblue":            Color{0.255, 0.412, 0.882},
	"saddlebrown":          Color{0.545, 0.271, 0.075},
	"salmon":               Color{0.980, 0.502, 0.447},
	"sandybrown":           Color{0.957, 0.643, 0.376},
	"seagreen":             Color{0.180, 0.545, 0.341},
	"seashell":             Color{1.000, 0.961, 0.933},
	"sienna":               Color{0.627, 0.322, 0.176},
	"silver":               Color{0.753, 0.753, 0.753},
	"skyblue":              Color{0.529, 0.808, 0.922},
	"slateblue":            Color{0.416, 0.353, 0.804},
	"slategray":            Color{0.439, 0.502, 0.565},
	"slategrey":            Color{0.439, 0.502, 0.565},
	"snow":                 Color{1.000, 0.980, 0.980},
	"springgreen":          Color{0.000, 1.000, 0.498},
	"steelblue":            Color{0.275, 0.510, 0.706},
	"tan":                  Color{0.824, 0.706, 0.549},
	"teal":                 Color{0.000, 0.502, 0.502},
	"thistle":              Color{0.847, 0.749, 0.847},
	"tomato":               Color{1.000, 0.388, 0.278},
	"turquoise":            Color{0.251, 0.878, 0.816},
	"violet":               Color{0.933, 0.510, 0.933},
	"wheat":                Color{0.961, 0.871, 0.702},
	"white":                Color{1.000, 1.000, 1.000},
	"whitesmoke":           Color{0.961, 0.961, 0.961},
	"yellow":               Color{1.000, 1.000, 0.000},
	"yellowgreen":          Color{0.604, 0.804, 0.196},
}
