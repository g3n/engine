// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package shader

func init() {
	AddShader("shaderScreenSpaceLineVertex", shaderScreenSpaceLineVertex)
	AddShader("shaderScreenSpaceLineFrag", shaderScreenSpaceLineFrag)
	AddShader("shaderScreenSpaceLineGeometry", shaderScreenSpaceLineGeometry)

	AddProgram("shaderScreenSpaceLine",
		"shaderScreenSpaceLineVertex",
		"shaderScreenSpaceLineFrag",
		"shaderScreenSpaceLineGeometry")
}

// Vertex Shader template
const shaderScreenSpaceLineVertex = `
#version {{.Version}}

{{template "attributes" .}}
{{template "material" .}}

// Model uniforms
uniform mat4 MVP;
uniform vec2 viewportSize;
uniform float thickness;

// Output for geometry shader
out vec3 o_color;
out vec2 o_viewportSize;
out float o_thickness; 

void main() {

    o_color = VertexColor;
    o_viewportSize = viewportSize;
    o_thickness = thickness;
    
    gl_Position = MVP * vec4(VertexPosition, 1.0);
}
`

// Geometry Shader template
const shaderScreenSpaceLineGeometry = `
#version {{.Version}}

layout(lines) in;
layout(triangle_strip, max_vertices = 4) out;

in vec3  o_color[2];
in vec2  o_viewportSize[2];
in float o_thickness[2];

out vec3 Color;

void main() {

    vec4 ndc0 = gl_in[0].gl_Position;
    vec4 ndc1 = gl_in[1].gl_Position;

		// calculate the line normal (A - B)
    vec2 direction = normalize( ndc1.xy - ndc0.xy );
    vec2 normal = vec2( -direction.y, direction.x );

    vec2 viewportSize = o_viewportSize[0];
    float thickness = o_thickness[0]; 

		// extrude & correct aspect ratio
    vec2 offset = ( thickness / viewportSize) * normal;

    gl_Position = vec4(ndc0.xy + offset * ndc0.w, ndc0.z, ndc0.w);
    Color = o_color[0];
    EmitVertex();

    gl_Position = vec4(ndc0.xy - offset * ndc0.w, ndc0.z, ndc0.w);
    Color = o_color[0];
    EmitVertex();

    gl_Position = vec4(ndc1.xy + offset * ndc1.w, ndc1.z, ndc1.w);
    Color = o_color[1];
    EmitVertex();

    gl_Position = vec4(ndc1.xy - offset * ndc1.w, ndc1.z, ndc1.w);
    Color = o_color[1];
    EmitVertex();

    EndPrimitive();
}
`

// Fragment Shader template
const shaderScreenSpaceLineFrag = `
#version {{.Version}}

in vec3 Color;
out vec4 FragColor;

void main() {
  FragColor = vec4(Color, 1.0);
}
`
