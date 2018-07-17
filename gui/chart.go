// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"fmt"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer/shaders"
	"math"
)

func init() {
	shaders.AddShader("shaderChartVertex", shaderChartVertex)
	shaders.AddShader("shaderChartFrag", shaderChartFrag)
	shaders.AddProgram("shaderChart", "shaderChartVertex", "shaderChartFrag")
}

//
// Chart implements a panel which can contain a title, an x scale,
// an y scale and several graphs
//
type Chart struct {
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
	fontSizeX  float64      // X scale label font size
	fontSizeY  float64      // Y scale label font size
	title      *Label       // Optional title label
	scaleX     *chartScaleX // X scale panel
	scaleY     *chartScaleY // Y scale panel
	labelsX    []*Label     // Array of scale X labels
	labelsY    []*Label     // Array of scale Y labels
	graphs     []*Graph     // Array of line graphs
}

const (
	deltaLine = 0.001 // Delta in NDC for lines over the boundary
)

// NewChart creates and returns a new chart panel with
// the specified dimensions in pixels.
func NewChart(width, height float32) *Chart {

	ch := new(Chart)
	ch.Init(width, height)
	return ch
}

// Init initializes a new chart with the specified width and height
// It is normally used to initialize a Chart embedded in a struct
func (ch *Chart) Init(width float32, height float32) {

	ch.Panel.Initialize(width, height)
	ch.left = 40
	ch.bottom = 20
	ch.top = 10
	ch.firstX = 0
	ch.stepX = 1
	ch.countStepX = 1
	ch.minY = -10.0
	ch.maxY = 10.0
	ch.autoY = false
	ch.formatX = "%v"
	ch.formatY = "%v"
	ch.fontSizeX = 14
	ch.fontSizeY = 14
	ch.Subscribe(OnResize, ch.onResize)
}

// SetTitle sets the chart title text and font size.
// To remove the title pass an empty string
func (ch *Chart) SetTitle(title string, size float64) {

	// Remove title
	if title == "" {
		if ch.title != nil {
			ch.Remove(ch.title)
			ch.title.Dispose()
			ch.title = nil
			ch.recalc()
		}
		return
	}

	// Sets title
	if ch.title == nil {
		ch.title = NewLabel(title)
		ch.title.SetColor4(math32.NewColor4("black"))
		ch.Add(ch.title)
	}
	ch.title.SetText(title)
	ch.title.SetFontSize(size)
	ch.recalc()
}

// SetMarginY sets the y scale margin
func (ch *Chart) SetMarginY(left float32) {

	ch.left = left
	ch.recalc()
}

// SetMarginX sets the x scale margin
func (ch *Chart) SetMarginX(bottom float32) {

	ch.bottom = bottom
	ch.recalc()
}

// SetFormatX sets the string format of the X scale labels
func (ch *Chart) SetFormatX(format string) {

	ch.formatX = format
	ch.updateLabelsX()
}

// SetFormatY sets the string format of the Y scale labels
func (ch *Chart) SetFormatY(format string) {

	ch.formatY = format
	ch.updateLabelsY()
}

// SetFontSizeX sets the font size for the x scale labels
func (ch *Chart) SetFontSizeX(size float64) {

	ch.fontSizeX = size
	for i := 0; i < len(ch.labelsX); i++ {
		ch.labelsX[i].SetFontSize(ch.fontSizeX)
	}
}

// SetFontSizeY sets the font size for the y scale labels
func (ch *Chart) SetFontSizeY(size float64) {

	ch.fontSizeY = size
	for i := 0; i < len(ch.labelsY); i++ {
		ch.labelsY[i].SetFontSize(ch.fontSizeY)
	}
}

// SetScaleX sets the X scale number of lines, lines color and label font size
func (ch *Chart) SetScaleX(lines int, color *math32.Color) {

	if ch.scaleX != nil {
		ch.ClearScaleX()
	}

	// Add scale lines
	ch.scaleX = newChartScaleX(ch, lines, color)
	ch.Add(ch.scaleX)

	// Add scale labels
	// The positions of the labels will be set by 'recalc()'
	value := ch.firstX
	for i := 0; i < lines; i++ {
		l := NewLabel(fmt.Sprintf(ch.formatX, value))
		l.SetColor4(math32.NewColor4("black"))
		l.SetFontSize(ch.fontSizeX)
		ch.Add(l)
		ch.labelsX = append(ch.labelsX, l)
		value += ch.stepX
	}
	ch.recalc()
}

// ClearScaleX removes the X scale if it was previously set
func (ch *Chart) ClearScaleX() {

	if ch.scaleX == nil {
		return
	}

	// Remove and dispose scale lines
	ch.Remove(ch.scaleX)
	ch.scaleX.Dispose()

	// Remove and dispose scale labels
	for i := 0; i < len(ch.labelsX); i++ {
		label := ch.labelsX[i]
		ch.Remove(label)
		label.Dispose()
	}
	ch.labelsX = ch.labelsX[0:0]
	ch.scaleX = nil
}

// SetScaleY sets the Y scale number of lines and color
func (ch *Chart) SetScaleY(lines int, color *math32.Color) {

	if ch.scaleY != nil {
		ch.ClearScaleY()
	}

	if lines < 2 {
		lines = 2
	}

	// Add scale lines
	ch.scaleY = newChartScaleY(ch, lines, color)
	ch.Add(ch.scaleY)

	// Add scale labels
	// The position of the labels will be set by 'recalc()'
	value := ch.minY
	step := (ch.maxY - ch.minY) / float32(lines-1)
	for i := 0; i < lines; i++ {
		l := NewLabel(fmt.Sprintf(ch.formatY, value))
		l.SetColor4(math32.NewColor4("black"))
		l.SetFontSize(ch.fontSizeY)
		ch.Add(l)
		ch.labelsY = append(ch.labelsY, l)
		value += step
	}
	ch.recalc()
}

// ClearScaleY removes the Y scale if it was previously set
func (ch *Chart) ClearScaleY() {

	if ch.scaleY == nil {
		return
	}

	// Remove and dispose scale lines
	ch.Remove(ch.scaleY)
	ch.scaleY.Dispose()

	// Remove and dispose scale labels
	for i := 0; i < len(ch.labelsY); i++ {
		label := ch.labelsY[i]
		ch.Remove(label)
		label.Dispose()
	}
	ch.labelsY = ch.labelsY[0:0]
	ch.scaleY = nil
}

// SetRangeX sets the X scale labels and range per step
// firstX is the value of first label of the x scale
// stepX is the step to be added to get the next x scale label
// countStepX is the number of elements of the data buffer for each line step
func (ch *Chart) SetRangeX(firstX float32, stepX float32, countStepX float32) {

	ch.firstX = firstX
	ch.stepX = stepX
	ch.countStepX = countStepX
	ch.updateGraphs()
}

// SetRangeY sets the minimum and maximum values of the y scale
func (ch *Chart) SetRangeY(min float32, max float32) {

	if ch.autoY {
		return
	}
	ch.minY = min
	ch.maxY = max
	ch.updateGraphs()
}

// SetRangeYauto sets the state of the auto
func (ch *Chart) SetRangeYauto(auto bool) {

	ch.autoY = auto
	if !auto {
		return
	}
	ch.updateGraphs()
}

// RangeY returns the current y range
func (ch *Chart) RangeY() (minY, maxY float32) {

	return ch.minY, ch.maxY
}

// AddLineGraph adds a line graph to the chart
func (ch *Chart) AddLineGraph(color *math32.Color, data []float32) *Graph {

	graph := newGraph(ch, color, data)
	ch.graphs = append(ch.graphs, graph)
	ch.Add(graph)
	ch.recalc()
	ch.updateGraphs()
	return graph
}

// RemoveGraph removes and disposes of the specified graph from the chart
func (ch *Chart) RemoveGraph(g *Graph) {

	ch.Remove(g)
	g.Dispose()
	for pos, current := range ch.graphs {
		if current == g {
			copy(ch.graphs[pos:], ch.graphs[pos+1:])
			ch.graphs[len(ch.graphs)-1] = nil
			ch.graphs = ch.graphs[:len(ch.graphs)-1]
			break
		}
	}
	if !ch.autoY {
		return
	}
	ch.updateGraphs()
}

// updateLabelsX updates the X scale labels text
func (ch *Chart) updateLabelsX() {

	if ch.scaleX == nil {
		return
	}
	pstep := (ch.ContentWidth() - ch.left) / float32(len(ch.labelsX))
	value := ch.firstX
	for i := 0; i < len(ch.labelsX); i++ {
		label := ch.labelsX[i]
		label.SetText(fmt.Sprintf(ch.formatX, value))
		px := ch.left + float32(i)*pstep
		label.SetPosition(px, ch.ContentHeight()-ch.bottom)
		value += ch.stepX
	}
}

// updateLabelsY updates the Y scale labels text and positions
func (ch *Chart) updateLabelsY() {

	if ch.scaleY == nil {
		return
	}

	th := float32(0)
	if ch.title != nil {
		th = ch.title.height
	}

	nlines := ch.scaleY.lines
	vstep := (ch.maxY - ch.minY) / float32(nlines-1)
	pstep := (ch.ContentHeight() - th - ch.top - ch.bottom) / float32(nlines-1)
	value := ch.minY
	for i := 0; i < nlines; i++ {
		label := ch.labelsY[i]
		label.SetText(fmt.Sprintf(ch.formatY, value))
		px := ch.left - 4 - label.Width()
		if px < 0 {
			px = 0
		}
		py := ch.ContentHeight() - ch.bottom - float32(i)*pstep
		label.SetPosition(px, py-label.Height()/2)
		value += vstep
	}
}

// calcRangeY calculates the minimum and maximum y values for all graphs
func (ch *Chart) calcRangeY() {

	if !ch.autoY || len(ch.graphs) == 0 {
		return
	}
	minY := float32(math.MaxFloat32)
	maxY := -float32(math.MaxFloat32)
	for g := 0; g < len(ch.graphs); g++ {
		graph := ch.graphs[g]
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
	ch.minY = minY
	ch.maxY = maxY
}

// updateGraphs should be called when the range the scales change or
// any graph data changes
func (ch *Chart) updateGraphs() {

	ch.calcRangeY()
	ch.updateLabelsX()
	ch.updateLabelsY()
	for i := 0; i < len(ch.graphs); i++ {
		g := ch.graphs[i]
		g.updateData()
	}
}

// onResize process OnResize events for this chart
func (ch *Chart) onResize(evname string, ev interface{}) {

	ch.recalc()
}

// recalc recalculates the positions of the inner panels
func (ch *Chart) recalc() {

	// Center title position
	if ch.title != nil {
		xpos := (ch.ContentWidth() - ch.title.width) / 2
		ch.title.SetPositionX(xpos)
	}

	// Recalc scale X and its labels
	if ch.scaleX != nil {
		ch.scaleX.recalc()
		ch.updateLabelsX()
	}

	// Recalc scale Y and its labels
	if ch.scaleY != nil {
		ch.scaleY.recalc()
		ch.updateLabelsY()
	}

	// Recalc graphs
	for i := 0; i < len(ch.graphs); i++ {
		g := ch.graphs[i]
		g.recalc()
		ch.SetTopChild(g)
	}
}

//
// chartScaleX is a panel with GL_LINES geometry which draws the chart X horizontal scale axis,
// vertical lines and line labels.
//
type chartScaleX struct {
	Panel                   // Embedded panel
	chart     *Chart        // Container chart
	lines     int           // Number of vertical lines
	mat       chartMaterial // Chart material
	uniBounds gls.Uniform   // Bounds uniform location cache
}

// newChartScaleX creates and returns a pointer to a new chartScaleX for the specified
// chart, number of lines and color
func newChartScaleX(chart *Chart, lines int, color *math32.Color) *chartScaleX {

	sx := new(chartScaleX)
	sx.chart = chart
	sx.lines = lines
	sx.uniBounds.Init("Bounds")

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
	geom.AddVBO(gls.NewVBO(positions).AddAttrib(gls.VertexPosition))

	// Initializes the panel graphic
	gr := graphic.NewGraphic(geom, gls.LINES)
	sx.mat.Init(color)
	gr.AddMaterial(sx, &sx.mat, 0, 0)
	sx.Panel.InitializeGraphic(chart.ContentWidth(), chart.ContentHeight(), gr)

	sx.recalc()
	return sx
}

// recalc recalculates the position and size of this scale inside its parent
func (sx *chartScaleX) recalc() {

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
func (sx *chartScaleX) RenderSetup(gs *gls.GLS, rinfo *core.RenderInfo) {

	// Sets model matrix
	var mm math32.Matrix4
	sx.SetModelMatrix(gs, &mm)

	// Transfer model matrix uniform
	location := sx.uniMatrix.Location(gs)
	gs.UniformMatrix4fv(location, 1, false, &mm[0])

	// Sets bounds in OpenGL window coordinates and transfer to shader
	_, _, _, height := gs.GetViewport()
	location = sx.uniBounds.Location(gs)
	gs.Uniform4f(location, sx.pospix.X, float32(height)-sx.pospix.Y, sx.width, sx.height)
}

//
// ChartScaleY is a panel with LINE geometry which draws the chart Y vertical scale axis,
// horizontal and labels.
//
type chartScaleY struct {
	Panel                   // Embedded panel
	chart     *Chart        // Container chart
	lines     int           // Number of horizontal lines
	mat       chartMaterial // Chart material
	uniBounds gls.Uniform   // Bounds uniform location cache
}

// newChartScaleY creates and returns a pointer to a new chartScaleY for the specified
// chart, number of lines and color
func newChartScaleY(chart *Chart, lines int, color *math32.Color) *chartScaleY {

	if lines < 2 {
		lines = 2
	}
	sy := new(chartScaleY)
	sy.chart = chart
	sy.lines = lines
	sy.uniBounds.Init("Bounds")

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
	geom.AddVBO(gls.NewVBO(positions).AddAttrib(gls.VertexPosition))

	// Initializes the panel with this graphic
	gr := graphic.NewGraphic(geom, gls.LINES)
	sy.mat.Init(color)
	gr.AddMaterial(sy, &sy.mat, 0, 0)
	sy.Panel.InitializeGraphic(chart.ContentWidth(), chart.ContentHeight(), gr)

	sy.recalc()
	return sy
}

// recalc recalculates the position and size of this scale inside its parent
func (sy *chartScaleY) recalc() {

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
func (sy *chartScaleY) RenderSetup(gs *gls.GLS, rinfo *core.RenderInfo) {

	// Sets model matrix
	var mm math32.Matrix4
	sy.SetModelMatrix(gs, &mm)

	// Transfer model matrix uniform
	location := sy.uniMatrix.Location(gs)
	gs.UniformMatrix4fv(location, 1, false, &mm[0])

	// Sets bounds in OpenGL window coordinates and transfer to shader
	_, _, _, height := gs.GetViewport()
	location = sy.uniBounds.Location(gs)
	gs.Uniform4f(location, sy.pospix.X, float32(height)-sy.pospix.Y, sy.width, sy.height)
}

//
// Graph is the GUI element that represents a single plotted function.
// A Chart has an array of Graph objects.
//
type Graph struct {
	Panel                   // Embedded panel
	chart     *Chart        // Container chart
	color     math32.Color  // Line color
	data      []float32     // Data y
	mat       chartMaterial // Chart material
	vbo       *gls.VBO
	positions math32.ArrayF32
	uniBounds gls.Uniform // Bounds uniform location cache
}

// newGraph creates and returns a pointer to a new graph for the specified chart
func newGraph(chart *Chart, color *math32.Color, data []float32) *Graph {

	lg := new(Graph)
	lg.uniBounds.Init("Bounds")
	lg.chart = chart
	lg.color = *color
	lg.data = data

	// Creates geometry and adds VBO with positions
	geom := geometry.NewGeometry()
	lg.positions = math32.NewArrayF32(0, 0)
	lg.vbo = gls.NewVBO(lg.positions).AddAttrib(gls.VertexPosition)
	geom.AddVBO(lg.vbo)

	// Initializes the panel with this graphic
	gr := graphic.NewGraphic(geom, gls.LINE_STRIP)
	lg.mat.Init(&lg.color)
	gr.AddMaterial(lg, &lg.mat, 0, 0)
	lg.Panel.InitializeGraphic(lg.chart.ContentWidth(), lg.chart.ContentHeight(), gr)

	lg.SetData(data)
	return lg
}

// SetColor sets the color of the graph
func (lg *Graph) SetColor(color *math32.Color) {

	lg.color = *color
}

// SetData sets the graph data
func (lg *Graph) SetData(data []float32) {

	lg.data = data
	lg.updateData()
}

// SetLineWidth sets the graph line width
func (lg *Graph) SetLineWidth(width float32) {

	lg.mat.SetLineWidth(width)
}

// updateData regenerates the lines for the current data
func (lg *Graph) updateData() {

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
	lg.SetChanged(true)
}

// recalc recalculates the position and width of the this panel
func (lg *Graph) recalc() {

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
func (lg *Graph) RenderSetup(gs *gls.GLS, rinfo *core.RenderInfo) {

	// Sets model matrix
	var mm math32.Matrix4
	lg.SetModelMatrix(gs, &mm)

	// Transfer model matrix uniform
	location := lg.uniMatrix.Location(gs)
	gs.UniformMatrix4fv(location, 1, false, &mm[0])

	// Sets bounds in OpenGL window coordinates and transfer to shader
	_, _, _, height := gs.GetViewport()
	location = lg.uniBounds.Location(gs)
	gs.Uniform4f(location, lg.pospix.X, float32(height)-lg.pospix.Y, lg.width, lg.height)
}

//
// Chart material
//
type chartMaterial struct {
	material.Material              // Embedded material
	color             math32.Color // emissive color
	uniColor          gls.Uniform  // color uniform location cache
}

func (cm *chartMaterial) Init(color *math32.Color) {

	cm.Material.Init()
	cm.SetShader("shaderChart")
	cm.SetShaderUnique(true)
	cm.uniColor.Init("MatColor")
	cm.color = *color
}

func (cm *chartMaterial) RenderSetup(gs *gls.GLS) {

	cm.Material.RenderSetup(gs)
	gs.Uniform3f(cm.uniColor.Location(gs), cm.color.R, cm.color.G, cm.color.B)
}

//
// Vertex Shader template
//
const shaderChartVertex = `
// Vertex attributes
#include <attributes>

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
