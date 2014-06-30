package camera

import (
	"github.com/Opioid/scout/core/scene"
	"github.com/Opioid/scout/base/math"
)

type Sample struct {
	Coordinates math.Vector2i
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