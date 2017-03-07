// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package shader

const shaderSpriteVertex = `
#version {{.Version}}

{{template "attributes" .}}

// Input uniforms
uniform mat4 MVP;

{{template "material" .}}

// Outputs for fragment shader
out vec3 Color;
out vec2 FragTexcoord;

void main() {

    // Applies transformation to vertex position
    gl_Position = MVP * vec4(VertexPosition, 1.0);

    // Outputs color
    Color = MatDiffuseColor;

    // Flips texture coordinate Y if requested.
    vec2 texcoord = VertexTexcoord;
    {{if .MatTexturesMax}}
    if (MatTexFlipY[0] > 0) {
        texcoord.y = 1 - texcoord.y;
    }
    {{ end }}
    FragTexcoord = texcoord;
}
`

//
// Fragment Shader template
//
const shaderSpriteFrag = `
#version {{.Version}}

{{template "material" .}}

// Inputs from vertex shader
in vec3 Color;
in vec2 FragTexcoord;

// Output
out vec4 FragColor;

void main() {

    // Combine all texture colors and opacity
    vec4 texCombined = vec4(1);
    {{if .MatTexturesMax }}
    for (int i = 0; i < {{.MatTexturesMax}}; i++) {
        vec4 texcolor = texture(MatTexture[i], FragTexcoord * MatTexRepeat[i] + MatTexOffset[i]);
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
