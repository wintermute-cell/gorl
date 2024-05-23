package gem

import (
	"gorl/fw/core/entities/proto"
	"gorl/fw/core/logging"
	"gorl/fw/core/render"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// ScreenspaceToCameraspace converts a point in screenspace to cameraspace.
// Uses the camera attached to the render layer the entity is on.
func ScreenspaceToCameraspace(entity proto.IEntity, point rl.Vector2) rl.Vector2 {
	entityNode, ok := gemInstance.nodeMap[entity]
	if !ok {
		logging.Error("Entity not found in gemInstance.nodeMap, could not convert screenspace to cameraspace.")
		return point
	}
	cam := render.GetRenderStage(entityNode.renderLayerIndex).Camera
	newPoint := render.PointToCameraSpace(cam, point)
	return newPoint
}
