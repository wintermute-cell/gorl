package scenes

import (
	"gorl/fw/core/entities/proto"
	"gorl/fw/modules/scenes"
)

// This checks at compile time if the interface is implemented
var _ scenes.Scene = (*FlowFieldScene)(nil)

// FlowField Scene
type FlowFieldScene struct {
	// Custom Fields
	// Add fields here for any state that the scene should keep track of
	// ...
	Root *proto.Entity
}

func (scn *FlowFieldScene) Init() {
	// Initialization logic for the scene
	// ...
}

func (scn *FlowFieldScene) Deinit() {
	// De-initialization logic for the scene
}

func (scn *FlowFieldScene) Update() {
	// Update logic for the scene
	// ...
}

func (scn *FlowFieldScene) FixedUpdate() {
	// FixedUpdate logic for the scene
	// ...
}

func (scn *FlowFieldScene) GetRoot() *proto.Entity {
	// Return the root entity for the scene
	return scn.Root
}
