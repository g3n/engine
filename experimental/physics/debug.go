// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package physics implements a basic physics engine.
package physics

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/experimental/collision"
)

// This file contains helpful infrastructure for debugging physics
type DebugHelper struct {

}

func ShowWorldFace(scene *core.Node, face []math32.Vector3, color *math32.Color) {

	if len(face) == 0 {
		return
	}

	vertices := math32.NewArrayF32(0, 16)
	for i := range face {
		vertices.AppendVector3(&face[i])
	}
	vertices.AppendVector3(&face[0])

	geom := geometry.NewGeometry()
	geom.AddVBO(gls.NewVBO(vertices).AddAttrib(gls.VertexPosition))

	mat := material.NewStandard(color)
	faceGraphic := graphic.NewLineStrip(geom, mat)
	scene.Add(faceGraphic)
}

func ShowPenAxis(scene *core.Node, axis *math32.Vector3) {//}, min, max float32) {

	vertices := math32.NewArrayF32(0, 16)

	size := float32(100)
	minPoint := axis.Clone().MultiplyScalar(size)
	maxPoint := axis.Clone().MultiplyScalar(-size)
	//vertices.AppendVector3(minPoint.Clone().SetX(minPoint.X - size))
	//vertices.AppendVector3(minPoint.Clone().SetX(minPoint.X + size))
	//vertices.AppendVector3(minPoint.Clone().SetY(minPoint.Y - size))
	//vertices.AppendVector3(minPoint.Clone().SetY(minPoint.Y + size))
	//vertices.AppendVector3(minPoint.Clone().SetZ(minPoint.Z - size))
	//vertices.AppendVector3(minPoint.Clone().SetZ(minPoint.Z + size))
	vertices.AppendVector3(minPoint)
	//vertices.AppendVector3(maxPoint.Clone().SetX(maxPoint.X - size))
	//vertices.AppendVector3(maxPoint.Clone().SetX(maxPoint.X + size))
	//vertices.AppendVector3(maxPoint.Clone().SetY(maxPoint.Y - size))
	//vertices.AppendVector3(maxPoint.Clone().SetY(maxPoint.Y + size))
	//vertices.AppendVector3(maxPoint.Clone().SetZ(maxPoint.Z - size))
	//vertices.AppendVector3(maxPoint.Clone().SetZ(maxPoint.Z + size))
	vertices.AppendVector3(maxPoint)

	geom := geometry.NewGeometry()
	geom.AddVBO(gls.NewVBO(vertices).AddAttrib(gls.VertexPosition))

	mat := material.NewStandard(&math32.Color{1,1,1})
	faceGraphic := graphic.NewLines(geom, mat)
	scene.Add(faceGraphic)
}

func ShowContact(scene *core.Node, contact *collision.Contact) {

	vertices := math32.NewArrayF32(0, 16)

	size := float32(0.0005)
	otherPoint := contact.Point.Clone().Add(contact.Normal.Clone().MultiplyScalar(-contact.Depth))
	vertices.AppendVector3(contact.Point.Clone().SetX(contact.Point.X - size))
	vertices.AppendVector3(contact.Point.Clone().SetX(contact.Point.X + size))
	vertices.AppendVector3(contact.Point.Clone().SetY(contact.Point.Y - size))
	vertices.AppendVector3(contact.Point.Clone().SetY(contact.Point.Y + size))
	vertices.AppendVector3(contact.Point.Clone().SetZ(contact.Point.Z - size))
	vertices.AppendVector3(contact.Point.Clone().SetZ(contact.Point.Z + size))
	vertices.AppendVector3(contact.Point.Clone())
	vertices.AppendVector3(otherPoint)

	geom := geometry.NewGeometry()
	geom.AddVBO(gls.NewVBO(vertices).AddAttrib(gls.VertexPosition))

	mat := material.NewStandard(&math32.Color{0,0,1})
	faceGraphic := graphic.NewLines(geom, mat)
	scene.Add(faceGraphic)
}