package math

import (
	"math"
)

func HemisphereSample_uniform(r, s float32)Vector3 {
	r1 := s * 2.0 * math.Pi
	r2 := 1.0 - r
	sp := Sqrt(1.0 - r2 * r2)

	return MakeVector3(Cos(r1) * sp, Sin(r1) * sp, r2)
}

func HemisphereSample_cos(r, s float32) Vector3 {
	r1 := s * 2.0 * math.Pi
	r2 := Sqrt(1.0 - r)
	sp := Sqrt(1.0 - r2 * r2)

	return MakeVector3(Cos(r1) * sp, Sin(r1) * sp, r2)
}