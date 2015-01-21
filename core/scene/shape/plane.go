package shape

import (
	"github.com/Opioid/scout/core/scene/shape/geometry"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/bounding"
	"github.com/Opioid/math32"
)

type Plane struct {
	aabb bounding.AABB
}

func NewPlane() *Plane {
	p := new(Plane)
	p.aabb = bounding.MakeAABB(math.MakeVector3(0.0, 0.0, 0.0), math.MakeVector3(0.0, 0.0, 0.0))
	return p
}

// works from both sides of the plane
func (p *Plane) Intersect(transformation *math.ComposedTransformation, ray *math.OptimizedRay, boundingMinT, boundingMaxT float32, 
						  intersection *geometry.Intersection) (bool, float32) {
	normal := transformation.Rotation.Row(2)

	d := -normal.Dot(transformation.Position)

	denom := normal.Dot(ray.Direction)

	numer := normal.Dot(ray.Origin) + d

	thit := -(numer / denom)
	
	if thit > ray.MinT && thit < ray.MaxT {
		intersection.Epsilon = 5e-4 * thit

		intersection.P = ray.Point(thit)
		intersection.T = transformation.Rotation.Row(0)
		intersection.B = transformation.Rotation.Row(1)
		intersection.N = normal

		u := transformation.ObjectToWorld.Row(0).Vector3().Dot(intersection.P)
		intersection.UV.X = u - math32.Floor(u)

		v := transformation.ObjectToWorld.Row(1).Vector3().Dot(intersection.P)
		intersection.UV.Y = v - math32.Floor(v)

		intersection.MaterialId = 0

		return true, thit
	} 

	return false, 0.0
}

// works from both sides of the plane
func (p *Plane) IntersectP(transformation *math.ComposedTransformation, ray *math.OptimizedRay, boundingMinT, boundingMaxT float32) bool {
	normal := transformation.Rotation.Row(2)

	d := -normal.Dot(transformation.Position)

	denom := normal.Dot(ray.Direction)

	numer := normal.Dot(ray.Origin) + d

	thit := -(numer / denom)
	
	if thit > ray.MinT && thit < ray.MaxT {
		return true
	} 

	return false
}

func (p *Plane) AABB() *bounding.AABB {
	return &p.aabb
}

func (p *Plane) IsComplex() bool {
	return false
}

func (p *Plane) IsFinite() bool {
	return false
}