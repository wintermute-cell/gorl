package gem

import rl "github.com/gen2brain/raylib-go/raylib"

// DebugDrawEntities draws a list of all entities by their name.
func DebugDrawEntities(position rl.Vector2, size int32) {
	rl.DrawText("Entities:", int32(position.X), int32(position.Y), size, rl.Lime)
	for lidx, layer := range gemInstance.entities {
		for eidx, entity := range layer {
			rl.DrawText(entity.GetName(), int32(position.X),
				int32(position.Y+float32(size*(int32(eidx)+1+int32(lidx)*int32(len(layer))))),
				size, rl.Lime)
		}
	}
}

// Helper function to recursively draw the hierarchy.
func drawHierarchyNode(node *hierarchyNode, position rl.Vector2, size int32, depth int32) rl.Vector2 {
	if node == nil {
		return position
	}

	// Calculate the position for the current entity
	xPos := int32(position.X) + depth*size
	yPos := int32(position.Y)

	// Draw the current entity name with indentation based on depth
	rl.DrawText(node.entity.GetName(), xPos, yPos, size, rl.Lime)

	// Update the position for the next entity
	nextPosition := rl.Vector2{X: position.X, Y: position.Y + float32(size)}

	// Recursively draw each child, updating the position
	for _, child := range node.children {
		nextPosition = drawHierarchyNode(child, nextPosition, size, depth+1)
	}

	// Return the updated position for the next sibling
	return rl.Vector2{X: position.X, Y: nextPosition.Y}
}

func DebugDrawHierarchy(position rl.Vector2, size int32) {
	rl.DrawText("Hierarchy:", int32(position.X), int32(position.Y), size, rl.Lime)
	if gemInstance.hierarchyRoot != nil {
		drawHierarchyNode(gemInstance.hierarchyRoot, rl.Vector2{X: position.X, Y: position.Y + float32(size)}, size, 1)
	}
}
