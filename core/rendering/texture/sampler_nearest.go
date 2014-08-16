package texture

import (
	"github.com/Opioid/scout/base/math"
)

type Sampler_nearest struct {
	texture *Texture2D
	address addressMode
}

func NewSampler_nearest(t *Texture2D, address addressMode) *Sampler_nearest {
	s := new(Sampler_nearest)
	s.texture = t
	s.address = address
	return s
}

func (s *Sampler_nearest) Sample(uv math.Vector2) math.Vector4 {
	auv := s.address.address2D(uv)
	x := int(auv.X * float32(s.texture.images[0].dimensions.X - 1))
	y := int(auv.Y * float32(s.texture.images[0].dimensions.Y - 1))
	return s.texture.images[0].at(x, y)
}