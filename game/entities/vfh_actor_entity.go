package entities

import (
	"gorl/fw/core/entities"
	input "gorl/fw/core/input/input_event"
	"gorl/fw/core/logging"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Ensure that VfhActorEntity implements IEntity.
var _ entities.IEntity = &VfhActorEntity{}

// VfhActor Entity
type VfhActorEntity struct {
	*entities.Entity // Required!

	// Custom Fields
	// Add fields here for any state that the entity should keep track of
	// ...
}

// NewVfhActorEntity creates a new instance of the VfhActorEntity.
func NewVfhActorEntity(position rl.Vector2) *VfhActorEntity {
	// NOTE: you can modify the constructor to take any parameters you need to
	// initialize the entity.
	new_ent := &VfhActorEntity{
		Entity: entities.NewEntity("VfhActorEntity", position, 0, rl.Vector2One()),
	}
	return new_ent
}

func (ent *VfhActorEntity) Init() {
	// Initialization logic for the entity
	// ...
}

func (ent *VfhActorEntity) Deinit() {
	// De-initialization logic for the entity
	// ...
}

func (ent *VfhActorEntity) Update() {
	// Update logic for the entity per frame
	// ...
}

func (ent *VfhActorEntity) Draw() {
	// Draw logic for the entity
	// ...
}

func (ent *VfhActorEntity) OnInputEvent(event *input.InputEvent) bool {
	// Logic to run when an input event is received.
	// Return false if the event was consumed and should not be propagated
	// further.
	// ...
	if event.Action == input.ActionMoveUp {
		logging.Debug("OOF")
	}
	return true
}
