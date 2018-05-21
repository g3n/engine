// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package equation

import (
	"github.com/g3n/engine/math32"
)

// Rotational is a rotational constraint equation.
// Works to keep the local vectors orthogonal to each other in world space.
type Rotational struct {
	Equation
	axisA    *math32.Vector3 // Local axis in A
	axisB    *math32.Vector3 // Local axis in B
	maxAngle float32        // Max angle
}

// NewRotational creates and returns a pointer to a new Rotational equation object.
func NewRotational(bodyA, bodyB IBody, maxForce float32) *Rotational {

	re := new(Rotational)

	re.axisA = math32.NewVector3(1, 0, 0)
	re.axisB = math32.NewVector3(0, 1, 0)
	re.maxAngle = math32.Pi / 2

	re.Equation.initialize(bodyA, bodyB, -maxForce, maxForce)

	return re
}

// SetAxisA sets the axis of body A.
func (re *Rotational) SetAxisA(axisA *math32.Vector3) {

	re.axisA = axisA
}

// AxisA returns the axis of body A.
func (re *Rotational) AxisA() math32.Vector3 {

	return *re.axisA
}

// SetAxisB sets the axis of body B.
func (re *Rotational) SetAxisB(axisB *math32.Vector3) {

	re.axisB = axisB
}

// AxisB returns the axis of body B.
func (re *Rotational) AxisB() math32.Vector3 {

	return *re.axisB
}

// SetAngle sets the maximum angle.
func (re *Rotational) SetMaxAngle(angle float32) {

	re.maxAngle = angle
}

// MaxAngle returns the maximum angle.
func (re *Rotational) MaxAngle() float32 {

	return re.maxAngle
}

// ComputeB
func (re *Rotational) ComputeB(h float32) float32 {

	// Calculate cross products
	nAnB := math32.NewVec3().CrossVectors(re.axisA, re.axisB)
	nBnA := math32.NewVec3().CrossVectors(re.axisB, re.axisA)

	// g = nA * nj
	// gdot = (nj x nA) * wi + (nA x nj) * wj
	// G = [0 nBnA 0 nAnB]
	// W = [vi wi vj wj]
	re.jeA.SetRotational(nBnA)
	re.jeB.SetRotational(nAnB)

	g := math32.Cos(re.maxAngle) - re.axisA.Dot(re.axisB)
	GW := re.ComputeGW()
	GiMf := re.ComputeGiMf()

	return -g*re.a - GW*re.b - h*GiMf
}
