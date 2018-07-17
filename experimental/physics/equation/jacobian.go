// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package equation

import "github.com/g3n/engine/math32"

// JacobianElement contains 6 entries, 3 spatial and 3 rotational degrees of freedom.
type JacobianElement struct {
	spatial    math32.Vector3
	rotational math32.Vector3
}

// SetSpatial sets the spatial component of the JacobianElement.
func (je *JacobianElement) SetSpatial(spatial *math32.Vector3) {

	je.spatial = *spatial
}

// Spatial returns the spatial component of the JacobianElement.
func (je *JacobianElement) Spatial() math32.Vector3 {

	return je.spatial
}

// Rotational sets the rotational component of the JacobianElement.
func (je *JacobianElement) SetRotational(rotational *math32.Vector3) {

	je.rotational = *rotational
}

// Rotational returns the rotational component of the JacobianElement.
func (je *JacobianElement) Rotational() math32.Vector3 {

	return je.rotational
}

// MultiplyElement multiplies the JacobianElement with another JacobianElement.
// None of the elements are changed.
func (je *JacobianElement) MultiplyElement(je2 *JacobianElement) float32 {

	return je.spatial.Dot(&je2.spatial) + je.rotational.Dot(&je2.rotational)
}

// MultiplyElement multiplies the JacobianElement with two vectors.
// None of the elements are changed.
func (je *JacobianElement) MultiplyVectors(spatial *math32.Vector3, rotational *math32.Vector3) float32 {

	return je.spatial.Dot(spatial) + je.rotational.Dot(rotational)
}
