#version 330 core
layout (location = 0) in vec3 aPos;
layout (location = 1) in vec2 aTexCoord;
layout (location = 2) in vec3 aNormal;

out vec2 TexCoord;

out vec3 FragPos;
out vec3 Normal;

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;

void main() 
{
    Normal = aNormal;
    FragPos = vec3(model * vec4(aPos,1.0));
    gl_Position = projection * view * vec4(FragPos,1.0f);
    TexCoord = vec2(aTexCoord.x,1.0 - aTexCoord.y);
}