package gem

import (
	"gorl/fw/core/entities/proto"
	"gorl/fw/core/logging"
	"slices"
)

type gem struct {
	root    *entityNode
	nodeMap map[proto.IEntity]*entityNode
}

type entityNode struct {
	entity           proto.IEntity
	parent           *entityNode
	children         []*entityNode
	renderLayerIndex int64
}

const DefaultLayer = 0

var gemInstance gem

// Init initializes the Gem. This should be called once at the start of the program.
func Init() {
	gemInstance = gem{
		root: &entityNode{
			entity:           &proto.Entity{Name: "GemRoot"},
			parent:           nil,
			children:         make([]*entityNode, 0),
			renderLayerIndex: 0,
		},
		nodeMap: make(map[proto.IEntity]*entityNode),
	}
	gemInstance.root.entity.SetDrawIndex(0)
	gemInstance.nodeMap[gemInstance.root.entity] = gemInstance.root
}

// traverseAndFlatten traverses the hierarchy, calling Update() and
// FixedUpdate() and returns a slice of layers each being a slice of entities.
//
// The tree is traversed in a depth-first manner, meaning that the children of
// a node are visited before the siblings of the node.
func traverseAndFlatten(node *entityNode, shouldCallUpdate, shouldCallFixedUpdate bool) map[int64][]proto.IEntity {

	layers := make(map[int64][]proto.IEntity) // this will be returned
	stack := make([]*entityNode, 0, 100)      // we need a stack to traverse the tree
	stack = append(stack, node)

	for len(stack) > 0 {
		// pop off the stack
		currentNode := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// add the entity to the correct layer
		if layer, ok := layers[currentNode.renderLayerIndex]; ok {
			layers[currentNode.renderLayerIndex] = append(layer, currentNode.entity)
		} else {
			layers[currentNode.renderLayerIndex] = []proto.IEntity{currentNode.entity}
		}

		// call the update functions
		if shouldCallUpdate {
			currentNode.entity.Update()
		}
		if shouldCallFixedUpdate {
			currentNode.entity.FixedUpdate()
		}

		// push children on the stack
		for _, childNode := range currentNode.children {
			stack = append(stack, childNode)
		}
	}

	// sort the layers by the entities GetDrawIndex
	for _, layer := range layers {
		slices.SortStableFunc(layer, func(a, b proto.IEntity) int {
			return int(a.GetDrawIndex() - b.GetDrawIndex())
		})
	}

	return layers
}

// Deinit deinitializes the Gem by calling Deinit on every entity. This should
// be called once at the end of the program.
func Deinit() {
	entities := traverseAndFlatten(gemInstance.root, false, false)
	for _, layer := range entities {
		for _, entity := range layer {
			entity.Deinit()
		}
	}
}

// GetRoot returns the root entity of the Gem. It is used as the root parent for all other entities.
// Only use this if you need to access the root entity directly, ergo you are not using scenes.
func GetRoot() proto.IEntity {
	return gemInstance.root.entity
}

// AddEntity adds an entity to the Gem, under the parent entity, on the specified layer.
//
// If you don't know what layer to use, use gem.DefaultLayer (alias for 0).
func AddEntity(parent, entity proto.IEntity, renderLayerIndex int64) {
	// TODO: not sure if this works with Entity2D for example. https://codefibershq.com/blog/golang-why-nil-is-not-always-nil
	if parent == (*proto.Entity)(nil) || entity == (*proto.Entity)(nil) {
		logging.Error("Parent or entity is nil. Parent: %v, Entity: %v", parent, entity)
		return
	}

	parentNode, ok := gemInstance.nodeMap[parent]
	if !ok {
		logging.Error("Parent %v not found in node map, not adding entity %v", parent, entity)
		return
	}

	// add to hierarchy
	node := &entityNode{
		entity:           entity,
		parent:           parentNode,
		children:         make([]*entityNode, 0),
		renderLayerIndex: renderLayerIndex,
	}
	gemInstance.nodeMap[entity] = node

	// add to the parents children list
	parentNode.children = append(parentNode.children, node)

	// attach hierarchical transform capabilities if the node implements IEntity2D
	tryAddHierarchicalTransformCallbacks(node, parentNode)

	entity.SetDrawIndex(parent.GetDrawIndex() + 1) // by default we draw children in front of their parents
	parent.OnChildAdded(entity)
	entity.Init()
	logging.Info("Added entity %v to parent %v", entity, parent)
}

// RemoveEntity removes an entity from the Gem.
// This will also remove the entity from the hierarchy,
func RemoveEntity(entity proto.IEntity) {
	if entity == (*proto.Entity)(nil) {
		logging.Error("Entity is nil")
		return
	}

	//
	// This function first tries to find the entity in the entities list and in
	// the hierarchy map. If this is successful, it will call the Deinit and
	// OnChildRemoved callbacks and will then remove the entity from the
	// entities list and the hierarchy map.
	//

	entityNode, ok := gemInstance.nodeMap[entity]
	if !ok {
		logging.Error("Entity %v not found in node map, cannot remove. (was it already removed before?)", entity)
		return
	}
	parentNode := entityNode.parent

	// recusively remove all children
	for _, child := range entityNode.children {
		RemoveEntity(child.entity)
	}

	// deinit and remove the entity itself
	entity.Deinit()
	delete(gemInstance.nodeMap, entity)
	parentNode.entity.OnChildRemoved(entity) // parent callback

	// remove the entity from the parent's children list
	// TODO: there could be a more performant way to do this
	idx := slices.Index(entityNode.parent.children, entityNode)
	parentNode.children = append(parentNode.children[:idx], parentNode.children[idx+1:]...)

	logging.Info("Removed entity %v", entity)
}

// ReParentEntity reparents an entity to a new parent. Will not trigger
// OnChildRemoved but will trigger OnChildAdded!
func ReParentEntity(entity, newParent proto.IEntity) {
	if entity == (*proto.Entity)(nil) || newParent == (*proto.Entity)(nil) {
		logging.Error("Entity or new parent is nil. Entity: %v, New Parent: %v", entity, newParent)
		return
	}

	// getting the relevant nodes
	entityNode, ok := gemInstance.nodeMap[entity]
	if !ok {
		logging.Error("Entity %v not found in node map, cannot reparent. (was it already removed before?)", entity)
		return
	}
	parentNode := entityNode.parent
	if parentNode == nil {
		logging.Error("Entity %v is already a root entity, cannot reparent", entity)
		return
	}
	newParentNode, ok := gemInstance.nodeMap[newParent]
	if !ok {
		logging.Error("New parent %v not found in node map, cannot reparent", newParent)
		return
	}

	// remove from old parent
	idx := slices.Index(parentNode.children, entityNode)
	parentNode.children = append(parentNode.children[:idx], parentNode.children[idx+1:]...)

	// add to new parent
	entityNode.parent = newParentNode
	newParentNode.children = append(newParentNode.children, entityNode)
	entity.SetDrawIndex(newParent.GetDrawIndex() + 1) // by default we draw children in front of their parents

	// adjust to the new parent's transform 2D space if applicable
	if ent2d, ok := entity.(proto.IEntity2D); ok {
		if parent2d, ok := newParentNode.entity.(proto.IEntity2D); ok {
			ent2d.SetPosition(parent2d.GetPosition())
			ent2d.SetScale(parent2d.GetScale())
			ent2d.SetRotation(parent2d.GetRotation())
		}
	}

	logging.Info("Reparented entity %v to %v", entity, newParent)
}

// GetChildren returns a slice of all children of the parent entity.
func GetChildren(parent proto.IEntity) []proto.IEntity {
	if parent == nil {
		logging.Error("Parent is nil")
		return nil
	}

	node, ok := gemInstance.nodeMap[parent]
	if !ok {
		logging.Error("Parent %v not found in node map", parent)
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

	node, ok := gemInstance.nodeMap[entity]
	if !ok {
		logging.Error("Entity %v not found in hierarchy map", entity)
		return nil
	}

	return node.parent.entity
}
