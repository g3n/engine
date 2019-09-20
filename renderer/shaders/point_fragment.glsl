precision highp float;

#include <material>

// Inputs from vertex shader
in vec3 Color;
flat in mat2 Rotation;

// Output
out vec4 FragColor;

void main() {

    // Compute final texture color
    vec4 texMixed = vec4(1);
    #if MAT_TEXTURES > 0
        vec2 pointCoord = Rotation * gl_PointCoord - vec2(0.5) + vec2(0.5);
        bool firstTex = true;
        if (MatTexVisible(0)) {
            vec4 texColor = texture(MatTexture[0], pointCoord * MatTexRepeat(0) + MatTexOffset(0));
            if (firstTex) {
                texMixed = texColor;
                firstTex = false;
            } else {
                texMixed = Blend(texMixed, texColor);
            }
        }
        #if MAT_TEXTURES > 1
            if (MatTexVisible(1)) {
                vec4 texColor = texture(MatTexture[1], pointCoord * MatTexRepeat(1) + MatTexOffset(1));
                if (firstTex) {
                    texMixed = texColor;
                    firstTex = false;
                } else {
                    texMixed = Blend(texMixed, texColor);
                }
            }
            #if MAT_TEXTURES > 2
                if (MatTexVisible(2)) {
                    vec4 texColor = texture(MatTexture[2], pointCoord * MatTexRepeat(2) + MatTexOffset(2));
                    if (firstTex) {
                        texMixed = texColor;
                        firstTex = false;
                    } else {
                        texMixed = Blend(texMixed, texColor);
                    }
                }
            #endif
        #endif
    #endif

    // Generates final color
    FragColor = min(vec4(Color, MatOpacity) * texMixed, vec4(1));
}
