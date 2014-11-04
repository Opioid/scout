package light 

import (
	"github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/base/math"
)

type Ambient struct {
	light
}

func NewAmbient() *Ambient {
	return &Ambient{}
}

func (l *Ambient) Samples(p math.Vector3, subsample uint32, sampler *sampler.ScrambledHammersley, samples *[]Sample) {
	result := Sample{}

	result.Energy = l.color

	*samples = append(*samples, result)
}