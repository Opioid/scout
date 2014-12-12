package sampler

import (
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/random"
)

type Sampler interface {
	SubSampler(start, end math.Vector2i, rng *random.Generator) Sampler 
	NumSamplesPerPixel() uint32

	Start() math.Vector2i

	Restart()
	Resize(start, end math.Vector2i)
	GenerateNewSample(sample *Sample) bool
}

type sampler struct {
	start math.Vector2i
	end   math.Vector2i
}

func (s *sampler) Start() math.Vector2i {
	return s.start
}