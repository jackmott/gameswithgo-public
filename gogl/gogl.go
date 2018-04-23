package gogl

import (
	"errors"
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"image/png"
	"io/ioutil"
	"os"
	"strings"
)

type ShaderID uint32
type ProgramID uint32
type BufferID uint32
type TextureID uint32

func MglTest() {
	x := mgl32.NewVecN(2)
	fmt.Println(x)
}

func GetVersion() string {
	return gl.GoStr(gl.GetString(gl.VERSION))
}

// Todo:  SDL2_image  stb_image Sean Barret
func LoadTextureAlpha(filename string) TextureID {
	infile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer infile.Close()

	img, err := png.Decode(infile)
	if err != nil {
		panic(err)
	}

	w := img.Bounds().Max.X
	h := img.Bounds().Max.Y

	pixels := make([]byte, w*h*4)
	bIndex := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			pixels[bIndex] = byte(r / 256)
			bIndex++
			pixels[bIndex] = byte(g / 256)
			bIndex++
			pixels[bIndex] = byte(b / 256)
			bIndex++
			pixels[bIndex] = byte(a / 256)
			bIndex++
		}
	}

	texture := GenBindTexture()
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(w), int32(h), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pixels))
	gl.GenerateMipmap(gl.TEXTURE_2D)
	return texture
}

func GenBindTexture() TextureID {
	var texId uint32
	gl.GenTextures(1, &texId)
	gl.BindTexture(gl.TEXTURE_2D, texId)
	return TextureID(texId)
}

func BindTexture(id TextureID) {
	gl.BindTexture(gl.TEXTURE_2D, uint32(id))
}

func LoadShader(path string, shaderType uint32) (ShaderID, error) {
	shaderFile, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, err
	}
	shaderFileStr := string(shaderFile)
	shaderId, err := CreateShader(shaderFileStr, shaderType)
	if err != nil {
		return 0, err
	}
	return shaderId, nil
}

func CreateShader(shaderSource string, shaderType uint32) (ShaderID, error) {
	shaderId := gl.CreateShader(shaderType)
	shaderSource = shaderSource + "\x00"
	csource, free := gl.Strs(shaderSource)
	gl.ShaderSource(shaderId, 1, csource, nil)
	free()
	gl.CompileShader(shaderId)
	var status int32
	gl.GetShaderiv(shaderId, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shaderId, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shaderId, logLength, nil, gl.Str(log))
		fmt.Println("Failed to compile shader: \n" + log)
		return 0, errors.New("Failed to compile shader")
	}
	return ShaderID(shaderId), nil
}

func CreateProgram(vertPath string, fragPath string) (ProgramID, error) {
	vert, err := LoadShader(vertPath, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}
	frag, err := LoadShader(fragPath, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}
	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, uint32(vert))
	gl.AttachShader(shaderProgram, uint32(frag))
	gl.LinkProgram(shaderProgram)

	var success int32
	gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &success)
	if success == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(shaderProgram, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(shaderProgram, logLength, nil, gl.Str(log))
		return 0, errors.New("Failed to link program: \n" + log)
	}
	gl.DeleteShader(uint32(vert))
	gl.DeleteShader(uint32(frag))

	//TODO finish hotloading shaders
	/*
		file, err := os.Stat(path)
		if err != nil {
			panic(err)
		}
		modTime := file.ModTime()
		loadedShaders = append(loadedShaders, shaderInfo{path, modTime})
	*/
	return ProgramID(shaderProgram), nil
}

func GenBindBuffer(target uint32) BufferID {
	var buffer uint32
	gl.GenBuffers(1, &buffer)
	gl.BindBuffer(target, buffer)
	return BufferID(buffer)
}

func GenBindVertexArray() BufferID {
	var VAO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.BindVertexArray(VAO)
	return BufferID(VAO)
}

func BindVertexArray(vaoID BufferID) {
	gl.BindVertexArray(uint32(vaoID))
}

func BufferDataFloat(target uint32, data []float32, usage uint32) {
	gl.BufferData(target, len(data)*4, gl.Ptr(data), usage)
}

func BufferDataInt(target uint32, data []uint32, usage uint32) {
	gl.BufferData(target, len(data)*4, gl.Ptr(data), usage)
}

func UnbindVertexArray() {
	gl.BindVertexArray(0)
}

func UseProgram(programID ProgramID) {
	gl.UseProgram(uint32(programID))
}

func TriangleNormal(p1, p2, p3 mgl32.Vec3) mgl32.Vec3 {
	U := p2.Sub(p1)
	V := p3.Sub(p1)
	return U.Cross(V).Normalize()
}
