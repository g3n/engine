//
// Lights uniforms
//

#if AMB_LIGHTS>0
    // Ambient lights color uniform
    uniform vec3 AmbientLightColor[AMB_LIGHTS];
#endif

#if DIR_LIGHTS>0
    // Directional lights uniform array. Each directional light uses 2 elements
    uniform vec3 DirLight[2*DIR_LIGHTS];
    // Macros to access elements inside the DirectionalLight uniform array
    #define DirLightColor(a)		DirLight[2*a]
    #define DirLightPosition(a)		DirLight[2*a+1]
#endif

#if POINT_LIGHTS>0
    // Point lights uniform array. Each point light uses 3 elements
    uniform vec3 PointLight[3*POINT_LIGHTS];
    // Macros to access elements inside the PointLight uniform array
    #define PointLightColor(a)			PointLight[3*a]
    #define PointLightPosition(a)		PointLight[3*a+1]
    #define PointLightLinearDecay(a)	PointLight[3*a+2].x
    #define PointLightQuadraticDecay(a)	PointLight[3*a+2].y
#endif

#if SPOT_LIGHTS>0
    // Spot lights uniforms. Each spot light uses 5 elements
    uniform vec3  SpotLight[5*SPOT_LIGHTS];
    // Macros to access elements inside the PointLight uniform array
    #define SpotLightColor(a)			SpotLight[5*a]
    #define SpotLightPosition(a)		SpotLight[5*a+1]
    #define SpotLightDirection(a)		SpotLight[5*a+2]
    #define SpotLightAngularDecay(a)	SpotLight[5*a+3].x
    #define SpotLightCutoffAngle(a)		SpotLight[5*a+3].y
    #define SpotLightLinearDecay(a)		SpotLight[5*a+3].z
    #define SpotLightQuadraticDecay(a)	SpotLight[5*a+4].x
#endif

