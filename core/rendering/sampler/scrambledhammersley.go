package sampler

import (
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/random"
)

type ScrambledHammersley struct {
	numSamplesPerIteration uint32
	numTotalSamples uint32

	currentSample uint32

	randomBits uint32

	samples []math.Vector2

	rng *random.Generator
}

func NewScrambledHammersley(numSamplesPerIteration uint32, rng *random.Generator) *ScrambledHammersley {
	s := new(ScrambledHammersley)
	s.allocateSamples(numSamplesPerIteration)
	s.rng = rng
	return s
}

func (s *ScrambledHammersley) allocateSamples(numSamplesPerIteration uint32) {
	s.numSamplesPerIteration = numSamplesPerIteration
	s.samples = make([]math.Vector2, numSamplesPerIteration)
}

func (s *ScrambledHammersley) Clone(rng *random.Generator) Sampler {
	ns := new(ScrambledHammersley)
	ns.allocateSamples(s.numSamplesPerIteration)
	ns.rng = rng
	return ns
}

func (s *ScrambledHammersley) NumSamplesPerIteration() uint32 {
	return s.numSamplesPerIteration
}

func (s *ScrambledHammersley) Restart(numIterations uint32) {
	s.currentSample = 0
	s.numTotalSamples = s.numSamplesPerIteration * numIterations
	s.randomBits = s.rng.RandomUint32()
}

func (s *ScrambledHammersley) GenerateNewSample(sample *math.Vector2) bool {
	if s.currentSample >= s.numSamplesPerIteration {
		return false
	}

//	*sample = math.ScrambledHammersley(s.currentSample + iteration * s.numSamplesPerIteration, s.numTotalSamples, s.randomBits)
	*sample = math.ScrambledHammersley(s.currentSample, s.numSamplesPerIteration, s.randomBits)

	s.currentSample++

	return true
}

func (s *ScrambledHammersley) GenerateSamples(iteration uint32) []math.Vector2 {
	for i := uint32(0); i < s.numSamplesPerIteration; i++ {
		s.samples[i] = math.ScrambledHammersley(i + iteration * s.numSamplesPerIteration, s.numTotalSamples, s.randomBits)
	}

	return s.samples
}