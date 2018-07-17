// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package constraint

import (
	"github.com/g3n/engine/experimental/physics/equation"
	"github.com/g3n/engine/math32"
)

// PointToPoint is an offset constraint.
// Connects two bodies at the specified offset points.
type PointToPoint struct {
	Constraint
	pivotA *math32.Vector3   // Pivot, defined locally in bodyA.
	pivotB *math32.Vector3   // Pivot, defined locally in bodyB.
	eqX    *equation.Contact
	eqY    *equation.Contact
	eqZ    *equation.Contact
}

// NewPointToPoint creates and returns a pointer to a new PointToPoint constraint object.
func NewPointToPoint(bodyA, bodyB IBody, pivotA, pivotB *math32.Vector3, maxForce float32) *PointToPoint {

	ptpc := new(PointToPoint)
	ptpc.initialize(bodyA, bodyB, pivotA, pivotB, maxForce)

	return ptpc
}

func (ptpc *PointToPoint) initialize(bodyA, bodyB IBody, pivotA, pivotB *math32.Vector3, maxForce float32) {

	ptpc.Constraint.initialize(bodyA, bodyB, true, true)

	ptpc.pivotA = pivotA // default is zero vec3
	ptpc.pivotB = pivotB // default is zero vec3

	ptpc.eqX = equation.NewContact(bodyA, bodyB, -maxForce, maxForce)
	ptpc.eqY = equation.NewContact(bodyA, bodyB, -maxForce, maxForce)
	ptpc.eqZ = equation.NewContact(bodyA, bodyB, -maxForce, maxForce)

	ptpc.eqX.SetNormal(&math32.Vector3{1,0,0})
	ptpc.eqY.SetNormal(&math32.Vector3{0,1,0})
	ptpc.eqZ.SetNormal(&math32.Vector3{0,0,1})

	ptpc.AddEquation(&ptpc.eqX.Equation)
	ptpc.AddEquation(&ptpc.eqY.Equation)
	ptpc.AddEquation(&ptpc.eqZ.Equation)
}

// Update updates the equations with data.
func (ptpc *PointToPoint) Update() {

	// Rotate the pivots to world space
	xRi := ptpc.pivotA.Clone().ApplyQuaternion(ptpc.bodyA.Quaternion())
	xRj := ptpc.pivotA.Clone().ApplyQuaternion(ptpc.bodyA.Quaternion())

	ptpc.eqX.SetRA(xRi)
	ptpc.eqX.SetRB(xRj)
	ptpc.eqY.SetRA(xRi)
	ptpc.eqY.SetRB(xRj)
	ptpc.eqZ.SetRA(xRi)
	ptpc.eqZ.SetRB(xRj)
}
