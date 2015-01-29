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

func (l *Directional) Sample(transformation *math.ComposedTransformation, p math.Vector3, time float32, subsample uint32, sampler sampler.Sampler) Sample {
	l.prop.TransformationAt(time, transformation)

	result := Sample{Energy: l.color, L: transformation.Rotation.Direction().Scale(-1.0), T: 1000.0}

	return result
}