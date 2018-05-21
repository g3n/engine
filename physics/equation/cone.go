// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package equation

import (
	"github.com/g3n/engine/math32"
)

// Cone is a cone constraint equation.
// Works to keep the given body world vectors aligned, or tilted within a given angle from each other.
type Cone struct {
	Equation
	axisA    *math32.Vector3 // Local axis in A
	axisB    *math32.Vector3 // Local axis in B
	angle    float32         // The "cone angle" to keep
}

// NewCone creates and returns a pointer to a new Cone equation object.
func NewCone(bodyA, bodyB IBody, axisA, axisB *math32.Vector3, angle, maxForce float32) *Cone {

	ce := new(Cone)

	ce.axisA = axisA // new Vec3(1, 0, 0)
	ce.axisB = axisB // new Vec3(0, 1, 0)
	ce.angle = angle // 0

	ce.Equation.initialize(bodyA, bodyB, -maxForce, maxForce)

	return ce
}

// SetAxisA sets the axis of body A.
func (ce *Cone) SetAxisA(axisA *math32.Vector3) {

	ce.axisA = axisA
}

// AxisA returns the axis of body A.
func (ce *Cone) AxisA() math32.Vector3 {

	return *ce.axisA
}

// SetAxisB sets the axis of body B.
func (ce *Cone) SetAxisB(axisB *math32.Vector3) {

	ce.axisB = axisB
}

// AxisB returns the axis of body B.
func (ce *Cone) AxisB() math32.Vector3 {

	return *ce.axisB
}

// SetAngle sets the cone angle.
func (ce *Cone) SetAngle(angle float32) {

	ce.angle = angle
}

// MaxAngle returns the cone angle.
func (ce *Cone) Angle() float32 {

	return ce.angle
}

// ComputeB
func (ce *Cone) ComputeB(h float32) float32 {

	// The angle between two vector is:
	// cos(theta) = a * b / (length(a) * length(b) = { len(a) = len(b) = 1 } = a * b

	// g = a * b
	// gdot = (b x a) * wi + (a x b) * wj
	// G = [0 bxa 0 axb]
	// W = [vi wi vj wj]
	ce.jeA.SetRotational(math32.NewVec3().CrossVectors(ce.axisB, ce.axisA))
	ce.jeB.SetRotational(math32.NewVec3().CrossVectors(ce.axisA, ce.axisB))

	g := math32.Cos(ce.angle) - ce.axisA.Dot(ce.axisB)
	GW := ce.ComputeGW()
	GiMf := ce.ComputeGiMf()

	return -g*ce.a - GW*ce.b - h*GiMf
}
