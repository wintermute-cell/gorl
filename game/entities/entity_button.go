package entities

import (
	"gorl/fw/core/entities/proto"
	input "gorl/fw/core/input/input_event"
	"gorl/fw/core/logging"
	"gorl/fw/core/render"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Button Entity
type ButtonEntity2D struct {
	// Required fields
	*proto.Entity2D

	// Custom Fields
	// Add fields here for any state that the entity should keep track of
	// ...
}

func NewButtonEntity2D(position rl.Vector2, rotation float32, scale rl.Vector2) *ButtonEntity2D {
	new_ent := &ButtonEntity2D{
		Entity2D: proto.NewEntity2D(position, rotation, scale),

		// Initialize custom fields here
		// ...
	}
	return new_ent
}

func (ent *ButtonEntity2D) Init() {
	// Initialization logic for the entity
	// ...
}

func (ent *ButtonEntity2D) Deinit() {
	// De-initialization logic for the entity
	// ...
}

func (ent *ButtonEntity2D) Update() {
	// Update logic for the entity
	// ...
}

func (ent *ButtonEntity2D) Draw() {
	//rl.DrawRectangleV(ent.GetPosition(), rl.NewVector2(100, 100), rl.Red)
	//rl.DrawText(ent.GetName(), int32(ent.GetPosition().X), int32(ent.GetPosition().Y), 20, rl.White)
	rotation := ent.GetRotation()
	rl.DrawRectanglePro(rl.NewRectangle(ent.GetPosition().X, ent.GetPosition().Y, 100*ent.GetScale().X, 100*ent.GetScale().Y), rl.NewVector2(0, 0), rotation, rl.Red)
	rl.DrawTextEx(rl.GetFontDefault(), ent.GetName(), rl.NewVector2(ent.GetPosition().X, ent.GetPosition().Y), 20, 0, rl.Green)
	rl.DrawRectangleLinesEx(rl.NewRectangle(ent.GetPosition().X, ent.GetPosition().Y, 100*ent.GetScale().X, 100*ent.GetScale().Y), 1, rl.Green)
}

func (ent *ButtonEntity2D) OnInputEvent(event *input.InputEvent) bool {
	mousePos := event.GetScreenSpaceMousePosition()
	mousePosCamSpace := render.PointToCameraSpace(, mousePos)
	if event.Action == input.ActionClickDown {
		hit := rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(ent.GetPosition().X, ent.GetPosition().Y, 100*ent.GetScale().X, 100*ent.GetScale().Y))
		if hit {
			logging.Debug("Button hit! %v", ent.GetName())
			render.SetCameraPosition(ent.GetPosition())
			return false
		}
	}

	return true
}
