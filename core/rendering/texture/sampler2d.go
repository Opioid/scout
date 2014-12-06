package texture

import (
	"github.com/Opioid/scout/base/math"
)

type Sampler2D interface {
	Sample(texture *Texture2D, uv math.Vector2) math.Vector4
	SampleLod(texture *Texture2D, uv math.Vector2, mipLevel float32) math.Vector4
}