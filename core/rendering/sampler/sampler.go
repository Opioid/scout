package sampler

import (
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/random"
)

type Sampler interface {
	Clone(rng *random.Generator) Sampler

	NumSamplesPerIteration() uint32

	Restart(numIterations uint32)

	GenerateCameraSample(offset math.Vector2, sample *CameraSample) bool

	GenerateSamples(iteration uint32, buffer []math.Vector2) []math.Vector2

	GenerateSample2D(index, iteration uint32) math.Vector2

	GenerateSample1D(index, iteration uint32) float32
}

type sampler struct {
	currentSample uint32
	numSamplesPerIteration uint32	
}

func (s *sampler) NumSamplesPerIteration() uint32 {
	return s.numSamplesPerIteration
}