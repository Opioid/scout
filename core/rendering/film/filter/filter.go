package filter

import (
	"github.com/Opioid/scout/base/math"
)

type Filter interface {
	Evaluate(p math.Vector2) float32
}