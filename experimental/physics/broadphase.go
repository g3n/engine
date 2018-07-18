// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package collision implements collision related algorithms and data structures.
package physics

import (
	"github.com/g3n/engine/experimental/physics/object"
)

// CollisionPair is a pair of bodies that may be colliding.
type CollisionPair struct {
	BodyA *object.Body
	BodyB *object.Body
}

// Broadphase is the base class for broadphase implementations.
type Broadphase struct {}

// NewBroadphase creates and returns a pointer to a new Broadphase.
func NewBroadphase() *Broadphase {

	b := new(Broadphase)
	return b
}

// FindCollisionPairs (naive implementation)
func (b *Broadphase) FindCollisionPairs(objects []*object.Body) []CollisionPair {

	pairs := make([]CollisionPair,0)

	for iA, bodyA := range objects {
		for _, bodyB := range objects[iA+1:] {
			if b.NeedTest(bodyA, bodyB) {
				BBa := bodyA.BoundingBox()
				BBb := bodyB.BoundingBox()
				if BBa.IsIntersectionBox(&BBb) {
					pairs = append(pairs, CollisionPair{bodyA, bodyB})
				}
			}
		}
	}

	return pairs
}

func (b *Broadphase) NeedTest(bodyA, bodyB *object.Body) bool {

	if !bodyA.CollidableWith(bodyB) || (bodyA.Sleeping() && bodyB.Sleeping()) {
		return false
	}

	return true
}
