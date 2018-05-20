// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package physics implements a basic physics engine.
package solver

import (
	"github.com/g3n/engine/physics/equation"
	"time"
	"github.com/g3n/engine/math32"
)

// ISolver is the interface type for all constraint solvers.
type ISolver interface {
	// Solve should return the number of iterations performed
	Solve(frameDelta time.Duration, nBodies int) int
}

// Solution represents a solver solution
type Solution struct {
	VelocityDeltas        []math32.Vector3
	AngularVelocityDeltas []math32.Vector3
}

// Constraint equation solver base class.
type Solver struct {
	equations []*equation.Equation // All equations to be solved
}

// AddEquation adds an equation to the solver.
func (s *Solver) AddEquation(eq *equation.Equation) {

	s.equations = append(s.equations, eq)
}

// RemoveEquation removes the specified equation from the solver.
// Returns true if found, false otherwise.
func (s *Solver) RemoveEquation(eq *equation.Equation) bool {

	for pos, current := range s.equations {
		if current == eq {
			copy(s.equations[pos:], s.equations[pos+1:])
			s.equations[len(s.equations)-1] = nil
			s.equations = s.equations[:len(s.equations)-1]
			return true
		}
	}
	return false
}

// ClearEquations removes all equations from the solver.
func (s *Solver) ClearEquations() {

	s.equations = s.equations[0:0]
}
