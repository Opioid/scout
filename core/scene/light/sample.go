package light

import (
	"github.com/Opioid/scout/base/math"
	_ "math"
)

type Sample struct {
	Energy math.Vector3
	L math.Vector3
	T float32
	Pdf float32
}