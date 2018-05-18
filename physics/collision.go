// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package physics

import (
	"github.com/g3n/engine/core"
)

// Particle represents a physics-driven particle.
type Collision struct {
	// TODO
}

// NewCollision creates and returns a pointer to a new Collision.
func  NewCollision(inode core.INode) *Collision {

	c := new(Collision)
	return c
}