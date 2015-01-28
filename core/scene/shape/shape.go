package shape

import (
	"github.com/Opioid/scout/core/scene/shape/geometry"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/bounding"
)

type Shape interface {
	Intersect(transformation *math.ComposedTransformation, ray, tray *math.OptimizedRay,
			  boundingMinT, boundingMaxT float32, intersection *geometry.Intersection) (bool, float32)

	IntersectP(transformation *math.ComposedTransformation, ray, tray *math.OptimizedRay,
			   boundingMinT, boundingMaxT float32) bool

	AABB() *bounding.AABB

	IsComplex() bool
	IsFinite() bool
}