// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package physics

import "github.com/g3n/engine/math32"

// RigidBody represents a physics-driven solid body.
type RigidBody struct {
	// TODO :)
	position math32.Vector3 // World position of the center of gravity
}

// NewRigidBody creates and returns a pointer to a new RigidBody.
func NewRigidBody() *RigidBody {

	b := new(RigidBody)
	return b
}