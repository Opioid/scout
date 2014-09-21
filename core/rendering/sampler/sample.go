package sampler

import (
	"github.com/Opioid/scout/base/math"
)

type Sample struct {
	Coordinates math.Vector2
}

func NewSample(x, y float32) *Sample {
	return &Sample{math.MakeVector2(x, y)}
}