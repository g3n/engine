package gltf

import (
	"fmt"

	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/texture"
)

func (g *GLTF) loadMaterialPBR(m *Material) (material.IMaterial, error) {

	// Currently simulating PBR material with our common materials
	pbr := m.PbrMetallicRoughness
	if pbr == nil {
		return nil, fmt.Errorf("PbrMetallicRoughness not supplied")
	}
	pm := material.NewPhong(&math32.Color{pbr.BaseColorFactor[0], pbr.BaseColorFactor[1], pbr.BaseColorFactor[2]})
	pm.SetAmbientColor(&math32.Color{1, 1, 1})
	pm.SetEmissiveColor(&math32.Color{0, 0, 0})
	//pm.SetSpecularColor(&math32.Color{0, 0, 0})
	//pm.SetShininess(0)
	//pm.SetOpacity(1)

	// BaseColorTexture
	var tex *texture.Texture2D
	var err error
	if pbr.BaseColorTexture != nil {
		tex, err = g.loadTextureInfo(pbr.BaseColorTexture)
		if err != nil {
			return nil, err
		}
		pm.AddTexture(tex)
	}

	return pm, nil
}
