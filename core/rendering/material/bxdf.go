package material

import (
	"github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/base/math"
)

type Bxdf interface {
	ImportanceSample(subsample uint32, sampler sampler.Sampler) math.Vector3

	Evaluate(l math.Vector3) math.Vector3
}