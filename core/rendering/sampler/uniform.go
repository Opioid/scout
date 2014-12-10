package sampler

import (
	"github.com/Opioid/scout/base/math"
	_ "fmt"
)

type Uniform struct {
	sampler
	samplesPerPixel math.Vector2i
	currentPixel math.Vector2i
	currentSample int
	offsets []math.Vector2
}

func NewUniform(start, end, samplesPerPixel math.Vector2i) *Uniform {
	u := new(Uniform)
	u.samplesPerPixel = samplesPerPixel
	u.offsets = make([]math.Vector2, samplesPerPixel.X * samplesPerPixel.Y)
	u.Resize(start, end)
	u.currentPixel = start
	return u
}

func (u *Uniform) Resize(start, end math.Vector2i) {
	u.start = start
	u.end = end

	ox := 1 / float32(u.samplesPerPixel.X)
	oy := 1 / float32(u.samplesPerPixel.Y)

	for y, i := int32(0), int32(0); y < u.samplesPerPixel.Y; y++ {
		for x := int32(0); x < u.samplesPerPixel.X; x++ {
			u.offsets[i] = math.MakeVector2((0.5 + float32(x)) * ox, (0.5 + float32(y)) * oy)
			i++
		}
	}
}

func (u *Uniform) Restart() {
	u.currentPixel = u.start
	u.currentSample = 0
}

func (u *Uniform) SubSampler(start, end math.Vector2i) Sampler {
	return NewUniform(start, end, u.samplesPerPixel)
}

func (u *Uniform) GenerateNewSample(s *Sample) bool {
	if u.currentPixel.X >= u.end.X {
		u.currentPixel.X = u.start.X
		u.currentPixel.Y++
	}

	if u.currentPixel.Y >= u.end.Y {
		return false
	}

	o := u.offsets[u.currentSample]

	s.Coordinates = math.MakeVector2(float32(u.currentPixel.X) + o.X, float32(u.currentPixel.Y) + o.Y)
	s.RelativeOffset = math.MakeVector2(o.X - 0.5, o.Y - 0.5)
	s.Id = uint32(u.currentSample)

	u.currentSample++

	if u.currentSample >= len(u.offsets) {
		u.currentSample = 0
		u.currentPixel.X++
	}

	return true
}

func (u *Uniform) NumSamplesPerPixel() uint32 {
	return uint32(u.samplesPerPixel.X * u.samplesPerPixel.Y)
}