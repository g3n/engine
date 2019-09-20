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
    // Alpha compositing (see here: https://ciechanow.ski/alpha-compositing/)
    vec4 Blend(vec4 texMixed, vec4 texColor) {
        texMixed.rgb *= texMixed.a;
        texColor.rgb *= texColor.a;
        texMixed = texColor + texMixed * (1 - texColor.a);
        if (texMixed.a > 0.0) {
            texMixed.rgb /= texMixed.a;
        }
        return texMixed;
    }
#endif
