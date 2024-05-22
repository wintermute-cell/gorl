package render

import (
	"gorl/fw/core/logging"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// A RenderStage is a stage in the rendering pipeline.
// It defines:
// - the target resolution at which the stage should render
type RenderStage struct {
	targetResolution rl.Vector2
	renderTexture    rl.RenderTexture2D

	camera               rl.Camera2D
	resolutionCorrection float32
	managedCameraZoom    float32
}

// Create a new RenderStage that renders at the given target resolution.
//
// The resolutionCorrection parameter is used to correct for the difference
// between the target resolution and the final screen resolution.
//
// If the stage renders at 2x the screen resolution, and you want it to behave
// the same as 1x resolution in terms of positioning and size, set the
// resolutionCorrection to 2.0 .
func NewRenderStage(targetResolution rl.Vector2, resolutionCorrection float32) *RenderStage {
	logging.Info("Creating new RenderStage with target resolution: %v", targetResolution)
	return &RenderStage{
		targetResolution: targetResolution,
		renderTexture: rl.LoadRenderTexture(
			int32(targetResolution.X),
			int32(targetResolution.Y),
		),
		camera: rl.NewCamera2D(
			rl.Vector2Zero(),
			rl.Vector2Zero(),
			0, resolutionCorrection,
		),
		resolutionCorrection: resolutionCorrection,
		managedCameraZoom:    1,
	}
}

func (rs *RenderStage) SetCameraTarget(target rl.Vector2) {
	// TODO: implement bounds clamping for the camera target
	rs.camera.Target = target
}

// GetCameraTarget returns the target of the camera for this RenderStage.
func (rs *RenderStage) GetCameraTarget() rl.Vector2 {
	return rs.camera.Target
}

// SetCameraOffset sets the offset of the camera for this RenderStage.
func (rs *RenderStage) SetCameraOffset(offset rl.Vector2) {
	rs.camera.Offset = rl.Vector2Scale(offset, rs.resolutionCorrection)
}

// GetCameraOffset returns the offset of the camera for this RenderStage.
func (rs *RenderStage) GetCameraOffset() rl.Vector2 {
	return rs.camera.Offset
}

// SetCameraZoom sets the zoom of the camera for this RenderStage.
func (rs *RenderStage) SetCameraZoom(zoom float32) {
	rs.managedCameraZoom = zoom
	rs.camera.Zoom = rs.managedCameraZoom * rs.resolutionCorrection
}

// GetCameraZoom returns the zoom of the camera for this RenderStage.
func (rs *RenderStage) GetCameraZoom() float32 {
	return rs.managedCameraZoom
}
