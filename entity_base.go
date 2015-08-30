package gogogamengine

import (
  "fmt"
  "math/rand"
  "errors"
  "strconv"
  "sync"
  "time"

)

type EntityCollection struct{
  Entities map[string]*Entity
  num_updates  uint64
  should_update bool
  entity_component_cache map[string]map[string]bool
}


func NewEntityCollection() *EntityCollection{
  var ECS = EntityCollection{}
  ECS.Entities = make(map[string]*Entity)
  ECS.entity_component_cache = make(map[string]map[string]bool)
  return &ECS
}

func (c *EntityCollection) AddToComponentCache(cID string, e *Entity){
  if _, ok := c.entity_component_cache[cID]; !ok{
    m := make(map[string]bool)
    c.entity_component_cache[cID] = m
  }
  c.entity_component_cache[cID][e.GetID()] = true
}

func (c *EntityCollection) RemoveFromComponentCache(cID string, e *Entity){
  if m, ok := c.entity_component_cache[cID]; ok{
    delete(m, e.GetID())
  }
}

func (c *EntityCollection) AddEntity(entity *Entity) error {
  //Do stuff for the Component
  if _, ok := c.Entities[entity.id]; ok{
    return errors.New("Entity with this id Already Exist!")
  }
  c.Entities[entity.id] = entity
  entity.ECS = c
  //immediately add componetns to component cache
  for compID, _ := range entity.components{
    c.AddToComponentCache(compID, entity)
  }
  return nil
}

func (c *EntityCollection) WithComponent(cID string) []*Entity {

//  fmt.Println("toReturn")
  toReturn := make([]*Entity, 0)
  for entityID, _ := range c.entity_component_cache[cID]{
    if entity, ok := c.Entities[entityID]; ok{
      if entity.should_delete == false{

        toReturn = append(toReturn, entity)
      }
    }
  }
  return toReturn

}


//Just keep tacking on Components to one entity class
type Entity struct{
  //Reference to EntityCollection
  ECS *EntityCollection

  id string
  components  map[string]Component
  world *World
  parent *Entity

  enabled bool
  should_delete bool
  last_update time.Time
  delta_time float64

  //Check everyting awakes and starts atleast once
  try_to_start bool
  try_to_awaken bool

  entityMutex *sync.RWMutex

}

func NewEntity() *Entity{
  return &Entity{
    id : strconv.Itoa( rand.Intn(1000000)),
    components:make(map[string]Component),
    entityMutex: &sync.RWMutex{}}
}

func (e *Entity) GetID() string{
  return e.id
}

func (e *Entity) AddComponent(component Component) *Component{


    e.entityMutex.RLock()
    defer e.entityMutex.RUnlock()

      e.components[component.GetID()] = component
      component.SetEntity(e)
      component.SetEnabled(true)
      component.Awake()
      if e.try_to_start{
        component.Start()
      }
      //check it has an Entity container
      if e.ECS != nil{

        e.ECS.AddToComponentCache(component.GetID(), e)
      }
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
  e.entityMutex.Lock()
  defer e.entityMutex.Unlock()
  delete(e.components, id)
  if e.ECS != nil{
    e.ECS.RemoveFromComponentCache(component.GetID(), e)

  }
  return componentptr
}

func (e *Entity) GetComponent(id string) (*Component, bool) {
  comp, ok := e.components[id]
  return &comp, ok
}

func (e *Entity) String() string{
  return fmt.Sprintf("%b", e.id)
}

func (e *Entity) SetEnabled(enabled bool){
    e.enabled = enabled
}

func (e *Entity) Awake(){
  e.SetEnabled(true)
  e.try_to_awaken = true
}

func (e *Entity) Start(bulk_start *sync.WaitGroup, finished_start chan<- string){
  //fmt.Println(bulk_start)
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
    e.should_delete = true
}

func (e *Entity) CheckDelete() bool{
  return e.should_delete
}
