//
// Physical material fragment shader
//

// Inputs from vertex shader
in vec4 Position;       // Vertex position in camera coordinates.
in vec3 Normal;         // Vertex normal in camera coordinates.
in vec3 CamDir;         // Direction from vertex to camera
in vec2 FragTexcoord;

// Material parameters uniform array
uniform vec4 Material[3];
// Macros to access elements inside the Material array
#define uBaseColor		    Material[0]
#define uEmissiveColor      Material[1]
#define uMetallicFactor     Material[2].x
#define uRoughnessFactor    Material[2].y

#include <lights>

// Final fragment color
out vec4 FragColor;

void main() {


    // Inverts the fragment normal if not FrontFacing
    vec3 fragNormal = Normal;
    if (!gl_FrontFacing) {
        fragNormal = -fragNormal;
    }


    // Final fragment color
    FragColor = uBaseColor;
}


