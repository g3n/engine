// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build wasm

package app

import (
	"fmt"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/window"
	"syscall/js"
	"time"
)

// Default canvas id
const canvasId = "g3n-canvas"

// Application
type Application struct {
	window.IWindow                    // Embedded WebGLCanvas
	keyState       *window.KeyState   // Keep track of keyboard state
	renderer       *renderer.Renderer // Renderer object
	frameStart     time.Time          // Frame start time
	frameDelta     time.Duration      // Duration of last frame
	exit           bool
	cbid           js.Value
}

// Application singleton
var app *Application

// App returns the Application singleton, creating it the first time.
func App() *Application {

	// Return singleton if already created
	if app != nil {
		return app
	}
	app = new(Application)
	// Initialize window
	err := window.Init(canvasId)
	if err != nil {
		panic(err)
	}
	app.IWindow = window.Get()
	// TODO audio setup here
	app.keyState = window.NewKeyState(app) // Create KeyState
	// Create renderer and add default shaders
	app.renderer = renderer.NewRenderer(app.Gls())
	err = app.renderer.AddDefaultShaders()
	if err != nil {
		panic(fmt.Errorf("AddDefaultShaders:%v", err))
	}
	return app
}

// Run starts the update loop.
// It calls the user-provided update function every frame.
func (app *Application) Run(update func(renderer *renderer.Renderer, deltaTime time.Duration)) {

	// Create channel so later we can prevent application from finishing while we wait for callbacks
	done := make(chan bool, 0)

	// Initialize frame time
	app.frameStart = time.Now()

	// Set up recurring calls to user's update function
	var tick js.Func
	tick = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// Update frame start and frame delta
		now := time.Now()
		app.frameDelta = now.Sub(app.frameStart)
		app.frameStart = now
		// Call user's update function
		update(app.renderer, app.frameDelta)
		// Set up new callback if not exiting
		if !app.exit {
			app.cbid = js.Global().Call("requestAnimationFrame", tick)
		} else {
			done <- true // Write to done channel to exit the app
		}
		return nil
	})
	app.cbid = js.Global().Call("requestAnimationFrame", tick)

	// Read from done channel
	// This channel will be empty (except when we want to exit the app)
	// It keeps the app from finishing while we wait for the next call to tick()
	<-done

	// Destroy the window
	app.IWindow.Destroy()
}

// Exit exits the app.
func (app *Application) Exit() {

	app.exit = true
}

// Renderer returns the application's renderer.
func (app *Application) Renderer() *renderer.Renderer {

	return app.renderer
}

// KeyState returns the application's KeyState.
func (app *Application) KeyState() *window.KeyState {

	return app.keyState
}
