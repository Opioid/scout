package shape

import (
	"github.com/Opioid/scout/core/scene/shape/geometry"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/bounding"
	"github.com/Opioid/math32"
	_ "fmt"
)

type Sphere struct {
	aabb bounding.AABB
}

func NewSphere() *Sphere {
	s := new(Sphere)
	s.aabb = bounding.MakeAABB(math.MakeVector3(-1.0, -1.0, -1.0), math.MakeVector3(1.0, 1.0, 1.0))
	return s
}

// Won't work from the inside!
func (s *Sphere) Intersect(transformation *math.ComposedTransformation, ray, tray *math.OptimizedRay,
						   boundingMinT, boundingMaxT float32, intersection *geometry.Intersection) (bool, float32) {
	v := ray.Origin.Sub(transformation.Position)
	b := -v.Dot(ray.Direction)
	radius := transformation.Scale.X
	det := (b * b) - v.Dot(v) + (radius * radius)

	if det > 0.0 {
		/*
		thit := b - math32.Sqrt(det)

		if thit > ray.MinT && thit < ray.MaxT {
			intersection.Epsilon = 5e-4 * thit

			intersection.P = ray.Point(thit)
			intersection.N = intersection.P.Sub(transformation.Position).Div(radius)

			t, b := math.CoordinateSystem(intersection.N)
			intersection.T = t
			intersection.B = b

			intersection.MaterialIndex = 0

			return true, thit
		}
		*/
	
		dist := math32.Sqrt(det)
		t0 := b - dist

		if t0 > ray.MinT && t0 < ray.MaxT {
			intersection.Epsilon = 5e-4 * t0

			intersection.P = ray.Point(t0)
			intersection.N = intersection.P.Sub(transformation.Position).Div(radius)

			t, b := math.CoordinateSystem(intersection.N)
			intersection.T = t
			intersection.B = b

			intersection.MaterialIndex = 0

			return true, t0
		}

		t1 := b + dist
		
		if t1 > ray.MinT && t1 < ray.MaxT {
			intersection.Epsilon = 5e-4 * t1

			intersection.P = ray.Point(t1)
			intersection.N = intersection.P.Sub(transformation.Position).Div(radius)

			t, b := math.CoordinateSystem(intersection.N)
			intersection.T = t
			intersection.B = b

			intersection.MaterialIndex = 0

			return true, t1
		}
		
	}

	return false, 0.0
}

func (s *Sphere) IntersectP(transformation *math.ComposedTransformation, ray, tray *math.OptimizedRay,
						    boundingMinT, boundingMaxT float32) bool {
	v := ray.Origin.Sub(transformation.Position)
	b := -v.Dot(ray.Direction)
	radius := transformation.Scale.X
	det := (b * b) - v.Dot(v) + (radius * radius)

	if det > 0.0 {
		dist := math32.Sqrt(det)
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