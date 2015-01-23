package light

import (
	"github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/base/math"
)

type Directional struct {
	light
}

func NewDirectional() *Directional {
	return &Directional{}
}

func (l *Directional) Samples(p math.Vector3, time float32, subsample, maxSamples uint32, sampler sampler.Sampler, samples []Sample) []Sample {
	samples = samples[:0]

	s := Sample{}

	transformation := l.prop.TransformationAt(time)

	s.L = transformation.Rotation.Direction().Scale(-1.0)
	s.Energy = l.color

	samples = append(samples, s)

	return samples
}

func (l *Directional) Sample(p math.Vector3, time float32, subsample uint32, sampler sampler.Sampler) Sample {
	transformation := l.prop.TransformationAt(time)

	result := Sample{Energy: l.color, L: transformation.Rotation.Direction().Scale(-1.0)}

	return result
}