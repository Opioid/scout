package shape

import (
	"github.com/Opioid/scout/base/math"
)

type Shape interface {
	Intersect(ray *math.Ray, thit *float32) bool
}