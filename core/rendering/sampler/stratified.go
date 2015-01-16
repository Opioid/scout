package sampler

import (
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/random"
	_ "fmt"
)

type Stratified struct {
	numStratifiedSamples math.Vector2i

	numSamples uint32

	currentSample uint32

	area math.Vector2

	// offsets never change after initialization
	offsets []math.Vector2

	// filled with the actual stratified samples every time GenerateSamples is called
	samples []math.Vector2

	rng *random.Generator
}

func NewStratified(numStratifiedSamples math.Vector2i, rng *random.Generator) *Stratified {
	s := new(Stratified)
	s.allocateSamples(numStratifiedSamples)
	s.rng = rng
	return s
}

func (s *Stratified) allocateSamples(numStratifiedSamples math.Vector2i) {
	s.numStratifiedSamples = numStratifiedSamples
	s.numSamples = uint32(numStratifiedSamples.X * numStratifiedSamples.Y)

	s.offsets = make([]math.Vector2, s.numSamples)
	s.samples = make([]math.Vector2, s.numSamples)

	s.area.X = 1.0 / float32(numStratifiedSamples.X)
	s.area.Y = 1.0 / float32(numStratifiedSamples.Y)

	for y, i := int32(0), int32(0); y < numStratifiedSamples.Y; y++ {
		for x := int32(0); x < numStratifiedSamples.X; x++ {
			s.offsets[i] = math.MakeVector2((0.5 + float32(x)) * s.area.X, (0.5 + float32(y)) * s.area.Y)
			i++
		}
	}
}

func (s *Stratified) Clone(rng *random.Generator) Sampler {
	ns := new(Stratified)
	ns.allocateSamples(s.numStratifiedSamples)
	ns.rng = rng
	return ns
}

func (s *Stratified) NumSamplesPerIteration() uint32 {
	return s.numSamples
}

func (s *Stratified) Restart(numIterations uint32) {
	s.currentSample = 0
}

func (s *Stratified) GenerateCameraSample(offset math.Vector2, sample *CameraSample) bool {
	if s.currentSample >= s.numSamples {
		return false
	}

	s2d := s.sample(s.currentSample)

	sample.Coordinates = offset.Add(s2d)
	sample.RelativeOffset = s2d.SubS(0.5)
	sample.LensUv = s2d//math.MakeVector2(s.rng.RandomFloat32(), s.rng.RandomFloat32())

	s.currentSample++
	
	return true
}

func (s *Stratified) GenerateSamples(iteration uint32) []math.Vector2 {
	for i := uint32(0); i < s.numSamples; i++ {
		s.samples[i] = s.sample(i)
	}

	return s.samples
}

func (s *Stratified) GenerateSample(index, iteration uint32) math.Vector2 {
	return s.sample(index)
}

func (s *Stratified) sample(id uint32) math.Vector2 {
	sample := s.offsets[id]
	sample.X += s.area.X * (s.rng.RandomFloat32() - 0.5)
	sample.Y += s.area.Y * (s.rng.RandomFloat32() - 0.5)
	return sample
}