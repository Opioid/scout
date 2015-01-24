package filter

import (
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/math32"
)

type Gaussian struct {
	width math.Vector2
	alpha float32
	exp   math.Vector2
}

func NewGaussian(width math.Vector2, alpha float32) *Gaussian {
	g := new(Gaussian)
	g.width = width
	g.alpha = alpha
	g.exp = math.MakeVector2(math.Exp(-alpha * width.X * width.X), math.Exp(-alpha * width.Y * width.Y))
	return g
}

func (g *Gaussian) Evaluate(p math.Vector2) float32 {
	return g.gaussian(p.X, g.exp.X) * g.gaussian(p.Y, g.exp.Y)
}

func (g *Gaussian) gaussian(d, exp float32) float32 {
	return math32.Max(0.0, math.Exp(-g.alpha * d * d) - exp)
}