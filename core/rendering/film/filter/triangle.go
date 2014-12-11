package filter

import (
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/math32"
)

type Triangle struct {
	width math.Vector2
}

func NewTriangle(width math.Vector2) *Triangle {
	t := new(Triangle)
	t.width = width
	return t
}

func (t *Triangle) Evaluate(p math.Vector2) float32 {
	return math32.Max(0, t.width.X - math32.Abs(p.X)) * math32.Max(0, t.width.Y - math32.Abs(p.Y))
}