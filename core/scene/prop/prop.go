package prop

import (
	"github.com/Opioid/scout/core/scene/entity"
	"github.com/Opioid/scout/core/scene/shape"
	"github.com/Opioid/scout/core/scene/shape/geometry"
	"github.com/Opioid/scout/core/scene/material"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/bounding"
)

type Prop struct {
	entity.Entity

	Shape shape.Shape

	AABB bounding.AABB	

	Materials []material.Material

	CastsShadow bool	
}

func NewProp() *Prop {
	p := new(Prop)
//	a.transformation.ObjectToWorld.SetIdentity()
	p.CastsShadow = true
	return p
}

func (p *Prop) SetTransformation(position, scale math.Vector3, rotation math.Quaternion) {
	p.Transformation.Set(position, scale, rotation)

	if p.Shape != nil {
		p.Shape.AABB().Transform(&p.Transformation.ObjectToWorld, &p.AABB)
	}
}

func (p *Prop) Intersect(ray *math.OptimizedRay, transformation *math.ComposedTransformation, intersection *geometry.Intersection) bool {
	var hit bool
	var boundingMinT, boundingMaxT float32

	if p.Shape.IsComplex() {
		hit, boundingMinT, boundingMaxT = p.AABB.Intersect(ray) 
		if !hit {
			return false
		}
	}

	p.TransformationAt(ray.Time, transformation)

	if hit, thit := p.Shape.Intersect(transformation, ray, boundingMinT, boundingMaxT, intersection); hit {
		ray.MaxT = thit

		return true		
	}
	
	return false
}

func (p *Prop) IntersectP(ray *math.OptimizedRay, transformation *math.ComposedTransformation) bool {
	var hit bool
	var boundingMinT, boundingMaxT float32

	if p.Shape.IsComplex() {
		hit, boundingMinT, boundingMaxT = p.AABB.Intersect(ray) 
		if !hit {
			return false
		}
	}

	p.TransformationAt(ray.Time, transformation)

	return p.Shape.IntersectP(transformation, ray, boundingMinT, boundingMaxT) 
}