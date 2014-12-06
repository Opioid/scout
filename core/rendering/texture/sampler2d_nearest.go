package texture

import (
	"github.com/Opioid/scout/base/math"
)

type Sampler2D_nearest struct {
	address addressMode
}

func NewSampler2D_nearest(address addressMode) *Sampler2D_nearest {
	s := new(Sampler2D_nearest)
	s.address = address
	return s
}

func (s *Sampler2D_nearest) Sample(texture *Texture2D, uv math.Vector2) math.Vector4 {
	auv := s.address.address2D(uv)

	d := texture.Image.Buffers[0].dimensions

	x := math.Mini(int32(auv.X * float32(d.X)), d.X - 1)
	y := math.Mini(int32(auv.Y * float32(d.Y)), d.Y - 1)

	return texture.Image.Buffers[0].At(x, y)
}

func (s *Sampler2D_nearest) SampleLod(texture *Texture2D, uv math.Vector2, mipLevel float32) math.Vector4 {
	auv := s.address.address2D(uv)

	b := &texture.Image.Buffers[int32(mipLevel)]

	x := int32(auv.X * float32(b.dimensions.X - 1) + 0.5)
	y := int32(auv.Y * float32(b.dimensions.Y - 1) + 0.5)

	return b.At(x, y)
}