package scenes

import "gorl/fw/core/entities"

type Scene struct {
	rootNode *entities.Entity
}

// GetRoot returns the root node of the scene.
func (s *Scene) GetRoot() *entities.Entity {
	if s.rootNode == nil {
		s.rootNode = &entities.Entity{Name: "SceneRoot"}
	}
	return s.rootNode
}
