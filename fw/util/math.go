package util

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// ============================================================================
//
// THIS FILE IS DEPRECATED, USE THE FUNCTIONS FROM core/math INSTEAD
// TODO: References to old functions need to be replaced with the new ones.
//
// ============================================================================

type number interface {
	int | int16 | int32 | int64 | uint | uint16 | uint32 | uint64 | float32 | float64
}

type signed_number interface {
	int | int16 | int32 | int64 | float32 | float64
}

// Max will return the maximum value between x and y
//
// Deprecated:
// Use functions from core/math instead.
func Max[T number](x, y T) T {
	return T(math.Max(float64(x), float64(y)))
}

// Min will return the minimum value between x and y.
func Min[T number](x, y T) T {
	return T(math.Min(float64(x), float64(y)))
}

// Clamps x between lower_bound and upper_bound, both inclusive.
// (Clamp will return at least lower_bound and at most upper_bound)
//
// Deprecated:
// Use functions from core/math instead.
func Clamp[T number](x, lower_bound, upper_bound T) T {
	v := Min(x, upper_bound)
	v = Max(v, lower_bound)
	return v
}

// Round x to the nearest integer, either down if x < .5 or up if x >= .5.
//
// Deprecated:
// Use functions from core/math instead.
func Round[T number](x T) T {
	integer, fraction := math.Modf(float64(x))
	v := x
	if fraction >= 0.5 {
		v = T(integer) + T(1.0)
	} else {
		v = T(integer)
	}
	return v
}

// Floor returns the largest integer value less than or equal to x.
func Floor[T signed_number](x T) T {
	return T(math.Floor(float64(x)))
}

// Ceil returns the smallest integer value greater than or equal to x.
func Ceil[T signed_number](x T) T {
	return T(math.Ceil(float64(x)))
}

// Vector2NormalizeSafe returns the normalized vector. If the input vector is
// zero, it returns a zero vector instead of (NaN, NaN).
//
// Deprecated:
// Use functions from core/math instead.
func Vector2NormalizeSafe(v rl.Vector2) rl.Vector2 {
	if v == rl.Vector2Zero() {
		return rl.Vector2Zero()
	} else {
		return rl.Vector2Normalize(v)
	}
}

// RotatePointAroundOrigin rotates a point around an origin by a given angle.
//
// Deprecated:
// Use functions from core/math instead.
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
