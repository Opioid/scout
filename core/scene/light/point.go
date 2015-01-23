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

func (l *Point) Samples(p math.Vector3, time float32, subsample, maxSamples uint32, sampler sampler.Sampler, samples []Sample) []Sample {
	samples = samples[:0]

	result := Sample{}

	transformation := l.prop.TransformationAt(time)

	v := transformation.Position.Sub(p)

	d := v.SquaredLength()
	i := 1.0 / d

	result.L = v.Div(math32.Sqrt(d))
	result.Energy = l.color.Scale(i * l.lumen)

	samples = append(samples, result)

	return samples
}

func (l *Point) Sample(p math.Vector3, time float32, subsample uint32, sampler sampler.Sampler) Sample {
	transformation := l.prop.TransformationAt(time)

	v := transformation.Position.Sub(p)

	d := v.SquaredLength()
	i := 1.0 / d

	result := Sample{Energy: l.color.Scale(i * l.lumen), L: v.Div(math32.Sqrt(d))}

	return result
}