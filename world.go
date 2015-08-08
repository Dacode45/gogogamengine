package goengine

import "sync"
import "fmt"

//Contains channels for comunication between systems
type SystemContainer struct{
  system System
  system_start chan bool
  system_started chan bool
  system_end chan bool
  force_end chan bool
  chan_open bool
}

func (sC *SystemContainer) End(){
  //if you intended for it to run
  if sC.chan_open{
    //if it has not ended already
    if !sC.system.Ended(){
      //if the system has started and is in main loop
      if sC.system.Started(){
        sC.force_end <- true
        <- sC.system_end
      }else{ //You're ending the system before it begins
        sC.force_end <- true
        <- sC.system_end
      }
      close(sC.system_start)
      close(sC.system_end)
      close(sC.force_end)
    }

  }
  fmt.Printf("%v : started = %v, ended = %v, ran = %v\n", sC.system.GetID(), sC.system.Started(), sC.system.Ended(), sC.chan_open)
}

func (sC *SystemContainer) Run(world *World){

  sC.system_start = make(chan bool)
  sC.system_started = make(chan bool)
  sC.system_end = make(chan bool)
  sC.force_end = make(chan bool)
  sC.chan_open = true //Channels were opened
  go sC.system.Run(world, sC.system_start, sC.system_started, sC.system_end, sC.force_end)
  sC.system_start <- true
  fmt.Println("Sending Start")
  <- sC.system_started
  fmt.Println("Recieved Start")
}

type World struct{
  //Bunch of local variables
  ECS *EntityCollection
  CCS *ComponentCollection
  root *Entity
  entity_update_limit uint

  ECSLock *sync.RWMutex

  systems map[string]*SystemContainer

  alive bool
}

const(
  GameSystemID = "GameSystem"
  RenderSystemID = "RenderSystem"
)


func NewWorld() *World{
  gSC := &SystemContainer{
    system: &GameSystem{},
  }
  rC := &SystemContainer{
    system: &Renderer{},
  }

  gSC.system.SetID(GameSystemID)
  rC.system.SetID(RenderSystemID)
  world := World{
    ECS:NewEntityCollection(),
    CCS:NewComponentCollection(),
    root:NewEntity(),
    ECSLock: &sync.RWMutex{},
    entity_update_limit:100,
    systems: make(map[string]*SystemContainer)}
  world.systems[GameSystemID] = gSC
  world.systems[RenderSystemID] = rC
  return &world
}

func (world *World) AddSystem(s *SystemContainer){
  dup, has_dup := world.systems[s.system.GetID()]
  if has_dup {
    dup.End()
  }
  world.systems[s.system.GetID()] = s
}

func (world *World) AddEntity(e *Entity){
  //fmt.Println("adding entity")
  err := world.ECS.AddEntity(e)
  if err != nil{
  //  fmt.Println(e)
  }
}

func (world *World) StartWorld(){
  world.alive = true

  //Run systems in a certain order first
  //Render system
  go world.systems[RenderSystemID].Run(world)
  go world.systems[GameSystemID].Run(world)
  essential_systems := []string{RenderSystemID, GameSystemID}
  for systemID, system := range world.systems{
    found_essential := false
    for _, essential := range essential_systems{
      //fmt.Println(systemID, essential)
      if systemID == essential{
        //We know it was activated.
        found_essential = true
      }
    }
    if !found_essential{
      fmt.Printf("Running %v\n", system.system.GetID())
      go system.Run(world)
    }
  }

  //Start all the systems

}

func (world *World) GetEntityUpdateLimit() uint{
  return world.entity_update_limit
}

func (world *World) Close(){
  world.alive = false
  for _, sC := range world.systems {
    sC.End()
  }
}

func (world *World) IsAlive() bool{
  return world.alive
}


func (world *World) String() string{
  return fmt.Sprintf("A World %v", world.alive)
}
//
// func (world *World) String() string{
//   return fmt.Sprintf("Hello World")
// }
