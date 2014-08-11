package scene

import (
	"github.com/Opioid/scout/core/scene/entity"
	"github.com/Opioid/scout/core/scene/shape"
	"github.com/Opioid/scout/base/math"
	_ "fmt"
)

type StaticProp struct {
	Prop
	Transformation entity.ComposedTransformation
}

func NewStaticProp(shape shape.Shape, material Material) *StaticProp {
	p := new(StaticProp)
	p.Shape = shape
	p.Material = material
	p.Transformation.ObjectToWorld.SetIdentity()
	return p
}

func (p *StaticProp) Intersect(ray *math.Ray, intersection *Intersection) bool {
	var thit, epsilon float32
	
	if !p.Shape.Intersect(&p.Transformation, ray, &thit, &epsilon, &intersection.Dg) {
		return false
	}

	intersection.Epsilon = epsilon
	ray.MaxT = thit

	return true
}

func (p *StaticProp) IntersectP(ray *math.Ray) bool {
	return p.Shape.IntersectP(&p.Transformation, ray) 
}