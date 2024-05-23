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

	Camera               rl.Camera2D
	resolutionCorrection float32
	managedCameraZoom    float32
}

var renderStages map[int64]*RenderStage = make(map[int64]*RenderStage)

// Create a new RenderStage that renders at the given target resolution.
// You can later refer to this RenderStage by the given key.
// If you don't know what key to use, use `gem.DefaultLayer`.
// The key will be returned by this function for convenience.
//
// The resolutionCorrection parameter is used to correct for the difference
// between the target resolution and the final screen resolution.
//
// If the stage renders at 2x the screen resolution, and you want it to behave
// the same as 1x resolution in terms of positioning and size, set the
// resolutionCorrection to 2.0 .
func CreateRenderStage(key int64, targetResolution rl.Vector2, resolutionCorrection float32) (layerKey int64) {
	layerKey = key
	if _, ok := renderStages[key]; ok {
		logging.Error("RenderStage with key %v already exists, choose another one.", key)
		return
	}

	logging.Info("Creating new RenderStage with target resolution: %v", targetResolution)
	renderStages[key] = &RenderStage{
		targetResolution: targetResolution,
		renderTexture: rl.LoadRenderTexture(
			int32(targetResolution.X),
			int32(targetResolution.Y),
		),
		Camera: rl.NewCamera2D(
			rl.Vector2Zero(),
			rl.Vector2Zero(),
			0, resolutionCorrection,
		),
		resolutionCorrection: resolutionCorrection,
		managedCameraZoom:    1,
	}
	return
}

// GetRenderStage returns the RenderStage with the given key.
// If the RenderStage does not exist, it will log an error and return nil.
func GetRenderStage(key int64) *RenderStage {
	if stage, ok := renderStages[key]; ok {
		return stage
	}
	logging.Error("RenderStage with key %v does not exist.", key)
	return nil
}

func (rs *RenderStage) SetCameraTarget(target rl.Vector2) {
	// TODO: implement bounds clamping for the camera target
	rs.Camera.Target = target
}

// GetCameraTarget returns the target of the camera for this RenderStage.
func (rs *RenderStage) GetCameraTarget() rl.Vector2 {
	return rs.Camera.Target
}

// SetCameraOffset sets the offset of the camera for this RenderStage.
func (rs *RenderStage) SetCameraOffset(offset rl.Vector2) {
	rs.Camera.Offset = rl.Vector2Scale(offset, rs.resolutionCorrection)
}

// GetCameraOffset returns the offset of the camera for this RenderStage.
func (rs *RenderStage) GetCameraOffset() rl.Vector2 {
	return rs.Camera.Offset
}

// SetCameraZoom sets the zoom of the camera for this RenderStage.
func (rs *RenderStage) SetCameraZoom(zoom float32) {
	rs.managedCameraZoom = zoom
	rs.Camera.Zoom = rs.managedCameraZoom * rs.resolutionCorrection
}

// GetCameraZoom returns the zoom of the camera for this RenderStage.
func (rs *RenderStage) GetCameraZoom() float32 {
	return rs.managedCameraZoom
}
