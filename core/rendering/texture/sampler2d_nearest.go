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
/*	auv := s.address.address2D(uv)

	x := int(auv.X * float32(s.texture.Image.Buffers[0].dimensions.X - 1) + 0.5)
	y := int(auv.Y * float32(s.texture.Image.Buffers[0].dimensions.Y - 1) + 0.5)

	return s.texture.Image.Buffers[0].At(x, y)
	*/


	auv := s.address.address2D(uv)

	d := s.texture.Image.Buffers[0].dimensions

	x := math.Mini(int(auv.X * float32(d.X)), d.X - 1)
	y := math.Mini(int(auv.Y * float32(d.Y)), d.Y - 1)

	return s.texture.Image.Buffers[0].At(x, y)
}

func (s *Sampler2D_nearest) SampleLod(uv math.Vector2, mipLevel int) math.Vector4 {
	auv := s.address.address2D(uv)

	b := &s.texture.Image.Buffers[mipLevel]

	x := int(auv.X * float32(b.dimensions.X - 1) + 0.5)
	y := int(auv.Y * float32(b.dimensions.Y - 1) + 0.5)

	return b.At(x, y)
}