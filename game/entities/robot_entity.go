package entities

import (
	"gorl/fw/core/entities"
	input "gorl/fw/core/input/input_event"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Ensure that RobotEntity implements IEntity.
var _ entities.IEntity = &RobotEntity{}

// Robot Entity
type RobotEntity struct {
	*entities.Entity // Required!

	// Custom Fields
	// Add fields here for any state that the entity should keep track of
	// ...
}

// NewRobotEntity creates a new instance of the RobotEntity.
func NewRobotEntity(position rl.Vector2, rotation float32, scale rl.Vector2) *RobotEntity {
	// NOTE: you can modify the constructor to take any parameters you need to
	// initialize the entity.
	new_ent := &RobotEntity{
		Entity: entities.NewEntity("RobotEntity", position, rotation, scale),

		// Initialize custom fields here
		// ...
	}
	return new_ent
}

func (ent *RobotEntity) Init() {
	// Initialization logic for the entity
	// ...
}

func (ent *RobotEntity) Deinit() {
	// De-initialization logic for the entity
	// ...
}

func (ent *RobotEntity) Update() {
	// Update logic for the entity per frame
	// ...
}

func (ent *RobotEntity) Draw() {
	// Draw logic for the entity
	// ...
}

func (ent *RobotEntity) OnInputEvent(event *input.InputEvent) bool {
	// Logic to run when an input event is received.
	// Return false if the event was consumed and should not be propagated
	// further.
	// ...
	return true
}
