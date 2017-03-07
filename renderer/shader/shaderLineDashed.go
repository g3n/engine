package shader


const shaderLineDashedVertex = `
#version {{.Version}}

{{?attributes}}

// Model uniforms
uniform mat4 MVP;

// Material uniforms
uniform float Scale;

// Outputs for fragment shader
out vec3 Color;
out float vLineDistance; 

void main() {

    vLineDistance = Scale * VertexDistance;
    Color = VertexColor;
    gl_Position = MVP * vec4(VertexPosition, 1.0);
}
`


//
// Fragment Shader template
//
const shaderLineDashedFrag = `
#version {{.Version}}

// Inputs from vertex shader
in vec3 Color;
in float vLineDistance;

// Material uniforms
uniform float DashSize;
uniform float TotalSize;

// Output
out vec4 FragColor;

void main() {

    if (mod(vLineDistance, TotalSize) > DashSize) {
        discard;
    }

    FragColor = vec4(Color, 1.0);
}

`
