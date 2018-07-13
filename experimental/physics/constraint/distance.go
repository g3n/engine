// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package constraint

import (
	"github.com/g3n/engine/experimental/physics/equation"
)

// Distance is a distance constraint.
// Constrains two bodies to be at a constant distance from each others center of mass.
type Distance struct {
	Constraint
	distance   float32 // Distance
	equation   *equation.Contact
}

// NewDistance creates and returns a pointer to a new Distance constraint object.
func NewDistance(bodyA, bodyB IBody, distance, maxForce float32) *Distance {

	dc := new(Distance)
	dc.initialize(bodyA, bodyB, true, true)

	// Default distance should be: bodyA.position.distanceTo(bodyB.position)
	// Default maxForce should be: 1e6

	dc.distance = distance

	dc.equation = equation.NewContact(bodyA, bodyB, -maxForce, maxForce) // Make it bidirectional
	dc.AddEquation(dc.equation)

	return dc
}

// Update updates the equation with data.
func (dc *Distance) Update() {

	halfDist := dc.distance * 0.5

	posA := dc.bodyA.Position()
	posB := dc.bodyB.Position()

	normal := posB.Sub(&posA)
	normal.Normalize()

	dc.equation.SetRA(normal.Clone().MultiplyScalar(halfDist))
	dc.equation.SetRB(normal.Clone().MultiplyScalar(-halfDist))
}
