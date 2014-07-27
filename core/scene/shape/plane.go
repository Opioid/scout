package shape

import (
	"github.com/Opioid/scout/core/scene/entity"
	"github.com/Opioid/scout/base/math"
)

type Plane struct {
	
}

// works from both sides of the plane
func (plane *Plane) Intersect(transformation *entity.ComposedTransformation, ray *math.Ray, thit *float32, epsilon *float32, dg *DifferentialGeometry) bool {
	normal := transformation.Matrix.Row(2).Vector3()

	d := -normal.Dot(transformation.Position)

	denom := normal.Dot(ray.Direction)

	numer := normal.Dot(ray.Origin) + d

	*thit = -(numer / denom)
	
	if *thit >= ray.MinT && *thit < ray.MaxT {
		*epsilon = 5e-4 * *thit

		dg.P = ray.Point(*thit)
		dg.Nn = normal

		return true
	} 

	return false
}

// works from both sides of the plane
func (plane *Plane) IntersectP(transformation *entity.ComposedTransformation, ray *math.Ray) bool {
	normal := transformation.Matrix.Row(2).Vector3()

	d := -normal.Dot(transformation.Position)

	denom := normal.Dot(ray.Direction)

	numer := normal.Dot(ray.Origin) + d

	thit := -(numer / denom)
	
	if thit >= ray.MinT && thit < ray.MaxT {
		return true
	} 

	return false
}