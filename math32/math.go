// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package math32 implements basic math functions which operate
// directly on float32 numbers without casting and contains
// types of common entities used in 3D Graphics such as vectors,
// matrices, quaternions and others.
package math32

import (
	"math"
)

const Pi = math.Pi
const degreeToRadiansFactor = math.Pi / 180
const radianToDegreesFactor = 180.0 / math.Pi

var Infinity = float32(math.Inf(1))

// DegToRad converts a number from degrees to radians
func DegToRad(degrees float32) float32 {

	return degrees * degreeToRadiansFactor
}

// RadToDeg converts a number from radians to degrees
func RadToDeg(radians float32) float32 {

	return radians * radianToDegreesFactor
}

// Clamp clamps x to the provided closed interval [a, b]
func Clamp(x, a, b float32) float32 {

	if x < a {
		return a
	}
	if x > b {
		return b
	}
	return x
}

// ClampInt clamps x to the provided closed interval [a, b]
func ClampInt(x, a, b int) int {

	if x < a {
		return a
	}
	if x > b {
		return b
	}
	return x
}

func Abs(v float32) float32 {
	return float32(math.Abs(float64(v)))
}

func Acos(v float32) float32 {
	return float32(math.Acos(float64(v)))
}

func Asin(v float32) float32 {
	return float32(math.Asin(float64(v)))
}

func Atan(v float32) float32 {
	return float32(math.Atan(float64(v)))
}

func Atan2(y, x float32) float32 {
	return float32(math.Atan2(float64(y), float64(x)))
}

func Ceil(v float32) float32 {
	return float32(math.Ceil(float64(v)))
}

func Cos(v float32) float32 {
	return float32(math.Cos(float64(v)))
}

func Floor(v float32) float32 {
	return float32(math.Floor(float64(v)))
}

func Inf(sign int) float32 {
	return float32(math.Inf(sign))
}

func Round(v float32) float32 {
	return Floor(v + 0.5)
}

func IsNaN(v float32) bool {
	return math.IsNaN(float64(v))
}

func Sin(v float32) float32 {
	return float32(math.Sin(float64(v)))
}

func Sqrt(v float32) float32 {
	return float32(math.Sqrt(float64(v)))
}

func Max(a, b float32) float32 {
	return float32(math.Max(float64(a), float64(b)))
}

func Min(a, b float32) float32 {
	return float32(math.Min(float64(a), float64(b)))
}

func Mod(a, b float32) float32 {
	return float32(math.Mod(float64(a), float64(b)))
}

func NaN() float32 {
	return float32(math.NaN())
}

func Pow(a, b float32) float32 {
	return float32(math.Pow(float64(a), float64(b)))
}

func Tan(v float32) float32 {
	return float32(math.Tan(float64(v)))
}
