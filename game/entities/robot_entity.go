package entities

import (
	"gorl/fw/core/entities"
	input "gorl/fw/core/input/input_event"
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Ensure that RobotEntity implements IEntity.
var _ entities.IEntity = &RobotEntity{}

// Robot Entity
type RobotEntity struct {
	*entities.Entity // Required!

	Velocity               rl.Vector2
	Acceleration           rl.Vector2
	CurrentTarget          rl.Vector2 // the next tile
	FinalTarget            rl.Vector2 // the target tile
	MaximumSpeed           float32
	MaximumForce           float32
	SlowDownDistance       float32
	ObstacleDetectionRange float32
	VectorToObstacle       rl.Vector2
	AvoidanceVelocity      rl.Vector2
	AvoidanceForce         float32

	Color rl.Color
}

// NewRobotEntity creates a new instance of the RobotEntity.
func NewRobotEntity() *RobotEntity {
	new_ent := &RobotEntity{
		Entity:                 entities.NewEntity("RobotEntity", rl.Vector2Zero(), 0, rl.Vector2One()),
		Velocity:               rl.Vector2Zero(),
		Acceleration:           rl.Vector2Zero(),
		CurrentTarget:          rl.Vector2Zero(),
		FinalTarget:            rl.Vector2Zero(),
		MaximumSpeed:           150,
		MaximumForce:           0.7,
		SlowDownDistance:       300,
		ObstacleDetectionRange: 150,
		VectorToObstacle:       rl.Vector2Zero(),
		AvoidanceVelocity:      rl.Vector2Zero(),
		AvoidanceForce:         0.65,
		Color:                  rl.NewColor(uint8(rand.Int()%255), uint8(rand.Int()%255), uint8(rand.Int()%255), 255),
	}
	return new_ent
}

func (ent *RobotEntity) Init() {
}

func (ent *RobotEntity) Deinit() {
}

func (ent *RobotEntity) FindClosestWall(obstacles []rl.Vector2) rl.Vector2 {
	// TODO:
	closestObstacle := rl.Vector2Zero()
	closestObstacleLength := math.MaxFloat64

	for _, obstacle := range obstacles {
		vecToObstacle := rl.Vector2Subtract(obstacle, ent.GetPosition())
		currentLength := float64(rl.Vector2Length(vecToObstacle))
		// if float32(closestObstacleLength) < 150 && currentLength < float32(closestObstacleLength) {
		if currentLength < closestObstacleLength {
			closestObstacle = vecToObstacle
			closestObstacleLength = float64(currentLength)
		}
	}

	// for some reason this works, but not if the check is above in the if clause of the loop
	if closestObstacleLength > float64(ent.ObstacleDetectionRange) {
		closestObstacle = rl.Vector2Zero()
	}

	// for debug purposes
	ent.VectorToObstacle = closestObstacle

	// NOTE: somehow this all just magically works yeey ?? \./

	// limit the steering by the AvoidanceForce
	closestObstacle = rl.Vector2ClampValue(closestObstacle, 0, float32(ent.AvoidanceForce))

	// this should be the force pushing the robot away from the obstacle
	// so we have to multiply it with -1
	return rl.Vector2Scale(closestObstacle, -1)
}

// like seek, but slows down the closer the target is
func (ent *RobotEntity) Arrive() rl.Vector2 {
	// find the target
	force := rl.Vector2Subtract(ent.CurrentTarget, ent.GetPosition())
	// limit the speed
	force = rl.Vector2ClampValue(force, float32(ent.MaximumSpeed), float32(ent.MaximumSpeed))

	// calculating distance to the target
	distanceToTarget := rl.Vector2Length(rl.Vector2Subtract(ent.FinalTarget, ent.GetPosition()))
	// at a certain threshold we want to slow down
	if distanceToTarget < ent.SlowDownDistance {
		scaleFactor := distanceToTarget / ent.SlowDownDistance
		force = rl.Vector2Scale(force, scaleFactor)
	}

	// calculate the steering
	force = rl.Vector2Subtract(force, ent.Velocity)
	// limit the steering by the MaximumForce
	force = rl.Vector2ClampValue(force, 0, ent.MaximumForce)

	return force
}

func (ent *RobotEntity) ApplyForce(force rl.Vector2) {
	ent.Acceleration = rl.Vector2Add(ent.Acceleration, force)
}

func (ent *RobotEntity) Update() {

	// MOVEMENT
	ent.ApplyForce(ent.Arrive())
	ent.ApplyForce(ent.AvoidanceVelocity)

	ent.Velocity = rl.Vector2Add(ent.Velocity, ent.Acceleration)
	ent.Velocity = rl.Vector2ClampValue(ent.Velocity, 0, ent.MaximumSpeed)
	ent.SetPosition(rl.Vector2Add(ent.GetPosition(), rl.Vector2Scale(ent.Velocity, rl.GetFrameTime())))
	ent.Acceleration = rl.Vector2Zero()
	// ent.AvoidanceVelocity = rl.Vector2Zero()
}

func (ent *RobotEntity) Draw() {
	rl.DrawCircleV(ent.GetPosition(), 10, ent.Color)
	// draw the velocity to see where the robot is going
	rl.DrawLineV(ent.GetPosition(), rl.Vector2Add(ent.GetPosition(), ent.Velocity), rl.Black)
	// draw avoidance velocity
	// rl.DrawLineV(ent.GetPosition(), rl.Vector2Add(ent.GetPosition(), ent.AvoidanceVelocity), rl.Red)
	// draw vector to closest obstacle
	rl.DrawLineV(ent.GetPosition(), rl.Vector2Add(ent.GetPosition(), ent.VectorToObstacle), rl.Red)
	// draw AvoidanceVelocity
	rl.DrawLineV(ent.GetPosition(), rl.Vector2Add(ent.GetPosition(), ent.AvoidanceVelocity), rl.Red)
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
