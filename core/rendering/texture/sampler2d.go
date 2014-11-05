package texture

import (
	"github.com/Opioid/scout/base/math"
)

type Sampler2D interface {
	Sample(uv math.Vector2) math.Vector4
	SampleLod(uv math.Vector2, mipLevel int) math.Vector4
}