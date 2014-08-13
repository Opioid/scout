package shape

import (
	"github.com/Opioid/scout/core/scene/entity"
	"github.com/Opioid/scout/base/math"
)

type plane struct {
	
}

// works from both sides of the plane
func (p *plane) Intersect(transformation *entity.ComposedTransformation, ray *math.Ray, thit *float32, epsilon *float32, dg *DifferentialGeometry) bool {
	normal := transformation.Rotation.Row(2)

	d := -normal.Dot(transformation.Position)

	denom := normal.Dot(ray.Direction)

	numer := normal.Dot(ray.Origin) + d

	*thit = -(numer / denom)
	
	if *thit >= ray.MinT && *thit < ray.MaxT {
		*epsilon = 5e-4 * *thit

		dg.P = ray.Point(*thit)
		dg.Nn = normal

		u := transformation.ObjectToWorld.Row(0).Vector3().Dot(dg.P)
		dg.UV.X = u - math.Floor(u)

		v := transformation.ObjectToWorld.Row(1).Vector3().Dot(dg.P)
		dg.UV.Y = v - math.Floor(v)

		return true
	} 

	return false
}

// works from both sides of the plane
func (p *plane) IntersectP(transformation *entity.ComposedTransformation, ray *math.Ray) bool {
	normal := transformation.Rotation.Row(2)

	d := -normal.Dot(transformation.Position)

	denom := normal.Dot(ray.Direction)

	numer := normal.Dot(ray.Origin) + d

	thit := -(numer / denom)
	
	if thit >= ray.MinT && thit < ray.MaxT {
		return true
	} 

	return false
}