package proto

import (
	input "gorl/fw/core/input/input_event"
)

// This checks at compile time if the interface is implemented
var _ IEntity = (*Entity)(nil)

// Base Entity
type Entity struct {
	// Required fields
	children  []IEntity
	parent    IEntity
	drawIndex int32
	Name      string

	// Custom Fields
	// Add fields here for any state that the entity should keep track of
	// ...
}

func NewEntity() *Entity {
	new_ent := &Entity{
		children:  []IEntity{},
		parent:    nil,
		drawIndex: 0,
		Name:      "",
	}

	return new_ent
}

func (ent *Entity) String() string {
	return ent.Name
}

// Init is called when the entity is added to the Gem.
func (ent *Entity) Init() {
	// Initialization logic for the entity
	// ...
}

// Deinit is called just before the entity is removed from the Gem.
func (ent *Entity) Deinit() {
	// De-initialization logic for the entity
	// ...
}

// Update is called every frame, before the draw logic.
func (ent *Entity) Update() {
	// Update logic for the entity
	// ...
}

// FixedUpdate is called every physics frame, if physics are enabled.
func (ent *Entity) FixedUpdate() {

}

// Draw is called every frame, after the update logic.
func (ent *Entity) Draw() {
	// Draw logic for the entity
	// ...
}

// FIXME: do we keep this?
func (ent *Entity) DrawGUI() {
	// GUI Draw logic for the entity
	// ...
}

// OnChildAdded is called every time an entity is added to the Gem with this
// Entity as the parent.
func (ent *Entity) OnChildAdded(child IEntity) {
	// Logic to run when a child is added to this entity
	// ...
}

// OnChildRemoved is called every time a child of this entity is removed from
// the Gem.
func (ent *Entity) OnChildRemoved(child IEntity) {
	// Logic to run when a child is removed from this entity
	// ...
}

// OnInputEvent is called by the event handling system when an input event occurs.
// The entity must decide if it should handle the event or not.
func (ent *Entity) OnInputEvent(*input.InputEvent) bool {
	// Logic to run when an input event is received.
	// Return false if the event was consumed and should not be propagated
	// further.
	// ...
	return true
}

// GetDrawIndex returns the draw index of this entity. Entities with a higher
// index are drawn in front of entities with a lower index.
func (ent *Entity) GetDrawIndex() int32 {
	return ent.drawIndex
}

// SetDrawIndex sets the draw index of this entity. Entities with a higher
// index are drawn in front of entities with a lower index.
func (ent *Entity) SetDrawIndex(index int32) {
	ent.drawIndex = index
}

// GetName returns the name of the entity.
// Returns "UnnamedEntity" if the entity has no name.
func (ent *Entity) GetName() string {
	if ent.Name == "" {
		ent.Name = "UnnamedEntity"
	}
	return ent.Name
}
