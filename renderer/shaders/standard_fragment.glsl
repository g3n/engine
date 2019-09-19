precision highp float;

// Inputs from vertex shader
in vec4 Position;     // Fragment position in camera coordinates
in vec3 Normal;       // Interpolated fragment normal in camera coordinates
in vec3 CamDir;       // Direction from fragment to camera
in vec2 FragTexcoord; // Fragment texture coordinates

#include <lights>
#include <material>
#include <phong_model>

// Final fragment color
out vec4 FragColor;

void main() {

    // Mix material color with textures colors
    vec4 texMixed = vec4(1);
    #if MAT_TEXTURES==1
        texMixed = MIX_TEXTURE(texMixed, FragTexcoord, 0);
    #elif MAT_TEXTURES==2
        texMixed = MIX_TEXTURE(texMixed, FragTexcoord, 0);
        texMixed = MIX_TEXTURE(texMixed, FragTexcoord, 1);
    #elif MAT_TEXTURES==3
        texMixed = MIX_TEXTURE(texMixed, FragTexcoord, 0);
        texMixed = MIX_TEXTURE(texMixed, FragTexcoord, 1);
        texMixed = MIX_TEXTURE(texMixed, FragTexcoord, 2);
    #endif

    // Combine material with texture colors
    vec4 matDiffuse = vec4(MatDiffuseColor, MatOpacity) * texMixed;
    vec4 matAmbient = vec4(MatAmbientColor, MatOpacity) * texMixed;

    // Normalize interpolated normal as it may have shrinked
    vec3 fragNormal = normalize(Normal);

    // Invert the fragment normal if not FrontFacing
    if (!gl_FrontFacing) {
        fragNormal = -fragNormal;
    }

    // Calculates the Ambient+Diffuse and Specular colors for this fragment using the Phong model.
    vec3 Ambdiff, Spec;
    phongModel(Position, fragNormal, CamDir, vec3(matAmbient), vec3(matDiffuse), Ambdiff, Spec);

    // Final fragment color
    FragColor = min(vec4(Ambdiff + Spec, matDiffuse.a), vec4(1.0));
}
