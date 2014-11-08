package texture

import (
	"github.com/Opioid/scout/base/math"
)

type Sampler2D_bilinear struct {
	texture *Texture2D
	address addressMode
}

func NewSampler2D_bilinear(t *Texture2D, address addressMode) *Sampler2D_bilinear {
	s := new(Sampler2D_bilinear)
	s.texture = t
	s.address = address
	return s
}

func (sampler *Sampler2D_bilinear) Sample(uv math.Vector2) math.Vector4 {
	auv := sampler.address.address2D(uv)

	d := sampler.texture.Image.Buffers[0].dimensions

	u := auv.X * float32(d.X) - 0.5
	v := auv.Y * float32(d.Y) - 0.5

	fu := math.Floor(u)
	fv := math.Floor(v)

	x := int(fu)
	y := int(fv)

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

func (s *Sampler2D_bilinear) SampleLod(uv math.Vector2, mipLevel int) math.Vector4 {
	auv := s.address.address2D(uv)

	b := &s.texture.Image.Buffers[mipLevel]

	x := int(auv.X * float32(b.dimensions.X - 1) + 0.5)
	y := int(auv.Y * float32(b.dimensions.Y - 1) + 0.5)
	
	return b.At(x, y)
}

func bilinear(c00, c01, c10, c11 math.Vector4, s, t float32) math.Vector4 {
	_s := 1.0 - s
	_t := 1.0 - t

	return (c00.Scale(_t).Add(c01.Scale(t))).Scale(_s).Add(
		   (c10.Scale(_t).Add(c11.Scale(t))).Scale(s))
}