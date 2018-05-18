// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package physics

import (
	"time"
)

type ICollidable interface {
	Collided() bool
	Static() bool
	//GetCollisions() []*Collision
}

type CollisionGroup struct {
	// TODO future
}

// Simulation represents a physics simulation.
type Simulation struct {
	forceFields []ForceField
	bodies      []*Body
	particles   []*Particle

	//dynamical   []int // non collidable but still affected by force fields
	//collidables []ICollidable

}

// NewSimulation creates and returns a pointer to a new physics simulation.
func NewSimulation() *Simulation {

	b := new(Simulation)
	return b
}

// AddForceField adds a force field to the simulation.
func (s *Simulation) AddForceField(ff ForceField) {

	s.forceFields = append(s.forceFields, ff)
}

// RemoveForceField removes the specified force field from the simulation.
// Returns true if found, false otherwise.
func (s *Simulation) RemoveForceField(ff ForceField) bool {

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

// AddBody adds a body to the simulation.
func (s *Simulation) AddBody(rb *Body) {

	s.bodies = append(s.bodies, rb)
}

// RemoveBody removes the specified body from the simulation.
// Returns true if found, false otherwise.
func (s *Simulation) RemoveBody(rb *Body) bool {

	for pos, current := range s.bodies {
		if current == rb {
			copy(s.bodies[pos:], s.bodies[pos+1:])
			s.bodies[len(s.bodies)-1] = nil
			s.bodies = s.bodies[:len(s.bodies)-1]
			return true
		}
	}
	return false
}

// AddParticle adds a particle to the simulation.
func (s *Simulation) AddParticle(rb *Particle) {

	s.particles = append(s.particles, rb)
}

// RemoveParticle removes the specified particle from the simulation.
// Returns true if found, false otherwise.
func (s *Simulation) RemoveParticle(rb *Particle) bool {

	for pos, current := range s.particles {
		if current == rb {
			copy(s.particles[pos:], s.particles[pos+1:])
			s.particles[len(s.particles)-1] = nil
			s.particles = s.particles[:len(s.particles)-1]
			return true
		}
	}
	return false
}

func (s *Simulation) applyForceFields(frameDelta time.Duration) {

	for _, ff := range s.forceFields {
		for _, rb := range s.bodies {
			pos := rb.GetNode().Position()
			force := ff.ForceAt(&pos)
			//log.Error("force: %v", force)
			//log.Error("mass: %v", rb.mass)
			//log.Error("frameDelta: %v", frameDelta.Seconds())
			vdiff := force.MultiplyScalar(float32(frameDelta.Seconds())/rb.mass)
			//log.Error("vdiff: %v", vdiff)
			rb.velocity.Add(vdiff)
		}
	}

	for _, ff := range s.forceFields {
		for _, rb := range s.particles {
			pos := rb.GetNode().Position()
			force := ff.ForceAt(&pos)
			//log.Error("force: %v", force)
			//log.Error("mass: %v", rb.mass)
			//log.Error("frameDelta: %v", frameDelta.Seconds())
			vdiff := force.MultiplyScalar(float32(frameDelta.Seconds())/rb.mass)
			//log.Error("vdiff: %v", vdiff)
			rb.velocity.Add(vdiff)
		}
	}
}

// updatePositions integrates the velocity of the objects to obtain the position in the next frame.
func (s *Simulation) updatePositions(frameDelta time.Duration) {

	for _, rb := range s.bodies {
		pos := rb.GetNode().Position()
		posDelta := rb.velocity
		posDelta.MultiplyScalar(float32(frameDelta.Seconds()))
		pos.Add(&posDelta)
		rb.GetNode().SetPositionVec(&pos)
	}

	for _, rb := range s.particles {
		pos := rb.GetNode().Position()
		posDelta := rb.velocity
		posDelta.MultiplyScalar(float32(frameDelta.Seconds()))
		pos.Add(&posDelta)
		rb.GetNode().SetPositionVec(&pos)
	}
}

//func (s *Simulation) CheckCollisions() []*Collision{//collidables []ICollidable) {
//
//	return
//}


func (s *Simulation) Backtrack() {

}

// Step steps the simulation.
func (s *Simulation) Step(frameDelta time.Duration) {

	// Check for collisions
	//collisions := s.CheckCollisions()
	//if len(collisions) > 0 {
	//	// Backtrack to 0 penetration
	//	s.Backtrack()
	//}

	// Apply static forces/inertia/impulses (only to objects that did not collide)
	s.applyForceFields(frameDelta)
	// s.applyDrag(frameDelta) // TODO

	// Apply impact forces/inertia/impulses to objects that collided
	//s.applyImpactForces(frameDelta)

	// Update object positions based on calculated final speeds
	s.updatePositions(frameDelta)

}
