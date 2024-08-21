package scenes

import (
	"gorl/fw/core/gem"
	"gorl/fw/modules/scenes"
	"gorl/game/entities"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// This checks at compile time if the interface is implemented
var _ scenes.IScene = (*PheroSceneScene)(nil)

// PheroScene Scene
type PheroSceneScene struct {
	scenes.Scene // Required!

	// Custom Fields
	// Add fields here for any state that the scene should keep track of
	// ...
}

func (scn *PheroSceneScene) Init() {
	// Initialization logic for the scene
	// ...
	cam := entities.NewCameraEntity(
		rl.Vector2Zero(),
		rl.Vector2Zero(),
	)
	gem.Append(scn.GetRoot(), cam)

	foodpiles := entities.NewFoodpilesEntity()
	gem.Append(scn.GetRoot(), foodpiles)

	actor := entities.NewAntbotsEntity(1000, rl.NewVector2(1920/2, 1080/2), 40, cam, foodpiles)
	gem.Append(scn.GetRoot(), actor)
}

func (scn *PheroSceneScene) Deinit() {
	// De-initialization logic for the scene
	// ...
}
