// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package renderer

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/renderer/shader"
)

type ShaderSpecs struct {
	Name             string             // Shader name
	Version          string             // GLSL version
	ShaderUnique     bool               // indicates if shader is independent of lights and textures
	UseLights        material.UseLights // Bitmask indicating which lights to consider
	AmbientLightsMax int                // Current number of ambient lights
	DirLightsMax     int                // Current Number of directional lights
	PointLightsMax   int                // Current Number of point lights
	SpotLightsMax    int                // Current Number of spot lights
	MatTexturesMax   int                // Current Number of material textures
}

type ProgSpecs struct {
	program *gls.Program // program object
	specs   ShaderSpecs  // associated specs
}

type Shaman struct {
	gs       *gls.GLS
	chunks   *template.Template            // template with all chunks
	shaders  map[string]*template.Template // maps shader name to its template
	proginfo map[string]shader.ProgramInfo // maps name of the program to ProgramInfo
	programs []ProgSpecs                   // list of compiled programs with specs
	specs    ShaderSpecs                   // Current shader specs
}

// NewShaman creates and returns a pointer to a new shader manager
func NewShaman(gs *gls.GLS) *Shaman {

	sm := new(Shaman)
	sm.Init(gs)
	return sm
}

// Init initializes the shander manager
func (sm *Shaman) Init(gs *gls.GLS) {

	sm.gs = gs
	sm.chunks = template.New("_chunks_")
	sm.shaders = make(map[string]*template.Template)
	sm.proginfo = make(map[string]shader.ProgramInfo)

	// Add "loop" function to chunks template
	// "loop" is used inside the shader templates to unroll loops.
	sm.chunks.Funcs(template.FuncMap{
		"loop": func(n int) []int {
			s := make([]int, n)
			for i := range s {
				s[i] = i
			}
			return s
		},
	})
}

// AddDefaultShaders adds to the shader manager all default
// shaders statically registered.
func (sm *Shaman) AddDefaultShaders() error {

	for name, source := range shader.Chunks() {
		err := sm.AddChunk(name, source)
		if err != nil {
			return err
		}
	}

	for name, source := range shader.Shaders() {
		err := sm.AddShader(name, source)
		if err != nil {
			return err
		}
	}

	for name, pinfo := range shader.Programs() {
		sm.proginfo[name] = pinfo
	}
	return nil
}

// AddChunk adds a shader chunk with the specified name and source code
func (sm *Shaman) AddChunk(name, source string) error {

	tmpl := sm.chunks.New(name)
	_, err := tmpl.Parse(source)
	if err != nil {
		return err
	}
	return nil
}

// AddShader adds a shader program with the specified name and source code
func (sm *Shaman) AddShader(name, source string) error {

	// Clone chunks template so any shader can use
	// any of the chunks
	tmpl, err := sm.chunks.Clone()
	if err != nil {
		return err
	}
	// Parses this shader template source
	_, err = tmpl.Parse(source)
	if err != nil {
		return err
	}

	sm.shaders[name] = tmpl
	return nil
}

// AddProgram adds a program with the specified name and associated vertex
// and fragment shaders names (previously registered)
// To specify other types of shaders for a program use SetProgramShader()
func (sm *Shaman) AddProgram(name, vertexName, fragName string) error {

	sm.proginfo[name] = shader.ProgramInfo{Vertex: vertexName, Frag: fragName}
	return nil
}

// SetProgramShader sets the shader type and name for a previously specified program name.
// Returns error if the specified program or shader name not found or
// if an invalid shader type was specified.
func (sm *Shaman) SetProgramShader(pname string, stype int, sname string) error {

	// Checks if program name is valid
	pinfo, ok := sm.proginfo[pname]
	if !ok {
		return fmt.Errorf("Program name:%s not found", pname)
	}

	// Checks if shader name is valid
	_, ok = sm.shaders[sname]
	if !ok {
		return fmt.Errorf("Shader name:%s not found", sname)
	}

	// Sets the program shader name for the specified type
	switch stype {
	case gls.VERTEX_SHADER:
		pinfo.Vertex = sname
	case gls.FRAGMENT_SHADER:
		pinfo.Frag = sname
	case gls.GEOMETRY_SHADER:
		pinfo.Geometry = sname
	default:
		return fmt.Errorf("Invalid shader type")
	}
	return nil
}

// SetProgram set the shader program to satisfy the specified specs.
// Returns an indication if the current shader has changed and a possible error
// when creating a new shader program.
// Receives a copy of the specs because it changes the fields which specify the
// number of lights depending on the UseLights flags.
func (sm *Shaman) SetProgram(s *ShaderSpecs) (bool, error) {

	// Checks material use lights bit mask
	specs := *s
	if (specs.UseLights & material.UseLightAmbient) == 0 {
		specs.AmbientLightsMax = 0
	}
	if (specs.UseLights & material.UseLightDirectional) == 0 {
		specs.DirLightsMax = 0
	}
	if (specs.UseLights & material.UseLightPoint) == 0 {
		specs.PointLightsMax = 0
	}
	if (specs.UseLights & material.UseLightSpot) == 0 {
		specs.SpotLightsMax = 0
	}

	// If current shader specs are the same as the specified specs, nothing to do.
	if sm.specs.Compare(&specs) {
		return false, nil
	}

	// Search for compiled program with the specified specs
	for _, pinfo := range sm.programs {
		if pinfo.specs.Compare(&specs) {
			sm.gs.UseProgram(pinfo.program)
			sm.specs = specs
			return true, nil
		}
	}

	// Generates new program with the specified specs
	prog, err := sm.GenProgram(&specs)
	if err != nil {
		return false, err
	}
	log.Debug("Created new shader:%v", specs.Name)

	// Save specs as current specs, adds new program to the list
	// and actives program
	sm.specs = specs
	sm.programs = append(sm.programs, ProgSpecs{prog, specs})
	sm.gs.UseProgram(prog)
	return true, nil
}

// Generates shader program from the specified specs
func (sm *Shaman) GenProgram(specs *ShaderSpecs) (*gls.Program, error) {

	// Get info for the specified shader program
	progInfo, ok := sm.proginfo[specs.Name]
	if !ok {
		return nil, fmt.Errorf("Program:%s not found", specs.Name)
	}

	// Sets the GLSL version string
	specs.Version = "330 core"

	// Get vertex shader compiled template
	vtempl, ok := sm.shaders[progInfo.Vertex]
	if !ok {
		return nil, fmt.Errorf("Vertex shader:%s template not found", progInfo.Vertex)
	}
	// Generates vertex shader source from template
	var sourceVertex bytes.Buffer
	err := vtempl.Execute(&sourceVertex, specs)
	if err != nil {
		return nil, err
	}

	// Get fragment shader compiled template
	fragTempl, ok := sm.shaders[progInfo.Frag]
	if !ok {
		return nil, fmt.Errorf("Fragment shader:%s template not found", progInfo.Frag)
	}
	// Generates fragment shader source from template
	var sourceFrag bytes.Buffer
	err = fragTempl.Execute(&sourceFrag, specs)
	if err != nil {
		return nil, err
	}

	// Checks for optional geometry shader compiled template
	var sourceGeom bytes.Buffer
	if progInfo.Geometry != "" {
		// Get geometry shader compiled template
		geomTempl, ok := sm.shaders[progInfo.Geometry]
		if !ok {
			return nil, fmt.Errorf("Geometry shader:%s template not found", progInfo.Geometry)
		}
		// Generates geometry shader source from template
		err = geomTempl.Execute(&sourceGeom, specs)
		if err != nil {
			return nil, err
		}
	}

	// Creates shader program
	prog := sm.gs.NewProgram()
	prog.AddShader(gls.VERTEX_SHADER, sourceVertex.String(), nil)
	prog.AddShader(gls.FRAGMENT_SHADER, sourceFrag.String(), nil)
	if progInfo.Geometry != "" {
		prog.AddShader(gls.GEOMETRY_SHADER, sourceGeom.String(), nil)
	}
	err = prog.Build()
	if err != nil {
		return nil, err
	}
	return prog, nil
}

// Compare compares two shaders specifications structures
func (ss *ShaderSpecs) Compare(other *ShaderSpecs) bool {

	if ss.Name != other.Name {
		return false
	}
	if other.ShaderUnique {
		return true
	}
	if ss.AmbientLightsMax == other.AmbientLightsMax &&
		ss.DirLightsMax == other.DirLightsMax &&
		ss.PointLightsMax == other.PointLightsMax &&
		ss.SpotLightsMax == other.SpotLightsMax &&
		ss.MatTexturesMax == other.MatTexturesMax {
		return true
	}
	return false
}
