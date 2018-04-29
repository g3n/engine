// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphic

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/texture"
)

// Skybox is the Graphic that represents a skybox
type Skybox struct {
	Graphic             // embedded graphic object
	uniMVM  gls.Uniform // model view matrix uniform location cache
	uniMVPM gls.Uniform // model view projection matrix uniform cache
	uniNM   gls.Uniform // normal matrix uniform cache
}

// SkyboxData contains the data necessary to locate the textures for a Skybox in a concise manner.
type SkyboxData struct {
	DirAndPrefix string
	Extension    string
	Suffixes     [6]string
}

// NewSkybox creates and returns a pointer to a Skybox with the specified textures.
func NewSkybox(data SkyboxData) (*Skybox, error) {

	skybox := new(Skybox)

	geom := geometry.NewCube(1)
	skybox.Graphic.Init(geom, gls.TRIANGLES)
	skybox.Graphic.SetCullable(false)

	for i := 0; i < 6; i++ {
		tex, err := texture.NewTexture2DFromImage(data.DirAndPrefix + data.Suffixes[i] + "." + data.Extension)
		if err != nil {
			return nil, err
		}
		matFace := material.NewStandard(math32.NewColor("white"))
		matFace.AddTexture(tex)
		matFace.SetSide(material.SideBack)
		matFace.SetUseLights(material.UseLightAmbient)

		// Disable writes to the depth buffer (call glDepthMask(GL_FALSE)).
		// This will cause every other object to draw over the skybox, making it always appear behind everything else.
		// It doesn't matter how small/big the skybox is as long as it's visible by the camera (within near/far planes).
		matFace.SetDepthMask(false)

		skybox.AddGroupMaterial(skybox, matFace, i)
	}

	// Creates uniforms
	skybox.uniMVM.Init("ModelViewMatrix")
	skybox.uniMVPM.Init("MVP")
	skybox.uniNM.Init("NormalMatrix")

	return skybox, nil
}

// RenderSetup is called by the engine before drawing the skybox geometry
// It is responsible to updating the current shader uniforms with
// the model matrices.
func (skybox *Skybox) RenderSetup(gs *gls.GLS, rinfo *core.RenderInfo) {

	var mvm math32.Matrix4
	mvm.Copy(&rinfo.ViewMatrix)

	// Clear translation
	mvm[12] = 0
	mvm[13] = 0
	mvm[14] = 0
	// mvm.ExtractRotation(&rinfo.ViewMatrix) // TODO <- ExtractRotation does not work as expected?

	// Transfer mvp uniform
	location := skybox.uniMVM.Location(gs)
	gs.UniformMatrix4fv(location, 1, false, &mvm[0])

	// Calculates model view projection matrix and updates uniform
	var mvpm math32.Matrix4
	mvpm.MultiplyMatrices(&rinfo.ProjMatrix, &mvm)
	location = skybox.uniMVPM.Location(gs)
	gs.UniformMatrix4fv(location, 1, false, &mvpm[0])

	// Calculates normal matrix and updates uniform
	var nm math32.Matrix3
	nm.GetNormalMatrix(&mvm)
	location = skybox.uniNM.Location(gs)
	gs.UniformMatrix3fv(location, 1, false, &nm[0])
}
