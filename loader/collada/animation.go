// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package collada

import (
	"fmt"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/math32"
	"strings"
)

// AnimationTarget contains all animation channels for an specific target node
type AnimationTarget struct {
	target   core.INode
	matrix   math32.Matrix4 // original node transformation matrix
	start    float32        // initial input offset value
	last     float32        // last input value
	minInput float32        // minimum input value for all channels
	maxInput float32        // maximum input value for all channels
	loop     bool           // animation loop flag
	rot      math32.Vector3 // rotation in XYZ Euler angles
	channels []*ChannelInstance
}

// A ChannelInstance associates an animation parameter channel to an interpolation sampler
type ChannelInstance struct {
	sampler *SamplerInstance
	action  ActionFunc
}

// SamplerInstance specifies the input key frames, output values for these key frames
// and interpolation information. It can be shared by more than one animation
type SamplerInstance struct {
	Input      []float32 // Input keys (usually time)
	Output     []float32 // Outputs values for the keys
	Interp     []string  // Names of interpolation functions for each key frame
	InTangent  []float32 // Origin tangents for Bezier interpolation
	OutTangent []float32 // End tangents for Bezier interpolation
}

// ActionFunc is the type for all functions that execute an specific parameter animation
type ActionFunc func(at *AnimationTarget, v float32)

// Reset resets the animation from the beginning
func (at *AnimationTarget) Reset() {

	at.last = at.start
	at.target.GetNode().SetMatrix(&at.matrix)
}

// SetLoop sets the state of the animation loop flag
func (at *AnimationTarget) SetLoop(loop bool) {

	at.loop = loop
}

// SetStart sets the initial offset value
func (at *AnimationTarget) SetStart(v float32) {

	at.start = v
}

// Update interpolates the specified input value for each animation target channel
// and executes its corresponding action function. Returns true if the input value
// is inside the key frames ranges or false otherwise.
func (at *AnimationTarget) Update(delta float32) bool {

	// Checks if input is less than minimum
	at.last = at.last + delta
	if at.last < at.minInput {
		return false
	}

	// Checks if input is greater than maximum
	if at.last > at.maxInput {
		if at.loop {
			at.Reset()
		} else {
			return false
		}
	}

	for i := 0; i < len(at.channels); i++ {
		ch := at.channels[i]
		// Get interpolated value
		v, ok := ch.sampler.Interpolate(at.last)
		if !ok {
			return false
		}
		// Call action func
		ch.action(at, v)
		// Sets final rotation
		at.target.GetNode().SetRotation(at.rot.X, at.rot.Y, at.rot.Z)
	}
	return true
}

// NewAnimationTargets creates and returns a map of all animation targets
// contained in the decoded Collada document and for the previously decoded scene.
// The map is indexed by the node loaderID.
func (d *Decoder) NewAnimationTargets(scene core.INode) (map[string]*AnimationTarget, error) {

	if d.dom.LibraryAnimations == nil {
		return nil, fmt.Errorf("No animations found")
	}

	// Maps target node to its animation target instance
	targetsMap := make(map[string]*AnimationTarget)

	// For each Collada animation element
	for _, ca := range d.dom.LibraryAnimations.Animation {

		// For each Collada channel for this animation
		for _, cc := range ca.Channel {

			// Separates the channel target in target id and target action
			parts := strings.Split(cc.Target, "/")
			if len(parts) < 2 {
				return nil, fmt.Errorf("Channel target invalid")
			}
			targetID := parts[0]
			targetAction := parts[1]

			// Get the target node object referenced by the target id from the specified scene.
			target := scene.GetNode().FindLoaderID(targetID)
			if target == nil {
				return nil, fmt.Errorf("Target node id:%s not found", targetID)
			}

			// Get reference to the AnimationTarget for this target in the local map
			// If not found creates the animation target and inserts in the map
			at := targetsMap[targetID]
			if at == nil {
				at = new(AnimationTarget)
				at.target = target
				at.matrix = target.GetNode().Matrix()
				targetsMap[targetID] = at
			}

			// Creates the sampler instance specified from the channel source
			si, err := NewSamplerInstance(ca, cc.Source)
			if err != nil {
				return nil, err
			}

			// Sets the action function from the target action
			var af ActionFunc
			switch targetAction {
			case "location.X":
				af = actionPositionX
			case "location.Y":
				af = actionPositionY
			case "location.Z":
				af = actionPositionZ
			case "rotationX.ANGLE":
				af = actionRotationX
			case "rotationY.ANGLE":
				af = actionRotationY
			case "rotationZ.ANGLE":
				af = actionRotationZ
			case "scale.X":
				af = actionScaleX
			case "scale.Y":
				af = actionScaleY
			case "scale.Z":
				af = actionScaleZ
			default:
				return nil, fmt.Errorf("Unsupported channel target action:%s", targetAction)
			}

			// Creates the channel instance for this sampler and target action and adds it
			// to the current AnimationTarget
			ci := &ChannelInstance{si, af}
			at.channels = append(at.channels, ci)
		}
	}
	// Set minimum and maximum input values for each animation target
	for _, at := range targetsMap {
		at.minInput = math32.Infinity
		at.maxInput = -math32.Infinity
		for _, ch := range at.channels {
			// First key frame input
			inp := ch.sampler.Input[0]
			if inp < at.minInput {
				at.minInput = inp
			}
			// Last key frame input
			inp = ch.sampler.Input[len(ch.sampler.Input)-1]
			if inp > at.maxInput {
				at.maxInput = inp
			}
		}
	}
	return targetsMap, nil
}

func actionPositionX(at *AnimationTarget, v float32) {

	at.target.GetNode().SetPositionX(v)
}

func actionPositionY(at *AnimationTarget, v float32) {

	at.target.GetNode().SetPositionY(v)
}

func actionPositionZ(at *AnimationTarget, v float32) {

	at.target.GetNode().SetPositionZ(v)
}

func actionRotationX(at *AnimationTarget, v float32) {

	at.rot.X = math32.DegToRad(v)
}

func actionRotationY(at *AnimationTarget, v float32) {

	at.rot.Y = math32.DegToRad(v)
}

func actionRotationZ(at *AnimationTarget, v float32) {

	at.rot.Z = math32.DegToRad(v)
}

func actionScaleX(at *AnimationTarget, v float32) {

	at.target.GetNode().SetScaleX(v)
}

func actionScaleY(at *AnimationTarget, v float32) {

	at.target.GetNode().SetScaleY(v)
}

func actionScaleZ(at *AnimationTarget, v float32) {

	at.target.GetNode().SetScaleZ(v)
}

// NewSamplerInstance creates and returns a pointer to a new SamplerInstance built
// with data from the specified Collada animation and URI
func NewSamplerInstance(ca *Animation, uri string) (*SamplerInstance, error) {

	id := strings.TrimPrefix(uri, "#")
	var cs *Sampler
	for _, current := range ca.Sampler {
		if current.Id == id {
			cs = current
			break
		}
	}
	if cs == nil {
		return nil, fmt.Errorf("Sampler:%s not found", id)
	}

	// Get sampler inputs
	si := new(SamplerInstance)
	for _, inp := range cs.Input {
		if inp.Semantic == "INPUT" {
			data, err := findSourceFloatArray(ca, inp.Source)
			if err != nil {
				return nil, err
			}
			si.Input = data
			continue
		}
		if inp.Semantic == "OUTPUT" {
			data, err := findSourceFloatArray(ca, inp.Source)
			if err != nil {
				return nil, err
			}
			si.Output = data
			continue
		}
		if inp.Semantic == "INTERPOLATION" {
			data, err := findSourceNameArray(ca, inp.Source)
			if err != nil {
				return nil, err
			}
			si.Interp = data
			continue
		}
		if inp.Semantic == "IN_TANGENT" {
			data, err := findSourceFloatArray(ca, inp.Source)
			if err != nil {
				return nil, err
			}
			si.InTangent = data
			continue
		}
		if inp.Semantic == "OUT_TANGENT" {
			data, err := findSourceFloatArray(ca, inp.Source)
			if err != nil {
				return nil, err
			}
			si.OutTangent = data
			continue
		}
	}
	return si, nil
}

// Interpolate returns the interpolated output and its validity
// for this sampler for the specified input.
func (si *SamplerInstance) Interpolate(inp float32) (float32, bool) {

	// Test limits
	if len(si.Input) < 2 {
		return 0, false
	}
	if inp < si.Input[0] {
		return 0, false
	}
	if inp > si.Input[len(si.Input)-1] {
		return 0, false
	}

	// Find key frame interval
	var idx int
	for idx = 0; idx < len(si.Input)-1; idx++ {
		if inp >= si.Input[idx] && inp < si.Input[idx+1] {
			break
		}
	}
	// Checks if interval was found
	if idx >= len(si.Input)-1 {
		return 0, false
	}

	switch si.Interp[idx] {
	case "STEP":
		return si.linearInterp(inp, idx), true
	case "LINEAR":
		return si.linearInterp(inp, idx), true
	case "BEZIER":
		return si.bezierInterp(inp, idx), true
	case "HERMITE":
		return si.linearInterp(inp, idx), true
	case "CARDINAL":
		return si.linearInterp(inp, idx), true
	case "BSPLINE":
		return si.linearInterp(inp, idx), true
	}

	return 0, false
}

func (si *SamplerInstance) linearInterp(inp float32, idx int) float32 {

	k1 := si.Input[idx]
	k2 := si.Input[idx+1]
	v1 := si.Output[idx]
	v2 := si.Output[idx+1]
	return v1 + (v2-v1)*(inp-k1)/(k2-k1)
}

func (si *SamplerInstance) bezierInterp(inp float32, idx int) float32 {

	p0 := si.Output[idx]
	p1 := si.Output[idx+1]
	c0 := si.OutTangent[2*idx+1]
	c1 := si.InTangent[2*(idx+1)+1]
	k1 := si.Input[idx]
	k2 := si.Input[idx+1]
	s := (inp - k1) / (k2 - k1)
	out := p0*math32.Pow(1-s, 3) + 3*c0*s*math32.Pow(1-s, 2) + 3*c1*s*s*(1-s) + p1*math32.Pow(s, 3)
	return out
}

func findSourceNameArray(ca *Animation, uri string) ([]string, error) {

	src := findSource(ca, uri)
	if src == nil {
		return nil, fmt.Errorf("Source:%s not found", uri)
	}
	na, ok := src.ArrayElement.(*NameArray)
	if !ok {
		return nil, fmt.Errorf("Source:%s is not NameArray", uri)
	}
	return na.Data, nil
}

func findSourceFloatArray(ca *Animation, uri string) ([]float32, error) {

	src := findSource(ca, uri)
	if src == nil {
		return nil, fmt.Errorf("Source:%s not found", uri)
	}
	fa, ok := src.ArrayElement.(*FloatArray)
	if !ok {
		return nil, fmt.Errorf("Source:%s is not FloatArray", uri)
	}
	return fa.Data, nil
}

func findSource(ca *Animation, uri string) *Source {

	id := strings.TrimPrefix(uri, "#")
	for _, src := range ca.Source {
		if src.Id == id {
			return src
		}
	}
	return nil
}
