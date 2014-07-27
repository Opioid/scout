package scene

import (
	"github.com/Opioid/scout/core/scene/entity"
	"github.com/Opioid/scout/core/scene/shape"
	"github.com/Opioid/scout/base/math"
	_ "fmt"
)

type StaticProp struct {
	Prop
	WorldTransformation entity.ComposedTransformation
}

func NewStaticProp(shape shape.Shape, material *Material) *StaticProp {
	p := new(StaticProp)
	p.Shape = shape
	p.Material = material
	return p
}

func (p *StaticProp) Intersect(ray *math.Ray, intersection *Intersection) bool {
	var thit, epsilon float32
	
	if !p.Shape.Intersect(&p.WorldTransformation, ray, &thit, &epsilon, &intersection.Dg) {
		return false
	}

	intersection.Epsilon = epsilon
	ray.MaxT = thit

	return true
}

func (p *StaticProp) IntersectP(ray *math.Ray) bool {
	return p.Shape.IntersectP(&p.WorldTransformation, ray) 
}

func (p *StaticProp) WorldPosition() math.Vector3 {
	return p.WorldTransformation.Matrix.Row(4).Vector3()
}

func (p *StaticProp) SetWorldTransformation(position, scale math.Vector3, rotation *math.Matrix3x3) {
	p.WorldTransformation.Position = position
	p.WorldTransformation.Scale = scale

	p.WorldTransformation.Matrix.SetBasis(rotation)
	p.WorldTransformation.Matrix.SetOrigin(position)
}