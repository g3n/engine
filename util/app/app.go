package app

import (
	"flag"
	"fmt"
	"runtime"

	"github.com/g3n/engine/audio/al"
	"github.com/g3n/engine/audio/ov"
	"github.com/g3n/engine/audio/vorbis"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/util/logger"
	"github.com/g3n/engine/window"
)

type App struct {
	win         window.IWindow       // Application window
	gl          *gls.GLS             // Application OpenGL state
	log         *logger.Logger       // Application logger
	renderer    *renderer.Renderer   // Renderer object
	camPersp    *camera.Perspective  // Perspective camera
	camOrtho    *camera.Orthographic // Orthographic camera
	audio       bool                 // Audio available
	vorbis      bool                 // Vorbis decoder available
	audioEFX    bool                 // Audio effect extension support available
	audioDev    *al.Device           // Audio player device
	captureDev  *al.Device           // Audio capture device
	frameRater  *FrameRater
	oCpuProfile *string
	scene       *core.Node // Node container for 3D tests
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

// appInstance points to the single Application instance
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
	log := logger.New(name, nil)
	log.AddWriter(logger.NewConsole(false))
	log.SetFormat(logger.FTIME | logger.FMICROS)
	log.SetLevel(ops.LogLevel)
	log.Info("%s v%d.%d starting", name, ops.VersionMajor, ops.VersionMinor)

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

	// Create OpenGL state
	gl, err := gls.New()
	if err != nil {
		return nil, err
	}

	// Creates application object
	app.win = win
	app.gl = gl

	// Creates perspective camera
	width, height := app.win.GetSize()
	aspect := float32(width) / float32(height)
	app.camPersp = camera.NewPerspective(65, aspect, 0.01, 1000)

	// Creates orthographic camera
	app.camOrtho = camera.NewOrthographic(-2, 2, 2, -2, 0.01, 100)
	app.camOrtho.SetPosition(0, 0, 3)
	app.camOrtho.LookAt(&math32.Vector3{0, 0, 0})
	app.camOrtho.SetZoom(1.0)

	// Creates renderer
	app.renderer = renderer.NewRenderer(gl)
	err = app.renderer.AddDefaultShaders()
	if err != nil {
		return nil, fmt.Errorf("Error from AddDefaulShaders:%v", err)
	}

	// Creates scene for 3D objects
	app.scene = core.NewNode()

	// Create frame rater
	app.frameRater = NewFrameRater(60)

	// Subscribe to window resize events
	app.win.Subscribe(window.OnWindowSize, func(evname string, ev interface{}) {
		app.OnWindowResize()
	})

	return app, nil

}

// Returns the single application instance
func Get() *App {

	return appInstance
}

// Window returns the application window
func (a *App) Window() window.IWindow {

	return a.win
}

func (a *App) Run() {

	for !a.win.ShouldClose() {
		// Starts measuring this frame
		a.frameRater.Start()

		// Poll input events and process them
		a.win.PollEvents()

		a.win.SwapBuffers()

		// Controls the frame rate and updates the FPS for the user
		a.frameRater.Wait()
	}
}

// OnWindowResize is called when the window resize event is received
func (app *App) OnWindowResize() {

	// Sets view port
	width, height := app.win.GetSize()
	app.gl.Viewport(0, 0, int32(width), int32(height))
	aspect := float32(width) / float32(height)

	// Sets camera aspect ratio
	app.camPersp.SetAspect(aspect)

	// Sets GUI root panel size
	//ctx.root.SetSize(float32(width), float32(height))
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
