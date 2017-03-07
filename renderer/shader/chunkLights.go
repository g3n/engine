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
// Directional lights uniforms
uniform vec3  DirLightColor[{{.DirLightsMax}}];
uniform vec3  DirLightPosition[{{.DirLightsMax}}];
{{end}}

{{if .PointLightsMax}}
// Point lights uniforms
uniform vec3  PointLightColor[{{.PointLightsMax}}];
uniform vec3  PointLightPosition[{{.PointLightsMax}}];
uniform float PointLightLinearDecay[{{.PointLightsMax}}];
uniform float PointLightQuadraticDecay[{{.PointLightsMax}}];
{{end}}

{{if .SpotLightsMax}}
// Spot lights uniforms
uniform vec3  SpotLightColor[{{.SpotLightsMax}}];
uniform vec3  SpotLightPosition[{{.SpotLightsMax}}];
uniform vec3  SpotLightDirection[{{.SpotLightsMax}}];
uniform float SpotLightAngularDecay[{{.SpotLightsMax}}];
uniform float SpotLightCutoffAngle[{{.SpotLightsMax}}];
uniform float SpotLightLinearDecay[{{.SpotLightsMax}}];
uniform float SpotLightQuadraticDecay[{{.SpotLightsMax}}];
{{end}}
`
