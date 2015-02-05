package light

import (
	"github.com/Opioid/scout/core/scene/prop"
	"github.com/Opioid/scout/core/scene/shape"
	"github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/math32"
	gomath "math"
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

const (
	hemispherePdf = 1.0 / (2.0 * gomath.Pi)
	hemisphereArea = 2.0 / gomath.Pi
)

func (l *Sphere) Sample(transformation *math.ComposedTransformation, p math.Vector3, time float32, subsample uint32, sampler sampler.Sampler) Sample {
	l.prop.TransformationAt(time, transformation)

	sample := sampler.GenerateSample2D(0, subsample)
	ls := math.SampleHemisphere_uniform(sample.X, sample.Y)

	z := p.Sub(transformation.Position).Normalized()
	cs := math.MakeCoordinateSystemMatrix3x3(z)

	n := cs.TransformVector3(ls)

	ws := transformation.Position.Add(n.Scale(transformation.Scale.X))

	v := ws.Sub(p)

	d := v.SquaredLength()
	t := math32.Sqrt(d)

	nDotV := n.Dot(v.Scale(-1.0).Normalized())

	if n.Dot(v.Normalized()) > 0.0 {
		// In this case no light will reach p, so we could make an early out
		d = 0.0
	}

	result := Sample{Energy: l.color.Scale(l.lumen), L: v.Div(t), T: t, Pdf: d / (math32.Abs(nDotV) * hemisphereArea)}

	return result
}