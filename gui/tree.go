// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/window"
)

// Tree is the tree structure GUI element.
type Tree struct {
	List               // Embedded list panel
	styles *TreeStyles // Pointer to styles
}

// TreeStyles contains the styling of all tree components for each valid GUI state.
type TreeStyles struct {
	List     *ListStyles     // Styles for the embedded list
	Node     *TreeNodeStyles // Styles for the node panel
	Padlevel float32         // Left padding indentation
}

// TreeNodeStyles contains a TreeNodeStyle for each valid GUI state.
type TreeNodeStyles struct {
	Normal TreeNodeStyle
}

// TreeNodeStyle contains the styling of a TreeNode.
type TreeNodeStyle struct {
	PanelStyle
	FgColor math32.Color4
	Icons   [2]string
}

// TreeNode is a tree node.
type TreeNode struct {
	Panel              // Embedded panel
	label    Label     // Node label
	icon     Label     // Node icon
	tree     *Tree     // Parent tree
	parNode  *TreeNode // Parent node
	items    []IPanel  // List of node items
	expanded bool      // Node expanded flag
}

// NewTree creates and returns a pointer to a new tree widget.
func NewTree(width, height float32) *Tree {

	t := new(Tree)
	t.Initialize(width, height)
	return t
}

// Initialize initializes the tree with the specified initial width and height
// It is normally used when the folder is embedded in another object.
func (t *Tree) Initialize(width, height float32) {

	t.List.initialize(true, width, height)
	t.SetStyles(&StyleDefault().Tree)
	t.List.Subscribe(OnKeyDown, t.onKey)
	t.List.Subscribe(OnKeyUp, t.onKey)
	t.List.Subscribe(OnCursor, t.onCursor)
}

// SetStyles sets the tree styles overriding the default style.
func (t *Tree) SetStyles(s *TreeStyles) {

	t.styles = s
	t.List.SetStyles(t.styles.List)
	t.update()
}

// InsertAt inserts a child panel at the specified position in the tree.
func (t *Tree) InsertAt(pos int, child IPanel) {

	t.List.InsertAt(pos, child)
}

// Add child panel to the end tree.
func (t *Tree) Add(ichild IPanel) {

	t.List.Add(ichild)
}

// InsertNodeAt inserts at the specified position a new tree node
// with the specified text at the end of this tree
// and returns pointer to the new node.
func (t *Tree) InsertNodeAt(pos int, text string) *TreeNode {

	n := newTreeNode(text, t, nil)
	n.update()
	n.recalc()
	t.List.InsertAt(pos, n)
	return n
}

// AddNode adds a new tree node with the specified text
// at the end of this tree and returns a pointer to the new node.
func (t *Tree) AddNode(text string) *TreeNode {

	n := newTreeNode(text, t, nil)
	n.update()
	n.recalc()
	t.List.Add(n)
	return n
}

// Remove removes the specified child from the tree or any
// of its children nodes.
func (t *Tree) Remove(child IPanel) {

	for idx := 0; idx < t.List.Len(); idx++ {
		curr := t.List.ItemAt(idx)
		if curr == child {
			node, ok := curr.(*TreeNode)
			if ok {
				node.remove()
			} else {
				t.List.Remove(child)
			}
			return
		}
		node, ok := curr.(*TreeNode)
		if ok {
			node.Remove(child)
		}
	}
}

// Selected returns the currently selected element or nil
func (t *Tree) Selected() IPanel {

	sel := t.List.Selected()
	if len(sel) == 0 {
		return nil
	}
	return sel[0]
}

// FindChild searches for the specified child in the tree and
// all its children. If found, returns the parent node and
// its position relative to the parent.
// If the parent is the tree returns nil as the parent
// If not found returns nil and -1
func (t *Tree) FindChild(child IPanel) (*TreeNode, int) {

	for idx := 0; idx < t.List.Len(); idx++ {
		curr := t.List.ItemAt(idx)
		if curr == child {
			return nil, idx
		}
		node, ok := curr.(*TreeNode)
		if ok {
			par, pos := node.FindChild(child)
			if pos >= 0 {
				return par, pos
			}
		}
	}
	return nil, -1
}

// onCursor receives subscribed cursor events over the tree
func (t *Tree) onCursor(evname string, ev interface{}) {

	// Do not propagate any cursor events
	t.root.StopPropagation(StopAll)
}

// onKey receives key down events for the embedded list
func (t *Tree) onKey(evname string, ev interface{}) {

	// Get selected item
	item := t.Selected()
	if item == nil {
		return
	}
	// If item is not a tree node, dispatch event to item
	node, ok := item.(*TreeNode)
	if !ok {
		item.SetRoot(t.root)
		item.GetPanel().Dispatch(evname, ev)
		return
	}
	// If not enter key pressed, ignore
	kev := ev.(*window.KeyEvent)
	if evname != OnKeyDown || kev.Keycode != window.KeyEnter {
		return
	}
	// Toggles the expansion state of the node
	node.expanded = !node.expanded
	node.update()
	node.updateItems()
}

//
// TreeNode methods
//

// newTreeNode creates and returns a pointer to a new TreeNode with
// the specified text, tree and parent node
func newTreeNode(text string, tree *Tree, parNode *TreeNode) *TreeNode {

	n := new(TreeNode)
	n.Panel.Initialize(0, 0)

	// Initialize node label
	n.label.initialize(text, StyleDefault().Font)
	n.Panel.Add(&n.label)

	// Create node icon
	n.icon.initialize("", StyleDefault().FontIcon)
	n.icon.SetFontSize(StyleDefault().Label.PointSize * 1.3)
	n.Panel.Add(&n.icon)

	// Subscribe to events
	n.Panel.Subscribe(OnMouseDown, n.onMouse)
	n.Panel.Subscribe(OnListItemResize, func(evname string, ev interface{}) {
		n.recalc()
	})
	n.tree = tree
	n.parNode = parNode

	n.update()
	n.recalc()
	return n
}

// Len returns the number of immediate children of this node
func (n *TreeNode) Len() int {

	return len(n.items)
}

// SetExpanded sets the expanded state of this node
func (n *TreeNode) SetExpanded(state bool) {

	n.expanded = state
	n.update()
	n.updateItems()
}

// FindChild searches for the specified child in this node and
// all its children. If found, returns the parent node and
// its position relative to the parent.
// If not found returns nil and -1
func (n *TreeNode) FindChild(child IPanel) (*TreeNode, int) {

	for pos, curr := range n.items {
		if curr == child {
			return n, pos
		}
		node, ok := curr.(*TreeNode)
		if ok {
			par, pos := node.FindChild(child)
			if par != nil {
				return par, pos
			}
		}
	}
	return nil, -1
}

// InsertAt inserts a child panel at the specified position in this node
// If the position is invalid, the function panics
func (n *TreeNode) InsertAt(pos int, child IPanel) {

	if pos < 0 || pos > len(n.items) {
		panic("TreeNode.InsertAt(): Invalid position")
	}
	// Insert item in the items array
	n.items = append(n.items, nil)
	copy(n.items[pos+1:], n.items[pos:])
	n.items[pos] = child
	if n.expanded {
		n.updateItems()
	}
}

// Add adds a child panel to this node
func (n *TreeNode) Add(child IPanel) {

	n.InsertAt(n.Len(), child)
}

// InsertNodeAt inserts a new node at the specified position in this node
// If the position is invalid, the function panics
func (n *TreeNode) InsertNodeAt(pos int, text string) *TreeNode {

	if pos < 0 || pos > len(n.items) {
		panic("TreeNode.InsertNodeAt(): Invalid position")
	}
	childNode := newTreeNode(text, n.tree, n)
	// Insert item in the items array
	n.items = append(n.items, nil)
	copy(n.items[pos+1:], n.items[pos:])
	n.items[pos] = childNode
	if n.expanded {
		n.updateItems()
	}
	return childNode
}

// AddNode adds a new node to this one and return its pointer
func (n *TreeNode) AddNode(text string) *TreeNode {

	return n.InsertNodeAt(n.Len(), text)
}

// Remove removes the specified child from this node or any
// of its children nodes
func (n *TreeNode) Remove(child IPanel) {

	for pos, curr := range n.items {
		if curr == child {
			copy(n.items[pos:], n.items[pos+1:])
			n.items[len(n.items)-1] = nil
			n.items = n.items[:len(n.items)-1]
			node, ok := curr.(*TreeNode)
			if ok {
				node.remove()
			} else {
				n.tree.List.Remove(curr)
			}
			n.updateItems()
			return
		}
		node, ok := curr.(*TreeNode)
		if ok {
			node.Remove(child)
		}
	}
}

// onMouse receives mouse button events over the tree node panel
func (n *TreeNode) onMouse(evname string, ev interface{}) {

	switch evname {
	case OnMouseDown:
		n.expanded = !n.expanded
		n.update()
		n.recalc()
		n.updateItems()
	default:
		return
	}
}

// level returns the level of this node from the start of the tree
func (n *TreeNode) level() int {

	level := 0
	parNode := n.parNode
	for parNode != nil {
		parNode = parNode.parNode
		level++
	}
	return level
}

// applyStyles applies the specified style to this tree node
func (n *TreeNode) applyStyle(s *TreeNodeStyle) {

	n.Panel.ApplyStyle(&s.PanelStyle)
	icode := 0
	if n.expanded {
		icode = 1
	}
	n.icon.SetText(string(s.Icons[icode]))
	n.icon.SetColor4(&s.FgColor)
	n.label.SetColor4(&s.FgColor)
}

// update updates this tree node style
func (n *TreeNode) update() {

	n.applyStyle(&n.tree.styles.Node.Normal)
}

// recalc recalculates the positions of the internal node panels
func (n *TreeNode) recalc() {

	// icon position
	n.icon.SetPosition(0, 0)

	// Label position and width
	n.label.SetPosition(n.icon.Width()+4, 0)
	n.Panel.SetContentHeight(n.label.Height())
	n.Panel.SetWidth(n.tree.ContentWidth())
}

// remove removes this node and all children from the tree list
func (n *TreeNode) remove() {

	n.tree.List.Remove(n)
	n.removeItems()
}

// removeItems removes this node children from the tree list
func (n *TreeNode) removeItems() {

	for _, ipanel := range n.items {
		// Remove item from scroller
		n.tree.List.Remove(ipanel)
		// If item is a node, remove all children
		node, ok := ipanel.(*TreeNode)
		if ok {
			node.removeItems()
			continue
		}
	}
}

// insert inserts this node and its expanded children in the tree list
// at the specified position
func (n *TreeNode) insert(pos int) int {

	n.update()
	n.tree.List.InsertAt(pos, n)
	var padLeft float32 = n.tree.styles.Padlevel * float32(n.level())
	n.tree.List.SetItemPadLeftAt(pos, padLeft)
	pos++
	return n.insertItems(pos)
}

// insertItems inserts this node items in the tree list
// at the specified position
func (n *TreeNode) insertItems(pos int) int {

	if !n.expanded {
		return pos
	}
	level := n.level() + 1
	var padLeft float32 = n.tree.styles.Padlevel * float32(level)
	for _, ipanel := range n.items {
		// Insert node and its children
		node, ok := ipanel.(*TreeNode)
		if ok {
			node.update()
			n.tree.List.InsertAt(pos, ipanel)
			n.tree.List.SetItemPadLeftAt(pos, padLeft)
			pos++
			pos = node.insertItems(pos)
			continue
		}
		// Insert item
		n.tree.List.InsertAt(pos, ipanel)
		n.tree.List.SetItemPadLeftAt(pos, padLeft)
		pos++
	}
	return pos
}

// updateItems updates this node items, removing or inserting them into the tree scroller
func (n *TreeNode) updateItems() {

	pos := n.tree.ItemPosition(n)
	if pos < 0 {
		return
	}
	n.removeItems()
	n.insertItems(pos + 1)
}
