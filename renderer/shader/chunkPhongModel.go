// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package shader

func init() {
	AddChunk("phong_model", chunkPhongModel)
}

const chunkPhongModel = `
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
    PointLightPosition[];
    PointLightLinearDecay[];
    PointLightQuadraticDecay[];
    MatSpecularColor
    MatShininess
*/
void phongModel(vec4 position, vec3 normal, vec3 camDir, vec3 matAmbient, vec3 matDiffuse, out vec3 ambdiff, out vec3 spec) {

    vec3 ambientTotal  = vec3(0.0);
    vec3 diffuseTotal  = vec3(0.0);
    vec3 specularTotal = vec3(0.0);

    {{if .AmbientLightsMax }}
    for (int i = 0; i < {{.AmbientLightsMax}}; i++) {
        ambientTotal += AmbientLightColor[i] * matAmbient;
    }
    {{ end }}

    {{if .DirLightsMax }}
    for (int i = 0; i < {{.DirLightsMax}}; i++) {
        // Diffuse reflection
        // DirLightPosition is the direction of the current light
        vec3 lightDirection = normalize(DirLightPosition[i]);
        // Calculates the dot product between the light direction and this vertex normal.
        float dotNormal = max(dot(lightDirection, normal), 0.0);
        diffuseTotal += DirLightColor[i] * matDiffuse * dotNormal;

        // Specular reflection
        // Calculates the light reflection vector 
        vec3 ref = reflect(-lightDirection, normal);
        if (dotNormal > 0.0) {
            specularTotal += DirLightColor[i] * MatSpecularColor * pow(max(dot(ref, camDir), 0.0), MatShininess);
        }
    }
    {{ end }}

    {{if .PointLightsMax }}
    for (int i = 0; i < {{.PointLightsMax}}; i++) {
        // Calculates the direction and distance from the current vertex to this point light.
        vec3 lightDirection = PointLightPosition[i] - vec3(position);
        float lightDistance = length(lightDirection);
        // Normalizes the lightDirection
        lightDirection = lightDirection / lightDistance;
        // Calculates the attenuation due to the distance of the light
        float attenuation = 1.0 / (1.0 + PointLightLinearDecay[i] * lightDistance +
            PointLightQuadraticDecay[i] * lightDistance * lightDistance);

        // Diffuse reflection
        float dotNormal = max(dot(lightDirection, normal), 0.0);
        diffuseTotal += PointLightColor[i] * matDiffuse * dotNormal * attenuation;
        
        // Specular reflection
        // Calculates the light reflection vector 
        vec3 ref = reflect(-lightDirection, normal);
        if (dotNormal > 0.0) {
            specularTotal += PointLightColor[i] * MatSpecularColor *
                pow(max(dot(ref, camDir), 0.0), MatShininess) * attenuation;
        }
    }
    {{ end }}

    {{if .SpotLightsMax }}
    for (int i = 0; i < {{.SpotLightsMax}}; i++) {
        // Calculates the direction and distance from the current vertex to this spot light.
        vec3 lightDirection = SpotLightPosition[i] - vec3(position);
        float lightDistance = length(lightDirection);
        lightDirection = lightDirection / lightDistance;

        // Calculates the attenuation due to the distance of the light
        float attenuation = 1.0 / (1.0 + SpotLightLinearDecay[i] * lightDistance +
            SpotLightQuadraticDecay[i] * lightDistance * lightDistance);

        // Calculates the angle between the vertex direction and spot direction
        // If this angle is greater than the cutoff the spotlight will not contribute
        // to the final color.
        float angle = acos(dot(-lightDirection, SpotLightDirection[i]));
        float cutoff = radians(clamp(SpotLightCutoffAngle[i], 0.0, 90.0));
        if (angle >= cutoff) {
            continue;
        }
        float spotFactor = pow(dot(-lightDirection, SpotLightDirection[i]), SpotLightAngularDecay[i]);

        // Diffuse reflection
        float dotNormal = max(dot(lightDirection, normal), 0.0);
        diffuseTotal += SpotLightColor[i] * matDiffuse * dotNormal * attenuation * spotFactor;

        // Specular reflection
        vec3 ref = reflect(-lightDirection, normal);
        if (dotNormal > 0.0) {
            specularTotal += SpotLightColor[i] * MatSpecularColor * pow(max(dot(ref, camDir), 0.0), MatShininess) * attenuation * spotFactor;
        }
    }
    {{ end }}

    // Sets output colors
    ambdiff = ambientTotal + MatEmissiveColor + diffuseTotal;
    spec = specularTotal;
}
`
