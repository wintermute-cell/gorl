package math

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// ============================================================================
// The vector2.go file contains utility functions for working with rl.Vector2.
// ============================================================================

// Vector2Angle returns the angle between two vectors in radians.
func Vector2Angle(v1, v2 rl.Vector2) float32 {
	cross := v1.X*v2.Y - v1.Y*v2.X
	dot := v1.X*v2.X + v1.Y*v2.Y
	return float32(math.Atan2(float64(cross), float64(dot)))
}

// Vector2Clamp restricts a vector within the limits specified by min and max vectors.
func Vector2Clamp(input, min, max rl.Vector2) rl.Vector2 {
	if input.X < min.X {
		input.X = min.X
	} else if input.X > max.X {
		input.X = max.X
	}
	if input.Y < min.Y {
		input.Y = min.Y
	} else if input.Y > max.Y {
		input.Y = max.Y
	}
	return input
}

// Vector2NormalizeSafe returns the normalized vector. If the input vector is
// zero, it returns a zero vector instead of (NaN, NaN).
func Vector2NormalizeSafe(v rl.Vector2) rl.Vector2 {
	if v == rl.Vector2Zero() {
		return rl.Vector2Zero()
	} else {
		return rl.Vector2Normalize(v)
	}
}
