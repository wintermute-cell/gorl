package entities

import (
	"gorl/fw/core/entities"
	input "gorl/fw/core/input/input_event"
	"gorl/fw/physics"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Ensure that EnvironmentEntity implements IEntity.
var _ entities.IEntity = &EnvironmentEntity{}

// Environment Entity
type EnvironmentEntity struct {
	*entities.Entity // Required!

	// Custom Fields
	// Add fields here for any state that the entity should keep track of
	// ...
	cols []*physics.Collider
}

// NewEnvironmentEntity creates a new instance of the EnvironmentEntity.
func NewEnvironmentEntity() *EnvironmentEntity {
	// NOTE: you can modify the constructor to take any parameters you need to
	// initialize the entity.
	new_ent := &EnvironmentEntity{
		Entity: entities.NewEntity("EnvironmentEntity", rl.Vector2Zero(), 0, rl.Vector2One()),

		//col: physics.NewConvexCollider(
		//	rl.Vector2Zero(),
		//	[]rl.Vector2{
		//		rl.NewVector2(0, 0),
		//		rl.NewVector2(0, 100),
		//		rl.NewVector2(100, 100),
		//		rl.NewVector2(100, 0),
		//	},
		//	physics.BodyTypeStatic,
		//),
	}

	return new_ent
}

func (ent *EnvironmentEntity) Init() {
	// Initialization logic for the entity
	// ...
}

func (ent *EnvironmentEntity) Deinit() {
	// De-initialization logic for the entity
	// ...
}

func (ent *EnvironmentEntity) Update() {
	// Update logic for the entity per frame
	// ...
}

func (ent *EnvironmentEntity) Draw() {
	//rl.DrawRectangleLines(
	//	int32(ent.col.GetPosition().X),
	//	int32(ent.col.GetPosition().Y),
	//	int32(100),
	//	int32(100),
	//	rl.Red,
	//)

	rl.DrawCircleV(rl.NewVector2(100, 100), 5, rl.Blue)
}

func (ent *EnvironmentEntity) OnInputEvent(event *input.InputEvent) bool {
	// Logic to run when an input event is received.
	// Return false if the event was consumed and should not be propagated
	// further.
	// ...
	return true
}
