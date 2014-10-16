package sampler

import (
	"github.com/Opioid/scout/base/math"
)

type Sample struct {
	Coordinates math.Vector2
	Id uint32
}

/*
func NewSample(x, y float32, id uint32) *Sample {
	return &Sample{math.MakeVector2(x, y), id}
}*/