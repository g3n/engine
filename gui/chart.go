package gui

import (
	"fmt"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer/shader"
)

func init() {
	shader.AddShader("shaderChartVertex", shaderChartVertex)
	shader.AddShader("shaderChartFrag", shaderChartFrag)
	shader.AddProgram("shaderChart", "shaderChartVertex", "shaderChartFrag")
}

// ChartLine implements a panel which can contain several line charts
type ChartLine struct {
	Panel                // Embedded panel
	grid    *ChartGrid   // Optional chart grid
	baseX   float32      // NDC x coordinate for y axis
	baseY   float32      // NDC y coordinate for x axis
	scaleX  *ChartScaleX // X scale objet
	firstX  float32      // First value for X label
	stepX   float32      // Step value for X label
	formatX string       // Format for X labels
	scaleY  *ChartScaleY // Y scale objet
	graphs  []*LineGraph // Array of line graphs
}

// NewChartLine creates and returns a new line chart panel with
// the specified dimensions in pixels.
func NewChartLine(width, height float32) *ChartLine {

	cl := new(ChartLine)
	cl.Panel.Initialize(width, height)
	cl.baseX = 0.1
	cl.baseY = -0.9
	cl.firstX = 0.0
	cl.stepX = 1.0
	cl.formatX = "%2.1f"
	return cl
}

// SetGrid sets the line chart grid with the specified number of
// grid lines and color
func (cl *ChartLine) SetGrid(xcount, ycount int, color *math32.Color) {

	if cl.grid != nil {
		cl.Node.Remove(cl.grid)
		cl.grid.Dispose()
	}
	cl.grid = NewChartGrid(&cl.Panel, xcount, ycount, color)
	cl.Node.Add(cl.grid)
}

// SetScaleX sets the line chart x scale number of lines and color
func (cl *ChartLine) SetScaleX(lines int, color *math32.Color) {

	if cl.scaleX != nil {
		cl.Node.Remove(cl.scaleX)
		cl.scaleX.Dispose()
	}
	cl.scaleX = newChartScaleX(cl, lines, color)
	cl.Node.Add(cl.scaleX)
}

// SetScaleY sets the line chart y scale number of lines and color
func (cl *ChartLine) SetScaleY(lines int, color *math32.Color) {

	if cl.scaleY != nil {
		cl.Node.Remove(cl.scaleY)
		cl.scaleY.Dispose()
	}
	cl.scaleY = newChartScaleY(cl, lines, color)
	cl.Node.Add(cl.scaleY)
}

// AddLine adds a line graph to the chart
func (cl *ChartLine) AddGraph(name, title string, color *math32.Color, data []float32) {

	graph := newLineGraph(&cl.Panel, name, title, color, data)
	cl.graphs = append(cl.graphs, graph)
	cl.Node.Add(graph)
}

//
//
//
// ChartScaleX
//
//
//
type ChartScaleX struct {
	graphic.Graphic                     // It is a Graphic
	chart           *ChartLine          // Container chart
	modelMatrixUni  gls.UniformMatrix4f // Model matrix uniform
}

func newChartScaleX(chart *ChartLine, lines int, color *math32.Color) *ChartScaleX {

	sx := new(ChartScaleX)
	sx.chart = chart

	// Generates grid lines using Normalized Device Coordinates and
	// considering that the parent panel model coordinates are:
	// 0,0,0           1,0,0
	// +---------------+
	// |               |
	// |               |
	// +---------------+
	// 0,-1,0          1,-1,0
	positions := math32.NewArrayF32(0, 0)

	// Appends scaleX bottom horizontal base line
	positions.Append(
		0, chart.baseY, 0, color.R, color.G, color.B, // line start vertex and color
		1, chart.baseY, 0, color.R, color.G, color.B, // line end vertex and color
	)
	// Appends scale X vertical lines
	step := 1 / (float32(lines) + 1)
	for i := 1; i < lines+1; i++ {
		nx := float32(i) * step
		positions.Append(
			nx, 0, 0, color.R, color.G, color.B, // line start vertex and color
			nx, chart.baseY, 0, color.R, color.G, color.B, // line end vertex and color
		)
		l := NewLabel(fmt.Sprintf(sx.chart.formatX, float32(i)))
		px, py := ndc2pix(&sx.chart.Panel, nx, chart.baseY)
		l.SetPosition(px, py)
		sx.chart.Add(l)
	}

	// Creates geometry using one interlaced VBO
	geom := geometry.NewGeometry()
	geom.AddVBO(gls.NewVBO().
		AddAttrib("VertexPosition", 3).
		AddAttrib("VertexColor", 3).
		SetBuffer(positions),
	)

	// Creates material
	mat := material.NewMaterial()
	mat.SetLineWidth(1.0)
	mat.SetShader("shaderChart")

	// Initializes the grid graphic
	sx.Graphic.Init(geom, gls.LINES)
	sx.AddMaterial(sx, mat, 0, 0)
	sx.modelMatrixUni.Init("ModelMatrix")

	return sx
}

// Converts panel ndc coordinates to relative pixels inside panel
func ndc2pix(p *Panel, nx, ny float32) (px, py float32) {

	w := p.ContentWidth()
	h := p.ContentHeight()
	return w * nx, -h * ny
}

// RenderSetup is called by the renderer before drawing this graphic
// Calculates the model matrix and transfer to OpenGL.
func (sx *ChartScaleX) RenderSetup(gs *gls.GLS, rinfo *core.RenderInfo) {

	// Set this model matrix the same as the chart panel
	var mm math32.Matrix4
	sx.chart.SetModelMatrix(gs, &mm)

	// Sets and transfer the model matrix uniform
	sx.modelMatrixUni.SetMatrix4(&mm)
	sx.modelMatrixUni.Transfer(gs)
}

//
//
//
// ChartScaleY
//
//
//
type ChartScaleY struct {
	graphic.Graphic                     // It is a Graphic
	chart           *ChartLine          // Container chart
	modelMatrixUni  gls.UniformMatrix4f // Model matrix uniform
}

func newChartScaleY(chart *ChartLine, lines int, color *math32.Color) *ChartScaleY {

	sy := new(ChartScaleY)
	sy.chart = chart

	// Generates grid lines using Normalized Device Coordinates and
	// considering that the parent panel model coordinates are:
	// 0,0,0           1,0,0
	// +---------------+
	// |               |
	// |               |
	// +---------------+
	// 0,-1,0          1,-1,0
	positions := math32.NewArrayF32(0, 0)

	// Appends scaleY left vertical axis line
	positions.Append(
		chart.baseX, 0, 0, color.R, color.G, color.B, // line start vertex and color
		chart.baseX, -1, 0, color.R, color.G, color.B, // line end vertex and color
	)
	// Appends scale horizontal lines starting from baseX
	step := 1 / (float32(lines) + 1)
	for i := 1; i < lines+1; i++ {
		ny := -float32(i) * step
		positions.Append(
			chart.baseX, ny, 0, color.R, color.G, color.B, // line start vertex and color
			1, ny, 0, color.R, color.G, color.B, // line end vertex and color
		)
		//		l := NewLabel(fmt.Sprintf(sx.chart.formatX, float32(i)))
		//		px, py := ndc2pix(&sx.chart.Panel, nx, baseY)
		//		l.SetPosition(px, py)
		//		sx.chart.Add(l)
	}

	// Creates geometry using one interlaced VBO
	geom := geometry.NewGeometry()
	geom.AddVBO(gls.NewVBO().
		AddAttrib("VertexPosition", 3).
		AddAttrib("VertexColor", 3).
		SetBuffer(positions),
	)

	// Creates material
	mat := material.NewMaterial()
	mat.SetLineWidth(1.0)
	mat.SetShader("shaderChart")

	// Initializes the grid graphic
	sy.Graphic.Init(geom, gls.LINES)
	sy.AddMaterial(sy, mat, 0, 0)
	sy.modelMatrixUni.Init("ModelMatrix")

	return sy
}

// RenderSetup is called by the renderer before drawing this graphic
// Calculates the model matrix and transfer to OpenGL.
func (sy *ChartScaleY) RenderSetup(gs *gls.GLS, rinfo *core.RenderInfo) {

	// Set this model matrix the same as the chart panel
	var mm math32.Matrix4
	sy.chart.SetModelMatrix(gs, &mm)

	// Sets and transfer the model matrix uniform
	sy.modelMatrixUni.SetMatrix4(&mm)
	sy.modelMatrixUni.Transfer(gs)
}

//
//
//
// ChartGrid
//
//
//

// ChartGrid implements a 2D grid used inside charts
type ChartGrid struct {
	graphic.Graphic                     // It is a Graphic
	chart           *Panel              // Container chart
	modelMatrixUni  gls.UniformMatrix4f // Model matrix uniform
}

// NewChartGrid creates and returns a pointer to a chart grid graphic for the
// specified parent chart and with the specified number of grid lines and color.
func NewChartGrid(chart *Panel, xcount, ycount int, color *math32.Color) *ChartGrid {

	cg := new(ChartGrid)
	cg.chart = chart

	// Generates grid lines using Normalized Device Coordinates and
	// considering that the parent panel model coordinates are:
	// 0,0,0           1,0,0
	// +---------------+
	// |               |
	// |               |
	// +---------------+
	// 0,-1,0          1,-1,0
	positions := math32.NewArrayF32(0, 0)
	xstep := 1 / (float32(xcount) + 1)
	for xi := 1; xi < xcount+1; xi++ {
		posx := float32(xi) * xstep
		positions.Append(
			posx, 0, 0, color.R, color.G, color.B, // line start vertex and color
			posx, -1, 0, color.R, color.G, color.B, // line end vertex and color
		)
	}
	ystep := 1 / (float32(ycount) + 1)
	for yi := 1; yi < ycount+1; yi++ {
		posy := -float32(yi) * ystep
		positions.Append(
			0, posy, 0, color.R, color.G, color.B, // line start vertex and color
			1, posy, 0, color.R, color.G, color.B, // line end vertex and color
		)
	}

	// Creates geometry using one interlaced VBO
	geom := geometry.NewGeometry()
	geom.AddVBO(
		gls.NewVBO().
			AddAttrib("VertexPosition", 3).
			AddAttrib("VertexColor", 3).
			SetBuffer(positions),
	)

	// Creates material
	mat := material.NewMaterial()
	mat.SetLineWidth(1.0)
	mat.SetShader("shaderChart")

	// Initializes the grid graphic
	cg.Graphic.Init(geom, gls.LINES)
	cg.AddMaterial(cg, mat, 0, 0)
	cg.modelMatrixUni.Init("ModelMatrix")
	return cg
}

// RenderSetup is called by the renderer before drawing this graphic
// Calculates the model matrix and transfer to OpenGL.
func (cg *ChartGrid) RenderSetup(gs *gls.GLS, rinfo *core.RenderInfo) {

	// Set this model matrix the same as the chart panel
	var mm math32.Matrix4
	cg.chart.SetModelMatrix(gs, &mm)

	// Sets and transfer the model matrix uniform
	cg.modelMatrixUni.SetMatrix4(&mm)
	cg.modelMatrixUni.Transfer(gs)
}

//
//
//
// LineGraph
//
//
//

// LineGraph implemens a 2D line graph
type LineGraph struct {
	graphic.Graphic                     // It is a Graphic
	chart           *Panel              // Container chart
	modelMatrixUni  gls.UniformMatrix4f // Model matrix uniform
	name            string              // Name id
	title           string              // Title string
	color           *math32.Color       // Line color
	data            []float32           // Data
}

// newLineGraph creates and returns a pointer to a line graph graphic for the
// specified parent chart
func newLineGraph(chart *Panel, name, title string, color *math32.Color, data []float32) *LineGraph {

	lg := new(LineGraph)
	lg.chart = chart
	lg.name = name
	lg.title = title
	lg.color = color
	lg.data = data

	// Generates graph lines using Normalized Device Coordinates and
	// considering that the parent panel model coordinates are:
	// 0,0,0           1,0,0
	// +---------------+
	// |               |
	// |               |
	// +---------------+
	// 0,-1,0          1,-1,0
	positions := math32.NewArrayF32(0, 0)
	for i := 0; i < len(data); i++ {
		px := float32(i) / float32(len(data))
		py := -1 + data[i]
		positions.Append(px, py, 0, color.R, color.G, color.B)
	}

	// Creates geometry using one interlaced VBO
	geom := geometry.NewGeometry()
	geom.AddVBO(gls.NewVBO().
		AddAttrib("VertexPosition", 3).
		AddAttrib("VertexColor", 3).
		SetBuffer(positions),
	)

	// Creates material
	mat := material.NewMaterial()
	mat.SetLineWidth(2.5)
	mat.SetShader("shaderChart")

	// Initializes the graphic
	lg.Graphic.Init(geom, gls.LINE_STRIP)
	lg.AddMaterial(lg, mat, 0, 0)
	lg.modelMatrixUni.Init("ModelMatrix")

	return lg
}

// RenderSetup is called by the renderer before drawing this graphic
// Calculates the model matrix and transfer to OpenGL.
func (lg *LineGraph) RenderSetup(gs *gls.GLS, rinfo *core.RenderInfo) {

	// Set this graphic model matrix the same as the container chart panel
	var mm math32.Matrix4
	lg.chart.SetModelMatrix(gs, &mm)

	// Sets and transfer the model matrix uniform
	lg.modelMatrixUni.SetMatrix4(&mm)
	lg.modelMatrixUni.Transfer(gs)
}

//
// Vertex Shader template
//
const shaderChartVertex = `
#version {{.Version}}

// Vertex attributes
{{template "attributes" .}}

// Input uniforms
uniform mat4 ModelMatrix;

// Outputs for fragment shader
out vec3 Color;

void main() {

    Color = VertexColor;

    // Set position
    vec4 pos = vec4(VertexPosition.xyz, 1);
    gl_Position = ModelMatrix * pos;
}
`

//
// Fragment Shader template
//
const shaderChartFrag = `
#version {{.Version}}

in vec3 Color;
out vec4 FragColor;

void main() {

    FragColor = vec4(Color, 1.0);
}
`
