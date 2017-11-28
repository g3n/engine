#include <attributes>

// Model uniforms
uniform mat4 MVP;

// Material uniforms
#include <material>

// Outputs for fragment shader
out vec3 Color;
flat out mat2 Rotation;

void main() {

    // Rotation matrix for fragment shader
    float rotSin = sin(MatPointRotationZ);
    float rotCos = cos(MatPointRotationZ);
    Rotation = mat2(rotCos, rotSin, - rotSin, rotCos);

    // Sets the vertex position
    vec4 pos = MVP * vec4(VertexPosition, 1.0);
    gl_Position = pos;

    // Sets the size of the rasterized point decreasing with distance
    gl_PointSize = (1.0 - pos.z / pos.w) * MatPointSize;

    // Outputs color
    Color = MatEmissiveColor;
}

