#include <attributes>

// Model uniforms
uniform mat4 ModelViewMatrix;
uniform mat3 NormalMatrix;
uniform mat4 MVP;

#include <material>
#include <morphtarget_vertex_declaration>
#include <bones_vertex_declaration>

// Output variables for Fragment shader
out vec4 Position;
out vec3 Normal;
out vec2 FragTexcoord;

void main() {

    // Transform vertex position to camera coordinates
    Position = ModelViewMatrix * vec4(VertexPosition, 1.0);

    // Transform vertex normal to camera coordinates
    Normal = normalize(NormalMatrix * VertexNormal);

    vec2 texcoord = VertexTexcoord;
#if MAT_TEXTURES > 0
    // Flip texture coordinate Y if requested.
    if (MatTexFlipY(0)) {
        texcoord.y = 1.0 - texcoord.y;
    }
#endif
    FragTexcoord = texcoord;
    vec3 vPosition = VertexPosition;
    mat4 finalWorld = mat4(1.0);
    #include <morphtarget_vertex>
    #include <bones_vertex>

    // Output projected and transformed vertex position
    gl_Position = MVP * finalWorld * vec4(vPosition, 1.0);
}
