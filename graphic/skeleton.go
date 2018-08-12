// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphic

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/math32"
	)

// Skeleton contains armature information.
type Skeleton struct {
	inverseBindMatrices []math32.Matrix4
	boneMatrices        []math32.Matrix4
	bones               []*core.Node
}

// NewSkeleton creates and returns a pointer to a new Skeleton.
func NewSkeleton() *Skeleton {

	sk := new(Skeleton)
	sk.boneMatrices = make([]math32.Matrix4, 0)
	sk.bones = make([]*core.Node, 0)
	return sk
}

// AddBone adds a bone to the skeleton along with an optional inverseBindMatrix.
func (sk *Skeleton) AddBone(node *core.Node, inverseBindMatrix *math32.Matrix4) {

	// Useful for debugging:
	//node.Add(NewAxisHelper(0.2))

	sk.bones = append(sk.bones, node)
	sk.boneMatrices = append(sk.boneMatrices, *math32.NewMatrix4())
	if inverseBindMatrix == nil {
		inverseBindMatrix = math32.NewMatrix4() // Identity matrix
	}

	sk.inverseBindMatrices = append(sk.inverseBindMatrices, *inverseBindMatrix)
}

// Bones returns the list of bones in the skeleton.
func (sk *Skeleton) Bones() []*core.Node {

	return sk.bones
}

// BoneMatrices calculates and returns the bone world matrices to be sent to the shader.
func (sk *Skeleton) BoneMatrices(invMat *math32.Matrix4) []math32.Matrix4 {

	// Update bone matrices based on inverseBindMatrices and the provided invMat
	for i := range sk.bones {
		bMat := sk.bones[i].MatrixWorld()
		bMat.MultiplyMatrices(&bMat, &sk.inverseBindMatrices[i])
		sk.boneMatrices[i].MultiplyMatrices(invMat, &bMat)
	}

	return sk.boneMatrices
}
