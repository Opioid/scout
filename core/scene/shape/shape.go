package shape

import (
	"github.com/Opioid/scout/base/math"
)

type Shape interface {
	Intersect(transformation *math.Transformation, ray *math.Ray, thit *float32) bool
}