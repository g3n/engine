// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math32

// Ray represents an oriented 3D line segment defined by an origin point and a direction vector.
type Ray struct {
	origin    Vector3
	direction Vector3
}

// NewRay creates and returns a pointer to a Ray object with
// the specified origin and direction vectors.
// If a nil pointer is supplied for any of the parameters,
// the zero vector will be used.
func NewRay(origin *Vector3, direction *Vector3) *Ray {

	ray := new(Ray)
	if origin != nil {
		ray.origin = *origin
	}
	if direction != nil {
		ray.direction = *direction
	}
	return ray
}

// Set sets the origin and direction vectors of this Ray.
func (ray *Ray) Set(origin, direction *Vector3) *Ray {

	ray.origin = *origin
	ray.direction = *direction
	return ray
}

// Copy copies other ray into this one.
func (ray *Ray) Copy(other *Ray) *Ray {

	*ray = *other
	return ray
}

// Origin returns a copy of this ray current origin.
func (ray *Ray) Origin() Vector3 {

	return ray.origin
}

// Direction returns a copy of this ray current direction.
func (ray *Ray) Direction() Vector3 {

	return ray.direction
}

// At calculates the point in the ray which is at the specified t distance from the origin
// along its direction.
// The calculated point is stored in optionalTarget, if not nil, and also returned.
func (ray *Ray) At(t float32, optionalTarget *Vector3) *Vector3 {

	var result *Vector3
	if optionalTarget != nil {
		result = optionalTarget
	} else {
		result = &Vector3{}
	}
	return result.Copy(&ray.direction).MultiplyScalar(t).Add(&ray.origin)
}

// Recast sets the new origin of the ray at the specified distance t
// from its origin along its direction.
func (ray *Ray) Recast(t float32) *Ray {

	var v1 Vector3
	ray.origin.Copy(ray.At(t, &v1))
	return ray
}

// ClosestPointToPoint calculates the point in the ray which is closest to the specified point.
// The calculated point is stored in optionalTarget, if not nil, and also returned.
func (ray *Ray) ClosestPointToPoint(point, optionalTarget *Vector3) *Vector3 {

	var result *Vector3
	if optionalTarget != nil {
		result = optionalTarget
	} else {
		result = NewVector3(0, 0, 0)
	}
	result.SubVectors(point, &ray.origin)
	directionDistance := result.Dot(&ray.direction)

	if directionDistance < 0 {
		return result.Copy(&ray.origin)
	}
	return result.Copy(&ray.direction).MultiplyScalar(directionDistance).Add(&ray.origin)
}

// DistanceToPoint returns the smallest distance
// from the ray direction vector to the specified point.
func (ray *Ray) DistanceToPoint(point *Vector3) float32 {

	return Sqrt(ray.DistanceSqToPoint(point))
}

// DistanceSqToPoint returns the smallest squared distance
// from the ray direction vector to the specified point.
// If the ray was pointed directly at the point this distance would be 0.
func (ray *Ray) DistanceSqToPoint(point *Vector3) float32 {

	var v1 Vector3

	directionDistance := v1.SubVectors(point, &ray.origin).Dot(&ray.direction)
	// point behind the ray
	if directionDistance < 0 {
		return ray.origin.DistanceTo(point)
	}
	v1.Copy(&ray.direction).MultiplyScalar(directionDistance).Add(&ray.origin)
	return v1.DistanceToSquared(point)
}

// DistanceSqToSegment returns the smallest squared distance
// from this ray to the line segment from v0 to v1.
// If optionalPointOnRay Vector3 is not nil,
// it is set with the coordinates of the point on the ray.
// if optionalPointOnSegment Vector3 is not nil,
// it is set with the coordinates of the point on the segment.
func (ray *Ray) DistanceSqToSegment(v0, v1, optionalPointOnRay, optionalPointOnSegment *Vector3) float32 {

	var segCenter Vector3
	var segDir Vector3
	var diff Vector3

	segCenter.Copy(v0).Add(v1).MultiplyScalar(0.5)
	segDir.Copy(v1).Sub(v0).Normalize()
	diff.Copy(&ray.origin).Sub(&segCenter)

	segExtent := v0.DistanceTo(v1) * 0.5
	a01 := -ray.direction.Dot(&segDir)
	b0 := diff.Dot(&ray.direction)
	b1 := -diff.Dot(&segDir)
	c := diff.LengthSq()
	det := Abs(1 - a01*a01)

	var s0, s1, sqrDist, extDet float32

	if det > 0 {

		// The ray and segment are not parallel.
		s0 = a01*b1 - b0
		s1 = a01*b0 - b1
		extDet = segExtent * det

		if s0 >= 0 {

			if s1 >= -extDet {

				if s1 <= extDet {
					// region 0
					// Minimum at interior points of ray and segment.
					invDet := 1 / det
					s0 *= invDet
					s1 *= invDet
					sqrDist = s0*(s0+a01*s1+2*b0) + s1*(a01*s0+s1+2*b1) + c

				} else {
					// region 1
					s1 = segExtent
					s0 = Max(0, -(a01*s1 + b0))
					sqrDist = -s0*s0 + s1*(s1+2*b1) + c
				}

			} else {
				// region 5
				s1 = -segExtent
				s0 = Max(0, -(a01*s1 + b0))
				sqrDist = -s0*s0 + s1*(s1+2*b1) + c

			}

		} else {

			if s1 <= -extDet {
				// region 4
				s0 = Max(0, -(-a01*segExtent + b0))
				if s0 > 0 {
					s1 = -segExtent
				} else {
					s1 = Min(Max(-segExtent, -b1), segExtent)
				}
				sqrDist = -s0*s0 + s1*(s1+2*b1) + c

			} else if s1 <= extDet {
				// region 3
				s0 = 0
				s1 = Min(Max(-segExtent, -b1), segExtent)
				sqrDist = s1*(s1+2*b1) + c

			} else {
				// region 2
				s0 = Max(0, -(a01*segExtent + b0))
				if s0 > 0 {
					s1 = segExtent
				} else {
					s1 = Min(Max(-segExtent, -b1), segExtent)
				}
				sqrDist = -s0*s0 + s1*(s1+2*b1) + c
			}
		}
	} else {

		// Ray and segment are parallel.
		if a01 > 0 {
			s1 = -segExtent
		} else {
			s1 = segExtent
		}
		s0 = Max(0, -(a01*s1 + b0))
		sqrDist = -s0*s0 + s1*(s1+2*b1) + c

	}

	if optionalPointOnRay != nil {
		optionalPointOnRay.Copy(&ray.direction).MultiplyScalar(s0).Add(&ray.origin)
	}

	if optionalPointOnSegment != nil {
		optionalPointOnSegment.Copy(&segDir).MultiplyScalar(s1).Add(&segCenter)
	}
	return sqrDist
}

// IsIntersectionSphere returns if this ray intersects with the specified sphere.
func (ray *Ray) IsIntersectionSphere(sphere *Sphere) bool {

	if ray.DistanceToPoint(&sphere.Center) <= sphere.Radius {
		return true
	}
	return false
}

// IntersectSphere calculates the point which is the intersection of this ray with the specified sphere.
// The calculated point is stored in optionalTarget, it not nil, and also returned.
// If no intersection is found the calculated point is set to nil.
func (ray *Ray) IntersectSphere(sphere *Sphere, optionalTarget *Vector3) *Vector3 {

	var v1 Vector3

	v1.SubVectors(&sphere.Center, &ray.origin)
	tca := v1.Dot(&ray.direction)
	d2 := v1.Dot(&v1) - tca*tca
	radius2 := sphere.Radius * sphere.Radius

	if d2 > radius2 {
		return nil
	}

	thc := Sqrt(radius2 - d2)

	// t0 = first intersect point - entrance on front of sphere
	t0 := tca - thc

	// t1 = second intersect point - exit point on back of sphere
	t1 := tca + thc

	// test to see if both t0 and t1 are behind the ray - if so, return null
	if t0 < 0 && t1 < 0 {
		return nil
	}

	// test to see if t0 is behind the ray:
	// if it is, the ray is inside the sphere, so return the second exit point scaled by t1,
	// in order to always return an intersect point that is in front of the ray.
	if t0 < 0 {
		return ray.At(t1, optionalTarget)
	}

	// else t0 is in front of the ray, so return the first collision point scaled by t0
	return ray.At(t0, optionalTarget)
}

// IsIntersectPlane returns if this ray intersects the specified plane.
func (ray *Ray) IsIntersectPlane(plane *Plane) bool {

	distToPoint := plane.DistanceToPoint(&ray.origin)
	if distToPoint == 0 {
		return true
	}

	denominator := plane.normal.Dot(&ray.direction)
	if denominator*distToPoint < 0 {
		return true
	}

	// ray origin is behind the plane (and is pointing behind it)
	return false
}

// DistanceToPlane returns the distance of this ray origin to its intersection point in the plane.
// If the ray does not intersects the plane, returns NaN.
func (ray *Ray) DistanceToPlane(plane *Plane) float32 {

	denominator := plane.normal.Dot(&ray.direction)
	if denominator == 0 {
		// line is coplanar, return origin
		if plane.DistanceToPoint(&ray.origin) == 0 {
			return 0
		}
		return NaN()
	}
	t := -(ray.origin.Dot(&plane.normal) + plane.constant) / denominator
	// Return if the ray never intersects the plane
	if t >= 0 {
		return t
	}
	return NaN()
}

// IntersectPlane calculates the point which is the intersection of this ray with the specified plane.
// The calculated point is stored in optionalTarget, if not nil, and also returned.
// If no intersection is found the calculated point is set to nil.
func (ray *Ray) IntersectPlane(plane *Plane, optionalTarget *Vector3) *Vector3 {

	t := ray.DistanceToPlane(plane)

	if t == NaN() {
		return nil
	}

	return ray.At(t, optionalTarget)

}

// IsIntersectionBox returns if this ray intersects the specified box.
func (ray *Ray) IsIntersectionBox(box *Box3) bool {

	var v Vector3

	if ray.IntersectBox(box, &v) != nil {
		return true
	}
	return false
}

// IntersectBox calculates the point which is the intersection of this ray with the specified box.
// The calculated point is stored in optionalTarget, it not nil, and also returned.
// If no intersection is found the calculated point is set to nil.
func (ray *Ray) IntersectBox(box *Box3, optionalTarget *Vector3) *Vector3 {

	// http://www.scratchapixel.com/lessons/3d-basic-lessons/lesson-7-intersecting-simple-shapes/ray-box-intersection/

	var tmin, tmax, tymin, tymax, tzmin, tzmax float32

	invdirx := 1 / ray.direction.X
	invdiry := 1 / ray.direction.Y
	invdirz := 1 / ray.direction.Z

	var origin = ray.origin

	if invdirx >= 0 {
		tmin = (box.Min.X - origin.X) * invdirx
		tmax = (box.Max.X - origin.X) * invdirx
	} else {
		tmin = (box.Max.X - origin.X) * invdirx
		tmax = (box.Min.X - origin.X) * invdirx
	}

	if invdiry >= 0 {
		tymin = (box.Min.Y - origin.Y) * invdiry
		tymax = (box.Max.Y - origin.Y) * invdiry
	} else {
		tymin = (box.Max.Y - origin.Y) * invdiry
		tymax = (box.Min.Y - origin.Y) * invdiry
	}

	if (tmin > tymax) || (tymin > tmax) {
		return nil
	}

	// These lines also handle the case where tmin or tmax is NaN
	// (result of 0 * Infinity). x !== x returns true if x is NaN

	if tymin > tmin || tmin != tmin {
		tmin = tymin
	}

	if tymax < tmax || tmax != tmax {
		tmax = tymax
	}

	if invdirz >= 0 {
		tzmin = (box.Min.Z - origin.Z) * invdirz
		tzmax = (box.Max.Z - origin.Z) * invdirz
	} else {
		tzmin = (box.Max.Z - origin.Z) * invdirz
		tzmax = (box.Min.Z - origin.Z) * invdirz
	}

	if (tmin > tzmax) || (tzmin > tmax) {
		return nil
	}

	if tzmin > tmin || tmin != tmin {
		tmin = tzmin
	}

	if tzmax < tmax || tmax != tmax {
		tmax = tzmax
	}

	//return point closest to the ray (positive side)

	if tmax < 0 {
		return nil
	}

	if tmin >= 0 {
		return ray.At(tmin, optionalTarget)
	}
	return ray.At(tmax, optionalTarget)
}

// IntersectTriangle returns if this ray intersects the triangle with the face
// defined by points a, b, c. Returns true if it intersects and sets the point
// parameter with the intersected point coordinates.
// If backfaceCulling is false it ignores the intersection if the face is not oriented
// in the ray direction.
func (ray *Ray) IntersectTriangle(a, b, c *Vector3, backfaceCulling bool, point *Vector3) bool {

	var diff Vector3
	var edge1 Vector3
	var edge2 Vector3
	var normal Vector3

	edge1.SubVectors(b, a)
	edge2.SubVectors(c, a)
	normal.CrossVectors(&edge1, &edge2)

	// Solve Q + t*D = b1*E1 + b2*E2 (Q = kDiff, D = ray direction,
	// E1 = kEdge1, E2 = kEdge2, N = Cross(E1,E2)) by
	//   |Dot(D,N)|*b1 = sign(Dot(D,N))*Dot(D,Cross(Q,E2))
	//   |Dot(D,N)|*b2 = sign(Dot(D,N))*Dot(D,Cross(E1,Q))
	//   |Dot(D,N)|*t = -sign(Dot(D,N))*Dot(Q,N)
	DdN := ray.direction.Dot(&normal)
	var sign float32

	if DdN > 0 {
		if backfaceCulling {
			return false
		}
		sign = 1
	} else if DdN < 0 {
		sign = -1
		DdN = -DdN
	} else {
		return false
	}

	diff.SubVectors(&ray.origin, a)
	DdQxE2 := sign * ray.direction.Dot(edge2.CrossVectors(&diff, &edge2))

	// b1 < 0, no intersection
	if DdQxE2 < 0 {
		return false
	}

	DdE1xQ := sign * ray.direction.Dot(edge1.Cross(&diff))
	// b2 < 0, no intersection
	if DdE1xQ < 0 {
		return false
	}

	// b1+b2 > 1, no intersection
	if DdQxE2+DdE1xQ > DdN {
		return false
	}

	// Line intersects triangle, check if ray does.
	QdN := -sign * diff.Dot(&normal)

	// t < 0, no intersection
	if QdN < 0 {
		return false
	}

	// Ray intersects triangle.
	ray.At(QdN/DdN, point)
	return true
}

// ApplyMatrix4 multiplies this ray origin and direction
// by the specified matrix4, basically transforming this ray coordinates.
func (ray *Ray) ApplyMatrix4(matrix4 *Matrix4) *Ray {

	ray.direction.Add(&ray.origin).ApplyMatrix4(matrix4)
	ray.origin.ApplyMatrix4(matrix4)
	ray.direction.Sub(&ray.origin)
	ray.direction.Normalize()
	return ray
}

// Equals returns if this ray is equal to other
func (ray *Ray) Equals(other *Ray) bool {

	return ray.origin.Equals(&other.origin) && ray.direction.Equals(&other.direction)
}

// Clone creates and returns a pointer to copy of this ray.
func (ray *Ray) Clone() *Ray {

	return NewRay(&ray.origin, &ray.direction)
}
