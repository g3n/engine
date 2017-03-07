// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package shader

func init() {
	AddShader("shaderStandardVertex", shaderStandardVertex)
	AddShader("shaderStandardFrag", shaderStandardFrag)
	AddProgram("shaderStandard", "shaderStandardVertex", "shaderStandardFrag")
}

//
// Vertex Shader template
//
const shaderStandardVertex = `
#version {{.Version}}

{{template "attributes" .}}

// Model uniforms
uniform mat4 ModelViewMatrix;
uniform mat3 NormalMatrix;
uniform mat4 MVP;

{{template "lights" .}}
{{template "material" .}}
{{template "phong_model" .}}


// Outputs for the fragment shader.
out vec3 ColorFrontAmbdiff;
out vec3 ColorFrontSpec;
out vec3 ColorBackAmbdiff;
out vec3 ColorBackSpec;
out vec2 FragTexcoord;

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

    vec2 texcoord = VertexTexcoord;
    {{if .MatTexturesMax }}
    // Flips texture coordinate Y if requested.
    if (MatTexFlipY[0] > 0) {
        texcoord.y = 1 - texcoord.y;
    }
    {{ end }}
    FragTexcoord = texcoord;

    gl_Position = MVP * vec4(VertexPosition, 1.0);
}
`

//
// Fragment Shader template
//
const shaderStandardFrag = `
#version {{.Version}}

{{template "material" .}}

// Inputs from Vertex shader
in vec3 ColorFrontAmbdiff;
in vec3 ColorFrontSpec;
in vec3 ColorBackAmbdiff;
in vec3 ColorBackSpec;
in vec2 FragTexcoord;

// Output
out vec4 FragColor;


void main() {

    vec4 texCombined = vec4(1);
    {{if .MatTexturesMax }}
    // Combine all texture colors and opacity
    for (int i = 0; i < {{.MatTexturesMax}}; i++) {
        if (MatTexVisible[i] == false) {
            continue;
        }
        vec4 texcolor = texture(MatTexture[i], FragTexcoord * MatTexRepeat[i] + MatTexOffset[i]);
        if (i == 0) {
            texCombined = texcolor;
        } else {
            texCombined = mix(texCombined, texcolor, texcolor.a);
        }
    }
    {{ end }}

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
