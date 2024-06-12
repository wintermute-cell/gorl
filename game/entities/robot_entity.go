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
	speed     float32
	direction rl.Vector2
}

// NewRobotEntity creates a new instance of the RobotEntity.
func NewRobotEntity() *RobotEntity {
	// NOTE: you can modify the constructor to take any parameters you need to
	// initialize the entity.
	new_ent := &RobotEntity{
		Entity: entities.NewEntity("RobotEntity", rl.Vector2Zero(), 0, rl.Vector2One()),

		// Initialize custom fields here
		// ...
	}
	return new_ent
}

func (ent *RobotEntity) Init() {
	ent.speed = 20
	ent.direction = rl.NewVector2(4, 1)
	ent.direction = rl.Vector2Normalize(ent.direction)
}

func (ent *RobotEntity) Deinit() {
}

func (ent *RobotEntity) Update() {
	// MOVEMENT
	distanceTraveled := rl.Vector2Scale(rl.Vector2Scale(ent.direction, ent.speed), rl.GetFrameTime())
	ent.SetPosition(rl.Vector2Add(ent.GetPosition(), distanceTraveled))
}

func (ent *RobotEntity) Draw() {
	rl.DrawCircleV(ent.GetPosition(), 20, rl.Green)
}

func (ent *RobotEntity) OnInputEvent(event *input.InputEvent) bool {
	return true
}

// TODO: Add a vector2 to the vectorpool and calculate the resulting direction.
// For now this is just like SetDirection(rl.Vector2) would be
func (ent *RobotEntity) AddDirectionVector(dir rl.Vector2) {
	ent.direction = dir
	ent.direction = rl.Vector2Normalize(dir)
}
