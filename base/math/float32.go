package math

import "math"

func Cos(x float32) float32 {
	return float32(math.Cos(float64(x)))
}

func Sin(x float32) float32 {
	return float32(math.Sin(float64(x)))
}

func Sqrt(x float32) float32 {
	return float32(math.Sqrt(float64(x)))
}

func InverseSqrt(x float32) float32 {
	return 1.0 / Sqrt(x)
}

func Pow(x, y float32) float32 {
	return float32(math.Pow(float64(x), float64(y)))
}

func Exp2(x float32) float32 {
	return float32(math.Exp2(float64(x)))
}

func Abs(x float32) float32 {
	return float32(math.Abs(float64(x)))
}

func Floor(x float32) float32 {
	return float32(math.Floor(float64(x)))
}

func Min(x, y float32) float32 {
	return float32(math.Min(float64(x), float64(y)))
}

func Max(x, y float32) float32 {
	return float32(math.Max(float64(x), float64(y)))
}

func Clamp(x, min, max float32) float32 {
	return float32(math.Min(math.Max(float64(x), float64(min)), float64(max)))
}

func DegreesToRadians(x float32) float32 {
	return x * math.Pi / 180.0;
}