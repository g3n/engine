//
// Lights uniforms
//

// Ambient lights uniforms
#if AMB_LIGHTS>0
uniform vec3 AmbientLightColor[AMB_LIGHTS];
#endif

// Directional lights uniform array. Each directional light uses 2 elements
#if DIR_LIGHTS>0
uniform vec3 DirLight[2*DIR_LIGHTS];
// Macros to access elements inside the DirectionalLight uniform array
#define DirLightColor(a)		DirLight[2*a]
#define DirLightPosition(a)		DirLight[2*a+1]
#endif

// Point lights uniform array. Each point light uses 3 elements
#if POINT_LIGHTS>0
uniform vec3 PointLight[3*POINT_LIGHTS];
// Macros to access elements inside the PointLight uniform array
#define PointLightColor(a)			PointLight[3*a]
#define PointLightPosition(a)		PointLight[3*a+1]
#define PointLightLinearDecay(a)	PointLight[3*a+2].x
#define PointLightQuadraticDecay(a)	PointLight[3*a+2].y
#endif


