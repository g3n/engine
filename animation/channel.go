// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package animation

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/geometry"
)

// A Channel associates an animation parameter channel to an interpolation sampler
type Channel struct {
	keyframes          math32.ArrayF32          // Input keys (usually time)
	values             math32.ArrayF32          // Outputs values for the keys
	interpType         InterpolationType        // Interpolation type
	interpAction       func(idx int, k float32) // Combined function for interpolation and update
	updateInterpAction func()                   // Function to update interpAction based on interpolation type
	inTangent          math32.ArrayF32          // Origin tangents for Spline interpolation
	outTangent         math32.ArrayF32          // End tangents for Spline interpolation
}

// SetBuffers sets the keyframe and value buffers.
func (c *Channel) SetBuffers(keyframes, values math32.ArrayF32) {

	c.keyframes = keyframes
	c.values = values
}

// Keyframes returns the keyframe buffer.
func (c *Channel) Keyframes() math32.ArrayF32 {

	return c.keyframes
}

// Values returns the value buffer.
func (c *Channel) Values() math32.ArrayF32 {

	return c.values
}

// SetInterpolationTangents sets the interpolation tangents.
func (c *Channel) SetInterpolationTangents(inTangent, outTangent math32.ArrayF32) {

	c.inTangent = inTangent
	c.outTangent = outTangent
}

// InterpolationTangents sets the interpolation tangents
func (c *Channel) InterpolationTangents() (inTangent, outTangent math32.ArrayF32) {

	return c.inTangent, c.outTangent
}

// SetInterpolationType sets the interpolation type for this channel.
func (c *Channel) SetInterpolationType(it InterpolationType) {

	// Don't update function if not needed
	if c.interpType == it {
		return
	}

	// Save interpolation type
	c.interpType = it

	// Call specialized function that updates the interpAction function
	c.updateInterpAction()
}

// InterpolationType returns the current interpolation type.
func (c *Channel) InterpolationType() InterpolationType {

	return c.interpType
}

// Update finds the keyframe preceding the specified time.
// Then, calls a stored function to interpolate the relevant values and update the target.
func (c *Channel) Update(time float32) {

	// Test limits
	if (len(c.keyframes) < 2) || (time < c.keyframes[0]) || (time > c.keyframes[len(c.keyframes)-1]) {
		return
	}

	// Find keyframe interval
	var idx int
	for idx = 0; idx < len(c.keyframes)-1; idx++ {
		if time >= c.keyframes[idx] && time < c.keyframes[idx+1] {
			break
		}
	}

	// Interpolate and update
	relativeDelta := (time-c.keyframes[idx])/(c.keyframes[idx+1]-c.keyframes[idx])
	c.interpAction(idx, relativeDelta)
}

// IChannel is the interface for all channel types.
type IChannel interface {
	Update(time float32)
	SetBuffers(keyframes, values math32.ArrayF32)
	Keyframes() math32.ArrayF32
	Values() math32.ArrayF32
	SetInterpolationType(it InterpolationType)
}

// NodeChannel is the IChannel for all node transforms.
type NodeChannel struct {
	Channel
	target core.INode
}

// PositionChannel is the animation channel for a node's position.
type PositionChannel NodeChannel

func NewPositionChannel(node core.INode) *PositionChannel {

	pc := new(PositionChannel)
	pc.target = node
	pc.updateInterpAction = func() {
		// Get node
		node := pc.target.GetNode()
		// Update interpolation function
		switch pc.interpType {
		case STEP:
			pc.interpAction = func(idx int, k float32) {
				var v math32.Vector3
				pc.values.GetVector3(idx*3, &v)
				node.SetPositionVec(&v)
			}
		case LINEAR:
			pc.interpAction = func(idx int, k float32) {
				var v1, v2 math32.Vector3
				pc.values.GetVector3(idx*3, &v1)
				pc.values.GetVector3((idx+1)*3, &v2)
				v1.Lerp(&v2, k)
				node.SetPositionVec(&v1)
			}
		case CUBICSPLINE: // TODO
			pc.interpAction = func(idx int, k float32) {
				var v1, v2 math32.Vector3
				pc.values.GetVector3(idx*3, &v1)
				pc.values.GetVector3((idx+1)*3, &v2)
				v1.Lerp(&v2, k)
				node.SetPositionVec(&v1)
			}
		}
	}
	pc.SetInterpolationType(LINEAR)
	return pc
}

// RotationChannel is the animation channel for a node's rotation.
type RotationChannel NodeChannel

func NewRotationChannel(node core.INode) *RotationChannel {

	rc := new(RotationChannel)
	rc.target = node
	rc.updateInterpAction = func() {
		// Get node
		node := rc.target.GetNode()
		// Update interpolation function
		switch rc.interpType {
		case STEP:
			rc.interpAction = func(idx int, k float32) {
				var q math32.Vector4
				rc.values.GetVector4(idx*4, &q)
				node.SetQuaternionVec(&q)
			}
		case LINEAR:
			rc.interpAction = func(idx int, k float32) {
				var q1, q2 math32.Vector4
				rc.values.GetVector4(idx*4, &q1)
				rc.values.GetVector4((idx+1)*4, &q2)
				quat1 := math32.NewQuaternion(q1.X, q1.Y, q1.Z, q1.W)
				quat2 := math32.NewQuaternion(q2.X, q2.Y, q2.Z, q2.W)
				quat1.Slerp(quat2, k)
				node.SetQuaternionQuat(quat1)
			}
		case CUBICSPLINE: // TODO
			rc.interpAction = func(idx int, k float32) {
				var q1, q2 math32.Vector4
				rc.values.GetVector4(idx*4, &q1)
				rc.values.GetVector4((idx+1)*4, &q2)
				quat1 := math32.NewQuaternion(q1.X, q1.Y, q1.Z, q1.W)
				quat2 := math32.NewQuaternion(q2.X, q2.Y, q2.Z, q2.W)
				quat1.Slerp(quat2, k)
				node.SetQuaternionQuat(quat1)
			}
		}
	}
	rc.SetInterpolationType(LINEAR)
	return rc
}

// ScaleChannel is the animation channel for a node's scale.
type ScaleChannel NodeChannel

func NewScaleChannel(node core.INode) *ScaleChannel {

	sc := new(ScaleChannel)
	sc.target = node
	sc.updateInterpAction = func() {
		// Get node
		node := sc.target.GetNode()
		// Update interpolation function
		switch sc.interpType {
		case STEP:
			sc.interpAction = func(idx int, k float32) {
				var v math32.Vector3
				sc.values.GetVector3(idx*3, &v)
				node.SetScaleVec(&v)
			}
		case LINEAR:
			sc.interpAction = func(idx int, k float32) {
				var v1, v2 math32.Vector3
				sc.values.GetVector3(idx*3, &v1)
				sc.values.GetVector3((idx+1)*3, &v2)
				v1.Lerp(&v2, k)
				node.SetScaleVec(&v1)
			}
		case CUBICSPLINE: // TODO
			sc.interpAction = func(idx int, k float32) {
				var v1, v2 math32.Vector3
				sc.values.GetVector3(idx*3, &v1)
				sc.values.GetVector3((idx+1)*3, &v2)
				v1.Lerp(&v2, k)
				node.SetScaleVec(&v1)
			}
		}
	}
	sc.SetInterpolationType(LINEAR)
	return sc
}

// MorphChannel is the IChannel for morph geometries.
type MorphChannel struct {
	Channel
	target *geometry.MorphGeometry
}

func NewMorphChannel(mg *geometry.MorphGeometry) *MorphChannel {

	mc := new(MorphChannel)
	mc.target = mg
	numWeights := len(mg.Weights())
	mc.updateInterpAction = func() {
		// Update interpolation function
		switch mc.interpType {
		case STEP:
			mc.interpAction = func(idx int, k float32) {
				start := idx*numWeights
				weights := mc.values[start:start+numWeights]
				mg.SetWeights(weights)
			}
		case LINEAR:
			mc.interpAction = func(idx int, k float32) {
				start1 := idx*numWeights
				start2 := (idx+1)*numWeights
				weights1 := mc.values[start1:start1+numWeights]
				weights2 := mc.values[start2:start2+numWeights]
				weightsNew := make([]float32, numWeights)
				for i := range weights1 {
					weightsNew[i] = weights1[i] + (weights2[i]-weights1[i])*k
				}
				mg.SetWeights(weightsNew)
			}
		case CUBICSPLINE: // TODO
			mc.interpAction = func(idx int, k float32) {
				start1 := idx*numWeights
				start2 := (idx+1)*numWeights
				weights1 := mc.values[start1:start1+numWeights]
				weights2 := mc.values[start2:start2+numWeights]
				weightsNew := make([]float32, numWeights)
				for i := range weights1 {
					weightsNew[i] = weights1[i] + (weights2[i]-weights1[i])*k
				}
				mg.SetWeights(weightsNew)
			}
		}
	}
	mc.SetInterpolationType(LINEAR)
	return mc
}

// InterpolationType specifies the interpolation type.
type InterpolationType string

// The various interpolation types.
const (
	STEP        = InterpolationType("STEP")          // The animated values remain constant to the output of the first keyframe, until the next keyframe.
	LINEAR      = InterpolationType("LINEAR")        // The animated values are linearly interpolated between keyframes. Spherical linear interpolation (slerp) is used to interpolate quaternions.
	CUBICSPLINE = InterpolationType("CUBICSPLINE")   // TODO
)
