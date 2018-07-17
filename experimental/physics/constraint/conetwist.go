// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package constraint

import (
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/experimental/physics/equation"
)

// ConeTwist constraint.
type ConeTwist struct {
	PointToPoint
	axisA      *math32.Vector3 // Rotation axis, defined locally in bodyA.
	axisB      *math32.Vector3 // Rotation axis, defined locally in bodyB.
	coneEq     *equation.Cone
	twistEq    *equation.Rotational
	angle      float32
	twistAngle float32
}

// NewConeTwist creates and returns a pointer to a new ConeTwist constraint object.
func NewConeTwist(bodyA, bodyB IBody, pivotA, pivotB, axisA, axisB *math32.Vector3, angle, twistAngle, maxForce float32) *ConeTwist {

	ctc := new(ConeTwist)

	// Default of pivots and axes should be vec3(0)

	ctc.initialize(bodyA, bodyB, pivotA, pivotB, maxForce)

	ctc.axisA = axisA
	ctc.axisB = axisB
	ctc.axisA.Normalize()
	ctc.axisB.Normalize()

	ctc.angle = angle
	ctc.twistAngle = twistAngle

	ctc.coneEq = equation.NewCone(bodyA, bodyB, ctc.axisA, ctc.axisB, ctc.angle, maxForce)

	ctc.twistEq = equation.NewRotational(bodyA, bodyB, maxForce)
	ctc.twistEq.SetAxisA(ctc.axisA)
	ctc.twistEq.SetAxisB(ctc.axisB)

	// Make the cone equation push the bodies toward the cone axis, not outward
	ctc.coneEq.SetMaxForce(0)
	ctc.coneEq.SetMinForce(-maxForce)

	// Make the twist equation add torque toward the initial position
	ctc.twistEq.SetMaxForce(0)
	ctc.twistEq.SetMinForce(-maxForce)

	ctc.AddEquation(ctc.coneEq)
	ctc.AddEquation(ctc.twistEq)

	return ctc
}

// Update updates the equations with data.
func (ctc *ConeTwist) Update() {

	ctc.PointToPoint.Update()

	// Update the axes to the cone constraint
	worldAxisA := ctc.bodyA.VectorToWorld(ctc.axisA)
	worldAxisB := ctc.bodyB.VectorToWorld(ctc.axisB)

	ctc.coneEq.SetAxisA(&worldAxisA)
	ctc.coneEq.SetAxisB(&worldAxisB)

	// Update the world axes in the twist constraint
	tA, _ := ctc.axisA.RandomTangents()
	worldTA := ctc.bodyA.VectorToWorld(tA)
	ctc.twistEq.SetAxisA(&worldTA)

	tB, _ := ctc.axisB.RandomTangents()
	worldTB := ctc.bodyB.VectorToWorld(tB)
	ctc.twistEq.SetAxisB(&worldTB)

	ctc.coneEq.SetAngle(ctc.angle)
	ctc.twistEq.SetMaxAngle(ctc.twistAngle)
}
