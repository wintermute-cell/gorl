package scenes

import (
	"gorl/fw/core/gem"
	"gorl/fw/modules/scenes"
	"gorl/game/entities"
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

	// for animation of the grid graph creation
	// mapEnt := entities.NewAnimMapToFfEntity()
	// gem.Append(scn.GetRoot(), mapEnt)

	testRobot := entities.NewRobotEntity()
	gem.Append(scn.GetRoot(), testRobot)

	gridGraph := entities.NewGridGraphEntity()
	gem.Append(scn.GetRoot(), gridGraph)

	// TODO: robots die direction geben und so
}

func (scn *FlowFieldScene) Deinit() {
	// De-initialization logic for the scene
	// ...
}
