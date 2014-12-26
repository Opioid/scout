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
//	p.transformation.ObjectToWorld.SetIdentity()
	return a
}

func (a *Actor) Intersect(ray *math.OptimizedRay, intersection *Intersection) bool {
	var boundingMinT, boundingMaxT float32

	if a.Shape.IsComplex() && !a.AABB.Intersect(ray, &boundingMinT, &boundingMaxT) {
		return false
	}

	transformation := a.TransformationAt(ray.Time)

	var thit, epsilon float32
	
	if !a.Shape.Intersect(&transformation, ray, boundingMinT, boundingMaxT, &thit, &epsilon, &intersection.Dg) {
		return false
	}

	intersection.Epsilon = epsilon
	ray.MaxT = thit

	return true
}