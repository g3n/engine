// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package object

import (
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/experimental/collision/shape"
)

// Body represents a physics-driven body.
type Body struct {
	*graphic.Graphic // TODO future - embed core.Node instead and calculate properties recursively

	material             *material.Material   // Physics material specifying friction and restitution
	index int
	name string

	// Mass properties
	mass       float32 // Total mass
	invMass    float32
	invMassEff float32 // Effective inverse mass

	// Rotational inertia and related properties
	rotInertia            *math32.Matrix3 // Angular mass i.e. moment of inertia in local coordinates
	invRotInertia         *math32.Matrix3 // Inverse rotational inertia in local coordinates
	invRotInertiaEff      *math32.Matrix3 // Effective inverse rotational inertia in local coordinates
	invRotInertiaWorld    *math32.Matrix3 // Inverse rotational inertia in world coordinates
	invRotInertiaWorldEff *math32.Matrix3 // Effective rotational inertia in world coordinates
	fixedRotation         bool            // Set to true if you don't want the body to rotate. Make sure to run .updateMassProperties() after changing this.

	// Position
	position       *math32.Vector3 // World position of the center of gravity (World space position of the body.)
	initPosition   *math32.Vector3 // Initial position of the body.
	prevPosition   *math32.Vector3 // Previous position
	interpPosition *math32.Vector3 // Interpolated position of the body.

	// Rotation
	quaternion       *math32.Quaternion // World space orientation of the body.
	initQuaternion   *math32.Quaternion
	prevQuaternion   *math32.Quaternion
	interpQuaternion *math32.Quaternion // Interpolated orientation of the body.

	// Linear and angular velocities
	velocity            *math32.Vector3 // Linear velocity (World space velocity of the body.)
	initVelocity        *math32.Vector3 // Initial linear velocity (World space velocity of the body.)
	angularVelocity     *math32.Vector3 // Angular velocity of the body, in world space. Think of the angular velocity as a vector, which the body rotates around. The length of this vector determines how fast (in radians per second) the body rotates.
	initAngularVelocity *math32.Vector3

	// Force and torque
	force  *math32.Vector3 // Linear force on the body in world space.
	torque *math32.Vector3 // World space rotational force on the body, around center of mass.

	// Damping and factors
	linearDamping  float32
	angularDamping float32
	linearFactor   *math32.Vector3 // Use this property to limit the motion along any world axis. (1,1,1) will allow motion along all axes while (0,0,0) allows none.
	angularFactor  *math32.Vector3 // Use this property to limit the rotational motion along any world axis. (1,1,1) will allow rotation along all axes while (0,0,0) allows none.

	// Body type and sleep settings
	bodyType        BodyType
	sleepState      BodySleepState // Current sleep state.
	allowSleep      bool           // If true, the body will automatically fall to sleep.
	sleepSpeedLimit float32        // If the speed (the norm of the velocity) is smaller than this value, the body is considered sleepy.
	sleepTimeLimit  float32        // If the body has been sleepy for this sleepTimeLimit seconds, it is considered sleeping.
	timeLastSleepy  float32
	wakeUpAfterNarrowphase bool

	// Collision settings
	colFilterGroup int  // Collision filter group
	colFilterMask  int  // Collision filter mask
	colResponse    bool // Whether to produce contact forces when in contact with other bodies. Note that contacts will be generated, but they will be disabled.

	aabb            *math32.Box3 // World space bounding box of the body and its shapes.
	aabbNeedsUpdate bool         // Indicates if the AABB needs to be updated before use.
	boundingRadius  float32      // Total bounding radius of the body (TODO including its shapes, relative to body.position.)

	// Cached geometry properties
	faces            [][3]math32.Vector3

	faceNormals      []math32.Vector3
	worldFaceNormals []math32.Vector3

	uniqueEdges      []math32.Vector3
	worldUniqueEdges []math32.Vector3

	shape shape.IShape

	// TODO future (for now a body is a single graphic with a single geometry)
	// shapes          []*Shape
	// shapeOffsets    []float32 // Position of each Shape in the body, given in local Body space.
	// shapeOrientations [] ?
}

// BodyType specifies how the body is affected during the simulation.
type BodyType int

const (
	// A static body does not move during simulation and behaves as if it has infinite mass.
	// Static bodies can be moved manually by setting the position of the body.
	// The velocity of a static body is always zero.
	// Static bodies do not collide with other static or kinematic bodies.
	Static = BodyType(iota)

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

// TODO Update simulation checks for BodyType to use bitwise operators ?

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

// TODO
type HullType int
const (
	Sphere = HullType(iota)
	Capsule
	Mesh // use mesh itself
)


// NewBody creates and returns a pointer to a new RigidBody.
// The igraphic's geometry *must* be convex.
func NewBody(igraphic graphic.IGraphic) *Body {

	b := new(Body)
	b.Graphic = igraphic.GetGraphic()
	b.bodyType = Dynamic

	// Rotational inertia and related properties
	b.rotInertia = math32.NewMatrix3()
	b.invRotInertia = math32.NewMatrix3()
	b.invRotInertiaEff = math32.NewMatrix3()
	b.invRotInertiaWorld = math32.NewMatrix3()
	b.invRotInertiaWorldEff = math32.NewMatrix3()

	// Position
	pos := b.GetNode().Position()
	b.position = pos.Clone()
	b.prevPosition = pos.Clone()
	b.interpPosition = pos.Clone()
	b.initPosition = pos.Clone()

	// Rotation
	quat := b.GetNode().Quaternion()
	b.quaternion = quat.Clone()
	b.prevQuaternion = quat.Clone()
	b.interpQuaternion = quat.Clone()
	b.initQuaternion = quat.Clone()

	// Linear and angular velocities
	b.velocity = math32.NewVec3()
	b.initVelocity = math32.NewVec3()
	b.angularVelocity = math32.NewVec3()
	b.initAngularVelocity = math32.NewVec3()

	// Force and torque
	b.force = math32.NewVec3()
	b.torque = math32.NewVec3()

	// Damping and factors
	b.linearDamping = 0.01
	b.angularDamping = 0.01
	b.linearFactor = math32.NewVector3(1, 1, 1)
	b.angularFactor = math32.NewVector3(1, 1, 1)

	// Sleep settings
	b.allowSleep = true
	b.sleepState = Awake
	b.sleepSpeedLimit = 0.1
	b.sleepTimeLimit = 1
	b.timeLastSleepy = 0

	// Collision filtering
	b.colFilterGroup = 1
	b.colFilterMask = -1

	//b.fixedRotation = true

	b.wakeUpAfterNarrowphase = false

	b.SetShape(shape.NewConvexHull(b.GetGeometry()))

	b.SetMass(1)
	b.UpdateMassProperties()
	b.UpdateEffectiveMassProperties()

	return b
}

// TODO future: modify this to be "AddShape" and keep track of list of shapes, their positions and orientations
// For now each body can only be a single shape or a single geometry
func (b *Body) SetShape(shape shape.IShape) {

	b.shape = shape
}

func (b *Body) Shape() shape.IShape {

	return b.shape
}

func (b *Body) BoundingBox() math32.Box3 {

	// TODO future allow multiple shapes
	mat4 := math32.NewMatrix4().Compose(b.position, b.quaternion, math32.NewVector3(1,1,1))
	localBB := b.shape.BoundingBox()
	worldBB := localBB.ApplyMatrix4(mat4)
	return *worldBB
}

func (b *Body) SetMass(mass float32) {

	// Do nothing if current mass is already the specified mass
	if mass == b.mass {
		return
	}

	// Set mass and update inverse mass
	b.mass = mass
	if b.mass > 0 {
		b.invMass = 1.0 / b.mass
	} else {
		// Body mass is zero - this means that the body is static
		b.invMass = 0
		b.bodyType = Static
	}

	b.UpdateMassProperties()
}

func (b *Body) SetIndex(i int) {

	b.index = i
}

func (b *Body) SetName(name string) {

	b.name = name
}

func (b *Body) Name() string {

	return b.name
}

func (b *Body) Index() int {

	return b.index
}

func (b *Body) Material() *material.Material {

	return b.material
}

func (b *Body) SetAllowSleep(state bool) {

	b.allowSleep = state
}

func (b *Body) AllowSleep() bool {

	return b.allowSleep
}

func (b *Body) SleepSpeedLimit() float32 {

	return b.sleepSpeedLimit
}

func (b *Body) SleepState() BodySleepState {

	return b.sleepState
}

// SetBodyType sets the body type.
func (b *Body) SetBodyType(bodyType BodyType) {

	// Do nothing if body is already of the specified bodyType
	if b.bodyType == bodyType {
		return
	}

	// If we want the body to be static we need to zero its mass
	if bodyType == Static {
		b.mass = 0
	}

	// Temporarily save original body type and update current body type
	origBodyType := b.bodyType
	b.bodyType = bodyType

	// If changed body type to or from Static then we need to update mass properties
	if origBodyType == Static || b.bodyType == Static {
		b.UpdateMassProperties()
	}
}

func (b *Body) BodyType() BodyType {

	return b. bodyType
}

func (b *Body) SetWakeUpAfterNarrowphase(state bool) {

	b.wakeUpAfterNarrowphase = state
}

func (b *Body) WakeUpAfterNarrowphase() bool {

	return b.wakeUpAfterNarrowphase
}

// ApplyVelocityDeltas adds the specified deltas to the body's linear and angular velocities.
func (b *Body) ApplyVelocityDeltas(linearD, angularD *math32.Vector3) {

	b.velocity.Add(linearD.Multiply(b.linearFactor))
	b.angularVelocity.Add(angularD.Multiply(b.angularFactor))
}

// ClearForces clears all forces on the body.
func (b *Body) ClearForces() {

	b.force.Zero()
	b.torque.Zero()
}

func (b *Body) InvMassEff() float32 {

	return b.invMassEff
}

func (b *Body) InvRotInertiaWorldEff() *math32.Matrix3 {

	return b.invRotInertiaWorldEff
}

func (b *Body) Position() math32.Vector3 {

	return *b.position
}

func (b *Body) Quaternion() *math32.Quaternion {

	return b.quaternion.Clone()
}

func (b *Body) SetVelocity(vel *math32.Vector3) {

	b.velocity = vel
}

func (b *Body) Velocity() math32.Vector3 {

	return *b.velocity
}

func (b *Body) SetAngularVelocity(vel *math32.Vector3) {

	b.angularVelocity = vel
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

func (b *Body) SetLinearDamping(d float32) {

	b.linearDamping = d
}

func (b *Body) LinearDamping() float32 {

	return b.linearDamping
}

func (b *Body) SetAngularDamping(d float32) {

	b.angularDamping = d
}

func (b *Body) AngularDamping() float32 {

	return b.angularDamping
}

func (b *Body) ApplyDamping(dt float32) {

	b.velocity.MultiplyScalar(math32.Pow(1.0 - b.linearDamping, dt))
	b.angularVelocity.MultiplyScalar(math32.Pow(1.0 - b.angularDamping, dt))
}

func (b *Body) SetLinearFactor(factor *math32.Vector3) {

	b.linearFactor = factor
}

func (b *Body) LinearFactor() math32.Vector3 {

	return *b.linearFactor
}

func (b *Body) SetAngularFactor(factor *math32.Vector3) {

	b.angularFactor = factor
}

func (b *Body) AngularFactor() math32.Vector3 {

	return *b.angularFactor
}

// SetFixedRotation specifies whether the body should rotate.
func (b *Body) SetFixedRotation(state bool) {

	// Do nothing if the fixedRotation flag already has the specified value
	if b.fixedRotation == state {
		return
	}

	// Set the fixedRotation flag and update mass properties
	b.fixedRotation = state
	b.UpdateMassProperties()
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
	b.velocity.Set(0, 0, 0)
	b.angularVelocity.Set(0, 0, 0)
	b.wakeUpAfterNarrowphase = false
}

// Called every timestep to update internal sleep timer and change sleep state if needed.
// time: The world time in seconds
func (b *Body) SleepTick(time float32) {

	if b.allowSleep {
		speedSquared := b.velocity.LengthSq() + b.angularVelocity.LengthSq()
		speedLimitSquared := math32.Pow(b.sleepSpeedLimit, 2)
		if b.sleepState == Awake && speedSquared < speedLimitSquared {
			b.sleepState = Sleepy
			b.timeLastSleepy = time
			b.Dispatch(SleepyEvent, nil)
		} else if b.sleepState == Sleepy && speedSquared > speedLimitSquared {
			b.WakeUp() // Wake up
		} else if b.sleepState == Sleepy && (time-b.timeLastSleepy) > b.sleepTimeLimit {
			b.Sleep() // Sleeping
			b.Dispatch(SleepEvent, nil)
		}
	}
}

// If checkSleeping is true then returns false if both bodies are currently sleeping.
func (b *Body) Sleeping() bool {

	return b.sleepState == Sleeping
}

// CollidableWith returns whether the body can collide with the specified body.
func (b *Body) CollidableWith(other *Body) bool {

	if (b.colFilterGroup & other.colFilterMask == 0) ||
		(other.colFilterGroup & b.colFilterMask == 0) ||
		(b.bodyType == Static) && (other.bodyType == Static) {
		return false
	}

	return true
}

func (b *Body) CollisionResponse() bool {

	return b.colResponse
}

// PointToLocal converts a world point to local body frame. TODO maybe move to Node
func (b *Body) PointToLocal(worldPoint *math32.Vector3) math32.Vector3 {

	return *worldPoint.Clone().Sub(b.position).ApplyQuaternion(b.quaternion.Conjugate())
}

// VectorToLocal converts a world vector to local body frame. TODO maybe move to Node
func (b *Body) VectorToLocal(worldVector *math32.Vector3) math32.Vector3 {

	return *worldVector.Clone().ApplyQuaternion(b.quaternion.Conjugate())
}

// PointToWorld converts a local point to world frame. TODO maybe move to Node
func (b *Body) PointToWorld(localPoint *math32.Vector3) math32.Vector3 {

	return *localPoint.Clone().ApplyQuaternion(b.quaternion).Add(b.position)
}

// VectorToWorld converts a local vector to world frame. TODO maybe move to Node
func (b *Body) VectorToWorld(localVector *math32.Vector3) math32.Vector3 {

	return *localVector.Clone().ApplyQuaternion(b.quaternion)
}

// UpdateEffectiveMassProperties
// If the body is sleeping, it should be immovable and thus have infinite mass during solve.
// This is solved by having a separate "effective mass" and other "effective" properties
func (b *Body) UpdateEffectiveMassProperties() {

	if b.sleepState == Sleeping || b.bodyType == Kinematic {
		b.invMassEff = 0
		b.invRotInertiaEff.Zero()
		b.invRotInertiaWorldEff.Zero()
	} else {
		b.invMassEff = b.invMass
		b.invRotInertiaEff.Copy(b.invRotInertia)
		b.invRotInertiaWorldEff.Copy(b.invRotInertiaWorld)
	}
}

// UpdateMassProperties
// Should be called whenever you change the body shape or mass.
func (b *Body) UpdateMassProperties() {

	if b.mass > 0 {
		b.invMass = 1.0 / b.mass
	} else {
		b.invMass = 0
	}

	if b.fixedRotation || b.bodyType == Static {
		b.rotInertia.Zero()
		b.invRotInertia.Zero()
	} else {
		*b.rotInertia = b.GetGeometry().RotationalInertia(b.mass)
		b.rotInertia.MultiplyScalar(10) // multiply by high density // TODO remove this ?
		b.invRotInertia.GetInverse(b.rotInertia) // Note: rotInertia is always positive definite and thus always invertible
	}

	b.UpdateInertiaWorld(true)
}

// Update .inertiaWorld and .invRotInertiaWorld
func (b *Body) UpdateInertiaWorld(force bool) {

	iRI := b.invRotInertia
	// If rotational inertia M = s*I, where I is identity and s a scalar, then
	//    R*M*R' = R*(s*I)*R' = s*R*I*R' = s*R*R' = s*I = M
	// where R is the rotation matrix.
	// In other words, we don't have to do the transformation if all diagonal entries are equal.
	if iRI[0] != iRI[4] || iRI[4] != iRI[8] || force {
		// iRIW = R * iRI * R'
		m1 := math32.NewMatrix3().MakeRotationFromQuaternion(b.quaternion)
		m2 := m1.Clone().Transpose()
		m2.Multiply(iRI)
		b.invRotInertiaWorld.MultiplyMatrices(m2, m1)
	}
}

// Forces from a force field need to be multiplied by mass.
func (b *Body) ApplyForceField(force *math32.Vector3) {

	b.force.Add(force.MultiplyScalar(b.mass))
}

// Apply force to a world point.
// This could for example be a point on the Body surface.
// Applying force this way will add to Body.force and Body.torque.
// relativePoint: A point relative to the center of mass to apply the force on.
func (b *Body) ApplyForce(force, relativePoint *math32.Vector3) {

	if b.bodyType != Dynamic { // Needed?
		return
	}

	// Add linear force
	b.force.Add(force) // TODO shouldn't rotational momentum be subtracted from linear momentum?

	// Add rotational force
	b.torque.Add(math32.NewVec3().CrossVectors(relativePoint, force))
}

// Apply force to a local point in the body.
// force: The force vector to apply, defined locally in the body frame.
// localPoint: A local point in the body to apply the force on.
func (b *Body) ApplyLocalForce(localForce, localPoint *math32.Vector3) {

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
	velo := impulse.Clone().MultiplyScalar(b.invMass)

	// Add linear impulse
	b.velocity.Add(velo)

	// Compute produced rotational impulse velocity
	rotVelo := math32.NewVec3().CrossVectors(r, impulse)
	rotVelo.ApplyMatrix3(b.invRotInertiaWorld)

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

	r := math32.NewVec3().SubVectors(worldPoint, b.position)
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
	e := b.invRotInertiaWorld
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
	b.quaternion.X += halfDt * (ax*bw + ay*bz - az*by)
	b.quaternion.Y += halfDt * (ay*bw + az*bx - ax*bz)
	b.quaternion.Z += halfDt * (az*bw + ax*by - ay*bx)
	b.quaternion.W += halfDt * (-ax*bx - ay*by - az*bz)

	// Normalize quaternion
	b.quaternion.Normalize()
	//if quatNormalize { // TODO future
	//	if quatNormalizeFast {
	//		b.quaternion.NormalizeFast()
	//	} else {
	//		b.quaternion.Normalize()
	//	}
	//}

	// Update position and rotation of Node (containing visual representation of the body)
	b.GetNode().SetPositionVec(b.position)
	b.GetNode().SetRotationQuat(b.quaternion)

	b.aabbNeedsUpdate = true

	// Update world inertia
	b.UpdateInertiaWorld(false)
}
