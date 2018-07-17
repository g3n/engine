// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package equation

import (
	"github.com/g3n/engine/math32"
)

// Contact is a contact/non-penetration constraint equation.
type Contact struct {
	Equation
	restitution float32         // "bounciness": u1 = -e*u0
	rA          *math32.Vector3 // World-oriented vector that goes from the center of bA to the contact point.
	rB          *math32.Vector3 // World-oriented vector that goes from the center of bB to the contact point.
	nA          *math32.Vector3 // Contact normal, pointing out of body A.
}

// NewContact creates and returns a pointer to a new Contact equation object.
func NewContact(bodyA, bodyB IBody, minForce, maxForce float32) *Contact {

	ce := new(Contact)

	// minForce default should be 0.

	ce.restitution = 0.5
	ce.rA = &math32.Vector3{0,0,0}
	ce.rB = &math32.Vector3{0,0,0}
	ce.nA = &math32.Vector3{0,0,0}

	ce.Equation.initialize(bodyA, bodyB, minForce, maxForce)

	return ce
}

func (ce *Contact) SetRestitution(r float32) {

	ce.restitution = r
}

func (ce *Contact) Restitution() float32 {

	return ce.restitution
}

func (ce *Contact) SetNormal(newNormal *math32.Vector3)  {

	ce.nA = newNormal
}

func (ce *Contact) Normal() math32.Vector3 {

	return *ce.nA
}

func (ce *Contact) SetRA(newRa *math32.Vector3)  {

	ce.rA = newRa
}

func (ce *Contact) RA() math32.Vector3 {

	return *ce.rA
}

func (ce *Contact) SetRB(newRb *math32.Vector3)  {

	ce.rB = newRb
}

func (ce *Contact) RB() math32.Vector3 {

	return *ce.rB
}

// ComputeB
func (ce *Contact) ComputeB(h float32) float32 {

	vA := ce.bA.Velocity()
	wA := ce.bA.AngularVelocity()

	vB := ce.bB.Velocity()
	wB := ce.bB.AngularVelocity()

	// Calculate cross products
	rnA := math32.NewVec3().CrossVectors(ce.rA, ce.nA)
	rnB := math32.NewVec3().CrossVectors(ce.rB, ce.nA)

	// g = xj+rB -(xi+rA)
	// G = [ -nA  -rnA  nA  rnB ]
	ce.jeA.SetSpatial(ce.nA.Clone().Negate())
	ce.jeA.SetRotational(rnA.Clone().Negate())
	ce.jeB.SetSpatial(ce.nA.Clone())
	ce.jeB.SetRotational(rnB.Clone())

	// Calculate the penetration vector
	posA := ce.bA.Position()
	posB := ce.bB.Position()
	penetrationVec := ce.rB.Clone().Add(&posB).Sub(ce.rA).Sub(&posA)

	g := ce.nA.Dot(penetrationVec)

	// Compute iteration
	ePlusOne := ce.restitution + 1
	GW := ePlusOne * vB.Dot(ce.nA) - ePlusOne * vA.Dot(ce.nA) + wB.Dot(rnB) - wA.Dot(rnA)
	GiMf := ce.ComputeGiMf()

	return -g*ce.a - GW*ce.b - h*GiMf
}

//// GetImpactVelocityAlongNormal returns the current relative velocity at the contact point.
//func (ce *Contact) GetImpactVelocityAlongNormal() float32 {
//
//	xi := ce.bA.Position().Add(ce.rA)
//	xj := ce.bB.Position().Add(ce.rB)
//
//	vi := ce.bA.GetVelocityAtWorldPoint(xi)
//	vj := ce.bB.GetVelocityAtWorldPoint(xj)
//
//	relVel := math32.NewVec3().SubVectors(vi, vj)
//
//    return ce.nA.Dot(relVel)
//}
