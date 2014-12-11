package filter

import (
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/math32"
)

type Mitchell struct {
	width math.Vector2
	invWidth math.Vector2
	b, c float32
}

func NewMitchell(width math.Vector2, b, c float32) *Mitchell {
	m := new(Mitchell)
	m.width = width
	m.invWidth = math.MakeVector2(1 / width.X, 1 / width.Y)
	m.b = b
	m.c = c
	return m
}

func (m *Mitchell) Evaluate(p math.Vector2) float32 {
	return m.mitchell1D(p.X * m.invWidth.X) * m.mitchell1D(p.Y * m.invWidth.Y)
}

func (m *Mitchell) mitchell1D(x float32) float32 {
	x = math32.Abs(2 * x)

	x2 := x * x
	
	if x > 1 {
		return ((-m.b -6 * m.c) * x2 * x + (6 * m.b + 30 * m.c) * x2 + 
				(-12 * m.b - 48 * m.c) * x + (8 * m.b + 24 * m.c)) * (1 / 6)
	} else {
	/*	return ((12 - 9 * m.b - 6 * m.c) * x2 * x + 
				(-18 + 12 * m.b + 6 * m.c) * x2 + 
				(6 - 2 * m.b)) * (1 / 6)*/
		return 1
	}
} 