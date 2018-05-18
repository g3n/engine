// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package physics

import (
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/core"
)

// Body represents a physics-driven solid body.
type Body struct {
	core.INode
	mass            float32        // Total mass
	velocity        math32.Vector3 // Linear velocity
	angularMass     math32.Matrix3 // Angular mass i.e. moment of inertia
	angularVelocity math32.Vector3 // Angular velocity
	position        math32.Vector3 // World position of the center of gravity
	static          bool // If true - the rigidBody does not move or rotate
}

// NewBody creates and returns a pointer to a new RigidBody.
func NewBody(inode core.INode) *Body {

	b := new(Body)
	b.INode = inode
	b.mass = 1
	return b
}