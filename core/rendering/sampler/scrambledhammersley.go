package sampler

import (
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/random"
)

type ScrambledHammersley struct {
	rng *random.Generator

	numSamples uint32

	currentSample uint32

	randomBits uint32
}

func MakeScrambledHammersley(rng *random.Generator) ScrambledHammersley {
	s := ScrambledHammersley{}
	s.rng = rng
	return s
}

func (s *ScrambledHammersley) Restart() {
	s.currentSample = 0

	s.randomBits = s.rng.RandomUint32()
}

func (s *ScrambledHammersley) Resize(numSamples uint32) {
	s.numSamples = numSamples
}

func (s *ScrambledHammersley) GenerateNewSample(sample *Sample) bool {
	if s.currentSample >= s.numSamples {
		return false
	}

	sample.Coordinates = math.ScrambledHammersley(uint32(s.currentSample), s.numSamples, s.randomBits)
/*
	sample.Coordinates = math.MakeVector2(
		math.ScrambledRadicalInverse_vdC(s.currentSample, s.randomBits), 
		math.RadicalInverse_S(s.currentSample, s.randomBits))
*/
	s.currentSample++

	return true
}