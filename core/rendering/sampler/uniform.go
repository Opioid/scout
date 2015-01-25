package sampler

import (
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/random"
)

type Uniform struct {
	sampler

	samples2d []math.Vector2
}

func NewUniform(samplesPerPixel math.Vector2i) *Uniform {
	u := new(Uniform)
	u.numSamplesPerIteration = uint32(samplesPerPixel.X * samplesPerPixel.Y)
	u.samples2d = make([]math.Vector2, u.numSamplesPerIteration)

	ax := 1.0 / float32(samplesPerPixel.X)
	ay := 1.0 / float32(samplesPerPixel.Y)

	for y, i := int32(0), int32(0); y < samplesPerPixel.Y; y++ {
		for x := int32(0); x < samplesPerPixel.X; x++ {
			u.samples2d[i] = math.MakeVector2((0.5 + float32(x)) * ax, (0.5 + float32(y)) * ay)
			i++
		}
	}	

	return u
}

func (u *Uniform) Clone(rng *random.Generator) Sampler {
	nu := new(Uniform)
	nu.numSamplesPerIteration = u.numSamplesPerIteration
	nu.samples2d = u.samples2d
	return nu
}

func (u *Uniform) Restart(numIterations uint32) {
	u.currentSample = 0
}

func (u *Uniform) GenerateCameraSample(offset math.Vector2, sample *CameraSample) bool {
	if u.currentSample >= u.numSamplesPerIteration {
		return false
	}

	s2d := u.samples2d[u.currentSample]

	sample.Coordinates = offset.Add(s2d)
	sample.RelativeOffset = s2d.SubS(0.5)

	u.currentSample++

	return true
}

func (u *Uniform) GenerateSamples(iteration uint32, buffer []math.Vector2) []math.Vector2 {
	copy(buffer, u.samples2d)
	return buffer
}

func (u *Uniform) GenerateSample(index, iteration uint32) math.Vector2 {
	return u.samples2d[index]
}

func (u *Uniform) GenerateSample1D(index, iteration uint32) float32 {
	return 0.5
}