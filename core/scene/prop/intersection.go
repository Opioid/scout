package prop

import (
	"github.com/Opioid/scout/core/scene/shape/geometry"
	"github.com/Opioid/scout/core/scene/material"
)

type Intersection struct {
	Prop *Prop
	Dg geometry.Differential
	Epsilon float32
}

func (i *Intersection) Material() material.Material {
	return i.Prop.Materials[i.Dg.MaterialId]
}