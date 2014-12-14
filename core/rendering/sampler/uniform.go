package sampler

import (
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/random"
	_ "fmt"
)

type Uniform struct {
	samplesPerPixel math.Vector2i
	numSamples uint32
	currentSample uint32
	offsets []math.Vector2
}

func NewUniform(samplesPerPixel math.Vector2i) *Uniform {
	u := new(Uniform)
	u.samplesPerPixel = samplesPerPixel
	u.numSamples = uint32(u.samplesPerPixel.X * u.samplesPerPixel.Y)
	u.offsets = make([]math.Vector2, samplesPerPixel.X * samplesPerPixel.Y)

	ox := 1.0 / float32(u.samplesPerPixel.X)
	oy := 1.0 / float32(u.samplesPerPixel.Y)

	for y, i := int32(0), int32(0); y < u.samplesPerPixel.Y; y++ {
		for x := int32(0); x < u.samplesPerPixel.X; x++ {
			u.offsets[i] = math.MakeVector2((0.5 + float32(x)) * ox, (0.5 + float32(y)) * oy)
			i++
		}
	}	
	return u
}

func (u *Uniform) Clone(rng *random.Generator) Sampler {
	return NewUniform(u.samplesPerPixel)
}

func (u *Uniform) Restart() {
	u.currentSample = 0
}

func (u *Uniform) GenerateNewSample(sample *math.Vector2) bool {
	if u.currentSample >= u.numSamples{
		return false
	}

	*sample = u.offsets[u.currentSample]

	u.currentSample++

	return true
}

func (u *Uniform) NumSamplesPerPixel() uint32 {
	return u.numSamples
}