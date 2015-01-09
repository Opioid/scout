package filter

import (
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/math32"
	_ "fmt"
)

type MitchellNetravali struct {
	width math.Vector2
	invWidth math.Vector2

	a0, a1, a2, b0, b1, b2, b3 float32
}

func NewMitchellNetravali(width math.Vector2, b, c float32) *MitchellNetravali {
	m := MitchellNetravali{}
	m.width = width
	m.invWidth = math.MakeVector2(1.0 / width.X, 1.0 / width.Y)

	m.a0 = 12.0 - 9.0 * b - 6.0 * c
	m.a1 = -18.0 + 12.0 * b + 6.0 * c
	m.a2 = 6.0 - 2.0 * b

	m.b0 = -b -6.0 * c
	m.b1 = 6.0 * b + 30.0 * c
	m.b2 = -12.0 * b - 48.0 * c
	m.b3 = 8.0 * b + 24.0 * c

	return &m
}

func (m *MitchellNetravali) Evaluate(p math.Vector2) float32 {
	return m.mitchellNetravali1D(p.X * m.invWidth.X) * m.mitchellNetravali1D(p.Y * m.invWidth.Y)
}

func (m *MitchellNetravali) mitchellNetravali1D(x float32) float32 {
	x = math32.Abs(x)

	x2 := x * x
	
	if x < 1.0 {
		return (m.a0 * x2 * x + m.a1 * x2 + m.a2) * (1.0 / 6.0)
	} else if x < 2.0 {
		return (m.b0 * x2 * x + m.b1 * x2 + m.b2 * x + m.b3) * (1.0 / 6.0)
	} else {
		return 0.0
	}
} 