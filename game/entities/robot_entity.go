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

	Velocity     rl.Vector2
	Acceleration rl.Vector2
	Target       rl.Vector2
	MaximumSpeed float32
	MaximumForce float32

	Color rl.Color
}

// NewRobotEntity creates a new instance of the RobotEntity.
func NewRobotEntity() *RobotEntity {
	new_ent := &RobotEntity{
		Entity:       entities.NewEntity("RobotEntity", rl.Vector2Zero(), 0, rl.Vector2One()),
		Velocity:     rl.Vector2Zero(),
		Acceleration: rl.Vector2Zero(),
		Target:       rl.Vector2Zero(),
		MaximumSpeed: 60,
		MaximumForce: 0.7,
		Color:        rl.NewColor(uint8(rand.Int()%255), uint8(rand.Int()%255), uint8(rand.Int()%255), 255),
	}
	return new_ent
}

func (ent *RobotEntity) Init() {
}

func (ent *RobotEntity) Deinit() {
}

func (ent *RobotEntity) Seek() {
	// find the target
	force := rl.Vector2Subtract(ent.Target, ent.GetPosition())
	// limit the speed
	force = rl.Vector2ClampValue(force, float32(ent.MaximumSpeed), float32(ent.MaximumSpeed))
	// calculate the steering
	force = rl.Vector2Subtract(force, ent.Velocity)
	// limit the steering by the MaximumForce
	force = rl.Vector2ClampValue(force, 0, ent.MaximumForce)
	ent.ApplyForce(force)
}

func (ent *RobotEntity) ApplyForce(force rl.Vector2) {
	ent.Acceleration = rl.Vector2Add(ent.Acceleration, force)
}

func (ent *RobotEntity) Update() {
	// MOVEMENT
	ent.Seek()
	ent.Velocity = rl.Vector2Add(ent.Velocity, ent.Acceleration)
	ent.Velocity = rl.Vector2ClampValue(ent.Velocity, 0, ent.MaximumSpeed)
	ent.SetPosition(rl.Vector2Add(ent.GetPosition(), rl.Vector2Scale(ent.Velocity, rl.GetFrameTime())))
	ent.Acceleration = rl.Vector2Zero()

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
