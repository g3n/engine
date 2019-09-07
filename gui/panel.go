// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
)

/*********************************************

 Panel areas:
 +------------------------------------------+
 |  Margin area                             |
 |  +------------------------------------+  |
 |  |  Border area                       |  |
 |  |  +------------------------------+  |  |
 |  |  | Padding area                 |  |  |
 |  |  |  +------------------------+  |  |  |
 |  |  |  | Content area           |  |  |  |
 |  |  |  |                        |  |  |  |
 |  |  |  |                        |  |  |  |
 |  |  |  +------------------------+  |  |  |
 |  |  |                              |  |  |
 |  |  +------------------------------+  |  |
 |  |                                    |  |
 |  +------------------------------------+  |
 |                                          |
 +------------------------------------------+

*********************************************/

// IPanel is the interface for all panel types
type IPanel interface {
	graphic.IGraphic
	GetPanel() *Panel
	Width() float32
	Height() float32
	Enabled() bool
	SetEnabled(bool)
	SetLayout(ILayout)
	InsideBorders(x, y float32) bool
	SetZLayerDelta(zLayerDelta int)
	ZLayerDelta() int

	// TODO these methods here should probably be defined in INode
	SetPosition(x, y float32)
	SetPositionX(x float32)
	SetPositionY(y float32)
	SetPositionZ(y float32)
}

// Panel is 2D rectangular graphic which by default has a quad (2 triangles) geometry.
// When using the default geometry, a panel has margins, borders, paddings
// and a content area. The content area can be associated with a texture
// It is the building block of most GUI widgets.
type Panel struct {
	*graphic.Graphic                    // Embedded graphic
	mat              *material.Material // panel material
	zLayerDelta      int                // Z-layer relative to parent

	bounded bool // Whether panel is bounded by its parent
	enabled bool // Whether event should be processed for this panel

	layout       ILayout     // current layout for children
	layoutParams interface{} // current layout parameters used by container panel

	marginSizes  RectBounds // external margin sizes in pixel coordinates
	borderSizes  RectBounds // border sizes in pixel coordinates
	paddingSizes RectBounds // padding sizes in pixel coordinates
	content      Rect       // current content rectangle in pixel coordinates

	// Absolute screen position and external size in pixels
	pospix math32.Vector3
	width  float32
	height float32

	xmin float32 // minimum absolute x this panel can use
	xmax float32 // maximum absolute x this panel can use
	ymin float32 // minimum absolute y this panel can use
	ymax float32 // maximum absolute y this panel can use

	// Uniforms sent to shader
	uniMatrix gls.Uniform // model matrix uniform location cache
	uniPanel  gls.Uniform // panel parameters uniform location cache
	udata     struct {    // Combined uniform data 8 * vec4
		bounds        math32.Vector4 // panel bounds in texture coordinates
		borders       math32.Vector4 // panel borders in texture coordinates
		paddings      math32.Vector4 // panel paddings in texture coordinates
		content       math32.Vector4 // panel content area in texture coordinates
		bordersColor  math32.Color4  // panel border color
		paddingsColor math32.Color4  // panel padding color
		contentColor  math32.Color4  // panel content color
		textureValid  float32        // texture valid flag (bool)
		dummy         [3]float32     // complete 8 * vec4
	}
}

// PanelStyle contains all the styling attributes of a Panel.
type PanelStyle struct {
	Margin      RectBounds
	Border      RectBounds
	Padding     RectBounds
	BorderColor math32.Color4
	BgColor     math32.Color4
}

// BasicStyle extends PanelStyle by adding a foreground color.
// Many GUI components can be styled using BasicStyle or redeclared versions thereof (e.g. ButtonStyle)
type BasicStyle struct {
	PanelStyle
	FgColor math32.Color4
}

// Quad geometry shared by ALL Panels
var panelQuadGeometry *geometry.Geometry

// NewPanel creates and returns a pointer to a new panel with the
// specified dimensions in pixels and a default quad geometry
func NewPanel(width, height float32) *Panel {

	p := new(Panel)
	p.Initialize(p, width, height)
	return p
}

// Initialize initializes this panel and is normally used by other types which embed a panel.
func (p *Panel) Initialize(ipan IPanel, width, height float32) { // TODO rename to Init

	p.width = width
	p.height = height

	// If first time, create panel quad geometry
	if panelQuadGeometry == nil {

		// Builds array with vertex positions and texture coordinates
		positions := math32.NewArrayF32(0, 20)
		positions.Append(
			0, 0, 0, 0, 1,
			0, -1, 0, 0, 0,
			1, -1, 0, 1, 0,
			1, 0, 0, 1, 1,
		)
		// Builds array of indices
		indices := math32.NewArrayU32(0, 6)
		indices.Append(0, 1, 2, 0, 2, 3)

		// Creates geometry
		geom := geometry.NewGeometry()
		geom.SetIndices(indices)
		geom.AddVBO(gls.NewVBO(positions).
			AddAttrib(gls.VertexPosition).
			AddAttrib(gls.VertexTexcoord),
		)
		panelQuadGeometry = geom
	}

	// Initialize material
	p.mat = material.NewMaterial()
	p.mat.SetUseLights(material.UseLightNone)
	p.mat.SetShader("panel")
	p.mat.SetShaderUnique(true)

	// For now set all panels as transparent by default
	// This means they are all rendered back to front, after and on top of everything else
	// TODO if we know a panel is opaque, setting it as such will improve rendering performance
	p.mat.SetTransparent(true)

	// Initialize graphic
	p.Graphic = graphic.NewGraphic(ipan, panelQuadGeometry.Incref(), gls.TRIANGLES)
	p.AddMaterial(p, p.mat, 0, 0)

	// Initialize uniforms location caches
	p.uniMatrix.Init("ModelMatrix")
	p.uniPanel.Init("Panel")

	// Set defaults
	p.udata.bordersColor = math32.Color4{0, 0, 0, 1}
	p.bounded = true
	p.enabled = true
	p.resize(width, height, true)
}

// InitializeGraphic initializes this panel with a different graphic
func (p *Panel) InitializeGraphic(width, height float32, gr *graphic.Graphic) {

	p.Graphic = gr
	p.width = width
	p.height = height

	// Initializes uniforms location caches
	p.uniMatrix.Init("ModelMatrix")
	p.uniPanel.Init("Panel")

	// Set defaults
	p.udata.bordersColor = math32.Color4{0, 0, 0, 1}
	p.bounded = true
	p.enabled = true
	p.resize(width, height, true)
}

// GetPanel satisfies the IPanel interface and
// returns pointer to this panel
func (p *Panel) GetPanel() *Panel {

	return p
}

// Material returns a pointer for this panel's Material
func (p *Panel) Material() *material.Material { // TODO remove - allow for setting and getting a single texture

	return p.mat
}

// SetTopChild sets the Z coordinate of the specified panel to
// be on top of all other children of this panel.
// The function does not check if the specified panel is a
// child of this one.
func (p *Panel) SetTopChild(ipan IPanel) {

	// Remove panel and if found appends to the end
	found := p.Remove(ipan)
	if found {
		p.Add(ipan)
		p.SetChanged(true)
	}
}

// SetZLayerDelta sets the Z-layer of the panel relative to its parent.
func (p *Panel) SetZLayerDelta(zLayerDelta int) {

	p.zLayerDelta = zLayerDelta
}

// ZLayerDelta returns the Z-layer of the panel relative to its parent.
func (p *Panel) ZLayerDelta() int {

	return p.zLayerDelta
}

// SetPosition sets this panel absolute position in pixel coordinates
// from left to right and from top to bottom of the screen.
func (p *Panel) SetPosition(x, y float32) {

	p.Node.SetPositionX(math32.Round(x))
	p.Node.SetPositionY(math32.Round(y))
}

// SetSize sets this panel external width and height in pixels.
func (p *Panel) SetSize(width, height float32) {

	if width < 0 {
		log.Warn("Invalid panel width:%v", width)
		width = 0
	}
	if height < 0 {
		log.Warn("Invalid panel height:%v", height)
		height = 0
	}
	p.resize(width, height, true)
}

// SetWidth sets this panel external width in pixels.
// The internal panel areas and positions are recalculated
func (p *Panel) SetWidth(width float32) {

	p.SetSize(width, p.height)
}

// SetHeight sets this panel external height in pixels.
// The internal panel areas and positions are recalculated
func (p *Panel) SetHeight(height float32) {

	p.SetSize(p.width, height)
}

// SetContentAspectWidth sets the width of the content area of the panel
// to the specified value and adjusts its height to keep the same aspect radio.
func (p *Panel) SetContentAspectWidth(width float32) {

	aspect := p.content.Width / p.content.Height
	height := width / aspect
	p.SetContentSize(width, height)
}

// SetContentAspectHeight sets the height of the content area of the panel
// to the specified value and adjusts its width to keep the same aspect ratio.
func (p *Panel) SetContentAspectHeight(height float32) {

	aspect := p.content.Width / p.content.Height
	width := height / aspect
	p.SetContentSize(width, height)
}

// Size returns this panel current external width and height in pixels
func (p *Panel) Size() (float32, float32) {

	return p.width, p.height
}

// Width returns the current panel external width in pixels
func (p *Panel) Width() float32 {

	return p.width
}

// Height returns the current panel external height in pixels
func (p *Panel) Height() float32 {

	return p.height
}

// ContentWidth returns the current width of the content area in pixels
func (p *Panel) ContentWidth() float32 {

	return p.content.Width
}

// ContentHeight returns the current height of the content area in pixels
func (p *Panel) ContentHeight() float32 {

	return p.content.Height
}

// SetMargins set this panel margin sizes in pixels
// and recalculates the panel external size
func (p *Panel) SetMargins(top, right, bottom, left float32) {

	p.marginSizes.Set(top, right, bottom, left)
	p.resize(p.calcWidth(), p.calcHeight(), true)
}

// SetMarginsFrom sets this panel margins sizes from the specified
// RectBounds pointer and recalculates the panel external size
func (p *Panel) SetMarginsFrom(src *RectBounds) {

	p.marginSizes = *src
	p.resize(p.calcWidth(), p.calcHeight(), true)
}

// Margins returns the current margin sizes in pixels
func (p *Panel) Margins() RectBounds {

	return p.marginSizes
}

// SetBorders sets this panel border sizes in pixels
// and recalculates the panel external size
func (p *Panel) SetBorders(top, right, bottom, left float32) {

	p.borderSizes.Set(top, right, bottom, left)
	p.resize(p.calcWidth(), p.calcHeight(), true)
}

// SetBordersFrom sets this panel border sizes from the specified
// RectBounds pointer and recalculates the panel size
func (p *Panel) SetBordersFrom(src *RectBounds) {

	p.borderSizes = *src
	p.resize(p.calcWidth(), p.calcHeight(), true)
}

// Borders returns this panel current border sizes
func (p *Panel) Borders() RectBounds {

	return p.borderSizes
}

// SetPaddings sets the panel padding sizes in pixels
func (p *Panel) SetPaddings(top, right, bottom, left float32) {

	p.paddingSizes.Set(top, right, bottom, left)
	p.resize(p.calcWidth(), p.calcHeight(), true)
}

// SetPaddingsFrom sets this panel padding sizes from the specified
// RectBounds pointer and recalculates the panel size
func (p *Panel) SetPaddingsFrom(src *RectBounds) {

	p.paddingSizes = *src
	p.resize(p.calcWidth(), p.calcHeight(), true)
}

// Paddings returns this panel padding sizes in pixels
func (p *Panel) Paddings() RectBounds {

	return p.paddingSizes
}

// SetBordersColor sets the color of this panel borders
// The borders opacity is set to 1.0 (full opaque)
func (p *Panel) SetBordersColor(color *math32.Color) {

	p.udata.bordersColor = math32.Color4{color.R, color.G, color.B, 1}
	p.SetChanged(true)
}

// SetBordersColor4 sets the color and opacity of this panel borders
func (p *Panel) SetBordersColor4(color *math32.Color4) {

	p.udata.bordersColor = *color
	p.SetChanged(true)
}

// BordersColor4 returns current border color
func (p *Panel) BordersColor4() math32.Color4 {

	return p.udata.bordersColor
}

// SetPaddingsColor sets the color of this panel paddings.
func (p *Panel) SetPaddingsColor(color *math32.Color) {

	p.udata.paddingsColor = math32.Color4{color.R, color.G, color.B, 1}
	p.SetChanged(true)
}

// SetColor sets the color of the panel paddings and content area
func (p *Panel) SetColor(color *math32.Color) *Panel {

	p.udata.paddingsColor = math32.Color4{color.R, color.G, color.B, 1}
	p.udata.contentColor = p.udata.paddingsColor
	p.SetChanged(true)
	return p
}

// SetColor4 sets the color of the panel paddings and content area
func (p *Panel) SetColor4(color *math32.Color4) *Panel {

	p.udata.paddingsColor = *color
	p.udata.contentColor = *color
	p.SetChanged(true)
	return p
}

// Color4 returns the current color of the panel content area
func (p *Panel) Color4() math32.Color4 {

	return p.udata.contentColor
}

// ApplyStyle applies the provided PanelStyle to the panel
func (p *Panel) ApplyStyle(ps *PanelStyle) {

	p.udata.bordersColor = ps.BorderColor
	p.udata.paddingsColor = ps.BgColor
	p.udata.contentColor = ps.BgColor
	p.marginSizes = ps.Margin
	p.borderSizes = ps.Border
	p.paddingSizes = ps.Padding
	p.resize(p.calcWidth(), p.calcHeight(), true)
}

// SetContentSize sets this panel content size to the specified dimensions.
// The external size of the panel may increase or decrease to acomodate
// the new content size.
func (p *Panel) SetContentSize(width, height float32) {

	p.setContentSize(width, height, true)
}

// SetContentWidth sets this panel content width to the specified dimension in pixels.
// The external size of the panel may increase or decrease to accommodate the new width
func (p *Panel) SetContentWidth(width float32) {

	p.SetContentSize(width, p.content.Height)
}

// SetContentHeight sets this panel content height to the specified dimension in pixels.
// The external size of the panel may increase or decrease to accommodate the new width
func (p *Panel) SetContentHeight(height float32) {

	p.SetContentSize(p.content.Width, height)
}

// MinWidth returns the minimum width of this panel (ContentWidth = 0)
func (p *Panel) MinWidth() float32 {

	return p.paddingSizes.Left + p.paddingSizes.Right +
		p.borderSizes.Left + p.borderSizes.Right +
		p.marginSizes.Left + p.marginSizes.Right
}

// MinHeight returns the minimum height of this panel (ContentHeight = 0)
func (p *Panel) MinHeight() float32 {

	return p.paddingSizes.Top + p.paddingSizes.Bottom +
		p.borderSizes.Top + p.borderSizes.Bottom +
		p.marginSizes.Top + p.marginSizes.Bottom
}

// Pospix returns this panel absolute coordinate in pixels
func (p *Panel) Pospix() math32.Vector3 {

	return p.pospix
}

// Add adds a child panel to this one
// This overrides the Node method to enforce that IPanels can only have IPanels as children
func (p *Panel) Add(ichild IPanel) *Panel {

	p.Node.Add(ichild)
	if p.layout != nil {
		p.layout.Recalc(p)
	}
	return p
}

// Remove removes the specified child from this panel
func (p *Panel) Remove(ichild IPanel) bool {

	res := p.Node.Remove(ichild)
	if res {
		if p.layout != nil {
			p.layout.Recalc(p)
		}
	}
	return res
}

// Bounded returns this panel bounded state
func (p *Panel) Bounded() bool {

	return p.bounded
}

// SetBounded sets this panel bounded state
func (p *Panel) SetBounded(bounded bool) {

	p.bounded = bounded
	p.SetChanged(true)
}

// UpdateMatrixWorld overrides the standard core.Node version which is called by
// the Engine before rendering the frame.
func (p *Panel) UpdateMatrixWorld() {

	par := p.Parent()
	if par == nil {
		p.updateBounds(nil)
		// Panel has parent
	} else {
		parpan, ok := par.(IPanel)
		if ok {
			p.updateBounds(parpan.GetPanel())
		} else {
			p.updateBounds(nil)
		}
	}
	// Update this panel children
	for _, ichild := range p.Children() {
		ichild.UpdateMatrixWorld()
	}
}

// ContainsPosition returns indication if this panel contains
// the specified screen position in pixels.
func (p *Panel) ContainsPosition(x, y float32) bool {

	if x < p.pospix.X || y < p.pospix.Y || x >= (p.pospix.X+p.width) || y >= (p.pospix.Y+p.height) {
		return false
	}
	return true
}

// InsideBorders returns indication if the specified screen
// position in pixels is inside the panel borders, including the borders width.
// Unlike "ContainsPosition" is does not consider the panel margins.
func (p *Panel) InsideBorders(x, y float32) bool {

	if x < (p.pospix.X+p.marginSizes.Left) || x >= (p.pospix.X+p.width-p.marginSizes.Right) ||
		y < (p.pospix.Y+p.marginSizes.Top) || y >= (p.pospix.Y+p.height-p.marginSizes.Bottom) {
		return false
	}
	return true
}

// Intersects returns if this panel intersects with the other panel
func (p *Panel) Intersects(other *Panel) bool {

	if p.pospix.X+p.width <= other.pospix.X || other.pospix.X+other.width <= p.pospix.X ||
		p.pospix.Y+p.height <= other.pospix.Y || other.pospix.Y+other.height <= p.pospix.Y {
		return false
	}
	return true
}

// SetEnabled sets the panel enabled state
// A disabled panel do not process key or mouse events.
func (p *Panel) SetEnabled(state bool) {

	p.enabled = state
	p.Dispatch(OnEnable, nil)
}

// Enabled returns the current enabled state of this panel
func (p *Panel) Enabled() bool {

	return p.enabled
}

// SetLayout sets the layout to use to position the children of this panel
// To remove the layout, call this function passing nil as parameter.
func (p *Panel) SetLayout(ilayout ILayout) {

	p.layout = ilayout
	if p.layout != nil {
		p.layout.Recalc(p)
	}
}

// Layout returns this panel current layout
func (p *Panel) Layout() ILayout {

	return p.layout
}

// SetLayoutParams sets the layout parameters for this panel
func (p *Panel) SetLayoutParams(params interface{}) {

	p.layoutParams = params
}

// LayoutParams returns this panel current layout parameters
func (p *Panel) LayoutParams() interface{} {

	return p.layoutParams
}

// ContentCoords converts the specified window absolute coordinates in pixels
// (as informed by OnMouse event) to this panel internal content area pixel coordinates
func (p *Panel) ContentCoords(wx, wy float32) (float32, float32) {

	cx := wx - p.pospix.X -
		p.paddingSizes.Left -
		p.borderSizes.Left -
		p.marginSizes.Left
	cy := wy - p.pospix.Y -
		p.paddingSizes.Top -
		p.borderSizes.Top -
		p.marginSizes.Top
	return cx, cy
}

// setContentSize is an internal version of SetContentSize() which allows
// to determine if the panel will recalculate its layout and dispatch event.
// It is normally used by layout managers when setting the panel content size
// to avoid another invocation of the layout manager.
func (p *Panel) setContentSize(width, height float32, dispatch bool) {

	// Calculates the new desired external width and height
	eWidth := width +
		p.paddingSizes.Left + p.paddingSizes.Right +
		p.borderSizes.Left + p.borderSizes.Right +
		p.marginSizes.Left + p.marginSizes.Right
	eHeight := height +
		p.paddingSizes.Top + p.paddingSizes.Bottom +
		p.borderSizes.Top + p.borderSizes.Bottom +
		p.marginSizes.Top + p.marginSizes.Bottom
	p.resize(eWidth, eHeight, dispatch)
}

// updateBounds is called by UpdateMatrixWorld() and calculates this panel
// bounds considering the bounds of its parent
func (p *Panel) updateBounds(par *Panel) {

	// If this panel has no parent, its pixel position is its Position
	if par == nil {
		p.pospix.X = p.Position().X
		p.pospix.Y = p.Position().Y
		// If this panel is bounded to its parent, its coordinates are relative
		// to the parent internal content rectangle.
	} else if p.bounded {
		p.pospix.X = p.Position().X + par.pospix.X + par.marginSizes.Left + par.borderSizes.Left + par.paddingSizes.Left
		p.pospix.Y = p.Position().Y + par.pospix.Y + par.marginSizes.Top + par.borderSizes.Top + par.paddingSizes.Top
		// Otherwise its coordinates are relative to the parent outer coordinates.
	} else {
		p.pospix.X = p.Position().X + par.pospix.X
		p.pospix.Y = p.Position().Y + par.pospix.Y
	}
	// Maximum x,y coordinates for this panel
	p.xmin = p.pospix.X
	p.ymin = p.pospix.Y
	p.xmax = p.pospix.X + p.width
	p.ymax = p.pospix.Y + p.height
	// Set default bounds to be entire panel texture
	p.udata.bounds = math32.Vector4{0, 0, 1, 1}
	// If this panel has no parent or is unbounded then the default bounds are correct
	if par == nil || !p.bounded {
		return
	}
	// From here on panel has parent and is bounded by parent
	// TODO why not use par.xmin, etc here ? shouldn't they be already updated?
	// Get the parent content area minimum and maximum absolute coordinates in pixels
	pxmin := par.pospix.X + par.marginSizes.Left + par.borderSizes.Left + par.paddingSizes.Left
	if pxmin < par.xmin {
		pxmin = par.xmin
	}
	pymin := par.pospix.Y + par.marginSizes.Top + par.borderSizes.Top + par.paddingSizes.Top
	if pymin < par.ymin {
		pymin = par.ymin
	}
	pxmax := par.pospix.X + par.width - (par.marginSizes.Right + par.borderSizes.Right + par.paddingSizes.Right)
	if pxmax > par.xmax {
		pxmax = par.xmax
	}
	pymax := par.pospix.Y + par.height - (par.marginSizes.Bottom + par.borderSizes.Bottom + par.paddingSizes.Bottom)
	if pymax > par.ymax {
		pymax = par.ymax
	}
	// Update this panel minimum x and y coordinates.
	if p.xmin < pxmin {
		p.xmin = pxmin
	}
	if p.ymin < pymin {
		p.ymin = pymin
	}
	// Update this panel maximum x and y coordinates.
	if p.xmax > pxmax {
		p.xmax = pxmax
	}
	if p.ymax > pymax {
		p.ymax = pymax
	}
	// If this panel is bounded to its parent, calculates the bounds
	// for clipping in texture coordinates
	if p.pospix.X < p.xmin {
		p.udata.bounds.X = (p.xmin - p.pospix.X) / p.width
	}
	if p.pospix.Y < p.ymin {
		p.udata.bounds.Y = (p.ymin - p.pospix.Y) / p.height
	}
	if p.pospix.X+p.width > p.xmax {
		p.udata.bounds.Z = (p.xmax - p.pospix.X) / p.width
	}
	if p.pospix.Y+p.height > p.ymax {
		p.udata.bounds.W = (p.ymax - p.pospix.Y) / p.height
	}
}

// calcWidth calculates the panel external width in pixels
func (p *Panel) calcWidth() float32 {

	return p.content.Width +
		p.paddingSizes.Left + p.paddingSizes.Right +
		p.borderSizes.Left + p.borderSizes.Right +
		p.marginSizes.Left + p.marginSizes.Right
}

// calcHeight calculates the panel external height in pixels
func (p *Panel) calcHeight() float32 {

	return p.content.Height +
		p.paddingSizes.Top + p.paddingSizes.Bottom +
		p.borderSizes.Top + p.borderSizes.Bottom +
		p.marginSizes.Top + p.marginSizes.Bottom
}

// resize tries to set the external size of the panel to the specified
// dimensions and recalculates the size and positions of the internal areas.
// The margins, borders and padding sizes are kept and the content
// area size is adjusted. So if the panel is decreased, its minimum
// size is determined by the margins, borders and paddings.
// Normally it should be called with dispatch=true to recalculate the
// panel layout and dispatch OnSize event.
func (p *Panel) resize(width, height float32, dispatch bool) {

	var padding Rect
	var border Rect

	width = math32.Round(width)
	height = math32.Round(height)

	// Adjust content width
	p.content.Width = width - p.marginSizes.Left - p.marginSizes.Right - p.borderSizes.Left - p.borderSizes.Right - p.paddingSizes.Left - p.paddingSizes.Right
	if p.content.Width < 0 {
		p.content.Width = 0
	}
	// Adjusts content height
	p.content.Height = height - p.marginSizes.Top - p.marginSizes.Bottom - p.borderSizes.Top - p.borderSizes.Bottom - p.paddingSizes.Top - p.paddingSizes.Bottom
	if p.content.Height < 0 {
		p.content.Height = 0
	}

	// Adjust other area widths
	padding.Width = p.paddingSizes.Left + p.content.Width + p.paddingSizes.Right
	border.Width = p.borderSizes.Left + padding.Width + p.borderSizes.Right
	// Adjust other area heights
	padding.Height = p.paddingSizes.Top + p.content.Height + p.paddingSizes.Bottom
	border.Height = p.borderSizes.Top + padding.Height + p.borderSizes.Bottom

	// Set area positions
	border.X = p.marginSizes.Left
	border.Y = p.marginSizes.Top
	padding.X = border.X + p.borderSizes.Left
	padding.Y = border.Y + p.borderSizes.Top
	p.content.X = padding.X + p.paddingSizes.Left
	p.content.Y = padding.Y + p.paddingSizes.Top

	// Set final panel dimensions (may be different from requested dimensions)
	p.width = p.marginSizes.Left + border.Width + p.marginSizes.Right
	p.height = p.marginSizes.Top + border.Height + p.marginSizes.Bottom

	// Update border uniform in texture coordinates (0,0 -> 1,1)
	p.udata.borders = math32.Vector4{
		float32(border.X) / float32(p.width),
		float32(border.Y) / float32(p.height),
		float32(border.Width) / float32(p.width),
		float32(border.Height) / float32(p.height),
	}
	// Update padding uniform in texture coordinates (0,0 -> 1,1)
	p.udata.paddings = math32.Vector4{
		float32(padding.X) / float32(p.width),
		float32(padding.Y) / float32(p.height),
		float32(padding.Width) / float32(p.width),
		float32(padding.Height) / float32(p.height),
	}
	// Update content uniform in texture coordinates (0,0 -> 1,1)
	p.udata.content = math32.Vector4{
		float32(p.content.X) / float32(p.width),
		float32(p.content.Y) / float32(p.height),
		float32(p.content.Width) / float32(p.width),
		float32(p.content.Height) / float32(p.height),
	}
	p.SetChanged(true) // TODO necessary?

	// Update layout and dispatch event
	if !dispatch {
		return
	}
	if p.layout != nil {
		p.layout.Recalc(p)
	}
	p.Dispatch(OnResize, nil)
}

// RenderSetup is called by the Engine before drawing the object
func (p *Panel) RenderSetup(gl *gls.GLS, rinfo *core.RenderInfo) {

	// Sets texture valid flag in uniforms
	// depending if the material has texture
	if p.mat.TextureCount() > 0 {
		p.udata.textureValid = 1
	} else {
		p.udata.textureValid = 0
	}

	// Sets model matrix
	var mm math32.Matrix4
	p.SetModelMatrix(gl, &mm)

	// Transfer model matrix uniform
	location := p.uniMatrix.Location(gl)
	gl.UniformMatrix4fv(location, 1, false, &mm[0])

	// Transfer panel parameters combined uniform
	location = p.uniPanel.Location(gl)
	const vec4count = 8
	gl.Uniform4fv(location, vec4count, &p.udata.bounds.X)
}

// SetModelMatrix calculates and sets the specified matrix with the model matrix for this panel
func (p *Panel) SetModelMatrix(gl *gls.GLS, mm *math32.Matrix4) {

	// Get scale of window (for HiDPI support)
	sX, sY := Manager().win.GetScale()

	// Get the current viewport width and height
	_, _, width, height := gl.GetViewport()

	// Compute common factors
	fX := 2 * float32(sX) / float32(width)
	fY := 2 * float32(sY) / float32(height)

	// Calculate the model matrix
	// Convert pixel coordinates to standard OpenGL clip coordinates and scale the quad for the viewport
	mm.Set(
		fX*float32(p.width), 0, 0, fX*p.pospix.X-1,
		0, fY*float32(p.height), 0, 1-fY*p.pospix.Y,
		0, 0, 1, p.Position().Z,
		0, 0, 0, 1,
	)
}
