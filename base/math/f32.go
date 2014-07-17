package math

import "math"

func Cos(s float32) float32 {
	return float32(math.Cos(float64(s)))
}

func Sin(s float32) float32 {
	return float32(math.Sin(float64(s)))
}

func Sqrt(s float32) float32 {
	return float32(math.Sqrt(float64(s)))
}

func DegreesToRadians(s float32) float32 {
	return s * math.Pi / 180.0;
}