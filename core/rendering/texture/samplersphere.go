package texture

import (
	"github.com/Opioid/scout/base/math"
)

type SamplerSphere interface {
	Texture() *Texture2D

	Sample(xyz math.Vector3) math.Vector4
}