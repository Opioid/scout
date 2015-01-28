package prop

import (
	"github.com/Opioid/scout/base/math"
)

type ScratchBuffer struct {
	Transformation math.ComposedTransformation
	Ray math.OptimizedRay
}