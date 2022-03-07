// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package physics

import (
	"github.com/g3n/engine/experimental/physics/equation"
	"github.com/g3n/engine/experimental/physics/solver"
	"github.com/g3n/engine/experimental/physics/constraint"
	"github.com/g3n/engine/experimental/collision"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/experimental/physics/object"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/experimental/collision/shape"
)

// Simulation represents a physics simulation.
type Simulation struct {
	scene       *core.Node
	forceFields []ForceField

	// Bodies under simulation
	bodies      []*object.Body  // Slice of bodies. May contain nil values.
	nilBodies   []int           // Array keeps track of which indices of the 'bodies' array are nil

	// Collision tracking
	collisionMatrix     collision.Matrix // Boolean triangular matrix indicating which pairs of bodies are colliding
	prevCollisionMatrix collision.Matrix // CollisionMatrix from the previous step.

	allowSleep  bool                // Makes bodies go to sleep when they've been inactive
	paused bool

	quatNormalizeSkip int  // How often to normalize quaternions. Set to 0 for every step, 1 for every second etc..
	                       // A larger value increases performance.
	                       // If bodies tend to explode, set to a smaller value (zero to be sure nothing can go wrong).
	quatNormalizeFast bool // Set to true to use fast quaternion normalization. It is often enough accurate to use. If bodies tend to explode, set to false.

	time float32      // The wall-clock time since simulation start
	stepnumber int    // Number of timesteps taken since start
	default_dt float32 // Default and last timestep sizes
	dt float32      // Currently / last used timestep. Is set to -1 if not available. This value is updated before each internal step, which means that it is "fresh" inside event callbacks.

	accumulator float32 // Time accumulator for interpolation. See http://gafferongames.com/game-physics/fix-your-timestep/

	broadphase  *Broadphase    // The broadphase algorithm to use, default is NaiveBroadphase
	narrowphase *Narrowphase   // The narrowphase algorithm to use
	solver      solver.ISolver // The solver algorithm to use, default is Gauss-Seidel

	constraints       []constraint.IConstraint  // All constraints

	materials         []*Material               // All added materials
	cMaterials        []*ContactMaterial

	//contactMaterialTable map[intPair]*ContactMaterial // Used to look up a ContactMaterial given two instances of Material.
	//defaultMaterial *Material
	defaultContactMaterial *ContactMaterial

	doProfiling      bool
}

// NewSimulation creates and returns a pointer to a new physics simulation.
func NewSimulation(scene *core.Node) *Simulation {

	s := new(Simulation)
	s.time = 0
	s.dt = -1
	s.default_dt = 1/60
	s.scene = scene

	// Set up broadphase, narrowphase, and solver
	s.broadphase = NewBroadphase()
	s.narrowphase = NewNarrowphase(s)
	s.solver = solver.NewGaussSeidel()

	s.collisionMatrix = collision.NewMatrix()
	s.prevCollisionMatrix = collision.NewMatrix()

	//s.contactMaterialTable = make(map[intPair]*ContactMaterial)
	//s.defaultMaterial = NewMaterial
	s.defaultContactMaterial = NewContactMaterial()

	return s
}

func (s *Simulation) Scene() *core.Node {

	return s.scene
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
func (s *Simulation) AddBody(body *object.Body, name string) {

	// Do nothing if body is already present
	for _, existingBody := range s.bodies {
		if existingBody == body {
			return // Do nothing
		}
	}

	// If there are any open/nil spots in the body slice - add the body to one of them
	// Else, just append to the end of the slice. Either way compute the body's index in the slice.
	var idx int
	nilLen := len(s.nilBodies)
	if nilLen > 0 {
		idx = s.nilBodies[nilLen]
		s.nilBodies = s.nilBodies[0:nilLen-1]
	} else {
		idx = len(s.bodies)
		s.bodies = append(s.bodies, body)

		// Initialize collision matrix values up to the current index (and set the colliding flag to false)
		s.collisionMatrix.Set(idx, idx, false)
		s.prevCollisionMatrix.Set(idx, idx, false)
	}

	body.SetIndex(idx)
	body.SetName(name)

	// TODO dispatch add-body event
	//s.Dispatch(AddBodyEvent, BodyEvent{body})
}

// RemoveBody removes the specified body from the simulation.
// Returns true if found, false otherwise.
func (s *Simulation) RemoveBody(body *object.Body) bool {

	for idx, current := range s.bodies {
		if current == body {
			s.bodies[idx] = nil
			// TODO dispatch remove-body event
			//s.Dispatch(AddBodyEvent, BodyEvent{body})
			return true
		}
	}
	return false
}

// Clean removes nil bodies from the bodies array, recalculates the body indices and updates the collision matrix.
//func (s *Simulation) Clean() {
//
//	// TODO Remove nil bodies from array
//	//copy(s.bodies[pos:], s.bodies[pos+1:])
//	//s.bodies[len(s.bodies)-1] = nil
//	//s.bodies = s.bodies[:len(s.bodies)-1]
//
//	// Recompute body indices (each body has a .index int property)
//	for i:=0; i<len(s.bodies); i++ {
//		s.bodies[i].SetIndex(i)
//	}
//
//	// TODO Update collision matrix
//
//}

// Bodies returns the slice of bodies under simulation.
// The slice may contain nil values!
func (s *Simulation) Bodies() []*object.Body{

	return s.bodies
}

func (s *Simulation) Step(frameDelta float32) {

	s.StepPlus(frameDelta, 0, 10)
}


// Step steps the simulation.
// maxSubSteps should be 10 by default
func (s *Simulation) StepPlus(frameDelta float32, timeSinceLastCalled float32, maxSubSteps int) {

	if s.paused {
		return
	}

	dt := frameDelta//float32(frameDelta.Seconds())

    //if timeSinceLastCalled == 0 { // Fixed, simple stepping

        s.internalStep(dt)

        // Increment time
        //s.time += dt

    //} else {
	//
    //    s.accumulator += timeSinceLastCalled
    //    var substeps = 0
    //    for s.accumulator >= dt && substeps < maxSubSteps {
    //        // Do fixed steps to catch up
    //        s.internalStep(dt)
    //        s.accumulator -= dt
    //        substeps++
    //    }
	//
    //    var t = (s.accumulator % dt) / dt
    //    for j := 0; j < len(s.bodies); j++ {
    //        var b = s.bodies[j]
    //        b.previousPosition.lerp(b.position, t, b.interpolatedPosition)
    //        b.previousQuaternion.slerp(b.quaternion, t, b.interpolatedQuaternion)
    //        b.previousQuaternion.normalize()
    //    }
    //    s.time += timeSinceLastCalled
    //}

}

// SetPaused sets the paused state of the simulation.
func (s *Simulation) SetPaused(state bool) {

	s.paused = state
}

// Paused returns the paused state of the simulation.
func (s *Simulation) Paused() bool {

	return s.paused
}

// ClearForces sets all body forces in the world to zero.
func (s *Simulation) ClearForces() {

	for i:=0; i < len(s.bodies); i++ {
		s.bodies[i].ClearForces()
	}
}

// AddConstraint adds a constraint to the simulation.
func (s *Simulation) AddConstraint(c constraint.IConstraint) {

	s.constraints = append(s.constraints, c)
}

func (s *Simulation) RemoveConstraint(c constraint.IConstraint) {

	// TODO
}

func (s *Simulation) AddMaterial(mat *Material) {

	s.materials = append(s.materials, mat)
}

func (s *Simulation) RemoveMaterial(mat *Material) {

	// TODO
}

// Adds a contact material to the simulation
func (s *Simulation) AddContactMaterial(cmat *ContactMaterial) {

	s.cMaterials = append(s.cMaterials, cmat)

	// TODO add contactMaterial materials to contactMaterialTable
	// s.contactMaterialTable.set(cmat.materials[0].id, cmat.materials[1].id, cmat)
}

// GetContactMaterial returns the contact material between the specified bodies.
func (s *Simulation) GetContactMaterial(bodyA, bodyB *object.Body) *ContactMaterial {

	var cm *ContactMaterial
	// TODO
	//if bodyA.material != nil && bodyB.material != nil {
	//	cm = s.contactMaterialTable.get(bodyA.material.id, bodyB.material.id)
	//	if cm == nil {
	//		cm = s.defaultContactMaterial
	//	}
	//} else {
	cm = s.defaultContactMaterial
	//}
	return cm
}


// Events =====================

type CollideEvent struct {
	body      *object.Body
	contactEq *equation.Contact
}

// TODO AddBodyEvent, RemoveBodyEvent
type ContactEvent struct {
	bodyA *object.Body
	bodyB *object.Body
}

const (
	BeginContactEvent = "physics.BeginContactEvent"
	EndContactEvent   = "physics.EndContactEvent"
	CollisionEv       = "physics.Collision"
)

// ===========================


// ApplySolution applies the specified solution to the bodies under simulation.
// The solution is a set of linear and angular velocity deltas for each body.
// This method alters the solution arrays.
func (s *Simulation) ApplySolution(sol *solver.Solution) {

	// Add results to velocity and angular velocity of bodies
	for i := 0; i < len(s.bodies); i++ {
		s.bodies[i].ApplyVelocityDeltas(&sol.VelocityDeltas[i], &sol.AngularVelocityDeltas[i])
	}
}

// Store old collision state info
func (s *Simulation) collisionMatrixTick() {

	s.prevCollisionMatrix = s.collisionMatrix
	s.collisionMatrix = collision.NewMatrix()

	lb := len(s.bodies)
	s.collisionMatrix.Set(lb, lb, false)

	// TODO verify that the matrices are indeed different
	//if s.prevCollisionMatrix == s.collisionMatrix {
	//	log.Error("SAME")
	//}

	// TODO
	//s.bodyOverlapKeeper.tick()
	//s.shapeOverlapKeeper.tick()
}

func (s *Simulation) uniqueBodiesFromPairs(pairs []CollisionPair) []*object.Body {

	bodiesUndergoingNarrowphase := make([]*object.Body, 0) // array of indices of s.bodies
	for _, pair := range pairs {
		alreadyAddedA := false
		alreadyAddedB := false
		for _, body := range bodiesUndergoingNarrowphase {
			if pair.BodyA == body {
				alreadyAddedA = true
				if alreadyAddedB {
					break
				}
			}
			if pair.BodyB == body {
				alreadyAddedB = true
				if alreadyAddedA {
					break
				}
			}
		}
		if !alreadyAddedA {
			bodiesUndergoingNarrowphase = append(bodiesUndergoingNarrowphase, pair.BodyA)
		}
		if !alreadyAddedB {
			bodiesUndergoingNarrowphase = append(bodiesUndergoingNarrowphase, pair.BodyB)
		}
	}

	return bodiesUndergoingNarrowphase
}

// TODO read https://gafferongames.com/post/fix_your_timestep/
func (s *Simulation) internalStep(dt float32) {

	s.dt = dt

	// Apply force fields (only to dynamic bodies
	for _, b := range s.bodies {
		if b.BodyType() == object.Dynamic {
			for _, ff := range s.forceFields {
				pos := b.Position()
				force := ff.ForceAt(&pos)
				b.ApplyForceField(&force)
			}
		}
	}

    // Find pairs of bodies that are potentially colliding (broadphase)
	pairs := s.broadphase.FindCollisionPairs(s.bodies)

	// Remove some pairs before proceeding to narrowphase based on constraints' colConn property
	// which specifies if constrained bodies should collide with one another
    s.prunePairs(pairs) // TODO review/implement

	// Precompute world normals/edges only for convex bodies that will undergo narrowphase
	for _, body := range s.uniqueBodiesFromPairs(pairs) {
		if ch, ok := body.Shape().(*shape.ConvexHull); ok{
			ch.ComputeWorldFaceNormalsAndUniqueEdges(body.Quaternion())
		}
	}

	// Switch collision matrices (to keep track of which collisions started/ended)
    s.collisionMatrixTick()

    // Resolve collisions and generate contact and friction equations
	contactEqs, frictionEqs := s.narrowphase.GenerateEquations(pairs)

	// Add all friction equations to solver
	for i := 0; i < len(frictionEqs); i++ {
		s.solver.AddEquation(frictionEqs[i])
	}

	// Add all contact equations to solver (and update some things)
	for i := 0; i < len(contactEqs); i++ {
		s.solver.AddEquation(contactEqs[i])
		s.updateSleepAndCollisionMatrix(contactEqs[i])
	}

    // Add all equations from user-added constraints to the solver
	userAddedEquations := 0
    for i := 0; i < len(s.constraints); i++ {
		s.constraints[i].Update()
        eqs := s.constraints[i].Equations()
        for j := 0; j < len(eqs); j++ {
			userAddedEquations++
            s.solver.AddEquation(eqs[j])
        }
    }

	// Emit events TODO implement
	s.emitContactEvents()
	// Wake up bodies
	// TODO why not wake bodies up inside s.updateSleepAndCollisionMatrix when setting the WakeUpAfterNarrowphase flag?
	// Maybe there we are only looking at bodies that belong to current contact equations...
	// and need to wake up all marked bodies
	for i := 0; i < len(s.bodies); i++ {
		bi := s.bodies[i]
		if bi != nil && bi.WakeUpAfterNarrowphase() {
			bi.WakeUp()
			bi.SetWakeUpAfterNarrowphase(false)
		}
	}

	// If we have any equations to solve
	if len(frictionEqs) + len(contactEqs) + userAddedEquations > 0 {
		// Update effective mass for all bodies
		for i := 0; i < len(s.bodies); i++ {
			s.bodies[i].UpdateEffectiveMassProperties()
		}
		// Solve the constrained system
		solution := s.solver.Solve(dt, len(s.bodies))
		// Apply linear and angular velocity deltas to bodies
		s.ApplySolution(solution)
		// Clear all equations added to the solver
		s.solver.ClearEquations()
	}

    // Apply damping (only to dynamic bodies)
    // See http://code.google.com/p/bullet/issues/detail?id=74 for details
    for _, body := range s.bodies {
        if body != nil && body.BodyType() == object.Dynamic {
			body.ApplyDamping(dt)
        }
    }

    // TODO s.Dispatch(World_step_preStepEvent)

	// Integrate the forces into velocities and the velocities into position deltas for all bodies
    // TODO future: quatNormalize := s.stepnumber % (s.quatNormalizeSkip + 1) == 0
    for _, body := range s.bodies {
		if body != nil {
			body.Integrate(dt, true, s.quatNormalizeFast)
		}
    }
    s.ClearForces()

    // TODO s.broadphase.dirty = true ?

    // Update world time
    s.time += dt
    s.stepnumber += 1

    // TODO s.Dispatch(World_step_postStepEvent)

    // Sleeping update
    if s.allowSleep {
        for i := 0; i < len(s.bodies); i++ {
            s.bodies[i].SleepTick(s.time)
        }
    }

}

// TODO - REVIEW THIS
func (s *Simulation) prunePairs(pairs []CollisionPair) []CollisionPair {

	// TODO There is probably a bug here when the same body can have multiple constraints and appear in multiple pairs

	//// Remove constrained pairs with collideConnected == false
	//pairIdxsToRemove := []int
	//for i := 0; i < len(s.constraints); i++ {
	//	c := s.constraints[i]
	//	cBodyA := s.bodies[c.BodyA().Index()]
	//	cBodyB := s.bodies[c.BodyB().Index()]
	//	if !c.CollideConnected() {
	//		for i := range pairs {
	//			if (pairs[i].BodyA == cBodyA && pairs[i].BodyB == cBodyB) ||
	//				(pairs[i].BodyA == cBodyB && pairs[i].BodyB == cBodyA) {
	//				pairIdxsToRemove = append(pairIdxsToRemove, i)
	//
	//			}
	//		}
	//	}
	//}
	//
	//// Remove pairs
	////var prunedPairs []CollisionPair
	//for i := range pairs {
	//	for _, idx := range pairIdxsToRemove {
	//		copy(pairs[i:], pairs[i+1:])
	//		//pairs[len(pairs)-1] = nil
	//		pairs = pairs[:len(pairs)-1]
	//	}
	//}

	return pairs
}

// generateContacts
func (s *Simulation) updateSleepAndCollisionMatrix(contactEq *equation.Contact) {

	// Get current collision indices
	bodyA := s.bodies[contactEq.BodyA().Index()]
	bodyB := s.bodies[contactEq.BodyB().Index()]

	// TODO future: update equations with physical material properties

	if bodyA.AllowSleep() && bodyA.BodyType() == object.Dynamic && bodyA.SleepState() == object.Sleeping && bodyB.SleepState() == object.Awake && bodyB.BodyType() != object.Static {
		velocityB := bodyB.Velocity()
		angularVelocityB := bodyB.AngularVelocity()
		speedSquaredB := velocityB.LengthSq() + angularVelocityB.LengthSq()
		speedLimitSquaredB := math32.Pow(bodyB.SleepSpeedLimit(), 2)
		if speedSquaredB >= speedLimitSquaredB*2 {
			bodyA.SetWakeUpAfterNarrowphase(true)
		}
	}

	if bodyB.AllowSleep() && bodyB.BodyType() == object.Dynamic && bodyB.SleepState() == object.Sleeping && bodyA.SleepState() == object.Awake && bodyA.BodyType() != object.Static {
		velocityA := bodyA.Velocity()
		angularVelocityA := bodyA.AngularVelocity()
		speedSquaredA := velocityA.LengthSq() + angularVelocityA.LengthSq()
		speedLimitSquaredA := math32.Pow(bodyA.SleepSpeedLimit(), 2)
		if speedSquaredA >= speedLimitSquaredA*2 {
			bodyB.SetWakeUpAfterNarrowphase(true)
		}
	}

	// Now we know that i and j are in contact. Set collision matrix state
	s.collisionMatrix.Set(bodyA.Index(), bodyB.Index(), true)

	if s.prevCollisionMatrix.Get(bodyA.Index(), bodyB.Index()) == false {
		// First contact!
		bodyA.Dispatch(CollisionEv, &CollideEvent{bodyB, contactEq})
		bodyB.Dispatch(CollisionEv, &CollideEvent{bodyA, contactEq})
	}

	// TODO this is only for events
	//s.bodyOverlapKeeper.set(bodyA.id, bodyB.id)
	//s.shapeOverlapKeeper.set(si.id, sj.id)

}

func (s *Simulation) emitContactEvents() {
	//TODO
}