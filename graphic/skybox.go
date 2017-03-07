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

type SkyboxData struct {
	DirAndPrefix string
	Extension    string
	Suffixes     [6]string
}

type Skybox struct {
	Graphic                     // embedded graphic object
	mvm     gls.UniformMatrix4f // model view matrix uniform
	mvpm    gls.UniformMatrix4f // model view projection matrix uniform
	nm      gls.UniformMatrix3f // normal matrix uniform
}

// NewSkybox creates and returns a pointer to a skybox with the specified textures
func NewSkybox(data SkyboxData) (*Skybox, error) {

	skybox := new(Skybox)

	geom := geometry.NewBox(50, 50, 50, 1, 1, 1)
	skybox.Graphic.Init(geom, gls.TRIANGLES)

	for i := 0; i < 6; i++ {
		tex, err := texture.NewTexture2DFromImage(data.DirAndPrefix + data.Suffixes[i] + "." + data.Extension)
		if err != nil {
			return nil, err
		}
		matFace := material.NewStandard(math32.NewColor(1, 1, 1))
		matFace.AddTexture(tex)
		matFace.SetSide(material.SideBack)
		matFace.SetUseLights(material.UseLightAmbient)
		skybox.AddGroupMaterial(skybox, matFace, i)
	}

	// Creates uniforms
	skybox.mvm.Init("ModelViewMatrix")
	skybox.mvpm.Init("MVP")
	skybox.nm.Init("NormalMatrix")

	return skybox, nil
}

// RenderSetup is called by the engine before drawing the skybox geometry
// It is responsible to updating the current shader uniforms with
// the model matrices.
func (skybox *Skybox) RenderSetup(gs *gls.GLS, rinfo *core.RenderInfo) {

	// TODO
	// Disable writes to the depth buffer (call glDepthMask(GL_FALSE)).
	// This will cause every other object to draw over the skybox, making it always appear "behind" everything else.
	// Since writes to the depth buffer are off, it doesn't matter how small the skybox is as long as it's larger than the camera's near clip plane.

	var mvm math32.Matrix4
	mvm.Copy(&rinfo.ViewMatrix)

	// Clear translation
	mvm[12] = 0
	mvm[13] = 0
	mvm[14] = 0
	// mvm.ExtractRotation(&rinfo.ViewMatrix) // TODO <- ExtractRotation does not work as expected?
	skybox.mvm.SetMatrix4(&mvm)
	skybox.mvm.Transfer(gs)

	// Calculates model view projection matrix and updates uniform
	var mvpm math32.Matrix4
	mvpm.MultiplyMatrices(&rinfo.ProjMatrix, &mvm)
	skybox.mvpm.SetMatrix4(&mvpm)
	skybox.mvpm.Transfer(gs)

	// Calculates normal matrix and updates uniform
	var nm math32.Matrix3
	nm.GetNormalMatrix(&mvm)
	skybox.nm.SetMatrix3(&nm)
	skybox.nm.Transfer(gs)

}
