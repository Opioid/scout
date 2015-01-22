package sampler

import (
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/random"
	_ "fmt"
)

type Hammersley struct {
	sampler

	numTotalSamples uint32	
}

func NewHammersley(numSamplesPerIteration uint32) *Hammersley {
	h := new(Hammersley)
	h.allocateSamples(numSamplesPerIteration)
	return h
}

func (h *Hammersley) allocateSamples(numSamplesPerIteration uint32) {
	h.numSamplesPerIteration = numSamplesPerIteration
	h.samples2d = make([]math.Vector2, numSamplesPerIteration)
}

func (h *Hammersley) Clone(rng *random.Generator) Sampler {
	return NewHammersley(h.numSamplesPerIteration)
}

func (h *Hammersley) Restart(numIterations uint32) {
	h.currentSample = 0
	h.numTotalSamples = h.numSamplesPerIteration * numIterations
}

func (h *Hammersley) GenerateCameraSample(offset math.Vector2, sample *CameraSample) bool {
	if h.currentSample >= h.numSamplesPerIteration {
		return false
	}

	s2d := math.Hammersley(h.currentSample, h.numSamplesPerIteration)
	
	sample.Coordinates = offset.Add(s2d)
	sample.RelativeOffset = s2d.SubS(0.5)
	sample.LensUv = s2d
	sample.Time = s2d.Y

	h.currentSample++

	return true
}

func (h *Hammersley) GenerateSamples(iteration uint32) []math.Vector2 {
	offset := iteration * h.numSamplesPerIteration

	for i := uint32(0); i < h.numSamplesPerIteration; i++ {
		h.samples2d[i] = math.Hammersley(i + offset, h.numTotalSamples)
	}

	return h.samples2d
}

func (h *Hammersley) GenerateSample(index, iteration uint32) math.Vector2 {
	return math.Hammersley(index + iteration * h.numSamplesPerIteration, h.numTotalSamples)
}