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
	"math"
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
	SetRoot(*Root)
	LostKeyFocus()
	TotalHeight() float32
}

// Panel is 2D rectangular graphic which by default has a quad (2 triangles) geometry.
// When using the default geometry, a panel has margins, borders, paddings
// and a content area. The content area can be associated wit a texture
// It is the building block of most GUI widgets.
type Panel struct {
	*graphic.Graphic                     // Embedded graphic
	root             *Root               // pointer to root container
	width            float32             // external width in pixels
	height           float32             // external height in pixels
	mat              *material.Material  // panel material
	marginSizes      BorderSizes         // external margin sizes in pixel coordinates
	borderSizes      BorderSizes         // border sizes in pixel coordinates
	paddingSizes     BorderSizes         // padding sizes in pixel coordinates
	content          Rect                // current content rectangle in pixel coordinates
	modelMatrixUni   gls.UniformMatrix4f // pointer to model matrix uniform
	borderColorUni   gls.Uniform4f       // pointer to border color uniform
	paddingColorUni  gls.Uniform4f       // pointer to padding color uniform
	contentColorUni  gls.Uniform4f       // pointer to content color uniform
	boundsUni        gls.Uniform4f       // pointer to bounds uniform (texture coordinates)
	borderUni        gls.Uniform4f       // pointer to border uniform (texture coordinates)
	paddingUni       gls.Uniform4f       // pointer to padding uniform (texture coordinates)
	contentUni       gls.Uniform4f       // pointer to content uniform (texture coordinates)
	pospix           math32.Vector3      // absolute position in pixels
	posclip          math32.Vector3      // position in clip (NDC) coordinates
	wclip            float32             // width in clip coordinates
	hclip            float32             // height in clip coordinates
	xmin             float32             // minimum absolute x this panel can use
	xmax             float32             // maximum absolute x this panel can use
	ymin             float32             // minimum absolute y this panel can use
	ymax             float32             // maximum absolute y this panel can use
	bounded          bool                // panel is bounded by its parent
	enabled          bool                // enable event processing
	cursorEnter      bool                // mouse enter dispatched
	layout           ILayout             // current layout for children
	layoutParams     interface{}         // current layout parameters used by container panel
}

const (
	deltaZ    = -0.00001      // delta Z for bounded panels
	deltaZunb = deltaZ * 1000 // delta Z for unbounded panels
)

// NewPanel creates and returns a pointer to a new panel with the
// specified dimensions in pixels and a default quad geometry
func NewPanel(width, height float32) *Panel {

	p := new(Panel)
	p.Initialize(width, height)
	return p
}

// Initialize initializes this panel and is normally used by other types which embed a panel.
func (p *Panel) Initialize(width, height float32) {

	p.width = width
	p.height = height

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
	geom.AddVBO(gls.NewVBO().
		AddAttrib("VertexPosition", 3).
		AddAttrib("VertexTexcoord", 2).
		SetBuffer(positions),
	)

	// Initialize material
	p.mat = material.NewMaterial()
	p.mat.SetShader("shaderPanel")

	// Initialize graphic
	p.Graphic = graphic.NewGraphic(geom, gls.TRIANGLES)
	p.AddMaterial(p, p.mat, 0, 0)

	// Creates and adds uniform
	p.modelMatrixUni.Init("ModelMatrix")
	p.borderColorUni.Init("BorderColor")
	p.paddingColorUni.Init("PaddingColor")
	p.contentColorUni.Init("ContentColor")
	p.boundsUni.Init("Bounds")
	p.borderUni.Init("Border")
	p.paddingUni.Init("Padding")
	p.contentUni.Init("Content")

	// Set defaults
	p.borderColorUni.Set(0, 0, 0, 1)
	p.bounded = true
	p.enabled = true
	p.resize(width, height)
}

// InitializeGraphic initializes this panel with a different graphic
func (p *Panel) InitializeGraphic(width, height float32, gr *graphic.Graphic) {

	p.Graphic = gr
	p.width = width
	p.height = height

	// Creates and adds uniform
	p.modelMatrixUni.Init("ModelMatrix")
	p.borderColorUni.Init("BorderColor")
	p.paddingColorUni.Init("PaddingColor")
	p.contentColorUni.Init("ContentColor")
	p.boundsUni.Init("Bounds")
	p.borderUni.Init("Border")
	p.paddingUni.Init("Padding")
	p.contentUni.Init("Content")

	// Set defaults
	p.borderColorUni.Set(0, 0, 0, 1)
	p.bounded = true
	p.enabled = true
	p.resize(width, height)
}

// GetPanel satisfies the IPanel interface and
// returns pointer to this panel
func (pan *Panel) GetPanel() *Panel {

	return pan
}

// SetRoot satisfies the IPanel interface
// Sets the pointer to the root panel for this panel and all its children
func (p *Panel) SetRoot(root *Root) {

	p.root = root
	for i := 0; i < len(p.Children()); i++ {
		cpan := p.Children()[i].(IPanel).GetPanel()
		cpan.SetRoot(root)
	}
}

// LostKeyFocus satisfies the IPanel interface and is called by gui root
// container when the panel loses the key focus
func (p *Panel) LostKeyFocus() {

}

// TotalHeight satisfies the IPanel interface and returns the total
// height of this panel considering visible not bounded children
func (p *Panel) TotalHeight() float32 {

	return p.Height()
}

// SetSelected satisfies the IPanel interface and is normally called
// by a list container to change the panel visual appearance
func (p *Panel) SetSelected2(state bool) {

}

// SetHighlighted satisfies the IPanel interface and is normally called
// by a list container to change the panel visual appearance
func (p *Panel) SetHighlighted2(state bool) {

}

// Material returns a pointer for this panel core.Material
func (p *Panel) Material() *material.Material {

	return p.mat
}

// Root returns the pointer for this panel root panel
func (p *Panel) Root() *Root {

	return p.root
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
	}
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
	p.resize(width, height)
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
	p.resize(p.calcWidth(), p.calcHeight())
}

// SetMarginsFrom sets this panel margins sizes from the specified
// BorderSizes pointer and recalculates the panel external size
func (p *Panel) SetMarginsFrom(src *BorderSizes) {

	p.marginSizes = *src
	p.resize(p.calcWidth(), p.calcHeight())
}

// Margins returns the current margin sizes in pixels
func (p *Panel) Margins() BorderSizes {

	return p.marginSizes
}

// SetBorders sets this panel border sizes in pixels
// and recalculates the panel external size
func (p *Panel) SetBorders(top, right, bottom, left float32) {

	p.borderSizes.Set(top, right, bottom, left)
	p.resize(p.calcWidth(), p.calcHeight())
}

// SetBordersFrom sets this panel border sizes from the specified
// BorderSizes pointer and recalculates the panel size
func (p *Panel) SetBordersFrom(src *BorderSizes) {

	p.borderSizes = *src
	p.resize(p.calcWidth(), p.calcHeight())
}

// Borders returns this panel current border sizes
func (p *Panel) Borders() BorderSizes {

	return p.borderSizes
}

// SetPaddings sets the panel padding sizes in pixels
func (p *Panel) SetPaddings(top, right, bottom, left float32) {

	p.paddingSizes.Set(top, right, bottom, left)
	p.resize(p.calcWidth(), p.calcHeight())
}

// SetPaddingsFrom sets this panel padding sizes from the specified
// BorderSizes pointer and recalculates the panel size
func (p *Panel) SetPaddingsFrom(src *BorderSizes) {

	p.paddingSizes = *src
	p.resize(p.calcWidth(), p.calcHeight())
}

// Paddings returns this panel padding sizes in pixels
func (p *Panel) Paddings() BorderSizes {

	return p.paddingSizes
}

// SetBordersColor sets the color of this panel borders
// The borders opacity is set to 1.0 (full opaque)
func (p *Panel) SetBordersColor(color *math32.Color) {

	p.borderColorUni.Set(color.R, color.G, color.B, 1)
}

// SetBordersColor4 sets the color and opacity of this panel borders
func (p *Panel) SetBordersColor4(color *math32.Color4) {

	p.borderColorUni.Set(color.R, color.G, color.B, color.A)
}

// BorderColor4 returns current border color
func (p *Panel) BordersColor4() math32.Color4 {

	var color math32.Color4
	p.borderColorUni.SetColor4(&color)
	return color
}

// SetPaddingsColor sets the color of this panel paddings.
func (p *Panel) SetPaddingsColor(color *math32.Color) {

	p.paddingColorUni.Set(color.R, color.G, color.B, 1.0)
}

// SetColor sets the color of the panel paddings and content area
func (p *Panel) SetColor(color *math32.Color) *Panel {

	p.paddingColorUni.Set(color.R, color.G, color.B, 1.0)
	p.contentColorUni.Set(color.R, color.G, color.B, 1.0)
	return p
}

// SetColor4 sets the color of the panel paddings and content area
func (p *Panel) SetColor4(color *math32.Color4) *Panel {

	p.paddingColorUni.Set(color.R, color.G, color.B, color.A)
	p.contentColorUni.Set(color.R, color.G, color.B, color.A)
	return p
}

// Color4 returns the current color of the panel content area
func (p *Panel) Color4() math32.Color4 {

	return p.contentColorUni.GetColor4()
}

// SetContentSize sets this panel content size to the specified dimensions.
// The external size of the panel may increase or decrease to acomodate
// the new content size.
func (p *Panel) SetContentSize(width, height float32) {

	// Calculates the new desired external width and height
	eWidth := width +
		p.paddingSizes.Left + p.paddingSizes.Right +
		p.borderSizes.Left + p.borderSizes.Right +
		p.marginSizes.Left + p.marginSizes.Right
	eHeight := height +
		p.paddingSizes.Top + p.paddingSizes.Bottom +
		p.borderSizes.Top + p.borderSizes.Bottom +
		p.marginSizes.Top + p.marginSizes.Bottom
	p.resize(eWidth, eHeight)
}

// SetContentWidth sets this panel content width to the specified dimension in pixels.
// The external size of the panel may increase or decrease to acomodate the new width
func (p *Panel) SetContentWidth(width float32) {

	p.SetContentSize(width, p.content.Height)
}

// SetContentHeight sets this panel content height to the specified dimension in pixels.
// The external size of the panel may increase or decrease to acomodate the new width
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

// Add adds a child panel to this one
func (p *Panel) Add(ichild IPanel) *Panel {

	p.Node.Add(ichild)
	node := ichild.GetPanel()
	node.SetParent(p)
	if p.root != nil {
		ichild.SetRoot(p.root)
		p.root.setZ(0, deltaZunb)
	}
	if p.layout != nil {
		p.layout.Recalc(p)
	}
	p.Dispatch(OnChild, nil)
	return p
}

// Remove removes the specified child from this panel
func (p *Panel) Remove(ichild IPanel) bool {

	res := p.Node.Remove(ichild)
	if res {
		if p.layout != nil {
			p.layout.Recalc(p)
		}
		p.Dispatch(OnChild, nil)
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
}

// UpdateMatrixWorld overrides the standard core.Node version which is called by
// the Engine before rendering the frame.
func (p *Panel) UpdateMatrixWorld() {

	// Panel has no parent should be the root panel
	par := p.Parent()
	if par == nil {
		p.updateBounds(nil)
		// Panel has parent
	} else {
		parpan := par.(*Panel)
		p.updateBounds(parpan)
	}
	// Update this panel children
	for _, ichild := range p.Children() {
		ichild.UpdateMatrixWorld()
	}
}

// ContainsPosition returns indication if this panel contains
// the specified screen position in pixels.
func (p *Panel) ContainsPosition(x, y float32) bool {

	if x < p.pospix.X || x >= (p.pospix.X+p.width) {
		return false
	}
	if y < p.pospix.Y || y >= (p.pospix.Y+p.height) {
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

// SetLayoutParams sets the layout parameters for this panel
func (p *Panel) SetLayoutParams(params interface{}) {

	p.layoutParams = params
}

// NDC2Pix converts the specified NDC coordinates (-1,1) to relative pixel coordinates
// for this panel content area.
// 0,0      1,0        0,0       w,0
// +--------+          +---------+
// |        | -------> |         |
// +--------+          +---------+
// 0,-1     1,-1       0,h       w,h
func (p *Panel) NDC2Pix(nx, ny float32) (x, y float32) {

	w := p.ContentWidth()
	h := p.ContentHeight()
	return w * nx, -h * ny
}

// Pix2NDC converts the specified relative pixel coordinates to NDC coordinates for this panel
// content area
// 0,0       w,0       0,0      1,0
// +---------+         +---------+
// |         | ------> |         |
// +---------+         +---------+
// 0,h       w,h       0,-1     1,-1
func (p *Panel) Pix2NDC(px, py float32) (nx, ny float32) {

	w := p.ContentWidth()
	h := p.ContentHeight()
	return px / w, -py / h
}

// setZ sets the Z coordinate for this panel and its children recursively
// starting at the specified z and zunb coordinates.
// The z coordinate is used for bound panels and zunb for unbounded panels.
// The z coordinate is set so panels added later are closer to the screen.
// All unbounded panels and its children are closer than any of the bounded panels.
func (p *Panel) setZ(z, zunb float32) (float32, float32) {

	// Bounded panel
	if p.bounded {
		p.SetPositionZ(z)
		z += deltaZ
		for _, ichild := range p.Children() {
			z, zunb = ichild.(IPanel).GetPanel().setZ(z, zunb)
		}
		return z, zunb
		// Unbounded panel
	} else {
		p.SetPositionZ(zunb)
		zchild := zunb + deltaZ
		zunb += deltaZunb
		for _, ichild := range p.Children() {
			_, zunb = ichild.(IPanel).GetPanel().setZ(zchild, zunb)
		}
		return z, zunb
	}
}

// updateBounds is called by UpdateMatrixWorld() and calculates this panel
// bounds considering the bounds of its parent
func (p *Panel) updateBounds(par *Panel) {

	// If no parent, it is the root panel
	if par == nil {
		p.pospix = p.Position()
		p.xmin = -math.MaxFloat32
		p.ymin = -math.MaxFloat32
		p.xmax = math.MaxFloat32
		p.ymax = math.MaxFloat32
		p.boundsUni.Set(0, 0, 1, 1)
		return
	}
	// If this panel is bounded to its parent, its coordinates are relative
	// to the parent internal content rectangle.
	if p.bounded {
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
	if p.bounded {
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
	}
	// Set default values for bounds in texture coordinates
	xmintex := float32(0.0)
	ymintex := float32(0.0)
	xmaxtex := float32(1.0)
	ymaxtex := float32(1.0)
	// If this panel is bounded to its parent, calculates the bounds
	// for clipping in texture coordinates
	if p.bounded {
		if p.pospix.X < p.xmin {
			xmintex = (p.xmin - p.pospix.X) / p.width
		}
		if p.pospix.Y < p.ymin {
			ymintex = (p.ymin - p.pospix.Y) / p.height
		}
		if p.pospix.X+p.width > p.xmax {
			xmaxtex = (p.xmax - p.pospix.X) / p.width
		}
		if p.pospix.Y+p.height > p.ymax {
			ymaxtex = (p.ymax - p.pospix.Y) / p.height
		}
	}
	// Sets bounds uniform
	p.boundsUni.Set(xmintex, ymintex, xmaxtex, ymaxtex)
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
func (p *Panel) resize(width, height float32) {

	var padding Rect
	var border Rect

	width = math32.Round(width)
	height = math32.Round(height)
	// Adjusts content width
	p.content.Width = width -
		p.marginSizes.Left - p.marginSizes.Right -
		p.borderSizes.Left - p.borderSizes.Right -
		p.paddingSizes.Left - p.paddingSizes.Right
	if p.content.Width < 0 {
		p.content.Width = 0
	}
	// Adjust other area widths
	padding.Width = p.paddingSizes.Left + p.content.Width + p.paddingSizes.Right
	border.Width = p.borderSizes.Left + padding.Width + p.borderSizes.Right

	// Adjusts content height
	p.content.Height = height -
		p.marginSizes.Top - p.marginSizes.Bottom -
		p.borderSizes.Top - p.borderSizes.Bottom -
		p.paddingSizes.Top - p.paddingSizes.Bottom
	if p.content.Height < 0 {
		p.content.Height = 0
	}
	// Adjust other area heights
	padding.Height = p.paddingSizes.Top + p.content.Height + p.paddingSizes.Bottom
	border.Height = p.borderSizes.Top + padding.Height + p.borderSizes.Bottom

	// Sets area positions
	border.X = p.marginSizes.Left
	border.Y = p.marginSizes.Top
	padding.X = border.X + p.borderSizes.Left
	padding.Y = border.Y + p.borderSizes.Top
	p.content.X = padding.X + p.paddingSizes.Left
	p.content.Y = padding.Y + p.paddingSizes.Top

	// Sets final panel dimensions (may be different from requested dimensions)
	p.width = p.marginSizes.Left + border.Width + p.marginSizes.Right
	p.height = p.marginSizes.Top + border.Height + p.marginSizes.Bottom

	// Updates border uniform in texture coordinates (0,0 -> 1,1)
	p.borderUni.Set(
		float32(border.X)/float32(p.width),
		float32(border.Y)/float32(p.height),
		float32(border.Width)/float32(p.width),
		float32(border.Height)/float32(p.height),
	)
	// Updates padding uniform in texture coordinates (0,0 -> 1,1)
	p.paddingUni.Set(
		float32(padding.X)/float32(p.width),
		float32(padding.Y)/float32(p.height),
		float32(padding.Width)/float32(p.width),
		float32(padding.Height)/float32(p.height),
	)
	// Updates content uniform in texture coordinates (0,0 -> 1,1)
	p.contentUni.Set(
		float32(p.content.X)/float32(p.width),
		float32(p.content.Y)/float32(p.height),
		float32(p.content.Width)/float32(p.width),
		float32(p.content.Height)/float32(p.height),
	)
	// Update layout and dispatch event
	if p.layout != nil {
		p.layout.Recalc(p)
	}
	p.Dispatch(OnResize, nil)
}

// RenderSetup is called by the Engine before drawing the object
func (p *Panel) RenderSetup(gl *gls.GLS, rinfo *core.RenderInfo) {

	// Sets model matrix
	var mm math32.Matrix4
	p.SetModelMatrix(gl, &mm)
	p.modelMatrixUni.SetMatrix4(&mm)

	// Transfer uniforms
	p.borderColorUni.Transfer(gl)
	p.paddingColorUni.Transfer(gl)
	p.contentColorUni.Transfer(gl)
	p.boundsUni.Transfer(gl)
	p.borderUni.Transfer(gl)
	p.paddingUni.Transfer(gl)
	p.contentUni.Transfer(gl)
	p.modelMatrixUni.Transfer(gl)
	//log.Error("panel:%p boundsUni.:%v", p, p.boundsUni)
	//log.Error("panel:%p borderUni.:%v", p, p.borderUni)
	//log.Error("panel:%p paddingUni:%v", p, p.paddingUni)
	//log.Error("panel:%p contentUni:%v", p, p.contentUni)
	//log.Error("panel:%p mm        :%v", p, p.modelMatrixUni)
}

// SetModelMatrix calculates and sets the specified matrix with the model matrix for this panel
func (p *Panel) SetModelMatrix(gl *gls.GLS, mm *math32.Matrix4) {

	// Get the current viewport width and height
	_, _, width, height := gl.GetViewport()
	fwidth := float32(width)
	fheight := float32(height)

	// Scale the quad for the viewport so it has fixed dimensions in pixels.
	p.wclip = 2 * float32(p.width) / fwidth
	p.hclip = 2 * float32(p.height) / fheight
	var scale math32.Vector3
	scale.Set(p.wclip, p.hclip, 1)

	// Convert absolute position in pixel coordinates from the top/left to
	// standard OpenGL clip coordinates of the quad center
	p.posclip.X = (p.pospix.X - fwidth/2) / (fwidth / 2)
	p.posclip.Y = -(p.pospix.Y - fheight/2) / (fheight / 2)
	p.posclip.Z = p.Position().Z

	// Calculates the model matrix
	var quat math32.Quaternion
	quat.SetIdentity()
	mm.Compose(&p.posclip, &quat, &scale)
}
