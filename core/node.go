// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/math32"
	"strings"
)

// INode is the interface for all node types.
type INode interface {
	IDispatcher
	GetNode() *Node
	UpdateMatrixWorld()
	Raycast(*Raycaster, *[]Intersect)
	Render(gs *gls.GLS)
	Dispose()
}

// Node represents an object in 3D space existing within a hierarchy.
type Node struct {
	Dispatcher             // Embedded event dispatcher
	parent     INode       // Parent node
	children   []INode     // Children nodes
	name       string      // Optional node name
	loaderID   string      // ID used by loader
	visible    bool        // Whether the node is visible
	changed    bool        // Whether the position/orientation/scale changed
	userData   interface{} // Generic user data

	// Spatial properties
	position    math32.Vector3    // Node position in 3D space (relative to parent)
	scale       math32.Vector3    // Node scale (relative to parent)
	direction   math32.Vector3    // Initial direction (relative to parent)
	rotation    math32.Vector3    // Node rotation specified in Euler angles (relative to parent)
	quaternion  math32.Quaternion // Node rotation specified as a Quaternion (relative to parent)
	matrix      math32.Matrix4    // Local transform matrix. Contains all position/rotation/scale information (relative to parent)
	matrixWorld math32.Matrix4    // World transform matrix. Contains all absolute position/rotation/scale information (i.e. relative to very top parent, generally the scene)
}

// NewNode returns a pointer to a new Node.
func NewNode() *Node {

	n := new(Node)
	n.Init()
	return n
}

// Init initializes the node.
// Normally called by other types which embed a Node.
func (n *Node) Init() {

	n.Dispatcher.Initialize()

	n.children = make([]INode, 0)
	n.visible = true
	n.changed = true

	// Initialize spatial properties
	n.position.Set(0, 0, 0)
	n.scale.Set(1, 1, 1)
	n.direction.Set(0, 0, 1)
	n.rotation.Set(0, 0, 0)
	n.quaternion.Set(0, 0, 0, 1)
	n.matrix.Identity()
	n.matrixWorld.Identity()
}

// GetNode satisfies the INode interface
// and returns a pointer to the embedded Node.
func (n *Node) GetNode() *Node {

	return n
}

// Raycast satisfies the INode interface.
func (n *Node) Raycast(rc *Raycaster, intersects *[]Intersect) {
}

// Render satisfies the INode interface.
func (n *Node) Render(gs *gls.GLS) {
}

// Dispose satisfies the INode interface.
func (n *Node) Dispose() {
}

// SetParent sets the parent.
func (n *Node) SetParent(iparent INode) {

	n.parent = iparent
}

// Parent returns the parent.
func (n *Node) Parent() INode {

	return n.parent
}

// SetName sets the (optional) name.
// The name can be used for debugging or other purposes.
func (n *Node) SetName(name string) {

	n.name = name
}

// Name returns the (optional) name.
func (n *Node) Name() string {

	return n.name
}

// SetLoaderID is normally used by external loaders, such as Collada,
// to assign an ID to the node with the ID value in the node description.
// Can be used to find other loaded nodes.
func (n *Node) SetLoaderID(id string) {

	n.loaderID = id
}

// LoaderID returns an optional ID set when this node was
// created by an external loader such as Collada.
func (n *Node) LoaderID() string {

	return n.loaderID
}

// SetVisible sets the visibility of the node.
func (n *Node) SetVisible(state bool) {

	n.visible = state
	n.changed = true
}

// Visible returns the visibility of the node.
func (n *Node) Visible() bool {

	return n.visible
}

// SetChanged sets the changed flag of the node.
func (n *Node) SetChanged(changed bool) {

	n.changed = changed
}

// Changed returns the changed flag of the node.
func (n *Node) Changed() bool {

	return n.changed
}

// SetUserData sets the generic user data associated to the node.
func (n *Node) SetUserData(data interface{}) {

	n.userData = data
}

// UserData returns the generic user data associated to the node.
func (n *Node) UserData() interface{} {

	return n.userData
}

// FindPath finds a node with the specified path starting with this node and
// searching in all its children recursively.
// A path is the sequence of the names from the first node to the desired node
// separated by the forward slash.
func (n *Node) FindPath(path string) INode {

	// Internal recursive function to find node
	var finder func(inode INode, path string) INode
	finder = func(inode INode, path string) INode {
		// Get first component of the path
		parts := strings.Split(path, "/")
		if len(parts) == 0 {
			return nil
		}
		first := parts[0]
		// Checks current node
		node := inode.GetNode()
		if node.name != first {
			return nil
		}
		// If the path has finished this is the desired node
		rest := strings.Join(parts[1:], "/")
		if rest == "" {
			return inode
		}
		// Otherwise search in this node children
		for _, ichild := range node.children {
			found := finder(ichild, rest)
			if found != nil {
				return found
			}
		}
		return nil
	}
	return finder(n, path)
}

// FindLoaderID looks in the specified node and in all its children
// for a node with the specified loaderID and if found returns it.
// Returns nil if not found.
func (n *Node) FindLoaderID(id string) INode {

	var finder func(parent INode, id string) INode
	finder = func(parent INode, id string) INode {
		pnode := parent.GetNode()
		if pnode.loaderID == id {
			return parent
		}
		for _, child := range pnode.children {
			found := finder(child, id)
			if found != nil {
				return found
			}
		}
		return nil
	}
	return finder(n, id)
}

// Children returns the list of children.
func (n *Node) Children() []INode {

	return n.children
}

// Add adds the specified node to the list of children and sets its parent pointer.
// If the specified node had a parent, the specified node is removed from the original parent's list of children.
func (n *Node) Add(ichild INode) *Node {

	n.setParentOf(ichild)
	n.children = append(n.children, ichild)
	return n
}

// AddAt adds the specified node to the list of children at the specified index and sets its parent pointer.
// If the specified node had a parent, the specified node is removed from the original parent's list of children.
func (n *Node) AddAt(idx int, ichild INode) {

	// Validate position
	if idx < 0 || idx > len(n.children) {
		panic("Node.AddAt: invalid position")
	}

	n.setParentOf(ichild)

	// Insert child in the specified position
	n.children = append(n.children, nil)
	copy(n.children[idx+1:], n.children[idx:])
	n.children[idx] = ichild
}

// setParentOf is used by Add and AddAt.
// It verifies that the node is not being added to itself and sets the parent pointer of the specified node.
// If the specified node had a parent, the specified node is removed from the original parent's list of children.
// It does not add the specified node to the list of children.
func (n *Node) setParentOf(ichild INode) {

	child := ichild.GetNode()
	if n == child {
		panic("Node.{Add,AddAt}: object can't be added as a child of itself")
	}
	// If the specified node already has a parent,
	// remove it from the original parent's list of children
	if child.parent != nil {
		child.parent.GetNode().Remove(ichild)
	}
	child.parent = n
}

// ChildAt returns the child at the specified index.
func (n *Node) ChildAt(idx int) INode {

	if idx < 0 || idx >= len(n.children) {
		return nil
	}
	return n.children[idx]
}

// ChildIndex returns the index of the specified child (-1 if not found).
func (n *Node) ChildIndex(ichild INode) int {

	for idx := 0; idx < len(n.children); idx++ {
		if n.children[idx] == ichild {
			return idx
		}
	}
	return -1
}

// Remove removes the specified INode from the list of children.
// Returns true if found or false otherwise.
func (n *Node) Remove(ichild INode) bool {

	for pos, current := range n.children {
		if current == ichild {
			copy(n.children[pos:], n.children[pos+1:])
			n.children[len(n.children)-1] = nil
			n.children = n.children[:len(n.children)-1]
			ichild.GetNode().parent = nil
			return true
		}
	}
	return false
}

// RemoveAt removes the child at the specified index.
func (n *Node) RemoveAt(idx int) INode {

	// Validate position
	if idx < 0 || idx >= len(n.children) {
		panic("Node.RemoveAt: invalid position")
	}

	child := n.children[idx]

	// Remove child from children list
	copy(n.children[idx:], n.children[idx+1:])
	n.children[len(n.children)-1] = nil
	n.children = n.children[:len(n.children)-1]

	return child
}

// RemoveAll removes all children.
func (n *Node) RemoveAll(recurs bool) {

	for pos, ichild := range n.children {
		n.children[pos] = nil
		ichild.GetNode().parent = nil
		if recurs {
			ichild.GetNode().RemoveAll(recurs)
		}
	}
	n.children = n.children[0:0]
}

// DisposeChildren removes and disposes of all children.
// If 'recurs' is true, call DisposeChildren on each child recursively.
func (n *Node) DisposeChildren(recurs bool) {

	for pos, ichild := range n.children {
		n.children[pos] = nil
		ichild.GetNode().parent = nil
		if recurs {
			ichild.GetNode().DisposeChildren(true)
		}
		ichild.Dispose()
	}
	n.children = n.children[0:0]
}

// SetPosition sets the position.
func (n *Node) SetPosition(x, y, z float32) {

	n.position.Set(x, y, z)
	n.changed = true
}

// SetPositionVec sets the position based on the specified vector pointer.
func (n *Node) SetPositionVec(vpos *math32.Vector3) {

	n.position = *vpos
	n.changed = true
}

// SetPositionX sets the X coordinate of the position.
func (n *Node) SetPositionX(x float32) {

	n.position.X = x
	n.changed = true
}

// SetPositionY sets the Y coordinate of the position.
func (n *Node) SetPositionY(y float32) {

	n.position.Y = y
	n.changed = true
}

// SetPositionZ sets the Z coordinate of the position.
func (n *Node) SetPositionZ(z float32) {

	n.position.Z = z
	n.changed = true
}

// Position returns the position as a vector.
func (n *Node) Position() math32.Vector3 {

	return n.position
}

// SetRotation sets the rotation in radians.
// The stored quaternion is updated accordingly.
func (n *Node) SetRotation(x, y, z float32) {

	n.rotation.Set(x, y, z)
	n.quaternion.SetFromEuler(&n.rotation)
	n.changed = true
}

// SetRotationVec sets the rotation in radians based on the specified vector pointer.
// The stored quaternion is updated accordingly.
func (n *Node) SetRotationVec(vrot *math32.Vector3) {

	n.rotation = *vrot
	n.quaternion.SetFromEuler(&n.rotation)
	n.changed = true
}

// SetRotationQuat sets the rotation based on the specified quaternion pointer.
// The stored quaternion is updated accordingly.
func (n *Node) SetRotationQuat(quat *math32.Quaternion) {

	n.quaternion = *quat
	n.changed = true
}

// SetRotationX sets the X rotation to the specified angle in radians.
// The stored quaternion is updated accordingly.
func (n *Node) SetRotationX(x float32) {

	n.rotation.X = x
	n.quaternion.SetFromEuler(&n.rotation)
	n.changed = true
}

// SetRotationY sets the Y rotation to the specified angle in radians.
// The stored quaternion is updated accordingly.
func (n *Node) SetRotationY(y float32) {

	n.rotation.Y = y
	n.quaternion.SetFromEuler(&n.rotation)
	n.changed = true
}

// SetRotationZ sets the Z rotation to the specified angle in radians.
// The stored quaternion is updated accordingly.
func (n *Node) SetRotationZ(z float32) {

	n.rotation.Z = z
	n.quaternion.SetFromEuler(&n.rotation)
	n.changed = true
}

// AddRotationX adds to the current X rotation the specified angle in radians.
// The stored quaternion is updated accordingly.
func (n *Node) AddRotationX(x float32) {

	n.rotation.X += x
	n.quaternion.SetFromEuler(&n.rotation)
	n.changed = true
}

// AddRotationY adds to the current Y rotation the specified angle in radians.
// The stored quaternion is updated accordingly.
func (n *Node) AddRotationY(y float32) {

	n.rotation.Y += y
	n.quaternion.SetFromEuler(&n.rotation)
	n.changed = true
}

// AddRotationZ adds to the current Z rotation the specified angle in radians.
// The stored quaternion is updated accordingly.
func (n *Node) AddRotationZ(z float32) {

	n.rotation.Z += z
	n.quaternion.SetFromEuler(&n.rotation)
	n.changed = true
}

// Rotation returns the current rotation.
func (n *Node) Rotation() math32.Vector3 {

	return n.rotation
}

// SetQuaternion sets the quaternion based on the specified quaternion unit multiples.
func (n *Node) SetQuaternion(x, y, z, w float32) {

	n.quaternion.Set(x, y, z, w)
	n.changed = true
}

// SetQuaternionQuat sets the quaternion based on the specified quaternion pointer.
func (n *Node) SetQuaternionQuat(q *math32.Quaternion) {

	n.quaternion = *q
	n.changed = true
}

// QuaternionMult multiplies the current quaternion by the specified quaternion.
func (n *Node) QuaternionMult(q *math32.Quaternion) {

	n.quaternion.Multiply(q)
	n.changed = true
}

// Quaternion returns the current quaternion.
func (n *Node) Quaternion() math32.Quaternion {

	return n.quaternion
}

// SetScale sets the scale.
func (n *Node) SetScale(x, y, z float32) {

	n.scale.Set(x, y, z)
	n.changed = true
}

// SetScaleVec sets the scale based on the specified vector pointer.
func (n *Node) SetScaleVec(scale *math32.Vector3) {

	n.scale = *scale
	n.changed = true
}

// SetScaleX sets the X scale.
func (n *Node) SetScaleX(sx float32) {

	n.scale.X = sx
	n.changed = true
}

// SetScaleY sets the Y scale.
func (n *Node) SetScaleY(sy float32) {

	n.scale.Y = sy
	n.changed = true
}

// SetScaleZ sets the Z scale.
func (n *Node) SetScaleZ(sz float32) {

	n.scale.Z = sz
	n.changed = true
}

// Scale returns the current scale.
func (n *Node) Scale() math32.Vector3 {

	return n.scale
}

// SetDirection sets the direction.
func (n *Node) SetDirection(x, y, z float32) {

	n.direction.Set(x, y, z)
	n.changed = true
}

// SetDirectionVec sets the direction based on a vector pointer.
func (n *Node) SetDirectionVec(vdir *math32.Vector3) {

	n.direction = *vdir
	n.changed = true
}

// Direction returns the direction.
func (n *Node) Direction() math32.Vector3 {

	return n.direction
}

// SetMatrix sets the local transformation matrix.
func (n *Node) SetMatrix(m *math32.Matrix4) {

	n.matrix = *m
	n.changed = true
}

// Matrix returns a copy of the local transformation matrix.
func (n *Node) Matrix() math32.Matrix4 {

	return n.matrix
}

// WorldPosition updates the world matrix and sets
// the specified vector to the current world position of this node.
func (n *Node) WorldPosition(result *math32.Vector3) {

	n.UpdateMatrixWorld()
	result.SetFromMatrixPosition(&n.matrixWorld)
}

// WorldQuaternion updates the world matrix and sets
// the specified quaternion to the current world quaternion of this node.
func (n *Node) WorldQuaternion(result *math32.Quaternion) {

	var position math32.Vector3
	var scale math32.Vector3
	n.UpdateMatrixWorld()
	n.matrixWorld.Decompose(&position, result, &scale)
}

// WorldRotation updates the world matrix and sets
// the specified vector to the current world rotation of this node in Euler angles.
func (n *Node) WorldRotation(result *math32.Vector3) {

	var quaternion math32.Quaternion
	n.WorldQuaternion(&quaternion)
	result.SetFromQuaternion(&quaternion)
}

// WorldScale updates the world matrix and sets
// the specified vector to the current world scale of this node.
func (n *Node) WorldScale(result *math32.Vector3) {

	var position math32.Vector3
	var quaternion math32.Quaternion
	n.UpdateMatrixWorld()
	n.matrixWorld.Decompose(&position, &quaternion, result)
}

// WorldDirection updates the world matrix and sets
// the specified vector to the current world direction of this node.
func (n *Node) WorldDirection(result *math32.Vector3) {

	var quaternion math32.Quaternion
	n.WorldQuaternion(&quaternion)
	*result = n.direction
	result.ApplyQuaternion(&quaternion)
}

// MatrixWorld returns a copy of the matrix world of this node.
func (n *Node) MatrixWorld() math32.Matrix4 {

	return n.matrixWorld
}

// UpdateMatrix updates (if necessary) the local transform matrix
// of this node based on its position, quaternion, and scale.
func (n *Node) UpdateMatrix() bool {

	if !n.changed {
		return false
	}
	n.matrix.Compose(&n.position, &n.quaternion, &n.scale)
	n.changed = false
	return true
}

// UpdateMatrixWorld updates the world transform matrix for this node and for all of its children.
func (n *Node) UpdateMatrixWorld() {

	if n.parent == nil {
		n.updateMatrixWorld(&n.matrix)
	} else {
		parent := n.parent.GetNode()
		n.updateMatrixWorld(&parent.matrixWorld)
	}
}

// updateMatrixWorld is used internally by UpdateMatrixWorld.
// If the local transform matrix has changed, this method updates it and also the world matrix of this node.
// Children are updated recursively. If any node has changed, then we update the world matrix
// of all of its descendants regardless if their local matrices have changed.
func (n *Node) updateMatrixWorld(parentMatrixWorld *math32.Matrix4) {

	// If the local transform matrix for this node has changed then we need to update the local
	// matrix for this node and also the world matrix for this and all subsequent nodes.
	if n.UpdateMatrix() {
		n.matrixWorld.MultiplyMatrices(parentMatrixWorld, &n.matrix)

		// Update matrices of children recursively, always updating the world matrix
		for _, ichild := range n.children {
			ichild.GetNode().updateMatrixWorldNoCheck(&n.matrixWorld)
		}
	} else {
		// Update matrices of children recursively, continuing to check for changes
		for _, ichild := range n.children {
			ichild.GetNode().updateMatrixWorld(&n.matrixWorld)
		}
	}
}

// updateMatrixWorldNoCheck is used internally by updateMatrixWorld.
// This method should be called when a node has changed since it always updates the matrix world.
func (n *Node) updateMatrixWorldNoCheck(parentMatrixWorld *math32.Matrix4) {

	// Update the local transform matrix (if necessary)
	n.UpdateMatrix()

	// Always update the matrix world since an ancestor of this node has changed
	// (i.e. and ancestor had its local transform matrix modified)
	n.matrixWorld.MultiplyMatrices(parentMatrixWorld, &n.matrix)

	// Update matrices of children recursively
	for _, ichild := range n.children {
		ichild.GetNode().updateMatrixWorldNoCheck(&n.matrixWorld)
	}
}