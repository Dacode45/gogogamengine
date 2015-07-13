package goengine

import "testing"
import "time"
import "fmt"

func TestCreateEntity(t *testing.T){
  e := NewEntity()
  fmt.Println(e)
}

func TestAddingEntity(t *testing.T){
  world := NewWorld()
  world.AddEntity(NewEntity())
}

//Test component
type Speak struct{
  BaseComponent
  say string
  count int;
}

func (s *Speak) Update(delta_time float64){
  if(s.count < 5){
    fmt.Println(s.say)
    s.count++
  }
}

func TestAddComponent(t *testing.T){
  fmt.Println("Starting TestAddComponent\n")
  world := NewWorld()
  e := NewEntity()
  world.AddEntity(e)
  speak := Speak{say:"Hello World"}
  e.AddComponent(&speak)
  //Fixing this requires reflection
}

func TestWorld(t *testing.T){
  fmt.Println("Starting TestWorld\n")
  world := NewWorld()
  e := NewEntity()
  speak := Speak{say:"Do It Work? It do!"}
  e.AddComponent(&speak)
  world.AddEntity(e)
  go world.StartWorld()
  time.Sleep(time.Second*2)
  fmt.Println("CLOSINGCLOSINGCLOSINGCLOSINGCLOSINGCLOSINGCLOSINGCLOSING")
  world.Close()
  fmt.Println("CLOSINGCLOSINGCLOSINGCLOSINGCLOSINGCLOSINGCLOSINGCLOSING")

}
