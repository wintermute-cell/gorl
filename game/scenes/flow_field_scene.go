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
		rl.NewVector2(
			float32(settings.CurrentSettings().RenderWidth),
			float32(settings.CurrentSettings().RenderHeight),
		),
		rl.Vector2Zero(),
		math.Flag0,
	)

	gem.Append(scn.GetRoot(), camera)

	gridGraph := entities.NewGridGraphEntity()
	gem.Append(scn.GetRoot(), gridGraph)
	// add some test robots
	testRobot1 := entities.NewRobotEntity()
	testRobot1.SetPosition(rl.NewVector2(float32(10)*40+20, float32(10)*40+20))
	testRobot2 := entities.NewRobotEntity()
	testRobot2.SetPosition(rl.NewVector2(float32(20)*40+20, float32(10)*40+20))
	testRobot3 := entities.NewRobotEntity()
	testRobot3.SetPosition(rl.NewVector2(float32(30)*40+20, float32(15)*40+20))
	testRobot4 := entities.NewRobotEntity()
	testRobot4.SetPosition(rl.NewVector2(float32(40)*40+20, float32(23)*40+20))
	testrobot5 := entities.NewRobotEntity()
	testrobot5.SetPosition(rl.NewVector2(float32(11)*40+20, float32(10)*40+20))

	gem.Append(gridGraph, testRobot1)
	// gem.Append(gridGraph, testRobot2)
	// gem.Append(gridGraph, testRobot3)
	// gem.Append(gridGraph, testRobot4)
	gem.Append(gridGraph, testrobot5)

	gridGraph.InitRobots()
	//
	// counter := 0
	// for k := range gridGraph.Gg.VertexMap {
	// 	counter++
	// 	if counter%10 == 0 {
	// 		testRobot := entities.NewRobotEntity()
	// 		testRobot.SetPosition(rl.Vector2Add(rl.Vector2Scale(k, 40), rl.NewVector2(20, 20)))
	// 		gem.Append(gridGraph, testRobot)
	// 	}
	// }

}

func (scn *FlowFieldScene) Deinit() {
	// De-initialization logic for the scene
	// ...
}
