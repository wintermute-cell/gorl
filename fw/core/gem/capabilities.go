package gem

import (
	"gorl/fw/core/entities/proto"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// tryAddHierarchicalTransformCallbacks tries to add the hierarchical transform
// callbacks to the entity. Will do nothing if the entity does not implement
// IEntity2D.
func tryAddHierarchicalTransformCallbacks(entityNode, parentNode *entityNode) {

	if ent2d, ok := entityNode.entity.(proto.IEntity2D); ok {

		// move the entitites transform to the parents local space, if applicable
		if parent2d, ok := parentNode.entity.(proto.IEntity2D); ok {
			ent2d.SetPosition(rl.Vector2Add(parent2d.GetPosition(), ent2d.GetPosition()))
			ent2d.SetScale(rl.Vector2Multiply(parent2d.GetScale(), ent2d.GetScale()))
			ent2d.SetRotation(parent2d.GetRotation() + ent2d.GetRotation())
		}

		// we register transform setter callback here to allow for tree
		// inheritance to the children.
		ent2d.AddOnSetPositionCallback(func(offset rl.Vector2) {
			for _, child := range entityNode.children {
				if childEntity2D, ok := child.entity.(proto.IEntity2D); ok {
					parentPos := ent2d.GetPosition()
					childPos := childEntity2D.GetPosition()
					relativePos := rl.Vector2Subtract(childPos, parentPos)
					childEntity2D.SetPosition(rl.Vector2Add(parentPos, rl.Vector2Add(relativePos, offset)))
				}
			}
		})

		ent2d.AddOnSetScaleCallback(func(offset rl.Vector2) {
			for _, child := range entityNode.children {
				if childEntity2D, ok := child.entity.(proto.IEntity2D); ok {
					childEntity2D.SetScale(rl.Vector2Multiply(childEntity2D.GetScale(), offset))
				}
			}
		})

		ent2d.AddOnSetRotationCallback(func(offset float32) {
			for _, child := range entityNode.children {
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

}
