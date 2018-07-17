// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package equation implements SPOOK equations based on
// the 2007 PhD thesis of Claude Lacoursière titled
// "Ghosts and Machines: Regularized Variational Methods for
// Interactive Simulations of Multibodies with Dry Frictional Contacts"
package equation

import (
	"github.com/g3n/engine/math32"
)

// IBody is the interface of all body types.
type IBody interface {
	Index() int
	Position() math32.Vector3
	Velocity() math32.Vector3
	AngularVelocity() math32.Vector3
	Force() math32.Vector3
	Torque() math32.Vector3
	InvMassEff() float32
	InvRotInertiaWorldEff() *math32.Matrix3
}

// IEquation is the interface type for all equations types.
type IEquation interface {
	SetBodyA(IBody)
	BodyA() IBody
	SetBodyB(IBody)
	BodyB() IBody
	JeA() JacobianElement
	JeB() JacobianElement
	SetEnabled(state bool)
	Enabled() bool
	MinForce() float32
	MaxForce() float32
	Eps() float32
	SetMultiplier(multiplier float32)
	ComputeB(h float32) float32
	ComputeC() float32
}

// Equation is a SPOOK constraint equation.
type Equation struct {
	id         int
	minForce   float32       // Minimum (read: negative max) force to be applied by the constraint.
	maxForce   float32       // Maximum (read: positive max) force to be applied by the constraint.
	bA         IBody // Body "i"
	bB         IBody // Body "j"
	a          float32       // SPOOK parameter
	b          float32       // SPOOK parameter
	eps        float32       // SPOOK parameter
	jeA        JacobianElement
	jeB        JacobianElement
	enabled    bool
	multiplier float32 // A number, proportional to the force added to the bodies.
}

// NewEquation creates and returns a pointer to a new Equation object.
func NewEquation(bi, bj IBody, minForce, maxForce float32) *Equation {

	e := new(Equation)
	e.initialize(bi, bj, minForce, maxForce)
	return e
}

func (e *Equation) initialize(bi, bj IBody, minForce, maxForce float32) {

	//e.id = Equation.id++
	e.minForce = minForce //-1e6
	e.maxForce = maxForce //1e6

	e.bA = bi
	e.bB = bj
	e.a = 0
	e.b = 0
	e.eps = 0
	e.enabled = true
	e.multiplier = 0

	// Set typical spook params (k, d, dt)
	e.SetSpookParams(1e7, 3, 1/60)
}

func (e *Equation) SetBodyA(ibody IBody) {

	e.bA = ibody
}

func (e *Equation) BodyA() IBody {

	return e.bA
}

func (e *Equation) SetBodyB(ibody IBody) {

	e.bB = ibody
}

func (e *Equation) BodyB() IBody {

	return e.bB
}

func (e *Equation) JeA() JacobianElement {

	return e.jeA
}

func (e *Equation) JeB() JacobianElement {

	return e.jeB
}

// SetMinForce sets the minimum force to be applied by the constraint.
func (e *Equation) SetMinForce(minForce float32) {

	e.minForce = minForce
}

// MinForce returns the minimum force to be applied by the constraint.
func (e *Equation) MinForce() float32 {

	return e.minForce
}

// SetMaxForce sets the maximum force to be applied by the constraint.
func (e *Equation) SetMaxForce(maxForce float32) {

	e.maxForce = maxForce
}

// MaxForce returns the maximum force to be applied by the constraint.
func (e *Equation) MaxForce() float32 {

	return e.maxForce
}

// Returns epsilon - the regularization constant which is multiplied by the identity matrix.
func (e *Equation) Eps() float32 {

	return e.eps
}

// SetMultiplier sets the multiplier.
func (e *Equation) SetMultiplier(multiplier float32) {

	e.multiplier = multiplier
}

// MaxForce returns the multiplier.
func (e *Equation) Multiplier() float32 {

	return e.multiplier
}

// SetEnable sets the enabled flag of the equation.
func (e *Equation) SetEnabled(state bool) {

	e.enabled = state
}

// Enabled returns the enabled flag of the equation.
func (e *Equation) Enabled() bool {

	return e.enabled
}

// SetSpookParams recalculates a, b, eps.
func (e *Equation) SetSpookParams(stiffness, relaxation float32, timeStep float32) {

	e.a = 4.0 / (timeStep * (1 + 4*relaxation))
	e.b = (4.0 * relaxation) / (1 + 4*relaxation)
	e.eps = 4.0 / (timeStep * timeStep * stiffness * (1 + 4*relaxation))
}

// ComputeB computes the RHS of the SPOOK equation.
func (e *Equation) ComputeB(h float32) float32 {

	GW := e.ComputeGW()
	Gq := e.ComputeGq()
	GiMf := e.ComputeGiMf()
	return -Gq*e.a - GW*e.b - GiMf*h
}

// ComputeGq computes G*q, where q are the generalized body coordinates.
func (e *Equation) ComputeGq() float32 {

	xi := e.bA.Position()
	xj := e.bB.Position()
	spatA := e.jeA.Spatial()
	spatB := e.jeB.Spatial()
	return (&spatA).Dot(&xi) + (&spatB).Dot(&xj)
}

// ComputeGW computes G*W, where W are the body velocities.
func (e *Equation) ComputeGW() float32 {

	vA := e.bA.Velocity()
	vB := e.bB.Velocity()
	wA := e.bA.AngularVelocity()
	wB := e.bB.AngularVelocity()
	return e.jeA.MultiplyVectors(&vA, &wA) + e.jeB.MultiplyVectors(&vB, &wB)
}

// ComputeGiMf computes G*inv(M)*f, where M is the mass matrix with diagonal blocks for each body, and f are the forces on the bodies.
func (e *Equation) ComputeGiMf() float32 {

	forceA := e.bA.Force()
	forceB := e.bB.Force()

	iMfA := forceA.MultiplyScalar(e.bA.InvMassEff())
	iMfB := forceB.MultiplyScalar(e.bB.InvMassEff())

	torqueA := e.bA.Torque()
	torqueB := e.bB.Torque()

	invIiTaui := torqueA.ApplyMatrix3(e.bA.InvRotInertiaWorldEff())
	invIjTauj := torqueB.ApplyMatrix3(e.bB.InvRotInertiaWorldEff())

	return e.jeA.MultiplyVectors(iMfA, invIiTaui) + e.jeB.MultiplyVectors(iMfB, invIjTauj)
}

// ComputeGiMGt computes G*inv(M)*G'.
func (e *Equation) ComputeGiMGt() float32 {

	rotA := e.jeA.Rotational()
	rotB := e.jeB.Rotational()
	rotAcopy := e.jeA.Rotational()
	rotBcopy := e.jeB.Rotational()

	result := e.bA.InvMassEff() + e.bB.InvMassEff()
	result += rotA.ApplyMatrix3(e.bA.InvRotInertiaWorldEff()).Dot(&rotAcopy)
	result += rotB.ApplyMatrix3(e.bB.InvRotInertiaWorldEff()).Dot(&rotBcopy)

	return result
}

// ComputeC computes the denominator part of the SPOOK equation: C = G*inv(M)*G' + eps.
func (e *Equation) ComputeC() float32 {

	return e.ComputeGiMGt() + e.eps
}
