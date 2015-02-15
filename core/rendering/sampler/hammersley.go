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
	sample.LensUv = math.MakeVector2(s2d.Y, s2d.X)
	sample.Time = s2d.Y

	h.currentSample++

	return true
}

func (h *Hammersley) GenerateSamples(iteration uint32, buffer []math.Vector2) []math.Vector2 {
	offset := iteration * h.numSamplesPerIteration

	for i := uint32(0); i < h.numSamplesPerIteration; i++ {
		buffer[i] = math.Hammersley(i + offset, h.numTotalSamples)
	}

	return buffer
}

func (h *Hammersley) GenerateSample2D(index, iteration uint32) math.Vector2 {
	return math.Hammersley(index + iteration * h.numSamplesPerIteration, h.numTotalSamples)
}

func (h *Hammersley) GenerateSample1D(index, iteration uint32) float32 {
	return 0.5
}