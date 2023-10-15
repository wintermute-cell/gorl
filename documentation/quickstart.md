<!-- LTeX: language=en-US -->
# Quickstart

This guide assumes you have already set up your development environment, and
are able to build a program. If you have not yet done that, see 
[Project Setup](/documentation/project-setup.md)

This guide will teach you to create a basic game using the `gorl` framework,
wherein a user-controlled green `Player` aims to collect red `Apples` to
gain points.

## Prerequisites

- Basic knowledge of the Go programming language.
- A configured Go development environment.
- A text editor or IDE of your choice.

## Step-by-Step Guide

### Step 1: Creating a Scene

Most objects in your game world will be represented as an `entity`. To group
multiple entities together into on coherent structure, gorl uses a `scene`.
So let's create one:

Using the bash script:
```bash
./scripts/create_scene.sh my_scene
```

Or by hand:
- Copy `pkg/scenes/scene_template.go` to `pkg/scenes/scene_my_scene.go`
- Replace every occurrence of `Template` with `MyScene`

Now, we register that scene with the scene manager (`scenes.Sm`):
```go
// in main.go

// scenes
my_scene_name := "my_scene"
scenes.Sm.RegisterScene(my_scene_name, &scenes.MyScene{})
scenes.Sm.EnableScene(my_scene_name)
```

Now, we also create a "root entity" for our scene. Entities are managed in a
tree structure, were each entity is the child of another. This "root entity" will be the parent of all entities in our scene.

```go
type MyScene struct {
    // ...
    root_entity *proto.BaseEntity
    // ...
}

// ...

func (scn *MyScene) Init() {
    // ...
	scn.root_entity = &proto.BaseEntity{Name: "MyRootEntity"}
    // ...
```

We create a `proto.BaseEntity` here, which is the simplest implementation of
an entity, from which all other entities extend. Since we just want to use this
entity as a parent, we don't need any specific entity logic, and the
`proto.BaseEntity` suffices.

### Step 2: Creating a player
First, we create new file for the player entity:

Using the bash script:
```bash
./scripts/create_entity2D.sh my_player
```

Or by hand:
- Copy `pkg/entities/entity2d_template.go` to `pkg/entities/entity_my_player.go`
- Replace every occurrence of `Template` with `MyPlayer`

In `pkg/scenes/scene_my_scene.go` create a player entity like this:
```go
// ...
// gem.AddEntity(gem.Root(), scn.scn_root_ent)

player := entities.NewMyPlayerEntity2D(
    rl.NewVector2(100, 100), // position
    0.0, // rotation
    rl.Vector2One(), // scale
    )
gem.AddEntity(
    scn.root_entity, // parented directly to the scene
    player, // our new entity
    )

// logging.Info("DevScene initialized.")
// ...
```

### Step 3: Adding logic to our player

Right now, our player entity exists, but is neither drawn, now does it
calculate anything. Let's fix that.
```go
// in pkg/entities/entity_my_player.go

func (ent *MyPlayerEntity2D) Update() {
	// Required update
	ent.BaseEntity2D.Update()

	// Update logic for the entity
    movement := input.GetMovementVector() // this is setup as WADS controls by default.
    move_speed := 2 * rl.GetFrameTime() // don't forget the delta time!
    ent.SetPosition(
        rl.Vector2Add(
            ent.GetPosition(), rl.Vector2Scale(movement, move_speed),
            ),
        )
}

func (ent *MyPlayerEntity2D) Draw() {
	// Draw logic for the entity
    rl.DrawCircleV(
        ent.GetPosition(), // position
        16, // radius
        rl.Green, // color
        )
}
```

You should now be able to build and run your game, and be presented with a
WASD controlled green circle!

### TODO, apples and scoring and stuff
