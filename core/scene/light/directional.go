package light

import (
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/random"
)

type Directional struct {
	light
}

func NewDirectional() *Directional {
	return &Directional{}
}

func (l *Directional) Vector(p math.Vector3) math.Vector3 {
	return l.entity.Transformation.Rotation.Direction().Scale(-1.0)
}

func (l *Directional) Light(p, color math.Vector3) math.Vector3 {
	return color.Mul(l.color)
}

func (l *Directional) Samples(p math.Vector3, rng *random.Generator, samples *[]Sample) {
	s := Sample{}

	s.L = l.entity.Transformation.Rotation.Direction().Scale(-1.0)
	s.Energy = l.color

	*samples = append(*samples, s)
}
