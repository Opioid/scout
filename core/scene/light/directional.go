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

func (l *Directional) Samples(p math.Vector3, sampler *sampler.Stratified, samples *[]Sample) {
	s := Sample{}

	s.L = l.entity.Transformation.Rotation.Direction().Scale(-1.0)
	s.Energy = l.color

	*samples = append(*samples, s)
}
