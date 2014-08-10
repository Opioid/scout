package texture

import (
	"github.com/Opioid/scout/base/math"
)

type Sampler_nearest struct {
	texture *Texture2D
}

func NewSampler_nearest(t *Texture2D) *Sampler_nearest {
	s := new(Sampler_nearest)
	s.texture = t
	return s
}

func (s *Sampler_nearest) Sample(uv math.Vector2) math.Vector4 {
	x := int(uv.X * float32(s.texture.images[0].dimensions.X - 1))
	y := int(uv.Y * float32(s.texture.images[0].dimensions.Y - 1))
	return s.texture.images[0].at(x, y)
}