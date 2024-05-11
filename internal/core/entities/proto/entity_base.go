package proto

import (
	"gorl/internal/util"
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

func (ent *Entity) String() string {
	return ent.Name
}

func (ent *Entity) Init() {
	// Initialization logic for the entity
	// ...
}

func (ent *Entity) Deinit() {
	// De-initialization logic for the entity
	// ...
}

func (ent *Entity) Update() {
	// Update logic for the entity
	// ...
}

func (ent *Entity) FixedUpdate() {

}

func (ent *Entity) Draw() {
	// Draw logic for the entity
	// ...
}

func (ent *Entity) DrawGUI() {
	// GUI Draw logic for the entity
	// ...
}

func (ent *Entity) OnChildAdded(child IEntity) {
	// Logic to run when a child is added to this entity
	// ...
}

func (ent *Entity) OnChildRemoved(child IEntity) {
	// Logic to run when a child is removed from this entity
	// ...
}

// AddChild adds a child to this entity.
func (ent *Entity) AddChild(child IEntity) {
	ent.children = append(ent.children, child)
}

// RemoveChild removes a child from this entity.
func (ent *Entity) RemoveChild(child IEntity) {
	idx := util.SliceIndex(ent.children, child)
	if idx > -1 {
		ent.children = util.SliceDelete(ent.children, idx, idx+1)
	}
}

// GetChildren returns the children of this entity.
func (ent *Entity) GetChildren() []IEntity {
	return ent.children
}

// GetParent sets the parent of this entity.
func (ent *Entity) GetParent() IEntity {
	return ent.parent
}

// SetParent sets the parent of this entity.
func (ent *Entity) SetParent(parent IEntity) {
	ent.parent = parent
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

func (ent *Entity) GetName() string {
	if ent.Name == "" {
		ent.Name = "UnnamedEntity"
	}
	return ent.Name
}
