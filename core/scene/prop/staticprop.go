package prop

import (
	"github.com/Opioid/scout/core/scene/entity"
	"github.com/Opioid/scout/core/scene/shape/geometry"
	"github.com/Opioid/scout/base/math"
	_ "github.com/Opioid/scout/base/math/bounding"
	_ "fmt"
)

type StaticProp struct {
	Prop
	transformation entity.ComposedTransformation
}

func NewStaticProp() *StaticProp {
	p := new(StaticProp)
	p.transformation.ObjectToWorld.SetIdentity()
	return p
}

func (p *StaticProp) Intersect(ray *math.OptimizedRay, intersection *geometry.Intersection) bool {
	var hit bool
	var boundingMinT, boundingMaxT float32

	if p.Shape.IsComplex() {
		hit, boundingMinT, boundingMaxT = p.AABB.Intersect(ray) 
		if !hit {
			return false
		}
	}

	if hit, thit := p.Shape.Intersect(&p.transformation, ray, boundingMinT, boundingMaxT, intersection); hit {
		ray.MaxT = thit

		return true
	}
	
	return false
}

func (p *StaticProp) IntersectP(ray *math.OptimizedRay) bool {
	var hit bool
	var boundingMinT, boundingMaxT float32

	if p.Shape.IsComplex() {
		hit, boundingMinT, boundingMaxT = p.AABB.Intersect(ray) 
		if !hit {
			return false
		}
	}

	return p.Shape.IntersectP(&p.transformation, ray, boundingMinT, boundingMaxT) 
}

func (p *StaticProp) SetTransformation(position, scale math.Vector3, rotation math.Quaternion) {
	p.transformation.Set(position, scale, rotation)

	p.Shape.AABB().Transform(&p.transformation.ObjectToWorld, &p.AABB)
}

func (p *StaticProp) Transformation() *entity.ComposedTransformation {
	return &p.transformation
}