// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package shader

func init() {
	AddShader("shaderPhongVertex", shaderPhongVertex)
	AddShader("shaderPhongFrag", shaderPhongFrag)
	AddProgram("shaderPhong", "shaderPhongVertex", "shaderPhongFrag")
}

//
// Vertex Shade template
//
var shaderPhongVertex = `
#version {{.Version}}

{{template "attributes" .}}

// Model uniforms
uniform mat4 ModelViewMatrix;
uniform mat3 NormalMatrix;
uniform mat4 MVP;

{{template "material" .}}

// Output variables for Fragment shader
out vec4 Position;
out vec3 Normal;
out vec3 CamDir;
out vec2 FragTexcoord;

void main() {

    // Transform this vertex position to camera coordinates.
    Position = ModelViewMatrix * vec4(VertexPosition, 1.0);

    // Transform this vertex normal to camera coordinates.
    Normal = normalize(NormalMatrix * VertexNormal);

    // Calculate the direction vector from the vertex to the camera
    // The camera is at 0,0,0
    CamDir = normalize(-Position.xyz);

    // Flips texture coordinate Y if requested.
    vec2 texcoord = VertexTexcoord;
    {{ if .MatTexturesMax }}
    if (MatTexFlipY(0)) {
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
var shaderPhongFrag = `
#version {{.Version}}

// Inputs from vertex shader
in vec4 Position;       // Vertex position in camera coordinates.
in vec3 Normal;         // Vertex normal in camera coordinates.
in vec3 CamDir;         // Direction from vertex to camera
in vec2 FragTexcoord;

{{template "lights" .}}
{{template "material" .}}
{{template "phong_model" .}}

// Final fragment color
out vec4 FragColor;

void main() {

    // Combine all texture colors
    vec4 texCombined = vec4(1);
    {{ range loop .MatTexturesMax }}
    if (MatTexVisible({{.}})) {
        vec4 texcolor = texture(MatTexture[{{.}}], FragTexcoord * MatTexRepeat({{.}}) + MatTexOffset({{.}}));
        if ({{.}} == 0) {
            texCombined = texcolor;
        } else {
            texCombined = mix(texCombined, texcolor, texcolor.a);
        }
    }
    {{ end }}

    // Combine material with texture colors
    vec4 matDiffuse = vec4(MatDiffuseColor, MatOpacity) * texCombined;
    vec4 matAmbient = vec4(MatAmbientColor, MatOpacity) * texCombined;

    // Inverts the fragment normal if not FrontFacing
    vec3 fragNormal = Normal;
    if (!gl_FrontFacing) {
        fragNormal = -fragNormal;
    }

    // Calculates the Ambient+Diffuse and Specular colors for this fragment using the Phong model.
    vec3 Ambdiff, Spec;
    phongModel(Position, fragNormal, CamDir, vec3(matAmbient), vec3(matDiffuse), Ambdiff, Spec);

    // Final fragment color
    FragColor = min(vec4(Ambdiff + Spec, matDiffuse.a), vec4(1.0));
}

`
