package render

import (
	"gorl/fw/core/logging"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type RenderSystem struct {
	screenResolution rl.Vector2
	renderTexture    rl.RenderTexture2D
	currentStage     *RenderStage
}

var rs *RenderSystem

func Init(screenResolution rl.Vector2) {
	rs = newRenderSystem(screenResolution)
}

// newRenderSystem creates a new RenderSystem with the given screen resolution.
func newRenderSystem(screenResolution rl.Vector2) *RenderSystem {
	logging.Info("Creating new RenderSystem with screen resolution: %v", screenResolution)
	renderSystem := &RenderSystem{
		screenResolution: screenResolution,
		renderTexture: rl.LoadRenderTexture(
			int32(screenResolution.X),
			int32(screenResolution.Y),
		),
	}
	return renderSystem
}

// EnableRenderStage sets the RenderSystem to render using the given RenderStage.
func EnableRenderStage(stage *RenderStage) {
	if stage == nil {
		logging.Error("Attempted to enable nil RenderStage.")
		return
	}

	// if there is a current stage, flush it to the render system's texture
	FlushRenderStage()

	rs.currentStage = stage
	rl.BeginTextureMode(rs.currentStage.renderTexture)
	rl.BeginMode2D(rs.currentStage.Camera)
	applyCameraTransformations(&stage.Camera)
}

// FlushRenderStage flushes the current RenderStage to the RenderSystem's
// texture. This is done automatically when switching stages, but can be called
// manually to finalize the last stage for this frame.
func FlushRenderStage() {
	// if a stage is already enabled, end it and flush its texture to the
	// render systems texture
	if rs.currentStage != nil {
		// if there is a current stage, we must be in it's texture mode.
		// end it.
		rl.EndMode2D()
		rl.EndTextureMode()

		// draw the current stage's texture to the render system's texture
		rl.BeginTextureMode(rs.renderTexture)
		rl.DrawTexturePro(
			rs.currentStage.renderTexture.Texture,
			rl.NewRectangle(0, 0, rs.currentStage.targetResolution.X, -rs.currentStage.targetResolution.Y),
			rl.NewRectangle(0, 0, rs.screenResolution.X, rs.screenResolution.Y),
			rl.NewVector2(0, 0), 0, rl.White)
		rl.EndTextureMode()

		rs.currentStage = nil
	}
}

// RenderToScreen renders the accumulated contents of all render stages used
// with the render system to the screen.
func RenderToScreen() {

	if rs.currentStage != nil {
		logging.Warning("RenderSystem.RenderToScreen called with active RenderStage. Call FlushRenderStage first.")
		FlushRenderStage()
	}

	rl.DrawTexturePro(
		rs.renderTexture.Texture,
		rl.NewRectangle(0, 0, rs.screenResolution.X, -rs.screenResolution.Y),
		rl.NewRectangle(0, 0, rs.screenResolution.X, rs.screenResolution.Y),
		rl.NewVector2(0, 0), 0, rl.White)

	// clear the render texture for the next frame
	rl.BeginTextureMode(rs.renderTexture)
	rl.ClearBackground(rl.Blank)
	rl.EndTextureMode()
}
