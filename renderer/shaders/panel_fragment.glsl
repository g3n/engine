//
// Fragment Shader template
//

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
uniform vec4 Panel[9];
#define Bounds			Panel[0]		  // panel bounds in texture coordinates
#define Border			Panel[1]		  // panel border in texture coordinates
#define Padding			Panel[2]		  // panel padding in texture coordinates
#define Content			Panel[3]		  // panel content area in texture coordinates
#define BorderColor		Panel[4]		  // panel border color
#define PaddingColor	Panel[5]		  // panel padding color
#define ContentColor	Panel[6]		  // panel content color
#define Roundness 	    Panel[7]		  // panel corner roundness
#define TextureValid	bool(Panel[8].x)  // texture valid flag
#define AspectRatio  	Panel[8].y        // panel aspect ratio

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

    if (Roundness == 0) {
        return true;
    }

    // Adjust fragment x coordinate multiplying by the aspect ratio
    float fragx = FragTexcoord.x * AspectRatio;
    vec2 frag = vec2(fragx, FragTexcoord.y);

    // Top left corner
    float radius = rect[3] * Roundness[0] / 2;
    float rx = rect[0]*AspectRatio + radius;
    float ry = rect[1] + radius;
    if (fragx <= rx && FragTexcoord.y <= ry) {
        vec2 center = vec2(rx, ry);
        float dist = distance(frag, center);
        if (dist < radius) {
            return true;
        }
        return false;
    }

    // Bottom left corner
    radius = rect[3] * Roundness[3] / 2;
    rx = rect[0]*AspectRatio + radius;
    ry = rect[1] + rect[3] - radius;
    if (fragx <= rx && FragTexcoord.y >= ry) {
        vec2 center = vec2(rx, ry);
        float dist = distance(frag, center);
        if (dist < radius) {
            return true;
        }
        return false;
    }

    // Top right corner
    radius = rect[2] * Roundness[1] / 2;
    rx = (rect[0] + rect[2])*AspectRatio - radius;
    ry = rect[1] + radius;
    if (fragx >= rx && FragTexcoord.y <= ry) {
        vec2 center = vec2(rx, ry);
        float dist = distance(frag, center);
        if (dist < radius) {
            return true;
        }
        return false;
    }

    // Bottom right corner
    radius = rect[3] * Roundness[2] / 2;
    rx = (rect[0] + rect[2])*AspectRatio - radius;
    ry = rect[1] + rect[3] - radius;
    if (fragx >= rx && FragTexcoord.y >= ry) {
        vec2 center = vec2(rx, ry);
        float dist = distance(frag, center);
        if (dist < radius) {
            return true;
        }
        return false;
    }

    // Fragment is inside the inner rectangle
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
            vec2 factor = vec2(1/Content[2], 1/Content[3]);
            vec2 texcoord = (FragTexcoord + offset) * factor;
            vec4 texColor = texture(MatTexture, texcoord * MatTexRepeat + MatTexOffset);
            // Mix content color with texture color ???
            //color = mix(color, texColor, texColor.a);
            color = texColor;
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

