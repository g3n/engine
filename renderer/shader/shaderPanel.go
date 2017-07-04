// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package shader

func init() {
	AddShader("shaderPanelVertex", shaderPanelVertex)
	AddShader("shaderPanelFrag", shaderPanelFrag)
	AddProgram("shaderPanel", "shaderPanelVertex", "shaderPanelFrag")
}

//
// Vertex Shader template
//
const shaderPanelVertex = `
#version {{.Version}}

// Vertex attributes
{{template "attributes" .}}

// Input uniforms
uniform mat4 ModelMatrix;

// Outputs for fragment shader
out vec2 FragTexcoord;


void main() {

    // Always flip texture coordinates
    vec2 texcoord = VertexTexcoord;
    texcoord.y = 1 - texcoord.y;
    FragTexcoord = texcoord;

    // Set position
    vec4 pos = vec4(VertexPosition.xyz, 1);
    gl_Position = ModelMatrix * pos;
}
`

//
// Fragment Shader template
//
const shaderPanelFrag = `
#version {{.Version}}

{{template "material" .}}

// Inputs from vertex shader
in vec2 FragTexcoord;

// Input uniform
uniform vec4 Panel[7];
const int bounds		= 0;              // index of uniform array for bounds coordinates
const int border		= 1;              // index of uniform array for border coordinates
const int padding		= 2;              // index of uniform array for padding coordinates
const int content		= 3;              // index of uniform array for content coordinates
const int borderColor  	= 4;              // index of uniform array for border color
const int paddingColor 	= 5;              // index of uniform array for padding color
const int contentColor	= 6;              // index of uniform array for content color

// Output
out vec4 FragColor;


/***
* Checks if current fragment texture coordinate is inside the
* supplied rectangle in texture coordinates:
* rect[0] - position x [0,1]
* rect[1] - position y [0,1]
* rect[2] - width [0,1]
* rect[3] - height [0,1]
*/
bool checkRect(vec4 rect) {

    if (FragTexcoord.x < rect[0]) {
        return false;
    }
    if (FragTexcoord.x > rect[0] + rect[2]) {
        return false;
    }
    if (FragTexcoord.y < rect[1]) {
        return false;
    }
    if (FragTexcoord.y > rect[1] + rect[3]) {
        return false;
    }
    return true;
}


void main() {

    // Discard fragment outside of received bounds
    // Bounds[0] - xmin
    // Bounds[1] - ymin
    // Bounds[2] - xmax
    // Bounds[3] - ymax
    if (FragTexcoord.x <= Panel[bounds][0] || FragTexcoord.x >= Panel[bounds][2]) {
        discard;
    }
    if (FragTexcoord.y <= Panel[bounds][1] || FragTexcoord.y >= Panel[bounds][3]) {
        discard;
    }

    // Check if fragment is inside content area
    if (checkRect(Panel[content])) {
        // If no texture, the color will be the material color.
        vec4 color = Panel[contentColor];
        {{ if .MatTexturesMax }}
            // Adjust texture coordinates to fit texture inside the content area
            vec2 offset = vec2(-Panel[content][0], -Panel[content][1]);
            vec2 factor = vec2(1/Panel[content][2], 1/Panel[content][3]);
            vec2 texcoord = (FragTexcoord + offset) * factor;
            color = texture(MatTexture[0], texcoord * MatTexRepeat(0) + MatTexOffset(0));
        {{ end }}
        if (color.a == 0) {
            discard;
        }
        FragColor = color;
        return;
    }

    // Checks if fragment is inside paddings area
    if (checkRect(Panel[padding])) {
        FragColor = Panel[paddingColor];
        return;
    }

    // Checks if fragment is inside borders area
    if (checkRect(Panel[border])) {
        FragColor = Panel[borderColor];
        return;
    }

    // Fragment is in margins area (always transparent)
    FragColor = vec4(1,1,1,0);
}

`
