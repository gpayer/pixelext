#version 330 core

out vec4 fragColor;

uniform sampler2D uTexture;
uniform vec4 uTexBounds;

// custom uniforms
uniform float uSpeed;
uniform float uTime;

void main() {
    vec2 t = gl_FragCoord.xy / uTexBounds.zw;
    vec3 influence = texture(uTexture, t).rgb;

    if (influence.r + influence.g + influence.b > 0.3) {
        t.y += cos(t.x * 40.0 + (uTime * uSpeed))*0.005;
        t.x += cos(t.y * 40.0 + (uTime * uSpeed))*0.01;
    }

    vec4 col = texture(uTexture, t);
    fragColor = vec4(col.rgb * vec3(0.6, 0.6, 1.2),col.a);
}
