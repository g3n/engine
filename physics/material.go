// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package physics implements a basic physics engine.
package physics

type Material struct {
	name        string
	friction    float32
	restitution float32
}

type ContactMaterial struct {
	mat1                       *Material
	mat2                       *Material
	friction                   float32
	restitution                float32
	contactEquationStiffness   float32
	contactEquationRelaxation  float32
	frictionEquationStiffness  float32
	frictionEquationRelaxation float32
}
