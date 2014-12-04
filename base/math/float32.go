package math

import "math"

func Cos(x float32) float32 {
	return float32(math.Cos(float64(x)))
}

func Sin(x float32) float32 {
	return float32(math.Sin(float64(x)))
}

func Acos(x float32) float32 {
	return float32(math.Acos(float64(x)))
}

// http://http.developer.nvidia.com/Cg/acos.html
// Absolute error <= 6.7e-5
func FastAcos(x float32) float32 {
	negate := float32(0)
	if x < 0 {
		negate = float32(1)
	}

	x = Absf(x)
	ret := float32(-0.0187293)
	ret *= x
	ret += 0.0742610
	ret *= x
	ret -= 0.2121144
	ret *= x
	ret += 1.5707288
	ret *= Sqrt(1 - x)
	ret -= 2 * negate * ret
	return negate * math.Pi + ret
}

func Atan2(x, y float32) float32 {
	return float32(math.Atan2(float64(x), float64(y)))
}

// http://http.developer.nvidia.com/Cg/atan2.html
func FastAtan2(x, y float32) float32 {
	ax := Absf(x)
	ay := Absf(y)  

	t0 := Maxf(ax, ay)
	t1 := Minf(ax, ay)
	t2 := 1 / t0
	t2 = t1 * t2

	t1 = t2 * t2
	t0 =         - 0.013480470
	t0 = t0 * t1 + 0.057477314
	t0 = t0 * t1 - 0.121239071
	t0 = t0 * t1 + 0.195635925
	t0 = t0 * t1 - 0.332994597
	t0 = t0 * t1 + 0.999995630
	t2 = t0 * t2

	if ax > ay {
		t2 = 1.570796327 - t2
	}

	if y < 0 {
		t2 = 3.141592654 - t2
	}

	if x < 0 {
		t2 = -t2
	}

	return t2
}

func Sqrt(x float32) float32 {
	return float32(math.Sqrt(float64(x)))
}

func InverseSqrt(x float32) float32 {
	return 1 / Sqrt(x)
}

func Pow(x, y float32) float32 {
	return float32(math.Pow(float64(x), float64(y)))
}

func Exp2(x float32) float32 {
	return float32(math.Exp2(float64(x)))
}

func Absf(x float32) float32 {
	return Maxf(x, -x)
}

func Floor(x float32) float32 {
	return float32(math.Floor(float64(x)))
}

func Minf(x, y float32) float32 {
	if x < y {
		return x
	} else {
		return y
	}
}

func Maxf(x, y float32) float32 {
	if x > y {
		return x
	} else {
		return y
	}
}

func Clampf(x, min, max float32) float32 {
	return Minf(Maxf(x, min), max)
}

func Saturate(x float32) float32 {
	return Clampf(x, 0, 1)
}

func DegreesToRadians(x float32) float32 {
	return x * math.Pi / 180;
}

func IsNaN(x float32) bool {
	return math.IsNaN(float64(x))
}