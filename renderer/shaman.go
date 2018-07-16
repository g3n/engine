// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package renderer

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/renderer/shaders"
	"strconv"
)

const GLSL_VERSION = "330 core"

// Regular expression to parse #include <name> [quantity] directive
var rexInclude *regexp.Regexp
const indexParameter = "{i}"

func init() {

	rexInclude = regexp.MustCompile(`#include\s+<(.*)>\s*(?:\[(.*)]|)`)
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
	Defines          gls.ShaderDefines  // Additional shader defines
}

// ProgSpecs represents a compiled shader program along with its specs
type ProgSpecs struct {
	program *gls.Program // program object
	specs   ShaderSpecs  // associated specs
}

// Shaman is the shader manager
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

// SetProgram sets the shader program to satisfy the specified specs.
// Returns an indication if the current shader has changed and a possible error
// when creating a new shader program.
// Receives a copy of the specs because it changes the fields which specify the
// number of lights depending on the UseLights flags.
func (sm *Shaman) SetProgram(s *ShaderSpecs) (bool, error) {

	// Checks material use lights bit mask
	var specs ShaderSpecs
	specs.copy(s)
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
	if sm.specs.equals(&specs) {
		return false, nil
	}

	// Search for compiled program with the specified specs
	for _, pinfo := range sm.programs {
		if pinfo.specs.equals(&specs) {
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

// GenProgram generates shader program from the specified specs
func (sm *Shaman) GenProgram(specs *ShaderSpecs) (*gls.Program, error) {

	// Get info for the specified shader program
	progInfo, ok := sm.proginfo[specs.Name]
	if !ok {
		return nil, fmt.Errorf("Program:%s not found", specs.Name)
	}

	// Sets the defines map
	defines := map[string]string{}
	defines["AMB_LIGHTS"] = strconv.Itoa(specs.AmbientLightsMax)
	defines["DIR_LIGHTS"] = strconv.Itoa(specs.DirLightsMax)
	defines["POINT_LIGHTS"] = strconv.Itoa(specs.PointLightsMax)
	defines["SPOT_LIGHTS"] = strconv.Itoa(specs.SpotLightsMax)
	defines["MAT_TEXTURES"] = strconv.Itoa(specs.MatTexturesMax)

	// Adds additional material and geometry defines from the specs parameter
	for name, value := range specs.Defines {
		defines[name] = value
	}

	// Get vertex shader source
	vertexSource, ok := sm.shadersm[progInfo.Vertex]
	if !ok {
		return nil, fmt.Errorf("Vertex shader:%s not found", progInfo.Vertex)
	}
	// Pre-process vertex shader source
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
	// Pre-process fragment shader source
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
		// Pre-process geometry shader source
		geomSource, err = sm.preprocess(geomSource, defines)
		if err != nil {
			return nil, err
		}
	}

	// Creates shader program
	prog := sm.gs.NewProgram()
	prog.AddShader(gls.VERTEX_SHADER, vertexSource)
	prog.AddShader(gls.FRAGMENT_SHADER, fragSource)
	if progInfo.Geometry != "" {
		prog.AddShader(gls.GEOMETRY_SHADER, geomSource)
	}
	err = prog.Build()
	if err != nil {
		return nil, err
	}

	return prog, nil
}


func (sm *Shaman) preprocess(source string, defines map[string]string) (string, error) {

	// If defines map supplied, generate prefix with glsl version directive first,
	// followed by "#define" directives
	var prefix = ""
	if defines != nil { // This is only true for the outer call
		prefix = fmt.Sprintf("#version %s\n", GLSL_VERSION)
		for name, value := range defines {
			prefix = prefix + fmt.Sprintf("#define %s %s\n", name, value)
		}
	}

	return sm.processIncludes(prefix + source, defines)
}


// preprocess preprocesses the specified source prefixing it with optional defines directives
// contained in "defines" parameter and replaces '#include <name>' directives
// by the respective source code of include chunk of the specified name.
// The included "files" are also processed recursively.
func (sm *Shaman) processIncludes(source string, defines map[string]string) (string, error) {

	// Find all string submatches for the "#include <name>" directive
	matches := rexInclude.FindAllStringSubmatch(source, 100)
	if len(matches) == 0 {
		return source, nil
	}

	// For each directive found, replace the name by the respective include chunk source code
	//var newSource = source
	for _, m := range matches {
		incFullMatch := m[0]
		incName := m[1]
		incQuantityVariable := m[2]

		// Get the source of the include chunk with the match <name>
		incSource := sm.includes[incName]
		if len(incSource) == 0 {
			return "", fmt.Errorf("Include:[%s] not found", incName)
		}

		// Preprocess the include chunk source code
		incSource, err := sm.processIncludes(incSource, defines)
		if err != nil {
			return "", err
		}

		// Skip line
		incSource = "\n" + incSource

		// Process include quantity variable if provided
		if incQuantityVariable != "" {
			incQuantityString, defined := defines[incQuantityVariable]
			if defined { // Only process #include if quantity variable is defined
				incQuantity, err := strconv.Atoi(incQuantityString)
				if err != nil {
					return "", err
				}
				// Check for iterated includes and populate index parameter
				if incQuantity > 0 {
					repeatedIncludeSource := ""
					for i := 0; i < incQuantity; i++ {
						// Replace all occurrences of the index parameter with the current index i.
						repeatedIncludeSource += strings.Replace(incSource, indexParameter, strconv.Itoa(i), -1)
					}
					incSource = repeatedIncludeSource
				}
			} else {
				incSource = ""
			}
		}

		// Replace all occurrences of the include directive with its processed source code
		source = strings.Replace(source, incFullMatch, incSource, -1)
	}
	return source, nil
}

// copy copies other spec into this
func (ss *ShaderSpecs) copy(other *ShaderSpecs) {

	*ss = *other
	if other.Defines != nil {
		ss.Defines = *gls.NewShaderDefines()
		ss.Defines.Add(&other.Defines)
	}
}

// equals compares two ShaderSpecs and returns true if they are effectively equal.
func (ss *ShaderSpecs) equals(other *ShaderSpecs) bool {

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
		ss.MatTexturesMax == other.MatTexturesMax &&
		ss.Defines.Equals(&other.Defines) {
		return true
	}
	return false
}
