// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package physics

import (
	"github.com/g3n/engine/experimental/physics/object"
	"github.com/g3n/engine/experimental/physics/equation"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/experimental/collision/shape"
)

// Narrowphase
type Narrowphase struct {
	simulation              *Simulation
	currentContactMaterial  *ContactMaterial
	enableFrictionReduction bool // If true friction is computed as average
	debugging               bool
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

// createFrictionEquationsFromContact
func (n *Narrowphase) createFrictionEquationsFromContact(contactEquation *equation.Contact) (*equation.Friction, *equation.Friction) {

	bodyA := n.simulation.bodies[contactEquation.BodyA().Index()]
	bodyB := n.simulation.bodies[contactEquation.BodyB().Index()]

	// TODO materials
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

	// Construct and set tangents
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

// TODO test this
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

// GenerateEquations is the Narrowphase entry point.
func (n *Narrowphase) GenerateEquations(pairs []CollisionPair) ([]*equation.Contact, []*equation.Friction) {

	// TODO don't "make" every time, simply re-slice
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
			//contactEqs, frictionEqs := n.Resolve(bodyA, bodyB)
			contactEqs, frictionEqs := n.ResolveCollision(bodyA, bodyB)
			allContactEqs = append(allContactEqs, contactEqs...)
			allFrictionEqs = append(allFrictionEqs, frictionEqs...)
		}
   }

   return allContactEqs, allFrictionEqs
}

// ResolveCollision figures out which implementation of collision detection and contact resolution to use depending on the shapes involved.
func (n *Narrowphase) ResolveCollision(bodyA, bodyB *object.Body) ([]*equation.Contact, []*equation.Friction) {

	shapeA := bodyA.Shape()
	shapeB := bodyB.Shape()
	posA := bodyA.Position()
	posB := bodyB.Position()
	quatA := bodyA.Quaternion()
	quatB := bodyB.Quaternion()

	switch sA := shapeA.(type) {
	case *shape.Sphere:
		switch sB := shapeB.(type) {
		case *shape.Sphere:
			return n.SphereSphere(bodyA, bodyB, sA, sB, &posA, &posB, quatA, quatB)
		case *shape.Plane:
			return n.SpherePlane(bodyA, bodyB, sA, sB, &posA, &posB, quatA, quatB)
		case *shape.ConvexHull:
			return n.SphereConvex(bodyA, bodyB, sA, sB, &posA, &posB, quatA, quatB)
		}
	case *shape.Plane:
		switch sB := shapeB.(type) {
		case *shape.Sphere:
			return n.SpherePlane(bodyB, bodyA, sB, sA, &posB, &posA, quatB, quatA)
		//case *shape.Plane: // plane-plane collision never happens...
		//	return n.PlanePlane(bodyA, bodyB, sA, sB, &posA, &posB, quatA, quatB)
		case *shape.ConvexHull:
			return n.PlaneConvex(bodyA, bodyB, sA, sB, &posA, &posB, quatA, quatB)
		}
	case *shape.ConvexHull:
		switch sB := shapeB.(type) {
		case *shape.Sphere:
			return n.SphereConvex(bodyB, bodyA, sB, sA, &posB, &posA, quatB, quatA)
		case *shape.Plane:
			return n.PlaneConvex(bodyB, bodyA, sB, sA, &posB, &posA, quatB, quatA)
		case *shape.ConvexHull:
			return n.ConvexConvex(bodyA, bodyB, sA, sB, &posA, &posB, quatA, quatB)
		}
	}

	return []*equation.Contact{}, []*equation.Friction{}
}

// SphereSphere resolves the collision between two spheres analytically.
func (n *Narrowphase) SphereSphere(bodyA, bodyB *object.Body, sphereA, sphereB *shape.Sphere, posA, posB *math32.Vector3, quatA, quatB *math32.Quaternion) ([]*equation.Contact, []*equation.Friction) {

	contactEqs := make([]*equation.Contact, 0, 1)
	frictionEqs := make([]*equation.Friction, 0, 2)

	radiusA := sphereA.Radius()
	radiusB := sphereB.Radius()

	if posA.DistanceToSquared(posB) > math32.Pow(radiusA + radiusB, 2) {
		// No collision
		return contactEqs, frictionEqs
	}

	// Find penetration axis
	penAxis := posB.Clone().Sub(posA).Normalize()

	// Create contact equation
	contactEq := equation.NewContact(bodyA, bodyB, 0, 1e6)
	contactEq.SetSpookParams(1e6, 3, n.simulation.dt)
	//contactEq.SetEnabled(sphereA.CollisionResponse() && sphereB.CollisionResponse()) // TODO
	contactEq.SetNormal(penAxis.Clone())
	contactEq.SetRA(penAxis.Clone().MultiplyScalar(radiusB))
	contactEq.SetRB(penAxis.Clone().MultiplyScalar(-radiusB))
	contactEqs = append(contactEqs, contactEq)

	// Create friction equations
	fEq1, fEq2 := n.createFrictionEquationsFromContact(contactEq)
	frictionEqs = append(frictionEqs, fEq1, fEq2)

	return contactEqs, frictionEqs
}

// SpherePlane resolves the collision between a sphere and a plane analytically.
func (n *Narrowphase) SpherePlane(bodyA, bodyB *object.Body, sphereA *shape.Sphere, planeB *shape.Plane, posA, posB *math32.Vector3, quatA, quatB *math32.Quaternion) ([]*equation.Contact, []*equation.Friction) {

	contactEqs := make([]*equation.Contact, 0)
	frictionEqs := make([]*equation.Friction, 0)

	sphereRadius := sphereA.Radius()
	localNormal := planeB.Normal()
	normal := localNormal.Clone().ApplyQuaternion(quatB).Negate().Normalize()

	// Project down sphere on plane
	point_on_plane_to_sphere := math32.NewVec3().SubVectors(posA, posB)
	plane_to_sphere_ortho := normal.Clone().MultiplyScalar(normal.Dot(point_on_plane_to_sphere))

	if -point_on_plane_to_sphere.Dot(normal) <= sphereRadius {

		//if justTest {
		//	return true
		//}

		// We will have one contact in this case
		contactEq := equation.NewContact(bodyA, bodyB, 0, 1e6)
		contactEq.SetSpookParams(1e6, 3, n.simulation.dt)
		contactEq.SetNormal(normal) // Normalize() might not be needed
		contactEq.SetRA(normal.Clone().MultiplyScalar(sphereRadius)) // Vector from sphere center to contact point
		contactEq.SetRB(math32.NewVec3().SubVectors(point_on_plane_to_sphere, plane_to_sphere_ortho)) // The sphere position projected to plane
		contactEqs = append(contactEqs, contactEq)

		// Create friction equations
		fEq1, fEq2 := n.createFrictionEquationsFromContact(contactEq)
		frictionEqs = append(frictionEqs, fEq1, fEq2)
	}

	return contactEqs, frictionEqs
}

// TODO The second half of this method is untested!!!
func (n *Narrowphase) SphereConvex(bodyA, bodyB *object.Body, sphereA *shape.Sphere, convexB *shape.ConvexHull, posA, posB *math32.Vector3, quatA, quatB *math32.Quaternion) ([]*equation.Contact, []*equation.Friction) {

	contactEqs := make([]*equation.Contact, 0)
	frictionEqs := make([]*equation.Friction, 0)

	// TODO
	//v3pool := this.v3pool
	//convex_to_sphere := math32.NewVec3().SubVectors(posA, posB)
    //normals := sj.faceNormals
    //faces := sj.faces
    //verts := sj.vertices
    //R :=     si.radius
    //penetrating_sides := []

    // COMMENTED OUT
	// if(convex_to_sphere.norm2() > si.boundingSphereRadius + sj.boundingSphereRadius){
	//     return;
	// }
	sphereRadius := sphereA.Radius()

	// First check if any vertex of the convex hull is inside the sphere
	done := false
	convexB.Geometry.ReadVertices(func(vertex math32.Vector3) bool {
		worldVertex := vertex.ApplyQuaternion(quatA).Add(posB)
		sphereToCorner := math32.NewVec3().SubVectors(worldVertex, posA)
		if sphereToCorner.LengthSq() < sphereRadius * sphereRadius {
			// Colliding! worldVertex is inside sphere.

			// Create contact equation
			contactEq := equation.NewContact(bodyA, bodyB, 0, 1e6)
			contactEq.SetSpookParams(1e6, 3, n.simulation.dt)
			//contactEq.SetEnabled(sphereA.CollisionResponse() && sphereB.CollisionResponse()) // TODO
			normalizedSphereToCorner := sphereToCorner.Clone().Normalize()
			contactEq.SetNormal(normalizedSphereToCorner)
			contactEq.SetRA(normalizedSphereToCorner.Clone().MultiplyScalar(sphereRadius))
			contactEq.SetRB(worldVertex.Clone().Sub(posB))
			contactEqs = append(contactEqs, contactEq)
			// Create friction equations
			fEq1, fEq2 := n.createFrictionEquationsFromContact(contactEq)
			frictionEqs = append(frictionEqs, fEq1, fEq2)
			// Set done flag
			done = true
			// Break out of loop
			return true
		}
		return false
	})
	if done {
		return contactEqs, frictionEqs
	}

	//Check side (plane) intersections TODO NOTE THIS IS UNTESTED
    convexFaces := convexB.Faces()
    convexWorldFaceNormals := convexB.WorldFaceNormals()
    for i := 0; i < len(convexFaces); i++ {
		worldNormal := convexWorldFaceNormals[i]
     	face := convexFaces[i]
     	// Get a world vertex from the face
     	var worldPoint = face[0].Clone().ApplyQuaternion(quatB).Add(posB)
     	// Get a point on the sphere, closest to the face normal
     	var worldSpherePointClosestToPlane = worldNormal.Clone().MultiplyScalar(-sphereRadius).Add(posA)
     	// Vector from a face point to the closest point on the sphere
     	var penetrationVec = math32.NewVec3().SubVectors(worldSpherePointClosestToPlane, worldPoint)
     	// The penetration. Negative value means overlap.
     	var penetration = penetrationVec.Dot(&worldNormal)
     	var worldPointToSphere = math32.NewVec3().SubVectors(posA, worldPoint)
     	if penetration < 0 && worldPointToSphere.Dot(&worldNormal) > 0 {
         	// Intersects plane. Now check if the sphere is inside the face polygon
         	worldFace := convexB.WorldFace(face, posB, quatB)
         	if n.pointBehindFace(worldFace, &worldNormal, posA) { // Is the sphere center behind the face (inside the convex polygon?
				// TODO NEVER GETTING INSIDE THIS IF STATEMENT!
				ShowWorldFace(n.simulation.Scene(), worldFace[:], &math32.Color{0,0,2})

				// if justTest {
             	//    return true
             	//}

				// Create contact equation
				contactEq := equation.NewContact(bodyA, bodyB, 0, 1e6)
				contactEq.SetSpookParams(1e6, 3, n.simulation.dt)
				//contactEq.SetEnabled(sphereA.CollisionResponse() && sphereB.CollisionResponse()) // TODO
				contactEq.SetNormal(worldNormal.Clone().Negate())
				contactEq.SetRA(worldNormal.Clone().MultiplyScalar(-sphereRadius))
				penetrationVec2 := worldNormal.Clone().MultiplyScalar(-penetration)
				penetrationSpherePoint := worldNormal.Clone().MultiplyScalar(-sphereRadius)
				contactEq.SetRB(posA.Clone().Sub(posB).Add(penetrationSpherePoint).Add(penetrationVec2))
				contactEqs = append(contactEqs, contactEq)
				// Create friction equations
				fEq1, fEq2 := n.createFrictionEquationsFromContact(contactEq)
				frictionEqs = append(frictionEqs, fEq1, fEq2)
				// Exit method (we only expect *one* face contact)
             	return contactEqs, frictionEqs
         	} else {
				// Edge?
             	for j := 0; j < len(worldFace); j++ {
             		// Get two world transformed vertices
                	v1 := worldFace[(j+1)%3].Clone()//.ApplyQuaternion(quatB).Add(posB)
                	v2 := worldFace[(j+2)%3].Clone()//.ApplyQuaternion(quatB).Add(posB)
                	// Construct edge vector
                	edge := math32.NewVec3().SubVectors(v2, v1)
                	// Construct the same vector, but normalized
                	edgeUnit := edge.Clone().Normalize()
                	// p is xi projected onto the edge
                	v1ToPosA := math32.NewVec3().SubVectors(posA, v1)
                	dot := v1ToPosA.Dot(edgeUnit)
					p := edgeUnit.Clone().MultiplyScalar(dot).Add(v1)
                	// Compute a vector from p to the center of the sphere
                	var posAtoP = math32.NewVec3().SubVectors(p, posA)
                	// Collision if the edge-sphere distance is less than the radius AND if p is in between v1 and v2
                	edgeL2 := edge.LengthSq()
                	patp2 := posAtoP.LengthSq()
                	if (dot > 0) && (dot*dot < edgeL2) && (patp2 < sphereRadius*sphereRadius) { // Collision if the edge-sphere distance is less than the radius
                	   // Edge contact!
                	   //if justTest {
                	   //    return true
                	   //}
						// Create contact equation
						contactEq := equation.NewContact(bodyA, bodyB, 0, 1e6)
						contactEq.SetSpookParams(1e6, 3, n.simulation.dt)
						//contactEq.SetEnabled(sphereA.CollisionResponse() && sphereB.CollisionResponse()) // TODO
						normal := p.Clone().Sub(posA).Normalize()
						contactEq.SetNormal(normal)
						contactEq.SetRA(normal.Clone().MultiplyScalar(sphereRadius))
						contactEq.SetRB(p.Clone().Sub(posB))
						contactEqs = append(contactEqs, contactEq)
						// Create friction equations
						//fEq1, fEq2 := n.createFrictionEquationsFromContact(contactEq)
						//frictionEqs = append(frictionEqs, fEq1, fEq2)
						// Exit method (we only expect *one* edge contact)
						return contactEqs, frictionEqs
                	}
             	}
			}
     	}
    }

	return contactEqs, frictionEqs
}

func (n *Narrowphase) pointBehindFace(worldFace [3]math32.Vector3, faceNormal, point *math32.Vector3) bool {

	pointInFace := worldFace[0].Clone()
	planeDelta := -pointInFace.Dot(faceNormal)
	depth := faceNormal.Dot(point) + planeDelta
	if depth > 0 {
		return false
	} else {
		return true
	}
}

// Checks if a given point is inside the polygon. I believe this is a generalization of checking whether a point is behind a face/plane when the face has more than 3 vertices.
//func (n *Narrowphase) pointInPolygon(verts []math32.Vector3, normal, p *math32.Vector3) bool {
//
//	firstTime := true
//    positiveResult := true // first value of positive result doesn't matter
//    N := len(verts)
//    for i := 0; i < N; i++ {
//        v := verts[i].Clone()
//        // Get edge to the next vertex
//        edge := math32.NewVec3().SubVectors(verts[(i+1)%N].Clone(), v)
//        // Get cross product between polygon normal and the edge
//        edgeCrossNormal := math32.NewVec3().CrossVectors(edge, normal)
//        // Get vector between point and current vertex
//        vertexToP := math32.NewVec3().SubVectors(p, v)
//        // This dot product determines which side of the edge the point is
//        r := edgeCrossNormal.Dot(vertexToP)
//        // If all such dot products have same sign, we are inside the polygon.
//        if firstTime || (r > 0 && positiveResult == true) || (r <= 0 && positiveResult == false) {
//            if firstTime {
//                positiveResult = r > 0
//				firstTime = false
//            }
//            continue
//        } else {
//            return false // Encountered some other sign. Exit.
//        }
//    }
//    // If we got here, all dot products were of the same sign.
//    return true
//}

// ConvexConvex implements collision detection and contact resolution between two convex hulls.
func (n *Narrowphase) ConvexConvex(bodyA, bodyB *object.Body, convexA, convexB *shape.ConvexHull, posA, posB *math32.Vector3, quatA, quatB *math32.Quaternion) ([]*equation.Contact, []*equation.Friction) {

	contactEqs := make([]*equation.Contact, 0)
	frictionEqs := make([]*equation.Friction, 0)

	// Check if colliding and find penetration axis
	penetrating, penAxis := convexA.FindPenetrationAxis(convexB, posA, posB, quatA, quatB)
	if !penetrating {
		return contactEqs, frictionEqs
	}

	if n.debugging {
		ShowPenAxis(n.simulation.Scene(), &penAxis) //, -1000, 1000)
		log.Error("Colliding (%v|%v) penAxis: %v", bodyA.Name(), bodyB.Name(), penAxis)
	}

	// Colliding! Find contacts.
	contacts := convexA.ClipAgainstHull(convexB, posA, posB, quatA, quatB, &penAxis, -100, 100)

	// For each contact found create a contact equation and the two associated friction equations
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
		contactEq.SetRA(contact.Normal.Clone().MultiplyScalar(-contact.Depth).Add(&contact.Point).Sub(posA))
		contactEq.SetRB(contact.Point.Clone().Sub(posB))
		contactEqs = append(contactEqs, contactEq)

		// If enableFrictionReduction is true then skip creating friction equations for individual contacts
		// We will create average friction equations later based on all contacts
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

	return contactEqs, frictionEqs
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


// Checks whether p is inside the polyhedra. Must be in local coords.
// The point lies outside of the convex hull of the other points if and only if
// the direction of all the vectors from it to those other points are on less than one half of a sphere around it.
// p is A point given in local coordinates
//func (n *Narrowphase) PointIsInside(p) {
//
//	verts := this.vertices
//	faces := this.faces
//	normals := this.faceNormals
//  positiveResult := null
//  N := this.faces.length
//  pointInside := ConvexPolyhedron_pointIsInside
//  this.getAveragePointLocal(pointInside)
//  for i := 0; i < N; i++ {
//      numVertices := this.faces[i].length
//      n := normals[i]
//      v := verts[faces[i][0]] // We only need one point in the face
//
//      // This dot product determines which side of the edge the point is
//      vToP := ConvexPolyhedron_vToP
//      p.vsub(v,vToP)
//      r1 := n.dot(vToP)
//
//      vToPointInside := ConvexPolyhedron_vToPointInside
//      pointInside.vsub(v,vToPointInside)
//      r2 := n.dot(vToPointInside)
//
//      if (r1<0 && r2>0) || (r1>0 && r2<0) {
//          return false // Encountered some other sign. Exit.
//      }
//  }
//
//  // If we got here, all dot products were of the same sign.
//  return positiveResult ? 1 : -1
//}

// TODO
func (n *Narrowphase) PlaneConvex(bodyA, bodyB *object.Body, planeA *shape.Plane, convexB *shape.ConvexHull, posA, posB *math32.Vector3, quatA, quatB *math32.Quaternion) ([]*equation.Contact, []*equation.Friction) {

	contactEqs := make([]*equation.Contact, 0)
	frictionEqs := make([]*equation.Friction, 0)

	//planeShape,
	//convexShape,
	//planePosition,
	//convexPosition,
	//planeQuat,
	//convexQuat,
	//planeBody,
	//convexBody,
	//si,
	//sj,
	//justTest) {
   //
	//// Simply return the points behind the plane.
   //worldVertex := planeConvex_v
   //worldNormal := planeConvex_normal
   //worldNormal.set(0,0,1)
   //planeQuat.vmult(worldNormal,worldNormal) // Turn normal according to plane orientation
   //
   //var numContacts = 0
   //var relpos = planeConvex_relpos
   //for i := 0; i < len(convexShape.vertices); i++ {
   //
   //    // Get world convex vertex
   //    worldVertex.copy(convexShape.vertices[i])
   //    convexQuat.vmult(worldVertex, worldVertex)
   //    convexPosition.vadd(worldVertex, worldVertex)
   //    worldVertex.vsub(planePosition, relpos)
   //
   //    var dot = worldNormal.dot(relpos)
   //    if dot <= 0.0 {
   //        if justTest {
   //            return true
   //        }
   //
   //        var r = this.createContactEquation(planeBody, convexBody, planeShape, convexShape, si, sj)
   //
   //        // Get vertex position projected on plane
   //        var projected = planeConvex_projected
   //        worldNormal.mult(worldNormal.dot(relpos),projected)
   //        worldVertex.vsub(projected, projected)
   //        projected.vsub(planePosition, r.ri) // From plane to vertex projected on plane
   //
   //        r.ni.copy(worldNormal) // Contact normal is the plane normal out from plane
   //
   //        // rj is now just the vector from the convex center to the vertex
   //        worldVertex.vsub(convexPosition, r.rj)
   //
   //        // Make it relative to the body
   //        r.ri.vadd(planePosition, r.ri)
   //        r.ri.vsub(planeBody.position, r.ri)
   //        r.rj.vadd(convexPosition, r.rj)
   //        r.rj.vsub(convexBody.position, r.rj)
   //
   //        this.result.push(r)
   //        numContacts++
   //        if !this.enableFrictionReduction {
   //            this.createFrictionEquationsFromContact(r, this.frictionResult)
   //        }
   //    }
   //}
   //
   //if this.enableFrictionReduction && numContacts {
   //    this.createFrictionFromAverage(numContacts)
   //}

   return contactEqs, frictionEqs
}