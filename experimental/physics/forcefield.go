// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package physics

import "github.com/g3n/engine/math32"

// ForceField represents a force field. A force is defined for every point.
type ForceField interface {
	ForceAt(pos *math32.Vector3) math32.Vector3
}

//
// ConstantForceField is a constant force field.
// It can be used to simulate surface gravity.
//
type ConstantForceField struct {
	force math32.Vector3
}

// NewConstantForceField creates and returns a pointer to a new ConstantForceField.
func NewConstantForceField(force *math32.Vector3) *ConstantForceField {

	g := new(ConstantForceField)
	g.force = *force
	return g
}

// SetForce sets the force of the force field.
func (g *ConstantForceField) SetForce(newDirection *math32.Vector3) {

	g.force = *newDirection
}

// Force returns the force of the force field.
func (g *ConstantForceField) Force() *math32.Vector3 {

	return &g.force
}

// ForceAt satisfies the ForceField interface and returns the force at the specified position.
func (g *ConstantForceField) ForceAt(pos *math32.Vector3) math32.Vector3 {

	return g.force
}

//
// AttractorForceField is a force field where all forces point to a single point.
// The force strength changes with the inverse distance squared.
// This can be used to model planetary attractions.
//
type AttractorForceField struct {
	position math32.Vector3
	mass     float32
}

// NewAttractorForceField creates and returns a pointer to a new AttractorForceField.
func NewAttractorForceField(position *math32.Vector3, mass float32) *AttractorForceField {

	pa := new(AttractorForceField)
	pa.position = *position
	pa.mass = mass
	return pa
}

// SetPosition sets the position of the AttractorForceField.
func (pa *AttractorForceField) SetPosition(newPosition *math32.Vector3) {

	pa.position = *newPosition
}

// Position returns the position of the AttractorForceField.
func (pa *AttractorForceField) Position() *math32.Vector3 {

	return &pa.position
}

// SetMass sets the mass of the AttractorForceField.
func (pa *AttractorForceField) SetMass(newMass float32) {

	pa.mass = newMass
}

// Mass returns the mass of the AttractorForceField.
func (pa *AttractorForceField) Mass() float32 {

	return pa.mass
}

// ForceAt satisfies the ForceField interface and returns the force at the specified position.
func (pa *AttractorForceField) ForceAt(pos *math32.Vector3) math32.Vector3 {

	dir := pos
	dir.Negate()
	dir.Add(&pa.position)
	dist := dir.Length()
	dir.Normalize()
	var val float32
	//log.Error("dist %v", dist)
	if dist > 0 {
		val = pa.mass/(dist*dist)
	} else {
		val = 0
	}
	val = math32.Min(val, 100) // TODO deal with instability
	//log.Error("%v", val)
	dir.MultiplyScalar(val) // TODO multiply by gravitational constant: 6.673×10−11 (N–m2)/kg2
	return *dir
}

//
// RepellerForceField is a force field where all forces point away from a single point.
// The force strength changes with the inverse distance squared.
//
type RepellerForceField struct {
	position math32.Vector3
	mass     float32
}

// NewRepellerForceField creates and returns a pointer to a new RepellerForceField.
func NewRepellerForceField(position *math32.Vector3, mass float32) *RepellerForceField {

	pr := new(RepellerForceField)
	pr.position = *position
	pr.mass = mass
	return pr
}

// SetPosition sets the position of the RepellerForceField.
func (pr *RepellerForceField) SetPosition(newPosition *math32.Vector3) {

	pr.position = *newPosition
}

// Position returns the position of the RepellerForceField.
func (pr *RepellerForceField) Position() *math32.Vector3 {

	return &pr.position
}

// SetMass sets the mass of the RepellerForceField.
func (pr *RepellerForceField) SetMass(newMass float32) {

	pr.mass = newMass
}

// Mass returns the mass of the RepellerForceField.
func (pr *RepellerForceField) Mass() float32 {

	return pr.mass
}

// ForceAt satisfies the ForceField interface and returns the force at the specified position.
func (pr *RepellerForceField) ForceAt(pos *math32.Vector3) math32.Vector3 {

	dir := pr.position
	dir.Negate()
	dir.Add(pos)
	dist := dir.Length()
	dir.Normalize()
	dir.MultiplyScalar(pr.mass/(dist*dist)) // TODO multiply by gravitational constant: 6.673×10−11 (N–m2)/kg2
	return dir
}
