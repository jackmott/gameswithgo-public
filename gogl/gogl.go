package gogl

import (
	"errors"
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"io/ioutil"
	_ "os"
	"strings"
	"time"
)

type ShaderID uint32
type ProgramID uint32
type VAOID uint32
type VBOID uint32

func GetVersion() string {
	return gl.GoStr(gl.GetString(gl.VERSION))
}

type programInfo struct {
	vertPath string
	fragPath string
	modified time.Time
}

var loadedShaders []programInfo

func CheckShadersForChanges() {
	/*
		for _, shaderInfo := range loadedShaders {
			file, err := os.Stat(oriInfo.path)
			if err != nil {
				panic(err)
			}
			modTime := file.ModTime()
			// check if greater than?
			if !modTime.Equal(shaderInfo.modified) {
				fmt.Println("Shader modified")
			}
		}
	*/
}

func LoadShader(path string, shaderType uint32) (ShaderID, error) {
	shaderFile, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
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

func GenBindBuffer(target uint32) VBOID {
	var VBO uint32
	gl.GenBuffers(1, &VBO)
	gl.BindBuffer(target, VBO)
	return VBOID(VBO)
}

func GenBindVertexArray() VAOID {
	var VAO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.BindVertexArray(VAO)
	return VAOID(VAO)
}

func BindVertexArray(vaoID VAOID) {
	gl.BindVertexArray(uint32(vaoID))
}

func BufferDataFloat(target uint32, data []float32, usage uint32) {
	gl.BufferData(target, len(data)*4, gl.Ptr(data), usage)
}

func UnbindVertexArray() {
	gl.BindVertexArray(0)
}

func UseProgram(programID ProgramID) {
	gl.UseProgram(uint32(programID))
}
