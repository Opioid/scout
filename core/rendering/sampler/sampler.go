package sampler

import (
	"github.com/Opioid/scout/base/math"
)

type Sampler interface {
	Restart()
	Resize(start, end math.Vector2i)
	SubSampler(start, end math.Vector2i) Sampler 
	GenerateNewSample(sample *Sample) bool
}

type base struct {
	start math.Vector2i
	end   math.Vector2i
}