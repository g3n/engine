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
	BoundingBox() math32.Box3
	Render(gs *gls.GLS)
	Clone() INode
	Dispose()
}

// Node represents an object in 3D space existing within a hierarchy.
type Node struct {
	Dispatcher                 // Embedded event dispatcher
	parent         INode       // Parent node
	children       []INode     // Children nodes
	name           string      // Optional node name
	loaderID       string      // ID used by loader
	visible        bool        // Whether the node is visible
	matNeedsUpdate bool        // Whether the the local matrix needs to be updated because position or scale has changed
	rotNeedsUpdate bool        // Whether the euler rotation and local matrix need to be updated because the quaternion has changed
	userData       interface{} // Generic user data

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

// BoundingBox satisfies the INode interface.
func (n *Node) BoundingBox() math32.Box3 {

	bbox := math32.Box3{n.position, n.position}
	for _, inode := range n.Children() {
		childBbox := inode.BoundingBox()
		bbox.Union(&childBbox)
	}
	return bbox
}

// Render satisfies the INode interface.
func (n *Node) Render(gs *gls.GLS) {
}

// Dispose satisfies the INode interface.
func (n *Node) Dispose() {
}

// Clone clones the Node and satisfies the INode interface.
func (n *Node) Clone() INode {

	clone := new(Node)

	// TODO clone Dispatcher?
	clone.Dispatcher.Initialize()

	clone.parent = n.parent
	clone.name = n.name + " (Clone)" // TODO append count?
	clone.loaderID = n.loaderID
	clone.visible = n.visible
	clone.userData = n.userData

	// Update matrix world and rotation if necessary
	n.UpdateMatrixWorld()
	n.Rotation()

	// Clone spatial properties
	clone.position = n.position
	clone.scale = n.scale
	clone.direction = n.direction
	clone.rotation = n.rotation
	clone.quaternion = n.quaternion
	clone.matrix = n.matrix
	clone.matrixWorld = n.matrixWorld
	clone.children = make([]INode, 0)

	// Clone children recursively
	for _, child := range n.children {
		clone.Add(child.Clone())
	}

	return clone
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
	n.matNeedsUpdate = true
}

// Visible returns the visibility of the node.
func (n *Node) Visible() bool {

	return n.visible
}

// SetChanged sets the matNeedsUpdate flag of the node.
func (n *Node) SetChanged(changed bool) {

	n.matNeedsUpdate = changed
}

// Changed returns the matNeedsUpdate flag of the node.
func (n *Node) Changed() bool {

	return n.matNeedsUpdate
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
	n.matNeedsUpdate = true
}

// SetPositionVec sets the position based on the specified vector pointer.
func (n *Node) SetPositionVec(vpos *math32.Vector3) {

	n.position = *vpos
	n.matNeedsUpdate = true
}

// SetPositionX sets the X coordinate of the position.
func (n *Node) SetPositionX(x float32) {

	n.position.X = x
	n.matNeedsUpdate = true
}

// SetPositionY sets the Y coordinate of the position.
func (n *Node) SetPositionY(y float32) {

	n.position.Y = y
	n.matNeedsUpdate = true
}

// SetPositionZ sets the Z coordinate of the position.
func (n *Node) SetPositionZ(z float32) {

	n.position.Z = z
	n.matNeedsUpdate = true
}

// Position returns the position as a vector.
func (n *Node) Position() math32.Vector3 {

	return n.position
}

// TranslateOnAxis translates the specified distance on the specified local axis.
func (n *Node) TranslateOnAxis(axis *math32.Vector3, dist float32) {

	v := math32.NewVec3().Copy(axis)
	v.ApplyQuaternion(&n.quaternion)
	v.MultiplyScalar(dist)
	n.position.Add(v)
	n.matNeedsUpdate = true
}

// TranslateX translates the specified distance on the local X axis.
func (n *Node) TranslateX(dist float32) {

	n.TranslateOnAxis(&math32.Vector3{1,0,0}, dist)
}

// TranslateY translates the specified distance on the local Y axis.
func (n *Node) TranslateY(dist float32) {

	n.TranslateOnAxis(&math32.Vector3{0,1,0}, dist)
}

// TranslateZ translates the specified distance on the local Z axis.
func (n *Node) TranslateZ(dist float32) {

	n.TranslateOnAxis(&math32.Vector3{0,0,1}, dist)
}

// SetRotation sets the global rotation in Euler angles (radians).
func (n *Node) SetRotation(x, y, z float32) {

	n.rotation.Set(x, y, z)
	n.quaternion.SetFromEuler(&n.rotation)
	n.matNeedsUpdate = true
}

// SetRotationVec sets the global rotation in Euler angles (radians) based on the specified vector pointer.
func (n *Node) SetRotationVec(vrot *math32.Vector3) {

	n.rotation = *vrot
	n.quaternion.SetFromEuler(&n.rotation)
	n.matNeedsUpdate = true
}

// SetRotationQuat sets the global rotation based on the specified quaternion pointer.
func (n *Node) SetRotationQuat(quat *math32.Quaternion) {

	n.quaternion = *quat
	n.rotNeedsUpdate = true
}

// SetRotationX sets the global X rotation to the specified angle in radians.
func (n *Node) SetRotationX(x float32) {

	if n.rotNeedsUpdate {
		n.rotation.SetFromQuaternion(&n.quaternion)
		n.rotNeedsUpdate = false
	}
	n.rotation.X = x
	n.quaternion.SetFromEuler(&n.rotation)
	n.matNeedsUpdate = true
}

// SetRotationY sets the global Y rotation to the specified angle in radians.
func (n *Node) SetRotationY(y float32) {

	if n.rotNeedsUpdate {
		n.rotation.SetFromQuaternion(&n.quaternion)
		n.rotNeedsUpdate = false
	}
	n.rotation.Y = y
	n.quaternion.SetFromEuler(&n.rotation)
	n.matNeedsUpdate = true
}

// SetRotationZ sets the global Z rotation to the specified angle in radians.
func (n *Node) SetRotationZ(z float32) {

	if n.rotNeedsUpdate {
		n.rotation.SetFromQuaternion(&n.quaternion)
		n.rotNeedsUpdate = false
	}
	n.rotation.Z = z
	n.quaternion.SetFromEuler(&n.rotation)
	n.matNeedsUpdate = true
}

// Rotation returns the current global rotation in Euler angles (radians).
func (n *Node) Rotation() math32.Vector3 {

	if n.rotNeedsUpdate {
		n.rotation.SetFromQuaternion(&n.quaternion)
		n.rotNeedsUpdate = false
	}
	return n.rotation
}

// RotateOnAxis rotates around the specified local axis the specified angle in radians.
func (n *Node) RotateOnAxis(axis *math32.Vector3, angle float32) {

	var rotQuat math32.Quaternion
	rotQuat.SetFromAxisAngle(axis, angle)
	n.QuaternionMult(&rotQuat)
}

// RotateX rotates around the local X axis the specified angle in radians.
func (n *Node) RotateX(x float32) {

	n.RotateOnAxis(&math32.Vector3{1,0,0}, x)
}

// RotateY rotates around the local Y axis the specified angle in radians.
func (n *Node) RotateY(y float32) {

	n.RotateOnAxis(&math32.Vector3{0,1,0}, y)
}

// RotateZ rotates around the local Z axis the specified angle in radians.
func (n *Node) RotateZ(z float32) {

	n.RotateOnAxis(&math32.Vector3{0,0,1}, z)
}

// SetQuaternion sets the quaternion based on the specified quaternion unit multiples.
func (n *Node) SetQuaternion(x, y, z, w float32) {

	n.quaternion.Set(x, y, z, w)
	n.rotNeedsUpdate = true
}

// SetQuaternionVec sets the quaternion based on the specified quaternion unit multiples vector.
func (n *Node) SetQuaternionVec(q *math32.Vector4) {

	n.quaternion.Set(q.X, q.Y, q.Z, q.W)
	n.rotNeedsUpdate = true
}

// SetQuaternionQuat sets the quaternion based on the specified quaternion pointer.
func (n *Node) SetQuaternionQuat(q *math32.Quaternion) {

	n.quaternion = *q
	n.rotNeedsUpdate = true
}

// QuaternionMult multiplies the current quaternion by the specified quaternion.
func (n *Node) QuaternionMult(q *math32.Quaternion) {

	n.quaternion.Multiply(q)
	n.rotNeedsUpdate = true
}

// Quaternion returns the current quaternion.
func (n *Node) Quaternion() math32.Quaternion {

	return n.quaternion
}

// SetScale sets the scale.
func (n *Node) SetScale(x, y, z float32) {

	n.scale.Set(x, y, z)
	n.matNeedsUpdate = true
}

// SetScaleVec sets the scale based on the specified vector pointer.
func (n *Node) SetScaleVec(scale *math32.Vector3) {

	n.scale = *scale
	n.matNeedsUpdate = true
}

// SetScaleX sets the X scale.
func (n *Node) SetScaleX(sx float32) {

	n.scale.X = sx
	n.matNeedsUpdate = true
}

// SetScaleY sets the Y scale.
func (n *Node) SetScaleY(sy float32) {

	n.scale.Y = sy
	n.matNeedsUpdate = true
}

// SetScaleZ sets the Z scale.
func (n *Node) SetScaleZ(sz float32) {

	n.scale.Z = sz
	n.matNeedsUpdate = true
}

// Scale returns the current scale.
func (n *Node) Scale() math32.Vector3 {

	return n.scale
}

// SetDirection sets the direction.
func (n *Node) SetDirection(x, y, z float32) {

	n.direction.Set(x, y, z)
	n.matNeedsUpdate = true
}

// SetDirectionVec sets the direction based on a vector pointer.
func (n *Node) SetDirectionVec(vdir *math32.Vector3) {

	n.direction = *vdir
	n.matNeedsUpdate = true
}

// Direction returns the direction.
func (n *Node) Direction() math32.Vector3 {

	return n.direction
}

// SetMatrix sets the local transformation matrix.
func (n *Node) SetMatrix(m *math32.Matrix4) {

	n.matrix = *m
	n.matrix.Decompose(&n.position, &n.quaternion, &n.scale)
	n.rotNeedsUpdate = true
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

	if !n.matNeedsUpdate && !n.rotNeedsUpdate {
		return false
	}
	n.matrix.Compose(&n.position, &n.quaternion, &n.scale)
	n.matNeedsUpdate = false
	return true
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
