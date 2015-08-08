package goengine

import "errors"

type System interface{
  Run(world *World, system_start chan bool, system_end chan<- bool, force_end <-chan bool )
  Started() bool
  Ended() bool
  GetID() string
  SetID(string) error
}

type BaseSystem struct{
  started bool
  ended bool
  id string
}

func (gS *BaseSystem) Started() bool{
  return gS.started
}

func (gS *BaseSystem) Ended() bool{
  return gS.ended
}

func (gS *BaseSystem) GetID() string{
  return gS.id
}

func (gS *BaseSystem) SetID(id string) error{
  if gS.id == ""{
    gS.id = id
    return nil
  }
  return errors.New("Cannot change System ID once set")
}
