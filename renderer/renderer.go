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
	gs          *gls.GLS
	shaman      Shaman                     // Internal shader manager
	scene       core.INode                 // Node containing 3D scene to render
	panelGui    gui.IPanel                 // Panel containing GUI to render
	panel3D     gui.IPanel                 // Panel which contains the 3D scene
	ambLights   []*light.Ambient           // Array of ambient lights for last scene
	dirLights   []*light.Directional       // Array of directional lights for last scene
	pointLights []*light.Point             // Array of point
	spotLights  []*light.Spot              // Array of spot lights for the scene
	others      []core.INode               // Other nodes (audio, players, etc)
	grmats      []*graphic.GraphicMaterial // Array of all graphic materials for scene
	rinfo       core.RenderInfo            // Preallocated Render info
	specs       ShaderSpecs                // Preallocated Shader specs
	redrawGui   bool
	needSwap    bool
	panList     []gui.IPanel
	panRendered int
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
	return r
}

func (r *Renderer) AddDefaultShaders() error {

	return r.shaman.AddDefaultShaders()
}

func (r *Renderer) AddChunk(name, source string) {

	r.shaman.AddChunk(name, source)
}

func (r *Renderer) AddShader(name, source string) {

	r.shaman.AddShader(name, source)
}

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

// SetScene sets the 3D scene to render
// If set to nil, no 3D scene will be rendered
func (r *Renderer) SetScene(scene core.INode) {

	r.scene = scene
}

func (r *Renderer) NeedSwap() bool {

	return r.needSwap
}

// Render renders the previously set Scene and Gui using the specified camera
func (r *Renderer) Render(icam camera.ICamera) error {

	r.redrawGui = false
	r.needSwap = false

	// Renders the 3D scene
	if r.scene != nil {
		err := r.renderScene(r.scene, icam)
		if err != nil {
			return err
		}
	}
	// Renders the Gui over the 3D scene
	if r.panelGui != nil {
		err := r.renderGui(icam)
		if err != nil {
			return err
		}
	}
	return nil
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
		r.needSwap = true
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
		}
		for idx, l := range r.dirLights {
			l.RenderSetup(r.gs, &r.rinfo, idx)
		}
		for idx, l := range r.pointLights {
			l.RenderSetup(r.gs, &r.rinfo, idx)
		}
		for idx, l := range r.spotLights {
			l.RenderSetup(r.gs, &r.rinfo, idx)
		}

		// Render this graphic material
		grmat.Render(r.gs, &r.rinfo)
	}
	return nil
}

// renderGui renders the Gui
func (r *Renderer) renderGui(icam camera.ICamera) error {

	// Updates panels bounds and relative positions
	parent := r.panelGui.GetPanel()
	parent.UpdateMatrixWorld()

	// If panel3D was defined and 3D scene is empty
	// Sets it renderable as the GUI background
	if r.panel3D != nil {
		if len(r.grmats) == 0 {
			r.panel3D.SetRenderable(true)
		} else {
			r.panel3D.SetRenderable(false)
		}
	}

	// Clears list of panels to render
	r.panList = r.panList[0:0]
	r.panRendered = 0
	// Redraw all GUI elements if necessary by appending the GUI panel to the render list
	if r.redrawGui || r.checkPanelUnbounded(r.panelGui) {
		r.panList = append(r.panList, r.panelGui)
	} else {
		r.buildPanelList(r.panelGui, true)
	}

	// If there are panels to render, disable the scissor test
	// which could have been set by the 3D scene renderer and
	// then clear the depth buffer, so the panels will be rendered
	// over the 3D scene.
	if len(r.panList) > 0 {
		r.gs.Disable(gls.SCISSOR_TEST)
		r.gs.Clear(gls.DEPTH_BUFFER_BIT)
		r.needSwap = true
	}

	// Render panels
	for i := 0; i < len(r.panList); i++ {
		err := r.renderPanel(r.panList[i])
		if err != nil {
			return err
		}
	}
	if r.panRendered > 0 {
		log.Error("panels rendered:%v", r.panRendered)
	}
	return nil
}

// buildPanelList builds list of panel materials that must be rendered.
func (r *Renderer) buildPanelList(ipan gui.IPanel, check3D bool) {

	// If panel is not visible, ignore it and all its children
	pan := ipan.GetPanel()
	if !pan.Visible() {
		return
	}

	// If this panel is not the main GUI panel checks if its over 3D.
	// If it is appends to the render list and we are done.
	// Otherwise do not check over 3D for its children because if the
	// parent is not over 3D than its children are also not over 3D.
	// The exception to this case is when the panel is unbounded,
	// which is verified in the next check.
	if ipan != r.panelGui {
		if check3D && r.checkPanelOver3D(ipan) {
			r.panList = append(r.panList, ipan)
			return
		} else {
			check3D = false
		}
	}

	// If panel is unbounded and is over 3D, appends to the render list
	if !pan.Bounded() && r.checkPanelOver3D(ipan) {
		r.panList = append(r.panList, ipan)
		return
	}
	// If this panel changed and is renderable or if any of its immediate children changed,
	// appends to the render list
	//if (pan.Changed() && pan.Renderable()) || r.checkPanelChildren(ipan) {
	if r.checkPanelChildren(ipan) {
		r.panList = append(r.panList, ipan)
		return
	}
	// Checks this panel children
	for _, ichild := range pan.Children() {
		r.buildPanelList(ichild.(gui.IPanel), check3D)
	}
}

// renderPanel renders the specified panel and all its children
// and then sets the panel as not changed.
func (r *Renderer) renderPanel(ipan gui.IPanel) error {

	// If panel not visible, ignore it and all its children
	pan := ipan.GetPanel()
	if !pan.Visible() {
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
		pan.SetChanged(false)
		r.panRendered++
	}
	// Renders this panel children
	for i := 0; i < len(pan.Children()); i++ {
		err := r.renderPanel(pan.Children()[i].(gui.IPanel))
		if err != nil {
			return err
		}
	}
	return nil
}

// checkPanelOver3D checks if the specified panel is over
// the area where the 3D scene will be rendered.
func (r *Renderer) checkPanelOver3D(ipan gui.IPanel) bool {

	if r.panel3D == nil || r.panel3D.Renderable() {
		return false
	}
	pan := ipan.GetPanel()
	if !pan.Visible() {
		return false
	}
	if r.panel3D.GetPanel().Intersects(pan) {
		//log.Error("panel..: %v %v %v", pan.Pospix(), pan.Width(), pan.Height())
		//panel3D := r.panel3D.GetPanel()
		//log.Error("panel3D: %v %v %v", panel3D.Pospix(), panel3D.Width(), panel3D.Height())
		return true
	}
	return false
}

// checkPanelUnbounded checks if the specified panel or any
// of its children has changed and it is unbounded
func (r *Renderer) checkPanelUnbounded(ipan gui.IPanel) bool {

	pan := ipan.GetPanel()
	if pan.Changed() && !pan.Bounded() {
		pan.SetChanged(false)
		return true
	}
	for _, ichild := range pan.Children() {
		if r.checkPanelUnbounded(ichild.(gui.IPanel)) {
			return true
		}
	}
	return false
}

// checkPanelChildren checks if any of this panel immediate children has changed
func (r *Renderer) checkPanelChildren(ipan gui.IPanel) bool {

	pan := ipan.GetPanel()
	for _, ichild := range pan.Children() {
		child := ichild.(gui.IPanel).GetPanel()
		if !child.Visible() || !child.Renderable() {
			continue
		}
		if child.Changed() {
			return true
		}
	}
	return false
}
