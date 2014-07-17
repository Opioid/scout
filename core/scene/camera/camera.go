package camera

import (
	"github.com/Opioid/scout/core/scene"
	"github.com/Opioid/scout/base/math"
)

type Sample struct {
	Coordinates math.Vector2
}

func NewSample(x, y float32) *Sample {
	return &Sample{math.Vector2{x, y}}
}

type Camera interface {
	Position() math.Vector3
	Film() Film
	GenerateRay(sample *Sample, ray *math.Ray)
}

type genericCamera struct {
	scene.Entity
	film Film
}