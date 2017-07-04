// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package shader

func init() {
	AddChunk("lights", chunkLights)
}

const chunkLights = `
{{if .AmbientLightsMax}}
// Ambient lights uniforms
uniform vec3 AmbientLightColor[{{.AmbientLightsMax}}];
{{end}}

{{if .DirLightsMax}}
// Directional lights uniform array. Each directional light uses 2 elements
uniform vec3  DirLight[2*{{.DirLightsMax}}];

// Macros to access elements inside the DirectionalLight uniform array
#define DirLightColor(a)		DirLight[2*a]
#define DirLightPosition(a)		DirLight[2*a+1]
{{end}}

{{if .PointLightsMax}}
// Point lights uniform array. Each point light uses 3 elements
uniform vec3  PointLight[3*{{.PointLightsMax}}];

// Macros to access elements inside the PointLight uniform array
#define PointLightColor(a)			PointLight[3*a]
#define PointLightPosition(a)		PointLight[3*a+1]
#define PointLightLinearDecay(a)	PointLight[3*a+2].x
#define PointLightQuadraticDecay(a)	PointLight[3*a+2].y
{{end}}

{{if .SpotLightsMax}}
// Spot lights uniforms. Each spot light uses 5 elements
uniform vec3  SpotLight[5*{{.SpotLightsMax}}];

// Macros to access elements inside the PointLight uniform array
#define SpotLightColor(a)			SpotLight[5*a]
#define SpotLightPosition(a)		SpotLight[5*a+1]
#define SpotLightDirection(a)		SpotLight[5*a+2]
#define SpotLightAngularDecay(a)	SpotLight[5*a+3].x
#define SpotLightCutoffAngle(a)		SpotLight[5*a+3].y
#define SpotLightLinearDecay(a)		SpotLight[5*a+3].z
#define SpotLightQuadraticDecay(a)	SpotLight[5*a+4].x
{{end}}
`
