//
// Fragment Shader template
//

precision highp float;

in vec3 Color;
out vec4 FragColor;

void main() {

    FragColor = vec4(Color, 1.0);
}

