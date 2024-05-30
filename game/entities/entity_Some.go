package entities

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	input "gorl/fw/core/input/input_event"
	"gorl/fw/core/entities/proto"
)

// Some Entity
type SomeEntity2D struct {
	// Required fields
	*proto.Entity2D

	// Custom Fields
	// Add fields here for any state that the entity should keep track of
	// ...
}

func NewSomeEntity2D(position rl.Vector2, rotation float32, scale rl.Vector2) *SomeEntity2D {
	new_ent := &SomeEntity2D{
		Entity2D: proto.NewEntity2D(position, rotation, scale),

		// Initialize custom fields here
		// ...
	}
	return new_ent
}

func (ent *SomeEntity2D) Init() {
	// Initialization logic for the entity
	// ...
}

func (ent *SomeEntity2D) Deinit() {
	// De-initialization logic for the entity
	// ...
}

func (ent *SomeEntity2D) Update() {
	// Update logic for the entity
	// ...
}

func (ent *SomeEntity2D) Draw() {
	// Draw logic for the entity
	// ...
}

func (ent *SomeEntity2D) OnInputEvent(event *input.InputEvent) bool {
	// Logic to run when an input event is received.
	// Return false if the event was consumed and should not be propagated
	// further.
	// ...
	return true
}
