package light

import (
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/random"
)

type Point struct {
	light
}

func NewPoint() *Point {
	return &Point{}
}

func (l *Point) Vector(p math.Vector3) math.Vector3 {
	return l.entity.Transformation.Position.Sub(p).Normalized()
}

func (l *Point) Light(p, color math.Vector3) math.Vector3 {
	d := l.entity.Transformation.Position.Sub(p).SquaredLength()
	i := 1.0 / d

	return color.Mul(l.color).Scale(i * l.lumen)
}

func (l *Point) Samples(p math.Vector3, rng *random.Generator, samples *[]Sample) {
	result := Sample{}

	v := l.entity.Transformation.Position.Sub(p)

	d := v.SquaredLength()
	i := 1.0 / d

	result.L = v.Div(math.Sqrt(d))
	result.Energy = l.color.Scale(i * l.lumen)

	*samples = append(*samples, result)
}