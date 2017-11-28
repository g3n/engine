// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package renderer

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/renderer/shaders"
)

// Regular expression to parse #include <name> directive
var rexInclude *regexp.Regexp

func init() {

	rexInclude = regexp.MustCompile(`#include\s+<(.*)>`)
}

// ShaderSpecs describes the specification of a compiled shader program
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
	includes map[string]string              // include files sources
	shadersm map[string]string              // maps shader name to its template
	proginfo map[string]shaders.ProgramInfo // maps name of the program to ProgramInfo
	programs []ProgSpecs                    // list of compiled programs with specs
	specs    ShaderSpecs                    // Current shader specs
}

// NewShaman creates and returns a pointer to a new shader manager
func NewShaman(gs *gls.GLS) *Shaman {

	sm := new(Shaman)
	sm.Init(gs)
	return sm
}

// Init initializes the shader manager
func (sm *Shaman) Init(gs *gls.GLS) {

	sm.gs = gs
	sm.includes = make(map[string]string)
	sm.shadersm = make(map[string]string)
	sm.proginfo = make(map[string]shaders.ProgramInfo)
}

// AddDefaultShaders adds to this shader manager all default
// include chunks, shaders and programs statically registered.
func (sm *Shaman) AddDefaultShaders() error {

	for _, name := range shaders.Includes() {
		sm.AddChunk(name, shaders.IncludeSource(name))
	}

	for _, name := range shaders.Shaders() {
		sm.AddShader(name, shaders.ShaderSource(name))
	}

	for _, name := range shaders.Programs() {
		sm.proginfo[name] = shaders.GetProgramInfo(name)
	}
	return nil
}

// AddChunk adds a shader chunk with the specified name and source code
func (sm *Shaman) AddChunk(name, source string) {

	sm.includes[name] = source
}

// AddShader adds a shader program with the specified name and source code
func (sm *Shaman) AddShader(name, source string) {

	sm.shadersm[name] = source
}

// AddProgram adds a program with the specified name and associated vertex
// and fragment shaders names (previously registered)
func (sm *Shaman) AddProgram(name, vertexName, fragName string, others ...string) {

	geomName := ""
	if len(others) > 0 {
		geomName = others[0]
	}
	sm.proginfo[name] = shaders.ProgramInfo{
		Vertex:   vertexName,
		Fragment: fragName,
		Geometry: geomName,
	}
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

	// Save specs as current specs, adds new program to the list and activates the program
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

	// Sets the defines map
	defines := map[string]string{}
	defines["GLSL_VERSION"] = "330 core"
	defines["AMB_LIGHTS"] = strconv.FormatUint(uint64(specs.AmbientLightsMax), 10)
	defines["DIR_LIGHTS"] = strconv.FormatUint(uint64(specs.DirLightsMax), 10)
	defines["POINT_LIGHTS"] = strconv.FormatUint(uint64(specs.PointLightsMax), 10)
	defines["SPOT_LIGHTS"] = strconv.FormatUint(uint64(specs.SpotLightsMax), 10)
	defines["MAT_TEXTURES"] = strconv.FormatUint(uint64(specs.MatTexturesMax), 10)

	// Get vertex shader source
	vertexSource, ok := sm.shadersm[progInfo.Vertex]
	if !ok {
		return nil, fmt.Errorf("Vertex shader:%s not found", progInfo.Vertex)
	}
	// Preprocess vertex shader source
	vertexSource, err := sm.preprocess(vertexSource, defines)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("vertexSource:%s\n", vertexSource)

	// Get fragment shader source
	fragSource, ok := sm.shadersm[progInfo.Fragment]
	if err != nil {
		return nil, fmt.Errorf("Fragment shader:%s not found", progInfo.Fragment)
	}
	// Preprocess fragment shader source
	fragSource, err = sm.preprocess(fragSource, defines)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("fragSource:%s\n", fragSource)

	// Checks for optional geometry shader compiled template
	var geomSource = ""
	if progInfo.Geometry != "" {
		// Get geometry shader source
		geomSource, ok = sm.shadersm[progInfo.Geometry]
		if !ok {
			return nil, fmt.Errorf("Geometry shader:%s not found", progInfo.Geometry)
		}
		// Preprocess geometry shader source
		geomSource, err = sm.preprocess(geomSource, defines)
		if err != nil {
			return nil, err
		}
	}

	// Creates shader program
	prog := sm.gs.NewProgram()
	prog.AddShader(gls.VERTEX_SHADER, vertexSource, nil)
	prog.AddShader(gls.FRAGMENT_SHADER, fragSource, nil)
	if progInfo.Geometry != "" {
		prog.AddShader(gls.GEOMETRY_SHADER, geomSource, nil)
	}
	err = prog.Build()
	if err != nil {
		return nil, err
	}
	return prog, nil
}

// preprocess preprocesses the specified source prefixing it with optional defines directives
// contained in "defines" parameter and replaces '#include <name>' directives
// by the respective source code of include chunk of the specified name.
// The included "files" are also processed recursively.
func (sm *Shaman) preprocess(source string, defines map[string]string) (string, error) {

	// If defines map supplied, generates prefix with glsl version directive first,
	// followed by "#define directives"
	var prefix = ""
	if defines != nil {
		prefix = fmt.Sprintf("#version %s\n", defines["GLSL_VERSION"])
		for name, value := range defines {
			if name == "GLSL_VERSION" {
				continue
			}
			prefix = prefix + fmt.Sprintf("#define %s %s\n", name, value)
		}
	}

	// Find all string submatches for the "#include <name>" directive
	matches := rexInclude.FindAllStringSubmatch(source, 100)
	if len(matches) == 0 {
		return prefix + source, nil
	}

	// For each directive found, replace the name by the respective include chunk source code
	var newSource = source
	for _, m := range matches {
		// Get the source of the include chunk with the match <name>
		incSource := sm.includes[m[1]]
		if len(incSource) == 0 {
			return "", fmt.Errorf("Include:[%s] not found", m[1])
		}
		// Preprocess the include chunk source code
		incSource, err := sm.preprocess(incSource, nil)
		if err != nil {
			return "", err
		}
		// Replace all occurances of the include directive with its processed source code
		newSource = strings.Replace(newSource, m[0], incSource, -1)
	}
	return prefix + newSource, nil
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
