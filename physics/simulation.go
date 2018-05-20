// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package physics

import (
	"time"
	"github.com/g3n/engine/physics/equation"
	"github.com/g3n/engine/physics/solver"
	"github.com/g3n/engine/physics/constraint"
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


	allowSleep   bool // Makes bodies go to sleep when they've been inactive
	contactEqs   []*equation.Contact // All the current contacts (instances of ContactEquation) in the world.
	frictionEqs  []*equation.Friction

	quatNormalizeSkip int // How often to normalize quaternions. Set to 0 for every step, 1 for every second etc..
	                      // A larger value increases performance.
	                      // If bodies tend to explode, set to a smaller value (zero to be sure nothing can go wrong).
	quatNormalizeFast bool // Set to true to use fast quaternion normalization. It is often enough accurate to use. If bodies tend to explode, set to false.

	time float32      // The wall-clock time since simulation start
	stepnumber int    // Number of timesteps taken since start
	default_dt float32 // Default and last timestep sizes

	accumulator float32 // Time accumulator for interpolation. See http://gafferongames.com/game-physics/fix-your-timestep/

	//broadphase IBroadphase // The broadphase algorithm to use. Default is NaiveBroadphase
	//narrowphase INarrowphase // The narrowphase algorithm to use
	solver solver.ISolver    // The solver algorithm to use. Default is GSSolver

	constraints       []*constraint.Constraint  // All constraints
	materials         []*Material               // All added materials
	cMaterials        []*ContactMaterial






	doProfiling      bool
}

// NewSimulation creates and returns a pointer to a new physics simulation.
func NewSimulation() *Simulation {

	s := new(Simulation)
	s.time = 0
	s.default_dt = 1/60

	//s.broadphase = NewNaiveBroadphase()
	//s.narrowphase = NewNarrowphase()
	s.solver = solver.NewGaussSeidel()

	return s
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
// @todo If the simulation has not yet started, why recrete and copy arrays for each body? Accumulate in dynamic arrays in this case.
// @todo Adding an array of bodies should be possible. This would save some loops too
func (s *Simulation) AddBody(body *Body) {

	// TODO only add if not already present
	s.bodies = append(s.bodies, body)

	//body.index = this.bodies.length
	//s.bodies.push(body)
	//body.simulation = s // TODO
	//body.initPosition.Copy(body.position)
	//body.initVelocity.Copy(body.velocity)
	//body.timeLastSleepy = s.time

	//if body instanceof Body { // TODO
	//	body.initAngularVelocity.Copy(body.angularVelocity)
	//	body.initQuaternion.Copy(body.quaternion)
	//}
	//
	//// TODO
	//s.collisionMatrix.setNumObjects(len(s.bodies))
	//s.addBodyEvent.body = body
	//s.idToBodyMap[body.id] = body
	//s.dispatchEvent(s.addBodyEvent)
}

// RemoveBody removes the specified body from the simulation.
// Returns true if found, false otherwise.
func (s *Simulation) RemoveBody(body *Body) bool {

	for pos, current := range s.bodies {
		if current == body {
			copy(s.bodies[pos:], s.bodies[pos+1:])
			s.bodies[len(s.bodies)-1] = nil
			s.bodies = s.bodies[:len(s.bodies)-1]

			body.simulation = nil

			// Recompute body indices (each body has a .index int property)
			for i:=0; i<len(s.bodies); i++ {
				s.bodies[i].index = i
			}

			// TODO
			//s.collisionMatrix.setNumObjects(len(s.bodies) - 1)
			//s.removeBodyEvent.body = body
			//delete s.idToBodyMap[body.id]
			//s.dispatchEvent(s.removeBodyEvent)

			return true
		}
	}

	return false
}

// Bodies returns the bodies under simulation.
func (s *Simulation) Bodies() []*Body{

	return s.bodies
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
		pos.Add(posDelta)
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

// ClearForces sets all body forces in the world to zero.
func (s *Simulation) ClearForces() {

	for i:=0; i < len(s.bodies); i++ {
		s.bodies[i].force.Set(0,0,0)
		s.bodies[i].torque.Set(0,0,0)
	}
}

// Add a constraint to the simulation.
func (s *Simulation) AddConstraint(c *constraint.Constraint) {

	s.constraints = append(s.constraints, c)
}

func (s *Simulation) RemoveConstraint(c *constraint.Constraint) {

	// TODO
}

func (s *Simulation) AddMaterial(mat *Material) {

	s.materials = append(s.materials, mat)
}

func (s *Simulation) RemoveMaterial(mat *Material) {

	// TODO
}

func (s *Simulation) AddContactMaterial(cmat *ContactMaterial) {

	s.cMaterials = append(s.cMaterials, cmat)

	// TODO add contactMaterial materials to contactMaterialTable
}


type ContactEvent struct {
	bodyA *Body
	bodyB *Body
	// shapeA
	// shapeB
}

const (
	BeginContactEvent = "physics.BeginContactEvent"
	EndContactEvent   = "physics.EndContactEvent"
)

func (s *Simulation) EmitContactEvents() {
	//TODO
}



func (s *Simulation) Solve() {

	//// Update solve mass for all bodies
	//if Neq > 0 {
	//	for i := 0; i < Nbodies; i++ {
	//		bodies[i].UpdateSolveMassProperties()
	//	}
	//}
	//
	//s.solver.Solve(frameDelta, len(bodies))

}


// Note - this method alters the solution arrays
func (s *Simulation) ApplySolution(solution *solver.Solution) {

	// Add results to velocity and angular velocity of bodies
	for i := 0; i < len(s.bodies); i++ {
		b := s.bodies[i]

		vDelta := solution.VelocityDeltas[i].Multiply(b.LinearFactor())
		b.AddToVelocity(vDelta)

		wDelta := solution.AngularVelocityDeltas[i].Multiply(b.AngularFactor())
		b.AddToAngularVelocity(wDelta)
	}
}
