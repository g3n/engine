// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package equation

import (
	"github.com/g3n/engine/math32"
)

// RotationalMotor is a rotational motor constraint equation.
// Tries to keep the relative angular velocity of the bodies to a given value.
type RotationalMotor struct {
	Equation // TODO maybe this should embed Rotational instead ?
	axisA       *math32.Vector3 // World oriented rotational axis
	axisB       *math32.Vector3 // World oriented rotational axis
	targetSpeed float32         // Target speed
}

// NewRotationalMotor creates and returns a pointer to a new RotationalMotor equation object.
func NewRotationalMotor(bodyA, bodyB IBody, maxForce float32) *RotationalMotor {

	re := new(RotationalMotor)
	re.Equation.initialize(bodyA, bodyB, -maxForce, maxForce)

	return re
}

// SetAxisA sets the axis of body A.
func (ce *RotationalMotor) SetAxisA(axisA *math32.Vector3) {

	ce.axisA = axisA
}

// AxisA returns the axis of body A.
func (ce *RotationalMotor) AxisA() math32.Vector3 {

	return *ce.axisA
}

// SetAxisB sets the axis of body B.
func (ce *RotationalMotor) SetAxisB(axisB *math32.Vector3) {

	ce.axisB = axisB
}

// AxisB returns the axis of body B.
func (ce *RotationalMotor) AxisB() math32.Vector3 {

	return *ce.axisB
}

// SetTargetSpeed sets the target speed.
func (ce *RotationalMotor) SetTargetSpeed(speed float32) {

	ce.targetSpeed = speed
}

// TargetSpeed returns the target speed.
func (ce *RotationalMotor) TargetSpeed() float32 {

	return ce.targetSpeed
}

// ComputeB
func (re *RotationalMotor) ComputeB(h float32) float32 {

	// g = 0
	// gdot = axisA * wi - axisB * wj
	// gdot = G * W = G * [vi wi vj wj]
	// =>
	// G = [0 axisA 0 -axisB]
	re.jeA.SetRotational(re.axisA.Clone())
	re.jeB.SetRotational(re.axisB.Clone().Negate())

	GW := re.ComputeGW() - re.targetSpeed
	GiMf := re.ComputeGiMf()

	return -GW*re.b - h*GiMf
}
