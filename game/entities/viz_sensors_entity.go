package entities

import (
	"gorl/fw/core/entities"
	input "gorl/fw/core/input/input_event"
	"gorl/fw/core/settings"
	"gorl/fw/util"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Ensure that VizSensorsEntity implements IEntity.
var _ entities.IEntity = &VizSensorsEntity{}

// VizSensors Entity
type VizSensorsEntity struct {
	*entities.Entity // Required!

	// Custom Fields
	// Add fields here for any state that the entity should keep track of
	// ...
	antSprite          rl.Texture2D
	sensorPoints       []rl.Vector2
	sensorDistanceGoal []float32
	sensorDistances    []float32
	walkTarget         rl.Vector2

	size float32
}

// NewVizSensorsEntity creates a new instance of the VizSensorsEntity.
func NewVizSensorsEntity() *VizSensorsEntity {
	pos := rl.NewVector2(float32(settings.CurrentSettings().RenderWidth)/2, (float32(settings.CurrentSettings().RenderHeight) - 200))

	sensorPoints := []rl.Vector2{}
	sensorDistances := []float32{}
	sensorDistanceGoal := []float32{}
	visionAngle := float32(100.0)
	sensorCount := 3
	anglePerSensor := (visionAngle / float32(sensorCount-1)) * rl.Deg2rad
	centerOffset := (visionAngle / 2) * rl.Deg2rad
	for i := 0; i < sensorCount; i++ {
		sensorPoints = append(sensorPoints, rl.Vector2Rotate(rl.NewVector2(0, -1), anglePerSensor*float32(i)-centerOffset))
		sensorDistances = append(sensorDistances, 20)
		sensorDistanceGoal = append(sensorDistanceGoal, 20)
	}

	new_ent := &VizSensorsEntity{
		Entity: entities.NewEntity("VizSensorsEntity", pos, 0, rl.Vector2One()),

		antSprite:          rl.LoadTexture("antbot.png"),
		sensorPoints:       sensorPoints,
		sensorDistances:    sensorDistances,
		sensorDistanceGoal: sensorDistanceGoal,
		walkTarget:         rl.Vector2Zero(),

		size: 10,
	}
	return new_ent
}

func (ent *VizSensorsEntity) Init() {
	// Initialization logic for the entity
	// ...
}

func (ent *VizSensorsEntity) Deinit() {
	// De-initialization logic for the entity
	// ...
}

func (ent *VizSensorsEntity) Update() {
	for idx, sd := range ent.sensorDistances {
		goal := ent.sensorDistanceGoal[idx]
		diff := goal - sd
		if diff > 0.1 || diff < -0.1 {
			ent.sensorDistances[idx] += diff * rl.GetFrameTime()
		} else {
			ent.sensorDistances[idx] = goal
		}
	}
}

func (ent *VizSensorsEntity) Draw() {
	// Draw logic for the entity
	// ...

	ent.walkTarget = rl.Vector2Zero()
	for idx, sensorPoint := range ent.sensorPoints {
		scaledSensor := rl.Vector2Scale(sensorPoint, ent.sensorDistances[idx]*ent.size)
		absPoint := rl.Vector2Add(ent.GetPosition(), scaledSensor)
		//rl.DrawLineV(ent.GetPosition(), absPoint, rl.Red)
		rl.DrawCircleV(absPoint, 5*ent.size, rl.Red)
		ent.walkTarget = rl.Vector2Add(ent.walkTarget, scaledSensor)
	}
	ent.walkTarget = util.Vector2NormalizeSafe(ent.walkTarget)
	ent.walkTarget = rl.Vector2Scale(ent.walkTarget, 50*ent.size)
	ent.walkTarget = rl.Vector2Add(ent.GetPosition(), ent.walkTarget)
	rl.DrawTexturePro(
		ent.antSprite,
		rl.NewRectangle(0, 0, float32(ent.antSprite.Width), float32(ent.antSprite.Height)),
		rl.NewRectangle(ent.GetPosition().X, ent.GetPosition().Y, float32(ent.antSprite.Width)*ent.size, float32(ent.antSprite.Height)*ent.size),
		rl.NewVector2((float32(ent.antSprite.Width)/2)*ent.size, (float32(ent.antSprite.Height)/2)*ent.size),
		0, rl.Black,
	)

	// draw walk target line
	//rl.DrawLineV(ent.GetPosition(), ent.walkTarget, rl.Blue)
	rl.DrawLineEx(ent.GetPosition(), ent.walkTarget, 2, rl.Lime)
	rl.DrawCircleV(ent.walkTarget, 1*ent.size, rl.Lime)
}

func (ent *VizSensorsEntity) OnInputEvent(event *input.InputEvent) bool {
	// Logic to run when an input event is received.
	// Return false if the event was consumed and should not be propagated
	// further.

	if event.Action == input.ActionClickDown || event.Action == input.ActionRightClickDown {
		fac := float32(1)
		if event.Action == input.ActionRightClickDown {
			fac = -1
		}
		mp := event.GetScreenSpaceMousePosition()
		for idx, sensorPoint := range ent.sensorPoints {
			scaledSensor := rl.Vector2Scale(sensorPoint, ent.sensorDistances[idx]*ent.size)
			absPoint := rl.Vector2Add(ent.GetPosition(), scaledSensor)
			if rl.CheckCollisionPointCircle(mp, absPoint, 5*ent.size) {
				ent.sensorDistanceGoal[idx] += 10 * fac
				return false
			}
		}

	}

	return true
}
