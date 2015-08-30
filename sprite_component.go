package gogogamengine

import "github.com/go-gl/gl/v4.1-core/gl"

type SpriteComponent struct {
	x                    float32
	y                    float32
	width                float32
	height               float32
	vaoID								 uint32
	vboID                uint32
	vertexData           [12]float32
	GraphicsInitialized  bool
	BaseComponent
}

func (sC *SpriteComponent) GetID() string {
	if sC.id == "" {
		sC.id = "SpriteComponent"
	}
	return sC.id
}
func (sC *SpriteComponent) SetPosition(x, y float32){
	sC.x = x
	sC.y = y
	if sC.GraphicsInitialized == true{
		sC.ReloadGraphics()
	}
}
//Sets the dimensions of Sprite
func (sC *SpriteComponent) SetDimensions(width, height float32) {
	sC.width = width
	sC.height = height
	if sC.GraphicsInitialized == true{
		sC.ReloadGraphics()
	}
}

//After Changing Position or Dimensions ReloadGraphics
func (sC *SpriteComponent) ReloadGraphics(){
	if sC.vaoID == 0{
		gl.GenVertexArrays(1, &sC.vaoID)
	}
	gl.BindVertexArray(sC.vaoID)
	if sC.vboID == 0{
		gl.GenBuffers(1, &sC.vboID)
	}

	//first triangel

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

	gl.BindBuffer(gl.ARRAY_BUFFER, sC.vboID)
	gl.BufferData(gl.ARRAY_BUFFER, len(sC.vertexData)*4, gl.Ptr(sC.vertexData), gl.STATIC_DRAW)
	//unbind buffer
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	sC.GraphicsInitialized = true
}

func (sC *SpriteComponent) Destroy() {
	if sC.vboID != 0 {
		gl.DeleteBuffers(1, &sC.vboID)
	}
}

func (sC *SpriteComponent) Update(delta_time float64) {
	if sC.entity.CheckDelete() {
		sC.Destroy()
	}
}
