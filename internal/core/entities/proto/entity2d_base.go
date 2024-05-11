package proto

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// This checks at compile time if the interface is implemented
var _ IEntity2D = (*Entity2D)(nil)

// Base Entity
type Entity2D struct {
	// Required fields
	Entity
	Transform Transform2D

	// Custom Fields
	// Add fields here for any state that the entity should keep track of
	// ...
}

func (ent *Entity2D) GetPosition() rl.Vector2 {
	p := ent.Transform.Position
	// if the parent is an entity2D too, use its position as a base
	switch e := ent.Entity.GetParent().(type) {
	case IEntity2D:
		p = rl.Vector2Add(p, e.GetPosition())
	}
	return p
}

func (ent *Entity2D) SetPosition(new_position rl.Vector2) {
	// if the parent is an entity2D too, subtract its position to make new_position relative
	switch e := ent.Entity.GetParent().(type) {
	case IEntity2D:
		new_position = rl.Vector2Subtract(new_position, e.GetPosition())
	}
	ent.Transform.Position = new_position
}

func (ent *Entity2D) GetScale() rl.Vector2 {
	return ent.Transform.Scale
}

func (ent *Entity2D) SetScale(new_scale rl.Vector2) {
	ent.Transform.Scale = new_scale
}

func (ent *Entity2D) GetRotation() float32 {
	return ent.Transform.Rotation
}

func (ent *Entity2D) SetRotation(new_rotation float32) {
	ent.Transform.Rotation = new_rotation
}

func (ent *Entity2D) GetTransform() *Transform2D {
	return &ent.Transform
}
