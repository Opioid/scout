package math

import (
	"github.com/Opioid/math32"
	"math"
)

func HemisphereSample_uniform(u, v float32)Vector3 {
	z := 1.0 - u
	r := math32.Sqrt(1.0 - z * z)
	phi := v * 2.0 * math.Pi

	return MakeVector3(Cos(phi) * r, Sin(phi) * r, z)
}

func HemisphereSample_cos(u, v float32) Vector3 {
	z := math32.Sqrt(1.0 - u)
	r := math32.Sqrt(1.0 - z * z)
	phi := v * 2.0 * math.Pi

	return MakeVector3(Cos(phi) * r, Sin(phi) * r, z)
}

func DiskSample_uniform(u, v float32) Vector3 {
	r := math32.Sqrt(u)
	theta := v * 2.0 * math.Pi
	
	return MakeVector3(Cos(theta) * r, Sin(theta) * r, 0.0)
}