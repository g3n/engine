// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package physics

import "github.com/g3n/engine/math32"

// ForceField represents a force field. A force is defined for every point.
type ForceField interface {
	ForceAt(pos *math32.Vector3) *math32.Vector3
}

//
// Constant is a constant force field.
// It can be used to simulate surface gravity.
//
type Constant struct {
	force math32.Vector3
}

// NewConstant creates and returns a pointer to a new Constant force field.
func NewConstant(acceleration float32) *Constant {

	g := new(Constant)
	g.force = math32.Vector3{0,0,-acceleration}
	return g
}

// SetForce sets the force of the force field.
func (g *Constant) SetForce(newDirection *math32.Vector3) {

	g.force = *newDirection
}

// Force returns the force of the force field.
func (g *Constant) Force() *math32.Vector3 {

	return &g.force
}

// ForceAt satisfies the ForceField interface and returns the force at the specified position.
func (g *Constant) ForceAt(pos *math32.Vector3) *math32.Vector3 {

	return &g.force
}

//
// PointAttractor is a force field where all forces point to a single point.
// The force strength changes with the inverse distance squared.
// This can be used to model planetary attractions.
//
type PointAttractor struct {
	position math32.Vector3
	mass     float32
}

// NewPointAttractor creates and returns a pointer to a new PointAttractor force field.
func NewPointAttractor(position *math32.Vector3, mass float32) *PointAttractor {

	pa := new(PointAttractor)
	pa.position = *position
	pa.mass = mass
	return pa
}

// SetPosition sets the position of the PointAttractor.
func (pa *PointAttractor) SetPosition(newPosition *math32.Vector3) {

	pa.position = *newPosition
}

// Position returns the position of the PointAttractor.
func (pa *PointAttractor) Position() *math32.Vector3 {

	return &pa.position
}

// SetMass sets the mass of the PointAttractor.
func (pa *PointAttractor) SetMass(newMass float32) {

	pa.mass = newMass
}

// Mass returns the mass of the PointAttractor.
func (pa *PointAttractor) Mass() float32 {

	return pa.mass
}

// ForceAt satisfies the ForceField interface and returns the force at the specified position.
func (pa *PointAttractor) ForceAt(pos *math32.Vector3) *math32.Vector3 {

	dir := pos
	dir.Negate()
	dir.Add(&pa.position)
	dist := dir.Length()
	dir.Normalize()
	dir.MultiplyScalar(pa.mass/(dist*dist))
	return dir
}

//
// PointRepeller is a force field where all forces point away from a single point.
// The force strength changes with the inverse distance squared.
//
type PointRepeller struct {
	position math32.Vector3
	mass     float32
}

// NewPointRepeller creates and returns a pointer to a new PointRepeller force field.
func NewPointRepeller(position *math32.Vector3, mass float32) *PointRepeller {

	pr := new(PointRepeller)
	pr.position = *position
	pr.mass = mass
	return pr
}

// SetPosition sets the position of the PointRepeller.
func (pr *PointRepeller) SetPosition(newPosition *math32.Vector3) {

	pr.position = *newPosition
}

// Position returns the position of the PointRepeller.
func (pr *PointRepeller) Position() *math32.Vector3 {

	return &pr.position
}

// SetMass sets the mass of the PointRepeller.
func (pr *PointRepeller) SetMass(newMass float32) {

	pr.mass = newMass
}

// Mass returns the mass of the PointRepeller.
func (pr *PointRepeller) Mass() float32 {

	return pr.mass
}

// ForceAt satisfies the ForceField interface and returns the force at the specified position.
func (pr *PointRepeller) ForceAt(pos *math32.Vector3) *math32.Vector3 {

	dir := pr.position
	dir.Negate()
	dir.Add(pos)
	dist := dir.Length()
	dir.Normalize()
	dir.MultiplyScalar(pr.mass/(dist*dist))
	return &dir
}
