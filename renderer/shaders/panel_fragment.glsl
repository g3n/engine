//
// Fragment Shader template
//

precision highp float;

// Texture uniforms
uniform sampler2D	MatTexture;
uniform vec2		MatTexinfo[3];

// Macros to access elements inside the MatTexinfo array
#define MatTexOffset		MatTexinfo[0]
#define MatTexRepeat		MatTexinfo[1]
#define MatTexFlipY	    	bool(MatTexinfo[2].x) // not used
#define MatTexVisible	    bool(MatTexinfo[2].y) // not used

// Inputs from vertex shader
in vec2 FragTexcoord;

// Input uniform
uniform vec4 Panel[8];
#define Bounds			Panel[0]		  // panel bounds in texture coordinates
#define Border			Panel[1]		  // panel border in texture coordinates
#define Padding			Panel[2]		  // panel padding in texture coordinates
#define Content			Panel[3]		  // panel content area in texture coordinates
#define BorderColor		Panel[4]		  // panel border color
#define PaddingColor	Panel[5]		  // panel padding color
#define ContentColor	Panel[6]		  // panel content color
#define TextureValid	bool(Panel[7].x)  // texture valid flag

// Output
out vec4 FragColor;


/***
* Checks if current fragment texture coordinate is inside the
* supplied rectangle in texture coordinates:
* rect[0] - position x [0,1]
* rect[1] - position y [0,1]
* rect[2] - width [0,1]
* rect[3] - height [0,1]
*/
bool checkRect(vec4 rect) {

    if (FragTexcoord.x < rect[0]) {
        return false;
    }
    if (FragTexcoord.x > rect[0] + rect[2]) {
        return false;
    }
    if (FragTexcoord.y < rect[1]) {
        return false;
    }
    if (FragTexcoord.y > rect[1] + rect[3]) {
        return false;
    }
    return true;
}


void main() {

    // Discard fragment outside of received bounds
    // Bounds[0] - xmin
    // Bounds[1] - ymin
    // Bounds[2] - xmax
    // Bounds[3] - ymax
    if (FragTexcoord.x <= Bounds[0] || FragTexcoord.x >= Bounds[2]) {
        discard;
    }
    if (FragTexcoord.y <= Bounds[1] || FragTexcoord.y >= Bounds[3]) {
        discard;
    }

    // Check if fragment is inside content area
    if (checkRect(Content)) {

        // If no texture, the color will be the material color.
        vec4 color = ContentColor;

		if (TextureValid) {
            // Adjust texture coordinates to fit texture inside the content area
            vec2 offset = vec2(-Content[0], -Content[1]);
            vec2 factor = vec2(1.0/Content[2], 1.0/Content[3]);
            vec2 texcoord = (FragTexcoord + offset) * factor;
            vec4 texColor = texture(MatTexture, texcoord * MatTexRepeat + MatTexOffset);

            // Mix content color with texture color.
            // Note that doing a simple linear interpolation (e.g. using mix()) is not correct!
            // The right formula can be found here: https://en.wikipedia.org/wiki/Alpha_compositing#Alpha_blending
            // For a more in-depth discussion: http://apoorvaj.io/alpha-compositing-opengl-blending-and-premultiplied-alpha.html#toc4

            // Pre-multiply the content color
            vec4 contentPre = ContentColor;
            contentPre.rgb *= contentPre.a;

            // Pre-multiply the texture color
            vec4 texPre = texColor;
            texPre.rgb *= texPre.a;

            // Combine colors the premultiplied final color
            color = texPre + contentPre * (1.0 - texPre.a);

            // Un-pre-multiply (pre-divide? :P)
            color.rgb /= color.a;
		}

        FragColor = color;
        return;
    }

    // Checks if fragment is inside paddings area
    if (checkRect(Padding)) {
        FragColor = PaddingColor;
        return;
    }

    // Checks if fragment is inside borders area
    if (checkRect(Border)) {
        FragColor = BorderColor;
        return;
    }

    // Fragment is in margins area (always transparent)
    FragColor = vec4(1,1,1,0);
}

