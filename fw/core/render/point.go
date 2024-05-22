package render

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// TODO: fix this description, now takes a camera as argument
//
// PointToCameraSpace converts a point in screen space to camera space.
// It uses the current render stage's camera to perform the conversion.
// If there is no current render stage, the point is returned unchanged and a
// warning is logged.
func PointToCameraSpace(camera rl.Camera2D, point rl.Vector2) rl.Vector2 {
	// Translate the point based on the camera's target
	translatedPoint := rl.Vector2Subtract(point, camera.Target)

	// Rotate the point around the camera's target
	angleRad := float64(camera.Rotation * (math.Pi / 180.0))
	sin, cos := math.Sincos(angleRad)
	rotatedPoint := rl.Vector2{
		X: translatedPoint.X*float32(cos) - translatedPoint.Y*float32(sin),
		Y: translatedPoint.X*float32(sin) + translatedPoint.Y*float32(cos),
	}

	// Apply the camera's zoom
	zoomedPoint := rl.Vector2Scale(rotatedPoint, 1.0/camera.Zoom)

	// Translate based on the camera's offset
	finalPoint := rl.Vector2Add(zoomedPoint, camera.Offset)

	return finalPoint

}
