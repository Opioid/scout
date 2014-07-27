package shape

import (
	"github.com/Opioid/scout/core/scene/entity"
	"github.com/Opioid/scout/base/math"
)

type Shape interface {
	Intersect(transformation *entity.ComposedTransformation, ray *math.Ray, thit *float32) bool

	IntersectP(transformation *entity.ComposedTransformation, ray *math.Ray, thit *float32) bool
}