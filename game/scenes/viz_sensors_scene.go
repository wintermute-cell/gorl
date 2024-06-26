package scenes

import (
	"gorl/fw/core/gem"
	"gorl/fw/modules/scenes"
	"gorl/game/entities"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// This checks at compile time if the interface is implemented
var _ scenes.IScene = (*VizSensorsScene)(nil)

// VizSensors Scene
type VizSensorsScene struct {
	scenes.Scene // Required!

	// Custom Fields
	// Add fields here for any state that the scene should keep track of
	// ...
}

func (scn *VizSensorsScene) Init() {

	cam := entities.NewCameraEntity(
		rl.Vector2Zero(),
		rl.Vector2Zero(),
	)
	gem.Append(scn.GetRoot(), cam)

	vizEntity := entities.NewVizSensorsEntity()
	gem.Append(scn.GetRoot(), vizEntity)
}

func (scn *VizSensorsScene) Deinit() {
	// De-initialization logic for the scene
	// ...
}
