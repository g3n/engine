// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math32

import (
	"strings"
)

// Color describes an RGB color
type Color struct {
	R float32
	G float32
	B float32
}

// NewColor creates and returns a pointer to a new Color
// with the specified web standard color name (case insensitive).
// Returns nil if the color name not found
func NewColor(name string) *Color {

	c, ok := mapColorNames[strings.ToLower(name)]
	if !ok {
		return nil
	}
	return &c
}

// ColorName returns a Color with the specified standard web color name (case insensitive).
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

// Add adds to each RGB component of this color the correspondent component of other color
// Returns pointer to this updated color
func (c *Color) Add(other *Color) *Color {

	c.R += other.R
	c.G += other.G
	c.B += other.B
	return c
}

// AddColors adds to each RGB component of this color the correspondent component of color1 and color2
// Returns pointer to this updated color
func (c *Color) AddColors(color1, color2 *Color) *Color {

	c.R = color1.R + color2.R
	c.G = color1.G + color2.G
	c.B = color1.B + color2.B
	return c
}

// AddScalar adds the specified scalar value to each RGB component of this color
// Returns pointer to this updated color
func (c *Color) AddScalar(s float32) *Color {

	c.R += s
	c.G += s
	c.B += s
	return c
}

// Multiply multiplies each RGB component of this color by other
// Returns pointer to this updated color
func (c *Color) Multiply(other *Color) *Color {

	c.R *= other.R
	c.G *= other.G
	c.B *= other.B
	return c
}

// MultiplyScalar multiplies each RGB component of this color by the specified scalar.
// Returns pointer to this updated color
func (c *Color) MultiplyScalar(v float32) *Color {

	c.R *= v
	c.G *= v
	c.B *= v
	return c
}

// Lerp linear sets this color as the linear interpolation of itself
// with the specified color for the specified alpha.
// Returns pointer to this updated color
func (c *Color) Lerp(color *Color, alpha float32) *Color {

	c.R += (color.R - c.R) * alpha
	c.G += (color.G - c.G) * alpha
	c.B += (color.B - c.B) * alpha
	return c
}

// Equals returns if this color is equal to other
func (c *Color) Equals(other *Color) bool {

	return (c.R == other.R) && (c.G == other.G) && (c.B == other.B)
}

// IsColorName returns if the specified name is valid color name
func IsColorName(name string) (Color, bool) {

	c, ok := mapColorNames[strings.ToLower(name)]
	return c, ok
}

// mapColorNames maps standard web color names to a Color with
// the standard web color's RGB component values
var mapColorNames = map[string]Color{
	"aliceblue":            {0.941, 0.973, 1.000},
	"antiquewhite":         {0.980, 0.922, 0.843},
	"aqua":                 {0.000, 1.000, 1.000},
	"aquamarine":           {0.498, 1.000, 0.831},
	"azure":                {0.941, 1.000, 1.000},
	"beige":                {0.961, 0.961, 0.863},
	"bisque":               {1.000, 0.894, 0.769},
	"black":                {0.000, 0.000, 0.000},
	"blanchedalmond":       {1.000, 0.922, 0.804},
	"blue":                 {0.000, 0.000, 1.000},
	"blueviolet":           {0.541, 0.169, 0.886},
	"brown":                {0.647, 0.165, 0.165},
	"burlywood":            {0.871, 0.722, 0.529},
	"cadetblue":            {0.373, 0.620, 0.627},
	"chartreuse":           {0.498, 1.000, 0.000},
	"chocolate":            {0.824, 0.412, 0.118},
	"coral":                {1.000, 0.498, 0.314},
	"cornflowerblue":       {0.392, 0.584, 0.929},
	"cornsilk":             {1.000, 0.973, 0.863},
	"crimson":              {0.863, 0.078, 0.235},
	"cyan":                 {0.000, 1.000, 1.000},
	"darkblue":             {0.000, 0.000, 0.545},
	"darkcyan":             {0.000, 0.545, 0.545},
	"darkgoldenrod":        {0.722, 0.525, 0.043},
	"darkgray":             {0.663, 0.663, 0.663},
	"darkgreen":            {0.000, 0.392, 0.000},
	"darkgrey":             {0.663, 0.663, 0.663},
	"darkkhaki":            {0.741, 0.718, 0.420},
	"darkmagenta":          {0.545, 0.000, 0.545},
	"darkolivegreen":       {0.333, 0.420, 0.184},
	"darkorange":           {1.000, 0.549, 0.000},
	"darkorchid":           {0.600, 0.196, 0.800},
	"darkred":              {0.545, 0.000, 0.000},
	"darksalmon":           {0.914, 0.588, 0.478},
	"darkseagreen":         {0.561, 0.737, 0.561},
	"darkslateblue":        {0.282, 0.239, 0.545},
	"darkslategray":        {0.184, 0.310, 0.310},
	"darkslategrey":        {0.184, 0.310, 0.310},
	"darkturquoise":        {0.000, 0.808, 0.820},
	"darkviolet":           {0.580, 0.000, 0.827},
	"deeppink":             {1.000, 0.078, 0.576},
	"deepskyblue":          {0.000, 0.749, 1.000},
	"dimgray":              {0.412, 0.412, 0.412},
	"dimgrey":              {0.412, 0.412, 0.412},
	"dodgerblue":           {0.118, 0.565, 1.000},
	"firebrick":            {0.698, 0.133, 0.133},
	"floralwhite":          {1.000, 0.980, 0.941},
	"forestgreen":          {0.133, 0.545, 0.133},
	"fuchsia":              {1.000, 0.000, 1.000},
	"gainsboro":            {0.863, 0.863, 0.863},
	"ghostwhite":           {0.973, 0.973, 1.000},
	"gold":                 {1.000, 0.843, 0.000},
	"goldenrod":            {0.855, 0.647, 0.125},
	"gray":                 {0.502, 0.502, 0.502},
	"green":                {0.000, 0.502, 0.000},
	"greenyellow":          {0.678, 1.000, 0.184},
	"grey":                 {0.502, 0.502, 0.502},
	"honeydew":             {0.941, 1.000, 0.941},
	"hotpink":              {1.000, 0.412, 0.706},
	"indianred":            {0.804, 0.361, 0.361},
	"indigo":               {0.294, 0.000, 0.510},
	"ivory":                {1.000, 1.000, 0.941},
	"khaki":                {0.941, 0.902, 0.549},
	"lavender":             {0.902, 0.902, 0.980},
	"lavenderblush":        {1.000, 0.941, 0.961},
	"lawngreen":            {0.486, 0.988, 0.000},
	"lemonchiffon":         {1.000, 0.980, 0.804},
	"lightblue":            {0.678, 0.847, 0.902},
	"lightcoral":           {0.941, 0.502, 0.502},
	"lightcyan":            {0.878, 1.000, 1.000},
	"lightgoldenrodyellow": {0.980, 0.980, 0.824},
	"lightgray":            {0.827, 0.827, 0.827},
	"lightgreen":           {0.565, 0.933, 0.565},
	"lightgrey":            {0.827, 0.827, 0.827},
	"lightpink":            {1.000, 0.714, 0.757},
	"lightsalmon":          {1.000, 0.627, 0.478},
	"lightseagreen":        {0.125, 0.698, 0.667},
	"lightskyblue":         {0.529, 0.808, 0.980},
	"lightslategray":       {0.467, 0.533, 0.600},
	"lightslategrey":       {0.467, 0.533, 0.600},
	"lightsteelblue":       {0.690, 0.769, 0.871},
	"lightyellow":          {1.000, 1.000, 0.878},
	"lime":                 {0.000, 1.000, 0.000},
	"limegreen":            {0.196, 0.804, 0.196},
	"linen":                {0.980, 0.941, 0.902},
	"magenta":              {1.000, 0.000, 1.000},
	"maroon":               {0.502, 0.000, 0.000},
	"mediumaquamarine":     {0.400, 0.804, 0.667},
	"mediumblue":           {0.000, 0.000, 0.804},
	"mediumorchid":         {0.729, 0.333, 0.827},
	"mediumpurple":         {0.576, 0.439, 0.859},
	"mediumseagreen":       {0.235, 0.702, 0.443},
	"mediumslateblue":      {0.482, 0.408, 0.933},
	"mediumspringgreen":    {0.000, 0.980, 0.604},
	"mediumturquoise":      {0.282, 0.820, 0.800},
	"mediumvioletred":      {0.780, 0.082, 0.522},
	"midnightblue":         {0.098, 0.098, 0.439},
	"mintcream":            {0.961, 1.000, 0.980},
	"mistyrose":            {1.000, 0.894, 0.882},
	"moccasin":             {1.000, 0.894, 0.710},
	"navajowhite":          {1.000, 0.871, 0.678},
	"navy":                 {0.000, 0.000, 0.502},
	"oldlace":              {0.992, 0.961, 0.902},
	"olive":                {0.502, 0.502, 0.000},
	"olivedrab":            {0.420, 0.557, 0.137},
	"orange":               {1.000, 0.647, 0.000},
	"orangered":            {1.000, 0.271, 0.000},
	"orchid":               {0.855, 0.439, 0.839},
	"palegoldenrod":        {0.933, 0.910, 0.667},
	"palegreen":            {0.596, 0.984, 0.596},
	"paleturquoise":        {0.686, 0.933, 0.933},
	"palevioletred":        {0.859, 0.439, 0.576},
	"papayawhip":           {1.000, 0.937, 0.835},
	"peachpuff":            {1.000, 0.855, 0.725},
	"peru":                 {0.804, 0.522, 0.247},
	"pink":                 {1.000, 0.753, 0.796},
	"plum":                 {0.867, 0.627, 0.867},
	"powderblue":           {0.690, 0.878, 0.902},
	"purple":               {0.502, 0.000, 0.502},
	"red":                  {1.000, 0.000, 0.000},
	"rosybrown":            {0.737, 0.561, 0.561},
	"royalblue":            {0.255, 0.412, 0.882},
	"saddlebrown":          {0.545, 0.271, 0.075},
	"salmon":               {0.980, 0.502, 0.447},
	"sandybrown":           {0.957, 0.643, 0.376},
	"seagreen":             {0.180, 0.545, 0.341},
	"seashell":             {1.000, 0.961, 0.933},
	"sienna":               {0.627, 0.322, 0.176},
	"silver":               {0.753, 0.753, 0.753},
	"skyblue":              {0.529, 0.808, 0.922},
	"slateblue":            {0.416, 0.353, 0.804},
	"slategray":            {0.439, 0.502, 0.565},
	"slategrey":            {0.439, 0.502, 0.565},
	"snow":                 {1.000, 0.980, 0.980},
	"springgreen":          {0.000, 1.000, 0.498},
	"steelblue":            {0.275, 0.510, 0.706},
	"tan":                  {0.824, 0.706, 0.549},
	"teal":                 {0.000, 0.502, 0.502},
	"thistle":              {0.847, 0.749, 0.847},
	"tomato":               {1.000, 0.388, 0.278},
	"turquoise":            {0.251, 0.878, 0.816},
	"violet":               {0.933, 0.510, 0.933},
	"wheat":                {0.961, 0.871, 0.702},
	"white":                {1.000, 1.000, 1.000},
	"whitesmoke":           {0.961, 0.961, 0.961},
	"yellow":               {1.000, 1.000, 0.000},
	"yellowgreen":          {0.604, 0.804, 0.196},
}
