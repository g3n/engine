//
// Vertex Shader
//
#include <attributes>

// Model uniforms
uniform mat4 ModelViewMatrix;
uniform mat3 NormalMatrix;
uniform mat4 MVP;

#include <material>
#include <morphtarget_vertex_declaration>

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
#if MAT_TEXTURES>0
    if (MatTexFlipY(0)) {
        texcoord.y = 1 - texcoord.y;
    }
#endif
    FragTexcoord = texcoord;
    vec3 vPosition = VertexPosition;
    #include <morphtarget_vertex> [MORPHTARGETS]

    gl_Position = MVP * vec4(vPosition, 1.0);
}

