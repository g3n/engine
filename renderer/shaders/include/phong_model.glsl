/***
 phong lighting model
 Parameters:
    position:   input vertex position in camera coordinates
    normal:     input vertex normal in camera coordinates
    camDir:     input camera directions
    matAmbient: input material ambient color
    matDiffuse: input material diffuse color
    ambdiff:    output ambient+diffuse color
    spec:       output specular color
 Uniforms:
    AmbientLightColor[]
    DiffuseLightColor[]
    DiffuseLightPosition[]
    PointLightColor[]
    PointLightPosition[]
    PointLightLinearDecay[]
    PointLightQuadraticDecay[]
    MatSpecularColor
    MatShininess
*****/
void phongModel(vec4 position, vec3 normal, vec3 camDir, vec3 matAmbient, vec3 matDiffuse, out vec3 ambdiff, out vec3 spec) {

    vec3 ambientTotal  = vec3(0.0);
    vec3 diffuseTotal  = vec3(0.0);
    vec3 specularTotal = vec3(0.0);

    bool noLights = true;
    const float EPS = 0.00001;

    float specular;

#if AMB_LIGHTS>0
    noLights = false;
    // Ambient lights
    for (int i = 0; i < AMB_LIGHTS; ++i) {
        ambientTotal += AmbientLightColor[i] * matAmbient;
    }
#endif

#if DIR_LIGHTS>0
    noLights = false;
    // Directional lights
    for (int i = 0; i < DIR_LIGHTS; ++i) {
        vec3 lightDirection = normalize(DirLightPosition(i)); // Vector from fragment to light source
        float dotNormal = dot(lightDirection, normal); // Dot product between light direction and fragment normal
        if (dotNormal > EPS) { // If the fragment is lit
            diffuseTotal += DirLightColor(i) * matDiffuse * dotNormal;

#ifdef BLINN
            specular = pow(max(dot(normal, normalize(lightDirection + camDir)), 0.0), MatShininess);
#else
            specular = pow(max(dot(reflect(-lightDirection, normal), camDir), 0.0), MatShininess);
#endif
            specularTotal += DirLightColor(i) * MatSpecularColor * specular;
        }
    }
#endif

#if POINT_LIGHTS>0
    noLights = false;
    // Point lights
    for (int i = 0; i < POINT_LIGHTS; ++i) {
        vec3 lightDirection = PointLightPosition(i) - vec3(position); // Vector from fragment to light source
        float lightDistance = length(lightDirection); // Distance from fragment to light source
        lightDirection = lightDirection / lightDistance; // Normalize lightDirection
        float dotNormal = dot(lightDirection, normal);  // Dot product between light direction and fragment normal
        if (dotNormal > EPS) { // If the fragment is lit
            float attenuation = 1.0 / (1.0 + lightDistance * (PointLightLinearDecay(i) + PointLightQuadraticDecay(i) * lightDistance));
            vec3 attenuatedColor = PointLightColor(i) * attenuation;
            diffuseTotal += attenuatedColor * matDiffuse * dotNormal;

#ifdef BLINN
            specular = pow(max(dot(normal, normalize(lightDirection + camDir)), 0.0), MatShininess);
#else
            specular = pow(max(dot(reflect(-lightDirection, normal), camDir), 0.0), MatShininess);
#endif
            specularTotal += attenuatedColor * MatSpecularColor * specular;
        }
    }
#endif

#if SPOT_LIGHTS>0
    noLights = false;
    for (int i = 0; i < SPOT_LIGHTS; ++i) {
        // Calculates the direction and distance from the current vertex to this spot light.
        vec3 lightDirection = SpotLightPosition(i) - vec3(position); // Vector from fragment to light source
        float lightDistance = length(lightDirection); // Distance from fragment to light source
        lightDirection = lightDirection / lightDistance; // Normalize lightDirection
        float angleDot = dot(-lightDirection, SpotLightDirection(i));
        float angle = acos(angleDot);
        float cutoff = radians(clamp(SpotLightCutoffAngle(i), 0.0, 90.0));
        if (angle < cutoff) { // Check if fragment is inside spotlight beam
            float dotNormal = dot(lightDirection, normal); // Dot product between light direction and fragment normal
            if (dotNormal > EPS) { // If the fragment is lit
                float attenuation = 1.0 / (1.0 + lightDistance * (SpotLightLinearDecay(i) + SpotLightQuadraticDecay(i) * lightDistance));
                float spotFactor = pow(angleDot, SpotLightAngularDecay(i));
                vec3 attenuatedColor = SpotLightColor(i) * attenuation * spotFactor;
                diffuseTotal += attenuatedColor * matDiffuse * dotNormal;

#ifdef BLINN
                specular = pow(max(dot(normal, normalize(lightDirection + camDir)), 0.0), MatShininess);
#else
                specular = pow(max(dot(reflect(-lightDirection, normal), camDir), 0.0), MatShininess);
#endif
                specularTotal += attenuatedColor * MatSpecularColor * specular;
            }
        }
    }
#endif
    if (noLights) {
        diffuseTotal = matDiffuse;
    }
    // Sets output colors
    ambdiff = ambientTotal + MatEmissiveColor + diffuseTotal;
    spec = specularTotal;
}
