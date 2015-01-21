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

func Sincos(x float32) (sin, cos float32) {
	sin64, cos64 := math.Sincos(float64(x))
	sin = float32(sin64)
	cos = float32(cos64)
	return
}

func Pow(x, y float32) float32 {
	return float32(math.Pow(float64(x), float64(y)))
}

func Exp(x float32) float32 {
	return float32(math.Exp(float64(x)))
}

func Exp2(x float32) float32 {
	return float32(math.Exp2(float64(x)))
}

func Saturate(x float32) float32 {
	return math32.Clamp(x, 0.0, 1.0)
}

func DegreesToRadians(x float32) float32 {
	return x * math.Pi / 180.0;
}

func IsInf(x float32) bool {
	return math.IsInf(float64(x), 0.0)
}