package math

import (
	"math"
)

// ============================================================================
// The generic.go file contains generic implementations and wrappers for golang
// stdlib math functions.
// ============================================================================

type number interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

type signed_number interface {
	int | int16 | int32 | int64 | float32 | float64
}

// Sign returns the sign of x (-1 if x < 0, 1 if x > 0, 0 if x == 0)
func Sign[T signed_number](x T) T {
	if x < 0 {
		return T(-1)
	} else if x > 0 {
		return T(1)
	} else {
		return T(0)
	}
}

// Abs returns the absolute value of x
func Abs[T number](x T) T {
	ret := x
	if x < 0 {
		ret = -x
	}
	return ret
}

// Clamps x between lower_bound and upper_bound, both inclusive.
// (Clamp will return at least lower_bound and at most upper_bound)
func Clamp[T number](x, lower_bound, upper_bound T) T {
	v := min(x, upper_bound)
	v = max(v, lower_bound)
	return v
}

// Round x to the nearest integer, either down if x < .5 or up if x >= .5.
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

// Pow returns x raised to the power of y.
func Pow[T number](x, y T) T {
	return T(math.Pow(float64(x), float64(y)))
}

// Lerp returns the linear interpolation between two numbers.
func Lerp[T number](a, b T, factor float32) T {
	return T(float32(a)*(1.0-factor) + (float32(b) * factor))
}
