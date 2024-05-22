package gem

import (
	"gorl/fw/core/entities/proto"
	"gorl/fw/core/logging"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
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

	// some IEntity2D specific stuff
	if ent2d, ok := entity.(proto.IEntity2D); ok {

		// move the entitites transform to the parents local space, if applicable
		if parent2d, ok := parent.(proto.IEntity2D); ok {
			ent2d.SetPosition(rl.Vector2Add(parent2d.GetPosition(), ent2d.GetPosition()))
			ent2d.SetScale(rl.Vector2Multiply(parent2d.GetScale(), ent2d.GetScale()))
			ent2d.SetRotation(parent2d.GetRotation() + ent2d.GetRotation())
		}

		// we register transform setter callback here to allow for tree
		// inheritance to the children.
		ent2d.AddOnSetPositionCallback(func(offset rl.Vector2) {
			for _, child := range node.children {
				if childEntity2D, ok := child.entity.(proto.IEntity2D); ok {
					parentPos := ent2d.GetPosition()
					childPos := childEntity2D.GetPosition()
					relativePos := rl.Vector2Subtract(childPos, parentPos)
					childEntity2D.SetPosition(rl.Vector2Add(parentPos, rl.Vector2Add(relativePos, offset)))
				}
			}
		})

		ent2d.AddOnSetScaleCallback(func(offset rl.Vector2) {
			for _, child := range node.children {
				if childEntity2D, ok := child.entity.(proto.IEntity2D); ok {
					childEntity2D.SetScale(rl.Vector2Multiply(childEntity2D.GetScale(), offset))
				}
			}
		})

		ent2d.AddOnSetRotationCallback(func(offset float32) {
			for _, child := range node.children {
				if childEntity2D, ok := child.entity.(proto.IEntity2D); ok {
					parentPos := ent2d.GetPosition()
					childPos := childEntity2D.GetPosition()
					relativePos := rl.Vector2Subtract(childPos, parentPos)

					// Apply rotation transformation relative to parent
					offsetToRad := offset * (math.Pi / 180)
					sin, cos := float32(math.Sin(float64(offsetToRad))), float32(math.Cos(float64(offsetToRad)))
					rotatedPos := rl.Vector2{
						X: relativePos.X*cos - relativePos.Y*sin,
						Y: relativePos.X*sin + relativePos.Y*cos,
					}
					childEntity2D.SetPosition(rl.Vector2Add(parentPos, rotatedPos))
					childEntity2D.SetRotation(childEntity2D.GetRotation() + offset)
				}
			}
		})
	}

	// callbacks
	entity.SetDrawIndex(parent.GetDrawIndex() + 1) // by default we draw children in front of their parents
	entity.Init()
	parent.OnChildAdded(entity)
	logging.Info("Added entity %v to parent %v", entity, parent)
}

// RemoveEntity removes an entity from the Gem.
// This will also remove the entity from the hierarchy,
func RemoveEntity(entity proto.IEntity) {
	if entity == nil {
		logging.Error("Entity is nil")
		return
	}

	//
	// This function first tries to find the entity in the entities list and in
	// the hierarchy map. If this is successful, it will call the Deinit and
	// OnChildRemoved callbacks and will then remove the entity from the
	// entities list and the hierarchy map.
	//

	// remove from entities list
	found := false
	foundInLayer := int64(-1)
	foundAtIndex := -1
	for lidx, layer := range gemInstance.entities {
		for eidx, e := range layer {
			if e == entity {
				// we won't remove the entity here because we are not yet sure
				// that we won't encounter an error.
				foundAtIndex = eidx
				foundInLayer = int64(lidx)
				logging.Debug("Found entity %v in layer %v at index %v", entity, foundInLayer, foundAtIndex)
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
		logging.Error("Entity %v not found in parent %v's children, cannot remove", entity, node.parent)
		return
	}

	// cut the entity out of the entities matrix. we do that before we remove
	// the children, because removing them will shift the indices.
	// since we still keep a reference until a bit later, we can still access
	// the entity and its children.
	gemInstance.entities[foundInLayer] = append(
		gemInstance.entities[foundInLayer][:foundAtIndex],
		gemInstance.entities[foundInLayer][foundAtIndex+1:]...,
	)

	// remove children
	for _, child := range node.children {
		// TODO: this will potentially call a lot of OnChildRemoved callbacks.
		// Is this intended or should only the directly removed entity call
		// that callback? To solve, we could create a private removeEntity func
		// that take an entity and a bool to indicate if it should call the
		// callback, and just redirect this function to that internal impl.
		RemoveEntity(child.entity)
	}

	// callbacks
	entity.Deinit()
	parentNode.entity.OnChildRemoved(entity)

	// remove the entity from the hierarchy, this should have been the last reference.
	parentNode.children = append(parentNode.children[:idx], parentNode.children[idx+1:]...)
	delete(gemInstance.hierarchyMap, entity)

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
	entity.SetDrawIndex(newParentNode.entity.GetDrawIndex() + 1) // by default we draw children in front of their parents

	// some IEntity2D specific stuff
	if ent2d, ok := entity.(proto.IEntity2D); ok {
		// move the entitites transform to the parents local space, if applicable
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
