package goengine

import (
	"io/ioutil"
	"runtime"
	"strings"
	"fmt"
	"errors"
  "github.com/golang/glog"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

//Render Game system that controls rendering and window operations
type Renderer struct {
	keyCallback glfw.KeyCallback
	BaseSystem
}

//Renderable any type that can be used with Render System
type Renderable interface {
	InitGraphics()
	GraphicsInitialized() bool
	Draw()
}

func (renderer *Renderer) Init() {
	runtime.LockOSThread()
}

func (renderer *Renderer) Run(world *World, system_start <-chan bool, system_started chan<- bool, system_end chan<- bool, force_end <-chan bool) {
	renderer.world = world
	renderer.Init()
	renderer.Loop()
}

func (renderer *Renderer) Loop(){

		defer glog.Flush()

		if err := glfw.Init(); err != nil {
			glog.Fatalln("GLFW failed to initialize:", err)
			glog.Info("GLFW failed to init")
		}

		glfw.WindowHint(glfw.ContextVersionMajor, 3)
		glfw.WindowHint(glfw.ContextVersionMinor, 2)
		glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
		glfw.WindowHint(glfw.Resizable, glfw.False)

		window, windowError := glfw.CreateWindow(800, 600, "OpenGL", nil, nil)

		if windowError != nil {
			glog.Info("oops")
			panic(windowError)
		}
		glog.Info("Buffer")
		defer func() {
			window.Destroy()
			glfw.Terminate()
			glog.Info("ending")
			renderer.ended = true
			//system_end <- true
		}()

		window.MakeContextCurrent()

		if err := gl.Init(); err != nil {
			panic(err)
		}

		renderer.keyCallback = func(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
			glog.Infof("key: %v, scancode: %v, action: %v, mods: %v", key, scancode, action, mods)
		}

		window.SetKeyCallback(renderer.keyCallback)

		//<-system_start
		renderer.started = true
		//system_started <- true
		gl.ClearDepth(1.0)

		var vboID uint32
		gl.GenBuffers(1, &vboID)
		var vertexData = []float32{
			1, 1, 0, 1, 0, 0, 0, 0, 1, 0, 1, 1,
		}
		gl.BindBuffer(gl.ARRAY_BUFFER, vboID)
		gl.BufferData(gl.ARRAY_BUFFER, len(vertexData)*4, gl.Ptr(vertexData), gl.STATIC_DRAW)
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)

		//get shaders
		vertexShader, fragmentShader := getShaders()
		program, err := newProgram(vertexShader, fragmentShader)
		program.UseProgram()

		if err != nil {
			panic(err)
		}

		//bind attributes

//	L:
		for !window.ShouldClose() {

				//renderer.drawGame(program)
				gl.EnableVertexAttribArray(0)
				gl.BindBuffer(gl.ARRAY_BUFFER, vboID)
				gl.VertexAttribPointer(
					0,
					3,
					gl.FLOAT,
					false,
					0,
					gl.PtrOffset(0))
				gl.DrawArrays(gl.TRIANGLES, 0, 3)
				gl.DisableVertexAttribArray(0);
				window.SwapBuffers()
				glfw.PollEvents()
			}


		if vboID != 0 {
			gl.DeleteBuffers(1, &vboID)
		}

}

//GLSLProgram Wrapper for shader programs
type GLSLProgram struct {
	program       uint32 //program
	numAttributes uint32
}

func newProgram(vertexSource string, fragSource string) (GLSLProgram, error) {

	gP := GLSLProgram{}
	vertexShader, err := compileShader(vertexSource, gl.VERTEX_SHADER)
	if err != nil {
		return gP, err
	}
	fragmentShader, err := compileShader(fragSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return gP, err
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return gP, errors.New(fmt.Sprintf("failed to link program: %v", log))
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	gP.program = program

	return gP, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csource := gl.Str(source)
	gl.ShaderSource(shader, 1, &csource, nil)
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE{
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}
	return shader, nil
}

func getShaders() (string, string) {
	colorVert, err := ioutil.ReadFile("shaders/colorShader.vert")
	if err != nil {
		//panic(err)
	}
	colorFrag, err := ioutil.ReadFile("shaders/colorShader.frag")
	if err != nil {
		//panic(err)
	}
	return string(colorVert) + "\x00", string(colorFrag) + "\x00"
}

func (program *GLSLProgram) UseProgram() {
	gl.UseProgram(program.program)
	for i := uint32(0); i < program.numAttributes; i++ {
		gl.EnableVertexAttribArray(i)
	}
}

func (program *GLSLProgram) UnUse() {
	gl.UseProgram(0)
	for i := uint32(0); i < program.numAttributes; i++ {
		gl.DisableVertexAttribArray(i)
	}
}

func (renderer *Renderer) drawGame(program GLSLProgram) {

}
