//
// Material properties uniform
//
uniform vec3 Material[5];
// Macros to access elements inside the MatTexinfo array
// Each texture uses 3 vec2 elements.
#define MatAmbientColor		Material[0]
#define MatDiffuseColor     Material[1]
#define MatEmissiveColor    Material[2]
#define MatSpecularColor    Material[3]
#define MatShininess        Material[4].x
#define MatOpacity          Material[4].y

#if MAT_TEXTURES > 0
    // Texture unit sampler array
    uniform sampler2D MatTexture[MAT_TEXTURES];
    // Texture parameters (3*vec2 per texture)
    uniform vec2 MatTexinfo[3*MAT_TEXTURES];
    // Macros to access elements inside the MatTexinfo array
    // Each texture uses 3 vec2 elements.
    #define MatTexOffset(a)		MatTexinfo[(3*a)]
    #define MatTexRepeat(a)		MatTexinfo[(3*a)+1]
    #define MatTexFlipY(a)		bool(MatTexinfo[(3*a)+2].x)
    #define MatTexVisible(a)	bool(MatTexinfo[(3*a)+2].y)
#endif


