package light

import (
	"github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/math32"
)

type Point struct {
	light
}

func NewPoint() *Point {
	return &Point{}
}

func (l *Point) Sample(transformation *math.ComposedTransformation, p math.Vector3, time float32, subsample uint32, sampler sampler.Sampler) Sample {
	l.prop.TransformationAt(time, transformation)

	v := transformation.Position.Sub(p)

	d := v.SquaredLength()
	i := 1.0 / d
	t := math32.Sqrt(d)

	result := Sample{Energy: l.color.Scale(i * l.lumen), L: v.Div(t), T: t, Pdf: 1.0}

	return result
}