package surrounding

import (
	"github.com/Opioid/scout/base/math"
)

type Surrounding interface {
	Sample(ray *math.OptimizedRay) math.Vector3

	SampleSecondary(ray *math.OptimizedRay) (math.Vector3, float32)

	SampleDiffuse(v math.Vector3) math.Vector3
	SampleSpecular(v math.Vector3, roughness float32) math.Vector3
}