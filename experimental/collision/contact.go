// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package collision

import "github.com/g3n/engine/math32"

// Contact describes a contact point, normal, and depth.
type Contact struct {
	Point  math32.Vector3
	Normal math32.Vector3
	Depth  float32
}
