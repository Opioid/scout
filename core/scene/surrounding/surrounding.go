package surrounding

import (
	"github.com/Opioid/scout/base/math"
)

type Surrounding interface {
	Sample(ray *math.OptimizedRay) math.Vector3
}