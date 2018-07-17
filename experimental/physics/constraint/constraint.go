// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package constraint implements physics constraints.
package constraint

import (
	"github.com/g3n/engine/experimental/physics/equation"
	"github.com/g3n/engine/math32"
)

type IBody interface {
	equation.IBody
	WakeUp()
	VectorToWorld(*math32.Vector3) math32.Vector3
	PointToLocal(*math32.Vector3) math32.Vector3
	VectorToLocal(*math32.Vector3) math32.Vector3
	Quaternion() *math32.Quaternion
}

type IConstraint interface {
	Update() // Update all the equations with data.
	Equations() []equation.IEquation
	CollideConnected() bool
	BodyA() IBody
	BodyB() IBody
}

// Constraint base struct.
type Constraint struct {
	equations []equation.IEquation // Equations to be solved in this constraint
	bodyA     IBody
	bodyB     IBody
	colConn   bool // Set to true if you want the bodies to collide when they are connected.
}

// NewConstraint creates and returns a pointer to a new Constraint object.
//func NewConstraint(bodyA, bodyB IBody, colConn, wakeUpBodies bool) *Constraint {
//
//	c := new(Constraint)
//	c.initialize(bodyA, bodyB, colConn, wakeUpBodies)
//	return c
//}

func (c *Constraint) initialize(bodyA, bodyB IBody, colConn, wakeUpBodies bool) {

	c.bodyA = bodyA
	c.bodyB = bodyB
	c.colConn = colConn // true

	if wakeUpBodies { // true
		if bodyA != nil {
			bodyA.WakeUp()
		}
		if bodyB != nil {
			bodyB.WakeUp()
		}
	}
}

// AddEquation adds an equation to the constraint.
func (c *Constraint) AddEquation(eq equation.IEquation) {

	c.equations = append(c.equations, eq)
}

// Equations returns the constraint's equations.
func (c *Constraint) Equations() []equation.IEquation {

	return c.equations
}

func (c *Constraint) CollideConnected() bool {

	return c.colConn
}

func (c *Constraint) BodyA() IBody {

	return c.bodyA
}

func (c *Constraint) BodyB() IBody {

	return c.bodyB
}

// SetEnable sets the enabled flag of the constraint equations.
func (c *Constraint) SetEnabled(state bool) {

	for i := range c.equations {
		c.equations[i].SetEnabled(state)
	}
}
