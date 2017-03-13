// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package shader

func init() {
	AddChunk("attributes", chunkAttributes)
}

const chunkAttributes = `
// Vertex attributes
layout(location = 0) in vec3  VertexPosition;
layout(location = 1) in vec3  VertexNormal;
layout(location = 2) in vec3  VertexColor;
layout(location = 3) in vec2  VertexTexcoord;
layout(location = 4) in float VertexDistance;
layout(location = 5) in vec4  VertexTexoffsets;
`
