package goengine

import (
	"fmt"
	"io/ioutil"
	"log"
	"runtime"
	"strings"

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

	if err := glfw.Init(); err != nil {
		log.Fatalln("GLFW failed to initialize:", err)
		fmt.Println("GLFW failed to init")
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Resizable, glfw.False)

	window, windowError := glfw.CreateWindow(800, 600, "OpenGL", nil, nil)

	if windowError != nil {
		fmt.Println("oops")
		panic(windowError)
	}
	fmt.Println("Buffer")
	defer func() {
		window.Destroy()
		glfw.Terminate()
		fmt.Println("ending")
		renderer.ended = true
		system_end <- true
	}()

	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	renderer.keyCallback = func(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		log.Printf("key: %v, scancode: %v, action: %v, mods: %v", key, scancode, action, mods)
	}

	window.SetKeyCallback(renderer.keyCallback)

	<-system_start
	renderer.started = true
	system_started <- true
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
	program.addAttribute("vertexPosition\x00")
	program.linkShaders()

	if err != nil {
		panic(err)
	}

	//bind attributes

L:
	for !window.ShouldClose() {

		select {
		case <-force_end:
			fmt.Println("Should End")
			break L
		default:
			renderer.drawGame(program)
			window.SwapBuffers()
			glfw.PollEvents()
		}
	}

	if vboID != 0 {
		gl.DeleteBuffers(1, &vboID)
	}
}

//GLSLProgram Wrapper for shader programs
type GLSLProgram struct {
	program       uint32
	vertexShader  uint32
	fragShader    uint32
	numAttributes uint32
}

func newProgram(vertexSource string, fragSource string) (GLSLProgram, error) {

	program := GLSLProgram{}
	vertexShader, err := compileShader(vertexSource, gl.VERTEX_SHADER)
	if err != nil {
		return program, err
	}
	fragShader, err := compileShader(fragSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return program, err
	}
	program.vertexShader = vertexShader
	program.fragShader = fragShader

	return program, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	if shader == 0 {
		return shader, fmt.Errorf("failed to create shader of type: %v", shaderType)
	}
	csource := gl.Str(source)
	gl.ShaderSource(shader, 1, &csource, nil)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))

		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		gl.DeleteShader(shader)
		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

func (program *GLSLProgram) linkShaders() {
	program.program = gl.CreateProgram()

	gl.AttachShader(program.program, program.vertexShader)
	gl.AttachShader(program.program, program.fragShader)

	gl.LinkProgram(program.program)

	var isLinked int32
	gl.GetProgramiv(program.program, gl.LINK_STATUS, &isLinked)
	if isLinked == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program.program, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))

		gl.GetProgramInfoLog(program.program, logLength, nil, gl.Str(log))

		gl.DeleteProgram(program.program)
		gl.DeleteShader(program.vertexShader)
		gl.DeleteShader(program.fragShader)
	}
	gl.DetachShader(program.program, program.vertexShader)
	gl.DetachShader(program.program, program.fragShader)
	gl.DeleteShader(program.vertexShader)
	gl.DeleteShader(program.fragShader)
}

func (program *GLSLProgram) addAttribute(attributeName string) {
	gl.BindAttribLocation(program.program, program.numAttributes, gl.Str(attributeName))
	program.numAttributes++
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

func (program *GLSLProgram) use() {
	gl.UseProgram(program.program)
	for i := uint32(0); i < program.numAttributes; i++ {
		gl.EnableVertexAttribArray(i)
	}
}

func (program *GLSLProgram) unuse() {
	gl.UseProgram(0)
	for i := uint32(0); i < program.numAttributes; i++ {
		gl.DisableVertexAttribArray(i)
	}
}

func (renderer *Renderer) drawGame(program GLSLProgram) {
	gl.ClearDepth(1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	program.use()
	entities := renderer.world.ECS.WithComponent("SpriteComponent")
	for _, entity := range entities {
		comp, exist := entity.GetComponent("SpriteComponent")
		if exist {
			sprite, ok := (*comp).(Renderable)
			if ok {
				if !sprite.GraphicsInitialized() {
					sprite.InitGraphics()
				}
				sprite.Draw()
			}
		}

	}

	program.unuse()

}
