package light

import (
	"github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/math32"
)

type Sphere struct {
	light
}

func NewSphere() *Sphere {
	l := Sphere{}
	return &l
}

func (l *Sphere) Samples(p math.Vector3, time float32, subsample, maxSamples uint32, sampler sampler.Sampler, samples []Sample) []Sample {
	samples = samples[:0]

	transformation := l.prop.TransformationAt(time)

	result := Sample{}

/*
	tsamples := sampler.GenerateSamples(subsample)

	for _, sample := range tsamples {
		ls := math.HemisphereSample_uniform(sample.X, sample.Y)
		ws := transformation.Rotation.TransformVector3(ls).Scale(l.radius)

		v := transformation.Position.Add(ws).Sub(p)

		d := v.SquaredLength()
		i := 1.0 / d

		result.L = v.Div(math32.Sqrt(d))
		result.Energy = l.color.Scale(i * l.lumen)

		samples = append(samples, result)
	}
*/

	for s := uint32(0); s < maxSamples; s++ {
		sample := sampler.GenerateSample(s, subsample)

		ls := math.HemisphereSample_uniform(sample.X, sample.Y)
		ws := transformation.Rotation.TransformVector3(ls).Scale(transformation.Scale.X)

		v := transformation.Position.Add(ws).Sub(p)

		d := v.SquaredLength()
		i := 1.0 / d

		result.L = v.Div(math32.Sqrt(d))
		result.Energy = l.color.Scale(i * l.lumen)

		samples = append(samples, result)		
	}

	return samples	
}