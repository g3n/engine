// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"strings"

	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/math32"
)

// Interface for all node types
type INode interface {
	GetNode() *Node
	UpdateMatrixWorld()
	Raycast(*Raycaster, *[]Intersect)
	Render(gs *gls.GLS)
	Dispose()
}

type Node struct {
	Dispatcher                    // Embedded event dispatcher
	loaderID    string            // ID used by loader
	name        string            // Optional node name
	position    math32.Vector3    // Node position, specified as a Vector3
	rotation    math32.Vector3    // Node rotation, specified in Euler angles.
	quaternion  math32.Quaternion // Node rotation, specified as a Quaternion.
	scale       math32.Vector3    // Node scale as a Vector3
	direction   math32.Vector3    // Initial direction
	matrix      math32.Matrix4    // Transform matrix relative to this node parent.
	matrixWorld math32.Matrix4    // Transform world matrix
	visible     bool              // Visible flag
	changed     bool              // Node position/orientation/scale changed
	parent      INode             // Parent node
	children    []INode           // Array with node children
	userData    interface{}       // Generic user data
}

// NewNode creates and returns a pointer to a new Node
func NewNode() *Node {

	n := new(Node)
	n.Init()
	return n
}

// Init initializes this Node
// It is normally use by other types which embed a Node
func (n *Node) Init() {

	n.Dispatcher.Initialize()
	n.position.Set(0, 0, 0)
	n.rotation.Set(0, 0, 0)
	n.quaternion.Set(0, 0, 0, 1)
	n.scale.Set(1, 1, 1)
	n.direction.Set(0, 0, 1)
	n.matrix.Identity()
	n.matrixWorld.Identity()
	n.children = make([]INode, 0)
	n.visible = true
	n.changed = true
}

// GetNode satisfies the INode interface and returns
// a pointer to the embedded Node
func (n *Node) GetNode() *Node {

	return n
}

// Raycast satisfies the INode interface
func (n *Node) Raycast(rc *Raycaster, intersects *[]Intersect) {
}

// Render satisfies the INode interface
func (n *Node) Render(gs *gls.GLS) {
}

// Dispose satisfies the INode interface
func (n *Node) Dispose() {
}

// SetLoaderID is normally used by external loaders, such as Collada,
// to assign an ID to the node with the ID value in the node description
// Can be used to find other loaded nodes.
func (n *Node) SetLoaderID(id string) {

	n.loaderID = id
}

// LoaderID returns an optional ID set when this node was
// created by an external loader such as Collada
func (n *Node) LoaderID() string {

	return n.loaderID
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

// FindLoaderID looks in the specified node and all its children
// for a node with the specifid loaderID and if found returns it.
// Returns nil if not found
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

// SetChanged sets this node changed flag
func (n *Node) SetChanged(changed bool) {

	n.changed = changed
}

// Changed returns this Node changed flag
func (n *Node) Changed() bool {

	return n.changed
}

// SetName set an option name for the node.
// This name can be used for debugging or other purposes.
func (n *Node) SetName(name string) {

	n.name = name
}

// Name returns current optional name for this node
func (n *Node) Name() string {

	return n.name
}

// SetPosition sets this node world position
func (n *Node) SetPosition(x, y, z float32) {

	n.position.Set(x, y, z)
	n.changed = true
}

// SetPositionVec sets this node position from the specified vector pointer
func (n *Node) SetPositionVec(vpos *math32.Vector3) {

	n.position = *vpos
	n.changed = true
}

// SetPositionX sets the x coordinate of this node position
func (n *Node) SetPositionX(x float32) {

	n.position.X = x
	n.changed = true
}

// SetPositionY sets the y coordinate of this node position
func (n *Node) SetPositionY(y float32) {

	n.position.Y = y
	n.changed = true
}

// SetPositionZ sets the z coordinate of this node position
func (n *Node) SetPositionZ(z float32) {

	n.position.Z = z
	n.changed = true
}

// Position returns the current node position as a vector
func (n *Node) Position() math32.Vector3 {

	return n.position
}

// SetRotation sets the three fields of the node rotation in radians
// The node quaternion is updated
func (n *Node) SetRotation(x, y, z float32) {

	n.rotation.Set(x, y, z)
	n.quaternion.SetFromEuler(&n.rotation)
	n.changed = true
}

// SetRotationX sets the x rotation angle in radians
// The node quaternion is updated
func (n *Node) SetRotationX(x float32) {

	n.rotation.X = x
	n.quaternion.SetFromEuler(&n.rotation)
	n.changed = true
}

// SetRotationY sets the y rotation angle in radians
// The node quaternion is updated
func (n *Node) SetRotationY(y float32) {

	n.rotation.Y = y
	n.quaternion.SetFromEuler(&n.rotation)
	n.changed = true
}

// SetRotationZ sets the z rotation angle in radians
// The node quaternion is updated
func (n *Node) SetRotationZ(z float32) {

	n.rotation.Z = z
	n.quaternion.SetFromEuler(&n.rotation)
	n.changed = true
}

// AddRotationX adds to the current rotation x coordinate in radians
// The node quaternion is updated
func (n *Node) AddRotationX(x float32) {

	n.rotation.X += x
	n.quaternion.SetFromEuler(&n.rotation)
	n.changed = true
}

// AddRotationY adds to the current rotation y coordinate in radians
// The node quaternion is updated
func (n *Node) AddRotationY(y float32) {

	n.rotation.Y += y
	n.quaternion.SetFromEuler(&n.rotation)
	n.changed = true
}

// AddRotationZ adds to the current rotation z coordinate in radians
// The node quaternion is updated
func (n *Node) AddRotationZ(z float32) {

	n.rotation.Z += z
	n.quaternion.SetFromEuler(&n.rotation)
	n.changed = true
}

// Rotation returns the current rotation
func (n *Node) Rotation() math32.Vector3 {

	return n.rotation
}

// SetQuaternion sets this node  quaternion with the specified fields
func (n *Node) SetQuaternion(x, y, z, w float32) {

	n.quaternion.Set(x, y, z, w)
	n.changed = true
}

// SetQuaternionQuat sets this node quaternion from the specified quaternion pointer
func (n *Node) SetQuaternionQuat(q *math32.Quaternion) {

	n.quaternion = *q
	n.changed = true
}

// QuaternionMult multiplies the quaternion by the specified quaternion
func (n *Node) QuaternionMult(q *math32.Quaternion) {

	n.quaternion.Multiply(q)
	n.changed = true
}

// Quaternion returns the current quaternion
func (n *Node) Quaternion() math32.Quaternion {

	return n.quaternion
}

// SetScale sets this node scale fields
func (n *Node) SetScale(x, y, z float32) {

	n.scale.Set(x, y, z)
	n.changed = true
}

// SetScaleVec sets this node scale from a pointer to a Vector3
func (n *Node) SetScaleVec(scale *math32.Vector3) {

	n.scale = *scale
	n.changed = true
}

// SetScaleX sets the X scale of this node
func (n *Node) SetScaleX(sx float32) {

	n.scale.X = sx
	n.changed = true
}

// SetScaleY sets the Y scale of this node
func (n *Node) SetScaleY(sy float32) {

	n.scale.Y = sy
	n.changed = true
}

// SetScaleZ sets the Z scale of this node
func (n *Node) SetScaleZ(sz float32) {

	n.scale.Z = sz
	n.changed = true
}

// Scale returns the current scale
func (n *Node) Scale() math32.Vector3 {

	return n.scale
}

// SetDirection sets this node initial direction vector
func (n *Node) SetDirection(x, y, z float32) {

	n.direction.Set(x, y, z)
	n.changed = true
}

// SetDirectionVec sets this node initial direction vector from a pointer to Vector3
func (n *Node) SetDirectionVec(vdir *math32.Vector3) {

	n.direction = *vdir
	n.changed = true
}

// Direction returns this node initial direction
func (n *Node) Direction() math32.Vector3 {

	return n.direction
}

// SetMatrix sets this node local transformation matrix
func (n *Node) SetMatrix(m *math32.Matrix4) {

	n.matrix = *m
	n.changed = true
}

// Matrix returns a copy of this node local transformation matrix
func (n *Node) Matrix() math32.Matrix4 {

	return n.matrix
}

// SetVisible sets the node visibility state
func (n *Node) SetVisible(state bool) {

	n.visible = state
	n.changed = true
}

// Visible returns the node visibility state
func (n *Node) Visible() bool {

	return n.visible
}

// WorldPosition updates this node world matrix and gets
// the current world position vector.
func (n *Node) WorldPosition(result *math32.Vector3) {

	n.UpdateMatrixWorld()
	result.SetFromMatrixPosition(&n.matrixWorld)
}

// WorldQuaternion sets the specified result quaternion with
// this node current world quaternion
func (n *Node) WorldQuaternion(result *math32.Quaternion) {

	var position math32.Vector3
	var scale math32.Vector3
	n.UpdateMatrixWorld()
	n.matrixWorld.Decompose(&position, result, &scale)
}

// WorldRotation sets the specified result vector with
// current world rotation of this node in Euler angles.
func (n *Node) WorldRotation(result *math32.Vector3) {

	var quaternion math32.Quaternion
	n.WorldQuaternion(&quaternion)
	result.SetFromQuaternion(&quaternion)
}

// WorldScale sets the specified result vector with
// the current world scale of this node
func (n *Node) WorldScale(result *math32.Vector3) {

	var position math32.Vector3
	var quaternion math32.Quaternion
	n.UpdateMatrixWorld()
	n.matrixWorld.Decompose(&position, &quaternion, result)
}

// WorldDirection updates this object world matrix and sets
// the current world direction.
func (n *Node) WorldDirection(result *math32.Vector3) {

	var quaternion math32.Quaternion
	n.WorldQuaternion(&quaternion)
	*result = n.direction
	result.ApplyQuaternion(&quaternion)
}

// MatrixWorld returns a copy of this node matrix world
func (n *Node) MatrixWorld() math32.Matrix4 {

	return n.matrixWorld
}

// UpdateMatrix updates this node local matrix transform from its
// current position, quaternion and scale.
func (n *Node) UpdateMatrix() {

	n.matrix.Compose(&n.position, &n.quaternion, &n.scale)
}

// UpdateMatrixWorld updates this node world transform matrix and of all its children
func (n *Node) UpdateMatrixWorld() {

	n.UpdateMatrix()
	if n.parent == nil {
		n.matrixWorld = n.matrix
	} else {
		parent := n.parent.GetNode()
		n.matrixWorld.MultiplyMatrices(&parent.matrixWorld, &n.matrix)
	}
	// Update this Node children matrices
	for _, ichild := range n.children {
		ichild.UpdateMatrixWorld()
	}
}

// SetParent sets this node parent
func (n *Node) SetParent(iparent INode) {

	n.parent = iparent
}

// Parent returns this node parent
func (n *Node) Parent() INode {

	return n.parent
}

// Children returns the list of this node children
func (n *Node) Children() []INode {

	return n.children
}

// Add adds the specified INode to this node list of children
func (n *Node) Add(ichild INode) *Node {

	child := ichild.GetNode()
	if n == child {
		panic("Node.Add: object can't be added as a child of itself")
		return nil
	}
	// If this child already has a parent,
	// removes it from this parent children list
	if child.parent != nil {
		child.parent.GetNode().Remove(ichild)
	}
	child.parent = n
	n.children = append(n.children, ichild)
	return n
}

// Remove removes the specified INode from this node list of children
// Returns true if found or false otherwise
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

// RemoveAll removes all children from this node
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

// DisposeChildren removes and disposes all children of this
// node and if 'recurs' is true for each of its children recursively.
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

// SetUserData sets this node associated generic user data
func (n *Node) SetUserData(data interface{}) {

	n.userData = data
}

// UserData returns this node associated generic user data
func (n *Node) UserData() interface{} {

	return n.userData
}
