package scenes

import (
	"gorl/fw/core/gem"
	"gorl/fw/core/math"
	"gorl/fw/core/settings"
	"gorl/fw/modules/scenes"
	"gorl/game/entities"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// This checks at compile time if the interface is implemented
var _ scenes.IScene = (*VfhScene)(nil)

// Vfh Scene
type VfhScene struct {
	scenes.Scene // Required!

	// Custom Fields
	// Add fields here for any state that the scene should keep track of
	// ...
}

func (scn *VfhScene) Init() {
	// Initialization logic for the scene
	// ...
	cameraEntity := entities.NewCameraEntity(
		rl.Vector2Zero(),
		rl.Vector2Zero(),
		rl.NewVector2(
			float32(settings.CurrentSettings().ScreenWidth),
			float32(settings.CurrentSettings().ScreenHeight),
		),
		rl.Vector2Zero(),
		math.Flag0,
	)
	gem.Append(gem.GetRoot(), cameraEntity)

	actor := entities.NewVfhActorEntity(rl.NewVector2(100, 100))
	gem.Append(scn.GetRoot(), actor)

	env := entities.NewEnvironmentEntity()
	gem.Append(scn.GetRoot(), env)

}

func (scn *VfhScene) Deinit() {
	// De-initialization logic for the scene
	// ...
}
