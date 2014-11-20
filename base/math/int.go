package math

func Mini(s, min int32) int32 {
	if s > min {
		return min
	} else {
		return s
	}
}

func Maxi(s, max int32) int32 {
	if s < max {
		return max
	} else {
		return s
	}
}

func Clampi(x, min, max int32) int32 {
	return Mini(Maxi(x, min), max)
}