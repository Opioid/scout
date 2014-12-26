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

func (l *Directional) Samples(p math.Vector3, subsample uint32, time float32, sampler *sampler.ScrambledHammersley, samples *[]Sample) {
	s := Sample{}

	transformation := l.entity.TransformationAt(time)

	s.L = transformation.Rotation.Direction().Scale(-1.0)
	s.Energy = l.color

	*samples = append(*samples, s)
}
