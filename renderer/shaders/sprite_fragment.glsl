//
// Fragment shader for sprite
//

precision highp float;

#include <material>

// Inputs from vertex shader
in vec3 Color;
in vec2 FragTexcoord;

// Output
out vec4 FragColor;

void main() {

    // Combine all texture colors and opacity
    vec4 texCombined = vec4(1);
#if MAT_TEXTURES>0
    for (int i = 0; i < {{.MatTexturesMax}}; i++) {
        vec4 texcolor = texture(MatTexture[i], FragTexcoord * MatTexRepeat(i) + MatTexOffset(i));
        if (i == 0) {
            texCombined = texcolor;
        } else {
            texCombined = mix(texCombined, texcolor, texcolor.a);
        }
    }
#endif

    // Combine material color with texture
    FragColor = min(vec4(Color, MatOpacity) * texCombined, vec4(1));
}

