package light

import (
	pkgsampler "github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/math32"
)

type Sphere struct {
	light
	radius float32
}

func NewSphere(radius float32) *Sphere {
	l := Sphere{}
	l.radius = radius
	return &l
}

func (l *Sphere) Samples(p math.Vector3, subsample uint32, time float32, sampler *pkgsampler.ScrambledHammersley, samples *[]Sample) {
	transformation := l.entity.TransformationAt(time)

	result := Sample{}

	tsamples := sampler.GenerateSamples(subsample)

	for _, sample := range tsamples {
		ls := math.HemisphereSample_uniform(sample.X, sample.Y)
		ws := transformation.Rotation.TransformVector3(ls).Scale(l.radius)

		v := transformation.Position.Add(ws).Sub(p)

		d := v.SquaredLength()
		i := 1.0 / d

		result.L = v.Div(math32.Sqrt(d))
		result.Energy = l.color.Scale(i * l.lumen)

		*samples = append(*samples, result)
	}
}