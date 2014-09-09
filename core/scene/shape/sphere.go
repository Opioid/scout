package shape

import (
	"github.com/Opioid/scout/core/scene/entity"
	"github.com/Opioid/scout/core/scene/shape/geometry"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/bounding"
)

type Sphere struct {
	aabb bounding.AABB
}

func NewSphere() *Sphere {
	s := new(Sphere)
	s.aabb = bounding.MakeAABB(math.Vector3{-1.0, -1.0, -1.0}, math.Vector3{1.0, 1.0, 1.0})
	return s
}

// Won't work from the inside!
func (s *Sphere) Intersect(transformation *entity.ComposedTransformation, ray *math.OptimizedRay, thit *float32, epsilon *float32, dg *geometry.Differential) bool {
	v := ray.Origin.Sub(transformation.Position)
	b := -v.Dot(ray.Direction)
	radius := transformation.Scale.X
	det := (b * b) - v.Dot(v) + (radius * radius)

	if det > 0.0 {
		*thit = b - math.Sqrt(det)

		if *thit > ray.MinT && *thit < ray.MaxT {
			*epsilon = 5e-4 * *thit

			dg.P = ray.Point(*thit)
			dg.N = dg.P.Sub(transformation.Position).Div(radius)

			return true
		} 
	}

	return false
}

func (s *Sphere) IntersectP(transformation *entity.ComposedTransformation, ray *math.OptimizedRay) bool {
	v := ray.Origin.Sub(transformation.Position)
	b := -v.Dot(ray.Direction)
	radius := transformation.Scale.X
	det := (b * b) - v.Dot(v) + (radius * radius)

	if det > 0.0 {
		dist := math.Sqrt(det)
		t0 := b - dist
		t1 := b + dist

		if t1 > ray.MinT && t0 < ray.MaxT {
			return true
		}

		if t0 > ray.MinT && t1 < ray.MaxT {
			return true
		}

	}

	return false	
}

func (s *Sphere) AABB() *bounding.AABB {
	return &s.aabb
}

func (s *Sphere) IsComplex() bool {
	return false
}

func (s *Sphere) IsFinite() bool {
	return true
}