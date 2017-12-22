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

type Renderer struct {
	gs          *gls.GLS
	shaman      Shaman                     // Internal shader manager
	scene       core.INode                 // Node containing 3D scene to render
	gui         gui.IPanel                 // Panel containing GUI to render
	ambLights   []*light.Ambient           // Array of ambient lights for last scene
	dirLights   []*light.Directional       // Array of directional lights for last scene
	pointLights []*light.Point             // Array of point
	spotLights  []*light.Spot              // Array of spot lights for the scene
	others      []core.INode               // Other nodes (audio, players, etc)
	grmats      []*graphic.GraphicMaterial // Array of all graphic materials for scene
	rinfo       core.RenderInfo            // Preallocated Render info
	specs       ShaderSpecs                // Preallocated Shader specs
}

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

// SetGui sets the gui panel which contains the Gui to render
// over the optional 3D scene.
// If set to nil, no Gui will be rendered
func (r *Renderer) SetGui(gui gui.IPanel) {

	r.gui = gui
}

// SetScene sets the Node which contains the scene to render
// If set to nil, no scene will be rendered
func (r *Renderer) SetScene(scene core.INode) {

	r.scene = scene
}

// Render renders the previously set Scene and Gui using the specified camera
func (r *Renderer) Render(icam camera.ICamera) error {

	// Renders the 3D scene
	if r.scene != nil {
		r.gs.Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		err := r.renderScene(r.scene, icam)
		if err != nil {
			return err
		}
	}
	// Renders the Gui over the 3D scene
	if r.gui != nil {
		r.gs.Clear(gls.DEPTH_BUFFER_BIT)
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

	// For each *GraphicMaterial
	for _, grmat := range r.grmats {
		//log.Debug("grmat:%v", grmat)
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
	parent := r.gui.GetPanel()
	parent.UpdateMatrixWorld()

	// Builds RenderInfo calls RenderSetup for all visible nodes
	icam.ViewMatrix(&r.rinfo.ViewMatrix)
	icam.ProjMatrix(&r.rinfo.ProjMatrix)

	var buildRenderList func(ipan gui.IPanel)
	buildRenderList = func(ipan gui.IPanel) {
		pan := ipan.GetPanel()
		// If panel is not visible, ignore
		if !pan.Visible() {
			return
		}
		// Get panel graphic materials
		gr := pan.GetGraphic()
		materials := gr.Materials()
		for i := 0; i < len(materials); i++ {
			r.grmats = append(r.grmats, &materials[i])
		}
		// Get this panel children
		for _, ichild := range pan.Children() {
			buildRenderList(ichild.(gui.IPanel))
		}
	}

	// Builds list of panel graphic materials to render
	r.grmats = r.grmats[0:0]
	buildRenderList(parent)

	// For each *GraphicMaterial
	for _, grmat := range r.grmats {
		mat := grmat.GetMaterial().GetMaterial()

		// Sets the shader specs for this material and sets shader program
		r.specs.Name = mat.Shader()
		r.specs.ShaderUnique = mat.ShaderUnique()
		_, err := r.shaman.SetProgram(&r.specs)
		if err != nil {
			return err
		}
		// Render this graphic material
		grmat.Render(r.gs, &r.rinfo)
	}

	return nil
}
