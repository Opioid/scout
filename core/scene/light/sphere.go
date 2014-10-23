package light

import (
	pkgsampler "github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/base/math"
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

func (l *Sphere) Samples(p math.Vector3, sampler *pkgsampler.ScrambledHammersley, samples *[]Sample) {
	result := Sample{}

	sample := pkgsampler.Sample{}
//	for i, len := 0, cap(*samples); i < len; i++ {
	for sampler.GenerateNewSample(&sample) {
		ls := math.HemisphereSample_uniform(sample.Coordinates.X, sample.Coordinates.Y)
		ws := l.entity.Transformation.Rotation.TransformVector3(ls).Scale(l.radius)

		v := l.entity.Transformation.Position.Add(ws).Sub(p)

		d := v.SquaredLength()
		i := 1.0 / d

		result.L = v.Div(math.Sqrt(d))
		result.Energy = l.color.Scale(i * l.lumen)

		*samples = append(*samples, result)
	}
}