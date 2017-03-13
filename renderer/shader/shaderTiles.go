// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package shader

const shaderTilesVertex = `
#version {{.Version}}

{{template "attributes" .}}

// Model uniforms
uniform mat4 ModelViewMatrix;
uniform mat3 NormalMatrix;
uniform mat4 MVP;

// Templates
{{template "uniLights" .}}
{{template "uniMaterial" .}}
{{template "funcPhongModel" .}}


// Outputs for the fragment shader.
out vec3 ColorFrontAmbdiff;
out vec3 ColorFrontSpec;
out vec3 ColorBackAmbdiff;
out vec3 ColorBackSpec;
out vec2 FragTexcoord;
flat out vec4 TexOffsets;


void main() {

    // Transform this vertex normal to camera coordinates.
    vec3 normal = normalize(NormalMatrix * VertexNormal);

    // Calculate this vertex position in camera coordinates
    vec4 position = ModelViewMatrix * vec4(VertexPosition, 1.0);

    // Calculate the direction vector from the vertex to the camera
    // The camera is at 0,0,0
    vec3 camDir = normalize(-position.xyz);

    // Calculates the vertex Ambient+Diffuse and Specular colors using the Phong model
    // for the front and back
    phongModel(position,  normal, camDir, MatAmbientColor, MatDiffuseColor, ColorFrontAmbdiff, ColorFrontSpec);
    phongModel(position, -normal, camDir, MatAmbientColor, MatDiffuseColor, ColorBackAmbdiff, ColorBackSpec);

    // Flips texture coordinate Y if requested.
    vec2 texcoord = VertexTexcoord;
    {{ if .MatTexturesMax }}
    if (MatTexFlipY[0] > 0) {
        texcoord.y = 1 - texcoord.y;
    }
    {{ end }}
    FragTexcoord = texcoord;
    TexOffsets = VertexTexoffsets;

    gl_Position = MVP * vec4(VertexPosition, 1.0);
}
`

//
// Fragment Shader template
//
const shaderTilesFrag = `
#version {{.Version}}

{{template "uniMaterial" .}}

// Inputs from Vertex shader
in vec3 ColorFrontAmbdiff;
in vec3 ColorFrontSpec;
in vec3 ColorBackAmbdiff;
in vec3 ColorBackSpec;
in vec2 FragTexcoord;
flat in vec4 TexOffsets;

// Output
out vec4 FragColor;


void main() {

    vec2 offset = vec2(TexOffsets.xy);
    vec2 repeat = vec2(TexOffsets.zw);
    vec4 texCombined = texture(MatTexture[0], FragTexcoord * repeat + offset);

//    vec4 texCombined = vec4(1);
//    {{ range loop .MatTexturesMax }}
//    // Combine all texture colors and opacity
//    {
//        vec4 texcolor = texture(MatTexture[{{.}}], FragTexcoord * MatTexRepeat[{{.}}] + MatTexOffset[{{.}}]);
//        if ({{.}} == 0) {
//            texCombined = texcolor;
//        } else {
//            texCombined = mix(texCombined, texcolor, texcolor.a);
//        }
//    }
//    {{ end }}

    vec4 colorAmbDiff;
    vec4 colorSpec;
    if (gl_FrontFacing) {
        colorAmbDiff = vec4(ColorFrontAmbdiff, MatOpacity);
        colorSpec = vec4(ColorFrontSpec, 0);
    } else {
        colorAmbDiff = vec4(ColorBackAmbdiff, MatOpacity);
        colorSpec = vec4(ColorBackSpec, 0);
    }
    FragColor = min(colorAmbDiff * texCombined + colorSpec, vec4(1));
}
`
