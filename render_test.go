package gogogamengine

import (
  "fmt"
  "os"
  "runtime"
  "strings"
  "io/ioutil"
  "path/filepath"
  "testing"

  "github.com/golang/glog"
  "github.com/go-gl/gl/v4.1-core/gl"
  "github.com/go-gl/glfw/v3.1/glfw"
)

const WindowHeight = 600
const WindowWidth = 800

func init(){
  runtime.LockOSThread()
}

type GLSLProgram struct{
  program uint32
  vertexShaderSource string
  fragmentShaderSource string
  linked bool
}
func(p *GLSLProgram) Use(){
  gl.UseProgram(p.program)
}
//TODO Implement changing shaders and relinking

func Run(){
  if err := glfw.Init(); err != nil{
    glog.Fatalln("failed to initialize glfw", err)
  }
  defer glfw.Terminate()

  setupWindowOptions()
  window, err := glfw.CreateWindow(WindowWidth, WindowHeight, "Game", nil, nil)
  if err != nil{
    panic(err)
  }
  window.MakeContextCurrent()

  //initilize Glow
  if err := gl.Init(); err != nil{
    panic(err)
  }

  version := gl.GoStr(gl.GetString(gl.VERSION))
  fmt.Println("OpenGL version", version)

  shaderSource, err := ReadShaders("colorShader")
  if err != nil{
    panic(err)
  }
  program, err := NewProgram(shaderSource)
  if err != nil{
    panic(err)
  }
  program.Use()

  sprite := &SpriteComponent{-.5,-.5,1,1}
  sprite.ReloadGraphics()


  vertAttrib := uint32(gl.GetAttribLocation(program.program, CStr("vertPosition")))
  gl.EnableVertexAttribArray(vertAttrib)
  gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

  gl.Enable(gl.DEPTH_TEST)
  gl.DepthFunc(gl.LESS)
  gl.ClearColor(1.0, 1.0, 1.0, 1.0)

  for !window.ShouldClose(){
    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

    program.Use()

    gl.BindVertexArray(sprite.vaoID)
    gl.DrawArrays(gl.TRIANGLES, 0, 2*3)

    window.SwapBuffers()
    glfw.PollEvents()
  }
}

func setupWindowOptions(){
  glfw.WindowHint(glfw.Resizable, glfw.False)
  glfw.WindowHint(glfw.ContextVersionMajor, 4)
  glfw.WindowHint(glfw.ContextVersionMinor, 1)
  glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
  glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
}

func NewProgram(shaderSource map[string]string)(*GLSLProgram, error){
  var hasVertexShader bool
  var hasFragmentShader bool

  vertexShaderSource, ok := shaderSource["vertex"]
  var vertexShader uint32
  if ok{

    vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
    if err != nil{
      return nil, err
    }
    hasVertexShader = true
  }

  fragmentShaderSource, ok := shaderSource["fragment"]
  var fragmentShader uint32
  if ok{

    fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
    if err != nil{
      return nil, err
    }
    hasFragmentShader = true
  }

  program := gl.CreateProgram()

  if !(hasVertexShader && hasFragmentShader){
    return nil, fmt.Errorf("Need vertex shader: %v, and fragment shader: %v", hasVertexShader, hasFragmentShader)
  }

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

    return nil, fmt.Errorf("Failed to link program(%v): %v",program, log)
  }

  gl.DeleteShader(vertexShader)
  gl.DeleteShader(fragmentShader)
  p := &GLSLProgram{program, vertexShaderSource, fragmentShaderSource, true}
  return p, nil
}

func compileShader(source string, shaderType uint32) (uint32, error){
  shader := gl.CreateShader(shaderType)

  csource := CStr(source)
  gl.ShaderSource(shader, 1, &csource, nil)
  gl.CompileShader(shader)

  var status int32
  gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
  if status == gl.FALSE{
    var logLength int32
    gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

    log := strings.Repeat("\x00", int(logLength+1))
    gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

    return 0, fmt.Errorf("Failed to compile %v: %v", source, log)
  }

  return shader, nil
}

func CStr(str string) *uint8{
  if !strings.HasSuffix(str, "\x00"){
    str = str+ "\x00"
  }
  return gl.Str(str)
}

//Looks for all shaders of with that name(vertex:.vert, fragment:.frag ) and returns a map
//with the stader sources. Only looks in shaders directory

var validShaderExt = map[string]string{".vert":"vertex", ".frag":"fragment"}
func ReadShaders(shadername string) (map[string]string, error){
  shaders := make(map[string]string)
  dirname := ".."+string(filepath.Separator)+"shaders"

  d, err := os.Open(dirname)
  if err != nil{
    return nil, err
  }
  defer d.Close()

  files, err := d.Readdir(-1)
  if err != nil{
    return nil, err
  }

  for _, file := range files{
    if file.Mode().IsRegular(){
      checkExt := filepath.Ext(file.Name())
      if shaderType, ok := validShaderExt[checkExt]; ok{
        if source, err := ioutil.ReadFile(file.Name()); err != nil{
          shaders[shaderType] = string(source)
        }
      }
    }
  }

  return shaders, nil
}

func TestRenderer(t *testing.T){
  Run()
}
