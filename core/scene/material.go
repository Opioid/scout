package scene

import (
	"github.com/Opioid/scout/base/math"
)

type Material interface {
	Evaluate(n, l, v math.Vector3) math.Vector3
	IsMirror() bool
}