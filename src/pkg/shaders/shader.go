package shaders

// VertexShader - the vertext shader for our view
var VertexShader = `
#version 330 core

in vec2 position;
in vec2 texture;

out vec2 Texture;

void main() {
	gl_Position = vec4(position.x, position.y, 0.0, 1.0);
	Texture = texture;
}
`

// FragmentShader - the fragment shader for our view
var FragmentShader = `
#version 330 core

in vec2 Texture;
out vec4 color;

uniform vec2 resolution;
uniform float time;

uniform sampler2D tex;

float rand(vec2 co) {
    return fract( sin( dot( co.xy, vec2(12.9898, 78.233) ) ) * 43758.5453);
}

void main(void) 
{
    // div by resolution
    
    //vec2 q = gl_FragCoord.xy / vec2(64 * 5, 32 * 5); 

    vec2 q = gl_FragCoord.xy / resolution;           //vec2(64 * 12, 32 * 12); // scale * 2

    // zoom in/out 
    vec2 uv = 0.5 + ( q - 0.5 ) * (0.98 + 0.006 * sin(0.9 * time));

	vec3 oric = texture(tex, vec2(q.x, 1.0 - q.y)).xyz;	
    vec3 col = oric;

    // contrast
    col = clamp( (col * 0.5) + (0.5 * col * col * 1.2), 0.0, 1.0);

    // vignette
    col *= 0.6 + 0.4 * 16.0 * uv.x * uv.y * (1.0 - uv.x) * (1.0 - uv.y);

    // color tint
    col *= vec3(0.9, 1.0, 0.7);

	// scanlines
	float crawlSpeed = 1.0f;
	float crawlSize  = 500.0f;
    col *= 0.8 + (0.2 * sin( (crawlSpeed * time) + (uv.y * crawlSize) ));

    // flicker
    col *= 1.0 - (0.07 * rand( vec2(time, tan(time) ) ) );

    //smoothen
    float comp = smoothstep( 0.2, 0.7, sin(time) );
	col = mix( col, oric, clamp( -2.0 + (2.0 * q.x) + (3.0 * comp), 0.0, 1.0) );
	color = vec4(col, 1.0);
}
`
