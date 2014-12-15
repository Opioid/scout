package sampler

import (
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/random"
)

type Uniform struct {
	numSamples uint32
	currentSample uint32
	samples []math.Vector2
}

func NewUniform(samplesPerPixel math.Vector2i) *Uniform {
	u := new(Uniform)
	u.numSamples = uint32(samplesPerPixel.X * samplesPerPixel.Y)
	u.samples = make([]math.Vector2, u.numSamples)

	ax := 1.0 / float32(samplesPerPixel.X)
	ay := 1.0 / float32(samplesPerPixel.Y)

	for y, i := int32(0), int32(0); y < samplesPerPixel.Y; y++ {
		for x := int32(0); x < samplesPerPixel.X; x++ {
			u.samples[i] = math.MakeVector2((0.5 + float32(x)) * ax, (0.5 + float32(y)) * ay)
			i++
		}
	}	

	return u
}

func (u *Uniform) Clone(rng *random.Generator) Sampler {
	nu := new(Uniform)
	nu.numSamples = u.numSamples
	nu.samples = u.samples
	return nu
}

func (u *Uniform) NumSamplesPerIteration() uint32 {
	return u.numSamples
}

func (u *Uniform) Restart(numIterations uint32) {
	u.currentSample = 0
}

func (u *Uniform) GenerateNewSample(sample *math.Vector2) bool {
	if u.currentSample >= u.numSamples {
		return false
	}

	*sample = u.samples[u.currentSample]

	u.currentSample++

	return true
}

func (u *Uniform) GenerateSamples(iteration uint32) []math.Vector2 {
	return u.samples
}