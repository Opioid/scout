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

func Min(x, min float32) float32 {
	if x > min {
		return min
	} else {
		return x
	}
}

func Max(x, max float32) float32 {
	if x < max {
		return max
	} else {
		return x
	}
}

func Clamp(x, min, max float32) float32 {
	if x < min {
		return min
	} else if x > max {
		return max
	} else {
		return x
	}
}

func DegreesToRadians(x float32) float32 {
	return x * math.Pi / 180.0;
}