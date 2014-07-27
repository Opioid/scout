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

func Abs(s float32) float32 {
	return float32(math.Abs(float64(s)))
}

func Min(s, min float32) float32 {
	if s > min {
		return min
	} else {
		return s
	}
}

func Max(s, max float32) float32 {
	if s < max {
		return max
	} else {
		return s
	}
}

func Clamp(s, min, max float32) float32 {
	if s < min {
		return min
	} else if s > max {
		return max
	} else {
		return s
	}
}

func DegreesToRadians(s float32) float32 {
	return s * math.Pi / 180.0;
}