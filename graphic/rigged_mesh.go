// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphic

import (
	"github.com/g3n/engine/gls"
	"strconv"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/math32"
)

// MaxBoneInfluencers is the maximum number of bone influencers per vertex.
const MaxBoneInfluencers = 4

// RiggedMesh is a Mesh associated with a skeleton.
type RiggedMesh struct {
	*Mesh    // Embedded mesh
	skeleton *Skeleton
	mBones   gls.Uniform
}

// NewRiggedMesh returns a new rigged mesh.
func NewRiggedMesh(mesh *Mesh) *RiggedMesh {

	rm := new(RiggedMesh)
	rm.Mesh = mesh
	rm.SetIGraphic(rm)
	rm.mBones.Init("mBones")
	rm.ShaderDefines.Set("BONE_INFLUENCERS", strconv.Itoa(MaxBoneInfluencers))
	rm.ShaderDefines.Set("TOTAL_BONES", "0")

	return rm
}

// SetSkeleton sets the skeleton used by the rigged mesh.
func (rm *RiggedMesh) SetSkeleton(sk *Skeleton) {

	rm.skeleton = sk
	rm.ShaderDefines.Set("TOTAL_BONES", strconv.Itoa(len(rm.skeleton.Bones())))
}

// SetSkeleton returns the skeleton used by the rigged mesh.
func (rm *RiggedMesh) Skeleton() *Skeleton {

	return rm.skeleton
}

// RenderSetup is called by the renderer before drawing the geometry.
func (rm *RiggedMesh) RenderSetup(gs *gls.GLS, rinfo *core.RenderInfo) {

	// Call base mesh's RenderSetup
	rm.Mesh.RenderSetup(gs, rinfo)

	// Get inverse matrix world
	var invMat math32.Matrix4
	node := rm.GetNode()
	nMW := node.MatrixWorld()
	err := invMat.GetInverse(&nMW)
	if err != nil {
		log.Error("Skeleton.BoneMatrices: inverting matrix failed!")
	}

	// Transfer bone matrices
	boneMatrices := rm.skeleton.BoneMatrices(&invMat)
	location := rm.mBones.Location(gs)
	gs.UniformMatrix4fv(location, int32(len(boneMatrices)), false, &boneMatrices[0][0])
}
