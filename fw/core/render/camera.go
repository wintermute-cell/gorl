package render

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type cameraTransformationBuffer struct {
	Position       []rl.Vector2
	PositionChange []rl.Vector2
	Offset         []rl.Vector2
	OffsetChange   []rl.Vector2
	Rotation       []float32
	RotationChange []float32
	Zoom           []float32
	ZoomChange     []float32
}

// Resets all slices in cameraTransformationBuffer to empty without reallocation
func (ctb *cameraTransformationBuffer) reset() {
	ctb.Position = ctb.Position[:0]
	ctb.PositionChange = ctb.PositionChange[:0]
	ctb.Offset = ctb.Offset[:0]
	ctb.OffsetChange = ctb.OffsetChange[:0]
	ctb.Rotation = ctb.Rotation[:0]
	ctb.RotationChange = ctb.RotationChange[:0]
	ctb.Zoom = ctb.Zoom[:0]
	ctb.ZoomChange = ctb.ZoomChange[:0]
}

var cameraTransformations cameraTransformationBuffer

func applyCameraTransformations(camera *rl.Camera2D) {
	for _, position := range cameraTransformations.Position {
		camera.Target = position
	}
	for _, positionChange := range cameraTransformations.PositionChange {
		camera.Target = rl.Vector2Add(camera.Target, positionChange)
	}
	for _, offset := range cameraTransformations.Offset {
		camera.Offset = offset
	}
	for _, offsetChange := range cameraTransformations.OffsetChange {
		camera.Offset = rl.Vector2Add(camera.Offset, offsetChange)
	}
	for _, rotation := range cameraTransformations.Rotation {
		camera.Rotation = rotation
	}
	for _, rotationChange := range cameraTransformations.RotationChange {
		camera.Rotation += rotationChange
	}
	for _, zoom := range cameraTransformations.Zoom {
		camera.Zoom = zoom
	}
	for _, zoomChange := range cameraTransformations.ZoomChange {
		camera.Zoom += zoomChange
	}
	cameraTransformations.reset()
}

// TODO: these transforms have to be applied to the correct camera / render stage
func SetCameraPosition(position rl.Vector2) {
	cameraTransformations.Position = append(cameraTransformations.Position, position)
}
