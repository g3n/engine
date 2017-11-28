//
// Material properties uniform
//
uniform vec3 Material[6];
// Macros to access elements inside the MatTexinfo array
// Each texture uses 3 vec2 elements.
#define MatAmbientColor		Material[0]
#define MatDiffuseColor     Material[1]
#define MatEmissiveColor    Material[2]
#define MatSpecularColor    Material[3]
#define MatShininess        Material[4].x
#define MatOpacity          Material[4].y
#define MatPointSize		Material[4].z
#define MatPointRotationZ	Material[5].x

#if MAT_TEXTURES > 0
    // Texture unit sampler array
    uniform sampler2D MatTexture[MAT_TEXTURES];
    // Texture parameters (3*vec2 per texture)
    uniform mat3 MatTexinfo[MAT_TEXTURES];
    // Macros to access elements inside the MatTexinfo array
    #define MatTexOffset(a)		MatTexinfo[a][0].xy
    #define MatTexRepeat(a)		MatTexinfo[a][1].xy
    #define MatTexFlipY(a)		bool(MatTexinfo[a][2].x)
    #define MatTexVisible(a)	bool(MatTexinfo[a][2].y)
#endif


