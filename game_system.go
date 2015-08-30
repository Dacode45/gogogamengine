package gogogamengine

import (
  "sync"
  "fmt"
)


type GameSystem struct{
  BaseSystem
}


func (gS *GameSystem) Run(world *World, system_start <-chan bool, system_started chan<- bool, system_end chan<- bool, force_end <-chan bool){
  gS.world = world
  gS.awakeProtocol()
  if <-system_start{
    all_entities_started := make(chan bool)
    go gS.startProtocol(gS.world.ECS.Entities, all_entities_started)
    //fmt.Println("Game System Awake")
     <- all_entities_started
     gS.started = true
    //fmt.Printf("Game System Started: %v", gS.started)
    system_started <- gS.started
    close(all_entities_started)
    //system_started <- ss
    gS.beforeUpdateProtocol()
    gS.CleanUp()

    L:
    for{ //change this into a select
      //fmt.Println("Update")
      select{
        case  <- force_end:
            break L
        default:
            if world.IsAlive(){
              //fmt.Println("Ending")
              update_protocol_finished := make(chan bool)
              //fmt.Println("hi")
              go gS.updateProtocol(gS.world.ECS.Entities, update_protocol_finished)
              <-update_protocol_finished
              gS.CleanUp() //death starts
            }else{
              break L
            }
      }
    }
  }
  gS.ended = true
  system_end <- gS.ended
}

func (gS *GameSystem) awakeProtocol(){
  //Protocol safe. no CONCURENCY
  /*THESE LINES WILL ALWAYS RUN FIRST AND DO NOT NEED CONCURENCY WRORIES*/
  //fmt.Println(gS.world.ECS.Entities)
  for _, e := range gS.world.ECS.Entities {
    e.Awake()
  }
}

//copy the map of pointers so real map can still be accessed
func (gS *GameSystem) startProtocol(entities map[string]*Entity, all_entities_started chan<- bool) {
  var bulk_start sync.WaitGroup
  bulk_start.Add(len(entities))
  finished_start := make(chan string, len(entities))
  //internal function for when everything has started
  finished_start_protocol := make(chan bool)
  go func(){
    bulk_start.Wait()
    close(finished_start)
  }()
  //go routine for handling enties that finished or failed to finish
  go func(){
    failed_starts := make([]string, 0)
    for id := range finished_start{
      if id != ""{
        failed_starts = append(failed_starts,id)
      }
    }

    //give each failed to start Component
    //one more chance to try incase they depended
    //on another entity starting before this did.
    if len(failed_starts) != 0{
      fmt.Println(failed_starts, "Failed to start again")
      var second_chance_bulk sync.WaitGroup
      second_chance_bulk.Add(len(failed_starts))
      second_chance_start := make(chan string, len(failed_starts))
      go func(){
        second_chance_bulk.Wait()
        close(second_chance_start)
      }()
      for _, id := range failed_starts{
        go entities[id].Start(&second_chance_bulk, second_chance_start)
      }
      //If they fail again, disable
      failed_second_starts := make([]string, 0)
      for id := range second_chance_start{
        if id != ""{
          failed_second_starts = append(failed_second_starts, id)
        }
      }

      //disable them
      if len(failed_second_starts) != 0{
        fmt.Println(failed_second_starts, " Failed to start again")
        for _, id := range failed_second_starts{
          entities[id].SetEnabled(false)
        }

        finished_start_protocol <- false
      }
    }
    finished_start_protocol <- true
  }()

  for _, e := range entities{
    go e.Start(&bulk_start, finished_start)
  }

  start_protocol_error :=  <- finished_start_protocol
  //fmt.Println("finished start protocol", start_protocol_error)
  all_entities_started <- start_protocol_error

}

func (gS *GameSystem) beforeUpdateProtocol(){

}

func (gS *GameSystem) updateProtocol(entities map[string]*Entity, finished_update_protocol chan<- bool){
  //fmt.Println("\n\nStarting New update protocol")
  var bulk_update sync.WaitGroup
  bulk_update.Add(len(entities))
  //Entity update limit controls how many actual updates are being run
  //all of them will be run but only so many go routines
  //will actually be iterating though the entity Component list
  can_update := make(chan bool, gS.world.GetEntityUpdateLimit())
  finished_update := make(chan string, len(entities))

  defer func(){
    bulk_update.Wait()
    close(finished_update)


  }()

  for _, e := range entities{
    go e.Update(&bulk_update, can_update, finished_update)
  }
  //func for handling finished update
  go func(){
    failed_updates := make([]string, 0)
    //fmt.Println("Waiting for finished update")
    for id := range finished_update{
      //allow more updates to happen
      //fmt.Println("Enity", id, "updated")

      can_update <-true
      if id != ""{
        failed_updates = append(failed_updates, id)
      }
    }
    close(can_update)
    //fmt.Println("finished update finished")

    //Do something with failed updates. Print for now
    //fmt.Println(failed_updates, "failed to update")
    finished_update_protocol <- true
  }()
  //now add values onto the can update list
  for i := uint(0); i < gS.world.GetEntityUpdateLimit(); i++{
    can_update <- true //will block when can update true i think
  }
  //fmt.Println("finished update protocol")


}

func (gS *GameSystem) CleanUp(){
  for id, e := range gS.world.ECS.Entities{
    if e.CheckDelete(){
      gS.world.ECSLock.Lock()
        delete(gS.world.ECS.Entities, id)
      gS.world.ECSLock.Unlock()
    }
  }
}
