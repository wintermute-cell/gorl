package scenes

import (
	"gorl/fw/core/gem"
	"gorl/fw/core/math"
	"gorl/fw/core/settings"
	"gorl/fw/modules/scenes"
	"gorl/game/entities"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// This checks at compile time if the interface is implemented
var _ scenes.IScene = (*FlowFieldScene)(nil)

// FlowField Scene
type FlowFieldScene struct {
	scenes.Scene // Required!

	// Custom Fields
	// Add fields here for any state that the scene should keep track of
	// ...
}

func (scn *FlowFieldScene) Init() {

	camera := entities.NewCameraEntity(
		rl.Vector2Zero(),
		rl.Vector2Zero(),
		rl.NewVector2(float32(settings.CurrentSettings().RenderWidth), float32(settings.CurrentSettings().RenderHeight)),
		rl.Vector2Zero(),
		math.Flag0,
	)
	gem.Append(scn.GetRoot(), camera)

	gridGraph := entities.NewGridGraphEntity()
	gem.Append(scn.GetRoot(), gridGraph)
	// add some test robots
	for i := range 10 {
		for j := range 10 {
			testRobot := entities.NewRobotEntity()
			testRobot.SetPosition(rl.NewVector2(float32(i+10)*40, float32(j)*40))
			testRobot.Color = rl.NewColor(uint8(200-rand.Int()%100), uint8(255-rand.Int()%150), 0, 255)
			gem.Append(gridGraph, testRobot)
		}
	}

}

func (scn *FlowFieldScene) Deinit() {
	// De-initialization logic for the scene
	// ...
}
