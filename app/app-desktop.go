// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !wasm

package app

import (
	"fmt"
	"github.com/g3n/engine/audio/al"
	"github.com/g3n/engine/audio/vorbis"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/window"
	"time"
)

// Desktop application defaults
const (
	title  = "G3N Application"
	width  = 800
	height = 600
)

// Application
type Application struct {
	window.IWindow                    // Embedded GlfwWindow
	keyState       *window.KeyState   // Keep track of keyboard state
	renderer       *renderer.Renderer // Renderer object
	audioDev       *al.Device         // Default audio device
	frameStart     time.Time          // Frame start time
	frameDelta     time.Duration      // Duration of last frame
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
	err := window.Init(width, height, title)
	if err != nil {
		panic(err)
	}
	app.IWindow = window.Get()
	app.openDefaultAudioDevice()           // Set up audio
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

	// Initialize frame time
	app.frameStart = time.Now()

	// Set up recurring calls to user's update function
	for true {
		// If Exit() was called or there was an attempt to close the window dispatch OnExit event for subscribers.
		// If no subscriber cancelled the event, terminate the application.
		if app.IWindow.(*window.GlfwWindow).ShouldClose() {
			app.Dispatch(OnExit, nil)
			// TODO allow for cancelling exit e.g. showing dialog asking the user if he/she wants to save changes
			// if exit was cancelled {
			//     app.IWindow.(*window.GlfwWindow).SetShouldClose(false)
			// } else {
			break
			// }
		}
		// Update frame start and frame delta
		now := time.Now()
		app.frameDelta = now.Sub(app.frameStart)
		app.frameStart = now
		// Call user's update function
		update(app.renderer, app.frameDelta)
		// Swap buffers and poll events
		app.IWindow.(*window.GlfwWindow).SwapBuffers()
		app.IWindow.(*window.GlfwWindow).PollEvents()
	}

	// Close default audio device
	if app.audioDev != nil {
		al.CloseDevice(app.audioDev)
	}
	// Destroy window
	app.Destroy()
}

// Exit requests to terminate the application
// Application will dispatch OnQuit events to registered subscribers which
// can cancel the process by calling CancelDispatch().
func (app *Application) Exit() {

	app.IWindow.(*window.GlfwWindow).SetShouldClose(true)
}

// Renderer returns the application's renderer.
func (app *Application) Renderer() *renderer.Renderer {

	return app.renderer
}

// KeyState returns the application's KeyState.
func (app *Application) KeyState() *window.KeyState {

	return app.keyState
}

// openDefaultAudioDevice opens the default audio device setting it to the current context
func (app *Application) openDefaultAudioDevice() error {

	// Opens default audio device
	var err error
	app.audioDev, err = al.OpenDevice("")
	if err != nil {
		return fmt.Errorf("opening OpenAL default device: %s", err)
	}
	// Check for OpenAL effects extension support
	var attribs []int
	if al.IsExtensionPresent("ALC_EXT_EFX") {
		attribs = []int{al.MAX_AUXILIARY_SENDS, 4}
	}
	// Create audio context
	acx, err := al.CreateContext(app.audioDev, attribs)
	if err != nil {
		return fmt.Errorf("creating OpenAL context: %s", err)
	}
	// Makes the context the current one
	err = al.MakeContextCurrent(acx)
	if err != nil {
		return fmt.Errorf("setting OpenAL context current: %s", err)
	}
	// Logs audio library versions
	fmt.Println("LOGGING")
	log.Info("%s version: %s", al.GetString(al.Vendor), al.GetString(al.Version))
	log.Info("%s", vorbis.VersionString())
	return nil
}
