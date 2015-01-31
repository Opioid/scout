package math

import (
	"github.com/Opioid/math32"
	"math"
)

func SampleDisk_uniform(u, v float32) (float32, float32) {
	r := math32.Sqrt(u)
	theta := v * 2.0 * math.Pi
	
	sintheta, costheta := Sincos(theta)

	return costheta * r, sintheta * r
}

func SampleDisk_concentric(u, v float32) (float32, float32) {
	sx := 2.0 * u - 1.0
	sy := 2.0 * v - 1.0

	var r, theta float32

	if sx >= -sy {
		if sx > sy {
			// handle first region of disk
			r = sx
			if sy > 0.0 {
				theta = sy / r
			} else {
				theta = 8.0 + sy / r
			}
		} else {
			// handle second region of disk
			r = sy
			theta = 2.0 - sx / r
		}
	} else {
		if sx <= sy {
			// handle third region of disk
			r = -sx
			theta = 4.0 - sy / r
		} else {
			// handle fourth region of disk
			r = -sy
			theta = 6.0 + sx / r
		}
	}

	theta *= math.Pi / 4.0

	sintheta, costheta := Sincos(theta)

	return costheta * r, sintheta * r
}

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

func SampleHemisphere_cos1(u, v float32) Vector3 {
	x, y := SampleDisk_concentric(u, v)
	z := math32.Sqrt(math32.Max(0.0, 1.0 - x * x - y * y))

	return MakeVector3(x, y, z)
}