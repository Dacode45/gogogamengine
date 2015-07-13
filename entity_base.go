package goengine

import (
  "fmt"
  "math/rand"
  "errors"
  "strconv"
  "sync"
  "time"

)

type entityCollection struct{
  Entities map[string]*Entity
  num_updates  uint64
  should_update bool
}


func NewEntityCollection() *entityCollection{
  var ECS = entityCollection{}
  ECS.Entities = make(map[string]*Entity)
  println("Entity Collection ", &ECS)
  return &ECS
}

func (c *entityCollection) AddEntity(entity *Entity) error {
  //Do stuff for the Component
  if _, ok := c.Entities[entity.id]; ok{
    return errors.New("Entity with this id Already Exist!")
  }
  c.Entities[entity.id] = entity
  return nil
}


//Just keep tacking on Components to one entity class
type Entity struct{
  id string
  components  map[string]Component
  world *World
  parent *Entity

  enabled bool
  can_enable bool
  should_delete bool
  last_update time.Time
  delta_time float64

  //Check everyting awakes and starts atleast once
  try_to_start bool
  try_to_awaken bool
}

var entityMutex = &sync.Mutex{}

func NewEntity() *Entity{
  return &Entity{id : strconv.Itoa( rand.Intn(1000000)), components:make(map[string]Component)}
}

func (e *Entity) GetId() string{
  return e.id
}

func (e *Entity) AddComponent(component Component) *Component{


  if c, ok := e.components[component.GetId()]; ok{
    //Component already attached
    return &c
  }

  e.components[component.GetId()] = component
  component.SetEntity(e)
  component.SetEnabled(true)
  component.Awake()
  return &component

}

func (e *Entity) AddComponentByString(id string, CCS ComponentCollection) *Component{
  component, ok := CCS.GetComponent(id)
  if !ok {
    return nil
  }

  return e.AddComponent(component)
}

//Returns Component to enable transfering of a Component
func (e *Entity) RemoveComponent(id string) *Component{
  component, ok := e.components[id]
  if !ok {
    return nil
  }
  componentptr := &component
  delete(e.components, id)
  return componentptr
}

func (e *Entity) String() string{
  return fmt.Sprintf("%b", e.id)
}

func (e *Entity) SetEnabled(enabled bool){
  entityMutex.Lock()
  if e.can_enable {
    e.enabled = enabled
  }
  entityMutex.Unlock()
}

func (e *Entity) Awake(){
  e.can_enable = true
  e.SetEnabled(true)
  e.try_to_awaken = true
}

func (e *Entity) Start(bulk_start *sync.WaitGroup, finished_start chan<- string){
  fmt.Println(bulk_start)
  e.try_to_start = true
  e.last_update = time.Now()
  no_crash := e.id
  defer func(){
    finished_start <- no_crash
    bulk_start.Done()
  }()


  for _, component := range e.components{
    if component.IsEnabled(){

      component.Start()
    }
  }

  no_crash = "" //if we get to this point no Component crashed and can nulify this
}


//If panic happens finished update will return the id of this object
func (e *Entity) Update(bulk_update *sync.WaitGroup, can_update <-chan bool, finished_update chan<- string){
      cU := <-can_update
      e.delta_time = time.Now().Sub(e.last_update).Seconds()
      e.last_update = time.Now()

      no_crash := e.id
      defer func(){
        //fmt.Println("Finished the update")
        finished_update<-no_crash
        bulk_update.Done()
      }()
      if cU && e.enabled{
        for _, component := range e.components{
          if component.IsEnabled(){
          //  fmt.Println("Updating component")
            component.Update(e.delta_time)
          }
        }
      }
      no_crash = ""
}

func (e *Entity) Delete(){
  if !e.should_delete{

    e.should_delete = true
  }
}

func (e *Entity) CheckDelete() bool{
  return e.should_delete
}
