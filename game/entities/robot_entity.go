package entities

import (
	"gorl/fw/core/entities"
	input "gorl/fw/core/input/input_event"
	"gorl/fw/util"
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Ensure that RobotEntity implements IEntity.
var _ entities.IEntity = &RobotEntity{}

// Robot Entity
type RobotEntity struct {
	*entities.Entity // Required!

	Velocity              rl.Vector2
	Acceleration          rl.Vector2
	CurrentTarget         rl.Vector2 // the next tile
	FinalTarget           rl.Vector2 // the target tile
	MaximumSpeed          float32
	MaximumForce          float32
	SlowDownDistance      float32
	WallDetectionRange    float32
	ClosestWall           rl.Vector2 // DEBUG
	WallAvoidanceVelocity rl.Vector2
	WallAvoidanceForce    float32

	RobotDetectionRange     float32
	RobotSeperationStrength float32
	RobotAlignmentStrength  float32
	RobotCohesionStrength   float32

	RobotSeperationVelocity rl.Vector2
	RobotAlignmentVelocity  rl.Vector2
	RobotCohesionVelocity   rl.Vector2

	SimSpeed float32

	HasCrashed bool
	Color      rl.Color
}

// NewRobotEntity creates a new instance of the RobotEntity.
func NewRobotEntity() *RobotEntity {
	new_ent := &RobotEntity{
		Entity:                entities.NewEntity("RobotEntity", rl.Vector2Zero(), 0, rl.Vector2One()),
		Velocity:              rl.Vector2Zero(),
		Acceleration:          rl.Vector2Zero(),
		CurrentTarget:         rl.Vector2Zero(),
		FinalTarget:           rl.Vector2Zero(),
		MaximumSpeed:          150,
		MaximumForce:          2.5,
		SlowDownDistance:      300,
		WallDetectionRange:    100,
		ClosestWall:           rl.Vector2Zero(),
		WallAvoidanceVelocity: rl.Vector2Zero(),
		WallAvoidanceForce:    1.0,

		RobotDetectionRange:     60,
		RobotSeperationStrength: 50,
		RobotAlignmentStrength:  20,
		RobotCohesionStrength:   0,

		RobotSeperationVelocity: rl.Vector2Zero(),
		RobotAlignmentVelocity:  rl.Vector2Zero(),
		RobotCohesionVelocity:   rl.Vector2Zero(),

		SimSpeed: 1,

		HasCrashed: false,
		Color:      rl.NewColor(uint8(rand.Int()%255), uint8(rand.Int()%255), uint8(rand.Int()%255), 255),
	}
	return new_ent
}

func (ent *RobotEntity) Init() {
}

func (ent *RobotEntity) Deinit() {
}

func (ent *RobotEntity) AvoidWall(obstacles []rl.Vector2) rl.Vector2 {
	force := rl.Vector2Zero()
	smallestDistance := math.MaxFloat64

	for _, obstacle := range obstacles {
		vecToObstacle := rl.Vector2Subtract(obstacle, ent.GetPosition())
		currentLength := float64(rl.Vector2Length(vecToObstacle))
		// if float32(closestObstacleLength) < 150 && currentLength < float32(closestObstacleLength) {
		if currentLength < smallestDistance {
			force = vecToObstacle
			smallestDistance = float64(currentLength)
		}
	}

	// for some reason this works, but not if the check is above in the if clause of the loop
	if smallestDistance > float64(ent.WallDetectionRange) {
		force = rl.Vector2Zero()
	}

	// for debug purposes
	ent.ClosestWall = force

	force = rl.Vector2Subtract(force, ent.WallAvoidanceVelocity)
	// limit the steering by the AvoidanceForce
	force = rl.Vector2ClampValue(force, float32(ent.WallAvoidanceForce), float32(ent.WallAvoidanceForce))

	// this should be the force pushing the robot away from the obstacle
	// so we have to multiply it with -1
	return rl.Vector2Scale(force, -1)
}

// ========================================================================================00
// JULIUS' FUNCTIONS
func (ent *RobotEntity) CalculateSeparationForce(nearbyRobots []*RobotEntity) rl.Vector2 {
	var force rl.Vector2
	for _, neighbor := range nearbyRobots {
		diff := rl.Vector2Subtract(ent.GetPosition(), neighbor.GetPosition())
		distance := rl.Vector2Length(diff)
		if distance < ent.RobotDetectionRange && distance > 0 { // separationThreshold is a defined constant
			pushForce := rl.Vector2Scale(util.Vector2NormalizeSafe(diff), ent.RobotSeperationStrength/distance)
			force = rl.Vector2Add(force, pushForce)
		}
	}
	return force
}

func (ent *RobotEntity) CalculateAlignmentForce(nearbyRobots []*RobotEntity) rl.Vector2 {
	var averageVelocity rl.Vector2
	var count int
	for _, neighbor := range nearbyRobots {
		averageVelocity = rl.Vector2Add(averageVelocity, neighbor.Velocity) // Assuming Velocity field exists
		count++
	}
	if count > 0 {
		averageVelocity = rl.Vector2Scale(averageVelocity, 1/float32(count))
		return rl.Vector2Scale(util.Vector2NormalizeSafe(averageVelocity), ent.RobotAlignmentStrength)
	}
	return rl.Vector2{}
}

func (ent *RobotEntity) CalculateCohesionForce(nearbyRobots []*RobotEntity) rl.Vector2 {
	var centerOfMass rl.Vector2
	var count int
	for _, neighbor := range nearbyRobots {
		centerOfMass = rl.Vector2Add(centerOfMass, neighbor.GetPosition())
		count++
	}
	if count > 0 {
		centerOfMass = rl.Vector2Scale(centerOfMass, 1/float32(count))
		toCenter := rl.Vector2Subtract(centerOfMass, ent.GetPosition())
		return rl.Vector2Scale(util.Vector2NormalizeSafe(toCenter), ent.RobotCohesionStrength)
	}
	return rl.Vector2{}
}

//==================================================================================================

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
	if !ent.HasCrashed {
		ent.ApplyForce(ent.Arrive())
		ent.ApplyForce(ent.WallAvoidanceVelocity)
		ent.ApplyForce(ent.RobotSeperationVelocity)
		ent.ApplyForce(ent.RobotAlignmentVelocity)

		ent.Velocity = rl.Vector2Add(ent.Velocity, ent.Acceleration)
		ent.Velocity = rl.Vector2ClampValue(ent.Velocity, 0, ent.MaximumSpeed)

		if ent.GetTilePosition() != ent.FinalTarget {

			scaledFrameTime := rl.GetFrameTime() * ent.SimSpeed
			ent.SetPosition(rl.Vector2Add(ent.GetPosition(), rl.Vector2Scale(ent.Velocity, scaledFrameTime)))
		}
	}

	ent.Acceleration = rl.Vector2Zero()
	// TODO: uncomment
	// ent.WallAvoidanceVelocity = rl.Vector2Zero()
	ent.RobotSeperationVelocity = rl.Vector2Zero()
	ent.RobotAlignmentVelocity = rl.Vector2Zero()
}

func (ent *RobotEntity) Draw() {
	rl.DrawCircleV(ent.GetPosition(), 10, ent.Color)
	// draw the velocity to see where the robot is going
	rl.DrawLineV(ent.GetPosition(), rl.Vector2Add(ent.GetPosition(), ent.Velocity), rl.Black)
	// rl.DrawLineV(ent.GetPosition(), rl.Vector2Add(ent.GetPosition(), ent.ClosestWall), rl.Red)
	rl.DrawLineV(ent.GetPosition(), rl.Vector2Add(ent.WallAvoidanceVelocity, ent.GetPosition()), rl.Green)
	// TODO: remove, this belongs in Update
	ent.WallAvoidanceVelocity = rl.Vector2Zero()

	// rl.DrawLineV(ent.GetPosition(), rl.Vector2Add(ent.GetPosition(), rl.Vector2Scale(ent.WallAvoidanceVelocity, 50)), rl.Red)
	// rl.DrawLineV(ent.GetPosition(), rl.Vector2Add(ent.GetPosition(), rl.Vector2Scale(ent.RobotSeperationVelocity, 50)), rl.Black)
	// draw wall detection range
	// rl.DrawCircleLines(int32(ent.GetPosition().X), int32(ent.GetPosition().Y), ent.WallDetectionRange, rl.Red)
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
