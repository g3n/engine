// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package physics

import (
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/graphic"
)

// Body represents a physics-driven body.
type Body struct {
	//core.INode // TODO instead of embedding INode - embed Node and have a method SetNode ?

	*graphic.Graphic

	mass            float32        // Total mass
	invMass         float32
	invMassSolve    float32

	velocity        *math32.Vector3 // Linear velocity (World space velocity of the body.)
	initVelocity    *math32.Vector3 // Initial linear velocity (World space velocity of the body.)
	vLambda         *math32.Vector3

	angularMass     *math32.Matrix3 // Angular mass i.e. moment of inertia

	inertia              *math32.Vector3
	invInertia           *math32.Vector3
	invInertiaSolve      *math32.Vector3
	invInertiaWorld      *math32.Matrix3
	invInertiaWorldSolve *math32.Matrix3

	fixedRotation    bool  // Set to true if you don't want the body to rotate. Make sure to run .updateMassProperties() after changing this.

	angularVelocity     *math32.Vector3 // Angular velocity of the body, in world space. Think of the angular velocity as a vector, which the body rotates around. The length of this vector determines how fast (in radians per second) the body rotates.
	initAngularVelocity *math32.Vector3
	wLambda             *math32.Vector3


	force           *math32.Vector3 // Linear force on the body in world space.
	torque          *math32.Vector3 // World space rotational force on the body, around center of mass.

	position        *math32.Vector3 // World position of the center of gravity (World space position of the body.)
	prevPosition    *math32.Vector3 // Previous position
	interpPosition  *math32.Vector3 // Interpolated position of the body.
	initPosition    *math32.Vector3 // Initial position of the body.

	quaternion       *math32.Quaternion // World space orientation of the body.
	initQuaternion   *math32.Quaternion
	prevQuaternion   *math32.Quaternion
	interpQuaternion *math32.Quaternion // Interpolated orientation of the body.

	bodyType        BodyType
	sleepState      BodySleepState // Current sleep state.
	allowSleep      bool           // If true, the body will automatically fall to sleep.
	sleepSpeedLimit float32        // If the speed (the norm of the velocity) is smaller than this value, the body is considered sleepy.
	sleepTimeLimit  float32        // If the body has been sleepy for this sleepTimeLimit seconds, it is considered sleeping.
	timeLastSleepy  float32

	simulation             *Simulation // Reference to the simulation the body is living in\
	collisionFilterGroup   int
	collisionFilterMask    int
	collisionResponse      bool // Whether to produce contact forces when in contact with other bodies. Note that contacts will be generated, but they will be disabled.

	wakeUpAfterNarrowphase bool
	material               *Material

	linearDamping          float32
	angularDamping         float32

	linearFactor           *math32.Vector3 // Use this property to limit the motion along any world axis. (1,1,1) will allow motion along all axes while (0,0,0) allows none.
	angularFactor          *math32.Vector3 // Use this property to limit the rotational motion along any world axis. (1,1,1) will allow rotation along all axes while (0,0,0) allows none.

	//aabb            *AABB   // World space bounding box of the body and its shapes.
	aabbNeedsUpdate   bool    // Indicates if the AABB needs to be updated before use.
	boundingRadius    float32 // Total bounding radius of the Body including its shapes, relative to body.position.

	// shapes          []*Shape
	// shapeOffsets    []float32 // Position of each Shape in the body, given in local Body space.
	// shapeOrientations [] ?

	index int
}

// BodyType specifies how the body is affected during the simulation.
type BodyType int

const (
	// A static body does not move during simulation and behaves as if it has infinite mass.
	// Static bodies can be moved manually by setting the position of the body.
	// The velocity of a static body is always zero.
	// Static bodies do not collide with other static or kinematic bodies.
	Static       = BodyType(iota)

	// A kinematic body moves under simulation according to its velocity.
	// They do not respond to forces.
	// They can be moved manually, but normally a kinematic body is moved by setting its velocity.
	// A kinematic body behaves as if it has infinite mass.
	// Kinematic bodies do not collide with other static or kinematic bodies.
	Kinematic

	// A dynamic body is fully simulated.
	// Can be moved manually by the user, but normally they move according to forces.
	// A dynamic body can collide with all body types.
	// A dynamic body always has finite, non-zero mass.
	Dynamic
)

// BodyStatus specifies
type BodySleepState int

const (
	Awake = BodySleepState(iota)
	Sleepy
	Sleeping
)

// Events
const (
	SleepyEvent  = "physics.SleepyEvent"  // Dispatched after a body has gone in to the sleepy state.
	SleepEvent   = "physics.SleepEvent"   // Dispatched after a body has fallen asleep.
	WakeUpEvent  = "physics.WakeUpEvent"  // Dispatched after a sleeping body has woken up.
	CollideEvent = "physics.CollideEvent" // Dispatched after two bodies collide. This event is dispatched on each of the two bodies involved in the collision.
)


// NewBody creates and returns a pointer to a new RigidBody.
func NewBody(igraphic graphic.IGraphic) *Body {

	b := new(Body)
	b.Graphic = igraphic.GetGraphic()

	// TODO mass setter/getter
	b.mass = 1 // cannon.js default is 0
	if b.mass > 0 {
		b.invMass = 1.0/b.mass
	} else {
		b.invMass = 0
	}
	b.bodyType = Dynamic // TODO auto set to Static if mass == 0

	b.collisionFilterGroup = 1
	b.collisionFilterMask = -1

	pos := igraphic.GetNode().Position()
	b.position 			= math32.NewVector3(0,0,0).Copy(&pos)
	b.prevPosition 		= math32.NewVector3(0,0,0).Copy(&pos)
	b.interpPosition 	= math32.NewVector3(0,0,0).Copy(&pos)
	b.initPosition 		= math32.NewVector3(0,0,0).Copy(&pos)

	quat := igraphic.GetNode().Quaternion()
	b.quaternion 		= math32.NewQuaternion(0,0,0,1).Copy(&quat)
	b.prevQuaternion 	= math32.NewQuaternion(0,0,0,1).Copy(&quat)
	b.interpQuaternion 	= math32.NewQuaternion(0,0,0,1).Copy(&quat)
	b.initQuaternion 	= math32.NewQuaternion(0,0,0,1).Copy(&quat)

	b.velocity = math32.NewVector3(0,0,0) // TODO copy options.velocity
	b.initVelocity = math32.NewVector3(0,0,0) // don't copy

	b.angularVelocity = math32.NewVector3(0,0,0)
	b.initAngularVelocity = math32.NewVector3(0,0,0)

	b.vLambda = math32.NewVector3(0,0,0)
	b.wLambda = math32.NewVector3(0,0,0)

	b.linearDamping = 0.01
	b.angularDamping = 0.01

	b.linearFactor = math32.NewVector3(1,1,1)
	b.angularFactor = math32.NewVector3(1,1,1)

	b.allowSleep = true
	b.sleepState = Awake
	b.sleepSpeedLimit = 0.1
	b.sleepTimeLimit = 1
	b.timeLastSleepy = 0

	b.force = math32.NewVector3(0,0,0)
	b.torque = math32.NewVector3(0,0,0)

	b.wakeUpAfterNarrowphase = false

	b.UpdateMassProperties()

	return b
}

func (b *Body) Position() math32.Vector3 {

	return *b.position
}

func (b *Body) SetVelocity(vel *math32.Vector3) {

	b.velocity = vel
}

func (b *Body) AddToVelocity(delta *math32.Vector3) {

	b.velocity.Add(delta)
}

func (b *Body) Velocity() math32.Vector3 {

	return *b.velocity
}

func (b *Body) SetAngularVelocity(vel *math32.Vector3) {

	b.angularVelocity = vel
}

func (b *Body) AddToAngularVelocity(delta *math32.Vector3) {

	b.angularVelocity.Add(delta)
}

func (b *Body) AngularVelocity() math32.Vector3 {

	return *b.angularVelocity
}

func (b *Body) Force() math32.Vector3 {

	return *b.force
}

func (b *Body) Torque() math32.Vector3 {

	return *b.torque
}

func (b *Body) SetVlambda(vLambda *math32.Vector3) {

	b.vLambda = vLambda
}

func (b *Body) AddToVlambda(delta *math32.Vector3) {

	b.vLambda.Add(delta)
}

func (b *Body) Vlambda() math32.Vector3 {

	return *b.vLambda
}

func (b *Body) SetWlambda(wLambda *math32.Vector3) {

	b.wLambda = wLambda
}

func (b *Body) AddToWlambda(delta *math32.Vector3) {

	b.wLambda.Add(delta)
}

func (b *Body) Wlambda() math32.Vector3 {

	return *b.wLambda
}

func (b *Body) InvMassSolve() float32 {

	return b.invMassSolve
}

func (b *Body) InvInertiaWorldSolve() *math32.Matrix3 {

	return b.invInertiaWorldSolve
}

func (b *Body) Quaternion() *math32.Quaternion {

	return b.quaternion
}

func (b *Body) LinearFactor() *math32.Vector3 {

	return b.linearFactor
}

func (b *Body) AngularFactor() *math32.Vector3 {

	return b.angularFactor
}

// WakeUp wakes the body up.
func (b *Body) WakeUp() {

	state := b.sleepState
	b.sleepState = Awake
	b.wakeUpAfterNarrowphase = false
	if state == Sleeping {
		b.Dispatch(WakeUpEvent, nil)
	}
}

// Sleep forces the body to sleep.
func (b *Body) Sleep() {

	b.sleepState = Sleeping
	b.velocity.Set(0,0,0)
	b.angularVelocity.Set(0,0,0)
	b.wakeUpAfterNarrowphase = false
}

// Called every timestep to update internal sleep timer and change sleep state if needed.
// time: The world time in seconds
func (b *Body) SleepTick(time float32) {

	if b.allowSleep {
		speedSquared := b.velocity.LengthSq() + b.angularVelocity.LengthSq()
		speedLimitSquared := math32.Pow(b.sleepSpeedLimit,2)
		if b.sleepState == Awake && speedSquared < speedLimitSquared {
			b.sleepState = Sleepy
			b.timeLastSleepy = time
			b.Dispatch(SleepyEvent, nil)
		} else if b.sleepState == Sleepy && speedSquared > speedLimitSquared {
			b.WakeUp() // Wake up
		} else if b.sleepState == Sleepy && (time - b.timeLastSleepy ) > b.sleepTimeLimit {
			b.Sleep() // Sleeping
			b.Dispatch(SleepEvent, nil)
		}
	}
}

// TODO maybe return vector instead of pointer in below methods

// PointToLocal converts a world point to local body frame. TODO maybe move to Node
func (b *Body) PointToLocal(worldPoint *math32.Vector3) math32.Vector3 {

	result := math32.NewVector3(0,0,0).SubVectors(worldPoint, b.position)
	conj := b.quaternion.Conjugate()
	result.ApplyQuaternion(conj)

	return *result
}

// VectorToLocal converts a world vector to local body frame. TODO maybe move to Node
func (b *Body) VectorToLocal(worldVector *math32.Vector3) math32.Vector3 {

	result := math32.NewVector3(0,0,0).Copy(worldVector)
	conj := b.quaternion.Conjugate()
	result.ApplyQuaternion(conj)

	return *result
}

// PointToWorld converts a local point to world frame. TODO maybe move to Node
func (b *Body) PointToWorld(localPoint *math32.Vector3) math32.Vector3 {

	result := math32.NewVector3(0,0,0).Copy(localPoint)
	result.ApplyQuaternion(b.quaternion)
	result.Add(b.position)

	return *result
}

// VectorToWorld converts a local vector to world frame. TODO maybe move to Node
func (b *Body) VectorToWorld(localVector *math32.Vector3) math32.Vector3 {

	result := math32.NewVector3(0,0,0).Copy(localVector)
	result.ApplyQuaternion(b.quaternion)

	return *result
}



func (b *Body) ComputeAABB() {
	// TODO
}


// UpdateSolveMassProperties
// If the body is sleeping, it should be immovable / have infinite mass during solve. We solve it by having a separate "solve mass".
func (b *Body) UpdateSolveMassProperties() {

	if b.sleepState == Sleeping || b.bodyType == Kinematic {
		b.invMassSolve = 0
		b.invInertiaSolve.Zero()
		b.invInertiaWorldSolve.Zero()
	} else {
		b.invMassSolve = b.invMass
		b.invInertiaSolve.Copy(b.invInertia)
		b.invInertiaWorldSolve.Copy(b.invInertiaWorld)
	}
}

// UpdateMassProperties // TODO
// Should be called whenever you change the body shape or mass.
func (b *Body) UpdateMassProperties() {

	// TODO getter of invMass ?
	if b.mass > 0 {
		b.invMass = 1.0/b.mass
	} else {
		b.invMass = 0
	}

	// Approximate with AABB box
	b.ComputeAABB()
	//halfExtents := math32.NewVector3(0,0,0).Set( // TODO
	//	(b.aabb.upperBound.x-b.aabb.lowerBound.x) / 2,
	//	(b.aabb.upperBound.y-b.aabb.lowerBound.y) / 2,
	//	(b.aabb.upperBound.z-b.aabb.lowerBound.z) / 2,
	//)
	//Box.CalculateInertia(halfExtents, b.mass, b.inertia) // TODO

	if b.fixedRotation {
		b.invInertia.Zero()
	} else {
		if b.inertia.X > 0 {
			b.invInertia.SetX(1/b.inertia.X)
		} else {
			b.invInertia.SetX(0)
		}
		if b.inertia.Y > 0 {
			b.invInertia.SetY(1/b.inertia.Y)
		} else {
			b.invInertia.SetY(0)
		}
		if b.inertia.Z > 0 {
			b.invInertia.SetZ(1/b.inertia.Z)
		} else {
			b.invInertia.SetZ(0)
		}
	}

	b.UpdateInertiaWorld(true)
}

// Update .inertiaWorld and .invInertiaWorld
func (b *Body) UpdateInertiaWorld(force bool) {

    I := b.invInertia
	// If angular mass M = s*I, where I is identity and s a scalar, then
	//    R*M*R' = R*(s*I)*R' = s*R*I*R' = s*R*R' = s*I = M
	// where R is the rotation matrix.
	// In other words, we don't have to do the transformation if all diagonal entries are equal.
    if I.X != I.Y || I.Y != I.Z || force {
    	//
    	// AngularMassWorld^(-1) = Rotation * AngularMassBody^(-1) * Rotation^(T)
    	//          3x3              3x3            3x3                  3x3
    	//
    	// Since AngularMassBodyTensor^(-1) is diagonal, then Rotation*AngularMassBodyTensor^(-1) is
    	// just scaling the columns of AngularMassBodyTensor by the diagonal components.
    	//
        m1 := math32.NewMatrix3()
        m2 := math32.NewMatrix3()

        m1.MakeRotationFromQuaternion(b.quaternion)
		m2.Copy(m1).Transpose()
        m1.ScaleColumns(I)

		b.invInertiaWorld.MultiplyMatrices(m1, m2)
    }
}

// Apply force to a world point.
// This could for example be a point on the Body surface.
// Applying force this way will add to Body.force and Body.torque.
// relativePoint: A point relative to the center of mass to apply the force on.
func (b *Body) ApplyForce(force, relativePoint *math32.Vector3) {

	if b.bodyType != Dynamic { // Needed?
		return
	}

	// Compute produced rotational force
	rotForce := math32.NewVector3(0,0,0)
	rotForce.CrossVectors(relativePoint, force)

	// Add linear force
	b.force.Add(force) // TODO shouldn't rotational momentum be subtracted from linear momentum?

	// Add rotational force
	b.torque.Add(rotForce)
}

// Apply force to a local point in the body.
// force: The force vector to apply, defined locally in the body frame.
// localPoint: A local point in the body to apply the force on.
func (b *Body) ApplyLocalForce(localForce, localPoint *math32.Vector3)  {

	if b.bodyType != Dynamic {
		return
	}

	// Transform the force vector to world space
	worldForce := b.VectorToWorld(localForce)
	relativePointWorld := b.VectorToWorld(localPoint)

	b.ApplyForce(&worldForce, &relativePointWorld)
}

// Apply impulse to a world point.
// This could for example be a point on the Body surface.
// An impulse is a force added to a body during a short period of time (impulse = force * time).
// Impulses will be added to Body.velocity and Body.angularVelocity.
// impulse: The amount of impulse to add.
// relativePoint: A point relative to the center of mass to apply the force on.
func (b *Body) ApplyImpulse(impulse, relativePoint *math32.Vector3) {

	if b.bodyType != Dynamic {
		return
	}

    // Compute point position relative to the body center
    r := relativePoint

    // Compute produced central impulse velocity
    velo := math32.NewVector3(0,0,0).Copy(impulse)
    velo.MultiplyScalar(b.invMass)

    // Add linear impulse
    b.velocity.Add(velo)

    // Compute produced rotational impulse velocity
	rotVelo := math32.NewVector3(0,0,0).CrossVectors(r, impulse)
	rotVelo.ApplyMatrix3(b.invInertiaWorld)

    // Add rotational Impulse
    b.angularVelocity.Add(rotVelo)
}

// Apply locally-defined impulse to a local point in the body.
// force: The force vector to apply, defined locally in the body frame.
// localPoint: A local point in the body to apply the force on.
func (b *Body) ApplyLocalImpulse(localImpulse, localPoint *math32.Vector3) {

	if b.bodyType != Dynamic {
		return
	}

	// Transform the force vector to world space
	worldImpulse := b.VectorToWorld(localImpulse)
	relativePointWorld := b.VectorToWorld(localPoint)

	b.ApplyImpulse(&worldImpulse, &relativePointWorld)
}

// Get world velocity of a point in the body.
func (b *Body) GetVelocityAtWorldPoint(worldPoint *math32.Vector3) *math32.Vector3 {

	r := math32.NewVector3(0,0,0)
	r.SubVectors(worldPoint, b.position)
	r.CrossVectors(b.angularVelocity, r)
	r.Add(b.velocity)

	return r
}

// Move the body forward in time.
// dt: Time step
// quatNormalize: Set to true to normalize the body quaternion
// quatNormalizeFast: If the quaternion should be normalized using "fast" quaternion normalization
func (b *Body) Integrate(dt float32, quatNormalize, quatNormalizeFast bool) {


    // Save previous position and rotation
    b.prevPosition.Copy(b.position)
    b.prevQuaternion.Copy(b.quaternion)

    // If static or sleeping - skip
    if !(b.bodyType == Dynamic || b.bodyType == Kinematic) || b.sleepState == Sleeping {
        return
    }

    // Integrate force over mass (acceleration) to obtain estimate for instantaneous velocities
    iMdt := b.invMass * dt
    b.velocity.X += b.force.X * iMdt * b.linearFactor.X
    b.velocity.Y += b.force.Y * iMdt * b.linearFactor.Y
    b.velocity.Z += b.force.Z * iMdt * b.linearFactor.Z

	// Integrate inverse angular mass times torque to obtain estimate for instantaneous angular velocities
    e := b.invInertiaWorld
    tx := b.torque.X * b.angularFactor.X
    ty := b.torque.Y * b.angularFactor.Y
    tz := b.torque.Z * b.angularFactor.Z
    b.angularVelocity.X += dt * (e[0]*tx + e[3]*ty + e[6]*tz)
    b.angularVelocity.Y += dt * (e[1]*tx + e[4]*ty + e[7]*tz)
    b.angularVelocity.Z += dt * (e[2]*tx + e[5]*ty + e[8]*tz)

	// Integrate velocity to obtain estimate for position
    b.position.X += b.velocity.X * dt
    b.position.Y += b.velocity.Y * dt
    b.position.Z += b.velocity.Z * dt

	// Integrate angular velocity to obtain estimate for rotation
	ax := b.angularVelocity.X * b.angularFactor.X
	ay := b.angularVelocity.Y * b.angularFactor.Y
	az := b.angularVelocity.Z * b.angularFactor.Z
	bx := b.quaternion.X
	by := b.quaternion.Y
	bz := b.quaternion.Z
	bw := b.quaternion.W
	halfDt := dt * 0.5
	b.quaternion.X += halfDt * (ax * bw + ay * bz - az * by)
	b.quaternion.Y += halfDt * (ay * bw + az * bx - ax * bz)
	b.quaternion.X += halfDt * (az * bw + ax * by - ay * bx)
	b.quaternion.W += halfDt * (- ax * bx - ay * by - az * bz)

	// Normalize quaternion
    if quatNormalize {
       if quatNormalizeFast {
			b.quaternion.NormalizeFast()
       } else {
			b.quaternion.Normalize()
       }
    }

    b.aabbNeedsUpdate = true  // TODO

    // Update world inertia
    b.UpdateInertiaWorld(false)
}
