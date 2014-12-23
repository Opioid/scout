package sampler

import (
	"github.com/Opioid/scout/base/math"
)

type CameraSample struct {
	Coordinates math.Vector2
	LensUv math.Vector2
	RelativeOffset math.Vector2
	Time float32
}