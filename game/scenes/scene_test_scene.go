package scenes

import (
	"gorl/fw/modules/scenes"
)

// This checks at compile time if the interface is implemented
var _ scenes.Scene = (*TestSceneScene)(nil)

// TestScene Scene
type TestSceneScene struct {
	// Custom Fields
	// Add fields here for any state that the scene should keep track of
	// ...
}

func (scn *TestSceneScene) Init() {
	// Initialization logic for the scene
	// ...
}

func (scn *TestSceneScene) Deinit() {
	// De-initialization logic for the scene
}

func (scn *TestSceneScene) DrawGUI() {
	// Draw the GUI for the scene
}

func (scn *TestSceneScene) Draw() {
	// Draw the scene
}
