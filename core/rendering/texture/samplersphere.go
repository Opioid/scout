package texture

import (
	"github.com/Opioid/scout/base/math"
)

type SamplerSphere interface {
	Sample(xyz math.Vector3) math.Vector4
}