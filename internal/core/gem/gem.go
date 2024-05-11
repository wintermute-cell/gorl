package gem

import (
	"gorl/internal/core/entities/proto"
	"gorl/internal/logging"
)

type gem struct {
	entities      [][]proto.IEntity
	hierarchyRoot *hierarchyNode
	hierarchyMap  map[proto.IEntity]*hierarchyNode
}

type hierarchyNode struct {
	entity   proto.IEntity
	parent   proto.IEntity
	children []*hierarchyNode
}

const DefaultLayer = 0

var gemInstance gem

// Init initializes the Gem. This should be called once at the start of the program.
func Init() {
	gemInstance = gem{
		entities: make([][]proto.IEntity, 1), // we preallocate the default layer
		hierarchyRoot: &hierarchyNode{
			entity:   &proto.Entity{Name: "GemRoot"},
			parent:   nil,
			children: make([]*hierarchyNode, 0),
		},
		hierarchyMap: make(map[proto.IEntity]*hierarchyNode),
	}
	gemInstance.hierarchyRoot.entity.SetDrawIndex(0)
	gemInstance.hierarchyMap[gemInstance.hierarchyRoot.entity] = gemInstance.hierarchyRoot
	gemInstance.entities[0] = make([]proto.IEntity, 0, 1000) // this is the default layer, prealloc some space since it's the most used
}

// Deinit deinitializes the Gem. This should be called once at the end of the program.
func Deinit() {
	for _, layer := range gemInstance.entities {
		for _, entity := range layer {
			entity.Deinit()
		}
	}
}

// GetRoot returns the root entity of the Gem. It is used as the root parent for all other entities.
// Only use this if you need to access the root entity directly, ergo you are not using scenes.
func GetRoot() proto.IEntity {
	return gemInstance.hierarchyRoot.entity
}

// AddEntity adds an entity to the Gem, under the parent entity, on the specified layer.
//
// If you don't know what layer to use, use gem.DefaultLayer (alias for 0).
func AddEntity(parent, entity proto.IEntity, layerIndex int64) {
	if parent == (*proto.Entity)(nil) || entity == (*proto.Entity)(nil) { // TODO: not sure if this works with Entity2D for example. https://codefibershq.com/blog/golang-why-nil-is-not-always-nil
		logging.Error("Parent or entity is nil. Parent: %v, Entity: %v", parent, entity)
		return
	}
	gemInstance.entities[layerIndex] = append(gemInstance.entities[layerIndex], entity)

	// add to hierarchy
	node := &hierarchyNode{
		entity:   entity,
		parent:   parent,
		children: make([]*hierarchyNode, 0),
	}
	gemInstance.hierarchyMap[entity] = node
	if _, ok := gemInstance.hierarchyMap[parent]; !ok {
		// undo adding the entity to the entities list
		gemInstance.entities[layerIndex] = gemInstance.entities[layerIndex][:len(gemInstance.entities[layerIndex])-1]
		logging.Error("Parent %v not found in hierarchy map, cannot add", parent)
		return
	}
	gemInstance.hierarchyMap[parent].children = append(gemInstance.hierarchyMap[parent].children, node)

	// callbacks
	entity.SetDrawIndex(parent.GetDrawIndex() + 1) // by default we draw children in front of their parents
	entity.Init()
	parent.OnChildAdded(entity)
	logging.Info("Added entity %v to parent %v", entity, parent)
}

func RemoveEntity(entity proto.IEntity) {
	if entity == nil {
		logging.Error("Entity is nil")
		return
	}

	// remove from entities list
	found := false
	var foundInLayer int64
	for _, layer := range gemInstance.entities {
		for i, e := range layer {
			if e == entity {
				layer = append(layer[:i], layer[i+1:]...)
				foundInLayer = int64(i)
				found = true
			}
		}
	}
	if !found {
		logging.Error("Entity %v not found", entity)
		return
	}

	// remove from hierarchy
	node, ok := gemInstance.hierarchyMap[entity]
	if !ok {
		logging.Error("Entity %v not found in hierarchy map, cannot remove", entity)
		return
	}
	if node.parent == nil {
		logging.Fatal("Entity %v is root or something went wrong (has no parent), cannot remove", entity)
		return
	}
	parentNode := gemInstance.hierarchyMap[node.parent]
	idx := -1
	for i, child := range parentNode.children {
		if child == node {
			idx = i
			break
		}
	}
	if idx == -1 {
		// undo removing the entity from the entities list
		gemInstance.entities[foundInLayer] = append(gemInstance.entities[foundInLayer], entity)
		logging.Error("Entity %v not found in parent %v's children, cannot remove", entity, node.parent)
		return
	}
	parentNode.children = append(parentNode.children[:idx], parentNode.children[idx+1:]...)
	delete(gemInstance.hierarchyMap, entity)

	// callbacks
	entity.Deinit()
	parentNode.entity.OnChildRemoved(entity)
	logging.Info("Removed entity %v", entity)
}

func ReParentEntity(entity, newParent proto.IEntity) {
	if entity == nil || newParent == nil {
		logging.Error("Entity or new parent is nil. Entity: %v, New Parent: %v", entity, newParent)
		return
	}

	// remove from hierarchy
	node, ok := gemInstance.hierarchyMap[entity]
	if !ok {
		logging.Error("Entity %v not found in hierarchy map, cannot reparent", entity)
		return
	}
	if node.parent == nil {
		logging.Fatal("Entity %v is root or something went wrong (has no parent), cannot reparent", entity)
		return
	}
	parentNode := gemInstance.hierarchyMap[node.parent]
	idx := -1
	for i, child := range parentNode.children {
		if child == node {
			idx = i
			break
		}
	}
	if idx == -1 {
		logging.Error("Entity %v not found in parent %v's children, cannot reparent", entity, node.parent)
		return
	}
	parentNode.children = append(parentNode.children[:idx], parentNode.children[idx+1:]...)
	delete(gemInstance.hierarchyMap, entity) // remove from map so we can re-add it later with the new parent

	// add to new parent
	node.parent = newParent
	newParentNode := gemInstance.hierarchyMap[newParent]
	newParentNode.children = append(newParentNode.children, node)
	gemInstance.hierarchyMap[entity] = node

	logging.Info("Reparented entity %v to %v", entity, newParent)
}

// GetChildren returns a slice of all children of the parent entity.
func GetChildren(parent proto.IEntity) []proto.IEntity {
	if parent == nil {
		logging.Error("Parent is nil")
		return nil
	}

	node, ok := gemInstance.hierarchyMap[parent]
	if !ok {
		logging.Error("Parent %v not found in hierarchy map", parent)
		return nil
	}

	children := make([]proto.IEntity, len(node.children))
	for i, child := range node.children {
		children[i] = child.entity
	}
	return children
}

// GetParent returns the parent of the entity.
func GetParent(entity proto.IEntity) proto.IEntity {
	if entity == nil {
		logging.Error("Entity is nil")
		return nil
	}

	node, ok := gemInstance.hierarchyMap[entity]
	if !ok {
		logging.Error("Entity %v not found in hierarchy map", entity)
		return nil
	}

	return node.parent
}
