package sampler

import (
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/random"
	_ "fmt"
)

type EMS struct {
	sampler

	rng *random.Generator
	randomBits uint32

	numTotalSamples uint32	
}

func NewEMS(numSamplesPerIteration uint32, rng *random.Generator) *EMS {
	s := new(EMS)
	s.rng = rng
	s.allocateSamples(numSamplesPerIteration)
	return s
}

func (s *EMS) allocateSamples(numSamplesPerIteration uint32) {
	s.numSamplesPerIteration = numSamplesPerIteration
}

func (s *EMS) Clone(rng *random.Generator) Sampler {
	return NewScrambledHammersley(s.numSamplesPerIteration, rng)
}

func (s *EMS) Restart(numIterations uint32) {
	s.currentSample = 0
	s.numTotalSamples = s.numSamplesPerIteration * numIterations
	s.randomBits = s.rng.RandomUint32()
}

func (s *EMS) GenerateCameraSample(offset math.Vector2, sample *CameraSample) bool {
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

func (s *EMS) GenerateSamples(iteration uint32, buffer []math.Vector2) []math.Vector2 {
	offset := iteration * s.numSamplesPerIteration

	for i := uint32(0); i < s.numSamplesPerIteration; i++ {
		index := i + offset
		buffer[i] = math.MakeVector2(math.ScrambledRadicalInverse_vdC(index, s.randomBits), math.RadicalInverse_S(index, s.randomBits))
	}

	return buffer
}

func (s *EMS) GenerateSample(index, iteration uint32) math.Vector2 {
	i := index + iteration * s.numSamplesPerIteration
	return math.MakeVector2(math.ScrambledRadicalInverse_vdC(i, s.randomBits), math.RadicalInverse_S(i, s.randomBits))
}