package shape

import (
	"github.com/Opioid/scout/core/scene/entity"
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
func (p *Plane) Intersect(transformation *entity.ComposedTransformation, ray *math.OptimizedRay, boundingMinT, boundingMaxT float32, 
						  thit *float32, epsilon *float32, dg *geometry.Differential) bool {
	normal := transformation.Rotation.Row(2)

	d := -normal.Dot(transformation.Position)

	denom := normal.Dot(ray.Direction)

	numer := normal.Dot(ray.Origin) + d

	*thit = -(numer / denom)
	
	if *thit > ray.MinT && *thit < ray.MaxT {
		*epsilon = 5e-4 * *thit

		dg.P = ray.Point(*thit)
		dg.T = transformation.Rotation.Row(0)
		dg.B = transformation.Rotation.Row(1)
		dg.N = normal

		u := transformation.ObjectToWorld.Row(0).Vector3().Dot(dg.P)
		dg.UV.X = u - math32.Floor(u)

		v := transformation.ObjectToWorld.Row(1).Vector3().Dot(dg.P)
		dg.UV.Y = v - math32.Floor(v)

		return true
	} 

	return false
}

// works from both sides of the plane
func (p *Plane) IntersectP(transformation *entity.ComposedTransformation, ray *math.OptimizedRay, boundingMinT, boundingMaxT float32) bool {
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