package gltf

import (
	"github.com/g3n/engine/material"
)

// TODO
// loadMaterialUnlit receives an interface value describing a KHR_materials_pbrSpecularGlossiness extension,
// decodes it and returns a Material closest to the specified description.
// The specification of this extension is at:
// https://github.com/KhronosGroup/glTF/tree/master/extensions/2.0/Khronos/KHR_materials_pbrSpecularGlossiness
func (g *GLTF) loadMaterialPbrSpecularGlossiness(ext interface{}) (material.IMaterial, error) {

	return nil, nil
}
