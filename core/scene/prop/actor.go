package prop

import (
	"github.com/Opioid/scout/core/scene/entity"
	"github.com/Opioid/scout/base/math"
)

type Actor struct {
	entity.Entity
	Prop
}

func NewActor() *Actor {
	a := new(Actor)
//	a.transformation.ObjectToWorld.SetIdentity()
	return a
}

func (a *Actor) Intersect(ray *math.OptimizedRay, intersection *Intersection) bool {
	var hit bool
	var boundingMinT, boundingMaxT float32

	if a.Shape.IsComplex() {
		hit, boundingMinT, boundingMaxT = a.AABB.Intersect(ray) 
		if !hit {
			return false
		}
	}

	transformation := a.TransformationAt(ray.Time)

	if hit, thit, epsilon := a.Shape.Intersect(&transformation, ray, boundingMinT, boundingMaxT, &intersection.Dg); hit {
		intersection.Epsilon = epsilon
		ray.MaxT = thit

		return true		
	}
	
	return false
}

func (a *Actor) IntersectP(ray *math.OptimizedRay) bool {
	var hit bool
	var boundingMinT, boundingMaxT float32

	if a.Shape.IsComplex() {
		hit, boundingMinT, boundingMaxT = a.AABB.Intersect(ray) 
		if !hit {
			return false
		}
	}

	transformation := a.TransformationAt(ray.Time)

	return a.Shape.IntersectP(&transformation, ray, boundingMinT, boundingMaxT) 
}