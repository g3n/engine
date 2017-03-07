// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package shader

func init() {
	AddChunk("attributes", chunkAttributes)
}

const chunkAttributes = `
// Vertex attributes
in layout(location = 0) vec3  VertexPosition;
in layout(location = 1) vec3  VertexNormal;
in layout(location = 2) vec3  VertexColor;
in layout(location = 3) vec2  VertexTexcoord;
in layout(location = 4) float VertexDistance;
in layout(location = 5) vec4  VertexTexoffsets;
`
