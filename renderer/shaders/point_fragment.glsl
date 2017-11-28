#include <material>

// Inputs from vertex shader
in vec3 Color;
flat in mat2 Rotation;

// Output
out vec4 FragColor;

void main() {

    vec4 texCombined = vec4(1);
    #if MAT_TEXTURES > 0
    // Combine all texture colors and opacity
    for (int i = 0; i < MAT_TEXTURES; i++) {
        vec2 pt = gl_PointCoord - vec2(0.5);
        vec4 texcolor = texture(MatTexture[i], (Rotation * pt + vec2(0.5)) * MatTexRepeat(i) + MatTexOffset(i));
        if (i  == 0) {
            texCombined = texcolor;
        } else {
            texCombined = mix(texCombined, texcolor, texcolor.a);
        }
    }
    #endif

    // Combine material color with texture
    FragColor = min(vec4(Color, MatOpacity) * texCombined, vec4(1));
}

