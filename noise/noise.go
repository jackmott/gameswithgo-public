package noise

// Experiment! draw some crazy stuff!
// Gist it next week and I'll show it off on stream

import (
	"math"
	"runtime"
	"sync"
)

// NoiseType indicates which noise MakeNoise will generate
type NoiseType int

const (
	FBM NoiseType = iota
	TURBULENCE
)

// Turbulence generates turbulent fractal noise
func Turbulence(x, y, frequency, lacunarity, gain float32, octaves int) float32 {
	var sum float32
	amplitude := float32(1)

	for i := 0; i < octaves; i++ {
		f := Snoise2(x*frequency, y*frequency) * amplitude
		if f < 0 {
			f = -1.0 * f
		}
		sum += f
		frequency *= lacunarity
		amplitude *= gain
	}
	return sum
}

// Fbm2 generates fractal brownian motion
func Fbm2(x, y, frequency, lacunarity, gain float32, octaves int) float32 {
	var sum float32
	amplitude := float32(1.0)
	for i := 0; i < octaves; i++ {
		sum += Snoise2(x*frequency, y*frequency) * amplitude
		frequency = frequency * lacunarity
		amplitude = amplitude * gain
	}
	return sum
}

// MakeNoise generates a 2d block of noise
func MakeNoise(noiseType NoiseType, frequency, lacunarity, gain float32, octaves, w, h int) (noise []float32, min, max float32) {

	noise = make([]float32, w*h)
	numRoutines := runtime.NumCPU()
	var wg sync.WaitGroup
	wg.Add(numRoutines)
	batchSize := len(noise) / numRoutines

	minChan := make(chan float32, numRoutines)
	maxChan := make(chan float32, numRoutines)

	for i := 0; i < numRoutines; i++ {
		go func(i int) {
			defer wg.Done()
			innerMin := float32(math.MaxFloat32)
			innerMax := float32(-math.MaxFloat32)
			start := i * batchSize
			end := start + batchSize - 1
			for j := start; j < end; j++ {
				x := j % w
				y := (j - x) / h
				if noiseType == TURBULENCE {
					noise[j] = Turbulence(float32(x), float32(y), frequency, lacunarity, gain, octaves)
				} else if noiseType == FBM {
					noise[j] = Fbm2(float32(x), float32(y), frequency, lacunarity, gain, octaves)
				}

				if noise[j] < innerMin {
					innerMin = noise[j]
				} else if noise[j] > innerMax {
					innerMax = noise[j]
				}

			}
			minChan <- innerMin
			maxChan <- innerMax
		}(i)
	}

	wg.Wait()
	close(minChan)
	close(maxChan)

	min = float32(math.MaxFloat32)
	for currentMin := range minChan {
		if currentMin < min {
			min = currentMin
		}
	}

	max = float32(-math.MaxFloat32)
	for currentMax := range maxChan {
		if currentMax > max {
			max = currentMax
		}
	}

	return noise, min, max
}

/* This code ported to Go from Stefan Gustavson's C implementation, his comments follow:
 * https://github.com/stegu/perlin-noise/blob/master/src/simplexnoise1234.c
 * SimplexNoise1234, Simplex noise with true analytic
 * derivative in 1D to 4D.
 *
 * Author: Stefan Gustavson, 2003-2005
 * Contact: stefan.gustavson@liu.se
 *
 *
 * This code was GPL licensed until February 2011.
 * As the original author of this code, I hereby
 * release it into the public domain.
 * Please feel free to use it for whatever you want.
 * Credit is appreciated where appropriate, and I also
 * appreciate being told where this code finds any use,
 * but you may do as you like.
 */

/*
 * This implementation is "Simplex Noise" as presented by
 * Ken Perlin at a relatively obscure and not often cited course
 * session "Real-Time Shading" at Siggraph 2001 (before real
 * time shading actually took off), under the title "hardware noise".
 * The 3D function is numerically equivalent to his Java reference
 * code available in the PDF course notes, although I re-implemented
 * it from scratch to get more readable code. The 1D, 2D and 4D cases
 * were implemented from scratch by me from Ken Perlin's text.
 *
 * This file has no dependencies on any other file, not even its own
 * header file. The header file is made for use by external code only.
 */

func fastFloor(x float32) int {
	if float32(int(x)) <= x {
		return int(x)
	}
	return int(x) - 1
}

// Static data

/*
 * Permutation table. This is just a random jumble of all numbers 0-255
 * This needs to be exactly the same for all instances on all platforms,
 * so it's easiest to just keep it as static explicit data.
 * This also removes the need for any initialisation of this class.
 *
 */
var perm = [256]uint8{151, 160, 137, 91, 90, 15,
	131, 13, 201, 95, 96, 53, 194, 233, 7, 225, 140, 36, 103, 30, 69, 142, 8, 99, 37, 240, 21, 10, 23,
	190, 6, 148, 247, 120, 234, 75, 0, 26, 197, 62, 94, 252, 219, 203, 117, 35, 11, 32, 57, 177, 33,
	88, 237, 149, 56, 87, 174, 20, 125, 136, 171, 168, 68, 175, 74, 165, 71, 134, 139, 48, 27, 166,
	77, 146, 158, 231, 83, 111, 229, 122, 60, 211, 133, 230, 220, 105, 92, 41, 55, 46, 245, 40, 244,
	102, 143, 54, 65, 25, 63, 161, 1, 216, 80, 73, 209, 76, 132, 187, 208, 89, 18, 169, 200, 196,
	135, 130, 116, 188, 159, 86, 164, 100, 109, 198, 173, 186, 3, 64, 52, 217, 226, 250, 124, 123,
	5, 202, 38, 147, 118, 126, 255, 82, 85, 212, 207, 206, 59, 227, 47, 16, 58, 17, 182, 189, 28, 42,
	223, 183, 170, 213, 119, 248, 152, 2, 44, 154, 163, 70, 221, 153, 101, 155, 167, 43, 172, 9,
	129, 22, 39, 253, 19, 98, 108, 110, 79, 113, 224, 232, 178, 185, 112, 104, 218, 246, 97, 228,
	251, 34, 242, 193, 238, 210, 144, 12, 191, 179, 162, 241, 81, 51, 145, 235, 249, 14, 239, 107,
	49, 192, 214, 31, 181, 199, 106, 157, 184, 84, 204, 176, 115, 121, 50, 45, 127, 4, 150, 254,
	138, 236, 205, 93, 222, 114, 67, 29, 24, 72, 243, 141, 128, 195, 78, 66, 215, 61, 156, 180}

//---------------------------------------------------------------------

func grad2(hash uint8, x, y float32) float32 {
	h := hash & 7 // Convert low 3 bits of hash code
	u := y
	v := 2 * x
	if h < 4 {
		u = x
		v = 2 * y
	} // into 8 simple gradient directions,
	// and compute the dot product with (x,y).

	if h&1 != 0 {
		u = -u
	}
	if h&2 != 0 {
		v = -v
	}
	return u + v
}

// 2D simplex noise
func Snoise2(x, y float32) float32 {

	const F2 float32 = 0.366025403 // F2 = 0.5*(sqrt(3.0)-1.0)
	const G2 float32 = 0.211324865 // G2 = (3.0-Math.sqrt(3.0))/6.0

	var n0, n1, n2 float32 // Noise contributions from the three corners

	// Skew the input space to determine which simplex cell we're in
	s := (x + y) * F2 // Hairy factor for 2D
	xs := x + s
	ys := y + s
	i := fastFloor(xs)
	j := fastFloor(ys)

	t := float32(i+j) * G2
	X0 := float32(i) - t // Unskew the cell origin back to (x,y) space
	Y0 := float32(j) - t
	x0 := x - X0 // The x,y distances from the cell origin
	y0 := y - Y0

	// For the 2D case, the simplex shape is an equilateral triangle.
	// Determine which simplex we are in.
	var i1, j1 uint8 // Offsets for second (middle) corner of simplex in (i,j) coords
	if x0 > y0 {
		i1 = 1
		j1 = 0
	} else { // lower triangle, XY order: (0,0)->(1,0)->(1,1)
		i1 = 0
		j1 = 1
	} // upper triangle, YX order: (0,0)->(0,1)->(1,1)

	// A step of (1,0) in (i,j) means a step of (1-c,-c) in (x,y), and
	// a step of (0,1) in (i,j) means a step of (-c,1-c) in (x,y), where
	// c = (3-sqrt(3))/6

	x1 := x0 - float32(i1) + G2 // Offsets for middle corner in (x,y) unskewed coords
	y1 := y0 - float32(j1) + G2
	x2 := x0 - 1.0 + 2.0*G2 // Offsets for last corner in (x,y) unskewed coords
	y2 := y0 - 1.0 + 2.0*G2

	// Wrap the integer indices at 256, to avoid indexing perm[] out of bounds
	ii := uint8(i)
	jj := uint8(j)

	// Calculate the contribution from the three corners
	t0 := 0.5 - x0*x0 - y0*y0
	if t0 < 0.0 {
		n0 = 0.0
	} else {
		t0 *= t0
		n0 = t0 * t0 * grad2(perm[ii+perm[jj]], x0, y0)
	}

	t1 := 0.5 - x1*x1 - y1*y1
	if t1 < 0.0 {
		n1 = 0.0
	} else {
		t1 *= t1
		n1 = t1 * t1 * grad2(perm[ii+i1+perm[jj+j1]], x1, y1)
	}

	t2 := 0.5 - x2*x2 - y2*y2
	if t2 < 0.0 {
		n2 = 0.0
	} else {
		t2 *= t2
		n2 = t2 * t2 * grad2(perm[ii+1+perm[jj+1]], x2, y2)
	}

	// Add contributions from each corner to get the final noise value.
	return (n0 + n1 + n2)
}
