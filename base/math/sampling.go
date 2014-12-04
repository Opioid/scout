package math

import (
	"github.com/Opioid/math32"
	"math"
)

func HemisphereSample_uniform(r, s float32)Vector3 {
	r1 := s * 2 * math.Pi
	r2 := 1 - r
	sp := math32.Sqrt(1 - r2 * r2)

	return MakeVector3(Cos(r1) * sp, Sin(r1) * sp, r2)
}

func HemisphereSample_cos(r, s float32) Vector3 {
	r1 := s * 2 * math.Pi
	r2 := math32.Sqrt(1 - r)
	sp := math32.Sqrt(1 - r2 * r2)

	return MakeVector3(Cos(r1) * sp, Sin(r1) * sp, r2)
}

func DiskSample_uniform(r, s float32)Vector3 {
	r1 := s * 2 * math.Pi
	r2 := 1 - r
	sp := math32.Sqrt(1 - r2 * r2)

	return MakeVector3(Cos(r1) * sp, Sin(r1) * sp, 0.0)
}