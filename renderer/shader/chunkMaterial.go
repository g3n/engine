// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package shader

func init() {
	AddChunk("material", chunkMaterial)
}

const chunkMaterial = `
// Material uniforms
uniform vec3	Material[6];

// Macros to access elements inside the Material uniform array
#define MatAmbientColor		Material[0]
#define MatDiffuseColor		Material[1]
#define MatSpecularColor	Material[2]
#define MatEmissiveColor	Material[3]
#define MatShininess		Material[4].x
#define MatOpacity			Material[4].y
#define MatPointSize		Material[4].z
#define MatPointRotationZ	Material[5].x

{{if .MatTexturesMax}}
// Textures uniforms
uniform sampler2D	MatTexture[{{.MatTexturesMax}}];
uniform mat3		MatTexinfo[{{.MatTexturesMax}}];

// Macros to access elements inside MatTexinfo uniform
#define MatTexOffset(a)		MatTexinfo[a][0].xy
#define MatTexRepeat(a)		MatTexinfo[a][1].xy
#define MatTexFlipY(a)		bool(MatTexinfo[a][2].x)
#define MatTexVisible(a)	bool(MatTexinfo[a][2].y)
{{ end }}
`
