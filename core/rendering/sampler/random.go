package sampler

import (
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/random"
	_ "fmt"
)

type Random struct {
	sampler

	rng *random.Generator
}

func NewRandom(numSamplesPerIteration uint32, rng *random.Generator) *Random {
	s := new(Random)
	s.rng = rng
	s.allocateSamples(numSamplesPerIteration)
	return s
}

func (s *Random) allocateSamples(numSamplesPerIteration uint32) {
	s.numSamplesPerIteration = numSamplesPerIteration
}

func (s *Random) Clone(rng *random.Generator) Sampler {
	return NewRandom(s.numSamplesPerIteration, rng)
}

func (s *Random) Restart(numIterations uint32) {
	s.currentSample = 0
}

func (s *Random) GenerateCameraSample(offset math.Vector2, sample *CameraSample) bool {
	if s.currentSample >= s.numSamplesPerIteration {
		return false
	}

	s2d := math.MakeVector2(s.rng.RandomFloat32(), s.rng.RandomFloat32())

	sample.Coordinates = offset.Add(s2d)
	sample.RelativeOffset = s2d.SubS(0.5)
	sample.LensUv = s2d
	sample.Time = s2d.Y

	s.currentSample++

	return true
}

func (s *Random) GenerateSamples(iteration uint32, buffer []math.Vector2) []math.Vector2 {
	for i := uint32(0); i < s.numSamplesPerIteration; i++ {
		buffer[i] = math.MakeVector2(s.rng.RandomFloat32(), s.rng.RandomFloat32())
	}

	return buffer
}

func (s *Random) GenerateSample(index, iteration uint32) math.Vector2 {
	return math.MakeVector2(s.rng.RandomFloat32(), s.rng.RandomFloat32())
}

func (s *Random) GenerateSample1D(index, iteration uint32) float32 {
	return s.rng.RandomFloat32()
}