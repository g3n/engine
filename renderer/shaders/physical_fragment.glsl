//
// Physically Based Shading of a microfacet surface material - Fragment Shader
// Modified from reference implementation at https://github.com/KhronosGroup/glTF-WebGL-PBR
//
// References:
// [1] Real Shading in Unreal Engine 4
//     http://blog.selfshadow.com/publications/s2013-shading-course/karis/s2013_pbs_epic_notes_v2.pdf
// [2] Physically Based Shading at Disney
//     http://blog.selfshadow.com/publications/s2012-shading-course/burley/s2012_pbs_disney_brdf_notes_v3.pdf
// [3] README.md - Environment Maps
//     https://github.com/KhronosGroup/glTF-WebGL-PBR/#environment-maps
// [4] "An Inexpensive BRDF Model for Physically based Rendering" by Christophe Schlick
//     https://www.cs.virginia.edu/~jdl/bib/appearance/analytic%20models/schlick94b.pdf

//#extension GL_EXT_shader_texture_lod: enable
//#extension GL_OES_standard_derivatives : enable

precision highp float;

//uniform vec3 u_LightDirection;
//uniform vec3 u_LightColor;

//#ifdef USE_IBL
//uniform samplerCube u_DiffuseEnvSampler;
//uniform samplerCube u_SpecularEnvSampler;
//uniform sampler2D u_brdfLUT;
//#endif

#ifdef HAS_BASECOLORMAP
uniform sampler2D uBaseColorSampler;
#endif
#ifdef HAS_METALROUGHNESSMAP
uniform sampler2D uMetallicRoughnessSampler;
#endif
#ifdef HAS_NORMALMAP
uniform sampler2D uNormalSampler;
//uniform float uNormalScale;
#endif
#ifdef HAS_EMISSIVEMAP
uniform sampler2D uEmissiveSampler;
#endif
#ifdef HAS_OCCLUSIONMAP
uniform sampler2D uOcclusionSampler;
uniform float uOcclusionStrength;
#endif

// Material parameters uniform array
uniform vec4 Material[3];
// Macros to access elements inside the Material array
#define uBaseColor		    Material[0]
#define uEmissiveColor      Material[1]
#define uMetallicFactor     Material[2].x
#define uRoughnessFactor    Material[2].y

#include <lights>

// Inputs from vertex shader
in vec3 Position;       // Vertex position in camera coordinates.
in vec3 Normal;         // Vertex normal in camera coordinates.
in vec3 CamDir;         // Direction from vertex to camera
in vec2 FragTexcoord;

// Final fragment color
out vec4 FragColor;

// Encapsulate the various inputs used by the various functions in the shading equation
// We store values in this struct to simplify the integration of alternative implementations
// of the shading terms, outlined in the Readme.MD Appendix.
struct PBRLightInfo
{
    float NdotL;                  // cos angle between normal and light direction
    float NdotV;                  // cos angle between normal and view direction
    float NdotH;                  // cos angle between normal and half vector
    float LdotH;                  // cos angle between light direction and half vector
    float VdotH;                  // cos angle between view direction and half vector
};

struct PBRInfo
{
    float perceptualRoughness;    // roughness value, as authored by the model creator (input to shader)
    float metalness;              // metallic value at the surface
    vec3 reflectance0;            // full reflectance color (normal incidence angle)
    vec3 reflectance90;           // reflectance color at grazing angle
    float alphaRoughness;         // roughness mapped to a more linear change in the roughness (proposed by [2])
    vec3 diffuseColor;            // color contribution from diffuse lighting
    vec3 specularColor;           // color contribution from specular lighting
};

const float M_PI = 3.141592653589793;
const float c_MinRoughness = 0.04;

vec4 SRGBtoLINEAR(vec4 srgbIn) {
//#ifdef MANUAL_SRGB
//    #ifdef SRGB_FAST_APPROXIMATION
//        vec3 linOut = pow(srgbIn.xyz,vec3(2.2));
//    #else //SRGB_FAST_APPROXIMATION
        vec3 bLess = step(vec3(0.04045),srgbIn.xyz);
        vec3 linOut = mix( srgbIn.xyz/vec3(12.92), pow((srgbIn.xyz+vec3(0.055))/vec3(1.055),vec3(2.4)), bLess );
//    #endif //SRGB_FAST_APPROXIMATION
        return vec4(linOut,srgbIn.w);
//#else //MANUAL_SRGB
//    return srgbIn;
//#endif //MANUAL_SRGB
}

// Find the normal for this fragment, pulling either from a predefined normal map
// or from the interpolated mesh normal and tangent attributes.
vec3 getNormal()
{
    // Retrieve the tangent space matrix
//#ifndef HAS_TANGENTS
    vec3 pos_dx = dFdx(Position);
    vec3 pos_dy = dFdy(Position);
    vec3 tex_dx = dFdx(vec3(FragTexcoord, 0.0));
    vec3 tex_dy = dFdy(vec3(FragTexcoord, 0.0));
    vec3 t = (tex_dy.t * pos_dx - tex_dx.t * pos_dy) / (tex_dx.s * tex_dy.t - tex_dy.s * tex_dx.t);

//#ifdef HAS_NORMALS
    vec3 ng = normalize(Normal);
//#else
//    vec3 ng = cross(pos_dx, pos_dy);
//#endif

    t = normalize(t - ng * dot(ng, t));
    vec3 b = normalize(cross(ng, t));
    mat3 tbn = mat3(t, b, ng);
//#else // HAS_TANGENTS
//    mat3 tbn = v_TBN;
//#endif

#ifdef HAS_NORMALMAP
    float uNormalScale = 1.0;
    vec3 n = texture(uNormalSampler, FragTexcoord).rgb;
    n = normalize(tbn * ((2.0 * n - 1.0) * vec3(uNormalScale, uNormalScale, 1.0)));
#else
    // The tbn matrix is linearly interpolated, so we need to re-normalize
    vec3 n = normalize(tbn[2].xyz);
#endif

    return n;
}

// Calculation of the lighting contribution from an optional Image Based Light source.
// Precomputed Environment Maps are required uniform inputs and are computed as outlined in [1].
// See our README.md on Environment Maps [3] for additional discussion.
vec3 getIBLContribution(PBRInfo pbrInputs, PBRLightInfo pbrLight, vec3 n, vec3 reflection)
{
    float mipCount = 9.0; // resolution of 512x512
    float lod = (pbrInputs.perceptualRoughness * mipCount);
    // retrieve a scale and bias to F0. See [1], Figure 3
    vec3 brdf = vec3(0.5,0.5,0.5);//SRGBtoLINEAR(texture(u_brdfLUT, vec2(pbrLight.NdotV, 1.0 - pbrInputs.perceptualRoughness))).rgb;
    vec3 diffuseLight = vec3(0.5,0.5,0.5);//SRGBtoLINEAR(textureCube(u_DiffuseEnvSampler, n)).rgb;

//#ifdef USE_TEX_LOD
//    vec3 specularLight = SRGBtoLINEAR(textureCubeLodEXT(u_SpecularEnvSampler, reflection, lod)).rgb;
//#else
    vec3 specularLight = vec3(0.5,0.5,0.5);//SRGBtoLINEAR(textureCube(u_SpecularEnvSampler, reflection)).rgb;
//#endif

    vec3 diffuse = diffuseLight * pbrInputs.diffuseColor;
    vec3 specular = specularLight * (pbrInputs.specularColor * brdf.x + brdf.y);

    // For presentation, this allows us to disable IBL terms
//    diffuse *= u_ScaleIBLAmbient.x;
//    specular *= u_ScaleIBLAmbient.y;

    return diffuse + specular;
}

// Basic Lambertian diffuse
// Implementation from Lambert's Photometria https://archive.org/details/lambertsphotome00lambgoog
// See also [1], Equation 1
vec3 diffuse(PBRInfo pbrInputs)
{
    return pbrInputs.diffuseColor / M_PI;
}

// The following equation models the Fresnel reflectance term of the spec equation (aka F())
// Implementation of fresnel from [4], Equation 15
vec3 specularReflection(PBRInfo pbrInputs, PBRLightInfo pbrLight)
{
    return pbrInputs.reflectance0 + (pbrInputs.reflectance90 - pbrInputs.reflectance0) * pow(clamp(1.0 - pbrLight.VdotH, 0.0, 1.0), 5.0);
}

// This calculates the specular geometric attenuation (aka G()),
// where rougher material will reflect less light back to the viewer.
// This implementation is based on [1] Equation 4, and we adopt their modifications to
// alphaRoughness as input as originally proposed in [2].
float geometricOcclusion(PBRInfo pbrInputs, PBRLightInfo pbrLight)
{
    float NdotL = pbrLight.NdotL;
    float NdotV = pbrLight.NdotV;
    float r = pbrInputs.alphaRoughness;

    float attenuationL = 2.0 * NdotL / (NdotL + sqrt(r * r + (1.0 - r * r) * (NdotL * NdotL)));
    float attenuationV = 2.0 * NdotV / (NdotV + sqrt(r * r + (1.0 - r * r) * (NdotV * NdotV)));
    return attenuationL * attenuationV;
}

// The following equation(s) model the distribution of microfacet normals across the area being drawn (aka D())
// Implementation from "Average Irregularity Representation of a Roughened Surface for Ray Reflection" by T. S. Trowbridge, and K. P. Reitz
// Follows the distribution function recommended in the SIGGRAPH 2013 course notes from EPIC Games [1], Equation 3.
float microfacetDistribution(PBRInfo pbrInputs, PBRLightInfo pbrLight)
{
    float roughnessSq = pbrInputs.alphaRoughness * pbrInputs.alphaRoughness;
    float f = (pbrLight.NdotH * roughnessSq - pbrLight.NdotH) * pbrLight.NdotH + 1.0;
    return roughnessSq / (M_PI * f * f);
}

vec3 pbrModel(PBRInfo pbrInputs, vec3 lightColor, vec3 lightDir) {

    vec3 n = getNormal();                             // normal at surface point
    vec3 v = normalize(CamDir);                       // Vector from surface point to camera
    vec3 l = normalize(lightDir);                     // Vector from surface point to light
    vec3 h = normalize(l+v);                          // Half vector between both l and v
    vec3 reflection = -normalize(reflect(v, n));

    float NdotL = clamp(dot(n, l), 0.001, 1.0);
    float NdotV = abs(dot(n, v)) + 0.001;
    float NdotH = clamp(dot(n, h), 0.0, 1.0);
    float LdotH = clamp(dot(l, h), 0.0, 1.0);
    float VdotH = clamp(dot(v, h), 0.0, 1.0);

    PBRLightInfo pbrLight = PBRLightInfo(
        NdotL,
        NdotV,
        NdotH,
        LdotH,
        VdotH
    );

    // Calculate the shading terms for the microfacet specular shading model
    vec3 F = specularReflection(pbrInputs, pbrLight);
    float G = geometricOcclusion(pbrInputs, pbrLight);
    float D = microfacetDistribution(pbrInputs, pbrLight);

    // Calculation of analytical lighting contribution
    vec3 diffuseContrib = (1.0 - F) * diffuse(pbrInputs);
    vec3 specContrib = F * G * D / (4.0 * NdotL * NdotV);
    // Obtain final intensity as reflectance (BRDF) scaled by the energy of the light (cosine law)
    vec3 color = NdotL * lightColor * (diffuseContrib + specContrib);

    return color;
}

void main() {

    float perceptualRoughness = uRoughnessFactor;
    float metallic = uMetallicFactor;

#ifdef HAS_METALROUGHNESSMAP
    // Roughness is stored in the 'g' channel, metallic is stored in the 'b' channel.
    // This layout intentionally reserves the 'r' channel for (optional) occlusion map data
    vec4 mrSample = texture(uMetallicRoughnessSampler, FragTexcoord);
    perceptualRoughness = mrSample.g * perceptualRoughness;
    metallic = mrSample.b * metallic;
#endif

    perceptualRoughness = clamp(perceptualRoughness, c_MinRoughness, 1.0);
    metallic = clamp(metallic, 0.0, 1.0);
    // Roughness is authored as perceptual roughness; as is convention,
    // convert to material roughness by squaring the perceptual roughness [2].
    float alphaRoughness = perceptualRoughness * perceptualRoughness;

    // The albedo may be defined from a base texture or a flat color
#ifdef HAS_BASECOLORMAP
    vec4 baseColor = SRGBtoLINEAR(texture(uBaseColorSampler, FragTexcoord)) * uBaseColor;
#else
    vec4 baseColor = uBaseColor;
#endif

    vec3 f0 = vec3(0.04);
    vec3 diffuseColor = baseColor.rgb * (vec3(1.0) - f0);
    diffuseColor *= 1.0 - metallic;

    vec3 specularColor = mix(f0, baseColor.rgb, uMetallicFactor);

    // Compute reflectance.
    float reflectance = max(max(specularColor.r, specularColor.g), specularColor.b);

    // For typical incident reflectance range (between 4% to 100%) set the grazing reflectance to 100% for typical fresnel effect.
    // For very low reflectance range on highly diffuse objects (below 4%), incrementally reduce grazing reflectance to 0%.
    float reflectance90 = clamp(reflectance * 25.0, 0.0, 1.0);
    vec3 specularEnvironmentR0 = specularColor.rgb;
    vec3 specularEnvironmentR90 = vec3(1.0, 1.0, 1.0) * reflectance90;

    PBRInfo pbrInputs = PBRInfo(
        perceptualRoughness,
        metallic,
        specularEnvironmentR0,
        specularEnvironmentR90,
        alphaRoughness,
        diffuseColor,
        specularColor
    );

//    vec3 normal = getNormal();
    vec3 color = vec3(0.0);

#if AMB_LIGHTS>0
    // Ambient lights
    for (int i = 0; i < AMB_LIGHTS; i++) {
        color += AmbientLightColor[i] * pbrInputs.diffuseColor;
    }
#endif

#if DIR_LIGHTS>0
    // Directional lights
    for (int i = 0; i < DIR_LIGHTS; i++) {
        // Diffuse reflection
        // DirLightPosition is the direction of the current light
        vec3 lightDirection = normalize(DirLightPosition(i));
        // PBR
        color += pbrModel(pbrInputs, DirLightColor(i), lightDirection);
    }
#endif

#if POINT_LIGHTS>0
    // Point lights
    for (int i = 0; i < POINT_LIGHTS; i++) {
        // Common calculations
        // Calculates the direction and distance from the current vertex to this point light.
        vec3 lightDirection = PointLightPosition(i) - vec3(Position);
        float lightDistance = length(lightDirection);
        // Normalizes the lightDirection
        lightDirection = lightDirection / lightDistance;
        // Calculates the attenuation due to the distance of the light
        float attenuation = 1.0 / (1.0 + PointLightLinearDecay(i) * lightDistance +
            PointLightQuadraticDecay(i) * lightDistance * lightDistance);
        vec3 attenuatedColor = PointLightColor(i) * attenuation;
        // PBR
        color += pbrModel(pbrInputs, attenuatedColor, lightDirection);
    }
#endif

#if SPOT_LIGHTS>0
    for (int i = 0; i < SPOT_LIGHTS; i++) {

        // Calculates the direction and distance from the current vertex to this spot light.
        vec3 lightDirection = SpotLightPosition(i) - vec3(Position);
        float lightDistance = length(lightDirection);
        lightDirection = lightDirection / lightDistance;

        // Calculates the attenuation due to the distance of the light
        float attenuation = 1.0 / (1.0 + SpotLightLinearDecay(i) * lightDistance +
            SpotLightQuadraticDecay(i) * lightDistance * lightDistance);

        // Calculates the angle between the vertex direction and spot direction
        // If this angle is greater than the cutoff the spotlight will not contribute
        // to the final color.
        float angle = acos(dot(-lightDirection, SpotLightDirection(i)));
        float cutoff = radians(clamp(SpotLightCutoffAngle(i), 0.0, 90.0));

        if (angle < cutoff) {
            float spotFactor = pow(dot(-lightDirection, SpotLightDirection(i)), SpotLightAngularDecay(i));
            vec3 attenuatedColor = SpotLightColor(i) * attenuation * spotFactor;
            // PBR
            color += pbrModel(pbrInputs, attenuatedColor, lightDirection);
        }
    }
#endif

    // Calculate lighting contribution from image based lighting source (IBL)
//#ifdef USE_IBL
//    color += getIBLContribution(pbrInputs, n, reflection);
//#endif

    // Apply optional PBR terms for additional (optional) shading
#ifdef HAS_OCCLUSIONMAP
    float ao = texture(uOcclusionSampler, FragTexcoord).r;
    color = mix(color, color * ao, 1.0);//, uOcclusionStrength);
#endif

#ifdef HAS_EMISSIVEMAP
    vec3 emissive = SRGBtoLINEAR(texture(uEmissiveSampler, FragTexcoord)).rgb * vec3(uEmissiveColor);
#else
    vec3 emissive = vec3(uEmissiveColor);
#endif
    color += emissive;

    // Base Color
//    FragColor = baseColor;

    // Normal
//    FragColor = vec4(n, 1.0);

    // Emissive Color
//    FragColor = vec4(emissive, 1.0);

    // F
//    color = F;

    // G
//    color = vec3(G);

    // D
//    color = vec3(D);

    // Specular
//    color = specContrib;

    // Diffuse
//    color = diffuseContrib;

    // Roughness
//    color = vec3(perceptualRoughness);

    // Metallic
//    color = vec3(metallic);

    // Final fragment color
    FragColor = vec4(pow(color,vec3(1.0/2.2)), baseColor.a);
}


