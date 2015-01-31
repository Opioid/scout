package light

import (
	"github.com/Opioid/scout/core/scene/prop"
	"github.com/Opioid/scout/core/scene/shape"
	"github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/math32"
)

type Sphere struct {
	light
}

func NewSphere(shape shape.Shape) *Sphere {
	l := Sphere{}
	l.prop.SetVisible(prop.IsLight, true)
	l.prop.Shape = shape
	return &l
}

func (l *Sphere) Sample(transformation *math.ComposedTransformation, p math.Vector3, time float32, subsample uint32, sampler sampler.Sampler) Sample {
	l.prop.TransformationAt(time, transformation)

	sample := sampler.GenerateSample2D(0, subsample)

	ls := math.SampleHemisphere_uniform(sample.X, sample.Y)
	ws := transformation.Rotation.TransformVector3(ls).Scale(transformation.Scale.X)

	v := transformation.Position.Add(ws).Sub(p)

	d := v.SquaredLength()
	i := 1.0 / d
	t := math32.Sqrt(d)

	result := Sample{Energy: l.color.Scale(i * l.lumen), L: v.Div(t), T: t}

	return result
}