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

//
//
// ChartLine implements a panel which can contain several line charts
//
//
type ChartLine struct {
	Panel                   // Embedded panel
	left       float32      // Left margin in pixels
	bottom     float32      // Bottom margin in pixels
	top        float32      // Top margin in pixels
	firstX     float32      // Value for the first x label
	stepX      float32      // Step for the next x label
	countStepX float32      // Number of values per x step
	minY       float32      // Minimum Y value
	maxY       float32      // Maximum Y value
	autoY      bool         // Auto range flag for Y values
	formatX    string       // String format for scale X labels
	formatY    string       // String format for scale Y labels
	title      *Label       // Optional title label
	scaleX     *ChartScaleX // X scale panel
	scaleY     *ChartScaleY // Y scale panel
	labelsX    []*Label     // Array of scale X labels
	labelsY    []*Label     // Array of scale Y labels
	graphs     []*LineGraph // Array of line graphs
}

const (
	deltaLine = 0.001 // Delta in NDC for lines over the boundary
)

// NewChartLine creates and returns a new line chart panel with
// the specified dimensions in pixels.
func NewChartLine(width, height float32) *ChartLine {

	cl := new(ChartLine)
	cl.Panel.Initialize(width, height)
	cl.left = 34
	cl.bottom = 20
	cl.top = 10
	cl.firstX = 0
	cl.stepX = 1
	cl.countStepX = 0
	cl.minY = -10.0
	cl.maxY = 10.0
	cl.autoY = false
	cl.formatX = "%v"
	cl.formatY = "%v"
	return cl
}

//func (cl *ChartLine) SetMargins(left, bottom float32) {
//
//	cl.baseX, cl.baseY = cl.Pix2NDC(left, bottom)
//	cl.recalc()
//}

// SetTitle sets the chart title
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

// SetFormatX sets the string format of the X scale labels
func (cl *ChartLine) SetFormatX(format string) {

	cl.formatX = format
	cl.updateLabelsX()
}

// SetFormatY sets the string format of the Y scale labels
func (cl *ChartLine) SetFormatY(format string) {

	cl.formatY = format
	cl.updateLabelsY()
}

// SetScaleX sets the X scale number of lines and color
func (cl *ChartLine) SetScaleX(lines int, color *math32.Color) {

	if cl.scaleX != nil {
		cl.ClearScaleX()
	}

	// Add scale lines
	cl.scaleX = newChartScaleX(cl, lines, color)
	cl.Add(cl.scaleX)

	// Add scale labels
	// The positions of the labels will be set by 'recalc()'
	value := cl.firstX
	for i := 0; i < lines; i++ {
		l := NewLabel(fmt.Sprintf(cl.formatX, value))
		cl.Add(l)
		cl.labelsX = append(cl.labelsX, l)
		value += cl.stepX
	}
	cl.recalc()
}

// ClearScaleX removes the X scale if it was previously set
func (cl *ChartLine) ClearScaleX() {

	if cl.scaleX == nil {
		return
	}

	// Remove and dispose scale lines
	cl.Remove(cl.scaleX)
	cl.scaleX.Dispose()

	// Remove and dispose scale labels
	for i := 0; i < len(cl.labelsX); i++ {
		label := cl.labelsX[i]
		cl.Remove(label)
		label.Dispose()
	}
	cl.labelsX = cl.labelsX[0:0]
	cl.scaleX = nil
}

// SetScaleY sets the Y scale number of lines and color
func (cl *ChartLine) SetScaleY(lines int, color *math32.Color) {

	if cl.scaleY != nil {
		cl.ClearScaleY()
	}

	if lines < 2 {
		lines = 2
	}

	// Add scale lines
	cl.scaleY = newChartScaleY(cl, lines, color)
	cl.Add(cl.scaleY)

	// Add scale labels
	// The position of the labels will be set by 'recalc()'
	value := cl.minY
	step := (cl.maxY - cl.minY) / float32(lines-1)
	for i := 0; i < lines; i++ {
		l := NewLabel(fmt.Sprintf(cl.formatY, value))
		cl.Add(l)
		cl.labelsY = append(cl.labelsY, l)
		value += step
	}
	cl.recalc()
}

// ClearScaleY removes the Y scale if it was previously set
func (cl *ChartLine) ClearScaleY() {

	if cl.scaleY == nil {
		return
	}

	// Remove and dispose scale lines
	cl.Remove(cl.scaleY)
	cl.scaleY.Dispose()

	// Remove and dispose scale labels
	for i := 0; i < len(cl.labelsY); i++ {
		label := cl.labelsY[i]
		cl.Remove(label)
		label.Dispose()
	}
	cl.labelsY = cl.labelsY[0:0]
	cl.scaleY = nil
}

// SetRangeX sets the X scale labels and range per step
// firstX is the value of first label of the x scale
// stepX is the step to be added to get the next x scale label
// countStepX is the number of elements of the data buffer for each line step
func (cl *ChartLine) SetRangeX(firstX float32, stepX float32, countStepX float32) {

	cl.firstX = firstX
	cl.stepX = stepX
	cl.countStepX = countStepX
	cl.updateGraphs()
}

// SetRangeY sets the minimum and maximum values of the y scale
func (cl *ChartLine) SetRangeY(min float32, max float32) {

	if cl.autoY {
		return
	}
	cl.minY = min
	cl.maxY = max
	cl.updateGraphs()
}

// SetRangeYauto sets the state of the auto
func (cl *ChartLine) SetRangeYauto(auto bool) {

	cl.autoY = auto
	if !auto {
		return
	}
	cl.updateGraphs()
}

// Returns the current y range
func (cl *ChartLine) RangeY() (minY, maxY float32) {

	return cl.minY, cl.maxY
}

// AddLine adds a line graph to the chart
func (cl *ChartLine) AddGraph(color *math32.Color, data []float32) *LineGraph {

	graph := newLineGraph(cl, color, data)
	cl.graphs = append(cl.graphs, graph)
	cl.Add(graph)
	cl.recalc()
	cl.updateGraphs()
	return graph
}

// RemoveGraph removes and disposes of the specified graph from the chart
func (cl *ChartLine) RemoveGraph(g *LineGraph) {

	cl.Remove(g)
	g.Dispose()
	for pos, current := range cl.graphs {
		if current == g {
			copy(cl.graphs[pos:], cl.graphs[pos+1:])
			cl.graphs[len(cl.graphs)-1] = nil
			cl.graphs = cl.graphs[:len(cl.graphs)-1]
			break
		}
	}
	if !cl.autoY {
		return
	}
	cl.updateGraphs()
}

// updateLabelsX updates the X scale labels text
func (cl *ChartLine) updateLabelsX() {

	if cl.scaleX == nil {
		return
	}
	pstep := (cl.ContentWidth() - cl.left) / float32(len(cl.labelsX))
	value := cl.firstX
	for i := 0; i < len(cl.labelsX); i++ {
		label := cl.labelsX[i]
		label.SetText(fmt.Sprintf(cl.formatX, value))
		px := cl.left + float32(i)*pstep
		label.SetPosition(px, cl.ContentHeight()-cl.bottom)
		value += cl.stepX
	}
}

// updateLabelsY updates the Y scale labels text and positions
func (cl *ChartLine) updateLabelsY() {

	if cl.scaleY == nil {
		return
	}

	th := float32(0)
	if cl.title != nil {
		th = cl.title.height
	}

	nlines := cl.scaleY.lines
	vstep := (cl.maxY - cl.minY) / float32(nlines-1)
	pstep := (cl.ContentHeight() - th - cl.top - cl.bottom) / float32(nlines-1)
	value := cl.minY
	for i := 0; i < nlines; i++ {
		label := cl.labelsY[i]
		label.SetText(fmt.Sprintf(cl.formatY, value))
		px := cl.left - 2 - label.Width()
		if px < 0 {
			px = 0
		}
		py := cl.ContentHeight() - cl.bottom - float32(i)*pstep
		label.SetPosition(px, py-label.Height()/2)
		value += vstep
	}
}

// calcRangeY calculates the minimum and maximum y values for all graphs
func (cl *ChartLine) calcRangeY() {

	if !cl.autoY || len(cl.graphs) == 0 {
		return
	}
	minY := float32(math.MaxFloat32)
	maxY := -float32(math.MaxFloat32)
	for g := 0; g < len(cl.graphs); g++ {
		graph := cl.graphs[g]
		for x := 0; x < len(graph.data); x++ {
			vy := graph.data[x]
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

// updateGraphs should be called when the range the scales change or
// any graph data changes
func (cl *ChartLine) updateGraphs() {

	cl.calcRangeY()
	cl.updateLabelsX()
	cl.updateLabelsY()
	for i := 0; i < len(cl.graphs); i++ {
		g := cl.graphs[i]
		g.updateData()
	}
}

// recalc recalculates the positions of the inner panels
func (cl *ChartLine) recalc() {

	// Center title position
	if cl.title != nil {
		xpos := (cl.ContentWidth() - cl.title.width) / 2
		cl.title.SetPositionX(xpos)
	}

	// Recalc scale X and its labels
	if cl.scaleX != nil {
		cl.scaleX.recalc()
		cl.updateLabelsX()
	}

	// Recalc scale Y and its labels
	if cl.scaleY != nil {
		cl.scaleY.recalc()
		cl.updateLabelsY()
	}

	// Recalc graphs
	for i := 0; i < len(cl.graphs); i++ {
		g := cl.graphs[i]
		g.recalc()
		cl.SetTopChild(g)
	}
}

//
//
// ChartScaleX is a panel with GL_LINES geometry which draws the chart X horizontal scale axis,
// vertical lines and line labels.
//
//
type ChartScaleX struct {
	Panel                // Embedded panel
	chart  *ChartLine    // Container chart
	lines  int           // Number of vertical lines
	bounds gls.Uniform4f // Bound uniform in OpenGL window coordinates
	mat    chartMaterial // Chart material
}

// newChartScaleX creates and returns a pointer to a new ChartScaleX for the specified
// chart, number of lines and color
func newChartScaleX(chart *ChartLine, lines int, color *math32.Color) *ChartScaleX {

	sx := new(ChartScaleX)
	sx.chart = chart
	sx.lines = lines
	sx.bounds.Init("Bounds")

	// Appends bottom horizontal line
	positions := math32.NewArrayF32(0, 0)
	positions.Append(0, -1+deltaLine, 0, 1, -1+deltaLine, 0)

	// Appends vertical lines
	step := 1 / float32(lines)
	for i := 0; i < lines; i++ {
		nx := float32(i) * step
		if i == 0 {
			nx += deltaLine
		}
		positions.Append(nx, 0, 0, nx, -1, 0)
	}

	// Creates geometry and adds VBO
	geom := geometry.NewGeometry()
	geom.AddVBO(gls.NewVBO().AddAttrib("VertexPosition", 3).SetBuffer(positions))

	// Initializes the panel graphic
	gr := graphic.NewGraphic(geom, gls.LINES)
	sx.mat.Init(color)
	gr.AddMaterial(sx, &sx.mat, 0, 0)
	sx.Panel.InitializeGraphic(chart.ContentWidth(), chart.ContentHeight(), gr)

	sx.recalc()
	return sx
}

// recalc recalculates the position and size of this scale inside its parent
func (sx *ChartScaleX) recalc() {

	py := sx.chart.top
	if sx.chart.title != nil {
		py += sx.chart.title.Height()
	}
	sx.SetPosition(sx.chart.left, py)
	sx.SetSize(sx.chart.ContentWidth()-sx.chart.left, sx.chart.ContentHeight()-py-sx.chart.bottom)
}

// RenderSetup is called by the renderer before drawing this graphic
// It overrides the original panel RenderSetup
// Calculates the model matrix and transfer to OpenGL.
func (sx *ChartScaleX) RenderSetup(gs *gls.GLS, rinfo *core.RenderInfo) {

	//log.Error("ChartScaleX RenderSetup:%v", sx.pospix)
	// Sets model matrix and transfer to shader
	var mm math32.Matrix4
	sx.SetModelMatrix(gs, &mm)
	sx.modelMatrixUni.SetMatrix4(&mm)
	sx.modelMatrixUni.Transfer(gs)

	// Sets bounds in OpenGL window coordinates and transfer to shader
	_, _, _, height := gs.GetViewport()
	sx.bounds.Set(sx.pospix.X, float32(height)-sx.pospix.Y, sx.width, sx.height)
	sx.bounds.Transfer(gs)
}

//
//
// ChartScaleY is a panel with LINE geometry which draws the chart Y vertical scale axis,
// horizontal and labels.
//
//
type ChartScaleY struct {
	Panel                // Embedded panel
	chart  *ChartLine    // Container chart
	lines  int           // Number of horizontal lines
	bounds gls.Uniform4f // Bound uniform in OpenGL window coordinates
	mat    chartMaterial // Chart material
}

// newChartScaleY creates and returns a pointer to a new ChartScaleY for the specified
// chart, number of lines and color
func newChartScaleY(chart *ChartLine, lines int, color *math32.Color) *ChartScaleY {

	if lines < 2 {
		lines = 2
	}
	sy := new(ChartScaleY)
	sy.chart = chart
	sy.lines = lines
	sy.bounds.Init("Bounds")

	// Appends left vertical line
	positions := math32.NewArrayF32(0, 0)
	positions.Append(0+deltaLine, 0, 0, 0+deltaLine, -1, 0)

	// Appends horizontal lines
	step := 1 / float32(lines-1)
	for i := 0; i < lines; i++ {
		ny := -1 + float32(i)*step
		if i == 0 {
			ny += deltaLine
		}
		if i == lines-1 {
			ny -= deltaLine
		}
		positions.Append(0, ny, 0, 1, ny, 0)
	}

	// Creates geometry and adds VBO
	geom := geometry.NewGeometry()
	geom.AddVBO(gls.NewVBO().AddAttrib("VertexPosition", 3).SetBuffer(positions))

	// Initializes the panel with this graphic
	gr := graphic.NewGraphic(geom, gls.LINES)
	sy.mat.Init(color)
	gr.AddMaterial(sy, &sy.mat, 0, 0)
	sy.Panel.InitializeGraphic(chart.ContentWidth(), chart.ContentHeight(), gr)

	sy.recalc()
	return sy
}

// recalc recalculates the position and size of this scale inside its parent
func (sy *ChartScaleY) recalc() {

	py := sy.chart.top
	if sy.chart.title != nil {
		py += sy.chart.title.Height()
	}
	sy.SetPosition(sy.chart.left, py)
	sy.SetSize(sy.chart.ContentWidth()-sy.chart.left, sy.chart.ContentHeight()-py-sy.chart.bottom)
}

// RenderSetup is called by the renderer before drawing this graphic
// It overrides the original panel RenderSetup
// Calculates the model matrix and transfer to OpenGL.
func (sy *ChartScaleY) RenderSetup(gs *gls.GLS, rinfo *core.RenderInfo) {

	//log.Error("ChartScaleY RenderSetup:%v", sy.pospix)
	// Sets model matrix and transfer to shader
	var mm math32.Matrix4
	sy.SetModelMatrix(gs, &mm)
	sy.modelMatrixUni.SetMatrix4(&mm)
	sy.modelMatrixUni.Transfer(gs)

	// Sets bounds in OpenGL window coordinates and transfer to shader
	_, _, _, height := gs.GetViewport()
	sy.bounds.Set(sy.pospix.X, float32(height)-sy.pospix.Y, sy.width, sy.height)
	sy.bounds.Transfer(gs)
}

//
//
// LineGraph
//
//
type LineGraph struct {
	Panel                   // Embedded panel
	chart     *ChartLine    // Container chart
	color     math32.Color  // Line color
	data      []float32     // Data y
	bounds    gls.Uniform4f // Bound uniform in OpenGL window coordinates
	mat       chartMaterial // Chart material
	vbo       *gls.VBO
	positions math32.ArrayF32
}

func newLineGraph(chart *ChartLine, color *math32.Color, y []float32) *LineGraph {

	lg := new(LineGraph)
	lg.bounds.Init("Bounds")
	lg.chart = chart
	lg.color = *color
	lg.data = y

	// Creates geometry and adds VBO with positions
	geom := geometry.NewGeometry()
	lg.vbo = gls.NewVBO().AddAttrib("VertexPosition", 3)
	lg.positions = math32.NewArrayF32(0, 0)
	lg.vbo.SetBuffer(lg.positions)
	geom.AddVBO(lg.vbo)

	// Initializes the panel with this graphic
	gr := graphic.NewGraphic(geom, gls.LINE_STRIP)
	lg.mat.Init(&lg.color)
	gr.AddMaterial(lg, &lg.mat, 0, 0)
	lg.Panel.InitializeGraphic(lg.chart.ContentWidth(), lg.chart.ContentHeight(), gr)

	lg.SetData(y)
	return lg
}

// SetColor sets the color of the graph
func (lg *LineGraph) SetColor(color *math32.Color) {

	lg.mat.color.SetColor(color)
}

// SetData sets the graph data
func (lg *LineGraph) SetData(data []float32) {

	lg.data = data
	lg.updateData()
}

// SetLineWidth sets the graph line width
func (lg *LineGraph) SetLineWidth(width float32) {

	lg.mat.SetLineWidth(width)
}

func (lg *LineGraph) updateData() {

	lines := 1
	if lg.chart.scaleX != nil {
		lines = lg.chart.scaleX.lines
	}
	step := 1.0 / (float32(lines) * lg.chart.countStepX)

	positions := math32.NewArrayF32(0, 0)
	rangeY := lg.chart.maxY - lg.chart.minY
	for i := 0; i < len(lg.data); i++ {
		px := float32(i) * step
		vy := lg.data[i]
		py := -1 + ((vy - lg.chart.minY) / rangeY)
		positions.Append(px, py, 0)
	}
	lg.vbo.SetBuffer(positions)
}

func (lg *LineGraph) recalc() {

	py := lg.chart.top
	if lg.chart.title != nil {
		py += lg.chart.title.Height()
	}
	px := lg.chart.left
	w := lg.chart.ContentWidth() - lg.chart.left
	h := lg.chart.ContentHeight() - py - lg.chart.bottom
	lg.SetPosition(px, py)
	lg.SetSize(w, h)
}

// RenderSetup is called by the renderer before drawing this graphic
// It overrides the original panel RenderSetup
// Calculates the model matrix and transfer to OpenGL.
func (lg *LineGraph) RenderSetup(gs *gls.GLS, rinfo *core.RenderInfo) {

	//log.Error("LineGraph RenderSetup:%v with/height: %v/%v", lg.posclip, lg.wclip, lg.hclip)
	// Sets model matrix and transfer to shader
	var mm math32.Matrix4
	lg.SetModelMatrix(gs, &mm)
	lg.modelMatrixUni.SetMatrix4(&mm)
	lg.modelMatrixUni.Transfer(gs)

	// Sets bounds in OpenGL window coordinates and transfer to shader
	_, _, _, height := gs.GetViewport()
	lg.bounds.Set(lg.pospix.X, float32(height)-lg.pospix.Y, lg.width, lg.height)
	lg.bounds.Transfer(gs)
}

//
//
// Chart material (for lines)
//
//
type chartMaterial struct {
	material.Material                // Embedded material
	color             *gls.Uniform3f // Emissive color uniform
}

func (cm *chartMaterial) Init(color *math32.Color) {

	cm.Material.Init()
	cm.SetShader("shaderChart")

	// Creates uniforms and adds to material
	cm.color = gls.NewUniform3f("MatColor")

	// Set initial values
	cm.color.SetColor(color)
}

func (cm *chartMaterial) RenderSetup(gs *gls.GLS) {

	cm.Material.RenderSetup(gs)
	cm.color.Transfer(gs)
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
uniform vec3 MatColor;

// Outputs for fragment shader
out vec3 Color;

void main() {

    Color = MatColor;

    // Set position
    vec4 pos = vec4(VertexPosition.xyz, 1);
	vec4 posclip = ModelMatrix * pos;
    gl_Position = posclip;
}
`

//
// Fragment Shader template
//
const shaderChartFrag = `
#version {{.Version}}

// Input uniforms from vertex shader
in vec3 Color;

// Input uniforms
uniform vec4 Bounds;

// Output
out vec4 FragColor;

void main() {

    // Discard fragment outside of the received bounds in OpenGL window pixel coordinates
    // Bounds[0] - x
    // Bounds[1] - y
    // Bounds[2] - width
    // Bounds[3] - height
    if (gl_FragCoord.x < Bounds[0] || gl_FragCoord.x > Bounds[0] + Bounds[2]) {
        discard;
    }
    if (gl_FragCoord.y > Bounds[1] || gl_FragCoord.y < Bounds[1] - Bounds[3]) {
        discard;
    }

    FragColor = vec4(Color, 1.0);
}
`
