// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This is a minimum G3N application showing how to create a window,
// a scene, add some 3D objects to the scene and render it.
// For more complete demos please see: https://github.com/g3n/g3nd
package main

import (
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/window"
	"math"
	"runtime"
)

func main() {

	// Creates window and OpenGL context
	win, err := window.New("glfw", 800, 600, "Hello G3N", false)
	if err != nil {
		panic(err)
	}

	// OpenGL functions must be executed in the same thread where
	// the context was created (by window.New())
	runtime.LockOSThread()

	// Create OpenGL state
	gs, err := gls.New()
	if err != nil {
		panic(err)
	}

	// Creates scene for 3D objects
	scene := core.NewNode()

	// Adds white ambient light to the scene
	ambLight := light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 0.5)
	scene.Add(ambLight)

	// Adds a perspective camera to the scene
	width, height := win.GetSize()
	aspect := float32(width) / float32(height)
	camera := camera.NewPerspective(65, aspect, 0.01, 1000)
	camera.SetPosition(0, 0, 5)

	// Add an axis helper
	axis := graphic.NewAxisHelper(2)
	scene.Add(axis)

	// Creates a wireframe sphere positioned at the center of the scene
	geom := geometry.NewSphere(2, 16, 16, 0, math.Pi*2, 0, math.Pi)
	mat := material.NewStandard(math32.NewColor(1, 1, 1))
	mat.SetSide(material.SideDouble)
	mat.SetWireframe(true)
	sphere := graphic.NewMesh(geom, mat)
	scene.Add(sphere)

	// Creates a renderer and adds default shaders
	rend := renderer.NewRenderer(gs)
	err = rend.AddDefaultShaders()
	if err != nil {
		panic(err)
	}

	// Sets window background color
	gs.ClearColor(0, 0, 0, 1.0)

	// Render loop
	for !win.ShouldClose() {

		// Clear buffers
		gs.Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)

		// Rotates the sphere a bit around the Z axis (up)
		sphere.AddRotationY(0.005)

		// Render the scene using the specified camera
		rend.Render(scene, camera)

		// Update window and checks for I/O events
		win.SwapBuffers()
		win.PollEvents()
	}
}
