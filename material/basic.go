// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package material

import ()

type Basic struct {
	Material // Embedded material
}

func NewBasic() *Basic {

	mb := new(Basic)
	mb.Material.Init()
	mb.SetShader("basic")
	return mb
}
