package prop

import (
	"github.com/Opioid/scout/core/scene/shape"
)

type Intersection struct {
	Prop *Prop
	Dg shape.DifferentialGeometry
	Epsilon float32
}