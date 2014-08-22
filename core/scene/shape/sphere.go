package shape

import (
	"github.com/Opioid/scout/core/scene/entity"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/bounding"
)

type Sphere struct {
	aabb bounding.AABB
}

func NewSphere() *Sphere {
	s := new(Sphere)
	s.aabb = bounding.AABB{math.Vector3{-1.0, -1.0, -1.0}, math.Vector3{1.0, 1.0, 1.0}}
	return s
}

// Won't work from the inside!
func (s *Sphere) Intersect(transformation *entity.ComposedTransformation, ray *math.Ray, thit *float32, epsilon *float32, dg *DifferentialGeometry) bool {
	v := ray.Origin.Sub(transformation.Position)
	b := -v.Dot(ray.Direction)
	radius := transformation.Scale.X
	det := (b * b) - v.Dot(v) + (radius * radius)

	if det > 0.0 {
		*thit = b - math.Sqrt(det)

		if *thit >= ray.MinT && *thit < ray.MaxT {
			*epsilon = 5e-4 * *thit

			dg.P = ray.Point(*thit)
			dg.Nn = dg.P.Sub(transformation.Position).Div(radius)

			return true
		} 
	}

	return false
}

// Won't work from the inside!
func (s *Sphere) IntersectP(transformation *entity.ComposedTransformation, ray *math.Ray) bool {
	v := ray.Origin.Sub(transformation.Position)
	b := -v.Dot(ray.Direction)
	radius := transformation.Scale.X
	det := (b * b) - v.Dot(v) + (radius * radius)

	if det > 0.0 {
		thit := b - math.Sqrt(det)

		if thit >= ray.MinT && thit < ray.MaxT {
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