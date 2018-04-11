package gogl

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"os"
	"time"
)

type Shader struct {
	id               ProgramID
	vertexPath       string
	fragmentPath     string
	vertexModified   time.Time
	fragmentModified time.Time
}

func NewShader(vertexPath string, fragmentPath string) (*Shader, error) {
	id, err := CreateProgram(vertexPath, fragmentPath)
	if err != nil {
		return nil, err
	}
	result := &Shader{id, vertexPath, fragmentPath, getModifiedTime(vertexPath), getModifiedTime(fragmentPath)}
	return result, nil
}

func (shader *Shader) Use() {
	UseProgram(shader.id)
}

func (shader *Shader) SetFloat(name string, f float32) {
	name_cstr := gl.Str(name + "\x00")
	location := gl.GetUniformLocation(uint32(shader.id), name_cstr)
	gl.Uniform1f(location, f)
}

func getModifiedTime(filePath string) time.Time {
	file, err := os.Stat(filePath)
	if err != nil {
		panic(err)
	}
	return file.ModTime()
}

func (shader *Shader) CheckShaderForChanges() {

	vertexModTime := getModifiedTime(shader.vertexPath)
	fragmentModTime := getModifiedTime(shader.fragmentPath)
	// check if greater than?
	if !vertexModTime.Equal(shader.vertexModified) ||
		!fragmentModTime.Equal(shader.fragmentModified) {
		id, err := CreateProgram(shader.vertexPath, shader.fragmentPath)
		if err != nil {
			fmt.Println(err)
		} else {
			gl.DeleteProgram(uint32(shader.id))
			shader.id = id
		}
	}

}
