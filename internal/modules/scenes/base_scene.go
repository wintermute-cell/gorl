package scenes

import (
	"gorl/internal/core/entities/proto"
	"gorl/internal/core/gem"
)

type BaseScene struct {
	rootNode *proto.Entity
}

func (s *BaseScene) GetRoot() *proto.Entity {
	if s.rootNode == nil {
		s.rootNode = &proto.Entity{Name: "SceneRoot"}
	}
	return s.rootNode
}

// Update calls the Update method of all entities in the scene root nodes subtree.
func (scn *BaseScene) Update() {
	children := gem.GetChildren(scn.rootNode)
	for len(children) > 0 {
		child := children[0]
		child.Update()
		// remove the first element and append its children to the end of the slice
		children = append(children[1:], gem.GetChildren(child)...)
	}
}

// FixedUpdate calls the FixedUpdate method of all entities in the scene root nodes subtree.
func (scn *BaseScene) FixedUpdate() {
	children := gem.GetChildren(scn.rootNode)
	for len(children) > 0 {
		child := children[0]
		child.FixedUpdate()
		// remove the first element and append its children to the end of the slice
		children = append(children[1:], gem.GetChildren(child)...)
	}
}
