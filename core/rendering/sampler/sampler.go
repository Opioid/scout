package sampler

import (
	"github.com/Opioid/scout/base/math"
)

type Sampler interface {
	Restart()
	Resize(start, end math.Vector2i)
	SubSampler(start, end math.Vector2i) Sampler 
	GenerateNewSample(sample *Sample) bool
	Start() math.Vector2i
	NumSamplesPerPixel() uint32
}

type sampler struct {
	start math.Vector2i
	end   math.Vector2i
}

func (s *sampler) Start() math.Vector2i {
	return s.start
}