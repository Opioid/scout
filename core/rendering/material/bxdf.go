package material

import (
	"github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/base/math"
)

type Bxdf interface {
	ImportanceSample(subsample uint32, sampler sampler.Sampler) (math.Vector3, math.Vector3, float32, float32)

	Evaluate(wi math.Vector3, NdotWi float32) (math.Vector3, float32)
}