package material

import (
	"github.com/Opioid/scout/base/math"
)

type Sample interface {
	Evaluate(l math.Vector3) math.Vector3
	Values() *Values
}