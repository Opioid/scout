package scene

import (
	"github.com/Opioid/scout/core/scene/entity"
	"github.com/Opioid/scout/core/scene/shape"
	"github.com/Opioid/scout/base/math"
	_ "github.com/Opioid/scout/base/math/bounding"
	_ "fmt"
)

type StaticProp struct {
	Prop
	transformation entity.ComposedTransformation
}

func NewStaticProp(shape shape.Shape, material Material) *StaticProp {
	p := new(StaticProp)
	p.Shape = shape
	p.Material = material
	p.transformation.ObjectToWorld.SetIdentity()
	return p
}

func (p *StaticProp) Intersect(ray *math.Ray, intersection *Intersection) bool {
	if p.Shape.IsComplex() && !p.AABB.Intersect(ray) {
		return false
	}

	var thit, epsilon float32
	
	if !p.Shape.Intersect(&p.transformation, ray, &thit, &epsilon, &intersection.Dg) {
		return false
	}

	intersection.Epsilon = epsilon
	ray.MaxT = thit

	return true
}

func (p *StaticProp) IntersectP(ray *math.Ray) bool {
	if p.Shape.IsComplex() && !p.AABB.Intersect(ray) {
		return false
	}

	return p.Shape.IntersectP(&p.transformation, ray) 
}

func (p *StaticProp) SetTransformation(position, scale math.Vector3, rotation math.Quaternion) {
	p.transformation.Set(position, scale, rotation)

	p.Shape.AABB().Transform(&p.transformation.ObjectToWorld, &p.AABB)
}

func (p *StaticProp) Transformation() *entity.ComposedTransformation {
	return &p.transformation
}