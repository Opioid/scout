package light

import (
	"github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/base/math"
)

type Point struct {
	light
}

func NewPoint() *Point {
	return &Point{}
}

func (l *Point) Samples(p math.Vector3, sampler *sampler.Stratified, samples *[]Sample) {
	result := Sample{}

	v := l.entity.Transformation.Position.Sub(p)

	d := v.SquaredLength()
	i := 1.0 / d

	result.L = v.Div(math.Sqrt(d))
	result.Energy = l.color.Scale(i * l.lumen)

	*samples = append(*samples, result)
}