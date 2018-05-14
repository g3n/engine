#include <material>

// GLSL 3.30 does not allow indexing texture sampler with non constant values.
// This macro is used to mix the texture with the specified index with the material color.
// It should be called for each texture index.
#define MIX_POINT_TEXTURE(i)                                                                                     \
    if (MatTexVisible(i)) {                                                                                      \
        vec2 pt = gl_PointCoord - vec2(0.5);                                                                     \
        vec4 texColor = texture(MatTexture[i], (Rotation * pt + vec2(0.5)) * MatTexRepeat(i) + MatTexOffset(i)); \
        if (i == 0) {                                                                                            \
            texMixed = texColor;                                                                                 \
        } else {                                                                                                 \
            texMixed = mix(texMixed, texColor, texColor.a);                                                      \
        }                                                                                                        \
    }

// Inputs from vertex shader
in vec3 Color;
flat in mat2 Rotation;

// Output
out vec4 FragColor;

void main() {

    // Mix material color with textures colors
    vec4 texMixed = vec4(1);
    #if MAT_TEXTURES==1
        MIX_POINT_TEXTURE(0)
    #elif MAT_TEXTURES==2
        MIX_POINT_TEXTURE(0)
        MIX_POINT_TEXTURE(1)
    #elif MAT_TEXTURES==3
        MIX_POINT_TEXTURE(0)
        MIX_POINT_TEXTURE(1)
        MIX_POINT_TEXTURE(2)
    #endif

    // Generates final color
    FragColor = min(vec4(Color, MatOpacity) * texMixed, vec4(1));
}

