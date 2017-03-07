// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package shader

func init() {
	AddChunk("material", chunkMaterial)
}

const chunkMaterial = `
// Material uniforms
uniform vec3      MatAmbientColor;
uniform vec3      MatDiffuseColor;
uniform vec3      MatSpecularColor;
uniform float     MatShininess;
uniform vec3      MatEmissiveColor;
uniform float     MatOpacity;

{{if .MatTexturesMax}}
uniform sampler2D MatTexture[{{.MatTexturesMax}}];
uniform vec2      MatTexRepeat[{{.MatTexturesMax}}];
uniform vec2      MatTexOffset[{{.MatTexturesMax}}];
uniform int       MatTexFlipY[{{.MatTexturesMax}}];
uniform bool      MatTexVisible[{{.MatTexturesMax}}];
{{ end }}
`
