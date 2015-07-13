package goengine

import "sync"
import "fmt"

type World struct{
  //Bunch of local variables
  ECS *entityCollection
  CCS *ComponentCollection
  root *Entity
  entity_update_limit uint

  ECSLock *sync.RWMutex

  gS GameSystem
  systems map[string]System

  World_Alive chan bool
}

func NewWorld() *World{
  world := World{
    ECS:NewEntityCollection(),
    CCS:NewComponentCollection(),
    root:NewEntity(),
    ECSLock: &sync.RWMutex{},
    entity_update_limit:100,
    World_Alive:make(chan bool)}
  return &world
}

func (world *World) AddEntity(e *Entity){
  fmt.Println("adding entity")
  world.ECSLock.Lock()
  err := world.ECS.AddEntity(e)
  if err != nil{
    fmt.Println(e)
  }
  world.ECSLock.Unlock()
}

func (world *World) StartWorld(){
  //start the game system
  gS_start := make(chan bool)
  gS_started := make(chan bool)
  gS_ended := make(chan bool)
  world.gS = GameSystem{}
  fmt.Println(world)
  go world.gS.Run(world, gS_start, gS_started, gS_ended)
  fmt.Println("hello")
  gS_start <-true
  <-gS_started
  <-gS_ended

  //Start all the systems

}

func (world *World) GetEntityUpdateLimit() uint{
  return world.entity_update_limit
}

func (world *World) Close(){
  world.World_Alive <- false
}
//
// func (world *World) String() string{
//   return fmt.Sprintf("Hello World")
// }
