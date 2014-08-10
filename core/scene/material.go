package scene

import (
	"github.com/Opioid/scout/core/scene/shape"
	"github.com/Opioid/scout/base/math"
)

type Material interface {
	Evaluate(dg *shape.DifferentialGeometry, l, v math.Vector3) math.Vector3
	IsMirror() bool
}