package shape

import (
	"github.com/Opioid/scout/core/scene/entity"
	"github.com/Opioid/scout/base/math"
)

type Plane struct {
	
}

// works from both sides of the plane
func (plane *Plane) Intersect(transformation *entity.ComposedTransformation, ray *math.Ray, thit *float32) bool {
	normal := transformation.Matrix.Row(2).Vector3()// math.Vector3{0.0, 0.0, -1.0}

	d := -normal.Dot(transformation.Position)

	denom := normal.Dot(ray.Direction)

	numer := normal.Dot(ray.Origin) + d

	*thit = -(numer / denom)
	
	if *thit >= 0.0 && *thit < ray.MaxT {
		return true
	} 

	return false
}

// works from both sides of the plane
func (plane *Plane) IntersectP(transformation *entity.ComposedTransformation, ray *math.Ray, thit *float32) bool {
	normal := transformation.Matrix.Row(2).Vector3()// math.Vector3{0.0, 0.0, -1.0}

	d := -normal.Dot(transformation.Position)

	denom := normal.Dot(ray.Direction)

	numer := normal.Dot(ray.Origin) + d

	*thit = -(numer / denom)
	
	if *thit >= 0.0 && *thit < ray.MaxT {
		return true
	} 

	return false
}