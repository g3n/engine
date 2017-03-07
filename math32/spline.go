// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math32

import (
//"math"
)

type Spline struct {
	points []Vector3
}

func NewSpline(points []Vector3) *Spline {

	this := new(Spline)
	this.points = make([]Vector3, len(points))
	copy(this.points, points)
	return this
}

func (this *Spline) InitFromArray(a []float32) {

	// PEND array of what ?
	//this.points = [];
	//for ( var i = 0; i < a.length; i ++ ) {
	//    this.points[ i ] = { x: a[ i ][ 0 ], y: a[ i ][ 1 ], z: a[ i ][ 2 ] };
	//}
}
