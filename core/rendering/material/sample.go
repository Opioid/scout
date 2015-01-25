package material

import (
	"github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/base/math"
)

type Sample interface {
	Evaluate(l math.Vector3) math.Vector3
	Values() *Values

	MonteCarloBxdf(subsample uint32, sampler sampler.Sampler) (Bxdf, float32)
}