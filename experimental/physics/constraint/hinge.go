// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package constraint

import (
	"github.com/g3n/engine/experimental/physics/equation"
	"github.com/g3n/engine/math32"
)

// Hinge constraint.
// Think of it as a door hinge.
// It tries to keep the door in the correct place and with the correct orientation.
type Hinge struct {
	PointToPoint
	axisA   *math32.Vector3           // Rotation axis, defined locally in bodyA.
	axisB   *math32.Vector3           // Rotation axis, defined locally in bodyB.
	rotEq1  *equation.Rotational
	rotEq2  *equation.Rotational
	motorEq *equation.RotationalMotor
}

// NewHinge creates and returns a pointer to a new Hinge constraint object.
func NewHinge(bodyA, bodyB IBody, pivotA, pivotB, axisA, axisB *math32.Vector3, maxForce float32) *Hinge {

	hc := new(Hinge)

	hc.initialize(bodyA, bodyB, pivotA, pivotB, maxForce)

	hc.axisA = axisA
	hc.axisB = axisB
	hc.axisA.Normalize()
	hc.axisB.Normalize()

	hc.rotEq1 = equation.NewRotational(bodyA, bodyB, maxForce)
	hc.rotEq2 = equation.NewRotational(bodyA, bodyB, maxForce)
	hc.motorEq = equation.NewRotationalMotor(bodyA, bodyB, maxForce)
	hc.motorEq.SetEnabled(false) // Not enabled by default

	hc.AddEquation(hc.rotEq1)
	hc.AddEquation(hc.rotEq2)
	hc.AddEquation(hc.motorEq)

	return hc
}

func (hc *Hinge) SetMotorEnabled(state bool) {

	hc.motorEq.SetEnabled(state)
}

func (hc *Hinge) SetMotorSpeed(speed float32) {

	hc.motorEq.SetTargetSpeed(speed)
}

func (hc *Hinge) SetMotorMaxForce(maxForce float32) {

	hc.motorEq.SetMaxForce(maxForce)
	hc.motorEq.SetMinForce(-maxForce)
}

// Update updates the equations with data.
func (hc *Hinge) Update() {

	hc.PointToPoint.Update()

	// Get world axes
	quatA := hc.bodyA.Quaternion()
	quatB := hc.bodyB.Quaternion()

	worldAxisA := hc.axisA.Clone().ApplyQuaternion(quatA)
	worldAxisB := hc.axisB.Clone().ApplyQuaternion(quatB)

	t1, t2 := worldAxisA.RandomTangents()
	hc.rotEq1.SetAxisA(t1)
	hc.rotEq2.SetAxisA(t2)
	hc.rotEq1.SetAxisB(worldAxisB)
	hc.rotEq2.SetAxisB(worldAxisB)

	if hc.motorEq.Enabled() {
		hc.motorEq.SetAxisA(hc.axisA.Clone().ApplyQuaternion(quatA))
		hc.motorEq.SetAxisB(hc.axisB.Clone().ApplyQuaternion(quatB))
	}
}
