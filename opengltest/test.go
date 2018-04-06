package main

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/veandco/go-sdl2/sdl"
	"strings"
	"unsafe"
)

const winWidth = 1280
const winHeight = 720

func main() {

	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("OpenGL", 200, 200,
		int32(winWidth), int32(winHeight), sdl.WINDOW_OPENGL)
	if err != nil {
		fmt.Println(err)
		return
	}
	window.GLCreateContext()
	defer window.Destroy()
	fmt.Println("A")
	gl.Init()
	fmt.Println("B")
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	vertices := []float32{
		0.5, 0.5, 0.0,
		0.5, -0.5, 0.0,
		-0.5, -0.5, 0.0,
		-0.5, 0.5, 0.0}

	var VBO uint32
	gl.GenBuffers(1, &VBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*int(unsafe.Sizeof(vertices[0])), gl.Ptr(vertices), gl.STATIC_DRAW)

	vertexShaderStr :=
		`#version 330 core
		layout (location = 0) in vec3 aPos;
	
		void main()
		{
			gl_Position = vec4(aPos.x, aPos.y, aPos.z, 1.0);
		}`
	vertexShader := gl.CreateShader(gl.VERTEX_SHADER)
	csources, free := gl.Strs(vertexShaderStr)
	gl.ShaderSource(vertexShader, 1, csources, nil)
	free()
	gl.CompileShader(vertexShader)
	var status int32
	gl.GetShaderiv(vertexShader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(vertexShader, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(vertexShader, logLength, nil, gl.Str(log))
		panic("failed to compile vertex shader \n" + log)
	}

	fragmentShaderStr := `
		#version 330 core
		out vec4 FragColor;

		void main()
		{
    		FragColor = vec4(1.0f, 0.5f, 0.2f, 1.0f);
		} `

	fragmentShader := gl.CreateShader(gl.VERTEX_SHADER)
	csources, free = gl.Strs(fragmentShaderStr)
	gl.ShaderSource(fragmentShader, 1, csources, nil)
	free()
	gl.CompileShader(fragmentShader)

	gl.GetShaderiv(fragmentShader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(fragmentShader, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(fragmentShader, logLength, nil, gl.Str(log))
		panic("failed to compile vertex shader \n" + log)
	}

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}
	}

}
