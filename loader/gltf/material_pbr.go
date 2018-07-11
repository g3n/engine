package gltf

import (
	"fmt"

	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/texture"
)

func (g *GLTF) loadMaterialPBR(m *Material) (material.IMaterial, error) {

	// Get pbr information
	pbr := m.PbrMetallicRoughness
	if pbr == nil {
		return nil, fmt.Errorf("PbrMetallicRoughness not supplied")
	}

	// Create new physically based material
	pm := material.NewPhysical()

	// TODO emmisive factor, emmissive map, occlusion, etc...

	// BaseColorFactor
	var baseColorFactor math32.Color4
	if pbr.BaseColorFactor != nil {
		baseColorFactor = math32.Color4{pbr.BaseColorFactor[0], pbr.BaseColorFactor[1], pbr.BaseColorFactor[2], pbr.BaseColorFactor[3]}
	} else {
		baseColorFactor = math32.Color4{1,1,1,1}
	}
	pm.SetBaseColorFactor(&baseColorFactor)

	// MetallicFactor
	var metallicFactor float32
	if pbr.MetallicFactor != nil {
		metallicFactor = *pbr.MetallicFactor
	} else {
		metallicFactor = 1
	}
	pm.SetMetallicFactor(metallicFactor)

	// RoughnessFactor
	var roughnessFactor float32
	if pbr.RoughnessFactor != nil {
		roughnessFactor = *pbr.RoughnessFactor
	} else {
		roughnessFactor = 1
	}
	pm.SetRoughnessFactor(roughnessFactor)

	// BaseColorTexture
	var tex *texture.Texture2D
	var err error
	if pbr.BaseColorTexture != nil {
		tex, err = g.loadTextureInfo(pbr.BaseColorTexture)
		if err != nil {
			return nil, err
		}
		pm.SetBaseColorMap(tex)
	}

	return pm, nil
}
