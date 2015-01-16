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

func (l *Point) Samples(p math.Vector3, time float32, subsample, maxSamples uint32, sampler sampler.Sampler, samples *[]Sample) {
	result := Sample{}

	transformation := l.entity.TransformationAt(time)

	v := transformation.Position.Sub(p)

	d := v.SquaredLength()
	i := 1 / d

	result.L = v.Div(math32.Sqrt(d))
	result.Energy = l.color.Scale(i * l.lumen)

	*samples = append(*samples, result)
}