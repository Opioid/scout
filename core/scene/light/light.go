package light

import (
	"github.com/Opioid/scout/core/scene/entity"
	"github.com/Opioid/scout/base/math"
)

type Type int

const (
	Directional Type = iota
	Point
	Spot
)

type Light struct {
	entity.Entity
	Type Type
	Color math.Vector3
}