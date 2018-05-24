// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package collision implements collision related algorithms and data structures.
package collision

import (
	"github.com/g3n/engine/physics/object"
)

// Broadphase is the base class for broadphase implementations.
type Broadphase struct {}

type Pair struct {
	BodyA *object.Body
	BodyB *object.Body
}

// NewBroadphase creates and returns a pointer to a new Broadphase.
func NewBroadphase() *Broadphase {

	b := new(Broadphase)
	return b
}

// FindCollisionPairs (naive implementation)
func (b *Broadphase) FindCollisionPairs(objects []*object.Body) []Pair {

	pairs := make([]Pair,0)

	for _, bodyA := range objects {
		for _, bodyB := range objects {
			if b.NeedTest(bodyA, bodyB) {
				BBa := bodyA.BoundingBox()
				BBb := bodyB.BoundingBox()
				if BBa.IsIntersectionBox(&BBb) {
					pairs = append(pairs, Pair{bodyA, bodyB})
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
