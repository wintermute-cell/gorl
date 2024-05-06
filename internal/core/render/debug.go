package render

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// DebugDrawStageViewports draws the render stage viewports relative to the
// screen resolution, taking into account the screen resolution and camera.
func DebugDrawStageViewports(position rl.Vector2, widgetScale int32, renderSystem *RenderSystem, renderStages []*RenderStage) {
	// Calculate scaled screen resolution
	screenRect := rl.Rectangle{
		X:      position.X,
		Y:      position.Y,
		Width:  float32(renderSystem.screenResolution.X) / float32(widgetScale),
		Height: float32(renderSystem.screenResolution.Y) / float32(widgetScale),
	}

	// Draw the main viewport
	rl.DrawRectangleRec(screenRect, rl.Fade(rl.Black, 0.5))

	colors := []rl.Color{
		rl.Red,
		rl.Green,
		rl.Blue,
		rl.Orange,
		rl.Purple,
		rl.Pink,
		rl.Yellow,
		rl.SkyBlue,
		rl.Lime,
	}

	// Loop through each render stage
	for idx, stage := range renderStages {
		// Calculate the viewport rectangle based on render resolution and camera zoom
		viewportRect := rl.Rectangle{
			X:      position.X + float32(stage.camera.Target.X-(stage.camera.Offset.X/stage.resolutionCorrection))*stage.camera.Zoom/float32(widgetScale),
			Y:      position.Y + float32(stage.camera.Target.Y-(stage.camera.Offset.Y/stage.resolutionCorrection))*stage.camera.Zoom/float32(widgetScale),
			Width:  float32(renderSystem.screenResolution.X) * stage.managedCameraZoom / float32(widgetScale),
			Height: float32(renderSystem.screenResolution.Y) * stage.managedCameraZoom / float32(widgetScale),
		}

		// Draw the viewport for this stage
		color := colors[idx%len(colors)]
		rl.DrawRectangleLinesEx(viewportRect, 2, color)
	}

}
