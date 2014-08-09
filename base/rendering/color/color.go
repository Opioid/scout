package color

import (
	"github.com/Opioid/scout/base/math"
)

// convert sRGB linear value to sRGB gamma value
func LinearToSrgb(c float32) float32 {
	if c <= 0.0 {
		return 0.0
	} else if c <  0.0031308 {
		return 12.92 * c
	} else if c < 1.0 {
		return 1.055 * math.Pow(c, 0.41666) - 0.055
	} else {
		return 1.0
	}
}