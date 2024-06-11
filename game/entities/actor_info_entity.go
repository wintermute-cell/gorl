package entities

import (
	"fmt"
	"gorl/fw/core/entities"
	input "gorl/fw/core/input/input_event"
	"gorl/fw/core/settings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Ensure that ActorInfoEntity implements IEntity.
var _ entities.IEntity = &ActorInfoEntity{}

// ActorInfo Entity
type ActorInfoEntity struct {
	*entities.Entity // Required!

	// Custom Fields
	// Add fields here for any state that the entity should keep track of
	// ...
	observed entities.IEntity
}

// NewActorInfoEntity creates a new instance of the ActorInfoEntity.
func NewActorInfoEntity() *ActorInfoEntity {
	// NOTE: you can modify the constructor to take any parameters you need to
	// initialize the entity.
	new_ent := &ActorInfoEntity{
		Entity: entities.NewEntity("ActorInfoEntity", rl.Vector2Zero(), 0, rl.Vector2One()),

		// Initialize custom fields here
		// ...
	}
	return new_ent
}

func (ent *ActorInfoEntity) Init() {
	// Initialization logic for the entity
	// ...
}

func (ent *ActorInfoEntity) Deinit() {
	// De-initialization logic for the entity
	// ...
}

func (ent *ActorInfoEntity) Update() {
	// Update logic for the entity per frame
	// ...
}

func (ent *ActorInfoEntity) Draw() {
	// Draw logic for the entity
	// ...
	textPos := rl.Vector2Add(ent.GetPosition(), rl.NewVector2(20, float32(settings.CurrentSettings().ScreenHeight-30)))
	if ent.observed == nil {
		rl.DrawText("Entity: None", int32(textPos.X), int32(textPos.Y), 20, rl.Lime)
		return
	}

	rl.DrawText(fmt.Sprintf("Entity: %v", "oof"), 20, int32(textPos.X), int32(textPos.Y), rl.Lime)
}

func (ent *ActorInfoEntity) OnInputEvent(event *input.InputEvent) bool {
	// Logic to run when an input event is received.
	// Return false if the event was consumed and should not be propagated
	// further.
	// ...
	return true
}

// GetObservedEntity returns the entity that this ActorInfoEntity is showing
// info for.
func (ent *ActorInfoEntity) SetObservedEntity(observed entities.IEntity) {
	ent.observed = observed
}
