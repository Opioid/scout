package light

import (
	"github.com/Opioid/scout/base/math"
)

type Directional struct {
	light
}

func NewDirectional() *Directional {
	return &Directional{}
}

func (l *Directional) Vector(p math.Vector3) math.Vector3 {
	return l.entity.Transformation.Rotation.Row(2).Scale(-1.0)
}

func (l *Directional) Light(p, color math.Vector3) math.Vector3 {
	return color.Mul(l.color)
}

