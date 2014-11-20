package texture

import (
	"github.com/Opioid/scout/base/math"
)

type Sampler2D_linear struct {
	texture *Texture2D
	address addressMode
}

func NewSampler2D_linear(t *Texture2D, address addressMode) *Sampler2D_linear {
	s := new(Sampler2D_linear)
	s.texture = t
	s.address = address
	return s
}

func (sampler *Sampler2D_linear) Sample(uv math.Vector2) math.Vector4 {
	auv := sampler.address.address2D(uv)

	d := sampler.texture.Image.Buffers[0].dimensions

	u := auv.X * float32(d.X) - 0.5
	v := auv.Y * float32(d.Y) - 0.5

	fu := math.Floor(u)
	fv := math.Floor(v)

	x := int32(fu)
	y := int32(fv)

	x1 := math.Mini(x + 1, d.X - 1)
	y1 := math.Mini(y + 1, d.Y - 1)

	x = math.Maxi(x, 0)
	y = math.Maxi(y, 0)	

	c00 := sampler.texture.Image.Buffers[0].At(x,  y)	
	c01 := sampler.texture.Image.Buffers[0].At(x,  y1)	
	c10 := sampler.texture.Image.Buffers[0].At(x1, y)	
	c11 := sampler.texture.Image.Buffers[0].At(x1, y1)

	s := u - fu
	t := v - fv

	return bilinear(c00, c01, c10, c11, s, t)
}

func (s *Sampler2D_linear) SampleLod(uv math.Vector2, mipLevel float32) math.Vector4 {
	auv := s.address.address2D(uv)

	b := &s.texture.Image.Buffers[int(mipLevel)]

	x := int32(auv.X * float32(b.dimensions.X - 1) + 0.5)
	y := int32(auv.Y * float32(b.dimensions.Y - 1) + 0.5)
	
	return b.At(x, y)
}

func bilinear(c00, c01, c10, c11 math.Vector4, s, t float32) math.Vector4 {
	_s := 1 - s
	_t := 1 - t

	return (c00.Scale(_t).Add(c01.Scale(t))).Scale(_s).Add(
		   (c10.Scale(_t).Add(c11.Scale(t))).Scale(s))
}