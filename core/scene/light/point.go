package light

import (
	"github.com/Opioid/scout/base/math"
	_ "math"
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

//	energy := l.color.Scale(l.lumen / (4.0 * gomath.Pi))

	return color.Mul(l.color).Scale(i * l.lumen)
}