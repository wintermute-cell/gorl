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
	transform           Transform2D
	setPositionCallback []func(new_position rl.Vector2)
	setScaleCallback    []func(new_scale rl.Vector2)
	setRotationCallback []func(new_rotation float32)

	// Custom Fields
	// Add fields here for any state that the entity should keep track of
	// ...
}

func NewEntity2D(position rl.Vector2, rotation float32, scale rl.Vector2) *Entity2D {
	new_ent := &Entity2D{
		Entity:    Entity{},
		transform: Transform2D{Position: position, Rotation: rotation, Scale: scale},
	}

	return new_ent
}

func (ent *Entity2D) GetPosition() rl.Vector2 {
	return ent.transform.Position
}

func (ent *Entity2D) SetPosition(new_position rl.Vector2) {
	//// if the parent is an entity2D too, subtract its position to make new_position relative
	//switch e := ent.Entity.GetParent().(type) {
	//case IEntity2D:
	//	new_position = rl.Vector2Subtract(new_position, e.GetPosition())
	//}
	//ent.transform.Position = new_position

	oldPosition := ent.transform.Position
	ent.transform.Position = new_position

	// calling this callback allows us to enhance the functionality of the
	// SetPosition method, without modifying the Entity2D itself. For example,
	// we can use this to apply tree inheritance of transforms through the GEM.
	for _, callback := range ent.setPositionCallback {
		if callback != nil {
			offset := rl.Vector2Subtract(new_position, oldPosition)
			callback(offset)
		}
	}
}

func (ent *Entity2D) GetScale() rl.Vector2 {
	return ent.transform.Scale
}

func (ent *Entity2D) SetScale(new_scale rl.Vector2) {
	oldScale := ent.transform.Scale
	ent.transform.Scale = new_scale

	// see comment in SetPosition for explanation
	for _, callback := range ent.setScaleCallback {
		if callback != nil {
			offset := rl.Vector2DivideV(new_scale, oldScale)
			callback(offset)
		}
	}
}

func (ent *Entity2D) GetRotation() float32 {
	return ent.transform.Rotation
}

func (ent *Entity2D) SetRotation(new_rotation float32) {
	oldRotation := ent.transform.Rotation
	ent.transform.Rotation = new_rotation

	// see comment in SetPosition for explanation
	for _, callback := range ent.setRotationCallback {
		if callback != nil {
			offset := new_rotation - oldRotation
			callback(offset)
		}
	}
}

// AddOnSetPositionCallback adds a callback that is called every time the
// position of the entity is set.
func (ent *Entity2D) AddOnSetPositionCallback(callback func(new_position rl.Vector2)) {
	ent.setPositionCallback = append(ent.setPositionCallback, callback)
}

// AddOnSetScaleCallback adds a callback that is called every time the
// scale of the entity is set.
func (ent *Entity2D) AddOnSetScaleCallback(callback func(new_scale rl.Vector2)) {
	ent.setScaleCallback = append(ent.setScaleCallback, callback)
}

// AddOnSetRotationCallback adds a callback that is called every time the
// rotation of the entity is set.
func (ent *Entity2D) AddOnSetRotationCallback(callback func(new_rotation float32)) {
	ent.setRotationCallback = append(ent.setRotationCallback, callback)
}

// GetTransform returns a pointer to the transform of the entity.
// You should not modify the transform directly, but use the SetPosition,
// SetScale and SetRotation methods instead.
//
// Modifying the transform directly will not trigger the callbacks, skipping
// things like tree inheritance of transforms.
//
// TODO: this might cause problems since we need the direct transform for
// animations, and animations then won't trigger the callbacks. But is this a
// problem or a feature?
func (ent *Entity2D) GetTransform() *Transform2D {
	return &ent.transform
}
