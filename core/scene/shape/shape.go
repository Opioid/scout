package shape

import (
	"github.com/Opioid/scout/core/scene/entity"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/bounding"
)

type Shape interface {
	Intersect(transformation *entity.ComposedTransformation, ray *math.OptimizedRay, thit *float32, epsilon *float32, dg *DifferentialGeometry) bool

	IntersectP(transformation *entity.ComposedTransformation, ray *math.OptimizedRay) bool

	AABB() *bounding.AABB

	IsComplex() bool
	IsFinite() bool
}