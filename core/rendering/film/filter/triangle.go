package filter

import (
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/math32"
)

type Triangle struct {
	width math.Vector2
}

func (t *Triangle) SetWidth(width math.Vector2) {
	t.width = width
}

func (t *Triangle) Evaluate(c math.Vector2) float32 {
	return math32.Max(0, t.width.X - math32.Abs(c.X)) * math32.Max(0, t.width.Y - math32.Abs(c.Y))
}