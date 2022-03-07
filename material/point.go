// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package material

import (
	"github.com/g3n/engine/math32"
)

// Point material is normally used for single point sprites
type Point struct {
	Standard // Embedded standard material
}

// NewPoint creates and returns a pointer to a new point material
func NewPoint(color *math32.Color) *Point {

	pm := new(Point)
	pm.Standard.Init("point", color)

	// Sets uniform's initial values
	pm.udata.emissive = *color
	pm.udata.psize = 1.0
	pm.udata.protationZ = 0
	return pm
}

// SetEmissiveColor sets the material emissive color
// The default is {0,0,0}
func (pm *Point) SetEmissiveColor(color *math32.Color) {

	pm.udata.emissive = *color
}

// SetSize sets the point size
func (pm *Point) SetSize(size float32) {

	pm.udata.psize = size
}

// SetRotationZ sets the point rotation around the Z axis.
func (pm *Point) SetRotationZ(rot float32) {

	pm.udata.protationZ = rot
}
