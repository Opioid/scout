package sampler

import (
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/random"
)

type Sampler interface {
	Clone(rng *random.Generator) Sampler

	NumSamplesPerIteration() uint32

	Restart(numIterations uint32)

	GenerateNewSample(offset math.Vector2, sample *CameraSample) bool

	GenerateSamples(iteration uint32) []math.Vector2
}