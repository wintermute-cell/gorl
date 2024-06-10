package entities

import (
	"gorl/fw/core/entities"
	input "gorl/fw/core/input/input_event"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Ensure that AnimSingleRobotToTargetEntity implements IEntity.
var _ entities.IEntity = &AnimSingleRobotToTargetEntity{}

// AnimSingleRobotToTarget Entity
type AnimSingleRobotToTargetEntity struct {
	*entities.Entity // Required!

	// Custom Fields
	// Add fields here for any state that the entity should keep track of
	// ...
}

// NewAnimSingleRobotToTargetEntity creates a new instance of the AnimSingleRobotToTargetEntity.
func NewAnimSingleRobotToTargetEntity(position rl.Vector2, rotation float32, scale rl.Vector2) *AnimSingleRobotToTargetEntity {
	// NOTE: you can modify the constructor to take any parameters you need to
	// initialize the entity.
	new_ent := &AnimSingleRobotToTargetEntity{
		Entity: entities.NewEntity("AnimSingleRobotToTargetEntity", position, rotation, scale),

		// Initialize custom fields here
		// ...
	}
	return new_ent
}

func (ent *AnimSingleRobotToTargetEntity) Init() {
	// Initialization logic for the entity
	// ...
}

func (ent *AnimSingleRobotToTargetEntity) Deinit() {
	// De-initialization logic for the entity
	// ...
}

func (ent *AnimSingleRobotToTargetEntity) Update() {
	// Update logic for the entity per frame
	// ...
}

func (ent *AnimSingleRobotToTargetEntity) Draw() {
	// Draw logic for the entity
	// ...
}

func (ent *AnimSingleRobotToTargetEntity) OnInputEvent(event *input.InputEvent) bool {
	// Logic to run when an input event is received.
	// Return false if the event was consumed and should not be propagated
	// further.
	// ...
	return true
}
