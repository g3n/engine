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

func NewContactMaterial() *ContactMaterial {

	cm := new(ContactMaterial)
	cm.friction = 0.3
	cm.restitution = 0.3
	cm.contactEquationStiffness = 1e7
	cm.contactEquationRelaxation = 3
	cm.frictionEquationStiffness = 1e7
	cm.frictionEquationRelaxation = 3
	return cm
}


//type intPair struct {
//	i int
//	j int
//}

//type ContactMaterialTable map[intPair]*ContactMaterial