package sampler

import (
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/random"
)

type Stratified struct {
	rng *random.Generator

	numSamples math.Vector2i
	step math.Vector2

	currentSample math.Vector2i
}

func MakeStratified(rng *random.Generator) Stratified {
	s := Stratified{}
	s.rng = rng
	return s
}

func (s *Stratified) Restart() {
	s.currentSample.X = 0
	s.currentSample.Y = 0
}

func (s *Stratified) Resize(numSamples math.Vector2i) {
	s.numSamples = numSamples

	s.step = math.MakeVector2(1.0 / float32(numSamples.X), 1.0 / float32(numSamples.Y))
}

func (s *Stratified) GenerateNewSample(sample *Sample) bool {
	if s.currentSample.X >= s.numSamples.X {
		s.currentSample.X = 0
		s.currentSample.Y++
	}

	if s.currentSample.Y >= s.numSamples.Y {
		return false
	}

	sample.Coordinates.X = (float32(s.currentSample.X) + s.rng.RandomFloat32()) * s.step.X
	sample.Coordinates.Y = (float32(s.currentSample.Y) + s.rng.RandomFloat32()) * s.step.Y

	s.currentSample.X++

	return true
}