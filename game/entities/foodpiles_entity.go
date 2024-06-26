package entities

import (
	"gorl/fw/core/entities"
	input "gorl/fw/core/input/input_event"
	"gorl/game/code"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Ensure that FoodpilesEntity implements IEntity.
var _ entities.IEntity = &FoodpilesEntity{}

// Foodpiles Entity
type FoodpilesEntity struct {
	*entities.Entity // Required!

	// Custom Fields
	Foodpiles *code.FoodPiles
}

// NewFoodpilesEntity creates a new instance of the FoodpilesEntity.
func NewFoodpilesEntity() *FoodpilesEntity {
	// NOTE: you can modify the constructor to take any parameters you need to
	// initialize the entity.
	new_ent := &FoodpilesEntity{
		Entity:    entities.NewEntity("FoodpilesEntity", rl.Vector2Zero(), 0, rl.Vector2One()),
		Foodpiles: code.NewFoodPiles(),
	}

	return new_ent
}

func (ent *FoodpilesEntity) Init() {
	// Initialization logic for the entity
	// ...
}

func (ent *FoodpilesEntity) Deinit() {
	// De-initialization logic for the entity
	// ...
}

func (ent *FoodpilesEntity) Update() {
	// Update logic for the entity per frame
	// ...
}

func (ent *FoodpilesEntity) Draw() {
	ent.Foodpiles.Draw()
}

func (ent *FoodpilesEntity) OnInputEvent(event *input.InputEvent) bool {
	// Logic to run when an input event is received.
	// Return false if the event was consumed and should not be propagated
	// further.
	// ...
	return true
}
