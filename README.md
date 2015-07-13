Okay How is this going to work from bottom up.

Every thing is an Entity, and an Entity is just a collection Components.
Components are Structs that have 4 functions
Awake (Called when attached to the world and has access to world functions)
Start (Called when the game starts and when renenabled)
Update (Iterate through all entity list and call this)

There are Two Main Components.
Transform
  Position: struct {x,y}
  Direction: struct {x,y,z}


Collision
  RigidBody: struct{width, height, dept}
  Velocity: struct{x,y,z}
  FixedUpdate() called after the physics step
  OnCollision Enter(collider)
  OnCollision Exit(collider)

Entities are controlled by systems.
4 Main sistyem

Game Runes the Awake Start and Update
Collisions handles physics

Render

UserInput Always listening.

Systems get global variables from the World Object

World Object is controlled by our World func that
is a single go routine syncronizing the systems control
of world state.
It has acess to the
global gameStart channel and will call it once all Entities have started.
It's forloop will
First clean up all entites that have been marked for deletion
Call start on all recently added entites


Call update on every Entity
<-wait on the update finish channel that all Entites emit
