package texture

import (
	"github.com/Opioid/scout/base/math"
)

type SamplerSphere interface {
	Sample(texture *Texture2D, xyz math.Vector3) math.Vector4

	SampleLod(texture *Texture2D, xyz math.Vector3, mipLevel float32) math.Vector4
}