//
// Vertex shader basic
//
#include <attributes>

// Model uniforms
uniform mat4 MVP;

// Final output color for fragment shader
out vec3 Color;

void main() {

    Color = VertexColor;
    gl_Position = MVP * vec4(VertexPosition, 1.0);
}


