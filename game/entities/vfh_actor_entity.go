package entities

import (
	"gorl/fw/core/entities"
	input "gorl/fw/core/input/input_event"
	"gorl/fw/core/logging"
	"gorl/fw/core/store"
	"gorl/fw/physics"
	"gorl/fw/util"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Ensure that VfhActorEntity implements IEntity.
var _ entities.IEntity = &VfhActorEntity{}

// VfhActor Entity
type VfhActorEntity struct {
	*entities.Entity // Required!

	// Custom Fields
	forward       rl.Vector2
	rayDirections []rl.Vector2
	rayHits       []physics.RaycastHit
	visionRange   float32
	goal          rl.Vector2
	mainCamera    *CameraEntity
}

// NewVfhActorEntity creates a new instance of the VfhActorEntity.
func NewVfhActorEntity(position rl.Vector2, rayAmount int32, visionAngle, visionRange float32) *VfhActorEntity {
	// NOTE: you can modify the constructor to take any parameters you need to
	// initialize the entity.
	new_ent := &VfhActorEntity{
		Entity:      entities.NewEntity("VfhActorEntity", position, 0, rl.Vector2One()),
		forward:     rl.NewVector2(1, 0),
		visionRange: visionRange,
	}

	anglePerRay := visionAngle / float32(rayAmount-1)
	new_ent.rayDirections = make([]rl.Vector2, rayAmount)
	for i := int32(0); i < rayAmount; i++ {
		centerOffset := (visionAngle / 2) * rl.Deg2rad
		new_ent.rayDirections[i] = rl.Vector2Rotate(rl.NewVector2(1, 0), (anglePerRay*rl.Deg2rad)*float32(i)-centerOffset)
	}
	new_ent.rayHits = make([]physics.RaycastHit, rayAmount, rayAmount)

	return new_ent
}

// SetGoal sets the goal for the actor to move towards.
func (ent *VfhActorEntity) SetGoal(goal rl.Vector2) {
	ent.goal = goal
}

func (ent *VfhActorEntity) Init() {
	// Initialization logic for the entity
	// ...
	mainCamera, ok := store.Get[*CameraEntity]()
	if ok {
		ent.mainCamera = mainCamera
	}
}

func (ent *VfhActorEntity) Deinit() {
	// De-initialization logic for the entity
	// ...
}

func (ent *VfhActorEntity) Update() {
	// Update logic for the entity per frame
	// ...

	// TODO: remove this
	ent.SetRotation(ent.GetRotation() + rl.GetFrameTime()*20)

	// make sure we have the camera
	if ent.mainCamera == nil {
		mainCamera, ok := store.Get[*CameraEntity]()
		if ok {
			ent.mainCamera = mainCamera
		}
	}

	//hit := physics.Raycast(ent.GetPosition(), rl.Vector2{X: -1, Y: 0}, 400, physics.CollisionCategoryAll)
	//if len(hit) > 0 {
	//	ent.hitPoint = hit[0].IntersectionPoint
	//} else {
	//	ent.hitPoint = rl.Vector2{}
	//}

	for idx, rayDir := range ent.rayDirections {
		rayDir = rl.Vector2Rotate(rayDir, ent.GetRotation()*rl.Deg2rad)
		hit := physics.Raycast(ent.GetPosition(), rayDir, ent.visionRange, physics.CollisionCategoryAll)
		if len(hit) > 0 {
			ent.rayHits[idx] = hit[0]
		} else {
			ent.rayHits[idx] = physics.RaycastHit{}
		}
	}
}

func (ent *VfhActorEntity) Draw() {
	// Draw logic for the entity
	// ...

	rl.DrawCircleV(ent.GetPosition(), 10, rl.Red)
	rl.DrawCircleV( // a little head to mark forward
		rl.Vector2Add(
			ent.GetPosition(),
			rl.Vector2Rotate(rl.NewVector2(10, 0), ent.GetRotation()*rl.Deg2rad),
		), 6, rl.Black)
	for idx, rayDir := range ent.rayDirections {
		rayDir = rl.Vector2Rotate(rayDir, ent.GetRotation()*rl.Deg2rad)
		distToHit := float32(math.MaxFloat32)
		if ent.rayHits[idx] != (physics.RaycastHit{}) {
			distToHit = rl.Vector2Distance(ent.GetPosition(), ent.rayHits[idx].IntersectionPoint)
		}
		scaledRay := rl.Vector2Scale(rayDir, util.Min(distToHit, ent.visionRange))
		rl.DrawLineV(ent.GetPosition(), rl.Vector2Add(ent.GetPosition(), scaledRay), rl.Green)
	}
	for _, hit := range ent.rayHits {
		if hit != (physics.RaycastHit{}) {
			rl.DrawCircleV(hit.IntersectionPoint, 5, rl.Blue)
		}
	}
}

func (ent *VfhActorEntity) OnInputEvent(event *input.InputEvent) bool {
	// Logic to run when an input event is received.
	// Return false if the event was consumed and should not be propagated
	// further.
	// ...

	worldCursor := ent.mainCamera.ScreenToWorld(event.GetScreenSpaceMousePosition())

	if event.Action == input.ActionClickDown && rl.CheckCollisionPointCircle(worldCursor, ent.GetPosition(), 10) {
		logging.Debug("Actor clicked!")
	}

	return true
}
