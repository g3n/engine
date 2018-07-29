// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package shape

import (
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/experimental/collision"
)

// ConvexHull is a convex triangle-based geometry used for collision detection and contact resolution.
type ConvexHull struct {
	geometry.Geometry

	// Cached geometry properties
	faces            [][3]math32.Vector3
	faceNormals      []math32.Vector3
	worldFaceNormals []math32.Vector3
	uniqueEdges      []math32.Vector3
	worldUniqueEdges []math32.Vector3

}

func NewConvexHull(geom *geometry.Geometry) *ConvexHull {

	ch := new(ConvexHull)
	// // TODO check if geometry is convex, panic if not
	//if !geom.IsConvex() {
	//	panic("geometry needs to be convex")
	//}
	// // TODO future: create function to break up geometry into convex shapes and add all shapes to body

	ch.Geometry = *geom

	// Perform single-time computations
	ch.computeFaceNormalsAndUniqueEdges()

	return ch
}

// Compute and store face normals and unique edges
func (ch *ConvexHull) computeFaceNormalsAndUniqueEdges() {

	ch.Geometry.ReadFaces(func(vA, vB, vC math32.Vector3) bool {

		// Store face vertices
		var face [3]math32.Vector3
		face[0] = vA
		face[1] = vB
		face[2] = vC
		ch.faces = append(ch.faces, face)

		// Compute edges
		edge1 := math32.NewVec3().SubVectors(&vB, &vA)
		edge2 := math32.NewVec3().SubVectors(&vC, &vB)
		edge3 := math32.NewVec3().SubVectors(&vA, &vC)

		// Compute and store face normal in b.faceNormals
		faceNormal := math32.NewVec3().CrossVectors(edge2, edge1)
		if faceNormal.Length() > 0 {
			faceNormal.Normalize().Negate()
		}
		ch.faceNormals = append(ch.faceNormals, *faceNormal)

		// Compare unique edges recorded so far with the three new face edges and store the unique ones
		tol := float32(1e-6)
		for p := 0; p < len(ch.uniqueEdges); p++ {
			ue := ch.uniqueEdges[p]
			if !ue.AlmostEquals(edge1, tol) {
				ch.uniqueEdges = append(ch.uniqueEdges, *edge1)
			}
			if !ue.AlmostEquals(edge2, tol) {
				ch.uniqueEdges = append(ch.uniqueEdges, *edge1)
			}
			if !ue.AlmostEquals(edge3, tol) {
				ch.uniqueEdges = append(ch.uniqueEdges, *edge1)
			}
		}

		return false
	})

	// Allocate space for worldFaceNormals and worldUniqueEdges
	ch.worldFaceNormals = make([]math32.Vector3, len(ch.faceNormals))
	ch.worldUniqueEdges = make([]math32.Vector3, len(ch.uniqueEdges))
}

// ComputeWorldFaceNormalsAndUniqueEdges
func (ch *ConvexHull) ComputeWorldFaceNormalsAndUniqueEdges(quat *math32.Quaternion) {

	// Re-compute world face normals from local face normals
	for i := 0; i < len(ch.faceNormals); i++ {
		ch.worldFaceNormals[i] = ch.faceNormals[i]
		ch.worldFaceNormals[i].ApplyQuaternion(quat)
	}
	// Re-compute world unique edges from local unique edges
	for i := 0; i < len(ch.uniqueEdges); i++ {
		ch.worldUniqueEdges[i] = ch.uniqueEdges[i]
		ch.worldUniqueEdges[i].ApplyQuaternion(quat)
	}
}

func (ch *ConvexHull) Faces() [][3]math32.Vector3 {

	return ch.faces
}

func (ch *ConvexHull) FaceNormals() []math32.Vector3 {

	return ch.faceNormals
}

func (ch *ConvexHull) WorldFaceNormals() []math32.Vector3 {

	return ch.worldFaceNormals
}

func (ch *ConvexHull) UniqueEdges() []math32.Vector3 {

	return ch.uniqueEdges
}

func (ch *ConvexHull) WorldUniqueEdges() []math32.Vector3 {

	return ch.worldUniqueEdges
}

// FindPenetrationAxis finds the penetration axis between two convex bodies.
// The normal points from bodyA to bodyB.
// Returns false if there is no penetration. If there is a penetration - returns true and the penetration axis.
func (ch *ConvexHull) FindPenetrationAxis(chB *ConvexHull, posA, posB *math32.Vector3, quatA, quatB *math32.Quaternion) (bool, math32.Vector3) {

	// Keep track of the smaller depth found so far
	// Note that the penetration axis is the one that causes
	// the smallest penetration depth when the two hulls are squished onto that axis!
	// (may seem a bit counter-intuitive)
	depthMin := math32.Inf(1)

	var penetrationAxis math32.Vector3
	var depth float32

	// Assume the geometries are penetrating.
	// As soon as (and if) we figure out that they are not, then return false.
	penetrating := true

	worldFaceNormalsA := ch.WorldFaceNormals()
	worldFaceNormalsB := ch.WorldFaceNormals()

	// Check world normals of body A
	for _, worldFaceNormal := range worldFaceNormalsA {
		// Check whether the face is colliding with geomB
		penetrating, depth = ch.TestPenetrationAxis(chB, &worldFaceNormal, posA, posB, quatA, quatB)
		if !penetrating {
			return false, penetrationAxis // penetrationAxis doesn't matter since not penetrating
		}
		if depth < depthMin {
			depthMin = depth
			penetrationAxis.Copy(&worldFaceNormal)
		}
	}

	// Check world normals of body B
	for _, worldFaceNormal := range worldFaceNormalsB {
		// Check whether the face is colliding with geomB
		penetrating, depth = ch.TestPenetrationAxis(chB, &worldFaceNormal, posA, posB, quatA, quatB)
		if !penetrating {
			return false, penetrationAxis // penetrationAxis doesn't matter since not penetrating
		}
		if depth < depthMin {
			depthMin = depth
			penetrationAxis.Copy(&worldFaceNormal)
		}
	}

	worldUniqueEdgesA := ch.WorldUniqueEdges()
	worldUniqueEdgesB := ch.WorldUniqueEdges()

	// Check all combinations of unique world edges
	for _, worldUniqueEdgeA := range worldUniqueEdgesA {
		for _, worldUniqueEdgeB := range worldUniqueEdgesB {
			// Cross edges
			edgeCross := math32.NewVec3().CrossVectors(&worldUniqueEdgeA, &worldUniqueEdgeB)
			// If the edges are not aligned
			tol := float32(1e-6)
			if edgeCross.Length() > tol { // Cross product is not close to zero
				edgeCross.Normalize()
				penetrating, depth = ch.TestPenetrationAxis(chB, edgeCross, posA, posB, quatA, quatB)
				if !penetrating {
					return false, penetrationAxis
				}
				if depth < depthMin {
					depthMin = depth
					penetrationAxis.Copy(edgeCross)
				}
			}
		}
	}

	deltaC := math32.NewVec3().SubVectors(posA, posB)
	if deltaC.Dot(&penetrationAxis) > 0.0 {
		penetrationAxis.Negate()
	}

	return true, penetrationAxis
}

// Both hulls are projected onto the axis and the overlap size (penetration depth) is returned if there is one.
// return {number} The overlap depth, or FALSE if no penetration.
func (ch *ConvexHull) TestPenetrationAxis(chB *ConvexHull, worldAxis, posA, posB *math32.Vector3, quatA, quatB *math32.Quaternion) (bool, float32) {

	maxA, minA := ch.ProjectOntoWorldAxis(worldAxis, posA, quatA)
	maxB, minB := chB.ProjectOntoWorldAxis(worldAxis, posB, quatB)

	if maxA < minB || maxB < minA {
		return false, 0 // Separated
	}

	d0 := maxA - minB
	d1 := maxB - minA

	if d0 < d1 {
		return true, d0
	} else {
		return true, d1
	}
}

// ProjectOntoWorldAxis projects the geometry onto the specified world axis.
func (ch *ConvexHull) ProjectOntoWorldAxis(worldAxis, pos *math32.Vector3, quat *math32.Quaternion) (float32, float32) {

	// Transform the world axis to local
	quatConj := quat.Conjugate()
	localAxis := worldAxis.Clone().ApplyQuaternion(quatConj)

	// Project onto the local axis
	max, min := ch.Geometry.ProjectOntoAxis(localAxis)

	// Offset to obtain values relative to world origin
	localOrigin := math32.NewVec3().Sub(pos).ApplyQuaternion(quatConj)
	add := localOrigin.Dot(localAxis)
	min -= add
	max -= add

	return max, min
}

// =====================================================================

//{array} result The an array of contact point objects, see clipFaceAgainstHull
func (ch *ConvexHull) ClipAgainstHull(chB *ConvexHull, posA, posB *math32.Vector3, quatA, quatB *math32.Quaternion, penAxis *math32.Vector3, minDist, maxDist float32) []collision.Contact {

	var contacts []collision.Contact

	// Invert penetration axis so it points from b to a
	invPenAxis := penAxis.Clone().Negate()

	// Find face of B that is closest (i.e. that is most aligned with the penetration axis)
	closestFaceBidx := -1
	dmax := math32.Inf(-1)
	worldFaceNormalsB := chB.WorldFaceNormals()
	for i, worldFaceNormal := range worldFaceNormalsB {
		// Note - normals must be pointing out of the body so that they align with the penetration axis in the line below
		d := worldFaceNormal.Dot(invPenAxis)
		if d > dmax {
			dmax = d
			closestFaceBidx = i
		}
	}

	// If found a closest face (sometimes we don't find one)
	if closestFaceBidx >= 0 {

		// Copy and transform face vertices to world coordinates
		faces := chB.Faces()
		worldClosestFaceB := ch.WorldFace(faces[closestFaceBidx], posB, quatB)

		// Clip the closest world face of B to only the portion that is inside the hull of A
		contacts = ch.clipFaceAgainstHull(posA, penAxis, quatA, worldClosestFaceB, minDist, maxDist)
	}

	return contacts
}

func (ch *ConvexHull) WorldFace(face [3]math32.Vector3, pos *math32.Vector3, quat *math32.Quaternion) [3]math32.Vector3 {

	var result [3]math32.Vector3
	result[0] = face[0]
	result[1] = face[1]
	result[2] = face[2]
	result[0].ApplyQuaternion(quat).Add(pos)
	result[1].ApplyQuaternion(quat).Add(pos)
	result[2].ApplyQuaternion(quat).Add(pos)
	return result
}

// Clip a face against a hull.
//@param {Array} worldVertsB1 An array of Vec3 with vertices in the world frame.
//@param Array result Array to store resulting contact points in. Will be objects with properties: point, depth, normal. These are represented in world coordinates.
func (ch *ConvexHull) clipFaceAgainstHull(posA, penAxis *math32.Vector3, quatA *math32.Quaternion, worldClosestFaceB [3]math32.Vector3, minDist, maxDist float32) []collision.Contact {

	contacts := make([]collision.Contact, 0)

	// Find the face of A with normal closest to the separating axis (i.e. that is most aligned with the penetration axis)
	closestFaceAidx := -1
	dmax := math32.Inf(-1)
	worldFaceNormalsA := ch.WorldFaceNormals()
	for i, worldFaceNormal := range worldFaceNormalsA {
		// Note - normals must be pointing out of the body so that they align with the penetration axis in the line below
		d := worldFaceNormal.Dot(penAxis)
		if d > dmax {
			dmax = d
			closestFaceAidx = i
		}
	}

	if closestFaceAidx < 0 {
		// Did not find any closest face...
		return contacts
	}

	//console.log("closest A: ",worldClosestFaceA);

	// Get the face and construct connected faces
	facesA := ch.Faces()
	//worldClosestFaceA := n.WorldFace(facesA[closestFaceAidx], bodyA)
	closestFaceA := facesA[closestFaceAidx]
	connectedFaces := make([]int, 0) // indexes of the connected faces
	for faceIdx := 0; faceIdx < len(facesA); faceIdx++ {
		// Skip worldClosestFaceA
		if faceIdx == closestFaceAidx {
			continue
		}
		// Test that face has not already been added
		for _, cfidx := range connectedFaces {
			if cfidx == faceIdx {
				continue
			}
		}
		face := facesA[faceIdx]
		// Loop through face vertices and see if any of them are also present in the closest face
		// If a vertex is shared and this connected face hasn't been recorded yet - record and break inner loop
		for pConnFaceVidx := 0; pConnFaceVidx < len(face); pConnFaceVidx++ {
			var goToNextFace bool
			// Test if face shares a vertex with closetFaceA - add it to connectedFaces if so and break out of both loops
			for closFaceVidx := 0; closFaceVidx < len(closestFaceA); closFaceVidx++ {
				if closestFaceA[closFaceVidx].Equals(&face[pConnFaceVidx]) {
					connectedFaces = append(connectedFaces, faceIdx)
					goToNextFace = true
					break
				}
			}
			if goToNextFace {
				break
			}
		}
	}


	worldClosestFaceA := ch.WorldFace(closestFaceA, posA, quatA)
	// DEBUGGING
	//if n.debugging {
	//	//log.Error("CONN-FACES: %v", len(connectedFaces))
	//	for _, fidx := range connectedFaces {
	//		wFace := n.WorldFace(facesA[fidx], bodyA)
	//		ShowWorldFace(n.simulation.Scene(), wFace[:], &math32.Color{0.8, 0.8, 0.8})
	//	}
	//	//log.Error("worldClosestFaceA: %v", worldClosestFaceA)
	//	//log.Error("worldClosestFaceB: %v", worldClosestFaceB)
	//	ShowWorldFace(n.simulation.Scene(), worldClosestFaceA[:], &math32.Color{2, 0, 0})
	//	ShowWorldFace(n.simulation.Scene(), worldClosestFaceB[:], &math32.Color{0, 2, 0})
	//}

	clippedFace := make([]math32.Vector3, len(worldClosestFaceB))
	for i, v := range worldClosestFaceB {
		clippedFace[i] = v
	}

	// TODO port simplified loop to cannon.js once done and verified
	// https://github.com/schteppe/cannon.js/issues/378
	// https://github.com/TheRohans/cannon.js/commit/62a1ce47a851b7045e68f7b120b9e4ecb0d91aab#r29106924
	// Iterate over connected faces and clip the planes associated with their normals
	for _, cfidx := range connectedFaces {
		connFace := facesA[cfidx]
		connFaceNormal := worldFaceNormalsA[cfidx]
		// Choose a vertex in the connected face and use it to find the plane constant
		worldFirstVertex := connFace[0].Clone().ApplyQuaternion(quatA).Add(posA)
		planeDelta := - worldFirstVertex.Dot(&connFaceNormal)
		clippedFace = ch.clipFaceAgainstPlane(clippedFace, connFaceNormal.Clone(), planeDelta)
	}

	// Plot clipped face
	//if n.debugging {
	//	log.Error("worldClosestFaceBClipped: %v", clippedFace)
	//	ShowWorldFace(n.simulation.Scene(), clippedFace, &math32.Color{0, 0, 2})
	//}

	closestFaceAnormal := worldFaceNormalsA[closestFaceAidx]
	worldFirstVertex := worldClosestFaceA[0].Clone()//.ApplyQuaternion(quatA).Add(&posA)
	planeDelta := -worldFirstVertex.Dot(&closestFaceAnormal)

	// For each vertex in the clipped face resolve its depth (relative to the closestFaceA) and create a contact
	for _, vertex := range clippedFace {
		depth := closestFaceAnormal.Dot(&vertex) + planeDelta
		// Cap distance
		if depth <= minDist {
			depth = minDist
		}
		if depth <= maxDist {
			if depth <= 0 {
				contacts = append(contacts, collision.Contact{
					Point: vertex,
					Normal: closestFaceAnormal,
					Depth: depth,
				})
			}
		}

	}

	return contacts
}


// clipFaceAgainstPlane clips the specified face against the back of the specified plane.
// This is used multiple times when finding the portion of a face of one convex hull that is inside another convex hull.
// planeNormal and planeConstant satisfy the plane equation n*x = p where n is the planeNormal and p is the planeConstant (and x is a point on the plane).
func (ch *ConvexHull) clipFaceAgainstPlane(face []math32.Vector3, planeNormal *math32.Vector3, planeConstant float32) []math32.Vector3 {

	// inVertices are the verts making up the face of hullB

	clippedFace := make([]math32.Vector3, 0)

	// If not a face (if an edge or a vertex) - don't clip it
	if len(face) < 2 {
		return face
	}

	firstVertex := face[len(face)-1]
	dotFirst := planeNormal.Dot(&firstVertex) + planeConstant

	for vi := 0; vi < len(face); vi++ {
		lastVertex := face[vi]
		dotLast := planeNormal.Dot(&lastVertex) + planeConstant
		if dotFirst < 0 { // Inside hull
			if dotLast < 0 { // Start < 0, end < 0, so output lastVertex
				clippedFace = append(clippedFace, lastVertex)
			} else { // Start < 0, end >= 0, so output intersection
				newv := firstVertex.Clone().Lerp(&lastVertex, dotFirst / (dotFirst - dotLast))
				clippedFace = append(clippedFace, *newv)
			}
		} else { // Outside hull
			if dotLast < 0 { // Start >= 0, end < 0 so output intersection and end
				newv := firstVertex.Clone().Lerp(&lastVertex, dotFirst / (dotFirst - dotLast))
				clippedFace = append(clippedFace, *newv)
				clippedFace = append(clippedFace, lastVertex)
			}
		}
		firstVertex = lastVertex
		dotFirst = dotLast
	}

	return clippedFace
}
