package sampler

import (
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/random"
	_ "fmt"
)

type ScrambledHammersley struct {
	rng *random.Generator
	randomBits uint32

	currentSample uint32
	numSamplesPerIteration uint32
	numTotalSamples uint32	

	samples []math.Vector2
}

func NewScrambledHammersley(numSamplesPerIteration uint32, rng *random.Generator) *ScrambledHammersley {
	s := new(ScrambledHammersley)
	s.rng = rng
	s.allocateSamples(numSamplesPerIteration)
	return s
}

func (s *ScrambledHammersley) allocateSamples(numSamplesPerIteration uint32) {
	s.numSamplesPerIteration = numSamplesPerIteration
	s.samples = make([]math.Vector2, numSamplesPerIteration)
}

func (s *ScrambledHammersley) Clone(rng *random.Generator) Sampler {
	return NewScrambledHammersley(s.numSamplesPerIteration, rng)
}

func (s *ScrambledHammersley) NumSamplesPerIteration() uint32 {
	return s.numSamplesPerIteration
}

func (s *ScrambledHammersley) Restart(numIterations uint32) {
	s.currentSample = 0
	s.numTotalSamples = s.numSamplesPerIteration * numIterations
	s.randomBits = s.rng.RandomUint32()
}

func (s *ScrambledHammersley) GenerateCameraSample(offset math.Vector2, sample *CameraSample) bool {
	if s.currentSample >= s.numSamplesPerIteration {
		return false
	}

	s2d := math.ScrambledHammersley(s.currentSample, s.numSamplesPerIteration, s.randomBits)

	sample.Coordinates = offset.Add(s2d)
	sample.RelativeOffset = s2d.SubS(0.5)
	sample.LensUv = s2d
	sample.Time = s2d.Y

	s.currentSample++

	return true
}

func (s *ScrambledHammersley) GenerateSamples(iteration uint32) []math.Vector2 {
	offset := iteration * s.numSamplesPerIteration

	for i := uint32(0); i < s.numSamplesPerIteration; i++ {
		s.samples[i] = math.ScrambledHammersley(i + offset, s.numTotalSamples, s.randomBits)
	}

	return s.samples
}

func (s *ScrambledHammersley) GenerateSample(index, iteration uint32) math.Vector2 {
	return math.ScrambledHammersley(index + iteration * s.numSamplesPerIteration, s.numTotalSamples, s.randomBits)
}