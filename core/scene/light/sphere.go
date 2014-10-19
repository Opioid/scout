package light

import (
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/random"
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

func (l *Sphere) Vector(p math.Vector3) math.Vector3 {
	return l.entity.Transformation.Position.Sub(p).Normalized()
}

func (l *Sphere) Light(p, color math.Vector3) math.Vector3 {
	d := l.entity.Transformation.Position.Sub(p).SquaredLength()
	i := 1.0 / d

	return color.Mul(l.color).Scale(i * l.lumen)
}


func (l *Sphere) Samples(p math.Vector3, rng *random.Generator, samples *[]Sample) {
	result := Sample{}

	for i, len := 0, cap(*samples); i < len; i++ {
		ls := math.HemisphereSample_uniform(rng.RandomFloat32(), rng.RandomFloat32())
		ws := l.entity.Transformation.Rotation.TransformVector3(ls).Scale(l.radius)

		v := l.entity.Transformation.Position.Add(ws).Sub(p)

		d := v.SquaredLength()
		i := 1.0 / d

		result.L = v.Div(math.Sqrt(d))
		result.Energy = l.color.Scale(i * l.lumen)

		*samples = append(*samples, result)

	}
}