package app

import (
	"flag"
	"fmt"
	"runtime"

	"github.com/g3n/engine/audio/al"
	"github.com/g3n/engine/audio/ov"
	"github.com/g3n/engine/audio/vorbis"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/camera/control"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/util/logger"
	"github.com/g3n/engine/window"
)

type App struct {
	win         window.IWindow        // Application window
	gl          *gls.GLS              // OpenGL state
	log         *logger.Logger        // Default application logger
	renderer    *renderer.Renderer    // Renderer object
	camPersp    *camera.Perspective   // Perspective camera
	camOrtho    *camera.Orthographic  // Orthographic camera
	camera      camera.ICamera        // Current camera
	orbit       *control.OrbitControl // Camera orbit controller
	ambLight    *light.Ambient        // Default ambient light
	audio       bool                  // Audio available
	vorbis      bool                  // Vorbis decoder available
	audioEFX    bool                  // Audio effect extension support available
	audioDev    *al.Device            // Audio player device
	captureDev  *al.Device            // Audio capture device
	frameRater  *FrameRater           // Render loop frame rater
	scene       *core.Node            // Node container for 3D tests
	guiroot     *gui.Root             // Gui root panel
	oCpuProfile *string
}

type Options struct {
	WinHeight    int  // Initial window height. Uses screen width if 0
	WinWidth     int  // Initial window width. Uses screen height if 0
	VersionMajor int  // Application version major
	VersionMinor int  // Application version minor
	LogLevel     int  // Initial log level
	EnableFlags  bool // Enable command line flags
	TargetFPS    uint // Desired FPS
}

// appInstance contains the pointer to the single Application instance
var appInstance *App

// Creates creates and returns the application object using the specified name for
// the window title and log messages
// This functions must be called only once.
func Create(name string, ops Options) (*App, error) {

	if appInstance != nil {
		return nil, fmt.Errorf("Application already created")
	}
	app := new(App)

	if ops.EnableFlags {
		app.oCpuProfile = flag.String("cpuprofile", "", "Activate cpu profiling writing profile to the specified file")
		flag.Parse()
	}

	// Creates application logger
	app.log = logger.New(name, nil)
	app.log.AddWriter(logger.NewConsole(false))
	app.log.SetFormat(logger.FTIME | logger.FMICROS)
	app.log.SetLevel(ops.LogLevel)
	app.log.Info("%s v%d.%d starting", name, ops.VersionMajor, ops.VersionMinor)

	// Window event handling must run on the main OS thread
	runtime.LockOSThread()

	// Creates window and sets it as the current context
	win, err := window.New("glfw", 10, 10, name, false)
	if err != nil {
		return nil, err
	}
	// Sets the window size
	swidth, sheight := win.GetScreenResolution(nil)
	if ops.WinWidth != 0 {
		swidth = ops.WinWidth
	}
	if ops.WinHeight != 0 {
		sheight = ops.WinHeight
	}
	win.SetSize(swidth, sheight)
	app.win = win

	// Create OpenGL state
	gl, err := gls.New()
	if err != nil {
		return nil, err
	}
	app.gl = gl
	cc := math32.NewColor("gray")
	app.gl.ClearColor(cc.R, cc.G, cc.B, 1)
	app.gl.Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)

	// Creates perspective camera
	width, height := app.win.GetSize()
	aspect := float32(width) / float32(height)
	app.camPersp = camera.NewPerspective(65, aspect, 0.01, 1000)

	// Creates orthographic camera
	app.camOrtho = camera.NewOrthographic(-2, 2, 2, -2, 0.01, 100)
	app.camOrtho.SetPosition(0, 0, 3)
	app.camOrtho.LookAt(&math32.Vector3{0, 0, 0})
	app.camOrtho.SetZoom(1.0)

	// Default camera is perspective
	app.camera = app.camPersp

	// Creates orbit camera control
	// It is important to do this after the root panel subscription
	// to avoid GUI events being propagated to the orbit control.
	app.orbit = control.NewOrbitControl(app.camera, app.win)

	// Creates scene for 3D objects
	app.scene = core.NewNode()

	// Creates gui root panel
	app.guiroot = gui.NewRoot(app.gl, app.win)

	// Creates renderer
	app.renderer = renderer.NewRenderer(gl)
	err = app.renderer.AddDefaultShaders()
	if err != nil {
		return nil, fmt.Errorf("Error from AddDefaulShaders:%v", err)
	}
	app.renderer.SetScene(app.scene)
	app.renderer.SetGui(app.guiroot)

	// Adds white ambient light to the scene
	app.ambLight = light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 0.5)
	app.scene.Add(app.ambLight)

	// Create frame rater
	app.frameRater = NewFrameRater(60)

	// Subscribe to window resize events
	app.win.Subscribe(window.OnWindowSize, func(evname string, ev interface{}) {
		app.OnWindowResize()
	})

	return app, nil
}

// App returns the application single instance or nil
// if the application was not created yet
func Get() *App {

	return appInstance
}

// Log returns the application logger
func (a *App) Log() *logger.Logger {

	return a.log
}

// Window returns the application window
func (a *App) Window() window.IWindow {

	return a.win
}

// Gui returns the current application Gui root panel
func (app *App) Gui() *gui.Root {

	return app.guiroot
}

// Scene returns the current application 3D scene
func (a *App) Scene() *core.Node {

	return a.scene
}

// SetScene sets the 3D scene to be rendered
func (a *App) SetScene(scene *core.Node) {

	a.renderer.SetScene(scene)
}

// SetGui sets the root panel of th3 gui to be rendered
func (app *App) SetGui(root *gui.Root) {

	app.guiroot = root
	app.renderer.SetGui(app.guiroot)
}

// CameraPersp returns the application perspective camera
func (app *App) CameraPersp() *camera.Perspective {

	return app.camPersp
}

// Camera returns the current application camera
func (app *App) Camera() camera.ICamera {

	return app.camera
}

func (app *App) Renderer() *renderer.Renderer {

	return app.renderer
}

// Runs runs the application render loop
func (app *App) Run() error {

	for !app.win.ShouldClose() {
		// Starts measuring this frame
		app.frameRater.Start()

		// Renders the current scene and/or gui
		rendered, err := app.renderer.Render(app.camera)
		if err != nil {
			return err
		}
		app.log.Error("render stats:%+v", app.renderer.Stats())

		// Poll input events and process them
		app.win.PollEvents()

		if rendered {
			app.win.SwapBuffers()
		}

		// Controls the frame rate and updates the FPS for the user
		app.frameRater.Wait()
	}
	return nil
}

// Quit ends the application
func (app *App) Quit() {

	app.win.SetShouldClose(true)
}

// OnWindowResize is called when the window resize event is received
func (app *App) OnWindowResize() {

	// Get window size and sets the viewport to the same size
	width, height := app.win.GetSize()
	app.gl.Viewport(0, 0, int32(width), int32(height))

	// Sets perspective camera aspect ratio
	aspect := float32(width) / float32(height)
	app.camPersp.SetAspect(aspect)
	app.log.Error("app window resize:%v", aspect)

	// Sets the size of GUI root panel size to the size of the screen
	if app.guiroot != nil {
		app.guiroot.SetSize(float32(width), float32(height))
	}
}

// LoadAudioLibs try to load audio libraries
func (a *App) LoadAudioLibs() error {

	// Try to load OpenAL
	err := al.Load()
	if err != nil {
		return err
	}

	// Opens default audio device
	a.audioDev, err = al.OpenDevice("")
	if a.audioDev == nil {
		return fmt.Errorf("Error: %s opening OpenAL default device", err)
	}

	// Checks for OpenAL effects extension support
	if al.IsExtensionPresent("ALC_EXT_EFX") {
		a.audioEFX = true
	}

	// Creates audio context with auxiliary sends
	var attribs []int
	if a.audioEFX {
		attribs = []int{al.MAX_AUXILIARY_SENDS, 4}
	}
	acx, err := al.CreateContext(a.audioDev, attribs)
	if err != nil {
		return fmt.Errorf("Error creating audio context:%s", err)
	}

	// Makes the context the current one
	err = al.MakeContextCurrent(acx)
	if err != nil {
		return fmt.Errorf("Error setting audio context current:%s", err)
	}
	//log.Info("%s version: %s", al.GetString(al.Vendor), al.GetString(al.Version))
	a.audio = true

	// Ogg Vorbis support
	err = ov.Load()
	if err == nil {
		a.vorbis = true
		vorbis.Load()
	}
	return nil
}
