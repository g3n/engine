// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package geometry

import (
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/math32"
	"strconv"
	"sort"
)

// MorphGeometry represents a base geometry and its morph targets.
type MorphGeometry struct {
	baseGeometry  *Geometry   // The base geometry
	targets       []*Geometry // The morph target geometries
	weights       []float32   // The weights for each morph target
	uniWeights    gls.Uniform // Texture unit uniform location cache
	morphGeom     *Geometry   // Cache of the last CPU-morphed geometry
}

// NumMorphInfluencers is the maximum number of active morph targets.
const NumMorphTargets = 8

// NewMorphGeometry creates and returns a pointer to a new MorphGeometry.
func NewMorphGeometry(baseGeometry *Geometry) *MorphGeometry {

	mg := new(MorphGeometry)
	mg.baseGeometry = baseGeometry

	mg.targets = make([]*Geometry, 0)
	mg.weights = make([]float32, 0)

	mg.baseGeometry.ShaderDefines.Set("MORPHTARGETS", strconv.Itoa(NumMorphTargets))
	mg.uniWeights.Init("morphTargetInfluences")
	return mg
}

// GetGeometry satisfies the IGeometry interface.
func (mg *MorphGeometry) GetGeometry() *Geometry {

	return mg.baseGeometry
}

// SetWeights sets the morph target weights.
func (mg *MorphGeometry) SetWeights(weights []float32) {

	if len(weights) != len(mg.weights) {
		panic("weights have invalid length")
	}
	mg.weights = weights
}

// Weights returns the morph target weights.
func (mg *MorphGeometry) Weights() []float32 {

	return mg.weights
}

// Weights returns the morph target weights.
func (mg *MorphGeometry) AddMorphTargets(morphTargets ...*Geometry) {

	mg.targets = append(mg.targets, morphTargets...)
	for range morphTargets {
		mg.weights = append(mg.weights, 0)
	}
}

// ActiveMorphTargets sorts the morph targets by weight and returns the top n morph targets with largest weight.
func (mg *MorphGeometry) ActiveMorphTargets() ([]*Geometry, []float32) {

	numTargets := len(mg.targets)
	if numTargets == 0 {
		return nil, nil
	}

	sortedMorphTargets := make([]*Geometry, numTargets)
	copy(sortedMorphTargets, mg.targets)
	sort.Slice(sortedMorphTargets, func(i, j int) bool {
		return mg.weights[i] > mg.weights[j]
	})

	sortedWeights := make([]float32, numTargets)
	copy(sortedWeights, mg.weights)
	sort.Slice(sortedWeights, func(i, j int) bool {
		return mg.weights[i] > mg.weights[j]
	})


	// TODO check current 0 weights

	//if len(mg.targets) < NumMorphTargets-1 {
		return sortedMorphTargets, sortedWeights
	//} else {
	//	return sortedMorphTargets[:NumMorphTargets-1]
	//}

}

// SetIndices sets the indices array for this geometry.
func (mg *MorphGeometry) SetIndices(indices math32.ArrayU32) {

	mg.baseGeometry.SetIndices(indices)
	for i := range mg.targets {
		mg.targets[i].SetIndices(indices)
	}
}

// ComputeMorphed computes a morphed geometry from the provided morph target weights.
// Note that morphing is usually computed by the GPU in shaders.
// This CPU implementation allows users to obtain an instance of a morphed geometry
// if so desired (loosing morphing ability).
func (mg *MorphGeometry) ComputeMorphed(weights []float32) *Geometry {

	morphed := NewGeometry()
	// TODO
	return morphed
}

// Dispose releases, if possible, OpenGL resources, C memory
// and VBOs associated with the base geometry and morph targets.
func (mg *MorphGeometry) Dispose() {

	mg.baseGeometry.Dispose()
	for i := range mg.targets {
		mg.targets[i].Dispose()
	}
}

// RenderSetup is called by the renderer before drawing the geometry.
func (mg *MorphGeometry) RenderSetup(gs *gls.GLS) {

	mg.baseGeometry.RenderSetup(gs)

	// Sort weights and find top 8 morph targets with largest current weight (8 is the max sent to shader)
	activeMorphTargets, activeWeights := mg.ActiveMorphTargets()

	for i, mt := range activeMorphTargets {

		mt.SetAttributeName(gls.VertexPosition, "MorphPosition"+strconv.Itoa(i))
		mt.SetAttributeName(gls.VertexNormal, "MorphNormal"+strconv.Itoa(i))
		//mt.SetAttributeName(vTangent, fmt.Sprintf("MorphTangent[%d]", i))

		// Transfer morphed geometry VBOs
		for _, vbo := range mt.VBOs() {
			vbo.Transfer(gs)
		}
	}

	// Transfer texture info combined uniform
	location := mg.uniWeights.Location(gs)
	gs.Uniform1fv(location, int32(len(activeWeights)), activeWeights)
}
