// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"github.com/g3n/engine/math32"
)

// RenderInfo is passed into Render/RenderSetup calls
type RenderInfo struct {
	ViewMatrix math32.Matrix4 // Current camera view matrix
	ProjMatrix math32.Matrix4 // Current camera projection matrix
}
