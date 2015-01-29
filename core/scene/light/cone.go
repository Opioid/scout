package light

import (
	"github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/base/math"
)

type Cone struct {
	light
}

func NewCone() *Cone {
	d := Cone{}
	return &d
}

func (l *Cone) Sample(transformation *math.ComposedTransformation, p math.Vector3, time float32, subsample uint32, sampler sampler.Sampler) Sample {
	l.prop.TransformationAt(time, transformation)

	sample := sampler.GenerateSample2D(0, subsample)

	ls := math.DiskSample_uniform(sample.X, sample.Y)
	ws := transformation.Rotation.TransformVector3(ls).Scale(transformation.Scale.X)

	v := transformation.Rotation.Direction().Scale(-1.0).Add(ws)

	result := Sample{Energy: l.color, L: v.Normalized(), T: 1000.0}

	return result
}
