//
// Material properties uniform
//

// Material parameters uniform array
uniform vec3 Material[6];
// Macros to access elements inside the Material array
#define MatAmbientColor		Material[0]
#define MatDiffuseColor     Material[1]
#define MatSpecularColor    Material[2]
#define MatEmissiveColor    Material[3]
#define MatShininess        Material[4].x
#define MatOpacity          Material[4].y
#define MatPointSize        Material[4].z
#define MatPointRotationZ   Material[5].x

#if MAT_TEXTURES > 0
    // Texture unit sampler array
    uniform sampler2D MatTexture[MAT_TEXTURES];
    // Texture parameters (3*vec2 per texture)
    uniform vec2 MatTexinfo[3*MAT_TEXTURES];
    // Macros to access elements inside the MatTexinfo array
    #define MatTexOffset(a)		MatTexinfo[(3*a)]
    #define MatTexRepeat(a)		MatTexinfo[(3*a)+1]
    #define MatTexFlipY(a)		bool(MatTexinfo[(3*a)+2].x)
    #define MatTexVisible(a)	bool(MatTexinfo[(3*a)+2].y)

// GLSL 3.30 does not allow indexing texture sampler with non constant values.
// This function is used to mix the texture with the specified index with the material color.
// It should be called for each texture index. It uses two externally defined variables:
// vec4 texColor
// vec4 texMixed
vec4 MIX_TEXTURE(vec4 texMixed, vec2 FragTexcoord, int i) {
    if (MatTexVisible(i)) {
        vec4 texColor = texture(MatTexture[i], FragTexcoord * MatTexRepeat(i) + MatTexOffset(i));
        if (i == 0) {
            texMixed = texColor;
        } else {
            texMixed = mix(texMixed, texColor, texColor.a);
        }
    }
    return texMixed;
}

#endif

// TODO for alpha blending dont use mix use implementation below (similar to one in panel shader)
//vec4 prevTexPre = texMixed;
//prevTexPre.rgb *= prevTexPre.a;
//vec4 currTexPre = texColor;
//currTexPre.rgb *= currTexPre.a;
//texMixed = currTexPre + prevTexPre * (1 - currTexPre.a);
//texMixed.rgb /= texMixed.a;
