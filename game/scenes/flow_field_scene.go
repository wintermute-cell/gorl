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
	// Initialization logic for the scene
	// ...
	ggEnt := entities.NewGridGraphEntity()
	gem.Append(scn.GetRoot(), ggEnt)
}

func (scn *FlowFieldScene) Deinit() {
	// De-initialization logic for the scene
	// ...
}
