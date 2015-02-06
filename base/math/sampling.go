package math

import (
	"github.com/Opioid/math32"
	"math"
)

func SampleDiskUniform(u, v float32) (float32, float32) {
	r := math32.Sqrt(u)
	theta := v * 2.0 * math.Pi
	
	sintheta, costheta := Sincos(theta)

	return costheta * r, sintheta * r
}

func SampleDiskConcentric(u, v float32) (float32, float32) {
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

func SampleHemisphereUniform(u, v float32) Vector3 {
	z := 1.0 - u
	r := math32.Sqrt(1.0 - z * z)
	phi := v * 2.0 * math.Pi

	sinphi, cosphi := Sincos(phi)

	return MakeVector3(cosphi * r, sinphi * r, z)
}

/*
func SampleHemisphere_cos(u, v float32) Vector3 {
	z := math32.Sqrt(1.0 - u)
	r := math32.Sqrt(1.0 - z * z)
	phi := v * 2.0 * math.Pi

	sinphi, cosphi := Sincos(phi)

	return MakeVector3(cosphi * r, sinphi * r, z)
}
*/

func SampleHemisphereCosine(u, v float32) Vector3 {
	x, y := SampleDiskConcentric(u, v)
	z := math32.Sqrt(math32.Max(0.0, 1.0 - x * x - y * y))

	return MakeVector3(x, y, z)
}

func SampleConeUniform(u, v, costhetamax float32) Vector3 {
	/*	
	float costheta = (1.f - u1) + u1 * costhetamax;
	float sintheta = sqrtf(1.f - costheta*costheta);
	float phi = u2 * 2.f * M_PI;
	return Vector(cosf(phi) * sintheta, sinf(phi) * sintheta, costheta);
	*/

	costheta := (1.0 - u) + u * costhetamax
	sintheta := math32.Sqrt(1.0 - costheta * costheta)
	phi := v * 2.0 * math32.Pi

	sinphi, cosphi := Sincos(phi)

	return MakeVector3(cosphi * sintheta, sinphi * sintheta, costheta)
}

func SampleOrientedConeUniform(u, v, costhetamax float32, x, y, z Vector3) Vector3 {
	/*
	float costheta = (1.f - u1) + u1 * costhetamax;
	float sintheta = sqrtf(1.f - costheta*costheta);
	float phi = u2 * 2.f * M_PI;
	return cosf(phi) * sintheta * x + sinf(phi) * sintheta * y + costheta * z;
	*/

	costheta := (1.0 - u) + u * costhetamax
	sintheta := math32.Sqrt(1.0 - costheta * costheta)
	phi := v * 2.0 * math32.Pi

	sinphi, cosphi := Sincos(phi)

	return x.Scale(cosphi * sintheta).Add(y.Scale(sinphi * sintheta)).Add(z.Scale(costheta))
}

func ConePdfUniform(costhetamax float32) float32 {
	return 1.0 / (2.0 * math32.Pi * (1.0 - costhetamax))
}