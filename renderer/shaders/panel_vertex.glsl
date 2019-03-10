//
// Vertex shader panel
//
#include <attributes>

// Model uniforms
uniform mat4 ModelMatrix;

// Outputs for fragment shader
out vec2 FragTexcoord;


void main() {

    // Always flip texture coordinates
    vec2 texcoord = VertexTexcoord;
    texcoord.y = 1.0 - texcoord.y;
    FragTexcoord = texcoord;

    // Set position
    vec4 pos = vec4(VertexPosition.xyz, 1);
    gl_Position = ModelMatrix * pos;
}

