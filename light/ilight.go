// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package light

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
)

// ILight is the interface that must be implemented for all light types.
type ILight interface {
	RenderSetup(gs *gls.GLS, rinfo *core.RenderInfo, idx int)
}
