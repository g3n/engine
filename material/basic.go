// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package material

// Basic is a simple material that uses the 'basic' shader
type Basic struct {
	Material // Embedded material
}

// NewBasic returns a pointer to a new Basic material
func NewBasic() *Basic {

	mb := new(Basic)
	mb.Material.Init()
	mb.SetShader("basic")
	return mb
}
