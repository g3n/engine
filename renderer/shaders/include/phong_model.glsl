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

#if AMB_LIGHTS>0
    // Ambient lights
    for (int i = 0; i < AMB_LIGHTS; ++i) {
        ambientTotal += AmbientLightColor[i] * matAmbient;
    }
#endif

#if DIR_LIGHTS>0
    // Directional lights
    for (int i = 0; i < DIR_LIGHTS; ++i) {
        vec3 lightDirection = normalize(DirLightPosition(i)); // Vector from fragment to light source
        float dotNormal = max(dot(lightDirection, normal), 0.0); // Dot product between light direction and fragment normal
        if (dotNormal > 0.0) { // If the fragment is lit
            diffuseTotal += DirLightColor(i) * matDiffuse * dotNormal;
            specularTotal += DirLightColor(i) * MatSpecularColor * pow(max(dot(reflect(-lightDirection, normal), camDir), 0.0), MatShininess);
        }
    }
#endif

#if POINT_LIGHTS>0
    // Point lights
    for (int i = 0; i < POINT_LIGHTS; ++i) {
        vec3 lightDirection = PointLightPosition(i) - vec3(position); // Vector from fragment to light source
        float lightDistance = length(lightDirection); // Distance from fragment to light source
        lightDirection = lightDirection / lightDistance; // Normalize lightDirection
        float dotNormal = max(dot(lightDirection, normal), 0.0);  // Dot product between light direction and fragment normal
        if (dotNormal > 0.0) { // If the fragment is lit
            float attenuation = 1.0 / (1.0 + PointLightLinearDecay(i) * lightDistance + PointLightQuadraticDecay(i) * lightDistance * lightDistance);
            vec3 attenuatedColor = PointLightColor(i) * attenuation;
            diffuseTotal += attenuatedColor * matDiffuse * dotNormal;
            specularTotal += attenuatedColor * MatSpecularColor * pow(max(dot(reflect(-lightDirection, normal), camDir), 0.0), MatShininess);
        }
    }
#endif

#if SPOT_LIGHTS>0
    for (int i = 0; i < SPOT_LIGHTS; ++i) {
        // Calculates the direction and distance from the current vertex to this spot light.
        vec3 lightDirection = SpotLightPosition(i) - vec3(position); // Vector from fragment to light source
        float lightDistance = length(lightDirection); // Distance from fragment to light source
        lightDirection = lightDirection / lightDistance; // Normalize lightDirection
        float angleDot = dot(-lightDirection, SpotLightDirection(i));
        float angle = acos(angleDot);
        float cutoff = radians(clamp(SpotLightCutoffAngle(i), 0.0, 90.0));
        if (angle < cutoff) { // Check if fragment is inside spotlight beam
            float dotNormal = max(dot(lightDirection, normal), 0.0); // Dot product between light direction and fragment normal
            if (dotNormal > 0.0) { // If the fragment is lit
                float attenuation = 1.0 / (1.0 + SpotLightLinearDecay(i) * lightDistance + SpotLightQuadraticDecay(i) * lightDistance * lightDistance);
                float spotFactor = pow(angleDot, SpotLightAngularDecay(i));
                vec3 attenuatedColor = SpotLightColor(i) * attenuation * spotFactor;
                diffuseTotal += attenuatedColor * matDiffuse * dotNormal;
                specularTotal += attenuatedColor * MatSpecularColor * pow(max(dot(reflect(-lightDirection, normal), camDir), 0.0), MatShininess);
            }
        }
    }
#endif

    // Sets output colors
    ambdiff = ambientTotal + MatEmissiveColor + diffuseTotal;
    spec = specularTotal;
}
