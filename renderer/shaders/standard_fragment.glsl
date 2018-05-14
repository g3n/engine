//
// Fragment Shader template
//
#include <material>

// Inputs from Vertex shader
in vec3 ColorFrontAmbdiff;
in vec3 ColorFrontSpec;
in vec3 ColorBackAmbdiff;
in vec3 ColorBackSpec;
in vec2 FragTexcoord;

// Output
out vec4 FragColor;


void main() {

    // Mix material color with textures colors
    vec4 texMixed = vec4(1);
    vec4 texColor;
    #if MAT_TEXTURES==1
        MIX_TEXTURE(0)
    #elif MAT_TEXTURES==2
        MIX_TEXTURE(0)
        MIX_TEXTURE(1)
    #elif MAT_TEXTURES==3
        MIX_TEXTURE(0)
        MIX_TEXTURE(1)
        MIX_TEXTURE(2)
    #endif

    vec4 colorAmbDiff;
    vec4 colorSpec;
    if (gl_FrontFacing) {
        colorAmbDiff = vec4(ColorFrontAmbdiff, MatOpacity);
        colorSpec = vec4(ColorFrontSpec, 0);
    } else {
        colorAmbDiff = vec4(ColorBackAmbdiff, MatOpacity);
        colorSpec = vec4(ColorBackSpec, 0);
    }
    FragColor = min(colorAmbDiff * texMixed + colorSpec, vec4(1));
}

