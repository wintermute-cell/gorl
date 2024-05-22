package entities

import (
	"gorl/fw/core/entities/proto"
	input "gorl/fw/core/input/input_event"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// TestEntity Entity
type TestEntityEntity struct {
	// Required fields
	*proto.Entity

	// Custom Fields
	// Add fields here for any state that the entity should keep track of
	// ...
}

func NewTestEntityEntity() *TestEntityEntity {
	new_ent := &TestEntityEntity{
		Entity: proto.NewEntity(),

		// Initialize custom fields here...
	}
	return new_ent
}

func (ent *TestEntityEntity) Init() {
	// Initialization logic for the entity
	// ...
}

func (ent *TestEntityEntity) Deinit() {
	// De-initialization logic for the entity
	// ...
}

func (ent *TestEntityEntity) Update() {
	// Update logic for the entity
	// ...
}

func (ent *TestEntityEntity) Draw() {
	// oscillate y with sin

	y := 100 + 50*math.Sin(rl.GetTime())

	rl.DrawCircleV(rl.NewVector2(100, float32(y)), 10, rl.Red)
}

func (ent *TestEntityEntity) OnInputEvent(event *input.InputEvent) bool {
	// Logic to run when an input event is received.
	// Return false if the event was consumed and should not be propagated
	// further.
	// ...
	return true
}
