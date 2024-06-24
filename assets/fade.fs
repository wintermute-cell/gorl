#version 330

out vec4 fragColor;

uniform vec2 BLUR;

uniform sampler2D texture0;
in vec2 fragTexCoord;

void main( void )
{
    vec4 col = texture2D( texture0, fragTexCoord );
    col.a = col.a * 0.95;
    fragColor = col;
}   
