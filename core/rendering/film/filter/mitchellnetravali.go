package filter

import (
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/math32"
)

type MitchellNetravali struct {
	width math.Vector2
	invWidth math.Vector2
	b, c float32
}

func NewMitchellNetravali(width math.Vector2, b, c float32) *MitchellNetravali {
	m := new(MitchellNetravali)
	m.width = width
	m.invWidth = math.MakeVector2(1.0 / width.X, 1.0 / width.Y)
	m.b = b
	m.c = c
	return m
}

func (m *MitchellNetravali) Evaluate(p math.Vector2) float32 {
	return m.mitchellNetravali1D(p.X * m.invWidth.X) * m.mitchellNetravali1D(p.Y * m.invWidth.Y)
}

func (m *MitchellNetravali) mitchellNetravali1D(x float32) float32 {
	x = math32.Abs(x)

	x2 := x * x
	
	if x < 1.0 {
		return ((12.0 - 9.0 * m.b - 6.0 * m.c) * x2 * x + 
				(-18.0 + 12.0 * m.b + 6.0 * m.c) * x2 + (6.0 - 2.0 * m.b)) * (1.0 / 6.0)
	} else if x < 2.0 {
		return ((-m.b -6.0 * m.c) * x2 * x + (6.0 * m.b + 30.0 * m.c) * x2 + 
				(-12.0 * m.b - 48.0 * m.c) * x + (8.0 * m.b + 24.0 * m.c)) * (1.0 / 6.0)
	} else {
		return 0.0
	}
} 