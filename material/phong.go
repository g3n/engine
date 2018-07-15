// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package material

import (
	"github.com/g3n/engine/math32"
)

// Phong material is identical to the Standard material but
// the calculation of the lighting model is done in the fragment shader.
type Phong struct {
	Standard // Embedded standard material
}

// NewPhong creates and returns a pointer to a new phong material
func NewPhong(color *math32.Color) *Phong {

	pm := new(Phong)
	pm.Standard.Init("phong", color)
	return pm
}
