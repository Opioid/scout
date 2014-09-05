package prop

import (
	"github.com/Opioid/scout/core/scene/shape/geometry"
)

type Intersection struct {
	Prop *Prop
	Dg geometry.Differential
	Epsilon float32
}