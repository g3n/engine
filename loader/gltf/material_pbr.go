package gltf

import (
	"fmt"

	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
)

func (g *GLTF) loadMaterialPBR(m *Material) (material.IMaterial, error) {

	// Get pbr information
	pbr := m.PbrMetallicRoughness
	if pbr == nil {
		return nil, fmt.Errorf("PbrMetallicRoughness not supplied")
	}

	// Create new physically based material
	pm := material.NewPhysical()

	// Double sided
	if m.DoubleSided {
		pm.SetSide(material.SideDouble)
	} else {
		pm.SetSide(material.SideFront)
	}

	var alphaMode string
	if len(m.AlphaMode) > 0{
		alphaMode = m.AlphaMode
	} else {
		alphaMode = "OPAQUE"
	}

	if alphaMode == "BLEND" {
		pm.SetTransparent(true)
	} else {
		pm.SetTransparent(false)
		if alphaMode == "MASK" {
			// TODO m.AlphaCutoff
			// pm.SetAlphaCutoff
		}
	}

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
		if pbr.MetallicRoughnessTexture != nil {
			metallicFactor = 1
		} else {
			metallicFactor = 0
		}
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

	// EmissiveFactor
	var emissiveFactor math32.Color
	if m.EmissiveFactor != nil {
		emissiveFactor = math32.Color{m.EmissiveFactor[0], m.EmissiveFactor[1], m.EmissiveFactor[2]}
	} else {
		if m.EmissiveTexture != nil {
			emissiveFactor = math32.Color{1, 1, 1}
		} else {
			emissiveFactor = math32.Color{0,0,0}
		}
	}
	pm.SetEmissiveFactor(&emissiveFactor)

	// BaseColorTexture
	if pbr.BaseColorTexture != nil {
		tex, err := g.LoadTexture(pbr.BaseColorTexture.Index)
		if err != nil {
			return nil, err
		}
		pm.SetBaseColorMap(tex)
	}

	// MetallicRoughnessTexture
	if pbr.MetallicRoughnessTexture != nil {
		tex, err := g.LoadTexture(pbr.MetallicRoughnessTexture.Index)
		if err != nil {
			return nil, err
		}
		pm.SetMetallicRoughnessMap(tex)
	}

	// NormalTexture
	if m.NormalTexture != nil {
		tex, err := g.LoadTexture(m.NormalTexture.Index)
		if err != nil {
			return nil, err
		}
		pm.SetNormalMap(tex)
	}

	// OcclusionTexture
	if m.OcclusionTexture != nil {
		tex, err := g.LoadTexture(m.OcclusionTexture.Index)
		if err != nil {
			return nil, err
		}
		pm.SetOcclusionMap(tex)
	}

	// EmissiveTexture
	if m.EmissiveTexture != nil {
		tex, err := g.LoadTexture(m.EmissiveTexture.Index)
		if err != nil {
			return nil, err
		}
		pm.SetEmissiveMap(tex)
	}

	return pm, nil
}
