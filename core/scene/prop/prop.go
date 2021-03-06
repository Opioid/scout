package prop

import (
	"github.com/Opioid/scout/core/scene/entity"
	"github.com/Opioid/scout/core/scene/shape"
	"github.com/Opioid/scout/core/scene/shape/geometry"
	"github.com/Opioid/scout/core/scene/material"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/bounding"
	_ "fmt"
)

type Prop struct {
	entity.Entity

	Shape shape.Shape

	AABB bounding.AABB	

	Materials []material.Material

	CastsShadow bool

	visibility uint8
}

func NewProp() *Prop {
	p := Prop{visibility: Primary | Secondary}
//	a.transformation.ObjectToWorld.SetIdentity()
	p.CastsShadow = true

	return &p
}

func (p *Prop) SetTransformation(position, scale math.Vector3, rotation math.Quaternion) {
	p.Transformation.Set(position, scale, rotation)

	if p.Shape != nil {
		p.Shape.AABB().Transform(&p.Transformation.ObjectToWorld, &p.AABB)
	}
}

func (p *Prop) Intersect(ray *math.OptimizedRay, scratch *ScratchBuffer, intersection *geometry.Intersection) bool {
	var hit bool
	var boundingMinT, boundingMaxT float32

	if p.Shape.IsComplex() {
		hit, boundingMinT, boundingMaxT = p.AABB.Intersect(ray) 
		if !hit {
			return false
		}
	}

/*
	p.TransformationAt(ray.Time, &scratch.Transformation)

	if hit, thit := p.Shape.Intersect(&scratch.Transformation, ray, &scratch.Ray, boundingMinT, boundingMaxT, intersection); hit {
		ray.MaxT = thit

		return true		
	}
	*/

	if hit, thit := p.Shape.Intersect(&p.Transformation, ray, &scratch.Ray, boundingMinT, boundingMaxT, intersection); hit {
		ray.MaxT = thit

		return true		
	}

	return false
}

func (p *Prop) IntersectP(ray *math.OptimizedRay, scratch *ScratchBuffer) bool {
	var hit bool
	var boundingMinT, boundingMaxT float32

	if p.Shape.IsComplex() {
		hit, boundingMinT, boundingMaxT = p.AABB.Intersect(ray) 
		if !hit {
			return false
		}
	}

	p.TransformationAt(ray.Time, &scratch.Transformation)

	return p.Shape.IntersectP(&scratch.Transformation, ray, &scratch.Ray, boundingMinT, boundingMaxT) 
}

func (p *Prop) IsVisible(flags uint8) bool {
	return p.visibility & flags != 0
}

func (p *Prop) SetVisible(flag uint8, visible bool) {
	if visible {
		p.visibility |= flag
	} else {
		p.visibility &= ^flag
	}
}