package goengine

import (
  "log"
  "fmt"
  "runtime"
  "github.com/go-gl/gl/v4.1-core/gl"
  "github.com/go-gl/glfw/v3.1/glfw"
)

type Renderer struct{
  keyCallback glfw.KeyCallback
  BaseSystem
}

type Renderable interface{
  InitGraphics()
  GraphicsInitialized() bool
  Draw()
}

func (renderer *Renderer) Init(){
    runtime.LockOSThread()
}

func (renderer *Renderer) Run(world *World, system_start <-chan bool, system_started chan<- bool, system_end chan<- bool, force_end <-chan bool ){
  renderer.world = world
  renderer.Init()

  if err:= glfw.Init(); err != nil{
    log.Fatalln("GLFW failed to initialize:", err)
    fmt.Println("GLFW failed to init")
  }

  glfw.WindowHint(glfw.ContextVersionMajor, 3)
  glfw.WindowHint(glfw.ContextVersionMinor, 2)
  glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
  glfw.WindowHint(glfw.Resizable, glfw.False)

  window, window_error := glfw.CreateWindow(800, 600, "OpenGL", nil, nil)

  if window_error != nil{
    fmt.Println("oops")
    panic (window_error)
  }
  fmt.Println("Buffer")
  defer func(){
    window.Destroy()
    glfw.Terminate()
    fmt.Println("ending")
    renderer.ended = true
    system_end <-true
  }()


  window.MakeContextCurrent()

  if err := gl.Init(); err != nil{
      panic(err)
  }

  renderer.keyCallback = func(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey){
    log.Println("key: %v, scancode: %v, action: %v, mods: %v", key, scancode, action, mods)
  }

  window.SetKeyCallback(renderer.keyCallback)

  <-system_start
  renderer.started = true
  system_started <- true
  gl.ClearDepth(1.0)

  var vboID uint32
  gl.GenBuffer(1, &vboID)
  var vertexData = []float32{
    1, 1, 0, 1, 0, 0, 0, 0, 1, 0, 1, 1,
  }
  gl.BindBuffer(gl.ARRAY_BUFFER, vboID)
  gl.BufferData(gl.ARRAY_BUFFER, len(vertexData)*4, gl.Ptr(vertexData), gl.STATIC_DRAW)

  gl.BindBuffer(gl.ARRAY_BUFFER, 0)

  L:
  for !window.ShouldClose(){

    select{
      case <- force_end:
        fmt.Println("Should End")
        break L
      default:
        //renderer.drawGame()
        gl.BindBuffer(gl.ARRAY_BUFFER, vboID)
        gl.EnableVertexAttribArray(0)

        gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))

        gl.DrawArrays(gl.TRIANGLES, 0, 2, 6)

        gl.DisableVertexAttribArray(0)
        window.SwapBuffers()
        glfw.PollEvents()
    }
  }

  if vboID != 0 {
    gl.DeleteBuffers(1, &vboID)
  }
}



func (renderer *Renderer) drawGame(){
  gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

  // entities := renderer.world.ECS.WithComponent("SpriteComponent")
  // for _, entity := range entities{
  //   comp, exist := entity.GetComponent("SpriteComponent")
  //   if exist{
  //     sprite, ok := (*comp).(Renderable)
  //     if ok{
  //       if !sprite.GraphicsInitialized(){
  //         sprite.InitGraphics()
  //       }
  //       sprite.Draw()
  //     }
  //   }

  //}

}
