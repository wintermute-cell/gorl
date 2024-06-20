package entities

import (
	"gorl/fw/core/entities"
	input "gorl/fw/core/input/input_event"
	"gorl/fw/core/logging"
	"gorl/fw/core/store"
	"gorl/fw/physics"
	"gorl/fw/util"
	"gorl/game/code/colorscheme"
	"image/color"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Ensure that VfhActorEntity implements IEntity.
var _ entities.IEntity = &VfhActorEntity{}

// VfhActor Entity
type VfhActorEntity struct {
	*entities.Entity // Required!

	// Custom Fields
	sprite        rl.Texture2D
	forward       rl.Vector2
	rayDirections []rl.Vector2
	rayHits       []physics.RaycastHit
	vfHistogram   []rl.Vector2
	adjustedVFH   []rl.Vector2
	visionRange   float32
	mainCamera    *CameraEntity
	collider      *physics.Collider
	isCrashed     bool
	optimalDir    rl.Vector2

	selfColor     color.RGBA
	rayColor      color.RGBA
	adjustedColor color.RGBA

	baseStation   *BaseStationEntity
	goalPath      []rl.Vector2
	goalDirection rl.Vector2

	// Animation flags
	hasPathfinding bool
}

// NewVfhActorEntity creates a new instance of the VfhActorEntity.
func NewVfhActorEntity(
	baseStation *BaseStationEntity,
	position rl.Vector2,
	rayAmount int32,
	visionAngle, visionRange float32,
) *VfhActorEntity {
	if rayAmount%2 == 0 {
		logging.Warning("Ray amount should be an odd number, so there is ray pointing straight forward! Adding 1 to the ray amount.")
		rayAmount = rayAmount + 1
	}
	if rayAmount < 3 {
		logging.Warning("Ray amount should be at least 3, setting it to 3.")
		rayAmount = 3
	}
	new_ent := &VfhActorEntity{
		sprite:        rl.LoadTexture("robot_small.png"),
		Entity:        entities.NewEntity("VfhActorEntity", position, 0, rl.Vector2One()),
		baseStation:   baseStation,
		forward:       rl.NewVector2(1, 0),
		visionRange:   visionRange,
		selfColor:     colorscheme.Colorscheme.Color01.ToRGBA(),
		rayColor:      colorscheme.Colorscheme.Color06.ToRGBA(),
		adjustedColor: colorscheme.Colorscheme.Color07.ToRGBA(),
		collider:      physics.NewCircleCollider(position, 6, physics.BodyTypeDynamic),
	}

	new_ent.collider.SetCategory(physics.CollisionCategoryPlayer)
	new_ent.collider.SetCallbacks(
		map[physics.CollisionCategory]physics.CollisionCallback{
			physics.CollisionCategoryEnvironment: func() {
				logging.Info("ACTOR CRASHED")
				new_ent.isCrashed = true
				new_ent.selfColor = colorscheme.Colorscheme.Color10.ToRGBA()
			},
		},
	)

	anglePerRay := visionAngle / float32(rayAmount-1)
	new_ent.rayDirections = make([]rl.Vector2, rayAmount)
	for i := int32(0); i < rayAmount; i++ {
		centerOffset := (visionAngle / 2) * rl.Deg2rad
		new_ent.rayDirections[i] = rl.Vector2Rotate(rl.NewVector2(1, 0), (anglePerRay*rl.Deg2rad)*float32(i)-centerOffset)
	}
	new_ent.rayHits = make([]physics.RaycastHit, rayAmount, rayAmount)
	new_ent.vfHistogram = make([]rl.Vector2, rayAmount, rayAmount)
	new_ent.adjustedVFH = make([]rl.Vector2, rayAmount, rayAmount)

	return new_ent
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

	if ent.isCrashed {
		return
	}

	ent.forward = rl.Vector2Rotate(rl.NewVector2(1, 0), ent.GetRotation()*rl.Deg2rad)

	const rayFadeSpeed = 0.5

	// make sure we have the camera
	if ent.mainCamera == nil {
		mainCamera, ok := store.Get[*CameraEntity]()
		if ok {
			ent.mainCamera = mainCamera
		}
	}

	for idx, rayDir := range ent.rayDirections {
		rayDir = rl.Vector2Rotate(rayDir, ent.GetRotation()*rl.Deg2rad)
		hits := physics.Raycast(ent.GetPosition(), rayDir, ent.visionRange, physics.CollisionCategoryEnvironment)
		if len(hits) > 0 {
			ent.rayHits[idx] = hits[0]
			distToIntersection := rl.Vector2Distance(ent.GetPosition(), hits[0].IntersectionPoint)
			ent.vfHistogram[idx] = rl.Vector2Scale(rayDir, distToIntersection)

			// Update the constructed map
			hitPos := hits[0].IntersectionPoint
			hitPos.X = util.Clamp(hitPos.X, 0, 1919)
			hitPos.Y = util.Clamp(hitPos.Y, 0, 1079)
			ent.baseStation.DiscoverPixel(int(hitPos.X), int(hitPos.Y))
		} else {
			ent.rayHits[idx] = physics.RaycastHit{}
			ent.vfHistogram[idx] = rl.Vector2Scale(rayDir, ent.visionRange)
		}
	}

	// Steering logic
	ent.goalDirection = rl.Vector2Zero()
	ent.goalPath = ent.baseStation.GetPath(ent.GetPosition(), ent.baseStation.GetPosition())
	if len(ent.goalPath) >= 3 {
		ent.goalDirection = util.Vector2NormalizeSafe(rl.Vector2Subtract(ent.goalPath[1], ent.GetPosition()))
		//rl.Vector2Normalize(rl.Vector2Subtract(ent.goalPath[2], ent.goalPath[0]))
	}
	ent.optimalDir = ent.VFHGoalOrientedCostFunction(ent.vfHistogram, ent.goalDirection)
	if rl.Vector2LengthSqr(ent.optimalDir) <= 0.000001 {
		ent.SetRotation(ent.GetRotation() + 1)
		// No direction to move towards, only rotate
		return
	}
	//angleToOptimal := rl.Vector2Angle(ent.forward, ent.optimalDir)
	angleToOptimal := util.Vector2AngleSmallest(ent.forward, ent.optimalDir)
	logging.Debug("Angle to optimal: %v", angleToOptimal*rl.Rad2deg)
	ent.SetRotation(ent.GetRotation() + angleToOptimal*rl.Rad2deg*(rl.GetFrameTime()*1))

	// Move forward every frame
	moveSpeed := float32(20)
	curPos := ent.GetPosition()
	moveDelta := rl.Vector2Scale(ent.forward, moveSpeed*rl.GetFrameTime())
	newPos := rl.Vector2Add(curPos, moveDelta)
	ent.SetPosition(newPos)
	ent.collider.SetPosition(newPos)
}

func (ent *VfhActorEntity) Draw() {

	// Draw the rays
	for _, rayDir := range ent.vfHistogram {
		rl.DrawLineV(
			ent.GetPosition(),
			rl.Vector2Add(ent.GetPosition(), rayDir),
			ent.rayColor,
		)
	}

	// Draw the adjusted histogram
	for _, rayDir := range ent.adjustedVFH {
		rl.DrawLineV(
			ent.GetPosition(),
			rl.Vector2Add(ent.GetPosition(), rayDir),
			ent.adjustedColor,
		)
	}

	// Draw the currently selected direction
	rl.DrawLineV(
		ent.GetPosition(),
		rl.Vector2Add(ent.GetPosition(), ent.optimalDir),
		rl.Red,
	)

	// Draw the direction towards the goal
	rl.DrawLineV(
		ent.GetPosition(),
		rl.Vector2Add(ent.GetPosition(), rl.Vector2Scale(ent.goalDirection, 50)),
		rl.Green,
	)

	// Draw the intersection points
	for _, hit := range ent.rayHits {
		if hit != (physics.RaycastHit{}) {
			rl.DrawCircleV(hit.IntersectionPoint, 2, ent.rayColor)
		}
	}

	// Draw the current path to the goal
	for idx, pos := range ent.goalPath {
		if idx <= 0 { // we skip the first line segment
			continue
		}
		rl.DrawLineV(ent.goalPath[idx-1], pos, rl.Green)
	}

	// Drawing the actor itself
	rl.DrawCircleV(ent.GetPosition(), 6, ent.selfColor)
	//rl.DrawTexturePro(
	//	ent.sprite,
	//	rl.NewRectangle(0, 0, float32(ent.sprite.Width), float32(ent.sprite.Height)),
	//	rl.NewRectangle(
	//		ent.GetPosition().X,
	//		ent.GetPosition().Y,
	//		float32(ent.sprite.Width)*2,
	//		float32(ent.sprite.Height)*2,
	//	),
	//	rl.NewVector2(float32(ent.sprite.Width), float32(ent.sprite.Height)),
	//	0, //ent.GetRotation(),
	//	rl.White,
	//)
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

	if event.Action == input.ActionNextAnimation {
		if !ent.hasPathfinding {
			ent.hasPathfinding = true
		}
	}

	return true
}

// VFHGoalOrientedCostFunction selects the most optimal direction to move towards, given a
// polar histogram of the environment. It uses direction centrality and alignment with the goal direction as weights.
// Might return a zero vector if no direction is suitable.
func (ent *VfhActorEntity) VFHGoalOrientedCostFunction(vfh []rl.Vector2, goalDirection rl.Vector2) rl.Vector2 {

	vfh = ent.EnlargementFunction(vfh, 10)
	ent.adjustedVFH = vfh

	mid := len(vfh) / 2
	selection := rl.Vector2Zero()
	selectionCost := float32(math.MaxFloat32)

	// delta returns the absolute angle between two vectors
	//delta := func(v1, v2 rl.Vector2) float32 {
	//	angleRad := rl.Vector2Angle(v1, v2)
	//	absRad := util.Abs(angleRad)
	//	return absRad
	//}
	delta := func(v1, v2 rl.Vector2) float32 {
		v := util.Vector2AngleSmallest(v1, v2)
		return util.Abs(v)
	}

	for i := mid; i >= 0; i-- {
		leftIdx := mid - i
		rightIdx := mid + i

		goalAngleLeft := delta(goalDirection, vfh[leftIdx])
		goalAngleRight := delta(goalDirection, vfh[rightIdx])

		left := vfh[leftIdx]
		leftLen := rl.Vector2Length(left)
		right := vfh[rightIdx]
		rightLen := rl.Vector2Length(right)

		if leftLen >= 0.000001 {
			if goalAngleLeft < selectionCost {
				selection = left
				selectionCost = goalAngleLeft
			}
		}
		if rightLen >= 0.000001 {
			if goalAngleRight < selectionCost {
				selection = right
				selectionCost = goalAngleRight
			}
		}
	}

	return selection
}

// The EnlargementFunction filters the histogram according to the VFH+ paper.
func (ent *VfhActorEntity) EnlargementFunction(vfh []rl.Vector2, robotRadius int32) []rl.Vector2 {

	// 1. Gather all detected obstacle points
	obstaclePoints := make([]rl.Vector2, 0)
	for idx := range vfh {
		if ent.rayHits[idx] != (physics.RaycastHit{}) { // This ray hit something
			obstaclePoints = append(obstaclePoints, ent.rayHits[idx].IntersectionPoint)
		}
	}

	// 2. Calculate an angular range `gamma` around each obstacle point
	// Calculated as in the paper: http://www.cs.cmu.edu/~iwan/papers/vfh+.pdf (page 2)
	gammas := make([]float32, 0)
	clearance := float64(robotRadius) * 1.0 // This is the minimum clearance between the robot and the obstacle, d_s in the paper
	for idx := range obstaclePoints {
		gammaRad := math.Asin(
			util.Clamp(float64(float64(robotRadius)+clearance)/float64(rl.Vector2Length(vfh[idx])), 0, 1),
		)
		gammas = append(gammas, float32(gammaRad))
	}

	for _, obs := range obstaclePoints {
		rl.DrawCircleV(obs, float32(robotRadius+int32(clearance)), rl.Red)
	}

	// 3. Iterate over vfh and set vfh[i]=0 the angle between vfh[i] and any
	// vector pointing to an obstacle point is \leq than the gamma of that
	// obstacle point.
vhfLoop:
	for idx := range vfh {
		for i := range obstaclePoints {
			// Calculate the angle between the obstacle point and the vfh vector
			obstacleDir := rl.Vector2Subtract(obstaclePoints[i], ent.GetPosition())
			angleRad := util.Abs(util.Vector2AngleSmallest(vfh[idx], obstacleDir))
			if angleRad <= gammas[i] {
				vfh[idx] = rl.Vector2Zero()
				continue vhfLoop
			} else {
				// Here we just leave the vfh vector as it is.
				// In the paper this is represented with a multiplication by 1.
			}
		}
	}

	return vfh
}
