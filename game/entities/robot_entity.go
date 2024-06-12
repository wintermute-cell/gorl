package entities

import (
	"gorl/fw/core/entities"
	input "gorl/fw/core/input/input_event"
	"math/rand"

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
	Color     rl.Color
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
	ent.speed = float32(rand.Int()%500 + 50)
	// TODO: remove
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
	rl.DrawCircleV(ent.GetPosition(), 10, ent.Color)
}

func (ent *RobotEntity) OnInputEvent(event *input.InputEvent) bool {
	return true
}

// Returns the tile position as a vector2
func (ent *RobotEntity) GetTilePosition() rl.Vector2 {
	tilePosition := rl.Vector2Scale(ent.GetPosition(), 1/40.0)
	tilePosition = rl.NewVector2(
		float32(int(tilePosition.X)),
		float32(int(tilePosition.Y)),
	)
	return tilePosition
}

// TODO: Add a vector2 to the vectorpool (?) and calculate the resulting direction.
// For now this is just like SetDirection(rl.Vector2) would be
func (ent *RobotEntity) AddDirectionVector(dir rl.Vector2) {
	// TODO: steering behaviour
	ent.direction = rl.Vector2Normalize(dir)
}
