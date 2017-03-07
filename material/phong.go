// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package material

import (
	"github.com/g3n/engine/math32"
)

type Phong struct {
	Standard // Embedded standard material
}

// NewPhong creates and returns a pointer to a new phong material
func NewPhong(color *math32.Color) *Phong {

	pm := new(Phong)
	pm.Standard.Init("shaderPhong", color)
	return pm
}
