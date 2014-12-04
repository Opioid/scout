package math

import (
	"github.com/Opioid/math32"
	"math"
)

func Cos(x float32) float32 {
	return float32(math.Cos(float64(x)))
}

func Sin(x float32) float32 {
	return float32(math.Sin(float64(x)))
}

func InverseSqrt(x float32) float32 {
	return 1 / math32.Sqrt(x)
}

func Pow(x, y float32) float32 {
	return float32(math.Pow(float64(x), float64(y)))
}

func Exp2(x float32) float32 {
	return float32(math.Exp2(float64(x)))
}

func Floor(x float32) float32 {
	return float32(math.Floor(float64(x)))
}

func Saturate(x float32) float32 {
	return math32.Clamp(x, 0, 1)
}

func DegreesToRadians(x float32) float32 {
	return x * math.Pi / 180;
}