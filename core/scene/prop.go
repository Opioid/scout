package scene

import (
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/core/scene/shape"
)

type Prop struct {
	Entity
	shape shape.Shape
}

func NewProp(shape shape.Shape) *Prop {
	p := new(Prop)
	p.shape = shape
	return p
}

func (p *Prop) Shape() shape.Shape {
	return p.shape
}

func (p *Prop) Intersect(ray *math.Ray, thit *float32) bool {
	return p.shape.Intersect(&p.Transformation, ray, thit)
}
