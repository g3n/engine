// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package shader

func init() {
	AddShader("shaderBasicVertex", shaderBasicVertex)
	AddShader("shaderBasicFrag", shaderBasicFrag)
	AddProgram("shaderBasic", "shaderBasicVertex", "shaderBasicFrag")
}

//
// Vertex Shader template
//
const shaderBasicVertex = `
#version {{.Version}}

{{template "attributes" .}}
{{template "material" .}}

// Model uniforms
uniform mat4 MVP;


// Final output color for fragment shader
out vec3 Color;

void main() {

    Color = VertexColor;
    gl_Position = MVP * vec4(VertexPosition, 1.0);
}
`

//
// Fragment Shader template
//
const shaderBasicFrag = `
#version {{.Version}}

in vec3 Color;
out vec4 FragColor;

void main() {

    FragColor = vec4(Color, 1.0);
}

`
