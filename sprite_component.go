package goengine

import(
    "fmt"
    "github.com/go-gl/gl/v4.1-core/gl"
)

type SpriteComponent struct{
  x float32
  y float32
  width float32
  height float32
  vboID uint32
  vertexData []float32
  graphics_initialized bool
  BaseComponent
}

func (sC *SpriteComponent) GetID() string{
  if sC.id == "" {
    sC.id = "SpriteComponent"
  }
  return sC.id
}

func (sC *SpriteComponent) SetDimensions(){

}

func (sC *SpriteComponent) GraphicsInitialized() bool{
  return sC.graphics_initialized
}

func (sC *SpriteComponent) InitGraphics(){
  if (sC.vboID == 0){
    gl.GenBuffers(1, &sC.vboID)
    fmt.Println(sC.vboID)
  }

  //first triangel
  sC.vertexData = [12]float32
  sC.vertexData[0] = sC.x + sC.width
  sC.vertexData[1] = sC.y + sC.height

  sC.vertexData[2] = sC.x
  sC.vertexData[3] = sC.y + sC.height

  sC.vertexData[4] = sC.x
  sC.vertexData[5] = sC.y
  //second triangle
  sC.vertexData[6] = sC.x
  sC.vertexData[7] = sC.y

  sC.vertexData[8] = sC.x + sC.width
  sC.vertexData[9] = sC.y

  sC.vertexData[10] = sC.x + sC.width
  sC.vertexData[11] = sC.y + sC.height

  gl.BindBuffer(gl.ARRAY_BUFFER, sC.vboID);
  gl.BufferData(gl.ARRAY_BUFFER, len(sC.vertexData)*4, gl.Ptr(sC.vertexData), gl.STATIC_DRAW)
  //unbind buffer
  gl.BindBuffer(gl.ARRAY_BUFFER, 0)

  sC.graphics_initialized = true
}

func (sC *SpriteComponent) Draw(){
  //fmt.Println("Drawing")
  gl.BindBuffer(gl.ARRAY_BUFFER, sC.vboID)
  //first index off array
  gl.EnableVertexAttribArray(0)
  //Pointing opengl to the start of our data
  //gl.VertexAttribPointer(index, size of 1 vertex (num elements), type, Normalize?, stride, pointer for interlieved)
  gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 0, gl.Ptr(nil))
  //gl.DrawArrays(mode (quads, triangles, points), first element, how many to draw)
  gl.DrawArrays(gl.TRIANGLES, 0, 6)
  gl.DisableVertexAttribArray(0)

  gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

func (sC *SpriteComponent) Destroy(){
  if (sC.vboID != 0){
    gl.DeleteBuffers(1, &sC.vboID)
  }
}

func (sC *SpriteComponent) Update(delta_time float64){
  if(sC.entity.CheckDelete()){
    sC.Destroy()
  }
}
