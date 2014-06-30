package math

import "math"

func Sqrt(s float32) float32 {
	return float32(math.Sqrt(float64(s)))
}