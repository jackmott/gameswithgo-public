#version 330 core
out vec4 FragColor;


in vec2 TexCoord;
in vec3 Normal;
in vec3 FragPos;



uniform sampler2D texture1;
uniform vec3 lightPos;
uniform vec3 lightColor; //color/brightness


void main() {
        
    // diffuse
    vec3 lightDir = normalize(lightPos - FragPos);
    float diff = max(dot(Normal,lightDir),0.0);
    vec3 diffuse = diff * lightColor;
    vec3 ambient = lightColor * 0.3;
    FragColor = vec4(ambient + diffuse,1.0) * texture(texture1,TexCoord);    
}