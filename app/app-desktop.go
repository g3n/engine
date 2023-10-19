// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !wasm
// +build !wasm

package app

import (
	"fmt"
	"time"

	"github.com/xackery/engine/renderer"
	"github.com/xackery/engine/window"
)

// Application
type Application struct {
	window.IWindow                    // Embedded GlfwWindow
	keyState       *window.KeyState   // Keep track of keyboard state
	renderer       *renderer.Renderer // Renderer object
	startTime      time.Time          // Application start time
	frameStart     time.Time          // Frame start time
	frameDelta     time.Duration      // Duration of last frame
}

// App returns the Application singleton, creating it the first time.
func App(width, height int, title string) *Application {

	// Return singleton if already created
	if a != nil {
		return a
	}
	a = new(Application)
	// Initialize window
	err := window.Init(width, height, title)
	if err != nil {
		panic(err)
	}
	a.IWindow = window.Get()
	a.keyState = window.NewKeyState(a) // Create KeyState
	// Create renderer and add default shaders
	a.renderer = renderer.NewRenderer(a.Gls())
	err = a.renderer.AddDefaultShaders()
	if err != nil {
		panic(fmt.Errorf("AddDefaultShaders:%v", err))
	}
	return a
}

// Run starts the update loop.
// It calls the user-provided update function every frame.
func (a *Application) Run(update func(rend *renderer.Renderer, deltaTime time.Duration)) {

	// Initialize start and frame time
	a.startTime = time.Now()
	a.frameStart = time.Now()

	// Set up recurring calls to user's update function
	for {
		// If Exit() was called or there was an attempt to close the window dispatch OnExit event for subscribers.
		// If no subscriber cancelled the event, terminate the application.
		if a.IWindow.(*window.GlfwWindow).ShouldClose() {
			a.Dispatch(OnExit, nil)
			// TODO allow for cancelling exit e.g. showing dialog asking the user if he/she wants to save changes
			// if exit was cancelled {
			//     a.IWindow.(*window.GlfwWindow).SetShouldClose(false)
			// } else {
			break
			// }
		}
		// Update frame start and frame delta
		now := time.Now()
		a.frameDelta = now.Sub(a.frameStart)
		a.frameStart = now
		// Call user's update function
		update(a.renderer, a.frameDelta)
		// Swap buffers and poll events
		a.IWindow.(*window.GlfwWindow).SwapBuffers()
		a.IWindow.(*window.GlfwWindow).PollEvents()
	}

	// Destroy window
	a.Destroy()
}

// Exit requests to terminate the application
// Application will dispatch OnQuit events to registered subscribers which
// can cancel the process by calling CancelDispatch().
func (a *Application) Exit() {

	a.IWindow.(*window.GlfwWindow).SetShouldClose(true)
}

// Renderer returns the application's renderer.
func (a *Application) Renderer() *renderer.Renderer {

	return a.renderer
}

// KeyState returns the application's KeyState.
func (a *Application) KeyState() *window.KeyState {

	return a.keyState
}

// RunTime returns the elapsed duration since the call to Run().
func (a *Application) RunTime() time.Duration {

	return time.Since(a.startTime)
}
