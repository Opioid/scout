package math

func Mini(s, min int) int {
	if s > min {
		return min
	} else {
		return s
	}
}

func Maxi(s, max int) int {
	if s < max {
		return max
	} else {
		return s
	}
}