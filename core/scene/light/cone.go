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

	x, y := math.SampleDisk_uniform(sample.X, sample.Y)

	ls := math.MakeVector3(x, y, 0.0)
	ws := transformation.Rotation.TransformVector3(ls).Scale(transformation.Scale.X)

	v := transformation.Rotation.Direction().Scale(-1.0).Add(ws)

	result := Sample{Energy: l.color, L: v.Normalized(), T: 1000.0}

	return result
}
