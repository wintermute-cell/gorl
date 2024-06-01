package entities

import (
	"gorl/fw/core/entities"
	input "gorl/fw/core/input/input_event"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Ensure that TemplateEntity implements IEntity.
var _ entities.IEntity = &TemplateEntity{}

// Template Entity
type TemplateEntity struct {
	// Required fields
	*entities.Entity

	// Custom Fields
	// Add fields here for any state that the entity should keep track of
	// ...
}

func NewTemplateEntity(position rl.Vector2, rotation float32, scale rl.Vector2) *TemplateEntity {
	new_ent := &TemplateEntity{
		Entity: entities.NewEntity("TemplateEntity", position, rotation, scale),

		// Initialize custom fields here
		// ...
	}
	return new_ent
}

func (ent *TemplateEntity) Init() {
	// Initialization logic for the entity
	// ...
}

func (ent *TemplateEntity) Deinit() {
	// De-initialization logic for the entity
	// ...
}

func (ent *TemplateEntity) Update() {
	// Update logic for the entity
	// ...
}

func (ent *TemplateEntity) Draw() {
	// Draw logic for the entity
	// ...
}

func (ent *TemplateEntity) OnInputEvent(event *input.InputEvent) bool {
	// Logic to run when an input event is received.
	// Return false if the event was consumed and should not be propagated
	// further.
	// ...
	return true
}
