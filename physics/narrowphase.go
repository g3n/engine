// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package physics

import (
	"github.com/g3n/engine/physics/object"
	"github.com/g3n/engine/physics/collision"
	"github.com/g3n/engine/physics/equation"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/physics/material"
)

// Narrowphase
type Narrowphase struct {
	simulation *Simulation
	currentContactMaterial *material.ContactMaterial

	enableFrictionReduction bool // If true friction is computed as average

	debugging bool
}

type Pair struct {
	bodyA *object.Body
	bodyB *object.Body
}

// NewNarrowphase creates and returns a pointer to a new Narrowphase.
func NewNarrowphase(simulation *Simulation) *Narrowphase {

	n := new(Narrowphase)
	n.simulation = simulation
	//n.enableFrictionReduction = true

	// FOR DEBUGGING
	//n.debugging = true

	return n
}


func (n *Narrowphase) GetContacts(pairs []collision.Pair) ([]*equation.Contact, []*equation.Friction) {

	allContactEqs := make([]*equation.Contact, 0)
	allFrictionEqs := make([]*equation.Friction, 0)

	for k := 0; k < len(pairs); k++ {

		// Get current collision bodies
		bodyA := pairs[k].BodyA
		bodyB := pairs[k].BodyB

		bodyTypeA := bodyA.BodyType()
		bodyTypeB := bodyB.BodyType()

		// For now these collisions are ignored
		// TODO future: just want to check for collision (in order to dispatch events) and not create equations
		justTest := (bodyTypeA == object.Kinematic) && (bodyTypeB == object.Static) ||
				    (bodyTypeA == object.Static)    && (bodyTypeB == object.Kinematic) ||
				    (bodyTypeA == object.Kinematic) && (bodyTypeB == object.Kinematic)

		// Get contacts
		if !justTest {
			_, contactEqs, frictionEqs := n.Resolve(bodyA, bodyB)
			allContactEqs = append(allContactEqs, contactEqs...)
			allFrictionEqs = append(allFrictionEqs, frictionEqs...)
		}
   }

   return allContactEqs, allFrictionEqs
}

// Convex - Convex collision detection
//func (n *Narrowphase) Resolve(si,sj,xi,xj,qi,qj,bi,bj,rsi,rsj,justTest) {
func (n *Narrowphase) Resolve(bodyA, bodyB *object.Body) (bool, []*equation.Contact, []*equation.Friction) {

	contactEqs := make([]*equation.Contact, 0)
	frictionEqs := make([]*equation.Friction, 0)

	// Check if colliding and find penetration axis
	penetrating, penAxis := n.FindPenetrationAxis(bodyA, bodyB)

	if penetrating {
		// Colliding! Find contacts.
		if n.debugging {
			ShowPenAxis(n.simulation.Scene(), &penAxis) //, -1000, 1000)
			log.Error("Colliding (%v|%v) penAxis: %v", bodyA.Name(), bodyB.Name(), penAxis)
		}
		contacts := n.ClipAgainstHull(bodyA, bodyB, &penAxis, -100, 100)
		//log.Error(" .... contacts: %v", contacts)

		posA := bodyA.Position()
		posB := bodyB.Position()

		for j := 0; j < len(contacts); j++ {

			contact := contacts[j]
			if n.debugging {
				ShowContact(n.simulation.Scene(), &contact) // TODO DEBUGGING
			}

			// Note - contact Normals point from B to A (contacts live in B)

			// Create contact equation and append it
			contactEq := equation.NewContact(bodyA, bodyB, 0, 1e6)
			contactEq.SetSpookParams(1e6, 3, n.simulation.dt)
			contactEq.SetEnabled(bodyA.CollisionResponse() && bodyB.CollisionResponse())
			contactEq.SetNormal(penAxis.Clone())

			//log.Error("contact.Depth: %v", contact.Depth)

			contactEq.SetRA(contact.Normal.Clone().MultiplyScalar(-contact.Depth).Add(&contact.Point).Sub(&posA))
			contactEq.SetRB(contact.Point.Clone().Sub(&posB))
			contactEqs = append(contactEqs, contactEq)

			// If enableFrictionReduction is true then skip creating friction equations for individual contacts
			// We will create average friction equations later based on all contacts
			// TODO
			if !n.enableFrictionReduction {
				fEq1, fEq2 := n.createFrictionEquationsFromContact(contactEq)
				frictionEqs = append(frictionEqs, fEq1, fEq2)
			}
		}

		// If enableFrictionReduction is true then we skipped creating friction equations for individual contacts
		// We now want to create average friction equations based on all contact points.
		// If we only have one contact however, then friction is small and we don't need to create the friction equations at all.
		// TODO
		if n.enableFrictionReduction && len(contactEqs) > 1 {
			//fEq1, fEq2 := n.createFrictionFromAverage(contactEqs)
			//frictionEqs = append(frictionEqs, fEq1, fEq2)
		}
	}

	return false, contactEqs, frictionEqs
}

func (n *Narrowphase) createFrictionEquationsFromContact(contactEquation *equation.Contact) (*equation.Friction, *equation.Friction) { //}, outArray) bool {

	bodyA := n.simulation.bodies[contactEquation.BodyA().Index()]
	bodyB := n.simulation.bodies[contactEquation.BodyB().Index()]

	// TODO
	// friction = n.currentContactMaterial.friction
	// if materials are defined then override: friction = matA.friction * matB.friction
	//var mug = friction * world.gravity.length()
	//var reducedMass = bodyA.InvMass() + bodyB.InvMass()
	//if reducedMass > 0 {
	//	reducedMass = 1/reducedMass
	//}
	slipForce := float32(0.5) //mug*reducedMass

	fricEq1 := equation.NewFriction(bodyA, bodyB, slipForce)
	fricEq2 := equation.NewFriction(bodyA, bodyB, slipForce)

	fricEq1.SetSpookParams(1e7, 3, n.simulation.dt)
	fricEq2.SetSpookParams(1e7, 3, n.simulation.dt)

	// Copy over the relative vectors
	cRA := contactEquation.RA()
	cRB := contactEquation.RB()
	fricEq1.SetRA(&cRA)
	fricEq1.SetRB(&cRB)
	fricEq2.SetRA(&cRA)
	fricEq2.SetRB(&cRB)

	// Construct tangents
	cNormal := contactEquation.Normal()
	t1, t2 := cNormal.RandomTangents()
	fricEq1.SetTangent(t1)
	fricEq2.SetTangent(t2)

	// Copy enabled state
	cEnabled := contactEquation.Enabled()
	fricEq1.SetEnabled(cEnabled)
	fricEq2.SetEnabled(cEnabled)

	return fricEq1, fricEq2
}

func (n *Narrowphase) createFrictionFromAverage(contactEqs []*equation.Contact) (*equation.Friction, *equation.Friction) {

	// The last contactEquation
	lastContactEq := contactEqs[len(contactEqs)-1]

	// Create a friction equation based on the last contact (we will modify it to take into account all contacts)
	fEq1, fEq2 := n.createFrictionEquationsFromContact(lastContactEq)
	if (fEq1 == nil && fEq2 == nil) || len(contactEqs) == 1 {
		return fEq1, fEq2
	}

	averageNormal := math32.NewVec3()
	averageContactPointA := math32.NewVec3()
	averageContactPointB := math32.NewVec3()

	bodyA := lastContactEq.BodyA()
	//bodyB := lastContactEq.BodyB()
	normal := lastContactEq.Normal()
	rA := lastContactEq.RA()
	rB := lastContactEq.RB()

	for _, cEq := range contactEqs {
		if cEq.BodyA() != bodyA {
			averageNormal.Add(&normal)
			averageContactPointA.Add(&rA)
			averageContactPointB.Add(&rB)
		} else {
			averageNormal.Sub(&normal)
			averageContactPointA.Add(&rB)
			averageContactPointB.Add(&rA)
		}
	}

	invNumContacts := float32(1) / float32(len(contactEqs))

	averageContactPointA.MultiplyScalar(invNumContacts)
	averageContactPointB.MultiplyScalar(invNumContacts)

	// Should be the same for both friction equations
	fEq1.SetRA(averageContactPointA)
	fEq1.SetRB(averageContactPointB)
	fEq2.SetRA(averageContactPointA)
	fEq2.SetRB(averageContactPointB)

	// Set tangents
	averageNormal.Normalize()
	t1, t2 := averageNormal.RandomTangents()
	fEq1.SetTangent(t1)
	fEq2.SetTangent(t2)

	return fEq1, fEq2
}


//
// Penetration Axis =============================================
//

// FindPenetrationAxis finds the penetration axis between two convex bodies.
// The normal points from bodyA to bodyB.
// Returns false if there is no penetration. If there is a penetration - returns true and the penetration axis.
func (n *Narrowphase) FindPenetrationAxis(bodyA, bodyB *object.Body) (bool, math32.Vector3) {

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

	worldFaceNormalsA := bodyA.WorldFaceNormals()
	worldFaceNormalsB := bodyB.WorldFaceNormals()

	// Check world normals of body A
	for _, worldFaceNormal := range worldFaceNormalsA {
		// Check whether the face is colliding with geomB
		penetrating, depth = n.TestPenetrationAxis(&worldFaceNormal, bodyA, bodyB)
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
		penetrating, depth = n.TestPenetrationAxis(&worldFaceNormal, bodyA, bodyB)
		if !penetrating {
			return false, penetrationAxis // penetrationAxis doesn't matter since not penetrating
		}
		if depth < depthMin {
			depthMin = depth
			penetrationAxis.Copy(&worldFaceNormal)
		}
	}

	worldUniqueEdgesA := bodyA.WorldUniqueEdges()
	worldUniqueEdgesB := bodyB.WorldUniqueEdges()

	// Check all combinations of unique world edges
	for _, worldUniqueEdgeA := range worldUniqueEdgesA {
		for _, worldUniqueEdgeB := range worldUniqueEdgesB {
			// Cross edges
			edgeCross := math32.NewVec3().CrossVectors(&worldUniqueEdgeA, &worldUniqueEdgeB)
			// If the edges are not aligned
			tol := float32(1e-6)
			if edgeCross.Length() > tol { // Cross product is not close to zero
				edgeCross.Normalize()
				penetrating, depth = n.TestPenetrationAxis(edgeCross, bodyA, bodyB)
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

	posA := bodyA.Position()
	posB := bodyB.Position()

	deltaC := math32.NewVec3().SubVectors(&posA, &posB)
   	if deltaC.Dot(&penetrationAxis) > 0.0 {
       	penetrationAxis.Negate()
   	}

   	return true, penetrationAxis
}

// Both hulls are projected onto the axis and the overlap size (penetration depth) is returned if there is one.
// return {number} The overlap depth, or FALSE if no penetration.
func (n *Narrowphase) TestPenetrationAxis(worldAxis *math32.Vector3, bodyA, bodyB *object.Body) (bool, float32) {

	maxA, minA := n.ProjectOntoWorldAxis(bodyA, worldAxis)
	maxB, minB := n.ProjectOntoWorldAxis(bodyB, worldAxis)

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
func (n *Narrowphase) ProjectOntoWorldAxis(body *object.Body, axis *math32.Vector3) (float32, float32) {

	// Transform the axis to local
	quatConj := body.Quaternion().Conjugate()
	localAxis := axis.Clone().ApplyQuaternion(quatConj)
	max, min := body.GetGeometry().ProjectOntoAxis(localAxis)

	// Offset to obtain values relative to world origin
	bodyPos := body.Position()
	localOrigin := math32.NewVec3().Sub(&bodyPos).ApplyQuaternion(quatConj)
	add := localOrigin.Dot(localAxis)
	min -= add
	max -= add

	return max, min
}

//
// Contact Finding =============================================
//

// Contact describes a contact point.
type Contact struct {
	Point  math32.Vector3
	Normal math32.Vector3
	Depth  float32
}

//{array} result The an array of contact point objects, see clipFaceAgainstHull
func (n *Narrowphase) ClipAgainstHull(bodyA, bodyB *object.Body, penAxis *math32.Vector3, minDist, maxDist float32) []Contact {

	var contacts []Contact

	// Invert penetration axis so it points from b to a
	invPenAxis := penAxis.Clone().Negate()

	// Find face of B that is closest (i.e. that is most aligned with the penetration axis)
	closestFaceBidx := -1
	dmax := math32.Inf(-1)
	worldFaceNormalsB := bodyB.WorldFaceNormals()
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
		faces := bodyB.Faces()
		worldClosestFaceB := n.WorldFace(faces[closestFaceBidx], bodyB)
		contacts = n.ClipFaceAgainstHull(penAxis, bodyA, worldClosestFaceB, minDist, maxDist)
	}

	return contacts
}

func (n *Narrowphase) WorldFace(face [3]math32.Vector3, body *object.Body) [3]math32.Vector3 {

	var result [3]math32.Vector3
	result[0] = face[0]
	result[1] = face[1]
	result[2] = face[2]

	pos := body.Position()
	result[0].ApplyQuaternion(body.Quaternion()).Add(&pos)
	result[1].ApplyQuaternion(body.Quaternion()).Add(&pos)
	result[2].ApplyQuaternion(body.Quaternion()).Add(&pos)
	return result
}

func (n *Narrowphase) WorldFaceNormal(normal *math32.Vector3, body *object.Body) math32.Vector3 {

	pos := body.Position()
	result := normal.Clone().ApplyQuaternion(body.Quaternion()).Add(&pos)
	return *result
}

// TODO move to geometry ?
// Clip a face against a hull.
//@param {Array} worldVertsB1 An array of Vec3 with vertices in the world frame.
//@param Array result Array to store resulting contact points in. Will be objects with properties: point, depth, normal. These are represented in world coordinates.
func (n *Narrowphase) ClipFaceAgainstHull(penAxis *math32.Vector3, bodyA *object.Body, worldClosestFaceB [3]math32.Vector3, minDist, maxDist float32) []Contact {

	contacts := make([]Contact, 0)

	// Find the face of A with normal closest to the separating axis (i.e. that is most aligned with the penetration axis)
	closestFaceAidx := -1
	dmax := math32.Inf(-1)
	worldFaceNormalsA := bodyA.WorldFaceNormals()
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
	facesA := bodyA.Faces()
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


	worldClosestFaceA := n.WorldFace(closestFaceA, bodyA)
	// DEBUGGING
	if n.debugging {
		//log.Error("CONN-FACES: %v", len(connectedFaces))
		for _, fidx := range connectedFaces {
			wFace := n.WorldFace(facesA[fidx], bodyA)
			ShowWorldFace(n.simulation.Scene(), wFace[:], &math32.Color{0.8, 0.8, 0.8})
		}
		//log.Error("worldClosestFaceA: %v", worldClosestFaceA)
		//log.Error("worldClosestFaceB: %v", worldClosestFaceB)
		ShowWorldFace(n.simulation.Scene(), worldClosestFaceA[:], &math32.Color{2, 0, 0})
		ShowWorldFace(n.simulation.Scene(), worldClosestFaceB[:], &math32.Color{0, 2, 0})
	}

	clippedFace := make([]math32.Vector3, len(worldClosestFaceB))
	for i, v := range worldClosestFaceB {
		clippedFace[i] = v
	}

	// TODO port simplified loop to cannon.js once done and verified
	// https://github.com/schteppe/cannon.js/issues/378
	// https://github.com/TheRohans/cannon.js/commit/62a1ce47a851b7045e68f7b120b9e4ecb0d91aab#r29106924
	posA := bodyA.Position()
	quatA := bodyA.Quaternion()
	// Iterate over connected faces and clip the planes associated with their normals
	for _, cfidx := range connectedFaces {
		connFace := facesA[cfidx]
		connFaceNormal := worldFaceNormalsA[cfidx]
		// Choose a vertex in the connected face and use it to find the plane constant
		worldFirstVertex := connFace[0].Clone().ApplyQuaternion(quatA).Add(&posA)
		planeDelta := - worldFirstVertex.Dot(&connFaceNormal)
		clippedFace = n.ClipFaceAgainstPlane(clippedFace, connFaceNormal.Clone(), planeDelta)
	}

	// Plot clipped face
	if n.debugging {
		log.Error("worldClosestFaceBClipped: %v", clippedFace)
		ShowWorldFace(n.simulation.Scene(), clippedFace, &math32.Color{0, 0, 2})
	}

	closestFaceAnormal := worldFaceNormalsA[closestFaceAidx]
	worldFirstVertex := worldClosestFaceA[0].Clone()//.ApplyQuaternion(quatA).Add(&posA)
	planeDelta := -worldFirstVertex.Dot(&closestFaceAnormal)

	for _, vertex := range clippedFace {
		depth := closestFaceAnormal.Dot(&vertex) + planeDelta
		// Cap distance
		if depth <= minDist {
			depth = minDist
		}
		if depth <= maxDist {
			if depth <= 0 {
				contacts = append(contacts, Contact{
					Point: vertex,
					Normal: closestFaceAnormal,
					Depth: depth,
				})
			}
		}

	}

	return contacts
}


// Clip a face in a hull against the back of a plane.
// @param {Number} planeConstant The constant in the mathematical plane equation
func (n *Narrowphase) ClipFaceAgainstPlane(face []math32.Vector3, planeNormal *math32.Vector3, planeConstant float32) []math32.Vector3 {

	// inVertices are the verts making up the face of hullB

	clippedFace := make([]math32.Vector3, 0)

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

//// TODO ?
//func (n *Narrowphase) GetAveragePointLocal(target) {
//
//	target = target || new Vec3()
//   n := this.vertices.length
//   verts := this.vertices
//   for i := 0; i < n; i++ {
//       target.vadd(verts[i],target)
//   }
//   target.mult(1/n,target)
//   return target
//}
//
//
//// Checks whether p is inside the polyhedra. Must be in local coords.
//// The point lies outside of the convex hull of the other points if and only if
//// the direction of all the vectors from it to those other points are on less than one half of a sphere around it.
//// p is A point given in local coordinates
//func (n *Narrowphase) PointIsInside(p) {
//
//	verts := this.vertices
//	faces := this.faces
//	normals := this.faceNormals
//   positiveResult := null
//   N := this.faces.length
//   pointInside := ConvexPolyhedron_pointIsInside
//   this.getAveragePointLocal(pointInside)
//   for i := 0; i < N; i++ {
//       numVertices := this.faces[i].length
//       n := normals[i]
//       v := verts[faces[i][0]] // We only need one point in the face
//
//       // This dot product determines which side of the edge the point is
//       vToP := ConvexPolyhedron_vToP
//       p.vsub(v,vToP)
//       r1 := n.dot(vToP)
//
//       vToPointInside := ConvexPolyhedron_vToPointInside
//       pointInside.vsub(v,vToPointInside)
//       r2 := n.dot(vToPointInside)
//
//       if (r1<0 && r2>0) || (r1>0 && r2<0) {
//           return false // Encountered some other sign. Exit.
//       }
//   }
//
//   // If we got here, all dot products were of the same sign.
//   return positiveResult ? 1 : -1
//}

// TODO
//func (n *Narrowphase) planevConvex(
//	planeShape,
//	convexShape,
//	planePosition,
//	convexPosition,
//	planeQuat,
//	convexQuat,
//	planeBody,
//	convexBody,
//	si,
//	sj,
//	justTest) {
//
//	// Simply return the points behind the plane.
//    worldVertex := planeConvex_v
//    worldNormal := planeConvex_normal
//    worldNormal.set(0,0,1)
//    planeQuat.vmult(worldNormal,worldNormal) // Turn normal according to plane orientation
//
//    var numContacts = 0
//    var relpos = planeConvex_relpos
//    for i := 0; i < len(convexShape.vertices); i++ {
//
//        // Get world convex vertex
//        worldVertex.copy(convexShape.vertices[i])
//        convexQuat.vmult(worldVertex, worldVertex)
//        convexPosition.vadd(worldVertex, worldVertex)
//        worldVertex.vsub(planePosition, relpos)
//
//        var dot = worldNormal.dot(relpos)
//        if dot <= 0.0 {
//            if justTest {
//                return true
//            }
//
//            var r = this.createContactEquation(planeBody, convexBody, planeShape, convexShape, si, sj)
//
//            // Get vertex position projected on plane
//            var projected = planeConvex_projected
//            worldNormal.mult(worldNormal.dot(relpos),projected)
//            worldVertex.vsub(projected, projected)
//            projected.vsub(planePosition, r.ri) // From plane to vertex projected on plane
//
//            r.ni.copy(worldNormal) // Contact normal is the plane normal out from plane
//
//            // rj is now just the vector from the convex center to the vertex
//            worldVertex.vsub(convexPosition, r.rj)
//
//            // Make it relative to the body
//            r.ri.vadd(planePosition, r.ri)
//            r.ri.vsub(planeBody.position, r.ri)
//            r.rj.vadd(convexPosition, r.rj)
//            r.rj.vsub(convexBody.position, r.rj)
//
//            this.result.push(r)
//            numContacts++
//            if !this.enableFrictionReduction {
//                this.createFrictionEquationsFromContact(r, this.frictionResult)
//            }
//        }
//    }
//
//    if this.enableFrictionReduction && numContacts {
//        this.createFrictionFromAverage(numContacts)
//    }
//}