// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package physics

type ICollidable interface {

}

// Simulation represents a physics simulation.
type Simulation struct {
	forceFields []*ForceField
	rigidBodies []*RigidBody
	particles   []*Particle

}

// NewSimulation creates and returns a pointer to a new physics simulation.
func NewSimulation() *Simulation {

	b := new(Simulation)
	return b
}

// AddForceField adds a force field to the simulation.
func (s *Simulation) AddForceField(ff *ForceField) {

	s.forceFields = append(s.forceFields, ff)
}

// RemoveForceField removes the specified force field from the simulation.
// Returns true if found or false otherwise.
func (s *Simulation) RemoveForceField(ff *ForceField) bool {

	for pos, current := range s.forceFields {
		if current == ff {
			copy(s.forceFields[pos:], s.forceFields[pos+1:])
			s.forceFields[len(s.forceFields)-1] = nil
			s.forceFields = s.forceFields[:len(s.forceFields)-1]
			return true
		}
	}
	return false
}

func (s *Simulation) CollisionBroadphase() {

}

func (s *Simulation) CheckCollisions(collidables []ICollidable) {

}

// Step steps the simulation.
func (s *Simulation) Step() {

	// Check for collisions (broad phase)
	//s.CollisionBroadphase()

	// Check for collisions (narrow phase)
	//s.CheckCollisions()

	// Apply forces/inertia/impulses


	// Update visual representation

}