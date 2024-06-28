package math

import (
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// ============================================================================
// The utility.go file contains utility functions for math operations that go
// beyond wrapping the stdlib as a generic.
// ============================================================================

// RandRange returns a random float32 between min and max.
// The returned value is in the range [min, max).
// e.g. min is inclusive, max is exclusive.
func RandRange(min, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

// ShortestLerp returns the shortest linear interpolation between two numbers.
func ShortestLerp(current, target, factor float32) float32 {
	// Calculate the difference
	difference := target - current

	// Calculate possible wrapped differences
	wrappedDifferencePlus := float64(difference + 360)
	wrappedDifferenceMinus := float64(difference - 360)

	// Check which one is the smallest in terms of absolute value
	if math.Abs(wrappedDifferencePlus) < math.Abs(float64(difference)) {
		difference = float32(wrappedDifferencePlus)
	} else if math.Abs(wrappedDifferenceMinus) < math.Abs(float64(difference)) {
		difference = float32(wrappedDifferenceMinus)
	}

	// Compute the lerped value
	lerped := current + difference*factor

	// Adjust the lerped value to be within the 0-360 range
	for lerped < 0 {
		lerped += 360
	}
	for lerped >= 360 {
		lerped -= 360
	}

	return lerped
}

// RotatePointAroundOrigin rotates a point around an origin by a given angle.
func RotatePointAroundOrigin(point, origin rl.Vector2, angle_deg float32) rl.Vector2 {
	angleRad := angle_deg * (math.Pi / 180) // Convert angle to radians

	cosAngle := float32(math.Cos(float64(angleRad)))
	sinAngle := float32(math.Sin(float64(angleRad)))

	// Translate point back to origin
	translated := rl.Vector2Subtract(point, origin)

	// Rotate point
	rotatedX := translated.X*cosAngle - translated.Y*sinAngle
	rotatedY := translated.X*sinAngle + translated.Y*cosAngle

	// Translate point back
	finalPoint := rl.NewVector2(rotatedX+origin.X, rotatedY+origin.Y)

	return finalPoint
}
