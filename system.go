package goengine


type System interface{
  Run(world *World, world_start chan<-bool, system_started <-chan bool)
}
