package sampler

import (
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/random"
)

type ScrambledHammersley struct {
	rng *random.Generator

	numSamples uint32
	numTotalSamples uint32

	currentSample uint32

	randomBits uint32

	samples []math.Vector2
}

func MakeScrambledHammersley(rng *random.Generator) ScrambledHammersley {
	s := ScrambledHammersley{}
	s.rng = rng
	return s
}

func (s *ScrambledHammersley) Restart(numSamplesPerPixel uint32) {
	s.currentSample = 0

	s.numTotalSamples = s.numSamples * numSamplesPerPixel

	s.randomBits = s.rng.RandomUint32()
}

func (s *ScrambledHammersley) Resize(numSamples uint32) {
	s.numSamples = numSamples

	s.samples = make([]math.Vector2, numSamples)
}
/*
func (s *ScrambledHammersley) GenerateNewSample(subsample uint32, sample *Sample) bool {
	if s.currentSample >= s.numSamples {
		return false
	}

	sample.Coordinates = math.ScrambledHammersley(uint32(s.currentSample) + subsample * s.numSamples, s.numTotalSamples, s.randomBits)

	s.currentSample++

	return true
}
*/
func (s *ScrambledHammersley) GenerateSamples(subsample uint32) []math.Vector2 {
	for i := uint32(0); i < s.numSamples; i++ {
		s.samples[i] = math.ScrambledHammersley(i + subsample * s.numSamples, s.numTotalSamples, s.randomBits)
	}

	return s.samples
}