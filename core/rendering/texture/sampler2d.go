package texture

import (
	"github.com/Opioid/scout/base/math"
)

type Sampler2D interface {
	Sample(texture *Texture2D, uv math.Vector2) math.Vector4
	Sample3(texture *Texture2D, uv math.Vector2) math.Vector3

	SampleLod(texture *Texture2D, uv math.Vector2, mipLevel float32) math.Vector4
	SampleLod3(texture *Texture2D, uv math.Vector2, mipLevel float32) math.Vector3
}