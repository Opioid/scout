package scene

import (
	"github.com/Opioid/scout/base/math"
)

type Intersection struct {
	Prop *Prop
	Dg DifferentialGeometry
}

type DifferentialGeometry struct {
	Point math.Vector3
}