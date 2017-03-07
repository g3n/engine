// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package shader

func init() {
	AddShader("shaderPointVertex", shaderPointVertex)
	AddShader("shaderPointFrag", shaderPointFrag)
	AddProgram("shaderPoint", "shaderPointVertex", "shaderPointFrag")
}

//
// Vertex Shader template
//
const shaderPointVertex = `
#version {{.Version}}

{{template "attributes" .}}

// Model uniforms
uniform mat4 MVP;

// Material uniforms
{{template "material" .}}

// Specific uniforms
uniform float PointSize;
uniform float RotationZ;

// Outputs for fragment shader
out vec3 Color;
flat out mat2 Rotation;

void main() {

    // Rotation matrix for fragment shader
    float rotSin = sin(RotationZ);
    float rotCos = cos(RotationZ);
    Rotation = mat2(rotCos, rotSin, - rotSin, rotCos);

    // Sets the vertex position
    vec4 pos = MVP * vec4(VertexPosition, 1.0);
    gl_Position = pos;

    // Sets the size of the rasterized point decreasing with distance
    gl_PointSize = (1.0 - pos.z / pos.w) * PointSize;

    // Outputs color
    Color = MatEmissiveColor;
}
`

//
// Fragment Shader template
//
const shaderPointFrag = `
#version {{.Version}}

{{template "material" .}}

// Inputs from vertex shader
in vec3 Color;
flat in mat2 Rotation;

// Output
out vec4 FragColor;

void main() {

    // Combine all texture colors and opacity
    vec4 texCombined = vec4(1);
    {{if .MatTexturesMax}}
    for (int i = 0; i < {{.MatTexturesMax}}; i++) {
        vec2 pt = gl_PointCoord - vec2(0.5);
        vec4 texcolor = texture(MatTexture[i], (Rotation * pt + vec2(0.5)) * MatTexRepeat[i] + MatTexOffset[i]);
        if (i == 0) {
            texCombined = texcolor;
        } else {
            texCombined = mix(texCombined, texcolor, texcolor.a);
        }
    }
    {{ end }}

    // Combine material color with texture
    FragColor = min(vec4(Color, MatOpacity) * texCombined, vec4(1));
}

`
