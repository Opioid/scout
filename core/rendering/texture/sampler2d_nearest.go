package texture

import (
	"github.com/Opioid/scout/base/math"
)

type Sampler2D_nearest struct {
	texture *Texture2D
	address addressMode
}

func NewSampler2D_nearest(t *Texture2D, address addressMode) *Sampler2D_nearest {
	s := new(Sampler2D_nearest)
	s.texture = t
	s.address = address
	return s
}

func (s *Sampler2D_nearest) Sample(uv math.Vector2) math.Vector4 {
	auv := s.address.address2D(uv)

	x := int(auv.X * float32(s.texture.image.buffers[0].dimensions.X - 1))
	y := int(auv.Y * float32(s.texture.image.buffers[0].dimensions.Y - 1))

	return s.texture.image.buffers[0].At(x, y)
}

func (s *Sampler2D_nearest) SampleLod(uv math.Vector2, mipLevel int) math.Vector4 {
	auv := s.address.address2D(uv)

	b := &s.texture.image.buffers[mipLevel]

	x := int(auv.X * float32(b.dimensions.X - 1))
	y := int(auv.Y * float32(b.dimensions.Y - 1))

	return b.At(x, y)
}