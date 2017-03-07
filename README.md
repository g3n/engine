# G3N - Go 3D Game Engine

G3N is a basic (for now!) OpenGL 3D Game Engine written in Go.
G3N was heavily inspired and based on the [three.js](https://threejs.org/) Javascript 3D library.
If you are curious about G3N and would like to see it in action
we recommend that you install the [G3N demo program](https://github.com/g3n/g3nd).

# Dependencies

The engine needs an OpenGL driver installed in the system and
on Unix like systems depends on some C libraries that can be installed using the platform package manager.
In all cases it is necessary to have gcc compatible C compiler installed.

* For Ubuntu/Debian-like Linux distributions, install `libgl1-mesa-dev` and `xorg-dev` packages.
* For CentOS/Fedora-like Linux distributions, install `libX11-devel libXcursor-devel libXrandr-devel libXinerama-devel mesa-libGL-devel libXi-devel` packages.
* Currently it was not tested on OS X. We encourage some feedback.
* For Windows we tested the build using the [mingw-w64](https://mingw-w64.org) toolchain.

G3N supports spatial audio using external libraries but loads these libraries
dinamically on demand, so you can install G3N and build a 3D application
(not using audio) without installing these libraries.

The following libraries are necessary for the optional audio support:

* For Ubuntu/Debian-like Linux distributions, install `libopenal1` and `libvorbisfile3`
* For CentOS/Fedora-like Linux distributions, install `libopenal1` and `libvorbisfile3 (to be verified)`
* Currently it was not tested on OS X. We encourage some feedback.
* For Windows its is necessary to install the following dlls:
  OpenAL32.dll, libogg.dll, libvorbis.dll and libvorbisfile.dll.

G3N was only tested with Go1.7.4+

# Installation

The following command will download the engine and all its dependencies, compile and
install the packages. Make sure your GOPATH is set correctly. 

`go get -u github.com/g3n/engine/...`

# Features

* Hierarchical scene graph. Any node can contain other nodes.
* Supports perspective and orthographic cameras. The camera can be controlled
  by the orbit control which allow zooming, rotation and panning using the mouse or keyboard.
* Suports ambient, directional, point and spot lights. Many lights can be added to the scene.
* Generators for primitive geometries such as: lines, box, sphere, cylinder and torus.
* Geometries can support multimaterials.
* Image textures can loaded from GIF, PNG or JPEG files and applied to materials.
* Loaders for the following 3D formats: Obj and Collada
* Text support allowing loading freetype fonts.
* Basic GUI supporting the widgets: label, image, button, checkbox, radiobutton,
  edit, scrollbar, slider, splitter, list, dropdown, tree, folder, window and layout managers
  (horizontal box, vertical box, grid, dock)
* Spatial audio support allowing playing sound from wave or Ogg Vorbis files.
* Users' applications can use their own vertex and fragment shaders.

# Basic application

The following code shows a basic G3N application 
([hellog3n](https://github.com/g3n/engine/tree/master/hellog3n))
which shows a wireframed sphere rotating.
You can install hellog3n using: `go get -u github.com/g3n/engine/hellog3n`

For more complex demos please see the [G3N demo program](https://github.com/g3n/g3nd).

```Go
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
```

<p align="center">
  <img style="float: right;" src="hellog3n/screenshot.png" alt="hellog3n Screenshot"/>
</p>

# To Do

G3N is a basic game engine. There is a lot of things to do.
We will soon insert here a list of the most important missing features.

# Documentation

For the engine API reference, please use
[![GoDoc](https://godoc.org/github.com/g3n/engine?status.svg)](https://godoc.org/github.com/g3n/engine).
We intend to add a user guide in the future.
Currently the best way to learn how to use the engine is to see the source code
of the demos from [G3ND](https://github.com/g3n/g3nd).

# Contributing

If you spot a bug or create a new feature you are encouraged to
send pull requests.


