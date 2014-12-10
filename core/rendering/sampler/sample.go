package sampler

import (
	"github.com/Opioid/scout/base/math"
)

type Sample struct {
	Coordinates math.Vector2
	RelativeOffset math.Vector2
	Id uint32
}