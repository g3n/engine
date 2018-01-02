// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package renderer

import (
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/light"
)

// Renderer renders a 3D scene and/or a 2D GUI on the current window.
type Renderer struct {
	gs           *gls.GLS
	shaman       Shaman                     // Internal shader manager
	stats        Stats                      // Renderer statistics
	scene        core.INode                 // Node containing 3D scene to render
	panelGui     gui.IPanel                 // Panel containing GUI to render
	panel3D      gui.IPanel                 // Panel which contains the 3D scene
	ambLights    []*light.Ambient           // Array of ambient lights for last scene
	dirLights    []*light.Directional       // Array of directional lights for last scene
	pointLights  []*light.Point             // Array of point
	spotLights   []*light.Spot              // Array of spot lights for the scene
	others       []core.INode               // Other nodes (audio, players, etc)
	grmats       []*graphic.GraphicMaterial // Array of all graphic materials for scene
	rinfo        core.RenderInfo            // Preallocated Render info
	specs        ShaderSpecs                // Preallocated Shader specs
	redrawGui    bool                       // Flag indicating the gui must be redrawn completely
	rendered     bool                       // Flag indicating if anything was rendered
	panList      []gui.IPanel               // list of panels to render
	frameBuffers int                        // Number of frame buffers
	frameCount   int                        // Current number of frame buffers to write
}

// Stats describes how many object types were rendered
// It is cleared at the start of each render
type Stats struct {
	Graphics int // Number of graphic objects rendered
	Lights   int // Number of lights rendered
	Panels   int // Number of Gui panels rendered
	Others   int // Number of other objects rendered
}

// NewRenderer creates and returns a pointer to a new Renderer
func NewRenderer(gs *gls.GLS) *Renderer {

	r := new(Renderer)
	r.gs = gs
	r.shaman.Init(gs)

	r.ambLights = make([]*light.Ambient, 0)
	r.dirLights = make([]*light.Directional, 0)
	r.pointLights = make([]*light.Point, 0)
	r.spotLights = make([]*light.Spot, 0)
	r.others = make([]core.INode, 0)
	r.grmats = make([]*graphic.GraphicMaterial, 0)
	r.panList = make([]gui.IPanel, 0)
	r.frameBuffers = 2
	return r
}

// AddDefaultShaders adds to this renderer's shader manager all default
// include chunks, shaders and programs statically registered.
func (r *Renderer) AddDefaultShaders() error {

	return r.shaman.AddDefaultShaders()
}

// AddChunk adds a shader chunk with the specified name and source code
func (r *Renderer) AddChunk(name, source string) {

	r.shaman.AddChunk(name, source)
}

// AddShader adds a shader program with the specified name and source code
func (r *Renderer) AddShader(name, source string) {

	r.shaman.AddShader(name, source)
}

// AddProgram adds a program with the specified name and associated vertex
// and fragment shaders names (previously registered)
func (r *Renderer) AddProgram(name, vertex, frag string, others ...string) {

	r.shaman.AddProgram(name, vertex, frag, others...)
}

// SetGui sets the gui panel which contains the Gui to render.
// If set to nil, no Gui will be rendered
func (r *Renderer) SetGui(gui gui.IPanel) {

	r.panelGui = gui
}

// SetGuiPanel3D sets the gui panel inside which the 3D scene is shown.
// This informs the renderer that the Gui elements over this panel
// must be redrawn even if they didn't change.
// This panel panel must not be renderable, otherwise it will cover the 3D scene.
func (r *Renderer) SetGuiPanel3D(panel3D gui.IPanel) {

	r.panel3D = panel3D
}

// Panel3D returns the current gui panel over the 3D scene
func (r *Renderer) Panel3D() gui.IPanel {

	return r.panel3D
}

// SetScene sets the 3D scene to render
// If set to nil, no 3D scene will be rendered
func (r *Renderer) SetScene(scene core.INode) {

	r.scene = scene
}

// Returns statistics
func (r *Renderer) Stats() Stats {

	return r.stats
}

// Render renders the previously set Scene and Gui using the specified camera
// Returns an indication if anything was rendered and an error
func (r *Renderer) Render(icam camera.ICamera) (bool, error) {

	r.redrawGui = false
	r.rendered = false
	r.stats = Stats{}

	// Renders the 3D scene
	if r.scene != nil {
		err := r.renderScene(r.scene, icam)
		if err != nil {
			return r.rendered, err
		}
	}
	// Renders the Gui over the 3D scene
	if r.panelGui != nil {
		err := r.renderGui()
		if err != nil {
			return r.rendered, err
		}
	}
	return r.rendered, nil
}

// renderScene renders the 3D scene using the specified camera
func (r *Renderer) renderScene(iscene core.INode, icam camera.ICamera) error {

	// Updates world matrices of all scene nodes
	iscene.UpdateMatrixWorld()
	scene := iscene.GetNode()

	// Builds RenderInfo calls RenderSetup for all visible nodes
	icam.ViewMatrix(&r.rinfo.ViewMatrix)
	icam.ProjMatrix(&r.rinfo.ProjMatrix)

	// Clear scene arrays
	r.ambLights = r.ambLights[0:0]
	r.dirLights = r.dirLights[0:0]
	r.pointLights = r.pointLights[0:0]
	r.spotLights = r.spotLights[0:0]
	r.others = r.others[0:0]
	r.grmats = r.grmats[0:0]

	// Internal function to classify a node and its children
	var classifyNode func(inode core.INode)
	classifyNode = func(inode core.INode) {

		// If node not visible, ignore
		node := inode.GetNode()
		if !node.Visible() {
			return
		}

		// Checks if node is a Graphic
		igr, ok := inode.(graphic.IGraphic)
		if ok {
			if igr.Renderable() {
				// Appends to list each graphic material for this graphic
				gr := igr.GetGraphic()
				materials := gr.Materials()
				for i := 0; i < len(materials); i++ {
					r.grmats = append(r.grmats, &materials[i])
				}
			}
			// Node is not a Graphic
		} else {
			// Checks if node is a Light
			il, ok := inode.(light.ILight)
			if ok {
				switch l := il.(type) {
				case *light.Ambient:
					r.ambLights = append(r.ambLights, l)
				case *light.Directional:
					r.dirLights = append(r.dirLights, l)
				case *light.Point:
					r.pointLights = append(r.pointLights, l)
				case *light.Spot:
					r.spotLights = append(r.spotLights, l)
				default:
					panic("Invalid light type")
				}
				// Other nodes
			} else {
				r.others = append(r.others, inode)
			}
		}

		// Classify node children
		for _, ichild := range node.Children() {
			classifyNode(ichild)
		}
	}

	// Classify all scene nodes
	classifyNode(scene)

	// Sets lights count in shader specs
	r.specs.AmbientLightsMax = len(r.ambLights)
	r.specs.DirLightsMax = len(r.dirLights)
	r.specs.PointLightsMax = len(r.pointLights)
	r.specs.SpotLightsMax = len(r.spotLights)

	// Render other nodes (audio players, etc)
	for i := 0; i < len(r.others); i++ {
		inode := r.others[i]
		if !inode.GetNode().Visible() {
			continue
		}
		r.others[i].Render(r.gs)
		r.stats.Others++
	}

	// If there is graphic material to render
	if len(r.grmats) > 0 {
		// If the 3D scene to draw is to be confined to user specified panel
		// sets scissor to avoid erasing gui elements outside of this panel
		if r.panel3D != nil {
			pos := r.panel3D.GetPanel().Pospix()
			width, height := r.panel3D.GetPanel().Size()
			_, _, _, viewheight := r.gs.GetViewport()
			r.gs.Enable(gls.SCISSOR_TEST)
			r.gs.Scissor(int32(pos.X), viewheight-int32(pos.Y)-int32(height), uint32(width), uint32(height))
		} else {
			r.gs.Disable(gls.SCISSOR_TEST)
			r.redrawGui = true
		}
		// Clears the area inside the current scissor
		r.gs.Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		r.rendered = true
	}

	// For each *GraphicMaterial
	for _, grmat := range r.grmats {
		mat := grmat.GetMaterial().GetMaterial()

		// Sets the shader specs for this material and sets shader program
		r.specs.Name = mat.Shader()
		r.specs.ShaderUnique = mat.ShaderUnique()
		r.specs.UseLights = mat.UseLights()
		r.specs.MatTexturesMax = mat.TextureCount()
		_, err := r.shaman.SetProgram(&r.specs)
		if err != nil {
			return err
		}

		// Setup lights (transfer lights uniforms)
		for idx, l := range r.ambLights {
			l.RenderSetup(r.gs, &r.rinfo, idx)
			r.stats.Lights++
		}
		for idx, l := range r.dirLights {
			l.RenderSetup(r.gs, &r.rinfo, idx)
			r.stats.Lights++
		}
		for idx, l := range r.pointLights {
			l.RenderSetup(r.gs, &r.rinfo, idx)
			r.stats.Lights++
		}
		for idx, l := range r.spotLights {
			l.RenderSetup(r.gs, &r.rinfo, idx)
			r.stats.Lights++
		}

		// Render this graphic material
		grmat.Render(r.gs, &r.rinfo)
		r.stats.Graphics++
	}
	return nil
}

// renderGui renders the Gui
func (r *Renderer) renderGui() error {

	// If no 3D scene was rendered sets Gui panels as renderable for background
	// User must define the colors
	if len(r.grmats) == 0 {
		r.panelGui.SetRenderable(true)
		if r.panel3D != nil {
			r.panel3D.SetRenderable(true)
		}
	} else {
		r.panelGui.SetRenderable(false)
		if r.panel3D != nil {
			r.panel3D.SetRenderable(false)
		}
	}

	// Clears list of panels to render
	r.panList = r.panList[0:0]
	// Redraw all GUI elements elements (panel3D == nil and 3D scene drawn)
	if r.redrawGui {
		r.appendPanel(r.panelGui)
		// Redraw GUI elements only if changed
		// Set the number of frame buffers to draw these changes
	} else if r.checkChanged(r.panelGui) {
		r.appendPanel(r.panelGui)
		r.frameCount = r.frameBuffers
		// No change, but need to update frame buffers
	} else if r.frameCount > 0 {
		r.appendPanel(r.panelGui)
		// No change, draw only panels over 3D if any
	} else {
		r.getPanelsOver3D()
	}
	if len(r.panList) == 0 {
		return nil
	}

	// Updates panels bounds and relative positions
	r.panelGui.GetPanel().UpdateMatrixWorld()
	// Disable the scissor test which could have been set by the 3D scene renderer
	// and then clear the depth buffer, so the panels will be rendered over the 3D scene.
	r.gs.Disable(gls.SCISSOR_TEST)
	r.gs.Clear(gls.DEPTH_BUFFER_BIT)

	// Render panels
	for i := 0; i < len(r.panList); i++ {
		err := r.renderPanel(r.panList[i])
		if err != nil {
			return err
		}
	}
	r.frameCount--
	r.rendered = true
	return nil
}

// getPanelsOver3D builds list of panels over 3D to be rendered
func (r *Renderer) getPanelsOver3D() {

	// If panel3D not set or renderable, nothing to do
	if r.panel3D == nil || r.panel3D.Renderable() {
		return
	}

	// Internal recursive function to check if any child of the
	// specified panel is unbounded and over 3D.
	// If it is, it is inserted in the list of panels to render.
	var checkUnbounded func(pan *gui.Panel)
	checkUnbounded = func(pan *gui.Panel) {

		for i := 0; i < len(pan.Children()); i++ {
			child := pan.Children()[i].(gui.IPanel).GetPanel()
			if !child.Bounded() && r.checkPanelOver3D(child) {
				r.appendPanel(child)
				continue
			}
			checkUnbounded(child)
		}
	}

	// For all children of the Gui, checks if it is over the 3D panel
	children := r.panelGui.GetPanel().Children()
	for i := 0; i < len(children); i++ {
		pan := children[i].(gui.IPanel).GetPanel()
		if !pan.Visible() {
			continue
		}
		if r.checkPanelOver3D(pan) {
			r.appendPanel(pan)
			continue
		}
		// Current child is not over 3D but can have an unbounded child which is
		checkUnbounded(pan)
	}
}

// renderPanel renders the specified panel and all its children
// and then sets the panel as not changed.
func (r *Renderer) renderPanel(ipan gui.IPanel) error {

	// If panel not visible, ignore it and all its children
	pan := ipan.GetPanel()
	if !pan.Visible() {
		pan.SetChanged(false)
		return nil
	}
	// If panel is renderable, renders it
	if pan.Renderable() {
		// Sets shader program for the panel's material
		grmat := pan.GetGraphic().Materials()[0]
		mat := grmat.GetMaterial().GetMaterial()
		r.specs.Name = mat.Shader()
		r.specs.ShaderUnique = mat.ShaderUnique()
		_, err := r.shaman.SetProgram(&r.specs)
		if err != nil {
			return err
		}
		// Render this panel's graphic material
		grmat.Render(r.gs, &r.rinfo)
		r.stats.Panels++
	}
	pan.SetChanged(false)
	// Renders this panel children
	for i := 0; i < len(pan.Children()); i++ {
		err := r.renderPanel(pan.Children()[i].(gui.IPanel))
		if err != nil {
			return err
		}
	}
	return nil
}

// appendPanel appends the specified panel to the list of panels to render.
// Currently there is no need to check for duplicates.
func (r *Renderer) appendPanel(ipan gui.IPanel) {

	r.panList = append(r.panList, ipan)
}

// checkChanged checks if the specified panel or any of its children is changed
func (r *Renderer) checkChanged(ipan gui.IPanel) bool {

	// Unbounded panels are checked even if not visible
	pan := ipan.GetPanel()
	if !pan.Bounded() && pan.Changed() {
		pan.SetChanged(false)
		return true
	}
	// Ignore invisible panel and its children
	if !pan.Visible() {
		return false
	}
	if pan.Changed() && pan.Renderable() {
		return true
	}
	for i := 0; i < len(pan.Children()); i++ {
		res := r.checkChanged(pan.Children()[i].(gui.IPanel))
		if res {
			return res
		}
	}
	return false
}

// checkPanelOver3D checks if the specified panel is over
// the area where the 3D scene will be rendered.
func (r *Renderer) checkPanelOver3D(ipan gui.IPanel) bool {

	pan := ipan.GetPanel()
	if !pan.Visible() {
		return false
	}
	if r.panel3D.GetPanel().Intersects(pan) {
		return true
	}
	return false
}
