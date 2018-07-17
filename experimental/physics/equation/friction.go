// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package equation

import (
	"github.com/g3n/engine/math32"
)

// Friction is a friction constraint equation.
type Friction struct {
	Equation
	rA *math32.Vector3 // World-oriented vector that goes from the center of bA to the contact point.
	rB *math32.Vector3 // World-oriented vector that starts in body j position and goes to the contact point.
	t  *math32.Vector3 // Contact tangent
}

// NewFriction creates and returns a pointer to a new Friction equation object.
// slipForce should be +-F_friction = +-mu * F_normal = +-mu * m * g
func NewFriction(bodyA, bodyB IBody, slipForce float32) *Friction {

	fe := new(Friction)

	fe.rA = math32.NewVec3()
	fe.rB = math32.NewVec3()
	fe.t = math32.NewVec3()

	fe.Equation.initialize(bodyA, bodyB, -slipForce, slipForce)

	return fe
}

func (fe *Friction) SetTangent(newTangent *math32.Vector3)  {

	fe.t = newTangent
}

func (fe *Friction) Tangent() math32.Vector3 {

	return *fe.t
}

func (fe *Friction) SetRA(newRa *math32.Vector3)  {

	fe.rA = newRa
}

func (fe *Friction) RA() math32.Vector3 {

	return *fe.rA
}

func (fe *Friction) SetRB(newRb *math32.Vector3)  {

	fe.rB = newRb
}

func (fe *Friction) RB() math32.Vector3 {

	return *fe.rB
}

// ComputeB
func (fe *Friction) ComputeB(h float32) float32 {

	// Calculate cross products
	rtA := math32.NewVec3().CrossVectors(fe.rA, fe.t)
	rtB := math32.NewVec3().CrossVectors(fe.rB, fe.t)

	// G = [-t -rtA t rtB]
	// And remember, this is a pure velocity constraint, g is always zero!
	fe.jeA.SetSpatial(fe.t.Clone().Negate())
	fe.jeA.SetRotational(rtA.Clone().Negate())
	fe.jeB.SetSpatial(fe.t.Clone())
	fe.jeB.SetRotational(rtB.Clone())

	var GW = fe.ComputeGW()
	var GiMf = fe.ComputeGiMf()

	return -GW*fe.b - h*GiMf
}
