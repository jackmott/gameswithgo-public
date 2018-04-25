package main

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/jackmott/gogl"
	"github.com/veandco/go-sdl2/sdl"
	"time"
)

const winWidth = 720
const winHeight = 720

func main() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}
	defer sdl.Quit()

	sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 3)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 3)

	window, err := sdl.CreateWindow("Hello Triangle", 200, 200, winWidth, winHeight, sdl.WINDOW_OPENGL)
	sdl.SetRelativeMouseMode(true)
	if err != nil {
		panic(err)
	}
	window.GLCreateContext()
	defer window.Destroy()

	gl.Init()
	gl.Enable(gl.DEPTH_TEST)
	gogl.MglTest()

	fmt.Println("OpenGL Version", gogl.GetVersion())

	shaderProgram, err := gogl.NewShader("shaders/hello.vert", "shaders/quadtexture.frag")
	if err != nil {
		panic(err)
	}

	texture := gogl.LoadTextureAlpha("assets/metalbox_full.png")

	vertices := []float32{
		-0.5, -0.5, -0.5, 0.0, 0.0,
		0.5, -0.5, -0.5, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		-0.5, 0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 0.0,

		-0.5, -0.5, 0.5, 0.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 1.0,
		0.5, 0.5, 0.5, 1.0, 1.0,
		-0.5, 0.5, 0.5, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,

		-0.5, 0.5, 0.5, 1.0, 0.0,
		-0.5, 0.5, -0.5, 1.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,
		-0.5, 0.5, 0.5, 1.0, 0.0,

		0.5, 0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, 0.5, 0.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0,

		-0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, -0.5, 1.0, 1.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,

		-0.5, 0.5, -0.5, 0.0, 1.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, 0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0,
		-0.5, 0.5, 0.5, 0.0, 0.0,
		-0.5, 0.5, -0.5, 0.0, 1.0}

	normals := make([]float32, 36*3)
	for tri := 0; tri < 12; tri++ {
		index := tri * 15
		p1 := mgl32.Vec3{vertices[index], vertices[index+1], vertices[index+2]}
		index += 5
		p2 := mgl32.Vec3{vertices[index], vertices[index+1], vertices[index+2]}
		index += 5
		p3 := mgl32.Vec3{vertices[index], vertices[index+1], vertices[index+2]}
		normal := gogl.TriangleNormal(p1, p2, p3)
		normals[tri*9] = normal.X()
		normals[tri*9+1] = normal.Y()
		normals[tri*9+2] = normal.Z()

		normals[tri*9+3] = normal.X()
		normals[tri*9+4] = normal.Y()
		normals[tri*9+5] = normal.Z()

		normals[tri*9+6] = normal.X()
		normals[tri*9+7] = normal.Y()
		normals[tri*9+8] = normal.Z()

	}

	cubePositions := []mgl32.Vec3{
		mgl32.Vec3{0.0, 0.0, 0.0},
		mgl32.Vec3{2.0, 5.0, -10.0},
		mgl32.Vec3{1.0, -5.0, 1.0}}

	VAO := gogl.GenBindVertexArray()
	gogl.GenBindBuffer(gl.ARRAY_BUFFER)
	gogl.BufferDataFloat(gl.ARRAY_BUFFER, vertices, gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, nil)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)
	gogl.GenBindBuffer(gl.ARRAY_BUFFER)
	gogl.BufferDataFloat(gl.ARRAY_BUFFER, normals, gl.STATIC_DRAW)
	gl.VertexAttribPointer(2, 3, gl.FLOAT, false, 3*4, nil)
	gl.EnableVertexAttribArray(2)
	gogl.UnbindVertexArray()

	keyboardState := sdl.GetKeyboardState()

	position := mgl32.Vec3{0.0, 0.0, 0.0}
	worldUp := mgl32.Vec3{0.0, 1.0, 0.0}

	camera := gogl.NewCamera(position, worldUp, -90.0, 0.0, 0.005, 0.1)
	var elapsedTime float32
	var mouseX, mouseY int32
	for {
		frameStart := time.Now()
		mouseX = 0
		mouseY = 0
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				return
			case *sdl.MouseMotionEvent:
				mouseX += e.XRel
				mouseY += e.YRel
			}
		}

		dir := gogl.Nowhere
		if keyboardState[sdl.SCANCODE_LEFT] != 0 {
			dir = gogl.Left
		}
		if keyboardState[sdl.SCANCODE_RIGHT] != 0 {
			dir = gogl.Right
		}
		if keyboardState[sdl.SCANCODE_UP] != 0 {
			dir = gogl.Forward
		}
		if keyboardState[sdl.SCANCODE_DOWN] != 0 {
			dir = gogl.Backward
		}

		camera.UpdateCamera(dir, elapsedTime, float32(mouseX), float32(mouseY))

		gl.ClearColor(0.0, 0.0, 0.0, 0.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		shaderProgram.Use()
		projectionMatrix := mgl32.Perspective(mgl32.DegToRad(45.0), float32(winWidth)/float32(winHeight), 0.1, 100.0)
		viewMatrix := camera.GetViewMatrix()
		shaderProgram.SetMat4("projection", projectionMatrix)
		shaderProgram.SetMat4("view", viewMatrix)
		shaderProgram.SetVec3("lightPos", mgl32.Vec3{0.0, 0.0, 1.0})
		shaderProgram.SetVec3("lightColor", mgl32.Vec3{1.0, 1.0, 1.0})
		shaderProgram.SetVec3("viewPos", camera.Position)
		gogl.BindTexture(texture)
		gogl.BindVertexArray(VAO)
		for i, pos := range cubePositions {
			modelMatrix := mgl32.Ident4()
			angle := 20.0 * float32(i)
			//todo normalize vec3?
			modelMatrix = mgl32.HomogRotate3D(mgl32.DegToRad(angle), mgl32.Vec3{1.0, 0.3, 0.5}).Mul4(modelMatrix)
			modelMatrix = mgl32.Translate3D(pos.X(), pos.Y(), pos.Z()).Mul4(modelMatrix)
			shaderProgram.SetMat4("model", modelMatrix)
			gl.DrawArrays(gl.TRIANGLES, 0, 36)
		}

		window.GLSwap()
		shaderProgram.CheckShaderForChanges()
		elapsedTime = float32(time.Since(frameStart).Seconds() * 1000)
	}
}
