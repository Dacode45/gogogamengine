package goengine

import(
  "fmt"
  "errors"
)

type ComponentCollection struct{
  Components map[string]Component

}

func NewComponentCollection() *ComponentCollection{
  var CCS = ComponentCollection{}
  CCS.Components = make(map[string]Component)
  return &CCS

}


func (c *ComponentCollection) RegisterComponent(comp Component) error {
  //Do stuff for the Component
  component := comp.(*BaseComponent)
  if _, ok := c.Components[component.id]; ok{
    return errors.New("Component with this name Already Exist!")
  }
  c.Components[component.id] = component
  return nil
}

func (cc *ComponentCollection) GetComponent(id string) (Component, bool){
  c, ok := cc.Components[id]
  return c, ok
}


type Component interface{

  GetId() string
  GetDesc() string
  GetEntity() *Entity
  SetEntity(e *Entity)
  Awake()
  Start()
  FixedUpdate()
  Update(delta_time float64)
  IsEnabled() bool
  SetEnabled(enabled bool)
  Register(World) error
}

type BaseComponent struct{
  id string
  description string
  entity *Entity
  enabled bool
}

//if program crashed finished_update is false


func (bC *BaseComponent) GetId() string{

  return bC.id
}

func (bC *BaseComponent) GetDesc() string{

  return bC.description
}

func (bC *BaseComponent) IsEnabled() bool{
  return bC.enabled
}

func (bC *BaseComponent) GetEntity() *Entity{
  return &(*bC.entity) //returns new pointer to entity
}

func (bC *BaseComponent) SetEntity(e *Entity){
  bC.entity = e
}

func (bC *BaseComponent) Awake(){

}
func (bC *BaseComponent) Start(){

}
func (bC *BaseComponent) FixedUpdate(){

}
func (bC *BaseComponent) Update(delta_time float64){
  fmt.Println(bC, "Updating")
}

func (bC *BaseComponent) SetEnabled(enabled bool){
    if !bC.enabled && enabled{ //go from not enabled to enabled.
      bC.enabled = enabled
      bC.Start()
    }else{

      bC.enabled = enabled
    }

}

func (bC *BaseComponent) Register(world World) error{
  //check to make sure it has all the right properties.
  if bC.id == "" || bC.description == ""{
    return errors.New("Component Must have an Id and Description")
  }
  err := world.CCS.RegisterComponent(bC)
  if err != nil{
    bC.SetEnabled(false)
    return err
  }
  bC.SetEnabled(true)
  return nil
}

type ExampleComponent struct{
  BaseComponent
  word string
}

func (ec *ExampleComponent) Awake(){
  ec.id = "Example"
  ec.word = "Hello World"
}

func (ec *ExampleComponent) SayHello(){
  fmt.Println(ec.word)
}

func (ec *ExampleComponent) Update(delta_time int){
  ec.SayHello()
}
