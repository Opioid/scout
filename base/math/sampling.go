package math

import (
	"github.com/Opioid/math32"
	"math"
)

func SampleHemisphere_uniform(u, v float32) Vector3 {
	z := 1.0 - u
	r := math32.Sqrt(1.0 - z * z)
	phi := v * 2.0 * math.Pi

	sinphi, cosphi := Sincos(phi)

	return MakeVector3(cosphi * r, sinphi * r, z)
}

func SampleHemisphere_cos(u, v float32) Vector3 {
	z := math32.Sqrt(1.0 - u)
	r := math32.Sqrt(1.0 - z * z)
	phi := v * 2.0 * math.Pi

	sinphi, cosphi := Sincos(phi)

	return MakeVector3(cosphi * r, sinphi * r, z)
}

func SampleDisk_uniform(u, v float32) Vector3 {
	r := math32.Sqrt(u)
	theta := v * 2.0 * math.Pi
	
	sintheta, costheta := Sincos(theta)

	return MakeVector3(costheta * r, sintheta * r, 0.0)
}