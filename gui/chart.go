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
	"math"
)

func init() {
	shader.AddShader("shaderChartVertex", shaderChartVertex)
	shader.AddShader("shaderChartFrag", shaderChartFrag)
	shader.AddProgram("shaderChart", "shaderChartVertex", "shaderChartFrag")
}

// ChartLine implements a panel which can contain several line charts
type ChartLine struct {
	Panel                // Embedded panel
	title   *Label       // Optional title label
	baseX   float32      // NDC x coordinate for y axis
	baseY   float32      // NDC y coordinate for x axis
	scaleX  *ChartScaleX // X scale panel
	scaleY  *ChartScaleY // Y scale panel
	offsetX int          // Initial offset in data buffers
	countX  int          // Count of data buffer points starting from offsetX
	firstX  float32      // Value of first data point to show
	stepX   float32      // Step to add to firstX for next data point
	minY    float32      // Minimum Y value
	maxY    float32      //
	autoY   bool         // Auto range flag for Y values
	graphs  []*LineGraph // Array of line graphs
}

// NewChartLine creates and returns a new line chart panel with
// the specified dimensions in pixels.
func NewChartLine(width, height float32) *ChartLine {

	cl := new(ChartLine)
	cl.Panel.Initialize(width, height)
	cl.baseX = 0.1
	cl.baseY = -0.9
	cl.offsetX = 0
	cl.countX = 10
	cl.firstX = 0.0
	cl.stepX = 1.0
	return cl
}

func (cl *ChartLine) SetTitle(title *Label) {

	if cl.title != nil {
		cl.Remove(cl.title)
		cl.title = nil
	}
	if title != nil {
		cl.Add(title)
		cl.title = title
	}
	cl.recalc()
}

// SetScaleX sets the line chart x scale number of lines and color
func (cl *ChartLine) SetScaleX(lines int, color *math32.Color) {

	if cl.scaleX != nil {
		cl.Remove(cl.scaleX)
		cl.scaleX.Dispose()
	}
	cl.scaleX = newChartScaleX(cl, lines, color)
	cl.Add(cl.scaleX)
}

// ClearScaleX removes the chart x scale if it was previously set
func (cl *ChartLine) ClearScaleX() {

	if cl.scaleX == nil {
		return
	}
	cl.Remove(cl.scaleX)
	cl.scaleX.Dispose()
}

// SetScaleY sets the line chart y scale number of lines and color
func (cl *ChartLine) SetScaleY(lines int, color *math32.Color) {

	if cl.scaleY != nil {
		cl.Remove(cl.scaleY)
		cl.scaleY.Dispose()
	}
	cl.scaleY = newChartScaleY(cl, lines, color)
	cl.Add(cl.scaleY)
}

// ClearScaleY removes the chart x scale if it was previously set
func (cl *ChartLine) ClearScaleY() {

	if cl.scaleY == nil {
		return
	}
	cl.Remove(cl.scaleY)
	cl.scaleY.Dispose()
}

// SetRangeX sets the interval of the data to be shown:
// offset is the start position in the Y data array.
// count is the number of data points to show, starting from the specified offset.
// first is the label for the first X point
// step is the value added to first for the next data point
func (cl *ChartLine) SetRangeX(offset int, count int, first float32, step float32) {

	cl.offsetX = offset
	cl.countX = count
	cl.firstX = first
	cl.stepX = step
}

func (cl *ChartLine) SetRangeY(min float32, max float32) {

}

func (cl *ChartLine) calcRangeY() {

	if cl.autoY {
		return
	}

	minY := float32(math.MaxFloat32)
	maxY := -float32(math.MaxFloat32)
	for g := 0; g < len(cl.graphs); g++ {
		graph := cl.graphs[g]
		for x := 0; x < cl.countX; x++ {
			if x+cl.offsetX >= len(graph.y) {
				break
			}
			vy := graph.y[x+cl.offsetX]
			if vy < minY {
				minY = vy
			}
			if vy > maxY {
				maxY = vy
			}
		}
	}
	cl.minY = minY
	cl.maxY = maxY
}

// recalc recalculates the positions of the inner panels
func (cl *ChartLine) recalc() {

	if cl.title != nil {
		xpos := (cl.ContentWidth() - cl.title.width) / 2
		cl.title.SetPositionX(xpos)
	}
	if cl.scaleX != nil {
		cl.scaleX.recalc()
	}
	if cl.scaleY != nil {
		cl.scaleY.recalc()
	}
}

//// AddLine adds a line graph to the chart
//func (cl *ChartLine) AddGraph(name, title string, color *math32.Color, data []float32) {
//
//	graph := newLineGraph(&cl.Panel, name, title, color, data)
//	cl.graphs = append(cl.graphs, graph)
//	cl.Node.Add(graph)
//}

// ChartScaleX is a panel with GL_LINES geometry which draws the chart X horizontal scale axis,
// vertical lines and line labels.
type ChartScaleX struct {
	Panel             // Embedded panel
	chart  *ChartLine // Container chart
	lines  int        // Number of vertical lines
	format string     // Labels format string
	model  *Label     // Model label to calculate height/width
	labels []*Label   // Array of scale labels
}

// newChartScaleX creates and returns a pointer to a new ChartScaleX for the specified
// chart, number of lines and color
func newChartScaleX(chart *ChartLine, lines int, color *math32.Color) *ChartScaleX {

	sx := new(ChartScaleX)
	sx.chart = chart
	sx.lines = lines
	sx.format = "%v"
	sx.model = NewLabel(" ")

	// Generates scale lines using Normalized Device Coordinates and
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
		chart.baseX, chart.baseY, 0, color.R, color.G, color.B, // line start vertex and color
		1, chart.baseY, 0, color.R, color.G, color.B, // line end vertex and color
	)
	// Appends scale X vertical lines
	startidx := 0
	if chart.scaleY != nil {
		startidx++
	}
	step := 1 / (float32(lines) + 1)
	for i := startidx; i < lines; i++ {
		nx := chart.baseX + float32(i)*step
		positions.Append(
			nx, 0, 0, color.R, color.G, color.B, // line start vertex and color
			nx, chart.baseY, 0, color.R, color.G, color.B, // line end vertex and color
		)
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

	// Initializes the panel with this graphic
	gr := graphic.NewGraphic(geom, gls.LINES)
	gr.AddMaterial(sx, mat, 0, 0)
	sx.Panel.InitializeGraphic(chart.ContentWidth(), chart.ContentHeight(), gr)

	// Add labels after the panel is initialized
	for i := 0; i < lines; i++ {
		nx := chart.baseX + float32(i)*step
		l := NewLabel(fmt.Sprintf(sx.format, float32(i)))
		px, py := ndc2pix(&sx.Panel, nx, chart.baseY)
		//log.Error("label x:%v y:%v", px, py)
		l.SetPosition(px, py)
		sx.Add(l)
		sx.labels = append(sx.labels, l)
	}
	sx.recalc()

	return sx
}

func (sx *ChartScaleX) updateLabels() {

	step := 1 / (float32(sx.lines) + 1)
	for i := 0; i < len(sx.labels); i++ {
		nx := sx.chart.baseX + float32(i)*step
		px, py := ndc2pix(&sx.Panel, nx, sx.chart.baseY)
		log.Error("label x:%v y:%v", px, py)
		l := sx.labels[i]
		l.SetPosition(px, py)
	}
}

func (sx *ChartScaleX) setLabelsText(x []float32) {

}

func (sx *ChartScaleX) recalc() {

	if sx.chart.title != nil {
		th := sx.chart.title.Height()
		sx.SetPosition(0, th)
		sx.SetHeight(sx.chart.ContentHeight() - th)
		sx.updateLabels()
	} else {
		sx.SetPosition(0, 0)
		sx.SetHeight(sx.chart.ContentHeight())
		sx.updateLabels()
	}
}

// RenderSetup is called by the renderer before drawing this graphic
// Calculates the model matrix and transfer to OpenGL.
func (sx *ChartScaleX) RenderSetup(gs *gls.GLS, rinfo *core.RenderInfo) {

	//log.Error("ChartScaleX RenderSetup:%p", sx)
	// Sets model matrix and transfer to shader
	var mm math32.Matrix4
	sx.SetModelMatrix(gs, &mm)
	sx.modelMatrixUni.SetMatrix4(&mm)
	sx.modelMatrixUni.Transfer(gs)
}

// ChartScaleY is a panel with LINE geometry which draws the chart Y vertical scale axis,
// horizontal and labels.
type ChartScaleY struct {
	Panel             // Embedded panel
	chart  *ChartLine // Container chart
	lines  int        // Number of horizontal lines
	format string     // Labels format string
	model  *Label     // Model label to calculate height/width
	labels []*Label   // Array of scale labels
}

// newChartScaleY creates and returns a pointer to a new ChartScaleY for the specified
// chart, number of lines and color
func newChartScaleY(chart *ChartLine, lines int, color *math32.Color) *ChartScaleY {

	sy := new(ChartScaleY)
	sy.chart = chart
	sy.lines = lines
	sy.format = "%v"
	sy.model = NewLabel(" ")

	// Generates scale lines using Normalized Device Coordinates and
	// considering that the parent panel model coordinates are:
	// 0,0,0    1,0,0
	// +--------+
	// |        |
	// +--------+
	// 0,-1,0   1,-1,0
	positions := math32.NewArrayF32(0, 0)
	// Appends scaleY left vertical line
	positions.Append(
		chart.baseX, chart.baseY, 0, color.R, color.G, color.B, // line start vertex and color
		chart.baseX, 0, 0, color.R, color.G, color.B, // line end vertex and color
	)
	// Appends scale horizontall lines
	startidx := 0
	if chart.scaleX != nil {
		startidx++
	}
	step := 1 / (float32(lines) + 1)
	for i := startidx; i < lines; i++ {
		ny := chart.baseY + float32(i)*step
		positions.Append(
			chart.baseX, ny, 0, color.R, color.G, color.B, // line start vertex and color
			1, ny, 0, color.R, color.G, color.B, // line end vertex and color
		)
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

	// Initializes the panel with this graphic
	gr := graphic.NewGraphic(geom, gls.LINES)
	gr.AddMaterial(sy, mat, 0, 0)
	sy.Panel.InitializeGraphic(chart.ContentWidth(), chart.ContentHeight(), gr)

	// Add labels after the panel is initialized
	for i := 0; i < lines; i++ {
		ny := chart.baseY + float32(i)*step
		l := NewLabel(fmt.Sprintf(sy.format, float32(i)))
		px, py := ndc2pix(&sy.Panel, 0, ny)
		py -= sy.model.Height() / 2
		//log.Error("label x:%v y:%v", px, py)
		l.SetPosition(px, py)
		sy.Add(l)
		sy.labels = append(sy.labels, l)
	}
	sy.recalc()
	return sy
}

func (sy *ChartScaleY) updateLabels() {

	step := 1 / (float32(sy.lines) + 1)
	for i := 0; i < len(sy.labels); i++ {
		ny := sy.chart.baseY + float32(i)*step
		px, py := ndc2pix(&sy.Panel, 0, ny)
		py -= sy.model.Height() / 2
		log.Error("label x:%v y:%v", px, py)
		l := sy.labels[i]
		l.SetPosition(px, py)
	}
}

func (sy *ChartScaleY) recalc() {

	if sy.chart.title != nil {
		th := sy.chart.title.Height()
		sy.SetPosition(0, th)
		sy.SetHeight(sy.chart.ContentHeight() - th)
		sy.updateLabels()
	} else {
		sy.SetPosition(0, 0)
		sy.SetHeight(sy.chart.ContentHeight())
		sy.updateLabels()
	}
}

// Converts panel ndc coordinates to relative pixels inside panel
func ndc2pix(p *Panel, nx, ny float32) (px, py float32) {

	w := p.ContentWidth()
	h := p.ContentHeight()
	return w * nx, -h * ny
}

//
// LineGraph
//
type LineGraph struct {
	Panel              // Embedded panel
	chart *ChartLine   // Container chart
	color math32.Color // Line color
	y     []float32    // Data y
}

func newLineGraph(chart *ChartLine, color *math32.Color, y []float32) *LineGraph {

	lg := new(LineGraph)
	lg.chart = chart
	lg.color = *color
	lg.y = y
	lg.setGeometry()
	return lg
}

func (lg *LineGraph) SetColor(color *math32.Color) {

}

func (lg *LineGraph) SetData(x, y []float32) {

	lg.y = y
}

func (lg *LineGraph) setGeometry() {

	lg.chart.calcRangeY()

	// Creates array for vertices and colors
	positions := math32.NewArrayF32(0, 0)
	origin := false
	step := 1.0 / float32(lg.chart.countX-1)
	rangeY := lg.chart.maxY - lg.chart.minY
	for i := 0; i < lg.chart.countX; i++ {
		x := i + lg.chart.offsetX
		if x >= len(lg.y) {
			break
		}
		// Get Y value and checks if it is inside the range
		vy := lg.y[x]
		if vy < lg.chart.minY || vy > lg.chart.maxY {
			continue
		}
		px := float32(i) * step
		if !origin {
			positions.Append(px, -1, 0, lg.color.R, lg.color.G, lg.color.B)
			origin = true
		}
		py := vy / rangeY
		positions.Append(px, py, 0, lg.color.R, lg.color.G, lg.color.B)
	}

	// Creates geometry using one interlaced VBO
	geom := geometry.NewGeometry()
	geom.AddVBO(gls.NewVBO().
		AddAttrib("VertexPosition", 3).
		AddAttrib("VertexColor", 3).
		SetBuffer(positions),
	)

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
