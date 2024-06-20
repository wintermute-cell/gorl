package scenes

import (
	"gorl/fw/core/gem"
	"gorl/fw/core/math"
	"gorl/fw/core/settings"
	"gorl/fw/core/store"

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
	cameraEntity := entities.NewCameraEntityEx(
		rl.Vector2Zero(),
		rl.Vector2Zero(),
		rl.NewVector2(
			float32(settings.CurrentSettings().ScreenWidth),
			float32(settings.CurrentSettings().ScreenHeight),
		),
		rl.Vector2Zero(),
		math.Flag0,
	)
	gem.Append(scn.GetRoot(), cameraEntity)
	store.Add(cameraEntity)

	//actorInfo := entities.NewActorInfoEntity()
	//gem.Append(cameraEntity, actorInfo)
	//store.Add(actorInfo)

	baseStation := entities.NewBaseStationEntity(rl.NewVector2(430, 650))

	actor := entities.NewVfhActorEntity(baseStation, rl.NewVector2(675, 840), 5, 120, 80)
	gem.Append(scn.GetRoot(), actor)

	//actor2 := entities.NewVfhActorEntity(baseStation, rl.NewVector2(1400, 820), 13, 120, 80)
	//gem.Append(scn.GetRoot(), actor2)

	//actor3 := entities.NewVfhActorEntity(baseStation, rl.NewVector2(1300, 620), 13, 120, 80)
	//gem.Append(scn.GetRoot(), actor3)

	env := entities.NewEnvironmentEntity()
	gem.Append(scn.GetRoot(), env)
	gem.Append(scn.GetRoot(), baseStation)
}

func (scn *VfhScene) Deinit() {
	// De-initialization logic for the scene
	// ...
}
