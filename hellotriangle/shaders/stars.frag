#version 330 core

#define OCTAVES 5
out vec4 FragColor;
in vec2 TexCoord;

uniform float x;
uniform float y;

vec3 mod289(vec3 x)
{
  return x - floor(x * (1.0 / 289.0)) * 289.0;
}

vec2 mod289(vec2 x) 
{
  return x - floor(x * (1.0 / 289.0)) * 289.0;
}

vec3 permute(vec3 x) 
{
  return mod289(((x*34.0)+1.0)*x);
}

// Simplex noise 
// https://github.com/ashima/webgl-noise 
// Copyright (C) 2011 Ashima Arts. All rights reserved.
float snoise(vec2 v)
  {
  const vec4 C = vec4(0.211324865405187,  // (3.0-sqrt(3.0))/6.0
                      0.366025403784439,  // 0.5*(sqrt(3.0)-1.0)
                     -0.577350269189626,  // -1.0 + 2.0 * C.x
                      0.024390243902439); // 1.0 / 41.0
// First corner
  vec2 i  = floor(v + dot(v, C.yy) );
  vec2 x0 = v -   i + dot(i, C.xx);

// Other corners
  vec2 i1;
  //i1.x = step( x0.y, x0.x ); // x0.x > x0.y ? 1.0 : 0.0
  //i1.y = 1.0 - i1.x;
  i1 = (x0.x > x0.y) ? vec2(1.0, 0.0) : vec2(0.0, 1.0);
  // x0 = x0 - 0.0 + 0.0 * C.xx ;
  // x1 = x0 - i1 + 1.0 * C.xx ;
  // x2 = x0 - 1.0 + 2.0 * C.xx ;
  vec4 x12 = x0.xyxy + C.xxzz;
  x12.xy -= i1;

// Permutations
  i = mod289(i); // Avoid truncation effects in permutation
  vec3 p = permute( permute( i.y + vec3(0.0, i1.y, 1.0 ))
		+ i.x + vec3(0.0, i1.x, 1.0 ));

  vec3 m = max(0.5 - vec3(dot(x0,x0), dot(x12.xy,x12.xy), dot(x12.zw,x12.zw)), 0.0);
  m = m*m ;
  m = m*m ;

// Gradients: 41 points uniformly over a line, mapped onto a diamond.
// The ring size 17*17 = 289 is close to a multiple of 41 (41*7 = 287)

  vec3 x = 2.0 * fract(p * C.www) - 1.0;
  vec3 h = abs(x) - 0.5;
  vec3 ox = floor(x + 0.5);
  vec3 a0 = x - ox;

// Normalise gradients implicitly by scaling m
// Approximation of: m *= inversesqrt( a0*a0 + h*h );
  m *= 1.79284291400159 - 0.85373472095314 * ( a0*a0 + h*h );

// Compute final noise value at P
  vec3 g;
  g.x  = a0.x  * x0.x  + h.x  * x0.y;
  g.yz = a0.yz * x12.xz + h.yz * x12.yw;
  return 130.0 * dot(m, g);
}

vec2 rand2(vec2 p)
{
    p = vec2(dot(p, vec2(12.9898,78.233)), dot(p, vec2(26.65125, 83.054543))); 
    return fract(sin(p) * 43758.5453);
}

float rand(vec2 p)
{
    return fract(sin(dot(p.xy ,vec2(54.90898,18.233))) * 4337.5453);
}

vec3 hsv2rgb(vec3 c)
{
    vec4 K = vec4(1.0, 2.0 / 3.0, 1.0 / 3.0, 3.0);
    vec3 p = abs(fract(c.xxx + K.xyz) * 6.0 - K.www);
    return c.z * mix(K.xxx, clamp(p - K.xxx, 0.0, 1.0), c.y);
}

// Thanks to David Hoskins https://www.shadertoy.com/view/4djGRh
float stars(in vec2 x, float numCells, float size, float br)
{
    vec2 n = x * numCells;
    vec2 f = floor(n);

	float d = 1.0e10;
    for (int i = -1; i <= 1; ++i)
    {
        for (int j = -1; j <= 1; ++j)
        {
            vec2 g = f + vec2(float(i), float(j));
			g = n - g - rand2(mod(g, numCells)) + rand(g);
            // Control size
            g *= 1. / (numCells * size);
			d = min(d, dot(g, g));
        }
    }

    return br * (smoothstep(.95, 1., (1. - sqrt(d))));
}

// Simple fractal noise
// persistence - A multiplier that determines how quickly the amplitudes diminish for 
// each successive octave.
// lacunarity - A multiplier that determines how quickly the frequency increases for 
// each successive octave.
float fractalNoise(in vec2 coord, in float persistence, in float lacunarity)
{    
    float n = 0.;
    float frequency = 1.;
    float amplitude = 1.;
    for (int o = 0; o < OCTAVES; ++o)
    {
        n += amplitude * snoise(coord * frequency);
        amplitude *= persistence;
        frequency *= lacunarity;
    }
    return n;
}

vec3 fractalNebula(in vec2 coord, vec3 color, float transparency)
{
    float n = fractalNoise(coord, .5, 2.);
    return n * color * transparency;
}

void main()
{
    //float resolution = max(iResolution.y, iResolution.y);
    
    vec2 coord = TexCoord.xy;
	vec2 pos = vec2(x,y);

    vec3 result = vec3(0.);

	
    vec2 nebula1pos = coord + pos *.010;
    vec2 nebula2pos = coord + pos * .015;
    vec2 star1pos = coord  + pos * .02;
    vec2 star2pos = coord  + pos *.05;
    vec2 star3pos = coord + pos *.1; 
            
    vec3 nebulaColor1 = hsv2rgb(vec3(.9+.005*sin(nebula1pos.x * 100.0 + nebula1pos.y), 0.5, .25));
	vec3 nebulaColor2 = hsv2rgb(vec3(.1+.001*sin(nebula2pos.x * 60.0 + nebula2pos.y), 1., .25));

    result += fractalNebula(nebula1pos, nebulaColor1, 1.);
    result += fractalNebula(nebula2pos + vec2(5., 7.2), nebulaColor2, .5);
    
    result += stars(star3pos, 4., 0.1, 2.) * vec3(.74, .74, .74);
    result += stars(star2pos, 8., 0.05, 1.) * vec3(.97, .74, .74);
    result += stars(star1pos, 16., 0.025, 0.5) * vec3(.9, .9, .95);
    
    FragColor = vec4(result, 1.);
   
}