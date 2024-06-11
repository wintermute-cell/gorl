package entities

import (
	"gorl/fw/core/entities"
	input "gorl/fw/core/input/input_event"
	"gorl/fw/core/logging"
	"gorl/fw/core/store"
	"gorl/fw/physics"
	"gorl/fw/util"
	"gorl/fw/util/easing"
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
	visionRange   float32
	goal          rl.Vector2
	mainCamera    *CameraEntity

	selfColor color.RGBA
	rayColor  color.RGBA

	constructedMapData []color.RGBA
	constructedMap     rl.Texture2D

	// Showcase trigger flags
	hasRays      bool
	rayFade      float32
	isRotating   bool
	isDrawingMap bool
	isHalting    bool
}

// NewVfhActorEntity creates a new instance of the VfhActorEntity.
func NewVfhActorEntity(position rl.Vector2, rayAmount int32, visionAngle, visionRange float32) *VfhActorEntity {
	// NOTE: you can modify the constructor to take any parameters you need to
	// initialize the entity.
	new_ent := &VfhActorEntity{
		sprite:      rl.LoadTexture("robot_small.png"),
		Entity:      entities.NewEntity("VfhActorEntity", position, 0, rl.Vector2One()),
		forward:     rl.NewVector2(1, 0),
		visionRange: visionRange,
		selfColor:   colorscheme.Colorscheme.Color01.ToRGBA(),
		rayColor:    colorscheme.Colorscheme.Color06.ToRGBA(),
	}

	// fill the image with transparent black pixels
	imgData := make([]byte, 1920*1080*4)
	emptyImg := rl.NewImage(imgData, 1920, 1080, 1, rl.UncompressedR8g8b8a8)
	new_ent.constructedMap = rl.LoadTextureFromImage(emptyImg)
	new_ent.constructedMapData = make([]color.RGBA, 1920*1080)

	new_ent.rayColor.A = 0 // We fade this in

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

	const rayFadeSpeed = 0.5
	if ent.hasRays {
		ent.rayFade = util.Min(ent.rayFade+rl.GetFrameTime(), 0.999)
		ent.rayColor.A = util.Min(255, uint8(easing.QuadOut(ent.rayFade, 1, 255, 1)))
	}

	if ent.isRotating && (!ent.isHalting) {
		ent.SetRotation(ent.GetRotation() + rl.GetFrameTime()*20)
	}

	// make sure we have the camera
	if ent.mainCamera == nil {
		mainCamera, ok := store.Get[*CameraEntity]()
		if ok {
			ent.mainCamera = mainCamera
		}
	}

	for idx, rayDir := range ent.rayDirections {
		rayDir = rl.Vector2Rotate(rayDir, ent.GetRotation()*rl.Deg2rad)
		hit := physics.Raycast(ent.GetPosition(), rayDir, ent.visionRange, physics.CollisionCategoryAll)
		if len(hit) > 0 {
			ent.rayHits[idx] = hit[0]

			// Update the constructed map
			if ent.isDrawingMap {
				hitPos := hit[0].IntersectionPoint
				hitPos.X = util.Clamp(hitPos.X, 0, 1919)
				hitPos.Y = util.Clamp(hitPos.Y, 0, 1079)
				idx := int(hitPos.Y)*1920 + int(hitPos.X)
				ent.constructedMapData[idx] = colorscheme.Colorscheme.Color04.ToRGBA()
			}
		} else {
			ent.rayHits[idx] = physics.RaycastHit{}
		}
	}

	// Make the texture match the map data
	rl.UpdateTexture(ent.constructedMap, ent.constructedMapData)
}

func (ent *VfhActorEntity) Draw() {

	// Draw the constructed map
	rl.DrawTexture(ent.constructedMap, 0, 0, rl.White)

	// Draw the rays
	for idx, rayDir := range ent.rayDirections {
		rayDir = rl.Vector2Rotate(rayDir, ent.GetRotation()*rl.Deg2rad)
		distToHit := float32(math.MaxFloat32)
		if ent.rayHits[idx] != (physics.RaycastHit{}) {
			distToHit = rl.Vector2Distance(ent.GetPosition(), ent.rayHits[idx].IntersectionPoint)
		}
		scaledRay := rl.Vector2Scale(rayDir, util.Min(distToHit, ent.visionRange))
		rl.DrawLineV(
			ent.GetPosition(),
			rl.Vector2Add(ent.GetPosition(), scaledRay),
			ent.rayColor,
		)
	}

	// Draw the intersection points
	for _, hit := range ent.rayHits {
		if hit != (physics.RaycastHit{}) {
			rl.DrawCircleV(hit.IntersectionPoint, 2, ent.rayColor)
		}
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

	// Draw constructed map
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
		if !ent.hasRays {
			ent.hasRays = true
		} else if !ent.isRotating {
			ent.isRotating = true
		} else if !ent.isDrawingMap {
			ent.isDrawingMap = true
		} else if !ent.isHalting {
			ent.isHalting = true
		}
	}

	return true
}
